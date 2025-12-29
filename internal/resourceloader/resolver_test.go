package resourceloader

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/powerpipe/internal/resourcecache"
	"github.com/turbot/powerpipe/internal/resourceindex"
)

func TestResolver_GetDependencies(t *testing.T) {
	index := resourceindex.NewResourceIndex()
	index.ModName = "testmod"

	// Add a control that references a query
	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "testmod.control.uses_query",
		ShortName: "uses_query",
		QueryRef:  "testmod.query.referenced",
		ModName:   "testmod",
	})

	// Add the referenced query
	index.Add(&resourceindex.IndexEntry{
		Type:      "query",
		Name:      "testmod.query.referenced",
		ShortName: "referenced",
		HasSQL:    true,
		ModName:   "testmod",
	})

	// Add a benchmark with children
	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "testmod.benchmark.parent",
		ShortName:  "parent",
		ChildNames: []string{"testmod.control.uses_query", "testmod.control.inline"},
		ModName:    "testmod",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("testmod", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	// Test control dependencies
	deps := resolver.GetDependencies("testmod.control.uses_query")
	require.Len(t, deps, 1)
	assert.Equal(t, "testmod.query.referenced", deps[0].To)
	assert.Equal(t, DepQuery, deps[0].Type)

	// Test benchmark dependencies
	deps = resolver.GetDependencies("testmod.benchmark.parent")
	require.Len(t, deps, 2)

	// Check child dependencies
	childNames := make(map[string]bool)
	for _, dep := range deps {
		childNames[dep.To] = true
		assert.Equal(t, DepChild, dep.Type)
	}
	assert.True(t, childNames["testmod.control.uses_query"])
	assert.True(t, childNames["testmod.control.inline"])

	// Test query has no dependencies
	deps = resolver.GetDependencies("testmod.query.referenced")
	assert.Len(t, deps, 0)
}

