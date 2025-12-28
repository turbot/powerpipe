# Task 3: Create Performance Benchmark Tests

## Objective

Create benchmark tests that measure mod loading performance with various sizes of mod collections, enabling quantitative comparison before and after optimizations.

## Context

- Need to measure performance with realistic workloads
- Benchmarks should scale from small to large mod setups
- Results should be reproducible and comparable
- Use Go's built-in benchmark framework plus timing instrumentation

## Dependencies

### Prerequisites
- Task 1 (Instrumentation) must be complete
- Task 2 (Mod Loading Tests) should be complete for fixture patterns

### Files to Create
- `internal/workspace/load_workspace_benchmark_test.go`
- `internal/dashboardserver/payload_benchmark_test.go`
- `testdata/mods/generated/` - Generated large mod fixtures
- `scripts/generate_test_mods.go` - Script to generate test fixtures

## Implementation Details

### 1. Create Mod Generator Script

`scripts/generate_test_mods.go`:

```go
//go:build ignore

package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: go run generate_test_mods.go <output_dir> <size>")
        fmt.Println("Sizes: small (10 dashboards), medium (50), large (200), xlarge (500)")
        os.Exit(1)
    }

    outputDir := os.Args[1]
    size := os.Args[2]

    counts := map[string]struct {
        dashboards int
        queries    int
        controls   int
        benchmarks int
    }{
        "small":  {10, 20, 30, 5},
        "medium": {50, 100, 150, 20},
        "large":  {200, 400, 500, 50},
        "xlarge": {500, 1000, 1500, 100},
    }

    c, ok := counts[size]
    if !ok {
        fmt.Printf("Unknown size: %s\n", size)
        os.Exit(1)
    }

    generateMod(outputDir, size, c.dashboards, c.queries, c.controls, c.benchmarks)
}

func generateMod(dir, name string, dashboards, queries, controls, benchmarks int) {
    os.MkdirAll(dir, 0755)

    // Generate mod.pp
    modContent := fmt.Sprintf(`mod "%s_test" {
  title = "%s Test Mod"
  description = "Generated mod for performance testing with %d dashboards, %d queries, %d controls, %d benchmarks"
}
`, name, strings.Title(name), dashboards, queries, controls, benchmarks)

    os.WriteFile(filepath.Join(dir, "mod.pp"), []byte(modContent), 0644)

    // Generate queries
    var queryContent strings.Builder
    for i := 0; i < queries; i++ {
        queryContent.WriteString(fmt.Sprintf(`
query "query_%d" {
  title = "Query %d"
  description = "Test query %d for performance benchmarking"
  sql = <<-EOQ
    SELECT
      %d as id,
      'value_%d' as name,
      now() as created_at
    FROM generate_series(1, 100)
  EOQ
}
`, i, i, i, i, i))
    }
    os.WriteFile(filepath.Join(dir, "queries.pp"), []byte(queryContent.String()), 0644)

    // Generate controls
    var controlContent strings.Builder
    for i := 0; i < controls; i++ {
        controlContent.WriteString(fmt.Sprintf(`
control "control_%d" {
  title = "Control %d"
  description = "Test control %d"
  sql = <<-EOQ
    SELECT
      'resource_%d' as resource,
      'ok' as status,
      'Control %d passed' as reason
  EOQ
  tags = {
    category = "test"
    index    = "%d"
  }
}
`, i, i, i, i, i, i))
    }
    os.WriteFile(filepath.Join(dir, "controls.pp"), []byte(controlContent.String()), 0644)

    // Generate benchmarks
    var benchmarkContent strings.Builder
    controlsPerBenchmark := controls / benchmarks
    for i := 0; i < benchmarks; i++ {
        children := make([]string, 0, controlsPerBenchmark)
        for j := 0; j < controlsPerBenchmark && (i*controlsPerBenchmark+j) < controls; j++ {
            children = append(children, fmt.Sprintf("control.control_%d", i*controlsPerBenchmark+j))
        }
        benchmarkContent.WriteString(fmt.Sprintf(`
benchmark "benchmark_%d" {
  title = "Benchmark %d"
  description = "Test benchmark %d"
  children = [
    %s
  ]
  tags = {
    category = "test"
  }
}
`, i, i, i, strings.Join(children, ",\n    ")))
    }
    os.WriteFile(filepath.Join(dir, "benchmarks.pp"), []byte(benchmarkContent.String()), 0644)

    // Generate dashboards
    var dashboardContent strings.Builder
    for i := 0; i < dashboards; i++ {
        dashboardContent.WriteString(fmt.Sprintf(`
