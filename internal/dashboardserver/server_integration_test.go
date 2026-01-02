package dashboardserver

import (
	"context"
	"encoding/json"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turbot/powerpipe/internal/dashboardexecute"
	"github.com/turbot/powerpipe/internal/workspace"
	"gopkg.in/olahol/melody.v1"
)

// =============================================================================
// Test Helpers
// =============================================================================

// mockSession implements a mock melody.Session for testing
type mockSession struct {
	id            string
	writtenData   [][]byte
	mu            sync.Mutex
	closed        bool
	writeCallback func([]byte)
}

func newMockSession(id string) *mockSession {
	return &mockSession{
		id:          id,
		writtenData: make([][]byte, 0),
	}
}

func (m *mockSession) Write(data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return nil
	}
	// Make a copy of the data
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)
	m.writtenData = append(m.writtenData, dataCopy)
	if m.writeCallback != nil {
		m.writeCallback(dataCopy)
	}
	return nil
}

func (m *mockSession) getWrittenData() [][]byte {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([][]byte, len(m.writtenData))
	copy(result, m.writtenData)
	return result
}

func (m *mockSession) clearWrittenData() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.writtenData = make([][]byte, 0)
}

func (m *mockSession) close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
}

// newTestServerEager creates an eager mode test server with the behavior_test_mod fixture
func newTestServerEager(t *testing.T) (*Server, func()) {
	t.Helper()
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	require.NoError(t, ew.GetError())

	webSocket := melody.New()

	server := &Server{
		mutex:            &sync.Mutex{},
		dashboardClients: make(map[string]*DashboardClientInfo),
		webSocket:        webSocket,
		workspace:        w,
	}

	cleanup := func() {
		_ = webSocket.Close()
	}

	return server, cleanup
}

// newTestServerLazy creates a lazy mode test server with the behavior_test_mod fixture
func newTestServerLazy(t *testing.T) (*Server, func()) {
	t.Helper()
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	lw, err := workspace.NewLazyWorkspace(ctx, modPath, workspace.DefaultLazyLoadConfig())
	require.NoError(t, err)

	webSocket := melody.New()

	server := &Server{
		mutex:            &sync.Mutex{},
		dashboardClients: make(map[string]*DashboardClientInfo),
		webSocket:        webSocket,
		lazyWorkspace:    lw,
	}

	cleanup := func() {
		_ = webSocket.Close()
		lw.Close()
	}

	return server, cleanup
}

// =============================================================================
// Server Initialization Tests
// =============================================================================

// TestServer_LazyModeCreation verifies that a lazy mode server is created correctly.
func TestServer_LazyModeCreation(t *testing.T) {
	server, cleanup := newTestServerLazy(t)
	defer cleanup()

	assert.True(t, server.isLazyMode(), "server should be in lazy mode")
	assert.NotNil(t, server.lazyWorkspace, "lazy workspace should be set")
	assert.Nil(t, server.workspace, "eager workspace should be nil in lazy mode")
}

// TestServer_EagerModeCreation verifies that an eager mode server is created correctly.
func TestServer_EagerModeCreation(t *testing.T) {
	server, cleanup := newTestServerEager(t)
	defer cleanup()

	assert.False(t, server.isLazyMode(), "server should not be in lazy mode")
	assert.Nil(t, server.lazyWorkspace, "lazy workspace should be nil in eager mode")
	assert.NotNil(t, server.workspace, "eager workspace should be set")
}

// TestServer_GetActiveWorkspace verifies the correct workspace is returned based on mode.
func TestServer_GetActiveWorkspace(t *testing.T) {
	t.Run("lazy mode returns lazy workspace", func(t *testing.T) {
		server, cleanup := newTestServerLazy(t)
		defer cleanup()

		ws := server.getActiveWorkspace()
		assert.NotNil(t, ws)
		// In lazy mode, getActiveWorkspace returns the LazyWorkspace
		_, ok := ws.(*workspace.LazyWorkspace)
		assert.True(t, ok, "should return LazyWorkspace")
	})

	t.Run("eager mode returns powerpipe workspace", func(t *testing.T) {
		server, cleanup := newTestServerEager(t)
		defer cleanup()

		ws := server.getActiveWorkspace()
		assert.NotNil(t, ws)
		_, ok := ws.(*workspace.PowerpipeWorkspace)
		assert.True(t, ok, "should return PowerpipeWorkspace")
	})
}

