# Task 22: Fix Dashboard Server Thread Safety Race Condition

**Status: COMPLETE**

## Objective

Fix the data race in `TestServer_ThreadSafety` where concurrent map access causes race conditions.

## Context

- Discovered when running `go test ./... -race`
- Test: `TestServer_ThreadSafety` in `internal/dashboardserver/server_integration_test.go`
- The test spawns 50 goroutines performing 100 operations each concurrently
- Race detector found 3 related race conditions

## Race Detection Output

```
==================
WARNING: DATA RACE
Read at 0x00c000ebe050 by goroutine 913:
  github.com/turbot/pipe-fittings/v2/workspace.(*Workspace).LoadLock()
      /Users/nathan/src/pipe-fittings/workspace/workspace.go:363 +0x40
  github.com/turbot/pipe-fittings/v2/workspace.(*Workspace).GetModResources()
      /Users/nathan/src/pipe-fittings/workspace/workspace.go:376 +0x3c
  github.com/turbot/powerpipe/internal/workspace.(*PowerpipeWorkspace).GetPowerpipeModResources()
      /Users/nathan/src/powerpipe/internal/workspace/powerpipe_workspace.go:140 +0x2c
  github.com/turbot/powerpipe/internal/dashboardserver.(*Server).buildAvailableDashboardsPayload()
      /Users/nathan/src/powerpipe/internal/dashboardserver/server.go:118 +0xa0
...

Previous write at 0x00c000ebe050 by goroutine 910:
  github.com/turbot/pipe-fittings/v2/workspace.(*Workspace).LoadLock()
      /Users/nathan/src/pipe-fittings/workspace/workspace.go:364 +0x74
...
==================
WARNING: DATA RACE
Write at 0x00c000dbb5f0 by goroutine 925:
  runtime.mapaccess2_faststr()
  github.com/turbot/powerpipe/internal/dashboardserver.(*Server).addDashboardClient()
      /Users/nathan/src/powerpipe/internal/dashboardserver/server.go:580 +0x6c

Previous read at 0x00c000dbb5f0 by goroutine 924:
  runtime.mapaccess1_faststr()
  github.com/turbot/powerpipe/internal/dashboardserver.(*Server).setDashboardInputsForSession()
      /Users/nathan/src/powerpipe/internal/dashboardserver/server.go:540 +0x268
==================
```

## Root Cause Analysis

### Issue 1: `getDashboardClients()` returns map reference without protection

The method `getDashboardClients()` (line 571-576) acquires a lock, returns the map, then releases the lock:

```go
func (s *Server) getDashboardClients() map[string]*DashboardClientInfo {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    return s.dashboardClients  // Map reference escapes the lock!
}
```

Then `setDashboardInputsForSession()` uses this map **outside** the lock:

```go
func (s *Server) setDashboardInputsForSession(sessionId string, inputs *dashboardexecute.InputValues) {
    dashboardClients := s.getDashboardClients()  // Lock acquired & released
    if sessionInfo, ok := dashboardClients[sessionId]; ok {  // Map access without lock!
        sessionInfo.DashboardInputs = inputs
    }
}
```

### Issue 2: `LoadLock()` race in pipe-fittings

This is a separate issue in the pipe-fittings library where `workspace.(*Workspace).LoadLock()` has a race condition on the `loaded` field. This is related to Task 16 (Schema Cache Race Condition).

## Priority

**HIGH** - This is causing test failures on every run with `-race` flag.

## Acceptance Criteria

- [ ] `TestServer_ThreadSafety` passes with `-race` flag
- [ ] No regression in existing dashboard server functionality
- [ ] All map access to `dashboardClients` is properly synchronized

## Proposed Solution

### Option 1: Make `setDashboardInputsForSession` hold the lock (Recommended)

```go
func (s *Server) setDashboardInputsForSession(sessionId string, inputs *dashboardexecute.InputValues) {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    if sessionInfo, ok := s.dashboardClients[sessionId]; ok {
        sessionInfo.DashboardInputs = inputs
    }
}
```

### Option 2: Return map copy from `getDashboardClients`

```go
func (s *Server) getDashboardClients() map[string]*DashboardClientInfo {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    // Return a copy
    result := make(map[string]*DashboardClientInfo, len(s.dashboardClients))
    for k, v := range s.dashboardClients {
        result[k] = v
    }
    return result
}
```

**Note**: Option 2 has performance implications and pointer aliasing issues.

## Files Modified

- `internal/dashboardserver/server.go` - Fixed `setDashboardInputsForSession` method
- `pipe-fittings/workspace/workspace.go` - Fixed `LoadLock` race condition

## Related Tasks

- Task 16: Schema Cache Race Condition - **RESOLVED** by fixing `LoadLock`
- Task 18: Event Handler Race Condition (separate Close() race)

## Implementation Steps

1. **Fix `setDashboardInputsForSession`**: Hold mutex during map access
2. **Verify**: Run `go test ./internal/dashboardserver/... -race`
3. **Full test**: Run `go test ./... -race` to ensure no regressions

## Notes

- The LC_DYSYMTAB warnings in the test output are macOS linker warnings related to Go's race detector - they are harmless and not code issues
- The `[ Error ]` log output is expected from error handling tests intentionally triggering errors
- The pipe-fittings `LoadLock` race is a separate issue that should be tracked in Task 16