dashboard "dashboard_%d" {
  title = "Dashboard %d"
  description = "Test dashboard %d for performance benchmarking"

  tags = {
    category = "test"
    index    = "%d"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_%d.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_%d.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_%d.sql
    }
  }
}
`, i, i, i, i, i%queries, i%queries, i%queries))
    }
    os.WriteFile(filepath.Join(dir, "dashboards.pp"), []byte(dashboardContent.String()), 0644)

    fmt.Printf("Generated %s mod in %s:\n", name, dir)
    fmt.Printf("  - %d dashboards\n", dashboards)
    fmt.Printf("  - %d queries\n", queries)
    fmt.Printf("  - %d controls\n", controls)
    fmt.Printf("  - %d benchmarks\n", benchmarks)
}
```

### 2. Create Workspace Loading Benchmarks

`internal/workspace/load_workspace_benchmark_test.go`:

```go
package workspace_test

import (
    "context"
    "os"
    "os/exec"
    "path/filepath"
    "testing"

    "github.com/turbot/powerpipe/internal/timing"
    "github.com/turbot/powerpipe/internal/workspace"
)

var benchmarkSizes = []string{"small", "medium", "large"}

func init() {
    // Enable timing for benchmarks
    os.Setenv("POWERPIPE_TIMING", "1")
}

// BenchmarkLoadWorkspace_Small benchmarks loading a small mod (10 dashboards)
func BenchmarkLoadWorkspace_Small(b *testing.B) {
    benchmarkLoadWorkspace(b, "small")
}

// BenchmarkLoadWorkspace_Medium benchmarks loading a medium mod (50 dashboards)
func BenchmarkLoadWorkspace_Medium(b *testing.B) {
    benchmarkLoadWorkspace(b, "medium")
}

// BenchmarkLoadWorkspace_Large benchmarks loading a large mod (200 dashboards)
func BenchmarkLoadWorkspace_Large(b *testing.B) {
    benchmarkLoadWorkspace(b, "large")
}

func benchmarkLoadWorkspace(b *testing.B, size string) {
    modPath := ensureGeneratedMod(b, size)
    ctx := context.Background()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        timing.Reset()
        _, ew := workspace.LoadWorkspacePromptingForVariables(ctx, modPath)
        if ew.GetError() != nil {
            b.Fatalf("Failed to load workspace: %v", ew.GetError())
        }
    }
    b.StopTimer()

    // Report timing breakdown
    if timing.IsEnabled() {
        b.Log(timing.ReportJSON())
    }
}

// BenchmarkLoadWorkspace_Parallel tests parallel loading capability
func BenchmarkLoadWorkspace_Parallel(b *testing.B) {
    modPath := ensureGeneratedMod(b, "medium")
    ctx := context.Background()

    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            _, ew := workspace.LoadWorkspacePromptingForVariables(ctx, modPath)
            if ew.GetError() != nil {
                b.Fatalf("Failed to load workspace: %v", ew.GetError())
            }
        }
    })
}

// BenchmarkFileIO measures just the file reading portion
func BenchmarkFileIO_Large(b *testing.B) {
    modPath := ensureGeneratedMod(b, "large")

    // List all .pp files
    files, _ := filepath.Glob(filepath.Join(modPath, "*.pp"))

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        for _, f := range files {
            _, err := os.ReadFile(f)
            if err != nil {
                b.Fatal(err)
            }
        }
    }
}

func ensureGeneratedMod(b *testing.B, size string) string {
    b.Helper()

    modPath := filepath.Join(testdataDir(), "mods", "generated", size)

    // Check if mod exists
    if _, err := os.Stat(filepath.Join(modPath, "mod.pp")); os.IsNotExist(err) {
        // Generate mod
        cmd := exec.Command("go", "run", "../../scripts/generate_test_mods.go", modPath, size)
        if err := cmd.Run(); err != nil {
            b.Skipf("Failed to generate test mod: %v", err)
        }
    }

    return modPath
}

func testdataDir() string {
    wd, _ := os.Getwd()
    return filepath.Join(wd, "..", "..", "testdata")
}
```

### 3. Create Payload Building Benchmarks

`internal/dashboardserver/payload_benchmark_test.go`:

```go
package dashboardserver_test

import (
    "context"
    "path/filepath"
    "testing"

    "github.com/turbot/powerpipe/internal/dashboardserver"
    "github.com/turbot/powerpipe/internal/workspace"
)

func BenchmarkBuildAvailableDashboardsPayload_Small(b *testing.B) {
    benchmarkPayload(b, "small")
}

func BenchmarkBuildAvailableDashboardsPayload_Medium(b *testing.B) {
    benchmarkPayload(b, "medium")
}

func BenchmarkBuildAvailableDashboardsPayload_Large(b *testing.B) {
    benchmarkPayload(b, "large")
}

func benchmarkPayload(b *testing.B, size string) {
    modPath := filepath.Join(testdataDir(), "mods", "generated", size)
    ctx := context.Background()

    w, ew := workspace.LoadWorkspacePromptingForVariables(ctx, modPath)
    if ew.GetError() != nil {
        b.Skipf("Failed to load workspace: %v", ew.GetError())
    }

    resources := w.GetPowerpipeModResources()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := dashboardserver.BuildAvailableDashboardsPayload(resources)
        if err != nil {
            b.Fatal(err)
        }
    }
}

// Benchmark JSON marshaling separately
func BenchmarkPayloadJSONMarshal_Large(b *testing.B) {
    // Similar setup, but measure JSON marshal time specifically
}

func testdataDir() string {
    // Implementation
}
```

### 4. Create Benchmark Runner Script

`scripts/run_benchmarks.sh`:

```bash
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

# Run benchmarks with memory profiling
go test -bench=. -benchmem -benchtime=5s \
    ./internal/workspace/... \
    ./internal/dashboardserver/... \
    -run=^$ \
    2>&1 | tee "$RESULTS_FILE"

# Parse results to JSON for comparison
go run ./scripts/parse_benchmark_results.go "$RESULTS_FILE" > "$JSON_FILE"

echo ""
echo "Benchmark complete. Results saved to:"
echo "  Text: $RESULTS_FILE"
echo "  JSON: $JSON_FILE"
```

### 5. Create Benchmark Comparison Script

`scripts/compare_benchmarks.go`:

```go
//go:build ignore

package main

import (
    "encoding/json"
    "fmt"
    "os"
)

type BenchmarkResult struct {
    Name       string  `json:"name"`
    NsPerOp    float64 `json:"ns_per_op"`
    BytesPerOp int64   `json:"bytes_per_op"`
    AllocsPerOp int64  `json:"allocs_per_op"`
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: go run compare_benchmarks.go <before.json> <after.json>")
        os.Exit(1)
    }

    before := loadResults(os.Args[1])
    after := loadResults(os.Args[2])

    fmt.Println("=== Benchmark Comparison ===")
    fmt.Printf("%-50s %15s %15s %10s\n", "Benchmark", "Before", "After", "Change")
    fmt.Println(strings.Repeat("-", 95))

    for name, b := range before {
        if a, ok := after[name]; ok {
            change := (a.NsPerOp - b.NsPerOp) / b.NsPerOp * 100
            fmt.Printf("%-50s %12.2fms %12.2fms %+9.1f%%\n",
                name,
                b.NsPerOp/1e6,
                a.NsPerOp/1e6,
                change)
        }
    }
}

func loadResults(path string) map[string]BenchmarkResult {
    // Implementation
}
```

## Acceptance Criteria

- [x] Mod generator script creates valid test mods
- [x] Generated mods for small/medium/large sizes exist
- [x] `BenchmarkLoadWorkspace_Small` runs and produces timing data
- [x] `BenchmarkLoadWorkspace_Medium` runs and produces timing data
- [x] `BenchmarkLoadWorkspace_Large` runs and produces timing data
- [x] `BenchmarkBuildAvailableDashboardsPayload_*` benchmarks work
- [x] Benchmark runner script produces consistent results
- [x] Benchmark comparison script can diff before/after results
- [x] Benchmark results include memory allocation data
- [x] Results are saved in both human-readable and JSON formats
- [x] Benchmarks complete in reasonable time (< 10 minutes total)

## Expected Benchmark Output

```
goos: darwin
goarch: arm64
pkg: github.com/turbot/powerpipe/internal/workspace
BenchmarkLoadWorkspace_Small-10           50     23456789 ns/op    12345678 B/op    123456 allocs/op
BenchmarkLoadWorkspace_Medium-10          20     56789012 ns/op    34567890 B/op    345678 allocs/op
BenchmarkLoadWorkspace_Large-10            5    234567890 ns/op   123456789 B/op   1234567 allocs/op

=== Timing Breakdown (Large) ===
LoadFileData:            45.23ms
ParseHclFiles:          123.45ms
Decoder.Decode:          56.78ms
verifyRuntimeDeps:       12.34ms
Total:                  237.80ms
```

## Notes

- Generated mods should not be committed to git (add to .gitignore)
- Benchmarks may need adjustment based on CI resources
- Consider adding CPU profiling (`-cpuprofile`) for deep analysis
- Memory benchmarks help identify allocation hotspots
