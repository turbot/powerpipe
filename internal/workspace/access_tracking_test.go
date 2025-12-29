package workspace

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turbot/pipe-fittings/v2/modconfig"
)

// AccessTracker records resource accesses for analysis of lazy loading patterns
type AccessTracker struct {
	mu       sync.Mutex
	accesses []AccessRecord
}

// AccessRecord represents a single resource access event
type AccessRecord struct {
	ResourceType string
	ResourceName string
	AccessType   string // "get", "iterate", "children", "walk"
	CallerFile   string
	CallerLine   int
}

// NewAccessTracker creates a new access tracker
func NewAccessTracker() *AccessTracker {
	return &AccessTracker{
		accesses: make([]AccessRecord, 0),
	}
}

// Record records an access event
func (t *AccessTracker) Record(resourceType, resourceName, accessType string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Get caller info (skip 2 frames: Record -> caller -> actual code)
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}

	t.accesses = append(t.accesses, AccessRecord{
		ResourceType: resourceType,
		ResourceName: resourceName,
		AccessType:   accessType,
		CallerFile:   filepath.Base(file),
		CallerLine:   line,
	})
}

// Report generates a string report of all accesses
func (t *AccessTracker) Report() string {
	t.mu.Lock()
	defer t.mu.Unlock()

	var b strings.Builder
	b.WriteString("Access Pattern Report\n")
	b.WriteString("=====================\n\n")

	// Group by access type
	byType := make(map[string][]AccessRecord)
	for _, a := range t.accesses {
		byType[a.AccessType] = append(byType[a.AccessType], a)
	}

	for accessType, records := range byType {
		b.WriteString(fmt.Sprintf("## %s accesses (%d)\n", accessType, len(records)))
		for _, r := range records {
			b.WriteString(fmt.Sprintf("  - %s.%s (from %s:%d)\n",
				r.ResourceType, r.ResourceName, r.CallerFile, r.CallerLine))
		}
		b.WriteString("\n")
	}

	return b.String()
}

// HasAccess checks if a specific access pattern was recorded
// Use "*" for resourceName to match any resource of the type
func (t *AccessTracker) HasAccess(resourceType, resourceName, accessType string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, a := range t.accesses {
		if a.ResourceType == resourceType && a.AccessType == accessType {
			if resourceName == "*" || a.ResourceName == resourceName {
				return true
			}
		}
	}
	return false
}

// CountAccesses counts accesses matching the criteria
// Use "*" to match any value for resourceType or resourceName
func (t *AccessTracker) CountAccesses(resourceType, resourceName, accessType string) int {
	t.mu.Lock()
	defer t.mu.Unlock()

	count := 0
	for _, a := range t.accesses {
		typeMatch := resourceType == "*" || a.ResourceType == resourceType
		nameMatch := resourceName == "*" || a.ResourceName == resourceName
		accessMatch := accessType == "*" || a.AccessType == accessType
		if typeMatch && nameMatch && accessMatch {
			count++
		}
	}
	return count
}

// GetAccessesByType returns all accesses of a given type
func (t *AccessTracker) GetAccessesByType(accessType string) []AccessRecord {
	t.mu.Lock()
	defer t.mu.Unlock()

	var result []AccessRecord
	for _, a := range t.accesses {
		if a.AccessType == accessType {
			result = append(result, a)
		}
	}
	return result
}

// Clear resets the tracker
func (t *AccessTracker) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.accesses = make([]AccessRecord, 0)
}

// trackChildAccess recursively tracks child access for a ModTreeItem
func trackChildAccess(tracker *AccessTracker, resource modconfig.ModTreeItem, depth int) {
	if depth > 10 {
		// Prevent infinite recursion
		return
	}

	for _, child := range resource.GetChildren() {
		tracker.Record(child.GetBlockType(), child.Name(), "children")
		if treeItem, ok := child.(modconfig.ModTreeItem); ok {
			trackChildAccess(tracker, treeItem, depth+1)
		}
	}
}

