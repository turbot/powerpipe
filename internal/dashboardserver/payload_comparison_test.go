package dashboardserver

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turbot/powerpipe/internal/workspace"
)

// comparisonTestdataDir returns the path to the testdata directory.
func comparisonTestdataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "testdata", "mods")
}

// skipIfModNotExists skips the test if the mod directory doesn't exist.
// Generated test mods are gitignored and only available for local testing.
func skipIfModNotExists(t *testing.T, modPath string) {
	if _, err := os.Stat(modPath); os.IsNotExist(err) {
		t.Skipf("Test mod not found (gitignored): %s", modPath)
	}
}

// TestDashboardListPayload_EagerVsLazy_Identical verifies that the dashboard list
// payload from lazy loading is identical to eager loading.
// This test should FAIL initially, proving the current gap between lazy and eager loading.
func TestDashboardListPayload_EagerVsLazy_Identical(t *testing.T) {
	modPath := filepath.Join(comparisonTestdataDir(), "generated", "medium")
	skipIfModNotExists(t, modPath)

	ctx := context.Background()

	// Load workspace eagerly
	eagerWs, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.Nil(t, ew.GetError(), "eager load should succeed")
	defer eagerWs.Close()

	// Build eager payload
	eagerPayloadBytes, err := buildAvailableDashboardsPayload(eagerWs.GetPowerpipeModResources())
	require.NoError(t, err)

	// Load workspace lazily
	lazyWs, err := workspace.LoadLazy(ctx, modPath)
	require.NoError(t, err)
	defer lazyWs.Close()

	// Build lazy payload from index
	lazyPayloadBytes, err := buildAvailableDashboardsPayloadFromIndex(lazyWs)
	require.NoError(t, err)

	// Parse both payloads for comparison
	var eagerPayload, lazyPayload AvailableDashboardsPayload
	require.NoError(t, json.Unmarshal(eagerPayloadBytes, &eagerPayload))
	require.NoError(t, json.Unmarshal(lazyPayloadBytes, &lazyPayload))

	// Compare dashboard counts
	assert.Equal(t, len(eagerPayload.Dashboards), len(lazyPayload.Dashboards),
		"should have same number of dashboards")

	// Compare benchmark counts
	assert.Equal(t, len(eagerPayload.Benchmarks), len(lazyPayload.Benchmarks),
		"should have same number of benchmarks")

	// Detailed comparison of each dashboard
	for name, eagerDash := range eagerPayload.Dashboards {
		t.Run("dashboard/"+name, func(t *testing.T) {
			lazyDash, exists := lazyPayload.Dashboards[name]
			require.True(t, exists, "lazy payload should contain dashboard %s", name)

			// Compare all fields
			assert.Equal(t, eagerDash.Title, lazyDash.Title,
				"titles should match for %s", name)
			assert.Equal(t, eagerDash.ShortName, lazyDash.ShortName,
				"short names should match for %s", name)
			assert.Equal(t, eagerDash.FullName, lazyDash.FullName,
				"full names should match for %s", name)
			assert.Equal(t, eagerDash.ModFullName, lazyDash.ModFullName,
				"mod full names should match for %s", name)

			// THIS IS THE KEY COMPARISON - tags must match exactly
			assert.Equal(t, eagerDash.Tags, lazyDash.Tags,
				"tags should match for dashboard %s", name)
		})
	}

	// Detailed comparison of each benchmark
	for name, eagerBench := range eagerPayload.Benchmarks {
		t.Run("benchmark/"+name, func(t *testing.T) {
			lazyBench, exists := lazyPayload.Benchmarks[name]
			require.True(t, exists, "lazy payload should contain benchmark %s", name)

			// Compare all fields
			assert.Equal(t, eagerBench.Title, lazyBench.Title,
				"titles should match for %s", name)
			assert.Equal(t, eagerBench.ShortName, lazyBench.ShortName,
				"short names should match for %s", name)
			assert.Equal(t, eagerBench.FullName, lazyBench.FullName,
				"full names should match for %s", name)
			assert.Equal(t, eagerBench.BenchmarkType, lazyBench.BenchmarkType,
				"benchmark types should match for %s", name)
			assert.Equal(t, eagerBench.IsTopLevel, lazyBench.IsTopLevel,
				"IsTopLevel should match for %s", name)

			// Tags must match exactly
			assert.Equal(t, eagerBench.Tags, lazyBench.Tags,
				"tags should match for benchmark %s", name)

			// Children count should match
			assert.Equal(t, len(eagerBench.Children), len(lazyBench.Children),
				"children count should match for %s", name)

			// Trunks should match
			assert.Equal(t, eagerBench.Trunks, lazyBench.Trunks,
				"trunks should match for %s", name)
		})
	}
}

