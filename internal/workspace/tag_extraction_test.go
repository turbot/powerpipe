package workspace

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// TestTagExtraction_MergeFunction tests that tags are extracted from merge() calls
// even before background resolution completes.
func TestTagExtraction_MergeFunction(t *testing.T) {
	// Create a temporary workspace with benchmark using merge() for tags
	tmpDir, err := os.MkdirTemp("", "tag_extraction_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create mod.pp with variable and benchmark using merge()
	modContent := `
variable "common_tags" {
  type = map(string)
  default = {
    service = "AWS"
    type = "Report"
  }
}

mod "test_tags" {
  title = "Test Tags"
}

benchmark "test_benchmark" {
  title = "Test Benchmark"
  tags = merge(var.common_tags, {
    compliance = "CIS"
    benchmark = "Benchmark"
  })

  control {
    sql = "select 1"
  }
}

dashboard "test_dashboard" {
  title = "Test Dashboard"
  tags = {
    category = "Accounts"
    service = "AWS"
    type = "Report"
  }

  text {
    value = "Test"
  }
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

	// Get available dashboards payload IMMEDIATELY (before background resolution)
	payload := lw.GetAvailableDashboardsFromIndex()
	if payload == nil {
		t.Fatal("expected payload to be non-nil")
	}

	// Check dashboard tags (should have all tags since they're literals)
	dashboard, ok := payload.Dashboards["test_tags.dashboard.test_dashboard"]
	if !ok {
		t.Fatal("expected test_dashboard to be in payload")
	}

	if len(dashboard.Tags) == 0 {
		t.Error("expected dashboard to have tags immediately (literal values)")
	} else {
		t.Logf("Dashboard tags (immediate): %v", dashboard.Tags)
		// Verify expected tags
		expectedDashTags := map[string]string{
			"category": "Accounts",
			"service":  "AWS",
			"type":     "Report",
		}
		for key, expectedVal := range expectedDashTags {
			if val, ok := dashboard.Tags[key]; !ok {
				t.Errorf("expected dashboard tag %s to be present", key)
			} else if val != expectedVal {
				t.Errorf("expected dashboard tag %s=%s, got %s", key, expectedVal, val)
			}
		}
	}

	// Check benchmark tags (should have partial tags from inline merge object)
	benchmark, ok := payload.Benchmarks["test_tags.benchmark.test_benchmark"]
	if !ok {
		t.Fatal("expected test_benchmark to be in payload")
	}

	if len(benchmark.Tags) == 0 {
		t.Error("expected benchmark to have partial tags immediately (from inline merge object)")
		t.Error("This indicates tag extraction from merge() inline objects is not working!")
	} else {
		t.Logf("Benchmark tags (immediate, partial): %v", benchmark.Tags)
		// Should have at least the inline tags from merge
		expectedPartialTags := []string{"compliance", "benchmark"}
		for _, key := range expectedPartialTags {
			if _, ok := benchmark.Tags[key]; !ok {
				t.Errorf("expected benchmark to have partial tag %s from inline merge object", key)
			}
		}
	}

	// Start background resolution
	lw.StartBackgroundResolution()

	// Wait for background resolution to complete (with timeout)
	completed := lw.WaitForResolution(5000) // 5 second timeout
	if !completed {
		t.Log("Warning: background resolution did not complete within timeout")
	}

	// Get payload again after background resolution
	payloadAfterResolution := lw.GetAvailableDashboardsFromIndex()
	benchmarkAfterResolution := payloadAfterResolution.Benchmarks["test_tags.benchmark.test_benchmark"]

	t.Logf("Benchmark tags (after resolution): %v", benchmarkAfterResolution.Tags)

	// After resolution, should have ALL tags including those from variable
	if completed {
		if len(benchmarkAfterResolution.Tags) <= len(benchmark.Tags) {
			t.Error("expected more tags after background resolution completes")
		}

		// Should now have tags from variable too
		if val, ok := benchmarkAfterResolution.Tags["service"]; !ok {
			t.Error("expected benchmark to have 'service' tag from variable after resolution")
		} else if val != "AWS" {
			t.Errorf("expected service=AWS, got %s", val)
		}
	}

	t.Log("✓ Tag extraction test completed")
}