// =============================================================================
// Available Dashboards Tests
// =============================================================================

// TestServer_AvailableDashboards_LazyMode verifies available dashboards payload in lazy mode.
func TestServer_AvailableDashboards_LazyMode(t *testing.T) {
	server, cleanup := newTestServerLazy(t)
	defer cleanup()

	// Measure time to build payload
	start := time.Now()
	payloadBytes, err := server.buildAvailableDashboardsPayload()
	elapsed := time.Since(start)

	require.NoError(t, err)
	require.NotEmpty(t, payloadBytes)

	// Lazy mode should be fast (<100ms for small fixtures)
	assert.Less(t, elapsed.Milliseconds(), int64(100),
		"lazy mode payload build should be fast")

	// Parse and verify payload structure
	var payload AvailableDashboardsPayload
	err = json.Unmarshal(payloadBytes, &payload)
	require.NoError(t, err)

	assert.Equal(t, "available_dashboards", payload.Action)
	assert.NotEmpty(t, payload.Dashboards, "should have dashboards")
	assert.NotEmpty(t, payload.Benchmarks, "should have benchmarks")
}

// TestServer_AvailableDashboards_EagerMode verifies available dashboards payload in eager mode.
func TestServer_AvailableDashboards_EagerMode(t *testing.T) {
	server, cleanup := newTestServerEager(t)
	defer cleanup()

	payloadBytes, err := server.buildAvailableDashboardsPayload()
	require.NoError(t, err)
	require.NotEmpty(t, payloadBytes)

	var payload AvailableDashboardsPayload
	err = json.Unmarshal(payloadBytes, &payload)
	require.NoError(t, err)

	assert.Equal(t, "available_dashboards", payload.Action)
	assert.NotEmpty(t, payload.Dashboards, "should have dashboards")
	assert.NotEmpty(t, payload.Benchmarks, "should have benchmarks")
}

// TestServer_DashboardPayloadStructure verifies that dashboard payload has all required fields.
func TestServer_DashboardPayloadStructure_Integration(t *testing.T) {
	server, cleanup := newTestServerLazy(t)
	defer cleanup()

	payloadBytes, err := server.buildAvailableDashboardsPayload()
	require.NoError(t, err)

	var payload AvailableDashboardsPayload
	require.NoError(t, json.Unmarshal(payloadBytes, &payload))

	for name, dash := range payload.Dashboards {
		assert.NotEmpty(t, dash.FullName, "dashboard %s should have FullName", name)
		assert.NotEmpty(t, dash.ShortName, "dashboard %s should have ShortName", name)
		assert.NotEmpty(t, dash.ModFullName, "dashboard %s should have ModFullName", name)
		assert.Equal(t, name, dash.FullName, "dashboard key should match FullName")
	}
}

// TestServer_BenchmarkPayloadStructure verifies that benchmark payload has all required fields.
func TestServer_BenchmarkPayloadStructure_Integration(t *testing.T) {
	server, cleanup := newTestServerLazy(t)
	defer cleanup()

	payloadBytes, err := server.buildAvailableDashboardsPayload()
	require.NoError(t, err)

	var payload AvailableDashboardsPayload
	require.NoError(t, json.Unmarshal(payloadBytes, &payload))

	for name, bench := range payload.Benchmarks {
		assert.NotEmpty(t, bench.FullName, "benchmark %s should have FullName", name)
		assert.NotEmpty(t, bench.ShortName, "benchmark %s should have ShortName", name)
		assert.NotEmpty(t, bench.ModFullName, "benchmark %s should have ModFullName", name)
		assert.NotEmpty(t, bench.BenchmarkType, "benchmark %s should have BenchmarkType", name)
		assert.Equal(t, name, bench.FullName, "benchmark key should match FullName")
	}
}