// TestDashboardPayload_AllFieldsPopulated verifies that all expected fields
// are populated in the lazy payload (not nil/empty where they shouldn't be).
func TestDashboardPayload_AllFieldsPopulated(t *testing.T) {
	modPath := filepath.Join(comparisonTestdataDir(), "generated", "medium")
	skipIfModNotExists(t, modPath)

	ctx := context.Background()

	lazyWs, err := workspace.LoadLazy(ctx, modPath)
	require.NoError(t, err)
	defer lazyWs.Close()

	lazyPayloadBytes, err := buildAvailableDashboardsPayloadFromIndex(lazyWs)
	require.NoError(t, err)

	var payload AvailableDashboardsPayload
	require.NoError(t, json.Unmarshal(lazyPayloadBytes, &payload))

	// Verify action field
	assert.Equal(t, "available_dashboards", payload.Action)

	// Check dashboards
	for name, dash := range payload.Dashboards {
		t.Run("dashboard/"+name, func(t *testing.T) {
			assert.NotEmpty(t, dash.FullName, "FullName must be present")
			assert.NotEmpty(t, dash.ShortName, "ShortName must be present")
			// Title should be populated if it exists in HCL
			// Tags should not be nil (may be empty map)
		})
	}

	// Check benchmarks
	for name, bench := range payload.Benchmarks {
		t.Run("benchmark/"+name, func(t *testing.T) {
			assert.NotEmpty(t, bench.FullName, "FullName must be present")
			assert.NotEmpty(t, bench.ShortName, "ShortName must be present")
			assert.NotEmpty(t, bench.BenchmarkType, "BenchmarkType must be present")
		})
	}
}

// TestBenchmarkPayload_ChildHierarchy verifies that benchmark child hierarchy
// is correctly built from the index.
func TestBenchmarkPayload_ChildHierarchy(t *testing.T) {
	modPath := filepath.Join(comparisonTestdataDir(), "generated", "medium")
	skipIfModNotExists(t, modPath)

	ctx := context.Background()

	// Load both workspaces
	eagerWs, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.Nil(t, ew.GetError())
	defer eagerWs.Close()

	lazyWs, err := workspace.LoadLazy(ctx, modPath)
	require.NoError(t, err)
	defer lazyWs.Close()

	// Build payloads
	eagerPayloadBytes, err := buildAvailableDashboardsPayload(eagerWs.GetPowerpipeModResources())
	require.NoError(t, err)

	lazyPayloadBytes, err := buildAvailableDashboardsPayloadFromIndex(lazyWs)
	require.NoError(t, err)

	var eagerPayload, lazyPayload AvailableDashboardsPayload
	require.NoError(t, json.Unmarshal(eagerPayloadBytes, &eagerPayload))
	require.NoError(t, json.Unmarshal(lazyPayloadBytes, &lazyPayload))

	// Compare hierarchy for each benchmark
	for name, eagerBench := range eagerPayload.Benchmarks {
		lazyBench := lazyPayload.Benchmarks[name]

		// Compare child hierarchy recursively
		t.Run(name, func(t *testing.T) {
			compareChildHierarchy(t, name, eagerBench.Children, lazyBench.Children)
		})
	}
}

