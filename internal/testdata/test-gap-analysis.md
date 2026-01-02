# Lazy Loading Test Gap Analysis

## Overview

This document provides a comprehensive analysis of test coverage for the lazy loading implementation in Powerpipe. It identifies gaps between existing tests and the functionality that needs testing, prioritized by risk.

## Test Inventory

### Existing Test Files

| Package | Test File | Coverage Focus |
|---------|-----------|----------------|
| `workspace` | `lazy_workspace_test.go` | Basic lazy loading, startup, on-demand loading, cache |
| `workspace` | `workspace_behavior_test.go` | Resource type accessibility, references, hierarchy |
| `resourceindex` | `scanner_test.go` | HCL scanning, pattern matching, line numbers |
| `resourceindex` | `index_test.go` | Index operations, concurrent access, payload building |
| `resourceloader` | `loader_test.go` | Load, cache hits, invalidation, preload |
| `resourceloader` | `resolver_test.go` | Dependencies, circular detection, topological order |
| `dashboardserver` | `server_lazy_test.go` | Lazy mode detection, payload conversion |

### Test Counts by Package

| Package | # Test Functions | # Benchmark Functions |
|---------|------------------|----------------------|
| `workspace` (lazy) | 13 | 0 |
| `resourceindex` | 24 + 19 = 43 | 0 |
| `resourceloader` | 9 + 11 = 20 | 0 |
| `dashboardserver` (lazy) | 8 | 1 |

---

## Coverage Heat Map

### lazy_workspace.go (537 lines)

| Function | Tested? | Coverage Level | Notes |
|----------|---------|----------------|-------|
| `NewLazyWorkspace()` | ✅ Partial | MEDIUM | Happy path only; error paths not tested |
| `GetWorkspaceForExecution()` | ❌ No | NONE | Eager loading fallback untested |
| `buildResourceIndex()` | ✅ Indirect | LOW | Tested via NewLazyWorkspace |
| `scanDependencyMods()` | ❌ No | NONE | Dependency mod scanning untested |
| `scanModInfo()` | ❌ No | NONE | mod.pp parsing edge cases untested |
| `preloadResources()` | ❌ No | NONE | Background preloading untested |
| `findByPattern()` | ❌ No | NONE | Pattern matching untested |
| `matchPattern()` | ❌ No | NONE | Wildcard matching untested |
| `GetResource()` | ✅ Partial | LOW | Only happy path via LoadResource |
| `LoadBenchmark()` | ❌ No | NONE | Benchmark loading untested |
| `LoadBenchmarkForExecution()` | ❌ No | NONE | Critical for execution, untested |
| `resolveChildrenRecursively()` | ❌ No | NONE | Child resolution untested |
| `InvalidateResource()` | ✅ Yes | HIGH | Good coverage |
| `InvalidateAll()` | ✅ Yes | HIGH | Good coverage |
| `CacheStats()` | ✅ Yes | HIGH | Good coverage |
| `IndexStats()` | ✅ Yes | HIGH | Good coverage |

### scanner.go (669 lines)

| Function | Tested? | Coverage Level | Notes |
|----------|---------|----------------|-------|
| `NewScanner()` | ✅ Yes | HIGH | Covered |
| `ScanFile()` | ✅ Yes | HIGH | Good coverage |
| `ScanFileWithOffsets()` | ✅ Yes | MEDIUM | Byte offsets tested |
| `scanReader()` | ✅ Indirect | MEDIUM | Via ScanBytes |
| `scanReaderWithOffsets()` | ✅ Indirect | MEDIUM | Via ScanFileWithOffsets |
| `processBlockLine()` | ✅ Indirect | MEDIUM | Via scan tests |
| `parseBlockStart()` | ✅ Indirect | MEDIUM | Tested via resource type tests |
| `parseAttribute()` | ✅ Indirect | LOW | Implicit in other tests |
| `finalizeBlock()` | ✅ Indirect | LOW | Implicit |
| `ScanDirectory()` | ✅ Yes | HIGH | Directory scanning covered |
| `ScanDirectoryParallel()` | ✅ Yes | MEDIUM | Parallel scanning covered |
| `ScanDirectoryWithModName()` | ❌ No | NONE | Mod-specific scanning untested |
| `MarkTopLevelResources()` | ✅ Yes | HIGH | Top-level marking covered |
| `SetParentNames()` | ✅ Yes | MEDIUM | Parent name setting covered |

