# Project Status: Lazy Loading Test Suite

## Current Phase: Test Failures Investigation

| Phase | Tasks | Status |
|-------|-------|--------|
| Phase 1: Foundation | Tasks 1-2 | **Complete** |
| Phase 2: Core Unit Tests | Tasks 3-6 | In Progress (Task 3 done) |
| Phase 3: Integration Tests | Tasks 7-9, **13** | Pending |
| Phase 4: Stress & Edge Cases | Tasks 10-11 | Pending |
| Phase 5: CLI Integration | Task 12 | Pending |
| Phase 6: Optional | Tasks 14-17 | In Progress (Task 14 done) |
| **Phase 7: Test Fixes** | **Tasks 19-22** | **Pending - BLOCKING** |

## Task Status

| Task | Description | Status | Notes |
|------|-------------|--------|-------|
| 1 | Analyze test gaps | Pending | Research task |
| 2 | Design test fixtures | **Complete** | Fixtures in testdata/mods/lazy-loading-tests/ |
| 3 | Scanner edge cases | **Complete** | 110 tests, documented 4 limitations |
| 4 | Lazy workspace transitions | **Complete** | 26 tests in lazy_workspace_transition_test.go |
| 5 | Dependency resolution | Pending | Parallel with 3,4,6 |
| 6 | Cache behavior | Pending | Parallel with 3-5 |
| 7 | WebSocket integration | Pending | Parallel with 8,9,13 |
| 8 | Cross-mod resource refs | Pending | Parallel with 7,9,13 |
| 9 | Benchmark hierarchy | Pending | Parallel with 7,8,13 |
| 10 | Concurrent access | **Complete** | 40+ tests across 3 packages, found pipe-fittings race |
| 11 | Error handling | Pending | Parallel with 10 |
| 12 | CLI integration | Pending | Final task |
| **13** | **Mod dependencies (CRITICAL)** | **Pending** | **Parallel with 7,8,9** |
| 14 | Scanner limitation fixes | **Complete** | Fixed 5 edge cases: escaped quotes, single-line blocks/tags, block comments, heredocs |
| 15 | Viper race condition fix | Pending | LOW priority, optional, requires pipe-fittings changes |
| **16** | **Schema cache race condition** | **Complete** | Fixed in Task 22 - `LoadLock` race in pipe-fittings |
| 17 | HCL syntax-based scanner | Pending | OPTIONAL, alternative to regex scanner using hclsyntax.ParseConfig |
| **18** | **Event handler race condition** | **Pending** | Requires changes to workspace_events.go |
| **19** | **Fix cache tests** | **Pending** | **HIGH priority** - Tests calling non-existent `newResourceCache()` |
| **20** | **Fix resourceloader panics** | **Pending** | **HIGH priority** - Nil pointer dereference in error_test.go |
| **21** | **Fix concurrent deadlock** | **Pending** | **CRITICAL** - `TestConcurrent_NoDeadlock` times out, sync.Once deadlock |
| **22** | **Fix dashboardserver thread safety race** | **Complete** | Fixed map race in powerpipe + LoadLock race in pipe-fittings |

## Execution Plan

### Recommended Execution Order

1. **Sequential**: Task 1 → Task 2 (foundation work)
2. **Parallel**: Tasks 3, 4, 5, 6 (independent unit tests)
3. **Parallel**: Tasks 7, 8, 9 (independent integration tests)
4. **Parallel**: Tasks 10, 11 (stress/error testing)
5. **Sequential**: Task 12 (CLI tests after all others)

### Time Estimates

- Phase 1: ~2-3 hours
- Phase 2: ~4-6 hours (parallel execution)
- Phase 3: ~3-4 hours (parallel execution)
- Phase 4: ~2-3 hours (parallel execution)
- Phase 5: ~2-3 hours

**Total**: ~13-19 hours of work

## Key Risk Areas Identified

1. **Hybrid Mode Transition** - CRITICAL
2. **Index vs Full Parse Mismatch** - CRITICAL
3. **Mod Dependency Resolution** - CRITICAL (silent failures, version conflicts)
4. **Benchmark Hierarchy/Trunks** - HIGH (already had bugs)
5. **Cross-Mod Resource References** - HIGH
6. **Cache Behavior** - MEDIUM-HIGH
7. **WebSocket State Management** - MEDIUM-HIGH
8. **Dependency Resolution** - MEDIUM
9. **Error Propagation** - MEDIUM

