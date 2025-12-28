# Task 8: Cache Available Dashboards Payload

## Objective

Cache the available dashboards payload to avoid rebuilding it on every WebSocket request, and pre-compute benchmark trunks during mod loading.

## Context

- `buildAvailableDashboardsPayload()` is called on every client connection
- The payload only changes when mod files change
- Benchmark trunk computation involves recursive tree traversal
- Pre-computing and caching can eliminate repeated work

## Dependencies

### Prerequisites
- Task 4 (Baseline Measurement) - Need baseline data for comparison

### Files to Modify
- `internal/dashboardserver/server.go`
- `internal/dashboardserver/payload.go`
- `internal/workspace/powerpipe_workspace.go` (for pre-computation)
- `internal/resources/benchmark.go` (for trunk caching)

## Implementation Details

### 1. Add Cached Payload to Server

```go
// internal/dashboardserver/server.go

type Server struct {
    mutex                   *sync.Mutex
    dashboardClients        map[string]*DashboardClientInfo
    webSocket               *melody.Melody
    workspace               *workspace.PowerpipeWorkspace
    defaultDatabase         connection.ConnectionStringProvider
    defaultSearchPathConfig backend.SearchPathConfig

    // Cached payloads
    payloadMutex             sync.RWMutex
    cachedAvailableDashboards []byte
    cachedServerMetadata      []byte
    cacheValid                bool
}

func NewServer(ctx context.Context, initData *initialisation.InitData, webSocket *melody.Melody) (*Server, error) {
    OutputWait(ctx, "Starting WorkspaceEvents Server")

    server := &Server{
        mutex:            &sync.Mutex{},
        dashboardClients: make(map[string]*DashboardClientInfo),
        webSocket:        webSocket,
        workspace:        initData.Workspace,
        defaultDatabase:  initData.DefaultDatabase,
        defaultSearchPathConfig: initData.DefaultSearchPathConfig,
    }

    // Pre-build cached payloads
    if err := server.buildCachedPayloads(ctx); err != nil {
        return nil, err
    }

    w := initData.Workspace
    w.RegisterDashboardEventHandler(ctx, server.HandleDashboardEvent)

    err := w.SetupWatcher(ctx, func(c context.Context, e error) {})
    OutputMessage(ctx, "WorkspaceEvents loaded")

    return server, err
}

func (s *Server) buildCachedPayloads(ctx context.Context) error {
    s.payloadMutex.Lock()
    defer s.payloadMutex.Unlock()

    // Build available dashboards payload
    workspaceResources := s.workspace.GetPowerpipeModResources()
    availablePayload, err := buildAvailableDashboardsPayload(workspaceResources)
    if err != nil {
        return err
    }
    s.cachedAvailableDashboards = availablePayload

    // Build server metadata payload
    metadataPayload, err := s.buildServerMetadataPayload(
        s.workspace.GetModResources(),
        &steampipeconfig.PipesMetadata{},
    )
    if err != nil {
        return err
    }
    s.cachedServerMetadata = metadataPayload

    s.cacheValid = true
    return nil
}

func (s *Server) invalidateCache() {
    s.payloadMutex.Lock()
    s.cacheValid = false
    s.payloadMutex.Unlock()
}

func (s *Server) getAvailableDashboardsPayload(ctx context.Context) ([]byte, error) {
    s.payloadMutex.RLock()
    if s.cacheValid {
        payload := s.cachedAvailableDashboards
        s.payloadMutex.RUnlock()
        return payload, nil
    }
    s.payloadMutex.RUnlock()

    // Rebuild cache
    if err := s.buildCachedPayloads(ctx); err != nil {
        return nil, err
    }

    s.payloadMutex.RLock()
    defer s.payloadMutex.RUnlock()
    return s.cachedAvailableDashboards, nil
}

func (s *Server) getServerMetadataPayload(ctx context.Context) ([]byte, error) {
    s.payloadMutex.RLock()
    if s.cacheValid {
        payload := s.cachedServerMetadata
        s.payloadMutex.RUnlock()
        return payload, nil
    }
    s.payloadMutex.RUnlock()

    // Rebuild cache
    if err := s.buildCachedPayloads(ctx); err != nil {
        return nil, err
    }

    s.payloadMutex.RLock()
    defer s.payloadMutex.RUnlock()
    return s.cachedServerMetadata, nil
}
```

### 2. Update Message Handler to Use Cache