// compareChildHierarchy recursively compares benchmark children.
func compareChildHierarchy(t *testing.T, parentName string, eagerChildren, lazyChildren []ModAvailableBenchmark) {
	t.Helper()

	assert.Equal(t, len(eagerChildren), len(lazyChildren),
		"child count mismatch for %s", parentName)

	// Build maps for easier comparison (order may differ)
	eagerMap := make(map[string]ModAvailableBenchmark)
	for _, child := range eagerChildren {
		eagerMap[child.FullName] = child
	}

	lazyMap := make(map[string]ModAvailableBenchmark)
	for _, child := range lazyChildren {
		lazyMap[child.FullName] = child
	}

	// Compare each child
	for name, eagerChild := range eagerMap {
		lazyChild, exists := lazyMap[name]
		if !assert.True(t, exists, "lazy missing child %s", name) {
			continue
		}

		assert.Equal(t, eagerChild.Title, lazyChild.Title,
			"title mismatch for child %s", name)
		assert.Equal(t, eagerChild.Tags, lazyChild.Tags,
			"tags mismatch for child %s", name)

		// Recurse into grandchildren
		if len(eagerChild.Children) > 0 || len(lazyChild.Children) > 0 {
			compareChildHierarchy(t, name, eagerChild.Children, lazyChild.Children)
		}
	}
}

// TestPayload_JSONEquivalence does a complete JSON comparison.
// This is the definitive test for payload equivalence.
func TestPayload_JSONEquivalence(t *testing.T) {
	modPath := filepath.Join(comparisonTestdataDir(), "generated", "small")
	skipIfModNotExists(t, modPath)

	ctx := context.Background()

	eagerWs, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.Nil(t, ew.GetError())
	defer eagerWs.Close()

	lazyWs, err := workspace.LoadLazy(ctx, modPath)
	require.NoError(t, err)
	defer lazyWs.Close()

	eagerPayloadBytes, err := buildAvailableDashboardsPayload(eagerWs.GetPowerpipeModResources())
	require.NoError(t, err)

	lazyPayloadBytes, err := buildAvailableDashboardsPayloadFromIndex(lazyWs)
	require.NoError(t, err)

	// Parse and re-marshal for normalized comparison
	var eagerPayload, lazyPayload AvailableDashboardsPayload
	require.NoError(t, json.Unmarshal(eagerPayloadBytes, &eagerPayload))
	require.NoError(t, json.Unmarshal(lazyPayloadBytes, &lazyPayload))

	// Normalize by re-marshaling with sorted keys
	eagerNorm, _ := json.MarshalIndent(eagerPayload, "", "  ")
	lazyNorm, _ := json.MarshalIndent(lazyPayload, "", "  ")

	// This is the ultimate test - byte-for-byte equivalence
	// Note: This will fail if there are any differences including:
	// - Missing tags
	// - Missing titles
	// - Different hierarchy
	// - Different ordering (maps are unordered in Go, but JSON output should be consistent)

	if string(eagerNorm) != string(lazyNorm) {
		t.Errorf("Payloads differ:\n\nEager:\n%s\n\nLazy:\n%s", eagerNorm, lazyNorm)
	}
}

// TestTagsNotNil verifies that tags are never nil (should be empty map if no tags).
func TestTagsNotNil(t *testing.T) {
	modPath := filepath.Join(comparisonTestdataDir(), "generated", "medium")
	skipIfModNotExists(t, modPath)

	ctx := context.Background()

	lazyWs, err := workspace.LoadLazy(ctx, modPath)
	require.NoError(t, err)
	defer lazyWs.Close()

	lazyPayloadBytes, err := buildAvailableDashboardsPayloadFromIndex(lazyWs)
	require.NoError(t, err)

	var payload AvailableDashboardsPayload
	require.NoError(t, json.Unmarshal(lazyPayloadBytes, &payload))

	// Check all dashboards
	for name, dash := range payload.Dashboards {
		// Tags should be a map (possibly empty), never nil
		// This ensures consistent JSON serialization
		assert.NotNil(t, dash.Tags, "dashboard %s tags should not be nil", name)
	}

	// Check all benchmarks
	for name, bench := range payload.Benchmarks {
		assert.NotNil(t, bench.Tags, "benchmark %s tags should not be nil", name)
	}
}

