package resourceindex

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// String Parsing Edge Cases
// =============================================================================

func TestScanner_TitleWithQuotes(t *testing.T) {
	// The HCL parser properly handles escaped characters and unescapes them.
	tests := []struct {
		name      string
		content   string
		wantName  string
		wantTitle string
	}{
		{
			name: "escaped double quotes are handled",
			content: `dashboard "quoted_title" {
    title = "Dashboard \"Quoted\" Title"
}`,
			wantName:  "testmod.dashboard.quoted_title",
			wantTitle: `Dashboard "Quoted" Title`, // HCL parser unescapes quotes
		},
		{
			name: "single quotes in title (literal)",
			content: `dashboard "single_quotes" {
    title = "Dashboard 'Single' Quotes"
}`,
			wantName:  "testmod.dashboard.single_quotes",
			wantTitle: "Dashboard 'Single' Quotes",
		},
		{
			name: "nested single and double quotes",
			content: `query "nested_quotes" {
    title = "He said 'Hello' to me"
    sql = "SELECT 1"
}`,
			wantName:  "testmod.query.nested_quotes",
			wantTitle: "He said 'Hello' to me",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewScanner("testmod")
			err := scanner.ScanBytes([]byte(tt.content), "test.pp")
			require.NoError(t, err)

			entry, ok := scanner.GetIndex().Get(tt.wantName)
			require.True(t, ok, "entry not found: %s", tt.wantName)
			assert.Equal(t, tt.wantTitle, entry.Title)
		})
	}
}

func TestScanner_TitleWithEscapes(t *testing.T) {
	// HCL parser unescapes \n and \t to actual newline and tab characters
	tests := []struct {
		name     string
		content  string
		wantName string
		wantTitle string
	}{
		{
			name: "newline escape",
			content: `dashboard "newline_title" {
    title = "Line1\nLine2"
}`,
			wantName:  "testmod.dashboard.newline_title",
			wantTitle: "Line1\nLine2", // HCL unescapes to actual newline
		},
		{
			name: "tab escape",
			content: `dashboard "tab_title" {
    title = "Col1\tCol2"
}`,
			wantName:  "testmod.dashboard.tab_title",
			wantTitle: "Col1\tCol2", // HCL unescapes to actual tab
		},
		{
			name: "multiple escapes",
			content: `query "multi_escape" {
    title = "Line1\nLine2\tTabbed"
    sql = "SELECT 1"
}`,
			wantName:  "testmod.query.multi_escape",
			wantTitle: "Line1\nLine2\tTabbed", // HCL unescapes both
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewScanner("testmod")
			err := scanner.ScanBytes([]byte(tt.content), "test.pp")
			require.NoError(t, err)

			entry, ok := scanner.GetIndex().Get(tt.wantName)
			require.True(t, ok, "entry not found: %s", tt.wantName)
			assert.Equal(t, tt.wantTitle, entry.Title)
		})
	}
}

func TestScanner_UnicodeContent(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		wantName    string
		wantTitle   string
		wantDesc    string
	}{
		{
			name: "emoji in title",
			content: `dashboard "emoji_dash" {
    title = "Dashboard 🚀 Rocket"
    description = "Has 🎉 emoji"
}`,
			wantName:  "testmod.dashboard.emoji_dash",
			wantTitle: "Dashboard 🚀 Rocket",
			wantDesc:  "Has 🎉 emoji",
		},
		{
			name: "chinese characters",
			content: `query "chinese" {
    title = "中文标题"
    description = "测试中文"
    sql = "SELECT 1"
}`,
			wantName:  "testmod.query.chinese",
			wantTitle: "中文标题",
			wantDesc:  "测试中文",
		},
		{
			name: "arabic characters",
			content: `query "arabic" {
    title = "مرحبا"
    sql = "SELECT 1"
}`,
			wantName:  "testmod.query.arabic",
			wantTitle: "مرحبا",
		},
		{
			name: "accented characters",
			content: `dashboard "accents" {
    title = "Café Résumé Naïve"
}`,
			wantName:  "testmod.dashboard.accents",
			wantTitle: "Café Résumé Naïve",
		},
		{
			name: "mixed scripts",
			content: `query "mixed_scripts" {
    title = "Hello 你好 مرحبا Привет"
    sql = "SELECT 1"
}`,
			wantName:  "testmod.query.mixed_scripts",
			wantTitle: "Hello 你好 مرحبا Привет",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewScanner("testmod")
			err := scanner.ScanBytes([]byte(tt.content), "test.pp")
			require.NoError(t, err)

			entry, ok := scanner.GetIndex().Get(tt.wantName)
			require.True(t, ok, "entry not found: %s", tt.wantName)
			assert.Equal(t, tt.wantTitle, entry.Title)
			if tt.wantDesc != "" {
				assert.Equal(t, tt.wantDesc, entry.Description)
			}
		})
	}
}

