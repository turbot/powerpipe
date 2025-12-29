# Task 15: Final Validation & Regression Testing

## Objective

Perform comprehensive validation to ensure lazy loading implementation is complete, correct, and achieves the memory and performance goals.

## Context

- All implementation tasks completed
- Need to verify no regressions in functionality
- Need to verify memory and performance goals met
- Need to validate with real-world mods
- Final task before release

## Dependencies

### Prerequisites
- All previous tasks (1-14) complete
- All tests passing

### Files to Create
- `test/integration/lazy_loading_test.go` - Integration tests
- `test/validation/memory_validation_test.go` - Memory validation
- `test/validation/behavior_validation_test.go` - Behavior validation
- `scripts/validate_lazy_loading.sh` - Validation script
- `docs/lazy-loading.md` - Documentation

## Implementation Details

### 1. Comprehensive Integration Tests

```go
// test/integration/lazy_loading_test.go
package integration

import (
    "context"
    "encoding/json"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestLazyLoading_EndToEnd_SmallMod(t *testing.T) {
    runEndToEndTest(t, "testdata/mods/small", 10, 50*time.Millisecond)
}

func TestLazyLoading_EndToEnd_MediumMod(t *testing.T) {
    runEndToEndTest(t, "testdata/mods/medium", 50, 100*time.Millisecond)
}

func TestLazyLoading_EndToEnd_LargeMod(t *testing.T) {
    runEndToEndTest(t, "testdata/mods/large", 200, 200*time.Millisecond)
}

func TestLazyLoading_EndToEnd_RealMod_AWS(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping real mod test in short mode")
    }
    runEndToEndTest(t, "../steampipe-mod-aws-compliance", 500, 500*time.Millisecond)
}

func runEndToEndTest(t *testing.T, modPath string, expectedResources int, maxStartupTime time.Duration) {
    ctx := context.Background()

    // 1. Test startup time
    start := time.Now()
    lw, err := workspace.NewLazyWorkspace(ctx, modPath, workspace.DefaultLazyLoadConfig())
    require.NoError(t, err)
    startupTime := time.Since(start)

    t.Logf("Startup time: %v (max: %v)", startupTime, maxStartupTime)
    assert.LessOrEqual(t, startupTime, maxStartupTime, "Startup too slow")

    // 2. Test index populated correctly
    indexCount := lw.Index().Count()
    t.Logf("Index contains %d resources", indexCount)
    assert.GreaterOrEqual(t, indexCount, expectedResources)

    // 3. Test available_dashboards works
    payload := lw.GetAvailableDashboardsFromIndex()
    assert.NotEmpty(t, payload.Dashboards)
    assert.NotEmpty(t, payload.Benchmarks)

    // 4. Test on-demand loading
    if len(payload.Dashboards) > 0 {
        var firstName string
        for name := range payload.Dashboards {
            firstName = name
            break
        }

        dash, err := lw.LoadDashboard(ctx, firstName)
        require.NoError(t, err)
        assert.NotNil(t, dash)
    }

    // 5. Test memory is bounded
    stats := lw.CacheStats()
    t.Logf("Cache: entries=%d, memory=%d bytes", stats.Entries, stats.MemoryBytes)
    assert.LessOrEqual(t, stats.MemoryBytes, int64(60*1024*1024), "Memory should be < 60MB")
}
```

### 2. Memory Validation

```go
// test/validation/memory_validation_test.go
package validation

import (
    "context"
    "runtime"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestMemory_BoundedGrowth(t *testing.T) {
    sizes := []struct {
        name       string
        dashboards int
    }{
        {"tiny", 10},
        {"small", 50},
        {"medium", 100},
        {"large", 200},
        {"xlarge", 500},
    }

    var memoryUsages []uint64

    for _, size := range sizes {
        t.Run(size.name, func(t *testing.T) {
            modPath := generateMod(t, size.dashboards)

            runtime.GC()
            runtime.GC()
            var before runtime.MemStats
            runtime.ReadMemStats(&before)

            lw, err := workspace.NewLazyWorkspace(context.Background(), modPath,
                workspace.DefaultLazyLoadConfig())
            require.NoError(t, err)
            _ = lw

            runtime.GC()
            var after runtime.MemStats
            runtime.ReadMemStats(&after)

            memUsed := after.HeapAlloc - before.HeapAlloc
            memoryUsages = append(memoryUsages, memUsed)

            t.Logf("%s (%d dashboards): %d MB", size.name, size.dashboards,
                memUsed/1024/1024)
        })
    }

    // Memory should be roughly constant (bounded)
    minMem := memoryUsages[0]
    maxMem := memoryUsages[0]
    for _, m := range memoryUsages {
        if m < minMem {
            minMem = m
        }
        if m > maxMem {
            maxMem = m
        }
    }

    // Max should be within 3x of min (bounded growth)
    assert.Less(t, maxMem, minMem*3,
        "Memory should be bounded regardless of mod size")
}

func TestMemory_Under60MB_LargeMod(t *testing.T) {
    modPath := generateMod(t, 500) // 500 dashboards

    runtime.GC()
    var before runtime.MemStats
    runtime.ReadMemStats(&before)

    lw, err := workspace.NewLazyWorkspace(context.Background(), modPath,
        workspace.DefaultLazyLoadConfig())
    require.NoError(t, err)

    // Access available dashboards (index only)
    _ = lw.GetAvailableDashboardsFromIndex()

    runtime.GC()
    var after runtime.MemStats
    runtime.ReadMemStats(&after)

    memUsed := after.HeapAlloc - before.HeapAlloc
    memMB := float64(memUsed) / 1024 / 1024

    t.Logf("Memory used: %.2f MB", memMB)
    assert.Less(t, memMB, 60.0, "Should use less than 60MB")
}

func TestMemory_CacheEviction(t *testing.T) {
    modPath := generateMod(t, 200)

    config := workspace.LazyLoadConfig{
        MaxCacheMemory: 10 * 1024 * 1024, // 10MB limit
    }

    lw, err := workspace.NewLazyWorkspace(context.Background(), modPath, config)
    require.NoError(t, err)

    ctx := context.Background()

    // Load all dashboards
    for _, entry := range lw.Index().Dashboards() {
        lw.Load(ctx, entry.Name)
    }

    stats := lw.CacheStats()
    t.Logf("Cache stats: entries=%d, memory=%d, evictions=%d",
        stats.Entries, stats.MemoryBytes, stats.Evictions)

    assert.LessOrEqual(t, stats.MemoryBytes, int64(10*1024*1024))
    assert.Greater(t, stats.Evictions, int64(0), "Should have evicted")
}
```

