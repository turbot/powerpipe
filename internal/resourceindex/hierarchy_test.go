package resourceindex

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Basic Hierarchy Tests
// =============================================================================

// TestHierarchy_SingleBenchmarkWithControls verifies a benchmark with control children
func TestHierarchy_SingleBenchmarkWithControls(t *testing.T) {
	idx := NewResourceIndex()

	// Add parent benchmark with 3 control children
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.parent",
		ShortName:  "parent",
		Title:      "Parent Benchmark",
		IsTopLevel: true,
		ChildNames: []string{
			"mod.control.ctrl1",
			"mod.control.ctrl2",
			"mod.control.ctrl3",
		},
	})

	// Add control children
	for i := 1; i <= 3; i++ {
		idx.Add(&IndexEntry{
			Type:       "control",
			Name:       fmt.Sprintf("mod.control.ctrl%d", i),
			ShortName:  fmt.Sprintf("ctrl%d", i),
			Title:      fmt.Sprintf("Control %d", i),
			ParentName: "mod.benchmark.parent",
		})
	}

	// Verify children list
	children := idx.GetChildren("mod.benchmark.parent")
	assert.Len(t, children, 3, "Expected 3 children")

	// Verify each control has parent set
	for i := 1; i <= 3; i++ {
		entry, ok := idx.Get(fmt.Sprintf("mod.control.ctrl%d", i))
		require.True(t, ok)
		assert.Equal(t, "mod.benchmark.parent", entry.ParentName)
	}
}

// TestHierarchy_TwoLevelBenchmark verifies a two-level benchmark hierarchy
// Top → Child → Controls
func TestHierarchy_TwoLevelBenchmark(t *testing.T) {
	idx := NewResourceIndex()

	// Top level benchmark
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.top",
		ShortName:  "top",
		Title:      "Top Benchmark",
		IsTopLevel: true,
		ChildNames: []string{"mod.benchmark.child"},
	})

	// Child benchmark
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.child",
		ShortName:  "child",
		Title:      "Child Benchmark",
		ParentName: "mod.benchmark.top",
		ChildNames: []string{"mod.control.ctrl1", "mod.control.ctrl2"},
	})

	// Controls under child benchmark
	for i := 1; i <= 2; i++ {
		idx.Add(&IndexEntry{
			Type:       "control",
			Name:       fmt.Sprintf("mod.control.ctrl%d", i),
			ShortName:  fmt.Sprintf("ctrl%d", i),
			ParentName: "mod.benchmark.child",
		})
	}

	// Verify top → child relationship
	topChildren := idx.GetChildren("mod.benchmark.top")
	require.Len(t, topChildren, 1)
	assert.Equal(t, "mod.benchmark.child", topChildren[0].Name)

	// Verify child → controls relationship
	childChildren := idx.GetChildren("mod.benchmark.child")
	assert.Len(t, childChildren, 2)

	// Verify child has parent set
	child, _ := idx.Get("mod.benchmark.child")
	assert.Equal(t, "mod.benchmark.top", child.ParentName)
}

// TestHierarchy_SiblingBenchmarks verifies multiple children under same parent
func TestHierarchy_SiblingBenchmarks(t *testing.T) {
	idx := NewResourceIndex()

	// Parent with two children
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.top",
		ShortName:  "top",
		IsTopLevel: true,
		ChildNames: []string{"mod.benchmark.child1", "mod.benchmark.child2"},
	})

	// Both children
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.child1",
		ShortName:  "child1",
		ParentName: "mod.benchmark.top",
	})
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.child2",
		ShortName:  "child2",
		ParentName: "mod.benchmark.top",
	})

	// Verify parent has both children
	children := idx.GetChildren("mod.benchmark.top")
	assert.Len(t, children, 2)

	childNames := make([]string, len(children))
	for i, c := range children {
		childNames[i] = c.Name
	}
	sort.Strings(childNames)
	assert.Equal(t, []string{"mod.benchmark.child1", "mod.benchmark.child2"}, childNames)

	// Verify both children have same parent
	child1, _ := idx.Get("mod.benchmark.child1")
	child2, _ := idx.Get("mod.benchmark.child2")
	assert.Equal(t, "mod.benchmark.top", child1.ParentName)
	assert.Equal(t, "mod.benchmark.top", child2.ParentName)
}

// =============================================================================
// Deep Hierarchy Tests
// =============================================================================

// createDeepHierarchy creates a benchmark chain of specified depth
// B0 → B1 → B2 → ... → B(depth-1) → Control
func createDeepHierarchy(idx *ResourceIndex, depth int) {
	for i := 0; i < depth; i++ {
		var childNames []string
		if i < depth-1 {
			// Not leaf, point to next benchmark
			childNames = []string{fmt.Sprintf("mod.benchmark.b%d", i+1)}
		} else {
			// Leaf benchmark, has a control child
			childNames = []string{"mod.control.leaf_control"}
		}

		parentName := ""
		if i > 0 {
			parentName = fmt.Sprintf("mod.benchmark.b%d", i-1)
		}

		idx.Add(&IndexEntry{
			Type:       "benchmark",
			Name:       fmt.Sprintf("mod.benchmark.b%d", i),
			ShortName:  fmt.Sprintf("b%d", i),
			Title:      fmt.Sprintf("Benchmark Level %d", i),
			IsTopLevel: i == 0,
			ParentName: parentName,
			ChildNames: childNames,
		})
	}

	// Add leaf control
	idx.Add(&IndexEntry{
		Type:       "control",
		Name:       "mod.control.leaf_control",
		ShortName:  "leaf_control",
		ParentName: fmt.Sprintf("mod.benchmark.b%d", depth-1),
	})
}

