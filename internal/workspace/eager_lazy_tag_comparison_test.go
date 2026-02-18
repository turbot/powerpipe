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
)

// TestEagerLazyTagComparison compares tag loading between eager and lazy modes
// on a real workspace to ensure parity with v1.4.3 behavior.
func TestEagerLazyTagComparison(t *testing.T) {
	ctx := context.Background()

	// Create a test workspace with mods
	tmpDir := t.TempDir()
	setupComprehensiveTestWorkspace(t, tmpDir)

	t.Run("Compare Eager vs Lazy Tag Coverage", func(t *testing.T) {
		// EAGER LOADING (v1.4.3 baseline behavior)
		t.Log("=== Loading with EAGER mode (v1.4.3 baseline) ===")
		eagerWorkspace, errAndWarnings := Load(ctx, tmpDir)
		require.NoError(t, errAndWarnings.GetError(), "eager loading should succeed")
		defer eagerWorkspace.Close()

		eagerStats := extractTagStatistics(eagerWorkspace.GetModResources())
		t.Logf("EAGER: %d dashboards (%.1f%% with tags, %.1f%% with service, %.1f%% with mod)",
			eagerStats.TotalDashboards,
			eagerStats.DashboardTagCoverage,
			eagerStats.DashboardServiceCoverage,
			eagerStats.DashboardModCoverage)
		t.Logf("EAGER: %d benchmarks (%.1f%% with tags, %.1f%% with service, %.1f%% with mod)",
			eagerStats.TotalBenchmarks,
			eagerStats.BenchmarkTagCoverage,
			eagerStats.BenchmarkServiceCoverage,
			eagerStats.BenchmarkModCoverage)

		// LAZY LOADING (new implementation)
		t.Log("\n=== Loading with LAZY mode (new implementation) ===")
		lazyWorkspace, err := LoadLazy(ctx, tmpDir)
		require.NoError(t, err, "lazy loading should succeed")
		defer lazyWorkspace.Close()

		// Wait for background resolution
		completed := lazyWorkspace.WaitForResolution(5 * time.Second)
		require.True(t, completed, "background resolution should complete")

		// Build payload from lazy workspace
		payload := lazyWorkspace.GetAvailableDashboardsFromIndex()
		lazyStats := extractTagStatisticsFromPayload(payload)
		t.Logf("LAZY:  %d dashboards (%.1f%% with tags, %.1f%% with service, %.1f%% with mod)",
			lazyStats.TotalDashboards,
			lazyStats.DashboardTagCoverage,
			lazyStats.DashboardServiceCoverage,
			lazyStats.DashboardModCoverage)
		t.Logf("LAZY:  %d benchmarks (%.1f%% with tags, %.1f%% with service, %.1f%% with mod)",
			lazyStats.TotalBenchmarks,
			lazyStats.BenchmarkTagCoverage,
			lazyStats.BenchmarkServiceCoverage,
			lazyStats.BenchmarkModCoverage)

		// ASSERTIONS: Lazy should match eager
		t.Log("\n=== Comparing Eager vs Lazy ===")

		// Resource counts should match
		assert.Equal(t, eagerStats.TotalDashboards, lazyStats.TotalDashboards,
			"dashboard count should match")
		assert.Equal(t, eagerStats.TotalBenchmarks, lazyStats.TotalBenchmarks,
			"benchmark count should match")

		// Tag coverage should match (allow 1% tolerance for rounding)
		assert.InDelta(t, eagerStats.DashboardTagCoverage, lazyStats.DashboardTagCoverage, 1.0,
			"dashboard tag coverage should match eager mode")
		assert.InDelta(t, eagerStats.BenchmarkTagCoverage, lazyStats.BenchmarkTagCoverage, 1.0,
			"benchmark tag coverage should match eager mode")

		// Service tag coverage should match
		assert.InDelta(t, eagerStats.DashboardServiceCoverage, lazyStats.DashboardServiceCoverage, 1.0,
			"dashboard service tag coverage should match eager mode")
		assert.InDelta(t, eagerStats.BenchmarkServiceCoverage, lazyStats.BenchmarkServiceCoverage, 1.0,
			"benchmark service tag coverage should match eager mode")

		// Mod tag coverage - lazy should be BETTER (we now add mod tag automatically)
		assert.GreaterOrEqual(t, lazyStats.DashboardModCoverage, eagerStats.DashboardModCoverage,
			"lazy mode should have same or better dashboard mod tag coverage")
		assert.GreaterOrEqual(t, lazyStats.BenchmarkModCoverage, eagerStats.BenchmarkModCoverage,
			"lazy mode should have same or better benchmark mod tag coverage")

		// Overall tag coverage
		eagerOverall := (eagerStats.DashboardsWithTags + eagerStats.BenchmarksWithTags) * 100.0 /
			float64(eagerStats.TotalDashboards+eagerStats.TotalBenchmarks)
		lazyOverall := (lazyStats.DashboardsWithTags + lazyStats.BenchmarksWithTags) * 100.0 /
			float64(lazyStats.TotalDashboards+lazyStats.TotalBenchmarks)

		t.Logf("\nOVERALL TAG COVERAGE: Eager=%.1f%%, Lazy=%.1f%%", eagerOverall, lazyOverall)
		assert.InDelta(t, eagerOverall, lazyOverall, 1.0,
			"overall tag coverage should match eager mode")

		t.Log("\n✅ LAZY MODE MATCHES EAGER MODE (v1.4.3 baseline)")
	})
}

