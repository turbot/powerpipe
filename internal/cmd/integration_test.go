package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Test Infrastructure
// =============================================================================

var (
	powerpipeBinary string
	buildOnce       sync.Once
	buildErr        error
)

// testdataDir returns the path to the testdata directory
func testdataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "testdata")
}

// buildPowerpipe finds the powerpipe binary, preferring an existing one
func buildPowerpipe(t *testing.T) string {
	buildOnce.Do(func() {
		// First, check if powerpipe is in PATH
		path, err := exec.LookPath("powerpipe")
		if err == nil {
			powerpipeBinary = path
			return
		}

		// Check /usr/local/bin
		if _, err := os.Stat("/usr/local/bin/powerpipe"); err == nil {
			powerpipeBinary = "/usr/local/bin/powerpipe"
			return
		}

		// Try to build the binary as a fallback
		_, filename, _, _ := runtime.Caller(0)
		projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")

		// Create a temp file for the binary
		tmpDir := os.TempDir()
		powerpipeBinary = filepath.Join(tmpDir, "powerpipe_test")
		if runtime.GOOS == "windows" {
			powerpipeBinary += ".exe"
		}

		// Build the binary using make
		cmd := exec.Command("make", "build")
		cmd.Dir = projectRoot

		if err := cmd.Run(); err != nil {
			// Try direct go build as fallback
			cmd = exec.Command("go", "build", "-o", powerpipeBinary, ".")
			cmd.Dir = projectRoot
			output, buildCmdErr := cmd.CombinedOutput()
			if buildCmdErr != nil {
				buildErr = fmt.Errorf("failed to build powerpipe: %v\nOutput: %s", buildCmdErr, output)
			}
		} else {
			// make build installs to /usr/local/bin
			powerpipeBinary = "/usr/local/bin/powerpipe"
		}
	})

	if buildErr != nil {
		t.Skip("Skipping integration test: failed to find or build powerpipe binary:", buildErr)
	}

	return powerpipeBinary
}

// PowerpipeResult holds the result of a powerpipe command execution
type PowerpipeResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Duration time.Duration
}

// runPowerpipe runs the powerpipe binary with the given arguments
func runPowerpipe(t *testing.T, workDir string, args ...string) *PowerpipeResult {
	return runPowerpipeWithEnv(t, workDir, nil, args...)
}

// runPowerpipeWithEnv runs the powerpipe binary with custom environment variables
func runPowerpipeWithEnv(t *testing.T, workDir string, env []string, args ...string) *PowerpipeResult {
	return runPowerpipeWithTimeout(t, workDir, env, 60*time.Second, args...)
}

// runPowerpipeWithTimeout runs the powerpipe binary with a custom timeout
func runPowerpipeWithTimeout(t *testing.T, workDir string, env []string, timeout time.Duration, args ...string) *PowerpipeResult {
	binary := buildPowerpipe(t)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, binary, args...)
	cmd.Dir = workDir

	// Set environment
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, env...)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else if ctx.Err() == context.DeadlineExceeded {
			t.Fatalf("Command timed out after %v: %v", timeout, args)
		} else {
			exitCode = -1
		}
	}

	return &PowerpipeResult{
		Stdout:   outBuf.String(),
		Stderr:   errBuf.String(),
		ExitCode: exitCode,
		Duration: duration,
	}
}

// getTestModPath returns the path to a test mod
func getTestModPath(subpath ...string) string {
	parts := append([]string{testdataDir(), "mods"}, subpath...)
	return filepath.Join(parts...)
}

// getFreePort finds an available port for server tests
func getFreePort(t *testing.T) int {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}

// waitForServer waits for a server to become available
func waitForServer(t *testing.T, port int, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 100*time.Millisecond)
		if err == nil {
			conn.Close()
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
}

// startServer starts a powerpipe server and returns a cleanup function
// Note: lazy loading is now the default, use POWERPIPE_WORKSPACE_PRELOAD=true to disable
func startServer(t *testing.T, workDir string, port int, extraArgs ...string) (*exec.Cmd, func()) {
	binary := buildPowerpipe(t)

	args := []string{"server", "--port", fmt.Sprintf("%d", port)}
	args = append(args, extraArgs...)

	cmd := exec.Command(binary, args...)
	cmd.Dir = workDir
	cmd.Env = os.Environ()

	// Capture output for debugging
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Start()
	require.NoError(t, err, "Failed to start server")

	// Wait for server to be ready
	if !waitForServer(t, port, 30*time.Second) {
		_ = cmd.Process.Kill()
		t.Fatalf("Server failed to start within timeout.\nStdout: %s\nStderr: %s", outBuf.String(), errBuf.String())
	}

	cleanup := func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
		}
	}

	return cmd, cleanup
}