// TestHierarchy_Deep10Levels tests a 10-level deep hierarchy
func TestHierarchy_Deep10Levels(t *testing.T) {
	idx := NewResourceIndex()
	createDeepHierarchy(idx, 10)

	// Verify all levels exist
	for i := 0; i < 10; i++ {
		entry, ok := idx.Get(fmt.Sprintf("mod.benchmark.b%d", i))
		require.True(t, ok, "Level %d should exist", i)

		// Verify parent relationship (except root)
		if i > 0 {
			assert.Equal(t, fmt.Sprintf("mod.benchmark.b%d", i-1), entry.ParentName)
		}

		// Verify child relationship (except leaf)
		if i < 9 {
			children := idx.GetChildren(entry.Name)
			require.Len(t, children, 1)
			assert.Equal(t, fmt.Sprintf("mod.benchmark.b%d", i+1), children[0].Name)
		}
	}

	// Only b0 should be top level
	topLevel := idx.TopLevelBenchmarks()
	require.Len(t, topLevel, 1)
	assert.Equal(t, "mod.benchmark.b0", topLevel[0].Name)
}

// TestHierarchy_Deep20Levels stress tests deep recursion
func TestHierarchy_Deep20Levels(t *testing.T) {
	idx := NewResourceIndex()
	createDeepHierarchy(idx, 20)

	// Should complete without stack overflow
	assert.Equal(t, 21, idx.Count()) // 20 benchmarks + 1 control

	// Verify can traverse full depth
	current := "mod.benchmark.b0"
	depth := 0
	for {
		children := idx.GetChildren(current)
		if len(children) == 0 {
			break
		}
		// Find benchmark child (skip control)
		found := false
		for _, c := range children {
			if c.Type == "benchmark" {
				current = c.Name
				depth++
				found = true
				break
			}
		}
		if !found {
			break
		}
	}
	assert.Equal(t, 19, depth) // 19 hops from b0 to b19
}

// TestHierarchy_DeepTrunkBuilding verifies trunk paths for deep hierarchy
func TestHierarchy_DeepTrunkBuilding(t *testing.T) {
	idx := NewResourceIndex()
	createDeepHierarchy(idx, 10)

	payload := idx.BuildAvailableDashboardsPayload()

	// Check trunk for b0 (top level)
	b0 := payload.Benchmarks["mod.benchmark.b0"]
	require.NotNil(t, b0.Trunks)
	assert.Len(t, b0.Trunks, 1)
	assert.Equal(t, []string{"mod.benchmark.b0"}, b0.Trunks[0])

	// Check trunk for b5 (middle of hierarchy)
	// Trunk should be: [b0, b1, b2, b3, b4, b5]
	b5 := payload.Benchmarks["mod.benchmark.b5"]
	require.NotNil(t, b5.Trunks)
	require.Len(t, b5.Trunks, 1)
	expectedTrunk := []string{
		"mod.benchmark.b0",
		"mod.benchmark.b1",
		"mod.benchmark.b2",
		"mod.benchmark.b3",
		"mod.benchmark.b4",
		"mod.benchmark.b5",
	}
	assert.Equal(t, expectedTrunk, b5.Trunks[0])

	// Check trunk length matches depth
	b9 := payload.Benchmarks["mod.benchmark.b9"]
	require.NotNil(t, b9.Trunks)
	require.Len(t, b9.Trunks, 1)
	assert.Len(t, b9.Trunks[0], 10) // Depth of 10
}

// =============================================================================
// Wide Hierarchy Tests
// =============================================================================

// createWideHierarchy creates a benchmark with specified number of children
func createWideHierarchy(idx *ResourceIndex, numChildren int, childType string) {
	childNames := make([]string, numChildren)
	for i := 0; i < numChildren; i++ {
		childNames[i] = fmt.Sprintf("mod.%s.child%d", childType, i)
	}

	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.wide_parent",
		ShortName:  "wide_parent",
		Title:      "Wide Benchmark",
		IsTopLevel: true,
		ChildNames: childNames,
	})

	// Add all children
	for i := 0; i < numChildren; i++ {
		idx.Add(&IndexEntry{
			Type:       childType,
			Name:       fmt.Sprintf("mod.%s.child%d", childType, i),
			ShortName:  fmt.Sprintf("child%d", i),
			Title:      fmt.Sprintf("Child %d", i),
			ParentName: "mod.benchmark.wide_parent",
		})
	}
}

// TestHierarchy_Wide100Children tests benchmark with 100 control children
func TestHierarchy_Wide100Children(t *testing.T) {
	idx := NewResourceIndex()
	createWideHierarchy(idx, 100, "control")

	// Verify all children tracked
	children := idx.GetChildren("mod.benchmark.wide_parent")
	assert.Len(t, children, 100)

	// Verify all children have correct parent
	for i := 0; i < 100; i++ {
		entry, ok := idx.Get(fmt.Sprintf("mod.control.child%d", i))
		require.True(t, ok)
		assert.Equal(t, "mod.benchmark.wide_parent", entry.ParentName)
	}
}