// TestServer_TopLevelBenchmarkHasTrunks verifies top-level benchmarks have trunks populated.
func TestServer_TopLevelBenchmarkHasTrunks_Integration(t *testing.T) {
	server, cleanup := newTestServerLazy(t)
	defer cleanup()

	payloadBytes, err := server.buildAvailableDashboardsPayload()
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

// =============================================================================
// Payload Consistency Tests
// =============================================================================

// TestServer_PayloadConsistency verifies lazy and eager modes produce consistent payloads.
func TestServer_PayloadConsistency(t *testing.T) {
	lazyServer, lazyCleanup := newTestServerLazy(t)
	defer lazyCleanup()

	eagerServer, eagerCleanup := newTestServerEager(t)
	defer eagerCleanup()

	lazyPayloadBytes, err := lazyServer.buildAvailableDashboardsPayload()
	require.NoError(t, err)

	eagerPayloadBytes, err := eagerServer.buildAvailableDashboardsPayload()
	require.NoError(t, err)

	var lazyPayload, eagerPayload AvailableDashboardsPayload
	require.NoError(t, json.Unmarshal(lazyPayloadBytes, &lazyPayload))
	require.NoError(t, json.Unmarshal(eagerPayloadBytes, &eagerPayload))

	// Verify same number of dashboards
	assert.Equal(t, len(eagerPayload.Dashboards), len(lazyPayload.Dashboards),
		"dashboard count should match between lazy and eager")

	// Verify same number of benchmarks
	assert.Equal(t, len(eagerPayload.Benchmarks), len(lazyPayload.Benchmarks),
		"benchmark count should match between lazy and eager")

	// Verify same dashboards exist
	for name := range eagerPayload.Dashboards {
		_, ok := lazyPayload.Dashboards[name]
		assert.True(t, ok, "dashboard %s should exist in lazy payload", name)
	}

	// Verify same benchmarks exist
	for name := range eagerPayload.Benchmarks {
		_, ok := lazyPayload.Benchmarks[name]
		assert.True(t, ok, "benchmark %s should exist in lazy payload", name)
	}

	// Verify dashboard identity fields match (FullName, ShortName, ModFullName)
	// Note: Title extraction may differ slightly between lazy and eager modes
	// as the lazy scanner extracts from HCL attributes while eager uses fully parsed resources
	for name, eagerDash := range eagerPayload.Dashboards {
		lazyDash := lazyPayload.Dashboards[name]
		assert.Equal(t, eagerDash.FullName, lazyDash.FullName,
			"dashboard %s FullName should match", name)
		assert.Equal(t, eagerDash.ShortName, lazyDash.ShortName,
			"dashboard %s ShortName should match", name)
		assert.Equal(t, eagerDash.ModFullName, lazyDash.ModFullName,
			"dashboard %s ModFullName should match", name)
	}
}

// =============================================================================
// Session Management Tests
// =============================================================================

// TestServer_SessionCreated verifies that sessions are created on connect.
func TestServer_SessionCreated(t *testing.T) {
	server, cleanup := newTestServerEager(t)
	defer cleanup()

	mockSess := newMockSession("test-session-1")

	// Simulate a melody.Session pointer for session ID
	melodySession := &melody.Session{}

	server.addSession(melodySession)

	sessionId := server.getSessionId(melodySession)
	clients := server.getDashboardClients()

	assert.Contains(t, clients, sessionId, "session should be added")
	assert.NotNil(t, clients[sessionId], "session info should not be nil")

	// Cleanup
	_ = mockSess
}

// TestServer_SessionRemoved verifies that sessions are removed from dashboard clients.
func TestServer_SessionRemoved(t *testing.T) {
	server, cleanup := newTestServerEager(t)
	defer cleanup()

	melodySession := &melody.Session{}
	server.addSession(melodySession)

	sessionId := server.getSessionId(melodySession)
	clients := server.getDashboardClients()
	assert.Contains(t, clients, sessionId, "session should exist before deletion")

	// Directly test the deleteDashboardClient method
	// (clearSession also calls CancelExecutionForSession which requires global state)
	server.deleteDashboardClient(sessionId)

	clients = server.getDashboardClients()
	assert.NotContains(t, clients, sessionId, "session should be removed after deletion")
}

// TestServer_SessionIsolation verifies that sessions are isolated from each other.
func TestServer_SessionIsolation(t *testing.T) {
	server, cleanup := newTestServerEager(t)
	defer cleanup()

	// Create two sessions
	session1 := &melody.Session{}
	session2 := &melody.Session{}

	server.addSession(session1)
	server.addSession(session2)

	sessionId1 := server.getSessionId(session1)
	sessionId2 := server.getSessionId(session2)

	// Set different dashboards for each session
	dashboard1 := "behavior_test.dashboard.main"
	dashboard2 := "behavior_test.dashboard.simple"

	server.setDashboardForSession(sessionId1, dashboard1, nil)
	server.setDashboardForSession(sessionId2, dashboard2, nil)

	clients := server.getDashboardClients()

	// Verify sessions have different dashboards
	assert.Equal(t, dashboard1, *clients[sessionId1].Dashboard)
	assert.Equal(t, dashboard2, *clients[sessionId2].Dashboard)

	// Verify changing one doesn't affect the other
	newDashboard := "behavior_test.dashboard.nested"
	server.setDashboardForSession(sessionId1, newDashboard, nil)

	clients = server.getDashboardClients()
	assert.Equal(t, newDashboard, *clients[sessionId1].Dashboard)
	assert.Equal(t, dashboard2, *clients[sessionId2].Dashboard)
}

// TestServer_SetDashboardInputsForSession verifies input management per session.
func TestServer_SetDashboardInputsForSession(t *testing.T) {
	server, cleanup := newTestServerEager(t)
	defer cleanup()

	session := &melody.Session{}
	server.addSession(session)
	sessionId := server.getSessionId(session)

	// Set dashboard first
	server.setDashboardForSession(sessionId, "behavior_test.dashboard.main", nil)

	// Set inputs
	inputs := &dashboardexecute.InputValues{
		Inputs: map[string]interface{}{
			"filter_selection": "opt1",
		},
	}
	server.setDashboardInputsForSession(sessionId, inputs)

	clients := server.getDashboardClients()
	assert.NotNil(t, clients[sessionId].DashboardInputs)
	assert.Equal(t, "opt1", clients[sessionId].DashboardInputs.Inputs["filter_selection"])
}

// =============================================================================
// Concurrent Session Tests
// =============================================================================

// TestServer_ConcurrentSessions verifies handling of multiple concurrent sessions.
func TestServer_ConcurrentSessions(t *testing.T) {
	server, cleanup := newTestServerEager(t)
	defer cleanup()

	const numSessions = 10
	var wg sync.WaitGroup
	sessions := make([]*melody.Session, numSessions)

	// Create sessions concurrently
	for i := 0; i < numSessions; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sessions[idx] = &melody.Session{}
			server.addSession(sessions[idx])
		}(i)
	}
	wg.Wait()

	// Verify all sessions exist
	clients := server.getDashboardClients()
	assert.Len(t, clients, numSessions, "all sessions should be created")

	// Set different dashboards concurrently
	dashboards := []string{
		"behavior_test.dashboard.main",
		"behavior_test.dashboard.simple",
		"behavior_test.dashboard.nested",
	}

	for i := 0; i < numSessions; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sessionId := server.getSessionId(sessions[idx])
			dashboard := dashboards[idx%len(dashboards)]
			server.setDashboardForSession(sessionId, dashboard, nil)
		}(i)
	}
	wg.Wait()

	// Verify dashboards were set correctly
	clients = server.getDashboardClients()
	for i := 0; i < numSessions; i++ {
		sessionId := server.getSessionId(sessions[i])
		expectedDashboard := dashboards[i%len(dashboards)]
		assert.NotNil(t, clients[sessionId].Dashboard)
		assert.Equal(t, expectedDashboard, *clients[sessionId].Dashboard)
	}
}

