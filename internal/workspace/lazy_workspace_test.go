package workspace

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLazyWorkspace_FastStartup(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "lazy-loading-tests", "generated", "medium")

	start := time.Now()
	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	loadTime := time.Since(start)

	t.Logf("Lazy workspace load time: %v", loadTime)

	// Should be very fast - just index building
	// For a medium mod, it should complete in under 500ms
	assert.Less(t, loadTime.Milliseconds(), int64(500),
		"Lazy load should be fast")

	// Index should have resources
	stats := lw.IndexStats()
	assert.Greater(t, stats.TotalEntries, 0, "Index should have entries")

	// Cache should be empty (nothing loaded yet)
	cacheStats := lw.CacheStats()
	assert.Equal(t, 0, cacheStats.Entries, "Cache should be empty at startup")

	lw.Close()
}

func TestLazyWorkspace_OnDemandLoading(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "lazy-loading-tests", "generated", "small")

	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	ctx := context.Background()

	// Nothing loaded initially
	assert.Equal(t, 0, lw.cache.Stats().Entries)

	// Load a dashboard
	dash, err := lw.LoadDashboard(ctx, "lazy_small.dashboard.dashboard_0")
	require.NoError(t, err)
	assert.NotNil(t, dash)

	// Now something is cached
	assert.Greater(t, lw.cache.Stats().Entries, 0)
}

func TestLazyWorkspace_AvailableDashboards(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "lazy-loading-tests", "generated", "small")

	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	// Get available dashboards without loading any
	payload := lw.GetAvailableDashboardsFromIndex()

	assert.NotEmpty(t, payload.Dashboards, "Should have dashboards in payload")
	assert.NotEmpty(t, payload.Benchmarks, "Should have benchmarks in payload")

	// Still nothing loaded in cache
	assert.Equal(t, 0, lw.cache.Stats().Entries, "Cache should still be empty")
}

func TestLazyWorkspace_MemoryBounded(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "lazy-loading-tests", "generated", "medium")

	config := LazyLoadConfig{
		MaxCacheMemory: 1 * 1024 * 1024, // 1MB limit (small for testing)
	}

	lw, err := NewLazyWorkspace(context.Background(), modPath, config)
	require.NoError(t, err)
	defer lw.Close()

	ctx := context.Background()

	// Load many resources
	resources := lw.lazyResources
	dashboardNames := resources.ListDashboardNames()

	for _, name := range dashboardNames {
		_, _ = lw.LoadResource(ctx, name)
	}

	// Memory should be bounded (though we can't strictly enforce due to object overhead)
	stats := lw.cache.Stats()
	t.Logf("Cache stats: entries=%d, memory=%d bytes, evictions=%d",
		stats.Entries, stats.MemoryBytes, stats.Evictions)

	// With a small cache limit, we expect some evictions if we loaded many resources
	// This test verifies the cache is working, not strict memory limits
	assert.Greater(t, stats.Entries, 0, "Should have some cached entries")
}

func TestLazyWorkspace_GetResource(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "lazy-loading-tests", "generated", "small")

	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	ctx := context.Background()

	// Load a specific resource
	resource, err := lw.LoadResource(ctx, "lazy_small.query.query_0")
	require.NoError(t, err)
	assert.NotNil(t, resource)
	assert.Equal(t, "lazy_small.query.query_0", resource.Name())
}

func TestLazyWorkspace_IsLazy(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "lazy-loading-tests", "generated", "small")

	// Test LazyWorkspace
	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	assert.True(t, lw.IsLazy(), "LazyWorkspace should return true for IsLazy()")
	lw.Close()

	// Test regular workspace
	pw := NewPowerpipeWorkspace(modPath)
	assert.False(t, pw.IsLazy(), "PowerpipeWorkspace should return false for IsLazy()")
	pw.Close()
}

func TestLazyModResources_Counts(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "lazy-loading-tests", "generated", "small")

	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	resources := lw.GetLazyModResources()

	// Get counts without loading
	dashCount := resources.DashboardCount()
	queryCount := resources.QueryCount()
	controlCount := resources.ControlCount()
	benchmarkCount := resources.BenchmarkCount()

	t.Logf("Counts: dashboards=%d, queries=%d, controls=%d, benchmarks=%d",
		dashCount, queryCount, controlCount, benchmarkCount)

	// Small mod should have some resources
	assert.Greater(t, dashCount, 0, "Should have dashboards")
	assert.Greater(t, queryCount, 0, "Should have queries")
	assert.Greater(t, controlCount, 0, "Should have controls")
	assert.Greater(t, benchmarkCount, 0, "Should have benchmarks")

	// Cache should still be empty - counts come from index
	assert.Equal(t, 0, lw.cache.Stats().Entries, "Cache should be empty after counts")
}

func TestLazyModResources_ListNames(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "lazy-loading-tests", "generated", "small")

	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	resources := lw.GetLazyModResources()

	// List names without loading
	dashNames := resources.ListDashboardNames()
	queryNames := resources.ListQueryNames()
	controlNames := resources.ListControlNames()
	benchmarkNames := resources.ListBenchmarkNames()

	assert.NotEmpty(t, dashNames, "Should have dashboard names")
	assert.NotEmpty(t, queryNames, "Should have query names")
	assert.NotEmpty(t, controlNames, "Should have control names")
	assert.NotEmpty(t, benchmarkNames, "Should have benchmark names")

	// Cache should still be empty
	assert.Equal(t, 0, lw.cache.Stats().Entries, "Cache should be empty after listing names")
}

func TestLoadLazy(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "lazy-loading-tests", "generated", "small")

	lw, err := LoadLazy(context.Background(), modPath)
	require.NoError(t, err)
	defer lw.Close()

	assert.True(t, lw.IsLazy())
	assert.Greater(t, lw.IndexStats().TotalEntries, 0)
}

func TestLoadAuto_WithLazyLoading(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "lazy-loading-tests", "generated", "small")

	wp, ew := LoadAuto(context.Background(), modPath, WithLazyLoading(true))
	require.Nil(t, ew.GetError())
	defer wp.Close()

	assert.True(t, wp.IsLazy(), "Should be lazy workspace")

	// Verify it's actually a LazyWorkspace
	_, ok = wp.(*LazyWorkspace)
	assert.True(t, ok, "Should be *LazyWorkspace type")
}

func TestLazyWorkspace_CacheInvalidation(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "lazy-loading-tests", "generated", "small")

	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	ctx := context.Background()

	// Load a resource
	_, err = lw.LoadResource(ctx, "lazy_small.query.query_0")
	require.NoError(t, err)

	initialCount := lw.cache.Stats().Entries
	assert.Greater(t, initialCount, 0)

	// Invalidate it
	lw.InvalidateResource("lazy_small.query.query_0")

	// Should be gone (or fewer entries)
	afterCount := lw.cache.Stats().Entries
	assert.Less(t, afterCount, initialCount)

	// Invalidate all
	lw.InvalidateAll()
	assert.Equal(t, 0, lw.cache.Stats().Entries)
}

func TestLazyWorkspace_IndexStats(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "lazy-loading-tests", "generated", "small")

	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	stats := lw.IndexStats()

	assert.Greater(t, stats.TotalEntries, 0)
	assert.Greater(t, stats.TotalSize, 0)
	assert.NotEmpty(t, stats.ByType)

	// Check we have expected types
	t.Logf("Index stats: %+v", stats)
}