type TagStatistics struct {
	TotalDashboards          int
	DashboardsWithTags       float64
	DashboardsWithService    float64
	DashboardsWithMod        float64
	DashboardTagCoverage     float64
	DashboardServiceCoverage float64
	DashboardModCoverage     float64
	TotalBenchmarks          int
	BenchmarksWithTags       float64
	BenchmarksWithService    float64
	BenchmarksWithMod        float64
	BenchmarkTagCoverage     float64
	BenchmarkServiceCoverage float64
	BenchmarkModCoverage     float64
}

func extractTagStatistics(modResources modconfig.ModResources) TagStatistics {
	stats := TagStatistics{}

	_ = modResources.WalkResources(func(item modconfig.HclResource) (bool, error) {
		tags := item.GetTags()
		blockType := item.GetBlockType()

		switch blockType {
		case "dashboard":
			stats.TotalDashboards++
			if len(tags) > 0 {
				stats.DashboardsWithTags++
			}
			if _, hasService := tags["service"]; hasService {
				stats.DashboardsWithService++
			}
			if _, hasMod := tags["mod"]; hasMod {
				stats.DashboardsWithMod++
			}

		case "benchmark":
			stats.TotalBenchmarks++
			if len(tags) > 0 {
				stats.BenchmarksWithTags++
			}
			if _, hasService := tags["service"]; hasService {
				stats.BenchmarksWithService++
			}
			if _, hasMod := tags["mod"]; hasMod {
				stats.BenchmarksWithMod++
			}
		}

		return true, nil
	})

	// Calculate percentages
	if stats.TotalDashboards > 0 {
		stats.DashboardTagCoverage = stats.DashboardsWithTags / float64(stats.TotalDashboards) * 100
		stats.DashboardServiceCoverage = stats.DashboardsWithService / float64(stats.TotalDashboards) * 100
		stats.DashboardModCoverage = stats.DashboardsWithMod / float64(stats.TotalDashboards) * 100
	}
	if stats.TotalBenchmarks > 0 {
		stats.BenchmarkTagCoverage = stats.BenchmarksWithTags / float64(stats.TotalBenchmarks) * 100
		stats.BenchmarkServiceCoverage = stats.BenchmarksWithService / float64(stats.TotalBenchmarks) * 100
		stats.BenchmarkModCoverage = stats.BenchmarksWithMod / float64(stats.TotalBenchmarks) * 100
	}

	return stats
}

