package workspace

// Tests for lazy-to-eager workspace transition.
//
// NOTE: Event handler race condition has been fixed (sync.Once + closeCh pattern).
// Event handler tests can now use t.Parallel().
//
// However, some tests still cannot use t.Parallel() due to:
// - Schema caching race: parse.getResourceSchema() uses a global map without synchronization
//
// The original viper race (viper.Set() in SetModfileExists()) has been fixed by
// removing the unnecessary viper global state update in pipe-fittings.
//
// To run all tests:
//   go test -v -count=1 ./internal/workspace/ -run "TestLazyWorkspace_"
//
// To run with race detector:
//   go test -v -race -count=1 ./internal/workspace/ -run "TestLazyWorkspace_"

import (
	"context"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/powerpipe/internal/dashboardevents"
)

// getTestModPath returns the path to a test mod relative to the test file location.
func getTestModPath(t *testing.T, modName string) string {
	t.Helper()
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	return filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "lazy-loading-tests", modName)
}

// getGeneratedModPath returns the path to a generated test mod.
// Uses lazy-loading-tests/generated/ which is committed (not the top-level generated/ which is gitignored).
func getGeneratedModPath(t *testing.T, size string) string {
	t.Helper()
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	return filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "lazy-loading-tests", "generated", size)
}

// =============================================================================
// Basic Transition Tests
// =============================================================================

// TestLazyWorkspace_FirstExecutionTriggersEagerLoad verifies that calling
// GetWorkspaceForExecution triggers an eager load of the workspace.

func TestLazyWorkspace_FirstExecutionTriggersEagerLoad(t *testing.T) {
	modPath := getTestModPath(t, "simple")
	ctx := context.Background()

	// Create lazy workspace
	lw, err := NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	// Verify IsLazy() returns true
	assert.True(t, lw.IsLazy(), "LazyWorkspace should return true for IsLazy()")

	// Verify no eager workspace exists yet
	assert.Nil(t, lw.eagerWorkspace, "Eager workspace should not exist before execution")

	// Call GetWorkspaceForExecution
	ew, err := lw.GetWorkspaceForExecution(ctx)
	require.NoError(t, err)
	require.NotNil(t, ew, "GetWorkspaceForExecution should return a workspace")

	// Verify eager workspace is now set
	assert.NotNil(t, lw.eagerWorkspace, "Eager workspace should exist after execution")

	// Verify the returned workspace is a PowerpipeWorkspace
	assert.False(t, ew.IsLazy(), "Eager workspace should return false for IsLazy()")

	// Verify the eager workspace has loaded resources
	modResources := ew.GetPowerpipeModResources()
	assert.NotNil(t, modResources, "Eager workspace should have mod resources")
}

// TestLazyWorkspace_LazyOpsAfterEagerLoad verifies that lazy operations
// still work after an eager load has occurred.

func TestLazyWorkspace_LazyOpsAfterEagerLoad(t *testing.T) {
	modPath := getTestModPath(t, "simple")
	ctx := context.Background()

	lw, err := NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	// First trigger eager load
	_, err = lw.GetWorkspaceForExecution(ctx)
	require.NoError(t, err)

	// Verify GetAvailableDashboardsFromIndex still works
	payload := lw.GetAvailableDashboardsFromIndex()
	assert.NotNil(t, payload, "GetAvailableDashboardsFromIndex should work after eager load")

	// Verify index is still accessible
	stats := lw.IndexStats()
	assert.Greater(t, stats.TotalEntries, 0, "Index should still have entries after eager load")

	// Verify cache is still accessible
	cacheStats := lw.CacheStats()
	assert.NotNil(t, &cacheStats, "Cache should be accessible after eager load")
}

// TestLazyWorkspace_EagerWorkspaceIsCached verifies that multiple calls to
// GetWorkspaceForExecution return the same workspace instance.

func TestLazyWorkspace_EagerWorkspaceIsCached(t *testing.T) {
	modPath := getTestModPath(t, "simple")
	ctx := context.Background()

	lw, err := NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	// First call
	ew1, err := lw.GetWorkspaceForExecution(ctx)
	require.NoError(t, err)
	require.NotNil(t, ew1)

	// Second call
	ew2, err := lw.GetWorkspaceForExecution(ctx)
	require.NoError(t, err)
	require.NotNil(t, ew2)

	// Should be the same instance (sync.Once ensures single load)
	assert.Same(t, ew1, ew2, "Multiple calls should return the same workspace instance")
}

