# Task 2: Memory Benchmarking Infrastructure

## Objective

Create a robust memory benchmarking infrastructure to measure and track memory usage throughout the lazy loading implementation. This provides evidence that changes are achieving the goal.

## Context

- Current baseline: 414 MB for large mod (200 dashboards, 400 queries, 500 controls)
- Target: < 60 MB with bounded growth regardless of mod size
- Need to measure both peak memory and steady-state memory
- Need to track memory at different points (load, idle, during execution)

## Dependencies

### Prerequisites
- None (this is a foundation task)

### Files to Create
- `internal/memprofile/profiler.go` - Memory profiling utilities
- `internal/memprofile/reporter.go` - Memory report generation
- `internal/workspace/workspace_memory_test.go` - Memory benchmarks
- `scripts/memory_benchmark.sh` - Automated memory benchmark script
- `benchmark_results/memory/` - Directory for results

### Files to Modify
- `internal/workspace/load_workspace_benchmark_test.go` - Add memory tracking

## Implementation Details

### 1. Memory Profiler Utilities

```go
// internal/memprofile/profiler.go
package memprofile

import (
    "fmt"
    "runtime"
    "time"
)

// MemorySnapshot captures memory state at a point in time
type MemorySnapshot struct {
    Timestamp     time.Time
    Label         string
    HeapAlloc     uint64  // Bytes allocated on heap
    HeapInuse     uint64  // Bytes in use by heap
    HeapObjects   uint64  // Number of allocated objects
    TotalAlloc    uint64  // Cumulative bytes allocated
    NumGC         uint32  // Number of GC cycles
    GCPauseTotal  uint64  // Total GC pause time (ns)
}

// TakeSnapshot captures current memory state
func TakeSnapshot(label string) MemorySnapshot {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    return MemorySnapshot{
        Timestamp:    time.Now(),
        Label:        label,
        HeapAlloc:    m.HeapAlloc,
        HeapInuse:    m.HeapInuse,
        HeapObjects:  m.HeapObjects,
        TotalAlloc:   m.TotalAlloc,
        NumGC:        m.NumGC,
        GCPauseTotal: m.PauseTotalNs,
    }
}

// ForceGC runs garbage collection and returns snapshot
func ForceGCAndSnapshot(label string) MemorySnapshot {
    runtime.GC()
    runtime.GC() // Run twice to ensure finalization
    return TakeSnapshot(label)
}

// MemoryTracker tracks memory over time
type MemoryTracker struct {
    snapshots []MemorySnapshot
    start     time.Time
}

func NewMemoryTracker() *MemoryTracker {
    return &MemoryTracker{
        start: time.Now(),
    }
}

func (t *MemoryTracker) Snapshot(label string) {
    t.snapshots = append(t.snapshots, TakeSnapshot(label))
}

func (t *MemoryTracker) SnapshotAfterGC(label string) {
    t.snapshots = append(t.snapshots, ForceGCAndSnapshot(label))
}

func (t *MemoryTracker) Report() *MemoryReport {
    return &MemoryReport{
        Duration:  time.Since(t.start),
        Snapshots: t.snapshots,
    }
}

// MemoryReport summarizes memory usage
type MemoryReport struct {
    Duration  time.Duration
    Snapshots []MemorySnapshot
}

func (r *MemoryReport) PeakHeapAlloc() uint64 {
    var peak uint64
    for _, s := range r.Snapshots {
        if s.HeapAlloc > peak {
            peak = s.HeapAlloc
        }
    }
    return peak
}

func (r *MemoryReport) FinalHeapAlloc() uint64 {
    if len(r.Snapshots) == 0 {
        return 0
    }
    return r.Snapshots[len(r.Snapshots)-1].HeapAlloc
}

func (r *MemoryReport) String() string {
    var b strings.Builder
    b.WriteString(fmt.Sprintf("Memory Report (duration: %v)\n", r.Duration))
    b.WriteString(fmt.Sprintf("Peak Heap: %s\n", formatBytes(r.PeakHeapAlloc())))
    b.WriteString(fmt.Sprintf("Final Heap: %s\n", formatBytes(r.FinalHeapAlloc())))
    b.WriteString("\nSnapshots:\n")
    for _, s := range r.Snapshots {
        b.WriteString(fmt.Sprintf("  %s: heap=%s objects=%d\n",
            s.Label, formatBytes(s.HeapAlloc), s.HeapObjects))
    }
    return b.String()
}

func formatBytes(b uint64) string {
    const unit = 1024
    if b < unit {
        return fmt.Sprintf("%d B", b)
    }
    div, exp := uint64(unit), 0
    for n := b / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    return fmt.Sprintf("%.2f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
```