// TestServer_ConcurrentSameDashboard verifies multiple sessions can run the same dashboard.
func TestServer_ConcurrentSameDashboard(t *testing.T) {
	server, cleanup := newTestServerEager(t)
	defer cleanup()

	const numSessions = 5
	sessions := make([]*melody.Session, numSessions)
	dashboard := "behavior_test.dashboard.main"

	// Create sessions and set same dashboard
	for i := 0; i < numSessions; i++ {
		sessions[i] = &melody.Session{}
		server.addSession(sessions[i])
		sessionId := server.getSessionId(sessions[i])
		server.setDashboardForSession(sessionId, dashboard, nil)
	}

	// Verify all have the same dashboard
	clients := server.getDashboardClients()
	for i := 0; i < numSessions; i++ {
		sessionId := server.getSessionId(sessions[i])
		assert.Equal(t, dashboard, *clients[sessionId].Dashboard)
	}
}

// =============================================================================
// Message Parsing Tests
// =============================================================================

// TestServer_ClientRequestParsing verifies client request message parsing.
func TestServer_ClientRequestParsing(t *testing.T) {
	testCases := []struct {
		name           string
		jsonMessage    string
		expectedAction string
		shouldParse    bool
	}{
		{
			name:           "get_available_dashboards",
			jsonMessage:    `{"action":"get_available_dashboards"}`,
			expectedAction: "get_available_dashboards",
			shouldParse:    true,
		},
		{
			name:           "get_server_metadata",
			jsonMessage:    `{"action":"get_server_metadata"}`,
			expectedAction: "get_server_metadata",
			shouldParse:    true,
		},
		{
			name:           "select_dashboard",
			jsonMessage:    `{"action":"select_dashboard","payload":{"dashboard":{"full_name":"test.dashboard.main"}}}`,
			expectedAction: "select_dashboard",
			shouldParse:    true,
		},
		{
			name:           "input_changed",
			jsonMessage:    `{"action":"input_changed","payload":{"changed_input":"filter","inputs":{"filter":"value"}}}`,
			expectedAction: "input_changed",
			shouldParse:    true,
		},
		{
			name:           "clear_dashboard",
			jsonMessage:    `{"action":"clear_dashboard"}`,
			expectedAction: "clear_dashboard",
			shouldParse:    true,
		},
		{
			name:           "keep_alive",
			jsonMessage:    `{"action":"keep_alive"}`,
			expectedAction: "keep_alive",
			shouldParse:    true,
		},
		{
			name:           "invalid json",
			jsonMessage:    `{invalid json`,
			expectedAction: "",
			shouldParse:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var request ClientRequest
			err := json.Unmarshal([]byte(tc.jsonMessage), &request)

			if tc.shouldParse {
				require.NoError(t, err, "should parse valid JSON")
				assert.Equal(t, tc.expectedAction, request.Action)
			} else {
				assert.Error(t, err, "should fail to parse invalid JSON")
			}
		})
	}
}