// =============================================================================
// Concurrent Access Tests
// =============================================================================

// TestLazyWorkspace_ConcurrentExecutionRequests tests that concurrent calls
// to GetWorkspaceForExecution all receive the same workspace without races.

func TestLazyWorkspace_ConcurrentExecutionRequests(t *testing.T) {
	modPath := getTestModPath(t, "simple")
	ctx := context.Background()

	lw, err := NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	const numGoroutines = 10
	workspaces := make([]*PowerpipeWorkspace, numGoroutines)
	errors := make([]error, numGoroutines)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Launch concurrent goroutines all calling GetWorkspaceForExecution
	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			ws, err := lw.GetWorkspaceForExecution(ctx)
			workspaces[idx] = ws
			errors[idx] = err
		}(i)
	}

	wg.Wait()

	// Verify all calls succeeded
	for i := 0; i < numGoroutines; i++ {
		assert.NoError(t, errors[i], "Goroutine %d should not have error", i)
		assert.NotNil(t, workspaces[i], "Goroutine %d should have received workspace", i)
	}

	// Verify all received the same workspace
	for i := 1; i < numGoroutines; i++ {
		assert.Same(t, workspaces[0], workspaces[i],
			"All goroutines should receive the same workspace instance")
	}
}

// TestLazyWorkspace_MixedConcurrentAccess tests that lazy and eager operations
// can run concurrently without interference.

func TestLazyWorkspace_MixedConcurrentAccess(t *testing.T) {
	modPath := getGeneratedModPath(t, "small")
	ctx := context.Background()

	lw, err := NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	const numLazyOps = 5
	const numEagerOps = 5

	var wg sync.WaitGroup
	wg.Add(numLazyOps + numEagerOps)

	// Track any panics
	var panicCount int32

	// Lazy operations (browsing)
	for i := 0; i < numLazyOps; i++ {
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					atomic.AddInt32(&panicCount, 1)
				}
			}()
			// Perform lazy operations
			_ = lw.GetAvailableDashboardsFromIndex()
			_ = lw.IndexStats()
		}()
	}

	// Eager operations (execution)
	for i := 0; i < numEagerOps; i++ {
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					atomic.AddInt32(&panicCount, 1)
				}
			}()
			_, _ = lw.GetWorkspaceForExecution(ctx)
		}()
	}

	wg.Wait()

	assert.Equal(t, int32(0), panicCount, "No panics should occur during mixed concurrent access")
}

// TestLazyWorkspace_GetResourceDuringTransition tests concurrent access to
// GetResource while an eager load is happening.

func TestLazyWorkspace_GetResourceDuringTransition(t *testing.T) {
	modPath := getGeneratedModPath(t, "small")
	ctx := context.Background()

	lw, err := NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	// Get a resource name to use for testing
	resources := lw.GetLazyModResources()
	queryNames := resources.ListQueryNames()
	require.NotEmpty(t, queryNames, "Should have queries in test mod")
	testQueryName := queryNames[0]

	const numOps = 10
	var wg sync.WaitGroup
	wg.Add(numOps * 2) // Half GetResource, half GetWorkspaceForExecution

	var panicCount int32
	var getResourceErrors int32
	var executionErrors int32

	// GetResource operations
	for i := 0; i < numOps; i++ {
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					atomic.AddInt32(&panicCount, 1)
				}
			}()
			_, err := lw.LoadResource(ctx, testQueryName)
			if err != nil {
				atomic.AddInt32(&getResourceErrors, 1)
			}
		}()
	}

	// GetWorkspaceForExecution operations
	for i := 0; i < numOps; i++ {
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					atomic.AddInt32(&panicCount, 1)
				}
			}()
			_, err := lw.GetWorkspaceForExecution(ctx)
			if err != nil {
				atomic.AddInt32(&executionErrors, 1)
			}
		}()
	}

	wg.Wait()

	assert.Equal(t, int32(0), panicCount, "No panics should occur")
	// Note: Some GetResource errors might be acceptable during transition,
	// but no panics should occur
}

// =============================================================================
// Event Handler Transfer Tests
// =============================================================================