### 2. Memory Benchmarks

```go
// internal/workspace/workspace_memory_test.go
package workspace

import (
    "runtime"
    "testing"

    "github.com/turbot/powerpipe/internal/memprofile"
)

func BenchmarkMemory_SmallMod(b *testing.B) {
    benchmarkModMemory(b, "testdata/mods/generated/small", 10)
}

func BenchmarkMemory_MediumMod(b *testing.B) {
    benchmarkModMemory(b, "testdata/mods/generated/medium", 50)
}

func BenchmarkMemory_LargeMod(b *testing.B) {
    benchmarkModMemory(b, "testdata/mods/generated/large", 200)
}

func BenchmarkMemory_XLargeMod(b *testing.B) {
    benchmarkModMemory(b, "testdata/mods/generated/xlarge", 500)
}

func benchmarkModMemory(b *testing.B, modPath string, expectedDashboards int) {
    b.Helper()

    // Ensure mod exists
    ensureTestMod(b, modPath)

    // Force GC before starting
    runtime.GC()
    runtime.GC()

    var initialMem runtime.MemStats
    runtime.ReadMemStats(&initialMem)

    b.ResetTimer()
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        ws, err := LoadWorkspace(context.Background(), modPath)
        if err != nil {
            b.Fatal(err)
        }

        // Verify load worked
        resources := ws.GetModResources()
        if len(resources.Dashboards) < expectedDashboards {
            b.Fatalf("Expected at least %d dashboards, got %d",
                expectedDashboards, len(resources.Dashboards))
        }

        // Don't cleanup between iterations to measure cumulative
    }

    b.StopTimer()

    // Force GC and measure final memory
    runtime.GC()
    runtime.GC()

    var finalMem runtime.MemStats
    runtime.ReadMemStats(&finalMem)

    heapGrowth := finalMem.HeapAlloc - initialMem.HeapAlloc
    b.ReportMetric(float64(heapGrowth)/float64(b.N), "heap-bytes/op")
    b.ReportMetric(float64(finalMem.HeapObjects-initialMem.HeapObjects)/float64(b.N), "heap-objects/op")
}

// TestMemoryScaling verifies memory grows linearly (or sublinearly) with mod size
func TestMemoryScaling(t *testing.T) {
    sizes := []struct {
        name       string
        path       string
        dashboards int
    }{
        {"small", "testdata/mods/generated/small", 10},
        {"medium", "testdata/mods/generated/medium", 50},
        {"large", "testdata/mods/generated/large", 200},
    }

    results := make(map[string]uint64)

    for _, size := range sizes {
        t.Run(size.name, func(t *testing.T) {
            ensureTestMod(t, size.path)

            runtime.GC()
            runtime.GC()

            var before runtime.MemStats
            runtime.ReadMemStats(&before)

            ws, err := LoadWorkspace(context.Background(), size.path)
            require.NoError(t, err)

            runtime.GC()

            var after runtime.MemStats
            runtime.ReadMemStats(&after)

            memUsed := after.HeapAlloc - before.HeapAlloc
            results[size.name] = memUsed

            t.Logf("%s: %s (%d dashboards)", size.name,
                formatBytes(memUsed), len(ws.GetModResources().Dashboards))
        })
    }

    // Log scaling factor
    if results["small"] > 0 && results["large"] > 0 {
        scalingFactor := float64(results["large"]) / float64(results["small"])
        dashboardScaling := float64(200) / float64(10) // 20x more dashboards
        t.Logf("Scaling factor: %.2fx memory for %.2fx dashboards",
            scalingFactor, dashboardScaling)

        // Memory should scale roughly linearly
        // Allow some overhead, but shouldn't be worse than O(n^2)
        if scalingFactor > dashboardScaling*1.5 {
            t.Errorf("Memory scaling worse than expected: %.2fx for %.2fx size increase",
                scalingFactor, dashboardScaling)
        }
    }
}

// TestMemoryBounded verifies memory is bounded after lazy loading
// This test will fail before lazy loading, pass after
func TestMemoryBounded(t *testing.T) {
    t.Skip("Enable after lazy loading implementation")

    // Load increasingly large mods
    sizes := []string{"small", "medium", "large", "xlarge"}
    var memoryUsages []uint64

    for _, size := range sizes {
        path := fmt.Sprintf("testdata/mods/generated/%s", size)
        ensureTestMod(t, path)

        mem := measureWorkspaceMemory(t, path)
        memoryUsages = append(memoryUsages, mem)
    }

    // With lazy loading, memory should be bounded
    maxMem := memoryUsages[0]
    for _, mem := range memoryUsages {
        if mem > maxMem {
            maxMem = mem
        }
    }

    // All should be within 2x of smallest (bounded by cache)
    for i, mem := range memoryUsages {
        if mem > maxMem*2 {
            t.Errorf("%s mod uses %s, expected bounded to ~%s",
                sizes[i], formatBytes(mem), formatBytes(maxMem))
        }
    }
}
```

