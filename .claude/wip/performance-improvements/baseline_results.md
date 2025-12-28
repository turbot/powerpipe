# Baseline Performance Results

**Date**: 2025-12-28
**Commit**: 4c880fde22ced34743d976aadf110707b5dc61e0
**Branch**: performance-improvements

## Test Environment

- **Processor**: Apple M4 Pro
- **Memory**: 48 GB
- **Disk**: SSD
- **OS**: macOS (Darwin 24.6.0)
- **Go Version**: go1.25.1 darwin/arm64

## Mod Sizes

| Size | Dashboards | Queries | Controls | Benchmarks | Files | Total Lines | Disk Size |
|------|------------|---------|----------|------------|-------|-------------|-----------|
| Small | 10 | 20 | 30 | 5 | 5 | 1,154 | 32K |
| Medium | 50 | 100 | 150 | 20 | 5 | 5,694 | 104K |
| Large | 200 | 400 | 500 | 50 | 5 | 20,904 | 364K |
| XLarge | 500 | 1000 | 1500 | 100 | 5 | 56,004 | 976K |

## Workspace Loading Times

| Size | Load Time | Memory Allocated | Allocations |
|------|-----------|------------------|-------------|
| Small | 10.13 ms | 16.38 MB | 121,317 |
| Medium | 67.15 ms | 133.08 MB | 589,952 |
| Large | 444.19 ms | 1,112.31 MB | 2,185,732 |
| XLarge | (not benchmarked) | - | - |

## Payload Building Times

| Size | Build Time | Memory Allocated | Allocations |
|------|------------|------------------|-------------|
| Small | 10.6 us | 12.67 KB | 125 |
| Medium | 43.8 us | 59.95 KB | 539 |
| Large | 137.7 us | 210.90 KB | 1,837 |

## Timing Breakdown (Large Mod)

From benchmark timing instrumentation:

| Operation | Time (ms) | % of Total |
|-----------|-----------|------------|
| workspace.Load | 444.19 | 100% |
| ├─ workspace.SetModfileExists | 0.01 | <0.1% |
| ├─ workspace.LoadExclusions | 0.01 | <0.1% |
| ├─ workspace.LoadWorkspaceMod | ~444 | ~99.9% |
| └─ workspace.verifyResourceRuntimeDependencies | 0.01 | <0.1% |

**Key Insight**: Nearly all time is spent in `LoadWorkspaceMod`, which calls into pipe-fittings for HCL parsing.

## CPU Profile Top Functions (Large Mod)

```
Duration: 11.39s, Total samples = 22.98s (201.83%)
      flat  flat%   sum%        cum   cum%
     4.58s 19.93% 19.93%      4.58s 19.93%  runtime.pthread_kill
     3.64s 15.84% 35.77%      3.64s 15.84%  runtime.madvise
     1.75s  7.62% 43.39%      1.75s  7.62%  runtime.usleep
     1.47s  6.40% 49.78%      1.47s  6.40%  runtime.pthread_cond_wait
     1.30s  5.66% 55.44%      1.30s  5.66%  internal/bytealg.IndexByteString
     0.81s  3.52% 58.96%      0.81s  3.52%  runtime.kevent
     0.74s  3.22% 62.18%      1.21s  5.27%  runtime.greyobject
     0.69s  3.00% 65.19%      3.44s 14.97%  runtime.scanobject
```

**GC Overhead**: ~40% of CPU time is in garbage collection (gcDrain, scanobject, greyobject).

## Memory Profile Top Allocators (Large Mod)

```
Total Allocated: 27,678.77 MB over all iterations
      flat  flat%   sum%        cum   cum%
 8720.01MB 31.50% 31.50%  8720.01MB 31.50%  strings.genSplit
 8668.80MB 31.32% 62.82% 17383.32MB 62.80%  pipe-fittings/parse.getSourceDefinition
 2418.90MB  8.74% 71.56%  3669.17MB 13.26%  go-cty/cty.ObjectVal
 1581.79MB  5.71% 77.28%  1581.79MB  5.71%  hcl/v2/hclsyntax.(*tokenAccum).emitToken
 1372.79MB  4.96% 82.24%  1372.79MB  4.96%  go-cty/cty.ObjectWithOptionalAttrs
  830.23MB  3.00% 85.24%   830.23MB  3.00%  bufio.(*Scanner).Scan
  553.57MB  2.00% 87.24%   553.57MB  2.00%  hcl/v2/gohcl.getFieldTags
```

