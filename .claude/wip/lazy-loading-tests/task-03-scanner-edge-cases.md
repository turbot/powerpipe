# Task 3: Index/Scanner Edge Case Tests

## Objective

Write comprehensive tests for the resource index scanner (`internal/resourceindex/scanner.go`) focusing on edge cases that could cause index vs full-parse mismatches.

## Context

- The scanner uses regex-based line scanning for speed
- It must produce results consistent with full HCL parsing
- Mismatches cause bugs where UI shows different data than execution uses
- Recent bugs (mod_full_name, trunks) were in this area

## Dependencies

- Task 2 (test fixtures)
- Files to test: `internal/resourceindex/scanner.go`, `internal/resourceindex/index.go`

## Acceptance Criteria

- [ ] Add tests to `internal/resourceindex/scanner_test.go`
- [ ] Cover all edge cases listed below
- [ ] Each test should verify scanner output matches expected
- [ ] Tests should run in <30 seconds total
- [ ] Zero flaky tests

## Test Cases to Implement

### String Parsing Edge Cases
```go
// Test: Titles with special characters
func TestScanner_TitleWithQuotes(t *testing.T)
// e.g., title = "Dashboard \"Quoted\" Title"

// Test: Titles with escape sequences
func TestScanner_TitleWithEscapes(t *testing.T)
// e.g., title = "Line1\nLine2\tTabbed"

// Test: Unicode in titles and names
func TestScanner_UnicodeContent(t *testing.T)
// e.g., title = "Dashboard "

// Test: Empty strings
func TestScanner_EmptyStrings(t *testing.T)
// e.g., title = ""
```

### Heredoc Edge Cases
```go
// Test: Heredoc SQL with special content
func TestScanner_HeredocWithQuotes(t *testing.T)
// SQL containing single/double quotes

// Test: Heredoc that looks like resource definition
func TestScanner_HeredocLooksLikeResource(t *testing.T)
// SQL containing "dashboard \"fake\" {"

// Test: Very large heredoc
func TestScanner_LargeHeredoc(t *testing.T)
// 100KB SQL heredoc

// Test: Nested heredoc markers
func TestScanner_NestedHeredocMarkers(t *testing.T)
// Content containing EOF within heredoc
```

### Whitespace and Formatting
```go
// Test: Various indentation styles
func TestScanner_IndentationVariations(t *testing.T)
// Tabs, spaces, mixed

// Test: Windows line endings
func TestScanner_WindowsLineEndings(t *testing.T)
// \r\n line endings

// Test: Mixed line endings
func TestScanner_MixedLineEndings(t *testing.T)
// Mix of \n and \r\n

// Test: Trailing whitespace
func TestScanner_TrailingWhitespace(t *testing.T)
// Lines ending with spaces/tabs
```

### Comment Handling
```go
// Test: Comments containing resource-like text
func TestScanner_CommentsWithResourceSyntax(t *testing.T)
// # dashboard "fake" {

// Test: Block comments
func TestScanner_BlockComments(t *testing.T)
// /* dashboard "fake" { */

// Test: Inline comments after values
func TestScanner_InlineComments(t *testing.T)
// title = "Real" # not this
```

### Children Array Parsing
```go
// Test: Multi-line children arrays
func TestScanner_MultiLineChildren(t *testing.T)
// children = [
//   benchmark.a,
//   benchmark.b
// ]

// Test: Single-line children
func TestScanner_SingleLineChildren(t *testing.T)
// children = [benchmark.a, benchmark.b]

// Test: Children with comments
func TestScanner_ChildrenWithComments(t *testing.T)
// children = [
//   benchmark.a, # comment
//   benchmark.b
// ]

// Test: Empty children array
func TestScanner_EmptyChildren(t *testing.T)
// children = []
```

### Byte Offset Accuracy
```go
// Test: Byte offset points to correct line
func TestScanner_ByteOffsetAccuracy(t *testing.T)
// Verify seeking to offset lands at resource

// Test: Byte offset with unicode
func TestScanner_ByteOffsetWithUnicode(t *testing.T)
// Multi-byte characters before resource

// Test: Byte offset with varying line lengths
func TestScanner_ByteOffsetVaryingLines(t *testing.T)
// Mix of short and long lines
```

### Resource Type Detection
```go
// Test: All resource types detected
func TestScanner_AllResourceTypes(t *testing.T)
// dashboard, benchmark, control, query, card, chart, etc.

// Test: Detection benchmark type
func TestScanner_DetectionBenchmarkType(t *testing.T)
// detection_benchmark vs benchmark

// Test: Custom/unknown types
func TestScanner_UnknownResourceType(t *testing.T)
// custom_resource "name" {} - should skip
```

### Hierarchy Building
```go
// Test: Parent-child relationship setting
func TestScanner_ParentChildRelationships(t *testing.T)
// After scan, children have parent set

// Test: Top-level detection accuracy
func TestScanner_TopLevelDetection(t *testing.T)
// Only true top-level marked

// Test: Orphan resources
func TestScanner_OrphanResources(t *testing.T)
// Resources with no parent
```

### Index vs Parse Comparison
```go
// Test: Scanner output matches full HCL parse
func TestScanner_MatchesFullParse(t *testing.T)
// For each test fixture:
// 1. Scan with scanner
// 2. Parse with full parser
// 3. Compare metadata

// Test: Benchmark hierarchy matches
func TestScanner_BenchmarkHierarchyMatchesParser(t *testing.T)
// Verify children lists match

// Test: Reference extraction matches
func TestScanner_QueryReferencesMatchParser(t *testing.T)
// Control's query refs match
```

## Implementation Notes

- Use table-driven tests where appropriate
- Create helper function to compare scanner entry vs parsed resource
- Log detailed diffs on failure for debugging
- Consider fuzzing for string parsing

## Output Files

- `internal/resourceindex/scanner_edge_test.go` - New edge case tests
- Updates to existing `scanner_test.go` if consolidation helps