func TestResolver_GetTransitiveDependencies(t *testing.T) {
	index := resourceindex.NewResourceIndex()
	index.ModName = "testmod"

	// Create a dependency chain:
	// benchmark.parent -> control.child -> query.leaf
	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "testmod.benchmark.parent",
		ShortName:  "parent",
		ChildNames: []string{"testmod.control.child"},
		ModName:    "testmod",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "testmod.control.child",
		ShortName: "child",
		QueryRef:  "testmod.query.leaf",
		ModName:   "testmod",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "query",
		Name:      "testmod.query.leaf",
		ShortName: "leaf",
		HasSQL:    true,
		ModName:   "testmod",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("testmod", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	deps := resolver.GetTransitiveDependencies("testmod.benchmark.parent")

	// Should include all three in dependency order (leaf first)
	require.Len(t, deps, 3)

	// Query should come before control (since control depends on query)
	queryIdx := indexOf(deps, "testmod.query.leaf")
	controlIdx := indexOf(deps, "testmod.control.child")
	benchmarkIdx := indexOf(deps, "testmod.benchmark.parent")

	assert.Less(t, queryIdx, controlIdx, "query should come before control")
	assert.Less(t, controlIdx, benchmarkIdx, "control should come before benchmark")
}

func TestResolver_GetDependencyOrder(t *testing.T) {
	index := resourceindex.NewResourceIndex()
	index.ModName = "testmod"

	// Set up resources with dependencies
	index.Add(&resourceindex.IndexEntry{
		Type:      "query",
		Name:      "testmod.query.q1",
		ShortName: "q1",
		HasSQL:    true,
		ModName:   "testmod",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "testmod.control.c1",
		ShortName: "c1",
		QueryRef:  "testmod.query.q1",
		ModName:   "testmod",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "testmod.benchmark.b1",
		ShortName:  "b1",
		ChildNames: []string{"testmod.control.c1"},
		ModName:    "testmod",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("testmod", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	names := []string{
		"testmod.benchmark.b1",
		"testmod.control.c1",
		"testmod.query.q1",
	}

	order, err := resolver.GetDependencyOrder(names)
	require.NoError(t, err)
	require.Len(t, order, 3)

	// Verify order: query -> control -> benchmark
	queryIdx := indexOf(order, "testmod.query.q1")
	controlIdx := indexOf(order, "testmod.control.c1")
	benchmarkIdx := indexOf(order, "testmod.benchmark.b1")

	assert.Less(t, queryIdx, controlIdx, "query should come before control")
	assert.Less(t, controlIdx, benchmarkIdx, "control should come before benchmark")
}

func TestResolver_CircularDependency(t *testing.T) {
	index := resourceindex.NewResourceIndex()
	index.ModName = "testmod"

	// Create circular dependency: a -> b -> c -> a
	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "testmod.benchmark.a",
		ShortName:  "a",
		ChildNames: []string{"testmod.benchmark.b"},
		ModName:    "testmod",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "testmod.benchmark.b",
		ShortName:  "b",
		ChildNames: []string{"testmod.benchmark.c"},
		ModName:    "testmod",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "testmod.benchmark.c",
		ShortName:  "c",
		ChildNames: []string{"testmod.benchmark.a"}, // Circular!
		ModName:    "testmod",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("testmod", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	// HasCircularDependency should detect the cycle
	assert.True(t, resolver.HasCircularDependency("testmod.benchmark.a"))
	assert.True(t, resolver.HasCircularDependency("testmod.benchmark.b"))
	assert.True(t, resolver.HasCircularDependency("testmod.benchmark.c"))

	// GetDependencyOrder should return an error
	names := []string{"testmod.benchmark.a", "testmod.benchmark.b", "testmod.benchmark.c"}
	_, err := resolver.GetDependencyOrder(names)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circular")
}

func TestResolver_NoCircularDependency(t *testing.T) {
	index := resourceindex.NewResourceIndex()
	index.ModName = "testmod"

	// Create a simple tree (no cycles)
	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "testmod.benchmark.root",
		ShortName:  "root",
		ChildNames: []string{"testmod.benchmark.child1", "testmod.benchmark.child2"},
		ModName:    "testmod",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "benchmark",
		Name:      "testmod.benchmark.child1",
		ShortName: "child1",
		ModName:   "testmod",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "benchmark",
		Name:      "testmod.benchmark.child2",
		ShortName: "child2",
		ModName:   "testmod",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("testmod", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	assert.False(t, resolver.HasCircularDependency("testmod.benchmark.root"))
}

func TestResolver_GetDependents(t *testing.T) {
	index := resourceindex.NewResourceIndex()
	index.ModName = "testmod"

	// Query used by multiple controls
	index.Add(&resourceindex.IndexEntry{
		Type:      "query",
		Name:      "testmod.query.shared",
		ShortName: "shared",
		HasSQL:    true,
		ModName:   "testmod",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "testmod.control.c1",
		ShortName: "c1",
		QueryRef:  "testmod.query.shared",
		ModName:   "testmod",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "testmod.control.c2",
		ShortName: "c2",
		QueryRef:  "testmod.query.shared",
		ModName:   "testmod",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("testmod", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	dependents := resolver.GetDependents("testmod.query.shared")
	assert.Len(t, dependents, 2)
	assert.Contains(t, dependents, "testmod.control.c1")
	assert.Contains(t, dependents, "testmod.control.c2")
}

func TestResolver_ResolveWithDependencies(t *testing.T) {
	modPath, mod := setupTestModWithDeps(t)
	loader := setupTestLoaderWithDeps(t, modPath, mod)
	resolver := NewDependencyResolver(loader.index, loader)

	ctx := context.Background()

	// Resolve a simple query first (no dependencies)
	err := resolver.ResolveWithDependencies(ctx, "testmod.query.q1")
	require.NoError(t, err)

	// Verify query is cached
	_, ok := loader.cache.GetResource("testmod.query.q1")
	assert.True(t, ok, "query should be cached")
}

func TestResolver_ResolveWithDependencies_Control(t *testing.T) {
	modPath, mod := setupTestModWithDeps(t)
	loader := setupTestLoaderWithDeps(t, modPath, mod)
	resolver := NewDependencyResolver(loader.index, loader)

	ctx := context.Background()

	// Resolve control - should load query dependency first
	// Note: The actual query loading may fail since it uses a reference,
	// but the resolver should still proceed and cache the control
	err := resolver.ResolveWithDependencies(ctx, "testmod.control.inline")
	require.NoError(t, err)

	// Verify control is cached
	_, ok := loader.cache.GetResource("testmod.control.inline")
	assert.True(t, ok, "control should be cached")
}

func TestResolver_BuildDependencyGraph(t *testing.T) {
	index := resourceindex.NewResourceIndex()
	index.ModName = "testmod"

	index.Add(&resourceindex.IndexEntry{
		Type:       "benchmark",
		Name:       "testmod.benchmark.b1",
		ShortName:  "b1",
		ChildNames: []string{"testmod.control.c1"},
		ModName:    "testmod",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "testmod.control.c1",
		ShortName: "c1",
		QueryRef:  "testmod.query.q1",
		ModName:   "testmod",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "query",
		Name:      "testmod.query.q1",
		ShortName: "q1",
		HasSQL:    true,
		ModName:   "testmod",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	mod := modconfig.NewMod("testmod", "/tmp", hcl.Range{})
	loader := NewLoader(index, cache, mod, "/tmp")
	resolver := NewDependencyResolver(index, loader)

	graph := resolver.BuildDependencyGraph()

	// Benchmark has 1 dependency (control)
	assert.Len(t, graph.Nodes["testmod.benchmark.b1"], 1)

	// Control has 1 dependency (query)
	assert.Len(t, graph.Nodes["testmod.control.c1"], 1)

	// Query has no dependencies
	assert.Len(t, graph.Nodes["testmod.query.q1"], 0)
}

// Helper functions

func indexOf(slice []string, item string) int {
	for i, s := range slice {
		if s == item {
			return i
		}
	}
	return -1
}

func setupTestModWithDeps(t testing.TB) (string, *modconfig.Mod) {
	tmpDir := t.TempDir()

	// Create mod.pp
	modContent := `mod "testmod" {
  title = "Test Mod"
}`
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "mod.pp"), []byte(modContent), 0644))

	// Create queries.pp
	queriesContent := `query "q1" {
  sql = "SELECT 1"
}

query "q2" {
  sql = "SELECT 2"
}`
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "queries.pp"), []byte(queriesContent), 0644))

	// Create controls.pp - one inline control (direct SQL), one referencing query
	controlsContent := `control "inline" {
  sql   = "SELECT 'ok' as status, 'test' as reason"
  title = "Inline Control"
}

control "child" {
  query = query.q1
  title = "Child Control"
}`
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "controls.pp"), []byte(controlsContent), 0644))

	mod := modconfig.NewMod("testmod", tmpDir, hcl.Range{})
	return tmpDir, mod
}

func setupTestLoaderWithDeps(t testing.TB, modPath string, mod *modconfig.Mod) *Loader {
	index := resourceindex.NewResourceIndex()
	index.ModName = "testmod"

	queriesFile := filepath.Join(modPath, "queries.pp")
	controlsFile := filepath.Join(modPath, "controls.pp")

	// Add query entries
	index.Add(&resourceindex.IndexEntry{
		Type:      "query",
		Name:      "testmod.query.q1",
		ShortName: "q1",
		FileName:  queriesFile,
		StartLine: 1,
		EndLine:   3,
		HasSQL:    true,
		ModName:   "testmod",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "query",
		Name:      "testmod.query.q2",
		ShortName: "q2",
		FileName:  queriesFile,
		StartLine: 5,
		EndLine:   7,
		HasSQL:    true,
		ModName:   "testmod",
	})

	// Add inline control entry (has direct SQL, no query reference)
	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "testmod.control.inline",
		ShortName: "inline",
		FileName:  controlsFile,
		StartLine: 1,
		EndLine:   4,
		HasSQL:    true,
		ModName:   "testmod",
	})

	// Add control entry with query reference
	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "testmod.control.child",
		ShortName: "child",
		FileName:  controlsFile,
		StartLine: 6,
		EndLine:   9,
		QueryRef:  "testmod.query.q1",
		ModName:   "testmod",
	})

	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())
	return NewLoader(index, cache, mod, modPath)
}
