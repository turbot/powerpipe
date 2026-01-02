# Task 15: Fix Viper Global State Race Condition

## Status: COMPLETED

## Objective

Address the race condition in pipe-fittings caused by viper global state usage during workspace loading.

## Root Cause Analysis

The race occurred in `workspace/workspace.go:SetModfileExists()` where `viper.Set(constants.ArgModLocation, ...)` was called during workspace loading. When multiple workspaces were loaded concurrently (e.g., in tests), they all accessed viper's global map without synchronization.

Viper's global functions are not thread-safe - concurrent Get/Set operations can cause panics.

## Solution Implemented

**Approach: Remove the viper.Set() call (Option 3)**

After analysis, the `viper.Set(ArgModLocation)` call was found to be redundant because:
1. Code reads `ArgModLocation` from viper **before** workspace load to get the starting path
2. After workspace load, code uses `workspace.Path` directly
3. The workspace already updates `w.Path` with the resolved mod directory

The fix removes the viper global state update while keeping the workspace's internal state (`w.Path`) correctly updated.

## Changes Made

### pipe-fittings/workspace/workspace.go

```go
// BEFORE:
func (w *Workspace) SetModfileExists() {
    modFile, err := FindModFilePath(w.Path)
    modFileExists := !errors.Is(err, ErrorNoModDefinition)

    if modFileExists {
        w.modFilePath = modFile
        // also set it in the viper config, so that it is available to whoever is using it
        viper.Set(constants.ArgModLocation, filepath.Dir(modFile))
        w.Path = filepath.Dir(modFile)
        w.Mod.SetFilePath(modFile)
    }
}

// AFTER:
func (w *Workspace) SetModfileExists() {
    modFile, err := FindModFilePath(w.Path)
    modFileExists := !errors.Is(err, ErrorNoModDefinition)

    if modFileExists {
        w.modFilePath = modFile
        // Update the workspace path to the actual mod directory (not the original working directory).
        // Note: We intentionally do NOT update viper here, as consumers should use w.Path directly
        // after workspace load. The viper ArgModLocation value is only used to determine the starting
        // path for workspace loading, not the resolved path. This avoids race conditions when multiple
        // workspaces are loaded concurrently (e.g., in tests).
        w.Path = filepath.Dir(modFile)
        w.Mod.SetFilePath(modFile)
    }
}
```

Also removed the unused `viper` import from `workspace/workspace.go`.

## Test Results

- The **viper race condition is fixed**
- Tests pass without `-race` flag
- Individual tests pass with `-race` flag

## Additional Race Conditions Discovered

When running multiple tests with `t.Parallel()` and `-race`, two additional race conditions were discovered (NOT related to viper):

### 1. Schema Caching Race (Task 16)
- **Location**: `pipe-fittings/parse/schema.go:getResourceSchema()`
- **Cause**: Global schema map accessed without synchronization
- **Impact**: Tests cannot use `t.Parallel()` with `-race`
- **Tracked in**: Task 16

### 2. Event Handler Race (Task 18)
- **Location**: `powerpipe/internal/workspace/workspace_events.go`
- **Cause**: Race between `Close()` and `handleDashboardEvent()`
- **Impact**: Event handler tests cannot run concurrently
- **Tracked in**: Task 18

## Files Modified

- `pipe-fittings/workspace/workspace.go` - Removed viper.Set() call
- `powerpipe/internal/workspace/lazy_workspace_transition_test.go` - Updated comments

## Acceptance Criteria

- [x] Identify all viper global state usage in pipe-fittings workspace code
- [x] Implement thread-safe alternative (removed unnecessary viper.Set())
- [ ] Tests can run with `t.Parallel()` and `-race` flag without detecting races
  - Note: Blocked by unrelated schema caching and event handler races
- [x] No breaking changes to pipe-fittings API

## Priority

**Low** - This was an optimization. The fix is complete but test parallelism remains blocked by unrelated races.