**Scanner Edge Cases NOT Tested:**
- Unicode characters in resource names
- Very long lines (>64KB)
- Files with mixed line endings (CRLF/LF)
- Deeply nested brace structures (>10 levels)
- Malformed heredocs
- Empty block labels
- Special characters in tags
- Escape sequences in strings

### index.go (239 lines)

| Function | Tested? | Coverage Level | Notes |
|----------|---------|----------------|-------|
| `NewResourceIndex()` | ✅ Yes | HIGH | Good coverage |
| `RegisterModName()` | ❌ No | NONE | Mod name registration untested |
| `ResolveModName()` | ❌ No | NONE | Mod name resolution untested |
| `Add()` | ✅ Yes | HIGH | Well tested |
| `Get()` | ✅ Yes | HIGH | Well tested |
| `GetByType()` | ✅ Yes | HIGH | Covered |
| `List()` | ✅ Yes | HIGH | Covered |
| `Remove()` | ✅ Yes | HIGH | Covered |
| `Dashboards()` | ✅ Yes | HIGH | Covered |
| `Benchmarks()` | ✅ Yes | HIGH | Covers both types |
| `TopLevelBenchmarks()` | ✅ Yes | MEDIUM | Covered |
| `GetChildren()` | ✅ Yes | HIGH | Hierarchy tests |
| `BuildAvailableDashboardsPayload()` | ✅ Yes | HIGH | Good payload tests |

### loader.go (300 lines)

| Function | Tested? | Coverage Level | Notes |
|----------|---------|----------------|-------|
| `NewLoader()` | ✅ Yes | HIGH | Covered |
| `SetResourceProvider()` | ❌ No | NONE | Not directly tested |
| `Load()` | ✅ Yes | HIGH | Good coverage |
| `LoadDashboard()` | ❌ No | NONE | Dashboard with children untested |
| `LoadBenchmark()` | ❌ No | NONE | Benchmark loading untested |
| `loadFromDisk()` | ✅ Indirect | MEDIUM | Via Load tests |
| `loadChildren()` | ❌ No | NONE | Child loading untested |
| `loadBenchmarkChildren()` | ❌ No | NONE | Benchmark child loading untested |
| `loadControlDependencies()` | ❌ No | NONE | Query loading for controls untested |
| `loadDetectionDependencies()` | ❌ No | NONE | Detection query loading untested |
| `Preload()` | ✅ Yes | MEDIUM | Parallel preload tested |
| `PreloadByType()` | ✅ Yes | MEDIUM | Type-based preload tested |
| `Stats()` | ✅ Yes | HIGH | Covered |
| `Clear()` | ✅ Yes | HIGH | Covered |
| `Invalidate()` | ✅ Yes | HIGH | Covered |

### parser.go (341 lines)

| Function | Tested? | Coverage Level | Notes |
|----------|---------|----------------|-------|
| `parseResource()` | ✅ Indirect | LOW | Implicit via loader tests |
| `readResourceBlock()` | ❌ No | NONE | Block reading untested |
| `readByLines()` | ❌ No | NONE | Line-based reading untested |
| `decodeResourceBlock()` | ❌ No | NONE | HCL decoding untested |
| `isNonCriticalDecodeError()` | ❌ No | NONE | Error classification untested |
| `createResource()` | ❌ No | NONE | Factory function untested |
| `decodeNestedBlocks()` | ❌ No | NONE | Nested block parsing untested |
| `isDashboardChildType()` | ❌ No | NONE | Type checking untested |

### resolver.go (302 lines)

| Function | Tested? | Coverage Level | Notes |
|----------|---------|----------------|-------|
| `NewDependencyResolver()` | ✅ Yes | HIGH | Covered |
| `ResolveWithDependencies()` | ✅ Partial | MEDIUM | Basic test, error paths untested |
| `GetDependencies()` | ✅ Yes | HIGH | Good coverage |
| `GetDependencyOrder()` | ✅ Yes | HIGH | Topological sort tested |
| `GetTransitiveDependencies()` | ✅ Yes | HIGH | Covered |
| `GetDependents()` | ✅ Yes | HIGH | Covered |
| `HasCircularDependency()` | ✅ Yes | HIGH | Cycle detection covered |
| `BuildDependencyGraph()` | ✅ Yes | MEDIUM | Graph building covered |

### preload.go (172 lines)

