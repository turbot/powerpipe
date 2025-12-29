#!/bin/bash
# Memory benchmarking script for powerpipe workspace loading
# Usage: ./scripts/memory_benchmark.sh [output_dir]

set -e

OUTPUT_DIR="${1:-./benchmark_results/memory/$(date +%Y%m%d_%H%M%S)}"
mkdir -p "$OUTPUT_DIR"

echo "Running memory benchmarks..."
echo "Results will be saved to $OUTPUT_DIR"
echo ""

# Get project root (directory containing go.mod)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROJECT_ROOT"

# Generate test mods if they don't exist
echo "Ensuring test mods are generated..."
for size in small medium large xlarge; do
    MOD_PATH="internal/testdata/mods/generated/$size"
    if [ ! -f "$MOD_PATH/mod.pp" ]; then
        echo "  Generating $size test mod..."
        go run scripts/generate_test_mods.go "$MOD_PATH" "$size"
    fi
done

echo ""
echo "Running memory benchmarks..."

# Run memory benchmarks with memory profiling
go test -bench=BenchmarkMemory -benchmem -benchtime=3x \
    -memprofile="$OUTPUT_DIR/mem.prof" \
    ./internal/workspace/... \
    -run='^$' \
    2>&1 | tee "$OUTPUT_DIR/benchmark.txt"

echo ""
echo "Generating memory profile analysis..."

# Generate memory profile analysis
if [ -f "$OUTPUT_DIR/mem.prof" ]; then
    go tool pprof -text -top=20 "$OUTPUT_DIR/mem.prof" > "$OUTPUT_DIR/mem_top.txt" 2>/dev/null || true
    go tool pprof -text -cum -top=20 "$OUTPUT_DIR/mem.prof" > "$OUTPUT_DIR/mem_cum.txt" 2>/dev/null || true
fi

echo ""
echo "Running memory scaling test..."

# Run scaling test
go test -v -run=TestMemoryScaling ./internal/workspace/... \
    2>&1 | tee "$OUTPUT_DIR/scaling.txt"

echo ""
echo "Running memory profile test..."

# Run detailed memory profile test
go test -v -run=TestMemoryProfile ./internal/workspace/... \
    2>&1 | tee "$OUTPUT_DIR/profile.txt"

# Extract key metrics
echo ""
echo "Extracting metrics..."

# Generate summary
cat > "$OUTPUT_DIR/summary.md" << 'SUMMARY_HEADER'
# Memory Benchmark Results

SUMMARY_HEADER

cat >> "$OUTPUT_DIR/summary.md" << EOF
**Date**: $(date)
**Commit**: $(git rev-parse HEAD 2>/dev/null || echo "unknown")
**Branch**: $(git branch --show-current 2>/dev/null || echo "unknown")

## Benchmark Results
\`\`\`
$(grep -E "^BenchmarkMemory|heap-bytes|heap-objects|final-heap" "$OUTPUT_DIR/benchmark.txt" 2>/dev/null || echo "No benchmark results")
\`\`\`

## Memory Scaling
\`\`\`
$(grep -E "small:|medium:|large:|Scaling" "$OUTPUT_DIR/scaling.txt" 2>/dev/null || echo "No scaling results")
\`\`\`

## Top Memory Allocators
\`\`\`
$(head -30 "$OUTPUT_DIR/mem_top.txt" 2>/dev/null || echo "No profile data")
\`\`\`

## Memory Profile Details
\`\`\`
$(grep -E "Peak heap:|Final heap:|Peak objects:|Memory Report" "$OUTPUT_DIR/profile.txt" 2>/dev/null || echo "No profile details")
\`\`\`
EOF

echo ""
echo "========================================"
echo "Memory Benchmark Complete"
echo "========================================"
echo ""
echo "Results saved to: $OUTPUT_DIR"
echo ""
echo "Summary:"
cat "$OUTPUT_DIR/summary.md"