// TestAccessPatterns_AvailableDashboards tests the access pattern for building available dashboards
// This simulates what happens during server startup
func TestAccessPatterns_AvailableDashboards(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "benchmark-only")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())
	require.NotNil(t, w)

	tracker := NewAccessTracker()
	res := w.GetPowerpipeModResources()

	// Simulate buildAvailableDashboardsPayload behavior
	// Track dashboard iteration
	for name := range res.Dashboards {
		tracker.Record("dashboard", name, "iterate")
	}

	// Track benchmark iteration
	for name, bench := range res.ControlBenchmarks {
		tracker.Record("benchmark", name, "iterate")
		trackChildAccess(tracker, bench, 0)
	}

	t.Log("\n" + tracker.Report())

	// Verify expected access patterns
	assert.True(t, tracker.HasAccess("benchmark", "*", "iterate"),
		"Should iterate benchmarks")

	// Count total iteration accesses
	benchmarkIterates := tracker.CountAccesses("benchmark", "*", "iterate")
	t.Logf("Total benchmark iterations: %d", benchmarkIterates)

	// Verify child access for benchmarks
	childAccesses := tracker.CountAccesses("benchmark", "*", "children")
	t.Logf("Total benchmark child accesses: %d", childAccesses)
}

// TestAccessPatterns_DashboardExecution tests access patterns during dashboard execution
func TestAccessPatterns_DashboardExecution(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "simple-mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())
	require.NotNil(t, w)

	tracker := NewAccessTracker()
	res := w.GetPowerpipeModResources()

	// Simulate executing a specific dashboard (not iterating all)
	dashboardName := "simple_test.dashboard.simple_dashboard"
	dash, exists := res.Dashboards[dashboardName]
	require.True(t, exists, "Dashboard should exist")

	tracker.Record("dashboard", dash.FullName, "get")

	// Track child access
	trackChildAccess(tracker, dash, 0)

	t.Log("\n" + tracker.Report())

	// Verify: should NOT iterate all dashboards, only access one
	dashIterates := tracker.CountAccesses("dashboard", "*", "iterate")
	assert.Equal(t, 0, dashIterates, "Should not iterate all dashboards")

	// Verify: should access one dashboard via get
	dashGets := tracker.CountAccesses("dashboard", dashboardName, "get")
	assert.Equal(t, 1, dashGets, "Should get exactly one dashboard")
}

// TestAccessPatterns_BenchmarkExecution tests access patterns during benchmark execution
func TestAccessPatterns_BenchmarkExecution(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "complex-mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())
	require.NotNil(t, w)

	tracker := NewAccessTracker()
	res := w.GetPowerpipeModResources()

	// Simulate executing a specific benchmark
	benchmarkName := "complex_test.benchmark.test_benchmark"
	bench, exists := res.ControlBenchmarks[benchmarkName]
	require.True(t, exists, "Benchmark should exist")

	tracker.Record("benchmark", bench.FullName, "get")

	// Track child access (simulating result group creation)
	trackChildAccess(tracker, bench, 0)

	t.Log("\n" + tracker.Report())

	// Verify: should NOT iterate all benchmarks
	benchIterates := tracker.CountAccesses("benchmark", "*", "iterate")
	assert.Equal(t, 0, benchIterates, "Should not iterate all benchmarks")

	// Verify: accessed child benchmarks and controls
	childBenchmarks := tracker.GetAccessesByType("children")
	t.Logf("Child accesses: %d", len(childBenchmarks))
	for _, access := range childBenchmarks {
		t.Logf("  - %s.%s", access.ResourceType, access.ResourceName)
	}
}

// TestAccessPatterns_WalkResources tests the impact of WalkResources
func TestAccessPatterns_WalkResources(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "complex-mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())
	require.NotNil(t, w)

	tracker := NewAccessTracker()
	res := w.GetPowerpipeModResources()

	// Walk all resources and track accesses
	err := res.WalkResources(func(item modconfig.HclResource) (bool, error) {
		tracker.Record(item.GetBlockType(), item.Name(), "walk")
		return true, nil
	})
	require.NoError(t, err)

	t.Log("\n" + tracker.Report())

	// Verify: walk touches every resource
	walkAccesses := tracker.GetAccessesByType("walk")
	t.Logf("Total resources walked: %d", len(walkAccesses))

	// Group by resource type
	byType := make(map[string]int)
	for _, a := range walkAccesses {
		byType[a.ResourceType]++
	}

	t.Log("\nResources by type:")
	for resType, count := range byType {
		t.Logf("  - %s: %d", resType, count)
	}
}

// TestAccessPatterns_QueryResolution tests query reference resolution patterns
func TestAccessPatterns_QueryResolution(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "complex-mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())
	require.NotNil(t, w)

	tracker := NewAccessTracker()
	res := w.GetPowerpipeModResources()

	// Simulate control query resolution
	controlName := "complex_test.control.test_control"
	control, exists := res.Controls[controlName]
	require.True(t, exists, "Control should exist")

	tracker.Record("control", control.FullName, "get")

	// Check if control has a query reference
	if control.Query != nil {
		tracker.Record("query", control.Query.FullName, "reference")
	}

	// Check if SQL property references a query
	if control.SQL != nil {
		// Try to resolve as query name
		parsedName, err := modconfig.ParseResourceName(*control.SQL)
		if err == nil {
			// Check if this is a query reference
			if _, queryExists := res.Queries[parsedName.ToResourceName()]; queryExists {
				tracker.Record("query", parsedName.ToResourceName(), "sql_reference")
			}
		}
	}

	t.Log("\n" + tracker.Report())
}

