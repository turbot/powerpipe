// Package resourceindex provides fast HCL scanning for resource metadata extraction.
// It uses the HCL syntax parser for correctness while avoiding full expression evaluation,
// enabling fast startup times for large workspaces.
package resourceindex

import (
	"context"
	"crypto/md5" //nolint:gosec // MD5 used for consistent hashing to match pipe-fittings format, not for security
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	filehelpers "github.com/turbot/go-kit/files"
	"github.com/turbot/powerpipe/internal/intern"
	"github.com/zclconf/go-cty/cty"
)

// ResourceTypes that should be indexed for lazy loading
var indexedTypes = map[string]bool{
	"dashboard":           true,
	"benchmark":           true,
	"control":             true,
	"query":               true,
	"card":                true,
	"chart":               true,
	"container":           true,
	"flow":                true,
	"graph":               true,
	"hierarchy":           true,
	"image":               true,
	"input":               true,
	"node":                true,
	"edge":                true,
	"table":               true,
	"text":                true,
	"category":            true,
	"detection":           true,
	"detection_benchmark": true,
	"variable":            true,
	"with":                true,
}

// Scanner extracts resource metadata from HCL files without full parsing.
// This uses the HCL syntax parser for correctness while avoiding full expression evaluation.
type Scanner struct {
	modName string
	modRoot string // Root directory of the current mod (for file() function resolution)
	index   *ResourceIndex
	mu      sync.Mutex // protects index during parallel scanning
}

// NewScanner creates a new scanner for extracting resource metadata.
func NewScanner(modName string) *Scanner {
	return &Scanner{
		modName: modName,
		index:   NewResourceIndex(),
	}
}

// SetModRoot sets the root directory for the current mod.
func (s *Scanner) SetModRoot(root string) {
	s.modRoot = root
}

// ScanFile extracts index entries from a single HCL file.
// Always captures byte offsets for source_definition retrieval.
func (s *Scanner) ScanFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return s.scanBytes(content, filePath, true)
}

// ScanFileWithOffsets extracts entries with byte offsets for efficient seeking.
func (s *Scanner) ScanFileWithOffsets(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return s.scanBytes(content, filePath, true)
}

// ScanBytes scans HCL content from a byte slice.
// Useful for testing without file I/O.
func (s *Scanner) ScanBytes(content []byte, filePath string) error {
	return s.scanBytes(content, filePath, false)
}

func (s *Scanner) scanBytes(content []byte, filePath string, withOffsets bool) error {
	// Syntax parse only - no expression evaluation
	file, diags := hclsyntax.ParseConfig(content, filePath, hcl.InitialPos)
	// Continue with partial results on errors - HCL parser provides good error recovery
	_ = diags

	if file == nil || file.Body == nil {
		return nil
	}

	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return nil
	}

	return s.processBody(body, filePath, withOffsets, content)
}

func (s *Scanner) processBody(body *hclsyntax.Body, filePath string, withOffsets bool, content []byte) error {
	for _, block := range body.Blocks {
		if !indexedTypes[block.Type] {
			continue
		}

		if len(block.Labels) == 0 {
			continue
		}

		// Use interned strings for commonly repeated values
		internedModName := intern.Intern(s.modName)
		internedType := intern.Intern(block.Type)
		internedShortName := intern.Intern(block.Labels[0])
		internedFileName := intern.Intern(filePath)

		// Build full name with interned components
		fullName := intern.Intern(internedModName + "." + internedType + "." + internedShortName)

		// Intern mod root if set
		internedModRoot := ""
		if s.modRoot != "" {
			internedModRoot = intern.Intern(s.modRoot)
		}

		// Build mod full name (mod.modname format)
		internedModFullName := intern.Intern("mod." + internedModName)

		entry := &IndexEntry{
			Type:        internedType,
			Name:        fullName,
			ShortName:   internedShortName,
			FileName:    internedFileName,
			StartLine:   block.TypeRange.Start.Line,
			EndLine:     block.Body.EndRange.Start.Line,
			ModName:     internedModName,
			ModFullName: internedModFullName,
			ModRoot:     internedModRoot,
		}

		// Calculate byte offsets if requested
		if withOffsets {
			entry.ByteOffset = int64(block.TypeRange.Start.Byte)
			entry.ByteLength = block.Body.EndRange.End.Byte - block.TypeRange.Start.Byte
		}

		// Extract attributes from block body
		s.extractAttributes(block.Body, entry)

		// Extract children from block body (including references for benchmarks)
		s.extractChildren(block.Body, entry, content)

		// Extract anonymous nested children for dashboards
		// These are inline container/card/etc blocks that get auto-generated names
		if block.Type == "dashboard" {
			s.extractAnonymousChildren(block.Body, entry)
		}

		// Extract tags from block body
		s.extractTags(block.Body, entry)

		// Set benchmark type
		if block.Type == "benchmark" {
			entry.BenchmarkType = intern.BenchmarkTypeControl
		} else if block.Type == "detection_benchmark" {
			entry.BenchmarkType = intern.BenchmarkTypeDetection
		}

		s.addEntry(entry)
	}

	return nil
}

