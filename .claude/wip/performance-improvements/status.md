# Performance Improvements Project - Status

**Project Start Date**: 2025-12-28
**Last Updated**: 2025-12-29

## Overall Progress

| Phase | Status | Notes |
|-------|--------|-------|
| Phase 1: Foundation | Complete | Tasks 1-4 complete |
| Phase 2: Optimizations | Complete | Tasks 5-7 complete, Task 8 skipped |
| Phase 3: Validation | Complete | Task 9 complete |

## Task Status

| Task | Title | Status | Assigned To | PR | Notes |
|------|-------|--------|-------------|-----|-------|
| 1 | Add Performance Instrumentation | Complete | - | - | Foundation for all measurements |
| 2 | Create Mod Loading Test Suite | Complete | - | - | Safety net before changes |
| 3 | Create Performance Benchmark Tests | Complete | - | - | Measurement framework |
| 4 | Baseline Performance Measurement | Complete | - | - | See baseline_results.md |
| 5 | Parallelize File I/O | Complete | - | - | 34% improvement for 100 files |
| 6 | Parallelize HCL Parsing | Complete | - | - | 58% improvement for 50 files |
| 7 | Optimize Database Client Creation | Complete | - | - | Concurrent with telemetry/modinstall |
| 8 | Cache Available Dashboards Payload | Skipped | - | - | Can be done separately |
| 9 | Final Performance Validation | Complete | - | - | See final_results.md |

## Performance Results Summary

### Baseline (2025-12-28)

| Mod Size | Load Time | Memory | Allocations |
|----------|-----------|--------|-------------|
| Small | 10.13 ms | 16.38 MB | 121,317 |
| Medium | 67.15 ms | 133.08 MB | 589,952 |
| Large | 444.19 ms | 1,112.31 MB | 2,185,732 |
| XLarge | - | - | - |

**Key Bottleneck**: `getSourceDefinition` string splitting accounts for 62.8% of allocations

### After Each Optimization

#### After Task 5 (Parallel File I/O)
| Benchmark | Sequential | Parallel | Improvement |
|-----------|------------|----------|-------------|
| LoadFileData (100 files) | 882,939 ns | 579,267 ns | 34% |

Note: Powerpipe workspace benchmarks unchanged because test mods have only 5 files.

#### After Task 6 (Parallel HCL Parsing)
| Benchmark | Sequential | Parallel | Improvement |
|-----------|------------|----------|-------------|
| ParseHclFiles (50 files) | 1,907,639 ns | 798,707 ns | 58% |

Note: Improvement primarily benefits large mods with many HCL files. Test mods have only 5 files so workspace benchmarks unchanged.

#### After Task 7 (Async DB Client)
| Component | Before | After | Improvement |
|-----------|--------|-------|-------------|
| Init with telemetry + modinstall + DB | Sequential | Concurrent | Up to 80% for slow DB |

Note: DB client creation now runs in parallel with telemetry init and mod installation. Improvement varies based on DB connection time - more significant with remote/slow databases (200-500ms savings possible).

#### After Task 8 (Payload Caching)
| Operation | Time | Improvement |
|-----------|------|-------------|
| First request | - | - |
| Cached request | - | - |

### Final Results
| Mod Size | Baseline | Final | Time Improvement | Memory Improvement |
|----------|----------|-------|------------------|-------------------|
| Small | 10.13 ms / 16.38 MB | 8.82 ms / 14.38 MB | 13% | 12% |
| Medium | 67.15 ms / 133.08 MB | 47.11 ms / 83.54 MB | 30% | 37% |
| Large | 444.19 ms / 1,112 MB | 239.85 ms / 413.70 MB | **46%** | **63%** |
| XLarge | - | - | - | - |

## Blockers & Issues

| Issue | Task | Description | Resolution |
|-------|------|-------------|------------|
| - | - | - | - |

## Decisions Log

| Date | Decision | Rationale |
|------|----------|-----------|
| - | - | - |

## Next Steps

