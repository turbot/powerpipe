# Project Plan: Comprehensive Lazy Loading Test Suite

## Overview

Build a robust, bug-hunting test suite for Powerpipe's lazy loading feature. The goal is to achieve high confidence before making lazy loading the default approach. This project focuses on finding bugs and edge cases, not just coverage metrics.

## Key Risk Areas Identified

Based on deep code analysis, these are the highest-risk areas where bugs are likely to lurk:

### 1. Hybrid Mode Transition (CRITICAL)
- The switch from lazy → eager loading when execution starts
- `GetWorkspaceForExecution()` uses `sync.Once` - what if it fails?
- Event handler copying between workspaces
- Race conditions between browsing and execution

### 2. Index vs Full Parse Mismatch (CRITICAL)
- Fast regex scanner may miss things that full HCL parser catches
- Byte offset accuracy for seeking back to resources
- Edge cases in title/description/tag extraction
- Malformed HCL that regex accepts but parser rejects

### 3. Benchmark Hierarchy & Trunks (HIGH)
- Already had bugs (mod_full_name, trunks missing) - proven fragile
- Recursive trunk building for deep hierarchies
- Parent-child relationship setting post-scan
- Top-level detection logic

### 4. Cross-Mod Dependencies (HIGH)
- Resources in one mod referencing another
- Mod name mapping and resolution
- Version-stripped paths causing collisions
- Dependency mod scanning order

### 4a. Mod Dependency Resolution (CRITICAL - Added)
- **Silent failure on missing deps** (resolver.go:72-79 skips without warning)
- **Transitive dependency chains** (A → B → C resolution)
- **Version conflict handling** (last-wins replacement, no merge)
- **Require block parsing** (version vs branch vs path vs tag)
- **Diamond dependencies** (same mod via multiple paths)
- **Circular mod dependencies** (ModA ↔ ModB)

### 5. Cache Behavior (MEDIUM-HIGH)
- LRU eviction under memory pressure
- Cache invalidation when files change
- Concurrent cache access patterns
- Cache corruption/partial state scenarios

### 6. Websocket Server State (MEDIUM-HIGH)
- Session management with execution state
- Input change handling mid-execution
- Execution cancellation and cleanup
- Multiple concurrent dashboard selections

### 7. Dependency Resolution (MEDIUM)
- Circular dependency detection
- Missing dependency handling
- Transitive dependency loading order
- Optional vs required dependencies

### 8. Error Propagation (MEDIUM)
- Partial load failures
- Index corruption recovery
- Eager load failure caching (sync.Once traps errors)
- Graceful degradation

## Task Breakdown

### Phase 1: Foundation & Research (Tasks 1-2)
1. **Task 1**: Analyze existing test gaps and create test inventory
2. **Task 2**: Design comprehensive test mod fixtures

### Phase 2: Core Unit Tests (Tasks 3-6)
3. **Task 3**: Index/Scanner edge case tests
4. **Task 4**: Lazy workspace transition tests
5. **Task 5**: Dependency resolution edge case tests
6. **Task 6**: Cache behavior tests

### Phase 3: Integration Tests (Tasks 7-9, 13)
7. **Task 7**: Websocket server integration tests
8. **Task 8**: Cross-mod resource reference tests
9. **Task 9**: Benchmark hierarchy tests
13. **Task 13**: Mod dependency edge cases (CRITICAL) - require blocks, version conflicts, transitive deps

### Phase 4: Stress & Edge Cases (Tasks 10-11)
10. **Task 10**: Concurrent access & race condition tests
11. **Task 11**: Error handling & recovery tests

### Phase 5: CLI Integration (Task 12)
12. **Task 12**: CLI command integration tests

### Phase 6: Optional Improvements (Tasks 14-17)
14. **Task 14**: Scanner limitation fixes (LOW priority) - fix edge cases found in Task 3 ✅ COMPLETE
15. **Task 15**: Viper race condition fix (LOW priority) - requires pipe-fittings changes
16. **Task 16**: Schema cache race condition (MEDIUM priority) - requires pipe-fittings changes
17. **Task 17**: HCL syntax-based scanner (OPTIONAL) - prototype alternative scanner using hclsyntax.ParseConfig

## Execution Strategy

**Sequential execution recommended** for phases due to dependencies:
- Phase 1 informs test fixture design used in later phases
- Phase 2 unit tests should pass before integration tests
- Phase 3 builds on Phase 2 foundations

**Within phases**, some tasks can run in parallel:
- Tasks 3-6 are independent (parallel OK)
- Tasks 7-9, 13 are independent (parallel OK)
- Tasks 10-11 are independent (parallel OK)

## Test Approach Philosophy

1. **Bug hunting over coverage**: Design tests to find bugs, not just hit lines
2. **Edge cases first**: Focus on boundaries, limits, and unusual inputs
3. **Realistic scenarios**: Model actual user workflows
4. **Regression prevention**: Tests should catch the bugs we've already fixed
5. **Deterministic & fast**: Tests must be reliable and quick to run

## Shared Test Fixtures

Create a comprehensive test mod structure in `/internal/testdata/mods/lazy-loading-tests/`:

```
lazy-loading-tests/
├── mod.pp                           # Main mod
├── simple/                          # Basic resources
├── complex-hierarchy/               # Deep nesting (10+ levels)
├── cross-references/                # Resources referencing each other
├── edge-cases/                      # Malformed, unicode, huge files
├── dependency-mod/                  # External mod simulation
└── generated/                       # Programmatically generated at scale
```

## Success Criteria

1. All existing tests continue to pass
2. New tests achieve >80% coverage of lazy loading code paths
3. At least 5 new bugs discovered and fixed
4. Test suite runs in <5 minutes
5. Tests are deterministic (no flaky tests)
6. Clear documentation of what each test validates

## Integration with Existing Tests

Align with existing test patterns in:
- `/internal/workspace/*_test.go`
- `/internal/resourceindex/*_test.go`
- `/internal/resourceloader/*_test.go`
- `/internal/dashboardserver/*_test.go`

Use existing fixtures where possible, extend as needed.

## Notes

- Websocket tests should use direct function calls, not browser automation
- Consider table-driven tests for parameter variations
- Use subtests for better failure isolation
- Profile memory usage in stress tests
