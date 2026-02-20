package resourceindex

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestScannerDiskRead tests that scanning from disk works correctly
func TestScannerDiskRead(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "scanner_disk_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Write test content to disk
	content := `mod "test" { title = "Test" }
benchmark "disk_test" {
  title = "Disk Test"
  tags = merge(var.common_tags, { source = "disk" })
  control { sql = "select 1" }
}
`
	filePath := filepath.Join(tmpDir, "test.pp")
	err = os.WriteFile(filePath, []byte(content), 0600)
	require.NoError(t, err)

	// Scan from disk using ScanFile
	scanner := NewScanner("test")
	scanner.SetModRoot(tmpDir)
	err = scanner.ScanFile(filePath)
	require.NoError(t, err)

	// Verify the benchmark was detected with unresolved tags
	entry, ok := scanner.GetIndex().Get("test.benchmark.disk_test")
	require.True(t, ok, "Should find benchmark in index")

	t.Logf("Scanned from disk: TagsResolved=%v, Tags=%v, UnresolvedRefs=%v",
		entry.TagsResolved, entry.Tags, entry.UnresolvedRefs)

	assert.False(t, entry.TagsResolved, "Should detect merge() from disk file")
	assert.Contains(t, entry.UnresolvedRefs, "tags", "Should mark tags as unresolved")
	assert.Equal(t, "disk", entry.Tags["source"], "Should extract inline tag value")
}
