package workspace

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turbot/pipe-fittings/v2/app_specific"
	"github.com/turbot/pipe-fittings/v2/app_specific_connection"
	"github.com/turbot/pipe-fittings/v2/connection"
	"github.com/turbot/pipe-fittings/v2/filepaths"
	"github.com/turbot/powerpipe/internal/powerpipeconfig"
)

func init() {
	// Register Powerpipe connection type constructors so that HCL evaluation
	// can resolve connection-typed variables (e.g. `type = connection.steampipe`).
	// Mirrors registerConnections() in internal/cmdconfig/app_specific.go.
	// Without this, compliance mods that declare connection variables panic with
	// a nil-pointer dereference when the cty type lookup returns cty.NilType.
	app_specific_connection.RegisterConnections(
		connection.NewMysqlConnection,
		connection.NewSteampipePgConnection,
		connection.NewPostgresConnection,
		connection.NewSqliteConnection,
		connection.NewDuckDbConnection,
		connection.NewTailpipeConnection,
	)
}

// loadTimingConnections loads PipelingConnections from the default Powerpipe
// config directory (~/.powerpipe/config/*.ppc).
// Returns an empty map when no config files exist — safe to pass to
// WithPipelingConnections; the eval context will simply have no connection
// keys and WithVariableValidation(false) suppresses missing-value errors.
func loadTimingConnections() map[string]connection.PipelingConnection {
	// EnsureConfigDir requires app_specific.InstallDir to be set.
	// In tests the CLI initialisation hasn't run, so set the default here.
	if app_specific.InstallDir == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			app_specific.InstallDir = home + "/.powerpipe"
		}
	}
	cfg, ew := powerpipeconfig.LoadPowerpipeConfig(filepaths.EnsureConfigDir())
	if ew.GetError() != nil || cfg == nil || cfg.PipelingConnections == nil {
		return map[string]connection.PipelingConnection{}
	}
	return cfg.PipelingConnections
}

// PipesStartupPhases captures the wall-clock durations and resource counts
// measured during a Pipes dashboard server startup sequence.
//
// The three phases map to what a Pipes user experiences after a pod restart:
//
//  1. WorkspaceLoad      — workspace is ready; server is unblocked
//  2. FirstDashboardList — UI can display the dashboard list (immediately after load)
//  3. TagsFullyResolved  — tags/grouping metadata is fully available in the index
type PipesStartupPhases struct {
	// Wall-clock durations (cumulative from T0)
	WorkspaceLoad      time.Duration
	FirstDashboardList time.Duration
	TagsFullyResolved  time.Duration

	// Resource counts — used to verify both modes see the same resources
	Dashboards int
	Benchmarks int

	// Tag health: percentage of resources with at least one tag
	TagCoveragePct float64

	// Heap memory allocated during the load (bytes)
	HeapAllocatedBytes uint64
}

// getPipesTimingWorkspacePath returns the workspace path for timing tests.
//
// Priority:
//  1. PIPES_TIMING_MOD_PATH env var — point to a real mod for production-scale numbers
//  2. Synthetic Pipes-like workspace sized by PIPES_TIMING_NUM_RESOURCES (default 50)
//
// Examples:
//
//	PIPES_TIMING_NUM_RESOURCES=750 go test ./internal/workspace -run TestPipesDashboardStartup... -v
//	PIPES_TIMING_MOD_PATH=~/.powerpipe/mods/github.com/turbot/steampipe-mod-aws-compliance \
//	  go test ./internal/workspace -run TestPipesDashboardStartup... -v
func getPipesTimingWorkspacePath(tb testing.TB) (string, func()) {
	tb.Helper()
	if p := os.Getenv("PIPES_TIMING_MOD_PATH"); p != "" {
		return p, func() {}
	}
	n := 50
	if s := os.Getenv("PIPES_TIMING_NUM_RESOURCES"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			n = v
		}
	}
	return createScaledPipesLikeWorkspace(tb, n)
}

// createScaledPipesLikeWorkspace creates a Pipes-realistic single-mod workspace with
// numDashboards dashboards and numDashboards benchmarks, all using variable-based tags
// (mimics real AWS compliance mod structure). Resources span multiple files, like a
// real compliance mod (one .pp file per service area).
//
// Using a flat workspace (no dep mods) keeps the test self-contained. The dep-mod
// Pipes scenario is covered separately by TestPipesScenario_LazyLoadingWithMultipleDependencyMods.
func createScaledPipesLikeWorkspace(tb testing.TB, numDashboards int) (string, func()) {
	tb.Helper()
	tmpDir, err := os.MkdirTemp("", "pipes_timing_test")
	require.NoError(tb, err)

	mainModContent := `mod "pipes_timing_test" {
  title = "Pipes Timing Test"
}

variable "common_tags" {
  type = map(string)
  default = {
    service = "aws"
    env     = "test"
    owner   = "platform-team"
  }
}

variable "common_dimensions" {
  default = ["account_id", "region"]
}
`
	require.NoError(tb, os.WriteFile(filepath.Join(tmpDir, "mod.pp"), []byte(mainModContent), 0600))

	// Spread resources across numFiles files (mimics real mods where each file
	// covers a service area: s3.pp, ec2.pp, iam.pp, etc.).
	const numFiles = 10
	perFile := (numDashboards + numFiles - 1) / numFiles

	for f := 0; f < numFiles; f++ {
		start := f * perFile
		end := start + perFile
		if end > numDashboards {
			end = numDashboards
		}
		if start >= end {
			break
		}
		content := buildTimingResourceFile(f, start, end)
		fname := fmt.Sprintf("service_%02d.pp", f+1)
		require.NoError(tb, os.WriteFile(filepath.Join(tmpDir, fname), []byte(content), 0600))
	}

	return tmpDir, func() { _ = os.RemoveAll(tmpDir) }
}