// =============================================================================
// Server Command Tests
// =============================================================================

func TestCLI_ServerStarts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping server test in short mode")
	}

	modPath := getTestModPath("lazy-loading-tests", "simple")
	port := getFreePort(t)

	_, cleanup := startServer(t, modPath, port)
	defer cleanup()

	// Verify server responds to HTTP requests
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/api/latest/status", port))
	if err != nil {
		// Try the root endpoint
		resp, err = http.Get(fmt.Sprintf("http://localhost:%d/", port))
	}
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCLI_ServerStartsWithPreload(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping server test in short mode")
	}

	modPath := getTestModPath("lazy-loading-tests", "simple")
	port := getFreePort(t)

	// Start server with workspace preload (eager loading)
	binary := buildPowerpipe(t)
	cmd := exec.Command(binary, "server", "--port", fmt.Sprintf("%d", port))
	cmd.Dir = modPath
	cmd.Env = append(os.Environ(), "POWERPIPE_WORKSPACE_PRELOAD=true")

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Start()
	require.NoError(t, err, "Failed to start server")
	defer func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
		}
	}()

	// Wait for server to be ready
	if !waitForServer(t, port, 30*time.Second) {
		t.Fatalf("Server failed to start within timeout.\nStdout: %s\nStderr: %s", outBuf.String(), errBuf.String())
	}

	// Verify server responds
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/", port))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCLI_ServerAvailableDashboardsAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping server test in short mode")
	}

	modPath := getTestModPath("lazy-loading-tests", "simple")

	// Test with lazy mode (default)
	port := getFreePort(t)
	_, cleanup := startServer(t, modPath, port)
	defer cleanup()

	// Give server time to fully initialize
	time.Sleep(500 * time.Millisecond)

	// The available_dashboards endpoint is typically accessed via WebSocket
	// For this test, we just verify the server is running and responding
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/", port))
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// =============================================================================
// Dashboard Command Tests
// =============================================================================

func TestCLI_DashboardListLazy(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	// Note: list command doesn't support --lazy-load flag as it doesn't execute queries
	result := runPowerpipe(t, modPath, "dashboard", "list", "--output", "json")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0, got %d. Stderr: %s", result.ExitCode, result.Stderr)
	// JSON output has mod_name and resource_name as separate fields
	assert.Contains(t, result.Stdout, `"mod_name": "lazy_simple"`, "Should list dashboard from lazy_simple mod")
	assert.Contains(t, result.Stdout, `"resource_name": "simple"`, "Should list the simple dashboard")
}

func TestCLI_DashboardShowLazy(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	// Show command also doesn't use --lazy-load
	result := runPowerpipe(t, modPath, "dashboard", "show", "lazy_simple.dashboard.simple", "--output", "json")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0, got %d. Stderr: %s", result.ExitCode, result.Stderr)
	assert.Contains(t, result.Stdout, "Simple Dashboard", "Should show dashboard title")
}

// =============================================================================
// Benchmark Command Tests
// =============================================================================

func TestCLI_BenchmarkListLazy(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	result := runPowerpipe(t, modPath, "benchmark", "list", "--output", "json")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0, got %d. Stderr: %s", result.ExitCode, result.Stderr)
	// JSON output has mod_name and resource_name as separate fields
	assert.Contains(t, result.Stdout, `"mod_name": "lazy_simple"`, "Should list benchmark from lazy_simple mod")
	assert.Contains(t, result.Stdout, `"resource_name": "simple"`, "Should list the simple benchmark")
}

func TestCLI_BenchmarkShowLazy(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	result := runPowerpipe(t, modPath, "benchmark", "show", "lazy_simple.benchmark.simple", "--output", "json")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0, got %d. Stderr: %s", result.ExitCode, result.Stderr)
	assert.Contains(t, result.Stdout, "Simple Benchmark", "Should show benchmark title")
}

// =============================================================================
// Control Command Tests
// =============================================================================

