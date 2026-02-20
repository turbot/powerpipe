package workspace

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/powerpipe/internal/resources"
)

// TestResourceAccess_GetResourceByName verifies that resources can be looked up
// by their fully qualified name using GetResource.
// This is a behavior test that must pass before AND after lazy loading implementation.
func TestResourceAccess_GetResourceByName(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	testCases := []struct {
		resourceType string
		name         string
		expectFound  bool
	}{
		// Dashboards
		{"dashboard", "behavior_test.dashboard.main", true},
		{"dashboard", "behavior_test.dashboard.nested", true},
		{"dashboard", "behavior_test.dashboard.simple", true},
		{"dashboard", "behavior_test.dashboard.nonexistent", false},

		// Queries
		{"query", "behavior_test.query.simple", true},
		{"query", "behavior_test.query.parameterized", true},
		{"query", "behavior_test.query.for_control", true},
		{"query", "behavior_test.query.nonexistent", false},

		// Controls
		{"control", "behavior_test.control.basic", true},
		{"control", "behavior_test.control.uses_query", true},
		{"control", "behavior_test.control.with_params", true},
		{"control", "behavior_test.control.nonexistent", false},

		// Benchmarks
		{"benchmark", "behavior_test.benchmark.top", true},
		{"benchmark", "behavior_test.benchmark.child_a", true},
		{"benchmark", "behavior_test.benchmark.flat", true},
		{"benchmark", "behavior_test.benchmark.nonexistent", false},
	}

	for _, tc := range testCases {
		t.Run(tc.resourceType+"/"+tc.name, func(t *testing.T) {
			parsed := &modconfig.ParsedResourceName{
				Mod:      "behavior_test",
				ItemType: tc.resourceType,
				Name:     tc.name[len("behavior_test."+tc.resourceType+"."):], // extract short name
			}

			resource, found := res.GetResource(parsed)

			if tc.expectFound {
				assert.True(t, found, "Resource should be found: %s", tc.name)
				assert.NotNil(t, resource)
				assert.Equal(t, tc.name, resource.Name())
			} else {
				assert.False(t, found, "Resource should not be found: %s", tc.name)
			}
		})
	}
}

// TestResourceAccess_GetResourceByType verifies that GetResource works for
// various resource types.
func TestResourceAccess_GetResourceByType(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Test different resource types
	resourceTypes := []struct {
		itemType   string
		name       string
		shouldFind bool
	}{
		{"dashboard", "main", true},
		{"query", "simple", true},
		{"control", "basic", true},
		{"benchmark", "top", true},
	}

	for _, rt := range resourceTypes {
		t.Run(rt.itemType, func(t *testing.T) {
			parsed := &modconfig.ParsedResourceName{
				Mod:      "behavior_test",
				ItemType: rt.itemType,
				Name:     rt.name,
			}

			resource, found := res.GetResource(parsed)
			assert.Equal(t, rt.shouldFind, found, "expected found=%v for %s.%s", rt.shouldFind, rt.itemType, rt.name)
			if rt.shouldFind {
				assert.NotNil(t, resource)
			}
		})
	}
}

