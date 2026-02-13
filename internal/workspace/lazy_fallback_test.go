package workspace

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// TestLazyLoading_EmptyIndexStillWorks tests that lazy loading works correctly
// even when the index has no dashboards or benchmarks. This is important for
// Pipes where we don't want to fall back to eager loading.
func TestLazyLoading_EmptyIndexStillWorks(t *testing.T) {
	// Create a temporary directory with just a mod.pp (no dashboard/benchmark resources)
	tmpDir, err := os.MkdirTemp("", "lazy_empty_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create minimal mod.pp
	modContent := `mod "test_empty" {
  title = "Test Empty Mod"
}`
	modPath := filepath.Join(tmpDir, "mod.pp")
	if err := os.WriteFile(modPath, []byte(modContent), 0600); err != nil {
		t.Fatalf("failed to write mod.pp: %v", err)
	}

	// Load lazy workspace
	ctx := context.Background()
	lw, err := NewLazyWorkspace(ctx, tmpDir, DefaultLazyLoadConfig())
	if err != nil {
		t.Fatalf("failed to create lazy workspace: %v", err)
	}
	defer lw.Close()

	// Verify index is empty (no dashboards or benchmarks)
	index := lw.GetIndex()
	dashboards := index.Dashboards()
	benchmarks := index.Benchmarks()

	if len(dashboards) != 0 {
		t.Errorf("expected 0 dashboards, got %d", len(dashboards))
	}
	if len(benchmarks) != 0 {
		t.Errorf("expected 0 benchmarks, got %d", len(benchmarks))
	}

	// NEW BEHAVIOR: Lazy workspace should work even with empty index
	// GetModResources should return valid mod resources
	modResources := lw.GetModResources()
	if modResources == nil {
		t.Error("expected GetModResources to return non-nil even with empty index")
	}

	t.Logf("Empty index test passed: lazy workspace works without fallback (dashboards=%d, benchmarks=%d)",
		len(dashboards), len(benchmarks))
}

// TestLazyLoading_GetModResourcesPopulatesDependencyMods tests that GetModResources
// returns dependency mods from the index without requiring eager loading.
func TestLazyLoading_GetModResourcesPopulatesDependencyMods(t *testing.T) {
	// Create a temporary workspace with main mod and a dependency
	tmpDir, err := os.MkdirTemp("", "lazy_mods_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create main mod.pp
	modContent := `mod "main_mod" {
  title = "Main Mod"
}

dashboard "test_dashboard" {
  title = "Test Dashboard"
  text {
    value = "Test"
  }
}`
	modPath := filepath.Join(tmpDir, "mod.pp")
	if err := os.WriteFile(modPath, []byte(modContent), 0600); err != nil {
		t.Fatalf("failed to write mod.pp: %v", err)
	}

	// Create a fake dependency mod structure
	depModDir := filepath.Join(tmpDir, ".powerpipe", "mods", "github.com", "turbot", "steampipe-mod-aws-insights@v1.0.0")
	if err := os.MkdirAll(depModDir, 0755); err != nil {
		t.Fatalf("failed to create dep mod dir: %v", err)
	}

	depModContent := `mod "aws_insights" {
  title = "AWS Insights"
}

dashboard "aws_dashboard" {
  title = "AWS Dashboard"
  text {
    value = "AWS Test"
  }
}`
	depModPath := filepath.Join(depModDir, "mod.pp")
	if err := os.WriteFile(depModPath, []byte(depModContent), 0600); err != nil {
		t.Fatalf("failed to write dep mod.pp: %v", err)
	}

	// Load lazy workspace
	ctx := context.Background()
	lw, err := NewLazyWorkspace(ctx, tmpDir, DefaultLazyLoadConfig())
	if err != nil {
		t.Fatalf("failed to create lazy workspace: %v", err)
	}
	defer lw.Close()

	// Get mod resources
	modResources := lw.GetModResources()
	if modResources == nil {
		t.Fatal("expected GetModResources to return non-nil")
	}

	// Verify main mod is present
	mainModFullName := lw.Mod.GetFullName()
	if _, ok := modResources.GetMods()[mainModFullName]; !ok {
		t.Errorf("main mod %s not found in mod resources", mainModFullName)
	}

	// Verify dependency mod is present
	depModFullName := "mod.aws_insights"
	depMod, ok := modResources.GetMods()[depModFullName]
	if !ok {
		t.Errorf("dependency mod %s not found in mod resources", depModFullName)
	} else {
		// Verify dependency mod has basic metadata
		if depMod.FullName != depModFullName {
			t.Errorf("expected dependency mod full name %s, got %s", depModFullName, depMod.FullName)
		}
		if depMod.ShortName != "aws_insights" {
			t.Errorf("expected dependency mod short name aws_insights, got %s", depMod.ShortName)
		}
	}

	t.Logf("GetModResources test passed: found %d mods (including dependencies)", len(modResources.GetMods()))
}

// TestLazyLoading_WithResources tests that lazy loading works correctly
// when resources ARE present.
func TestLazyLoading_WithResources(t *testing.T) {
	// Create a temporary directory with dashboard and benchmark resources
	tmpDir, err := os.MkdirTemp("", "lazy_with_resources_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create mod.pp with dashboard and benchmark
	modContent := `mod "test_with_resources" {
  title = "Test Mod with Resources"
}

dashboard "test_dashboard" {
  title = "Test Dashboard"

  text {
    value = "Hello from test dashboard"
  }
}

benchmark "test_benchmark" {
  title = "Test Benchmark"
  children = []
}
`
	modPath := filepath.Join(tmpDir, "mod.pp")
	if err := os.WriteFile(modPath, []byte(modContent), 0600); err != nil {
		t.Fatalf("failed to write mod.pp: %v", err)
	}

	// Load lazy workspace
	ctx := context.Background()
	lw, err := NewLazyWorkspace(ctx, tmpDir, DefaultLazyLoadConfig())
	if err != nil {
		t.Fatalf("failed to create lazy workspace: %v", err)
	}
	defer lw.Close()

	// Verify resources are indexed
	index := lw.GetIndex()
	dashboards := index.Dashboards()
	benchmarks := index.Benchmarks()

	if len(dashboards) != 1 {
		t.Errorf("expected 1 dashboard, got %d", len(dashboards))
	}
	if len(benchmarks) != 1 {
		t.Errorf("expected 1 benchmark, got %d", len(benchmarks))
	}

	// Verify GetModResources works
	modResources := lw.GetModResources()
	if modResources == nil {
		t.Error("expected GetModResources to return non-nil")
	}

	t.Logf("With resources test passed: dashboards=%d, benchmarks=%d",
		len(dashboards), len(benchmarks))
}