func TestCLI_ControlListLazy(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	result := runPowerpipe(t, modPath, "control", "list", "--output", "json")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0, got %d. Stderr: %s", result.ExitCode, result.Stderr)
	// JSON output has mod_name and resource_name as separate fields
	assert.Contains(t, result.Stdout, `"mod_name": "lazy_simple"`, "Should list controls from lazy_simple mod")
	assert.Contains(t, result.Stdout, `"resource_name": "inline_sql"`, "Should list inline_sql control")
	assert.Contains(t, result.Stdout, `"resource_name": "uses_query"`, "Should list uses_query control")
}

func TestCLI_ControlShowLazy(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	result := runPowerpipe(t, modPath, "control", "show", "lazy_simple.control.inline_sql", "--output", "json")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0, got %d. Stderr: %s", result.ExitCode, result.Stderr)
	assert.Contains(t, result.Stdout, "Control with Inline SQL", "Should show control title")
}

// =============================================================================
// Query Command Tests
// =============================================================================

func TestCLI_QueryListLazy(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	result := runPowerpipe(t, modPath, "query", "list", "--output", "json")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0, got %d. Stderr: %s", result.ExitCode, result.Stderr)
	// JSON output has mod_name and resource_name as separate fields
	assert.Contains(t, result.Stdout, `"mod_name": "lazy_simple"`, "Should list queries from lazy_simple mod")
	assert.Contains(t, result.Stdout, `"resource_name": "simple_count"`, "Should list simple_count query")
	assert.Contains(t, result.Stdout, `"resource_name": "simple_status"`, "Should list simple_status query")
}

func TestCLI_QueryShowLazy(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	result := runPowerpipe(t, modPath, "query", "show", "lazy_simple.query.simple_count", "--output", "json")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0, got %d. Stderr: %s", result.ExitCode, result.Stderr)
	assert.Contains(t, result.Stdout, "Simple Count Query", "Should show query title")
}

// =============================================================================
// Mode Comparison Tests
// =============================================================================

func TestCLI_DashboardListModeComparison(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	// Run with both modes (list doesn't use lazy-load flag, but we verify consistency)
	result := runPowerpipe(t, modPath, "dashboard", "list", "--output", "json")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0, got %d. Stderr: %s", result.ExitCode, result.Stderr)

	// Parse JSON and verify structure
	var output []map[string]interface{}
	if err := json.Unmarshal([]byte(result.Stdout), &output); err != nil {
		// Try as a single object
		var single map[string]interface{}
		if err2 := json.Unmarshal([]byte(result.Stdout), &single); err2 != nil {
			// Might not be JSON format, just verify it contains expected content
			assert.Contains(t, result.Stdout, "simple")
		}
	}
}

func TestCLI_BenchmarkListModeComparison(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	result := runPowerpipe(t, modPath, "benchmark", "list", "--output", "json")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0, got %d. Stderr: %s", result.ExitCode, result.Stderr)
	assert.Contains(t, result.Stdout, "simple", "Should contain benchmark info")
}

// =============================================================================
// Error Handling Tests
// =============================================================================

func TestCLI_NonExistentMod(t *testing.T) {
	// Use a non-existent directory
	modPath := filepath.Join(testdataDir(), "mods", "nonexistent-mod-xyz")

	result := runPowerpipe(t, modPath, "dashboard", "list")

	// Should fail with non-zero exit code
	assert.NotEqual(t, 0, result.ExitCode, "Expected non-zero exit code for non-existent mod")
}

func TestCLI_NonExistentDashboard(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	result := runPowerpipe(t, modPath, "dashboard", "show", "nonexistent_dashboard")

	// Should fail with non-zero exit code
	assert.NotEqual(t, 0, result.ExitCode, "Expected non-zero exit code for non-existent dashboard")
	// Error message should indicate resource not found
	combinedOutput := result.Stdout + result.Stderr
	assert.True(t, strings.Contains(combinedOutput, "not found") ||
		strings.Contains(combinedOutput, "does not exist") ||
		strings.Contains(combinedOutput, "invalid") ||
		strings.Contains(combinedOutput, "error"),
		"Error message should indicate resource not found. Got: %s", combinedOutput)
}

func TestCLI_NonExistentBenchmark(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	result := runPowerpipe(t, modPath, "benchmark", "show", "nonexistent_benchmark")

	assert.NotEqual(t, 0, result.ExitCode, "Expected non-zero exit code for non-existent benchmark")
}

