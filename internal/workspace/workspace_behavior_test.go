package workspace

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turbot/powerpipe/internal/resources"
)

// TestWorkspaceLoading_AllResourceTypesAccessible verifies that all resource types
// are properly loaded and accessible from the workspace.
// This is a behavior test that must pass before AND after lazy loading implementation.
func TestWorkspaceLoading_AllResourceTypesAccessible(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())
	require.NotNil(t, w)

	res := w.GetPowerpipeModResources()
	require.NotNil(t, res)

	// Verify mod metadata
	assert.Equal(t, "behavior_test", w.Mod.ShortName)
	assert.Equal(t, "Behavior Test Mod", *w.Mod.Title)

	// Verify Variables loaded
	assert.Len(t, res.Variables, 3, "should have 3 variables")
	assert.Contains(t, res.Variables, "behavior_test.var.region")
	assert.Contains(t, res.Variables, "behavior_test.var.count")
	assert.Contains(t, res.Variables, "behavior_test.var.enabled")

	// Verify variable values
	regionVar := res.Variables["behavior_test.var.region"]
	require.NotNil(t, regionVar)
	assert.Equal(t, "us-east-1", regionVar.ValueGo)

	// Verify Locals loaded (locals block creates individual local entries)
	assert.GreaterOrEqual(t, len(res.Locals), 2, "should have at least 2 locals")

	// Verify Queries loaded
	assert.GreaterOrEqual(t, len(res.Queries), 4, "should have at least 4 queries")
	assert.Contains(t, res.Queries, "behavior_test.query.simple")
	assert.Contains(t, res.Queries, "behavior_test.query.parameterized")
	assert.Contains(t, res.Queries, "behavior_test.query.for_control")
	assert.Contains(t, res.Queries, "behavior_test.query.for_table")

	// Verify query content
	simpleQuery := res.Queries["behavior_test.query.simple"]
	require.NotNil(t, simpleQuery)
	require.NotNil(t, simpleQuery.SQL)
	assert.Equal(t, "SELECT 1 as value", *simpleQuery.SQL)
	assert.NotNil(t, simpleQuery.Title)
	assert.Equal(t, "Simple Query", *simpleQuery.Title)

	// Verify parameterized query
	paramQuery := res.Queries["behavior_test.query.parameterized"]
	require.NotNil(t, paramQuery)
	params := paramQuery.GetParams()
	assert.Len(t, params, 2, "parameterized query should have 2 params")

	// Verify Controls loaded
	assert.GreaterOrEqual(t, len(res.Controls), 6, "should have at least 6 controls")
	assert.Contains(t, res.Controls, "behavior_test.control.basic")
	assert.Contains(t, res.Controls, "behavior_test.control.uses_query")
	assert.Contains(t, res.Controls, "behavior_test.control.with_params")

	// Verify control properties
	basicControl := res.Controls["behavior_test.control.basic"]
	require.NotNil(t, basicControl)
	assert.NotNil(t, basicControl.Title)
	assert.Equal(t, "Basic Control", *basicControl.Title)
	assert.NotNil(t, basicControl.SQL)

	// Verify Benchmarks loaded
	assert.GreaterOrEqual(t, len(res.ControlBenchmarks), 4, "should have at least 4 benchmarks")
	assert.Contains(t, res.ControlBenchmarks, "behavior_test.benchmark.top")
	assert.Contains(t, res.ControlBenchmarks, "behavior_test.benchmark.child_a")
	assert.Contains(t, res.ControlBenchmarks, "behavior_test.benchmark.child_b")
	assert.Contains(t, res.ControlBenchmarks, "behavior_test.benchmark.flat")

	// Verify Dashboards loaded
	assert.GreaterOrEqual(t, len(res.Dashboards), 8, "should have at least 8 dashboards")
	assert.Contains(t, res.Dashboards, "behavior_test.dashboard.main")
	assert.Contains(t, res.Dashboards, "behavior_test.dashboard.nested")
	assert.Contains(t, res.Dashboards, "behavior_test.dashboard.with_graph")
	assert.Contains(t, res.Dashboards, "behavior_test.dashboard.with_flow")
	assert.Contains(t, res.Dashboards, "behavior_test.dashboard.with_hierarchy")
	assert.Contains(t, res.Dashboards, "behavior_test.dashboard.with_categories")
	assert.Contains(t, res.Dashboards, "behavior_test.dashboard.with_inputs")
	assert.Contains(t, res.Dashboards, "behavior_test.dashboard.simple")

	// Verify dashboard properties
	mainDash := res.Dashboards["behavior_test.dashboard.main"]
	require.NotNil(t, mainDash)
	assert.NotNil(t, mainDash.Title)
	assert.Equal(t, "Main Dashboard", *mainDash.Title)

	// Verify Cards loaded (from dashboards)
	assert.NotEmpty(t, res.DashboardCards, "should have dashboard cards")

	// Verify Charts loaded (from dashboards)
	assert.NotEmpty(t, res.DashboardCharts, "should have dashboard charts")

	// Verify Containers loaded (from dashboards)
	assert.NotEmpty(t, res.DashboardContainers, "should have dashboard containers")

	// Verify Tables loaded (from dashboards)
	assert.NotEmpty(t, res.DashboardTables, "should have dashboard tables")

	// Verify Texts loaded (from dashboards)
	assert.NotEmpty(t, res.DashboardTexts, "should have dashboard texts")

	// Verify Images loaded (from dashboards)
	assert.NotEmpty(t, res.DashboardImages, "should have dashboard images")

	// Verify Graphs loaded (from dashboards)
	assert.NotEmpty(t, res.DashboardGraphs, "should have dashboard graphs")

	// Verify Flows loaded (from dashboards)
	assert.NotEmpty(t, res.DashboardFlows, "should have dashboard flows")

	// Verify Hierarchies loaded (from dashboards)
	assert.NotEmpty(t, res.DashboardHierarchies, "should have dashboard hierarchies")

	// Note: Anonymous nodes/edges/categories within graphs may not be in top-level maps
	// They are accessible via the graph's GetChildren() method instead
	// The important thing is that graphs, flows, and hierarchies exist

	// Verify Inputs loaded
	assert.NotEmpty(t, res.DashboardInputs, "should have dashboard inputs")
}

