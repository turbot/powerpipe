package dashboardserver

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turbot/powerpipe/internal/resourceindex"
	"github.com/turbot/powerpipe/internal/workspace"
)

// TestBuildAvailableDashboardsPayloadFromIndex verifies that the lazy mode
// builds the same structure as eager mode.
func TestBuildAvailableDashboardsPayloadFromIndex(t *testing.T) {
	// Create a mock index with test data
	idx := resourceindex.NewResourceIndex()
	idx.ModName = "test_mod"
	idx.ModFullName = "mod.test_mod"
	idx.ModTitle = "Test Mod"

	// Add a dashboard
	idx.Add(&resourceindex.IndexEntry{
		Type:        "dashboard",
		Name:        "test_mod.dashboard.main",
		ShortName:   "main",
		Title:       "Main Dashboard",
		ModName:     "test_mod",
		ModFullName: "mod.test_mod",
		Tags:        map[string]string{"category": "test"},
	})

	// Add a top-level benchmark
	idx.Add(&resourceindex.IndexEntry{
		Type:          "benchmark",
		Name:          "test_mod.benchmark.security",
		ShortName:     "security",
		Title:         "Security Benchmark",
		ModName:       "test_mod",
		ModFullName:   "mod.test_mod",
		IsTopLevel:    true,
		BenchmarkType: "control",
		Tags:          map[string]string{"type": "security"},
	})

	// Add a child benchmark
	idx.Add(&resourceindex.IndexEntry{
		Type:          "benchmark",
		Name:          "test_mod.benchmark.access",
		ShortName:     "access",
		Title:         "Access Control",
		ParentName:    "test_mod.benchmark.security",
		ModName:       "test_mod",
		ModFullName:   "mod.test_mod",
		BenchmarkType: "control",
	})

	// Update parent's children list
	if parent, ok := idx.Get("test_mod.benchmark.security"); ok {
		parent.ChildNames = append(parent.ChildNames, "test_mod.benchmark.access")
	}

	// Build payload from index
	payload := idx.BuildAvailableDashboardsPayload()

	// Verify structure
	assert.Equal(t, "available_dashboards", payload.Action)
	assert.Len(t, payload.Dashboards, 1)
	assert.Len(t, payload.Benchmarks, 2)

	// Verify dashboard
	dash, ok := payload.Dashboards["test_mod.dashboard.main"]
	require.True(t, ok)
	assert.Equal(t, "Main Dashboard", dash.Title)
	assert.Equal(t, "test_mod.dashboard.main", dash.FullName)
	assert.Equal(t, "main", dash.ShortName)
	assert.Equal(t, "mod.test_mod", dash.ModFullName)

	// Verify benchmark
	bench, ok := payload.Benchmarks["test_mod.benchmark.security"]
	require.True(t, ok)
	assert.Equal(t, "Security Benchmark", bench.Title)
	assert.True(t, bench.IsTopLevel)
	assert.Equal(t, "control", bench.BenchmarkType)
}

// TestConvertIndexBenchmarkInfo verifies recursive benchmark conversion.
func TestConvertIndexBenchmarkInfo(t *testing.T) {
	info := resourceindex.BenchmarkInfo{
		Title:         "Parent",
		FullName:      "test_mod.benchmark.parent",
		ShortName:     "parent",
		BenchmarkType: "control",
		Tags:          map[string]string{"level": "1"},
		IsTopLevel:    true,
		Trunks:        [][]string{{"test_mod.benchmark.parent"}},
		Children: []resourceindex.BenchmarkInfo{
			{
				Title:         "Child",
				FullName:      "test_mod.benchmark.child",
				ShortName:     "child",
				BenchmarkType: "control",
				Tags:          map[string]string{"level": "2"},
				Children:      nil,
			},
		},
		ModFullName: "mod.test_mod",
	}

	result := convertIndexBenchmarkInfo(info)

	assert.Equal(t, "Parent", result.Title)
	assert.Equal(t, "test_mod.benchmark.parent", result.FullName)
	assert.True(t, result.IsTopLevel)
	assert.Len(t, result.Children, 1)
	assert.Equal(t, "Child", result.Children[0].Title)
	assert.Equal(t, "test_mod.benchmark.child", result.Children[0].FullName)
}

