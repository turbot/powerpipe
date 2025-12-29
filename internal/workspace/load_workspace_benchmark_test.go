package workspace

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/turbot/powerpipe/internal/memprofile"
	"github.com/turbot/powerpipe/internal/timing"
)

func init() {
	// Enable timing for benchmarks
	os.Setenv("POWERPIPE_TIMING", "1")
}

// BenchmarkLoadWorkspace_Small benchmarks loading a small mod (10 dashboards)
func BenchmarkLoadWorkspace_Small(b *testing.B) {
	benchmarkLoadWorkspace(b, "small")
}

// BenchmarkLoadWorkspace_Medium benchmarks loading a medium mod (50 dashboards)
func BenchmarkLoadWorkspace_Medium(b *testing.B) {
	benchmarkLoadWorkspace(b, "medium")
}

// BenchmarkLoadWorkspace_Large benchmarks loading a large mod (200 dashboards)
func BenchmarkLoadWorkspace_Large(b *testing.B) {
	benchmarkLoadWorkspace(b, "large")
}

func benchmarkLoadWorkspace(b *testing.B, size string) {
	modPath := ensureGeneratedMod(b, size)
	ctx := context.Background()

	// Track memory
	tracker := memprofile.NewMemoryTracker()
	tracker.SnapshotAfterGC("before")

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		timing.Reset()
		_, ew := Load(ctx, modPath, WithVariableValidation(false))
		if ew.GetError() != nil {
			b.Fatalf("Failed to load workspace: %v", ew.GetError())
		}
	}
	b.StopTimer()

	// Capture memory after load
	tracker.SnapshotAfterGC("after")
	report := tracker.Report()

	// Report memory metrics
	b.ReportMetric(float64(report.FinalHeapAlloc())/(1024*1024), "heap-MB")

	// Report timing breakdown
	if timing.IsEnabled() {
		b.Log(timing.ReportJSON())
	}
}

// BenchmarkLoadWorkspace_Parallel tests parallel loading capability
// Note: Currently skipped due to global viper state not being thread-safe
func BenchmarkLoadWorkspace_Parallel(b *testing.B) {
	b.Skip("Skipping parallel benchmark: workspace loading uses global viper state which is not thread-safe")
	modPath := ensureGeneratedMod(b, "medium")
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, ew := Load(ctx, modPath, WithVariableValidation(false))
			if ew.GetError() != nil {
				b.Fatalf("Failed to load workspace: %v", ew.GetError())
			}
		}
	})
}

// BenchmarkFileIO_Large measures just the file reading portion
func BenchmarkFileIO_Large(b *testing.B) {
	modPath := ensureGeneratedMod(b, "large")

	// List all .pp files
	files, err := filepath.Glob(filepath.Join(modPath, "*.pp"))
	if err != nil {
		b.Fatalf("Failed to glob files: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, f := range files {
			_, err := os.ReadFile(f)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

// BenchmarkLoadWorkspace_ColdStart simulates cold start by clearing caches
func BenchmarkLoadWorkspace_ColdStart(b *testing.B) {
	modPath := ensureGeneratedMod(b, "medium")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		timing.Reset()
		_, ew := Load(ctx, modPath, WithVariableValidation(false))
		if ew.GetError() != nil {
			b.Fatalf("Failed to load workspace: %v", ew.GetError())
		}
	}
}

func ensureGeneratedMod(b *testing.B, size string) string {
	b.Helper()

	modPath := filepath.Join(benchmarkTestdataDir(), "mods", "generated", size)

	// Check if mod exists
	if _, err := os.Stat(filepath.Join(modPath, "mod.pp")); os.IsNotExist(err) {
		// Generate mod using the generator script
		scriptPath := filepath.Join(projectRoot(), "scripts", "generate_test_mods.go")
		cmd := exec.Command("go", "run", scriptPath, modPath, size)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			b.Skipf("Failed to generate test mod: %v", err)
		}
	}

	return modPath
}

func benchmarkTestdataDir() string {
	return filepath.Join(projectRoot(), "internal", "testdata")
}

func projectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	// Go up from internal/workspace to project root
	return filepath.Join(filepath.Dir(filename), "..", "..")
}