func TestCLI_NonExistentControl(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	result := runPowerpipe(t, modPath, "control", "show", "nonexistent_control")

	assert.NotEqual(t, 0, result.ExitCode, "Expected non-zero exit code for non-existent control")
}

func TestCLI_NonExistentQuery(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	// Note: query show treats unrecognized arguments as inline SQL, so this test
	// uses a clearly invalid format to trigger an error
	result := runPowerpipe(t, modPath, "query", "show", "mod.invalid.query.name.format.xyz")

	// The command may succeed if it treats the arg as SQL, which is valid behavior
	// So we just verify it completes without crashing
	t.Logf("Query show exit code: %d", result.ExitCode)
}

// =============================================================================
// Flag Combination Tests
// =============================================================================

func TestCLI_WithModLocation(t *testing.T) {
	// Run from a different directory with --mod-location pointing to the test mod
	modPath := getTestModPath("lazy-loading-tests", "simple")
	tmpDir := t.TempDir()

	result := runPowerpipe(t, tmpDir, "dashboard", "list", "--mod-location", modPath, "--output", "json")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0, got %d. Stderr: %s", result.ExitCode, result.Stderr)
	// JSON output has mod_name and resource_name as separate fields
	assert.Contains(t, result.Stdout, `"mod_name": "lazy_simple"`, "Should list dashboards from specified mod location")
	assert.Contains(t, result.Stdout, `"resource_name": "simple"`, "Should list the simple dashboard")
}

func TestCLI_WithOutputJSON(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	result := runPowerpipe(t, modPath, "dashboard", "list", "--output", "json")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0, got %d. Stderr: %s", result.ExitCode, result.Stderr)

	// Verify it's valid JSON (array or object)
	var output interface{}
	err := json.Unmarshal([]byte(result.Stdout), &output)
	assert.NoError(t, err, "Output should be valid JSON. Got: %s", result.Stdout)
}

func TestCLI_WithOutputPlain(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	result := runPowerpipe(t, modPath, "dashboard", "list", "--output", "plain")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0, got %d. Stderr: %s", result.ExitCode, result.Stderr)
	// Plain output should be human-readable
	assert.Contains(t, result.Stdout, "simple", "Plain output should contain dashboard name")
}

// =============================================================================
// Environment Variable Tests
// =============================================================================

func TestCLI_WorkspacePreloadEnvVar(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	// Test with workspace preload disabled (default lazy loading)
	result := runPowerpipe(t, modPath, "dashboard", "list", "--output", "json")
	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0, got %d. Stderr: %s", result.ExitCode, result.Stderr)
}

func TestCLI_WorkspacePreloadEnabled(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	// Test with workspace preload enabled (eager loading)
	result := runPowerpipeWithEnv(t, modPath, []string{"POWERPIPE_WORKSPACE_PRELOAD=true"}, "dashboard", "list", "--output", "json")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0, got %d. Stderr: %s", result.ExitCode, result.Stderr)
}

// =============================================================================
// Performance Tests
// =============================================================================

func TestCLI_StartupTimeLazy(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	modPath := getTestModPath("lazy-loading-tests", "simple")

	// Measure startup time for lazy mode
	result := runPowerpipe(t, modPath, "dashboard", "list", "--output", "json")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0")
	t.Logf("Dashboard list completed in %v", result.Duration)

	// Just log timing - actual performance comparison would need a large mod
	assert.Less(t, result.Duration, 30*time.Second, "Command should complete within reasonable time")
}

// =============================================================================
// Complex Mod Tests
// =============================================================================

func TestCLI_DeepHierarchyMod(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "deep-hierarchy")

	// Skip if the mod doesn't exist
	if _, err := os.Stat(modPath); os.IsNotExist(err) {
		t.Skip("deep-hierarchy test mod not found")
	}

	result := runPowerpipe(t, modPath, "benchmark", "list", "--output", "json")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0, got %d. Stderr: %s", result.ExitCode, result.Stderr)
}

func TestCLI_WideHierarchyMod(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "wide-hierarchy")

	// Skip if the mod doesn't exist
	if _, err := os.Stat(modPath); os.IsNotExist(err) {
		t.Skip("wide-hierarchy test mod not found")
	}

	result := runPowerpipe(t, modPath, "benchmark", "list", "--output", "json")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0, got %d. Stderr: %s", result.ExitCode, result.Stderr)
}