1. [x] Task 1: Add Performance Instrumentation - Complete
2. [x] Task 2: Create Mod Loading Test Suite - Complete
3. [x] Task 3: Create Performance Benchmark Tests - Complete
4. [x] Task 4: Run Baseline Performance Measurement - Complete
5. [x] Task 5: Parallelize File I/O (pipe-fittings) - Complete
6. [x] Task 6: Parallelize HCL Parsing (pipe-fittings) - Complete
7. [x] Task 7: Optimize Database Client Creation (Powerpipe) - Complete
8. [ ] Task 8: Cache Available Dashboards Payload (Powerpipe) - Skipped
9. [x] Task 9: Final Performance Validation - Complete

## Project Complete

**Goal Achieved**: 46% faster load time, 63% less memory for large mods

Remaining work:
- [ ] Create pipe-fittings PR for parallel I/O, parallel parsing, and getSourceDefinition fix
- [ ] Optionally implement Task 8 (payload caching) for WebSocket performance

## Additional Optimization Analysis (2025-12-28)

See `additional_optimizations.md` for detailed analysis of further optimization opportunities.

### New Top Memory Allocators (Post-Optimization)
| Allocator | % of Memory | Priority |
|-----------|-------------|----------|
| cty.ObjectVal | 23.82% | HIGH |
| hclsyntax.emitToken | 15.42% | MEDIUM |
| cty.ObjectWithOptionalAttrs | 13.50% | HIGH |
| bufio.Scanner.Scan | 8.03% | LOW |
| gohcl.getFieldTags | 5.28% | MEDIUM |

### HIGH Priority Optimizations Identified
1. **CTY type reflection caching** (pipe-fittings) - 20-30% reduction in ObjectVal allocations
2. **Connection value map optimization** (pipe-fittings) - Reduce ObjectVal per connection
3. **Credential environment map caching** (pipe-fittings) - Eliminate repeated cty.StringVal

### MEDIUM Priority Optimizations Identified
4. PowerpipeModResources map pre-allocation
5. Equals() function optimization with content hashing
6. Incremental HCL parsing for workspace reload
7. Execution tree object pooling
8. Available dashboards payload caching (Task 8)

### Expected Additional Impact (if all HIGH priority implemented)
| Metric | Current | Expected |
|--------|---------|----------|
| Large Mod Load | 239.85 ms | ~160-180 ms |
| Large Mod Memory | 413.70 MB | ~280-320 MB |

**Combined total from baseline**: ~65-70% faster, ~75-80% less memory

## Lazy Loading Dashboard Server Fix (2025-12-29)

### Issue
Dashboard server in lazy-load mode was showing blank dashboards with error:
```
does not define either a 'sql' property or a 'query' property
```

### Root Cause
In `pipe-fittings/parse/decode_body.go`, the `resolveReferences` function wasn't properly handling nil pointer fields when resolving query references like `query = query.some_query`.

The code checked:
```go
if fieldVal.Type().Kind() == reflect.Pointer && !fieldVal.IsNil() {
    fieldVal = fieldVal.Elem()
}
if fieldVal.Kind() == reflect.Struct { // This failed for nil pointers
```

When `Query *Query` was nil, `fieldVal.Kind()` remained `reflect.Pointer`, so the reference resolution was skipped.

### Fix
Modified `resolveReferences` to check the TYPE rather than VALUE:
- For pointer fields, check if pointed-to type is a struct implementing HclResource
- For non-pointer struct fields, check if address type implements HclResource
- This works regardless of whether the field value is nil

### Additional Fixes Made
1. **Nested block decoding** (`parser.go`): Added `decodeNestedBlocks` to parse inline children (card, table, container) in dashboards
2. **Mod name mapping** (`index.go`, `lazy_workspace.go`): Added mapping from full mod paths (`github.com/turbot/...`) to short names (`aws_insights`)
3. **Server error handling** (`api.go`): Fixed silent port binding failure that was using `log.Fatalf` in goroutine

### Verification
- ✅ Dashboard with inline SQL cards executes correctly
- ✅ Dashboard with query references (`query = query.some_query`) resolves and executes
- ✅ Nested containers with cards work properly
- ✅ No errors in server logs

## Notes

- Phase 1 tasks must be completed sequentially
- Phase 2 tasks can be done in parallel using git worktrees
- All changes should be atomic and reversible
- Keep PRs focused and reviewable
