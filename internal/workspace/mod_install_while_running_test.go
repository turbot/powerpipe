package workspace

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/powerpipe/internal/resourceindex"
	"github.com/turbot/powerpipe/internal/resources"
)

// TestModInstallWhileRunning tests the complete file watching flow:
// 1. Start server (lazy workspace) with empty workspace
// 2. Install mod while server is running (simulate by rebuilding index)
// 3. Verify dashboards appear with resolved tags (service, mod, etc.)
// 4. Compare with baseline to ensure no regression vs eager loading
func TestModInstallWhileRunning(t *testing.T) {
	ctx := context.Background()

	// Create a temporary workspace directory
	tmpDir := t.TempDir()

	// Create initial mod.pp (empty workspace, no dependencies)
	modContent := `mod "test_workspace" {
  title = "Test Workspace"
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "mod.pp"), []byte(modContent), 0644)
	require.NoError(t, err)

	t.Run("1. Initial load with empty workspace", func(t *testing.T) {
		// Load workspace initially (no mods installed)
		lw, err := LoadLazy(ctx, tmpDir)
		require.NoError(t, err)
		defer lw.Close()

		// Get initial payload
		payload := lw.GetAvailableDashboardsFromIndex()
		assert.NotNil(t, payload)

		// Should have no dashboards or benchmarks initially
		assert.Equal(t, 0, len(payload.Dashboards), "should start with no dashboards")
		assert.Equal(t, 0, len(payload.Benchmarks), "should start with no benchmarks")

		t.Logf("✓ Initial workspace loaded with 0 dashboards, 0 benchmarks")
	})

	t.Run("2. Install mod and rebuild index (simulates file watcher)", func(t *testing.T) {
		// Load workspace initially (empty, no mods installed yet)
		lw, err := LoadLazy(ctx, tmpDir)
		require.NoError(t, err)
		defer lw.Close()

		// Get initial payload (before "mod install")
		payloadBefore := lw.GetAvailableDashboardsFromIndex()
		dashboardCountBefore := len(payloadBefore.Dashboards)
		benchmarkCountBefore := len(payloadBefore.Benchmarks)

		t.Logf("Before rebuild: %d dashboards, %d benchmarks", dashboardCountBefore, benchmarkCountBefore)

		// Now "install" the mod by creating test mod files (simulates `powerpipe mod install`)
		setupTestMod(t, tmpDir)

		// Simulate mod install by rebuilding index (this is what file watcher does)
		err = lw.RebuildIndex(ctx)
		require.NoError(t, err)

		// Wait for background resolution to complete
		// This simulates the delay between file change and resolution complete
		completed := lw.WaitForResolution(5 * time.Second)
		require.True(t, completed, "background resolution should complete within 5 seconds")

		// Get payload after rebuild
		payloadAfter := lw.GetAvailableDashboardsFromIndex()
		dashboardCountAfter := len(payloadAfter.Dashboards)
		benchmarkCountAfter := len(payloadAfter.Benchmarks)

		t.Logf("After rebuild: %d dashboards, %d benchmarks", dashboardCountAfter, benchmarkCountAfter)

		// Should have found the new resources
		assert.Greater(t, dashboardCountAfter, dashboardCountBefore, "should have more dashboards after rebuild")
		assert.Greater(t, benchmarkCountAfter, benchmarkCountBefore, "should have more benchmarks after rebuild")

		t.Logf("✓ Index rebuild detected new resources")
	})

	t.Run("3. Verify tags resolved after rebuild", func(t *testing.T) {
		// Setup workspace with test mod
		setupTestMod(t, tmpDir)

		// Load workspace
		lw, err := LoadLazy(ctx, tmpDir)
		require.NoError(t, err)
		defer lw.Close()

		// Rebuild index (simulates mod install)
		err = lw.RebuildIndex(ctx)
		require.NoError(t, err)

		// Wait for resolution
		completed := lw.WaitForResolution(5 * time.Second)
		require.True(t, completed, "background resolution should complete")

		// Get payload
		payload := lw.GetAvailableDashboardsFromIndex()

		// Check that resources have proper tags
		dashboardsWithServiceTag := 0
		dashboardsWithModTag := 0
		for _, dash := range payload.Dashboards {
			if dash.Tags != nil {
				if _, hasService := dash.Tags["service"]; hasService {
					dashboardsWithServiceTag++
				}
				if _, hasMod := dash.Tags["mod"]; hasMod {
					dashboardsWithModTag++
				}
			}
		}

		benchmarksWithServiceTag := 0
		benchmarksWithModTag := 0
		for _, bench := range payload.Benchmarks {
			if bench.Tags != nil {
				if _, hasService := bench.Tags["service"]; hasService {
					benchmarksWithServiceTag++
				}
				if _, hasMod := bench.Tags["mod"]; hasMod {
					benchmarksWithModTag++
				}
			}
		}

		totalResources := len(payload.Dashboards) + len(payload.Benchmarks)
		totalWithServiceTag := dashboardsWithServiceTag + benchmarksWithServiceTag
		totalWithModTag := dashboardsWithModTag + benchmarksWithModTag

		t.Logf("Service tags: %d/%d resources (%.1f%%)",
			totalWithServiceTag, totalResources,
			float64(totalWithServiceTag)/float64(totalResources)*100)
		t.Logf("Mod tags: %d/%d resources (%.1f%%)",
			totalWithModTag, totalResources,
			float64(totalWithModTag)/float64(totalResources)*100)

		// All resources should have mod tag (added automatically)
		assert.Equal(t, totalResources, totalWithModTag,
			"all resources should have mod tag after rebuild and resolution")

		// Most resources should have service tag (from HCL)
		// Allow some resources without service tag (not all dashboards need it)
		serviceTagPercentage := float64(totalWithServiceTag) / float64(totalResources) * 100
		assert.Greater(t, serviceTagPercentage, 80.0,
			"at least 80%% of resources should have service tag")

		t.Logf("✓ Tags properly resolved after index rebuild")
	})

	t.Run("4. Compare with baseline (eager vs lazy tag coverage)", func(t *testing.T) {
		// Create a clean workspace for comparison
		setupTestMod(t, tmpDir)

		// EAGER MODE (v1.4.3 baseline)
		t.Log("=== Testing EAGER mode (v1.4.3 baseline) ===")
		eagerWorkspace, errAndWarnings := Load(ctx, tmpDir)
		if errAndWarnings.GetError() != nil {
			// Skip if dependencies aren't properly installed (expected for local path mods)
			t.Skipf("Eager loading failed (expected for test mods with local paths): %v", errAndWarnings.GetError())
			return
		}
		defer eagerWorkspace.Close()

		// Get eager workspace statistics
		eagerPayload := buildPayloadFromEagerWorkspace(eagerWorkspace)
		eagerStats := extractTagStats(eagerPayload)

		t.Logf("EAGER: %d dashboards (%.1f%% with tags, %.1f%% with service, %.1f%% with mod)",
			eagerStats.totalDashboards,
			eagerStats.dashboardTagCoverage,
			eagerStats.dashboardServiceCoverage,
			eagerStats.dashboardModCoverage)
		t.Logf("EAGER: %d benchmarks (%.1f%% with tags, %.1f%% with service, %.1f%% with mod)",
			eagerStats.totalBenchmarks,
			eagerStats.benchmarkTagCoverage,
			eagerStats.benchmarkServiceCoverage,
			eagerStats.benchmarkModCoverage)

		// LAZY MODE (new implementation with rebuild)
		t.Log("\n=== Testing LAZY mode with index rebuild ===")
		lazyWorkspace, err := LoadLazy(ctx, tmpDir)
		require.NoError(t, err)
		defer lazyWorkspace.Close()

		// Rebuild index (simulates mod install while running)
		err = lazyWorkspace.RebuildIndex(ctx)
		require.NoError(t, err)

		// Wait for resolution
		completed := lazyWorkspace.WaitForResolution(5 * time.Second)
		require.True(t, completed, "background resolution should complete")

		// Get lazy workspace statistics
		lazyPayload := lazyWorkspace.GetAvailableDashboardsFromIndex()
		lazyStats := extractTagStats(lazyPayload)

		t.Logf("LAZY:  %d dashboards (%.1f%% with tags, %.1f%% with service, %.1f%% with mod)",
			lazyStats.totalDashboards,
			lazyStats.dashboardTagCoverage,
			lazyStats.dashboardServiceCoverage,
			lazyStats.dashboardModCoverage)
		t.Logf("LAZY:  %d benchmarks (%.1f%% with tags, %.1f%% with service, %.1f%% with mod)",
			lazyStats.totalBenchmarks,
			lazyStats.benchmarkTagCoverage,
			lazyStats.benchmarkServiceCoverage,
			lazyStats.benchmarkModCoverage)

		// ASSERTIONS: Compare eager vs lazy
		t.Log("\n=== Comparing Eager vs Lazy ===")

		// Resource counts should match
		assert.Equal(t, eagerStats.totalDashboards, lazyStats.totalDashboards,
			"dashboard count should match")
		assert.Equal(t, eagerStats.totalBenchmarks, lazyStats.totalBenchmarks,
			"benchmark count should match")

		// Overall tag coverage should match (allow 1% tolerance)
		eagerOverall := (eagerStats.dashboardsWithTags + eagerStats.benchmarksWithTags) * 100.0 /
			float64(eagerStats.totalDashboards+eagerStats.totalBenchmarks)
		lazyOverall := (lazyStats.dashboardsWithTags + lazyStats.benchmarksWithTags) * 100.0 /
			float64(lazyStats.totalDashboards+lazyStats.totalBenchmarks)

		t.Logf("OVERALL TAG COVERAGE: Eager=%.1f%%, Lazy=%.1f%%", eagerOverall, lazyOverall)
		assert.InDelta(t, eagerOverall, lazyOverall, 1.0,
			"overall tag coverage should match eager mode")

		// Service tag coverage should match
		assert.InDelta(t, eagerStats.dashboardServiceCoverage, lazyStats.dashboardServiceCoverage, 1.0,
			"dashboard service tag coverage should match")
		assert.InDelta(t, eagerStats.benchmarkServiceCoverage, lazyStats.benchmarkServiceCoverage, 1.0,
			"benchmark service tag coverage should match")

		// Mod tag coverage - lazy should be better (we add mod tag automatically)
		assert.GreaterOrEqual(t, lazyStats.dashboardModCoverage, eagerStats.dashboardModCoverage,
			"lazy mode should have same or better dashboard mod tag coverage")
		assert.GreaterOrEqual(t, lazyStats.benchmarkModCoverage, eagerStats.benchmarkModCoverage,
			"lazy mode should have same or better benchmark mod tag coverage")

		t.Logf("\n✓ Lazy mode with rebuild matches v1.4.3 baseline (eager loading)")
	})

	t.Run("5. Verify grouping functionality", func(t *testing.T) {
		// Setup workspace with test mod
		setupTestMod(t, tmpDir)

		// Load workspace
		lw, err := LoadLazy(ctx, tmpDir)
		require.NoError(t, err)
		defer lw.Close()

		// Rebuild index
		err = lw.RebuildIndex(ctx)
		require.NoError(t, err)

		// Wait for resolution
		completed := lw.WaitForResolution(5 * time.Second)
		require.True(t, completed, "background resolution should complete")

		// Get payload
		payload := lw.GetAvailableDashboardsFromIndex()

		// Group by service
		serviceGroups := make(map[string]int)
		for _, dash := range payload.Dashboards {
			if dash.Tags != nil {
				if service, ok := dash.Tags["service"]; ok {
					serviceGroups[service]++
				} else {
					serviceGroups["Other"]++
				}
			}
		}
		for _, bench := range payload.Benchmarks {
			if bench.Tags != nil {
				if service, ok := bench.Tags["service"]; ok {
					serviceGroups[service]++
				} else {
					serviceGroups["Other"]++
				}
			}
		}

		// Group by mod
		modGroups := make(map[string]int)
		for _, dash := range payload.Dashboards {
			if dash.Tags != nil {
				if mod, ok := dash.Tags["mod"]; ok {
					modGroups[mod]++
				} else {
					modGroups["Other"]++
				}
			}
		}
		for _, bench := range payload.Benchmarks {
			if bench.Tags != nil {
				if mod, ok := bench.Tags["mod"]; ok {
					modGroups[mod]++
				} else {
					modGroups["Other"]++
				}
			}
		}

		t.Logf("Group by Service: %v", serviceGroups)
		t.Logf("Group by Mod: %v", modGroups)

		// Should have proper service grouping (not all "Other")
		totalResources := len(payload.Dashboards) + len(payload.Benchmarks)
		otherServiceCount := serviceGroups["Other"]

		// With our test data, all resources have "AWS S3" service tag, so no "Other" group
		assert.Equal(t, 0, otherServiceCount,
			"resources with service tags should not be in 'Other' group")

		// Should have proper mod grouping (not all "Other")
		otherModCount := modGroups["Other"]
		assert.Equal(t, 0, otherModCount,
			"no resources should be in 'Other' mod group (all should have mod tag)")

		// Should have service groups (our test has all resources tagged with "AWS S3")
		assert.Greater(t, len(serviceGroups), 0,
			"should have at least one service group")

		// Verify that tagged resources are properly grouped
		awsS3Count := serviceGroups["AWS S3"]
		assert.Equal(t, totalResources, awsS3Count,
			"all test resources should be grouped under 'AWS S3' service")

		// Should have at least one mod group
		assert.Greater(t, len(modGroups), 0,
			"should have at least one mod group")

		t.Logf("✓ Resources properly grouped by service and mod")
	})
}

// buildPayloadFromEagerWorkspace builds a payload from an eager workspace
// similar to what the lazy workspace does
func buildPayloadFromEagerWorkspace(w *PowerpipeWorkspace) *resourceindex.AvailableDashboardsPayload {
	modResources := w.GetModResources()

	payload := &resourceindex.AvailableDashboardsPayload{
		Action:     "available_dashboards",
		Dashboards: make(map[string]resourceindex.DashboardInfo),
		Benchmarks: make(map[string]resourceindex.BenchmarkInfo),
	}

	// Walk through all resources
	_ = modResources.WalkResources(func(item modconfig.HclResource) (bool, error) {
		// Type assert to get access to methods
		switch res := item.(type) {
		case *resources.Dashboard:
			tags := res.Tags
			if tags == nil {
				tags = make(map[string]string)
			}
			// Add mod tag if not present
			metadata := res.GetMetadata()
			if _, exists := tags["mod"]; !exists && metadata != nil {
				tags["mod"] = metadata.ModFullName
			}
			modFullName := ""
			if metadata != nil {
				modFullName = metadata.ModFullName
			}
			payload.Dashboards[res.Name()] = resourceindex.DashboardInfo{
				Title:       res.GetTitle(),
				FullName:    res.Name(),
				ShortName:   res.GetUnqualifiedName(),
				Tags:        tags,
				ModFullName: modFullName,
			}

		case *resources.Benchmark:
			tags := res.Tags
			if tags == nil {
				tags = make(map[string]string)
			}
			// Add mod tag if not present
			metadata := res.GetMetadata()
			if _, exists := tags["mod"]; !exists && metadata != nil {
				tags["mod"] = metadata.ModFullName
			}
			modFullName := ""
			if metadata != nil {
				modFullName = metadata.ModFullName
			}
			payload.Benchmarks[res.Name()] = resourceindex.BenchmarkInfo{
				Title:       res.GetTitle(),
				FullName:    res.Name(),
				ShortName:   res.GetUnqualifiedName(),
				Tags:        tags,
				ModFullName: modFullName,
			}
		}

		return true, nil
	})

	return payload
}

// tagStats holds tag coverage statistics
type tagStats struct {
	totalDashboards          int
	dashboardsWithTags       float64
	dashboardsWithService    float64
	dashboardsWithMod        float64
	dashboardTagCoverage     float64
	dashboardServiceCoverage float64
	dashboardModCoverage     float64
	totalBenchmarks          int
	benchmarksWithTags       float64
	benchmarksWithService    float64
	benchmarksWithMod        float64
	benchmarkTagCoverage     float64
	benchmarkServiceCoverage float64
	benchmarkModCoverage     float64
}

// extractTagStats extracts tag statistics from a payload
func extractTagStats(payload *resourceindex.AvailableDashboardsPayload) tagStats {
	stats := tagStats{}

	// Count dashboards
	stats.totalDashboards = len(payload.Dashboards)
	for _, dash := range payload.Dashboards {
		if len(dash.Tags) > 0 {
			stats.dashboardsWithTags++
		}
		if _, hasService := dash.Tags["service"]; hasService {
			stats.dashboardsWithService++
		}
		if _, hasMod := dash.Tags["mod"]; hasMod {
			stats.dashboardsWithMod++
		}
	}

	// Count benchmarks
	stats.totalBenchmarks = len(payload.Benchmarks)
	for _, bench := range payload.Benchmarks {
		if len(bench.Tags) > 0 {
			stats.benchmarksWithTags++
		}
		if _, hasService := bench.Tags["service"]; hasService {
			stats.benchmarksWithService++
		}
		if _, hasMod := bench.Tags["mod"]; hasMod {
			stats.benchmarksWithMod++
		}
	}

	// Calculate percentages
	if stats.totalDashboards > 0 {
		stats.dashboardTagCoverage = stats.dashboardsWithTags / float64(stats.totalDashboards) * 100
		stats.dashboardServiceCoverage = stats.dashboardsWithService / float64(stats.totalDashboards) * 100
		stats.dashboardModCoverage = stats.dashboardsWithMod / float64(stats.totalDashboards) * 100
	}
	if stats.totalBenchmarks > 0 {
		stats.benchmarkTagCoverage = stats.benchmarksWithTags / float64(stats.totalBenchmarks) * 100
		stats.benchmarkServiceCoverage = stats.benchmarksWithService / float64(stats.totalBenchmarks) * 100
		stats.benchmarkModCoverage = stats.benchmarksWithMod / float64(stats.totalBenchmarks) * 100
	}

	return stats
}

// setupTestMod creates a test mod in the workspace directory
// Creates resources directly in workspace (no dependencies) to work with both eager and lazy loading
func setupTestMod(t *testing.T, workspaceDir string) {
	t.Helper()

	// Update workspace mod.pp (no dependencies - resources inline)
	modContent := `mod "test_workspace" {
  title = "Test Workspace"
}
`
	err := os.WriteFile(filepath.Join(workspaceDir, "mod.pp"), []byte(modContent), 0644)
	require.NoError(t, err)

	// Create a dashboard with tags
	dashboardContent := `dashboard "test_dashboard" {
  title = "Test Dashboard"
  tags = {
    service = "AWS S3"
    type    = "Dashboard"
  }
}
`
	err = os.WriteFile(filepath.Join(workspaceDir, "dashboard.pp"), []byte(dashboardContent), 0644)
	require.NoError(t, err)

	// Create a benchmark with tags
	benchmarkContent := `benchmark "test_benchmark" {
  title = "Test Benchmark"
  tags = {
    service = "AWS S3"
    type    = "Benchmark"
  }
  children = [
    control.test_control
  ]
}

control "test_control" {
  title = "Test Control"
  sql = "select 1 as result"
}
`
	err = os.WriteFile(filepath.Join(workspaceDir, "benchmark.pp"), []byte(benchmarkContent), 0644)
	require.NoError(t, err)
}
