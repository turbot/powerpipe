package resourceloader

import (
	"bytes"
	"context"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	filehelpers "github.com/turbot/go-kit/files"
	"github.com/turbot/pipe-fittings/v2/funcs"
	"github.com/zclconf/go-cty/cty"
)

// Byte patterns for fast pre-scanning.
// These are used to quickly filter which files need full HCL parsing.
var (
	// Pattern to detect variable blocks: matches `variable` at start of line
	variablePattern = []byte("variable")
	// Pattern to detect locals blocks: `locals {` or `locals{`
	localsPattern = []byte("locals")
)

// EvalContextBuilder builds an HCL evaluation context from workspace variables and locals.
// This enables lazy loading to resolve expressions like `tags = local.common_tags`.
//
// Optimization: Uses a fast byte-level pre-scan to identify which files actually contain
// variable or locals definitions, avoiding expensive HCL parsing for the majority of
// files in large mods (which are typically just benchmarks/controls).
type EvalContextBuilder struct {
	workspacePath string
	variables     map[string]cty.Value // var.xxx -> value
	locals        map[string]cty.Value // local.xxx -> value
}

// NewEvalContextBuilder creates a new eval context builder for a workspace.
func NewEvalContextBuilder(workspacePath string) *EvalContextBuilder {
	return &EvalContextBuilder{
		workspacePath: workspacePath,
		variables:     make(map[string]cty.Value),
		locals:        make(map[string]cty.Value),
	}
}

// Build scans the workspace for variables and locals, evaluates them,
// and returns an hcl.EvalContext ready for use in resource parsing.
func (b *EvalContextBuilder) Build(ctx context.Context) (*hcl.EvalContext, error) {
	// List all .pp and .sp files
	files, err := b.listFiles(ctx)
	if err != nil {
		return nil, err
	}

	// Fast pre-scan: identify which files contain variables or locals
	// This avoids expensive HCL parsing for files that don't need it
	varFiles, localFiles, err := b.prescanFiles(files)
	if err != nil {
		return nil, err
	}

	// NOTE: We don't return early even if main workspace has no variables/locals
	// because we still need to scan dependency mods which may have variables/locals
	// that are referenced by resources in those mods

	// First pass: collect all variable defaults (only from files with variables)
	for filePath, content := range varFiles {
		if err := b.parseVariables(filePath, content); err != nil {
			// Continue on error - some files may have syntax issues
			continue
		}
	}

	// Build eval context with variables for locals evaluation
	evalCtx := &hcl.EvalContext{
		Functions: funcs.ContextFunctions(b.workspacePath),
		Variables: map[string]cty.Value{
			"var": cty.ObjectVal(b.variables),
		},
	}

	// Second pass: evaluate locals (only from files with locals)
	for filePath, content := range localFiles {
		if err := b.parseLocals(filePath, content, evalCtx); err != nil {
			// Continue on error
			continue
		}
	}

	// Third pass: scan dependency mods for their variables and locals
	// This is critical for Pipes scenarios where benchmarks from dependency mods
	// (like aws_compliance, aws_insights) use variables defined in those mods.
	// Ignore errors - we can still use variables from main workspace even if dependency scanning fails
	_ = b.ScanDependencyMods(ctx)

	// Build final eval context with both variables and locals
	finalCtx := &hcl.EvalContext{
		Functions: funcs.ContextFunctions(b.workspacePath),
		Variables: map[string]cty.Value{
			"var":   cty.ObjectVal(b.variables),
			"local": cty.ObjectVal(b.locals),
		},
	}

	return finalCtx, nil
}

// prescanFiles performs a fast byte-level scan to identify which files
// contain variable or locals definitions. Returns maps of filepath -> content
// for files that need full parsing.
//
// This is much faster than HCL parsing because:
// 1. bytes.Contains is very fast (uses optimized assembly on most platforms)
// 2. Most files in large mods don't contain variables or locals
// 3. We reuse the file content we already read for the actual parsing
func (b *EvalContextBuilder) prescanFiles(files []string) (varFiles, localFiles map[string][]byte, err error) {
	varFiles = make(map[string][]byte)
	localFiles = make(map[string][]byte)

	for _, filePath := range files {
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue // Skip files we can't read
		}

		// Fast check for variable definitions
		if bytes.Contains(content, variablePattern) {
			varFiles[filePath] = content
		}

		// Fast check for locals definitions
		if bytes.Contains(content, localsPattern) {
			localFiles[filePath] = content
		}
	}

	return varFiles, localFiles, nil
}

// listFiles returns all .pp and .sp files in the workspace.
func (b *EvalContextBuilder) listFiles(ctx context.Context) ([]string, error) {
	listOpts := &filehelpers.ListOptions{
		Flags:   filehelpers.FilesRecursive,
		Include: []string{"**/*.pp", "**/*.sp"},
		Exclude: []string{".*/**"}, // Skip hidden directories
	}
	return filehelpers.ListFilesWithContext(ctx, b.workspacePath, listOpts)
}

// parseVariables extracts variable blocks and their defaults from file content.
func (b *EvalContextBuilder) parseVariables(filePath string, content []byte) error {
	file, diags := hclsyntax.ParseConfig(content, filePath, hcl.InitialPos)
	if diags.HasErrors() || file == nil || file.Body == nil {
		return nil // Continue with partial results
	}

	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return nil
	}

	for _, block := range body.Blocks {
		if block.Type != "variable" || len(block.Labels) == 0 {
			continue
		}

		varName := block.Labels[0]

		// Look for default attribute
		if defaultAttr, ok := block.Body.Attributes["default"]; ok {
			// Create a simple eval context with just functions for evaluating the default
			evalCtx := &hcl.EvalContext{
				Functions: funcs.ContextFunctions(b.workspacePath),
			}

			val, diags := defaultAttr.Expr.Value(evalCtx)
			if !diags.HasErrors() {
				b.variables[varName] = val
			}
		}
	}

	return nil
}