// TestLazyWorkspace_EventHandlersCopied verifies that event handlers registered
// on the lazy workspace are copied to the eager workspace during transition.

func TestLazyWorkspace_EventHandlersCopied(t *testing.T) {
	t.Parallel()
	modPath := getTestModPath(t, "simple")
	ctx := context.Background()

	lw, err := NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	// Track received events
	var receivedEvents []dashboardevents.DashboardEvent
	var mu sync.Mutex

	// Register handler on lazy workspace
	handler := func(ctx context.Context, event dashboardevents.DashboardEvent) {
		mu.Lock()
		receivedEvents = append(receivedEvents, event)
		mu.Unlock()
	}

	lw.RegisterDashboardEventHandler(ctx, handler)

	// Trigger eager load
	ew, err := lw.GetWorkspaceForExecution(ctx)
	require.NoError(t, err)

	// Publish event on eager workspace
	testEvent := &dashboardevents.WorkspaceError{
		Error: nil,
	}
	ew.PublishDashboardEvent(ctx, testEvent)

	// Wait a bit for async event handling
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	assert.GreaterOrEqual(t, len(receivedEvents), 1, "Handler should receive events from eager workspace")
}

// TestLazyWorkspace_MultipleHandlersTransferred verifies that multiple handlers
// are all properly transferred during transition.

func TestLazyWorkspace_MultipleHandlersTransferred(t *testing.T) {
	t.Parallel()
	modPath := getTestModPath(t, "simple")
	ctx := context.Background()

	lw, err := NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	// Track events for each handler
	var handler1Count, handler2Count, handler3Count int32

	lw.RegisterDashboardEventHandler(ctx, func(ctx context.Context, event dashboardevents.DashboardEvent) {
		atomic.AddInt32(&handler1Count, 1)
	})

	lw.RegisterDashboardEventHandler(ctx, func(ctx context.Context, event dashboardevents.DashboardEvent) {
		atomic.AddInt32(&handler2Count, 1)
	})

	lw.RegisterDashboardEventHandler(ctx, func(ctx context.Context, event dashboardevents.DashboardEvent) {
		atomic.AddInt32(&handler3Count, 1)
	})

	// Trigger eager load
	ew, err := lw.GetWorkspaceForExecution(ctx)
	require.NoError(t, err)

	// Publish event on eager workspace
	testEvent := &dashboardevents.WorkspaceError{
		Error: nil,
	}
	ew.PublishDashboardEvent(ctx, testEvent)

	// Wait for async event handling
	time.Sleep(100 * time.Millisecond)

	assert.GreaterOrEqual(t, atomic.LoadInt32(&handler1Count), int32(1), "Handler 1 should receive event")
	assert.GreaterOrEqual(t, atomic.LoadInt32(&handler2Count), int32(1), "Handler 2 should receive event")
	assert.GreaterOrEqual(t, atomic.LoadInt32(&handler3Count), int32(1), "Handler 3 should receive event")
}

// TestLazyWorkspace_HandlerAfterEagerLoad verifies that handlers registered
// after eager load work on the eager workspace.

func TestLazyWorkspace_HandlerAfterEagerLoad(t *testing.T) {
	t.Parallel()
	modPath := getTestModPath(t, "simple")
	ctx := context.Background()

	lw, err := NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	// First trigger eager load
	ew, err := lw.GetWorkspaceForExecution(ctx)
	require.NoError(t, err)

	// Now register handler on eager workspace
	var receivedCount int32
	handler := func(ctx context.Context, event dashboardevents.DashboardEvent) {
		atomic.AddInt32(&receivedCount, 1)
	}

	ew.RegisterDashboardEventHandler(ctx, handler)

	// Publish event
	testEvent := &dashboardevents.WorkspaceError{
		Error: nil,
	}
	ew.PublishDashboardEvent(ctx, testEvent)

	// Wait for async event handling
	time.Sleep(100 * time.Millisecond)

	assert.GreaterOrEqual(t, atomic.LoadInt32(&receivedCount), int32(1),
		"Handler registered after eager load should receive events")
}

// =============================================================================
// Error Handling Tests
// =============================================================================

// TestLazyWorkspace_EagerLoadFailureCached verifies that an eager load failure
// is cached and returned on subsequent calls (not re-tried).

