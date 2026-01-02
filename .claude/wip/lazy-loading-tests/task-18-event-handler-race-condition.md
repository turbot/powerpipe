# Task 18: Fix Event Handler Race Condition in Powerpipe (Optional)

## Objective

Address the race condition in powerpipe's workspace event handling between `Close()` and `handleDashboardEvent()`.

## Context

- Discovered during Task 15 (Viper Race Condition Fix) investigation
- The race occurs in `internal/workspace/workspace_events.go`
- Race between `Close()` writing/closing and `handleDashboardEvent()` goroutine reading
- This affects event handler tests when running with `t.Parallel()` and `-race`

## Race Detection Output

```
WARNING: DATA RACE
Write at 0x00c000aeed60 by goroutine 733:
  github.com/turbot/powerpipe/internal/workspace.(*PowerpipeWorkspace).Close()
      /Users/nathan/src/powerpipe/internal/workspace/powerpipe_workspace.go:47 +0xc4
  github.com/turbot/powerpipe/internal/workspace.(*LazyWorkspace).Close()
      /Users/nathan/src/powerpipe/internal/workspace/lazy_workspace.go:533 +0x38

Previous read at 0x00c000aeed60 by goroutine 739:
  github.com/turbot/powerpipe/internal/workspace.(*PowerpipeWorkspace).handleDashboardEvent()
      /Users/nathan/src/powerpipe/internal/workspace/workspace_events.go:58 +0x40
  github.com/turbot/powerpipe/internal/workspace.(*PowerpipeWorkspace).RegisterDashboardEventHandler.gowrap1()
      /Users/nathan/src/powerpipe/internal/workspace/workspace_events.go:43 +0x48
```

## Dependencies

- This is a powerpipe issue (not pipe-fittings)
- Independent of Task 15 (viper race - fixed)
- Independent of Task 16 (schema cache race)

## Acceptance Criteria

- [ ] Event handler tests pass with `t.Parallel()` and `-race` flag
- [ ] No breaking changes to workspace event API
- [ ] Graceful shutdown of event handlers during Close()

## Files to Investigate

- `internal/workspace/workspace_events.go` - Event handler registration and dispatch
- `internal/workspace/powerpipe_workspace.go` - Close() method

## Proposed Solutions

### Option 1: Use sync.Once for Close + Context Cancellation
```go
type PowerpipeWorkspace struct {
    // ...
    closeOnce sync.Once
    closeCh   chan struct{}
}

func (w *PowerpipeWorkspace) Close() {
    w.closeOnce.Do(func() {
        close(w.closeCh)
        // Wait for handlers to finish or timeout
    })
}

func (w *PowerpipeWorkspace) handleDashboardEvent(ctx context.Context, event DashboardEvent) {
    select {
    case <-w.closeCh:
        return // Workspace closing, don't process
    default:
        // Process event
    }
}
```

### Option 2: RWMutex to Protect Event Channel Access
```go
type PowerpipeWorkspace struct {
    // ...
    eventMu sync.RWMutex
    closed  bool
}

func (w *PowerpipeWorkspace) Close() {
    w.eventMu.Lock()
    w.closed = true
    w.eventMu.Unlock()
    // ...
}

func (w *PowerpipeWorkspace) handleDashboardEvent(...) {
    w.eventMu.RLock()
    if w.closed {
        w.eventMu.RUnlock()
        return
    }
    w.eventMu.RUnlock()
    // Process event
}
```

### Option 3: Atomic Flag
```go
type PowerpipeWorkspace struct {
    // ...
    closed atomic.Bool
}

func (w *PowerpipeWorkspace) Close() {
    if w.closed.Swap(true) {
        return // Already closed
    }
    // ...
}

func (w *PowerpipeWorkspace) handleDashboardEvent(...) {
    if w.closed.Load() {
        return
    }
    // Process event
}
```

## Recommended Approach

**Option 1 (sync.Once + Context)** is recommended because:
1. Idiomatic Go pattern for graceful shutdown
2. Works well with context cancellation
3. Ensures Close() is only executed once
4. Provides clean signal to handlers to stop

## Implementation Steps

1. **Analyze**: Review current event handling implementation
2. **Design**: Choose synchronization approach
3. **Implement**: Add close coordination
4. **Test**: Enable `t.Parallel()` in event handler tests
5. **Verify**: Run with `-race` flag

## Affected Tests

These tests in `lazy_workspace_transition_test.go` cannot use `t.Parallel()`:
- `TestLazyWorkspace_EventHandlersCopied`
- `TestLazyWorkspace_MultipleHandlersTransferred`
- `TestLazyWorkspace_HandlerAfterEagerLoad`
- `TestLazyWorkspace_PublishEventBeforeAndAfterTransition`

## Priority

**Low** - This is an optimization for test parallelism. The race only manifests in tests, not in production (single workspace per process with proper lifecycle management).

## Related

- Task 15: Viper Race Condition Fix (completed - discovered this issue)
- Task 16: Schema Cache Race Condition (separate pipe-fittings issue)

## Notes

- This race only affects tests that register event handlers and close workspaces concurrently
- In production, workspace lifecycle is managed by a single goroutine
- Fix would enable full test parallelization with `-race` flag (along with Task 16)
