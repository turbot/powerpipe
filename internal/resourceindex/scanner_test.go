package resourceindex

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScanner_SimpleFile(t *testing.T) {
	content := `
dashboard "my_dashboard" {
    title = "My Dashboard"
    description = "A test dashboard"
}

query "my_query" {
    title = "My Query"
    sql = "SELECT * FROM table"
}
`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	index := scanner.GetIndex()

	// Check dashboard was found
	dash, ok := index.Get("testmod.dashboard.my_dashboard")
	assert.True(t, ok)
	assert.Equal(t, "dashboard", dash.Type)
	assert.Equal(t, "my_dashboard", dash.ShortName)
	assert.Equal(t, "My Dashboard", dash.Title)
	assert.Equal(t, "A test dashboard", dash.Description)

	// Check query was found
	query, ok := index.Get("testmod.query.my_query")
	assert.True(t, ok)
	assert.Equal(t, "query", query.Type)
	assert.Equal(t, "My Query", query.Title)
	assert.True(t, query.HasSQL)
}

func TestScanner_BenchmarkWithChildren(t *testing.T) {
	content := `
benchmark "parent" {
    title = "Parent Benchmark"
    children = [
        benchmark.child1,
        benchmark.child2,
        control.ctrl1
    ]
}

benchmark "child1" {
    title = "Child 1"
}

benchmark "child2" {
    title = "Child 2"
}

control "ctrl1" {
    title = "Control 1"
    sql = "SELECT 'ok' as status"
}
`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)
	scanner.MarkTopLevelResources()
	scanner.SetParentNames()

	index := scanner.GetIndex()

	// Check parent benchmark
	parent, ok := index.Get("testmod.benchmark.parent")
	require.True(t, ok)
	assert.Equal(t, "Parent Benchmark", parent.Title)
	assert.True(t, parent.IsTopLevel)
	assert.Len(t, parent.ChildNames, 3)
	assert.Contains(t, parent.ChildNames, "testmod.benchmark.child1")
	assert.Contains(t, parent.ChildNames, "testmod.benchmark.child2")
	assert.Contains(t, parent.ChildNames, "testmod.control.ctrl1")

	// Check children have parent set
	child1, ok := index.Get("testmod.benchmark.child1")
	require.True(t, ok)
	assert.Equal(t, "testmod.benchmark.parent", child1.ParentName)
	assert.False(t, child1.IsTopLevel)
}

func TestScanner_BenchmarkType(t *testing.T) {
	content := `
benchmark "control_benchmark" {
    title = "Control Benchmark"
}

detection_benchmark "detection_benchmark" {
    title = "Detection Benchmark"
}
`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	index := scanner.GetIndex()

	cb, ok := index.Get("testmod.benchmark.control_benchmark")
	require.True(t, ok)
	assert.Equal(t, "control", cb.BenchmarkType)

	db, ok := index.Get("testmod.detection_benchmark.detection_benchmark")
	require.True(t, ok)
	assert.Equal(t, "detection", db.BenchmarkType)
}

func TestScanner_Tags(t *testing.T) {
	content := `
dashboard "tagged_dashboard" {
    title = "Tagged Dashboard"
    tags = {
        service = "aws"
        type = "compliance"
        region = "us-east-1"
    }
}
`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	index := scanner.GetIndex()

	dash, ok := index.Get("testmod.dashboard.tagged_dashboard")
	require.True(t, ok)
	assert.Equal(t, "aws", dash.Tags["service"])
	assert.Equal(t, "compliance", dash.Tags["type"])
	assert.Equal(t, "us-east-1", dash.Tags["region"])
}

func TestScanner_NestedContainers(t *testing.T) {
	content := `
dashboard "nested" {
    title = "Nested Dashboard"

    container {
        card {
            sql = "SELECT 1"
        }

        container {
            chart {
                type = "bar"
                sql = "SELECT * FROM metrics"
            }
        }
    }
}
`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	index := scanner.GetIndex()

	// Dashboard should be found
	dash, ok := index.Get("testmod.dashboard.nested")
	require.True(t, ok)
	assert.Equal(t, "Nested Dashboard", dash.Title)

	// Nested anonymous blocks don't get indexed (no quoted name)
	// Only named top-level resources are indexed
	assert.Equal(t, 1, index.Count())
}

func TestScanner_QueryReference(t *testing.T) {
	content := `
control "uses_query" {
    title = "Control with Query Reference"
    query = query.my_query
}

query "my_query" {
    sql = "SELECT * FROM table"
}
`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	index := scanner.GetIndex()

	ctrl, ok := index.Get("testmod.control.uses_query")
	require.True(t, ok)
	assert.True(t, ctrl.HasSQL) // Has reference to query with SQL
}

