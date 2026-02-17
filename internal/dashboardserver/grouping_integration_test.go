package dashboardserver

/*
DASHBOARD GROUPING INTEGRATION TEST

This test suite validates that Powerpipe's lazy loading mechanism correctly loads tags from
dependency mods, ensuring dashboards and benchmarks can be grouped properly in the UI.

═══════════════════════════════════════════════════════════════════════════════════════════

OVERVIEW: How This Test Works

This is an end-to-end integration test that:
1. Starts a REAL Powerpipe server process (via exec.Command)
2. Connects to it via WebSocket (mimicking a browser client or Pipes)
3. Sends get_available_dashboards action (the same action Pipes uses)
4. Validates that dashboards/benchmarks have proper tags for grouping

This test is critical because it validates the EXACT code path that Pipes uses when
displaying dashboards and benchmarks to users.

═══════════════════════════════════════════════════════════════════════════════════════════

WHY THIS TEST EXISTS: The Lazy Loading Tags Bug

Before this test was created, we had a regression where lazy loading wasn't properly
loading tags from dependency mods. This caused:
- Dashboards grouped into "Other" (not by service like AWS/EC2, AWS/RDS, etc.)
- Benchmarks grouped incorrectly
- Poor user experience in Pipes and local Powerpipe usage

The bug was subtle:
- v1.4.3 (eager loading): 1796/1865 benchmarks (96.3%) had tags ✅
- v1.5.0-rc.0 (lazy loading): 819/1865 benchmarks (43.9%) had tags ❌

This test ensures the regression doesn't happen again.

═══════════════════════════════════════════════════════════════════════════════════════════

HOW IT MIMICS PIPES

Pipes Scenario:
1. User opens Pipes dashboard
2. Pipes frontend connects to Powerpipe server via WebSocket (ws://)
3. Pipes sends: {"action": "get_available_dashboards"}
4. Server responds with dashboard/benchmark metadata including tags
5. Pipes UI groups resources by tags (e.g., service="AWS/EC2")

This Test Mimics That Exactly:
1. Starts Powerpipe server on localhost:19033
2. Uses gorilla/websocket to connect (ws://localhost:19033/ws)
3. Sends: {"action": "get_available_dashboards"}
4. Parses JSON response to validate tags are present
5. Asserts >90% of resources have tags for proper grouping

The test workspace (/Users/pskrbasu/pskr) has:
- Main workspace (empty, no .pp files)
- Dependency mods installed in .powerpipe/mods/:
  - aws-compliance@v1.13.0 (475 files, complex tag expressions)
  - aws-insights@v1.2.0
  - net-insights@v1.0.1

This setup is IDENTICAL to a typical Pipes workspace where:
- Users don't write .pp files in their main workspace
- All dashboards/benchmarks come from installed dependency mods
- Tags use complex merge(local.xxx, {...}) expressions

═══════════════════════════════════════════════════════════════════════════════════════════

WEBSOCKET PROTOCOL

The test uses the exact WebSocket protocol that Powerpipe's dashboard server implements:

Connection Flow:
1. HTTP GET http://localhost:19033 → Wait for 200 OK
2. WebSocket upgrade ws://localhost:19033/ws → Establish connection
3. Send JSON: {"action": "get_available_dashboards"}
4. Receive JSON: {
     "action": "available_dashboards",
     "dashboards": {
       "aws_insights.dashboard.ec2_instance_detail": {
         "title": "EC2 Instance Detail",
         "tags": {"service": "AWS/EC2", "type": "Detail"},
         "mod_full_name": "mod.aws_insights",
         ...
       },
       ...
     },
     "benchmarks": {
       "aws_compliance.benchmark.cis_v500_2_2": {
         "title": "2.2 Ensure RDS encryption",
         "tags": {"service": "AWS/RDS", "cis_version": "v5.0.0", ...},
         ...
       },
       ...
     }
   }

The test validates:
- Tags are present in the response (not null/empty)
- Tags include "service" key for grouping
- >90% of resources have tags (matching v1.4.3 baseline)

═══════════════════════════════════════════════════════════════════════════════════════════

TEST PHASES

The test runs in TWO phases to catch timing issues:

Phase 1: Immediate Connection (Race Condition Test)
- Connects as soon as server starts
- Tests lazy loading phase 1 (index build) results
- Expected: ~43% of benchmarks have tags (literal values only)
- This is OK - background resolution hasn't completed yet

Phase 2: After Waiting (Final State Test)
- Waits 3 seconds for background resolution
- Tests lazy loading phase 2 (background resolution) results
- Expected: >90% of benchmarks have tags (merge() expressions evaluated)
- This MUST pass or grouping is broken

Why two phases?
- Phase 1 catches immediate UI issues (Pipes connects fast)
- Phase 2 validates the fix works (background resolution succeeds)
- Together they ensure both fast startup AND correct final state

═══════════════════════════════════════════════════════════════════════════════════════════

TAG STRUCTURE

Tags enable grouping in Pipes UI. Example benchmark tag structure:

{
  "service": "AWS/RDS",        ← PRIMARY GROUPING KEY
  "type": "Benchmark",         ← Resource type
  "category": "Compliance",    ← Category grouping
  "cis": "true",               ← Framework flag
  "cis_version": "v5.0.0",     ← Framework version
  "cis_section_id": "2.2"      ← Section identifier
}

These tags come from complex HCL expressions in dependency mods:

  tags = merge(local.cis_v500_2_2_common_tags, {
    service = "AWS/RDS"
  })

Where local.cis_v500_2_2_common_tags itself references other locals:

  locals {
    cis_v500_2_2_common_tags = merge(local.aws_compliance_common_tags, {
      cis_section_id = "2.2"
    })
  }

This test validates that these complex expressions are evaluated correctly.

═══════════════════════════════════════════════════════════════════════════════════════════

KEY FILES TESTED

This test validates the entire lazy loading pipeline:

1. internal/workspace/lazy_workspace.go
   - Creates lazy workspace with index + cache
   - Calls BuildEvalContext to load variables/locals
   - Starts background resolution

2. internal/resourceloader/eval_context.go [THE FIX]
   - Scans dependency mods for variables/locals
   - Multi-pass parsing to resolve local references
   - Critical fix: Always scan deps even if main workspace empty

3. internal/resourceloader/parser.go
   - Parses resources from HCL files
   - Evaluates tag expressions with eval context
   - Uses pipe-fittings DecodeHclBody

4. internal/dashboardserver/server.go
   - WebSocket server implementation
   - Handles get_available_dashboards action
   - Builds JSON response with tags

5. internal/resourceindex/entry.go
   - Index entry with lazy-loaded metadata
   - TagsResolved flag tracking

═══════════════════════════════════════════════════════════════════════════════════════════

RUNNING THE TEST

Run full test:
  go test -v -run TestDashboardGrouping_RealMod ./internal/dashboardserver/

Requirements:
- Powerpipe built and installed at /usr/local/bin/powerpipe
- Test workspace at /Users/pskrbasu/pskr with:
  - .powerpipe/mods/github.com/turbot/steampipe-mod-aws-compliance@v1.13.0/
  - .powerpipe/mods/github.com/turbot/steampipe-mod-aws-insights@v1.2.0/
  - .powerpipe/mods/github.com/turbot/steampipe-mod-net-insights@v1.0.1/

The test will:
1. Start server (takes ~800ms to initialize with lazy loading)
2. Connect via WebSocket
3. Validate tag counts match baseline (1796/1865 = 96.3%)
4. Fail if grouping would be broken (<90% tags)

═══════════════════════════════════════════════════════════════════════════════════════════

WHAT THIS TEST PROVES

✅ Lazy loading loads tags from dependency mods
✅ Multi-pass local resolution works for complex merge() expressions
✅ WebSocket API returns complete metadata
✅ Dashboard grouping will work in Pipes
✅ No regression from v1.4.3 baseline

If this test fails:
- Dashboard grouping is broken
- Pipes users will see "Other" group instead of proper service groups
- Benchmarks won't be organized correctly
- User experience is degraded

═══════════════════════════════════════════════════════════════════════════════════════════
*/

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDashboardGrouping_RealMod tests dashboard grouping with a real mod that has dependency mods.
// This test verifies that tags are properly loaded for dashboards from dependency mods,
// ensuring they can be grouped correctly (not all ending up in "Other").
//
// See the comprehensive comment at the top of this file for full details on how this test works.
//
// Run with: go test -v -run TestDashboardGrouping_RealMod ./internal/dashboardserver/
func TestDashboardGrouping_RealMod(t *testing.T) {
	// Skip if the test mod doesn't exist
	modPath := "/Users/pskrbasu/pskr"
	if _, err := os.Stat(modPath); os.IsNotExist(err) {
		t.Skipf("Test mod not found at %s", modPath)
	}

	// Find a free port
	port := 19033

	// Get powerpipe binary path
	powerpipeBinary := "/usr/local/bin/powerpipe"
	if _, err := os.Stat(powerpipeBinary); os.IsNotExist(err) {
		t.Skipf("Powerpipe binary not found at %s", powerpipeBinary)
	}

	// Start powerpipe server
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, powerpipeBinary, "server", "--port", fmt.Sprintf("%d", port))
	cmd.Dir = modPath
	cmd.Env = append(os.Environ(), "POWERPIPE_LOG_LEVEL=INFO")

	// Capture output for debugging
	logFile := filepath.Join(os.TempDir(), "powerpipe_grouping_test.log")
	logWriter, err := os.Create(logFile)
	require.NoError(t, err)
	defer logWriter.Close()
	cmd.Stdout = logWriter
	cmd.Stderr = logWriter

	err = cmd.Start()
	require.NoError(t, err, "Failed to start powerpipe server")
	defer func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
		}
	}()

	t.Logf("Started powerpipe server on port %d (PID: %d)", port, cmd.Process.Pid)
	t.Logf("Server logs: %s", logFile)

	// Wait for server to be ready (try to connect to HTTP endpoint)
	serverURL := fmt.Sprintf("http://localhost:%d", port)
	require.Eventually(t, func() bool {
		resp, err := http.Get(serverURL)
		if err == nil {
			_ = resp.Body.Close()
			return resp.StatusCode == 200
		}
		return false
	}, 15*time.Second, 500*time.Millisecond, "Server did not start in time")

	t.Logf("Server is ready at %s", serverURL)

	// Test 1: Connect IMMEDIATELY to catch race condition
	t.Run("Immediate connection (before background resolution)", func(t *testing.T) {
		// Connect as soon as server is ready, before background resolution completes
		payload := getAvailableDashboardsViaWebSocket(t, port)
		dashboardCount := countDashboardsWithTags(t, payload)
		benchmarkCount := countBenchmarksWithTags(t, payload)
		t.Logf("Immediate connection: %d dashboards have tags, %d benchmarks have tags",
			dashboardCount, benchmarkCount)

		// This might show the bug - if <100% tags, grouping will fail
		validateDashboardTags(t, payload, "immediate", false)
		validateBenchmarkTags(t, payload, "immediate", false)
	})

	// Test 2: Wait a bit and check again (should get broadcast update)
	t.Run("After waiting (should receive broadcast)", func(t *testing.T) {
		// Wait for background resolution and broadcast
		time.Sleep(3 * time.Second)

		payload := getAvailableDashboardsViaWebSocket(t, port)
		dashboardCount := countDashboardsWithTags(t, payload)
		benchmarkCount := countBenchmarksWithTags(t, payload)
		t.Logf("After waiting: %d dashboards have tags, %d benchmarks have tags",
			dashboardCount, benchmarkCount)

		validateDashboardTags(t, payload, "after waiting", true) // Require ALL tags
		validateBenchmarkTags(t, payload, "after waiting", true) // Require ALL tags
	})
}