// TestWorkspaceLoading_ResourceReferencesResolved verifies that cross-resource
// references are properly resolved during workspace loading.
func TestWorkspaceLoading_ResourceReferencesResolved(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Verify control that references a query has the query resolved
	ctrl := res.Controls["behavior_test.control.uses_query"]
	require.NotNil(t, ctrl, "control.uses_query should exist")

	query := ctrl.GetQuery()
	assert.NotNil(t, query, "Control's query reference should be resolved")
}

// TestWorkspaceLoading_NestedResourcesAccessible verifies that nested resources
// (like dashboard children) are properly accessible.
func TestWorkspaceLoading_NestedResourcesAccessible(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Test nested dashboard
	nestedDash := res.Dashboards["behavior_test.dashboard.nested"]
	require.NotNil(t, nestedDash)

	children := nestedDash.GetChildren()
	assert.NotEmpty(t, children, "nested dashboard should have children")

	// Verify we can traverse to nested children
	for _, child := range children {
		if container, ok := child.(*resources.DashboardContainer); ok {
			nestedChildren := container.GetChildren()
			assert.NotEmpty(t, nestedChildren, "container should have nested children")
		}
	}
}

// TestWorkspaceLoading_BenchmarkHierarchy verifies that benchmark hierarchy
// with nested children is properly resolved.
func TestWorkspaceLoading_BenchmarkHierarchy(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Verify top-level benchmark
	topBenchmark := res.ControlBenchmarks["behavior_test.benchmark.top"]
	require.NotNil(t, topBenchmark)

	topChildren := topBenchmark.GetChildren()
	assert.Len(t, topChildren, 3, "top benchmark should have 3 children (2 benchmarks + 1 control)")

	// Verify all children are resolved (not nil)
	for _, child := range topChildren {
		assert.NotNil(t, child, "benchmark child should not be nil")
	}

	// Verify child benchmarks have their children resolved
	childA := res.ControlBenchmarks["behavior_test.benchmark.child_a"]
	require.NotNil(t, childA)
	assert.Len(t, childA.GetChildren(), 2, "child_a should have 2 control children")

	childB := res.ControlBenchmarks["behavior_test.benchmark.child_b"]
	require.NotNil(t, childB)
	assert.Len(t, childB.GetChildren(), 1, "child_b should have 1 control child")

	// Verify child types are correct
	for _, child := range childA.GetChildren() {
		_, isControl := child.(*resources.Control)
		assert.True(t, isControl, "child_a children should be controls")
	}
}