func (s *Scanner) extractAttributes(body *hclsyntax.Body, entry *IndexEntry) {
	for name, attr := range body.Attributes {
		switch name {
		case "title":
			entry.Title, entry.TitleResolved = extractStringWithResolution(attr.Expr)
		case "description":
			entry.Description, entry.DescriptionResolved = extractStringWithResolution(attr.Expr)
		case "category":
			entry.Category, _ = extractStringWithResolution(attr.Expr)
		case "documentation":
			entry.Documentation, _ = extractStringWithResolution(attr.Expr)
		case "display":
			entry.Display, _ = extractStringWithResolution(attr.Expr)
		case "width":
			if width, ok := extractIntLiteral(attr.Expr); ok {
				entry.Width = &width
			}
		case "sql":
			entry.HasSQL = true
			// Try to extract literal SQL value
			if sqlText, resolved := extractStringWithResolution(attr.Expr); resolved && sqlText != "" {
				entry.SQL = sqlText
			}
			// Check if sql is a reference like query.xxx.sql
			if ref := extractSQLQueryReference(attr.Expr, s.modName); ref != "" {
				entry.QueryRef = intern.Intern(ref)
			}
		case "query":
			entry.HasSQL = true
			// Try to extract query reference
			if ref := extractReference(attr.Expr); ref != "" {
				entry.QueryRef = intern.Intern(s.modName + "." + ref)
			}
		}
	}

	// If title/description attributes weren't present, mark as resolved (to empty)
	if _, exists := body.Attributes["title"]; !exists {
		entry.TitleResolved = true
	}
	if _, exists := body.Attributes["description"]; !exists {
		entry.DescriptionResolved = true
	}
}

func (s *Scanner) extractChildren(body *hclsyntax.Body, entry *IndexEntry, content []byte) {
	attr, ok := body.Attributes["children"]
	if !ok {
		return
	}

	// Extract source definition for the children attribute
	var childrenSourceDef string
	startLine := attr.SrcRange.Start.Line
	endLine := attr.SrcRange.End.Line
	if len(content) > 0 {
		childrenSourceDef = extractSourceDefinition(content, attr.SrcRange.Start.Byte, attr.SrcRange.End.Byte)
	}

	// Children is typically a tuple expression: [benchmark.a, control.b]
	// References can be:
	//   - "type.name" (local reference, needs mod prefix)
	//   - "mod.type.name" (cross-mod reference, already qualified)
	if tuple, ok := attr.Expr.(*hclsyntax.TupleConsExpr); ok {
		for _, elem := range tuple.Exprs {
			if ref := extractReference(elem); ref != "" {
				// Count parts to determine if it's a local or cross-mod reference
				parts := strings.Count(ref, ".") + 1
				var fullRef string
				if parts == 2 {
					// Local reference (type.name) - prefix with current mod
					fullRef = intern.Intern(s.modName + "." + ref)
				} else {
					// Cross-mod reference (mod.type.name) - already qualified
					fullRef = intern.Intern(ref)
				}
				entry.ChildNames = append(entry.ChildNames, fullRef)

				// Create a reference entry for benchmarks
				if entry.Type == "benchmark" || entry.Type == "detection_benchmark" {
					referenceFrom := entry.Type + "." + entry.ShortName
					reference := Reference{
						AutoGenerated:    false,
						EndLineNumber:    endLine,
						FileName:         entry.FileName,
						FromAttribute:    "children",
						FromBlockName:    entry.ShortName,
						FromBlockType:    entry.Type,
						IsAnonymous:      false,
						ModName:          entry.ModName,
						ReferenceFrom:    referenceFrom,
						ReferenceTo:      ref,
						ResourceName:     generateResourceHash(ref, referenceFrom, entry.Type, entry.ShortName, "children"),
						SourceDefinition: childrenSourceDef,
						StartLineNumber:  startLine,
					}
					entry.References = append(entry.References, reference)
				}
			}
		}
	}
}