```go
// internal/dashboardserver/server.go

func (s *Server) handleMessageFunc(ctx context.Context) func(session *melody.Session, msg []byte) {
    return func(session *melody.Session, msg []byte) {
        sessionId := s.getSessionId(session)

        var request ClientRequest
        if err := json.Unmarshal(msg, &request); err != nil {
            slog.Warn("failed to marshal message", "error", err.Error())
            return
        }

        if request.Action != "keep_alive" {
            slog.Debug("handleMessageFunc", "message", string(msg))
        }

        switch request.Action {
        case "get_server_metadata":
            payload, err := s.getServerMetadataPayload(ctx)
            if err != nil {
                OutputError(ctx, sperr.WrapWithMessage(err, "error getting server metadata"))
                return
            }
            _ = session.Write(payload)

        case "get_available_dashboards":
            payload, err := s.getAvailableDashboardsPayload(ctx)
            if err != nil {
                OutputError(ctx, sperr.WrapWithMessage(err, "error getting available dashboards"))
                return
            }
            _ = session.Write(payload)

        // ... rest of switch cases unchanged
        }
    }
}
```

### 3. Invalidate Cache on Dashboard Changes

```go
// internal/dashboardserver/server.go

func (s *Server) HandleDashboardEvent(ctx context.Context, event dashboardevents.DashboardEvent) {
    // ... existing code ...

    switch e := event.(type) {
    case *dashboardevents.DashboardChanged:
        slog.Debug("DashboardChanged event")

        // Invalidate cache on any dashboard change
        s.invalidateCache()

        // ... rest of existing handling
    }
}
```

### 4. Pre-compute Benchmark Trunks

```go
// internal/resources/benchmark.go

type Benchmark struct {
    // ... existing fields ...

    // Pre-computed trunk paths
    computedTrunks [][]string
    trunksComputed bool
}

// ComputeTrunks pre-computes trunk paths for this benchmark
func (b *Benchmark) ComputeTrunks() {
    if b.trunksComputed {
        return
    }

    b.computedTrunks = b.buildTrunks([]string{b.FullName})
    b.trunksComputed = true
}

func (b *Benchmark) buildTrunks(currentTrunk []string) [][]string {
    trunks := [][]string{currentTrunk}

    for _, child := range b.GetChildren() {
        if childBenchmark, ok := child.(*Benchmark); ok {
            childTrunk := append([]string{}, currentTrunk...)
            childTrunk = append(childTrunk, childBenchmark.FullName)
            trunks = append(trunks, childBenchmark.buildTrunks(childTrunk)...)
        }
    }

    return trunks
}

// GetTrunks returns pre-computed trunks or computes them on demand
func (b *Benchmark) GetTrunks() [][]string {
    if !b.trunksComputed {
        b.ComputeTrunks()
    }
    return b.computedTrunks
}
```

### 5. Trigger Pre-computation After Mod Load

```go
// internal/workspace/powerpipe_workspace.go

func (w *PowerpipeWorkspace) precomputeBenchmarkTrunks() {
    resources := w.GetPowerpipeModResources()

    // Pre-compute control benchmark trunks
    for _, benchmark := range resources.ControlBenchmarks {
        benchmark.ComputeTrunks()
    }

    // Pre-compute detection benchmark trunks
    for _, benchmark := range resources.DetectionBenchmarks {
        benchmark.ComputeTrunks()
    }
}
```

### 6. Use Pre-computed Trunks in Payload Builder