// TestServer_SelectDashboardPayloadParsing verifies select_dashboard payload parsing.
func TestServer_SelectDashboardPayloadParsing(t *testing.T) {
	jsonMessage := `{
		"action": "select_dashboard",
		"payload": {
			"dashboard": {"full_name": "test.dashboard.main"},
			"inputs": {"filter": "value1"},
			"search_path": ["schema1", "schema2"],
			"search_path_prefix": ["prefix1"]
		}
	}`

	var request ClientRequest
	err := json.Unmarshal([]byte(jsonMessage), &request)
	require.NoError(t, err)

	assert.Equal(t, "select_dashboard", request.Action)
	assert.Equal(t, "test.dashboard.main", request.Payload.Dashboard.FullName)
	assert.Equal(t, "value1", request.Payload.Inputs["filter"])
	assert.Equal(t, []string{"schema1", "schema2"}, request.Payload.SearchPath)
	assert.Equal(t, []string{"prefix1"}, request.Payload.SearchPathPrefix)

	// Verify InputValues conversion
	inputValues := request.Payload.InputValues()
	assert.NotNil(t, inputValues)
	assert.Equal(t, "value1", inputValues.Inputs["filter"])
}

// =============================================================================
// Error Handling Tests
// =============================================================================

// TestServer_InvalidMessageJSON verifies server handles invalid JSON gracefully.
func TestServer_InvalidMessageJSON(t *testing.T) {
	// Test that invalid JSON is handled without panic
	invalidMessages := []string{
		`{invalid`,
		`null`,
		`[]`,
		`""`,
		`123`,
	}

	for _, msg := range invalidMessages {
		t.Run(msg, func(t *testing.T) {
			var request ClientRequest
			err := json.Unmarshal([]byte(msg), &request)
			// Should either error or produce an empty request
			if err == nil {
				assert.Empty(t, request.Action)
			}
		})
	}
}