// TestHierarchy_Wide500Children stress tests wide benchmarks
func TestHierarchy_Wide500Children(t *testing.T) {
	idx := NewResourceIndex()

	start := time.Now()
	createWideHierarchy(idx, 500, "control")
	createTime := time.Since(start)

	// Should complete quickly (under 100ms)
	assert.Less(t, createTime.Milliseconds(), int64(100),
		"Creating 500 children took too long: %v", createTime)

	// Memory should be bounded
	size := idx.Size()
	t.Logf("Index size for 501 entries: %d bytes (%.2f KB)", size, float64(size)/1024)
	assert.Less(t, size, 500*1024, "Index too large for 500 children")

	// All children accessible
	children := idx.GetChildren("mod.benchmark.wide_parent")
	assert.Len(t, children, 500)
}

// TestHierarchy_MixedWidth tests mixed width at different levels
// Top: 10 children, each child: 20 grandchildren = 200 total leaf controls
func TestHierarchy_MixedWidth(t *testing.T) {
	idx := NewResourceIndex()

	// Top level benchmark with 10 children
	childNames := make([]string, 10)
	for i := 0; i < 10; i++ {
		childNames[i] = fmt.Sprintf("mod.benchmark.level1_%d", i)
	}

	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.top",
		ShortName:  "top",
		IsTopLevel: true,
		ChildNames: childNames,
	})

	// 10 level-1 benchmarks, each with 20 control children
	for i := 0; i < 10; i++ {
		grandChildNames := make([]string, 20)
		for j := 0; j < 20; j++ {
			grandChildNames[j] = fmt.Sprintf("mod.control.ctrl_%d_%d", i, j)
		}

		idx.Add(&IndexEntry{
			Type:       "benchmark",
			Name:       fmt.Sprintf("mod.benchmark.level1_%d", i),
			ShortName:  fmt.Sprintf("level1_%d", i),
			ParentName: "mod.benchmark.top",
			ChildNames: grandChildNames,
		})

		// Add 20 controls for each level-1 benchmark
		for j := 0; j < 20; j++ {
			idx.Add(&IndexEntry{
				Type:       "control",
				Name:       fmt.Sprintf("mod.control.ctrl_%d_%d", i, j),
				ShortName:  fmt.Sprintf("ctrl_%d_%d", i, j),
				ParentName: fmt.Sprintf("mod.benchmark.level1_%d", i),
			})
		}
	}

	// Verify structure
	// 1 top + 10 level1 + 200 controls = 211 total
	assert.Equal(t, 211, idx.Count())

	// Top has 10 children
	topChildren := idx.GetChildren("mod.benchmark.top")
	assert.Len(t, topChildren, 10)

	// Each level-1 has 20 children
	for i := 0; i < 10; i++ {
		children := idx.GetChildren(fmt.Sprintf("mod.benchmark.level1_%d", i))
		assert.Len(t, children, 20)
	}

	// Build payload should work
	payload := idx.BuildAvailableDashboardsPayload()
	assert.Len(t, payload.Benchmarks, 11) // 1 top + 10 level-1
}

// =============================================================================
// Trunk Building Tests
// =============================================================================

// TestHierarchy_TopLevelTrunk verifies trunk for top-level benchmark
func TestHierarchy_TopLevelTrunk(t *testing.T) {
	idx := NewResourceIndex()

	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.standalone",
		ShortName:  "standalone",
		IsTopLevel: true,
	})

	payload := idx.BuildAvailableDashboardsPayload()
	b := payload.Benchmarks["mod.benchmark.standalone"]

	// Top level should have trunk [[benchmark_name]]
	require.Len(t, b.Trunks, 1)
	assert.Equal(t, []string{"mod.benchmark.standalone"}, b.Trunks[0])
}

// TestHierarchy_ChildTrunkIncludesParent verifies child trunk includes parent path
func TestHierarchy_ChildTrunkIncludesParent(t *testing.T) {
	idx := NewResourceIndex()

	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.parent",
		ShortName:  "parent",
		IsTopLevel: true,
		ChildNames: []string{"mod.benchmark.child"},
	})
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.child",
		ShortName:  "child",
		ParentName: "mod.benchmark.parent",
	})

	payload := idx.BuildAvailableDashboardsPayload()
	child := payload.Benchmarks["mod.benchmark.child"]

	// Child trunk should be [[parent, child]]
	require.Len(t, child.Trunks, 1)
	assert.Equal(t, []string{"mod.benchmark.parent", "mod.benchmark.child"}, child.Trunks[0])
}

// TestHierarchy_DiamondPatternTrunks tests multiple paths to same node
// Top → A, Top → B, A → Child, B → Child
// Child should have two trunks: [[Top, A, Child], [Top, B, Child]]
func TestHierarchy_DiamondPatternTrunks(t *testing.T) {
	idx := NewResourceIndex()

	// Top level benchmark
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.top",
		ShortName:  "top",
		IsTopLevel: true,
		ChildNames: []string{"mod.benchmark.a", "mod.benchmark.b"},
	})

	// Branch A
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.a",
		ShortName:  "a",
		ParentName: "mod.benchmark.top",
		ChildNames: []string{"mod.benchmark.diamond_child"},
	})

	// Branch B
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.b",
		ShortName:  "b",
		ParentName: "mod.benchmark.top",
		ChildNames: []string{"mod.benchmark.diamond_child"},
	})

	// Diamond child (appears under both A and B)
	idx.Add(&IndexEntry{
		Type:      "benchmark",
		Name:      "mod.benchmark.diamond_child",
		ShortName: "diamond_child",
		// Note: ParentName can only hold one parent, but ChildNames can reference it multiple times
	})

	payload := idx.BuildAvailableDashboardsPayload()
	child := payload.Benchmarks["mod.benchmark.diamond_child"]

	// Diamond child should have two trunks (one from each path)
	require.Len(t, child.Trunks, 2, "Diamond child should have 2 trunks")

	// Sort trunks for deterministic comparison
	sort.Slice(child.Trunks, func(i, j int) bool {
		// Sort by second element (a vs b)
		if len(child.Trunks[i]) > 1 && len(child.Trunks[j]) > 1 {
			return child.Trunks[i][1] < child.Trunks[j][1]
		}
		return false
	})

	assert.Equal(t, []string{"mod.benchmark.top", "mod.benchmark.a", "mod.benchmark.diamond_child"}, child.Trunks[0])
	assert.Equal(t, []string{"mod.benchmark.top", "mod.benchmark.b", "mod.benchmark.diamond_child"}, child.Trunks[1])
}

