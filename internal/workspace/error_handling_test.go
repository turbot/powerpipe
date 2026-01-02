package workspace

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turbot/pipe-fittings/v2/modconfig"
)

// =============================================================================
// Lazy Workspace Error Tests
// =============================================================================

func TestError_IndexBuildFailure_NonexistentDirectory(t *testing.T) {
	// Test that workspace creation fails gracefully for non-existent directory
	_, err := NewLazyWorkspace(context.Background(), "/nonexistent/path/that/does/not/exist", DefaultLazyLoadConfig())

	assert.Error(t, err, "Should error on non-existent directory")
	// Should not panic, just return error
}

func TestError_IndexBuildFailure_FileInsteadOfDirectory(t *testing.T) {
	// Create a temp file (not directory)
	tmpFile, err := os.CreateTemp("", "not-a-dir-*.txt")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Try to create workspace with file path
	_, err = NewLazyWorkspace(context.Background(), tmpFile.Name(), DefaultLazyLoadConfig())

	assert.Error(t, err, "Should error when path is a file, not directory")
}

func TestError_EagerLoadFailureCached(t *testing.T) {
	// Test that errors from GetWorkspaceForExecution are cached via sync.Once
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)

	// Use a mod with syntax errors - this should fail during eager load
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "error-conditions", "invalid-syntax")

	// Lazy workspace creation should succeed (just builds index)
	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err, "Lazy workspace creation should succeed even with broken files")
	defer lw.Close()

	ctx := context.Background()

	// First call to GetWorkspaceForExecution should fail
	_, err1 := lw.GetWorkspaceForExecution(ctx)
	assert.Error(t, err1, "First eager load should fail due to syntax errors")

	// Second call should return same error (cached by sync.Once)
	_, err2 := lw.GetWorkspaceForExecution(ctx)
	assert.Error(t, err2, "Second call should also fail")

	// Errors should be the same (cached)
	if err1 != nil && err2 != nil {
		assert.Equal(t, err1.Error(), err2.Error(), "Error should be cached")
	}
}

func TestError_ResourceNotFoundInIndex(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "generated", "small")

	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	ctx := context.Background()

	// Try to load non-existent resource
	_, err = lw.LoadResource(ctx, "small_test.query.nonexistent_resource")
	assert.Error(t, err, "Should error for non-existent resource")
	assert.Contains(t, err.Error(), "not found", "Error should indicate resource not found")
}

func TestError_DashboardNotFound(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "generated", "small")

	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	ctx := context.Background()

	// Try to load non-existent dashboard
	_, err = lw.LoadDashboard(ctx, "small_test.dashboard.nonexistent_dashboard")
	assert.Error(t, err, "Should error for non-existent dashboard")
}

func TestError_BenchmarkNotFound(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "generated", "small")

	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	ctx := context.Background()

	// Try to load non-existent benchmark
	_, err = lw.LoadBenchmark(ctx, "small_test.benchmark.nonexistent_benchmark")
	assert.Error(t, err, "Should error for non-existent benchmark")
}

func TestError_GetResource_ParsedNameNil(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "generated", "small")

	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	// GetResource with a resource name that doesn't exist
	parsedName := &modconfig.ParsedResourceName{
		ItemType: "query",
		Name:     "nonexistent",
	}
	resource, found := lw.GetResource(parsedName)
	assert.False(t, found, "Should not find non-existent resource")
	assert.Nil(t, resource, "Resource should be nil")
}

// =============================================================================
// Panic Prevention Tests
// =============================================================================

