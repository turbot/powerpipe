package resourceloader

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

// TestEvalContext_EmptyMapsInCty tests that creating cty.ObjectVal with empty maps doesn't cause issues
func TestEvalContext_EmptyMapsInCty(t *testing.T) {
	// Create empty maps
	emptyVars := make(map[string]cty.Value)
	emptyLocals := make(map[string]cty.Value)

	// Try to create cty.ObjectVal with empty maps
	varObj := cty.ObjectVal(emptyVars)
	localObj := cty.ObjectVal(emptyLocals)

	if varObj.IsNull() {
		t.Error("cty.ObjectVal with empty map returned null")
	}

	if localObj.IsNull() {
		t.Error("cty.ObjectVal with empty map returned null")
	}

	t.Logf("Empty maps work fine in cty.ObjectVal")
}

// TestEvalContext_RealWorkspaceWithAwsCompliance tests building eval context
// with actual aws-compliance mod structure
func TestEvalContext_RealWorkspaceWithAwsCompliance(t *testing.T) {
	// Create workspace structure mimicking user's scenario
	tmpDir := t.TempDir()

	// Create mod.pp
	modContent := `mod "local" {
  title = "test"
}`
	if err := os.WriteFile(filepath.Join(tmpDir, "mod.pp"), []byte(modContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create .powerpipe/mods structure with simplified aws-compliance
	modsDir := filepath.Join(tmpDir, ".powerpipe", "mods")
	awsComplianceDir := filepath.Join(modsDir, "github.com", "turbot", "steampipe-mod-aws-compliance@v1.13.0")
	fsDir := filepath.Join(awsComplianceDir, "foundational_security")

	if err := os.MkdirAll(fsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create mod.pp for aws-compliance
	depModContent := `mod "aws-compliance" {
  title = "AWS Compliance"
}`
	if err := os.WriteFile(filepath.Join(awsComplianceDir, "mod.pp"), []byte(depModContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a locals file with complex expressions but no file() calls
	// (Based on my grep, aws-compliance doesn't have file() in locals)
	localsContent := `locals {
  foundational_security_guardduty_common_tags = {
    service = "AWS/GuardDuty"
  }

  audit_manager_control_tower_common_tags = {
    service = "AWS/AuditManager"
  }
}`
	if err := os.WriteFile(filepath.Join(fsDir, "locals.pp"), []byte(localsContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Try to build eval context - this mimics what happens during lazy workspace creation
	ctx := context.Background()
	builder := NewEvalContextBuilder(tmpDir)

	evalCtx, err := builder.Build(ctx)

	if err != nil {
		t.Errorf("Build() failed: %v", err)
	}

	if evalCtx == nil {
		t.Fatal("Build() returned nil eval context")
	}

	t.Logf("Successfully built eval context")
	t.Logf("Resolved %d variables, %d locals", len(builder.variables), len(builder.locals))

	// Check that we got some locals from the dependency mod
	if len(builder.locals) < 2 {
		t.Errorf("Expected at least 2 locals from dependency mod, got %d", len(builder.locals))
	}
}