## Files to Create

### Test Files
- `internal/resourceindex/scanner_edge_test.go`
- `internal/resourceindex/hierarchy_test.go`
- `internal/resourceindex/mod_discovery_test.go`
- `internal/workspace/lazy_workspace_transition_test.go`
- `internal/workspace/cross_mod_test.go`
- `internal/workspace/mod_dependency_test.go`
- `internal/workspace/concurrent_test.go`
- `internal/workspace/error_handling_test.go`
- `internal/resourceloader/resolver_edge_test.go`
- `internal/resourceloader/transitive_deps_test.go`
- `internal/resourcecache/cache_edge_test.go`
- `internal/dashboardserver/server_integration_test.go`
- `internal/cmd/integration_test.go`

### Test Fixtures
- `internal/testdata/mods/lazy-loading-tests/` (full structure)
- `internal/testdata/test-gap-analysis.md`

## Notes

- Tests should pass with `-race` flag when run individually
- **Known limitation**: Concurrent GetResource calls trigger schema cache race in pipe-fittings (see Task 16)
- **Known limitation**: Event handler tests trigger race on Close() (see Task 18)
- Focus on bug hunting, not just coverage
- Each test should have clear purpose documentation
- Use table-driven tests where appropriate

## Race Conditions Status

| Race Condition | Location | Task | Status |
|----------------|----------|------|--------|
| Viper Global State | pipe-fittings `workspace.SetModfileExists()` | Task 15 | **FIXED** |
| Schema Cache / LoadLock | pipe-fittings `workspace.LoadLock()` | Task 16/22 | **FIXED** |
| Dashboard Clients Map | powerpipe `dashboardserver.setDashboardInputsForSession()` | Task 22 | **FIXED** |
| Event Handler | powerpipe `workspace_events.go` | Task 18 | Pending |

Task 18 must be completed for full test parallelization with `-race` flag.

## Test Failure Summary (2026-01-02)

Running `go test ./... -v` revealed the following failures:

### 1. Resource Cache Tests (Task 19) - HIGH Priority
**Package**: `internal/resourcecache`
**Error**: `panic: interface conversion: interface {} is nil, not *resourcecache.resourceCache`
**Tests affected**: 7 tests (`TestCache_BasicSetGet`, `TestCache_EdgeCase_*`, `TestConcurrent_*`)
**Root cause**: Tests call `newResourceCache()` which doesn't exist or returns wrong type

### 2. Resource Loader Tests (Task 20) - HIGH Priority
**Package**: `internal/resourceloader`
**Error**: `panic: runtime error: invalid memory address or nil pointer dereference`
**Tests affected**: 5+ tests (`TestErrorHandling_*`)
**Root cause**: Nil pointer dereference at error_test.go:33

### 3. Concurrent Workspace Tests (Task 21) - CRITICAL
**Package**: `internal/workspace`
**Error**: `panic: test timed out after 1m0s` (deadlock)
**Tests affected**: `TestConcurrent_NoDeadlock`, `TestConcurrent_RaceConditions`
**Root cause**: sync.Once deadlock in `LazyWorkspace.GetWorkspaceForExecution` (lazy_workspace.go:151)

### 4. Dashboard Server Thread Safety (Task 22) - ~~HIGH~~ **FIXED**
**Package**: `internal/dashboardserver`
**Error**: `testing.go:1617: race detected during execution of test`
**Tests affected**: `TestServer_ThreadSafety`
**Root cause**:
- `getDashboardClients()` returns map reference that escapes the lock
- `setDashboardInputsForSession()` accesses map without holding mutex
- Also triggers pipe-fittings `LoadLock()` race (related to Task 16)
**Fix**:
- Fixed `setDashboardInputsForSession` to hold mutex in powerpipe
- Changed `loadLock` from `*sync.Mutex` to `sync.Mutex` in pipe-fittings

### Recommended Fix Order
1. **Task 21** (CRITICAL) - Fix deadlock first as it may indicate real implementation bug
2. ~~**Task 22** (HIGH) - Fix dashboardserver map race - simple fix~~ **DONE**
3. **Task 19** (HIGH) - Fix cache tests
4. **Task 20** (HIGH) - Fix resourceloader tests