// TestResourceAccess_WalkResources verifies that WalkResources visits all
// resources in the mod.
func TestResourceAccess_WalkResources(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	var visited []string
	err := res.WalkResources(func(r modconfig.HclResource) (bool, error) {
		visited = append(visited, r.Name())
		return true, nil
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, visited)

	// Verify expected resources were visited
	assert.Contains(t, visited, "behavior_test.dashboard.main")
	assert.Contains(t, visited, "behavior_test.query.simple")
	assert.Contains(t, visited, "behavior_test.control.basic")
	assert.Contains(t, visited, "behavior_test.benchmark.top")
}

// TestResourceAccess_WalkResources_EarlyExit verifies that WalkResources
// can be stopped early by returning false.
func TestResourceAccess_WalkResources_EarlyExit(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	visitCount := 0
	maxVisits := 5

	err := res.WalkResources(func(r modconfig.HclResource) (bool, error) {
		visitCount++
		// Stop after maxVisits
		return visitCount < maxVisits, nil
	})

	assert.NoError(t, err)
	assert.Equal(t, maxVisits, visitCount, "should stop after %d visits", maxVisits)
}

// TestResourceAccess_QueryProviders verifies that QueryProviders returns
// all resources that can provide SQL queries.
func TestResourceAccess_QueryProviders(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	providers := res.QueryProviders()
	assert.NotEmpty(t, providers, "should have query providers")

	// Track types of query providers found
	foundQuery := false
	foundControl := false
	foundCard := false
	foundChart := false
	foundTable := false

	for _, p := range providers {
		switch p.(type) {
		case *resources.Query:
			foundQuery = true
		case *resources.Control:
			foundControl = true
		case *resources.DashboardCard:
			foundCard = true
		case *resources.DashboardChart:
			foundChart = true
		case *resources.DashboardTable:
			foundTable = true
		}
	}

	assert.True(t, foundQuery, "should have Query providers")
	assert.True(t, foundControl, "should have Control providers")
	assert.True(t, foundCard, "should have DashboardCard providers")
	assert.True(t, foundChart, "should have DashboardChart providers")
	assert.True(t, foundTable, "should have DashboardTable providers")
}

// TestResourceAccess_QueryProviders_HaveSQL verifies that QueryProviders
// return resources with SQL or Query references.
func TestResourceAccess_QueryProviders_HaveSQL(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	providers := res.QueryProviders()

	// Count providers with SQL or Query
	withSQLOrQuery := 0
	for _, p := range providers {
		sql := p.GetSQL()
		query := p.GetQuery()
		if (sql != nil && *sql != "") || query != nil {
			withSQLOrQuery++
		}
	}

	// Most query providers should have SQL or query reference
	assert.Greater(t, withSQLOrQuery, 0, "should have providers with SQL or Query")
}

// TestResourceAccess_TopLevelResources verifies that TopLevelResources
// returns only resources from the current mod.
func TestResourceAccess_TopLevelResources(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	topLevel := res.TopLevelResources()
	require.NotNil(t, topLevel)

	// Cast to PowerpipeModResources
	topRes, ok := topLevel.(*resources.PowerpipeModResources)
	require.True(t, ok)

	// Verify it has resources
	assert.NotEmpty(t, topRes.Dashboards, "should have top-level dashboards")
	assert.NotEmpty(t, topRes.Queries, "should have top-level queries")
	assert.NotEmpty(t, topRes.Controls, "should have top-level controls")
	assert.NotEmpty(t, topRes.ControlBenchmarks, "should have top-level benchmarks")
}

// TestResourceAccess_DirectMapAccess verifies that resources can be
// accessed directly via map lookup.
func TestResourceAccess_DirectMapAccess(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Direct map access should work
	dashboard := res.Dashboards["behavior_test.dashboard.main"]
	assert.NotNil(t, dashboard)
	assert.Equal(t, "behavior_test.dashboard.main", dashboard.Name())

	query := res.Queries["behavior_test.query.simple"]
	assert.NotNil(t, query)
	assert.Equal(t, "behavior_test.query.simple", query.Name())

	control := res.Controls["behavior_test.control.basic"]
	assert.NotNil(t, control)
	assert.Equal(t, "behavior_test.control.basic", control.Name())

	benchmark := res.ControlBenchmarks["behavior_test.benchmark.top"]
	assert.NotNil(t, benchmark)
	assert.Equal(t, "behavior_test.benchmark.top", benchmark.Name())
}

// TestResourceAccess_DashboardInputsNested verifies that dashboard inputs
// are organized by dashboard name.
func TestResourceAccess_DashboardInputsNested(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// DashboardInputs is nested by dashboard name
	assert.NotEmpty(t, res.DashboardInputs, "should have dashboard inputs")

	// Inputs should be organized by dashboard
	for dashName, inputs := range res.DashboardInputs {
		assert.NotEmpty(t, dashName, "dashboard name should not be empty")
		assert.NotEmpty(t, inputs, "inputs map should not be empty for %s", dashName)
	}
}

// TestResourceAccess_ResourceNames verifies that resource names follow
// the expected pattern.
func TestResourceAccess_ResourceNames(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// All resources should have names in format: mod_name.type.name
	for name, dashboard := range res.Dashboards {
		assert.Contains(t, name, "behavior_test.dashboard.")
		assert.Equal(t, name, dashboard.Name())
	}

	for name, query := range res.Queries {
		assert.Contains(t, name, "behavior_test.query.")
		assert.Equal(t, name, query.Name())
	}

	for name, control := range res.Controls {
		assert.Contains(t, name, "behavior_test.control.")
		assert.Equal(t, name, control.Name())
	}

	for name, benchmark := range res.ControlBenchmarks {
		assert.Contains(t, name, "behavior_test.benchmark.")
		assert.Equal(t, name, benchmark.Name())
	}
}

// TestResourceAccess_ResourceCount verifies that we can count all resources.
func TestResourceAccess_ResourceCount(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Count total resources via WalkResources
	totalFromWalk := 0
	err := res.WalkResources(func(r modconfig.HclResource) (bool, error) {
		totalFromWalk++
		return true, nil
	})
	assert.NoError(t, err)

	// Count via direct map access
	totalFromMaps := len(res.Dashboards) +
		len(res.Queries) +
		len(res.Controls) +
		len(res.ControlBenchmarks) +
		len(res.Variables) +
		len(res.Locals) +
		len(res.DashboardCards) +
		len(res.DashboardCharts) +
		len(res.DashboardContainers) +
		len(res.DashboardTables) +
		len(res.DashboardTexts) +
		len(res.DashboardImages) +
		len(res.DashboardGraphs) +
		len(res.DashboardFlows) +
		len(res.DashboardHierarchies) +
		len(res.DashboardNodes) +
		len(res.DashboardEdges) +
		len(res.DashboardCategories) +
		len(res.GlobalDashboardInputs) +
		len(res.Mods)

	// Count nested inputs
	for _, inputs := range res.DashboardInputs {
		totalFromMaps += len(inputs)
	}

	// WalkResources should visit at least as many resources as in the maps
	// (it may visit more due to Mods being counted)
	assert.GreaterOrEqual(t, totalFromWalk, totalFromMaps-len(res.Mods),
		"WalkResources should visit all resources")
}

// TestResourceAccess_EmptyMapAccess verifies that accessing non-existent
// resources returns nil/zero values.
func TestResourceAccess_EmptyMapAccess(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Accessing non-existent keys should return nil
	assert.Nil(t, res.Dashboards["nonexistent"])
	assert.Nil(t, res.Queries["nonexistent"])
	assert.Nil(t, res.Controls["nonexistent"])
	assert.Nil(t, res.ControlBenchmarks["nonexistent"])
}