// parseLocals extracts locals blocks and evaluates them from file content.
func (b *EvalContextBuilder) parseLocals(filePath string, content []byte, evalCtx *hcl.EvalContext) error {
	file, diags := hclsyntax.ParseConfig(content, filePath, hcl.InitialPos)
	if diags.HasErrors() {
		// Parse errors - skip this file
		return nil
	}
	if file == nil || file.Body == nil {
		return nil
	}

	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return nil
	}

	localsCount := 0
	for _, block := range body.Blocks {
		if block.Type != "locals" {
			continue
		}

		localsCount++

		// Evaluate each attribute in the locals block
		for name, attr := range block.Body.Attributes {
			val, diags := attr.Expr.Value(evalCtx)
			if !diags.HasErrors() {
				b.locals[name] = val
			}
			// If evaluation fails, skip this local (it may have unresolved references)
			// It might be resolved in a later pass
		}
	}

	return nil
}

// BuildEvalContext is a convenience function to build an eval context for a workspace.
func BuildEvalContext(ctx context.Context, workspacePath string) (*hcl.EvalContext, error) {
	builder := NewEvalContextBuilder(workspacePath)
	return builder.Build(ctx)
}

// MergeEvalContext merges functions into an existing eval context.
// This is useful when the loader needs to add functions to a pre-built context.
func MergeEvalContext(base *hcl.EvalContext, rootPath string) *hcl.EvalContext {
	if base == nil {
		return &hcl.EvalContext{
			Functions: funcs.ContextFunctions(rootPath),
		}
	}

	// If base has no functions, add them
	if base.Functions == nil {
		base.Functions = funcs.ContextFunctions(rootPath)
	}

	return base
}

// ScanDependencyMods scans dependency mods for variables and locals.
// This is needed because dependency mods can also define variables/locals that
// resources in those mods may reference.
func (b *EvalContextBuilder) ScanDependencyMods(ctx context.Context) error {
	modsDir := filepath.Join(b.workspacePath, ".powerpipe", "mods")
	info, err := os.Stat(modsDir)
	if os.IsNotExist(err) || (err == nil && !info.IsDir()) {
		return nil // No dependency mods directory
	}
	if err != nil {
		return err // Unexpected error accessing mods directory
	}

	err = filepath.Walk(modsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil //nolint:nilerr // Intentionally skip inaccessible paths and continue walking
		}

		// Look for mod.pp or mod.sp files
		if info.IsDir() || (info.Name() != "mod.pp" && info.Name() != "mod.sp") {
			return nil
		}

		// Found a mod - scan its directory
		modDir := filepath.Dir(path)

		files, err := b.listFilesInDir(ctx, modDir)
		if err != nil {
			return nil //nolint:nilerr // Skip mods we can't list files for
		}

		// Pre-scan files in this mod
		varFiles, localFiles, err := b.prescanFiles(files)
		if err != nil {
			return nil //nolint:nilerr // Skip mods we can't prescan
		}

		// Parse variables first
		for filePath, content := range varFiles {
			_ = b.parseVariables(filePath, content)
		}

		// Build intermediate context for locals evaluation
		// IMPORTANT: Need to include both existing variables AND existing locals
		// because locals in dependency mods may reference other locals
		evalCtx := &hcl.EvalContext{
			Functions: funcs.ContextFunctions(modDir),
			Variables: map[string]cty.Value{
				"var":   cty.ObjectVal(b.variables),
				"local": cty.ObjectVal(b.locals),
			},
		}

		// Parse locals - need to make multiple passes because locals can reference each other
		// across different files. Keep parsing until no new locals are added.
		// CRITICAL: All files must be parsed in each pass because locals in different files
		// reference each other (e.g., ec2.pp references locals from all_controls.pp)
		maxPasses := 10  // Increase to handle deep dependency chains
		for pass := 0; pass < maxPasses; pass++ {
			localsBefore := len(b.locals)

			// Update eval context with current locals for this pass
			if pass > 0 {
				evalCtx.Variables["local"] = cty.ObjectVal(b.locals)
			}

			// Parse ALL files in this pass - this allows cross-file local references
			for filePath, content := range localFiles {
				_ = b.parseLocals(filePath, content, evalCtx)
			}

			localsAfter := len(b.locals)
			localsAdded := localsAfter - localsBefore

			// If no new locals were added, we're done
			if localsAdded == 0 {
				break
			}

			// Update eval context for next pass
			evalCtx.Variables["local"] = cty.ObjectVal(b.locals)
		}

		return filepath.SkipDir // Don't recurse into subdirectories of this mod
	})

	return err
}

// listFilesInDir returns all .pp and .sp files in a specific directory.
func (b *EvalContextBuilder) listFilesInDir(ctx context.Context, dir string) ([]string, error) {
	listOpts := &filehelpers.ListOptions{
		Flags:   filehelpers.FilesRecursive,
		Include: []string{"**/*.pp", "**/*.sp"},
		Exclude: []string{".*/**"},
	}
	return filehelpers.ListFilesWithContext(ctx, dir, listOpts)
}
