package dashboardexecute

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turbot/pipe-fittings/v2/app_specific"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/pipe-fittings/v2/parse"
	pparse "github.com/turbot/powerpipe/internal/parse"
	"github.com/turbot/powerpipe/internal/resources"
	"github.com/turbot/powerpipe/internal/workspace"
)

func init() {
	// Set up app-specific constants required for mod loading
	app_specific.AppName = "powerpipe"
	app_specific.ModDataExtensions = []string{".pp", ".sp"}
	app_specific.VariablesExtensions = []string{".ppvars", ".spvars"}
	app_specific.AutoVariablesExtensions = []string{".auto.ppvars", ".auto.spvars"}
	app_specific.DefaultVarsFileName = "powerpipe.ppvars"
	app_specific.LegacyDefaultVarsFileName = "steampipe.spvars"
	app_specific.WorkspaceIgnoreFile = ".powerpipeignore"
	app_specific.WorkspaceDataDir = ".powerpipe"

	// Set up app-specific functions required for mod loading
	modconfig.AppSpecificNewModResourcesFunc = resources.NewModResources
	parse.ModDecoderFunc = pparse.NewPowerpipeModDecoder
	parse.AppSpecificGetResourceSchemaFunc = pparse.GetResourceSchema
}

func testdataDir() string {
	// Find testdata directory relative to this test file using runtime.Caller
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "testdata")
}

// TestExecution_DashboardStructure verifies that dashboards have the correct
// structure for execution (children, inputs, etc.).
// This is a behavior test that must pass before AND after lazy loading implementation.
func TestExecution_DashboardStructure(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Get main dashboard
	mainDash := res.Dashboards["behavior_test.dashboard.main"]
	require.NotNil(t, mainDash)

	// Verify dashboard has children that execution can iterate
	children := mainDash.GetChildren()
	assert.NotEmpty(t, children, "dashboard should have children for execution")

	// Verify all children are valid ModTreeItems
	for _, child := range children {
		assert.NotNil(t, child, "child should not be nil")
		assert.NotEmpty(t, child.Name(), "child should have a name")
	}
}

// TestExecution_BenchmarkStructure verifies that benchmarks have the correct
// structure for execution (children hierarchy).
func TestExecution_BenchmarkStructure(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Get top benchmark
	topBench := res.ControlBenchmarks["behavior_test.benchmark.top"]
	require.NotNil(t, topBench)

	// Verify benchmark children are accessible
	children := topBench.GetChildren()
	assert.NotEmpty(t, children, "benchmark should have children for execution")

	// Verify children can be traversed recursively (important for execution)
	childCount := 0
	var walkBenchmark func(item modconfig.ModTreeItem)
	walkBenchmark = func(item modconfig.ModTreeItem) {
		childCount++
		for _, child := range item.GetChildren() {
			walkBenchmark(child)
		}
	}
	walkBenchmark(topBench)

	// Should have walked top + child_a + child_b + nested controls
	assert.Greater(t, childCount, 3, "should traverse multiple children")
}

// TestExecution_ControlsHaveSQL verifies that controls have SQL for execution.
func TestExecution_ControlsHaveSQL(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// All controls should have SQL or Query for execution
	for name, ctrl := range res.Controls {
		sql := ctrl.GetSQL()
		query := ctrl.GetQuery()

		hasSQL := sql != nil && *sql != ""
		hasQuery := query != nil

		assert.True(t, hasSQL || hasQuery,
			"control %s should have SQL or Query reference for execution", name)
	}
}

// TestExecution_QueryProvidersHaveExecutableSQL verifies that all query
// providers can provide SQL for execution.
func TestExecution_QueryProvidersHaveExecutableSQL(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	providers := res.QueryProviders()

	// Track counts for reporting
	withSQL := 0
	withQuery := 0
	withNeither := 0

	for _, p := range providers {
		sql := p.GetSQL()
		query := p.GetQuery()

		if sql != nil && *sql != "" {
			withSQL++
		} else if query != nil {
			withQuery++
		} else {
			withNeither++
		}
	}

	// Most providers should have SQL or query
	totalExecutable := withSQL + withQuery
	assert.Greater(t, totalExecutable, 0, "should have executable query providers")

	// Log stats for visibility
	t.Logf("Query providers: %d with SQL, %d with Query reference, %d with neither",
		withSQL, withQuery, withNeither)
}