// buildTimingResourceFile generates a .pp file with dashboards and benchmarks
// using variable-based tags (merge() calls trigger background resolution).
func buildTimingResourceFile(fileIdx, startIdx, endIdx int) string {
	var sb strings.Builder
	for i := startIdx; i < endIdx; i++ {
		fmt.Fprintf(&sb, `dashboard "dashboard_%03d" {
  title = "Service %02d Dashboard %03d"
  tags  = var.common_tags

  text {
    value = "Service area %02d, resource %03d"
  }
}

benchmark "benchmark_%03d" {
  title    = "Service %02d Benchmark %03d"
  tags     = merge(var.common_tags, { service_area = "area_%02d", category = "compliance" })
  children = []
}

`, i, fileIdx, i, fileIdx, i, i, fileIdx, i, fileIdx)
	}
	return sb.String()
}

// measureEagerPhases measures Pipes startup phases using full eager (v1.4.3-style) loading.
//
// In eager mode all three phases collapse to a single Load() call:
// all HCL is parsed and all tags are resolved before Load() returns.
//
// Returns (phases, true) on success or (zero, false) when eager loading is not possible
// (e.g. real workspaces with connection-typed variables require the full app context).
func measureEagerPhases(tb testing.TB, ctx context.Context, workspacePath string) (PipesStartupPhases, bool) {
	tb.Helper()

	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	heapBefore := m.HeapAlloc

	t0 := time.Now()
	ws, errAndWarnings := Load(ctx, workspacePath,
		WithVariableValidation(false),
		WithPipelingConnections(loadTimingConnections()),
	)
	tLoad := time.Since(t0)
	if err := errAndWarnings.GetError(); err != nil {
		tb.Logf("Eager load unavailable: %v", err)
		return PipesStartupPhases{}, false
	}
	defer ws.Close()

	// All resources are in memory the moment Load() returns —
	// there is no separate "first list" or "resolution" phase.
	res := ws.GetPowerpipeModResources()
	tList := time.Since(t0)
	tTags := tList

	runtime.ReadMemStats(&m)
	heapDelta := uint64(0)
	if m.HeapAlloc > heapBefore {
		heapDelta = m.HeapAlloc - heapBefore
	}

	dashCount := len(res.Dashboards)
	benchCount := len(res.ControlBenchmarks) + len(res.DetectionBenchmarks)

	total := dashCount + benchCount
	withTags := 0
	for _, d := range res.Dashboards {
		if len(d.GetTags()) > 0 {
			withTags++
		}
	}
	for _, b := range res.ControlBenchmarks {
		if len(b.GetTags()) > 0 {
			withTags++
		}
	}
	for _, b := range res.DetectionBenchmarks {
		if len(b.GetTags()) > 0 {
			withTags++
		}
	}
	tagCovPct := 0.0
	if total > 0 {
		tagCovPct = float64(withTags) / float64(total) * 100.0
	}

	return PipesStartupPhases{
		WorkspaceLoad:      tLoad,
		FirstDashboardList: tList,
		TagsFullyResolved:  tTags,
		Dashboards:         dashCount,
		Benchmarks:         benchCount,
		TagCoveragePct:     tagCovPct,
		HeapAllocatedBytes: heapDelta,
	}, true
}