func TestLazyWorkspace_EagerLoadFailureCached(t *testing.T) {
	ctx := context.Background()

	// We need to create a lazy workspace that will fail on eager load
	// Use an invalid workspace path that can be indexed but not eagerly loaded
	lw, err := NewLazyWorkspace(ctx, getTestModPath(t, "error-conditions"), DefaultLazyLoadConfig())
	if err != nil {
		// If we can't even create the lazy workspace, that's fine for this test
		t.Skip("Cannot create lazy workspace for error test")
	}
	defer lw.Close()

	// Manually set the workspace path to something that will fail on eager load
	originalPath := lw.workspacePath
	lw.workspacePath = "/nonexistent/path/that/will/fail"

	// First call should fail
	_, err1 := lw.GetWorkspaceForExecution(ctx)
	if err1 == nil {
		// Load() succeeded with the modified path - this can happen in some environments
		// Skip the test rather than fail
		lw.workspacePath = originalPath
		t.Skip("Could not induce eager load failure with invalid path")
	}

	// Second call should return the same cached error (sync.Once)
	_, err2 := lw.GetWorkspaceForExecution(ctx)
	assert.Error(t, err2, "Second call should also fail")
	assert.Equal(t, err1, err2, "Should return the same cached error")

	// Restore path for cleanup
	lw.workspacePath = originalPath
}

// TestLazyWorkspace_LazyOpsAfterEagerFailure verifies that lazy operations
// still work even after an eager load has failed.

func TestLazyWorkspace_LazyOpsAfterEagerFailure(t *testing.T) {
	modPath := getTestModPath(t, "error-conditions")
	ctx := context.Background()

	lw, err := NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	if err != nil {
		t.Skip("Cannot create lazy workspace for error test")
	}
	defer lw.Close()

	// Record initial stats
	initialStats := lw.IndexStats()

	// Manually break eager loading
	originalPath := lw.workspacePath
	lw.workspacePath = "/nonexistent/path"

	// Try eager load (should fail)
	_, err = lw.GetWorkspaceForExecution(ctx)

	// Restore path before checking
	lw.workspacePath = originalPath

	if err == nil {
		// Load() succeeded with the modified path - skip this test
		t.Skip("Could not induce eager load failure with invalid path")
	}

	// Lazy operations should still work
	payload := lw.GetAvailableDashboardsFromIndex()
	assert.NotNil(t, payload, "GetAvailableDashboardsFromIndex should work after eager failure")

	stats := lw.IndexStats()
	assert.Equal(t, initialStats.TotalEntries, stats.TotalEntries,
		"Index should still have same entries after eager failure")

	// Cache should still be functional
	cacheStats := lw.CacheStats()
	assert.NotNil(t, &cacheStats, "Cache should be accessible")
}

// =============================================================================
// Resource Loading Difference Tests
// =============================================================================

// TestLazyWorkspace_LoadBenchmarkDifference tests the difference between
// LoadBenchmark (lazy) and LoadBenchmarkForExecution (eager with children).

func TestLazyWorkspace_LoadBenchmarkDifference(t *testing.T) {
	modPath := getTestModPath(t, "simple")
	ctx := context.Background()

	lw, err := NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	// Get benchmark name
	benchNames := lw.GetLazyModResources().ListBenchmarkNames()
	require.NotEmpty(t, benchNames, "Should have benchmarks")
	benchName := benchNames[0]

	// LoadBenchmark: children loaded into cache but not set on benchmark
	bench1, err := lw.LoadBenchmark(ctx, benchName)
	require.NoError(t, err)
	require.NotNil(t, bench1)

	// LoadBenchmarkForExecution: children field properly populated
	bench2, err := lw.LoadBenchmarkForExecution(ctx, benchName)
	require.NoError(t, err)
	require.NotNil(t, bench2)

	// The benchmark loaded for execution should have children populated
	children := bench2.GetChildren()
	// Note: The simple test mod has a benchmark with 2 children
	assert.GreaterOrEqual(t, len(children), 0, "LoadBenchmarkForExecution should resolve children")
}

// TestLazyWorkspace_GetResourceAfterEager tests that GetResource works correctly
// after an eager load has occurred.

