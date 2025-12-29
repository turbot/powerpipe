# Task 10: Dashboard Server Integration

## Objective

Integrate lazy loading into the dashboard server, ensuring the UI receives the same data while resources are loaded on-demand.

## Context

- Dashboard server provides WebSocket API to UI
- `available_dashboards` payload must work from index only
- Dashboard execution must load resources on-demand
- Server must handle concurrent dashboard executions efficiently

## Dependencies

### Prerequisites
- Task 9 (Workspace Integration) - Lazy workspace
- Task 1 (Behavior Tests) - Server behavior tests

### Files to Modify
- `internal/dashboardserver/server.go` - Use lazy workspace
- `internal/dashboardserver/available_dashboards.go` - Use index
- `internal/dashboardexecute/dashboard_tree_run_impl.go` - Lazy resource access

### Files to Create
- `internal/dashboardserver/server_lazy_test.go` - Integration tests

## Implementation Details

### 1. Server with Lazy Workspace

```go
// internal/dashboardserver/server.go modifications

// Server holds a lazy workspace when enabled
type Server struct {
    // ... existing fields

    // Lazy workspace (when lazy loading enabled)
    lazyWorkspace *workspace.LazyWorkspace

    // Flag for lazy mode
    lazyLoadingEnabled bool
}

// NewServer creates server with lazy loading option
func NewServer(ctx context.Context, opts ServerOptions) (*Server, error) {
    s := &Server{
        // ... existing init
    }

    if opts.LazyLoading {
        lw, err := workspace.NewLazyWorkspace(ctx, opts.WorkspacePath,
            workspace.DefaultLazyLoadConfig())
        if err != nil {
            return nil, err
        }
        s.lazyWorkspace = lw
        s.lazyLoadingEnabled = true
    } else {
        // Existing eager loading
        ws, err := workspace.LoadWorkspace(ctx, opts.WorkspacePath)
        if err != nil {
            return nil, err
        }
        s.workspace = ws
    }

    return s, nil
}

// getWorkspace returns the appropriate workspace
func (s *Server) getWorkspace() workspace.WorkspaceProvider {
    if s.lazyLoadingEnabled {
        return s.lazyWorkspace
    }
    return s.workspace
}
```

### 2. Available Dashboards from Index

```go
// internal/dashboardserver/available_dashboards.go modifications

func (s *Server) buildAvailableDashboardsPayload() ([]byte, error) {
    if s.lazyLoadingEnabled {
        // Use index directly - no parsing needed!
        payload := s.lazyWorkspace.GetAvailableDashboardsFromIndex()
        return json.Marshal(payload)
    }

    // Existing implementation for non-lazy mode
    return s.buildAvailableDashboardsPayloadEager()
}

// This is the key optimization - available_dashboards now comes from
// the index without loading any HCL resources!
```

### 3. Dashboard Execution with Lazy Loading

```go
// internal/dashboardexecute/lazy_execute.go
package dashboardexecute

import (
    "context"

    "github.com/turbot/powerpipe/internal/workspace"
)

// ExecuteDashboardLazy executes a dashboard using lazy loading
func ExecuteDashboardLazy(ctx context.Context, lw *workspace.LazyWorkspace,
    sessionId string, dashboardName string, inputs map[string]interface{}) error {

    // Load the dashboard (and its children) on demand
    dash, err := lw.LoadDashboard(ctx, dashboardName)
    if err != nil {
        return fmt.Errorf("loading dashboard %s: %w", dashboardName, err)
    }

    // Now execute with the loaded dashboard
    executor := NewDashboardExecutor(lw)
    return executor.Execute(ctx, sessionId, dash, inputs)
}

// DashboardExecutor executes dashboards with lazy resource access
type DashboardExecutor struct {
    lw *workspace.LazyWorkspace
}

func NewDashboardExecutor(lw *workspace.LazyWorkspace) *DashboardExecutor {
    return &DashboardExecutor{lw: lw}
}

func (e *DashboardExecutor) Execute(ctx context.Context, sessionId string,
    dash *modconfig.Dashboard, inputs map[string]interface{}) error {

    // Create execution tree
    tree := e.buildExecutionTree(ctx, dash)

    // Execute (existing logic)
    return tree.Execute(ctx, sessionId, inputs)
}

func (e *DashboardExecutor) buildExecutionTree(ctx context.Context,
    dash *modconfig.Dashboard) *DashboardTreeRun {

    // Children are already loaded by LoadDashboard
    // Build tree as usual
    return NewDashboardTreeRun(dash, e.lw, nil)
}
```

