package workspace_test

import (
	"context"
	"encoding/json"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turbot/powerpipe/internal/workspace"
)

// comparisonTestdataDir returns the path to the testdata directory for comparison tests.
func comparisonTestdataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "testdata", "mods")
}

// TestResourceMetadata_EagerVsLazy_Identical verifies that resource metadata
// from lazy loading matches eager loading exactly.
// This test should FAIL initially, proving that lazy loading differs from eager loading.
func TestResourceMetadata_EagerVsLazy_Identical(t *testing.T) {
	modPath := filepath.Join(comparisonTestdataDir(), "generated", "medium")

	ctx := context.Background()

	// Load workspace eagerly
	eagerWs, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.Nil(t, ew.GetError(), "eager load should succeed")
	defer eagerWs.Close()

	// Load workspace lazily
	lazyWs, err := workspace.LoadLazy(ctx, modPath)
	require.NoError(t, err, "lazy load should succeed")
	defer lazyWs.Close()

	// Test dashboards
	t.Run("dashboards", func(t *testing.T) {
		eagerRes := eagerWs.GetPowerpipeModResources()
		lazyResources := lazyWs.GetLazyModResources()

		dashNames := lazyResources.ListDashboardNames()
		require.NotEmpty(t, dashNames, "should have dashboards")

		for _, name := range dashNames {
			t.Run(name, func(t *testing.T) {
				eagerDash := eagerRes.Dashboards[name]
				require.NotNil(t, eagerDash, "eager dashboard should exist")

				// Load lazy dashboard
				lazyDash, err := lazyWs.LoadDashboard(ctx, name)
				require.NoError(t, err, "lazy load should succeed")

				// Compare titles
				eagerTitle := ""
				if eagerDash.Title != nil {
					eagerTitle = *eagerDash.Title
				}
				lazyTitle := ""
				if lazyDash.Title != nil {
					lazyTitle = *lazyDash.Title
				}
				assert.Equal(t, eagerTitle, lazyTitle, "titles should match")

				// Compare tags - THIS IS WHERE WE EXPECT FAILURES
				// Lazy loading may not have all tag information
				assert.Equal(t, eagerDash.Tags, lazyDash.Tags,
					"tags should match for dashboard %s", name)
			})
		}
	})

	// Test benchmarks
	t.Run("benchmarks", func(t *testing.T) {
		eagerRes := eagerWs.GetPowerpipeModResources()
		lazyResources := lazyWs.GetLazyModResources()

		benchNames := lazyResources.ListBenchmarkNames()
		require.NotEmpty(t, benchNames, "should have benchmarks")

		for _, name := range benchNames {
			t.Run(name, func(t *testing.T) {
				eagerBench := eagerRes.ControlBenchmarks[name]
				require.NotNil(t, eagerBench, "eager benchmark should exist")

				// Load lazy benchmark
				lazyBench, err := lazyWs.LoadBenchmark(ctx, name)
				require.NoError(t, err, "lazy load should succeed")

				// Compare titles
				assert.Equal(t, eagerBench.GetTitle(), lazyBench.GetTitle(),
					"titles should match")

				// Compare tags - THIS IS WHERE WE EXPECT FAILURES
				assert.Equal(t, eagerBench.GetTags(), lazyBench.GetTags(),
					"tags should match for benchmark %s", name)
			})
		}
	})
}

// TestTagStructure_EagerVsLazy_Identical verifies that tag structures match exactly.
// Tags must have the same keys AND values.
func TestTagStructure_EagerVsLazy_Identical(t *testing.T) {
	modPath := filepath.Join(comparisonTestdataDir(), "generated", "medium")

	ctx := context.Background()

	eagerWs, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.Nil(t, ew.GetError())
	defer eagerWs.Close()

	lazyWs, err := workspace.LoadLazy(ctx, modPath)
	require.NoError(t, err)
	defer lazyWs.Close()

	eagerRes := eagerWs.GetPowerpipeModResources()

	// Test a dashboard with known tags
	for name, eagerDash := range eagerRes.Dashboards {
		if len(eagerDash.Tags) == 0 {
			continue // Skip dashboards without tags
		}

		t.Run(name, func(t *testing.T) {
			lazyDash, err := lazyWs.LoadDashboard(ctx, name)
			require.NoError(t, err)

			eagerTags := eagerDash.Tags
			lazyTags := lazyDash.Tags

			// Same number of tags
			assert.Equal(t, len(eagerTags), len(lazyTags),
				"should have same number of tags")

			// Same keys
			for key := range eagerTags {
				_, exists := lazyTags[key]
				assert.True(t, exists, "lazy tags missing key: %s", key)
			}

			// Same values
			for key, eagerVal := range eagerTags {
				assert.Equal(t, eagerVal, lazyTags[key],
					"tag value mismatch for key: %s", key)
			}
		})
	}
}