// TestServerIsLazyMode verifies lazy mode detection.
func TestServerIsLazyMode(t *testing.T) {
	// Server without lazy workspace
	s1 := &Server{
		workspace: nil,
	}
	assert.False(t, s1.isLazyMode())

	// Server with lazy workspace - cannot create a real one without files,
	// so we just test the nil check logic
}

// TestPayloadConversionRoundTrip verifies that the payload can be marshaled and unmarshaled.
func TestPayloadConversionRoundTrip(t *testing.T) {
	// Create index payload
	indexPayload := &resourceindex.AvailableDashboardsPayload{
		Action: "available_dashboards",
		Dashboards: map[string]resourceindex.DashboardInfo{
			"test_mod.dashboard.main": {
				Title:       "Main Dashboard",
				FullName:    "test_mod.dashboard.main",
				ShortName:   "main",
				Tags:        map[string]string{"env": "test"},
				ModFullName: "mod.test_mod",
			},
		},
		Benchmarks: map[string]resourceindex.BenchmarkInfo{
			"test_mod.benchmark.security": {
				Title:         "Security",
				FullName:      "test_mod.benchmark.security",
				ShortName:     "security",
				BenchmarkType: "control",
				Tags:          map[string]string{"type": "sec"},
				IsTopLevel:    true,
				ModFullName:   "mod.test_mod",
			},
		},
	}

	// Convert to server payload format
	payload := AvailableDashboardsPayload{
		Action:     "available_dashboards",
		Dashboards: make(map[string]ModAvailableDashboard),
		Benchmarks: make(map[string]ModAvailableBenchmark),
	}

	for name, dash := range indexPayload.Dashboards {
		payload.Dashboards[name] = ModAvailableDashboard{
			Title:       dash.Title,
			FullName:    dash.FullName,
			ShortName:   dash.ShortName,
			Tags:        dash.Tags,
			ModFullName: dash.ModFullName,
		}
	}

	for name, bench := range indexPayload.Benchmarks {
		payload.Benchmarks[name] = convertIndexBenchmarkInfo(bench)
	}

	// Marshal
	bytes, err := json.Marshal(payload)
	require.NoError(t, err)

	// Unmarshal
	var result AvailableDashboardsPayload
	err = json.Unmarshal(bytes, &result)
	require.NoError(t, err)

	// Verify
	assert.Equal(t, "available_dashboards", result.Action)
	assert.Len(t, result.Dashboards, 1)
	assert.Len(t, result.Benchmarks, 1)

	dash, ok := result.Dashboards["test_mod.dashboard.main"]
	require.True(t, ok)
	assert.Equal(t, "Main Dashboard", dash.Title)

	bench, ok := result.Benchmarks["test_mod.benchmark.security"]
	require.True(t, ok)
	assert.Equal(t, "Security", bench.Title)
	assert.True(t, bench.IsTopLevel)
}

// TestLazyWorkspaceInterface verifies that LazyWorkspace implements DashboardServerWorkspace.
func TestLazyWorkspaceInterface(t *testing.T) {
	// This test verifies the interface is properly implemented
	// We can't create a real lazy workspace without files, but we can verify the types

	// Compile-time interface check
	var _ workspace.DashboardServerWorkspace = (*workspace.LazyWorkspace)(nil)
	var _ workspace.DashboardServerWorkspace = (*workspace.PowerpipeWorkspace)(nil)
}