### 4. Benchmark Execution with Lazy Loading

```go
// internal/controlexecute/lazy_execute.go
package controlexecute

import (
    "context"

    "github.com/turbot/powerpipe/internal/workspace"
)

// ExecuteBenchmarkLazy executes a benchmark using lazy loading
func ExecuteBenchmarkLazy(ctx context.Context, lw *workspace.LazyWorkspace,
    benchmarkName string) (*ExecutionTree, error) {

    // Load benchmark with all children and control queries
    bench, err := lw.LoadBenchmark(ctx, benchmarkName)
    if err != nil {
        return nil, fmt.Errorf("loading benchmark %s: %w", benchmarkName, err)
    }

    // Execute with loaded resources
    return ExecuteBenchmark(ctx, lw, bench)
}
```

### 5. WebSocket Handler Updates

```go
// internal/dashboardserver/websocket_handler.go modifications

func (s *Server) handleDashboardSelect(ctx context.Context, sessionId string,
    msg DashboardSelectMessage) error {

    dashboardName := msg.DashboardName

    if s.lazyLoadingEnabled {
        // Load dashboard on-demand
        return ExecuteDashboardLazy(ctx, s.lazyWorkspace, sessionId,
            dashboardName, msg.Inputs)
    }

    // Existing eager execution
    return s.executeDashboardEager(ctx, sessionId, dashboardName, msg.Inputs)
}

func (s *Server) handleBenchmarkRun(ctx context.Context, sessionId string,
    msg BenchmarkRunMessage) error {

    benchmarkName := msg.BenchmarkName

    if s.lazyLoadingEnabled {
        tree, err := ExecuteBenchmarkLazy(ctx, s.lazyWorkspace, benchmarkName)
        if err != nil {
            return err
        }
        return s.sendBenchmarkResults(sessionId, tree)
    }

    // Existing eager execution
    return s.runBenchmarkEager(ctx, sessionId, benchmarkName)
}
```

### 6. Dashboard Metadata

```go
// internal/dashboardserver/dashboard_metadata.go modifications

func (s *Server) buildDashboardMetadataPayload(ctx context.Context,
    dashboardName string) ([]byte, error) {

    if s.lazyLoadingEnabled {
        // Load just this dashboard (not all dashboards)
        dash, err := s.lazyWorkspace.LoadDashboard(ctx, dashboardName)
        if err != nil {
            return nil, err
        }
        return buildMetadataPayload(dash)
    }

    // Existing implementation
    return s.buildMetadataPayloadEager(dashboardName)
}
```

### 7. Tests

