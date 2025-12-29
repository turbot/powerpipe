# Task 1: Comprehensive Behavior Tests

## Objective

Create a comprehensive test suite that captures all current Powerpipe behavior before making architectural changes. These tests serve as a safety net to ensure lazy loading doesn't break existing functionality.

## Context

- This is the **first and most critical task** - we cannot safely refactor without tests
- Current test coverage may not cover all resource access patterns
- We need tests that verify behavior, not implementation details
- Tests should pass before AND after lazy loading implementation

## Dependencies

### Prerequisites
- None (this is a foundation task)

### Files to Create
- `internal/workspace/workspace_behavior_test.go` - Workspace loading behavior
- `internal/workspace/resource_access_test.go` - Resource access patterns
- `internal/dashboardserver/server_behavior_test.go` - Server API behavior
- `internal/dashboardexecute/execution_behavior_test.go` - Dashboard execution
- `testdata/mods/behavior_test_mod/` - Test mod with all resource types

### Files to Modify
- None (pure additions)

## Implementation Details

### 1. Create Comprehensive Test Mod

Create a test mod with ALL resource types and edge cases:

```hcl
# testdata/mods/behavior_test_mod/mod.pp
mod "behavior_test" {
  title = "Behavior Test Mod"
}

# Include all resource types:
# - Dashboards with all panel types
# - Benchmarks with nested children
# - Controls with various query patterns
# - Queries with args and params
# - Variables with defaults and overrides
# - Inputs (global and dashboard-scoped)
# - Categories, flows, graphs, hierarchies
# - Cross-references between resources
# - Resources that reference other resources
```

### 2. Workspace Loading Tests

```go
// internal/workspace/workspace_behavior_test.go

func TestWorkspaceLoading_AllResourceTypesAccessible(t *testing.T) {
    // Load workspace
    ws := loadTestWorkspace(t, "behavior_test_mod")

    // Verify all resource types are accessible
    resources := ws.GetModResources()

    // Dashboards
    assert.NotEmpty(t, resources.Dashboards)
    for name, dash := range resources.Dashboards {
        assert.NotEmpty(t, dash.FullName)
        assert.Equal(t, name, dash.FullName)
        // Verify all dashboard properties accessible
        _ = dash.GetTitle()
        _ = dash.GetDescription()
        _ = dash.GetTags()
        _ = dash.GetChildren()
    }

    // Controls
    assert.NotEmpty(t, resources.Controls)
    for _, ctrl := range resources.Controls {
        _ = ctrl.GetSQL()
        _ = ctrl.GetArgs()
        _ = ctrl.GetParams()
    }

    // Benchmarks
    assert.NotEmpty(t, resources.ControlBenchmarks)
    for _, bench := range resources.ControlBenchmarks {
        children := bench.GetChildren()
        // Verify children are resolved
        for _, child := range children {
            assert.NotNil(t, child)
        }
    }

    // ... repeat for all 20+ resource types
}

func TestWorkspaceLoading_ResourceReferencesResolved(t *testing.T) {
    ws := loadTestWorkspace(t, "behavior_test_mod")
    resources := ws.GetModResources()

    // Find a control that references a query
    ctrl := resources.Controls["behavior_test.control.uses_query"]
    assert.NotNil(t, ctrl)

    // Verify the query reference is resolved
    query := ctrl.GetQuery()
    assert.NotNil(t, query, "Query reference should be resolved")
    assert.NotEmpty(t, query.GetSQL())
}

func TestWorkspaceLoading_NestedResourcesAccessible(t *testing.T) {
    ws := loadTestWorkspace(t, "behavior_test_mod")
    resources := ws.GetModResources()

    // Dashboard with nested containers
    dash := resources.Dashboards["behavior_test.dashboard.nested"]
    children := dash.GetChildren()

    for _, child := range children {
        if container, ok := child.(*resources.DashboardContainer); ok {
            // Verify nested children accessible
            nestedChildren := container.GetChildren()
            assert.NotEmpty(t, nestedChildren)
        }
    }
}

func TestWorkspaceLoading_VariableResolution(t *testing.T) {
    ws := loadTestWorkspace(t, "behavior_test_mod")
    resources := ws.GetModResources()

    // Variables should have values
    for _, v := range resources.Variables {
        // Default should be set
        assert.False(t, v.Default.IsNull())
    }
}
```