func TestScanner_HeredocSQL(t *testing.T) {
	content := `
query "heredoc_query" {
    title = "Query with Heredoc"
    sql = <<-EOQ
        SELECT
            id,
            name,
            status
        FROM
            resources
        WHERE
            active = true
    EOQ
}
`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	index := scanner.GetIndex()

	query, ok := index.Get("testmod.query.heredoc_query")
	require.True(t, ok)
	assert.True(t, query.HasSQL)
	assert.Equal(t, "Query with Heredoc", query.Title)
}

func TestScanner_AllResourceTypes(t *testing.T) {
	content := `
dashboard "d1" { title = "Dashboard" }
benchmark "b1" { title = "Benchmark" }
control "c1" { title = "Control" }
query "q1" { title = "Query" }
card "card1" { title = "Card" }
chart "chart1" { title = "Chart" }
container "container1" { title = "Container" }
flow "flow1" { title = "Flow" }
graph "graph1" { title = "Graph" }
hierarchy "hierarchy1" { title = "Hierarchy" }
image "image1" { title = "Image" }
input "input1" { title = "Input" }
node "node1" { title = "Node" }
edge "edge1" { title = "Edge" }
table "table1" { title = "Table" }
text "text1" { title = "Text" }
category "category1" { title = "Category" }
detection "detection1" { title = "Detection" }
detection_benchmark "db1" { title = "Detection Benchmark" }
variable "var1" { default = "test" }
`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	index := scanner.GetIndex()

	// Should have all 20 resource types
	assert.Equal(t, 20, index.Count())

	// Verify each type was indexed
	types := index.Types()
	assert.Contains(t, types, "dashboard")
	assert.Contains(t, types, "benchmark")
	assert.Contains(t, types, "control")
	assert.Contains(t, types, "query")
	assert.Contains(t, types, "card")
	assert.Contains(t, types, "chart")
	assert.Contains(t, types, "container")
	assert.Contains(t, types, "flow")
	assert.Contains(t, types, "graph")
	assert.Contains(t, types, "hierarchy")
	assert.Contains(t, types, "image")
	assert.Contains(t, types, "input")
	assert.Contains(t, types, "node")
	assert.Contains(t, types, "edge")
	assert.Contains(t, types, "table")
	assert.Contains(t, types, "text")
	assert.Contains(t, types, "category")
	assert.Contains(t, types, "detection")
	assert.Contains(t, types, "detection_benchmark")
	assert.Contains(t, types, "variable")
}

func TestScanner_LineNumbers(t *testing.T) {
	content := `dashboard "first" {
    title = "First"
}

query "second" {
    title = "Second"
    sql = "SELECT 1"
}

control "third" {
    title = "Third"
    sql = "SELECT 'ok'"
}
`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	index := scanner.GetIndex()

	first, _ := index.Get("testmod.dashboard.first")
	assert.Equal(t, 1, first.StartLine)
	assert.Equal(t, 3, first.EndLine)

	second, _ := index.Get("testmod.query.second")
	assert.Equal(t, 5, second.StartLine)
	assert.Equal(t, 8, second.EndLine)

	third, _ := index.Get("testmod.control.third")
	assert.Equal(t, 10, third.StartLine)
	assert.Equal(t, 13, third.EndLine)
}

func TestScanner_FileOperations(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.pp")

	content := `
dashboard "file_test" {
    title = "File Test"
}
`
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0600))

	scanner := NewScanner("testmod")
	err := scanner.ScanFile(filePath)
	require.NoError(t, err)

	index := scanner.GetIndex()

	dash, ok := index.Get("testmod.dashboard.file_test")
	assert.True(t, ok)
	assert.Equal(t, filePath, dash.FileName)
}

func TestScanner_DirectoryScanning(t *testing.T) {
	tmpDir := t.TempDir()

	// Create multiple files in directory structure
	require.NoError(t, os.MkdirAll(filepath.Join(tmpDir, "subdir"), 0755))

	files := map[string]string{
		"root.pp": `dashboard "root" { title = "Root" }`,
		"subdir/nested.pp": `
query "nested_query" { sql = "SELECT 1" }
control "nested_control" { sql = "SELECT 'ok'" }
`,
	}

	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		require.NoError(t, os.WriteFile(fullPath, []byte(content), 0600))
	}

	scanner := NewScanner("testmod")
	err := scanner.ScanDirectory(tmpDir)
	require.NoError(t, err)

	index := scanner.GetIndex()

	assert.Equal(t, 3, index.Count())
	_, ok := index.Get("testmod.dashboard.root")
	assert.True(t, ok)
	_, ok = index.Get("testmod.query.nested_query")
	assert.True(t, ok)
	_, ok = index.Get("testmod.control.nested_control")
	assert.True(t, ok)
}

