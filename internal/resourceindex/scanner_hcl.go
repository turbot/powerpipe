package resourceindex

import (
	"os"

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
	if diags.HasErrors() {
		// Continue with partial results - extract what we can from valid portions
		// The HCL parser provides good error recovery
	}

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
			entry.Title = extractStringLiteral(attr.Expr)
		case "description":
			entry.Description = extractStringLiteral(attr.Expr)
		case "sql":
			entry.HasSQL = true
		case "query":
			entry.HasSQL = true
			// Try to extract query reference
			if ref := extractReference(attr.Expr); ref != "" {
				entry.QueryRef = intern.Intern(s.modName + "." + ref)
			}
		}
	}
}

func (s *Scanner) extractHCLChildren(body *hclsyntax.Body, entry *IndexEntry) {
	attr, ok := body.Attributes["children"]
	if !ok {
		return
	}

	// Children is typically a tuple expression: [benchmark.a, control.b]
	if tuple, ok := attr.Expr.(*hclsyntax.TupleConsExpr); ok {
		for _, elem := range tuple.Exprs {
			if ref := extractReference(elem); ref != "" {
				fullRef := intern.Intern(s.modName + "." + ref)
				entry.ChildNames = append(entry.ChildNames, fullRef)
			}
		}
	}
}

func (s *Scanner) extractHCLTags(body *hclsyntax.Body, entry *IndexEntry) {
	// Tags can be an attribute with object value
	if attr, ok := body.Attributes["tags"]; ok {
		if obj, ok := attr.Expr.(*hclsyntax.ObjectConsExpr); ok {
			entry.Tags = make(map[string]string)
			for _, item := range obj.Items {
				// Key and value should both be extractable
				key := extractObjectKey(item.KeyExpr)
				val := extractStringLiteral(item.ValueExpr)
				if key != "" {
					// Intern tag keys (often repeated: service, category, etc.)
					entry.Tags[intern.Intern(key)] = val
				}
			}
		}
	}

	// Tags can also be a nested block (less common but valid HCL)
	for _, block := range body.Blocks {
		if block.Type == "tags" {
			entry.Tags = make(map[string]string)
			for name, attr := range block.Body.Attributes {
				val := extractStringLiteral(attr.Expr)
				entry.Tags[intern.Intern(name)] = val
			}
		}
	}
}

// extractStringLiteral attempts to get a string value from an HCL expression
// without full evaluation. Only works for literal strings.
func extractStringLiteral(expr hcl.Expression) string {
	switch e := expr.(type) {
	case *hclsyntax.TemplateExpr:
		// Simple string "foo" is a TemplateExpr with one LiteralValueExpr part
		if len(e.Parts) == 1 {
			if lit, ok := e.Parts[0].(*hclsyntax.LiteralValueExpr); ok {
				if lit.Val.Type() == cty.String {
					return lit.Val.AsString()
				}
			}
		}
		// For template expressions with multiple parts (string interpolation),
		// we can try to concatenate the literal parts
		var result string
		for _, part := range e.Parts {
			if lit, ok := part.(*hclsyntax.LiteralValueExpr); ok {
				if lit.Val.Type() == cty.String {
					result += lit.Val.AsString()
				}
			}
			// Non-literal parts (variables, function calls) are skipped
			// This gives us partial values for mixed templates
		}
		return result

	case *hclsyntax.LiteralValueExpr:
		if e.Val.Type() == cty.String {
			return e.Val.AsString()
		}

	case *hclsyntax.TemplateWrapExpr:
		// Unwrap and recurse
		return extractStringLiteral(e.Wrapped)
	}

	// Non-literal expressions (variables, function calls) return empty
	// This is fine - we'll get the real value on full parse
	return ""
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

// extractReference extracts a "type.name" reference from an expression.
// References like benchmark.child are ScopeTraversalExpr.
func extractReference(expr hcl.Expression) string {
	if trav, ok := expr.(*hclsyntax.ScopeTraversalExpr); ok {
		if len(trav.Traversal) >= 2 {
			root := trav.Traversal.RootName()
			if attr, ok := trav.Traversal[1].(hcl.TraverseAttr); ok {
				return root + "." + attr.Name
			}
		}
	}
	return ""
}