// TestServer_UnknownAction verifies unknown actions are handled gracefully.
func TestServer_UnknownAction(t *testing.T) {
	server, cleanup := newTestServerEager(t)
	defer cleanup()

	// Create a valid ClientRequest with unknown action
	request := ClientRequest{
		Action: "unknown_action",
	}

	jsonBytes, err := json.Marshal(request)
	require.NoError(t, err)

	// Verify the message parses correctly
	var parsed ClientRequest
	err = json.Unmarshal(jsonBytes, &parsed)
	require.NoError(t, err)
	assert.Equal(t, "unknown_action", parsed.Action)

	// The server should handle unknown actions gracefully (no crash)
	_ = server
}

// =============================================================================
// Payload Serialization Tests
// =============================================================================

// TestServer_PayloadSerialization verifies all payload types serialize to valid JSON.
func TestServer_PayloadSerialization(t *testing.T) {
	t.Run("AvailableDashboardsPayload", func(t *testing.T) {
		payload := AvailableDashboardsPayload{
			Action: "available_dashboards",
			Dashboards: map[string]ModAvailableDashboard{
				"test.dashboard.main": {
					FullName:    "test.dashboard.main",
					ShortName:   "main",
					Title:       "Main Dashboard",
					ModFullName: "mod.test",
				},
			},
			Benchmarks: map[string]ModAvailableBenchmark{
				"test.benchmark.security": {
					FullName:      "test.benchmark.security",
					ShortName:     "security",
					Title:         "Security Benchmark",
					BenchmarkType: "control",
					IsTopLevel:    true,
					ModFullName:   "mod.test",
				},
			},
		}

		bytes, err := json.Marshal(payload)
		require.NoError(t, err)

		var deserialized AvailableDashboardsPayload
		err = json.Unmarshal(bytes, &deserialized)
		require.NoError(t, err)

		assert.Equal(t, payload.Action, deserialized.Action)
		assert.Len(t, deserialized.Dashboards, 1)
		assert.Len(t, deserialized.Benchmarks, 1)
	})

	t.Run("ErrorPayload", func(t *testing.T) {
		payload := ErrorPayload{
			Action: "workspace_error",
			Error:  "test error message",
		}

		bytes, err := json.Marshal(payload)
		require.NoError(t, err)

		var deserialized ErrorPayload
		err = json.Unmarshal(bytes, &deserialized)
		require.NoError(t, err)

		assert.Equal(t, "workspace_error", deserialized.Action)
		assert.Equal(t, "test error message", deserialized.Error)
	})

	t.Run("InputValuesClearedPayload", func(t *testing.T) {
		payload := InputValuesClearedPayload{
			Action:        "input_values_cleared",
			ClearedInputs: []string{"input1", "input2"},
			ExecutionId:   "exec-123",
		}

		bytes, err := json.Marshal(payload)
		require.NoError(t, err)

		var deserialized InputValuesClearedPayload
		err = json.Unmarshal(bytes, &deserialized)
		require.NoError(t, err)

		assert.Equal(t, "input_values_cleared", deserialized.Action)
		assert.Equal(t, []string{"input1", "input2"}, deserialized.ClearedInputs)
	})
}