// measureLazyPhases measures Pipes startup phases using lazy (v1.5.0-style) loading.
//
// Uses NewLazyWorkspace directly (bypassing the 200ms initial-resolution wait in LoadLazy)
// to measure pure index-build latency — this reflects the actual Pipes server behavior
// where the dashboard list is served as soon as the index is built, with tags resolving
// progressively in the background.
//
// Phase 1: NewLazyWorkspace() — fast HCL scan builds the resource index
// Phase 2: GetAvailableDashboardsFromIndex() — instant; index already in memory
// Phase 3: WaitForResolution() — background goroutines finish resolving variable references
func measureLazyPhases(tb testing.TB, ctx context.Context, workspacePath string) PipesStartupPhases {
	tb.Helper()

	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	heapBefore := m.HeapAlloc

	t0 := time.Now()
	// Use NewLazyWorkspace directly to measure pure index-build time without
	// the 200ms initial-resolution wait that LoadLazy adds for CLI use.
	// (The Pipes server serves dashboards immediately after index build.)
	lw, err := NewLazyWorkspace(ctx, workspacePath, DefaultLazyLoadConfig())
	tLoad := time.Since(t0)
	require.NoError(tb, err, "lazy workspace (index build)")
	defer lw.Close()

	// Start background resolution (LoadLazy would do this too).
	lw.StartBackgroundResolution()

	// Phase 2: dashboard list from index (no additional disk I/O needed)
	payload := lw.GetAvailableDashboardsFromIndex()
	tList := time.Since(t0)

	// Phase 3: wait for background variable/template resolution to finish
	lw.WaitForResolution(30 * time.Second)
	tTags := time.Since(t0)

	// Rebuild payload after full resolution to get accurate tag coverage.
	// Tags in the initial payload may be unresolved (empty) for complex expressions.
	resolvedPayload := lw.GetAvailableDashboardsFromIndex()

	runtime.ReadMemStats(&m)
	heapDelta := uint64(0)
	if m.HeapAlloc > heapBefore {
		heapDelta = m.HeapAlloc - heapBefore
	}

	dashCount := len(payload.Dashboards)
	benchCount := len(payload.Benchmarks)

	total := dashCount + benchCount
	withTags := 0
	for _, d := range resolvedPayload.Dashboards {
		if len(d.Tags) > 0 {
			withTags++
		}
	}
	for _, b := range resolvedPayload.Benchmarks {
		if len(b.Tags) > 0 {
			withTags++
		}
	}
	tagCovPct := 0.0
	if total > 0 {
		tagCovPct = float64(withTags) / float64(total) * 100.0
	}

	return PipesStartupPhases{
		WorkspaceLoad:      tLoad,
		FirstDashboardList: tList,
		TagsFullyResolved:  tTags,
		Dashboards:         dashCount,
		Benchmarks:         benchCount,
		TagCoveragePct:     tagCovPct,
		HeapAllocatedBytes: heapDelta,
	}
}

// bestPhases returns the run with the minimum WorkspaceLoad time.
// Using the minimum across runs removes GC and scheduling jitter.
func bestPhases(runs []PipesStartupPhases) PipesStartupPhases {
	best := runs[0]
	for _, r := range runs[1:] {
		if r.WorkspaceLoad < best.WorkspaceLoad {
			best = r
		}
	}
	return best
}

// TestPipesDashboardStartup_TimingComparison benchmarks Pipes-realistic startup
// timing for eager (v1.4.3) vs lazy (v1.5.0) workspace loading.
//
// Each mode is measured 3 times; the best (minimum WorkspaceLoad) run is reported.
// The synthetic workspace (100 dashboards + 100 benchmarks = 200 resources) is small
// enough to run quickly in CI while still demonstrating the phase structure.
//
// For production-scale numbers, point to a large compliance mod:
//
//	PIPES_TIMING_MOD_PATH=~/.powerpipe/mods/github.com/turbot/steampipe-mod-aws-compliance \
//	  go test ./internal/workspace -run TestPipesDashboardStartup_TimingComparison -v
func TestPipesDashboardStartup_TimingComparison(t *testing.T) {
	ctx := context.Background()
	workspacePath, cleanup := getPipesTimingWorkspacePath(t)
	defer cleanup()

	const runs = 3

	// Eager: 3 runs, take best (minimum WorkspaceLoad).
	// May be unavailable for real workspaces that require the full app context
	// (e.g. connection-typed variables resolved via Steampipe workspace profiles).
	var eagerRuns []PipesStartupPhases
	eagerAvailable := true
	for i := 0; i < runs; i++ {
		phases, ok := measureEagerPhases(t, ctx, workspacePath)
		if !ok {
			eagerAvailable = false
			break
		}
		eagerRuns = append(eagerRuns, phases)
	}

	var eager PipesStartupPhases
	if eagerAvailable {
		eager = bestPhases(eagerRuns)
	}

	// Lazy: 3 runs, take best
	var lazyRuns []PipesStartupPhases
	for i := 0; i < runs; i++ {
		lazyRuns = append(lazyRuns, measureLazyPhases(t, ctx, workspacePath))
	}
	lazy := bestPhases(lazyRuns)

	// Human-readable tabular output
	printPipesTimingComparison(t, workspacePath, eager, lazy, eagerAvailable)

	if eagerAvailable {
		// Correctness guard: both modes must see identical resource counts
		assert.Equal(t, eager.Dashboards, lazy.Dashboards,
			"eager and lazy should discover the same number of dashboards")
		assert.Equal(t, eager.Benchmarks, lazy.Benchmarks,
			"eager and lazy should discover the same number of benchmarks")

		// Performance assertion: lazy index build must be faster than eager full load.
		// This is the key Pipes improvement: users see the dashboard list in <50ms not 15s.
		// Note: lazy measurement uses NewLazyWorkspace directly (index build only) which
		// reflects the actual Pipes server behavior (no initial-resolution wait).
		assert.Less(t, lazy.FirstDashboardList, eager.FirstDashboardList,
			"lazy loading should make the first dashboard list available faster than eager loading")
	}
}