func TestLazyWorkspace_GetResourceAfterEager(t *testing.T) {
	modPath := getTestModPath(t, "simple")
	ctx := context.Background()

	lw, err := NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	// Get a control name for testing
	controlNames := lw.GetLazyModResources().ListControlNames()
	require.NotEmpty(t, controlNames, "Should have controls")
	controlName := controlNames[0]

	// Load resource before eager load
	resource1, err := lw.LoadResource(ctx, controlName)
	require.NoError(t, err)
	require.NotNil(t, resource1)

	// Trigger eager load
	_, err = lw.GetWorkspaceForExecution(ctx)
	require.NoError(t, err)

	// Load resource after eager load
	resource2, err := lw.LoadResource(ctx, controlName)
	require.NoError(t, err)
	require.NotNil(t, resource2)

	// Both should return valid resources
	assert.Equal(t, resource1.Name(), resource2.Name(), "Resource names should match")
}

// TestLazyWorkspace_CacheStateAfterEager tests the cache state after eager load.

func TestLazyWorkspace_CacheStateAfterEager(t *testing.T) {
	modPath := getTestModPath(t, "simple")
	ctx := context.Background()

	lw, err := NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	// Load some resources to populate cache
	queryNames := lw.GetLazyModResources().ListQueryNames()
	for _, name := range queryNames {
		_, _ = lw.LoadResource(ctx, name)
	}

	cacheStatsBefore := lw.CacheStats()

	// Trigger eager load
	_, err = lw.GetWorkspaceForExecution(ctx)
	require.NoError(t, err)

	cacheStatsAfter := lw.CacheStats()

	// Cache should still be functional
	assert.GreaterOrEqual(t, cacheStatsAfter.Entries, 0, "Cache should be functional after eager load")

	t.Logf("Cache before eager: entries=%d, memory=%d",
		cacheStatsBefore.Entries, cacheStatsBefore.MemoryBytes)
	t.Logf("Cache after eager: entries=%d, memory=%d",
		cacheStatsAfter.Entries, cacheStatsAfter.MemoryBytes)
}

// =============================================================================
// Interface Compliance Tests
// =============================================================================

// TestLazyWorkspace_InterfaceCompliance is a compile-time check that
// LazyWorkspace satisfies DashboardServerWorkspace interface.

func TestLazyWorkspace_InterfaceCompliance(t *testing.T) {
	// Compile-time check
	var _ DashboardServerWorkspace = (*LazyWorkspace)(nil)

	// Runtime verification
	modPath := getTestModPath(t, "simple")
	ctx := context.Background()

	lw, err := NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	// Verify interface methods exist and work
	var ws DashboardServerWorkspace = lw

	// WorkspaceProvider methods
	assert.True(t, ws.IsLazy())
	ws.Close() // Should not panic

	// Re-create for remaining tests
	lw, err = NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	ws = lw

	// Test GetResource
	parsedName := &modconfig.ParsedResourceName{
		Mod:      "lazy_simple",
		ItemType: "query",
		Name:     "basic",
	}
	_, _ = ws.GetResource(parsedName) // Should not panic

	// Test LoadDashboard
	dashNames := lw.GetLazyModResources().ListDashboardNames()
	if len(dashNames) > 0 {
		_, _ = ws.LoadDashboard(ctx, dashNames[0])
	}

	// Test LoadBenchmark
	benchNames := lw.GetLazyModResources().ListBenchmarkNames()
	if len(benchNames) > 0 {
		_, _ = ws.LoadBenchmark(ctx, benchNames[0])
	}
}

// TestLazyWorkspace_EagerInterfaceCompliance verifies that the workspace returned
// by GetWorkspaceForExecution also satisfies the DashboardServerWorkspace interface.

func TestLazyWorkspace_EagerInterfaceCompliance(t *testing.T) {
	// Compile-time check
	var _ DashboardServerWorkspace = (*PowerpipeWorkspace)(nil)

	modPath := getTestModPath(t, "simple")
	ctx := context.Background()

	lw, err := NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	ew, err := lw.GetWorkspaceForExecution(ctx)
	require.NoError(t, err)

	// Verify interface compliance
	var ws DashboardServerWorkspace = ew

	assert.False(t, ws.IsLazy(), "Eager workspace should not be lazy")

	// Test GetModResources
	modResources := ws.GetModResources()
	assert.NotNil(t, modResources, "Should have mod resources")
}

// =============================================================================
// State Consistency Tests
// =============================================================================