// TestWorkspaceLoading_DashboardChildren verifies that dashboard children
// (cards, charts, containers, etc.) are properly accessible.
func TestWorkspaceLoading_DashboardChildren(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Get main dashboard
	mainDash := res.Dashboards["behavior_test.dashboard.main"]
	require.NotNil(t, mainDash)

	children := mainDash.GetChildren()
	assert.NotEmpty(t, children, "main dashboard should have children")

	// Track child types found
	foundInput := false
	foundContainer := false
	foundText := false
	foundImage := false

	for _, child := range children {
		switch child.(type) {
		case *resources.DashboardInput:
			foundInput = true
		case *resources.DashboardContainer:
			foundContainer = true
		case *resources.DashboardText:
			foundText = true
		case *resources.DashboardImage:
			foundImage = true
		}
	}

	assert.True(t, foundInput, "should have input child")
	assert.True(t, foundContainer, "should have container child")
	assert.True(t, foundText, "should have text child")
	assert.True(t, foundImage, "should have image child")
}

// TestWorkspaceLoading_GraphComponents verifies that graph visualizations
// with nodes, edges, and categories are properly loaded.
func TestWorkspaceLoading_GraphComponents(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Get graph dashboard
	graphDash := res.Dashboards["behavior_test.dashboard.with_graph"]
	require.NotNil(t, graphDash)

	children := graphDash.GetChildren()
	assert.NotEmpty(t, children, "graph dashboard should have children")

	// Find the graph
	var graph *resources.DashboardGraph
	for _, child := range children {
		if g, ok := child.(*resources.DashboardGraph); ok {
			graph = g
			break
		}
	}
	require.NotNil(t, graph, "should find a graph in the dashboard")

	// Verify graph has nodes
	graphChildren := graph.GetChildren()
	assert.NotEmpty(t, graphChildren, "graph should have node/edge children")
}

// TestWorkspaceLoading_FlowComponents verifies that flow visualizations
// are properly loaded.
func TestWorkspaceLoading_FlowComponents(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Get flow dashboard
	flowDash := res.Dashboards["behavior_test.dashboard.with_flow"]
	require.NotNil(t, flowDash)

	children := flowDash.GetChildren()
	assert.NotEmpty(t, children, "flow dashboard should have children")

	// Find the flow
	var flow *resources.DashboardFlow
	for _, child := range children {
		if f, ok := child.(*resources.DashboardFlow); ok {
			flow = f
			break
		}
	}
	require.NotNil(t, flow, "should find a flow in the dashboard")
}

// TestWorkspaceLoading_HierarchyComponents verifies that hierarchy visualizations
// are properly loaded.
func TestWorkspaceLoading_HierarchyComponents(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Get hierarchy dashboard
	hierarchyDash := res.Dashboards["behavior_test.dashboard.with_hierarchy"]
	require.NotNil(t, hierarchyDash)

	children := hierarchyDash.GetChildren()
	assert.NotEmpty(t, children, "hierarchy dashboard should have children")

	// Find the hierarchy
	var hierarchy *resources.DashboardHierarchy
	for _, child := range children {
		if h, ok := child.(*resources.DashboardHierarchy); ok {
			hierarchy = h
			break
		}
	}
	require.NotNil(t, hierarchy, "should find a hierarchy in the dashboard")
}

// TestWorkspaceLoading_Categories verifies that graph categories are accessible
// through the graph resources.
func TestWorkspaceLoading_Categories(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Graphs with categories exist
	assert.NotEmpty(t, res.DashboardGraphs, "should have graphs")

	// Categories are defined within graphs - verify by checking graph children
	graphFound := false
	for _, graph := range res.DashboardGraphs {
		children := graph.GetChildren()
		if len(children) > 0 {
			graphFound = true
		}
	}
	assert.True(t, graphFound, "should have graphs with children")
}

