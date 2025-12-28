# Task 2: Create Mod Loading Test Suite

## Objective

Create a comprehensive test suite for mod loading functionality to ensure correctness is maintained during performance optimizations.

## Context

- Performance optimizations will change core loading code paths
- Need safety net tests before making changes
- Tests should cover various mod configurations
- Tests should verify resource parsing correctness

## Dependencies

### Files to Create
- `internal/workspace/load_workspace_test.go` - Workspace loading tests
- `internal/parse/mod_decoder_test.go` - Decoder tests
- `testdata/mods/` - Test mod fixtures

### Prerequisites
- Task 1 (Instrumentation) should be complete to enable timing in tests

## Implementation Details

### 1. Create Test Fixtures

Create test mods in `testdata/mods/`:

#### `testdata/mods/simple-mod/mod.pp`
```hcl
mod "simple_test" {
  title = "Simple Test Mod"
}

query "simple_query" {
  sql = "SELECT 1"
}

dashboard "simple_dashboard" {
  title = "Simple Dashboard"

  card {
    sql = "SELECT 1 as value"
  }
}
```

#### `testdata/mods/complex-mod/mod.pp`
```hcl
mod "complex_test" {
  title = "Complex Test Mod"
}

variable "region" {
  type    = string
  default = "us-east-1"
}

locals {
  common_tags = {
    test = "true"
  }
}

query "parameterized_query" {
  sql = "SELECT * FROM table WHERE region = $1"
  param "region" {
    default = var.region
  }
}

dashboard "complex_dashboard" {
  title = "Complex Dashboard"

  input "selection" {
    type = "select"
    sql  = "SELECT DISTINCT name FROM options"
  }

  container {
    card {
      sql = query.parameterized_query.sql
      args = [self.input.selection.value]
    }

    chart {
      type = "bar"
      sql  = "SELECT * FROM metrics"
    }
  }
}

control "test_control" {
  title = "Test Control"
  sql   = "SELECT 'ok' as status"
}

benchmark "test_benchmark" {
  title    = "Test Benchmark"
  children = [control.test_control]
}
```

#### `testdata/mods/with-dependencies/mod.pp`
```hcl
mod "with_deps" {
  title = "Mod With Dependencies"

  require {
    mod "github.com/turbot/steampipe-mod-aws-insights" {
      version = ">=0.1.0"
    }
  }
}
```

#### `testdata/mods/benchmark-only/mod.pp`
```hcl
mod "benchmark_only" {
  title = "Benchmark Only Mod"
}

control "ctrl_1" {
  sql = "SELECT 'pass'"
}

control "ctrl_2" {
  sql = "SELECT 'pass'"
}

benchmark "parent" {
  title = "Parent Benchmark"
  children = [
    benchmark.child_1,
    benchmark.child_2
  ]
}

benchmark "child_1" {
  title = "Child 1"
  children = [control.ctrl_1]
}

benchmark "child_2" {
  title = "Child 2"
  children = [control.ctrl_2]
}
```

### 2. Create Workspace Loading Tests

`internal/workspace/load_workspace_test.go`:

```go
package workspace_test

import (
    "context"
    "os"
    "path/filepath"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/turbot/powerpipe/internal/resources"
    "github.com/turbot/powerpipe/internal/workspace"
)

func TestLoadSimpleMod(t *testing.T) {
    ctx := context.Background()
    modPath := filepath.Join(testdataDir(), "mods", "simple-mod")

    w, ew := workspace.Load(ctx, modPath)
    require.NoError(t, ew.GetError())
    require.NotNil(t, w)

    // Verify mod loaded
    assert.Equal(t, "simple_test", w.Mod.ShortName)

    // Verify resources loaded
    res := w.GetPowerpipeModResources()
    assert.Len(t, res.Queries, 1)
    assert.Len(t, res.Dashboards, 1)
    assert.Contains(t, res.Queries, "simple_test.query.simple_query")
    assert.Contains(t, res.Dashboards, "simple_test.dashboard.simple_dashboard")
}

func TestLoadComplexMod(t *testing.T) {
    ctx := context.Background()
    modPath := filepath.Join(testdataDir(), "mods", "complex-mod")

    w, ew := workspace.Load(ctx, modPath)
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

    // Verify dashboard structure
    dashboard := res.Dashboards["complex_test.dashboard.complex_dashboard"]
    require.NotNil(t, dashboard)
    assert.NotEmpty(t, dashboard.GetChildren())
}

func TestLoadBenchmarkHierarchy(t *testing.T) {
    ctx := context.Background()
    modPath := filepath.Join(testdataDir(), "mods", "benchmark-only")

    w, ew := workspace.Load(ctx, modPath)
    require.NoError(t, ew.GetError())

    res := w.GetPowerpipeModResources()

    // Verify benchmark hierarchy
    parent := res.ControlBenchmarks["benchmark_only.benchmark.parent"]
    require.NotNil(t, parent)
    assert.Len(t, parent.GetChildren(), 2)

    // Verify children are correct type
    for _, child := range parent.GetChildren() {
        _, ok := child.(*resources.Benchmark)
        assert.True(t, ok, "child should be a Benchmark")
    }
}

func TestLoadModWithoutModfile(t *testing.T) {
    ctx := context.Background()
    modPath := t.TempDir()

    // Create a simple query file without a mod.pp
    queryFile := filepath.Join(modPath, "query.pp")
    os.WriteFile(queryFile, []byte(`query "test" { sql = "SELECT 1" }`), 0644)

    w, ew := workspace.Load(ctx, modPath)
    require.NoError(t, ew.GetError())

    // Should create default mod
    assert.Equal(t, "local", w.Mod.ShortName)
}

func TestLoadModWithVariables(t *testing.T) {
    ctx := context.Background()
    modPath := filepath.Join(testdataDir(), "mods", "complex-mod")

    w, ew := workspace.LoadWorkspacePromptingForVariables(ctx, modPath)
    require.NoError(t, ew.GetError())

    res := w.GetPowerpipeModResources()

    // Verify variable is loaded with default value
    v := res.Variables["complex_test.var.region"]
    require.NotNil(t, v)
    assert.Equal(t, "us-east-1", v.ValueGo)
}

func TestLoadModResourceCount(t *testing.T) {
    ctx := context.Background()
    modPath := filepath.Join(testdataDir(), "mods", "complex-mod")

    w, ew := workspace.Load(ctx, modPath)
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
}

func TestLoadModDashboardChildren(t *testing.T) {
    ctx := context.Background()
    modPath := filepath.Join(testdataDir(), "mods", "complex-mod")

    w, ew := workspace.Load(ctx, modPath)
    require.NoError(t, ew.GetError())

    res := w.GetPowerpipeModResources()
    dashboard := res.Dashboards["complex_test.dashboard.complex_dashboard"]

    // Verify dashboard has children
    children := dashboard.GetChildren()
    assert.NotEmpty(t, children)

    // Verify input was added
    assert.NotEmpty(t, res.DashboardInputs)
}

func TestLoadModIdempotent(t *testing.T) {
    ctx := context.Background()
    modPath := filepath.Join(testdataDir(), "mods", "simple-mod")

    // Load twice
    w1, _ := workspace.Load(ctx, modPath)
    w2, _ := workspace.Load(ctx, modPath)

    res1 := w1.GetPowerpipeModResources()
    res2 := w2.GetPowerpipeModResources()

    // Should produce identical results
    assert.Equal(t, len(res1.Queries), len(res2.Queries))
    assert.Equal(t, len(res1.Dashboards), len(res2.Dashboards))
}

func testdataDir() string {
    // Find testdata directory relative to this test file
    wd, _ := os.Getwd()
    return filepath.Join(wd, "..", "..", "testdata")
}
```

### 3. Create Resource Verification Tests

`internal/resources/mod_resources_test.go`:

```go
package resources_test

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/turbot/pipe-fittings/v2/modconfig"
    "github.com/turbot/powerpipe/internal/resources"
)

func TestGetResourceByParsedName(t *testing.T) {
    mod := modconfig.NewMod("test", "/tmp", hcl.Range{})
    res := resources.NewModResources(mod).(*resources.PowerpipeModResources)

    // Add a query
    query := &resources.Query{
        QueryProviderImpl: resources.QueryProviderImpl{
            FullName: "test.query.my_query",
        },
    }
    res.AddResource(query)

    // Retrieve by parsed name
    parsed, _ := modconfig.ParseResourceName("query.my_query")
    found, ok := res.GetResource(parsed)

    assert.True(t, ok)
    assert.Equal(t, query, found)
}

func TestAddResourceDuplicateDetection(t *testing.T) {
    mod := modconfig.NewMod("test", "/tmp", hcl.Range{})
    res := resources.NewModResources(mod).(*resources.PowerpipeModResources)

    query1 := &resources.Query{
        QueryProviderImpl: resources.QueryProviderImpl{
            FullName: "test.query.duplicate",
        },
    }
    query2 := &resources.Query{
        QueryProviderImpl: resources.QueryProviderImpl{
            FullName: "test.query.duplicate",
        },
    }

    diags1 := res.AddResource(query1)
    assert.Empty(t, diags1)

    diags2 := res.AddResource(query2)
    assert.NotEmpty(t, diags2, "should detect duplicate")
}

func TestWalkResources(t *testing.T) {
    mod := modconfig.NewMod("test", "/tmp", hcl.Range{})
    res := resources.NewModResources(mod).(*resources.PowerpipeModResources)

    // Add various resources
    res.AddResource(&resources.Query{...})
    res.AddResource(&resources.Control{...})
    res.AddResource(&resources.Dashboard{...})

    var count int
    res.WalkResources(func(item modconfig.HclResource) (bool, error) {
        count++
        return true, nil
    })

    assert.Equal(t, 4, count) // 3 resources + 1 mod
}
```

### 4. Create Parse Decoder Tests

`internal/parse/mod_decoder_test.go`:

```go
package parse_test

import (
    "testing"

    "github.com/hashicorp/hcl/v2"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/turbot/powerpipe/internal/parse"
)

func TestDecodeDashboard(t *testing.T) {
    hclContent := `
dashboard "test" {
    title = "Test Dashboard"

    card {
        sql = "SELECT 1"
    }
}
`
    decoder := parse.NewPowerpipeModDecoder()
    // Test decode logic...
}

func TestDecodeBenchmark(t *testing.T) {
    hclContent := `
control "c1" {
    sql = "SELECT 'pass'"
}

benchmark "test" {
    title = "Test Benchmark"
    children = [control.c1]
}
`
    // Test benchmark decoding...
}

func TestDecodeDetectionBenchmark(t *testing.T) {
    hclContent := `
detection "d1" {
    sql = "SELECT * FROM events"
}

benchmark "test" {
    type = "detection"
    title = "Detection Benchmark"
    children = [detection.d1]
}
`
    // Test detection benchmark decoding...
}
```

## Acceptance Criteria

- [ ] Test fixtures created for simple, complex, and benchmark mods
- [ ] `TestLoadSimpleMod` passes
- [ ] `TestLoadComplexMod` passes
- [ ] `TestLoadBenchmarkHierarchy` passes
- [ ] `TestLoadModWithoutModfile` passes
- [ ] `TestLoadModWithVariables` passes
- [ ] `TestLoadModResourceCount` passes
- [ ] `TestLoadModDashboardChildren` passes
- [ ] `TestLoadModIdempotent` passes
- [ ] Resource retrieval tests pass
- [ ] Parse decoder tests pass
- [ ] All tests run in < 30 seconds
- [ ] Tests are deterministic (no flakiness)

## Notes

- Use `t.Parallel()` where possible for faster test runs
- Clean up temp directories in tests
- Consider table-driven tests for edge cases
- Tests should not require external dependencies (database, network)
