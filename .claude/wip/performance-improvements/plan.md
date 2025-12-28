# Project Plan: Powerpipe Mod Loading Performance Improvements

## Overview

Powerpipe can be slow to start when configured with many mods and dashboards. This project systematically improves mod loading performance through instrumentation, testing, and targeted optimizations.

**Goal**: Reduce startup time to dashboard list visibility by 50%+ for large mod setups.

## Task Breakdown

1. **Task 1**: Add Performance Instrumentation (Status: Pending)
2. **Task 2**: Create Mod Loading Test Suite (Status: Pending)
3. **Task 3**: Create Performance Benchmark Tests (Status: Pending)
4. **Task 4**: Baseline Performance Measurement (Status: Pending)
5. **Task 5**: Parallelize File I/O (Status: Pending)
6. **Task 6**: Parallelize HCL Parsing (Status: Pending)
7. **Task 7**: Optimize Database Client Creation (Status: Pending)
8. **Task 8**: Cache Available Dashboards Payload (Status: Pending)
9. **Task 9**: Final Performance Validation (Status: Pending)

## Execution Strategy

**Phase 1 - Foundation (Sequential)**: Tasks 1-4 must run sequentially
- Task 1: Instrumentation (required for measuring everything)
- Task 2: Regression tests (safety net before changes)
- Task 3: Performance benchmarks (measurement framework)
- Task 4: Baseline measurement (data collection)

**Phase 2 - Optimizations (Can be Parallel)**: Tasks 5-8
- Tasks 5 & 6: File I/O and HCL parsing (in pipe-fittings, can be parallel)
- Tasks 7 & 8: DB client and payload caching (in powerpipe, can be parallel)

**Phase 3 - Validation (Sequential)**: Task 9
- Final validation after all optimizations

## Dependencies

```
Task 1 (Instrumentation)
    ↓
Task 2 (Mod Loading Tests)
    ↓
Task 3 (Performance Benchmarks)
    ↓
Task 4 (Baseline Measurement)
    ↓
┌───────────────────────────────────────┐
│  Tasks 5, 6, 7, 8 (Optimizations)     │
│  Can run in parallel with worktrees   │
└───────────────────────────────────────┘
    ↓
Task 9 (Final Validation)
```

## Shared Resources

### Files Modified by Multiple Tasks
- `internal/workspace/load_workspace.go` - Tasks 1, 5
- `internal/initialisation/init_data.go` - Tasks 1, 7
- `internal/dashboardserver/payload.go` - Tasks 1, 8
- `internal/dashboardserver/server.go` - Tasks 1, 8

### Coordination Strategy
- Phase 1 tasks are sequential, no conflicts
- Phase 2 tasks modify different files, can use worktrees
- Each optimization task includes its own performance measurement

## Integration Plan

1. Each task creates its own feature branch
2. Tasks include unit tests for changes
3. Performance results documented in task completion notes
4. Final task (9) validates all improvements work together

## Performance Measurement Points

Key timing checkpoints to instrument:
1. `LoadWorkspacePromptingForVariables()` - Total workspace load
2. `LoadFileData()` - File I/O time
3. `ParseHclFiles()` - HCL parsing time
4. `Decoder.Decode()` - Block decoding time
5. `loadModDependenciesAsync()` - Dependency loading
6. `NewDbClient()` - Database connection time
7. `buildAvailableDashboardsPayload()` - Payload building time
8. Total time from `main()` to "Dashboard server started" message

## Success Criteria

- [ ] All existing tests pass
- [ ] No regression in mod loading correctness
- [ ] 50%+ improvement in startup time for large mod setups
- [ ] Performance improvements documented with before/after data
- [ ] New tests prevent future regressions

## Notes

- Changes to `pipe-fittings` library will require separate PR
- Focus first on Powerpipe-side optimizations
- Measure after each change to validate improvement
- Keep changes atomic and reversible
