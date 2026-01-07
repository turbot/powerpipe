package workspace_test

import (
	"context"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/turbot/pipe-fittings/v2/app_specific"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/pipe-fittings/v2/parse"
	pparse "github.com/turbot/powerpipe/internal/parse"
	"github.com/turbot/powerpipe/internal/resources"
	"github.com/turbot/powerpipe/internal/workspace"
)

func init() {
	// Set up app-specific constants required for mod loading
	app_specific.AppName = "powerpipe"
	app_specific.ModDataExtensions = []string{".pp", ".sp"}
	app_specific.VariablesExtensions = []string{".ppvars", ".spvars"}
	app_specific.AutoVariablesExtensions = []string{".auto.ppvars", ".auto.spvars"}
	app_specific.DefaultVarsFileName = "powerpipe.ppvars"
	app_specific.LegacyDefaultVarsFileName = "steampipe.spvars"
	app_specific.WorkspaceIgnoreFile = ".powerpipeignore"
	app_specific.WorkspaceDataDir = ".powerpipe"

	// Set up app-specific functions required for mod loading
	modconfig.AppSpecificNewModResourcesFunc = resources.NewModResources
	parse.ModDecoderFunc = pparse.NewPowerpipeModDecoder
	parse.AppSpecificGetResourceSchemaFunc = pparse.GetResourceSchema
}

// BenchmarkModPath points to the large test mod for benchmarking.
// This can be overridden by setting BENCHMARK_MOD_PATH environment variable.
func getBenchmarkModPath() string {
	if path := os.Getenv("BENCHMARK_MOD_PATH"); path != "" {
		return path
	}
	// Default to the performance test mod
	return "/Users/nathan/src/powerpipe-performance-test"
}

// BenchmarkEagerLoad measures full eager workspace loading time.
// This loads all HCL resources at startup.
func BenchmarkEagerLoad(b *testing.B) {
	modPath := getBenchmarkModPath()

	// Check if mod path exists
	if _, err := os.Stat(modPath); os.IsNotExist(err) {
		b.Skipf("Benchmark mod path does not exist: %s", modPath)
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ws, errAndWarnings := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
		if errAndWarnings.GetError() != nil {
			b.Fatalf("Failed to load workspace: %v", errAndWarnings.GetError())
		}
		ws.Close()
	}
}

// BenchmarkLazyLoad measures lazy workspace loading time.
// This only builds the resource index at startup.
func BenchmarkLazyLoad(b *testing.B) {
	modPath := getBenchmarkModPath()

	// Check if mod path exists
	if _, err := os.Stat(modPath); os.IsNotExist(err) {
		b.Skipf("Benchmark mod path does not exist: %s", modPath)
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ws, err := workspace.LoadLazy(ctx, modPath)
		if err != nil {
			b.Fatalf("Failed to load lazy workspace: %v", err)
		}
		ws.Close()
	}
}

// BenchmarkEagerLoadMemory measures memory consumption during eager loading.
func BenchmarkEagerLoadMemory(b *testing.B) {
	modPath := getBenchmarkModPath()

	if _, err := os.Stat(modPath); os.IsNotExist(err) {
		b.Skipf("Benchmark mod path does not exist: %s", modPath)
	}

	ctx := context.Background()

	var m runtime.MemStats

	// Force GC and get baseline
	runtime.GC()
	runtime.ReadMemStats(&m)
	beforeAlloc := m.TotalAlloc
	beforeHeap := m.HeapAlloc

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ws, errAndWarnings := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
		if errAndWarnings.GetError() != nil {
			b.Fatalf("Failed to load workspace: %v", errAndWarnings.GetError())
		}

		runtime.ReadMemStats(&m)
		afterAlloc := m.TotalAlloc
		afterHeap := m.HeapAlloc

		b.ReportMetric(float64(afterAlloc-beforeAlloc)/1024/1024, "MB_total_alloc")
		b.ReportMetric(float64(afterHeap-beforeHeap)/1024/1024, "MB_heap")
		b.ReportMetric(float64(m.HeapAlloc)/1024/1024, "MB_heap_current")

		ws.Close()

		// Reset for next iteration
		runtime.GC()
		runtime.ReadMemStats(&m)
		beforeAlloc = m.TotalAlloc
		beforeHeap = m.HeapAlloc
	}
}