// printPipesTimingComparison writes the formatted timing table to the test log.
// When eagerAvailable is false (real workspace requires full app context), eager columns show "N/A".
func printPipesTimingComparison(t *testing.T, workspacePath string, eager, lazy PipesStartupPhases, eagerAvailable bool) {
	t.Helper()

	eagerDur := func(d time.Duration) string {
		if !eagerAvailable {
			return "N/A"
		}
		return formatDurationMS(d)
	}
	speedupStr := func(e, l time.Duration) string {
		if !eagerAvailable {
			return "N/A"
		}
		if l <= 0 {
			return "∞"
		}
		return fmt.Sprintf("%.1fx", float64(e)/float64(l))
	}

	resourceLine := fmt.Sprintf("%d dashboards | %d benchmarks", lazy.Dashboards, lazy.Benchmarks)
	if eagerAvailable && eager.Dashboards == lazy.Dashboards && eager.Benchmarks == lazy.Benchmarks {
		resourceLine += " ✓"
	}

	tagLine := fmt.Sprintf("%.1f%% lazy", lazy.TagCoveragePct)
	if eagerAvailable {
		tagLine = fmt.Sprintf("%.1f%% eager  | %s", eager.TagCoveragePct, tagLine)
		if eager.TagCoveragePct == lazy.TagCoveragePct {
			tagLine += " ✓"
		}
	}

	memLine := fmt.Sprintf("%.1f MB lazy", float64(lazy.HeapAllocatedBytes)/1024/1024)
	if eagerAvailable && eager.HeapAllocatedBytes > 0 && lazy.HeapAllocatedBytes > 0 {
		ratio := float64(eager.HeapAllocatedBytes) / float64(lazy.HeapAllocatedBytes)
		memLine = fmt.Sprintf("%.1f MB eager | %s  (%.1fx less)",
			float64(eager.HeapAllocatedBytes)/1024/1024, memLine, ratio)
	}

	t.Logf("")
	t.Logf("=== Pipes Dashboard Startup: v1.4.3 (Eager) vs v1.5.0 (Lazy) ===")
	t.Logf("Workspace: %s", workspacePath)
	if !eagerAvailable {
		t.Logf("Note: Eager load unavailable — real workspaces with connection-typed variables")
		t.Logf("      require the full app context (Steampipe workspace profile + connections).")
		t.Logf("      Lazy loading works without connections because it only scans HCL metadata.")
	}
	t.Logf("")
	t.Logf("%-35s  %-14s  %-14s  %s",
		"Phase", "Eager (v1.4.3)", "Lazy (v1.5.0)", "Speedup")
	t.Logf("%-35s  %-14s  %-14s  %s",
		strings.Repeat("─", 35), strings.Repeat("─", 14),
		strings.Repeat("─", 14), strings.Repeat("─", 7))
	t.Logf("%-35s  %-14s  %-14s  %s",
		"1. Workspace Load",
		eagerDur(eager.WorkspaceLoad), formatDurationMS(lazy.WorkspaceLoad),
		speedupStr(eager.WorkspaceLoad, lazy.WorkspaceLoad))
	t.Logf("%-35s  %-14s  %-14s  %s",
		"2. First Dashboard List Available",
		eagerDur(eager.FirstDashboardList), formatDurationMS(lazy.FirstDashboardList),
		speedupStr(eager.FirstDashboardList, lazy.FirstDashboardList))
	t.Logf("   (users can browse immediately)")
	t.Logf("%-35s  %-14s  %-14s  %s",
		"3. All Tags Fully Resolved",
		eagerDur(eager.TagsFullyResolved), formatDurationMS(lazy.TagsFullyResolved),
		speedupStr(eager.TagsFullyResolved, lazy.TagsFullyResolved))
	t.Logf("   (UI grouping dropdowns ready)")
	t.Logf("")
	t.Logf("Resources:                   %s", resourceLine)
	t.Logf("Tag Coverage:                %s", tagLine)
	t.Logf("Memory (heap allocated):     %s", memLine)
	t.Logf("")
	t.Logf("Lazy measurement uses index-build time (NewLazyWorkspace, no initial-resolution wait).")
	t.Logf("This reflects the actual Pipes server behavior: serve dashboards immediately,")
	t.Logf("resolve tags progressively in the background.")
	t.Logf("")
	t.Logf("NOTE: Run with PIPES_TIMING_MOD_PATH=/path/to/aws-compliance for production-scale numbers.")
	t.Logf("Expected production speedup: 30-40x (15-20s eager → 300-500ms lazy index build)")
}

// formatDurationMS formats a duration as "NNN ms" with thousands separator for readability.
func formatDurationMS(d time.Duration) string {
	ms := d.Milliseconds()
	if ms < 1000 {
		return fmt.Sprintf("%d ms", ms)
	}
	return fmt.Sprintf("%d,%03d ms", ms/1000, ms%1000)
}