| Function | Tested? | Coverage Level | Notes |
|----------|---------|----------------|-------|
| `DefaultPreloadOptions()` | ❌ No | NONE | Untested |
| `PreloadWithDependencies()` | ❌ No | NONE | Complex preload untested |
| `PreloadBenchmark()` | ❌ No | NONE | Untested |
| `PreloadDashboard()` | ❌ No | NONE | Untested |

### server.go - Lazy Mode Paths

| Function | Tested? | Coverage Level | Notes |
|----------|---------|----------------|-------|
| `NewServerWithLazyWorkspace()` | ❌ No | NONE | Server creation untested |
| `isLazyMode()` | ✅ Partial | LOW | Only nil check tested |
| `getActiveWorkspace()` | ❌ No | NONE | Workspace selection untested |
| `buildAvailableDashboardsPayload()` | ✅ Partial | MEDIUM | Index path tested, not integration |

---

## Prioritized Gap List

### HIGH Priority (Core Functionality Gaps)

| ID | Component | Gap Description | Risk |
|----|-----------|-----------------|------|
| H1 | `lazy_workspace.go:150` | `GetWorkspaceForExecution()` - Eager loading fallback never tested | HIGH |
| H2 | `lazy_workspace.go:430` | `LoadBenchmarkForExecution()` - Critical for check execution | HIGH |
| H3 | `loader.go:88` | `LoadBenchmark()` - Benchmark loading with children | HIGH |
| H4 | `lazy_workspace.go:447` | `resolveChildrenRecursively()` - Child resolution for execution | HIGH |
| H5 | `preload.go:27` | `PreloadWithDependencies()` - Complex parallel preload | HIGH |
| H6 | `parser.go:117` | `decodeResourceBlock()` - HCL block decoding | HIGH |
| H7 | `parser.go:240` | `decodeNestedBlocks()` - Dashboard nested children | HIGH |

### MEDIUM Priority (Edge Cases and Error Paths)

| ID | Component | Gap Description | Risk |
|----|-----------|-----------------|------|
| M1 | `lazy_workspace.go:82` | `NewLazyWorkspace()` error paths - mod.pp failures | MEDIUM |
| M2 | `lazy_workspace.go:203` | `scanDependencyMods()` - Dependency scanning | MEDIUM |
| M3 | `lazy_workspace.go:297` | `preloadResources()` - Background preloading | MEDIUM |
| M4 | `scanner.go:574` | `ScanDirectoryWithModName()` - Multi-mod scanning | MEDIUM |
| M5 | `loader.go:128` | `loadChildren()` - Child loading errors | MEDIUM |
| M6 | `loader.go:151` | `loadBenchmarkChildren()` - Recursive benchmark children | MEDIUM |
| M7 | `index.go:42` | `RegisterModName()` / `ResolveModName()` | MEDIUM |
| M8 | `parser.go:56` | `readResourceBlock()` - Byte offset seeking | MEDIUM |
| M9 | `parser.go:80` | `readByLines()` - Line-based block reading | MEDIUM |
| M10 | `parser.go:171` | `isNonCriticalDecodeError()` - Error classification | MEDIUM |

### LOW Priority (Nice to Have)

| ID | Component | Gap Description | Risk |
|----|-----------|-----------------|------|
| L1 | `lazy_workspace.go:319` | `findByPattern()` / `matchPattern()` - Wildcard matching | LOW |
| L2 | `scanner.go` | Unicode edge cases in resource names | LOW |
| L3 | `scanner.go` | Very large files (>1MB) with many resources | LOW |
| L4 | `scanner.go` | Mixed line endings (CRLF/LF) | LOW |
| L5 | `loader.go:188` | `loadControlDependencies()` | LOW |
| L6 | `loader.go:197` | `loadDetectionDependencies()` | LOW |
| L7 | `server.go:72` | `NewServerWithLazyWorkspace()` | LOW |
| L8 | `preload.go` | `DefaultPreloadOptions()` | LOW |

---

## Recommended Test Additions

### Phase 1: Critical Path Testing (HIGH Priority)

1. **`TestLazyWorkspace_GetWorkspaceForExecution`**
   - Test eager workspace is loaded on first execution request
   - Test error caching (second call returns cached error)
   - Test event handler transfer to eager workspace
   - Location: `workspace/lazy_workspace_test.go`

2. **`TestLazyWorkspace_LoadBenchmarkForExecution`**
   - Test benchmark with nested children is properly resolved
   - Test Children field is populated (not just cached)
   - Test parent relationships are set
   - Test with 3+ level deep hierarchies
   - Location: `workspace/lazy_workspace_test.go`