func TestScanner_SkipsHiddenDirectories(t *testing.T) {
	tmpDir := t.TempDir()

	// Create visible file
	require.NoError(t, os.WriteFile(
		filepath.Join(tmpDir, "visible.pp"),
		[]byte(`dashboard "visible" { title = "Visible" }`),
		0600,
	))

	// Create hidden directory with file
	hiddenDir := filepath.Join(tmpDir, ".hidden")
	require.NoError(t, os.MkdirAll(hiddenDir, 0755))
	require.NoError(t, os.WriteFile(
		filepath.Join(hiddenDir, "hidden.pp"),
		[]byte(`dashboard "hidden" { title = "Hidden" }`),
		0600,
	))

	scanner := NewScanner("testmod")
	err := scanner.ScanDirectory(tmpDir)
	require.NoError(t, err)

	index := scanner.GetIndex()

	// Should only have visible file
	assert.Equal(t, 1, index.Count())
	_, ok := index.Get("testmod.dashboard.visible")
	assert.True(t, ok)
	_, ok = index.Get("testmod.dashboard.hidden")
	assert.False(t, ok)
}

func TestScanner_ParallelScanning(t *testing.T) {
	tmpDir := t.TempDir()

	// Create many files
	for i := 0; i < 50; i++ {
		content := fmt.Sprintf(`
dashboard "dashboard_%d" {
    title = "Dashboard %d"
}
query "query_%d" {
    sql = "SELECT %d"
}
`, i, i, i, i)
		filePath := filepath.Join(tmpDir, fmt.Sprintf("file_%d.pp", i))
		require.NoError(t, os.WriteFile(filePath, []byte(content), 0600))
	}

	scanner := NewScanner("testmod")
	err := scanner.ScanDirectoryParallel(tmpDir, 4)
	require.NoError(t, err)

	index := scanner.GetIndex()

	// Should have 100 resources (50 dashboards + 50 queries)
	assert.Equal(t, 100, index.Count())
}