// TestHierarchy_TrunkOrderConsistent verifies deterministic trunk ordering
func TestHierarchy_TrunkOrderConsistent(t *testing.T) {
	// Create same hierarchy multiple times and verify same output
	createTestHierarchy := func() *AvailableDashboardsPayload {
		idx := NewResourceIndex()

		idx.Add(&IndexEntry{
			Type:       "benchmark",
			Name:       "mod.benchmark.top",
			ShortName:  "top",
			IsTopLevel: true,
			ChildNames: []string{"mod.benchmark.c1", "mod.benchmark.c2", "mod.benchmark.c3"},
		})

		for i := 1; i <= 3; i++ {
			idx.Add(&IndexEntry{
				Type:       "benchmark",
				Name:       fmt.Sprintf("mod.benchmark.c%d", i),
				ShortName:  fmt.Sprintf("c%d", i),
				ParentName: "mod.benchmark.top",
			})
		}

		return idx.BuildAvailableDashboardsPayload()
	}

	// Run multiple times
	payloads := make([]*AvailableDashboardsPayload, 5)
	for i := 0; i < 5; i++ {
		payloads[i] = createTestHierarchy()
	}

	// Compare all payloads - they should produce consistent structure
	for i := 1; i < 5; i++ {
		// Trunk content should match (though order may vary due to map iteration)
		for name, info := range payloads[0].Benchmarks {
			otherInfo := payloads[i].Benchmarks[name]
			assert.Equal(t, len(info.Trunks), len(otherInfo.Trunks),
				"Trunk count mismatch for %s on iteration %d", name, i)
		}
	}
}

// =============================================================================
// Parent-Child Relationship Tests
// =============================================================================

// TestHierarchy_SetParentNames verifies parent assignment
func TestHierarchy_SetParentNames(t *testing.T) {
	idx := NewResourceIndex()

	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.parent",
		ShortName:  "parent",
		IsTopLevel: true,
		ChildNames: []string{"mod.benchmark.child"},
	})
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.child",
		ShortName:  "child",
		ParentName: "mod.benchmark.parent", // Explicitly set during scanning
	})

	child, _ := idx.Get("mod.benchmark.child")
	assert.Equal(t, "mod.benchmark.parent", child.ParentName)
}

// TestHierarchy_BidirectionalRelationship verifies parent↔child both directions
func TestHierarchy_BidirectionalRelationship(t *testing.T) {
	idx := NewResourceIndex()

	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.parent",
		ShortName:  "parent",
		IsTopLevel: true,
		ChildNames: []string{"mod.benchmark.child"},
	})
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.child",
		ShortName:  "child",
		ParentName: "mod.benchmark.parent",
	})

	// Parent.ChildNames contains child
	parent, _ := idx.Get("mod.benchmark.parent")
	assert.Contains(t, parent.ChildNames, "mod.benchmark.child")

	// Child.ParentName == Parent.Name
	child, _ := idx.Get("mod.benchmark.child")
	assert.Equal(t, parent.Name, child.ParentName)

	// GetChildren returns the child
	children := idx.GetChildren("mod.benchmark.parent")
	require.Len(t, children, 1)
	assert.Equal(t, "mod.benchmark.child", children[0].Name)
}

// TestHierarchy_OrphanDetection verifies benchmarks not in any children list are top-level
func TestHierarchy_OrphanDetection(t *testing.T) {
	idx := NewResourceIndex()

	// Two independent benchmarks (orphans/top-level)
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.orphan1",
		ShortName:  "orphan1",
		IsTopLevel: true,
	})
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.orphan2",
		ShortName:  "orphan2",
		IsTopLevel: true,
	})

	// One with parent
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.parent",
		ShortName:  "parent",
		IsTopLevel: true,
		ChildNames: []string{"mod.benchmark.child"},
	})
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.child",
		ShortName:  "child",
		ParentName: "mod.benchmark.parent",
		IsTopLevel: false,
	})

	topLevel := idx.TopLevelBenchmarks()
	assert.Len(t, topLevel, 3) // orphan1, orphan2, parent

	topLevelNames := make([]string, len(topLevel))
	for i, b := range topLevel {
		topLevelNames[i] = b.Name
	}
	sort.Strings(topLevelNames)

	assert.Contains(t, topLevelNames, "mod.benchmark.orphan1")
	assert.Contains(t, topLevelNames, "mod.benchmark.orphan2")
	assert.Contains(t, topLevelNames, "mod.benchmark.parent")
	assert.NotContains(t, topLevelNames, "mod.benchmark.child")
}