// measureLazyProductionPath measures workspace startup using the production CLI lazy path
// (LoadLazy), which is what Powerpipe uses when POWERPIPE_WORKSPACE_PRELOAD is not set.
//
// Unlike measureLazyPhases (which calls NewLazyWorkspace directly), this function uses
// LoadLazy and therefore includes the 200 ms initial-resolution wait that the CLI adds
// so that literal tags are ready before the first dashboard list is served.
//
// Phases measured (cumulative from T0):
//  1. WorkspaceLoad      — LoadLazy() returns (index built + 200ms wait)
//  2. FirstDashboardList — GetAvailableDashboardsFromIndex() called (near-instant after Load)
//  3. TagsFullyResolved  — WaitForResolution() completes (all variable references resolved)
func measureLazyProductionPath(tb testing.TB, ctx context.Context, workspacePath string) PipesStartupPhases {
	tb.Helper()

	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	heapBefore := m.HeapAlloc

	t0 := time.Now()
	lw, err := LoadLazy(ctx, workspacePath)
	tLoad := time.Since(t0)
	require.NoError(tb, err, "LoadLazy (production path)")
	defer lw.Close()

	// Phase 2: dashboard list — already in memory after LoadLazy returns
	payload := lw.GetAvailableDashboardsFromIndex()
	tList := time.Since(t0)

	// Phase 3: wait for any background resolution still in progress
	lw.WaitForResolution(30 * time.Second)
	tTags := time.Since(t0)

	// Rebuild payload after full resolution for accurate tag coverage.
	resolvedPayload := lw.GetAvailableDashboardsFromIndex()

	runtime.ReadMemStats(&m)
	heapDelta := uint64(0)
	if m.HeapAlloc > heapBefore {
		heapDelta = m.HeapAlloc - heapBefore
	}

	dashCount := len(payload.Dashboards)
	benchCount := len(payload.Benchmarks)

	total := dashCount + benchCount
	withTags := 0
	for _, d := range resolvedPayload.Dashboards {
		if len(d.Tags) > 0 {
			withTags++
		}
	}
	for _, b := range resolvedPayload.Benchmarks {
		if len(b.Tags) > 0 {
			withTags++
		}
	}
	tagCovPct := 0.0
	if total > 0 {
		tagCovPct = float64(withTags) / float64(total) * 100.0
	}

	return PipesStartupPhases{
		WorkspaceLoad:      tLoad,
		FirstDashboardList: tList,
		TagsFullyResolved:  tTags,
		Dashboards:         dashCount,
		Benchmarks:         benchCount,
		TagCoveragePct:     tagCovPct,
		HeapAllocatedBytes: heapDelta,
	}
}

// TestWorkspaceStartup_ProductionPathComparison compares the two code paths that the
// production CLI takes based on POWERPIPE_WORKSPACE_PRELOAD:
//
//   - Lazy path  (env var NOT set):    LoadLazy — index build + 200ms initial-resolution wait
//   - Eager path (env var = "true"):   Load     — full HCL parse + all variable resolution
//
// The test always runs both paths and prints a side-by-side comparison. The
// POWERPIPE_WORKSPACE_PRELOAD column shows which path production would take.
//
// Run with lazy active (default):
//
//	go test ./internal/workspace -run TestWorkspaceStartup_ProductionPathComparison -v -count=1
//
// Run with eager active (mirrors POWERPIPE_WORKSPACE_PRELOAD=true):
//
//	POWERPIPE_WORKSPACE_PRELOAD=true \
//	  go test ./internal/workspace -run TestWorkspaceStartup_ProductionPathComparison -v -count=1
//
// Use PIPES_TIMING_MOD_PATH for a real workspace:
//
//	PIPES_TIMING_MOD_PATH=/path/to/workspace \
//	  go test ./internal/workspace -run TestWorkspaceStartup_ProductionPathComparison -v -count=1
func TestWorkspaceStartup_ProductionPathComparison(t *testing.T) {
	ctx := context.Background()
	workspacePath, cleanup := getPipesTimingWorkspacePath(t)
	defer cleanup()

	// Determine which production path is "active" based on the env var.
	// This does not change which code paths the test exercises — both always run.
	preloadEnabled := os.Getenv("POWERPIPE_WORKSPACE_PRELOAD") == "true" ||
		os.Getenv("POWERPIPE_WORKSPACE_PRELOAD") == "1" ||
		os.Getenv("POWERPIPE_WORKSPACE_PRELOAD") == "yes"

	const runs = 3

	// Lazy path: LoadLazy (production default — index build + 200ms initial-resolution wait)
	var lazyProdRuns []PipesStartupPhases
	for i := 0; i < runs; i++ {
		lazyProdRuns = append(lazyProdRuns, measureLazyProductionPath(t, ctx, workspacePath))
	}
	lazyProd := bestPhases(lazyProdRuns)

	// Eager path: Load (production path when POWERPIPE_WORKSPACE_PRELOAD=true)
	var eagerRuns []PipesStartupPhases
	eagerAvailable := true
	for i := 0; i < runs; i++ {
		phases, ok := measureEagerPhases(t, ctx, workspacePath)
		if !ok {
			eagerAvailable = false
			break
		}
		eagerRuns = append(eagerRuns, phases)
	}
	var eager PipesStartupPhases
	if eagerAvailable {
		eager = bestPhases(eagerRuns)
	}

	printProductionPathComparison(t, workspacePath, lazyProd, eager, eagerAvailable, preloadEnabled)

	if eagerAvailable {
		assert.Equal(t, lazyProd.Dashboards, eager.Dashboards,
			"lazy and eager paths should discover the same number of dashboards")
		assert.Equal(t, lazyProd.Benchmarks, eager.Benchmarks,
			"lazy and eager paths should discover the same number of benchmarks")
	}
}