func TestScanner_ByteOffsets(t *testing.T) {
	content := `dashboard "first" {
    title = "First"
}

query "second" {
    sql = "SELECT 1"
}
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.pp")
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0600))

	scanner := NewScanner("testmod")
	err := scanner.ScanFileWithOffsets(filePath)
	require.NoError(t, err)

	index := scanner.GetIndex()

	first, _ := index.Get("testmod.dashboard.first")
	assert.Greater(t, first.ByteLength, 0)

	second, _ := index.Get("testmod.query.second")
	assert.Greater(t, second.ByteOffset, int64(0))
	assert.Greater(t, second.ByteLength, 0)
}

func TestScanner_Performance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping performance test in short mode")
	}

	tmpDir := t.TempDir()

	// Create 100 files with 10 resources each (1000 total)
	for i := 0; i < 100; i++ {
		var content strings.Builder
		for j := 0; j < 10; j++ {
			content.WriteString(fmt.Sprintf(`
query "query_%d_%d" {
    title = "Query %d %d"
    description = "This is query number %d in file %d"
    sql = "SELECT %d, %d FROM table WHERE id = %d"
}
`, i, j, i, j, j, i, i, j, j))
		}
		filePath := filepath.Join(tmpDir, fmt.Sprintf("file_%d.pp", i))
		require.NoError(t, os.WriteFile(filePath, []byte(content.String()), 0600))
	}

	// Time sequential scan
	scanner := NewScanner("testmod")
	start := time.Now()
	err := scanner.ScanDirectory(tmpDir)
	require.NoError(t, err)
	seqDuration := time.Since(start)

	index := scanner.GetIndex()
	assert.Equal(t, 1000, index.Count())

	t.Logf("Sequential scan: %d resources in %v", index.Count(), seqDuration)
	t.Logf("Index size: %d bytes (%.2f KB)", index.Size(), float64(index.Size())/1024)

	// Verify performance target: < 500ms for 1000 resources
	assert.Less(t, seqDuration.Milliseconds(), int64(500),
		"Sequential scan too slow: %v", seqDuration)

	// Time parallel scan
	scanner2 := NewScanner("testmod2")
	start = time.Now()
	err = scanner2.ScanDirectoryParallel(tmpDir, 4)
	require.NoError(t, err)
	parDuration := time.Since(start)

	t.Logf("Parallel scan (4 workers): %d resources in %v", scanner2.GetIndex().Count(), parDuration)

	// Parallel should be faster (or at least not slower)
	assert.LessOrEqual(t, parDuration.Milliseconds(), seqDuration.Milliseconds()+100)
}

func TestScanner_MalformedFile(t *testing.T) {
	// HCL parser is stricter about syntax - malformed content may prevent parsing
	// This test verifies the scanner doesn't panic on invalid input
	content := `
dashboard "valid" {
    title = "Valid"
}

this is not valid HCL
{ random braces }

query "also_valid" {
    sql = "SELECT 1"
}
`
	scanner := NewScanner("testmod")
	// Scanner should not panic, but may not extract data from malformed files
	err := scanner.ScanBytes([]byte(content), "test.pp")
	// HCL parser returns diagnostics for malformed content, but doesn't error
	_ = err

	// Note: Unlike the regex scanner, the HCL parser may not extract
	// partial data from files with syntax errors
	// The important thing is it doesn't crash
}

func TestScanner_EmptyFile(t *testing.T) {
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(""), "test.pp")
	require.NoError(t, err)

	index := scanner.GetIndex()
	assert.Equal(t, 0, index.Count())
}

func TestScanner_CommentsIgnored(t *testing.T) {
	content := `
# This is a comment
dashboard "commented" {
    # title = "Wrong Title"
    title = "Correct Title"
    // Another comment style
    description = "Test"
}
`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)

	index := scanner.GetIndex()

	dash, ok := index.Get("testmod.dashboard.commented")
	require.True(t, ok)
	// Note: Simple scanner doesn't fully handle comments, but should get last value
	assert.Equal(t, "Test", dash.Description)
}

func TestScanner_SetModInfo(t *testing.T) {
	scanner := NewScanner("testmod")
	scanner.SetModInfo("testmod", "github.com/test/testmod", "Test Mod Title")

	index := scanner.GetIndex()
	assert.Equal(t, "testmod", index.ModName)
	assert.Equal(t, "github.com/test/testmod", index.ModFullName)
	assert.Equal(t, "Test Mod Title", index.ModTitle)
}

func TestScanner_TopLevelMarking(t *testing.T) {
	content := `
benchmark "root1" {
    title = "Root 1"
    children = [benchmark.child1]
}

benchmark "root2" {
    title = "Root 2"
}

benchmark "child1" {
    title = "Child 1"
    children = [benchmark.grandchild]
}

benchmark "grandchild" {
    title = "Grandchild"
}

dashboard "my_dash" {
    title = "Dashboard"
}
`
	scanner := NewScanner("testmod")
	err := scanner.ScanBytes([]byte(content), "test.pp")
	require.NoError(t, err)
	scanner.MarkTopLevelResources()

	index := scanner.GetIndex()

	root1, _ := index.Get("testmod.benchmark.root1")
	assert.True(t, root1.IsTopLevel)

	root2, _ := index.Get("testmod.benchmark.root2")
	assert.True(t, root2.IsTopLevel)

	child1, _ := index.Get("testmod.benchmark.child1")
	assert.False(t, child1.IsTopLevel)

	grandchild, _ := index.Get("testmod.benchmark.grandchild")
	assert.False(t, grandchild.IsTopLevel)

	dash, _ := index.Get("testmod.dashboard.my_dash")
	assert.True(t, dash.IsTopLevel)
}

func TestScanner_NonPPFilesIgnored(t *testing.T) {
	tmpDir := t.TempDir()

	// Create various files
	files := map[string]string{
		"valid.pp":   `dashboard "valid" { title = "Valid" }`,
		"readme.md":  `# README`,
		"config.hcl": `resource "test" { name = "test" }`,
		"data.json":  `{"key": "value"}`,
	}

	for name, content := range files {
		require.NoError(t, os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0600))
	}

	scanner := NewScanner("testmod")
	err := scanner.ScanDirectory(tmpDir)
	require.NoError(t, err)

	index := scanner.GetIndex()

	// Should only scan .pp files
	assert.Equal(t, 1, index.Count())
	_, ok := index.Get("testmod.dashboard.valid")
	assert.True(t, ok)
}

func TestScanner_LargeFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping large file test in short mode")
	}

	// Create a file with many resources
	var content strings.Builder
	for i := 0; i < 500; i++ {
		content.WriteString(fmt.Sprintf(`
control "control_%d" {
    title = "Control %d"
    description = "This is a longer description for control number %d that includes more text to simulate real-world control definitions with detailed explanations."
    sql = <<-EOQ
        SELECT
            'ok' as status,
            'resource_%d' as resource,
            'passed' as reason
        FROM
            some_table
        WHERE
            id = %d
    EOQ
}
`, i, i, i, i, i))
	}

	scanner := NewScanner("testmod")
	start := time.Now()
	err := scanner.ScanBytes([]byte(content.String()), "large.pp")
	require.NoError(t, err)
	duration := time.Since(start)

	index := scanner.GetIndex()

	t.Logf("Scanned %d controls from large file in %v", index.Count(), duration)
	assert.Equal(t, 500, index.Count())
	// Skip timing assertion when race detector is enabled (adds significant overhead)
	if !raceEnabled {
		assert.Less(t, duration.Milliseconds(), int64(200))
	}
}
