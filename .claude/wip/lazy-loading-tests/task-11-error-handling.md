# Task 11: Error Handling & Recovery Tests

## Objective

Write tests that verify correct error handling, graceful degradation, and recovery from failure scenarios in lazy loading.

## Context

- Errors can occur at many points: scanning, parsing, loading, execution
- sync.Once traps errors - important to handle correctly
- Partial failures need graceful handling
- User-facing errors should be clear and actionable
- System should recover gracefully where possible

## Dependencies

- Tasks 3-9 (unit/integration tests should pass)
- Task 6 (error condition fixtures)

## Acceptance Criteria

- [ ] Add tests to `internal/workspace/error_handling_test.go`
- [ ] Test error scenarios for each major component
- [ ] Verify error messages are clear and helpful
- [ ] Test recovery after errors
- [ ] Ensure no panics or crashes

## Test Cases to Implement

### Scanner Errors
```go
// Test: Invalid HCL syntax
func TestError_InvalidHCLSyntax(t *testing.T)
// File with syntax errors
// Scanner should skip or report
// Not crash

// Test: File read permission denied
func TestError_FilePermissionDenied(t *testing.T)
// File with no read permissions
// Graceful error handling

// Test: File disappears during scan
func TestError_FileDisappearsDuringScan(t *testing.T)
// File deleted mid-scan
// Handle gracefully

// Test: Malformed UTF-8 content
func TestError_MalformedUTF8(t *testing.T)
// Invalid byte sequences
// Don't crash, skip or report
```

### Parser Errors
```go
// Test: Parse error for single resource
func TestError_SingleResourceParseError(t *testing.T)
// One resource invalid, others valid
// Load valid resources, report invalid

// Test: Missing required field
func TestError_MissingRequiredField(t *testing.T)
// Resource missing required field (e.g., sql in query)
// Clear error message

// Test: Invalid reference
func TestError_InvalidReference(t *testing.T)
// Reference to non-existent resource
// Error identifies what's missing

// Test: Type mismatch in reference
func TestError_TypeMismatchReference(t *testing.T)
// Benchmark child pointing to query
// Clear error about type mismatch
```

### Lazy Workspace Errors
```go
// Test: Index build failure
func TestError_IndexBuildFailure(t *testing.T)
// Workspace directory issues
// Clear error, workspace not created

// Test: Eager load failure cached
func TestError_EagerLoadFailureCached(t *testing.T)
// First GetWorkspaceForExecution fails
// Error cached in sync.Once
// Subsequent calls return same error

// Test: Partial resource load failure
func TestError_PartialResourceLoadFailure(t *testing.T)
// Some resources load, some fail
// Partial functionality available

// Test: Cache lookup after load failure
func TestError_CacheLookupAfterFailure(t *testing.T)
// Resource failed to load
// Cache lookup handles missing entry
```

### Dependency Resolution Errors
```go
// Test: Missing dependency error message
func TestError_MissingDependencyMessage(t *testing.T)
// Control needs query that doesn't exist
// Error: "control.X requires query.Y which was not found"

// Test: Circular dependency error message
func TestError_CircularDependencyMessage(t *testing.T)
// A → B → C → A cycle
// Error shows cycle: "circular dependency: A → B → C → A"

// Test: Transitive dependency failure
func TestError_TransitiveDependencyFailure(t *testing.T)
// A → B → C, C fails to load
// Error propagates back to A
```

### Execution Errors
```go
// Test: Dashboard not found
func TestError_DashboardNotFound(t *testing.T)
// Request non-existent dashboard
// Clear "dashboard not found" error

// Test: Execution timeout
func TestError_ExecutionTimeout(t *testing.T)
// Dashboard takes too long
// Context cancelled, error sent to client

// Test: SQL execution error
func TestError_SQLExecutionError(t *testing.T)
// Query has SQL error
// Error propagated to control
// Execution continues for other controls

// Test: Input validation error
func TestError_InputValidationError(t *testing.T)
// Invalid input value
// Clear validation message
```

### Server Errors
```go
// Test: Workspace error broadcast
func TestError_WorkspaceErrorBroadcast(t *testing.T)
// Workspace file watcher error
// All clients receive error

// Test: Session-specific error
func TestError_SessionSpecificError(t *testing.T)
// Error in one session
// Other sessions unaffected

// Test: Invalid client message
func TestError_InvalidClientMessage(t *testing.T)
// Malformed JSON from client
// Server doesn't crash
// Error logged

// Test: Unknown action type
func TestError_UnknownActionType(t *testing.T)
// Unknown action in client message
// Graceful handling
```

### Recovery Scenarios
```go
// Test: Recovery after file fixed
func TestError_RecoveryAfterFileFix(t *testing.T)
// File has error
// File fixed
// Subsequent load succeeds

// Test: Cache invalidation after error
func TestError_CacheInvalidationAfterError(t *testing.T)
// Resource failed to load
// Cache cleared
// Retry succeeds

// Test: Session recovery after execution error
func TestError_SessionRecoveryAfterError(t *testing.T)
// Execution fails
// Same session can run another dashboard

// Test: Server recovery after many errors
func TestError_ServerRecoveryAfterManyErrors(t *testing.T)
// Multiple consecutive errors
// Server remains functional
```

### Error Message Quality
```go
// Test: Error includes file location
func TestError_IncludesFileLocation(t *testing.T)
// Parse error
// Message includes file:line

// Test: Error includes resource name
func TestError_IncludesResourceName(t *testing.T)
// Resource-specific error
// Message identifies which resource

// Test: Error suggestions
func TestError_IncludesSuggestions(t *testing.T)
// Common mistakes
// Error suggests fix (if possible)

// Test: Error doesn't leak internals
func TestError_NoInternalLeakage(t *testing.T)
// User-facing errors
// No stack traces, internal paths
```

### Panic Prevention
```go
// Test: No panic on nil resource
func TestError_NoPanicOnNil(t *testing.T)
// GetResource returns nil
// Callers handle nil without panic

// Test: No panic on empty index
func TestError_NoPanicOnEmptyIndex(t *testing.T)
// Empty workspace
// Operations don't panic

// Test: No panic on concurrent error
func TestError_NoPanicOnConcurrentError(t *testing.T)
// Error during concurrent operations
// Clean error handling, no panic
```

### Error Aggregation
```go
// Test: Multiple errors collected
func TestError_MultipleErrorsCollected(t *testing.T)
// Several resources have errors
// All errors reported, not just first

// Test: Error count limits
func TestError_ErrorCountLimits(t *testing.T)
// Many errors (>100)
// Limited output to avoid spam
// Summary provided
```

## Test Fixture Requirements

Create error-inducing fixtures:
```
internal/testdata/mods/error-conditions/
├── invalid-syntax/
│   ├── mod.pp
│   └── broken.pp           # Syntax error
├── missing-refs/
│   ├── mod.pp
│   └── bad-refs.pp         # References to non-existent
├── circular-deps/
│   ├── mod.pp
│   └── circular.pp         # Circular benchmark hierarchy
└── partial-valid/
    ├── mod.pp
    ├── valid1.pp           # Valid resources
    └── invalid2.pp         # Invalid resources
```

## Implementation Notes

- Test both error occurrence and error message content
- Use require.Error and assert.ErrorContains
- Create helper to inject errors (mock file system)
- Test that errors are logged appropriately
- Verify cleanup happens after errors

## Output Files

- `internal/workspace/error_handling_test.go`
- `internal/resourceloader/error_test.go`
- `internal/dashboardserver/error_test.go`
