# Task 17: HCL Syntax-Based Scanner (Optional)

## Objective

Prototype an alternative scanner implementation using `hclsyntax.ParseConfig` instead of regex-based parsing. This would provide correctness guarantees by leveraging the canonical HCL parser while avoiding the expensive expression evaluation step.

## Priority

**OPTIONAL** - The current regex-based scanner works and has been hardened with edge case fixes. This task explores whether an HCL-based approach offers better maintainability and correctness at acceptable performance cost.

## Context

### Current Regex Approach

The regex-based scanner (`internal/resourceindex/scanner.go`) has required fixes for:
1. Escaped quotes in attribute values
2. Single-line block attributes
3. Single-line tags
4. Block comments (`/* ... */`)
5. Heredoc content parsing

Each edge case required custom state tracking and regex patterns. More edge cases likely exist (nested heredocs, template expressions, multi-line strings, etc.).

### HCL Library Layers

The HCL library provides multiple parsing levels:

| Level | Function | What It Does | Performance |
|-------|----------|--------------|-------------|
| Syntax Parse | `hclsyntax.ParseConfig` | Tokenize + AST, no evaluation | Medium |
| Full Parse | `hclparse.Parser.ParseHCL` | Above + decode to Go types | Slow |
| Full Decode | `gohcl.DecodeBody` | Above + expression evaluation | Slowest |

The **syntax parse** is the sweet spot - it handles all syntax correctly (comments, heredocs, escapes, strings) but skips expensive expression evaluation.

### What We Need for the Index

| Field | Source | Evaluation Needed? |
|-------|--------|-------------------|
| Block type | `block.Type` | No |
| Block name | `block.Labels[0]` | No |
| Title | `block.Body.Attributes["title"]` | String literal only |
| Description | `block.Body.Attributes["description"]` | String literal only |
| Tags | `block.Body.Blocks` (tags block) | String literals only |
| HasSQL | Attribute existence check | No |
| Children | `block.Body.Attributes["children"]` | Reference extraction |
| Line numbers | `block.TypeRange`, `block.EndRange()` | No |
| Byte offsets | Range positions | No |

For string literals like `title = "My Title"`, we can extract values directly from AST nodes without full evaluation.

## Proposed Implementation

### Core Scanner Function

```go
package resourceindex

import (
    "os"

    "github.com/hashicorp/hcl/v2"
    "github.com/hashicorp/hcl/v2/hclsyntax"
    "github.com/zclconf/go-cty/cty"
)

// ScanFileHCL scans a file using HCL syntax parsing instead of regex.
// This handles all HCL syntax edge cases correctly but is potentially slower.
func (s *Scanner) ScanFileHCL(filePath string) error {
    content, err := os.ReadFile(filePath)
    if err != nil {
        return err
    }

    // Syntax parse only - no expression evaluation
    file, diags := hclsyntax.ParseConfig(content, filePath, hcl.InitialPos)
    if diags.HasErrors() {
        // Log warning but continue - extract what we can
        // Or: return diags for strict mode
    }

    body, ok := file.Body.(*hclsyntax.Body)
    if !ok {
        return nil
    }

    return s.processHCLBody(body, filePath)
}

func (s *Scanner) processHCLBody(body *hclsyntax.Body, filePath string) error {
    for _, block := range body.Blocks {
        if !indexedTypes[block.Type] {
            continue
        }

        if len(block.Labels) == 0 {
            continue
        }

        entry := &IndexEntry{
            Type:      intern.Intern(block.Type),
            ShortName: intern.Intern(block.Labels[0]),
            FileName:  intern.Intern(filePath),
            StartLine: block.TypeRange.Start.Line,
            EndLine:   block.Body.EndRange.Start.Line,
            ModName:   intern.Intern(s.modName),
        }

        // Build full name
        entry.Name = intern.Intern(s.modName + "." + block.Type + "." + block.Labels[0])
        entry.ModFullName = intern.Intern("mod." + s.modName)

        // Extract attributes from block body
        s.extractHCLAttributes(block.Body, entry)

        // Extract children from nested blocks
        s.extractHCLChildren(block.Body, entry)

        // Extract tags
        s.extractHCLTags(block.Body, entry)

        s.addEntry(entry)
    }

    return nil
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
    case *hclsyntax.LiteralValueExpr:
        if e.Val.Type() == cty.String {
            return e.Val.AsString()
        }
    }
    // Non-literal expressions (variables, function calls) return empty
    // This is fine - we'll get the real value on full parse
    return ""
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
    // Tags can be an attribute with object value or a nested block
    if attr, ok := body.Attributes["tags"]; ok {
        if obj, ok := attr.Expr.(*hclsyntax.ObjectConsExpr); ok {
            entry.Tags = make(map[string]string)
            for _, item := range obj.Items {
                // Key and value should both be extractable
                key := extractStringLiteral(item.KeyExpr)
                val := extractStringLiteral(item.ValueExpr)
                if key != "" {
                    entry.Tags[intern.Intern(key)] = val
                }
            }
        }
    }
}

// extractReference extracts a "type.name" reference from an expression
func extractReference(expr hcl.Expression) string {
    // References like benchmark.child are ScopeTraversalExpr
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
```