```go
// internal/dashboardserver/server_lazy_test.go
package dashboardserver

import (
    "context"
    "encoding/json"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestServer_LazyAvailableDashboards(t *testing.T) {
    server := setupLazyServer(t)
    defer server.Shutdown()

    // Request available dashboards
    start := time.Now()
    payload, err := server.buildAvailableDashboardsPayload()
    elapsed := time.Since(start)
    require.NoError(t, err)

    var response AvailableDashboardsPayload
    err = json.Unmarshal(payload, &response)
    require.NoError(t, err)

    // Should be fast (from index)
    assert.Less(t, elapsed.Milliseconds(), int64(10))

    // Should have dashboards
    assert.NotEmpty(t, response.Dashboards)

    // No resources should be loaded yet
    assert.Equal(t, 0, server.lazyWorkspace.CacheStats().Entries)
}

func TestServer_LazyDashboardExecution(t *testing.T) {
    server := setupLazyServer(t)
    defer server.Shutdown()

    ctx := context.Background()
    sessionId := "test-session"

    // Execute dashboard
    err := server.handleDashboardSelect(ctx, sessionId, DashboardSelectMessage{
        DashboardName: "testmod.dashboard.main",
    })
    require.NoError(t, err)

    // Dashboard should now be cached
    stats := server.lazyWorkspace.CacheStats()
    assert.Greater(t, stats.Entries, 0)
}

func TestServer_LazyBenchmarkExecution(t *testing.T) {
    server := setupLazyServer(t)
    defer server.Shutdown()

    ctx := context.Background()
    sessionId := "test-session"

    // Run benchmark
    err := server.handleBenchmarkRun(ctx, sessionId, BenchmarkRunMessage{
        BenchmarkName: "testmod.benchmark.simple",
    })
    require.NoError(t, err)
}

func TestServer_ConcurrentDashboards(t *testing.T) {
    server := setupLazyServer(t)
    defer server.Shutdown()

    ctx := context.Background()

    // Execute multiple dashboards concurrently
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            sessionId := fmt.Sprintf("session-%d", id)
            dashName := fmt.Sprintf("testmod.dashboard.dash_%d", id%5)

            err := server.handleDashboardSelect(ctx, sessionId, DashboardSelectMessage{
                DashboardName: dashName,
            })
            assert.NoError(t, err)
        }(i)
    }

    wg.Wait()
}

func TestServer_PayloadEquivalence(t *testing.T) {
    // Ensure lazy and eager produce same payloads
    lazyServer := setupLazyServer(t)
    eagerServer := setupEagerServer(t)
    defer lazyServer.Shutdown()
    defer eagerServer.Shutdown()

    lazyPayload, err := lazyServer.buildAvailableDashboardsPayload()
    require.NoError(t, err)

    eagerPayload, err := eagerServer.buildAvailableDashboardsPayload()
    require.NoError(t, err)

    // Parse and compare
    var lazyResp, eagerResp AvailableDashboardsPayload
    json.Unmarshal(lazyPayload, &lazyResp)
    json.Unmarshal(eagerPayload, &eagerResp)

    assert.Equal(t, len(lazyResp.Dashboards), len(eagerResp.Dashboards))
    assert.Equal(t, len(lazyResp.Benchmarks), len(eagerResp.Benchmarks))

    // Check each dashboard matches
    for name, lazyDash := range lazyResp.Dashboards {
        eagerDash, ok := eagerResp.Dashboards[name]
        require.True(t, ok)
        assert.Equal(t, lazyDash.Title, eagerDash.Title)
        assert.Equal(t, lazyDash.FullName, eagerDash.FullName)
    }
}
```

## Acceptance Criteria

- [x] Server can start with lazy loading enabled
- [x] `available_dashboards` uses index (no resource loading)
- [x] Dashboard execution loads resources on-demand
- [x] Benchmark execution loads resources on-demand
- [x] Concurrent dashboard executions work correctly
- [x] Payloads match between lazy and eager modes
- [ ] Server handles cache eviction during long sessions (needs LazyWorkspace file watcher)
- [x] WebSocket handlers updated for lazy loading
- [x] All server behavior tests from Task 1 pass
- [x] Memory usage bounded during server operation

## Implementation Notes

### Completed Changes

1. **Extended `WorkspaceProvider` interface** to `DashboardServerWorkspace`:
   - Added `RegisterDashboardEventHandler`, `SetupWatcher`, `GetModResources`, `PublishDashboardEvent`
   - Both `PowerpipeWorkspace` and `LazyWorkspace` implement this interface

2. **Updated `Server` struct**:
   - Added `lazyWorkspace` field for lazy loading mode
   - Added `isLazyMode()` helper method
   - Added `getActiveWorkspace()` to return the appropriate workspace
   - Added `getWorkspaceForExecution()` for dashboard execution
   - Added `NewServerWithLazyWorkspace()` constructor for lazy mode

3. **Updated payload building**:
   - Added `buildAvailableDashboardsPayloadFromIndex()` for lazy mode
   - Added `convertIndexBenchmarkInfo()` for benchmark conversion
   - Server method `buildAvailableDashboardsPayload()` delegates to appropriate builder

4. **Updated WebSocket handlers**:
   - All handlers now use `getActiveWorkspace()` for workspace access
   - Dashboard execution uses `getWorkspaceForExecution()` for the executor

5. **Tests**:
   - Added `server_lazy_test.go` with 8 tests covering lazy mode functionality
   - All existing tests continue to pass

## Notes

- Consider adding lazy loading toggle as server flag
- May need to handle session cleanup (release cached resources)
- Watch for WebSocket message ordering during load
- Monitor cache hit rates in production