func TestScanner_EmptyStrings(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantName  string
		wantTitle string
		wantDesc  string
	}{
		{
			name: "empty title",
			content: `dashboard "empty_title" {
    title = ""
    description = "Has a description"
}`,
			wantName:  "testmod.dashboard.empty_title",
			wantTitle: "",
			wantDesc:  "Has a description",
		},
		{
			name: "empty description",
			content: `query "empty_desc" {
    title = "Has Title"
    description = ""
    sql = "SELECT 1"
}`,
			wantName:  "testmod.query.empty_desc",
			wantTitle: "Has Title",
			wantDesc:  "",
		},
		{
			name: "both empty",
			content: `control "both_empty" {
    title = ""
    description = ""
    sql = "SELECT 'ok'"
}`,
			wantName:  "testmod.control.both_empty",
			wantTitle: "",
			wantDesc:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewScanner("testmod")
			err := scanner.ScanBytes([]byte(tt.content), "test.pp")
			require.NoError(t, err)

			entry, ok := scanner.GetIndex().Get(tt.wantName)
			require.True(t, ok, "entry not found: %s", tt.wantName)
			assert.Equal(t, tt.wantTitle, entry.Title)
			assert.Equal(t, tt.wantDesc, entry.Description)
		})
	}
}

// =============================================================================
// Heredoc Edge Cases
// =============================================================================

func TestScanner_HeredocWithQuotes(t *testing.T) {
	content := `query "heredoc_quotes" {
    title = "Heredoc with Quotes"
    sql = <<-EOQ
        SELECT
            "column_name",
            'string_value',
            "table"."field"
        FROM
            "schema"."table"
        WHERE
            name = 'test'
    EOQ
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.query.heredoc_quotes")
	require.True(t, ok)
	assert.True(t, entry.HasSQL)
	assert.Equal(t, "Heredoc with Quotes", entry.Title)
}

func TestScanner_HeredocLooksLikeResource(t *testing.T) {
	// Heredoc content that looks like a resource definition should not be parsed
	content := `query "heredoc_fake_resource" {
    title = "Heredoc with Fake Resource"
    sql = <<-EOQ
        -- This SQL contains text that looks like HCL
        SELECT 'dashboard "fake" {' as example
        FROM (
            SELECT 'benchmark "not_real" {' as text
        )
        WHERE 'query "also_fake" { sql = "SELECT 1" }' IS NOT NULL
    EOQ
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	index := scanner.GetIndex()

	// Should only have the real query, not the fake ones inside the heredoc
	assert.Equal(t, 1, index.Count())

	entry, ok := index.Get("testmod.query.heredoc_fake_resource")
	require.True(t, ok)
	assert.True(t, entry.HasSQL)

	// Fake resources should not exist
	_, ok = index.Get("testmod.dashboard.fake")
	assert.False(t, ok)
	_, ok = index.Get("testmod.benchmark.not_real")
	assert.False(t, ok)
}

func TestScanner_HeredocUnquotedFakeResource(t *testing.T) {
	// Test heredoc with unquoted fake resource that looks exactly like HCL
	// The fake resource inside heredoc should NOT be parsed because the
	// blockStartRegex requires the pattern to start at the beginning of line
	// (with optional whitespace). The indentation in heredocs means the
	// fake resource line starts with spaces, making it less likely to match.
	content := `query "real" {
    title = "Real Query"
    sql = <<-EOQ
        -- Fake resource below (not in quotes)
        dashboard "fake" {
            title = "Should Not Parse"
        }
    EOQ
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	index := scanner.GetIndex()

	// The real query should be found
	_, ok := index.Get("testmod.query.real")
	assert.True(t, ok, "real query should be found")

	// The fake dashboard inside heredoc should NOT be found
	// The regex `^\s*(\w+)\s+"([^"]+)"` does match indented lines,
	// so this test verifies the current behavior
	_, hasFake := index.Get("testmod.dashboard.fake")
	assert.False(t, hasFake, "fake dashboard in heredoc should not be indexed")
}