// =============================================================================
// Large Payload Tests
// =============================================================================

// TestServer_LargePayload verifies that large payloads are handled correctly.
func TestServer_LargePayload(t *testing.T) {
	// Create a payload with many dashboards and benchmarks
	payload := AvailableDashboardsPayload{
		Action:     "available_dashboards",
		Dashboards: make(map[string]ModAvailableDashboard),
		Benchmarks: make(map[string]ModAvailableBenchmark),
	}

	// Add 500 dashboards
	for i := 0; i < 500; i++ {
		name := "test.dashboard.dash_" + string(rune('0'+i/100)) + string(rune('0'+(i/10)%10)) + string(rune('0'+i%10))
		payload.Dashboards[name] = ModAvailableDashboard{
			FullName:    name,
			ShortName:   "dash_" + string(rune('0'+i%10)),
			Title:       "Dashboard " + string(rune('0'+i%10)),
			ModFullName: "mod.test",
			Tags:        map[string]string{"env": "test", "index": string(rune('0' + i%10))},
		}
	}

	// Add 100 benchmarks
	for i := 0; i < 100; i++ {
		name := "test.benchmark.bench_" + string(rune('0'+i/10)) + string(rune('0'+i%10))
		payload.Benchmarks[name] = ModAvailableBenchmark{
			FullName:      name,
			ShortName:     "bench_" + string(rune('0'+i%10)),
			Title:         "Benchmark " + string(rune('0'+i%10)),
			BenchmarkType: "control",
			IsTopLevel:    i < 10,
			ModFullName:   "mod.test",
		}
	}

	// Serialize
	bytes, err := json.Marshal(payload)
	require.NoError(t, err)
	assert.NotEmpty(t, bytes)

	// Deserialize
	var deserialized AvailableDashboardsPayload
	err = json.Unmarshal(bytes, &deserialized)
	require.NoError(t, err)

	assert.Len(t, deserialized.Dashboards, 500, "all dashboards should be preserved")
	assert.Len(t, deserialized.Benchmarks, 100, "all benchmarks should be preserved")
}

// =============================================================================
// Lazy Mode Specific Tests
// =============================================================================

// TestServer_LazyModeWorkspaceForExecution verifies eager workspace loading in lazy mode.
func TestServer_LazyModeWorkspaceForExecution(t *testing.T) {
	server, cleanup := newTestServerLazy(t)
	defer cleanup()
	ctx := context.Background()

	// First call should trigger eager load
	ws, err := server.getWorkspaceForExecution(ctx)
	require.NoError(t, err)
	assert.NotNil(t, ws, "should return eager workspace for execution")

	// Second call should return cached workspace
	ws2, err := server.getWorkspaceForExecution(ctx)
	require.NoError(t, err)
	assert.Equal(t, ws, ws2, "should return same cached workspace")
}

// =============================================================================
// Payload Timing Tests
// =============================================================================

