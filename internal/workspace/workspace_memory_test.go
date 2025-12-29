package workspace

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/turbot/powerpipe/internal/memprofile"
)

// BenchmarkMemory_Small benchmarks memory usage for loading a small mod (10 dashboards)
func BenchmarkMemory_Small(b *testing.B) {
	benchmarkModMemory(b, "small", 10)
}

// BenchmarkMemory_Medium benchmarks memory usage for loading a medium mod (50 dashboards)
func BenchmarkMemory_Medium(b *testing.B) {
	benchmarkModMemory(b, "medium", 50)
}

// BenchmarkMemory_Large benchmarks memory usage for loading a large mod (200 dashboards)
func BenchmarkMemory_Large(b *testing.B) {
	benchmarkModMemory(b, "large", 200)
}

// BenchmarkMemory_XLarge benchmarks memory usage for loading an xlarge mod (500 dashboards)
func BenchmarkMemory_XLarge(b *testing.B) {
	benchmarkModMemory(b, "xlarge", 500)
}

func benchmarkModMemory(b *testing.B, size string, expectedDashboards int) {
	b.Helper()

	modPath := ensureGeneratedModMem(b, size)
	ctx := context.Background()

	// Force GC before starting
	runtime.GC()
	runtime.GC()

	var initialMem runtime.MemStats
	runtime.ReadMemStats(&initialMem)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ws, ew := Load(ctx, modPath, WithVariableValidation(false))
		if ew.GetError() != nil {
			b.Fatal(ew.GetError())
		}

		// Verify load worked
		res := ws.GetPowerpipeModResources()
		if len(res.Dashboards) < expectedDashboards {
			b.Fatalf("Expected at least %d dashboards, got %d",
				expectedDashboards, len(res.Dashboards))
		}

		// Don't cleanup between iterations to measure cumulative memory
	}

	b.StopTimer()

	// Force GC and measure final memory
	runtime.GC()
	runtime.GC()

	var finalMem runtime.MemStats
	runtime.ReadMemStats(&finalMem)

	heapGrowth := finalMem.HeapAlloc - initialMem.HeapAlloc
	objectGrowth := finalMem.HeapObjects - initialMem.HeapObjects
	b.ReportMetric(float64(heapGrowth)/float64(b.N), "heap-bytes/op")
	b.ReportMetric(float64(objectGrowth)/float64(b.N), "heap-objects/op")
	b.ReportMetric(float64(finalMem.HeapAlloc)/(1024*1024), "final-heap-MB")
}

// TestMemoryScaling verifies memory grows linearly (or sublinearly) with mod size
func TestMemoryScaling(t *testing.T) {
	sizes := []struct {
		name       string
		size       string
		dashboards int
	}{
		{"small", "small", 10},
		{"medium", "medium", 50},
		{"large", "large", 200},
	}

	results := make(map[string]uint64)
	ctx := context.Background()

	for _, size := range sizes {
		t.Run(size.name, func(t *testing.T) {
			modPath := ensureGeneratedModMem(t, size.size)

			runtime.GC()
			runtime.GC()

			var before runtime.MemStats
			runtime.ReadMemStats(&before)

			ws, ew := Load(ctx, modPath, WithVariableValidation(false))
			require.NoError(t, ew.GetError())

			runtime.GC()

			var after runtime.MemStats
			runtime.ReadMemStats(&after)

			memUsed := after.HeapAlloc - before.HeapAlloc
			results[size.name] = memUsed

			res := ws.GetPowerpipeModResources()
			t.Logf("%s: %s (%d dashboards, %d queries, %d controls)",
				size.name,
				memprofile.FormatBytes(memUsed),
				len(res.Dashboards),
				len(res.Queries),
				len(res.Controls))
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
	ctx := context.Background()

	for _, size := range sizes {
		modPath := ensureGeneratedModMem(t, size)
		mem := measureWorkspaceMemory(t, ctx, modPath)
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
				sizes[i], memprofile.FormatBytes(mem), memprofile.FormatBytes(maxMem))
		}
	}
}

// TestMemoryProfile runs a detailed memory profile of workspace loading
func TestMemoryProfile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory profile in short mode")
	}

	modPath := ensureGeneratedModMem(t, "large")
	ctx := context.Background()

	tracker := memprofile.NewMemoryTracker()

	// Initial snapshot
	tracker.SnapshotAfterGC("initial")

	// Load workspace
	ws, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	tracker.Snapshot("after-load")

	// Access resources
	res := ws.GetPowerpipeModResources()
	_ = len(res.Dashboards)

	tracker.Snapshot("after-access")

	// Force GC
	tracker.SnapshotAfterGC("after-gc")

	// Generate report
	report := tracker.Report()
	t.Log("\n" + report.String())

	// Report key metrics
	t.Logf("Peak heap: %s", memprofile.FormatBytes(report.PeakHeapAlloc()))
	t.Logf("Final heap: %s", memprofile.FormatBytes(report.FinalHeapAlloc()))
	t.Logf("Peak objects: %d", report.PeakHeapObjects())
}

