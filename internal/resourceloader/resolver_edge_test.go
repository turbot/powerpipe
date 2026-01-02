package resourceloader

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/powerpipe/internal/resourcecache"
	"github.com/turbot/powerpipe/internal/resourceindex"
)

// =============================================================================
// Test Helpers
// =============================================================================

// testGraphBuilder helps construct dependency graphs for testing
type testGraphBuilder struct {
	index *resourceindex.ResourceIndex
}

func newTestGraphBuilder(modName string) *testGraphBuilder {
	index := resourceindex.NewResourceIndex()
	index.ModName = modName
	return &testGraphBuilder{index: index}
}

func (b *testGraphBuilder) addResource(resType, shortName string) *testGraphBuilder {
	name := fmt.Sprintf("%s.%s.%s", b.index.ModName, resType, shortName)
	b.index.Add(&resourceindex.IndexEntry{
		Type:      resType,
		Name:      name,
		ShortName: shortName,
		ModName:   b.index.ModName,
	})
	return b
}

func (b *testGraphBuilder) addResourceWithChildren(resType, shortName string, children []string) *testGraphBuilder {
	name := fmt.Sprintf("%s.%s.%s", b.index.ModName, resType, shortName)
	b.index.Add(&resourceindex.IndexEntry{
		Type:       resType,
		Name:       name,
		ShortName:  shortName,
		ChildNames: children,
		ModName:    b.index.ModName,
	})
	return b
}

func (b *testGraphBuilder) addResourceWithQuery(resType, shortName, queryRef string) *testGraphBuilder {
	name := fmt.Sprintf("%s.%s.%s", b.index.ModName, resType, shortName)
	b.index.Add(&resourceindex.IndexEntry{
		Type:      resType,
		Name:      name,
		ShortName: shortName,
		QueryRef:  queryRef,
		ModName:   b.index.ModName,
	})
	return b
}

func (b *testGraphBuilder) addResourceWithInputs(resType, shortName string, inputs []string) *testGraphBuilder {
	name := fmt.Sprintf("%s.%s.%s", b.index.ModName, resType, shortName)
	b.index.Add(&resourceindex.IndexEntry{
		Type:       resType,
		Name:       name,
		ShortName:  shortName,
		InputNames: inputs,
		ModName:    b.index.ModName,
	})
	return b
}

func (b *testGraphBuilder) addResourceFull(resType, shortName string, children []string, queryRef string, inputs []string) *testGraphBuilder {
	name := fmt.Sprintf("%s.%s.%s", b.index.ModName, resType, shortName)
	b.index.Add(&resourceindex.IndexEntry{
		Type:       resType,
		Name:       name,
		ShortName:  shortName,
		ChildNames: children,
		QueryRef:   queryRef,
		InputNames: inputs,
		ModName:    b.index.ModName,
	})
	return b
}

func (b *testGraphBuilder) buildResolver() *DependencyResolver {
	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod(b.index.ModName, "/tmp", hcl.Range{})
	loader := NewLoader(b.index, cache, mod, "/tmp")
	return NewDependencyResolver(b.index, loader)
}

func (b *testGraphBuilder) build() *resourceindex.ResourceIndex {
	return b.index
}

// =============================================================================
// Complex Dependency Graphs Tests
// =============================================================================

// TestResolver_DiamondDependency tests the diamond dependency pattern:
//
//	  A
//	 / \
//	B   C
//	 \ /
//	  D
//
// D should appear only once in the dependency order
func TestResolver_DiamondDependency(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	// D is the leaf (query)
	builder.addResource("query", "d")

	// B and C both depend on D
	builder.addResourceWithChildren("benchmark", "b", []string{"testmod.query.d"})
	builder.addResourceWithChildren("benchmark", "c", []string{"testmod.query.d"})

	// A depends on both B and C
	builder.addResourceWithChildren("benchmark", "a", []string{
		"testmod.benchmark.b",
		"testmod.benchmark.c",
	})

	resolver := builder.buildResolver()

	// Get dependency order starting from A
	order, err := resolver.GetDependencyOrder([]string{"testmod.benchmark.a"})
	require.NoError(t, err)

	// D should appear exactly once
	dCount := 0
	for _, name := range order {
		if name == "testmod.query.d" {
			dCount++
		}
	}
	assert.Equal(t, 1, dCount, "D should appear exactly once in dependency order")

	// Verify order: D should come before B, C, and A
	dIdx := indexOf(order, "testmod.query.d")
	bIdx := indexOf(order, "testmod.benchmark.b")
	cIdx := indexOf(order, "testmod.benchmark.c")
	aIdx := indexOf(order, "testmod.benchmark.a")

	assert.Less(t, dIdx, bIdx, "D should come before B")
	assert.Less(t, dIdx, cIdx, "D should come before C")
	assert.Less(t, bIdx, aIdx, "B should come before A")
	assert.Less(t, cIdx, aIdx, "C should come before A")
}