func TestError_NoPanicOnEmptyWorkspace(t *testing.T) {
	// Create a temp directory with empty mod
	tmpDir, err := os.MkdirTemp("", "empty-mod-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create minimal mod.pp
	modContent := `mod "empty_mod" { title = "Empty" }`
	err = os.WriteFile(filepath.Join(tmpDir, "mod.pp"), []byte(modContent), 0600)
	require.NoError(t, err)

	// Should not panic
	lw, err := NewLazyWorkspace(context.Background(), tmpDir, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	// Operations on empty workspace should not panic
	ctx := context.Background()

	// Load non-existent should error, not panic
	_, err = lw.LoadResource(ctx, "empty_mod.query.test")
	assert.Error(t, err)

	// Stats should work
	stats := lw.IndexStats()
	assert.Equal(t, 0, stats.TotalEntries, "Empty mod should have no entries")

	cacheStats := lw.CacheStats()
	assert.Equal(t, 0, cacheStats.Entries, "Cache should be empty")

	// GetAvailableDashboardsFromIndex should not panic
	payload := lw.GetAvailableDashboardsFromIndex()
	assert.NotNil(t, payload)
	assert.Empty(t, payload.Dashboards)
	assert.Empty(t, payload.Benchmarks)
}

func TestError_NoPanicOnNilResource(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "generated", "small")

	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	// GetResource should return false, not panic
	parsedName := &modconfig.ParsedResourceName{
		ItemType: "dashboard",
		Name:     "nonexistent",
	}
	resource, found := lw.GetResource(parsedName)
	assert.False(t, found)
	assert.Nil(t, resource)
}

func TestError_NoPanicOnConcurrentErrors(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "generated", "small")

	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	ctx := context.Background()
	var wg sync.WaitGroup

	// Concurrently try to load non-existent resources
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// These should all error, but not panic
			_, _ = lw.LoadResource(ctx, "small_test.query.nonexistent")
			_, _ = lw.LoadDashboard(ctx, "small_test.dashboard.nonexistent")
			_, _ = lw.LoadBenchmark(ctx, "small_test.benchmark.nonexistent")
		}(i)
	}

	wg.Wait() // Should complete without panic
}

// =============================================================================
// Cache Error Handling Tests
// =============================================================================

func TestError_CacheLookupAfterInvalidation(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "generated", "small")

	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	ctx := context.Background()

	// Load a resource
	resourceName := "small_test.query.query_0"
	resource1, err := lw.LoadResource(ctx, resourceName)
	require.NoError(t, err)
	assert.NotNil(t, resource1)

	// Invalidate it
	lw.InvalidateResource(resourceName)

	// Loading again should work (reloads from disk)
	resource2, err := lw.LoadResource(ctx, resourceName)
	require.NoError(t, err)
	assert.NotNil(t, resource2)
}

func TestError_CacheInvalidationOnNonexistent(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "generated", "small")

	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	// Invalidating non-existent resource should not panic
	lw.InvalidateResource("nonexistent.resource.name")

	// InvalidateAll should also work
	lw.InvalidateAll()
}

// =============================================================================
// Partial Load Tests
// =============================================================================

func TestError_PartialValidResources(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "error-conditions", "partial-valid")

	// Create lazy workspace - this should succeed with partial-valid mod
	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	ctx := context.Background()

	// Valid resources should load
	valid, err := lw.LoadResource(ctx, "error_partial_valid.query.valid_query_1")
	require.NoError(t, err)
	assert.NotNil(t, valid)

	valid2, err := lw.LoadResource(ctx, "error_partial_valid.query.valid_query_2")
	require.NoError(t, err)
	assert.NotNil(t, valid2)

	// The control with bad reference should exist in index but may fail during execution
	// The scanner just indexes metadata - it doesn't validate references
	stats := lw.IndexStats()
	assert.Greater(t, stats.TotalEntries, 0, "Index should have entries")
}

// =============================================================================
// Circular Dependency Tests
// =============================================================================

func TestError_CircularDependencyInIndex(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "error-conditions", "circular-deps")

	// Lazy workspace creation should succeed (just builds index)
	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	// Index should have the resources
	stats := lw.IndexStats()
	assert.Greater(t, stats.TotalEntries, 0, "Index should have entries")

	// The resolver should exist
	resolver := lw.GetResolver()
	assert.NotNil(t, resolver)

	// Note: The scanner doesn't fully parse children arrays from HCL,
	// so circular dependencies defined in files may not be detected.
	// This test verifies the workspace loads without panic.
	// For proper circular dependency testing, see resourceloader/error_test.go
	// which uses manually constructed indexes with proper ChildNames.
}

