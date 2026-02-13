package dashboardserver

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turbot/powerpipe/internal/dashboardevents"
	"github.com/turbot/powerpipe/internal/resourceindex"
	"gopkg.in/olahol/melody.v1"
)

// =============================================================================
// Server Error Handling Tests
// =============================================================================

func TestError_ServerNilWorkspace(t *testing.T) {
	// Server with nil workspaces should handle gracefully
	s := &Server{
		mutex:            &sync.Mutex{},
		dashboardClients: make(map[string]*DashboardClientInfo),
		workspace:        nil,
		lazyWorkspace:    nil,
	}

	// Mode detection should not panic
	assert.False(t, s.isLazyMode())

	// getActiveWorkspace returns nil for nil server
	// This would panic if used, but the check itself shouldn't
}

func TestError_ServerClientSessionManagement(t *testing.T) {
	// Test that client session operations are thread-safe
	s := &Server{
		mutex:            &sync.Mutex{},
		dashboardClients: make(map[string]*DashboardClientInfo),
	}

	// Add and remove clients concurrently
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sessionId := string(rune('a' + id%26))
			clientInfo := &DashboardClientInfo{}

			// Add
			s.addDashboardClient(sessionId, clientInfo)

			// Get
			clients := s.getDashboardClients()
			_ = clients // Just access, don't hold

			// Delete
			s.deleteDashboardClient(sessionId)
		}(i)
	}

	wg.Wait() // Should complete without race or panic
}

func TestError_WritePayloadToNonexistentSession(t *testing.T) {
	s := &Server{
		mutex:            &sync.Mutex{},
		dashboardClients: make(map[string]*DashboardClientInfo),
	}

	// Writing to non-existent session should not panic
	s.writePayloadToSession("nonexistent-session-id", []byte("test payload"))
}

// =============================================================================
// Event Handling Error Tests
// =============================================================================

func TestError_HandleDashboardEvent_WorkspaceError(t *testing.T) {
	webSocket := melody.New()
	s := &Server{
		mutex:            &sync.Mutex{},
		dashboardClients: make(map[string]*DashboardClientInfo),
		webSocket:        webSocket,
	}

	ctx := context.Background()

	// Create workspace error event
	wsError := &dashboardevents.WorkspaceError{
		Error: errors.New("test workspace error"),
	}

	// Should not panic when handling error event
	s.HandleDashboardEvent(ctx, wsError)
}

func TestError_HandleDashboardEvent_ExecutionError(t *testing.T) {
	webSocket := melody.New()
	s := &Server{
		mutex:            &sync.Mutex{},
		dashboardClients: make(map[string]*DashboardClientInfo),
		webSocket:        webSocket,
	}

	ctx := context.Background()

	// Create execution error event
	execError := &dashboardevents.ExecutionError{
		Session: "test-session",
		Error:   errors.New("test execution error"),
	}

	// Should not panic when handling error event
	s.HandleDashboardEvent(ctx, execError)
}

// =============================================================================
// Payload Build Error Tests
// =============================================================================

func TestError_BuildPayloadFromEmptyIndex(t *testing.T) {
	// Empty index should produce valid but empty payload
	idx := resourceindex.NewResourceIndex()
	idx.ModName = "empty_mod"
	idx.ModFullName = "mod.empty_mod"

	payload := idx.BuildAvailableDashboardsPayload()

	assert.Equal(t, "available_dashboards", payload.Action)
	assert.Empty(t, payload.Dashboards)
	assert.Empty(t, payload.Benchmarks)
}

func TestError_BuildPayloadWithMissingFields(t *testing.T) {
	// Index with entries missing some fields
	idx := resourceindex.NewResourceIndex()
	idx.ModName = "test_mod"
	idx.ModFullName = "mod.test_mod"

	// Entry with minimal fields
	idx.Add(&resourceindex.IndexEntry{
		Type:      "dashboard",
		Name:      "test_mod.dashboard.minimal",
		ShortName: "minimal",
		// Title is empty
		// Tags is nil
	})

	// Entry with empty name
	idx.Add(&resourceindex.IndexEntry{
		Type:      "benchmark",
		Name:      "test_mod.benchmark.no_title",
		ShortName: "no_title",
		// Title empty
	})

	// Should not panic
	payload := idx.BuildAvailableDashboardsPayload()
	assert.NotNil(t, payload)
}

func TestError_PayloadMarshalError(t *testing.T) {
	// Test that payload marshaling handles edge cases
	payload := AvailableDashboardsPayload{
		Action:     "available_dashboards",
		Dashboards: make(map[string]ModAvailableDashboard),
		Benchmarks: make(map[string]ModAvailableBenchmark),
	}

	// Add dashboard with empty/nil maps
	payload.Dashboards["test.dashboard.empty"] = ModAvailableDashboard{
		Title:     "",
		FullName:  "test.dashboard.empty",
		ShortName: "empty",
		Tags:      nil, // nil tags
	}

	// Should marshal without error
	bytes, err := json.Marshal(payload)
	require.NoError(t, err)
	assert.NotEmpty(t, bytes)

	// Should unmarshal back
	var result AvailableDashboardsPayload
	err = json.Unmarshal(bytes, &result)
	require.NoError(t, err)
}

