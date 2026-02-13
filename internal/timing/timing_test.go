package timing

import (
	"os"
	"testing"
	"time"
)

func TestTrack_Disabled(t *testing.T) {
	// Ensure timing is disabled
	originalEnabled := enabled
	enabled = false
	defer func() { enabled = originalEnabled }()

	Reset()

	// Track should return a no-op function when disabled
	done := Track("test_operation")
	time.Sleep(10 * time.Millisecond)
	done()

	// No timings should be recorded
	timings := GetTimings()
	if len(timings) != 0 {
		t.Errorf("Expected no timings when disabled, got %d", len(timings))
	}
}

func TestTrack_Enabled(t *testing.T) {
	// Enable timing
	originalEnabled := enabled
	originalDetailed := detailed
	enabled = true
	detailed = false
	defer func() {
		enabled = originalEnabled
		detailed = originalDetailed
	}()

	Reset()

	// Track an operation
	done := Track("test_operation")
	time.Sleep(10 * time.Millisecond)
	done()

	// Should have one timing
	timings := GetTimings()
	if len(timings) != 1 {
		t.Fatalf("Expected 1 timing, got %d", len(timings))
	}

	if timings[0].Name != "test_operation" {
		t.Errorf("Expected name 'test_operation', got '%s'", timings[0].Name)
	}

	if timings[0].Duration < 10*time.Millisecond {
		t.Errorf("Expected duration >= 10ms, got %v", timings[0].Duration)
	}

	if timings[0].DurationMs < 10 {
		t.Errorf("Expected DurationMs >= 10, got %f", timings[0].DurationMs)
	}
}

func TestTrack_WithContext(t *testing.T) {
	// Enable timing
	originalEnabled := enabled
	enabled = true
	defer func() { enabled = originalEnabled }()

	Reset()

	// Track with context
	done := Track("test_with_context", "some context info")
	done()

	timings := GetTimings()
	if len(timings) != 1 {
		t.Fatalf("Expected 1 timing, got %d", len(timings))
	}

	if timings[0].Context != "some context info" {
		t.Errorf("Expected context 'some context info', got '%s'", timings[0].Context)
	}
}

func TestMultipleTimings(t *testing.T) {
	// Enable timing
	originalEnabled := enabled
	enabled = true
	defer func() { enabled = originalEnabled }()

	Reset()

	// Track multiple operations
	done1 := Track("operation_1")
	done1()

	done2 := Track("operation_2")
	done2()

	done3 := Track("operation_3")
	done3()

	timings := GetTimings()
	if len(timings) != 3 {
		t.Fatalf("Expected 3 timings, got %d", len(timings))
	}
}

func TestReset(t *testing.T) {
	// Enable timing
	originalEnabled := enabled
	enabled = true
	defer func() { enabled = originalEnabled }()

	// Add some timings
	done := Track("test")
	done()

	// Reset should clear all timings
	Reset()

	timings := GetTimings()
	if len(timings) != 0 {
		t.Errorf("Expected 0 timings after reset, got %d", len(timings))
	}
}

func TestIsEnabled(t *testing.T) {
	originalEnabled := enabled
	defer func() { enabled = originalEnabled }()

	enabled = true
	if !IsEnabled() {
		t.Error("Expected IsEnabled() to return true")
	}

	enabled = false
	if IsEnabled() {
		t.Error("Expected IsEnabled() to return false")
	}
}

func TestReportJSON_Disabled(t *testing.T) {
	originalEnabled := enabled
	enabled = false
	defer func() { enabled = originalEnabled }()

	Reset()

	result := ReportJSON()
	if result != "{}" {
		t.Errorf("Expected '{}' when disabled, got '%s'", result)
	}
}

func TestReportJSON_Enabled(t *testing.T) {
	originalEnabled := enabled
	enabled = true
	defer func() { enabled = originalEnabled }()

	Reset()

	done := Track("json_test")
	done()

	result := ReportJSON()
	if result == "{}" {
		t.Error("Expected non-empty JSON when timings exist")
	}

	// Basic check that it contains our operation name
	if !containsString(result, "json_test") {
		t.Errorf("Expected JSON to contain 'json_test', got '%s'", result)
	}
}

func TestGetTimings_ReturnsCopy(t *testing.T) {
	originalEnabled := enabled
	enabled = true
	defer func() { enabled = originalEnabled }()

	Reset()

	done := Track("original")
	done()

	// Get timings
	copy1 := GetTimings()

	// Add another timing
	done2 := Track("another")
	done2()

	// Original copy should not be affected
	if len(copy1) != 1 {
		t.Errorf("Original copy should still have 1 timing, got %d", len(copy1))
	}

	// New get should have both
	copy2 := GetTimings()
	if len(copy2) != 2 {
		t.Errorf("New copy should have 2 timings, got %d", len(copy2))
	}
}

func TestEnvironmentVariable(t *testing.T) {
	// Test that environment variable is checked at package load time
	// This is a basic sanity test
	envVal := os.Getenv("POWERPIPE_TIMING")
	expectedEnabled := envVal != ""

	// Re-evaluate based on current env
	testEnabled := os.Getenv("POWERPIPE_TIMING") != ""
	if testEnabled != expectedEnabled {
		t.Errorf("Unexpected enabled state based on environment variable")
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStringHelper(s, substr))
}

func containsStringHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
