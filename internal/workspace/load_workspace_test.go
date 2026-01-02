package workspace

import (
	"context"
	"os"
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

func TestLoadSimpleMod(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "simple-mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())
	require.NotNil(t, w)

	// Verify mod loaded
	assert.Equal(t, "simple_test", w.Mod.ShortName)

	// Verify resources loaded
	res := w.GetPowerpipeModResources()
	assert.Len(t, res.Queries, 1, "should have 1 query")
	assert.Len(t, res.Dashboards, 1, "should have 1 dashboard")
	assert.Contains(t, res.Queries, "simple_test.query.simple_query")
	assert.Contains(t, res.Dashboards, "simple_test.dashboard.simple_dashboard")

	// Verify query content
	query := res.Queries["simple_test.query.simple_query"]
	require.NotNil(t, query)
	assert.NotNil(t, query.SQL)
	assert.Equal(t, "SELECT 1", *query.SQL)

	// Verify dashboard has children (the card)
	dashboard := res.Dashboards["simple_test.dashboard.simple_dashboard"]
	require.NotNil(t, dashboard)
	assert.NotEmpty(t, dashboard.GetChildren(), "dashboard should have children")
}

func TestLoadComplexMod(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "complex-mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())
	require.NotNil(t, w)

	res := w.GetPowerpipeModResources()

	// Verify all resource types
	assert.Len(t, res.Variables, 1, "should have 1 variable")
	assert.Len(t, res.Locals, 1, "should have 1 local")
	assert.Len(t, res.Queries, 1, "should have 1 query")
	assert.Len(t, res.Dashboards, 1, "should have 1 dashboard")
	assert.Len(t, res.Controls, 1, "should have 1 control")
	assert.Len(t, res.ControlBenchmarks, 1, "should have 1 benchmark")

	// Verify variable is loaded with default value
	v := res.Variables["complex_test.var.region"]
	require.NotNil(t, v)
	assert.Equal(t, "us-east-1", v.ValueGo)

	// Verify dashboard structure
	dashboard := res.Dashboards["complex_test.dashboard.complex_dashboard"]
	require.NotNil(t, dashboard)
	assert.NotEmpty(t, dashboard.GetChildren(), "dashboard should have children")

	// Verify control
	control := res.Controls["complex_test.control.test_control"]
	require.NotNil(t, control)
	assert.NotNil(t, control.Title)
	assert.Equal(t, "Test Control", *control.Title)

	// Verify benchmark
	benchmark := res.ControlBenchmarks["complex_test.benchmark.test_benchmark"]
	require.NotNil(t, benchmark)
	assert.Equal(t, 1, len(benchmark.GetChildren()), "benchmark should have 1 child")
}

func TestLoadBenchmarkHierarchy(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "benchmark-only")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Verify benchmark hierarchy
	parent := res.ControlBenchmarks["benchmark_only.benchmark.parent"]
	require.NotNil(t, parent)
	assert.Len(t, parent.GetChildren(), 2, "parent benchmark should have 2 children")

	// Verify children are correct type
	for _, child := range parent.GetChildren() {
		_, ok := child.(*resources.Benchmark)
		assert.True(t, ok, "child should be a Benchmark")
	}

	// Verify child benchmarks
	child1 := res.ControlBenchmarks["benchmark_only.benchmark.child_1"]
	require.NotNil(t, child1)
	assert.Len(t, child1.GetChildren(), 1, "child_1 should have 1 control child")

	child2 := res.ControlBenchmarks["benchmark_only.benchmark.child_2"]
	require.NotNil(t, child2)
	assert.Len(t, child2.GetChildren(), 1, "child_2 should have 1 control child")

	// Verify controls
	assert.Len(t, res.Controls, 2, "should have 2 controls")
	assert.Contains(t, res.Controls, "benchmark_only.control.ctrl_1")
	assert.Contains(t, res.Controls, "benchmark_only.control.ctrl_2")
}

func TestLoadModWithoutModfile(t *testing.T) {
	ctx := context.Background()
	modPath := t.TempDir()

	// Create a simple query file without a mod.pp
	queryFile := filepath.Join(modPath, "query.pp")
	err := os.WriteFile(queryFile, []byte(`query "test" { sql = "SELECT 1" }`), 0600)
	require.NoError(t, err)

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	// Should create default mod
	assert.Equal(t, "local", w.Mod.ShortName)

	// Query should still be loaded
	res := w.GetPowerpipeModResources()
	assert.Len(t, res.Queries, 1)
}

func TestLoadModResourceCount(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "complex-mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Count all resources to ensure nothing is missed
	totalResources := len(res.Queries) +
		len(res.Controls) +
		len(res.ControlBenchmarks) +
		len(res.Dashboards) +
		len(res.Variables) +
		len(res.Locals)

	assert.Greater(t, totalResources, 0, "should have loaded resources")
	// Expected: 1 query + 1 control + 1 benchmark + 1 dashboard + 1 variable + 1 local = 6
	assert.Equal(t, 6, totalResources, "should have exactly 6 top-level resources")
}

func TestLoadModIdempotent(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "simple-mod")

	// Load twice
	w1, ew1 := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew1.GetError())

	w2, ew2 := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew2.GetError())

	res1 := w1.GetPowerpipeModResources()
	res2 := w2.GetPowerpipeModResources()

	// Should produce identical results
	assert.Equal(t, len(res1.Queries), len(res2.Queries))
	assert.Equal(t, len(res1.Dashboards), len(res2.Dashboards))
	assert.Equal(t, len(res1.Controls), len(res2.Controls))
	assert.Equal(t, len(res1.ControlBenchmarks), len(res2.ControlBenchmarks))
}

func TestLoadModDashboardInputs(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "complex-mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()
	dashboard := res.Dashboards["complex_test.dashboard.complex_dashboard"]
	require.NotNil(t, dashboard)

	// Verify inputs are initialized
	assert.NotEmpty(t, res.DashboardInputs, "should have dashboard inputs")
}

func TestModResourceNames(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "simple-mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Verify resource names follow the pattern: mod_name.resource_type.resource_name
	for name := range res.Queries {
		assert.Contains(t, name, "simple_test.query.", "query name should contain mod prefix")
	}

	for name := range res.Dashboards {
		assert.Contains(t, name, "simple_test.dashboard.", "dashboard name should contain mod prefix")
	}
}

func TestLoadModWithOptions(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "simple-mod")

	// Test with variable validation disabled
	w, ew := Load(ctx, modPath, WithVariableValidation(false), WithLateBinding(true))
	require.NoError(t, ew.GetError())
	require.NotNil(t, w)

	// Verify options were applied
	assert.True(t, w.SupportLateBinding)
	assert.False(t, w.ValidateVariables)
}

func TestLoadModWithBlockTypeInclusions(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "benchmark-only")

	// Load with only benchmark and control types
	w, ew := Load(ctx, modPath, WithVariableValidation(false), WithBlockType([]string{"benchmark", "control"}))
	require.NoError(t, ew.GetError())
	require.NotNil(t, w)

	res := w.GetPowerpipeModResources()

	// Verify benchmarks and controls are loaded
	assert.NotEmpty(t, res.ControlBenchmarks, "should have benchmarks")
	assert.NotEmpty(t, res.Controls, "should have controls")
}