// extractSourceDefinition extracts the source text from content using byte offsets.
// It includes leading whitespace from the line start for proper indentation.
func extractSourceDefinition(content []byte, start, end int) string {
	if start < 0 || end > len(content) || start >= end {
		return ""
	}

	// Find the start of the line to include leading whitespace
	lineStart := start
	for lineStart > 0 && content[lineStart-1] != '\n' {
		lineStart--
	}

	return string(content[lineStart:end])
}

// generateResourceHash generates a short hash for the resource name.
// This matches the format used by eager loading (MD5 hash of reference string).
func generateResourceHash(to, from, blockType, blockName, attribute string) string {
	// Match the exact format from pipe-fittings ResourceReference.String()
	str := fmt.Sprintf("To: %s\nFrom: %s\nBlockType: %s\nBlockName: %s\nAttribute: %s",
		to, from, blockType, blockName, attribute)

	hash := md5.Sum([]byte(str)) //nolint:gosec // MD5 used for consistent hashing, not security
	return hex.EncodeToString(hash[:])[:8]
}

// dashboardChildTypes defines block types that can be nested children in dashboards
var dashboardChildTypes = map[string]bool{
	"container": true,
	"card":      true,
	"chart":     true,
	"flow":      true,
	"graph":     true,
	"hierarchy": true,
	"image":     true,
	"input":     true,
	"table":     true,
	"text":      true,
}

// extractAnonymousChildren scans nested blocks inside a dashboard and generates
// anonymous child names using the same naming convention as eager loading.
// Format: {mod_name}.{child_type}.dashboard_{parent_shortname}_anonymous_{child_type}_{index}
func (s *Scanner) extractAnonymousChildren(body *hclsyntax.Body, entry *IndexEntry) {
	// Track indices per block type (for anonymous naming)
	typeIndices := make(map[string]int)

	// Process nested blocks in order
	for _, block := range body.Blocks {
		if !dashboardChildTypes[block.Type] {
			continue
		}

		// Check if it's an anonymous block (no labels) or named block
		var childName string
		if len(block.Labels) > 0 {
			// Named child: mod_name.type.name
			childName = intern.Intern(s.modName + "." + block.Type + "." + block.Labels[0])
		} else {
			// Anonymous child: mod_name.type.dashboard_{parent}_anonymous_{type}_{index}
			idx := typeIndices[block.Type]
			typeIndices[block.Type] = idx + 1
			childName = intern.Intern(fmt.Sprintf("%s.%s.dashboard_%s_anonymous_%s_%d",
				s.modName, block.Type, entry.ShortName, block.Type, idx))
		}

		entry.ChildNames = append(entry.ChildNames, childName)
	}
}

func (s *Scanner) extractTags(body *hclsyntax.Body, entry *IndexEntry) {
	// Default to resolved if no tags attribute exists
	entry.TagsResolved = true

	// Tags can be an attribute with object value
	if attr, ok := body.Attributes["tags"]; ok {
		entry.Tags, entry.TagsResolved, entry.UnresolvedRefs = extractTagsComplete(attr.Expr)
	}

	// Tags can also be a nested block (less common but valid HCL)
	for _, block := range body.Blocks {
		if block.Type == "tags" {
			entry.Tags = make(map[string]string)
			entry.TagsResolved = true
			for name, attr := range block.Body.Attributes {
				val, resolved := extractStringWithResolution(attr.Expr)
				entry.Tags[intern.Intern(name)] = val
				if !resolved {
					entry.TagsResolved = false
					entry.UnresolvedRefs = append(entry.UnresolvedRefs, "tag:"+name)
				}
			}
		}
	}
}

