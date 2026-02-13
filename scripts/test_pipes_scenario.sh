#!/bin/bash

# Script to test the Pipes lazy loading scenario
# This reproduces the exact issue that was happening in Pipes and verifies the fix works

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$REPO_ROOT"

echo "======================================"
echo "Testing Pipes Lazy Loading Scenario"
echo "======================================"
echo ""
echo "This test reproduces the bug where:"
echo "  - installed_mods was empty in server_metadata"
echo "  - Dashboards were all grouped under 'Other'"
echo "  - Tags weren't populated on benchmarks"
echo ""
echo "Running tests..."
echo ""

# Run the Pipes scenario tests
go test -v -run TestPipesScenario ./internal/workspace/ 2>&1 | tee /tmp/pipes_scenario_test.log

# Check if tests passed
if [ ${PIPESTATUS[0]} -eq 0 ]; then
    echo ""
    echo "======================================"
    echo "✓ ALL TESTS PASSED!"
    echo "======================================"
    echo ""
    echo "The fix is working correctly:"
    echo "  ✓ Lazy loading works without eager fallback"
    echo "  ✓ installed_mods properly populated"
    echo "  ✓ Server metadata correctly built"
    echo "  ✓ Dashboards grouped by mod"
    echo "  ✓ Tags populated on resources"
    echo ""
    echo "The workspace is ready for deployment to Pipes!"
    exit 0
else
    echo ""
    echo "======================================"
    echo "✗ TESTS FAILED!"
    echo "======================================"
    echo ""
    echo "See /tmp/pipes_scenario_test.log for details"
    exit 1
fi
