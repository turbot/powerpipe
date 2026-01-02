# Task 14: Fix Scanner Limitations

## Objective

Fix known limitations in the regex-based scanner (`internal/resourceindex/scanner.go`) that were documented during edge case testing in Task 3.

## Context

The scanner uses fast regex-based parsing for performance. During edge case testing, we documented several limitations where the scanner produces different results than a full HCL parser would. While these edge cases are rare in practice, fixing them improves correctness and reduces potential index vs full-parse mismatches.

## Dependencies

- Task 3 (scanner edge case tests) - **COMPLETED**
- Tests already exist that document current behavior

## Priority

**LOW** - These are edge cases that rarely occur in real-world mod files. Fix if time permits or if real-world issues are reported.

## Limitations to Fix

### 1. Escaped Quotes in Attribute Values (LOW)

**Problem**: The regex `"([^"]*)"` stops at the first `"` character, even if escaped.

**Current behavior**:
```hcl
title = "Dashboard \"Quoted\" Title"
# Extracted as: "Dashboard \"  (truncated)
```

**Fix approach**: Update `attrStringRegex` to handle escaped characters:
```go
// Before
var attrStringRegex = regexp.MustCompile(`^\s*(\w+)\s*=\s*"([^"]*)"`)

// After - handle escaped chars
var attrStringRegex = regexp.MustCompile(`^\s*(\w+)\s*=\s*"((?:[^"\\]|\\.)*)"`)
```

**Test to update**: `TestScanner_TitleWithQuotes/escaped_double_quotes_truncate_value`

---

### 2. Single-Line Block Attributes (LOW)

**Problem**: When a block opens and closes on the same line, attributes aren't captured.

**Current behavior**:
```hcl
dashboard "d" { title = "Title" }
# Resource indexed, but title is empty
```

**Root cause**: `processBlockLine` is called after block start detection. When block closes on same line, depth becomes 0 and block is finalized before attributes are processed.

**Fix approach**: In `scanReader` and `scanReaderWithOffsets`, after detecting a block start with braces on the same line, extract and process the content between `{` and `}`:
```go
if blockStart := s.parseBlockStart(line); blockStart != nil {
    // ... existing code ...

    // If single-line block, extract content between braces
    if openBraces > 0 && closeBraces > 0 {
        braceStart := strings.Index(line, "{")
        braceEnd := strings.LastIndex(line, "}")
        if braceStart < braceEnd {
            innerContent := line[braceStart+1 : braceEnd]
            s.processBlockLine(innerContent, block)
        }
    }
}
```

**Test to update**: `TestScanner_SingleLineBlock`

---

### 3. Single-Line Tags (LOW)

**Problem**: Tags on a single line aren't parsed.

**Current behavior**:
```hcl
tags = { service = "aws" category = "security" }
# Tags map is empty
```

**Root cause**: Scanner looks for `tags` followed by `{` to start parsing, then expects tag entries on subsequent lines. When closing `}` is on same line, no entries are captured.

**Fix approach**: When detecting tags block start, check if it closes on same line and parse inline:
```go
if strings.HasPrefix(trimmed, "tags") && strings.Contains(line, "{") {
    // Check if single-line tags
    if strings.Contains(line, "}") {
        // Extract content between { and }
        start := strings.Index(line, "{")
        end := strings.LastIndex(line, "}")
        if start < end {
            tagsContent := line[start+1 : end]
            // Parse key = "value" pairs from tagsContent
            s.parseSingleLineTags(tagsContent, block)
        }
    } else {
        block.inTags = true
    }
    return
}
```

**Test to update**: `TestScanner_TagsSingleLine`

---

### 4. Block Comments (MEDIUM)

**Problem**: Content inside `/* ... */` block comments is still scanned, potentially matching resource patterns.

**Current behavior**:
```hcl
/* dashboard "commented_out" {
    title = "Should Not Parse"
} */
# Resource may be indexed depending on line structure
```

**Fix approach**: Add block comment state tracking:
```go
type blockState struct {
    // ... existing fields ...
    inBlockComment bool
}

func (s *Scanner) scanReader(r io.Reader, filePath string) error {
    // ... existing code ...
    inBlockComment := false

    for scanner.Scan() {
        line := scanner.Text()

        // Track block comments
        if strings.Contains(line, "/*") {
            inBlockComment = true
        }
        if inBlockComment {
            if strings.Contains(line, "*/") {
                inBlockComment = false
            }
            continue // Skip lines in block comments
        }

        // ... rest of processing ...
    }
}
```

**Note**: This simple approach doesn't handle `/* */` on same line or nested comments, but covers common cases.

**Test to update**: `TestScanner_BlockComments`

---

## Implementation Steps

1. [ ] Fix escaped quotes regex
2. [ ] Update `TestScanner_TitleWithQuotes` to expect correct behavior
3. [ ] Fix single-line block attribute parsing
4. [ ] Update `TestScanner_SingleLineBlock` to expect correct behavior
5. [ ] Fix single-line tags parsing
6. [ ] Update `TestScanner_TagsSingleLine` to expect correct behavior
7. [ ] Add block comment tracking
8. [ ] Update `TestScanner_BlockComments` to expect correct behavior
9. [ ] Run full test suite to ensure no regressions
10. [ ] Benchmark to ensure performance isn't significantly impacted

## Acceptance Criteria

- [ ] All four limitations fixed
- [ ] Updated tests pass with correct expected values
- [ ] No regressions in existing scanner tests
- [ ] Performance benchmark shows <10% slowdown (scanner is in hot path)
- [ ] Code is well-commented explaining the regex patterns

## Performance Considerations

The scanner is performance-critical - it runs on every file during mod loading. Any fixes should:

1. Avoid adding significant regex complexity
2. Minimize additional string allocations
3. Keep the hot path (normal multi-line blocks) fast
4. Only do extra work for edge cases when detected

Consider benchmarking before/after with:
```bash
go test -bench=. ./internal/resourceindex/...
```

## Output Files

- Modified: `internal/resourceindex/scanner.go`
- Modified: `internal/resourceindex/scanner_edge_test.go` (update expectations)