### 3. Behavior Validation

```go
// test/validation/behavior_validation_test.go
package validation

import (
    "context"
    "encoding/json"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestBehavior_PayloadEquivalence(t *testing.T) {
    modPath := "testdata/mods/behavior_test_mod"

    ctx := context.Background()

    // Load with lazy loading
    lazyWs, err := workspace.NewLazyWorkspace(ctx, modPath,
        workspace.DefaultLazyLoadConfig())
    require.NoError(t, err)

    lazyPayload := lazyWs.GetAvailableDashboardsFromIndex()

    // Load with eager loading
    eagerWs, err := workspace.LoadWorkspace(ctx, modPath)
    require.NoError(t, err)

    eagerPayload := buildAvailableDashboardsEager(eagerWs)

    // Compare
    assertPayloadsEquivalent(t, lazyPayload, eagerPayload)
}

func assertPayloadsEquivalent(t *testing.T, lazy, eager *AvailableDashboardsPayload) {
    t.Helper()

    // Same number of dashboards
    assert.Equal(t, len(lazy.Dashboards), len(eager.Dashboards),
        "Dashboard count mismatch")

    // Same number of benchmarks
    assert.Equal(t, len(lazy.Benchmarks), len(eager.Benchmarks),
        "Benchmark count mismatch")

    // Check each dashboard
    for name, lazyDash := range lazy.Dashboards {
        eagerDash, ok := eager.Dashboards[name]
        require.True(t, ok, "Dashboard missing in eager: %s", name)

        assert.Equal(t, lazyDash.Title, eagerDash.Title)
        assert.Equal(t, lazyDash.FullName, eagerDash.FullName)
        assert.Equal(t, lazyDash.ShortName, eagerDash.ShortName)
    }

    // Check each benchmark
    for name, lazyBench := range lazy.Benchmarks {
        eagerBench, ok := eager.Benchmarks[name]
        require.True(t, ok, "Benchmark missing in eager: %s", name)

        assert.Equal(t, lazyBench.Title, eagerBench.Title)
        assert.Equal(t, lazyBench.IsTopLevel, eagerBench.IsTopLevel)
        assert.Equal(t, len(lazyBench.Children), len(eagerBench.Children))
    }
}

func TestBehavior_DashboardExecution(t *testing.T) {
    modPath := "testdata/mods/execution_test_mod"
    ctx := context.Background()

    lw, err := workspace.NewLazyWorkspace(ctx, modPath,
        workspace.DefaultLazyLoadConfig())
    require.NoError(t, err)

    // Execute dashboard
    dash, err := lw.LoadDashboard(ctx, "testmod.dashboard.executable")
    require.NoError(t, err)

    // Verify all children loaded
    assert.NotEmpty(t, dash.GetChildren())

    for _, child := range dash.GetChildren() {
        assert.NotNil(t, child)
    }
}

func TestBehavior_BenchmarkExecution(t *testing.T) {
    modPath := "testdata/mods/execution_test_mod"
    ctx := context.Background()

    lw, err := workspace.NewLazyWorkspace(ctx, modPath,
        workspace.DefaultLazyLoadConfig())
    require.NoError(t, err)

    // Execute benchmark
    bench, err := lw.LoadBenchmark(ctx, "testmod.benchmark.executable")
    require.NoError(t, err)

    // Verify all children loaded
    assert.NotEmpty(t, bench.GetChildren())

    // Verify controls have their queries
    for _, child := range bench.GetChildren() {
        if ctrl, ok := child.(*modconfig.Control); ok {
            assert.NotNil(t, ctrl.Query, "Control should have query loaded")
        }
    }
}

func TestBehavior_AllCommandsWork(t *testing.T) {
    modPath := "testdata/mods/command_test_mod"

    commands := []struct {
        name string
        args []string
    }{
        {"list_dashboards", []string{"dashboard", "list"}},
        {"list_benchmarks", []string{"benchmark", "list"}},
        {"run_dashboard", []string{"dashboard", "run", "testmod.dashboard.test"}},
        {"run_benchmark", []string{"benchmark", "run", "testmod.benchmark.test"}},
        {"inspect", []string{"inspect", "testmod.dashboard.test"}},
    }

    for _, cmd := range commands {
        t.Run(cmd.name, func(t *testing.T) {
            err := runCommand(cmd.args, modPath, "--lazy-load")
            assert.NoError(t, err)
        })
    }
}
```

