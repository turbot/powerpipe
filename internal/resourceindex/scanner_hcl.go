package resourceindex

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/turbot/powerpipe/internal/intern"
	"github.com/zclconf/go-cty/cty"
)

// ScanFileHCL scans a file using HCL syntax parsing instead of regex.
// This handles all HCL syntax edge cases correctly but may be slower.
func (s *Scanner) ScanFileHCL(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return s.scanBytesHCL(content, filePath, false)
}

// ScanFileHCLWithOffsets scans a file using HCL syntax parsing and tracks byte offsets.
func (s *Scanner) ScanFileHCLWithOffsets(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return s.scanBytesHCL(content, filePath, true)
}

// ScanBytesHCL scans HCL content from a byte slice using the HCL parser.
func (s *Scanner) ScanBytesHCL(content []byte, filePath string) error {
	return s.scanBytesHCL(content, filePath, false)
}

func (s *Scanner) scanBytesHCL(content []byte, filePath string, withOffsets bool) error {
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

	return s.processHCLBody(body, filePath, withOffsets)
}

func (s *Scanner) processHCLBody(body *hclsyntax.Body, filePath string, withOffsets bool) error {
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
		s.extractHCLAttributes(block.Body, entry)

		// Extract children from block body
		s.extractHCLChildren(block.Body, entry)

		// Extract tags from block body
		s.extractHCLTags(block.Body, entry)

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

func (s *Scanner) extractHCLAttributes(body *hclsyntax.Body, entry *IndexEntry) {
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

func (s *Scanner) extractHCLChildren(body *hclsyntax.Body, entry *IndexEntry) {
	attr, ok := body.Attributes["children"]
	if !ok {
		return
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
			}
		}
	}
}

func (s *Scanner) extractHCLTags(body *hclsyntax.Body, entry *IndexEntry) {
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