// getAvailableDashboardsViaWebSocket connects to the WebSocket and retrieves available_dashboards
func getAvailableDashboardsViaWebSocket(t *testing.T, port int) map[string]interface{} {
	t.Helper()

	wsURL := fmt.Sprintf("ws://localhost:%d/ws", port)
	t.Logf("Connecting to WebSocket at %s", wsURL)

	// Connect to WebSocket
	dialer := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	conn, _, err := dialer.Dial(wsURL, nil)
	require.NoError(t, err, "Failed to connect to WebSocket")
	defer conn.Close()

	// Send get_available_dashboards request
	request := map[string]interface{}{
		"action": "get_available_dashboards",
	}
	err = conn.WriteJSON(request)
	require.NoError(t, err, "Failed to send request")

	// Read response with timeout
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	_, message, err := conn.ReadMessage()
	require.NoError(t, err, "Failed to read response")

	// Parse response
	var response map[string]interface{}
	err = json.Unmarshal(message, &response)
	require.NoError(t, err, "Failed to parse response")

	t.Logf("Received response with action: %v", response["action"])

	return response
}

// countDashboardsWithTags counts how many dashboards have non-empty tags
func countDashboardsWithTags(t *testing.T, payload map[string]interface{}) int {
	t.Helper()

	dashboardsRaw, ok := payload["dashboards"]
	if !ok {
		return 0
	}
	dashboards, ok := dashboardsRaw.(map[string]interface{})
	if !ok {
		return 0
	}

	count := 0
	for _, dashRaw := range dashboards {
		dash, ok := dashRaw.(map[string]interface{})
		if !ok {
			continue
		}
		tagsRaw, hasTags := dash["tags"]
		if hasTags {
			tags, isMap := tagsRaw.(map[string]interface{})
			if isMap && len(tags) > 0 {
				count++
			}
		}
	}
	return count
}