// printProductionPathComparison writes the production-path timing table to the test log.
func printProductionPathComparison(t *testing.T, workspacePath string,
	lazy, eager PipesStartupPhases, eagerAvailable, preloadEnabled bool) {
	t.Helper()

	eagerDur := func(d time.Duration) string {
		if !eagerAvailable {
			return "N/A"
		}
		return formatDurationMS(d)
	}
	ratioStr := func(eager, lazy time.Duration) string {
		if !eagerAvailable || lazy <= 0 {
			return "N/A"
		}
		return fmt.Sprintf("%.1fx", float64(eager)/float64(lazy))
	}

	activeMode := "lazy (POWERPIPE_WORKSPACE_PRELOAD not set)"
	if preloadEnabled {
		activeMode = "eager (POWERPIPE_WORKSPACE_PRELOAD=true)"
	}

	resourceLine := fmt.Sprintf("%d dashboards | %d benchmarks", lazy.Dashboards, lazy.Benchmarks)
	if eagerAvailable && eager.Dashboards == lazy.Dashboards && eager.Benchmarks == lazy.Benchmarks {
		resourceLine += " ✓"
	}

	t.Logf("")
	t.Logf("=== Workspace Startup: Production Path Comparison ===")
	t.Logf("Workspace:                   %s", workspacePath)
	t.Logf("Active production path:      %s", activeMode)
	if !eagerAvailable {
		t.Logf("Note: Eager load unavailable — real workspaces with connection-typed variables")
		t.Logf("      require the full app context (Steampipe workspace profile + connections).")
	}
	t.Logf("")
	t.Logf("%-35s  %-16s  %-16s  %s",
		"Phase", "Lazy (default)", "Eager (preload)", "Eager/Lazy")
	t.Logf("%-35s  %-16s  %-16s  %s",
		strings.Repeat("─", 35), strings.Repeat("─", 16),
		strings.Repeat("─", 16), strings.Repeat("─", 9))
	t.Logf("%-35s  %-16s  %-16s  %s",
		"1. Workspace Load",
		formatDurationMS(lazy.WorkspaceLoad), eagerDur(eager.WorkspaceLoad),
		ratioStr(eager.WorkspaceLoad, lazy.WorkspaceLoad))
	t.Logf("%-35s  %-16s  %-16s  %s",
		"2. First Dashboard List",
		formatDurationMS(lazy.FirstDashboardList), eagerDur(eager.FirstDashboardList),
		ratioStr(eager.FirstDashboardList, lazy.FirstDashboardList))
	t.Logf("%-35s  %-16s  %-16s  %s",
		"3. All Tags Fully Resolved",
		formatDurationMS(lazy.TagsFullyResolved), eagerDur(eager.TagsFullyResolved),
		ratioStr(eager.TagsFullyResolved, lazy.TagsFullyResolved))
	t.Logf("")
	t.Logf("Resources:  %s", resourceLine)
	t.Logf("Memory:     %.1f MB lazy | %s eager",
		float64(lazy.HeapAllocatedBytes)/1024/1024,
		func() string {
			if !eagerAvailable {
				return "N/A"
			}
			return fmt.Sprintf("%.1f MB", float64(eager.HeapAllocatedBytes)/1024/1024)
		}())
	t.Logf("")
	t.Logf("NOTE: Lazy path uses LoadLazy (index build + 200ms initial-resolution wait).")
	t.Logf("      This is the production CLI path, not the Pipes server path")
	t.Logf("      (which uses NewLazyWorkspace and serves dashboards before the wait).")
	t.Logf("      See TestPipesDashboardStartup_TimingComparison for Pipes-server timings.")
}

// localPerfResult holds measurements for a single mode (eager or lazy) across all
// four metrics in the local performance table.
type localPerfResult struct {
	serverStartup time.Duration
	heapBytes     uint64
	benchmarkList time.Duration
	dashboardList time.Duration
	numBenchmarks int
	numDashboards int
	numFiles      int
}

