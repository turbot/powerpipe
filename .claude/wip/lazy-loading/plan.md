# Lazy Loading Architecture - Project Plan

## Overview

Transform Powerpipe's memory model from "load everything upfront" to "lazy loading with LRU cache" to achieve **bounded memory usage regardless of mod size**.

### Goals
- **Primary**: Memory usage bounded by cache size, not mod size
- **Secondary**: Faster initial startup (index-only load)
- **Tertiary**: Maintain all existing functionality with no regressions

### Current State (Baseline from Task 2)
- Large mod (200 dashboards): 14 MB heap, 250K objects, 394 MB total allocations
- XLarge mod (500 dashboards): ~35 MB heap
- Memory scales at 13.6x for 20x dashboards (sub-linear, which is good)
- All resources parsed and held in memory
- HCL AST (`Remain` fields) kept unnecessarily
- Full source text stored in every resource

### Target State
- **Primary Goal**: Faster startup via index-only initial load
- **Secondary Goal**: Reduced allocation churn (394 MB → <100 MB)
- **Tertiary Goal**: Bounded memory with LRU cache
- Lightweight index always in memory (~1 KB per 100 resources)
- On-demand parsing only when resources are executed

## Task Breakdown

| Task | Title | Status | Phase | Dependencies |
|------|-------|--------|-------|--------------|
| 1 | Comprehensive Behavior Tests | Pending | Foundation | - |
| 2 | Memory Benchmarking Infrastructure | Pending | Foundation | - |
| 3 | Resource Access Pattern Analysis | Pending | Foundation | - |
| 4 | Resource Index Design & Implementation | Pending | Core | 1, 2 |
| 5 | File Scanner for Index Building | Pending | Core | 4 |
| 6 | LRU Cache Implementation | Pending | Core | 4 |
| 7 | On-Demand Resource Parser | Pending | Core | 5, 6 |
| 8 | Dependency Resolution for Lazy Loading | Pending | Core | 7 |
| 9 | Workspace Integration | Pending | Integration | 8 |
| 10 | Dashboard Server Integration | Pending | Integration | 9 |
| 11 | CLI Commands Integration | Pending | Integration | 9 |
| 12 | Post-Parse Cleanup (Remain fields) | Pending | Optimization | 7 |
| 13 | Lazy Source Definition Loading | Pending | Optimization | 9 |
| 14 | String Interning | Pending | Optimization | 9 |
| 15 | Final Validation & Regression Testing | Pending | Validation | 10, 11, 12, 13, 14 |

## Execution Strategy

**Mixed**: Foundation tasks in parallel, then Core tasks sequentially, Integration tasks can partially parallelize.

```
Phase 1 (Foundation) - Parallel:
  [Task 1] ─┬─► [Task 4+]
  [Task 2] ─┤
  [Task 3] ─┘

Phase 2 (Core) - Sequential:
  [Task 4] → [Task 5] → [Task 6] → [Task 7] → [Task 8]

Phase 3 (Integration) - Mixed:
  [Task 9] → [Task 10] ─┬─► [Task 15]
            [Task 11] ─┤
            [Task 12] ─┤
            [Task 13] ─┤
            [Task 14] ─┘
```

## Repository Split

This project spans two repositories. Changes must be coordinated.

### pipe-fittings (shared library)
Used by Powerpipe, Steampipe, and Flowpipe. Changes here affect all products.

**Files with `Remain hcl.Body` (8 files):**
- `modconfig/mod.go`
- `modconfig/variable.go`
- `modconfig/local.go`
- `modconfig/resource_metadata.go`
- `modconfig/hcl_resource_impl.go`
- `modconfig/mod_tree_item_impl.go`
- `modconfig/resource_with_metadata_impl.go`
- `modconfig/interfaces.go`

**Key files to modify:**
- `modconfig/resource_metadata.go` - Contains `SourceDefinition` field (Task 13)
- `parse/parser.go` - Add post-parse cleanup hook (Task 12)
- `modconfig/*.go` - Add `ClearRemain()` methods (Task 12)

### powerpipe (this repository)
Powerpipe-specific changes. No impact on other products.

**Files with `Remain hcl.Body` (25 files):**
- `internal/resources/dashboard.go`
- `internal/resources/control.go`
- `internal/resources/query.go`
- `internal/resources/control_benchmark.go`
- `internal/resources/dashboard_card.go`
- `internal/resources/dashboard_chart.go`
- `internal/resources/dashboard_container.go`
- `internal/resources/dashboard_flow.go`
- `internal/resources/dashboard_graph.go`
- `internal/resources/dashboard_hierarchy.go`
- `internal/resources/dashboard_image.go`
- `internal/resources/dashboard_input.go`
- `internal/resources/dashboard_table.go`
- `internal/resources/dashboard_text.go`
- `internal/resources/dashboard_category.go`
- `internal/resources/dashboard_node.go`
- `internal/resources/dashboard_edge.go`
- `internal/resources/detection.go`
- `internal/resources/detection_benchmark.go`
- And 6 more implementation files

**New packages to create:**
- `internal/resourceindex/` - Resource index (Task 4)
- `internal/resourcecache/` - LRU cache (Task 6)
- `internal/resourceloader/` - On-demand loading (Task 7)

### Change Summary by Task

