#!/bin/bash
# Measures real-world server startup times for eager vs lazy loading
#
# Usage: ./benchmark-loading.sh [mod_path] [runs]
#   mod_path: Path to the mod to benchmark (default: /Users/nathan/src/powerpipe-performance-test)
#   runs: Number of runs for each test (default: 3)
#
# Environment variables:
#   POWERPIPE_WORKSPACE_PRELOAD=true  Forces eager loading
#   POWERPIPE_WORKSPACE_PRELOAD=false Forces lazy loading (default)

set -e

# Use gtimeout on macOS, timeout on Linux
if command -v gtimeout &> /dev/null; then
    TIMEOUT_CMD="gtimeout"
elif command -v timeout &> /dev/null; then
    TIMEOUT_CMD="timeout"
else
    echo "ERROR: timeout command not found. Install coreutils (brew install coreutils on macOS)"
    exit 1
fi

# Use gdate on macOS for millisecond precision, date on Linux
if command -v gdate &> /dev/null; then
    DATE_CMD="gdate"
else
    DATE_CMD="date"
fi

# Function to get current time in milliseconds
get_time_ms() {
    $DATE_CMD +%s%3N
}

MOD_PATH="${1:-/Users/nathan/src/powerpipe-performance-test}"
RUNS="${2:-3}"
PORT=9999

echo "=== Powerpipe Loading Benchmarks ==="
echo "Date: $(date)"
echo "Mod path: $MOD_PATH"
echo "Runs: $RUNS"
echo ""

# Check if mod path exists
if [ ! -d "$MOD_PATH" ]; then
    echo "ERROR: Mod path does not exist: $MOD_PATH"
    exit 1
fi

# Count files in the mod
FILE_COUNT=$(find "$MOD_PATH" -name "*.pp" -o -name "*.sp" 2>/dev/null | wc -l | tr -d ' ')
echo "Files in mod: $FILE_COUNT"
echo ""

# Function to measure server startup time
measure_startup() {
    local mode="$1"
    local preload_setting="$2"

    echo "--- $mode Loading ---"

    for i in $(seq 1 $RUNS); do
        # Kill any existing powerpipe server
        pkill -f "powerpipe server" 2>/dev/null || true
        sleep 1

        # Start the server and capture output
        POWERPIPE_WORKSPACE_PRELOAD="$preload_setting" \
        $TIMEOUT_CMD 120 powerpipe server \
            --mod-location "$MOD_PATH" \
            --port "$PORT" \
            2>&1 &
        SERVER_PID=$!

        # Wait for server to start and capture timing
        local start_time=$(get_time_ms)
        local ready=false
        local timeout_count=0

        while [ $timeout_count -lt 600 ]; do
            # Check if server is responding
            if curl -s "http://localhost:$PORT/api/v0/available_dashboards" > /dev/null 2>&1; then
                local end_time=$(get_time_ms)
                local elapsed=$((end_time - start_time))
                echo "  Run $i: ${elapsed}ms (server ready)"
                ready=true
                break
            fi

            # Check if server process died
            if ! kill -0 $SERVER_PID 2>/dev/null; then
                echo "  Run $i: Server exited unexpectedly"
                break
            fi

            sleep 0.1
            timeout_count=$((timeout_count + 1))
        done

        if [ "$ready" = false ]; then
            echo "  Run $i: TIMEOUT"
        fi

        # Clean up
        kill $SERVER_PID 2>/dev/null || true
        wait $SERVER_PID 2>/dev/null || true
        sleep 1
    done
    echo ""
}

# Run with eager loading
measure_startup "Eager" "true"

# Run with lazy loading
measure_startup "Lazy" "false"

echo "=== Benchmark Complete ==="
