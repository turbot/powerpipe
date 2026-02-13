package dashboardserver

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/turbot/pipe-fittings/v2/app_specific"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/pipe-fittings/v2/parse"
	pparse "github.com/turbot/powerpipe/internal/parse"
	ppresources "github.com/turbot/powerpipe/internal/resources"
	"github.com/turbot/powerpipe/internal/timing"
	"github.com/turbot/powerpipe/internal/workspace"
)

func init() {
	// Enable timing for benchmarks
	os.Setenv("POWERPIPE_TIMING", "1")

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
	modconfig.AppSpecificNewModResourcesFunc = ppresources.NewModResources
	parse.ModDecoderFunc = pparse.NewPowerpipeModDecoder
	parse.AppSpecificGetResourceSchemaFunc = pparse.GetResourceSchema
}

// BenchmarkBuildAvailableDashboardsPayload_Small benchmarks payload building for small mod
func BenchmarkBuildAvailableDashboardsPayload_Small(b *testing.B) {
	benchmarkPayload(b, "small")
}

// BenchmarkBuildAvailableDashboardsPayload_Medium benchmarks payload building for medium mod
func BenchmarkBuildAvailableDashboardsPayload_Medium(b *testing.B) {
	benchmarkPayload(b, "medium")
}

// BenchmarkBuildAvailableDashboardsPayload_Large benchmarks payload building for large mod
func BenchmarkBuildAvailableDashboardsPayload_Large(b *testing.B) {
	benchmarkPayload(b, "large")
}

func benchmarkPayload(b *testing.B, size string) {
	modPath := ensureGeneratedMod(b, size)
	ctx := context.Background()

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	if ew.GetError() != nil {
		b.Skipf("Failed to load workspace: %v", ew.GetError())
	}

	resources := w.GetPowerpipeModResources()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		timing.Reset()
		_, err := buildAvailableDashboardsPayload(resources)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()

	if timing.IsEnabled() {
		b.Log(timing.ReportJSON())
	}
}

// BenchmarkPayloadJSONMarshal_Large measures JSON marshaling time specifically
func BenchmarkPayloadJSONMarshal_Large(b *testing.B) {
	modPath := ensureGeneratedMod(b, "large")
	ctx := context.Background()

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	if ew.GetError() != nil {
		b.Skipf("Failed to load workspace: %v", ew.GetError())
	}

	resources := w.GetPowerpipeModResources()

	// Build payload once outside timing loop
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := buildAvailableDashboardsPayload(resources)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkPayloadRepeatedCalls measures repeated payload building (cache potential)
func BenchmarkPayloadRepeatedCalls(b *testing.B) {
	modPath := ensureGeneratedMod(b, "medium")
	ctx := context.Background()

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	if ew.GetError() != nil {
		b.Skipf("Failed to load workspace: %v", ew.GetError())
	}

	resources := w.GetPowerpipeModResources()

	// Call multiple times in each iteration to simulate caching benefit
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			_, err := buildAvailableDashboardsPayload(resources)
			if err != nil {
				b.Fatal(err)
			}
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
	// Go up from internal/dashboardserver to project root
	return filepath.Join(filepath.Dir(filename), "..", "..")
}
