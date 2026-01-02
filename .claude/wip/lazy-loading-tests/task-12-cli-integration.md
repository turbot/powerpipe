# Task 12: CLI Command Integration Tests

## Objective

Write integration tests that verify lazy loading works correctly when invoked through CLI commands (server, dashboard, benchmark, etc.).

## Context

- CLI commands are the user entry point
- `--lazy-load` flag enables lazy loading
- Commands must work identically in eager and lazy modes
- Tests should verify end-to-end behavior
- Use subprocess testing pattern

## Dependencies

- All prior tasks (1-11) should be complete
- Files to test: `internal/cmd/*.go`

## Acceptance Criteria

- [ ] Add tests to `internal/cmd/integration_test.go`
- [ ] Test all major CLI commands with lazy loading
- [ ] Verify output matches between eager and lazy modes
- [ ] Test error cases through CLI
- [ ] Tests should be runnable in CI

## Test Cases to Implement

### Server Command
```go
// Test: Server starts with --lazy-load
func TestCLI_ServerStartsLazy(t *testing.T)
// powerpipe server --lazy-load
// Verify server starts
// Check endpoint responds

// Test: Server auto mode (default)
func TestCLI_ServerAutoMode(t *testing.T)
// powerpipe server (no flag)
// Verify mode selection based on workspace size

// Test: Server serves dashboards in lazy mode
func TestCLI_ServerServesDashboardsLazy(t *testing.T)
// Start server with --lazy-load
// HTTP request to get available dashboards
// Verify response contains dashboards

// Test: Server WebSocket works in lazy mode
func TestCLI_ServerWebSocketLazy(t *testing.T)
// Start server with --lazy-load
// Connect via WebSocket
// Send/receive messages
```

### Dashboard Command
```go
// Test: Dashboard run with lazy loading
func TestCLI_DashboardRunLazy(t *testing.T)
// powerpipe dashboard run dashboard.name --lazy-load
// Verify execution completes
// Output matches eager mode

// Test: Dashboard list with lazy loading
func TestCLI_DashboardListLazy(t *testing.T)
// powerpipe dashboard list --lazy-load
// Verify list output
// All dashboards shown

// Test: Dashboard show with lazy loading
func TestCLI_DashboardShowLazy(t *testing.T)
// powerpipe dashboard show dashboard.name --lazy-load
// Verify metadata shown
// Matches eager mode output
```

### Benchmark Command
```go
// Test: Benchmark run with lazy loading
func TestCLI_BenchmarkRunLazy(t *testing.T)
// powerpipe benchmark run benchmark.name --lazy-load
// Verify execution completes
// Controls run correctly

// Test: Benchmark list with lazy loading
func TestCLI_BenchmarkListLazy(t *testing.T)
// powerpipe benchmark list --lazy-load
// Verify list shows all benchmarks
// Hierarchy preserved

// Test: Benchmark with cross-mod controls
func TestCLI_BenchmarkCrossModLazy(t *testing.T)
// Benchmark includes dep mod controls
// All controls execute
```

### Control Command
```go
// Test: Control run with lazy loading
func TestCLI_ControlRunLazy(t *testing.T)
// powerpipe control run control.name --lazy-load
// Control executes
// Query resolved correctly

// Test: Control list with lazy loading
func TestCLI_ControlListLazy(t *testing.T)
// powerpipe control list --lazy-load
// All controls listed
```

### Query Command
```go
// Test: Query run with lazy loading
func TestCLI_QueryRunLazy(t *testing.T)
// powerpipe query run query.name --lazy-load
// Query executes
// Results returned

// Test: Query list with lazy loading
func TestCLI_QueryListLazy(t *testing.T)
// powerpipe query list --lazy-load
// All queries listed
```

### Mode Comparison Tests
```go
// Test: Dashboard list matches between modes
func TestCLI_DashboardListModeComparison(t *testing.T)
// Run with --lazy-load
// Run without (eager)
// Output should match

// Test: Benchmark run output matches
func TestCLI_BenchmarkRunModeComparison(t *testing.T)
// Same benchmark in both modes
// Same results (modulo timing)

// Test: Server available dashboards match
func TestCLI_ServerAvailableDashboardsModeComparison(t *testing.T)
// Start server in each mode
// Compare available_dashboards payload
// Should match exactly
```

