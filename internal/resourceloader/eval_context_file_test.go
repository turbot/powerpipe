package resourceloader

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// TestEvalContext_DependencyModWithFileFunction tests the scenario where:
// 1. Start with empty workspace (just mod.pp)
// 2. Install dependency mod (aws-compliance) that has locals with file() calls
// 3. Build eval context - should NOT crash even if file paths in locals are relative
//
// This reproduces the issue where powerpipe crashes with:
// "Invalid value for "path" parameter: no file exists at ./foundational_security/docs/..."
func TestEvalContext_DependencyModWithFileFunction(t *testing.T) {
	// Create temporary workspace
	tmpDir := t.TempDir()

	// Create main mod.pp (empty workspace)
	modContent := `mod "local" {
  title = "test_workspace"
}`
	if err := os.WriteFile(filepath.Join(tmpDir, "mod.pp"), []byte(modContent), 0600); err != nil {
		t.Fatal(err)
	}

	// Create .powerpipe/mods directory structure
	modsDir := filepath.Join(tmpDir, ".powerpipe", "mods")
	depModDir := filepath.Join(modsDir, "github.com", "test-org", "test-mod@v1.0.0")
	if err := os.MkdirAll(depModDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create dependency mod with mod.pp
	depModContent := `mod "test-mod" {
  title = "Test Mod"
}`
	if err := os.WriteFile(filepath.Join(depModDir, "mod.pp"), []byte(depModContent), 0600); err != nil {
		t.Fatal(err)
	}

	// Create a locals file with file() function referencing a non-existent file
	// This mimics aws-compliance mod structure where locals reference docs files
	localsContent := `locals {
  # This file() call references a relative path that doesn't exist
  # Should not crash during eval context building
  test_doc = file("./docs/test_document.md")

  # Another local that might reference the one above
  test_metadata = {
    documentation = local.test_doc
    title = "Test"
  }
}`
	if err := os.WriteFile(filepath.Join(depModDir, "test_resource.pp"), []byte(localsContent), 0600); err != nil {
		t.Fatal(err)
	}

	// Create another locals file that references locals from the first file
	// This tests multi-pass parsing with file() errors
	locals2Content := `locals {
  # Reference the local from test_resource.pp
  derived_metadata = merge(local.test_metadata, {
    extra = "data"
  })
}`
	if err := os.WriteFile(filepath.Join(depModDir, "other_resource.pp"), []byte(locals2Content), 0600); err != nil {
		t.Fatal(err)
	}

	// Try to build eval context - this should NOT crash
	// Even though file() references non-existent files
	ctx := context.Background()
	builder := NewEvalContextBuilder(tmpDir)

	evalCtx, err := builder.Build(ctx)

	// We should get an eval context, not a crash
	if err != nil {
		t.Errorf("Build() returned error: %v", err)
	}

	if evalCtx == nil {
		t.Fatal("Build() returned nil eval context")
	}

	// NOTE: Functions is intentionally nil in the eval context from Build()
	// because functions are added later by the loader with the correct base path
	// for each file being parsed (see eval_context.go line 97-101)

	// Check that we have variables map
	if evalCtx.Variables == nil {
		t.Error("Eval context has no variables")
	}

	// The key test: we should NOT have crashed, even though:
	// 1. Dependency mod has locals with file() calls
	// 2. The file paths don't exist
	// 3. Multi-pass parsing tries to evaluate these locals
	//
	// Expected behavior:
	// - Locals that can't be evaluated should be skipped gracefully
	// - Build() should complete without crashing
	// - Some locals may not be resolved, but that's OK for index building

	t.Logf("Successfully built eval context with %d variables and %d locals",
		len(builder.variables), len(builder.locals))

	// If we get here without crashing, the test passes
	// The actual number of resolved locals doesn't matter -
	// what matters is that we don't crash on file() errors
}

// TestEvalContext_DependencyModWithMissingDocs tests with actual aws-compliance pattern
// where locals reference markdown documentation files that may or may not exist
func TestEvalContext_DependencyModWithMissingDocs(t *testing.T) {
	tmpDir := t.TempDir()

	// Create main workspace
	modContent := `mod "local" {
  title = "pskr"
}`
	if err := os.WriteFile(filepath.Join(tmpDir, "mod.pp"), []byte(modContent), 0600); err != nil {
		t.Fatal(err)
	}

	// Create dependency mod structure mimicking aws-compliance
	modsDir := filepath.Join(tmpDir, ".powerpipe", "mods")
	depModDir := filepath.Join(modsDir, "github.com", "turbot", "steampipe-mod-aws-compliance@v1.13.0")
	foundationalDir := filepath.Join(depModDir, "foundational_security")

	if err := os.MkdirAll(foundationalDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create mod.pp
	depModContent := `mod "aws-compliance" {
  title = "AWS Compliance"
}`
	if err := os.WriteFile(filepath.Join(depModDir, "mod.pp"), []byte(depModContent), 0600); err != nil {
		t.Fatal(err)
	}

	// Create guardduty.pp with file() call like the actual error
	// This mimics line 72 in foundational_security/guardduty.pp
	guarddutyContent := `locals {
  foundational_security_guardduty_common_tags = {
    service = "AWS/GuardDuty"
  }

  # This file() call on line 72 references a docs file that doesn't exist
  # In actual aws-compliance, this would be:
  # documentation = file("./foundational_security/docs/foundational_security_guardduty_7.md")
  guardduty_7_docs = file("./foundational_security/docs/foundational_security_guardduty_7.md")
}`
	if err := os.WriteFile(filepath.Join(foundationalDir, "guardduty.pp"), []byte(guarddutyContent), 0600); err != nil {
		t.Fatal(err)
	}

	// Try to build eval context
	ctx := context.Background()
	builder := NewEvalContextBuilder(tmpDir)

	// This should NOT crash with "no file exists" error
	evalCtx, err := builder.Build(ctx)

	if err != nil {
		t.Errorf("Build() crashed with error: %v", err)
		t.Errorf("This reproduces the issue where installing aws-compliance crashes the server")
	}

	if evalCtx == nil {
		t.Fatal("Build() returned nil eval context")
	}

	t.Logf("Build completed without crash - resolved %d locals", len(builder.locals))

	// The test passes if we don't crash
	// In v1.4.3 this worked, so it should work now too
}
