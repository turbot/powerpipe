package resourceindex

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestScannerRegressionScenario tests scanning with the exact pattern from the regression test
func TestScannerRegressionScenario(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "scanner_regression")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create main mod with 10 benchmarks using EXACT format from regression test
	modContent := `mod "test" { title = "Test" }
variable "common_tags" { default = { service = "AWS", category = "Compliance" } }
`
	for i := 0; i < 10; i++ {
		modContent += "benchmark \"bench_" + strconv.Itoa(i) + "\" {\n"
		modContent += "  title = \"Benchmark " + strconv.Itoa(i) + "\"\n"
		modContent += "  tags = merge(var.common_tags, { benchmark = \"true\", index = \"" + strconv.Itoa(i) + "\" })\n"
		modContent += "  control { sql = \"select 1\" }\n"
		modContent += "}\n"
	}

	filePath := filepath.Join(tmpDir, "mod.pp")
	err = os.WriteFile(filePath, []byte(modContent), 0600)
	require.NoError(t, err)

	// Scan using ScanDirectoryParallel like buildResourceIndex does
	scanner := NewScanner("test")
	scanner.SetModRoot(tmpDir)
	err = scanner.ScanDirectoryParallel(tmpDir, 0)
	require.NoError(t, err)

	// Check all benchmarks
	index := scanner.GetIndex()
	benchmarks := index.Benchmarks()
	require.Len(t, benchmarks, 10, "should have 10 benchmarks")

	// Check first benchmark in detail
	firstBenchmark := benchmarks[0]
	t.Logf("First benchmark: %s", firstBenchmark.Name)
	t.Logf("  TagsResolved: %v", firstBenchmark.TagsResolved)
	t.Logf("  Tags: %v (len=%d)", firstBenchmark.Tags, len(firstBenchmark.Tags))
	t.Logf("  UnresolvedRefs: %v", firstBenchmark.UnresolvedRefs)
	t.Logf("  NeedsResolution: %v", firstBenchmark.NeedsResolution())

	// Verify all benchmarks have unresolved tags
	for i, bench := range benchmarks {
		t.Logf("Benchmark %d: %s - TagsResolved=%v, Tags=%v, UnresolvedRefs=%v",
			i, bench.Name, bench.TagsResolved, bench.Tags, bench.UnresolvedRefs)

		assert.False(t, bench.TagsResolved,
			"Benchmark %s should NOT be marked as resolved (has merge())", bench.Name)
		assert.Contains(t, bench.UnresolvedRefs, "tags",
			"Benchmark %s should have 'tags' in UnresolvedRefs", bench.Name)
		assert.True(t, bench.NeedsResolution(),
			"Benchmark %s should need resolution", bench.Name)

		// Should extract partial tags from inline object
		assert.Equal(t, "true", bench.Tags["benchmark"],
			"Should extract benchmark tag from inline object")
		// Verify index tag exists and is numeric (exact value may vary due to ordering)
		assert.NotEmpty(t, bench.Tags["index"],
			"Should extract index tag from inline object")
	}
}