### 3. Resource Access Pattern Tests

```go
// internal/workspace/resource_access_test.go

func TestResourceAccess_GetResourceByName(t *testing.T) {
    ws := loadTestWorkspace(t, "behavior_test_mod")
    resources := ws.GetModResources()

    testCases := []struct {
        resourceType string
        name         string
    }{
        {"dashboard", "behavior_test.dashboard.main"},
        {"query", "behavior_test.query.simple"},
        {"control", "behavior_test.control.basic"},
        {"benchmark", "behavior_test.benchmark.top"},
    }

    for _, tc := range testCases {
        t.Run(tc.resourceType+"/"+tc.name, func(t *testing.T) {
            parsed, _ := modconfig.ParseResourceName(tc.name)
            resource, found := resources.GetResource(parsed)
            assert.True(t, found, "Resource should be found: %s", tc.name)
            assert.NotNil(t, resource)
            assert.Equal(t, tc.name, resource.Name())
        })
    }
}

func TestResourceAccess_WalkResources(t *testing.T) {
    ws := loadTestWorkspace(t, "behavior_test_mod")
    resources := ws.GetModResources()

    var visited []string
    err := resources.WalkResources(func(r modconfig.HclResource) (bool, error) {
        visited = append(visited, r.Name())
        return true, nil
    })

    assert.NoError(t, err)
    assert.NotEmpty(t, visited)

    // Verify all expected resources were visited
    assert.Contains(t, visited, "behavior_test.dashboard.main")
    assert.Contains(t, visited, "behavior_test.query.simple")
}

func TestResourceAccess_QueryProviders(t *testing.T) {
    ws := loadTestWorkspace(t, "behavior_test_mod")
    resources := ws.GetModResources()

    providers := resources.QueryProviders()
    assert.NotEmpty(t, providers)

    for _, p := range providers {
        // All query providers should have SQL or Query reference
        sql := p.GetSQL()
        query := p.GetQuery()
        assert.True(t, sql != "" || query != nil,
            "QueryProvider %s should have SQL or Query", p.Name())
    }
}
```

### 4. Dashboard Server Behavior Tests

```go
// internal/dashboardserver/server_behavior_test.go

func TestServer_GetAvailableDashboards(t *testing.T) {
    server := setupTestServer(t, "behavior_test_mod")
    defer server.Shutdown()

    // Request available dashboards
    payload, err := server.getAvailableDashboardsPayload()
    assert.NoError(t, err)

    var response AvailableDashboardsPayload
    err = json.Unmarshal(payload, &response)
    assert.NoError(t, err)

    // Verify dashboards are present
    assert.NotEmpty(t, response.Dashboards)

    // Verify each dashboard has required fields
    for name, dash := range response.Dashboards {
        assert.NotEmpty(t, dash.FullName)
        assert.Equal(t, name, dash.FullName)
    }

    // Verify benchmarks are present
    assert.NotEmpty(t, response.Benchmarks)

    // Verify benchmark hierarchy (trunks)
    for _, bench := range response.Benchmarks {
        if bench.IsTopLevel {
            assert.NotEmpty(t, bench.Trunks)
        }
    }
}

func TestServer_DashboardMetadata(t *testing.T) {
    server := setupTestServer(t, "behavior_test_mod")
    defer server.Shutdown()

    // Get a dashboard
    dash := server.workspace.GetModResources().(*resources.PowerpipeModResources).
        Dashboards["behavior_test.dashboard.main"]

    payload, err := server.buildDashboardMetadataPayload(dash)
    assert.NoError(t, err)

    var response DashboardMetadataPayload
    err = json.Unmarshal(payload, &response)
    assert.NoError(t, err)

    assert.Equal(t, "dashboard_metadata", response.Action)
}
```

### 5. Dashboard Execution Tests