### Byte Offset Support

```go
func (s *Scanner) ScanFileHCLWithOffsets(filePath string) error {
    content, err := os.ReadFile(filePath)
    if err != nil {
        return err
    }

    file, diags := hclsyntax.ParseConfig(content, filePath, hcl.InitialPos)
    if diags.HasErrors() {
        // Handle gracefully
    }

    body, ok := file.Body.(*hclsyntax.Body)
    if !ok {
        return nil
    }

    for _, block := range body.Blocks {
        if !indexedTypes[block.Type] || len(block.Labels) == 0 {
            continue
        }

        // HCL provides byte offsets in Range
        startOffset := block.TypeRange.Start.Byte
        endOffset := block.Body.EndRange.End.Byte

        entry := &IndexEntry{
            // ... other fields ...
            ByteOffset: int64(startOffset),
            ByteLength: endOffset - startOffset,
        }

        // ... rest of extraction ...
    }

    return nil
}
```

## Implementation Steps

1. [ ] Create `scanner_hcl.go` with HCL-based implementation
2. [ ] Implement `extractStringLiteral` for various expression types
3. [ ] Implement `extractReference` for traversal expressions
4. [ ] Add comprehensive tests comparing regex vs HCL output
5. [ ] Benchmark both implementations:
   ```bash
   go test -bench=BenchmarkScanner -benchmem ./internal/resourceindex/...
   ```
6. [ ] Test with edge cases that were problematic for regex:
   - Escaped quotes in strings
   - Single-line blocks
   - Single-line tags
   - Block comments
   - Heredoc content with HCL-like syntax
   - Nested heredocs
   - Template expressions
7. [ ] Compare correctness on real-world mods
8. [ ] Document performance characteristics
9. [ ] Decide: replace regex scanner or keep as fallback

## Acceptance Criteria

- [ ] HCL-based scanner produces identical output to regex scanner for valid HCL
- [ ] HCL-based scanner correctly handles all edge cases without custom code
- [ ] Performance is within acceptable bounds (< 3x slower than regex)
- [ ] Memory usage is acceptable for large mods
- [ ] Graceful handling of malformed HCL (partial results or clear errors)

## Performance Targets

| Metric | Regex Scanner | HCL Scanner Target |
|--------|---------------|-------------------|
| 1000 resources | ~10ms | < 30ms |
| Memory per resource | ~300 bytes | < 1KB |
| Large mod (10K resources) | ~100ms | < 300ms |

## Decision Criteria

**Adopt HCL scanner if:**
- Performance is within 3x of regex scanner
- All edge cases handled correctly without custom code
- Simpler codebase (delete regex patterns and state tracking)

**Keep regex scanner if:**
- HCL scanner is > 5x slower
- Memory usage is prohibitive for large mods
- Edge cases in real mods that HCL parser doesn't handle

**Hybrid approach if:**
- HCL scanner is 3-5x slower
- Cache parsed ASTs for unchanged files
- Use HCL for initial scan, regex for incremental updates

## Files to Create/Modify

- Create: `internal/resourceindex/scanner_hcl.go`
- Create: `internal/resourceindex/scanner_hcl_test.go`
- Modify: `internal/resourceindex/scanner_test.go` (add comparison tests)
- Potentially delete: Regex patterns and state tracking code in `scanner.go`

## Dependencies

Already in use:
- `github.com/hashicorp/hcl/v2`
- `github.com/hashicorp/hcl/v2/hclsyntax`
- `github.com/zclconf/go-cty/cty`

## Notes

- The `hclsyntax.ParseConfig` function is specifically designed for syntax-level parsing without evaluation
- Expression types like `TemplateExpr`, `LiteralValueExpr`, `ScopeTraversalExpr` can be inspected without evaluation
- Non-literal expressions (variables, functions) will return empty strings, which is fine - the real value comes on full parse
- The HCL parser's error recovery may allow partial extraction even from malformed files