// =============================================================================
// Client Request Error Tests
// =============================================================================

func TestError_InvalidClientMessage(t *testing.T) {
	// Test handling of malformed JSON
	invalidMessages := []string{
		"",
		"not json",
		"{invalid json}",
		`{"action": 123}`, // Wrong type
		`{"action": "unknown_action"}`,
	}

	for _, msg := range invalidMessages {
		var request ClientRequest
		err := json.Unmarshal([]byte(msg), &request)
		// Some will fail to unmarshal, which is expected
		// The important thing is none cause panic
		_ = err
	}
}

func TestError_ClientRequestMissingPayload(t *testing.T) {
	// Request with action but missing payload
	jsonMsg := `{"action": "select_dashboard"}`

	var request ClientRequest
	err := json.Unmarshal([]byte(jsonMsg), &request)
	require.NoError(t, err)
	assert.Equal(t, "select_dashboard", request.Action)
	// Payload will be zero value, which handlers must handle
}

// =============================================================================
// Concurrent Error Tests
// =============================================================================

func TestError_ConcurrentEventHandling(t *testing.T) {
	webSocket := melody.New()
	s := &Server{
		mutex:            &sync.Mutex{},
		dashboardClients: make(map[string]*DashboardClientInfo),
		webSocket:        webSocket,
	}

	ctx := context.Background()
	var wg sync.WaitGroup

	// Send many events concurrently
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Different event types
			switch id % 3 {
			case 0:
				s.HandleDashboardEvent(ctx, &dashboardevents.WorkspaceError{
					Error: errors.New("concurrent error"),
				})
			case 1:
				s.HandleDashboardEvent(ctx, &dashboardevents.ExecutionError{
					Session: "session-" + string(rune('a'+id%26)),
					Error:   errors.New("exec error"),
				})
			case 2:
				s.HandleDashboardEvent(ctx, &dashboardevents.InputValuesCleared{
					Session:       "session-" + string(rune('a'+id%26)),
					ClearedInputs: []string{"input1"},
				})
			}
		}(i)
	}

	wg.Wait() // Should complete without race or panic
}

// =============================================================================
// Payload Builder Error Tests
// =============================================================================

func TestError_BuildWorkspaceErrorPayload(t *testing.T) {
	// Test building error payload
	event := &dashboardevents.WorkspaceError{
		Error: errors.New("test error message"),
	}

	payload, err := buildWorkspaceErrorPayload(event)
	require.NoError(t, err)
	assert.NotNil(t, payload)

	// Verify payload contains error
	var result map[string]interface{}
	err = json.Unmarshal(payload, &result)
	require.NoError(t, err)

	assert.Equal(t, "workspace_error", result["action"])
}

func TestError_BuildExecutionErrorPayload(t *testing.T) {
	event := &dashboardevents.ExecutionError{
		Session: "test-session",
		Error:   errors.New("execution failed"),
	}

	payload, err := buildExecutionErrorPayload(event)
	require.NoError(t, err)
	assert.NotNil(t, payload)

	var result map[string]interface{}
	err = json.Unmarshal(payload, &result)
	require.NoError(t, err)

	assert.Equal(t, "execution_error", result["action"])
}

// =============================================================================
// Recovery Tests
// =============================================================================

func TestError_ServerRecoveryAfterErrors(t *testing.T) {
	webSocket := melody.New()
	s := &Server{
		mutex:            &sync.Mutex{},
		dashboardClients: make(map[string]*DashboardClientInfo),
		webSocket:        webSocket,
	}

	ctx := context.Background()

	// Generate many errors
	for i := 0; i < 100; i++ {
		s.HandleDashboardEvent(ctx, &dashboardevents.WorkspaceError{
			Error: errors.New("repeated error"),
		})
	}

	// Server should still be functional
	assert.NotNil(t, s.dashboardClients)
	assert.NotNil(t, s.mutex)

	// Should be able to add/remove clients
	s.addDashboardClient("recovery-test", &DashboardClientInfo{})
	clients := s.getDashboardClients()
	assert.Contains(t, clients, "recovery-test")
	s.deleteDashboardClient("recovery-test")
}

func TestError_SessionRecoveryAfterExecutionError(t *testing.T) {
	s := &Server{
		mutex:            &sync.Mutex{},
		dashboardClients: make(map[string]*DashboardClientInfo),
	}

	sessionId := "recovery-session"
	s.addDashboardClient(sessionId, &DashboardClientInfo{
		// Note: Session is nil here, but that's OK - we're testing the server's
		// client tracking, not actual message writing
	})

	// Verify session was added
	clients := s.getDashboardClients()
	_, exists := clients[sessionId]
	assert.True(t, exists, "Session should exist after being added")

	// Delete and re-add to simulate recovery
	s.deleteDashboardClient(sessionId)
	s.addDashboardClient(sessionId, &DashboardClientInfo{})

	// Session should still be valid
	clients = s.getDashboardClients()
	_, exists = clients[sessionId]
	assert.True(t, exists, "Session should exist after recovery")
}