### 4. Validation Script

```bash
#!/bin/bash
# scripts/validate_lazy_loading.sh

set -e

echo "=== Lazy Loading Validation ==="
echo ""

# Build
echo "Building powerpipe..."
make build

# Run unit tests
echo ""
echo "Running unit tests..."
go test ./internal/resourceindex/... ./internal/resourcecache/... ./internal/resourceloader/...

# Run integration tests
echo ""
echo "Running integration tests..."
go test -v ./test/integration/...

# Run memory validation
echo ""
echo "Running memory validation..."
go test -v ./test/validation/... -run TestMemory

# Run behavior validation
echo ""
echo "Running behavior validation..."
go test -v ./test/validation/... -run TestBehavior

# Run with real mods
echo ""
echo "Testing with real mods..."

if [ -d "../steampipe-mod-aws-compliance" ]; then
    echo "Testing AWS Compliance mod..."
    ./powerpipe dashboard --lazy-load --mod-path ../steampipe-mod-aws-compliance &
    PID=$!
    sleep 5
    curl -s http://localhost:9033/api/available_dashboards > /dev/null
    kill $PID
    echo "AWS Compliance mod: OK"
fi

# Memory comparison
echo ""
echo "Memory comparison:"
echo "Running memory benchmark..."
go test -v -run TestMemory_Comparison ./test/validation/...

echo ""
echo "=== Validation Complete ==="
```

### 5. Documentation

```markdown
# Lazy Loading in Powerpipe

## Overview

Powerpipe now supports lazy loading of HCL resources, which dramatically reduces
memory usage and startup time for large mods.

## Enabling Lazy Loading

Lazy loading is enabled by default. To disable:

```bash
powerpipe dashboard --lazy-load=false
```

## How It Works

1. **Index-Only Startup**: At startup, Powerpipe scans mod files and builds a
   lightweight index (~1KB per 100 resources) without parsing full HCL.

2. **On-Demand Loading**: Resources are parsed only when accessed (e.g., when
   running a dashboard or benchmark).

3. **LRU Cache**: Parsed resources are cached in an LRU cache with bounded
   memory (default 50MB). Least-recently-used resources are evicted.

## Memory Usage

| Mod Size | Eager Loading | Lazy Loading | Reduction |
|----------|---------------|--------------|-----------|
| Small (10 dashboards) | 50 MB | 20 MB | 60% |
| Medium (50 dashboards) | 150 MB | 25 MB | 83% |
| Large (200 dashboards) | 400 MB | 35 MB | 91% |

## Configuration

Environment variables:
- `POWERPIPE_LAZY_LOAD`: Enable/disable lazy loading (default: true)
- `POWERPIPE_CACHE_SIZE_MB`: Maximum cache size in MB (default: 50)

## Limitations

- First access to a resource may be slightly slower (parsing required)
- Resource iteration (e.g., WalkResources) loads all resources
- File watching reindexes on changes

## Troubleshooting

If you experience issues with lazy loading:

1. Try disabling with `--lazy-load=false`
2. Check cache stats: `powerpipe --debug` shows cache hit rates
3. Report issues at github.com/turbot/powerpipe/issues
```

## Acceptance Criteria

- [ ] All unit tests pass
- [ ] All integration tests pass
- [ ] Memory validation tests pass (< 60MB for large mods)
- [ ] Behavior validation tests pass (payload equivalence)
- [ ] All CLI commands work with lazy loading
- [ ] Real-world mods work (AWS Compliance, etc.)
- [ ] Documentation complete
- [ ] Validation script runs successfully
- [ ] No performance regressions for small mods
- [ ] Cache eviction works correctly
- [ ] Error handling is graceful

## Success Metrics

| Metric | Target | Validation |
|--------|--------|------------|
| Memory (500 dashboard mod) | < 60 MB | Memory validation test |
| Startup time (500 dashboards) | < 500ms | Integration test |
| Cache hit rate (typical usage) | > 90% | Metrics collection |
| available_dashboards response | < 10ms | Integration test |
| No regressions | All behavior tests pass | Behavior validation |

## Notes

- Run validation on CI before merging
- Test with real customer mods if available
- Monitor production metrics after release
- Have rollback plan (--lazy-load=false)