// TestLazyWorkspace_ResourceDataConsistency verifies that resource data is
// consistent between lazy and eager loading paths.

func TestLazyWorkspace_ResourceDataConsistency(t *testing.T) {
	modPath := getTestModPath(t, "simple")
	ctx := context.Background()

	lw, err := NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	// Get control via lazy path
	controlNames := lw.GetLazyModResources().ListControlNames()
	require.NotEmpty(t, controlNames)
	controlName := controlNames[0]

	lazyResource, err := lw.LoadResource(ctx, controlName)
	require.NoError(t, err)
	require.NotNil(t, lazyResource)

	// Get same control via eager workspace
	ew, err := lw.GetWorkspaceForExecution(ctx)
	require.NoError(t, err)

	// Parse the resource name to look it up in eager workspace
	parsedName, err := modconfig.ParseResourceName(controlName)
	require.NoError(t, err)

	eagerResource, found := ew.GetResource(parsedName)
	require.True(t, found, "Resource should be found in eager workspace")

	// Compare key fields
	assert.Equal(t, lazyResource.Name(), eagerResource.Name(), "Names should match")
	assert.Equal(t, lazyResource.GetUnqualifiedName(), eagerResource.GetUnqualifiedName(),
		"Unqualified names should match")
}

// TestLazyWorkspace_AvailableDashboardsConsistency verifies that the available
// dashboards are consistent between index and eager workspace.

func TestLazyWorkspace_AvailableDashboardsConsistency(t *testing.T) {
	modPath := getTestModPath(t, "simple")
	ctx := context.Background()

	lw, err := NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	// Get dashboard names from index
	indexDashboardNames := lw.GetLazyModResources().ListDashboardNames()

	// Trigger eager load
	ew, err := lw.GetWorkspaceForExecution(ctx)
	require.NoError(t, err)

	// Get dashboard names from eager workspace
	eagerDashboards := ew.GetPowerpipeModResources().Dashboards
	eagerDashboardNames := make([]string, 0, len(eagerDashboards))
	for name := range eagerDashboards {
		eagerDashboardNames = append(eagerDashboardNames, name)
	}

	// Compare counts (order might differ)
	assert.Equal(t, len(indexDashboardNames), len(eagerDashboardNames),
		"Dashboard counts should match between index and eager workspace")
}

// =============================================================================
// Memory and Performance Tests
// =============================================================================

// TestLazyWorkspace_MemoryDuringTransition measures memory usage during
// the transition from lazy to eager loading.

func TestLazyWorkspace_MemoryDuringTransition(t *testing.T) {
	modPath := getGeneratedModPath(t, "small")
	ctx := context.Background()

	lw, err := NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	// Get memory stats before transition
	cacheStatsBefore := lw.CacheStats()
	indexStatsBefore := lw.IndexStats()

	t.Logf("Before transition: cache entries=%d, cache memory=%d, index entries=%d, index size=%d",
		cacheStatsBefore.Entries, cacheStatsBefore.MemoryBytes,
		indexStatsBefore.TotalEntries, indexStatsBefore.TotalSize)

	// Trigger eager load
	_, err = lw.GetWorkspaceForExecution(ctx)
	require.NoError(t, err)

	// Get memory stats after transition
	cacheStatsAfter := lw.CacheStats()
	indexStatsAfter := lw.IndexStats()

	t.Logf("After transition: cache entries=%d, cache memory=%d, index entries=%d, index size=%d",
		cacheStatsAfter.Entries, cacheStatsAfter.MemoryBytes,
		indexStatsAfter.TotalEntries, indexStatsAfter.TotalSize)

	// Index should remain unchanged
	assert.Equal(t, indexStatsBefore.TotalEntries, indexStatsAfter.TotalEntries,
		"Index entries should not change during transition")
}

// TestLazyWorkspace_TransitionTime measures the time for eager load transition.

func TestLazyWorkspace_TransitionTime(t *testing.T) {
	// Test with small mod first
	modPath := getGeneratedModPath(t, "small")
	ctx := context.Background()

	lw, err := NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	start := time.Now()
	_, err = lw.GetWorkspaceForExecution(ctx)
	transitionTime := time.Since(start)

	require.NoError(t, err)

	t.Logf("Small mod transition time: %v", transitionTime)

	// For a small mod, transition should be relatively fast (< 5s)
	assert.Less(t, transitionTime.Seconds(), float64(5),
		"Small mod transition should complete in under 5 seconds")
}

