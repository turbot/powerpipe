package resourceindex

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndex_AddAndGet(t *testing.T) {
	idx := NewResourceIndex()

	entry := &IndexEntry{
		Type:      "dashboard",
		Name:      "mymod.dashboard.test",
		ShortName: "test",
		Title:     "Test Dashboard",
		FileName:  "/path/to/mod/dashboards.pp",
		StartLine: 10,
		EndLine:   50,
	}

	idx.Add(entry)

	// Get by name
	got, ok := idx.Get("mymod.dashboard.test")
	assert.True(t, ok)
	assert.Equal(t, "Test Dashboard", got.Title)

	// Get by type
	dashboards := idx.Dashboards()
	assert.Len(t, dashboards, 1)
	assert.Equal(t, "mymod.dashboard.test", dashboards[0].Name)
}

func TestIndex_GetNotFound(t *testing.T) {
	idx := NewResourceIndex()

	got, ok := idx.Get("nonexistent")
	assert.False(t, ok)
	assert.Nil(t, got)
}

func TestIndex_List(t *testing.T) {
	idx := NewResourceIndex()

	idx.Add(&IndexEntry{Type: "dashboard", Name: "mod.dashboard.d1", ShortName: "d1"})
	idx.Add(&IndexEntry{Type: "query", Name: "mod.query.q1", ShortName: "q1"})
	idx.Add(&IndexEntry{Type: "control", Name: "mod.control.c1", ShortName: "c1"})

	entries := idx.List()
	assert.Len(t, entries, 3)
}

func TestIndex_Remove(t *testing.T) {
	idx := NewResourceIndex()

	entry := &IndexEntry{
		Type:      "dashboard",
		Name:      "mymod.dashboard.test",
		ShortName: "test",
		Title:     "Test Dashboard",
	}
	idx.Add(entry)

	// Verify it exists
	_, ok := idx.Get("mymod.dashboard.test")
	assert.True(t, ok)
	assert.Equal(t, 1, idx.Count())

	// Remove it
	removed := idx.Remove("mymod.dashboard.test")
	assert.True(t, removed)

	// Verify it's gone
	_, ok = idx.Get("mymod.dashboard.test")
	assert.False(t, ok)
	assert.Equal(t, 0, idx.Count())

	// Verify dashboards list is empty
	dashboards := idx.Dashboards()
	assert.Len(t, dashboards, 0)

	// Remove nonexistent returns false
	removed = idx.Remove("nonexistent")
	assert.False(t, removed)
}

func TestIndex_BenchmarkHierarchy(t *testing.T) {
	idx := NewResourceIndex()

	// Add parent benchmark
	parent := &IndexEntry{
		Type:       "benchmark",
		Name:       "mymod.benchmark.parent",
		ShortName:  "parent",
		Title:      "Parent Benchmark",
		IsTopLevel: true,
		ChildNames: []string{"mymod.benchmark.child1", "mymod.benchmark.child2"},
	}
	idx.Add(parent)

	// Add children
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mymod.benchmark.child1",
		ShortName:  "child1",
		Title:      "Child 1",
		ParentName: "mymod.benchmark.parent",
	})
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mymod.benchmark.child2",
		ShortName:  "child2",
		Title:      "Child 2",
		ParentName: "mymod.benchmark.parent",
	})

	// Get children
	children := idx.GetChildren("mymod.benchmark.parent")
	assert.Len(t, children, 2)

	// Top level benchmarks
	topLevel := idx.TopLevelBenchmarks()
	assert.Len(t, topLevel, 1)
	assert.Equal(t, "mymod.benchmark.parent", topLevel[0].Name)
}

func TestIndex_GetChildrenMissingParent(t *testing.T) {
	idx := NewResourceIndex()

	children := idx.GetChildren("nonexistent")
	assert.Nil(t, children)
}

func TestIndex_GetChildrenNoChildren(t *testing.T) {
	idx := NewResourceIndex()

	idx.Add(&IndexEntry{
		Type:      "benchmark",
		Name:      "mymod.benchmark.leaf",
		ShortName: "leaf",
	})

	children := idx.GetChildren("mymod.benchmark.leaf")
	assert.Nil(t, children)
}

func TestIndex_GetByType(t *testing.T) {
	idx := NewResourceIndex()

	idx.Add(&IndexEntry{Type: "dashboard", Name: "mod.dashboard.d1", ShortName: "d1"})
	idx.Add(&IndexEntry{Type: "dashboard", Name: "mod.dashboard.d2", ShortName: "d2"})
	idx.Add(&IndexEntry{Type: "query", Name: "mod.query.q1", ShortName: "q1"})
	idx.Add(&IndexEntry{Type: "control", Name: "mod.control.c1", ShortName: "c1"})

	dashboards := idx.GetByType("dashboard")
	assert.Len(t, dashboards, 2)

	queries := idx.GetByType("query")
	assert.Len(t, queries, 1)

	controls := idx.GetByType("control")
	assert.Len(t, controls, 1)

	// Nonexistent type
	none := idx.GetByType("nonexistent")
	assert.Nil(t, none)
}