// TestLocalPerformanceComparison measures eager vs lazy/phased loading for local
// (non-Pipes) use against a real workspace, producing a table in the style of the
// release-notes performance results section.
//
// Four metrics are measured:
//
//	Server Startup  — LoadLazy (lazy) or Load (eager): time until workspace is ready
//	Memory          — heap allocated after server startup
//	benchmark list  — NewLazyWorkspace + index lookup (lazy) or Load (eager)
//	dashboard list  — same pattern as benchmark list
//
// The benchmark/dashboard list measurements use NewLazyWorkspace (no 200ms CLI wait)
// to show the fastest achievable list-command latency.  Server Startup uses LoadLazy
// (with the 200ms wait) to reflect the dashboard-server startup path.
//
// Run against a workspace with aws-compliance installed:
//
//	LOCAL_PERF_MOD_PATH=/Users/pskrbasu/pskr \
//	  go test ./internal/workspace -run TestLocalPerformanceComparison -v -count=1
func TestLocalPerformanceComparison(t *testing.T) {
	modPath := os.Getenv("LOCAL_PERF_MOD_PATH")
	if modPath == "" {
		modPath = os.Getenv("PIPES_TIMING_MOD_PATH")
	}
	if modPath == "" {
		t.Skip("Set LOCAL_PERF_MOD_PATH=/path/to/workspace to run this test")
	}

	ctx := context.Background()
	const runs = 3

	// --- Eager measurements (3 runs, take best) ---
	var eagerRuns []localPerfResult
	eagerAvail := true
	for i := 0; i < runs; i++ {
		r, ok := measureEagerLocal(t, ctx, modPath)
		if !ok {
			eagerAvail = false
			break
		}
		eagerRuns = append(eagerRuns, r)
	}
	var bestEager localPerfResult
	if eagerAvail {
		bestEager = bestLocalResult(eagerRuns)
	}

	// --- Lazy measurements (3 runs, take best) ---
	var lazyRuns []localPerfResult
	for i := 0; i < runs; i++ {
		lazyRuns = append(lazyRuns, measureLazyLocal(t, ctx, modPath))
	}
	bestLazy := bestLocalResult(lazyRuns)

	printLocalPerfTable(t, modPath, bestEager, bestLazy, eagerAvail)
}

// measureEagerLocal measures all four local-perf metrics using eager loading (Load).
// Server startup, benchmark list, and dashboard list all collapse to the same Load()
// call — resources are available immediately after it returns.
func measureEagerLocal(tb testing.TB, ctx context.Context, modPath string) (localPerfResult, bool) {
	tb.Helper()

	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	heapBefore := m.HeapAlloc

	t0 := time.Now()
	ws, ew := Load(ctx, modPath,
		WithVariableValidation(false),
		WithPipelingConnections(loadTimingConnections()),
	)
	serverStartup := time.Since(t0)
	if ew.GetError() != nil {
		tb.Logf("Eager load unavailable: %v", ew.GetError())
		return localPerfResult{}, false
	}
	defer ws.Close()

	runtime.ReadMemStats(&m)
	heapDelta := uint64(0)
	if m.HeapAlloc > heapBefore {
		heapDelta = m.HeapAlloc - heapBefore
	}

	res := ws.GetPowerpipeModResources()
	// After Load(), listing is instant — resources already in memory.
	listTime := time.Since(t0)

	return localPerfResult{
		serverStartup: serverStartup,
		heapBytes:     heapDelta,
		benchmarkList: listTime,
		dashboardList: listTime,
		numBenchmarks: len(res.ControlBenchmarks) + len(res.DetectionBenchmarks),
		numDashboards: len(res.Dashboards),
	}, true
}

// measureLazyLocal measures all four local-perf metrics using lazy loading.
//
//   - Server Startup: LoadLazy (production CLI path, includes 200ms initial-resolution wait)
//   - benchmark/dashboard list: NewLazyWorkspace + GetAvailableDashboardsFromIndex
//     (index-only path, no wait — fastest achievable list latency)
func measureLazyLocal(tb testing.TB, ctx context.Context, modPath string) localPerfResult {
	tb.Helper()

	// --- Server Startup: LoadLazy ---
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	heapBefore := m.HeapAlloc

	t0 := time.Now()
	lw, err := LoadLazy(ctx, modPath)
	serverStartup := time.Since(t0)
	require.NoError(tb, err, "LoadLazy")
	lw.Close()

	runtime.ReadMemStats(&m)
	heapDelta := uint64(0)
	if m.HeapAlloc > heapBefore {
		heapDelta = m.HeapAlloc - heapBefore
	}

	// --- benchmark list / dashboard list: NewLazyWorkspace (index only, no wait) ---
	t1 := time.Now()
	lw2, err := NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	require.NoError(tb, err, "NewLazyWorkspace")
	payload := lw2.GetAvailableDashboardsFromIndex()
	listTime := time.Since(t1)
	lw2.Close()

	return localPerfResult{
		serverStartup: serverStartup,
		heapBytes:     heapDelta,
		benchmarkList: listTime,
		dashboardList: listTime,
		numBenchmarks: len(payload.Benchmarks),
		numDashboards: len(payload.Dashboards),
	}
}

// bestLocalResult returns the run with the minimum serverStartup time.
func bestLocalResult(runs []localPerfResult) localPerfResult {
	best := runs[0]
	for _, r := range runs[1:] {
		if r.serverStartup < best.serverStartup {
			best = r
		}
	}
	return best
}