func (s *Scanner) addEntry(entry *IndexEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.index.Add(entry)
}

// listPowerpipeFiles returns all .pp and .sp files in a directory using go-kit.
// This uses the same file listing logic as pipe-fittings for consistency.
func listPowerpipeFiles(ctx context.Context, dirPath string) ([]string, error) {
	// Validate directory exists and is a directory
	info, err := os.Stat(dirPath)
	if err != nil {
		return nil, fmt.Errorf("accessing directory %s: %w", dirPath, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dirPath)
	}

	listOpts := &filehelpers.ListOptions{
		Flags:   filehelpers.FilesRecursive,
		Include: []string{"**/*.pp", "**/*.sp"},
		Exclude: []string{".*/**"}, // Skip hidden directories
	}
	return filehelpers.ListFilesWithContext(ctx, dirPath, listOpts)
}

// ScanDirectory scans all .pp and .sp files in a directory recursively.
func (s *Scanner) ScanDirectory(dirPath string) error {
	return s.ScanDirectoryWithContext(context.Background(), dirPath)
}

// ScanDirectoryWithContext scans all .pp and .sp files with context support.
func (s *Scanner) ScanDirectoryWithContext(ctx context.Context, dirPath string) error {
	files, err := listPowerpipeFiles(ctx, dirPath)
	if err != nil {
		return err
	}

	for _, filePath := range files {
		if err := s.ScanFile(filePath); err != nil {
			return err
		}
	}
	return nil
}

// ScanDirectoryParallel scans files in parallel for faster indexing.
func (s *Scanner) ScanDirectoryParallel(dirPath string, workers int) error {
	return s.ScanDirectoryParallelWithContext(context.Background(), dirPath, workers)
}

// ScanDirectoryParallelWithContext scans files in parallel with context support.
func (s *Scanner) ScanDirectoryParallelWithContext(ctx context.Context, dirPath string, workers int) error {
	// Use go-kit for consistent file listing with pipe-fittings
	files, err := listPowerpipeFiles(ctx, dirPath)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return nil
	}

	// Set default workers
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	if workers > len(files) {
		workers = len(files)
	}

	// Process in parallel
	fileChan := make(chan string, len(files))
	errChan := make(chan error, 1)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range fileChan {
				if err := s.ScanFile(filePath); err != nil {
					select {
					case errChan <- fmt.Errorf("scanning %s: %w", filePath, err):
					default:
					}
					return
				}
			}
		}()
	}

	// Send files to workers
	for _, f := range files {
		fileChan <- f
	}
	close(fileChan)

	// Wait for completion
	wg.Wait()
	close(errChan)

	// Return first error if any
	if err := <-errChan; err != nil {
		return err
	}

	return nil
}

// ScanDirectoryWithModName scans a directory with a specific mod name.
// This is used for scanning dependency mods where each mod has its own name.
func (s *Scanner) ScanDirectoryWithModName(dirPath, modName string) error {
	// Save current mod name and root
	originalModName := s.modName
	originalModRoot := s.modRoot

	// Set the new mod name and root for this directory
	s.modName = modName
	s.modRoot = dirPath

	// Use go-kit for consistent file listing
	files, err := listPowerpipeFiles(context.Background(), dirPath)
	if err != nil {
		// Restore original mod name and root before returning
		s.modName = originalModName
		s.modRoot = originalModRoot
		return err
	}

	for _, filePath := range files {
		if err := s.ScanFile(filePath); err != nil {
			// Restore original mod name and root before returning
			s.modName = originalModName
			s.modRoot = originalModRoot
			return err
		}
	}

	// Restore original mod name and root
	s.modName = originalModName
	s.modRoot = originalModRoot

	return nil
}

// GetIndex returns the built index.
func (s *Scanner) GetIndex() *ResourceIndex {
	return s.index
}