## Identified Bottlenecks

### 1. Source Definition Line Counting (High Impact)
- **Location**: `pipe-fittings/parse.getSourceDefinition`
- **Impact**: 62.8% of memory allocations (17.4 GB cumulative)
- **Root Cause**: Uses `strings.Split` to count lines for every parsed resource, creating massive string slice allocations
- **Optimization Opportunity**: Replace with byte counting or single-pass line counter

### 2. Garbage Collection Overhead (High Impact)
- **Location**: Runtime GC
- **Impact**: ~40% of CPU time
- **Root Cause**: High allocation rate from string operations and cty object creation
- **Optimization Opportunity**: Reduce allocations in hot paths, object pooling

### 3. HCL Token Emission (Medium Impact)
- **Location**: `hclsyntax.(*tokenAccum).emitToken`
- **Impact**: 5.71% of memory (1.58 GB)
- **Root Cause**: Token allocations during HCL lexing
- **Optimization Opportunity**: This is in hashicorp/hcl, limited opportunity unless we cache parsed results

### 4. CTY Object Creation (Medium Impact)
- **Location**: `cty.ObjectVal`, `cty.ObjectWithOptionalAttrs`
- **Impact**: 13.7% of memory (3.8 GB combined)
- **Root Cause**: Creating cty values for every resource attribute
- **Optimization Opportunity**: Lazy evaluation, caching common patterns

### 5. Field Tag Reflection (Low Impact)
- **Location**: `gohcl.getFieldTags`
- **Impact**: 2% of memory
- **Root Cause**: Reflection operations during struct decoding
- **Optimization Opportunity**: Cache field tag results (may already be done)

## Scaling Characteristics

| Metric | Small→Medium | Medium→Large | Factor |
|--------|--------------|--------------|--------|
| Load Time | 6.6x | 6.6x | Linear (O(n)) |
| Memory | 8.1x | 8.4x | Linear-ish |
| Allocations | 4.9x | 3.7x | Sub-linear |
| Lines of Code | 4.9x | 3.7x | - |

**Observations**:
- Load time scales approximately linearly with mod size (good)
- Memory usage scales slightly worse than linear
- No super-linear bottlenecks detected
- The dominant cost is HCL parsing and string operations

## Key Insights

1. **`getSourceDefinition` is the primary bottleneck** - accounting for 62.8% of allocations via repeated calls to `strings.Split` for line counting
2. **GC pressure is significant** - High allocation rate causes ~40% of CPU to be spent in garbage collection
3. **Parsing is sequential** - All HCL files are parsed one at a time, no parallelism
4. **Memory-bound more than CPU-bound** - The string allocations dominate the profile
5. **Payload building is negligible** - Only 137.7 us for large mod vs 444 ms for loading (0.03%)

## Optimization Priority

Based on baseline measurements, prioritize optimizations:

1. **High Impact**: Fix `getSourceDefinition` string splitting - 62.8% of allocations
2. **High Impact**: Parallelize HCL file parsing - currently sequential
3. **Medium Impact**: Reduce cty object allocations - 13.7% of memory
4. **Medium Impact**: Cache parsed HCL results between reloads
5. **Low Impact**: Payload caching - only 0.03% of startup time

## Raw Data

- Workspace loading benchmarks: `benchmark_results/baseline/workspace_load.txt`
- Payload building benchmarks: `benchmark_results/baseline/payload_build.txt`
- CPU profile: `benchmark_results/baseline/cpu.prof`
- Memory profile: `benchmark_results/baseline/mem.prof`
- CPU profile analysis: `benchmark_results/baseline/cpu_profile_top.txt`
- Memory profile analysis: `benchmark_results/baseline/mem_profile_top.txt`

## Summary

The baseline measurements reveal that mod loading is dominated by HCL parsing operations, with `getSourceDefinition`'s line-counting logic being the single largest contributor to memory allocations. The high allocation rate leads to significant GC overhead.

For a large mod (200 dashboards, 400 queries, 500 controls):
- **Load time**: 444 ms
- **Memory**: 1.1 GB per load operation
- **Primary bottleneck**: String allocations in source definition parsing

These measurements establish the "before" baseline for optimization work in Tasks 5-8.
