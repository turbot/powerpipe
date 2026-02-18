package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanModInfo_OpengraphTitle(t *testing.T) {
	// Create a temporary directory with a test mod.pp file
	// that has both a top-level title and an opengraph title
	tmpDir := t.TempDir()
	modContent := `mod "test_mod" {
  title         = "Test Mod"
  description   = "A test mod"

  opengraph {
    title       = "Powerpipe Mod for Test Mod"
    description = "A test mod for Powerpipe"
  }
}
`
	modPath := filepath.Join(tmpDir, "mod.pp")
	if err := os.WriteFile(modPath, []byte(modContent), 0644); err != nil {
		t.Fatalf("Failed to write test mod.pp: %v", err)
	}

	// Call scanModInfo
	modName, modFullName, modTitle, err := scanModInfo(tmpDir)
	if err != nil {
		t.Fatalf("scanModInfo failed: %v", err)
	}

	// Verify results
	if modName != "test_mod" {
		t.Errorf("Expected modName 'test_mod', got '%s'", modName)
	}
	if modFullName != "mod.test_mod" {
		t.Errorf("Expected modFullName 'mod.test_mod', got '%s'", modFullName)
	}

	// This is the critical check: should capture the top-level title, not opengraph title
	if modTitle != "Test Mod" {
		t.Errorf("Expected modTitle 'Test Mod', got '%s' - scanModInfo is capturing title from wrong block (likely opengraph)", modTitle)
	}

	// Explicitly verify it did NOT capture the opengraph title
	if modTitle == "Powerpipe Mod for Test Mod" {
		t.Error("scanModInfo incorrectly captured title from opengraph block instead of top-level mod block")
	}

	t.Logf("✓ Successfully captured correct title: '%s'", modTitle)
}

func TestScanModInfo_RealAwsComplianceMod(t *testing.T) {
	// Test with the actual AWS Compliance mod if it exists
	modPath := "/Users/pskrbasu/pskr/.powerpipe/mods/github.com/turbot/steampipe-mod-aws-compliance@v1.13.0"

	// Skip if the mod doesn't exist (won't fail in CI)
	if _, err := os.Stat(filepath.Join(modPath, "mod.pp")); os.IsNotExist(err) {
		t.Skip("AWS Compliance mod not found, skipping")
	}

	modName, modFullName, modTitle, err := scanModInfo(modPath)
	if err != nil {
		t.Fatalf("scanModInfo failed: %v", err)
	}

	t.Logf("ModName: %s", modName)
	t.Logf("ModFullName: %s", modFullName)
	t.Logf("ModTitle: %s", modTitle)

	// Verify it extracted the correct title
	if modTitle != "AWS Compliance" {
		t.Errorf("Expected modTitle 'AWS Compliance', got '%s'", modTitle)
	}

	// Ensure it did NOT capture the opengraph title
	if modTitle == "Powerpipe Mod for AWS Compliance" {
		t.Error("scanModInfo incorrectly captured opengraph title instead of top-level title")
	}

	t.Logf("✓ Correctly captured AWS Compliance mod title")
}