```go
// internal/dashboardexecute/execution_behavior_test.go

func TestExecution_DashboardRun(t *testing.T) {
    ws := loadTestWorkspace(t, "behavior_test_mod")
    executor := setupTestExecutor(t, ws)

    dash := ws.GetModResources().(*resources.PowerpipeModResources).
        Dashboards["behavior_test.dashboard.main"]

    ctx := context.Background()
    sessionId := "test-session"

    err := executor.ExecuteDashboard(ctx, sessionId, dash, nil, ws)
    assert.NoError(t, err)

    // Verify execution completed
    // (specific assertions depend on dashboard content)
}

func TestExecution_BenchmarkRun(t *testing.T) {
    ws := loadTestWorkspace(t, "behavior_test_mod")
    executor := setupTestExecutor(t, ws)

    bench := ws.GetModResources().(*resources.PowerpipeModResources).
        ControlBenchmarks["behavior_test.benchmark.top"]

    ctx := context.Background()
    sessionId := "test-session"

    err := executor.ExecuteDashboard(ctx, sessionId, bench, nil, ws)
    assert.NoError(t, err)
}

func TestExecution_NestedBenchmark(t *testing.T) {
    ws := loadTestWorkspace(t, "behavior_test_mod")
    executor := setupTestExecutor(t, ws)

    // Execute a benchmark with nested children
    bench := ws.GetModResources().(*resources.PowerpipeModResources).
        ControlBenchmarks["behavior_test.benchmark.nested"]

    ctx := context.Background()

    // Track all controls executed
    var executedControls []string
    // ... setup hooks to track execution

    err := executor.ExecuteDashboard(ctx, "test", bench, nil, ws)
    assert.NoError(t, err)

    // Verify all nested controls were executed
    assert.Contains(t, executedControls, "behavior_test.control.nested_1")
    assert.Contains(t, executedControls, "behavior_test.control.nested_2")
}
```

### 6. Edge Case Tests

```go
func TestEdgeCases_EmptyMod(t *testing.T) {
    ws := loadTestWorkspace(t, "empty_mod")
    resources := ws.GetModResources()

    assert.True(t, resources.Empty())
}

func TestEdgeCases_CircularReferences(t *testing.T) {
    // Should handle gracefully or error clearly
    ws := loadTestWorkspace(t, "circular_ref_mod")
    resources := ws.GetModResources()

    // Verify no infinite loops occurred
    assert.NotNil(t, resources)
}

func TestEdgeCases_MissingReferences(t *testing.T) {
    // Resources referencing non-existent resources
    ws, err := loadTestWorkspaceWithErrors(t, "missing_ref_mod")

    // Should either error or handle gracefully
    if err != nil {
        assert.Contains(t, err.Error(), "reference")
    }
}

func TestEdgeCases_LargeResourceCount(t *testing.T) {
    // Test with 1000+ resources
    ws := loadTestWorkspace(t, "large_mod")
    resources := ws.GetModResources()

    // All resources should be accessible
    count := 0
    resources.WalkResources(func(r modconfig.HclResource) (bool, error) {
        count++
        return true, nil
    })

    assert.Greater(t, count, 1000)
}
```

## Acceptance Criteria

- [ ] Test mod created with ALL resource types (20+)
- [ ] Test mod includes cross-resource references
- [ ] Test mod includes nested resources (benchmarks, containers)
- [ ] Workspace loading tests verify all resource types accessible
- [ ] Resource access tests verify GetResource, WalkResources, QueryProviders
- [ ] Server tests verify available_dashboards and dashboard_metadata payloads
- [ ] Execution tests verify dashboard and benchmark runs complete
- [ ] Edge case tests for empty, circular refs, missing refs, large mods
- [ ] All tests pass on current codebase (before any changes)
- [ ] Tests run in < 30 seconds total
- [ ] Tests are deterministic (no flaky tests)

## Notes

- These tests must pass BEFORE and AFTER lazy loading implementation
- Focus on behavior, not implementation details
- Don't test internal data structures that will change
- Use table-driven tests where appropriate
- Consider adding golden file tests for JSON payloads