// =============================================================================
// Panic Prevention Tests
// =============================================================================

func TestError_NoPanicOnNilWebSocket(t *testing.T) {
	s := &Server{
		mutex:            &sync.Mutex{},
		dashboardClients: make(map[string]*DashboardClientInfo),
		webSocket:        nil, // Nil WebSocket
	}

	// Operations that might use webSocket should be careful
	// At minimum, mode checks should not panic
	assert.False(t, s.isLazyMode())
}

func TestError_NoPanicOnEmptyClientMap(t *testing.T) {
	s := &Server{
		mutex:            &sync.Mutex{},
		dashboardClients: make(map[string]*DashboardClientInfo),
	}

	// Operations on empty client map should not panic
	clients := s.getDashboardClients()
	assert.Empty(t, clients)

	s.writePayloadToSession("nonexistent", []byte("test"))
	s.deleteDashboardClient("nonexistent")
}

// =============================================================================
// Error Message Quality Tests
// =============================================================================

func TestError_WorkspaceErrorMessagePreserved(t *testing.T) {
	specificMessage := "very specific error message for testing 12345"
	event := &dashboardevents.WorkspaceError{
		Error: errors.New(specificMessage),
	}

	payload, err := buildWorkspaceErrorPayload(event)
	require.NoError(t, err)

	// Error message should be in payload
	payloadStr := string(payload)
	assert.Contains(t, payloadStr, specificMessage,
		"Error message should be preserved in payload")
}

func TestError_ExecutionErrorMessagePreserved(t *testing.T) {
	specificMessage := "execution error with unique identifier 67890"
	event := &dashboardevents.ExecutionError{
		Session: "test-session",
		Error:   errors.New(specificMessage),
	}

	payload, err := buildExecutionErrorPayload(event)
	require.NoError(t, err)

	payloadStr := string(payload)
	assert.Contains(t, payloadStr, specificMessage,
		"Error message should be preserved in payload")
}

// =============================================================================
// Index Conversion Error Tests
// =============================================================================

func TestError_ConvertEmptyBenchmarkInfo(t *testing.T) {
	// Empty benchmark info
	emptyInfo := resourceindex.BenchmarkInfo{}

	// Should not panic
	result := convertIndexBenchmarkInfo(emptyInfo)
	assert.NotNil(t, result)
	assert.Empty(t, result.Title)
	assert.Empty(t, result.Children)
}

func TestError_ConvertDeeplyNestedBenchmarks(t *testing.T) {
	// Create deeply nested benchmark structure
	depth := 50
	var buildNested func(level int) resourceindex.BenchmarkInfo
	buildNested = func(level int) resourceindex.BenchmarkInfo {
		info := resourceindex.BenchmarkInfo{
			Title:         "Level " + string(rune('0'+level)),
			FullName:      "mod.benchmark.level_" + string(rune('0'+level)),
			ShortName:     "level_" + string(rune('0'+level)),
			BenchmarkType: "control",
		}
		if level < depth {
			info.Children = []resourceindex.BenchmarkInfo{buildNested(level + 1)}
		}
		return info
	}

	root := buildNested(0)

	// Should not panic on deep recursion
	result := convertIndexBenchmarkInfo(root)
	assert.NotNil(t, result)

	// Verify depth
	current := result
	count := 0
	for len(current.Children) > 0 {
		count++
		current = current.Children[0]
	}
	assert.Equal(t, depth, count, "All levels should be converted")
}

// =============================================================================
// Payload Structure Tests
// =============================================================================

func TestError_PayloadStructureOnError(t *testing.T) {
	// Verify error payloads have correct structure
	wsError := &dashboardevents.WorkspaceError{
		Error: errors.New("test"),
	}
	payload, _ := buildWorkspaceErrorPayload(wsError)

	var result map[string]interface{}
	err := json.Unmarshal(payload, &result)
	require.NoError(t, err)

	// Should have action field
	assert.Contains(t, result, "action")

	// Should have error field
	assert.Contains(t, result, "error")
}

func TestError_EmptyPayloadFields(t *testing.T) {
	// Create payload with empty fields
	payload := AvailableDashboardsPayload{
		Action:     "",
		Dashboards: nil,
		Benchmarks: nil,
	}

	// Should marshal without panic
	bytes, err := json.Marshal(payload)
	require.NoError(t, err)

	// Should unmarshal
	var result AvailableDashboardsPayload
	err = json.Unmarshal(bytes, &result)
	require.NoError(t, err)
}