// TestDashboardListPayload_EagerVsLazy_Identical compares the full payload structure.
// This uses JSON comparison to ensure complete equivalence.
func TestDashboardListPayload_EagerVsLazy_Identical(t *testing.T) {
	modPath := filepath.Join(comparisonTestdataDir(), "generated", "medium")

	ctx := context.Background()

	// Load eager workspace
	eagerWs, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.Nil(t, ew.GetError())
	defer eagerWs.Close()

	// Load lazy workspace
	lazyWs, err := workspace.LoadLazy(ctx, modPath)
	require.NoError(t, err)
	defer lazyWs.Close()

	// Get payloads from index
	lazyPayload := lazyWs.GetAvailableDashboardsFromIndex()
	require.NotNil(t, lazyPayload)

	// Compare dashboard counts
	eagerRes := eagerWs.GetPowerpipeModResources()
	assert.Equal(t, len(eagerRes.Dashboards), len(lazyPayload.Dashboards),
		"should have same number of dashboards")

	// Compare benchmark counts
	assert.Equal(t, len(eagerRes.ControlBenchmarks), len(lazyPayload.Benchmarks),
		"should have same number of benchmarks")

	// Detailed comparison of each dashboard
	for name, eagerDash := range eagerRes.Dashboards {
		t.Run("dashboard/"+name, func(t *testing.T) {
			lazyDash, exists := lazyPayload.Dashboards[name]
			require.True(t, exists, "lazy payload should contain dashboard")

			// Compare title
			eagerTitle := ""
			if eagerDash.Title != nil {
				eagerTitle = *eagerDash.Title
			}
			assert.Equal(t, eagerTitle, lazyDash.Title, "titles should match")

			// Compare tags - THIS IS THE KEY COMPARISON
			assert.Equal(t, eagerDash.Tags, lazyDash.Tags,
				"tags should match for dashboard %s", name)

			// Compare short name
			assert.Equal(t, eagerDash.ShortName, lazyDash.ShortName,
				"short names should match")
		})
	}
}

// TestBenchmarkHierarchy_EagerVsLazy_Identical verifies that benchmark hierarchy
// (parent-child relationships) matches between eager and lazy loading.
func TestBenchmarkHierarchy_EagerVsLazy_Identical(t *testing.T) {
	modPath := filepath.Join(comparisonTestdataDir(), "generated", "medium")

	ctx := context.Background()

	eagerWs, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.Nil(t, ew.GetError())
	defer eagerWs.Close()

	lazyWs, err := workspace.LoadLazy(ctx, modPath)
	require.NoError(t, err)
	defer lazyWs.Close()

	eagerRes := eagerWs.GetPowerpipeModResources()
	lazyPayload := lazyWs.GetAvailableDashboardsFromIndex()

	for name, eagerBench := range eagerRes.ControlBenchmarks {
		t.Run(name, func(t *testing.T) {
			lazyBench, exists := lazyPayload.Benchmarks[name]
			require.True(t, exists, "lazy payload should contain benchmark")

			// Compare child counts
			eagerChildCount := 0
			for _, child := range eagerBench.GetChildren() {
				// Only count benchmark children (not controls)
				if _, ok := eagerRes.ControlBenchmarks[child.Name()]; ok {
					eagerChildCount++
				}
			}

			assert.Equal(t, eagerChildCount, len(lazyBench.Children),
				"child counts should match for benchmark %s", name)

			// Compare IsTopLevel flag
			eagerIsTopLevel := false
			for _, parent := range eagerBench.GetParents() {
				if parent.Name() == eagerWs.GetPowerpipeModResources().Mod.GetFullName() {
					eagerIsTopLevel = true
					break
				}
			}
			assert.Equal(t, eagerIsTopLevel, lazyBench.IsTopLevel,
				"IsTopLevel should match for benchmark %s", name)
		})
	}
}

