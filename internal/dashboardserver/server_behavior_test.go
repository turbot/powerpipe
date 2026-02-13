package dashboardserver

import (
	"context"
	"encoding/json"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turbot/powerpipe/internal/resources"
	"github.com/turbot/powerpipe/internal/workspace"
)

func testdataDir() string {
	// Find testdata directory relative to this test file using runtime.Caller
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "testdata")
}

// TestServer_GetAvailableDashboards verifies that the available dashboards payload
// is correctly built from workspace resources.
// This is a behavior test that must pass before AND after lazy loading implementation.
func TestServer_GetAvailableDashboards(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Build the available dashboards payload
	payloadBytes, err := buildAvailableDashboardsPayload(res)
	require.NoError(t, err)
	require.NotEmpty(t, payloadBytes)

	// Parse the payload
	var payload AvailableDashboardsPayload
	err = json.Unmarshal(payloadBytes, &payload)
	require.NoError(t, err)

	// Verify payload structure
	assert.Equal(t, "available_dashboards", payload.Action)

	// Verify dashboards are present
	assert.NotEmpty(t, payload.Dashboards, "payload should contain dashboards")

	// Verify expected dashboards exist
	mainDash, ok := payload.Dashboards["behavior_test.dashboard.main"]
	assert.True(t, ok, "should have main dashboard")
	assert.Equal(t, "behavior_test.dashboard.main", mainDash.FullName)
	assert.Equal(t, "main", mainDash.ShortName)
	assert.Equal(t, "Main Dashboard", mainDash.Title)

	// Verify benchmarks are present
	assert.NotEmpty(t, payload.Benchmarks, "payload should contain benchmarks")

	// Verify expected benchmarks exist
	topBench, ok := payload.Benchmarks["behavior_test.benchmark.top"]
	assert.True(t, ok, "should have top benchmark")
	assert.Equal(t, "behavior_test.benchmark.top", topBench.FullName)
	assert.Equal(t, "top", topBench.ShortName)
	assert.Equal(t, "Top Level Benchmark", topBench.Title)
	assert.True(t, topBench.IsTopLevel, "top benchmark should be top-level")
}

// TestServer_DashboardPayloadFields verifies that dashboard payload fields
// are correctly populated.
func TestServer_DashboardPayloadFields(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	payloadBytes, err := buildAvailableDashboardsPayload(res)
	require.NoError(t, err)

	var payload AvailableDashboardsPayload
	err = json.Unmarshal(payloadBytes, &payload)
	require.NoError(t, err)

	// Check all dashboards have required fields
	for name, dash := range payload.Dashboards {
		assert.NotEmpty(t, dash.FullName, "dashboard %s should have FullName", name)
		assert.NotEmpty(t, dash.ShortName, "dashboard %s should have ShortName", name)
		assert.NotEmpty(t, dash.ModFullName, "dashboard %s should have ModFullName", name)
		assert.Equal(t, name, dash.FullName, "dashboard key should match FullName")
	}
}

// TestServer_BenchmarkPayloadFields verifies that benchmark payload fields
// are correctly populated.
func TestServer_BenchmarkPayloadFields(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	payloadBytes, err := buildAvailableDashboardsPayload(res)
	require.NoError(t, err)

	var payload AvailableDashboardsPayload
	err = json.Unmarshal(payloadBytes, &payload)
	require.NoError(t, err)

	// Check all benchmarks have required fields
	for name, bench := range payload.Benchmarks {
		assert.NotEmpty(t, bench.FullName, "benchmark %s should have FullName", name)
		assert.NotEmpty(t, bench.ShortName, "benchmark %s should have ShortName", name)
		assert.NotEmpty(t, bench.ModFullName, "benchmark %s should have ModFullName", name)
		assert.Equal(t, name, bench.FullName, "benchmark key should match FullName")
		assert.NotEmpty(t, bench.BenchmarkType, "benchmark %s should have BenchmarkType", name)
	}
}

// TestServer_BenchmarkHierarchy verifies that benchmark hierarchy (trunks)
// is correctly represented in the payload.
func TestServer_BenchmarkHierarchy(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	payloadBytes, err := buildAvailableDashboardsPayload(res)
	require.NoError(t, err)

	var payload AvailableDashboardsPayload
	err = json.Unmarshal(payloadBytes, &payload)
	require.NoError(t, err)

	// Find top-level benchmark
	topBench, ok := payload.Benchmarks["behavior_test.benchmark.top"]
	require.True(t, ok, "should have top benchmark")

	// Verify top-level benchmark has trunks
	assert.True(t, topBench.IsTopLevel, "top benchmark should be top-level")
	assert.NotEmpty(t, topBench.Trunks, "top-level benchmark should have trunks")

	// Verify children are populated
	assert.NotEmpty(t, topBench.Children, "top benchmark should have children")
}

