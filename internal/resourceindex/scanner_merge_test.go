package resourceindex

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestScannerMergeDetection tests that the scanner correctly detects merge() calls in tags
func TestScannerMergeDetection(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		expectResolved bool
		expectRefs     []string
	}{
		{
			name: "merge with var reference",
			content: `mod "test" { title = "Test" }
benchmark "test1" {
  title = "Test 1"
  tags = merge(var.common_tags, { extra = "value" })
  control { sql = "select 1" }
}`,
			expectResolved: false,
			expectRefs:     []string{"tags"},
		},
		{
			name: "literal tags object",
			content: `mod "test" { title = "Test" }
benchmark "test2" {
  title = "Test 2"
  tags = { service = "aws", env = "prod" }
  control { sql = "select 1" }
}`,
			expectResolved: true,
			expectRefs:     nil,
		},
		{
			name: "var reference without merge",
			content: `mod "test" { title = "Test" }
benchmark "test3" {
  title = "Test 3"
  tags = var.common_tags
  control { sql = "select 1" }
}`,
			expectResolved: false,
			expectRefs:     []string{"tags"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewScanner("test")
			err := scanner.ScanBytes([]byte(tt.content), "test.pp")
			require.NoError(t, err)

			// Find the benchmark entry
			var entry *IndexEntry
			for _, e := range scanner.GetIndex().List() {
				if e.Type == "benchmark" {
					entry = e
					break
				}
			}
			require.NotNil(t, entry, "Should find benchmark entry")

			t.Logf("%s: TagsResolved=%v, UnresolvedRefs=%v",
				tt.name, entry.TagsResolved, entry.UnresolvedRefs)

			assert.Equal(t, tt.expectResolved, entry.TagsResolved,
				"TagsResolved should be %v", tt.expectResolved)

			if tt.expectRefs != nil {
				for _, ref := range tt.expectRefs {
					assert.Contains(t, entry.UnresolvedRefs, ref,
						"Should have unresolved ref: %s", ref)
				}
			}
		})
	}
}
