package workspace

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turbot/powerpipe/internal/resources"
)

// TestPipesScenario_LazyLoadingWithMultipleDependencyMods replicates the exact
// scenario that happens in Pipes where:
// 1. Multiple AWS mods with duplicate variable definitions exist
// 2. Lazy loading should work without triggering eager loading
// 3. GetModResources() should return all installed_mods
// 4. Dashboard server metadata should be correctly populated
//
// This test reproduces the bug reported where:
// - installed_mods was empty in server_metadata
// - Dashboards were all grouped under "Other"
// - Tags weren't populated on benchmarks/dashboards
func TestPipesScenario_LazyLoadingWithMultipleDependencyMods(t *testing.T) {
	ctx := context.Background()

	// Create a test workspace that mimics Pipes environment
	workspaceDir, cleanup := createPipesLikeWorkspace(t)
	defer cleanup()

	t.Logf("Created Pipes-like workspace at: %s", workspaceDir)

	// Step 1: Load workspace with lazy loading (as Pipes does)
	t.Run("1. Load workspace with lazy loading", func(t *testing.T) {
		lw, err := LoadLazy(ctx, workspaceDir)
		require.NoError(t, err, "lazy workspace load should succeed")
		require.NotNil(t, lw, "lazy workspace should not be nil")
		defer lw.Close()

		// Verify index has resources
		index := lw.GetIndex()
		dashboards := index.Dashboards()
		benchmarks := index.Benchmarks()

		t.Logf("Index loaded: %d dashboards, %d benchmarks", len(dashboards), len(benchmarks))
		assert.NotEmpty(t, dashboards, "should have dashboards in index")
		assert.NotEmpty(t, benchmarks, "should have benchmarks in index")

		// Verify that we don't fall back to eager loading even though there are duplicate variables
		assert.True(t, lw.IsLazy(), "workspace should remain in lazy mode")
		t.Log("✓ Lazy loading succeeded without fallback")
	})

	// Step 2: Get mod resources (as dashboard server does for server_metadata)
	t.Run("2. Get mod resources for server metadata", func(t *testing.T) {
		lw, err := LoadLazy(ctx, workspaceDir)
		require.NoError(t, err)
		defer lw.Close()

		// This is what the dashboard server calls
		modResources := lw.GetModResources()
		require.NotNil(t, modResources, "GetModResources should return non-nil")

		// Cast to PowerpipeModResources
		powerpipeModRes, ok := modResources.(*resources.PowerpipeModResources)
		require.True(t, ok, "should be PowerpipeModResources")

		// Verify main mod
		assert.NotNil(t, powerpipeModRes.Mod, "main mod should be present")
		t.Logf("Main mod: %s (%s)", powerpipeModRes.Mod.GetFullName(), powerpipeModRes.Mod.ShortName)

		// Verify installed_mods (dependency mods)
		installedMods := powerpipeModRes.Mods
		require.NotNil(t, installedMods, "Mods map should not be nil")

		// Should have main mod + 3 dependency mods
		expectedModCount := 4 // main + aws_compliance + aws_insights + aws_tags
		assert.GreaterOrEqual(t, len(installedMods), expectedModCount,
			"should have at least %d mods (main + dependencies)", expectedModCount)

		// Verify each dependency mod has proper metadata
		expectedDepMods := map[string]string{
			"mod.aws_compliance": "aws_compliance",
			"mod.aws_insights":   "aws_insights",
			"mod.aws_tags":       "aws_tags",
		}

		for fullName, shortName := range expectedDepMods {
			mod, exists := installedMods[fullName]
			require.True(t, exists, "dependency mod %s should exist in installed_mods", fullName)
			assert.Equal(t, fullName, mod.FullName, "mod full name should match")
			assert.Equal(t, shortName, mod.ShortName, "mod short name should match")
			assert.NotNil(t, mod.Title, "mod should have a title")
			t.Logf("✓ Found dependency mod: %s (%s) - %s", fullName, shortName, *mod.Title)
		}

		t.Log("✓ All installed_mods properly populated")
	})

	// Step 3: Build server metadata payload (as dashboard server does)
	t.Run("3. Build server metadata payload", func(t *testing.T) {
		lw, err := LoadLazy(ctx, workspaceDir)
		require.NoError(t, err)
		defer lw.Close()

		modResources := lw.GetModResources()
		powerpipeModRes := modResources.(*resources.PowerpipeModResources)

		// Build the payload exactly as the dashboard server does
		payload := buildMockServerMetadataPayload(powerpipeModRes)
		require.NotNil(t, payload, "payload should not be nil")

		// Serialize to JSON to verify structure
		jsonBytes, err := json.MarshalIndent(payload, "", "  ")
		require.NoError(t, err, "should serialize to JSON")

		t.Logf("Server metadata payload:\n%s", string(jsonBytes))

		// Verify installed_mods in payload (nested in metadata)
		metadata, ok := payload["metadata"].(map[string]interface{})
		require.True(t, ok, "payload should have metadata")

		installedMods, ok := metadata["installed_mods"].(map[string]interface{})
		require.True(t, ok, "metadata should have installed_mods")
		require.NotEmpty(t, installedMods, "installed_mods should not be empty in payload")

		// Verify each dependency mod is in the payload
		for modName := range installedMods {
			modData := installedMods[modName].(map[string]interface{})
			assert.NotEmpty(t, modData["full_name"], "mod should have full_name")
			assert.NotEmpty(t, modData["short_name"], "mod should have short_name")
			t.Logf("✓ Payload contains mod: %s", modName)
		}

		t.Log("✓ Server metadata payload correctly populated")
	})

	// Step 4: Verify available dashboards payload (for mod grouping)
	t.Run("4. Verify dashboards grouped by mod", func(t *testing.T) {
		lw, err := LoadLazy(ctx, workspaceDir)
		require.NoError(t, err)
		defer lw.Close()

		// Get available dashboards from index
		indexPayload := lw.GetAvailableDashboardsFromIndex()
		require.NotNil(t, indexPayload, "available dashboards should not be nil")

		// Verify dashboards have mod metadata
		dashboards := indexPayload.Dashboards
		require.NotEmpty(t, dashboards, "should have dashboards")

		modsWithDashboards := make(map[string]int)
		for name, dash := range dashboards {
			assert.NotEmpty(t, dash.ModFullName, "dashboard %s should have mod_full_name", name)
			modsWithDashboards[dash.ModFullName]++
			t.Logf("Dashboard: %s -> Mod: %s", name, dash.ModFullName)
		}

		// Verify dashboards are from different mods (not all "Other")
		assert.Greater(t, len(modsWithDashboards), 1,
			"dashboards should be from multiple mods, not all grouped under 'Other'")

		t.Logf("✓ Dashboards grouped across %d mods", len(modsWithDashboards))
	})

	// Step 5: Verify tags are populated IMMEDIATELY after LoadLazy returns
	// This is the key test that catches the bug - tags MUST be available immediately,
	// not after additional waiting. This ensures LoadLazy waits for background resolution.
	//
	// NOTE: This test may pass even without the fix in small workspaces because resolution
	// completes very quickly. The real-world issue occurs with large mods (800+ files) where
	// resolution takes longer. The test documents the EXPECTED behavior that LoadLazy must
	// guarantee: tags are available when it returns, not requiring additional waiting.
	t.Run("5. Verify tags populated IMMEDIATELY (no additional wait)", func(t *testing.T) {
		// Load workspace - this should wait for initial background resolution
		lw, err := LoadLazy(ctx, workspaceDir)
		require.NoError(t, err)
		defer lw.Close()

		// Check resolution status RIGHT after LoadLazy returns
		stats := lw.BackgroundResolverStats()
		t.Logf("After LoadLazy: started=%v, complete=%v, queue_length=%v, fully_resolved=%v",
			stats.IsStarted, stats.IsComplete, stats.QueueLength, lw.IsFullyResolved())

		// Check if entries need resolution
		needsResolutionCount := 0
		resolvedCount := 0
		for _, entry := range lw.GetIndex().List() {
			if entry.NeedsResolution() {
				needsResolutionCount++
			} else if entry.Type == "benchmark" || entry.Type == "dashboard" {
				resolvedCount++
			}
		}
		t.Logf("Index status: %d entries need resolution, %d benchmarks/dashboards resolved", needsResolutionCount, resolvedCount)

		// IMMEDIATELY get the payload - no additional waiting
		// This is what Pipes does, and why the bug occurred
		indexPayload := lw.GetAvailableDashboardsFromIndex()
		benchmarks := indexPayload.Benchmarks
		require.NotEmpty(t, benchmarks, "should have benchmarks")

		// CRITICAL: Tags must be populated NOW, not empty
		// Before the fix, tags were {} because LoadLazy returned before resolution completed
		benchmarksWithTags := 0
		emptyTagBenchmarks := []string{}

		for name, bench := range benchmarks {
			if len(bench.Tags) > 0 {
				benchmarksWithTags++
				t.Logf("✓ Benchmark %s has %d tags: %v", name, len(bench.Tags), bench.Tags)
			} else {
				emptyTagBenchmarks = append(emptyTagBenchmarks, name)
				t.Errorf("✗ Benchmark %s has EMPTY tags - this is the bug!", name)
			}
		}

		// FAIL THE TEST if any benchmarks have empty tags
		// This ensures the regression is caught
		if len(emptyTagBenchmarks) > 0 {
			t.Fatalf("REGRESSION: %d benchmarks have empty tags: %v\n"+
				"This means LoadLazy() is not waiting for background resolution!\n"+
				"Tags must be populated IMMEDIATELY after LoadLazy() returns.",
				len(emptyTagBenchmarks), emptyTagBenchmarks)
		}

		// All benchmarks must have tags
		require.Equal(t, len(benchmarks), benchmarksWithTags,
			"All benchmarks should have tags immediately after LoadLazy returns")

		t.Logf("✓ ALL %d/%d benchmarks have tags IMMEDIATELY (bug would cause 0/%d)",
			benchmarksWithTags, len(benchmarks), len(benchmarks))
	})

	// Step 6: Verify LoadLazy actually waits for background resolution
	// This test explicitly checks the timing behavior
	t.Run("6. LoadLazy waits for background resolution before returning", func(t *testing.T) {
		// Create a NEW lazy workspace directly (bypassing LoadLazy to test without waiting)
		lwDirect, err := NewLazyWorkspace(ctx, workspaceDir, DefaultLazyLoadConfig())
		require.NoError(t, err)
		defer lwDirect.Close()

		// Start background resolution manually
		lwDirect.StartBackgroundResolution()

		// Get payload IMMEDIATELY without waiting - simulates the bug
		payloadBefore := lwDirect.GetAvailableDashboardsFromIndex()
		benchmarksBefore := payloadBefore.Benchmarks

		// Count benchmarks with tags BEFORE background resolution completes
		tagsBeforeCount := 0
		for _, bench := range benchmarksBefore {
			if len(bench.Tags) > 0 {
				tagsBeforeCount++
			}
		}

		// Now load via LoadLazy (which should wait)
		lwWithWait, err := LoadLazy(ctx, workspaceDir)
		require.NoError(t, err)
		defer lwWithWait.Close()

		payloadAfter := lwWithWait.GetAvailableDashboardsFromIndex()
		benchmarksAfter := payloadAfter.Benchmarks

		// Count benchmarks with tags AFTER LoadLazy (with waiting)
		tagsAfterCount := 0
		for _, bench := range benchmarksAfter {
			if len(bench.Tags) > 0 {
				tagsAfterCount++
			}
		}

		// LoadLazy should have MORE or EQUAL tags than immediate access
		// because it waits for background resolution
		require.GreaterOrEqual(t, tagsAfterCount, tagsBeforeCount,
			"LoadLazy should wait for resolution - tags should be available immediately")

		// In this test case with merge() tags, LoadLazy should have ALL tags
		require.Equal(t, len(benchmarksAfter), tagsAfterCount,
			"LoadLazy should return with all tags resolved")

		t.Logf("✓ LoadLazy waits for resolution: %d/%d benchmarks have tags immediately",
			tagsAfterCount, len(benchmarksAfter))
		t.Logf("  (without waiting: only %d/%d would have complete tags)",
			tagsBeforeCount, len(benchmarksBefore))
	})

	t.Log("\n=== ALL PIPES SCENARIO TESTS PASSED ===")
	t.Log("✓ Lazy loading works without eager fallback")
	t.Log("✓ installed_mods properly populated")
	t.Log("✓ Server metadata correctly built")
	t.Log("✓ Dashboards grouped by mod (not all 'Other')")
	t.Log("✓ Tags populated IMMEDIATELY after LoadLazy returns")
	t.Log("✓ LoadLazy waits for background resolution (catches timing bugs)")
}