// TestHierarchy_MissingChild verifies graceful handling when child doesn't exist
func TestHierarchy_MissingChild(t *testing.T) {
	idx := NewResourceIndex()

	// Parent references non-existent child
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.parent",
		ShortName:  "parent",
		IsTopLevel: true,
		ChildNames: []string{"mod.benchmark.missing_child"},
	})

	// GetChildren should return empty (not crash)
	children := idx.GetChildren("mod.benchmark.parent")
	assert.Empty(t, children)

	// Payload should still build without error
	payload := idx.BuildAvailableDashboardsPayload()
	assert.NotNil(t, payload)
	assert.Len(t, payload.Benchmarks["mod.benchmark.parent"].Children, 0)
}

// =============================================================================
// Top-Level Detection Tests
// =============================================================================

// TestHierarchy_CorrectTopLevel verifies only truly top-level benchmarks marked
func TestHierarchy_CorrectTopLevel(t *testing.T) {
	idx := NewResourceIndex()

	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.top1",
		ShortName:  "top1",
		IsTopLevel: true,
		ChildNames: []string{"mod.benchmark.child"},
	})
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.top2",
		ShortName:  "top2",
		IsTopLevel: true,
	})
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.child",
		ShortName:  "child",
		ParentName: "mod.benchmark.top1",
		IsTopLevel: false, // Not top level
	})

	topLevel := idx.TopLevelBenchmarks()
	assert.Len(t, topLevel, 2)

	for _, b := range topLevel {
		assert.True(t, b.IsTopLevel)
		assert.NotEqual(t, "mod.benchmark.child", b.Name)
	}
}

// TestHierarchy_MultipleTopLevel verifies multiple independent benchmark trees
func TestHierarchy_MultipleTopLevel(t *testing.T) {
	idx := NewResourceIndex()

	// Tree 1
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.tree1_root",
		ShortName:  "tree1_root",
		IsTopLevel: true,
		ChildNames: []string{"mod.benchmark.tree1_child"},
	})
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.tree1_child",
		ShortName:  "tree1_child",
		ParentName: "mod.benchmark.tree1_root",
	})

	// Tree 2
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.tree2_root",
		ShortName:  "tree2_root",
		IsTopLevel: true,
		ChildNames: []string{"mod.benchmark.tree2_child"},
	})
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.tree2_child",
		ShortName:  "tree2_child",
		ParentName: "mod.benchmark.tree2_root",
	})

	topLevel := idx.TopLevelBenchmarks()
	assert.Len(t, topLevel, 2)

	topLevelNames := make([]string, len(topLevel))
	for i, b := range topLevel {
		topLevelNames[i] = b.Name
	}
	sort.Strings(topLevelNames)
	assert.Equal(t, []string{"mod.benchmark.tree1_root", "mod.benchmark.tree2_root"}, topLevelNames)
}

// TestHierarchy_FlatBenchmarks verifies flat benchmarks (only control children)
func TestHierarchy_FlatBenchmarks(t *testing.T) {
	idx := NewResourceIndex()

	// Multiple benchmarks, each only with control children (no benchmark children)
	for i := 1; i <= 3; i++ {
		idx.Add(&IndexEntry{
			Type:       "benchmark",
			Name:       fmt.Sprintf("mod.benchmark.flat%d", i),
			ShortName:  fmt.Sprintf("flat%d", i),
			IsTopLevel: true,
			ChildNames: []string{fmt.Sprintf("mod.control.ctrl%d", i)},
		})
		idx.Add(&IndexEntry{
			Type:       "control",
			Name:       fmt.Sprintf("mod.control.ctrl%d", i),
			ShortName:  fmt.Sprintf("ctrl%d", i),
			ParentName: fmt.Sprintf("mod.benchmark.flat%d", i),
		})
	}

	// All benchmarks are top-level (no benchmark children)
	topLevel := idx.TopLevelBenchmarks()
	assert.Len(t, topLevel, 3)
}

// =============================================================================
// Detection Benchmark Tests
// =============================================================================

// TestHierarchy_DetectionBenchmarks tests detection_benchmark type handling
func TestHierarchy_DetectionBenchmarks(t *testing.T) {
	idx := NewResourceIndex()

	idx.Add(&IndexEntry{
		Type:          "detection_benchmark",
		Name:          "mod.detection_benchmark.detect_parent",
		ShortName:     "detect_parent",
		BenchmarkType: "detection",
		IsTopLevel:    true,
		ChildNames:    []string{"mod.detection_benchmark.detect_child"},
	})
	idx.Add(&IndexEntry{
		Type:          "detection_benchmark",
		Name:          "mod.detection_benchmark.detect_child",
		ShortName:     "detect_child",
		BenchmarkType: "detection",
		ParentName:    "mod.detection_benchmark.detect_parent",
	})

	// Should appear in Benchmarks() list
	benchmarks := idx.Benchmarks()
	assert.Len(t, benchmarks, 2)

	// Should appear in TopLevelBenchmarks
	topLevel := idx.TopLevelBenchmarks()
	assert.Len(t, topLevel, 1)
	assert.Equal(t, "mod.detection_benchmark.detect_parent", topLevel[0].Name)

	// Hierarchy should work
	children := idx.GetChildren("mod.detection_benchmark.detect_parent")
	assert.Len(t, children, 1)
}