func TestScanner_BlockCommentInsideString(t *testing.T) {
	// Test that block comment markers inside strings don't trigger comment handling.
	// The scanner checks for /* at the line level, but the attribute parsing happens
	// after block comment checking. Since /* appears in the middle of a string value,
	// it doesn't affect the block comment state (block comment check is done first,
	// and the line doesn't START with /*).
	content := `dashboard "url_test" {
    title = "URL: http://example.com/* path */"
    description = "Contains /* and */ in string"
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.url_test")
	require.True(t, ok, "dashboard should be found")
	// The title should be correctly parsed even with /* in the string
	assert.Equal(t, "URL: http://example.com/* path */", entry.Title)
	assert.Equal(t, "Contains /* and */ in string", entry.Description)
}

func TestScanner_LargeHeredoc(t *testing.T) {
	// Create a 100KB SQL heredoc
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString(`query "large_heredoc" {
    title = "Large Heredoc Query"
    sql = <<-EOQ
        SELECT `)

	// Generate many columns
	for i := 0; i < 1000; i++ {
		if i > 0 {
			sqlBuilder.WriteString(",\n            ")
		}
		sqlBuilder.WriteString("column_")
		sqlBuilder.WriteString(strings.Repeat("abcdefghij", 10)) // ~100 chars per column name
		sqlBuilder.WriteString(" AS c_")
		sqlBuilder.WriteString(string(rune('a' + (i % 26))))
	}

	sqlBuilder.WriteString(`
        FROM
            large_table
    EOQ
}`)

	content := sqlBuilder.String()
	require.Greater(t, len(content), 100*1024, "content should be > 100KB")

	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.query.large_heredoc")
	require.True(t, ok)
	assert.True(t, entry.HasSQL)
	assert.Equal(t, "Large Heredoc Query", entry.Title)
}

func TestScanner_NestedHeredocMarkers(t *testing.T) {
	// Content containing EOF-like strings within heredoc
	content := `query "nested_markers" {
    title = "Nested Markers"
    sql = <<-EOQ
        SELECT
            'EOQ is not the end' as note,
            'Neither is <<-EOQ' as fake_start,
            '    EOQ' as indented_marker
        FROM table
    EOQ
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.query.nested_markers")
	require.True(t, ok)
	assert.True(t, entry.HasSQL)
}

func TestScanner_MultipleHeredocsInFile(t *testing.T) {
	content := `query "first_heredoc" {
    title = "First"
    sql = <<-EOQ
        SELECT 1
    EOQ
}

query "second_heredoc" {
    title = "Second"
    sql = <<-EOT
        SELECT 2
    EOT
}

control "third_heredoc" {
    title = "Third"
    sql = <<-SQL
        SELECT 'ok' as status
    SQL
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	index := scanner.GetIndex()
	assert.Equal(t, 3, index.Count())

	first, _ := index.Get("testmod.query.first_heredoc")
	require.NotNil(t, first)
	assert.True(t, first.HasSQL)

	second, _ := index.Get("testmod.query.second_heredoc")
	require.NotNil(t, second)
	assert.True(t, second.HasSQL)

	third, _ := index.Get("testmod.control.third_heredoc")
	require.NotNil(t, third)
	assert.True(t, third.HasSQL)
}

// =============================================================================
// Whitespace and Formatting Edge Cases
// =============================================================================

func TestScanner_IndentationVariations(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantName string
	}{
		{
			name: "tab indentation",
			content: "dashboard \"tab_indent\" {\n\ttitle = \"Tab Indented\"\n}",
			wantName: "testmod.dashboard.tab_indent",
		},
		{
			name: "space indentation",
			content: "dashboard \"space_indent\" {\n    title = \"Space Indented\"\n}",
			wantName: "testmod.dashboard.space_indent",
		},
		{
			name: "mixed indentation",
			content: "dashboard \"mixed_indent\" {\n\t    title = \"Mixed Indented\"\n}",
			wantName: "testmod.dashboard.mixed_indent",
		},
		{
			name: "no indentation",
			content: "dashboard \"no_indent\" {\ntitle = \"No Indent\"\n}",
			wantName: "testmod.dashboard.no_indent",
		},
		{
			name: "deep indentation",
			content: "dashboard \"deep_indent\" {\n                title = \"Deep Indent\"\n}",
			wantName: "testmod.dashboard.deep_indent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewScanner("testmod")
			err := scanner.ScanBytes([]byte(tt.content), "test.pp")
			require.NoError(t, err)

			_, ok := scanner.GetIndex().Get(tt.wantName)
			assert.True(t, ok, "entry not found: %s", tt.wantName)
		})
	}
}

func TestScanner_WindowsLineEndings(t *testing.T) {
	// Content with \r\n line endings
	content := "dashboard \"windows\" {\r\n    title = \"Windows Line Endings\"\r\n    description = \"Uses CRLF\"\r\n}\r\n"

	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.windows")
	require.True(t, ok)
	assert.Equal(t, "Windows Line Endings", entry.Title)
	assert.Equal(t, "Uses CRLF", entry.Description)
}

func TestScanner_MixedLineEndings(t *testing.T) {
	// Mix of Unix (\n) and Windows (\r\n) line endings
	content := "dashboard \"mixed_endings\" {\n    title = \"Mixed Endings\"\r\n    description = \"Unix and Windows\"\n}\r\n"

	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.mixed_endings")
	require.True(t, ok)
	assert.Equal(t, "Mixed Endings", entry.Title)
}

func TestScanner_TrailingWhitespace(t *testing.T) {
	// Lines with trailing spaces
	content := "dashboard \"trailing\" {   \n    title = \"Trailing Spaces\"   \n    description = \"Has trailing whitespace\"  \t  \n}   "

	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.trailing")
	require.True(t, ok)
	assert.Equal(t, "Trailing Spaces", entry.Title)
}

func TestScanner_NoSpaceAroundEquals(t *testing.T) {
	content := `query "compact" {
title="Compact"
description="No spaces around equals"
sql="SELECT 1"
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.query.compact")
	require.True(t, ok)
	assert.Equal(t, "Compact", entry.Title)
	assert.Equal(t, "No spaces around equals", entry.Description)
}

func TestScanner_ExcessiveWhitespace(t *testing.T) {
	content := `dashboard    "extra_spaces"    {

    title     =     "Extra Spaces"

    description     =     "Lots of whitespace"

}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.extra_spaces")
	require.True(t, ok)
	assert.Equal(t, "Extra Spaces", entry.Title)
}

// =============================================================================
// Comment Handling Edge Cases
// =============================================================================

func TestScanner_CommentsWithResourceSyntax(t *testing.T) {
	content := `# dashboard "fake_in_comment" {
#     title = "Should Not Parse"
# }

dashboard "real" {
    title = "Real Dashboard"
}

// benchmark "another_fake" {
//     title = "Also should not parse"
// }`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	index := scanner.GetIndex()

	// Only the real dashboard should be found
	assert.Equal(t, 1, index.Count())

	entry, ok := index.Get("testmod.dashboard.real")
	require.True(t, ok)
	assert.Equal(t, "Real Dashboard", entry.Title)

	// Fake resources in comments should not exist
	_, ok = index.Get("testmod.dashboard.fake_in_comment")
	assert.False(t, ok)
	_, ok = index.Get("testmod.benchmark.another_fake")
	assert.False(t, ok)
}

func TestScanner_BlockComments(t *testing.T) {
	content := `/* dashboard "block_commented" {
    title = "Should Not Parse"
} */

dashboard "after_block" {
    title = "After Block Comment"
}

/*
benchmark "multi_line_block" {
    title = "Also not parsed"
    children = [control.fake]
}
*/`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	index := scanner.GetIndex()

	// Only the real dashboard should be indexed (not the commented-out ones)
	assert.Equal(t, 1, index.Count())

	entry, ok := index.Get("testmod.dashboard.after_block")
	require.True(t, ok)
	assert.Equal(t, "After Block Comment", entry.Title)

	// Commented-out resources should NOT be indexed
	_, ok = index.Get("testmod.dashboard.block_commented")
	assert.False(t, ok, "block-commented dashboard should not be indexed")

	_, ok = index.Get("testmod.benchmark.multi_line_block")
	assert.False(t, ok, "block-commented benchmark should not be indexed")
}

func TestScanner_InlineComments(t *testing.T) {
	content := `dashboard "inline_comments" {
    title = "Real Title" # not this
    description = "Real Description" // also not this
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.inline_comments")
	require.True(t, ok)
	// Scanner may or may not strip inline comments - document current behavior
	assert.Contains(t, entry.Title, "Real Title")
}

func TestScanner_CommentedAttribute(t *testing.T) {
	content := `dashboard "commented_attr" {
    # title = "Commented Out"
    title = "Active Title"
    # description = "Also Commented"
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.commented_attr")
	require.True(t, ok)
	assert.Equal(t, "Active Title", entry.Title)
}

// =============================================================================
// Children Array Parsing Edge Cases
// =============================================================================

func TestScanner_MultiLineChildren(t *testing.T) {
	content := `benchmark "multi_line_children" {
    title = "Multi-line Children"
    children = [
        benchmark.child1,
        benchmark.child2,
        control.ctrl1,
        control.ctrl2
    ]
}

benchmark "child1" { title = "Child 1" }
benchmark "child2" { title = "Child 2" }
control "ctrl1" { sql = "SELECT 'ok'" }
control "ctrl2" { sql = "SELECT 'ok'" }`

	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.benchmark.multi_line_children")
	require.True(t, ok)
	require.Len(t, entry.ChildNames, 4)
	assert.Contains(t, entry.ChildNames, "testmod.benchmark.child1")
	assert.Contains(t, entry.ChildNames, "testmod.benchmark.child2")
	assert.Contains(t, entry.ChildNames, "testmod.control.ctrl1")
	assert.Contains(t, entry.ChildNames, "testmod.control.ctrl2")
}

func TestScanner_SingleLineChildren(t *testing.T) {
	content := `benchmark "single_line" {
    title = "Single Line Children"
    children = [benchmark.a, benchmark.b, control.c]
}

benchmark "a" { title = "A" }
benchmark "b" { title = "B" }
control "c" { sql = "SELECT 'ok'" }`

	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.benchmark.single_line")
	require.True(t, ok)
	require.Len(t, entry.ChildNames, 3)
	assert.Contains(t, entry.ChildNames, "testmod.benchmark.a")
	assert.Contains(t, entry.ChildNames, "testmod.benchmark.b")
	assert.Contains(t, entry.ChildNames, "testmod.control.c")
}

func TestScanner_ChildrenWithComments(t *testing.T) {
	content := `benchmark "children_with_comments" {
    title = "Children with Comments"
    children = [
        benchmark.child1, # first child
        benchmark.child2, // second child
        # benchmark.commented_out,
        control.ctrl1
    ]
}

benchmark "child1" { title = "Child 1" }
benchmark "child2" { title = "Child 2" }
control "ctrl1" { sql = "SELECT 'ok'" }`

	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.benchmark.children_with_comments")
	require.True(t, ok)
	// Scanner should find the real children, behavior with comments may vary
	assert.GreaterOrEqual(t, len(entry.ChildNames), 3)
}

func TestScanner_EmptyChildren(t *testing.T) {
	content := `benchmark "empty_children" {
    title = "Empty Children Array"
    children = []
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.benchmark.empty_children")
	require.True(t, ok)
	assert.Empty(t, entry.ChildNames)
}

func TestScanner_ChildrenWithTrailingComma(t *testing.T) {
	content := `benchmark "trailing_comma" {
    title = "Trailing Comma"
    children = [
        benchmark.child1,
        benchmark.child2,
    ]
}

benchmark "child1" { title = "Child 1" }
benchmark "child2" { title = "Child 2" }`

	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.benchmark.trailing_comma")
	require.True(t, ok)
	assert.Len(t, entry.ChildNames, 2)
}

// =============================================================================
// Byte Offset Accuracy Tests
// =============================================================================

func TestScanner_ByteOffsetAccuracy(t *testing.T) {
	content := `dashboard "first" {
    title = "First"
}

query "second" {
    title = "Second"
    sql = "SELECT 1"
}
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.pp")
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))

	scanner := NewScanner("testmod")
	err := scanner.ScanFileWithOffsets(filePath)
	require.NoError(t, err)

	index := scanner.GetIndex()

	// Read file content
	fileContent, err := os.ReadFile(filePath)
	require.NoError(t, err)

	first, _ := index.Get("testmod.dashboard.first")
	require.NotNil(t, first)

	// Verify offset points to start of "dashboard" keyword
	assert.GreaterOrEqual(t, first.ByteOffset, int64(0))
	assert.Less(t, first.ByteOffset, int64(len(fileContent)))

	// Extract the content at the offset
	extracted := string(fileContent[first.ByteOffset:first.ByteOffset+int64(first.ByteLength)])
	assert.Contains(t, extracted, "dashboard")
	assert.Contains(t, extracted, "first")

	second, _ := index.Get("testmod.query.second")
	require.NotNil(t, second)
	assert.Greater(t, second.ByteOffset, first.ByteOffset)

	extracted = string(fileContent[second.ByteOffset:second.ByteOffset+int64(second.ByteLength)])
	assert.Contains(t, extracted, "query")
	assert.Contains(t, extracted, "second")
}

func TestScanner_ByteOffsetWithUnicode(t *testing.T) {
	// Unicode characters before the resource affect byte offset
	content := `# 中文注释 Unicode comment 🚀

dashboard "after_unicode" {
    title = "After Unicode"
}
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.pp")
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))

	scanner := NewScanner("testmod")
	err := scanner.ScanFileWithOffsets(filePath)
	require.NoError(t, err)

	fileContent, err := os.ReadFile(filePath)
	require.NoError(t, err)

	entry, _ := scanner.GetIndex().Get("testmod.dashboard.after_unicode")
	require.NotNil(t, entry)

	// Verify byte offset is correct even with unicode
	assert.GreaterOrEqual(t, entry.ByteOffset, int64(0))
	assert.Less(t, entry.ByteOffset+int64(entry.ByteLength), int64(len(fileContent))+1)

	// Seeking to offset should land at dashboard
	extracted := string(fileContent[entry.ByteOffset:entry.ByteOffset+int64(entry.ByteLength)])
	assert.Contains(t, extracted, "dashboard")
}

func TestScanner_ByteOffsetVaryingLines(t *testing.T) {
	// Mix of very short and very long lines
	content := `x
xx
xxx
xxxx
xxxxx

dashboard "after_varying" {
    title = "After Varying Lines"
}

` + strings.Repeat("y", 200) + `

query "at_end" {
    sql = "SELECT 1"
}
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.pp")
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))

	scanner := NewScanner("testmod")
	err := scanner.ScanFileWithOffsets(filePath)
	require.NoError(t, err)

	fileContent, err := os.ReadFile(filePath)
	require.NoError(t, err)

	index := scanner.GetIndex()

	dash, _ := index.Get("testmod.dashboard.after_varying")
	require.NotNil(t, dash)
	extracted := string(fileContent[dash.ByteOffset:dash.ByteOffset+int64(dash.ByteLength)])
	assert.Contains(t, extracted, "dashboard")

	query, _ := index.Get("testmod.query.at_end")
	require.NotNil(t, query)
	assert.Greater(t, query.ByteOffset, dash.ByteOffset)
}