func extractTagStatisticsFromPayload(payload *resourceindex.AvailableDashboardsPayload) TagStatistics {
	stats := TagStatistics{}

	// Count dashboards
	stats.TotalDashboards = len(payload.Dashboards)
	for _, dash := range payload.Dashboards {
		if len(dash.Tags) > 0 {
			stats.DashboardsWithTags++
		}
		if _, hasService := dash.Tags["service"]; hasService {
			stats.DashboardsWithService++
		}
		if _, hasMod := dash.Tags["mod"]; hasMod {
			stats.DashboardsWithMod++
		}
	}

	// Count benchmarks
	stats.TotalBenchmarks = len(payload.Benchmarks)
	for _, bench := range payload.Benchmarks {
		if len(bench.Tags) > 0 {
			stats.BenchmarksWithTags++
		}
		if _, hasService := bench.Tags["service"]; hasService {
			stats.BenchmarksWithService++
		}
		if _, hasMod := bench.Tags["mod"]; hasMod {
			stats.BenchmarksWithMod++
		}
	}

	// Calculate percentages
	if stats.TotalDashboards > 0 {
		stats.DashboardTagCoverage = stats.DashboardsWithTags / float64(stats.TotalDashboards) * 100
		stats.DashboardServiceCoverage = stats.DashboardsWithService / float64(stats.TotalDashboards) * 100
		stats.DashboardModCoverage = stats.DashboardsWithMod / float64(stats.TotalDashboards) * 100
	}
	if stats.TotalBenchmarks > 0 {
		stats.BenchmarkTagCoverage = stats.BenchmarksWithTags / float64(stats.TotalBenchmarks) * 100
		stats.BenchmarkServiceCoverage = stats.BenchmarksWithService / float64(stats.TotalBenchmarks) * 100
		stats.BenchmarkModCoverage = stats.BenchmarksWithMod / float64(stats.TotalBenchmarks) * 100
	}

	return stats
}

// setupComprehensiveTestWorkspace creates a workspace with resources directly
// (no dependencies) to work with both eager and lazy loading
func setupComprehensiveTestWorkspace(t *testing.T, workspaceDir string) {
	t.Helper()

	// Main workspace mod.pp - no dependencies, all resources inline
	modContent := `mod "test_workspace" {
  title = "Test Workspace"
}
`
	err := os.WriteFile(filepath.Join(workspaceDir, "mod.pp"), []byte(modContent), 0644)
	require.NoError(t, err)

	// Dashboard WITH service tags
	dashboard1 := `dashboard "tagged_dashboard" {
  title = "Tagged Dashboard"
  tags = {
    service = "AWS S3"
    type    = "Report"
    env     = "prod"
  }
}
`
	err = os.WriteFile(filepath.Join(workspaceDir, "tagged_dashboard.pp"), []byte(dashboard1), 0644)
	require.NoError(t, err)

	// Dashboard WITHOUT service tags
	dashboard2 := `dashboard "untagged_dashboard" {
  title = "Untagged Dashboard"
  tags = {
    type = "Summary"
  }
}
`
	err = os.WriteFile(filepath.Join(workspaceDir, "untagged_dashboard.pp"), []byte(dashboard2), 0644)
	require.NoError(t, err)

	// Benchmark WITH service tags
	benchmark1 := `benchmark "tagged_benchmark" {
  title = "Tagged Benchmark"
  tags = {
    service = "AWS EC2"
    type    = "Compliance"
    framework = "CIS"
  }
  children = [
    control.tagged_control
  ]
}

control "tagged_control" {
  title = "Tagged Control"
  sql = "select 1 as result"
}
`
	err = os.WriteFile(filepath.Join(workspaceDir, "tagged_benchmark.pp"), []byte(benchmark1), 0644)
	require.NoError(t, err)

	// Benchmark WITHOUT service tags
	benchmark2 := `benchmark "untagged_benchmark" {
  title = "Untagged Benchmark"
  tags = {
    category = "Security"
  }
  children = [
    control.untagged_control
  ]
}

control "untagged_control" {
  title = "Untagged Control"
  sql = "select 1 as result"
}
`
	err = os.WriteFile(filepath.Join(workspaceDir, "untagged_benchmark.pp"), []byte(benchmark2), 0644)
	require.NoError(t, err)
}