### 3. Memory Benchmark Script

```bash
#!/bin/bash
# scripts/memory_benchmark.sh

set -e

OUTPUT_DIR="benchmark_results/memory/$(date +%Y%m%d_%H%M%S)"
mkdir -p "$OUTPUT_DIR"

echo "Running memory benchmarks..."

# Generate test mods if needed
go run scripts/generate_test_mods.go

# Run benchmarks with memory profiling
go test -bench=BenchmarkMemory -benchmem -benchtime=3x \
    -memprofile="$OUTPUT_DIR/mem.prof" \
    ./internal/workspace/... \
    2>&1 | tee "$OUTPUT_DIR/benchmark.txt"

# Generate memory profile analysis
go tool pprof -text -top=20 "$OUTPUT_DIR/mem.prof" > "$OUTPUT_DIR/mem_top.txt"
go tool pprof -text -cum -top=20 "$OUTPUT_DIR/mem.prof" > "$OUTPUT_DIR/mem_cum.txt"

# Run scaling test
go test -v -run=TestMemoryScaling ./internal/workspace/... \
    2>&1 | tee "$OUTPUT_DIR/scaling.txt"

# Generate summary
cat > "$OUTPUT_DIR/summary.md" << EOF
# Memory Benchmark Results

**Date**: $(date)
**Commit**: $(git rev-parse HEAD)

## Benchmark Results
\`\`\`
$(cat "$OUTPUT_DIR/benchmark.txt" | grep -E "^Benchmark|heap-bytes")
\`\`\`

## Top Memory Allocators
\`\`\`
$(head -30 "$OUTPUT_DIR/mem_top.txt")
\`\`\`

## Scaling Test
\`\`\`
$(cat "$OUTPUT_DIR/scaling.txt" | grep -E "small:|medium:|large:|Scaling")
\`\`\`
EOF

echo "Results saved to $OUTPUT_DIR"
echo ""
cat "$OUTPUT_DIR/summary.md"
```

### 4. Continuous Memory Tracking

```go
// internal/memprofile/continuous.go
package memprofile

import (
    "context"
    "time"
)

// ContinuousTracker samples memory at regular intervals
type ContinuousTracker struct {
    interval  time.Duration
    snapshots []MemorySnapshot
    cancel    context.CancelFunc
    done      chan struct{}
}

func NewContinuousTracker(interval time.Duration) *ContinuousTracker {
    return &ContinuousTracker{
        interval: interval,
        done:     make(chan struct{}),
    }
}

func (t *ContinuousTracker) Start(ctx context.Context) {
    ctx, t.cancel = context.WithCancel(ctx)

    go func() {
        defer close(t.done)
        ticker := time.NewTicker(t.interval)
        defer ticker.Stop()

        for {
            select {
            case <-ctx.Done():
                return
            case <-ticker.C:
                t.snapshots = append(t.snapshots,
                    TakeSnapshot(time.Now().Format("15:04:05.000")))
            }
        }
    }()
}

func (t *ContinuousTracker) Stop() *MemoryReport {
    if t.cancel != nil {
        t.cancel()
        <-t.done
    }
    return &MemoryReport{Snapshots: t.snapshots}
}
```

## Acceptance Criteria

- [ ] `memprofile` package created with snapshot and tracking utilities
- [ ] Memory benchmarks for small, medium, large, xlarge mods
- [ ] Memory scaling test verifies linear or sublinear growth
- [ ] Benchmark script generates comprehensive reports
- [ ] Results include heap allocation, object count, GC stats
- [ ] Baseline results captured in `benchmark_results/memory/baseline/`
- [ ] Memory profile can identify top allocators
- [ ] Tests can run in CI (no special requirements)
- [ ] Documentation explains how to run and interpret benchmarks

## Notes

- Run benchmarks multiple times for consistency
- Consider adding memory limit tests (fail if > X MB)
- Profile should capture both peak and steady-state memory
- Compare results before/after each phase of implementation