// TestResolver_DeepChain tests a deep dependency chain (7+ levels):
// A → B → C → D → E → F → G
// Loading order should be: G, F, E, D, C, B, A
func TestResolver_DeepChain(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	// Create chain: a -> b -> c -> d -> e -> f -> g
	depth := 10
	for i := depth - 1; i >= 0; i-- {
		name := string(rune('a' + i))
		if i == depth-1 {
			// Leaf node (no children)
			builder.addResource("benchmark", name)
		} else {
			// Has child pointing to next node
			childName := string(rune('a' + i + 1))
			builder.addResourceWithChildren("benchmark", name, []string{
				fmt.Sprintf("testmod.benchmark.%s", childName),
			})
		}
	}

	resolver := builder.buildResolver()

	// Get dependency order starting from 'a'
	order, err := resolver.GetDependencyOrder([]string{"testmod.benchmark.a"})
	require.NoError(t, err)
	require.Len(t, order, depth)

	// Verify order: later letters (deeper in chain) should come first
	for i := 0; i < depth-1; i++ {
		currName := fmt.Sprintf("testmod.benchmark.%s", string(rune('a'+i)))
		nextName := fmt.Sprintf("testmod.benchmark.%s", string(rune('a'+i+1)))
		currIdx := indexOf(order, currName)
		nextIdx := indexOf(order, nextName)
		assert.Greater(t, currIdx, nextIdx, "%s should come after %s", currName, nextName)
	}
}

// TestResolver_WideDependencies tests a resource with many dependencies:
// A depends on B1, B2, B3, ..., B50
func TestResolver_WideDependencies(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	width := 50
	children := make([]string, width)

	// Create 50 leaf nodes
	for i := 0; i < width; i++ {
		name := fmt.Sprintf("b%d", i)
		builder.addResource("query", name)
		children[i] = fmt.Sprintf("testmod.query.%s", name)
	}

	// A depends on all of them
	builder.addResourceWithChildren("benchmark", "a", children)

	resolver := builder.buildResolver()

	// Get dependency order
	order, err := resolver.GetDependencyOrder([]string{"testmod.benchmark.a"})
	require.NoError(t, err)
	require.Len(t, order, width+1) // 50 children + A

	// A should come last
	aIdx := indexOf(order, "testmod.benchmark.a")
	assert.Equal(t, width, aIdx, "A should be the last element")

	// All B nodes should come before A
	for i := 0; i < width; i++ {
		bName := fmt.Sprintf("testmod.query.b%d", i)
		bIdx := indexOf(order, bName)
		assert.Less(t, bIdx, aIdx, "%s should come before A", bName)
	}
}

// TestResolver_MixedGraph tests a realistic mixed graph with varying depth and width
//
//	        root_benchmark
//	       /      |       \
//	  bench1   bench2   control1 (→ query1)
//	    |      /    \
//	control2  c3    c4 (→ query2)
//	(→ query1)
func TestResolver_MixedGraph(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	// Leaf queries
	builder.addResource("query", "query1")
	builder.addResource("query", "query2")

	// Controls that reference queries
	builder.addResourceWithQuery("control", "control1", "testmod.query.query1")
	builder.addResourceWithQuery("control", "control2", "testmod.query.query1")
	builder.addResource("control", "c3")
	builder.addResourceWithQuery("control", "c4", "testmod.query.query2")

	// Intermediate benchmarks
	builder.addResourceWithChildren("benchmark", "bench1", []string{"testmod.control.control2"})
	builder.addResourceWithChildren("benchmark", "bench2", []string{
		"testmod.control.c3",
		"testmod.control.c4",
	})

	// Root benchmark
	builder.addResourceWithChildren("benchmark", "root_benchmark", []string{
		"testmod.benchmark.bench1",
		"testmod.benchmark.bench2",
		"testmod.control.control1",
	})

	resolver := builder.buildResolver()

	// Get dependency order
	order, err := resolver.GetDependencyOrder([]string{"testmod.benchmark.root_benchmark"})
	require.NoError(t, err)

	// Verify key ordering constraints
	// query1 must come before control1, control2
	q1Idx := indexOf(order, "testmod.query.query1")
	c1Idx := indexOf(order, "testmod.control.control1")
	c2Idx := indexOf(order, "testmod.control.control2")
	assert.Less(t, q1Idx, c1Idx)
	assert.Less(t, q1Idx, c2Idx)

	// query2 must come before c4
	q2Idx := indexOf(order, "testmod.query.query2")
	c4Idx := indexOf(order, "testmod.control.c4")
	assert.Less(t, q2Idx, c4Idx)

	// bench1 must come before root
	b1Idx := indexOf(order, "testmod.benchmark.bench1")
	rootIdx := indexOf(order, "testmod.benchmark.root_benchmark")
	assert.Less(t, b1Idx, rootIdx)

	// control2 must come before bench1
	assert.Less(t, c2Idx, b1Idx)
}

// =============================================================================
// Circular Dependency Detection Tests
// =============================================================================