```go
// internal/dashboardserver/payload.go

func buildAvailableDashboardsPayload(workspaceResources *resources.PowerpipeModResources) ([]byte, error) {
    payload := AvailableDashboardsPayload{
        Action:     "available_dashboards",
        Dashboards: make(map[string]ModAvailableDashboard),
        Benchmarks: make(map[string]ModAvailableBenchmark),
        Snapshots:  workspaceResources.Snapshots,
    }

    if workspaceResources.Mod != nil {
        topLevelResources := resources.GetModResources(workspaceResources.Mod)

        // Dashboards (unchanged)
        for _, dashboard := range topLevelResources.Dashboards {
            // ... existing code
        }

        // Benchmarks - use pre-computed trunks
        for _, benchmark := range topLevelResources.ControlBenchmarks {
            if benchmark.IsAnonymous() {
                continue
            }

            isTopLevel := isTopLevelBenchmark(benchmark)

            availableBenchmark := ModAvailableBenchmark{
                Title:         benchmark.GetTitle(),
                FullName:      benchmark.FullName,
                ShortName:     benchmark.ShortName,
                BenchmarkType: "control",
                Tags:          benchmark.Tags,
                IsTopLevel:    isTopLevel,
                Children:      buildBenchmarkChildrenFromPrecomputed(benchmark),
                Trunks:        benchmark.GetTrunks(), // Use pre-computed
                ModFullName:   benchmark.Mod.GetFullName(),
            }

            payload.Benchmarks[benchmark.FullName] = availableBenchmark
        }

        // Similar for detection benchmarks...
    }

    return json.Marshal(payload)
}

func isTopLevelBenchmark(benchmark *resources.Benchmark) bool {
    for _, parent := range benchmark.GetParents() {
        if _, ok := parent.(*modconfig.Mod); ok {
            return true
        }
    }
    return false
}

func buildBenchmarkChildrenFromPrecomputed(benchmark *resources.Benchmark) []ModAvailableBenchmark {
    var children []ModAvailableBenchmark
    for _, child := range benchmark.GetChildren() {
        if childBenchmark, ok := child.(*resources.Benchmark); ok {
            children = append(children, ModAvailableBenchmark{
                Title:         childBenchmark.GetTitle(),
                FullName:      childBenchmark.FullName,
                ShortName:     childBenchmark.ShortName,
                BenchmarkType: "control",
                Tags:          childBenchmark.Tags,
                Trunks:        childBenchmark.GetTrunks(),
                Children:      buildBenchmarkChildrenFromPrecomputed(childBenchmark),
            })
        }
    }
    return children
}
```

### 7. Add Tests

```go
// internal/dashboardserver/payload_test.go

func TestCachedAvailableDashboardsPayload(t *testing.T) {
    // Setup server with workspace
    server := setupTestServer(t)

    // First call should build cache
    start := time.Now()
    payload1, err := server.getAvailableDashboardsPayload(context.Background())
    firstCallTime := time.Since(start)
    require.NoError(t, err)

    // Second call should use cache (much faster)
    start = time.Now()
    payload2, err := server.getAvailableDashboardsPayload(context.Background())
    secondCallTime := time.Since(start)
    require.NoError(t, err)

    // Payloads should be identical
    assert.Equal(t, payload1, payload2)

    // Second call should be significantly faster
    assert.Less(t, secondCallTime, firstCallTime/10,
        "cached call should be >10x faster")

    t.Logf("First call: %v, Cached call: %v", firstCallTime, secondCallTime)
}

func TestCacheInvalidationOnDashboardChange(t *testing.T) {
    server := setupTestServer(t)

    // Get initial payload
    payload1, _ := server.getAvailableDashboardsPayload(context.Background())

    // Invalidate cache
    server.invalidateCache()

    // Next call should rebuild
    payload2, _ := server.getAvailableDashboardsPayload(context.Background())

    // Payloads should still match (same workspace)
    assert.Equal(t, payload1, payload2)
}

func TestPrecomputedBenchmarkTrunks(t *testing.T) {
    benchmark := &resources.Benchmark{
        FullName: "mod.benchmark.parent",
        // Add child benchmarks...
    }

    // Compute trunks
    benchmark.ComputeTrunks()

    // Should be computed
    assert.True(t, benchmark.trunksComputed)
    assert.NotEmpty(t, benchmark.GetTrunks())

    // Second call should return cached
    trunks := benchmark.GetTrunks()
    assert.NotEmpty(t, trunks)
}
```

## Acceptance Criteria

- [ ] Available dashboards payload is cached on server startup
- [ ] Server metadata payload is cached on server startup
- [ ] Cache is invalidated on dashboard change events
- [ ] Cached payload is returned for subsequent requests
- [ ] Benchmark trunks are pre-computed during mod load
- [ ] Pre-computed trunks are used in payload building
- [ ] Unit tests verify caching behavior
- [ ] Unit tests verify cache invalidation
- [ ] Thread-safe cache access (RWMutex)
- [ ] Performance improvement measured and documented

## Expected Performance Improvement

| Operation | Baseline | After | Improvement |
|-----------|----------|-------|-------------|
| First `get_available_dashboards` | 50ms | 50ms | 0% (building) |
| Subsequent requests | 50ms | <1ms | >98% |
| With many clients | N*50ms | 50ms + N*<1ms | ~98% |

## Notes

- Cache should be thread-safe for concurrent WebSocket connections
- Consider cache TTL for very long-running servers (optional)
- Pre-computation adds slight overhead to initial load
- Total benefit depends on number of client connections
