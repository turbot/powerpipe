package resourceloader

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/turbot/pipe-fittings/v2/funcs"
)

// TestFileFunction_BasePath tests that file() function uses correct base path
// when called from different directories
func TestFileFunction_BasePath(t *testing.T) {
	// Create temporary mod structure
	tmpDir := t.TempDir()

	// Create a docs subdirectory and a file in it
	docsDir := filepath.Join(tmpDir, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		t.Fatal(err)
	}

	testContent := "# Test Document\nThis is a test markdown file."
	docsFile := filepath.Join(docsDir, "test.md")
	if err := os.WriteFile(docsFile, []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create an HCL expression that uses file() with relative path
	hclCode := `test_value = file("./docs/test.md")`

	file, diags := hclsyntax.ParseConfig([]byte(hclCode), "test.hcl", hcl.InitialPos)
	if diags.HasErrors() {
		t.Fatalf("Failed to parse HCL: %v", diags)
	}

	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		t.Fatal("Failed to get HCL body")
	}

	// Create eval context with file() function pointing to tmpDir as base
	evalCtx := &hcl.EvalContext{
		Functions: funcs.ContextFunctions(tmpDir),
	}

	// Try to evaluate the expression
	attr, exists := body.Attributes["test_value"]
	if !exists {
		t.Fatal("Attribute not found")
	}

	val, diags := attr.Expr.Value(evalCtx)
	if diags.HasErrors() {
		t.Errorf("Failed to evaluate file() expression: %v", diags)
		t.Errorf("This means file() function is not using the base path correctly")
	}

	if val.IsNull() {
		t.Error("file() returned null value")
	}

	// Check that we got the file content
	result := val.AsString()
	if result != testContent {
		t.Errorf("file() returned wrong content.\nExpected: %s\nGot: %s", testContent, result)
	}

	t.Logf("file() function correctly resolved path relative to base: %s", tmpDir)
}

// TestFileFunction_SubdirectoryBase tests file() when the HCL file is in a subdirectory
// This mimics the actual guardduty.pp scenario:
// - File location: .../foundational_security/guardduty.pp
// - file() call: file("./foundational_security/docs/...")
// - Base should be: mod root, not guardduty.pp location
func TestFileFunction_SubdirectoryBase(t *testing.T) {
	// Create mod structure:
	// mod_root/
	//   foundational_security/
	//     guardduty.pp (with file() call)
	//     docs/
	//       test.md

	modRoot := t.TempDir()
	fsDir := filepath.Join(modRoot, "foundational_security")
	docsDir := filepath.Join(fsDir, "docs")

	if err := os.MkdirAll(docsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create the docs file
	testContent := "# GuardDuty Documentation"
	docsFile := filepath.Join(docsDir, "test.md")
	if err := os.WriteFile(docsFile, []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create guardduty.pp with file() call
	// The path is relative to MOD ROOT, not to guardduty.pp location
	hclCode := `test_value = file("./foundational_security/docs/test.md")`

	file, diags := hclsyntax.ParseConfig([]byte(hclCode), "guardduty.pp", hcl.InitialPos)
	if diags.HasErrors() {
		t.Fatalf("Failed to parse HCL: %v", diags)
	}

	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		t.Fatal("Failed to get HCL body")
	}

	// CRITICAL: The base path for funcs.ContextFunctions should be modRoot,
	// NOT the directory containing guardduty.pp
	evalCtx := &hcl.EvalContext{
		Functions: funcs.ContextFunctions(modRoot), // Use mod root as base
	}

	// Try to evaluate the expression
	attr, exists := body.Attributes["test_value"]
	if !exists {
		t.Fatal("Attribute not found")
	}

	val, diags := attr.Expr.Value(evalCtx)
	if diags.HasErrors() {
		t.Errorf("Failed to evaluate file() expression: %v", diags)
		t.Errorf("This reproduces the issue where file() can't find docs in subdirectory")
		t.Errorf("Expected path: %s", docsFile)

		// This is the bug - file() is looking for the file relative to the wrong base
		return
	}

	if val.IsNull() {
		t.Error("file() returned null value")
	}

	// Check that we got the file content
	result := val.AsString()
	if result != testContent {
		t.Errorf("file() returned wrong content.\nExpected: %s\nGot: %s", testContent, result)
	}

	t.Logf("file() function correctly resolved subdirectory path relative to mod root: %s", modRoot)
}