func TestIndex_Types(t *testing.T) {
	idx := NewResourceIndex()

	idx.Add(&IndexEntry{Type: "dashboard", Name: "mod.dashboard.d1", ShortName: "d1"})
	idx.Add(&IndexEntry{Type: "query", Name: "mod.query.q1", ShortName: "q1"})
	idx.Add(&IndexEntry{Type: "control", Name: "mod.control.c1", ShortName: "c1"})

	types := idx.Types()
	assert.Len(t, types, 3)
	assert.Contains(t, types, "dashboard")
	assert.Contains(t, types, "query")
	assert.Contains(t, types, "control")
}

func TestIndex_BenchmarksMixed(t *testing.T) {
	idx := NewResourceIndex()

	idx.Add(&IndexEntry{Type: "benchmark", Name: "mod.benchmark.b1", ShortName: "b1"})
	idx.Add(&IndexEntry{Type: "detection_benchmark", Name: "mod.detection_benchmark.db1", ShortName: "db1"})

	benchmarks := idx.Benchmarks()
	assert.Len(t, benchmarks, 2)
}

func TestIndex_Size(t *testing.T) {
	idx := NewResourceIndex()

	// Add many entries
	for i := 0; i < 1000; i++ {
		idx.Add(&IndexEntry{
			Type:      "query",
			Name:      fmt.Sprintf("mymod.query.query_%d", i),
			ShortName: fmt.Sprintf("query_%d", i),
			Title:     fmt.Sprintf("Query %d", i),
			FileName:  "/path/to/queries.pp",
			StartLine: i * 10,
			EndLine:   i*10 + 9,
		})
	}

	// Index should be small
	size := idx.Size()
	t.Logf("Index size for 1000 entries: %d bytes (%.2f KB)", size, float64(size)/1024)

	// Should be less than 500KB for 1000 entries
	assert.Less(t, size, 500*1024, "Index too large")

	// Check count
	assert.Equal(t, 1000, idx.Count())
}

func TestIndex_Stats(t *testing.T) {
	idx := NewResourceIndex()

	idx.Add(&IndexEntry{Type: "dashboard", Name: "mod.dashboard.d1", ShortName: "d1"})
	idx.Add(&IndexEntry{Type: "dashboard", Name: "mod.dashboard.d2", ShortName: "d2"})
	idx.Add(&IndexEntry{Type: "query", Name: "mod.query.q1", ShortName: "q1"})

	stats := idx.Stats()
	assert.Equal(t, 3, stats.TotalEntries)
	assert.Equal(t, 2, stats.ByType["dashboard"])
	assert.Equal(t, 1, stats.ByType["query"])
	assert.Greater(t, stats.TotalSize, 0)
}

func TestIndex_ConcurrentAccess(t *testing.T) {
	idx := NewResourceIndex()

	// Pre-populate
	for i := 0; i < 100; i++ {
		idx.Add(&IndexEntry{
			Type:      "query",
			Name:      fmt.Sprintf("mod.query.q%d", i),
			ShortName: fmt.Sprintf("q%d", i),
		})
	}

	var wg sync.WaitGroup
	// Concurrent reads
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				idx.Get(fmt.Sprintf("mod.query.q%d", j%100))
				idx.Dashboards()
				idx.Count()
				idx.Size()
			}
		}(i)
	}

	// Concurrent writes
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				idx.Add(&IndexEntry{
					Type:      "dashboard",
					Name:      fmt.Sprintf("mod.dashboard.d%d_%d", i, j),
					ShortName: fmt.Sprintf("d%d_%d", i, j),
				})
			}
		}(i)
	}

	wg.Wait()

	// Should have original 100 queries + 5*20 dashboards
	assert.Equal(t, 200, idx.Count())
}

func TestIndex_AvailableDashboardsPayload(t *testing.T) {
	idx := NewResourceIndex()

	// Add dashboards
	idx.Add(&IndexEntry{
		Type:        "dashboard",
		Name:        "mymod.dashboard.main",
		ShortName:   "main",
		Title:       "Main Dashboard",
		Tags:        map[string]string{"service": "aws"},
		ModFullName: "mymod",
	})

	// Add benchmarks with hierarchy
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mymod.benchmark.cis",
		ShortName:  "cis",
		Title:      "CIS Benchmark",
		IsTopLevel: true,
		ChildNames: []string{"mymod.benchmark.cis_1"},
	})
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mymod.benchmark.cis_1",
		ShortName:  "cis_1",
		Title:      "CIS 1.x",
		ParentName: "mymod.benchmark.cis",
	})

	payload := idx.BuildAvailableDashboardsPayload()

	assert.Equal(t, "available_dashboards", payload.Action)
	assert.Len(t, payload.Dashboards, 1)
	assert.Equal(t, "Main Dashboard", payload.Dashboards["mymod.dashboard.main"].Title)
	assert.Equal(t, "aws", payload.Dashboards["mymod.dashboard.main"].Tags["service"])

	assert.Len(t, payload.Benchmarks, 2)
	assert.True(t, payload.Benchmarks["mymod.benchmark.cis"].IsTopLevel)
	assert.Len(t, payload.Benchmarks["mymod.benchmark.cis"].Children, 1)
	assert.Equal(t, "CIS 1.x", payload.Benchmarks["mymod.benchmark.cis"].Children[0].Title)
}

