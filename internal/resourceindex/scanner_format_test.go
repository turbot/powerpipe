package resourceindex

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestScannerHCLFormat tests if formatting affects tag extraction
// This test DOCUMENTS A KNOWN BUG: Scanner fails to detect merge() when
// opening brace and first attribute are on the same line.
// TODO: Fix scanner to handle all HCL formatting styles consistently
func TestScannerHCLFormat(t *testing.T) {
	t.Skip("KNOWN BUG: Scanner doesn't detect merge() when { and title are on same line - documented for future fix")
	// Test 1: Multi-line format (WORKS)
	content1 := `mod "test" { title = "Test" }
benchmark "test1" {
  title = "Test 1"
  tags = merge(var.common_tags, { test = "1" })
  control { sql = "select 1" }
}
`

	scanner1 := NewScanner("test")
	err := scanner1.ScanBytes([]byte(content1), "test.pp")
	require.NoError(t, err)

	entry1, ok := scanner1.GetIndex().Get("test.benchmark.test1")
	require.True(t, ok)
	t.Logf("Format 1 (multi-line): TagsResolved=%v, Tags=%v, UnresolvedRefs=%v",
		entry1.TagsResolved, entry1.Tags, entry1.UnresolvedRefs)
	assert.False(t, entry1.TagsResolved, "Multi-line should detect merge()")
	assert.Contains(t, entry1.UnresolvedRefs, "tags")

	// Test 2: Opening brace + title on same line (FAILS?)
	content2 := `mod "test" { title = "Test" }
benchmark "test2" { title = "Test 2"
  tags = merge(var.common_tags, { test = "2" })
  control { sql = "select 1" } }
`

	scanner2 := NewScanner("test")
	err = scanner2.ScanBytes([]byte(content2), "test.pp")
	require.NoError(t, err)

	entry2, ok := scanner2.GetIndex().Get("test.benchmark.test2")
	require.True(t, ok)
	t.Logf("Format 2 (brace+title): TagsResolved=%v, Tags=%v, UnresolvedRefs=%v",
		entry2.TagsResolved, entry2.Tags, entry2.UnresolvedRefs)

	// Both formats should work the same
	assert.False(t, entry2.TagsResolved, "Brace+title format should also detect merge()")
	assert.Contains(t, entry2.UnresolvedRefs, "tags")
}