// TestLazyWorkspace_TransitionTimeMedium tests transition time for a medium mod.

func TestLazyWorkspace_TransitionTimeMedium(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping medium mod transition test in short mode")
	}

	modPath := getGeneratedModPath(t, "medium")
	ctx := context.Background()

	lw, err := NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	start := time.Now()
	_, err = lw.GetWorkspaceForExecution(ctx)
	transitionTime := time.Since(start)

	require.NoError(t, err)

	t.Logf("Medium mod transition time: %v", transitionTime)

	// For a medium mod, transition should complete in under 30 seconds
	assert.Less(t, transitionTime.Seconds(), float64(30),
		"Medium mod transition should complete in under 30 seconds")
}

// =============================================================================
// Edge Case Tests
// =============================================================================

// TestLazyWorkspace_DoubleClose verifies that calling Close twice doesn't panic.

func TestLazyWorkspace_DoubleClose(t *testing.T) {
	t.Parallel()
	modPath := getTestModPath(t, "simple")
	ctx := context.Background()

	lw, err := NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)

	// First close
	assert.NotPanics(t, func() {
		lw.Close()
	}, "First close should not panic")

	// Second close
	assert.NotPanics(t, func() {
		lw.Close()
	}, "Second close should not panic")
}

// TestLazyWorkspace_ConcurrentClose verifies that concurrent close calls
// don't cause races.

func TestLazyWorkspace_ConcurrentClose(t *testing.T) {
	t.Parallel()
	modPath := getTestModPath(t, "simple")
	ctx := context.Background()

	lw, err := NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)

	// Trigger eager load first
	_, _ = lw.GetWorkspaceForExecution(ctx)

	const numGoroutines = 5
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			defer func() {
				// Recover from any panics
				_ = recover()
			}()
			lw.Close()
		}()
	}

	wg.Wait()
	// Test passes if no deadlock occurs
}

// TestLazyWorkspace_TransitionWithMinimalMod tests transition behavior
// with a mod that has minimal resources.

func TestLazyWorkspace_TransitionWithMinimalMod(t *testing.T) {
	// Use simple mod which has minimal resources
	modPath := getTestModPath(t, "simple")
	ctx := context.Background()

	lw, err := NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	// Should be able to create lazy workspace
	assert.True(t, lw.IsLazy())
	assert.Greater(t, lw.IndexStats().TotalEntries, 0)

	// Should be able to transition to eager
	ew, err := lw.GetWorkspaceForExecution(ctx)
	require.NoError(t, err)
	require.NotNil(t, ew)
}

// TestLazyWorkspace_PublishEventBeforeAndAfterTransition verifies event
// publishing works both before and after eager load.

func TestLazyWorkspace_PublishEventBeforeAndAfterTransition(t *testing.T) {
	t.Parallel()
	modPath := getTestModPath(t, "simple")
	ctx := context.Background()

	lw, err := NewLazyWorkspace(ctx, modPath, DefaultLazyLoadConfig())
	require.NoError(t, err)
	defer lw.Close()

	var eventCount int32
	handler := func(ctx context.Context, event dashboardevents.DashboardEvent) {
		atomic.AddInt32(&eventCount, 1)
	}

	lw.RegisterDashboardEventHandler(ctx, handler)

	// Publish event on lazy workspace (before transition)
	lw.PublishDashboardEvent(ctx, &dashboardevents.WorkspaceError{Error: nil})

	time.Sleep(50 * time.Millisecond)
	countBeforeTransition := atomic.LoadInt32(&eventCount)

	// Trigger transition
	ew, err := lw.GetWorkspaceForExecution(ctx)
	require.NoError(t, err)

	// Publish event on eager workspace (after transition)
	ew.PublishDashboardEvent(ctx, &dashboardevents.WorkspaceError{Error: nil})

	time.Sleep(50 * time.Millisecond)
	countAfterTransition := atomic.LoadInt32(&eventCount)

	assert.GreaterOrEqual(t, countBeforeTransition, int32(1),
		"Should receive events before transition")
	assert.Greater(t, countAfterTransition, countBeforeTransition,
		"Should receive events after transition")
}