func TestIndex_AvailableDashboardsPayload_DeepHierarchy(t *testing.T) {
	idx := NewResourceIndex()

	// Create a 3-level benchmark hierarchy
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.root",
		ShortName:  "root",
		Title:      "Root",
		IsTopLevel: true,
		ChildNames: []string{"mod.benchmark.level1"},
	})
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.level1",
		ShortName:  "level1",
		Title:      "Level 1",
		ParentName: "mod.benchmark.root",
		ChildNames: []string{"mod.benchmark.level2"},
	})
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.level2",
		ShortName:  "level2",
		Title:      "Level 2",
		ParentName: "mod.benchmark.level1",
	})

	payload := idx.BuildAvailableDashboardsPayload()

	// Verify hierarchy
	root := payload.Benchmarks["mod.benchmark.root"]
	require.Len(t, root.Children, 1)

	level1 := root.Children[0]
	assert.Equal(t, "Level 1", level1.Title)
	require.Len(t, level1.Children, 1)

	level2 := level1.Children[0]
	assert.Equal(t, "Level 2", level2.Title)
	assert.Len(t, level2.Children, 0)

	// Verify trunks are recorded
	assert.NotEmpty(t, payload.Benchmarks["mod.benchmark.level1"].Trunks)
	assert.NotEmpty(t, payload.Benchmarks["mod.benchmark.level2"].Trunks)
}

func TestIndex_AvailableDashboardsPayload_MixedChildren(t *testing.T) {
	idx := NewResourceIndex()

	// Benchmark with both benchmark and control children
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.parent",
		ShortName:  "parent",
		IsTopLevel: true,
		ChildNames: []string{"mod.benchmark.child", "mod.control.ctrl1"},
	})
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.child",
		ShortName:  "child",
		ParentName: "mod.benchmark.parent",
	})
	idx.Add(&IndexEntry{
		Type:       "control",
		Name:       "mod.control.ctrl1",
		ShortName:  "ctrl1",
		ParentName: "mod.benchmark.parent",
	})

	payload := idx.BuildAvailableDashboardsPayload()

	parent := payload.Benchmarks["mod.benchmark.parent"]
	// Only benchmark children should be included, not controls
	assert.Len(t, parent.Children, 1)
	assert.Equal(t, "mod.benchmark.child", parent.Children[0].FullName)
}

func TestIndexEntry_Size(t *testing.T) {
	entry := &IndexEntry{
		Type:      "dashboard",
		Name:      "mymod.dashboard.test_dashboard",
		ShortName: "test_dashboard",
		Title:     "Test Dashboard with a Longer Title",
		Description: "This is a longer description that adds some " +
			"additional bytes to the entry size calculation",
		FileName:  "/path/to/mod/dashboards.pp",
		StartLine: 10,
		EndLine:   50,
		Tags: map[string]string{
			"service": "aws",
			"type":    "compliance",
		},
		ChildNames: []string{"mymod.dashboard.child1", "mymod.dashboard.child2"},
	}

	size := entry.Size()
	t.Logf("Entry size: %d bytes", size)

	// Should be reasonable
	assert.Greater(t, size, 100)
	assert.Less(t, size, 1000)
}

func TestIndex_Empty(t *testing.T) {
	idx := NewResourceIndex()

	assert.Equal(t, 0, idx.Count())
	assert.Equal(t, 0, idx.Size())
	assert.Empty(t, idx.List())
	assert.Empty(t, idx.Types())
	assert.Empty(t, idx.Dashboards())
	assert.Empty(t, idx.Benchmarks())
	assert.Empty(t, idx.Queries())
	assert.Empty(t, idx.Controls())
	assert.Empty(t, idx.TopLevelBenchmarks())

	stats := idx.Stats()
	assert.Equal(t, 0, stats.TotalEntries)
	assert.Equal(t, 0, stats.TotalSize)
	assert.Empty(t, stats.ByType)

	payload := idx.BuildAvailableDashboardsPayload()
	assert.Equal(t, "available_dashboards", payload.Action)
	assert.Empty(t, payload.Dashboards)
	assert.Empty(t, payload.Benchmarks)
}