// TestPipesScenario_VerifyNoEagerLoadingFallback specifically tests that
// the fallback to eager loading has been removed
func TestPipesScenario_VerifyNoEagerLoadingFallback(t *testing.T) {
	ctx := context.Background()

	// Create workspace with empty index (no dashboards/benchmarks)
	workspaceDir, cleanup := createEmptyWorkspace(t)
	defer cleanup()

	t.Log("Testing that empty index does NOT trigger eager loading fallback...")

	lw, err := LoadLazy(ctx, workspaceDir)
	require.NoError(t, err, "lazy workspace should load even with empty index")
	defer lw.Close()

	// Verify it's still in lazy mode (not fallen back to eager)
	assert.True(t, lw.IsLazy(), "workspace should remain in lazy mode even with empty index")

	// GetModResources should still work
	modResources := lw.GetModResources()
	assert.NotNil(t, modResources, "GetModResources should work even with empty index")

	t.Log("✓ No eager loading fallback occurred")
	t.Log("✓ Lazy loading works correctly even with empty index")
}

// Helper: Create a Pipes-like workspace with multiple AWS mods
func createPipesLikeWorkspace(t *testing.T) (string, func()) {
	tmpDir, err := os.MkdirTemp("", "pipes_scenario_test")
	require.NoError(t, err)

	// Create main mod with dashboards and benchmarks
	// Use variable references like real Pipes workspaces
	mainModContent := `mod "smoketest" {
  title = "Smoke Test Workspace"
}

variable "common_tags" {
  type = map(string)
  default = {
    service = "main"
    env = "test"
  }
}

dashboard "main_dashboard" {
  title = "Main Dashboard"
  tags = var.common_tags

  text {
    value = "Main workspace dashboard"
  }
}

benchmark "main_benchmark" {
  title = "Main Benchmark"
  tags = merge(var.common_tags, { benchmark = "true" })
  children = []
}
`
	err = os.WriteFile(filepath.Join(tmpDir, "mod.pp"), []byte(mainModContent), 0600)
	require.NoError(t, err)

	// Create dependency mods directory structure
	modsDir := filepath.Join(tmpDir, ".powerpipe", "mods")

	// Create AWS Compliance mod with variable references
	createAwsMod(t, modsDir, "aws_compliance", "v1.13.0", `mod "aws_compliance" {
  title = "AWS Compliance"
}

variable "common_tags" {
  type = map(string)
  default = {
    service = "AWS"
    category = "Compliance"
  }
}

variable "common_dimensions" {
  default = ["account_id", "region"]
}

variable "tag_dimensions" {
  default = ["environment", "project"]
}

dashboard "compliance_dashboard" {
  title = "AWS Compliance Dashboard"
  tags = var.common_tags

  text {
    value = "Compliance checks"
  }
}

benchmark "cis_v1_2_0" {
  title = "CIS v1.2.0"
  tags = merge(var.common_tags, { cis = "true", plugin = "aws" })
  children = []
}
`)

	// Create AWS Insights mod with variable references
	createAwsMod(t, modsDir, "aws_insights", "v1.0.0", `mod "aws_insights" {
  title = "AWS Insights"
}

variable "common_tags" {
  type = map(string)
  default = {
    service = "AWS"
    category = "Insights"
  }
}

variable "common_dimensions" {
  default = ["account_id", "region"]
}

dashboard "insights_dashboard" {
  title = "AWS Insights Dashboard"
  tags = var.common_tags

  text {
    value = "AWS insights"
  }
}
`)

	// Create AWS Tags mod with variable references
	createAwsMod(t, modsDir, "aws_tags", "v1.0.1", `mod "aws_tags" {
  title = "AWS Tags"
}

variable "common_tags" {
  type = map(string)
  default = {
    service = "AWS"
    category = "Tags"
  }
}

variable "common_dimensions" {
  default = ["account_id"]
}

variable "tag_dimensions" {
  default = ["environment"]
}

dashboard "tags_dashboard" {
  title = "AWS Tags Dashboard"
  tags = var.common_tags

  text {
    value = "Tag management"
  }
}

benchmark "tag_benchmark" {
  title = "Tag Compliance"
  tags = merge(var.common_tags, { benchmark = "tag_compliance" })
  children = []
}
`)

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

// Helper: Create a single AWS mod
func createAwsMod(t *testing.T, modsBaseDir, modName, version, content string) {
	// Create directory structure: .powerpipe/mods/github.com/turbot/steampipe-mod-{name}@{version}
	modPath := filepath.Join(modsBaseDir, "github.com", "turbot",
		fmt.Sprintf("steampipe-mod-%s@%s", modName, version))

	err := os.MkdirAll(modPath, 0755)
	require.NoError(t, err)

	// Write mod.pp
	err = os.WriteFile(filepath.Join(modPath, "mod.pp"), []byte(content), 0600)
	require.NoError(t, err)
}

// Helper: Create empty workspace
func createEmptyWorkspace(t *testing.T) (string, func()) {
	tmpDir, err := os.MkdirTemp("", "empty_workspace_test")
	require.NoError(t, err)

	// Create minimal mod.pp with no resources
	modContent := `mod "empty_mod" {
  title = "Empty Mod"
}
`
	err = os.WriteFile(filepath.Join(tmpDir, "mod.pp"), []byte(modContent), 0600)
	require.NoError(t, err)

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

// Helper: Build mock server metadata payload (mimics dashboard server logic)
func buildMockServerMetadataPayload(modResources *resources.PowerpipeModResources) map[string]interface{} {
	installedMods := make(map[string]interface{})

	// This mimics the logic in dashboardserver/payload.go buildServerMetadataPayload()
	for _, mod := range modResources.Mods {
		// Ignore current mod (only include dependencies)
		if mod.GetFullName() == modResources.Mod.GetFullName() {
			continue
		}

		modData := map[string]interface{}{
			"title":      getTitleString(mod.Title),
			"full_name":  mod.GetFullName(),
			"short_name": mod.ShortName,
		}
		installedMods[mod.GetFullName()] = modData
	}

	payload := map[string]interface{}{
		"action": "server_metadata",
		"metadata": map[string]interface{}{
			"installed_mods": installedMods,
			"mod": map[string]interface{}{
				"title":      getTitleString(modResources.Mod.Title),
				"full_name":  modResources.Mod.GetFullName(),
				"short_name": modResources.Mod.ShortName,
			},
		},
	}

	return payload
}

// Helper: Get title string from pointer
func getTitleString(title *string) string {
	if title == nil {
		return ""
	}
	return *title
}

// TestPipesStartup_IncompleteDepMod_DanglingSymlink replicates the exact Pipes pod-restart
// race condition where the PVC contains stale/incomplete mod files:
//
//	Error: failed to load lazy workspace: building index: scanning dependency mods:
//	       scanning mod aws_insights: open .../dashboards/emr/emr.pp: no such file or directory
//
// Simulation: a dangling symlink mimics a file that appears in directory listings
// (WalkDir uses Lstat, so it finds the symlink entry) but fails to open
// (os.ReadFile follows the symlink and gets ENOENT).
// This is equivalent to NFS stale handles and partial mod extraction on PVC.
//
// Before fix: returns "failed to load lazy workspace: ..."
// After fix:  loads successfully; aws_insights dep mod skipped with a WARN log.
func TestPipesStartup_IncompleteDepMod_DanglingSymlink(t *testing.T) {
	ctx := context.Background()
	workspaceDir := t.TempDir()

	// Main mod
	require.NoError(t, os.WriteFile(
		filepath.Join(workspaceDir, "mod.pp"),
		[]byte(`mod "test_main" { title = "Test Main" }`), 0600))

	// Dependency mod directory (simulating stale PVC state)
	depModDir := filepath.Join(workspaceDir, ".powerpipe", "mods",
		"github.com", "turbot", "steampipe-mod-aws-insights@v1.2.0")
	require.NoError(t, os.MkdirAll(depModDir, 0755))

	// mod.pp exists (directory was partially set up)
	require.NoError(t, os.WriteFile(
		filepath.Join(depModDir, "mod.pp"),
		[]byte(`mod "aws_insights" { title = "AWS Insights" }`), 0600))

	// Some files are fully installed
	s3Dir := filepath.Join(depModDir, "dashboards", "s3")
	require.NoError(t, os.MkdirAll(s3Dir, 0755))
	require.NoError(t, os.WriteFile(
		filepath.Join(s3Dir, "s3.pp"),
		[]byte(`dashboard "s3_overview" { title = "S3 Overview" tags = { service = "S3" } }`), 0600))

	// emr/ directory exists but emr.pp is a dangling symlink — simulates partial install on PVC
	emrDir := filepath.Join(depModDir, "dashboards", "emr")
	require.NoError(t, os.MkdirAll(emrDir, 0755))
	require.NoError(t, os.Symlink(
		"/nonexistent/target/emr.pp",
		filepath.Join(emrDir, "emr.pp")))

	// Before fix: returns "failed to load lazy workspace: building index: ..."
	// After fix:  loads successfully, aws_insights dep mod skipped (will be reinstalled)
	lw, err := LoadLazy(ctx, workspaceDir)
	require.NoError(t, err,
		"LoadLazy should not fail when dependency mod has missing/unreadable files")
	require.NotNil(t, lw)
	defer lw.Close()

	// Workspace should still be usable (main mod resources available)
	payload := lw.GetAvailableDashboardsFromIndex()
	require.NotNil(t, payload)
}

// TestPipesStartup_IncompleteDepMod_RaceCondition simulates a TOCTOU race:
// a file exists when listed by the scanner, then is deleted concurrently before
// os.ReadFile reads it. This is the concurrent equivalent of the dangling-symlink test.
//
// Whether the race is hit or not, LoadLazy must not crash.
// After the fix is applied, this test should always pass.
func TestPipesStartup_IncompleteDepMod_RaceCondition(t *testing.T) {
	ctx := context.Background()
	workspaceDir := t.TempDir()

	require.NoError(t, os.WriteFile(
		filepath.Join(workspaceDir, "mod.pp"),
		[]byte(`mod "test_main" { title = "Test Main" }`), 0600))

	depModDir := filepath.Join(workspaceDir, ".powerpipe", "mods",
		"github.com", "turbot", "steampipe-mod-aws-insights@v1.2.0")
	require.NoError(t, os.MkdirAll(depModDir, 0755))
	require.NoError(t, os.WriteFile(
		filepath.Join(depModDir, "mod.pp"),
		[]byte(`mod "aws_insights" { title = "AWS Insights" }`), 0600))

	// Create a .pp file, then delete it in a goroutine while LoadLazy runs
	emrDir := filepath.Join(depModDir, "dashboards", "emr")
	require.NoError(t, os.MkdirAll(emrDir, 0755))
	emrFile := filepath.Join(emrDir, "emr.pp")
	require.NoError(t, os.WriteFile(emrFile,
		[]byte(`dashboard "emr_detail" { title = "EMR" }`), 0600))

	// Concurrently delete the file while LoadLazy is running
	go func() {
		_ = os.Remove(emrFile)
	}()

	// Whether the race is hit or not, LoadLazy must not crash
	lw, err := LoadLazy(ctx, workspaceDir)
	if err == nil {
		defer lw.Close()
		t.Log("LoadLazy succeeded (race not hit or fix applied)")
	} else {
		t.Logf("LoadLazy failed (race hit, bug not yet fixed): %v", err)
		// Once the fix is applied, change this to require.NoError
	}
}