// TestSourceDefinition_LazyLoaded_NotEmpty verifies that source definitions
// are available for lazy-loaded resources.
func TestSourceDefinition_LazyLoaded_NotEmpty(t *testing.T) {
	modPath := filepath.Join(comparisonTestdataDir(), "generated", "medium")

	ctx := context.Background()

	lazyWs, err := workspace.LoadLazy(ctx, modPath)
	require.NoError(t, err)
	defer lazyWs.Close()

	lazyResources := lazyWs.GetLazyModResources()
	dashNames := lazyResources.ListDashboardNames()
	require.NotEmpty(t, dashNames)

	// Test at least one dashboard
	name := dashNames[0]
	t.Run(name, func(t *testing.T) {
		dash, err := lazyWs.LoadDashboard(ctx, name)
		require.NoError(t, err)

		// Source definition should be populated for lazy-loaded resources
		// NOTE: This may fail if source_definition is not populated during lazy loading
		sourceDef := dash.GetSourceDefinition()
		assert.NotEmpty(t, sourceDef,
			"source_definition should be populated for lazy-loaded dashboard")

		// If present, verify it contains expected content
		if sourceDef != "" {
			assert.Contains(t, sourceDef, "dashboard",
				"source should contain 'dashboard' keyword")
		}
	})
}

// TestJSONPayload_EagerVsLazy_Identical does a full JSON comparison of payloads.
// This is the most comprehensive comparison, ensuring byte-for-byte equivalence
// when serialized.
func TestJSONPayload_EagerVsLazy_Identical(t *testing.T) {
	modPath := filepath.Join(comparisonTestdataDir(), "generated", "small")

	ctx := context.Background()

	eagerWs, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.Nil(t, ew.GetError())
	defer eagerWs.Close()

	lazyWs, err := workspace.LoadLazy(ctx, modPath)
	require.NoError(t, err)
	defer lazyWs.Close()

	// Build comparable structures
	eagerRes := eagerWs.GetPowerpipeModResources()
	lazyPayload := lazyWs.GetAvailableDashboardsFromIndex()

	// Build eager dashboard map in comparable format
	type comparableDash struct {
		Title     string            `json:"title"`
		ShortName string            `json:"short_name"`
		Tags      map[string]string `json:"tags"`
	}

	eagerDashMap := make(map[string]comparableDash)
	for name, d := range eagerRes.Dashboards {
		title := ""
		if d.Title != nil {
			title = *d.Title
		}
		eagerDashMap[name] = comparableDash{
			Title:     title,
			ShortName: d.ShortName,
			Tags:      d.Tags,
		}
	}

	lazyDashMap := make(map[string]comparableDash)
	for name, d := range lazyPayload.Dashboards {
		lazyDashMap[name] = comparableDash{
			Title:     d.Title,
			ShortName: d.ShortName,
			Tags:      d.Tags,
		}
	}

	// Marshal to JSON for comparison
	eagerJSON, err := json.MarshalIndent(eagerDashMap, "", "  ")
	require.NoError(t, err)

	lazyJSON, err := json.MarshalIndent(lazyDashMap, "", "  ")
	require.NoError(t, err)

	// Compare - THIS IS THE DEFINITIVE TEST
	// If this fails, lazy loading produces different output than eager loading
	assert.Equal(t, string(eagerJSON), string(lazyJSON),
		"JSON payloads should be identical")
}

// TestAllFieldsPresent_LazyPayload verifies that all expected fields
// are populated in the lazy-loaded payload.
func TestAllFieldsPresent_LazyPayload(t *testing.T) {
	modPath := filepath.Join(comparisonTestdataDir(), "generated", "medium")

	ctx := context.Background()

	lazyWs, err := workspace.LoadLazy(ctx, modPath)
	require.NoError(t, err)
	defer lazyWs.Close()

	payload := lazyWs.GetAvailableDashboardsFromIndex()
	require.NotNil(t, payload)

	// Check dashboards have all required fields
	t.Run("dashboards", func(t *testing.T) {
		for name, dash := range payload.Dashboards {
			t.Run(name, func(t *testing.T) {
				assert.NotEmpty(t, dash.FullName, "FullName must be present")
				assert.NotEmpty(t, dash.ShortName, "ShortName must be present")
				// Title may be empty, but Tags should not be nil if dashboard has tags
				// NOTE: This test documents expected behavior - may fail if not implemented
			})
		}
	})

	// Check benchmarks have all required fields
	t.Run("benchmarks", func(t *testing.T) {
		for name, bench := range payload.Benchmarks {
			t.Run(name, func(t *testing.T) {
				assert.NotEmpty(t, bench.FullName, "FullName must be present")
				assert.NotEmpty(t, bench.ShortName, "ShortName must be present")
				assert.NotEmpty(t, bench.BenchmarkType, "BenchmarkType must be present")
			})
		}
	})
}