// TestWorkspaceLoading_VariableValues verifies that variable default values
// are properly resolved.
func TestWorkspaceLoading_VariableValues(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Test string variable
	regionVar := res.Variables["behavior_test.var.region"]
	require.NotNil(t, regionVar)
	assert.Equal(t, "us-east-1", regionVar.ValueGo)

	// Test number variable
	countVar := res.Variables["behavior_test.var.count"]
	require.NotNil(t, countVar)
	assert.Equal(t, 10, countVar.ValueGo)

	// Test bool variable
	enabledVar := res.Variables["behavior_test.var.enabled"]
	require.NotNil(t, enabledVar)
	assert.Equal(t, true, enabledVar.ValueGo)
}

// TestWorkspaceLoading_ResourceMetadata verifies that resource metadata
// (title, description, tags) is properly loaded.
func TestWorkspaceLoading_ResourceMetadata(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Test query metadata
	simpleQuery := res.Queries["behavior_test.query.simple"]
	require.NotNil(t, simpleQuery)
	assert.NotNil(t, simpleQuery.Title)
	assert.Equal(t, "Simple Query", *simpleQuery.Title)
	assert.NotNil(t, simpleQuery.Description)
	assert.Equal(t, "A simple SELECT query", *simpleQuery.Description)
	assert.NotEmpty(t, simpleQuery.Tags, "query should have tags")

	// Test dashboard metadata
	mainDash := res.Dashboards["behavior_test.dashboard.main"]
	require.NotNil(t, mainDash)
	assert.NotNil(t, mainDash.Title)
	assert.Equal(t, "Main Dashboard", *mainDash.Title)
	assert.NotNil(t, mainDash.Description)
	assert.NotEmpty(t, mainDash.Tags, "dashboard should have tags")

	// Test control metadata
	basicControl := res.Controls["behavior_test.control.basic"]
	require.NotNil(t, basicControl)
	assert.NotNil(t, basicControl.Title)
	assert.Equal(t, "Basic Control", *basicControl.Title)
	assert.NotNil(t, basicControl.Description)
	assert.NotEmpty(t, basicControl.Tags, "control should have tags")
}

// TestWorkspaceLoading_ControlSeverity verifies that control severity is loaded.
func TestWorkspaceLoading_ControlSeverity(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	basicControl := res.Controls["behavior_test.control.basic"]
	require.NotNil(t, basicControl)
	assert.NotNil(t, basicControl.Severity)
	assert.Equal(t, "low", *basicControl.Severity)

	usesQueryControl := res.Controls["behavior_test.control.uses_query"]
	require.NotNil(t, usesQueryControl)
	assert.NotNil(t, usesQueryControl.Severity)
	assert.Equal(t, "medium", *usesQueryControl.Severity)
}

// TestWorkspaceLoading_QueryParams verifies that query parameters are properly loaded.
func TestWorkspaceLoading_QueryParams(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	paramQuery := res.Queries["behavior_test.query.parameterized"]
	require.NotNil(t, paramQuery)

	params := paramQuery.GetParams()
	require.Len(t, params, 2, "should have 2 parameters")

	// Verify param names
	paramNames := make([]string, len(params))
	for i, p := range params {
		paramNames[i] = p.ShortName
	}
	assert.Contains(t, paramNames, "region")
	assert.Contains(t, paramNames, "min_count")
}

// TestWorkspaceLoading_Idempotent verifies that loading the same workspace
// multiple times produces identical results.
func TestWorkspaceLoading_Idempotent(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	// Load twice
	w1, ew1 := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew1.GetError())

	w2, ew2 := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew2.GetError())

	res1 := w1.GetPowerpipeModResources()
	res2 := w2.GetPowerpipeModResources()

	// Verify counts match
	assert.Equal(t, len(res1.Queries), len(res2.Queries), "query count should match")
	assert.Equal(t, len(res1.Controls), len(res2.Controls), "control count should match")
	assert.Equal(t, len(res1.ControlBenchmarks), len(res2.ControlBenchmarks), "benchmark count should match")
	assert.Equal(t, len(res1.Dashboards), len(res2.Dashboards), "dashboard count should match")
	assert.Equal(t, len(res1.Variables), len(res2.Variables), "variable count should match")
	assert.Equal(t, len(res1.DashboardCards), len(res2.DashboardCards), "card count should match")
	assert.Equal(t, len(res1.DashboardCharts), len(res2.DashboardCharts), "chart count should match")
}