| Task | powerpipe | pipe-fittings | Notes |
|------|-----------|---------------|-------|
| 1-3 | Tests only | None | Foundation |
| 4 | New package | None | Index is Powerpipe-only |
| 5 | New package | None | Scanner is Powerpipe-only |
| 6 | New package | None | Cache is Powerpipe-only |
| 7 | New package | Optional helpers | Core in Powerpipe |
| 8 | New code | None | Resolver is Powerpipe-only |
| 9-11 | Modifications | None | Integration is Powerpipe-only |
| 12 | 25 files | 8 files | Split cleanup |
| 13 | Wrapper only | `ResourceMetadata` | Core change in pipe-fittings |
| 14 | New package | None | Interning is Powerpipe-only |
| 15 | Tests only | None | Validation |

## HCL Variable Resolution Compatibility

### Analysis Summary

The lazy loading approach is **fully compatible** with HCL variable parsing and expansion. Here's why:

### How HCL Parsing Works in pipe-fittings

1. **Variables parsed first**: `ModVariableMap` is populated in an initial pass before resource parsing
2. **EvalContext built incrementally**: As resources are decoded, they're added to `ReferenceTypeValueMap`
3. **Dependencies tracked in DAG**: `topsort.Graph` ensures resources decode in dependency order
4. **Multiple decode passes**: Parser iterates until all `UnresolvedBlocks` are resolved

### Why Lazy Loading Works

| Operation | Requirement | Lazy Loading Approach |
|-----------|-------------|----------------------|
| `available_dashboards` | Name, title, tags, children | ✅ Index-only (no HCL eval needed) |
| Execute dashboard | Full parse + resolved refs | ✅ Parse on-demand with full EvalContext |
| Execute benchmark | Hierarchy + control queries | ✅ Parse benchmark tree on-demand |

**Key insight**: We're not avoiding parsing - we're **deferring** it. When execution happens:
1. Variables are already resolved (parsed upfront)
2. Full parsing infrastructure is used
3. Dependencies resolve normally via the DAG

### Critical Design Decisions

1. **Variables MUST be parsed upfront** - they're needed for EvalContext
2. **Index includes parent-child relationships** - scanner extracts `children = [...]` without full parsing
3. **On-demand parsing uses existing infrastructure** - reuses `ModParseContext`, `Decoder`, dependency resolution

### Benchmark Children Handling

The `available_dashboards` payload needs benchmark hierarchy. Two options:

- **Option A (Chosen)**: Index stores `ChildNames` - scanner extracts from HCL without full parsing
- **Option B**: Parse benchmarks for hierarchy - more expensive, not needed

The scanner (Task 5) will extract `children = [benchmark.child1, ...]` using regex, avoiding full HCL parsing.

## Dependencies

### Task Dependencies
- Tasks 1-3: Independent (can run in parallel)
- Task 4: Requires Tasks 1-3 complete (need tests before changing architecture)
- Tasks 5-8: Sequential chain (each builds on previous)
- Task 9: Requires Task 8 (core must be complete)
- Tasks 10-14: Require Task 9 (workspace integration must work first)
- Task 15: Requires all previous tasks

### External Dependencies
- **pipe-fittings**: Tasks 12, 13 require changes to pipe-fittings
- **hashicorp/hcl**: No changes needed (we work around it)

## Shared Resources

### Files Modified by Multiple Tasks

| File | Tasks | Coordination |
|------|-------|--------------|
| `internal/workspace/powerpipe_workspace.go` | 9, 10, 11 | Task 9 owns structure, 10-11 add methods |
| `internal/resources/mod_resources.go` | 4, 9 | Task 4 adds index, Task 9 integrates |
| `pipe-fittings/parse/parser.go` | 5, 7 | Task 5 adds scanner, Task 7 adds on-demand parsing |
| `pipe-fittings/modconfig/mod.go` | 4, 8 | Task 4 adds index field, Task 8 adds lazy resolution |

### New Files

| File | Task | Purpose |
|------|------|---------|
| `internal/resourceindex/index.go` | 4 | Resource index structure |
| `internal/resourceindex/scanner.go` | 5 | File scanner for building index |
| `internal/resourcecache/cache.go` | 6 | LRU cache implementation |
| `internal/resourcecache/loader.go` | 7 | On-demand resource loader |

## Integration Plan

### Phase Gates
Each phase has a gate requiring:
1. All phase tasks complete
2. All tests passing
3. Memory benchmark shows expected improvement
4. No regressions in existing functionality

### Rollback Strategy
- Feature flag: `POWERPIPE_LAZY_LOADING=true/false`
- Default: `false` during development, `true` after validation
- Allows quick rollback if issues found in production

### Backwards Compatibility
- All existing CLI commands work unchanged
- All existing API endpoints work unchanged
- Snapshot format unchanged
- Mod format unchanged

## Success Metrics

| Metric | Baseline | Target | Validation |
|--------|----------|--------|------------|
| Heap memory (large mod) | 14 MB | < 10 MB | Benchmark |
| Total allocations | 394 MB | < 100 MB | Benchmark |
| Initial load time (large) | 241 ms | < 100 ms | Benchmark |
| Initial load time (xlarge) | 1049 ms | < 200 ms | Benchmark |
| `available_dashboards` latency | ~100 ms | < 10 ms | Benchmark |
| First resource access | 0 ms (preloaded) | < 10 ms | Benchmark |
| All tests passing | 100% | 100% | CI |
| Cache hit rate | N/A | > 90% | Telemetry |

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Code assumes all resources loaded | High | Task 3 identifies all access patterns |
| Reference resolution complexity | Medium | Task 8 dedicated to this |
| Cache thrashing with many resources | Low | Configurable cache size |
| Performance regression for small mods | Low | Skip lazy loading for < 100 resources |

## Notes

- This is a significant architectural change - prioritize testing
- pipe-fittings changes should be in separate PR, merged first
- Consider adding metrics/telemetry for cache hit rates
- Document new architecture for future maintainers