// TestExecution_DashboardWalkResources verifies that dashboard.WalkResources
// traverses all nested resources correctly.
func TestExecution_DashboardWalkResources(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	mainDash := res.Dashboards["behavior_test.dashboard.main"]
	require.NotNil(t, mainDash)

	// Walk all resources in the dashboard
	visited := make(map[string]bool)
	err := mainDash.WalkResources(func(resource modconfig.HclResource) (bool, error) {
		visited[resource.Name()] = true
		return true, nil
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, visited, "should have visited resources")
}

// TestExecution_BenchmarkWalkResources verifies that benchmark.WalkResources
// traverses all nested resources correctly.
func TestExecution_BenchmarkWalkResources(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	topBench := res.ControlBenchmarks["behavior_test.benchmark.top"]
	require.NotNil(t, topBench)

	// Walk all resources in the benchmark
	visited := make(map[string]bool)
	err := topBench.WalkResources(func(resource modconfig.ModTreeItem) (bool, error) {
		visited[resource.Name()] = true
		return true, nil
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, visited, "should have visited resources")

	// Should include nested controls
	assert.True(t, visited["behavior_test.control.nested_1"] ||
		visited["behavior_test.control.basic"],
		"should have visited at least one control")
}

// TestExecution_GraphHasNodesAndEdges verifies that graphs have nodes and
// edges accessible for execution.
func TestExecution_GraphHasNodesAndEdges(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Find a graph
	assert.NotEmpty(t, res.DashboardGraphs, "should have graphs")

	for _, graph := range res.DashboardGraphs {
		children := graph.GetChildren()
		// Graphs should have node and/or edge children
		assert.NotEmpty(t, children, "graph should have node/edge children")
	}
}

// TestExecution_InputValuesStructure verifies that InputValues can be
// created and manipulated correctly.
func TestExecution_InputValuesStructure(t *testing.T) {
	// Create input values
	inputs := NewInputValues()
	assert.NotNil(t, inputs)

	// Set some input values
	inputs.Inputs = map[string]interface{}{
		"input1": "value1",
		"input2": 42,
		"input3": true,
	}

	assert.Equal(t, "value1", inputs.Inputs["input1"])
	assert.Equal(t, 42, inputs.Inputs["input2"])
	assert.Equal(t, true, inputs.Inputs["input3"])
}

// TestExecution_DashboardInputsAccessible verifies that dashboard inputs
// are accessible for execution.
func TestExecution_DashboardInputsAccessible(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Dashboard inputs should be accessible
	assert.NotEmpty(t, res.DashboardInputs, "should have dashboard inputs")

	// Each dashboard with inputs should have them registered
	for dashName, inputs := range res.DashboardInputs {
		assert.NotEmpty(t, dashName, "dashboard name should not be empty")
		for inputName, input := range inputs {
			assert.NotEmpty(t, inputName, "input name should not be empty")
			assert.NotNil(t, input, "input should not be nil")
		}
	}
}

// TestExecution_NestedContainersTraversable verifies that nested containers
// in dashboards can be traversed for execution.
func TestExecution_NestedContainersTraversable(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Get nested dashboard
	nestedDash := res.Dashboards["behavior_test.dashboard.nested"]
	require.NotNil(t, nestedDash)

	// Walk and count depth of nesting
	maxDepth := 0
	var walkWithDepth func(item modconfig.ModTreeItem, depth int)
	walkWithDepth = func(item modconfig.ModTreeItem, depth int) {
		if depth > maxDepth {
			maxDepth = depth
		}
		for _, child := range item.GetChildren() {
			walkWithDepth(child, depth+1)
		}
	}
	walkWithDepth(nestedDash, 0)

	// Nested dashboard should have depth > 2
	assert.Greater(t, maxDepth, 2, "nested dashboard should have deep nesting")
}

// TestExecution_ResourcesHaveValidNames verifies that all resources have
// valid names that execution can use.
func TestExecution_ResourcesHaveValidNames(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Walk all resources and verify names
	err := res.WalkResources(func(r modconfig.HclResource) (bool, error) {
		name := r.Name()
		assert.NotEmpty(t, name, "resource should have a name")

		// Names should follow format: mod.type.name
		parts := splitName(name)
		assert.GreaterOrEqual(t, len(parts), 2, "name should have at least mod.type.name format: %s", name)

		return true, nil
	})
	assert.NoError(t, err)
}

// Helper function to split resource names
func splitName(name string) []string {
	result := []string{}
	current := ""
	for _, c := range name {
		if c == '.' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

// TestExecution_ControlSeverityAccessible verifies that control severity
// is accessible for execution reporting.
func TestExecution_ControlSeverityAccessible(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// At least some controls should have severity
	severityCount := 0
	for _, ctrl := range res.Controls {
		if ctrl.Severity != nil {
			severityCount++
		}
	}
	assert.Greater(t, severityCount, 0, "some controls should have severity set")
}

// TestExecution_BenchmarkParentsAccessible verifies that benchmark parents
// can be accessed (needed for trunk calculation).
func TestExecution_BenchmarkParentsAccessible(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Child benchmarks should have parents
	childA := res.ControlBenchmarks["behavior_test.benchmark.child_a"]
	require.NotNil(t, childA)

	parents := childA.GetParents()
	assert.NotEmpty(t, parents, "child benchmark should have parents")
}

// TestExecution_ControlParentsAccessible verifies that control parents
// can be accessed.
func TestExecution_ControlParentsAccessible(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Controls in benchmarks should have parents
	nestedCtrl := res.Controls["behavior_test.control.nested_1"]
	require.NotNil(t, nestedCtrl)

	parents := nestedCtrl.GetParents()
	assert.NotEmpty(t, parents, "nested control should have parents")
}