// BenchmarkLazyLoadMemory measures memory consumption during lazy loading.
func BenchmarkLazyLoadMemory(b *testing.B) {
	modPath := getBenchmarkModPath()

	if _, err := os.Stat(modPath); os.IsNotExist(err) {
		b.Skipf("Benchmark mod path does not exist: %s", modPath)
	}

	ctx := context.Background()

	var m runtime.MemStats

	// Force GC and get baseline
	runtime.GC()
	runtime.ReadMemStats(&m)
	beforeAlloc := m.TotalAlloc
	beforeHeap := m.HeapAlloc

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ws, err := workspace.LoadLazy(ctx, modPath)
		if err != nil {
			b.Fatalf("Failed to load lazy workspace: %v", err)
		}

		runtime.ReadMemStats(&m)
		afterAlloc := m.TotalAlloc
		afterHeap := m.HeapAlloc

		b.ReportMetric(float64(afterAlloc-beforeAlloc)/1024/1024, "MB_total_alloc")
		b.ReportMetric(float64(afterHeap-beforeHeap)/1024/1024, "MB_heap")
		b.ReportMetric(float64(m.HeapAlloc)/1024/1024, "MB_heap_current")

		ws.Close()

		// Reset for next iteration
		runtime.GC()
		runtime.ReadMemStats(&m)
		beforeAlloc = m.TotalAlloc
		beforeHeap = m.HeapAlloc
	}
}

// BenchmarkIndexBuild measures just the ResourceIndex construction time.
func BenchmarkIndexBuild(b *testing.B) {
	modPath := getBenchmarkModPath()

	if _, err := os.Stat(modPath); os.IsNotExist(err) {
		b.Skipf("Benchmark mod path does not exist: %s", modPath)
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Load lazy workspace which only builds the index
		ws, err := workspace.LoadLazy(ctx, modPath, workspace.WithLazyLoadConfig(workspace.LazyLoadConfig{
			MaxCacheMemory: 50 * 1024 * 1024,
			EnablePreload:  false,
		}))
		if err != nil {
			b.Fatalf("Failed to build index: %v", err)
		}
		ws.Close()
	}
}

// BenchmarkDashboardListFromIndex measures building the dashboard list payload from the index.
func BenchmarkDashboardListFromIndex(b *testing.B) {
	modPath := getBenchmarkModPath()

	if _, err := os.Stat(modPath); os.IsNotExist(err) {
		b.Skipf("Benchmark mod path does not exist: %s", modPath)
	}

	ctx := context.Background()

	// Load lazy workspace first
	ws, err := workspace.LoadLazy(ctx, modPath)
	if err != nil {
		b.Fatalf("Failed to load lazy workspace: %v", err)
	}
	defer ws.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		payload := ws.GetAvailableDashboardsFromIndex()
		if payload == nil {
			b.Fatal("Failed to get dashboard list")
		}
		// Report some metrics
		if i == 0 {
			b.ReportMetric(float64(len(payload.Dashboards)), "dashboards")
			b.ReportMetric(float64(len(payload.Benchmarks)), "benchmarks")
		}
	}
}

// BenchmarkDashboardLoadOnDemand measures loading a single dashboard on-demand.
func BenchmarkDashboardLoadOnDemand(b *testing.B) {
	modPath := getBenchmarkModPath()

	if _, err := os.Stat(modPath); os.IsNotExist(err) {
		b.Skipf("Benchmark mod path does not exist: %s", modPath)
	}

	ctx := context.Background()

	// Load lazy workspace first
	ws, err := workspace.LoadLazy(ctx, modPath)
	if err != nil {
		b.Fatalf("Failed to load lazy workspace: %v", err)
	}
	defer ws.Close()

	// Get a dashboard name from the index
	payload := ws.GetAvailableDashboardsFromIndex()
	if len(payload.Dashboards) == 0 {
		b.Skip("No dashboards found in mod")
	}

	// Pick first dashboard
	var dashName string
	for name := range payload.Dashboards {
		dashName = name
		break
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ws.InvalidateAll() // Clear cache for fair measurement
		dash, err := ws.LoadDashboard(ctx, dashName)
		if err != nil {
			b.Fatalf("Failed to load dashboard: %v", err)
		}
		if dash == nil {
			b.Fatal("Dashboard is nil")
		}
	}
}