// TestAccessPatterns_ResourceLookup tests the GetResource pattern
func TestAccessPatterns_ResourceLookup(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "complex-mod")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())
	require.NotNil(t, w)

	tracker := NewAccessTracker()

	// Test various resource lookups
	testCases := []struct {
		name         string
		expectedType string
	}{
		{"complex_test.query.parameterized_query", "query"},
		{"complex_test.dashboard.complex_dashboard", "dashboard"},
		{"complex_test.control.test_control", "control"},
		{"complex_test.benchmark.test_benchmark", "benchmark"},
	}

	for _, tc := range testCases {
		parsedName, err := modconfig.ParseResourceName(tc.name)
		require.NoError(t, err)

		resource, found := w.GetResource(parsedName)
		if found {
			tracker.Record(tc.expectedType, resource.Name(), "get")
		}
	}

	t.Log("\n" + tracker.Report())

	// Verify all lookups were recorded
	getAccesses := tracker.CountAccesses("*", "*", "get")
	assert.Equal(t, len(testCases), getAccesses, "Should have one get per test case")
}

// TestLazyLoadingRequirements documents what lazy loading must support
func TestLazyLoadingRequirements(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "benchmark-only")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())
	require.NotNil(t, w)

	res := w.GetPowerpipeModResources()

	// Document what information is needed at different stages

	t.Log("=== LAZY LOADING REQUIREMENTS ===\n")

	// 1. Startup - Available Dashboards
	t.Log("1. Server Startup (Available Dashboards List)")
	t.Log("   Required fields for each dashboard/benchmark:")
	t.Log("   - FullName")
	t.Log("   - ShortName")
	t.Log("   - Title")
	t.Log("   - Tags")
	t.Log("   - Mod.GetFullName()")
	t.Log("   - GetParents() - for IsTopLevel check")
	t.Log("   - GetChildren() - for hierarchy")
	t.Log("")

	// 2. Count resources that would need index entries
	indexEntries := len(res.Dashboards) + len(res.ControlBenchmarks) + len(res.DetectionBenchmarks)
	t.Logf("   Total index entries needed: %d\n", indexEntries)

	// 3. Execution - Full resource needed
	t.Log("2. Dashboard/Benchmark Execution")
	t.Log("   Required: Full resource definition including:")
	t.Log("   - SQL statements")
	t.Log("   - Param definitions")
	t.Log("   - Display properties")
	t.Log("   - All children (recursive)")
	t.Log("")

	// 4. File Watcher - Diff operation
	t.Log("3. File Watcher (Resource Change Detection)")
	t.Log("   Requires comparison of ALL resources")
	t.Log("   Consider: Index-based diff for initial check")
	t.Log("")

	// Count total resources
	totalResources := 0
	_ = res.WalkResources(func(item modconfig.HclResource) (bool, error) {
		totalResources++
		return true, nil
	})
	t.Logf("   Total resources in workspace: %d\n", totalResources)
}

// TestAccessPatterns_ChildRelationships documents parent-child relationships
func TestAccessPatterns_ChildRelationships(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "benchmark-only")

	w, ew := Load(ctx, modPath, WithVariableValidation(false))
	require.NoError(t, ew.GetError())
	require.NotNil(t, w)

	res := w.GetPowerpipeModResources()

	t.Log("=== PARENT-CHILD RELATIONSHIPS ===\n")

	// Document benchmark hierarchy
	for name, bench := range res.ControlBenchmarks {
		children := bench.GetChildren()
		if len(children) > 0 {
			t.Logf("Benchmark: %s", name)
			for _, child := range children {
				t.Logf("  └─ %s: %s", child.GetBlockType(), child.Name())
			}
			t.Log("")
		}
	}

	// Document dashboard hierarchy
	for name, dash := range res.Dashboards {
		children := dash.GetChildren()
		if len(children) > 0 {
			t.Logf("Dashboard: %s", name)
			for _, child := range children {
				t.Logf("  └─ %s: %s", child.GetBlockType(), child.Name())
			}
			t.Log("")
		}
	}
}