// SetModInfo sets mod-level information on the index.
func (s *Scanner) SetModInfo(modName, modFullName, modTitle string) {
	// Intern mod names as they're repeated across all resources
	s.index.ModName = intern.Intern(modName)
	s.index.ModFullName = intern.Intern(modFullName)
	// Titles are usually unique, don't intern
	s.index.ModTitle = modTitle
}

// MarkTopLevelResources marks resources as top-level based on parent references.
// This should be called after scanning is complete.
func (s *Scanner) MarkTopLevelResources() {
	s.index.mu.Lock()
	defer s.index.mu.Unlock()

	// Build set of all child names
	childNames := make(map[string]bool)
	for _, entry := range s.index.entries {
		for _, child := range entry.ChildNames {
			childNames[child] = true
		}
	}

	// Mark entries that are not children of anything as top-level
	for _, entry := range s.index.entries {
		// Only mark benchmarks and dashboards as top-level
		if entry.Type == "benchmark" || entry.Type == "detection_benchmark" || entry.Type == "dashboard" {
			if !childNames[entry.Name] {
				entry.IsTopLevel = true
			}
		}
	}
}

// SetParentNames sets ParentName and ParentNames on child entries based on ChildNames.
// This should be called after scanning is complete.
// Controls can be children of multiple benchmarks, so we track all parents.
func (s *Scanner) SetParentNames() {
	s.index.mu.Lock()
	defer s.index.mu.Unlock()

	for _, entry := range s.index.entries {
		for _, childName := range entry.ChildNames {
			if child, ok := s.index.entries[childName]; ok {
				// Track all parents for multi-path support
				child.ParentNames = append(child.ParentNames, entry.Name)
				// Also set single ParentName for backwards compatibility (use first parent)
				if child.ParentName == "" {
					child.ParentName = entry.Name
				}
			}
		}
	}
}

// ComputePaths computes the full hierarchical paths for all entries.
// This should be called after SetParentNames.
func (s *Scanner) ComputePaths() {
	s.index.mu.Lock()
	defer s.index.mu.Unlock()

	for _, entry := range s.index.entries {
		// Use visited set to prevent cycles
		visited := make(map[string]bool)
		entry.Paths = s.buildPathsForEntry(entry, visited)
	}
}

// buildPathsForEntry builds all paths for an entry by traversing up the parent hierarchy.
// visited tracks entries already in the current path to prevent cycles.
func (s *Scanner) buildPathsForEntry(entry *IndexEntry, visited map[string]bool) [][]string {
	// Check for cycle
	if visited[entry.Name] {
		// Cycle detected, return simple path to break recursion
		return [][]string{{entry.ModFullName, entry.Name}}
	}
	visited[entry.Name] = true
	defer func() { visited[entry.Name] = false }()

	// If no parents, simple path: [mod, resource]
	if len(entry.ParentNames) == 0 {
		return [][]string{{entry.ModFullName, entry.Name}}
	}

	// Build a path for each parent, recursively including their ancestry
	var paths [][]string
	for _, parentName := range entry.ParentNames {
		parent, ok := s.index.entries[parentName]
		if !ok {
			// Parent not found, use simple path
			paths = append(paths, []string{entry.ModFullName, parentName, entry.Name})
			continue
		}

		// Get parent's paths and extend each with this entry
		parentPaths := s.buildPathsForEntry(parent, visited)
		for _, parentPath := range parentPaths {
			// Extend parent path with this entry's name
			fullPath := make([]string, len(parentPath)+1)
			copy(fullPath, parentPath)
			fullPath[len(parentPath)] = entry.Name
			paths = append(paths, fullPath)
		}
	}

	return paths
}

// -----------------------------------------------------------------------------
// HCL Expression Extraction Helpers
// -----------------------------------------------------------------------------

