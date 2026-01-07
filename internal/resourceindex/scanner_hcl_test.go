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
// HCL Scanner Basic Tests
// =============================================================================

func TestScannerHCL_BasicDashboard(t *testing.T) {
	content := `dashboard "basic" {
    title = "Basic Dashboard"
    description = "A basic test dashboard"
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.basic")
	require.True(t, ok)
	assert.Equal(t, "dashboard", entry.Type)
	assert.Equal(t, "basic", entry.ShortName)
	assert.Equal(t, "Basic Dashboard", entry.Title)
	assert.Equal(t, "A basic test dashboard", entry.Description)
	assert.Equal(t, "testmod", entry.ModName)
	assert.Equal(t, "mod.testmod", entry.ModFullName)
}

func TestScannerHCL_QueryWithSQL(t *testing.T) {
	content := `query "my_query" {
    title = "My Query"
    sql = "SELECT * FROM table"
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.query.my_query")
	require.True(t, ok)
	assert.True(t, entry.HasSQL)
}

func TestScannerHCL_ControlWithQueryRef(t *testing.T) {
	content := `control "my_control" {
    title = "My Control"
    query = query.shared_query
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.control.my_control")
	require.True(t, ok)
	assert.True(t, entry.HasSQL)
	assert.Equal(t, "testmod.query.shared_query", entry.QueryRef)
}

func TestScannerHCL_BenchmarkWithChildren(t *testing.T) {
	content := `benchmark "parent" {
    title = "Parent Benchmark"
    children = [
        benchmark.child1,
        control.ctrl1,
        control.ctrl2
    ]
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.benchmark.parent")
	require.True(t, ok)
	assert.Equal(t, "control", entry.BenchmarkType)
	require.Len(t, entry.ChildNames, 3)
	assert.Contains(t, entry.ChildNames, "testmod.benchmark.child1")
	assert.Contains(t, entry.ChildNames, "testmod.control.ctrl1")
	assert.Contains(t, entry.ChildNames, "testmod.control.ctrl2")
}

func TestScannerHCL_DetectionBenchmark(t *testing.T) {
	content := `detection_benchmark "detect" {
    title = "Detection Benchmark"
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.detection_benchmark.detect")
	require.True(t, ok)
	assert.Equal(t, "detection", entry.BenchmarkType)
}

func TestScannerHCL_Tags(t *testing.T) {
	content := `dashboard "tagged" {
    title = "Tagged Dashboard"
    tags = {
        service = "aws"
        category = "security"
    }
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.tagged")
	require.True(t, ok)
	assert.Equal(t, "aws", entry.Tags["service"])
	assert.Equal(t, "security", entry.Tags["category"])
}

// =============================================================================
// Edge Case Tests - HCL Parser Advantages
// =============================================================================

func TestScannerHCL_EscapedQuotes(t *testing.T) {
	content := `dashboard "escaped" {
    title = "Dashboard \"Quoted\" Title"
    description = "Contains \"nested\" quotes"
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.escaped")
	require.True(t, ok)
	// HCL parser properly handles escaped quotes
	assert.Equal(t, `Dashboard "Quoted" Title`, entry.Title)
	assert.Equal(t, `Contains "nested" quotes`, entry.Description)
}

func TestScannerHCL_SingleLineBlock(t *testing.T) {
	content := `dashboard "single" { title = "Single Line" }`

	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.single")
	require.True(t, ok)
	assert.Equal(t, "Single Line", entry.Title)
}

func TestScannerHCL_SingleLineTags(t *testing.T) {
	// Note: HCL requires commas between object items on a single line
	content := `dashboard "inline_tags" {
    title = "Inline Tags"
    tags = { service = "aws", category = "compliance" }
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.inline_tags")
	require.True(t, ok)
	assert.Equal(t, "aws", entry.Tags["service"])
	assert.Equal(t, "compliance", entry.Tags["category"])
}

func TestScannerHCL_BlockComments(t *testing.T) {
	content := `/* This is a block comment
dashboard "fake" {
    title = "Should Not Parse"
}
*/

dashboard "real" {
    title = "Real Dashboard"
}

/*
benchmark "also_fake" {
    title = "Also not parsed"
}
*/
`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	index := scanner.GetIndex()

	// Only the real dashboard should be indexed
	assert.Equal(t, 1, index.Count())

	_, ok := index.Get("testmod.dashboard.real")
	assert.True(t, ok)

	_, ok = index.Get("testmod.dashboard.fake")
	assert.False(t, ok)

	_, ok = index.Get("testmod.benchmark.also_fake")
	assert.False(t, ok)
}

func TestScannerHCL_Heredoc(t *testing.T) {
	content := `query "heredoc" {
    title = "Heredoc Query"
    sql = <<-EOQ
        SELECT
            "column_name",
            'string_value'
        FROM
            "schema"."table"
        WHERE
            name = 'test'
    EOQ
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.query.heredoc")
	require.True(t, ok)
	assert.True(t, entry.HasSQL)
	assert.Equal(t, "Heredoc Query", entry.Title)
}

func TestScannerHCL_HeredocWithFakeResources(t *testing.T) {
	// HCL parser correctly skips content inside heredocs
	content := `query "real" {
    title = "Real Query"
    sql = <<-EOQ
        -- This SQL contains text that looks like HCL
        SELECT 'dashboard "fake" {' as example
        FROM (
            SELECT 'benchmark "not_real" {' as text
        )
    EOQ
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	index := scanner.GetIndex()

	// Only the real query should be indexed
	assert.Equal(t, 1, index.Count())

	_, ok := index.Get("testmod.query.real")
	assert.True(t, ok)

	_, ok = index.Get("testmod.dashboard.fake")
	assert.False(t, ok)
}

func TestScannerHCL_Unicode(t *testing.T) {
	content := `dashboard "unicode" {
    title = "Dashboard 🚀 Rocket"
    description = "中文描述"
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.unicode")
	require.True(t, ok)
	assert.Equal(t, "Dashboard 🚀 Rocket", entry.Title)
	assert.Equal(t, "中文描述", entry.Description)
}

func TestScannerHCL_BracesInStrings(t *testing.T) {
	// HCL parser correctly handles braces inside strings
	content := `dashboard "braces" {
    title = "Title with {braces} inside"
    description = "More {nested {braces}}"
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.braces")
	require.True(t, ok)
	assert.Equal(t, "Title with {braces} inside", entry.Title)
	assert.Equal(t, "More {nested {braces}}", entry.Description)
}

func TestScannerHCL_CommentMarksInStrings(t *testing.T) {
	// HCL parser correctly handles comment markers inside strings
	content := `dashboard "comments" {
    title = "URL: http://example.com/* path */"
    description = "Contains /* and */ in string"
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.comments")
	require.True(t, ok)
	assert.Equal(t, "URL: http://example.com/* path */", entry.Title)
	assert.Equal(t, "Contains /* and */ in string", entry.Description)
}

// =============================================================================
// Byte Offset Tests
// =============================================================================

func TestScannerHCL_ByteOffsetAccuracy(t *testing.T) {
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
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0600))

	scanner := NewScanner("testmod")
	err := scanner.ScanFileHCLWithOffsets(filePath)
	require.NoError(t, err)

	fileContent, err := os.ReadFile(filePath)
	require.NoError(t, err)

	index := scanner.GetIndex()

	first, _ := index.Get("testmod.dashboard.first")
	require.NotNil(t, first)

	// Verify offset points to start of "dashboard" keyword
	assert.GreaterOrEqual(t, first.ByteOffset, int64(0))
	assert.Less(t, first.ByteOffset, int64(len(fileContent)))

	// Extract content at offset
	extracted := string(fileContent[first.ByteOffset : first.ByteOffset+int64(first.ByteLength)])
	assert.Contains(t, extracted, "dashboard")
	assert.Contains(t, extracted, "first")

	second, _ := index.Get("testmod.query.second")
	require.NotNil(t, second)
	assert.Greater(t, second.ByteOffset, first.ByteOffset)

	extracted = string(fileContent[second.ByteOffset : second.ByteOffset+int64(second.ByteLength)])
	assert.Contains(t, extracted, "query")
	assert.Contains(t, extracted, "second")
}

// =============================================================================
// All Resource Types Test
// =============================================================================

func TestScannerHCL_AllResourceTypes(t *testing.T) {
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
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	index := scanner.GetIndex()

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

// =============================================================================
// Error Handling Tests
// =============================================================================

func TestScannerHCL_MalformedHCL(t *testing.T) {
	// Partial HCL - missing closing brace
	content := `dashboard "partial" {
    title = "Partial"
`
	scanner := NewScanner("testmod")
	// Should not panic, may or may not extract partial data
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	// Error is acceptable but should not crash
	_ = err
}

func TestScannerHCL_EmptyFile(t *testing.T) {
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(""), "test.pp")
	require.NoError(t, err)
	assert.Equal(t, 0, scanner.GetIndex().Count())
}

func TestScannerHCL_CommentsOnly(t *testing.T) {
	content := `# Just comments
// More comments
/* Block comment */
`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)
	assert.Equal(t, 0, scanner.GetIndex().Count())
}

// =============================================================================
// Benchmark Tests
// =============================================================================

func BenchmarkScanner_HCL(b *testing.B) {
	content := generateBenchmarkContent(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scanner := NewScanner("testmod")
		_ = scanner.ScanBytesHCL(content, "bench.pp")
	}
}

func BenchmarkScanner_HCL_Large(b *testing.B) {
	content := generateBenchmarkContent(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scanner := NewScanner("testmod")
		_ = scanner.ScanBytesHCL(content, "bench.pp")
	}
}

func BenchmarkScanner_HCL_WithHeredocs(b *testing.B) {
	content := generateBenchmarkContentWithHeredocs(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scanner := NewScanner("testmod")
		_ = scanner.ScanBytesHCL(content, "bench.pp")
	}
}

func generateBenchmarkContent(count int) []byte {
	var buf bytes.Buffer
	for i := 0; i < count; i++ {
		buf.WriteString(`query "q`)
		buf.WriteString(strings.Repeat("0", 4))
		buf.WriteString(`" {
    title = "Query `)
		buf.WriteString(strings.Repeat("0", 4))
		buf.WriteString(`"
    description = "Description for query number `)
		buf.WriteString(strings.Repeat("0", 4))
		buf.WriteString(`"
    sql = "SELECT * FROM table_`)
		buf.WriteString(strings.Repeat("0", 4))
		buf.WriteString(`"
    tags = {
        service = "aws"
        category = "security"
    }
}

`)
	}
	return buf.Bytes()
}

func generateBenchmarkContentWithHeredocs(count int) []byte {
	var buf bytes.Buffer
	for i := 0; i < count; i++ {
		buf.WriteString(`query "q`)
		buf.WriteString(strings.Repeat("0", 4))
		buf.WriteString(`" {
    title = "Query with Heredoc"
    sql = <<-EOQ
        SELECT
            column1,
            column2,
            column3
        FROM
            schema.table
        WHERE
            status = 'active'
            AND created_at > now() - interval '7 days'
        ORDER BY
            created_at DESC
        LIMIT 100
    EOQ
}

`)
	}
	return buf.Bytes()
}

// =============================================================================
// Additional HCL-Specific Tests
// =============================================================================

func TestScannerHCL_LineNumbers(t *testing.T) {
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
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	dash, _ := scanner.GetIndex().Get("testmod.dashboard.line2")
	assert.Equal(t, 2, dash.StartLine, "StartLine mismatch for dashboard")

	query, _ := scanner.GetIndex().Get("testmod.query.line6")
	assert.Equal(t, 6, query.StartLine, "StartLine mismatch for query")
}

func TestScannerHCL_MultipleHeredocs(t *testing.T) {
	content := `query "first" {
    title = "First"
    sql = <<-EOQ
        SELECT 1
    EOQ
}

query "second" {
    title = "Second"
    sql = <<-EOT
        SELECT 2
    EOT
}
`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	index := scanner.GetIndex()
	assert.Equal(t, 2, index.Count())

	first, _ := index.Get("testmod.query.first")
	assert.True(t, first.HasSQL)
	assert.Equal(t, "First", first.Title)

	second, _ := index.Get("testmod.query.second")
	assert.True(t, second.HasSQL)
	assert.Equal(t, "Second", second.Title)
}

func TestScannerHCL_BenchmarkTypes(t *testing.T) {
	content := `
benchmark "control_bench" {
    title = "Control Benchmark"
}

detection_benchmark "detect_bench" {
    title = "Detection Benchmark"
}
`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	cb, _ := scanner.GetIndex().Get("testmod.benchmark.control_bench")
	assert.Equal(t, "control", cb.BenchmarkType)

	db, _ := scanner.GetIndex().Get("testmod.detection_benchmark.detect_bench")
	assert.Equal(t, "detection", db.BenchmarkType)
}

func TestScannerHCL_QueryRefs(t *testing.T) {
	content := `
control "with_ref" {
    title = "Control with Ref"
    query = query.shared
}

control "with_sql" {
    title = "Control with SQL"
    sql = "SELECT 'ok'"
}

query "shared" {
    sql = "SELECT 'ok'"
}
`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	ref, _ := scanner.GetIndex().Get("testmod.control.with_ref")
	assert.True(t, ref.HasSQL)
	assert.Equal(t, "testmod.query.shared", ref.QueryRef)

	sql, _ := scanner.GetIndex().Get("testmod.control.with_sql")
	assert.True(t, sql.HasSQL)
	assert.Empty(t, sql.QueryRef)
}

// =============================================================================
// Enhanced Extraction Tests - Resolution Tracking
// =============================================================================

func TestScannerHCL_TitleResolution_Literal(t *testing.T) {
	content := `dashboard "literal_title" {
    title = "My Dashboard Title"
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.literal_title")
	require.True(t, ok)
	assert.Equal(t, "My Dashboard Title", entry.Title)
	assert.True(t, entry.TitleResolved, "Literal title should be marked as resolved")
}

func TestScannerHCL_TitleResolution_Variable(t *testing.T) {
	content := `dashboard "var_title" {
    title = var.dashboard_title
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.var_title")
	require.True(t, ok)
	assert.Empty(t, entry.Title, "Variable title should be empty initially")
	assert.False(t, entry.TitleResolved, "Variable title should NOT be marked as resolved")
}

func TestScannerHCL_TitleResolution_Interpolation(t *testing.T) {
	content := `dashboard "interp_title" {
    title = "Dashboard for ${var.service}"
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.interp_title")
	require.True(t, ok)
	// Should extract the literal parts
	assert.Contains(t, entry.Title, "Dashboard for ")
	assert.False(t, entry.TitleResolved, "Interpolated title should NOT be marked as resolved")
}

func TestScannerHCL_TitleResolution_NoTitle(t *testing.T) {
	content := `dashboard "no_title" {
    description = "No title here"
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.no_title")
	require.True(t, ok)
	assert.Empty(t, entry.Title)
	assert.True(t, entry.TitleResolved, "Absent title should be marked as resolved (to empty)")
}

func TestScannerHCL_DescriptionResolution_Literal(t *testing.T) {
	content := `dashboard "literal_desc" {
    description = "A simple description"
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.literal_desc")
	require.True(t, ok)
	assert.Equal(t, "A simple description", entry.Description)
	assert.True(t, entry.DescriptionResolved, "Literal description should be resolved")
}

func TestScannerHCL_DescriptionResolution_Variable(t *testing.T) {
	content := `dashboard "var_desc" {
    description = var.common_description
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.var_desc")
	require.True(t, ok)
	assert.Empty(t, entry.Description)
	assert.False(t, entry.DescriptionResolved, "Variable description should NOT be resolved")
}

func TestScannerHCL_TagsResolution_LiteralTags(t *testing.T) {
	content := `dashboard "literal_tags" {
    title = "Dashboard"
    tags = {
        service  = "AWS"
        category = "Security"
        type     = "Dashboard"
    }
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.literal_tags")
	require.True(t, ok)
	assert.Equal(t, "AWS", entry.Tags["service"])
	assert.Equal(t, "Security", entry.Tags["category"])
	assert.Equal(t, "Dashboard", entry.Tags["type"])
	assert.True(t, entry.TagsResolved, "Literal tags should be marked as resolved")
	assert.Empty(t, entry.UnresolvedRefs, "No unresolved refs for literal tags")
}

func TestScannerHCL_TagsResolution_VariableTags(t *testing.T) {
	content := `dashboard "var_tags" {
    title = "Dashboard"
    tags = var.common_tags
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.var_tags")
	require.True(t, ok)
	assert.False(t, entry.TagsResolved, "Variable tags should NOT be marked as resolved")
	assert.Contains(t, entry.UnresolvedRefs, "tags", "Should list 'tags' as unresolved")
}

func TestScannerHCL_TagsResolution_MergeTags(t *testing.T) {
	content := `dashboard "merge_tags" {
    title = "Dashboard"
    tags = merge(var.common_tags, {
        service = "AWS"
        custom  = "value"
    })
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.merge_tags")
	require.True(t, ok)
	// Should extract inline literals even when merge is used
	assert.Equal(t, "AWS", entry.Tags["service"])
	assert.Equal(t, "value", entry.Tags["custom"])
	// But should mark as needing resolution
	assert.False(t, entry.TagsResolved, "merge() tags should NOT be marked as resolved")
	assert.Contains(t, entry.UnresolvedRefs, "tags")
}

func TestScannerHCL_TagsResolution_LocalReference(t *testing.T) {
	content := `dashboard "local_tags" {
    title = "Dashboard"
    tags = local.aws_compliance_common_tags
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.local_tags")
	require.True(t, ok)
	assert.False(t, entry.TagsResolved, "local. tags should NOT be marked as resolved")
	assert.Contains(t, entry.UnresolvedRefs, "tags")
}

func TestScannerHCL_TagsResolution_NoTags(t *testing.T) {
	content := `dashboard "no_tags" {
    title = "No Tags Dashboard"
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.no_tags")
	require.True(t, ok)
	assert.True(t, entry.TagsResolved, "Absent tags should be marked as resolved")
	assert.Empty(t, entry.Tags)
}

func TestScannerHCL_TagsResolution_MixedValues(t *testing.T) {
	content := `dashboard "mixed_tags" {
    title = "Dashboard"
    tags = {
        service  = "AWS"
        dynamic  = var.dynamic_tag
        static   = "StaticValue"
    }
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.mixed_tags")
	require.True(t, ok)
	assert.Equal(t, "AWS", entry.Tags["service"])
	assert.Equal(t, "StaticValue", entry.Tags["static"])
	assert.Empty(t, entry.Tags["dynamic"], "Variable tag value should be empty")
	assert.False(t, entry.TagsResolved, "Mixed tags with variables should NOT be resolved")
	assert.Contains(t, entry.UnresolvedRefs, "tag:dynamic")
}

func TestScannerHCL_EnhancedMetadata_Category(t *testing.T) {
	content := `dashboard "with_category" {
    title    = "Dashboard"
    category = "Compliance"
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.with_category")
	require.True(t, ok)
	assert.Equal(t, "Compliance", entry.Category)
}

func TestScannerHCL_EnhancedMetadata_Documentation(t *testing.T) {
	content := `control "with_docs" {
    title         = "Control"
    documentation = "https://docs.example.com/control"
    sql           = "SELECT 'ok'"
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.control.with_docs")
	require.True(t, ok)
	assert.Equal(t, "https://docs.example.com/control", entry.Documentation)
}

func TestScannerHCL_EnhancedMetadata_Width(t *testing.T) {
	content := `card "with_width" {
    title = "Card"
    width = 4
    sql   = "SELECT 1"
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.card.with_width")
	require.True(t, ok)
	require.NotNil(t, entry.Width)
	assert.Equal(t, 4, *entry.Width)
}

func TestScannerHCL_EnhancedMetadata_Display(t *testing.T) {
	content := `control "with_display" {
    title   = "Control"
    display = "none"
    sql     = "SELECT 'ok'"
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.control.with_display")
	require.True(t, ok)
	assert.Equal(t, "none", entry.Display)
}

func TestScannerHCL_NeedsResolution_Method(t *testing.T) {
	// Fully resolved entry
	content1 := `dashboard "resolved" {
    title       = "Resolved Dashboard"
    description = "A description"
    tags = {
        service = "AWS"
    }
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content1), "test.pp")
	require.NoError(t, err)

	resolved, _ := scanner.GetIndex().Get("testmod.dashboard.resolved")
	assert.False(t, resolved.NeedsResolution(), "Fully literal entry should not need resolution")
	assert.True(t, resolved.IsFullyResolved())

	// Entry needing resolution
	content2 := `dashboard "unresolved" {
    title = var.title
    tags = var.tags
}`
	scanner2 := NewScanner("testmod2")
	err = scanner2.ScanBytesHCL([]byte(content2), "test2.pp")
	require.NoError(t, err)

	unresolved, _ := scanner2.GetIndex().Get("testmod2.dashboard.unresolved")
	assert.True(t, unresolved.NeedsResolution(), "Entry with variables should need resolution")
	assert.False(t, unresolved.IsFullyResolved())
}

func TestScannerHCL_GetUnresolvedFields_Method(t *testing.T) {
	content := `dashboard "partial" {
    title       = var.title
    description = "Literal description"
    tags = var.common_tags
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, _ := scanner.GetIndex().Get("testmod.dashboard.partial")
	fields := entry.GetUnresolvedFields()

	assert.Contains(t, fields, "title")
	assert.NotContains(t, fields, "description") // Description is resolved
	assert.Contains(t, fields, "tags")
}

func TestScannerHCL_TagsBlock_LiteralValues(t *testing.T) {
	// Test tags as a block instead of attribute
	content := `dashboard "tags_block" {
    title = "Dashboard"

    tags {
        service  = "AWS"
        category = "Security"
    }
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.tags_block")
	require.True(t, ok)
	assert.Equal(t, "AWS", entry.Tags["service"])
	assert.Equal(t, "Security", entry.Tags["category"])
	assert.True(t, entry.TagsResolved)
}

func TestScannerHCL_TagsBlock_VariableValue(t *testing.T) {
	content := `dashboard "tags_block_var" {
    title = "Dashboard"

    tags {
        service  = "AWS"
        dynamic  = var.tag_value
    }
}`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytesHCL([]byte(content), "test.pp")
	require.NoError(t, err)

	entry, ok := scanner.GetIndex().Get("testmod.dashboard.tags_block_var")
	require.True(t, ok)
	assert.Equal(t, "AWS", entry.Tags["service"])
	assert.False(t, entry.TagsResolved)
	assert.Contains(t, entry.UnresolvedRefs, "tag:dynamic")
}