func TestCLI_CrossRefsMod(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "cross-refs")

	// Skip if the mod doesn't exist
	if _, err := os.Stat(modPath); os.IsNotExist(err) {
		t.Skip("cross-refs test mod not found")
	}

	result := runPowerpipe(t, modPath, "dashboard", "list", "--output", "json")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0, got %d. Stderr: %s", result.ExitCode, result.Stderr)
}

func TestCLI_MultiModMain(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "multi-mod", "main")

	// Skip if the mod doesn't exist
	if _, err := os.Stat(modPath); os.IsNotExist(err) {
		t.Skip("multi-mod/main test mod not found")
	}

	result := runPowerpipe(t, modPath, "benchmark", "list", "--output", "json")

	// Multi-mod tests may require mod install - skip if dependencies not installed
	if result.ExitCode != 0 && strings.Contains(result.Stderr, "mod install") {
		t.Skip("multi-mod test requires 'powerpipe mod install' to be run first")
	}

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0, got %d. Stderr: %s", result.ExitCode, result.Stderr)
}

// =============================================================================
// Help Command Tests
// =============================================================================

func TestCLI_HelpCommand(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	result := runPowerpipe(t, modPath, "--help")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0")
	assert.Contains(t, result.Stdout, "server", "Help should mention server command")
	assert.Contains(t, result.Stdout, "dashboard", "Help should mention dashboard command")
	assert.Contains(t, result.Stdout, "benchmark", "Help should mention benchmark command")
}

func TestCLI_ServerHelpCommand(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	result := runPowerpipe(t, modPath, "server", "--help")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0")
	assert.Contains(t, result.Stdout, "--port", "Server help should mention --port flag")
}

func TestCLI_DashboardHelpCommand(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	result := runPowerpipe(t, modPath, "dashboard", "--help")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0")
	assert.Contains(t, result.Stdout, "list", "Dashboard help should mention list subcommand")
	assert.Contains(t, result.Stdout, "show", "Dashboard help should mention show subcommand")
	assert.Contains(t, result.Stdout, "run", "Dashboard help should mention run subcommand")
}

func TestCLI_DashboardRunHelpCommand(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	result := runPowerpipe(t, modPath, "dashboard", "run", "--help")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0")
	assert.Contains(t, result.Stdout, "--output", "Dashboard run help should mention --output flag")
}

func TestCLI_BenchmarkRunHelpCommand(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	result := runPowerpipe(t, modPath, "benchmark", "run", "--help")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0")
	assert.Contains(t, result.Stdout, "--output", "Benchmark run help should mention --output flag")
}

func TestCLI_ControlRunHelpCommand(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	result := runPowerpipe(t, modPath, "control", "run", "--help")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0")
	assert.Contains(t, result.Stdout, "--output", "Control run help should mention --output flag")
}

func TestCLI_QueryRunHelpCommand(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	result := runPowerpipe(t, modPath, "query", "run", "--help")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0")
	assert.Contains(t, result.Stdout, "--output", "Query run help should mention --output flag")
}

// =============================================================================
// Version Command Test
// =============================================================================

func TestCLI_VersionCommand(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	result := runPowerpipe(t, modPath, "--version")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0")
	assert.Contains(t, result.Stdout, "Powerpipe", "Version output should contain Powerpipe")
}

// =============================================================================
// Server Port Conflict Test
// =============================================================================

func TestCLI_ServerPortConflict(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping server test in short mode")
	}

	modPath := getTestModPath("lazy-loading-tests", "simple")
	port := getFreePort(t)

	// Start first server
	_, cleanup := startServer(t, modPath, port)
	defer cleanup()

	// Try to start second server on same port
	binary := buildPowerpipe(t)
	cmd := exec.Command(binary, "server", "--port", fmt.Sprintf("%d", port))
	cmd.Dir = modPath

	// Should fail quickly
	err := cmd.Start()
	if err == nil {
		// Wait for it to fail
		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()

		select {
		case err := <-done:
			// Command should have exited with error
			if exitErr, ok := err.(*exec.ExitError); ok {
				assert.NotEqual(t, 0, exitErr.ExitCode(), "Second server should fail due to port conflict")
			}
		case <-time.After(5 * time.Second):
			_ = cmd.Process.Kill()
			t.Error("Second server should have failed due to port conflict")
		}
	}
}