// =============================================================================
// Error Message Quality Tests
// =============================================================================

func TestError_IncludesResourceName(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "generated", "small")

	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	ctx := context.Background()

	// Try to load non-existent resource
	resourceName := "small_test.query.very_specific_name_for_test"
	_, err = lw.LoadResource(ctx, resourceName)

	assert.Error(t, err)
	// Error should include the resource name
	assert.Contains(t, err.Error(), "very_specific_name_for_test",
		"Error message should include the resource name")
}

func TestError_NoInternalLeakage(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "generated", "small")

	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	ctx := context.Background()

	// Try to load non-existent resource
	_, err = lw.LoadResource(ctx, "small_test.query.nonexistent")

	if err != nil {
		errMsg := err.Error()
		// Should not contain internal paths or stack traces
		assert.NotContains(t, errMsg, "goroutine", "Error should not contain stack traces")
		assert.NotContains(t, errMsg, "panic", "Error should not mention panic")
	}
}

// =============================================================================
// Context Cancellation Tests
// =============================================================================

func TestError_ContextCancellation(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "generated", "small")

	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Operations with cancelled context should fail gracefully
	// Note: Some operations may not check context, which is acceptable
	// The important thing is they don't panic

	// Try loading - may or may not fail depending on implementation
	resource, err := lw.LoadResource(ctx, "small_test.query.query_0")
	// Either succeeds (if context not checked) or fails with context error
	if err != nil {
		assert.True(t,
			strings.Contains(err.Error(), "context") || err == context.Canceled,
			"Error should be context-related if failed: %v", err)
	} else {
		assert.NotNil(t, resource)
	}
}

// =============================================================================
// Recovery Tests
// =============================================================================

func TestError_RecoveryAfterCacheClear(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "generated", "small")

	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	ctx := context.Background()

	// Load some resources
	resource1, err := lw.LoadResource(ctx, "small_test.query.query_0")
	require.NoError(t, err)
	assert.NotNil(t, resource1)

	// Clear the cache
	lw.InvalidateAll()
	assert.Equal(t, 0, lw.CacheStats().Entries)

	// Should be able to reload
	resource2, err := lw.LoadResource(ctx, "small_test.query.query_0")
	require.NoError(t, err)
	assert.NotNil(t, resource2)
}

func TestError_SessionRecoveryAfterError(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "generated", "small")

	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	ctx := context.Background()

	// First, cause an error by trying to load non-existent
	_, err = lw.LoadResource(ctx, "small_test.query.nonexistent")
	assert.Error(t, err)

	// Should be able to load valid resources afterwards
	resource, err := lw.LoadResource(ctx, "small_test.query.query_0")
	require.NoError(t, err)
	assert.NotNil(t, resource)

	// And load dashboard
	dash, err := lw.LoadDashboard(ctx, "small_test.dashboard.dashboard_0")
	require.NoError(t, err)
	assert.NotNil(t, dash)
}

// =============================================================================
// Multiple Errors Tests
// =============================================================================

func TestError_MultipleConsecutiveErrors(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modPath := filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "generated", "small")

	lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	ctx := context.Background()

	// Generate many errors
	for i := 0; i < 100; i++ {
		_, _ = lw.LoadResource(ctx, "small_test.query.nonexistent")
	}

	// Workspace should still be functional
	stats := lw.IndexStats()
	assert.Greater(t, stats.TotalEntries, 0, "Index should still have entries")

	// Should still be able to load valid resources
	resource, err := lw.LoadResource(ctx, "small_test.query.query_0")
	require.NoError(t, err)
	assert.NotNil(t, resource)
}
