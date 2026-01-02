# Task 10: Concurrent Access & Race Condition Tests

## Objective

Write stress tests that exercise concurrent access patterns to find race conditions, deadlocks, and data corruption issues.

## Context

- Lazy loading has shared state (cache, index)
- Multiple goroutines may access simultaneously
- Dashboard server handles concurrent sessions
- Race conditions are notoriously hard to reproduce
- Go's race detector is essential

## Dependencies

- Tasks 3-9 (all prior tests should pass)
- All source files with shared state

## Acceptance Criteria

- [ ] Add tests to `internal/workspace/concurrent_test.go`
- [ ] All tests pass with `-race` flag
- [ ] Test at least 5 different concurrent access patterns
- [ ] No deadlocks detected
- [ ] Tests run in <60 seconds total

## Test Cases to Implement

### Concurrent Resource Access
```go
// Test: Concurrent GetResource calls
func TestConcurrent_GetResource(t *testing.T)
// 100 goroutines each calling GetResource
// Mix of same and different resources
// No races, correct values returned

// Test: Concurrent cache access
func TestConcurrent_CacheAccess(t *testing.T)
// Heavy concurrent reads and writes
// LRU eviction happening
// No corruption

// Test: Concurrent index access
func TestConcurrent_IndexAccess(t *testing.T)
// Many goroutines querying index
// No races in index data structures
```

### Concurrent Workspace Operations
```go
// Test: Concurrent GetWorkspaceForExecution
func TestConcurrent_GetWorkspaceForExecution(t *testing.T)
// 50 goroutines calling simultaneously
// sync.Once should ensure single load
// All get same workspace

// Test: Browse during eager load
func TestConcurrent_BrowseDuringEagerLoad(t *testing.T)
// Some goroutines browsing (GetAvailableDashboards)
// Other goroutines trigger eager load
// No races or deadlocks

// Test: LoadBenchmark during LoadBenchmarkForExecution
func TestConcurrent_LoadBenchmarkVariants(t *testing.T)
// One goroutine calls LoadBenchmark
// Another calls LoadBenchmarkForExecution
// Same benchmark, different methods
```

### Concurrent Dependency Resolution
```go
// Test: Concurrent dependency resolution
func TestConcurrent_DependencyResolution(t *testing.T)
// Multiple benchmarks resolving dependencies
// Shared dependency graph
// Correct resolution for all

// Test: Concurrent circular detection
func TestConcurrent_CircularDetection(t *testing.T)
// Many concurrent HasCircularDependency calls
// Different starting points
// Correct results
```

### Concurrent Server Sessions
```go
// Test: Many concurrent sessions
func TestConcurrent_ManySessions(t *testing.T)
// 50 concurrent sessions
// Each doing different operations
// No cross-session contamination

// Test: Concurrent select_dashboard
func TestConcurrent_SelectDashboard(t *testing.T)
// 20 sessions select dashboards simultaneously
// Some same dashboard, some different
// All executions complete correctly

// Test: Session connect/disconnect storm
func TestConcurrent_SessionStorm(t *testing.T)
// Rapid connect/disconnect cycles
// Interleaved with operations
// No leaked resources
```

### Concurrent File Operations
```go
// Test: Concurrent load from disk
func TestConcurrent_LoadFromDisk(t *testing.T)
// Many goroutines loading different resources
// File I/O concurrent
// No corruption

// Test: Concurrent byte offset seeking
func TestConcurrent_ByteOffsetSeeking(t *testing.T)
// Multiple resources from same file
// Concurrent seeking to different offsets
// Correct data read
```

### Stress Patterns
```go
// Test: Read-heavy workload
func TestConcurrent_ReadHeavy(t *testing.T)
// 95% reads, 5% writes
// Common dashboard browsing pattern
// No degradation

// Test: Write-heavy workload
func TestConcurrent_WriteHeavy(t *testing.T)
// Heavy cache updates
// Many evictions
// Correct behavior

// Test: Bursty workload
func TestConcurrent_Bursty(t *testing.T)
// Periods of high activity
// Followed by idle
// Recovery between bursts
```

### Deadlock Detection
```go
// Test: No deadlock under contention
func TestConcurrent_NoDeadlock(t *testing.T)
// Lock ordering violations could cause deadlock
// Heavy concurrent access
// Completes within timeout (30s)

// Test: Lock hierarchy verification
func TestConcurrent_LockHierarchy(t *testing.T)
// Operations that acquire multiple locks
// Verify consistent ordering
```

### Memory Safety
```go
// Test: No data races in maps
func TestConcurrent_MapSafety(t *testing.T)
// Concurrent map access
// Must use sync.Map or mutex
// Race detector clean

// Test: No slice corruption
func TestConcurrent_SliceSafety(t *testing.T)
// Slices used in concurrent context
// No append races

// Test: Pointer safety
func TestConcurrent_PointerSafety(t *testing.T)
// Shared pointers
// No use-after-free patterns
```

### Event System Concurrency
```go
// Test: Concurrent event publishing
func TestConcurrent_EventPublishing(t *testing.T)
// Multiple executions publishing events
// Events routed correctly
// No lost events

// Test: Event handler during registration
func TestConcurrent_EventHandlerRegistration(t *testing.T)
// Publishing while handlers being added
// No missed handlers
```

### Recovery from Concurrent Issues
```go
// Test: Recovery after partial failure
func TestConcurrent_RecoveryAfterFailure(t *testing.T)
// One goroutine fails
// Others continue correctly
// System recovers

// Test: Timeout handling
func TestConcurrent_TimeoutHandling(t *testing.T)
// Context cancellation during concurrent ops
// Clean shutdown
// No leaked goroutines
```

## Implementation Guidelines

### Test Structure
```go
func TestConcurrent_Pattern(t *testing.T) {
    // Setup workspace
    ws := setupTestWorkspace(t, "fixture")

    // Create wait group for goroutines
    var wg sync.WaitGroup

    // Channel for errors
    errCh := make(chan error, 100)

    // Launch concurrent goroutines
    for i := 0; i < 50; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            // Do concurrent operation
            if err := doConcurrentOp(ws, id); err != nil {
                errCh <- err
            }
        }(i)
    }

    // Wait for completion
    wg.Wait()
    close(errCh)

    // Check for errors
    for err := range errCh {
        t.Errorf("concurrent error: %v", err)
    }
}
```

### Running with Race Detector
```bash
go test -race -timeout 120s ./internal/workspace/... -run Concurrent
```

### Goroutine Leak Detection
```go
func TestConcurrent_NoGoroutineLeaks(t *testing.T) {
    before := runtime.NumGoroutine()

    // Run concurrent operations
    runConcurrentTest(t)

    // Allow goroutines to settle
    time.Sleep(100 * time.Millisecond)

    after := runtime.NumGoroutine()
    if after > before + 5 { // Small tolerance
        t.Errorf("goroutine leak: before=%d, after=%d", before, after)
    }
}
```

## Notes

- Always run with `-race` flag
- Use reasonable timeouts to catch deadlocks
- Profile CPU/memory during stress tests
- Consider using t.Parallel() for test parallelism
- Monitor for goroutine leaks

## Output Files

- `internal/workspace/concurrent_test.go`
- `internal/resourcecache/concurrent_test.go`
- `internal/dashboardserver/concurrent_test.go`