// =============================================================================
// Resource Type Detection Tests
// =============================================================================

func TestScanner_AllResourceTypesEdgeCases(t *testing.T) {
	content := `
dashboard "d" { title = "D" }
benchmark "b" { title = "B" }
control "c" { sql = "SELECT 'ok'" }
query "q" { sql = "SELECT 1" }
card "card" { sql = "SELECT 1" }
chart "chart" { sql = "SELECT 1" }
container "cont" { title = "Container" }
flow "f" { title = "Flow" }
graph "g" { title = "Graph" }
hierarchy "h" { title = "Hierarchy" }
image "img" { title = "Image" }
input "inp" { title = "Input" }
node "n" { title = "Node" }
edge "e" { title = "Edge" }
table "tbl" { sql = "SELECT 1" }
text "txt" { value = "Text" }
category "cat" { title = "Category" }
detection "det" { sql = "SELECT 1" }
detection_benchmark "db" { title = "Detection Benchmark" }
variable "var" { default = "x" }
with "w" { sql = "SELECT 1" }
`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	index := scanner.GetIndex()

	// All indexed types should be present
	expectedTypes := []string{
		"dashboard", "benchmark", "control", "query", "card", "chart",
		"container", "flow", "graph", "hierarchy", "image", "input",
		"node", "edge", "table", "text", "category", "detection",
		"detection_benchmark", "variable", "with",
	}

	for _, typ := range expectedTypes {
		entries := index.GetByType(typ)
		assert.NotEmpty(t, entries, "no entries for type: %s", typ)
	}
}