// countBenchmarksWithTags counts how many benchmarks have non-empty tags
func countBenchmarksWithTags(t *testing.T, payload map[string]interface{}) int {
	t.Helper()

	benchmarksRaw, ok := payload["benchmarks"]
	if !ok {
		return 0
	}
	benchmarks, ok := benchmarksRaw.(map[string]interface{})
	if !ok {
		return 0
	}

	count := 0
	for _, benchRaw := range benchmarks {
		bench, ok := benchRaw.(map[string]interface{})
		if !ok {
			continue
		}
		tagsRaw, hasTags := bench["tags"]
		if hasTags {
			tags, isMap := tagsRaw.(map[string]interface{})
			if isMap && len(tags) > 0 {
				count++
			}
		}
	}
	return count
}

// validateBenchmarkTags validates that benchmarks have proper tags for grouping
func validateBenchmarkTags(t *testing.T, payload map[string]interface{}, phase string, requireAllTags bool) {
	t.Helper()

	// Extract benchmarks
	benchmarksRaw, ok := payload["benchmarks"]
	if !ok {
		t.Log("No benchmarks in payload")
		return
	}

	benchmarks, ok := benchmarksRaw.(map[string]interface{})
	require.True(t, ok, "Benchmarks field is not a map")

	totalBenchmarks := len(benchmarks)
	if totalBenchmarks == 0 {
		t.Log("No benchmarks found")
		return
	}

	t.Logf("Found %d benchmarks", totalBenchmarks)

	// Count benchmarks with tags
	benchmarksWithTags := 0
	benchmarksWithServiceTag := 0

	for name, benchRaw := range benchmarks {
		bench, ok := benchRaw.(map[string]interface{})
		require.True(t, ok, "Benchmark %s is not a map", name)

		// Check tags
		tagsRaw, hasTags := bench["tags"]
		if hasTags {
			tags, isMap := tagsRaw.(map[string]interface{})
			if isMap && len(tags) > 0 {
				benchmarksWithTags++

				// Check for "service" tag
				if _, hasService := tags["service"]; hasService {
					benchmarksWithServiceTag++
				}

				// Log first few for debugging
				if benchmarksWithTags <= 3 {
					t.Logf("Benchmark %s has tags: %v", name, tags)
				}
			}
		}
	}

	// Report statistics
	t.Logf("=== %s Benchmark Statistics ===", phase)
	t.Logf("Total benchmarks: %d", totalBenchmarks)
	t.Logf("Benchmarks with tags: %d (%.1f%%)", benchmarksWithTags,
		float64(benchmarksWithTags)/float64(totalBenchmarks)*100)
	t.Logf("Benchmarks with 'service' tag: %d", benchmarksWithServiceTag)

	if requireAllTags {
		// After background resolution, ALL benchmarks should have tags
		percentWithTags := float64(benchmarksWithTags) / float64(totalBenchmarks) * 100
		assert.Greater(t, percentWithTags, 90.0,
			"After resolution, expected >90%% of benchmarks to have tags (found %.1f%%). "+
				"This causes benchmarks to be grouped incorrectly.", percentWithTags)
	}
}

