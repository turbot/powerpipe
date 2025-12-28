# Task 1: Add Performance Instrumentation

## Objective

Add comprehensive timing instrumentation to the mod loading flow to enable accurate performance measurement and identify bottlenecks.

## Context

- This is the foundation task - all subsequent performance work depends on it
- Existing `utils.LogTime()` calls are sparse and inconsistent
- Need structured, parseable timing output for benchmarking
- Instrumentation should be toggleable via environment variable

## Dependencies

### Files to Modify
- `internal/workspace/load_workspace.go` - Add timing to Load functions
- `internal/initialisation/init_data.go` - Add timing to Init flow
- `internal/dashboardserver/server.go` - Add timing to server startup
- `internal/dashboardserver/payload.go` - Add timing to payload building
- `internal/cmd/server.go` - Add overall startup timing

### New Files to Create
- `internal/timing/timing.go` - New timing utilities package

### External Dependencies
- None (pure Go implementation)

## Implementation Details

### 1. Create Timing Package

Create `internal/timing/timing.go`:
```go
package timing

import (
    "encoding/json"
    "fmt"
    "os"
    "sort"
    "sync"
    "time"
)

var (
    enabled  = os.Getenv("POWERPIPE_TIMING") != ""
    detailed = os.Getenv("POWERPIPE_TIMING") == "detailed"
    mu       sync.Mutex
    timings  []TimingEntry
)

type TimingEntry struct {
    Name      string        `json:"name"`
    Duration  time.Duration `json:"duration_ns"`
    DurationMs float64      `json:"duration_ms"`
    StartTime time.Time     `json:"start_time"`
    Context   string        `json:"context,omitempty"`
}

// Track returns a function to call when operation completes
func Track(name string, context ...string) func() {
    if !enabled {
        return func() {}
    }
    start := time.Now()
    ctx := ""
    if len(context) > 0 {
        ctx = context[0]
    }
    return func() {
        duration := time.Since(start)
        mu.Lock()
        timings = append(timings, TimingEntry{
            Name:       name,
            Duration:   duration,
            DurationMs: float64(duration.Nanoseconds()) / 1e6,
            StartTime:  start,
            Context:    ctx,
        })
        mu.Unlock()
        if detailed {
            fmt.Fprintf(os.Stderr, "[TIMING] %s: %.2fms\n", name, float64(duration.Nanoseconds())/1e6)
        }
    }
}

// Report outputs all collected timings
func Report() {
    if !enabled || len(timings) == 0 {
        return
    }
    mu.Lock()
    defer mu.Unlock()

    // Sort by start time
    sort.Slice(timings, func(i, j int) bool {
        return timings[i].StartTime.Before(timings[j].StartTime)
    })

    fmt.Fprintln(os.Stderr, "\n=== Performance Timing Report ===")
    var total time.Duration
    for _, t := range timings {
        fmt.Fprintf(os.Stderr, "%-50s %10.2fms\n", t.Name, t.DurationMs)
        total += t.Duration
    }
    fmt.Fprintf(os.Stderr, "%-50s %10.2fms\n", "TOTAL (sum)", float64(total.Nanoseconds())/1e6)
    fmt.Fprintln(os.Stderr, "=================================")
}

// ReportJSON outputs timings as JSON for programmatic processing
func ReportJSON() string {
    if !enabled || len(timings) == 0 {
        return "{}"
    }
    mu.Lock()
    defer mu.Unlock()

    data, _ := json.MarshalIndent(timings, "", "  ")
    return string(data)
}

// Reset clears collected timings
func Reset() {
    mu.Lock()
    timings = nil
    mu.Unlock()
}

// IsEnabled returns whether timing is enabled
func IsEnabled() bool {
    return enabled
}
```

### 2. Add Instrumentation Points

#### `internal/workspace/load_workspace.go`
```go
import "github.com/turbot/powerpipe/internal/timing"

func LoadWorkspacePromptingForVariables(...) {
    defer timing.Track("LoadWorkspacePromptingForVariables")()
    // existing code...
}

func Load(...) {
    defer timing.Track("Load")()
    // Add around key operations:
    defer timing.Track("SetModfileExists")()
    // ...
    defer timing.Track("LoadExclusions")()
    // ...
    defer timing.Track("LoadWorkspaceMod")()
    // ...
    defer timing.Track("verifyResourceRuntimeDependencies")()
}
```

#### `internal/initialisation/init_data.go`
```go
func NewInitData[T](...) {
    defer timing.Track("NewInitData")()
    // ...
}

func (i *InitData) Init(...) {
    defer timing.Track("InitData.Init")()

    defer timing.Track("telemetry.Init")()
    // ...
    defer timing.Track("modinstaller.InstallWorkspaceDependencies")()
    // ...
    defer timing.Track("db_client.GetDefaultDatabaseConfig")()
    // ...
    defer timing.Track("db_client.NewDbClient")()
    // ...
    defer timing.Track("validateModRequirementsRecursively")()
    // ...
    defer timing.Track("NewDashboardExecutor")()
}
```

#### `internal/dashboardserver/server.go`
```go
func NewServer(...) {
    defer timing.Track("dashboardserver.NewServer")()
    // ...
}

func (s *Server) InitAsync(...) {
    defer timing.Track("Server.InitAsync")()
    // ...
}
```

#### `internal/dashboardserver/payload.go`
```go
func buildAvailableDashboardsPayload(...) {
    defer timing.Track("buildAvailableDashboardsPayload")()
    // ...
}

func (s *Server) buildServerMetadataPayload(...) {
    defer timing.Track("buildServerMetadataPayload")()
    // ...
}
```

#### `internal/cmd/server.go`
```go
func runServerCmd(...) {
    overallStart := time.Now()
    defer func() {
        timing.Track("runServerCmd.total")()
        timing.Report()
    }()
    // ... at end before <-ctx.Done()
}
```

### 3. Environment Variables

- `POWERPIPE_TIMING=1` - Enable timing, output summary at end
- `POWERPIPE_TIMING=detailed` - Enable timing, output each measurement as it happens
- `POWERPIPE_TIMING=json` - Output JSON format for programmatic parsing

## Acceptance Criteria

- [ ] `timing` package created with Track, Report, ReportJSON functions
- [ ] All key mod loading functions instrumented
- [ ] All key initialization functions instrumented
- [ ] All key payload building functions instrumented
- [ ] Environment variable control works correctly
- [ ] Timing output is parseable for benchmark tests
- [ ] No performance impact when timing is disabled
- [ ] Unit tests for timing package
- [ ] Example output documented

## Expected Output Example

```
=== Performance Timing Report ===
LoadWorkspacePromptingForVariables                      1523.45ms
  SetModfileExists                                         1.23ms
  LoadExclusions                                           0.45ms
  LoadWorkspaceMod                                      1450.67ms
  verifyResourceRuntimeDependencies                       70.10ms
InitData.Init                                            345.67ms
  telemetry.Init                                           5.23ms
  db_client.NewDbClient                                  289.45ms
  validateModRequirementsRecursively                      45.67ms
dashboardserver.NewServer                                 12.34ms
buildAvailableDashboardsPayload                           45.67ms
runServerCmd.total                                      1932.45ms
=================================
```

## Notes

- Keep timing overhead minimal (nanosecond-level when disabled)
- Use `sync.Mutex` for thread safety
- Consider adding memory profiling hooks for future use
- Timing package should be reusable across Powerpipe