// TestServer_BenchmarkChildren verifies that benchmark children are
// correctly nested in the payload.
func TestServer_BenchmarkChildren(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	payloadBytes, err := buildAvailableDashboardsPayload(res)
	require.NoError(t, err)

	var payload AvailableDashboardsPayload
	err = json.Unmarshal(payloadBytes, &payload)
	require.NoError(t, err)

	// Find top-level benchmark
	topBench := payload.Benchmarks["behavior_test.benchmark.top"]

	// Count children - top should have child_a and child_b benchmarks
	childBenchmarkCount := 0
	for _, child := range topBench.Children {
		if child.FullName == "behavior_test.benchmark.child_a" ||
			child.FullName == "behavior_test.benchmark.child_b" {
			childBenchmarkCount++
		}
	}
	assert.Equal(t, 2, childBenchmarkCount, "top benchmark should have 2 child benchmarks")
}

// TestServer_DashboardTags verifies that dashboard tags are included in the payload.
func TestServer_DashboardTags(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	payloadBytes, err := buildAvailableDashboardsPayload(res)
	require.NoError(t, err)

	var payload AvailableDashboardsPayload
	err = json.Unmarshal(payloadBytes, &payload)
	require.NoError(t, err)

	// Main dashboard should have tags
	mainDash := payload.Dashboards["behavior_test.dashboard.main"]
	assert.NotEmpty(t, mainDash.Tags, "main dashboard should have tags")
}

// TestServer_BenchmarkTags verifies that benchmark tags are included in the payload.
func TestServer_BenchmarkTags(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	payloadBytes, err := buildAvailableDashboardsPayload(res)
	require.NoError(t, err)

	var payload AvailableDashboardsPayload
	err = json.Unmarshal(payloadBytes, &payload)
	require.NoError(t, err)

	// Top benchmark should have tags
	topBench := payload.Benchmarks["behavior_test.benchmark.top"]
	assert.NotEmpty(t, topBench.Tags, "top benchmark should have tags")
}

// TestServer_PayloadDeterministic verifies that building the payload multiple times
// produces consistent results.
func TestServer_PayloadDeterministic(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	// Build payload multiple times
	payload1, err := buildAvailableDashboardsPayload(res)
	require.NoError(t, err)

	payload2, err := buildAvailableDashboardsPayload(res)
	require.NoError(t, err)

	// Parse both payloads
	var p1, p2 AvailableDashboardsPayload
	require.NoError(t, json.Unmarshal(payload1, &p1))
	require.NoError(t, json.Unmarshal(payload2, &p2))

	// Verify counts match
	assert.Equal(t, len(p1.Dashboards), len(p2.Dashboards), "dashboard count should be consistent")
	assert.Equal(t, len(p1.Benchmarks), len(p2.Benchmarks), "benchmark count should be consistent")

	// Verify same dashboards exist
	for name := range p1.Dashboards {
		_, ok := p2.Dashboards[name]
		assert.True(t, ok, "dashboard %s should exist in both payloads", name)
	}

	// Verify same benchmarks exist
	for name := range p1.Benchmarks {
		_, ok := p2.Benchmarks[name]
		assert.True(t, ok, "benchmark %s should exist in both payloads", name)
	}
}

// TestServer_ResourceCountsMatch verifies that the payload contains the same
// number of resources as the workspace.
func TestServer_ResourceCountsMatch(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	payloadBytes, err := buildAvailableDashboardsPayload(res)
	require.NoError(t, err)

	var payload AvailableDashboardsPayload
	require.NoError(t, json.Unmarshal(payloadBytes, &payload))

	// Get top-level resources to compare
	topRes := resources.GetModResources(res.Mod)

	// Dashboard count should match
	assert.Equal(t, len(topRes.Dashboards), len(payload.Dashboards),
		"payload dashboard count should match workspace")

	// Non-anonymous benchmark count should match
	nonAnonymousBenchmarks := 0
	for _, b := range topRes.ControlBenchmarks {
		if !b.IsAnonymous() {
			nonAnonymousBenchmarks++
		}
	}
	assert.Equal(t, nonAnonymousBenchmarks, len(payload.Benchmarks),
		"payload benchmark count should match workspace (excluding anonymous)")
}

// TestServer_AllDashboardsHaveModFullName verifies that all dashboards
// have a valid ModFullName reference.
func TestServer_AllDashboardsHaveModFullName(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	payloadBytes, err := buildAvailableDashboardsPayload(res)
	require.NoError(t, err)

	var payload AvailableDashboardsPayload
	require.NoError(t, json.Unmarshal(payloadBytes, &payload))

	for name, dash := range payload.Dashboards {
		assert.NotEmpty(t, dash.ModFullName, "dashboard %s should have ModFullName", name)
		assert.Equal(t, "mod.behavior_test", dash.ModFullName,
			"dashboard %s ModFullName should be mod.behavior_test", name)
	}
}

// TestServer_TopLevelBenchmarksHaveTrunks verifies that all top-level benchmarks
// have properly populated trunks.
func TestServer_TopLevelBenchmarksHaveTrunks(t *testing.T) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	res := w.GetPowerpipeModResources()

	payloadBytes, err := buildAvailableDashboardsPayload(res)
	require.NoError(t, err)

	var payload AvailableDashboardsPayload
	require.NoError(t, json.Unmarshal(payloadBytes, &payload))

	topLevelCount := 0
	for name, bench := range payload.Benchmarks {
		if bench.IsTopLevel {
			topLevelCount++
			assert.NotEmpty(t, bench.Trunks, "top-level benchmark %s should have trunks", name)
		}
	}

	assert.Greater(t, topLevelCount, 0, "should have at least one top-level benchmark")
}
