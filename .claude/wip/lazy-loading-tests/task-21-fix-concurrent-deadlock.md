# Task 21: Fix Concurrent Test Deadlock

## Status: COMPLETED

## Priority: CRITICAL

## Problem

Concurrent workspace tests are timing out with deadlock:

```
--- FAIL: TestConcurrent_NoDeadlock (60.01s)
panic: test timed out after 1m0s
```

The stack trace shows many goroutines stuck at:
- `sync.(*Once).doSlow` - waiting on mutex
- `(*LazyWorkspace).GetWorkspaceForExecution` - calling sync.Once.Do

This suggests a deadlock in the LazyWorkspace's sync.Once initialization.

## Root Cause Analysis

The deadlock pattern shows:
1. Multiple goroutines call `GetWorkspaceForExecution` concurrently
2. They all reach `sync.Once.Do()` at line 151 of `lazy_workspace.go`
3. The sync.Once's underlying mutex creates contention
4. If the initialization function (`initOnce`) itself tries to acquire additional locks, it could deadlock

Likely causes:
1. The initialization function in sync.Once is blocking on something
2. Re-entrant lock acquisition during initialization
3. The test setup doesn't properly initialize the workspace

## Files to Investigate

- `internal/workspace/concurrent_test.go` - especially `TestConcurrent_NoDeadlock`
- `internal/workspace/lazy_workspace.go:151` - the GetWorkspaceForExecution method
- Look for what happens inside the sync.Once.Do function

## Tasks

1. [ ] Review `GetWorkspaceForExecution` implementation
2. [ ] Identify what the sync.Once initialization does
3. [ ] Check if initialization requires locks that could deadlock
4. [ ] Fix test or implementation to prevent deadlock
5. [ ] Verify concurrent tests pass with -race flag

## Stack Trace Analysis

```
github.com/turbot/powerpipe/internal/workspace.(*LazyWorkspace).GetWorkspaceForExecution(0x140008b4000, {0x1061fb7f8?, 0x10748dbc0?})
    /Users/nathan/src/powerpipe/internal/workspace/lazy_workspace.go:151 +0x58
github.com/turbot/powerpipe/internal/workspace.TestConcurrent_NoDeadlock.func1(0x1a)
    /Users/nathan/src/powerpipe/internal/workspace/concurrent_test.go:577 +0x130
```

The test at line 577 spawns goroutines that call GetWorkspaceForExecution, which blocks at line 151.

## Notes

- This is a CRITICAL issue as it affects concurrent access to LazyWorkspace
- May indicate a real bug in the LazyWorkspace implementation, not just the test
- Need to understand the locking/sync strategy of LazyWorkspace

## Resolution

### Root Cause
The actual error was **NOT** a deadlock in `sync.Once`, but rather a **concurrent map writes** panic in `resolveChildrenRecursively`:

```
fatal error: concurrent map writes
github.com/turbot/pipe-fittings/v2/modconfig.(*ModTreeItemImpl).AddParent(...)
```

When multiple goroutines call `LoadBenchmarkForExecution` concurrently for the same benchmark:
1. They all load the benchmark from cache
2. They all call `resolveChildrenRecursively`
3. They all call `AddParent` on the same child resources
4. `AddParent` modifies a map without synchronization
5. This causes the concurrent map writes panic

### Fix
Added per-benchmark synchronization using `sync.Map` and `sync.Once`:

1. Added `resolvedBenchmarks sync.Map` field to `LazyWorkspace`
2. Added `benchmarkResolution` struct with a `sync.Once`
3. Modified `LoadBenchmarkForExecution` to use per-benchmark `sync.Once`

This ensures that child resolution for any given benchmark happens exactly once, even when multiple goroutines call `LoadBenchmarkForExecution` concurrently.

### Changes Made
- `internal/workspace/lazy_workspace.go`:
  - Added `resolvedBenchmarks sync.Map` field
  - Added `benchmarkResolution` struct
  - Modified `LoadBenchmarkForExecution` to use per-benchmark sync.Once

### Testing
All 14 concurrent tests pass with `-race` flag:
```
go test -v -race -timeout 120s ./internal/workspace/... -run "TestConcurrent_"
```