### Error Handling via CLI
```go
// Test: Invalid mod error message
func TestCLI_InvalidModError(t *testing.T)
// powerpipe server --lazy-load in invalid mod
// Clear error message
// Non-zero exit code

// Test: Non-existent resource error
func TestCLI_NonExistentResource(t *testing.T)
// powerpipe dashboard run nonexistent --lazy-load
// "dashboard not found" error

// Test: Execution error
func TestCLI_ExecutionError(t *testing.T)
// Dashboard with SQL error
// Error propagated to CLI output
```

### Flag Combinations
```go
// Test: --lazy-load with --mod-location
func TestCLI_LazyWithModLocation(t *testing.T)
// Specify different mod directory
// Lazy loading works

// Test: --lazy-load with --output
func TestCLI_LazyWithOutput(t *testing.T)
// Different output formats (json, csv)
// Lazy loading produces correct output

// Test: --lazy-load with --where
func TestCLI_LazyWithWhere(t *testing.T)
// Filter controls/benchmarks
// Filtering works in lazy mode

// Test: --lazy-load with --tag
func TestCLI_LazyWithTag(t *testing.T)
// Tag-based filtering
// Works in lazy mode
```

### Timing and Performance
```go
// Test: Server startup time in lazy mode
func TestCLI_LazyStartupTime(t *testing.T)
// Large workspace
// Measure startup time
// Should be faster than eager

// Test: First execution time in lazy mode
func TestCLI_LazyFirstExecutionTime(t *testing.T)
// First dashboard execution
// Includes eager load overhead
// Subsequent faster
```

### Environment Variables
```go
// Test: POWERPIPE_LAZY_LOAD env var
func TestCLI_LazyLoadEnvVar(t *testing.T)
// Set POWERPIPE_LAZY_LOAD=true
// Verify lazy loading enabled

// Test: Flag overrides env var
func TestCLI_FlagOverridesEnvVar(t *testing.T)
// Env var set one way
// Flag set opposite
// Flag wins
```

## Test Infrastructure

### Subprocess Testing Pattern
```go
func runPowerpipe(t *testing.T, args ...string) (stdout, stderr string, exitCode int) {
    // Build powerpipe if needed
    binary := buildPowerpipe(t)

    // Create command
    cmd := exec.Command(binary, args...)
    cmd.Dir = testModDirectory(t)

    // Capture output
    var outBuf, errBuf bytes.Buffer
    cmd.Stdout = &outBuf
    cmd.Stderr = &errBuf

    // Run
    err := cmd.Run()

    return outBuf.String(), errBuf.String(), cmd.ProcessState.ExitCode()
}

// Helper to compare outputs
func assertOutputMatches(t *testing.T, lazy, eager string) {
    // Normalize outputs (timestamps, etc.)
    lazyNorm := normalizeOutput(lazy)
    eagerNorm := normalizeOutput(eager)

    assert.Equal(t, lazyNorm, eagerNorm)
}
```

### Test Mod Setup
```go
func setupCLITestMod(t *testing.T, fixture string) string {
    // Copy fixture to temp directory
    tmpDir := t.TempDir()
    copyFixture(t, fixture, tmpDir)
    return tmpDir
}
```

### Server Testing
```go
func startTestServer(t *testing.T, lazy bool, port int) *exec.Cmd {
    args := []string{"server", "--port", strconv.Itoa(port)}
    if lazy {
        args = append(args, "--lazy-load")
    }

    cmd := exec.Command(powerpipeBinary, args...)
    cmd.Start()

    // Wait for server to be ready
    waitForServer(t, port)

    t.Cleanup(func() {
        cmd.Process.Kill()
    })

    return cmd
}
```

## Notes

- Build powerpipe binary once per test run
- Use unique ports for server tests
- Clean up processes reliably
- Consider timeout for all subprocess calls
- Test on all supported platforms if possible

## Output Files

- `internal/cmd/integration_test.go`
- `internal/cmd/test_helpers.go`