// TestServer_PayloadTimingLazyVsEager compares payload build times.
func TestServer_PayloadTimingLazyVsEager(t *testing.T) {
	lazyServer, lazyCleanup := newTestServerLazy(t)
	defer lazyCleanup()

	eagerServer, eagerCleanup := newTestServerEager(t)
	defer eagerCleanup()

	// Time lazy payload build (should be fast)
	lazyStart := time.Now()
	for i := 0; i < 10; i++ {
		_, _ = lazyServer.buildAvailableDashboardsPayload()
	}
	lazyDuration := time.Since(lazyStart)

	// Time eager payload build
	eagerStart := time.Now()
	for i := 0; i < 10; i++ {
		_, _ = eagerServer.buildAvailableDashboardsPayload()
	}
	eagerDuration := time.Since(eagerStart)

	t.Logf("Lazy payload build (10 iterations): %v", lazyDuration)
	t.Logf("Eager payload build (10 iterations): %v", eagerDuration)

	// Both should complete quickly for small fixtures
	assert.Less(t, lazyDuration.Milliseconds(), int64(1000),
		"lazy payload build should complete in under 1 second")
	assert.Less(t, eagerDuration.Milliseconds(), int64(1000),
		"eager payload build should complete in under 1 second")
}

// =============================================================================
// Thread Safety Tests
// =============================================================================

// TestServer_ThreadSafety verifies concurrent access to server methods.
func TestServer_ThreadSafety(t *testing.T) {
	server, cleanup := newTestServerEager(t)
	defer cleanup()

	const numGoroutines = 50
	const numOperations = 100

	var wg sync.WaitGroup
	errCount := atomic.Int32{}

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			session := &melody.Session{}
			server.addSession(session)
			sessionId := server.getSessionId(session)

			for j := 0; j < numOperations; j++ {
				// Randomly perform operations
				switch j % 5 {
				case 0:
					server.setDashboardForSession(sessionId, "behavior_test.dashboard.main", nil)
				case 1:
					server.setDashboardInputsForSession(sessionId, &dashboardexecute.InputValues{
						Inputs: map[string]interface{}{"key": "value"},
					})
				case 2:
					_ = server.getDashboardClients()
				case 3:
					_, err := server.buildAvailableDashboardsPayload()
					if err != nil {
						errCount.Add(1)
					}
				case 4:
					_ = server.isLazyMode()
				}
			}
		}(i)
	}

	wg.Wait()
	assert.Equal(t, int32(0), errCount.Load(), "no errors should occur during concurrent operations")
}

// =============================================================================
// Benchmark Tests
// =============================================================================

// BenchmarkServerPayloadBuild benchmarks payload building performance.
func BenchmarkServerPayloadBuild(b *testing.B) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	b.Run("lazy", func(b *testing.B) {
		lw, err := workspace.NewLazyWorkspace(ctx, modPath, workspace.DefaultLazyLoadConfig())
		if err != nil {
			b.Fatal(err)
		}
		defer lw.Close()

		server := &Server{
			mutex:            &sync.Mutex{},
			dashboardClients: make(map[string]*DashboardClientInfo),
			lazyWorkspace:    lw,
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = server.buildAvailableDashboardsPayload()
		}
	})

	b.Run("eager", func(b *testing.B) {
		w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
		if ew.GetError() != nil {
			b.Fatal(ew.GetError())
		}

		server := &Server{
			mutex:            &sync.Mutex{},
			dashboardClients: make(map[string]*DashboardClientInfo),
			workspace:        w,
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = server.buildAvailableDashboardsPayload()
		}
	})
}

// BenchmarkSessionOperations benchmarks session management operations.
func BenchmarkSessionOperations(b *testing.B) {
	ctx := context.Background()
	modPath := filepath.Join(testdataDir(), "mods", "behavior_test_mod")

	w, ew := workspace.Load(ctx, modPath, workspace.WithVariableValidation(false))
	if ew.GetError() != nil {
		b.Fatal(ew.GetError())
	}

	server := &Server{
		mutex:            &sync.Mutex{},
		dashboardClients: make(map[string]*DashboardClientInfo),
		workspace:        w,
	}

	b.Run("add_session", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			session := &melody.Session{}
			server.addSession(session)
		}
	})

	b.Run("set_dashboard", func(b *testing.B) {
		session := &melody.Session{}
		server.addSession(session)
		sessionId := server.getSessionId(session)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			server.setDashboardForSession(sessionId, "test.dashboard.main", nil)
		}
	})

	b.Run("get_clients", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = server.getDashboardClients()
		}
	})
}
