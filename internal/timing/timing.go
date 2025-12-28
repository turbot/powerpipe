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
	jsonMode = os.Getenv("POWERPIPE_TIMING") == "json"
	mu       sync.Mutex
	timings  []TimingEntry
)

// TimingEntry represents a single timing measurement
type TimingEntry struct {
	Name       string        `json:"name"`
	Duration   time.Duration `json:"duration_ns"`
	DurationMs float64       `json:"duration_ms"`
	StartTime  time.Time     `json:"start_time"`
	Context    string        `json:"context,omitempty"`
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

	if jsonMode {
		fmt.Fprintln(os.Stderr, ReportJSON())
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

	// Sort by start time before output
	sort.Slice(timings, func(i, j int) bool {
		return timings[i].StartTime.Before(timings[j].StartTime)
	})

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

// GetTimings returns a copy of the collected timings
func GetTimings() []TimingEntry {
	mu.Lock()
	defer mu.Unlock()
	result := make([]TimingEntry, len(timings))
	copy(result, timings)
	return result
}