// TestMemoryPerResourceType measures memory usage per resource type
func TestMemoryPerResourceType(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping detailed memory test in short mode")
	}

	modPath := ensureGeneratedModMem(t, "large")
	ctx := context.Background()

	runtime.GC()
	runtime.GC()

	var before runtime.MemStats
	runtime.ReadMemStats(&before)

	ws, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	var after runtime.MemStats
	runtime.ReadMemStats(&after)

	res := ws.GetPowerpipeModResources()

	totalMem := after.HeapAlloc - before.HeapAlloc
	totalObjects := len(res.Dashboards) + len(res.Queries) + len(res.Controls)

	if totalObjects > 0 {
		avgPerObject := totalMem / uint64(totalObjects)
		t.Logf("Total memory: %s", memprofile.FormatBytes(totalMem))
		t.Logf("Total objects: %d", totalObjects)
		t.Logf("Avg memory per object: %s", memprofile.FormatBytes(avgPerObject))
		t.Logf("  Dashboards: %d", len(res.Dashboards))
		t.Logf("  Queries: %d", len(res.Queries))
		t.Logf("  Controls: %d", len(res.Controls))
		t.Logf("  Benchmarks: %d", len(res.ControlBenchmarks))
	}
}

// TestMemoryLimit fails if memory exceeds a threshold
func TestMemoryLimit(t *testing.T) {
	const maxMemoryMB = 500 // Maximum allowed memory in MB

	modPath := ensureGeneratedModMem(t, "large")
	ctx := context.Background()

	mem := measureWorkspaceMemory(t, ctx, modPath)
	memMB := float64(mem) / (1024 * 1024)

	t.Logf("Memory used: %.2f MB (limit: %d MB)", memMB, maxMemoryMB)

	if memMB > maxMemoryMB {
		t.Errorf("Memory usage %.2f MB exceeds limit of %d MB", memMB, maxMemoryMB)
	}
}

// measureWorkspaceMemory measures the memory used to load a workspace
func measureWorkspaceMemory(t testing.TB, ctx context.Context, modPath string) uint64 {
	t.Helper()

	runtime.GC()
	runtime.GC()

	var before runtime.MemStats
	runtime.ReadMemStats(&before)

	_, ew := Load(ctx, modPath, WithVariableValidation(false))
	if ew.GetError() != nil {
		t.Fatal(ew.GetError())
	}

	runtime.GC()

	var after runtime.MemStats
	runtime.ReadMemStats(&after)

	return after.HeapAlloc - before.HeapAlloc
}

// ensureGeneratedModMem ensures a generated test mod exists
// Named differently to avoid conflict with the one in load_workspace_benchmark_test.go
func ensureGeneratedModMem(tb testing.TB, size string) string {
	tb.Helper()

	modPath := filepath.Join(benchmarkTestdataDir(), "mods", "generated", size)

	// Check if mod exists
	if _, err := os.Stat(filepath.Join(modPath, "mod.pp")); os.IsNotExist(err) {
		// Generate mod using the generator script
		scriptPath := filepath.Join(projectRoot(), "scripts", "generate_test_mods.go")
		cmd := exec.Command("go", "run", scriptPath, modPath, size)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			tb.Skipf("Failed to generate test mod: %v", err)
		}
	}

	return modPath
}