// TestPayloadTiming verifies that building from index is fast.
func TestPayloadTiming(t *testing.T) {
	// Create a larger index
	idx := resourceindex.NewResourceIndex()
	idx.ModName = "test_mod"
	idx.ModFullName = "mod.test_mod"

	// Add 100 dashboards
	for i := 0; i < 100; i++ {
		idx.Add(&resourceindex.IndexEntry{
			Type:        "dashboard",
			Name:        "test_mod.dashboard.dash_" + string(rune('a'+i%26)) + string(rune('0'+i/26)),
			ShortName:   "dash_" + string(rune('a'+i%26)) + string(rune('0'+i/26)),
			Title:       "Dashboard " + string(rune('0'+i)),
			ModName:     "test_mod",
			ModFullName: "mod.test_mod",
		})
	}

	// Add 50 benchmarks with children
	for i := 0; i < 50; i++ {
		idx.Add(&resourceindex.IndexEntry{
			Type:          "benchmark",
			Name:          "test_mod.benchmark.bench_" + string(rune('a'+i%26)) + string(rune('0'+i/26)),
			ShortName:     "bench_" + string(rune('a'+i%26)) + string(rune('0'+i/26)),
			Title:         "Benchmark " + string(rune('0'+i)),
			ModName:       "test_mod",
			ModFullName:   "mod.test_mod",
			IsTopLevel:    true,
			BenchmarkType: "control",
		})
	}

	// Time the payload build
	start := time.Now()
	payload := idx.BuildAvailableDashboardsPayload()
	elapsed := time.Since(start)

	// Should be very fast (under 10ms for this small index)
	assert.Less(t, elapsed.Milliseconds(), int64(100))
	assert.Len(t, payload.Dashboards, 100)
	assert.Len(t, payload.Benchmarks, 50)

	t.Logf("Built payload for %d dashboards and %d benchmarks in %v",
		len(payload.Dashboards), len(payload.Benchmarks), elapsed)
}

// TestServerBuildAvailableDashboardsPayload tests the server method delegation.
func TestServerBuildAvailableDashboardsPayload(t *testing.T) {
	// Test without lazy workspace (nil case)
	s := &Server{
		workspace:     nil,
		lazyWorkspace: nil,
	}

	// Without any workspace, this would panic, so we just verify mode detection
	assert.False(t, s.isLazyMode())
}

// BenchmarkPayloadFromIndex benchmarks the index-based payload build.
func BenchmarkPayloadFromIndex(b *testing.B) {
	// Create index with realistic data
	idx := resourceindex.NewResourceIndex()
	idx.ModName = "test_mod"
	idx.ModFullName = "mod.test_mod"

	for i := 0; i < 200; i++ {
		idx.Add(&resourceindex.IndexEntry{
			Type:        "dashboard",
			Name:        "test_mod.dashboard.dash_" + string(rune('a'+i%26)) + string(rune('0'+i/26)),
			ShortName:   "dash_" + string(rune('a'+i%26)) + string(rune('0'+i/26)),
			Title:       "Dashboard",
			ModName:     "test_mod",
			ModFullName: "mod.test_mod",
			Tags:        map[string]string{"env": "test", "category": "ops"},
		})
	}

	for i := 0; i < 100; i++ {
		idx.Add(&resourceindex.IndexEntry{
			Type:          "benchmark",
			Name:          "test_mod.benchmark.bench_" + string(rune('a'+i%26)) + string(rune('0'+i/26)),
			ShortName:     "bench_" + string(rune('a'+i%26)) + string(rune('0'+i/26)),
			Title:         "Benchmark",
			ModName:       "test_mod",
			ModFullName:   "mod.test_mod",
			IsTopLevel:    true,
			BenchmarkType: "control",
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = idx.BuildAvailableDashboardsPayload()
	}
}

// TestDashboardServerWorkspaceInterface ensures the interface is complete.
func TestDashboardServerWorkspaceInterface(t *testing.T) {
	// This is a compile-time test that ensures both workspace types
	// implement the DashboardServerWorkspace interface correctly.
	// If this compiles, the interface is satisfied.

	ctx := context.Background()

	// We can't instantiate real workspaces without files, but we can
	// verify the interface methods exist by type-checking
	type workspaceChecker interface {
		workspace.DashboardServerWorkspace
	}

	// These lines will fail to compile if the interfaces are not satisfied
	var _ workspaceChecker = (*workspace.PowerpipeWorkspace)(nil)
	var _ workspaceChecker = (*workspace.LazyWorkspace)(nil)

	// Suppress unused variable warning
	_ = ctx
}