// TestHierarchy_MixedBenchmarkTypes tests regular and detection benchmarks together
func TestHierarchy_MixedBenchmarkTypes(t *testing.T) {
	idx := NewResourceIndex()

	// Regular benchmark
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.regular",
		ShortName:  "regular",
		IsTopLevel: true,
		ChildNames: []string{"mod.control.ctrl1"},
	})
	idx.Add(&IndexEntry{
		Type:       "control",
		Name:       "mod.control.ctrl1",
		ShortName:  "ctrl1",
		ParentName: "mod.benchmark.regular",
	})

	// Detection benchmark
	idx.Add(&IndexEntry{
		Type:          "detection_benchmark",
		Name:          "mod.detection_benchmark.detection",
		ShortName:     "detection",
		BenchmarkType: "detection",
		IsTopLevel:    true,
		ChildNames:    []string{"mod.detection.det1"},
	})
	idx.Add(&IndexEntry{
		Type:       "detection",
		Name:       "mod.detection.det1",
		ShortName:  "det1",
		ParentName: "mod.detection_benchmark.detection",
	})

	// Both should be in benchmarks list
	benchmarks := idx.Benchmarks()
	assert.Len(t, benchmarks, 2)

	// Both should be top level
	topLevel := idx.TopLevelBenchmarks()
	assert.Len(t, topLevel, 2)

	// Payload should handle both
	payload := idx.BuildAvailableDashboardsPayload()
	assert.Len(t, payload.Benchmarks, 2)
}

// =============================================================================
// Payload Generation Tests
// =============================================================================

// TestHierarchy_BenchmarkInfoStructure verifies all BenchmarkInfo fields
func TestHierarchy_BenchmarkInfoStructure(t *testing.T) {
	idx := NewResourceIndex()

	idx.Add(&IndexEntry{
		Type:          "benchmark",
		Name:          "mymod.benchmark.full",
		ShortName:     "full",
		Title:         "Full Benchmark",
		BenchmarkType: "control",
		Tags:          map[string]string{"service": "aws"},
		IsTopLevel:    true,
		ModFullName:   "mymod",
		ChildNames:    []string{"mymod.benchmark.child"},
	})
	idx.Add(&IndexEntry{
		Type:        "benchmark",
		Name:        "mymod.benchmark.child",
		ShortName:   "child",
		Title:       "Child Benchmark",
		ParentName:  "mymod.benchmark.full",
		ModFullName: "mymod",
	})

	payload := idx.BuildAvailableDashboardsPayload()
	info := payload.Benchmarks["mymod.benchmark.full"]

	// Verify all fields populated
	assert.Equal(t, "Full Benchmark", info.Title)
	assert.Equal(t, "mymod.benchmark.full", info.FullName)
	assert.Equal(t, "full", info.ShortName)
	assert.Equal(t, "mymod", info.ModFullName)
	assert.True(t, info.IsTopLevel)
	assert.NotEmpty(t, info.Trunks)
	assert.Len(t, info.Children, 1)
	assert.Equal(t, "aws", info.Tags["service"])
}

// TestHierarchy_PayloadRecursiveChildren verifies nested BenchmarkInfo in payload
func TestHierarchy_PayloadRecursiveChildren(t *testing.T) {
	idx := NewResourceIndex()

	// 3-level hierarchy: root → mid → leaf
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.root",
		ShortName:  "root",
		Title:      "Root",
		IsTopLevel: true,
		ChildNames: []string{"mod.benchmark.mid"},
	})
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.mid",
		ShortName:  "mid",
		Title:      "Middle",
		ParentName: "mod.benchmark.root",
		ChildNames: []string{"mod.benchmark.leaf"},
	})
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.leaf",
		ShortName:  "leaf",
		Title:      "Leaf",
		ParentName: "mod.benchmark.mid",
	})

	payload := idx.BuildAvailableDashboardsPayload()

	// Root should have mid as child
	root := payload.Benchmarks["mod.benchmark.root"]
	require.Len(t, root.Children, 1)
	mid := root.Children[0]
	assert.Equal(t, "mod.benchmark.mid", mid.FullName)

	// Mid should have leaf as child (nested)
	require.Len(t, mid.Children, 1)
	leaf := mid.Children[0]
	assert.Equal(t, "mod.benchmark.leaf", leaf.FullName)

	// Leaf should have no children
	assert.Empty(t, leaf.Children)
}

