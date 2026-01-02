package resourceloader

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/powerpipe/internal/resourcecache"
	"github.com/turbot/powerpipe/internal/resourceindex"
)

// =============================================================================
// Transitive Dependency Resolution Tests
// =============================================================================

func TestTransitiveDeps_TwoLevel(t *testing.T) {
	// Test: Main -> DepA -> DepLeaf (two-level transitive)
	index := resourceindex.NewResourceIndex()

	// Main mod benchmark -> control -> query chain
	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "main.benchmark.top",
		ShortName:  "top",
		ChildNames: []string{"dep_a.control.middle"},
		ModName:    "main",
	})

	// Dep A control references dep_leaf query
	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "dep_a.control.middle",
		ShortName: "middle",
		QueryRef:  "dep_leaf.query.bottom",
		ModName:   "dep_a",
	})

	// Dep leaf query (bottom of chain)
	index.Add(&resourceindex.IndexEntry{
		Type:      "query",
		Name:      "dep_leaf.query.bottom",
		ShortName: "bottom",
		HasSQL:    true,
		ModName:   "dep_leaf",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("main", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	// Get transitive dependencies of the benchmark
	deps := resolver.GetTransitiveDependencies("main.benchmark.top")

	// Should include all three: query, control, benchmark
	assert.Len(t, deps, 3, "Should have 3 dependencies in chain")

	// Verify order: leaf first, then middle, then top
	queryIdx := indexOf(deps, "dep_leaf.query.bottom")
	controlIdx := indexOf(deps, "dep_a.control.middle")
	benchmarkIdx := indexOf(deps, "main.benchmark.top")

	assert.Less(t, queryIdx, controlIdx, "Query should come before control")
	assert.Less(t, controlIdx, benchmarkIdx, "Control should come before benchmark")
}

func TestTransitiveDeps_ThreeLevel(t *testing.T) {
	// Test: A -> B -> C -> D (three-level transitive)
	index := resourceindex.NewResourceIndex()

	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "mod_a.benchmark.a",
		ShortName:  "a",
		ChildNames: []string{"mod_b.benchmark.b"},
		ModName:    "mod_a",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "mod_b.benchmark.b",
		ShortName:  "b",
		ChildNames: []string{"mod_c.benchmark.c"},
		ModName:    "mod_b",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "mod_c.benchmark.c",
		ShortName:  "c",
		ChildNames: []string{"mod_d.control.d"},
		ModName:    "mod_c",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "mod_d.control.d",
		ShortName: "d",
		HasSQL:    true,
		ModName:   "mod_d",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("mod_a", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	deps := resolver.GetTransitiveDependencies("mod_a.benchmark.a")

	// Should include all four
	assert.Len(t, deps, 4)

	// Verify order (D should come first as it's the leaf)
	dIdx := indexOf(deps, "mod_d.control.d")
	cIdx := indexOf(deps, "mod_c.benchmark.c")
	bIdx := indexOf(deps, "mod_b.benchmark.b")
	aIdx := indexOf(deps, "mod_a.benchmark.a")

	assert.Less(t, dIdx, cIdx)
	assert.Less(t, cIdx, bIdx)
	assert.Less(t, bIdx, aIdx)
}

// =============================================================================
// Cross-Mod Dependency Tests
// =============================================================================

func TestTransitiveDeps_CrossModDependencies(t *testing.T) {
	// Test: Resources from different mods in dependency chain
	index := resourceindex.NewResourceIndex()

	// Main mod control using dep mod query
	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "main.control.cross_mod",
		ShortName: "cross_mod",
		QueryRef:  "dep.query.shared",
		ModName:   "main",
	})

	// Dep mod query
	index.Add(&resourceindex.IndexEntry{
		Type:      "query",
		Name:      "dep.query.shared",
		ShortName: "shared",
		HasSQL:    true,
		ModName:   "dep",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("main", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	deps := resolver.GetDependencies("main.control.cross_mod")

	require.Len(t, deps, 1)
	assert.Equal(t, "dep.query.shared", deps[0].To)
	assert.Equal(t, DepQuery, deps[0].Type)
}

func TestTransitiveDeps_MultipleCrossModRefs(t *testing.T) {
	// Test: Benchmark with children from multiple mods
	index := resourceindex.NewResourceIndex()

	index.Add(&resourceindex.IndexEntry{
		Type:      "benchmark",
		Name:      "main.benchmark.multi",
		ShortName: "multi",
		ChildNames: []string{
			"mod_a.control.control_a",
			"mod_b.control.control_b",
			"mod_c.control.control_c",
		},
		ModName: "main",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "mod_a.control.control_a",
		ShortName: "control_a",
		HasSQL:    true,
		ModName:   "mod_a",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "mod_b.control.control_b",
		ShortName: "control_b",
		HasSQL:    true,
		ModName:   "mod_b",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "mod_c.control.control_c",
		ShortName: "control_c",
		HasSQL:    true,
		ModName:   "mod_c",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("main", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	deps := resolver.GetDependencies("main.benchmark.multi")

	// Should have 3 children from 3 different mods
	assert.Len(t, deps, 3)

	childNames := make(map[string]bool)
	for _, dep := range deps {
		childNames[dep.To] = true
	}

	assert.True(t, childNames["mod_a.control.control_a"])
	assert.True(t, childNames["mod_b.control.control_b"])
	assert.True(t, childNames["mod_c.control.control_c"])
}

// =============================================================================
// Diamond Dependency Tests
// =============================================================================

func TestTransitiveDeps_DiamondPattern(t *testing.T) {
	// Test: Main -> Left, Right; Left -> Shared; Right -> Shared
	index := resourceindex.NewResourceIndex()

	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "main.benchmark.top",
		ShortName:  "top",
		ChildNames: []string{"left.control.left", "right.control.right"},
		ModName:    "main",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "left.control.left",
		ShortName: "left",
		QueryRef:  "shared.query.shared",
		ModName:   "left",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "right.control.right",
		ShortName: "right",
		QueryRef:  "shared.query.shared",
		ModName:   "right",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "query",
		Name:      "shared.query.shared",
		ShortName: "shared",
		HasSQL:    true,
		ModName:   "shared",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("main", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	deps := resolver.GetTransitiveDependencies("main.benchmark.top")

	// Should include all 4 unique resources (shared only once)
	assert.Len(t, deps, 4, "Diamond should include 4 unique resources")

	// Verify shared is only included once
	sharedCount := 0
	for _, dep := range deps {
		if dep == "shared.query.shared" {
			sharedCount++
		}
	}
	assert.Equal(t, 1, sharedCount, "Shared resource should appear exactly once")

	// Shared should come before left and right
	sharedIdx := indexOf(deps, "shared.query.shared")
	leftIdx := indexOf(deps, "left.control.left")
	rightIdx := indexOf(deps, "right.control.right")
	topIdx := indexOf(deps, "main.benchmark.top")

	assert.Less(t, sharedIdx, leftIdx, "Shared should come before left")
	assert.Less(t, sharedIdx, rightIdx, "Shared should come before right")
	assert.Less(t, leftIdx, topIdx, "Left should come before top")
	assert.Less(t, rightIdx, topIdx, "Right should come before top")
}

func TestTransitiveDeps_DiamondDependencyOrder(t *testing.T) {
	// Test: GetDependencyOrder handles diamond correctly
	index := resourceindex.NewResourceIndex()

	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "main.benchmark.diamond",
		ShortName:  "diamond",
		ChildNames: []string{"left.benchmark.left", "right.benchmark.right"},
		ModName:    "main",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "left.benchmark.left",
		ShortName:  "left",
		ChildNames: []string{"shared.control.shared"},
		ModName:    "left",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "right.benchmark.right",
		ShortName:  "right",
		ChildNames: []string{"shared.control.shared"},
		ModName:    "right",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "shared.control.shared",
		ShortName: "shared",
		HasSQL:    true,
		ModName:   "shared",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("main", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	// Get dependency order
	names := []string{"main.benchmark.diamond"}
	order, err := resolver.GetDependencyOrder(names)

	require.NoError(t, err)
	require.Len(t, order, 4)

	// Shared should come first (no dependencies)
	assert.Equal(t, "shared.control.shared", order[0], "Shared should be first")
}

// =============================================================================
// Circular Dependency Tests (Cross-Mod)
// =============================================================================

func TestTransitiveDeps_CircularCrossMod(t *testing.T) {
	// Test: Circular dependency across mods (A -> B -> A)
	index := resourceindex.NewResourceIndex()

	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "mod_a.benchmark.a",
		ShortName:  "a",
		ChildNames: []string{"mod_b.benchmark.b"},
		ModName:    "mod_a",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "mod_b.benchmark.b",
		ShortName:  "b",
		ChildNames: []string{"mod_a.benchmark.a"}, // Circular!
		ModName:    "mod_b",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("mod_a", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	// HasCircularDependency should detect the cycle
	assert.True(t, resolver.HasCircularDependency("mod_a.benchmark.a"))

	// GetDependencyOrder should return error
	names := []string{"mod_a.benchmark.a"}
	_, err := resolver.GetDependencyOrder(names)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circular")
}

func TestTransitiveDeps_LongCrossModCircular(t *testing.T) {
	// Test: A -> B -> C -> D -> A (long circular chain across mods)
	index := resourceindex.NewResourceIndex()

	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "mod_a.benchmark.a",
		ShortName:  "a",
		ChildNames: []string{"mod_b.benchmark.b"},
		ModName:    "mod_a",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "mod_b.benchmark.b",
		ShortName:  "b",
		ChildNames: []string{"mod_c.benchmark.c"},
		ModName:    "mod_b",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "mod_c.benchmark.c",
		ShortName:  "c",
		ChildNames: []string{"mod_d.benchmark.d"},
		ModName:    "mod_c",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "mod_d.benchmark.d",
		ShortName:  "d",
		ChildNames: []string{"mod_a.benchmark.a"}, // Circular back to A
		ModName:    "mod_d",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("mod_a", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	// All nodes in the cycle should be detected as circular
	assert.True(t, resolver.HasCircularDependency("mod_a.benchmark.a"))
	assert.True(t, resolver.HasCircularDependency("mod_b.benchmark.b"))
	assert.True(t, resolver.HasCircularDependency("mod_c.benchmark.c"))
	assert.True(t, resolver.HasCircularDependency("mod_d.benchmark.d"))
}

// =============================================================================
// Missing Cross-Mod Dependency Tests
// =============================================================================

func TestTransitiveDeps_MissingCrossModDep(t *testing.T) {
	// Test: Control references query from non-existent mod
	index := resourceindex.NewResourceIndex()

	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "main.control.broken",
		ShortName: "broken",
		QueryRef:  "nonexistent_mod.query.missing",
		ModName:   "main",
	})

	// Note: nonexistent_mod.query.missing is NOT in the index

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("main", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	// GetDependencies should still return the reference
	deps := resolver.GetDependencies("main.control.broken")
	require.Len(t, deps, 1)
	assert.Equal(t, "nonexistent_mod.query.missing", deps[0].To)

	// GetTransitiveDependencies handles missing gracefully
	transDeps := resolver.GetTransitiveDependencies("main.control.broken")
	// Should have at least the control itself
	assert.Contains(t, transDeps, "main.control.broken")
}

func TestTransitiveDeps_PartiallyMissingDeps(t *testing.T) {
	// Test: Some dependencies exist, some don't
	index := resourceindex.NewResourceIndex()

	index.Add(&resourceindex.IndexEntry{
		Type:      "benchmark",
		Name:      "main.benchmark.partial",
		ShortName: "partial",
		ChildNames: []string{
			"exists.control.good",
			"missing.control.bad",
		},
		ModName: "main",
	})

	// Only one child exists
	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "exists.control.good",
		ShortName: "good",
		HasSQL:    true,
		ModName:   "exists",
	})

	// missing.control.bad is NOT added

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("main", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	// GetDependencies shows all references (whether they exist or not)
	deps := resolver.GetDependencies("main.benchmark.partial")
	assert.Len(t, deps, 2)

	// GetTransitiveDependencies handles gracefully
	transDeps := resolver.GetTransitiveDependencies("main.benchmark.partial")
	// Should include existing resources
	assert.Contains(t, transDeps, "exists.control.good")
	assert.Contains(t, transDeps, "main.benchmark.partial")
}

// =============================================================================
// Wide Dependency Tests
// =============================================================================

func TestTransitiveDeps_WideDependencies(t *testing.T) {
	// Test: Benchmark with many children from many mods
	index := resourceindex.NewResourceIndex()

	childNames := []string{}
	for i := 0; i < 10; i++ {
		modName := fmt.Sprintf("mod_%d", i)
		controlName := fmt.Sprintf("mod_%d.control.control_%d", i, i)
		shortName := fmt.Sprintf("control_%d", i)
		childNames = append(childNames, controlName)

		index.Add(&resourceindex.IndexEntry{
			Type:      "control",
			Name:      controlName,
			ShortName: shortName,
			HasSQL:    true,
			ModName:   modName,
		})
	}

	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "main.benchmark.wide",
		ShortName:  "wide",
		ChildNames: childNames,
		ModName:    "main",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("main", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	// Should get all 10 children
	deps := resolver.GetDependencies("main.benchmark.wide")
	assert.Len(t, deps, 10)

	// Transitive should include benchmark + all 10 children
	transDeps := resolver.GetTransitiveDependencies("main.benchmark.wide")
	assert.Len(t, transDeps, 11)
}

// =============================================================================
// Dependency Order Tests
// =============================================================================

func TestTransitiveDeps_DependencyOrderCrossMod(t *testing.T) {
	// Test: Ordering with cross-mod dependencies
	index := resourceindex.NewResourceIndex()

	// Chain: benchmark -> control (different mod) -> query (another mod)
	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "main.benchmark.top",
		ShortName:  "top",
		ChildNames: []string{"dep_a.control.middle"},
		ModName:    "main",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "dep_a.control.middle",
		ShortName: "middle",
		QueryRef:  "dep_b.query.bottom",
		ModName:   "dep_a",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "query",
		Name:      "dep_b.query.bottom",
		ShortName: "bottom",
		HasSQL:    true,
		ModName:   "dep_b",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("main", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	names := []string{"main.benchmark.top"}
	order, err := resolver.GetDependencyOrder(names)

	require.NoError(t, err)
	require.Len(t, order, 3)

	// Verify order: query first, then control, then benchmark
	assert.Equal(t, "dep_b.query.bottom", order[0])
	assert.Equal(t, "dep_a.control.middle", order[1])
	assert.Equal(t, "main.benchmark.top", order[2])
}

// =============================================================================
// Resolve With Dependencies Tests (Integration)
// =============================================================================

func TestTransitiveDeps_ResolveWithDependencies(t *testing.T) {
	// Test: ResolveWithDependencies handles cross-mod refs
	index := resourceindex.NewResourceIndex()

	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "main.control.test",
		ShortName: "test",
		QueryRef:  "dep.query.helper",
		ModName:   "main",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "query",
		Name:      "dep.query.helper",
		ShortName: "helper",
		HasSQL:    true,
		ModName:   "dep",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("main", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	ctx := context.Background()

	// ResolveWithDependencies should not panic on cross-mod deps
	// It may error due to missing files, but shouldn't crash
	err := resolver.ResolveWithDependencies(ctx, "main.control.test")
	// Error expected due to no actual files
	t.Logf("ResolveWithDependencies error (expected): %v", err)
}

func TestTransitiveDeps_CircularResolution(t *testing.T) {
	// Test: ResolveWithDependencies detects cycles
	// Note: Due to implementation details, circular detection happens
	// in GetDependencyOrder/HasCircularDependency before loading
	index := resourceindex.NewResourceIndex()

	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "a.benchmark.a",
		ShortName:  "a",
		ChildNames: []string{"b.benchmark.b"},
		ModName:    "a",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "b.benchmark.b",
		ShortName:  "b",
		ChildNames: []string{"a.benchmark.a"}, // Circular
		ModName:    "b",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("a", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	// HasCircularDependency should detect the cycle
	assert.True(t, resolver.HasCircularDependency("a.benchmark.a"),
		"Should detect circular dependency")

	// GetDependencyOrder should also detect the cycle
	_, err := resolver.GetDependencyOrder([]string{"a.benchmark.a"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "circular")
}

// =============================================================================
// Dependency Graph Tests
// =============================================================================

func TestTransitiveDeps_BuildDependencyGraphCrossMod(t *testing.T) {
	// Test: BuildDependencyGraph with cross-mod refs
	index := resourceindex.NewResourceIndex()

	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "main.benchmark.root",
		ShortName:  "root",
		ChildNames: []string{"dep.control.child"},
		ModName:    "main",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "dep.control.child",
		ShortName: "child",
		QueryRef:  "dep.query.leaf",
		ModName:   "dep",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "query",
		Name:      "dep.query.leaf",
		ShortName: "leaf",
		HasSQL:    true,
		ModName:   "dep",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("main", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	graph := resolver.BuildDependencyGraph()

	// Verify graph structure
	assert.Len(t, graph.Nodes["main.benchmark.root"], 1)
	assert.Len(t, graph.Nodes["dep.control.child"], 1)
	assert.Len(t, graph.Nodes["dep.query.leaf"], 0)

	// Verify cross-mod edge
	rootDeps := graph.Nodes["main.benchmark.root"]
	assert.Equal(t, "dep.control.child", rootDeps[0].To)
}

// =============================================================================
// GetDependents Tests (Reverse Lookup)
// =============================================================================

func TestTransitiveDeps_GetDependentsCrossMod(t *testing.T) {
	// Test: GetDependents finds cross-mod dependents
	index := resourceindex.NewResourceIndex()

	// Shared query used by multiple controls from different mods
	index.Add(&resourceindex.IndexEntry{
		Type:      "query",
		Name:      "shared.query.common",
		ShortName: "common",
		HasSQL:    true,
		ModName:   "shared",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "mod_a.control.uses_common",
		ShortName: "uses_common",
		QueryRef:  "shared.query.common",
		ModName:   "mod_a",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "mod_b.control.also_uses",
		ShortName: "also_uses",
		QueryRef:  "shared.query.common",
		ModName:   "mod_b",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("shared", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	dependents := resolver.GetDependents("shared.query.common")

	assert.Len(t, dependents, 2)
	assert.Contains(t, dependents, "mod_a.control.uses_common")
	assert.Contains(t, dependents, "mod_b.control.also_uses")
}

// =============================================================================
// Benchmark Tests
// =============================================================================

func BenchmarkTransitiveDeps_ThreeLevelChain(b *testing.B) {
	index := resourceindex.NewResourceIndex()

	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "a.benchmark.a",
		ShortName:  "a",
		ChildNames: []string{"b.control.b"},
		ModName:    "a",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "b.control.b",
		ShortName: "b",
		QueryRef:  "c.query.c",
		ModName:   "b",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "query",
		Name:      "c.query.c",
		ShortName: "c",
		HasSQL:    true,
		ModName:   "c",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("a", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = resolver.GetTransitiveDependencies("a.benchmark.a")
	}
}

func BenchmarkTransitiveDeps_DiamondPattern(b *testing.B) {
	index := resourceindex.NewResourceIndex()

	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "main.benchmark.top",
		ShortName:  "top",
		ChildNames: []string{"left.control.left", "right.control.right"},
		ModName:    "main",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "left.control.left",
		ShortName: "left",
		QueryRef:  "shared.query.shared",
		ModName:   "left",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "right.control.right",
		ShortName: "right",
		QueryRef:  "shared.query.shared",
		ModName:   "right",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "query",
		Name:      "shared.query.shared",
		ShortName: "shared",
		HasSQL:    true,
		ModName:   "shared",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("main", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = resolver.GetDependencyOrder([]string{"main.benchmark.top"})
	}
}

func BenchmarkTransitiveDeps_WideChildren(b *testing.B) {
	index := resourceindex.NewResourceIndex()

	childNames := []string{}
	for i := 0; i < 50; i++ {
		name := fmt.Sprintf("mod_%d.control.c_%d", i, i)
		childNames = append(childNames, name)
		index.Add(&resourceindex.IndexEntry{
			Type:      "control",
			Name:      name,
			ShortName: fmt.Sprintf("c_%d", i),
			HasSQL:    true,
			ModName:   fmt.Sprintf("mod_%d", i),
		})
	}

	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "main.benchmark.wide",
		ShortName:  "wide",
		ChildNames: childNames,
		ModName:    "main",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("main", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = resolver.GetTransitiveDependencies("main.benchmark.wide")
	}
}