// TestMeasureLoadTimes runs a simple timing test and outputs results.
// This is useful for quick manual testing.
// Skip in short mode since it requires an external mod with dependencies.
func TestMeasureLoadTimes(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	modPath := getBenchmarkModPath()

	if _, err := os.Stat(modPath); os.IsNotExist(err) {
		t.Skipf("Benchmark mod path does not exist: %s", modPath)
	}

	ctx := context.Background()
	runs := 3

	t.Log("=== Workspace Loading Performance Test ===")
	t.Logf("Mod path: %s", modPath)
	t.Log("")

	// Measure eager loading
	t.Log("--- Eager Loading ---")
	var eagerTimes []time.Duration
	for i := 0; i < runs; i++ {
		start := time.Now()
		ws, errAndWarnings := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
		elapsed := time.Since(start)
		if errAndWarnings.GetError() != nil {
			t.Fatalf("Eager load failed: %v", errAndWarnings.GetError())
		}
		ws.Close()
		eagerTimes = append(eagerTimes, elapsed)
		t.Logf("  Run %d: %v", i+1, elapsed)
	}

	// Measure lazy loading
	t.Log("")
	t.Log("--- Lazy Loading ---")
	var lazyTimes []time.Duration
	for i := 0; i < runs; i++ {
		start := time.Now()
		ws, err := workspace.LoadLazy(ctx, modPath)
		elapsed := time.Since(start)
		if err != nil {
			t.Fatalf("Lazy load failed: %v", err)
		}
		ws.Close()
		lazyTimes = append(lazyTimes, elapsed)
		t.Logf("  Run %d: %v", i+1, elapsed)
	}

	// Calculate averages
	var eagerAvg, lazyAvg time.Duration
	for _, d := range eagerTimes {
		eagerAvg += d
	}
	eagerAvg /= time.Duration(runs)

	for _, d := range lazyTimes {
		lazyAvg += d
	}
	lazyAvg /= time.Duration(runs)

	t.Log("")
	t.Log("--- Summary ---")
	t.Logf("Eager Loading Average: %v", eagerAvg)
	t.Logf("Lazy Loading Average:  %v", lazyAvg)
	if eagerAvg > 0 {
		speedup := float64(eagerAvg) / float64(lazyAvg)
		t.Logf("Speedup: %.2fx", speedup)
	}
}

// TestMeasureMemory measures memory consumption differences.
// Skip in short mode since it requires an external mod with dependencies.
func TestMeasureMemory(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory test in short mode")
	}

	modPath := getBenchmarkModPath()

	if _, err := os.Stat(modPath); os.IsNotExist(err) {
		t.Skipf("Benchmark mod path does not exist: %s", modPath)
	}

	ctx := context.Background()
	var m runtime.MemStats

	t.Log("=== Memory Consumption Test ===")
	t.Logf("Mod path: %s", modPath)
	t.Log("")

	// Measure eager loading memory
	t.Log("--- Eager Loading ---")
	runtime.GC()
	runtime.ReadMemStats(&m)
	beforeHeap := m.HeapAlloc

	ws, errAndWarnings := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	if errAndWarnings.GetError() != nil {
		t.Fatalf("Eager load failed: %v", errAndWarnings.GetError())
	}

	runtime.ReadMemStats(&m)
	eagerHeap := m.HeapAlloc - beforeHeap
	t.Logf("  Heap allocated: %.2f MB", float64(eagerHeap)/1024/1024)
	t.Logf("  Current heap:   %.2f MB", float64(m.HeapAlloc)/1024/1024)
	ws.Close()

	// Force GC
	runtime.GC()

	// Measure lazy loading memory
	t.Log("")
	t.Log("--- Lazy Loading ---")
	runtime.ReadMemStats(&m)
	beforeHeap = m.HeapAlloc

	lws, err := workspace.LoadLazy(ctx, modPath)
	if err != nil {
		t.Fatalf("Lazy load failed: %v", err)
	}

	runtime.ReadMemStats(&m)
	lazyHeap := m.HeapAlloc - beforeHeap
	t.Logf("  Heap allocated: %.2f MB", float64(lazyHeap)/1024/1024)
	t.Logf("  Current heap:   %.2f MB", float64(m.HeapAlloc)/1024/1024)
	t.Logf("  Index entries:  %d", lws.IndexStats().TotalEntries)
	lws.Close()

	// Summary
	t.Log("")
	t.Log("--- Summary ---")
	t.Logf("Eager heap: %.2f MB", float64(eagerHeap)/1024/1024)
	t.Logf("Lazy heap:  %.2f MB", float64(lazyHeap)/1024/1024)
	if eagerHeap > lazyHeap {
		reduction := float64(eagerHeap-lazyHeap) / float64(eagerHeap) * 100
		t.Logf("Memory reduction: %.1f%%", reduction)
	} else {
		t.Log("Note: Memory comparison affected by GC - use benchmark functions for accurate measurements")
	}
}