// TestPayload_RealTagsFromMod verifies that tags are properly extracted from
// a mod with actual tag definitions using the phased_loading_comparison_mod.
func TestPayload_RealTagsFromMod(t *testing.T) {
	modPath := filepath.Join(comparisonTestdataDir(), "..", "..", "..",
		"tests", "acceptance", "test_data", "mods", "phased_loading_comparison_mod")

	ctx := context.Background()

	// Load workspace lazily
	lazyWs, err := workspace.LoadLazy(ctx, modPath)
	require.NoError(t, err)
	defer lazyWs.Close()

	// Build payload from index
	lazyPayloadBytes, err := buildAvailableDashboardsPayloadFromIndex(lazyWs)
	require.NoError(t, err)

	var payload AvailableDashboardsPayload
	require.NoError(t, json.Unmarshal(lazyPayloadBytes, &payload))

	// Find dashboard_with_tags
	dashWithTags, ok := payload.Dashboards["phased_loading_comparison.dashboard.dashboard_with_tags"]
	require.True(t, ok, "should have dashboard_with_tags")

	// Verify tags are populated
	assert.Equal(t, "Dashboard With Tags", dashWithTags.Title, "title should be populated")
	require.NotNil(t, dashWithTags.Tags, "tags should not be nil")
	assert.Equal(t, "test_service", dashWithTags.Tags["service"], "service tag should match")
	assert.Equal(t, "comparison", dashWithTags.Tags["category"], "category tag should match")
	assert.Equal(t, "acceptance_test", dashWithTags.Tags["type"], "type tag should match")

	// Find benchmark_with_tags
	benchWithTags, ok := payload.Benchmarks["phased_loading_comparison.benchmark.benchmark_with_tags"]
	require.True(t, ok, "should have benchmark_with_tags")

	// Verify benchmark tags are populated
	assert.Equal(t, "Benchmark With Tags", benchWithTags.Title, "benchmark title should be populated")
	require.NotNil(t, benchWithTags.Tags, "benchmark tags should not be nil")
	assert.Equal(t, "test_service", benchWithTags.Tags["service"], "benchmark service tag should match")
	assert.Equal(t, "comparison", benchWithTags.Tags["category"], "benchmark category tag should match")
}

// TestPayload_TagsMatchEagerLoading verifies that tags from lazy loading
// exactly match tags from eager loading for the phased_loading_comparison_mod.
func TestPayload_TagsMatchEagerLoading(t *testing.T) {
	modPath := filepath.Join(comparisonTestdataDir(), "..", "..", "..",
		"tests", "acceptance", "test_data", "mods", "phased_loading_comparison_mod")

	ctx := context.Background()

	// Load workspace eagerly
	eagerWs, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.Nil(t, ew.GetError(), "eager load should succeed")
	defer eagerWs.Close()

	// Build eager payload
	eagerPayloadBytes, err := buildAvailableDashboardsPayload(eagerWs.GetPowerpipeModResources())
	require.NoError(t, err)

	// Load workspace lazily
	lazyWs, err := workspace.LoadLazy(ctx, modPath)
	require.NoError(t, err)
	defer lazyWs.Close()

	// Build lazy payload
	lazyPayloadBytes, err := buildAvailableDashboardsPayloadFromIndex(lazyWs)
	require.NoError(t, err)

	// Parse payloads
	var eagerPayload, lazyPayload AvailableDashboardsPayload
	require.NoError(t, json.Unmarshal(eagerPayloadBytes, &eagerPayload))
	require.NoError(t, json.Unmarshal(lazyPayloadBytes, &lazyPayload))

	// Compare dashboard tags
	for name, eagerDash := range eagerPayload.Dashboards {
		lazyDash, ok := lazyPayload.Dashboards[name]
		require.True(t, ok, "lazy should have dashboard %s", name)
		assert.Equal(t, eagerDash.Tags, lazyDash.Tags,
			"tags should match for dashboard %s\neager: %v\nlazy: %v",
			name, eagerDash.Tags, lazyDash.Tags)
	}

	// Compare benchmark tags
	for name, eagerBench := range eagerPayload.Benchmarks {
		lazyBench, ok := lazyPayload.Benchmarks[name]
		require.True(t, ok, "lazy should have benchmark %s", name)
		assert.Equal(t, eagerBench.Tags, lazyBench.Tags,
			"tags should match for benchmark %s\neager: %v\nlazy: %v",
			name, eagerBench.Tags, lazyBench.Tags)
	}
}
