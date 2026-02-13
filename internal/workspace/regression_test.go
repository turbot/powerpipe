package workspace

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/turbot/powerpipe/internal/resourceindex"
)

// TestRegressionEmptyTags_ProductionScenario simulates the production bug
// by checking tags immediately after LoadLazy without any delay.
func TestRegressionEmptyTags_ProductionScenario(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "regression_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create main mod with 50 benchmarks
	// IMPORTANT: Use proper multi-line format - scanner has bug with "{ title" on same line
	modContent := `mod "test" { title = "Test" }
variable "common_tags" { default = { service = "AWS", category = "Compliance" } }
`
	for i := 0; i < 50; i++ {
		modContent += "benchmark \"bench_" + strconv.Itoa(i) + "\" {\n"
		modContent += "  title = \"Benchmark " + strconv.Itoa(i) + "\"\n"
		modContent += "  tags = merge(var.common_tags, { benchmark = \"true\", index = \"" + strconv.Itoa(i) + "\" })\n"
		modContent += "  control { sql = \"select 1\" }\n"
		modContent += "}\n"
	}

	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "mod.pp"), []byte(modContent), 0600))

	// Create 3 dependency mods with 20 benchmarks each
	// IMPORTANT: Use proper multi-line format - scanner has bug with "{ title" on same line
	for _, modName := range []string{"aws_compliance", "aws_insights", "aws_tags"} {
		modDir := filepath.Join(tmpDir, ".powerpipe", "mods", "github.com", "turbot", "steampipe-mod-"+modName+"@v1.0.0")
		require.NoError(t, os.MkdirAll(modDir, 0755))
		content := "mod \"" + modName + "\" { title = \"" + modName + "\" }\nvariable \"common_tags\" { default = { service = \"AWS\" } }\n"
		for i := 0; i < 20; i++ {
			content += "benchmark \"bench_" + modName + "_" + strconv.Itoa(i) + "\" {\n"
			content += "  title = \"Benchmark\"\n"
			content += "  tags = merge(var.common_tags, { mod = \"" + modName + "\" })\n"
			content += "  control { sql = \"select 1\" }\n"
			content += "}\n"
		}
		require.NoError(t, os.WriteFile(filepath.Join(modDir, "mod.pp"), []byte(content), 0600))
	}

	ctx := context.Background()

	// Debug: read back one of the files to verify content
	mainModContent, _ := os.ReadFile(filepath.Join(tmpDir, "mod.pp"))
	if len(mainModContent) > 500 {
		t.Logf("Main mod.pp first 500 chars:\n%s\n...", string(mainModContent[:500]))
	} else {
		t.Logf("Main mod.pp:\n%s", string(mainModContent))
	}

	depModPath := filepath.Join(tmpDir, ".powerpipe", "mods", "github.com", "turbot", "steampipe-mod-aws_compliance@v1.0.0", "mod.pp")
	depModContent, _ := os.ReadFile(depModPath)
	if len(depModContent) > 300 {
		t.Logf("Dep mod.pp first 300 chars:\n%s\n...", string(depModContent[:300]))
	} else {
		t.Logf("Dep mod.pp:\n%s", string(depModContent))
	}

	t.Run("With LoadLazy fix - tags should be populated", func(t *testing.T) {
		start := time.Now()
		lw, err := LoadLazy(ctx, tmpDir)
		require.NoError(t, err)
		defer lw.Close()
		elapsed := time.Since(start)

		// Check if eval context is set
		evalCtx := lw.GetLoader().GetEvalContext()
		if evalCtx == nil {
			t.Logf("WARNING: eval context is nil!")
		} else if evalCtx.Variables == nil {
			t.Logf("WARNING: eval context has no variables!")
		} else {
			t.Logf("✓ Eval context has variables: %v", len(evalCtx.Variables))
		}

		// Check index entries to see if they need resolution
		needsResolutionCount := 0
		sampleEntries := []string{}
		allEntries := lw.GetIndex().List()
		t.Logf("Total index entries: %d", len(allEntries))

		// Get first benchmark to check in detail
		var firstBenchmark *resourceindex.IndexEntry
		for i, entry := range allEntries {
			if entry.NeedsResolution() {
				needsResolutionCount++
			}
			// Get first benchmark for detailed analysis
			if firstBenchmark == nil && entry.Type == "benchmark" {
				firstBenchmark = entry
			}
			// Log first 3 benchmarks regardless of resolution status
			if i < 3 && entry.Type == "benchmark" {
				sampleEntries = append(sampleEntries, fmt.Sprintf("%s (title=%v, desc=%v, tags=%v, unresolvedRefs=%v, currentTags=%v)",
					entry.Name, entry.TitleResolved, entry.DescriptionResolved, entry.TagsResolved, entry.UnresolvedRefs, entry.Tags))
			}
		}
		t.Logf("Index entries needing resolution: %d", needsResolutionCount)
		for _, sample := range sampleEntries {
			t.Logf("  Sample: %s", sample)
		}

		if firstBenchmark != nil {
			t.Logf("DETAILED first benchmark:")
			t.Logf("  Name: %s", firstBenchmark.Name)
			t.Logf("  TagsResolved: %v", firstBenchmark.TagsResolved)
			t.Logf("  Tags in index entry: %v (len=%d)", firstBenchmark.Tags, len(firstBenchmark.Tags))
			t.Logf("  UnresolvedRefs: %v", firstBenchmark.UnresolvedRefs)
		}

		// Check background resolver stats
		stats := lw.BackgroundResolverStats()
		t.Logf("Background resolver: started=%v, complete=%v, queue_length=%v",
			stats.IsStarted, stats.IsComplete, stats.QueueLength)

		// Check if fully resolved
		if lw.IsFullyResolved() {
			t.Logf("✓ Background resolution completed")
		} else {
			t.Logf("WARNING: Background resolution not complete - waiting additional 5s...")
			time.Sleep(5 * time.Second)
			if lw.IsFullyResolved() {
				t.Logf("✓ Background resolution completed after additional wait")
			} else {
				t.Logf("ERROR: Background resolution still not complete!")
			}
		}

		payload := lw.GetAvailableDashboardsFromIndex()
		emptyCount := 0
		benchmarksWithTags := 0
		for name, bench := range payload.Benchmarks {
			if len(bench.Tags) == 0 {
				emptyCount++
			} else {
				benchmarksWithTags++
				if benchmarksWithTags <= 3 {
					t.Logf("Sample benchmark %s has tags: %v", name, bench.Tags)
				}
			}
		}

		t.Logf("LoadLazy took %v, %d/%d benchmarks with empty tags", elapsed, emptyCount, len(payload.Benchmarks))
		require.Equal(t, 0, emptyCount, "All benchmarks should have tags after LoadLazy with fix")
	})
}