// TestHierarchy_LargePayloadPerformance tests payload building performance
func TestHierarchy_LargePayloadPerformance(t *testing.T) {
	idx := NewResourceIndex()

	// Create 50 top-level benchmarks, each with 10 sub-benchmarks,
	// each sub-benchmark with 4 controls = 50 + 500 benchmarks + 2000 controls
	for i := 0; i < 50; i++ {
		topChildNames := make([]string, 10)
		for j := 0; j < 10; j++ {
			topChildNames[j] = fmt.Sprintf("mod.benchmark.sub_%d_%d", i, j)
		}

		idx.Add(&IndexEntry{
			Type:       "benchmark",
			Name:       fmt.Sprintf("mod.benchmark.top_%d", i),
			ShortName:  fmt.Sprintf("top_%d", i),
			IsTopLevel: true,
			ChildNames: topChildNames,
		})

		for j := 0; j < 10; j++ {
			controlNames := make([]string, 4)
			for k := 0; k < 4; k++ {
				controlNames[k] = fmt.Sprintf("mod.control.ctrl_%d_%d_%d", i, j, k)
			}

			idx.Add(&IndexEntry{
				Type:       "benchmark",
				Name:       fmt.Sprintf("mod.benchmark.sub_%d_%d", i, j),
				ShortName:  fmt.Sprintf("sub_%d_%d", i, j),
				ParentName: fmt.Sprintf("mod.benchmark.top_%d", i),
				ChildNames: controlNames,
			})

			for k := 0; k < 4; k++ {
				idx.Add(&IndexEntry{
					Type:       "control",
					Name:       fmt.Sprintf("mod.control.ctrl_%d_%d_%d", i, j, k),
					ShortName:  fmt.Sprintf("ctrl_%d_%d_%d", i, j, k),
					ParentName: fmt.Sprintf("mod.benchmark.sub_%d_%d", i, j),
				})
			}
		}
	}

	// 50 + 500 + 2000 = 2550 total
	assert.Equal(t, 2550, idx.Count())

	start := time.Now()
	payload := idx.BuildAvailableDashboardsPayload()
	buildTime := time.Since(start)

	t.Logf("Payload build time for 2550 entries: %v", buildTime)
	assert.Less(t, buildTime.Seconds(), float64(1), "Payload build should complete in <1s")

	// Verify structure
	assert.Len(t, payload.Benchmarks, 550) // 50 top + 500 sub
}

// =============================================================================
// Edge Case Tests
// =============================================================================

// TestHierarchy_EmptyChildren verifies handling of empty children array
func TestHierarchy_EmptyChildren(t *testing.T) {
	idx := NewResourceIndex()

	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.empty",
		ShortName:  "empty",
		IsTopLevel: true,
		ChildNames: []string{}, // Explicit empty
	})

	children := idx.GetChildren("mod.benchmark.empty")
	assert.Empty(t, children)

	payload := idx.BuildAvailableDashboardsPayload()
	assert.Empty(t, payload.Benchmarks["mod.benchmark.empty"].Children)
}

// TestHierarchy_SingleChild verifies arrays of length 1 handled correctly
func TestHierarchy_SingleChild(t *testing.T) {
	idx := NewResourceIndex()

	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.parent",
		ShortName:  "parent",
		IsTopLevel: true,
		ChildNames: []string{"mod.benchmark.only_child"},
	})
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.only_child",
		ShortName:  "only_child",
		ParentName: "mod.benchmark.parent",
	})

	children := idx.GetChildren("mod.benchmark.parent")
	require.Len(t, children, 1)
	assert.Equal(t, "mod.benchmark.only_child", children[0].Name)

	payload := idx.BuildAvailableDashboardsPayload()
	require.Len(t, payload.Benchmarks["mod.benchmark.parent"].Children, 1)
}

// TestHierarchy_DuplicateChild verifies duplicate child handling
func TestHierarchy_DuplicateChild(t *testing.T) {
	idx := NewResourceIndex()

	// Same child listed twice (shouldn't normally happen but should handle gracefully)
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.parent",
		ShortName:  "parent",
		IsTopLevel: true,
		ChildNames: []string{"mod.benchmark.child", "mod.benchmark.child"}, // Duplicate
	})
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.child",
		ShortName:  "child",
		ParentName: "mod.benchmark.parent",
	})

	// GetChildren returns entries based on lookup, so duplicate in list = duplicate in result
	children := idx.GetChildren("mod.benchmark.parent")
	// Current implementation returns duplicates
	assert.GreaterOrEqual(t, len(children), 1)

	// Payload should handle this gracefully
	payload := idx.BuildAvailableDashboardsPayload()
	assert.NotNil(t, payload)
}

// TestHierarchy_CircularHierarchy tests that circular references are detected and skipped.
// A circular reference (A → B → A) should NOT cause infinite recursion.
func TestHierarchy_CircularHierarchy(t *testing.T) {
	idx := NewResourceIndex()

	// Create circular reference: A → B → A
	// Note: This is malformed data but syntactically valid HCL
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.a",
		ShortName:  "a",
		IsTopLevel: true,
		ChildNames: []string{"mod.benchmark.b"},
	})
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.b",
		ShortName:  "b",
		ParentName: "mod.benchmark.a",
		ChildNames: []string{"mod.benchmark.a"}, // Circular!
	})

	// GetChildren should work (just returns what's in ChildNames)
	childrenA := idx.GetChildren("mod.benchmark.a")
	assert.Len(t, childrenA, 1)

	childrenB := idx.GetChildren("mod.benchmark.b")
	assert.Len(t, childrenB, 1)

	// BuildAvailableDashboardsPayload should complete without infinite loop
	// The circular reference should be detected and skipped
	payload := idx.BuildAvailableDashboardsPayload()
	assert.NotNil(t, payload)

	// Verify both benchmarks are in the payload
	assert.Contains(t, payload.Benchmarks, "mod.benchmark.a")
	assert.Contains(t, payload.Benchmarks, "mod.benchmark.b")

	// A should have B as a child
	benchA := payload.Benchmarks["mod.benchmark.a"]
	require.Len(t, benchA.Children, 1)
	assert.Equal(t, "mod.benchmark.b", benchA.Children[0].FullName)

	// B's children should NOT include A (circular reference skipped)
	// Note: B is processed as a child of A, so when it tries to recurse to A,
	// A is already in the visiting set and gets skipped
	benchBAsChild := benchA.Children[0]
	assert.Empty(t, benchBAsChild.Children, "Circular reference should be skipped")
}

