# Task 7: WebSocket Server Integration Tests

## Objective

Write integration tests for the dashboard server's WebSocket handling, focusing on lazy mode behavior, message handling, and state management.

## Context

- Dashboard server uses Melody library for WebSocket
- Two modes: eager and lazy (hybrid)
- Must handle concurrent sessions
- State management critical for correct execution
- Testing functions directly, not through browser

## Dependencies

- Task 2 (test fixtures)
- Tasks 3-6 (unit tests should pass)
- Files to test: `internal/dashboardserver/server.go`, `internal/dashboardserver/payload.go`

## Acceptance Criteria

- [ ] Add tests to `internal/dashboardserver/server_integration_test.go`
- [ ] Test all WebSocket message types
- [ ] Verify lazy mode behavior
- [ ] Test session state management
- [ ] Test concurrent session handling

## Test Cases to Implement

### Server Initialization
```go
// Test: Lazy mode server creation
func TestServer_LazyModeCreation(t *testing.T)
// NewServerWithLazyWorkspace creates properly
// isLazyMode() returns true

// Test: Eager mode server creation
func TestServer_EagerModeCreation(t *testing.T)
// NewServer creates properly
// isLazyMode() returns false

// Test: Server starts and stops cleanly
func TestServer_StartStop(t *testing.T)
// Start server
// Verify listening
// Stop server
// Verify shutdown
```

### Message Handling: get_available_dashboards
```go
// Test: Available dashboards in lazy mode
func TestServer_AvailableDashboards_LazyMode(t *testing.T)
// Create lazy server
// Call handler for get_available_dashboards
// Verify payload built from index
// Fast response (<100ms)

// Test: Available dashboards in eager mode
func TestServer_AvailableDashboards_EagerMode(t *testing.T)
// Create eager server
// Call handler
// Verify full payload

// Test: Dashboard payload structure
func TestServer_DashboardPayloadStructure(t *testing.T)
// Verify all required fields present:
// full_name, short_name, title, tags, mod_full_name

// Test: Benchmark payload structure
func TestServer_BenchmarkPayloadStructure(t *testing.T)
// Verify: full_name, title, children, trunks, is_top_level
```

### Message Handling: select_dashboard
```go
// Test: Select dashboard triggers eager load in lazy mode
func TestServer_SelectDashboard_TriggersEagerLoad(t *testing.T)
// Lazy server
// select_dashboard message
// Verify eager workspace used

// Test: Select dashboard with inputs
func TestServer_SelectDashboard_WithInputs(t *testing.T)
// Include input values in request
// Verify inputs passed to executor

// Test: Select dashboard with search path
func TestServer_SelectDashboard_WithSearchPath(t *testing.T)
// Include search path override
// Verify passed correctly

// Test: Select non-existent dashboard
func TestServer_SelectDashboard_NotFound(t *testing.T)
// Request dashboard that doesn't exist
// Verify error response
```

### Message Handling: input_changed
```go
// Test: Input change during execution
func TestServer_InputChanged_DuringExecution(t *testing.T)
// Start dashboard execution
// Send input_changed
// Verify re-execution triggered

// Test: Input change clears dependent inputs
func TestServer_InputChanged_ClearsDependents(t *testing.T)
// Inputs with dependencies
// Change parent input
// Verify dependent inputs cleared

// Test: Input change event sent to client
func TestServer_InputChanged_EventSent(t *testing.T)
// Verify InputValuesCleared event
```

### Message Handling: Other Messages
```go
// Test: get_server_metadata
func TestServer_ServerMetadata(t *testing.T)
// Verify mod info, version, etc.

// Test: clear_dashboard
func TestServer_ClearDashboard(t *testing.T)
// Active execution
// clear_dashboard
// Execution cancelled

// Test: keep_alive
func TestServer_KeepAlive(t *testing.T)
// Verify session kept alive
```

### Session Management
```go
// Test: Session created on connect
func TestServer_SessionCreated(t *testing.T)
// Simulate connect
// Verify DashboardClientInfo created

// Test: Session removed on disconnect
func TestServer_SessionRemoved(t *testing.T)
// Connect then disconnect
// Verify session cleaned up

// Test: Execution cancelled on disconnect
func TestServer_ExecutionCancelledOnDisconnect(t *testing.T)
// Start execution
// Disconnect
// Verify execution cancelled

// Test: Session state isolation
func TestServer_SessionIsolation(t *testing.T)
// Two sessions
// Different dashboards
// No cross-contamination
```

### Concurrent Sessions
```go
// Test: Multiple concurrent sessions
func TestServer_ConcurrentSessions(t *testing.T)
// 10 concurrent sessions
// Each running different dashboard
// All complete correctly

// Test: Concurrent same dashboard
func TestServer_ConcurrentSameDashboard(t *testing.T)
// 5 sessions running same dashboard
// Independent executions

// Test: Session limit (if any)
func TestServer_SessionLimit(t *testing.T)
// Many sessions
// Verify behavior at limit
```

### Event Broadcasting
```go
// Test: Events sent to correct session
func TestServer_EventRoutedToSession(t *testing.T)
// Multiple sessions
// Event for one
// Only that session receives

// Test: Broadcast to all sessions
func TestServer_BroadcastEvent(t *testing.T)
// workspace_error should broadcast
// All sessions receive

// Test: Event ordering
func TestServer_EventOrdering(t *testing.T)
// execution_started before leaf_node_updated
// leaf_node_updated before execution_complete
```

### Error Handling
```go
// Test: Invalid message JSON
func TestServer_InvalidMessageJSON(t *testing.T)
// Send malformed JSON
// Server doesn't crash
// Error logged

// Test: Unknown action type
func TestServer_UnknownAction(t *testing.T)
// Send unknown action
// Graceful handling

// Test: Execution error
func TestServer_ExecutionError(t *testing.T)
// Dashboard with error
// ExecutionError event sent

// Test: Workspace error
func TestServer_WorkspaceError(t *testing.T)
// Simulate workspace error
// Broadcast to clients
```

### Payload Generation
```go
// Test: Lazy payload matches eager payload
func TestServer_PayloadConsistency(t *testing.T)
// Same workspace
// Build payload in lazy mode
// Build payload in eager mode
// Should match (modulo execution-specific fields)

// Test: Large payload handling
func TestServer_LargePayload(t *testing.T)
// 500 dashboards
// Payload generates correctly
// No truncation

// Test: Payload JSON serialization
func TestServer_PayloadSerialization(t *testing.T)
// All payload types serialize to valid JSON
// Can be deserialized
```

## Implementation Notes

- Don't use actual WebSocket connections - test handler functions directly
- Mock Melody session for session tests
- Use test workspace with fixtures
- Verify event types and payloads

## Test Infrastructure

Create helper functions:
```go
// Create test server with fixture
func newTestServer(t *testing.T, fixture string, lazy bool) *Server

// Create mock session
func newMockSession(t *testing.T) *mockSession

// Simulate message handling
func (s *Server) handleTestMessage(session, action, payload) (events []Event, err error)
```

## Output Files

- `internal/dashboardserver/server_integration_test.go`
- `internal/dashboardserver/server_test_helpers.go`