3. **`TestLoader_LoadBenchmark`**
   - Test loading benchmark with child controls
   - Test loading benchmark with nested benchmarks
   - Test query dependencies are resolved
   - Location: `resourceloader/loader_test.go`

4. **`TestLoader_PreloadWithDependencies`**
   - Test parallel preloading respects dependency order
   - Test progress callback is called
   - Test error callback works
   - Test context cancellation
   - Location: `resourceloader/loader_test.go` or new `preload_test.go`

5. **`TestParser_DecodeResourceBlock`**
   - Test decoding each resource type
   - Test with variable references (should be non-critical)
   - Test nested blocks in dashboards
   - Test error handling for malformed HCL
   - Location: `resourceloader/parser_test.go` (new file)

### Phase 2: Edge Cases (MEDIUM Priority)

6. **`TestLazyWorkspace_NewLazyWorkspace_Errors`**
   - Test with missing mod.pp
   - Test with malformed mod.pp
   - Test with unreadable directory
   - Location: `workspace/lazy_workspace_test.go`

7. **`TestLazyWorkspace_ScanDependencyMods`**
   - Test with dependency mods in .powerpipe/mods
   - Test mod name mapping works
   - Test versioned mod paths (@v1.2.0)
   - Location: `workspace/lazy_workspace_test.go`

8. **`TestScanner_WithModName`**
   - Test `ScanDirectoryWithModName()` preserves mod context
   - Test resources have correct mod prefix
   - Location: `resourceindex/scanner_test.go`

9. **`TestIndex_ModNameMapping`**
   - Test `RegisterModName()` and `ResolveModName()`
   - Test with multiple dependency mods
   - Location: `resourceindex/index_test.go`

10. **`TestParser_ReadResourceBlock`**
    - Test byte offset reading
    - Test line-based fallback
    - Test with heredocs spanning many lines
    - Location: `resourceloader/parser_test.go`

### Phase 3: Robustness (LOW Priority)

11. **`TestLazyWorkspace_PatternMatching`**
    - Test `findByPattern()` with wildcards
    - Test edge cases: *, prefix*, *suffix
    - Location: `workspace/lazy_workspace_test.go`

12. **`TestScanner_UnicodeNames`**
    - Test resource names with unicode characters
    - Test tags with unicode values
    - Location: `resourceindex/scanner_test.go`

13. **`TestLoader_ControlDependencies`**
    - Test control with query reference
    - Test control with inline SQL (no dependency)
    - Location: `resourceloader/loader_test.go`

---

## Existing Tests to Verify/Extend

These tests exist but may need extension:

1. `TestLazyWorkspace_OnDemandLoading` - Add error path testing
2. `TestResolver_ResolveWithDependencies` - Add cancellation testing
3. `TestScanner_Performance` - Add memory usage assertions
4. `TestIndex_ConcurrentAccess` - Add more stress testing

---

## Test Infrastructure Needs

1. **Test Fixtures**: Need more complex test mods with:
   - Deep benchmark hierarchies (4+ levels)
   - Cross-mod references (dependency mods)
   - Large numbers of resources (1000+)
   - Malformed/edge case HCL files

2. **Helper Functions**: Consider adding:
   - `setupLazyWorkspaceWithMod()` - Create temp mod with specified structure
   - `assertBenchmarkChildrenResolved()` - Deep verify children
   - `waitForPreload()` - Helper for async preload tests

3. **Mocking**: May need mocks for:
   - File system operations (simulate read errors)
   - HCL parser (inject parse errors)
   - Cache (simulate eviction pressure)

---

## Summary

| Priority | Total Gaps | Coverage % |
|----------|------------|------------|
| HIGH | 7 | ~40% of critical paths tested |
| MEDIUM | 10 | ~50% of edge cases tested |
| LOW | 8 | ~30% of nice-to-haves tested |

**Overall Assessment**: The lazy loading implementation has reasonable test coverage for basic happy paths, but lacks:
- Error path testing
- Execution-mode testing (GetWorkspaceForExecution, LoadBenchmarkForExecution)
- Complex preloading scenarios
- Parser/decoder unit tests
- Integration tests between components

**Recommended Action**: Start with HIGH priority gaps (Phase 1) as these cover the critical execution paths that users will rely on for `powerpipe check` and `powerpipe dashboard` commands.