// printLocalPerfTable prints the performance results table in the release-notes style.
func printLocalPerfTable(t *testing.T, modPath string, eager, lazy localPerfResult, eagerAvail bool) {
	t.Helper()

	numResources := lazy.numBenchmarks + lazy.numDashboards
	if eagerAvail && eager.numBenchmarks+eager.numDashboards > numResources {
		numResources = eager.numBenchmarks + eager.numDashboards
	}

	eagerMS := func(d time.Duration) string {
		if !eagerAvail {
			return "N/A"
		}
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	improvement := func(e, l time.Duration) string {
		if !eagerAvail || l <= 0 {
			return "N/A"
		}
		return fmt.Sprintf("%.1fx faster", float64(e)/float64(l))
	}
	memImprovement := func(e, l uint64) string {
		if !eagerAvail || e == 0 || l == 0 {
			return "N/A"
		}
		pct := (1.0 - float64(l)/float64(e)) * 100.0
		return fmt.Sprintf("%.0f%% reduction", pct)
	}

	eagerMem := "N/A"
	if eagerAvail {
		eagerMem = fmt.Sprintf("%dMB", eager.heapBytes/1024/1024)
	}
	lazyMem := fmt.Sprintf("%dMB", lazy.heapBytes/1024/1024)

	sep := strings.Repeat("─", 68)

	t.Logf("")
	t.Logf("Performance Results")
	t.Logf("%s", sep)
	t.Logf("Tested with: %s (%d benchmarks/controls, %d dashboards)",
		modPath, numResources, lazy.numDashboards)
	t.Logf("")
	t.Logf("%-22s  %-12s  %-14s  %s", "Metric", "Eager", "Lazy/Phased", "Improvement")
	t.Logf("%-22s  %-12s  %-14s  %s",
		strings.Repeat("─", 22), strings.Repeat("─", 12),
		strings.Repeat("─", 14), strings.Repeat("─", 20))
	t.Logf("%-22s  %-12s  %-14s  %s",
		"Server Startup",
		eagerMS(eager.serverStartup),
		fmt.Sprintf("%dms", lazy.serverStartup.Milliseconds()),
		improvement(eager.serverStartup, lazy.serverStartup))
	t.Logf("%-22s  %-12s  %-14s  %s",
		"Memory at Startup",
		eagerMem, lazyMem,
		memImprovement(eager.heapBytes, lazy.heapBytes))
	t.Logf("%-22s  %-12s  %-14s  %s",
		"benchmark list",
		eagerMS(eager.benchmarkList),
		fmt.Sprintf("%dms", lazy.benchmarkList.Milliseconds()),
		improvement(eager.benchmarkList, lazy.benchmarkList))
	t.Logf("%-22s  %-12s  %-14s  %s",
		"dashboard list",
		eagerMS(eager.dashboardList),
		fmt.Sprintf("%dms", lazy.dashboardList.Milliseconds()),
		improvement(eager.dashboardList, lazy.dashboardList))
	t.Logf("%s", sep)
	t.Logf("")
	t.Logf("Server Startup: LoadLazy (lazy, includes 200ms initial-resolution wait)")
	t.Logf("               vs Load (eager, full HCL parse + variable resolution)")
	t.Logf("benchmark/dashboard list: NewLazyWorkspace + index lookup (lazy, no wait)")
	t.Logf("               vs Load + count (eager)")
}

// BenchmarkPipesDashboardStartup_Eager measures eager (v1.4.3-style) startup latency.
// Run with -benchtime=3x for stable results:
//
//	go test ./internal/workspace -bench BenchmarkPipesDashboardStartup_Eager -benchtime=3x -v
func BenchmarkPipesDashboardStartup_Eager(b *testing.B) {
	workspacePath, cleanup := getPipesTimingWorkspacePath(b)
	defer cleanup()

	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ws, errAndWarnings := Load(ctx, workspacePath, WithVariableValidation(false))
		if errAndWarnings.GetError() != nil {
			b.Fatalf("eager load failed: %v", errAndWarnings.GetError())
		}
		res := ws.GetPowerpipeModResources()
		if i == 0 {
			b.ReportMetric(float64(len(res.Dashboards)), "dashboards")
			b.ReportMetric(float64(len(res.ControlBenchmarks)), "benchmarks")
		}
		ws.Close()
	}
}

// BenchmarkPipesDashboardStartup_Lazy measures lazy (v1.5.0-style) index-build latency
// including the time to first usable dashboard list (no initial-resolution wait).
// Run with -benchtime=3x for stable results:
//
//	go test ./internal/workspace -bench BenchmarkPipesDashboardStartup_Lazy -benchtime=3x -v
func BenchmarkPipesDashboardStartup_Lazy(b *testing.B) {
	workspacePath, cleanup := getPipesTimingWorkspacePath(b)
	defer cleanup()

	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		lw, err := NewLazyWorkspace(ctx, workspacePath, DefaultLazyLoadConfig())
		if err != nil {
			b.Fatalf("lazy workspace (index build) failed: %v", err)
		}
		payload := lw.GetAvailableDashboardsFromIndex()
		if i == 0 {
			b.ReportMetric(float64(len(payload.Dashboards)), "dashboards")
			b.ReportMetric(float64(len(payload.Benchmarks)), "benchmarks")
		}
		lw.Close()
	}
}