func TestScanner_DetectionBenchmarkType(t *testing.T) {
	content := `
benchmark "control_bench" {
    title = "Control Benchmark"
}

detection_benchmark "detect_bench" {
    title = "Detection Benchmark"
}
`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	index := scanner.GetIndex()

	cb, _ := index.Get("testmod.benchmark.control_bench")
	require.NotNil(t, cb)
	assert.Equal(t, "benchmark", cb.Type)
	assert.Equal(t, "control", cb.BenchmarkType)

	db, _ := index.Get("testmod.detection_benchmark.detect_bench")
	require.NotNil(t, db)
	assert.Equal(t, "detection_benchmark", db.Type)
	assert.Equal(t, "detection", db.BenchmarkType)
}

func TestScanner_UnknownResourceType(t *testing.T) {
	content := `
custom_resource "should_skip" {
    title = "Not Indexed"
}

unknown_type "also_skip" {
    value = "ignored"
}

dashboard "real" {
    title = "Real Dashboard"
}
`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	index := scanner.GetIndex()

	// Only the known type should be indexed
	assert.Equal(t, 1, index.Count())

	_, ok := index.Get("testmod.dashboard.real")
	assert.True(t, ok)

	_, ok = index.Get("testmod.custom_resource.should_skip")
	assert.False(t, ok)

	_, ok = index.Get("testmod.unknown_type.also_skip")
	assert.False(t, ok)
}