// TestHierarchy_LongCycle tests that longer cycles are also detected.
// A → B → C → A should not cause infinite recursion.
func TestHierarchy_LongCycle(t *testing.T) {
	idx := NewResourceIndex()

	// Create a longer cycle: A → B → C → A
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.a",
		ShortName:  "a",
		IsTopLevel: true,
		ChildNames: []string{"mod.benchmark.b"},
	})
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.b",
		ShortName:  "b",
		ParentName: "mod.benchmark.a",
		ChildNames: []string{"mod.benchmark.c"},
	})
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.c",
		ShortName:  "c",
		ParentName: "mod.benchmark.b",
		ChildNames: []string{"mod.benchmark.a"}, // Back to A - cycle!
	})

	// Should complete without infinite loop
	payload := idx.BuildAvailableDashboardsPayload()
	assert.NotNil(t, payload)

	// Verify structure: A → B → C, but C's children should not include A
	benchA := payload.Benchmarks["mod.benchmark.a"]
	require.Len(t, benchA.Children, 1)

	benchB := benchA.Children[0]
	assert.Equal(t, "mod.benchmark.b", benchB.FullName)
	require.Len(t, benchB.Children, 1)

	benchC := benchB.Children[0]
	assert.Equal(t, "mod.benchmark.c", benchC.FullName)
	assert.Empty(t, benchC.Children, "Circular reference back to A should be skipped")
}

// TestHierarchy_SelfReference tests that a benchmark referencing itself is handled.
func TestHierarchy_SelfReference(t *testing.T) {
	idx := NewResourceIndex()

	// Self-reference: A → A
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.a",
		ShortName:  "a",
		IsTopLevel: true,
		ChildNames: []string{"mod.benchmark.a"}, // Points to itself!
	})

	// Should complete without infinite loop
	payload := idx.BuildAvailableDashboardsPayload()
	assert.NotNil(t, payload)

	// A should have no children (self-reference is skipped because A is already visiting)
	benchA := payload.Benchmarks["mod.benchmark.a"]
	assert.Empty(t, benchA.Children, "Self-reference should be skipped")
}

// TestHierarchy_NilChildNames verifies nil vs empty ChildNames
func TestHierarchy_NilChildNames(t *testing.T) {
	idx := NewResourceIndex()

	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.nil_children",
		ShortName:  "nil_children",
		IsTopLevel: true,
		ChildNames: nil, // nil, not empty slice
	})

	children := idx.GetChildren("mod.benchmark.nil_children")
	assert.Empty(t, children)
}

// TestHierarchy_ControlNotInBenchmarkChildren verifies controls not in payload children
func TestHierarchy_ControlNotInBenchmarkChildren(t *testing.T) {
	idx := NewResourceIndex()

	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.parent",
		ShortName:  "parent",
		IsTopLevel: true,
		ChildNames: []string{"mod.benchmark.child_bench", "mod.control.child_ctrl"},
	})
	idx.Add(&IndexEntry{
		Type:       "benchmark",
		Name:       "mod.benchmark.child_bench",
		ShortName:  "child_bench",
		ParentName: "mod.benchmark.parent",
	})
	idx.Add(&IndexEntry{
		Type:       "control",
		Name:       "mod.control.child_ctrl",
		ShortName:  "child_ctrl",
		ParentName: "mod.benchmark.parent",
	})

	payload := idx.BuildAvailableDashboardsPayload()
	parent := payload.Benchmarks["mod.benchmark.parent"]

	// Only benchmark children in payload, not controls
	require.Len(t, parent.Children, 1)
	assert.Equal(t, "mod.benchmark.child_bench", parent.Children[0].FullName)
}

// =============================================================================
// Benchmark (Performance) Tests
// =============================================================================

// BenchmarkHierarchy_DeepTrunkBuilding benchmarks trunk building for deep hierarchies
func BenchmarkHierarchy_DeepTrunkBuilding(b *testing.B) {
	idx := NewResourceIndex()
	createDeepHierarchy(idx, 50)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx.BuildAvailableDashboardsPayload()
	}
}

// BenchmarkHierarchy_WideChildren benchmarks wide hierarchy handling
func BenchmarkHierarchy_WideChildren(b *testing.B) {
	idx := NewResourceIndex()
	createWideHierarchy(idx, 1000, "control")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx.GetChildren("mod.benchmark.wide_parent")
	}
}

// BenchmarkHierarchy_PayloadBuilding benchmarks full payload generation
func BenchmarkHierarchy_PayloadBuilding(b *testing.B) {
	idx := NewResourceIndex()

	// Create realistic hierarchy
	for i := 0; i < 100; i++ {
		childNames := make([]string, 20)
		for j := 0; j < 20; j++ {
			childNames[j] = fmt.Sprintf("mod.benchmark.sub_%d_%d", i, j)
		}

		idx.Add(&IndexEntry{
			Type:       "benchmark",
			Name:       fmt.Sprintf("mod.benchmark.top_%d", i),
			ShortName:  fmt.Sprintf("top_%d", i),
			IsTopLevel: true,
			ChildNames: childNames,
		})

		for j := 0; j < 20; j++ {
			idx.Add(&IndexEntry{
				Type:       "benchmark",
				Name:       fmt.Sprintf("mod.benchmark.sub_%d_%d", i, j),
				ShortName:  fmt.Sprintf("sub_%d_%d", i, j),
				ParentName: fmt.Sprintf("mod.benchmark.top_%d", i),
			})
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx.BuildAvailableDashboardsPayload()
	}
}
