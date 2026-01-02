# Task 1: Analyze Existing Test Gaps and Create Test Inventory

## Objective

Create a detailed inventory of existing test coverage and identify specific gaps that need to be filled. This task produces a comprehensive gap analysis document that guides all subsequent testing work.

## Context

- Powerpipe has existing tests across multiple packages
- Recent lazy loading changes added new code paths
- Some areas already well-tested, others have gaps
- Need to understand baseline before adding tests

## Dependencies

- None (this is a research task)

## Acceptance Criteria

- [ ] Complete inventory of all existing lazy loading related tests
- [ ] Identify specific untested code paths in:
  - `internal/workspace/lazy_workspace.go`
  - `internal/resourceindex/scanner.go`
  - `internal/resourceindex/index.go`
  - `internal/resourceloader/loader.go`
  - `internal/resourceloader/resolver.go`
  - `internal/dashboardserver/server.go` (lazy mode paths)
- [ ] Document coverage gaps with specific function/line references
- [ ] Prioritize gaps by risk (HIGH/MEDIUM/LOW)
- [ ] Create `/internal/testdata/test-gap-analysis.md` with findings

## Approach

1. **Catalog existing tests** by reading all `*_test.go` files in relevant packages
2. **Map test coverage** to source code functions
3. **Identify untested branches** especially:
   - Error handling paths
   - Edge case conditions
   - Concurrent access scenarios
4. **Document findings** in structured format

## Key Areas to Analyze

### Lazy Workspace (`lazy_workspace.go`)
- [ ] `NewLazyWorkspace()` - error paths
- [ ] `GetWorkspaceForExecution()` - race conditions, error caching
- [ ] `GetResource()` - cache miss paths
- [ ] `LoadBenchmarkForExecution()` vs `LoadBenchmark()` - subtle differences
- [ ] `preloadResources()` - background loading

### Scanner (`scanner.go`)
- [ ] Regex edge cases (special characters, unicode)
- [ ] Byte offset accuracy with various line endings
- [ ] Malformed HCL handling
- [ ] Very large files (>1MB)
- [ ] Deeply nested structures

### Index (`index.go`)
- [ ] Concurrent access patterns
- [ ] Hierarchy building edge cases
- [ ] Mod name collision handling

### Loader (`loader.go`)
- [ ] Cache eviction under pressure
- [ ] Partial load failures
- [ ] Dependency resolution failures

### Server (`server.go`)
- [ ] Lazy mode message handlers
- [ ] Session state management
- [ ] Execution cancellation

## Output

Create file: `/internal/testdata/test-gap-analysis.md` containing:
1. Test inventory table
2. Coverage heat map (by function)
3. Prioritized gap list
4. Recommended test additions

## Notes

- Focus on lazy loading specific code, not general workspace loading
- Note any existing tests that might be broken or incomplete
- Identify opportunities to refactor/consolidate tests
