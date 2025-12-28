#!/bin/bash
set -e

OUTPUT_DIR="${1:-./benchmark_results}"
mkdir -p "$OUTPUT_DIR"

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
RESULTS_FILE="$OUTPUT_DIR/benchmark_${TIMESTAMP}.txt"
JSON_FILE="$OUTPUT_DIR/benchmark_${TIMESTAMP}.json"

echo "Running performance benchmarks..."
echo "Results will be saved to $RESULTS_FILE"

# Enable timing
export POWERPIPE_TIMING=1

# Get project root (directory containing go.mod)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROJECT_ROOT"

# Generate test mods if they don't exist
echo "Ensuring test mods are generated..."
for size in small medium large; do
    MOD_PATH="internal/testdata/mods/generated/$size"
    if [ ! -f "$MOD_PATH/mod.pp" ]; then
        echo "  Generating $size test mod..."
        go run scripts/generate_test_mods.go "$MOD_PATH" "$size"
    fi
done

echo ""
echo "Running benchmarks..."

# Run benchmarks with memory profiling
# -benchtime=3s for reasonable runtime
# -count=3 for statistical significance
go test -bench=. -benchmem -benchtime=3s -count=1 \
    ./internal/workspace/... \
    ./internal/dashboardserver/... \
    -run='^$' \
    2>&1 | tee "$RESULTS_FILE"

# Parse results to JSON for comparison
echo ""
echo "Parsing results to JSON..."
go run scripts/parse_benchmark_results.go "$RESULTS_FILE" > "$JSON_FILE" 2>/dev/null || echo "Warning: Could not parse to JSON"

echo ""
echo "Benchmark complete. Results saved to:"
echo "  Text: $RESULTS_FILE"
echo "  JSON: $JSON_FILE"