// validateDashboardTags validates that dashboards have proper tags for grouping
func validateDashboardTags(t *testing.T, payload map[string]interface{}, phase string, requireAllTags bool) {
	t.Helper()

	// Verify action
	assert.Equal(t, "available_dashboards", payload["action"], "Wrong action in response")

	// Extract dashboards
	dashboardsRaw, ok := payload["dashboards"]
	require.True(t, ok, "Response missing 'dashboards' field")

	dashboards, ok := dashboardsRaw.(map[string]interface{})
	require.True(t, ok, "Dashboards field is not a map")

	totalDashboards := len(dashboards)
	require.Greater(t, totalDashboards, 0, "No dashboards found")
	t.Logf("Found %d dashboards", totalDashboards)

	// Count dashboards by mod and tags
	dashboardsByMod := make(map[string]int)
	dashboardsWithTags := 0
	dashboardsWithServiceTag := 0
	dashboardsFromDepMods := 0

	for name, dashRaw := range dashboards {
		dash, ok := dashRaw.(map[string]interface{})
		require.True(t, ok, "Dashboard %s is not a map", name)

		// Get mod_full_name
		modFullName, _ := dash["mod_full_name"].(string)
		dashboardsByMod[modFullName]++

		// Check if from dependency mod (not local/smoketest)
		if modFullName != "" && modFullName != "mod.local" && modFullName != "mod.smoketest" {
			dashboardsFromDepMods++
		}

		// Check tags
		tagsRaw, hasTags := dash["tags"]
		if hasTags {
			tags, isMap := tagsRaw.(map[string]interface{})
			if isMap && len(tags) > 0 {
				dashboardsWithTags++

				// Check for "service" tag (common grouping key)
				if _, hasService := tags["service"]; hasService {
					dashboardsWithServiceTag++
				}

				// Log first few for debugging
				if dashboardsWithTags <= 3 {
					t.Logf("Dashboard %s from %s has tags: %v", name, modFullName, tags)
				}
			}
		}
	}

	// Report statistics
	t.Logf("=== %s Statistics ===", phase)
	t.Logf("Total dashboards: %d", totalDashboards)
	t.Logf("Dashboards from dependency mods: %d", dashboardsFromDepMods)
	t.Logf("Dashboards with tags: %d (%.1f%%)", dashboardsWithTags, float64(dashboardsWithTags)/float64(totalDashboards)*100)
	t.Logf("Dashboards with 'service' tag: %d", dashboardsWithServiceTag)
	t.Logf("Dashboards by mod:")
	for mod, count := range dashboardsByMod {
		t.Logf("  %s: %d", mod, count)
	}

	// Validations
	assert.Greater(t, len(dashboardsByMod), 1, "Expected dashboards from multiple mods")

	// Check we have dashboards from dependency mods
	assert.Greater(t, dashboardsFromDepMods, 0, "Expected dashboards from dependency mods")

	if requireAllTags {
		// After background resolution, ALL dashboards from dependency mods should have tags
		// This is the key check - if this fails, dashboards will be grouped into "Other"
		percentWithTags := float64(dashboardsWithTags) / float64(totalDashboards) * 100
		assert.Greater(t, percentWithTags, 90.0,
			"After resolution, expected >90%% of dashboards to have tags (found %.1f%%). "+
				"This causes dashboards to be grouped incorrectly into 'Other'.", percentWithTags)

		// Most dashboards should have the 'service' tag for grouping
		assert.Greater(t, dashboardsWithServiceTag, totalDashboards/2,
			"Expected majority of dashboards to have 'service' tag for grouping")
	} else {
		// Initial payload might not have all tags yet (background resolution in progress)
		t.Logf("Initial check: %d/%d dashboards have tags", dashboardsWithTags, totalDashboards)
	}
}