// extractTagsComplete extracts the full tags map with resolution tracking.
// It handles literal object expressions, variable references, and merge() calls.
func extractTagsComplete(expr hcl.Expression) (map[string]string, bool, []string) {
	tags := make(map[string]string)
	unresolvedRefs := []string{}
	allResolved := true

	switch e := expr.(type) {
	case *hclsyntax.ObjectConsExpr:
		// tags = { key = "value", ... }
		for _, item := range e.Items {
			key := extractObjectKey(item.KeyExpr)
			if key == "" {
				continue
			}

			value, resolved := extractExpressionValue(item.ValueExpr)
			tags[intern.Intern(key)] = value
			if !resolved {
				allResolved = false
				unresolvedRefs = append(unresolvedRefs, "tag:"+key)
			}
		}

	case *hclsyntax.ScopeTraversalExpr:
		// tags = var.common_tags - needs full resolution
		allResolved = false
		unresolvedRefs = append(unresolvedRefs, "tags")

	case *hclsyntax.FunctionCallExpr:
		// tags = merge(var.common_tags, { ... }) - needs resolution
		allResolved = false
		unresolvedRefs = append(unresolvedRefs, "tags")

		// Try to extract any inline object arguments
		for _, arg := range e.Args {
			if obj, ok := arg.(*hclsyntax.ObjectConsExpr); ok {
				for _, item := range obj.Items {
					key := extractObjectKey(item.KeyExpr)
					if key != "" {
						value, _ := extractExpressionValue(item.ValueExpr)
						tags[intern.Intern(key)] = value
					}
				}
			}
		}

	default:
		// Unknown expression type - mark as unresolved
		allResolved = false
		unresolvedRefs = append(unresolvedRefs, "tags")
	}

	return tags, allResolved, unresolvedRefs
}

// extractStringLiteral attempts to get a string value from an HCL expression
// without full evaluation. Only works for literal strings.
func extractStringLiteral(expr hcl.Expression) string {
	val, _ := extractStringWithResolution(expr)
	return val
}

// extractStringWithResolution extracts a string value and reports if it was fully resolved.
// Returns (value, resolved) where resolved=true means the value is complete (no interpolation/variables).
func extractStringWithResolution(expr hcl.Expression) (string, bool) {
	switch e := expr.(type) {
	case *hclsyntax.TemplateExpr:
		// Simple string "foo" is a TemplateExpr with one LiteralValueExpr part
		if len(e.Parts) == 1 {
			if lit, ok := e.Parts[0].(*hclsyntax.LiteralValueExpr); ok {
				if lit.Val.Type() == cty.String {
					return lit.Val.AsString(), true
				}
			}
		}
		// For template expressions with multiple parts (string interpolation),
		// we can try to concatenate the literal parts, but mark as unresolved
		var result string
		hasNonLiteral := false
		for _, part := range e.Parts {
			if lit, ok := part.(*hclsyntax.LiteralValueExpr); ok {
				if lit.Val.Type() == cty.String {
					result += lit.Val.AsString()
				}
			} else {
				// Non-literal parts (variables, function calls)
				hasNonLiteral = true
			}
		}
		return result, !hasNonLiteral

	case *hclsyntax.LiteralValueExpr:
		if e.Val.Type() == cty.String {
			return e.Val.AsString(), true
		}
		return "", true

	case *hclsyntax.TemplateWrapExpr:
		// Unwrap and recurse
		return extractStringWithResolution(e.Wrapped)

	case *hclsyntax.ScopeTraversalExpr:
		// Variable reference like var.title - needs resolution
		return "", false

	case *hclsyntax.FunctionCallExpr:
		// Function call - needs resolution
		return "", false

	case *hclsyntax.ConditionalExpr:
		// Conditional expression - needs resolution
		return "", false
	}

	// Unknown expression type - mark as unresolved
	return "", false
}

// extractExpressionValue extracts a value from any expression, returning (value, resolved).
// Works for strings, booleans, and numbers.
func extractExpressionValue(expr hcl.Expression) (string, bool) {
	switch e := expr.(type) {
	case *hclsyntax.LiteralValueExpr:
		if e.Val.Type() == cty.String {
			return e.Val.AsString(), true
		}
		if e.Val.Type() == cty.Bool {
			if e.Val.True() {
				return "true", true
			}
			return "false", true
		}
		if e.Val.Type() == cty.Number {
			f, _ := e.Val.AsBigFloat().Float64()
			return strings.TrimRight(strings.TrimRight(formatFloat(f), "0"), "."), true
		}
		return "", true

	case *hclsyntax.TemplateExpr:
		return extractStringWithResolution(e)

	case *hclsyntax.TemplateWrapExpr:
		return extractStringWithResolution(e.Wrapped)

	case *hclsyntax.ScopeTraversalExpr:
		// Variable reference - return placeholder
		return "", false

	default:
		return "", false
	}
}