// =============================================================================
// Hierarchy Building Tests
// =============================================================================

func TestScanner_ParentChildRelationships(t *testing.T) {
	content := `
benchmark "parent" {
    title = "Parent"
    children = [
        benchmark.child,
        control.ctrl
    ]
}

benchmark "child" {
    title = "Child"
    children = [
        control.grandchild_ctrl
    ]
}

control "ctrl" { sql = "SELECT 'ok'" }
control "grandchild_ctrl" { sql = "SELECT 'ok'" }
`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)
	scanner.SetParentNames()

	index := scanner.GetIndex()

	child, _ := index.Get("testmod.benchmark.child")
	require.NotNil(t, child)
	assert.Equal(t, "testmod.benchmark.parent", child.ParentName)

	ctrl, _ := index.Get("testmod.control.ctrl")
	require.NotNil(t, ctrl)
	assert.Equal(t, "testmod.benchmark.parent", ctrl.ParentName)

	grandchild, _ := index.Get("testmod.control.grandchild_ctrl")
	require.NotNil(t, grandchild)
	assert.Equal(t, "testmod.benchmark.child", grandchild.ParentName)
}

func TestScanner_TopLevelDetection(t *testing.T) {
	content := `
benchmark "top1" {
    title = "Top 1"
    children = [benchmark.nested]
}

benchmark "top2" {
    title = "Top 2"
}

benchmark "nested" {
    title = "Nested"
}

dashboard "dash1" {
    title = "Dashboard 1"
}

dashboard "dash2" {
    title = "Dashboard 2"
}

control "standalone" {
    sql = "SELECT 'ok'"
}
`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)
	scanner.MarkTopLevelResources()

	index := scanner.GetIndex()

	// Top-level benchmarks
	top1, _ := index.Get("testmod.benchmark.top1")
	assert.True(t, top1.IsTopLevel)

	top2, _ := index.Get("testmod.benchmark.top2")
	assert.True(t, top2.IsTopLevel)

	// Nested benchmark should not be top-level
	nested, _ := index.Get("testmod.benchmark.nested")
	assert.False(t, nested.IsTopLevel)

	// Dashboards are always top-level
	dash1, _ := index.Get("testmod.dashboard.dash1")
	assert.True(t, dash1.IsTopLevel)

	dash2, _ := index.Get("testmod.dashboard.dash2")
	assert.True(t, dash2.IsTopLevel)

	// Controls are not marked as top-level (only dashboards and benchmarks)
	ctrl, _ := index.Get("testmod.control.standalone")
	assert.False(t, ctrl.IsTopLevel)
}

func TestScanner_OrphanResources(t *testing.T) {
	// Resources that are not children of anything
	content := `
benchmark "root" {
    title = "Root"
    children = [control.included]
}

control "included" { sql = "SELECT 'ok'" }
control "orphan1" { sql = "SELECT 'ok'" }
control "orphan2" { sql = "SELECT 'ok'" }
query "standalone" { sql = "SELECT 1" }
`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)
	scanner.SetParentNames()

	index := scanner.GetIndex()

	included, _ := index.Get("testmod.control.included")
	assert.Equal(t, "testmod.benchmark.root", included.ParentName)

	orphan1, _ := index.Get("testmod.control.orphan1")
	assert.Empty(t, orphan1.ParentName)

	orphan2, _ := index.Get("testmod.control.orphan2")
	assert.Empty(t, orphan2.ParentName)

	standalone, _ := index.Get("testmod.query.standalone")
	assert.Empty(t, standalone.ParentName)
}