// TestResolver_SimpleCycle tests simple cycle detection: A → B → A
func TestResolver_SimpleCycle(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	builder.addResourceWithChildren("benchmark", "a", []string{"testmod.benchmark.b"})
	builder.addResourceWithChildren("benchmark", "b", []string{"testmod.benchmark.a"})

	resolver := builder.buildResolver()

	// HasCircularDependency should detect the cycle
	assert.True(t, resolver.HasCircularDependency("testmod.benchmark.a"))
	assert.True(t, resolver.HasCircularDependency("testmod.benchmark.b"))

	// GetDependencyOrder should return an error
	_, err := resolver.GetDependencyOrder([]string{"testmod.benchmark.a"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circular")
}

// TestResolver_LongCycle tests a longer cycle: A → B → C → D → E → A
func TestResolver_LongCycle(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	cycleLen := 5
	for i := 0; i < cycleLen; i++ {
		curr := string(rune('a' + i))
		next := string(rune('a' + (i+1)%cycleLen))
		builder.addResourceWithChildren("benchmark", curr, []string{
			fmt.Sprintf("testmod.benchmark.%s", next),
		})
	}

	resolver := builder.buildResolver()

	// All nodes in the cycle should report circular dependency
	for i := 0; i < cycleLen; i++ {
		name := fmt.Sprintf("testmod.benchmark.%s", string(rune('a'+i)))
		assert.True(t, resolver.HasCircularDependency(name), "cycle should be detected from %s", name)
	}

	// GetDependencyOrder should detect the cycle
	names := make([]string, cycleLen)
	for i := 0; i < cycleLen; i++ {
		names[i] = fmt.Sprintf("testmod.benchmark.%s", string(rune('a'+i)))
	}
	_, err := resolver.GetDependencyOrder(names)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circular")
}

// TestResolver_CycleInSubtree tests a cycle in a subtree that doesn't involve the root:
// root → A → B, but B → C → D → C (cycle in subtree)
func TestResolver_CycleInSubtree(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	// Cycle: c → d → c
	builder.addResourceWithChildren("benchmark", "c", []string{"testmod.benchmark.d"})
	builder.addResourceWithChildren("benchmark", "d", []string{"testmod.benchmark.c"})

	// b depends on c (connects to cycle)
	builder.addResourceWithChildren("benchmark", "b", []string{"testmod.benchmark.c"})

	// a depends on b
	builder.addResourceWithChildren("benchmark", "a", []string{"testmod.benchmark.b"})

	// root depends on a
	builder.addResourceWithChildren("benchmark", "root", []string{"testmod.benchmark.a"})

	resolver := builder.buildResolver()

	// Root should detect the cycle in its subtree
	assert.True(t, resolver.HasCircularDependency("testmod.benchmark.root"))

	// GetDependencyOrder should fail
	_, err := resolver.GetDependencyOrder([]string{"testmod.benchmark.root"})
	assert.Error(t, err)
}

// TestResolver_MultipleCycles tests multiple independent cycles:
// A → B → A (cycle 1)
// C → D → C (cycle 2, independent)
func TestResolver_MultipleCycles(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	// Cycle 1: a ↔ b
	builder.addResourceWithChildren("benchmark", "a", []string{"testmod.benchmark.b"})
	builder.addResourceWithChildren("benchmark", "b", []string{"testmod.benchmark.a"})

	// Cycle 2: c ↔ d (independent)
	builder.addResourceWithChildren("benchmark", "c", []string{"testmod.benchmark.d"})
	builder.addResourceWithChildren("benchmark", "d", []string{"testmod.benchmark.c"})

	resolver := builder.buildResolver()

	// Each cycle should be detected independently
	assert.True(t, resolver.HasCircularDependency("testmod.benchmark.a"))
	assert.True(t, resolver.HasCircularDependency("testmod.benchmark.c"))

	// Resolving either should fail
	_, err1 := resolver.GetDependencyOrder([]string{"testmod.benchmark.a"})
	assert.Error(t, err1)

	_, err2 := resolver.GetDependencyOrder([]string{"testmod.benchmark.c"})
	assert.Error(t, err2)
}

// TestResolver_SelfReference tests direct self-reference: A → A
func TestResolver_SelfReference(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	builder.addResourceWithChildren("benchmark", "a", []string{"testmod.benchmark.a"})

	resolver := builder.buildResolver()

	// Self-reference should be detected as circular
	assert.True(t, resolver.HasCircularDependency("testmod.benchmark.a"))

	// GetDependencyOrder should fail
	_, err := resolver.GetDependencyOrder([]string{"testmod.benchmark.a"})
	assert.Error(t, err)
}

// TestResolver_NoCycleFalsePositive tests that a complex DAG doesn't false-positive as a cycle
//
//	  A
//	 /|\
//	B C D
//	|/ \|
//	E   F
//	 \ /
//	  G
//
// Multiple paths to G, but NO cycle
func TestResolver_NoCycleFalsePositive(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	// Leaf
	builder.addResource("query", "g")

	// E and F both depend on G
	builder.addResourceWithChildren("benchmark", "e", []string{"testmod.query.g"})
	builder.addResourceWithChildren("benchmark", "f", []string{"testmod.query.g"})

	// B depends on E
	builder.addResourceWithChildren("benchmark", "b", []string{"testmod.benchmark.e"})

	// C depends on E
	builder.addResourceWithChildren("benchmark", "c", []string{"testmod.benchmark.e"})

	// D depends on F
	builder.addResourceWithChildren("benchmark", "d", []string{"testmod.benchmark.f"})

	// A depends on B, C, D
	builder.addResourceWithChildren("benchmark", "a", []string{
		"testmod.benchmark.b",
		"testmod.benchmark.c",
		"testmod.benchmark.d",
	})

	resolver := builder.buildResolver()

	// Should NOT be detected as circular
	assert.False(t, resolver.HasCircularDependency("testmod.benchmark.a"))

	// GetDependencyOrder should succeed
	order, err := resolver.GetDependencyOrder([]string{"testmod.benchmark.a"})
	require.NoError(t, err)
	assert.NotEmpty(t, order)

	// G should come first since it's the deepest leaf
	assert.Equal(t, "testmod.query.g", order[0])
}

// =============================================================================
// Missing Dependencies Tests
// =============================================================================

// TestResolver_MissingDirectDep tests handling of missing direct dependency
func TestResolver_MissingDirectDep(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	// Control references a query that doesn't exist
	builder.addResourceWithQuery("control", "orphan", "testmod.query.nonexistent")

	resolver := builder.buildResolver()

	// GetDependencies should return the missing reference
	deps := resolver.GetDependencies("testmod.control.orphan")
	require.Len(t, deps, 1)
	assert.Equal(t, "testmod.query.nonexistent", deps[0].To)

	// GetDependencyOrder includes missing dependency names in the graph
	// (they act as leaf nodes with no dependencies of their own)
	order, err := resolver.GetDependencyOrder([]string{"testmod.control.orphan"})
	require.NoError(t, err)
	assert.Len(t, order, 2) // orphan + nonexistent placeholder

	// The missing dependency should come before orphan (leaf first)
	nonexistentIdx := indexOf(order, "testmod.query.nonexistent")
	orphanIdx := indexOf(order, "testmod.control.orphan")
	assert.Less(t, nonexistentIdx, orphanIdx, "missing dep should come before orphan")
}

// TestResolver_MissingTransitiveDep tests missing transitive dependency:
// A → B → C (C missing)
func TestResolver_MissingTransitiveDep(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	// B references missing C
	builder.addResourceWithChildren("benchmark", "b", []string{"testmod.benchmark.c"})

	// A depends on B
	builder.addResourceWithChildren("benchmark", "a", []string{"testmod.benchmark.b"})

	resolver := builder.buildResolver()

	// GetTransitiveDependencies should still work and include missing deps
	deps := resolver.GetTransitiveDependencies("testmod.benchmark.a")

	// Should include A, B, and C (missing deps are included as placeholders)
	assert.Contains(t, deps, "testmod.benchmark.a")
	assert.Contains(t, deps, "testmod.benchmark.b")
	assert.Contains(t, deps, "testmod.benchmark.c")

	// GetDependencyOrder includes missing deps as leaf nodes
	order, err := resolver.GetDependencyOrder([]string{"testmod.benchmark.a"})
	require.NoError(t, err)
	assert.Len(t, order, 3) // a, b, and c (missing but referenced)

	// Missing C should come first (leaf), then B, then A
	cIdx := indexOf(order, "testmod.benchmark.c")
	bIdx := indexOf(order, "testmod.benchmark.b")
	aIdx := indexOf(order, "testmod.benchmark.a")
	assert.Less(t, cIdx, bIdx, "C should come before B")
	assert.Less(t, bIdx, aIdx, "B should come before A")
}

// TestResolver_OptionalDepMissing tests that optional/missing deps are handled gracefully
func TestResolver_OptionalDepMissing(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	// Dashboard with input that doesn't exist
	builder.addResourceWithInputs("dashboard", "dash", []string{"testmod.input.missing"})

	resolver := builder.buildResolver()

	// Should report the dependency
	deps := resolver.GetDependencies("testmod.dashboard.dash")
	require.Len(t, deps, 1)
	assert.Equal(t, DepInput, deps[0].Type)

	// Ordering includes missing deps as placeholders
	order, err := resolver.GetDependencyOrder([]string{"testmod.dashboard.dash"})
	require.NoError(t, err)
	assert.Len(t, order, 2) // dash + missing input

	// Missing input should come before dashboard
	inputIdx := indexOf(order, "testmod.input.missing")
	dashIdx := indexOf(order, "testmod.dashboard.dash")
	assert.Less(t, inputIdx, dashIdx)
}

// TestResolver_InlineDependency tests that inline resources are handled correctly
// Inline resources are defined within their parent and shouldn't be loaded separately
// NOTE: The resolver includes all referenced names in the dependency order, even if
// they don't exist in the index (they appear as leaf nodes). This is by design -
// it allows the loader to handle inline resources appropriately.
func TestResolver_InlineDependency(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	// Dashboard with child that isn't in the index (inline)
	builder.addResourceWithChildren("dashboard", "dash", []string{
		"testmod.card.inline_card", // This card is defined inline, not in index
	})

	resolver := builder.buildResolver()

	// GetDependencies returns the reference even if not in index
	deps := resolver.GetDependencies("testmod.dashboard.dash")
	require.Len(t, deps, 1)

	// GetDependencyOrder includes the reference (inline resources act as leaf nodes)
	// The loader handles these appropriately during actual resource loading
	order, err := resolver.GetDependencyOrder([]string{"testmod.dashboard.dash"})
	require.NoError(t, err)
	assert.Len(t, order, 2) // dash + inline_card placeholder

	// Inline card reference comes before dashboard
	cardIdx := indexOf(order, "testmod.card.inline_card")
	dashIdx := indexOf(order, "testmod.dashboard.dash")
	assert.Less(t, cardIdx, dashIdx)
}

// =============================================================================
// Cross-Mod Dependencies Tests
// =============================================================================

// TestResolver_CrossModReference tests reference to a dependency mod resource
func TestResolver_CrossModReference(t *testing.T) {
	index := resourceindex.NewResourceIndex()
	index.ModName = "mainmod"

	// Add resource from dependency mod
	index.Add(&resourceindex.IndexEntry{
		Type:      "query",
		Name:      "depmod.query.shared_query",
		ShortName: "shared_query",
		ModName:   "depmod",
	})

	// Main mod control references dep mod query
	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "mainmod.control.uses_dep",
		ShortName: "uses_dep",
		QueryRef:  "depmod.query.shared_query",
		ModName:   "mainmod",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("mainmod", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	// GetDependencies should find cross-mod reference
	deps := resolver.GetDependencies("mainmod.control.uses_dep")
	require.Len(t, deps, 1)
	assert.Equal(t, "depmod.query.shared_query", deps[0].To)

	// GetDependencyOrder should include both
	order, err := resolver.GetDependencyOrder([]string{"mainmod.control.uses_dep"})
	require.NoError(t, err)
	require.Len(t, order, 2)

	// Query should come before control
	qIdx := indexOf(order, "depmod.query.shared_query")
	cIdx := indexOf(order, "mainmod.control.uses_dep")
	assert.Less(t, qIdx, cIdx)
}

// TestResolver_CrossModCycle tests circular dependency detection across mods:
// mainmod.A → depmod.B → mainmod.C → mainmod.A
func TestResolver_CrossModCycle(t *testing.T) {
	index := resourceindex.NewResourceIndex()
	index.ModName = "mainmod"

	// mainmod.a → depmod.b
	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "mainmod.benchmark.a",
		ShortName:  "a",
		ChildNames: []string{"depmod.benchmark.b"},
		ModName:    "mainmod",
	})

	// depmod.b → mainmod.c
	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "depmod.benchmark.b",
		ShortName:  "b",
		ChildNames: []string{"mainmod.benchmark.c"},
		ModName:    "depmod",
	})

	// mainmod.c → mainmod.a (creates cycle)
	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "mainmod.benchmark.c",
		ShortName:  "c",
		ChildNames: []string{"mainmod.benchmark.a"},
		ModName:    "mainmod",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("mainmod", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	// Should detect the cross-mod cycle
	assert.True(t, resolver.HasCircularDependency("mainmod.benchmark.a"))

	// GetDependencyOrder should fail
	_, err := resolver.GetDependencyOrder([]string{"mainmod.benchmark.a"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circular")
}

// TestResolver_ModNameResolution tests that short and full names resolve correctly
func TestResolver_ModNameResolution(t *testing.T) {
	index := resourceindex.NewResourceIndex()
	index.ModName = "mymod"

	// Add entries with full mod names
	index.Add(&resourceindex.IndexEntry{
		Type:        "query",
		Name:        "mymod.query.q1",
		ShortName:   "q1",
		ModName:     "mymod",
		ModFullName: "github.com/org/mymod",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "mymod.control.c1",
		ShortName: "c1",
		QueryRef:  "mymod.query.q1", // Uses short mod name
		ModName:   "mymod",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("mymod", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	// Should resolve short name correctly
	deps := resolver.GetDependencies("mymod.control.c1")
	require.Len(t, deps, 1)
	assert.Equal(t, "mymod.query.q1", deps[0].To)
}

// =============================================================================
// Dependency Types Tests
// =============================================================================

// TestResolver_ChildDependencies tests child dependency resolution
func TestResolver_ChildDependencies(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	builder.addResource("control", "c1")
	builder.addResource("control", "c2")
	builder.addResource("control", "c3")

	builder.addResourceWithChildren("benchmark", "parent", []string{
		"testmod.control.c1",
		"testmod.control.c2",
		"testmod.control.c3",
	})

	resolver := builder.buildResolver()

	deps := resolver.GetDependencies("testmod.benchmark.parent")
	require.Len(t, deps, 3)

	for _, dep := range deps {
		assert.Equal(t, DepChild, dep.Type)
		assert.Equal(t, "testmod.benchmark.parent", dep.From)
	}
}

// TestResolver_QueryReferenceDep tests query reference dependency
func TestResolver_QueryReferenceDep(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	builder.addResource("query", "shared")
	builder.addResourceWithQuery("control", "c1", "testmod.query.shared")

	resolver := builder.buildResolver()

	deps := resolver.GetDependencies("testmod.control.c1")
	require.Len(t, deps, 1)
	assert.Equal(t, DepQuery, deps[0].Type)
	assert.Equal(t, "testmod.query.shared", deps[0].To)

	// Query should be loaded before control
	order, err := resolver.GetDependencyOrder([]string{"testmod.control.c1"})
	require.NoError(t, err)
	assert.Equal(t, "testmod.query.shared", order[0])
	assert.Equal(t, "testmod.control.c1", order[1])
}

// TestResolver_InputDependency tests input dependency for dashboards
func TestResolver_InputDependency(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	builder.addResource("input", "region")
	builder.addResource("input", "account")
	builder.addResourceWithInputs("dashboard", "dash", []string{
		"testmod.input.region",
		"testmod.input.account",
	})

	resolver := builder.buildResolver()

	deps := resolver.GetDependencies("testmod.dashboard.dash")
	require.Len(t, deps, 2)

	for _, dep := range deps {
		assert.Equal(t, DepInput, dep.Type)
	}

	// Inputs should come before dashboard
	order, err := resolver.GetDependencyOrder([]string{"testmod.dashboard.dash"})
	require.NoError(t, err)

	dashIdx := indexOf(order, "testmod.dashboard.dash")
	assert.Equal(t, 2, dashIdx, "dashboard should be last")
}

// TestResolver_MixedDepTypes tests resource with multiple dependency types
func TestResolver_MixedDepTypes(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	// Create various dependency targets
	builder.addResource("query", "q1")
	builder.addResource("input", "input1")
	builder.addResource("control", "child1")

	// Dashboard with children, query, and inputs
	builder.addResourceFull("dashboard", "complex",
		[]string{"testmod.control.child1"},
		"testmod.query.q1",
		[]string{"testmod.input.input1"},
	)

	resolver := builder.buildResolver()

	deps := resolver.GetDependencies("testmod.dashboard.complex")
	require.Len(t, deps, 3)

	// Verify we have all three dependency types
	hasChild := false
	hasQuery := false
	hasInput := false
	for _, dep := range deps {
		switch dep.Type {
		case DepChild:
			hasChild = true
		case DepQuery:
			hasQuery = true
		case DepInput:
			hasInput = true
		}
	}
	assert.True(t, hasChild, "should have child dependency")
	assert.True(t, hasQuery, "should have query dependency")
	assert.True(t, hasInput, "should have input dependency")
}

// =============================================================================
// Loading Order Tests
// =============================================================================

// TestResolver_TopologicalSort tests that topological sort is correct
func TestResolver_TopologicalSort(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	// Create a graph where each resource must come after its dependencies
	// q1 has no deps
	builder.addResource("query", "q1")
	// c1 depends on q1
	builder.addResourceWithQuery("control", "c1", "testmod.query.q1")
	// b1 depends on c1
	builder.addResourceWithChildren("benchmark", "b1", []string{"testmod.control.c1"})
	// root depends on b1
	builder.addResourceWithChildren("benchmark", "root", []string{"testmod.benchmark.b1"})

	resolver := builder.buildResolver()

	order, err := resolver.GetDependencyOrder([]string{"testmod.benchmark.root"})
	require.NoError(t, err)

	// For each resource, verify all its dependencies come before it
	for i, name := range order {
		deps := resolver.GetDependencies(name)
		for _, dep := range deps {
			depIdx := indexOf(order, dep.To)
			if depIdx == -1 {
				continue // Dependency not in index (missing)
			}
			assert.Less(t, depIdx, i, "dependency %s should come before %s", dep.To, name)
		}
	}
}

// TestResolver_OrderStability tests that the topological order is always valid
// Note: The exact order of independent nodes (e.g., two queries with no dependencies)
// may vary due to Go map iteration order. This test verifies that dependencies
// always come before their dependents, which is the invariant that matters.
func TestResolver_OrderStability(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	// Create a moderately complex graph
	builder.addResource("query", "q1")
	builder.addResource("query", "q2")
	builder.addResourceWithQuery("control", "c1", "testmod.query.q1")
	builder.addResourceWithQuery("control", "c2", "testmod.query.q2")
	builder.addResourceWithChildren("benchmark", "b1", []string{
		"testmod.control.c1",
		"testmod.control.c2",
	})

	resolver := builder.buildResolver()

	// Run multiple times and verify the topological property is maintained
	for i := 0; i < 10; i++ {
		order, err := resolver.GetDependencyOrder([]string{"testmod.benchmark.b1"})
		require.NoError(t, err)
		require.Len(t, order, 5)

		// Verify topological property: all dependencies come before dependents
		// q1 before c1
		q1Idx := indexOf(order, "testmod.query.q1")
		c1Idx := indexOf(order, "testmod.control.c1")
		assert.Less(t, q1Idx, c1Idx, "q1 should come before c1 (iteration %d)", i)

		// q2 before c2
		q2Idx := indexOf(order, "testmod.query.q2")
		c2Idx := indexOf(order, "testmod.control.c2")
		assert.Less(t, q2Idx, c2Idx, "q2 should come before c2 (iteration %d)", i)

		// c1 and c2 before b1
		b1Idx := indexOf(order, "testmod.benchmark.b1")
		assert.Less(t, c1Idx, b1Idx, "c1 should come before b1 (iteration %d)", i)
		assert.Less(t, c2Idx, b1Idx, "c2 should come before b1 (iteration %d)", i)
	}
}

// TestResolver_ParallelSafeOrder tests that independent deps can be identified for parallel loading
func TestResolver_ParallelSafeOrder(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	// Create resources that are independent (can be loaded in parallel)
	// q1, q2, q3 are all independent
	builder.addResource("query", "q1")
	builder.addResource("query", "q2")
	builder.addResource("query", "q3")

	// bench depends on all three
	builder.addResourceWithChildren("benchmark", "bench", []string{
		"testmod.query.q1",
		"testmod.query.q2",
		"testmod.query.q3",
	})

	resolver := builder.buildResolver()

	order, err := resolver.GetDependencyOrder([]string{"testmod.benchmark.bench"})
	require.NoError(t, err)

	// All queries should come before benchmark
	// The queries are independent and could be loaded in parallel
	benchIdx := indexOf(order, "testmod.benchmark.bench")
	for _, qName := range []string{"testmod.query.q1", "testmod.query.q2", "testmod.query.q3"} {
		qIdx := indexOf(order, qName)
		assert.Less(t, qIdx, benchIdx, "%s should come before benchmark", qName)
	}
}

// =============================================================================
// Performance Tests
// =============================================================================

// TestResolver_LargeGraphPerformance tests that large graphs resolve quickly
func TestResolver_LargeGraphPerformance(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	// Create a large graph with 1000 nodes
	nodeCount := 1000

	// Create leaf nodes (500 queries)
	for i := 0; i < nodeCount/2; i++ {
		builder.addResource("query", fmt.Sprintf("q%d", i))
	}

	// Create controls that reference queries (500 controls)
	for i := 0; i < nodeCount/2; i++ {
		builder.addResourceWithQuery("control", fmt.Sprintf("c%d", i),
			fmt.Sprintf("testmod.query.q%d", i%250))
	}

	// Create benchmarks that contain controls
	for i := 0; i < 50; i++ {
		children := make([]string, 10)
		for j := 0; j < 10; j++ {
			children[j] = fmt.Sprintf("testmod.control.c%d", i*10+j)
		}
		builder.addResourceWithChildren("benchmark", fmt.Sprintf("b%d", i), children)
	}

	// Root benchmark
	allBenchmarks := make([]string, 50)
	for i := 0; i < 50; i++ {
		allBenchmarks[i] = fmt.Sprintf("testmod.benchmark.b%d", i)
	}
	builder.addResourceWithChildren("benchmark", "root", allBenchmarks)

	resolver := builder.buildResolver()

	// Measure resolution time
	start := time.Now()
	order, err := resolver.GetDependencyOrder([]string{"testmod.benchmark.root"})
	duration := time.Since(start)

	require.NoError(t, err)
	assert.NotEmpty(t, order)

	// Should complete in well under 1 second
	assert.Less(t, duration, time.Second, "large graph resolution should be fast, took %v", duration)

	t.Logf("Large graph (%d nodes) resolved in %v", len(order), duration)
}

// TestResolver_LargeGraphMemory tests that large graph doesn't use excessive memory
func TestResolver_LargeGraphMemory(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	// Create a graph designed to test visited set efficiency
	// Deep chain + wide branches
	depth := 100
	width := 50

	// Create a chain of depth 100
	for i := depth - 1; i >= 0; i-- {
		name := fmt.Sprintf("chain%d", i)
		if i == depth-1 {
			builder.addResource("benchmark", name)
		} else {
			builder.addResourceWithChildren("benchmark", name, []string{
				fmt.Sprintf("testmod.benchmark.chain%d", i+1),
			})
		}
	}

	// Add wide branches at various points
	for i := 0; i < depth; i += 10 {
		branchChildren := make([]string, width)
		for j := 0; j < width; j++ {
			leafName := fmt.Sprintf("leaf_%d_%d", i, j)
			builder.addResource("query", leafName)
			branchChildren[j] = fmt.Sprintf("testmod.query.%s", leafName)
		}

		// Update chain node to include branch
		chainName := fmt.Sprintf("testmod.benchmark.chain%d", i)
		entry, _ := builder.index.Get(chainName)
		if entry != nil {
			entry.ChildNames = append(entry.ChildNames, branchChildren...)
		}
	}

	resolver := builder.buildResolver()

	// This should complete without stack overflow or memory issues
	order, err := resolver.GetDependencyOrder([]string{"testmod.benchmark.chain0"})
	require.NoError(t, err)
	assert.NotEmpty(t, order)
}

// =============================================================================
// Error Messages Tests
// =============================================================================

// TestResolver_CircularErrorMessage tests that circular dependency error is descriptive
func TestResolver_CircularErrorMessage(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	// Create cycle: a → b → c → a
	builder.addResourceWithChildren("benchmark", "a", []string{"testmod.benchmark.b"})
	builder.addResourceWithChildren("benchmark", "b", []string{"testmod.benchmark.c"})
	builder.addResourceWithChildren("benchmark", "c", []string{"testmod.benchmark.a"})

	resolver := builder.buildResolver()

	_, err := resolver.GetDependencyOrder([]string{"testmod.benchmark.a"})
	require.Error(t, err)

	// Error should mention "circular"
	assert.Contains(t, err.Error(), "circular")

	// Error should include the involved resources
	errStr := err.Error()
	// At least some of the cycle members should be mentioned
	containsAny := strings.Contains(errStr, "benchmark.a") ||
		strings.Contains(errStr, "benchmark.b") ||
		strings.Contains(errStr, "benchmark.c")
	assert.True(t, containsAny, "error should mention cycle members: %s", errStr)
}

// TestResolver_MissingDepErrorMessage tests error messages for missing dependencies
func TestResolver_MissingDepErrorMessage(t *testing.T) {
	index := resourceindex.NewResourceIndex()
	index.ModName = "testmod"

	// Control that references non-existent query
	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "testmod.control.orphan",
		ShortName: "orphan",
		QueryRef:  "testmod.query.nonexistent",
		ModName:   "testmod",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("testmod", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	// ResolveWithDependencies for a missing resource should error
	ctx := context.Background()
	err := resolver.ResolveWithDependencies(ctx, "testmod.query.nonexistent")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// =============================================================================
// Edge Cases
// =============================================================================

// TestResolver_EmptyGraph tests resolver with no resources
func TestResolver_EmptyGraph(t *testing.T) {
	builder := newTestGraphBuilder("testmod")
	resolver := builder.buildResolver()

	// Empty input should return empty output
	order, err := resolver.GetDependencyOrder([]string{})
	require.NoError(t, err)
	assert.Empty(t, order)

	// Non-existent resource
	deps := resolver.GetDependencies("testmod.benchmark.nonexistent")
	assert.Empty(t, deps)
}

// TestResolver_SingleNode tests resolver with a single node
func TestResolver_SingleNode(t *testing.T) {
	builder := newTestGraphBuilder("testmod")
	builder.addResource("query", "alone")

	resolver := builder.buildResolver()

	order, err := resolver.GetDependencyOrder([]string{"testmod.query.alone"})
	require.NoError(t, err)
	assert.Len(t, order, 1)
	assert.Equal(t, "testmod.query.alone", order[0])
}

// TestResolver_DisconnectedComponents tests multiple disconnected subgraphs
func TestResolver_DisconnectedComponents(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	// Component 1: a1 → b1
	builder.addResource("query", "b1")
	builder.addResourceWithChildren("benchmark", "a1", []string{"testmod.query.b1"})

	// Component 2: a2 → b2 (disconnected from component 1)
	builder.addResource("query", "b2")
	builder.addResourceWithChildren("benchmark", "a2", []string{"testmod.query.b2"})

	resolver := builder.buildResolver()

	// Request both roots
	order, err := resolver.GetDependencyOrder([]string{
		"testmod.benchmark.a1",
		"testmod.benchmark.a2",
	})
	require.NoError(t, err)
	assert.Len(t, order, 4)

	// Each component should maintain proper order
	b1Idx := indexOf(order, "testmod.query.b1")
	a1Idx := indexOf(order, "testmod.benchmark.a1")
	assert.Less(t, b1Idx, a1Idx)

	b2Idx := indexOf(order, "testmod.query.b2")
	a2Idx := indexOf(order, "testmod.benchmark.a2")
	assert.Less(t, b2Idx, a2Idx)
}

// TestResolver_DuplicateInputNames tests handling of duplicate names in input
func TestResolver_DuplicateInputNames(t *testing.T) {
	builder := newTestGraphBuilder("testmod")
	builder.addResource("query", "q1")
	builder.addResourceWithChildren("benchmark", "b1", []string{"testmod.query.q1"})

	resolver := builder.buildResolver()

	// Pass duplicate names
	order, err := resolver.GetDependencyOrder([]string{
		"testmod.benchmark.b1",
		"testmod.benchmark.b1", // duplicate
		"testmod.query.q1",     // already included as dep
	})
	require.NoError(t, err)

	// Should deduplicate
	assert.Len(t, order, 2)

	// Count occurrences
	counts := make(map[string]int)
	for _, name := range order {
		counts[name]++
	}
	assert.Equal(t, 1, counts["testmod.benchmark.b1"])
	assert.Equal(t, 1, counts["testmod.query.q1"])
}

// TestResolver_GetDependents_MultipleRefs tests GetDependents with resources referenced by many
func TestResolver_GetDependents_MultipleRefs(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	// Shared query used by many controls
	builder.addResource("query", "shared")

	numControls := 20
	for i := 0; i < numControls; i++ {
		builder.addResourceWithQuery("control", fmt.Sprintf("c%d", i), "testmod.query.shared")
	}

	resolver := builder.buildResolver()

	dependents := resolver.GetDependents("testmod.query.shared")
	assert.Len(t, dependents, numControls)

	// Verify all controls are listed
	sort.Strings(dependents)
	for i := 0; i < numControls; i++ {
		expected := fmt.Sprintf("testmod.control.c%d", i)
		assert.Contains(t, dependents, expected)
	}
}

// TestResolver_BuildDependencyGraph_Complete tests complete graph building
func TestResolver_BuildDependencyGraph_Complete(t *testing.T) {
	builder := newTestGraphBuilder("testmod")

	builder.addResource("query", "q1")
	builder.addResource("query", "q2")
	builder.addResourceWithQuery("control", "c1", "testmod.query.q1")
	builder.addResourceWithQuery("control", "c2", "testmod.query.q2")
	builder.addResourceWithChildren("benchmark", "b1", []string{
		"testmod.control.c1",
		"testmod.control.c2",
	})

	resolver := builder.buildResolver()

	graph := resolver.BuildDependencyGraph()
	assert.NotNil(t, graph)
	assert.Len(t, graph.Nodes, 5) // 2 queries + 2 controls + 1 benchmark

	// Verify benchmark has 2 child deps
	assert.Len(t, graph.Nodes["testmod.benchmark.b1"], 2)

	// Verify controls have query deps
	assert.Len(t, graph.Nodes["testmod.control.c1"], 1)
	assert.Len(t, graph.Nodes["testmod.control.c2"], 1)

	// Verify queries have no deps
	assert.Len(t, graph.Nodes["testmod.query.q1"], 0)
	assert.Len(t, graph.Nodes["testmod.query.q2"], 0)
}