// extractIntLiteral attempts to extract an integer value from an expression.
func extractIntLiteral(expr hcl.Expression) (int, bool) {
	switch e := expr.(type) {
	case *hclsyntax.LiteralValueExpr:
		if e.Val.Type() == cty.Number {
			f, _ := e.Val.AsBigFloat().Float64()
			return int(f), true
		}
	case *hclsyntax.TemplateWrapExpr:
		return extractIntLiteral(e.Wrapped)
	}
	return 0, false
}

// formatFloat formats a float64, returning an integer string if the value is whole
func formatFloat(f float64) string {
	// Check if it's a whole number
	if f == float64(int64(f)) {
		return fmt.Sprintf("%d", int64(f))
	}
	return fmt.Sprintf("%g", f)
}

// extractObjectKey extracts a key from an object key expression.
// Object keys can be identifiers, literals, or parenthesized expressions.
func extractObjectKey(expr hcl.Expression) string {
	switch e := expr.(type) {
	case *hclsyntax.ObjectConsKeyExpr:
		// Try to get the traversal name (for unquoted identifiers)
		if trav, ok := e.Wrapped.(*hclsyntax.ScopeTraversalExpr); ok {
			return trav.Traversal.RootName()
		}
		// Try as a literal string
		return extractStringLiteral(e.Wrapped)

	case *hclsyntax.ScopeTraversalExpr:
		// Direct traversal (identifier)
		return e.Traversal.RootName()

	case *hclsyntax.LiteralValueExpr:
		if e.Val.Type() == cty.String {
			return e.Val.AsString()
		}

	case *hclsyntax.TemplateExpr:
		return extractStringLiteral(e)
	}

	return ""
}

// extractSQLQueryReference extracts a query reference from an sql expression like query.xxx.sql
// Returns the full query name if the expression is of the form query.name.sql or mod.query.name.sql
func extractSQLQueryReference(expr hcl.Expression, modName string) string {
	if trav, ok := expr.(*hclsyntax.ScopeTraversalExpr); ok {
		if len(trav.Traversal) >= 3 {
			// Extract all parts
			parts := make([]string, 0, len(trav.Traversal))
			parts = append(parts, trav.Traversal.RootName())
			for i := 1; i < len(trav.Traversal); i++ {
				if attr, ok := trav.Traversal[i].(hcl.TraverseAttr); ok {
					parts = append(parts, attr.Name)
				}
			}

			// Check if last part is "sql"
			if parts[len(parts)-1] == "sql" {
				// Remove the "sql" part
				parts = parts[:len(parts)-1]

				// Now we have either:
				// - query.name (2 parts) - local reference
				// - mod.query.name (3 parts) - cross-mod reference
				if len(parts) == 2 && parts[0] == "query" {
					// Local query reference: query.name -> modName.query.name
					return modName + ".query." + parts[1]
				} else if len(parts) >= 3 {
					// Cross-mod reference: mod.query.name -> mod.query.name (already qualified)
					return strings.Join(parts, ".")
				}
			}
		}
	}
	return ""
}

// extractReference extracts a reference from an expression.
// References can be:
//   - "type.name" (e.g., benchmark.child, control.my_control)
//   - "mod.type.name" (e.g., dependency_1.control.version for cross-mod references)
func extractReference(expr hcl.Expression) string {
	if trav, ok := expr.(*hclsyntax.ScopeTraversalExpr); ok {
		if len(trav.Traversal) >= 2 {
			// Build the full reference from all traversal parts
			parts := make([]string, 0, len(trav.Traversal))
			parts = append(parts, trav.Traversal.RootName())
			for i := 1; i < len(trav.Traversal); i++ {
				if attr, ok := trav.Traversal[i].(hcl.TraverseAttr); ok {
					parts = append(parts, attr.Name)
				}
			}
			return strings.Join(parts, ".")
		}
	}
	return ""
}