// =============================================================================
// ModFullName and ModRoot Tests
// =============================================================================

func TestScanner_ModFullName(t *testing.T) {
	content := `dashboard "test" {
    title = "Test"
}`
	scanner := NewScanner("mymod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("mymod.dashboard.test")
	require.True(t, ok)

	assert.Equal(t, "mymod", entry.ModName)
	assert.Equal(t, "mod.mymod", entry.ModFullName)
}

func TestScanner_ModRoot(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.pp")
	content := `dashboard "test" { title = "Test" }`
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))

	scanner := NewScanner("mymod")
	scanner.SetModRoot(tmpDir)
	err := scanner.ScanFile(filePath)
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("mymod.dashboard.test")
	require.True(t, ok)
	assert.Equal(t, tmpDir, entry.ModRoot)
}

// =============================================================================
// Single Line Block Tests
// =============================================================================

func TestScanner_SingleLineBlock(t *testing.T) {
	// Single-line blocks are detected by the scanner and attributes on the
	// same line are now captured. The scanner extracts content between { and }
	// for single-line blocks and processes attributes from that content.
	tests := []struct {
		name        string
		content     string
		wantName    string
		wantIndexed bool
		wantTitle   string
	}{
		{
			name:        "simple single line - title captured",
			content:     `dashboard "single" { title = "Single Line" }`,
			wantName:    "testmod.dashboard.single",
			wantIndexed: true,
			wantTitle:   "Single Line",
		},
		{
			name:        "multi-line - title captured",
			content:     "dashboard \"multiline\" {\n    title = \"Multi Line\"\n}",
			wantName:    "testmod.dashboard.multiline",
			wantIndexed: true,
			wantTitle:   "Multi Line",
		},
		{
			name:        "single line with multiple attributes",
			content:     `query "compact" { title = "Compact" sql = "SELECT 1" }`,
			wantName:    "testmod.query.compact",
			wantIndexed: true,
			wantTitle:   "Compact",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewScanner("testmod")
			err := scanner.ScanBytes([]byte(tt.content), "test.pp")
			require.NoError(t, err)

			entry, ok := scanner.GetIndex().Get(tt.wantName)
			if tt.wantIndexed {
				require.True(t, ok, "entry not found: %s", tt.wantName)
				assert.Equal(t, tt.wantTitle, entry.Title)
			} else {
				assert.False(t, ok)
			}
		})
	}
}

// =============================================================================
// Query Reference Tests
// =============================================================================

func TestScanner_QueryReferenceEdgeCases(t *testing.T) {
	content := `
control "with_query_ref" {
    title = "Control with Query Reference"
    query = query.shared_query
}

control "with_inline_sql" {
    title = "Control with Inline SQL"
    sql = "SELECT 'ok' as status"
}

query "shared_query" {
    sql = "SELECT 'ok' as status"
}
`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	index := scanner.GetIndex()

	withRef, _ := index.Get("testmod.control.with_query_ref")
	require.NotNil(t, withRef)
	assert.True(t, withRef.HasSQL)
	assert.Equal(t, "testmod.query.shared_query", withRef.QueryRef)

	withInline, _ := index.Get("testmod.control.with_inline_sql")
	require.NotNil(t, withInline)
	assert.True(t, withInline.HasSQL)
	assert.Empty(t, withInline.QueryRef)
}

// =============================================================================
// Tags Tests
// =============================================================================

func TestScanner_TagsMultiLine(t *testing.T) {
	content := `
dashboard "tagged" {
    title = "Tagged Dashboard"
    tags = {
        service = "aws"
        category = "security"
        region = "us-east-1"
    }
}
`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, _ := scanner.GetIndex().Get("testmod.dashboard.tagged")
	require.NotNil(t, entry)
	assert.Equal(t, "aws", entry.Tags["service"])
	assert.Equal(t, "security", entry.Tags["category"])
	assert.Equal(t, "us-east-1", entry.Tags["region"])
}

