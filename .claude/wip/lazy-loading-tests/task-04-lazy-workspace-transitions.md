# Task 4: Lazy Workspace Transition Tests

## Objective

Write comprehensive tests for the lazy-to-eager workspace transition, focusing on the hybrid loading approach and potential race conditions.

## Context

- `LazyWorkspace` uses `sync.Once` for eager loading
- First execution triggers `GetWorkspaceForExecution()`
- Event handlers must be copied between workspaces
- This is a critical path - bugs here affect all executions

## Dependencies

- Task 2 (test fixtures)
- Files to test: `internal/workspace/lazy_workspace.go`

## Acceptance Criteria

- [ ] Add tests to `internal/workspace/lazy_workspace_transition_test.go`
- [ ] Test all transition scenarios
- [ ] Verify thread-safety with concurrent tests
- [ ] Test error handling during transition
- [ ] No race conditions detected with `-race` flag

## Test Cases to Implement

### Basic Transition
```go
// Test: First call triggers eager load
func TestLazyWorkspace_FirstExecutionTriggersEagerLoad(t *testing.T)
// 1. Create lazy workspace
// 2. Verify IsLazy() == true
// 3. Call GetWorkspaceForExecution()
// 4. Verify eager workspace returned
// 5. Verify subsequent calls return same workspace

// Test: Lazy operations still work after eager load
func TestLazyWorkspace_LazyOpsAfterEagerLoad(t *testing.T)
// GetAvailableDashboardsFromIndex still works

// Test: Multiple calls return same workspace
func TestLazyWorkspace_EagerWorkspaceIsCached(t *testing.T)
// sync.Once ensures single load
```

### Concurrent Access
```go
// Test: Concurrent execution requests
func TestLazyWorkspace_ConcurrentExecutionRequests(t *testing.T)
// Launch 10 goroutines calling GetWorkspaceForExecution()
// All should get same workspace
// No panics or races

// Test: Mixed lazy/eager concurrent access
func TestLazyWorkspace_MixedConcurrentAccess(t *testing.T)
// Some goroutines browsing (lazy)
// Some goroutines executing (triggers eager)
// Should not interfere

// Test: Race between GetResource and GetWorkspaceForExecution
func TestLazyWorkspace_GetResourceDuringTransition(t *testing.T)
// One goroutine getting resource lazily
// Another triggering eager load
// No race condition
```

### Event Handler Transfer
```go
// Test: Event handlers copied to eager workspace
func TestLazyWorkspace_EventHandlersCopied(t *testing.T)
// 1. Register handler on lazy workspace
// 2. Trigger eager load
// 3. Publish event on eager workspace
// 4. Verify handler receives event

// Test: Multiple handlers all transferred
func TestLazyWorkspace_MultipleHandlersTransferred(t *testing.T)
// Register multiple handlers, verify all work

// Test: Handler registration after eager load
func TestLazyWorkspace_HandlerAfterEagerLoad(t *testing.T)
// Register handler after transition
// Should work on eager workspace
```

### Error Handling
```go
// Test: Eager load failure is cached
func TestLazyWorkspace_EagerLoadFailureCached(t *testing.T)
// Create workspace with invalid mod
// First GetWorkspaceForExecution() fails
// Second call returns same error (not re-tried)

// Test: Partial eager load failure
func TestLazyWorkspace_PartialEagerLoadFailure(t *testing.T)
// Some resources parse, some fail
// Verify error handling

// Test: Lazy operations work after eager failure
func TestLazyWorkspace_LazyOpsAfterEagerFailure(t *testing.T)
// Browsing should still work from index
```

### Resource Loading Differences
```go
// Test: LoadBenchmark vs LoadBenchmarkForExecution
func TestLazyWorkspace_LoadBenchmarkDifference(t *testing.T)
// LoadBenchmark: children loaded but not set
// LoadBenchmarkForExecution: children field populated
// Verify the difference

// Test: GetResource returns cached after eager
func TestLazyWorkspace_GetResourceAfterEager(t *testing.T)
// After eager load, GetResource should use eager workspace

// Test: Cache state after eager load
func TestLazyWorkspace_CacheStateAfterEager(t *testing.T)
// Verify cache relationship with eager workspace
```

### Workspace Interface Compliance
```go
// Test: LazyWorkspace satisfies DashboardServerWorkspace
func TestLazyWorkspace_InterfaceCompliance(t *testing.T)
// Compile-time check + runtime verification

// Test: Eager workspace from GetWorkspaceForExecution satisfies interface
func TestLazyWorkspace_EagerInterfaceCompliance(t *testing.T)
```

### State Consistency
```go
// Test: Resource data consistent between lazy and eager
func TestLazyWorkspace_ResourceDataConsistency(t *testing.T)
// Get resource via lazy path
// Get same resource via eager path
// Compare all fields

// Test: Benchmark hierarchy consistent
func TestLazyWorkspace_HierarchyConsistency(t *testing.T)
// Build hierarchy from index
// Build hierarchy from eager
// Should match

// Test: Available dashboards consistent
func TestLazyWorkspace_AvailableDashboardsConsistency(t *testing.T)
// From index vs from eager workspace
```

### Memory and Performance
```go
// Test: Memory usage during transition
func TestLazyWorkspace_MemoryDuringTransition(t *testing.T)
// Measure memory before/after eager load
// Should be bounded

// Test: Transition time is acceptable
func TestLazyWorkspace_TransitionTime(t *testing.T)
// For medium mod, transition < 5s
// For large mod, transition < 30s
```

## Implementation Notes

- Use `t.Parallel()` for concurrent tests
- Run with `-race` flag in CI
- Create helper to setup lazy workspace with test fixture
- Use mock workspace for error injection

## Output Files

- `internal/workspace/lazy_workspace_transition_test.go`