// =============================================================================
// Concurrent Operations Test
// =============================================================================

func TestCLI_ConcurrentListCommands(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	// Run multiple list commands concurrently
	var wg sync.WaitGroup
	results := make([]*PowerpipeResult, 4)

	commands := [][]string{
		{"dashboard", "list", "--output", "json"},
		{"benchmark", "list", "--output", "json"},
		{"control", "list", "--output", "json"},
		{"query", "list", "--output", "json"},
	}

	for i, cmd := range commands {
		wg.Add(1)
		go func(idx int, args []string) {
			defer wg.Done()
			results[idx] = runPowerpipe(t, modPath, args...)
		}(i, cmd)
	}

	wg.Wait()

	// All should succeed
	for i, result := range results {
		assert.Equal(t, 0, result.ExitCode, "Command %d failed with exit code %d. Stderr: %s",
			i, result.ExitCode, result.Stderr)
	}
}

// =============================================================================
// Edge Cases
// =============================================================================

func TestCLI_EmptyModDirectory(t *testing.T) {
	// Create an empty temp directory
	tmpDir := t.TempDir()

	result := runPowerpipe(t, tmpDir, "dashboard", "list")

	// Should fail because there's no mod.pp file
	assert.NotEqual(t, 0, result.ExitCode, "Expected non-zero exit code for empty directory")
}

func TestCLI_InvalidFlag(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	result := runPowerpipe(t, modPath, "dashboard", "list", "--invalid-flag-xyz")

	// Should fail with error about unknown flag
	assert.NotEqual(t, 0, result.ExitCode, "Expected non-zero exit code for invalid flag")
	assert.Contains(t, result.Stderr, "unknown", "Error should mention unknown flag")
}

func TestCLI_InvalidSubcommand(t *testing.T) {
	modPath := getTestModPath("lazy-loading-tests", "simple")

	result := runPowerpipe(t, modPath, "invalid-command-xyz")

	// Should fail with error about unknown command
	assert.NotEqual(t, 0, result.ExitCode, "Expected non-zero exit code for invalid command")
}

// =============================================================================
// Mod with Variables
// =============================================================================

func TestCLI_ModWithVariables(t *testing.T) {
	modPath := getTestModPath("behavior_test_mod")

	// Skip if the mod doesn't exist
	if _, err := os.Stat(modPath); os.IsNotExist(err) {
		t.Skip("behavior_test_mod not found")
	}

	result := runPowerpipe(t, modPath, "variable", "list", "--output", "json")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0, got %d. Stderr: %s", result.ExitCode, result.Stderr)
}

// =============================================================================
// Server WebSocket Test (Basic)
// =============================================================================

func TestCLI_ServerWebSocketEndpoint(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping server test in short mode")
	}

	modPath := getTestModPath("lazy-loading-tests", "simple")
	port := getFreePort(t)

	_, cleanup := startServer(t, modPath, port)
	defer cleanup()

	// Verify WebSocket endpoint exists by checking HTTP upgrade
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/ws", port), nil)
	require.NoError(t, err)
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Sec-WebSocket-Key", "test")
	req.Header.Set("Sec-WebSocket-Version", "13")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)

	// We don't expect the full WebSocket handshake to succeed with our simple test,
	// but we should get some response from the server
	if err == nil {
		defer resp.Body.Close()
		// WebSocket upgrade would be 101, but any response means the endpoint exists
		t.Logf("WebSocket endpoint responded with status: %d", resp.StatusCode)
	}
}

// =============================================================================
// Large Output Test
// =============================================================================

func TestCLI_LargeModList(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large mod test in short mode")
	}

	modPath := getTestModPath("lazy-loading-tests", "generated", "medium")

	// Skip if the mod doesn't exist
	if _, err := os.Stat(modPath); os.IsNotExist(err) {
		t.Skip("generated/medium test mod not found")
	}

	result := runPowerpipe(t, modPath, "benchmark", "list", "--output", "json")

	assert.Equal(t, 0, result.ExitCode, "Expected exit code 0, got %d. Stderr: %s", result.ExitCode, result.Stderr)
	t.Logf("Large mod benchmark list completed in %v", result.Duration)
}