func TestScanner_TagsSingleLine(t *testing.T) {
	// Single-line tags are now supported by the scanner.
	// Note: HCL requires commas between object items on a single line.
	content := `dashboard "inline_tags" {
    title = "Inline Tags"
    tags = { service = "aws", category = "compliance" }
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.inline_tags")
	require.True(t, ok, "dashboard not found")
	assert.Equal(t, "Inline Tags", entry.Title)
	// Single-line tags are now parsed correctly
	assert.Equal(t, "aws", entry.Tags["service"])
	assert.Equal(t, "compliance", entry.Tags["category"])
}

func TestScanner_EmptyTags(t *testing.T) {
	content := `dashboard "empty_tags" {
    title = "Empty Tags"
    tags = {}
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, _ := scanner.GetIndex().Get("testmod.dashboard.empty_tags")
	require.NotNil(t, entry)
	// Empty tags map or nil are both acceptable
}

// =============================================================================
// Concurrent Safety Tests
// =============================================================================

func TestScanner_ConcurrentScanning(t *testing.T) {
	tmpDir := t.TempDir()

	// Create multiple files
	for i := 0; i < 20; i++ {
		content := []byte(`
dashboard "dash_` + string(rune('a'+i)) + `" {
    title = "Dashboard ` + string(rune('A'+i)) + `"
}
query "query_` + string(rune('a'+i)) + `" {
    sql = "SELECT ` + string(rune('0'+i%10)) + `"
}
`)
		filePath := filepath.Join(tmpDir, "file_"+string(rune('a'+i))+".pp")
		require.NoError(t, os.WriteFile(filePath, content, 0644))
	}

	scanner := NewScanner("testmod")
	err := scanner.ScanDirectoryParallel(tmpDir, 8)
	require.NoError(t, err)

	index := scanner.GetIndex()
	assert.Equal(t, 40, index.Count()) // 20 dashboards + 20 queries
}

// =============================================================================
// Edge Cases from Test Fixtures
// =============================================================================

func TestScanner_FixtureEdgeCases(t *testing.T) {
	fixtureDir := "../../internal/testdata/mods/lazy-loading-tests/edge-cases"

	if _, err := os.Stat(fixtureDir); os.IsNotExist(err) {
		t.Skip("test fixtures not found")
	}

	scanner := NewScanner("edge_cases")
	err := scanner.ScanDirectory(fixtureDir)
	require.NoError(t, err)

	index := scanner.GetIndex()

	// Verify some expected entries exist from fixtures
	assert.Greater(t, index.Count(), 0)

	// Check unicode fixture entries
	if entry, ok := index.Get("edge_cases.query.query_with_emoji_title"); ok {
		assert.Contains(t, entry.Title, "🚀")
	}

	// Check special chars fixture entries
	if entry, ok := index.Get("edge_cases.query.query_with_underscore_name"); ok {
		assert.Equal(t, "Query With Underscore", entry.Title)
	}
}

// =============================================================================
// Brace Counting Edge Cases
// =============================================================================

func TestScanner_BraceCountingInStrings(t *testing.T) {
	content := `dashboard "braces_in_string" {
    title = "Title with {braces} inside"
    description = "More {nested {braces}}"
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.braces_in_string")
	require.True(t, ok)
	assert.Equal(t, "Title with {braces} inside", entry.Title)
}

func TestScanner_BraceCountingInSQL(t *testing.T) {
	content := `query "sql_with_braces" {
    title = "SQL with Braces"
    sql = "SELECT jsonb_build_object('key', value) FROM (SELECT 1 as value) t"
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.query.sql_with_braces")
	require.True(t, ok)
	assert.True(t, entry.HasSQL)
}

// =============================================================================
// File and Directory Operations
// =============================================================================

func TestScanner_NonExistentFile(t *testing.T) {
	scanner := NewScanner("testmod")
	err := scanner.ScanFile("/nonexistent/path/file.pp")
	assert.Error(t, err)
}

func TestScanner_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	scanner := NewScanner("testmod")
	err := scanner.ScanDirectory(tmpDir)
	require.NoError(t, err)

	assert.Equal(t, 0, scanner.GetIndex().Count())
}

func TestScanner_OnlyNonPPFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create non-.pp files
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "readme.md"), []byte("# README"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "config.json"), []byte("{}"), 0644))

	scanner := NewScanner("testmod")
	err := scanner.ScanDirectory(tmpDir)
	require.NoError(t, err)

	assert.Equal(t, 0, scanner.GetIndex().Count())
}

// =============================================================================
// Line Number Tests
// =============================================================================

func TestScanner_LineNumberAccuracy(t *testing.T) {
	content := `
dashboard "line2" {
    title = "Starts at line 2"
}

query "line6" {
    title = "Starts at line 6"
    sql = "SELECT 1"
}
`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	index := scanner.GetIndex()

	dash, _ := index.Get("testmod.dashboard.line2")
	require.NotNil(t, dash)
	assert.Equal(t, 2, dash.StartLine)
	assert.Equal(t, 4, dash.EndLine)

	query, _ := index.Get("testmod.query.line6")
	require.NotNil(t, query)
	assert.Equal(t, 6, query.StartLine)
	assert.Equal(t, 9, query.EndLine)
}

// =============================================================================
// Stress Tests
// =============================================================================

func TestScanner_ManyResourcesInOneFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping stress test in short mode")
	}

	var buf bytes.Buffer
	for i := 0; i < 1000; i++ {
		buf.WriteString(`query "q`)
		buf.WriteString(strings.Repeat("0", 4-len(string(rune('0'+i%10)))))
		buf.WriteString(string(rune('0' + i/1000%10)))
		buf.WriteString(string(rune('0' + i/100%10)))
		buf.WriteString(string(rune('0' + i/10%10)))
		buf.WriteString(string(rune('0' + i%10)))
		buf.WriteString(`" { sql = "SELECT `)
		buf.WriteString(string(rune('0' + i%10)))
		buf.WriteString(`" }
`)
	}

	scanner := NewScanner("testmod")
	err := scanner.ScanBytes(buf.Bytes(), "stress.pp")
	require.NoError(t, err)

	assert.Equal(t, 1000, scanner.GetIndex().Count())
}
