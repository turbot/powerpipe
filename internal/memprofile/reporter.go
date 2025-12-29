package memprofile

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// MemoryReport summarizes memory usage from tracked snapshots
type MemoryReport struct {
	Duration  time.Duration
	Snapshots []MemorySnapshot
}

// PeakHeapAlloc returns the maximum heap allocation across all snapshots
func (r *MemoryReport) PeakHeapAlloc() uint64 {
	var peak uint64
	for _, s := range r.Snapshots {
		if s.HeapAlloc > peak {
			peak = s.HeapAlloc
		}
	}
	return peak
}

// FinalHeapAlloc returns the heap allocation from the last snapshot
func (r *MemoryReport) FinalHeapAlloc() uint64 {
	if len(r.Snapshots) == 0 {
		return 0
	}
	return r.Snapshots[len(r.Snapshots)-1].HeapAlloc
}

// PeakHeapObjects returns the maximum number of heap objects across all snapshots
func (r *MemoryReport) PeakHeapObjects() uint64 {
	var peak uint64
	for _, s := range r.Snapshots {
		if s.HeapObjects > peak {
			peak = s.HeapObjects
		}
	}
	return peak
}

// TotalAllocated returns the total bytes allocated (cumulative)
func (r *MemoryReport) TotalAllocated() uint64 {
	if len(r.Snapshots) < 2 {
		return 0
	}
	first := r.Snapshots[0]
	last := r.Snapshots[len(r.Snapshots)-1]
	return last.TotalAlloc - first.TotalAlloc
}

// GCCycles returns the number of GC cycles that occurred
func (r *MemoryReport) GCCycles() uint32 {
	if len(r.Snapshots) < 2 {
		return 0
	}
	first := r.Snapshots[0]
	last := r.Snapshots[len(r.Snapshots)-1]
	return last.NumGC - first.NumGC
}

// String returns a human-readable summary of the memory report
func (r *MemoryReport) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Memory Report (duration: %v)\n", r.Duration))
	b.WriteString(fmt.Sprintf("Peak Heap: %s\n", FormatBytes(r.PeakHeapAlloc())))
	b.WriteString(fmt.Sprintf("Final Heap: %s\n", FormatBytes(r.FinalHeapAlloc())))
	b.WriteString(fmt.Sprintf("Peak Objects: %d\n", r.PeakHeapObjects()))
	b.WriteString(fmt.Sprintf("Total Allocated: %s\n", FormatBytes(r.TotalAllocated())))
	b.WriteString(fmt.Sprintf("GC Cycles: %d\n", r.GCCycles()))
	b.WriteString("\nSnapshots:\n")
	for _, s := range r.Snapshots {
		b.WriteString(fmt.Sprintf("  %s: heap=%s objects=%d\n",
			s.Label, FormatBytes(s.HeapAlloc), s.HeapObjects))
	}
	return b.String()
}

// Summary returns a brief one-line summary
func (r *MemoryReport) Summary() string {
	return fmt.Sprintf("peak=%s final=%s objects=%d gc=%d",
		FormatBytes(r.PeakHeapAlloc()),
		FormatBytes(r.FinalHeapAlloc()),
		r.PeakHeapObjects(),
		r.GCCycles())
}

// ReportJSON represents the JSON structure for memory reports
type ReportJSON struct {
	Duration       string         `json:"duration"`
	PeakHeapBytes  uint64         `json:"peak_heap_bytes"`
	FinalHeapBytes uint64         `json:"final_heap_bytes"`
	PeakHeapMB     float64        `json:"peak_heap_mb"`
	FinalHeapMB    float64        `json:"final_heap_mb"`
	PeakObjects    uint64         `json:"peak_objects"`
	TotalAllocated uint64         `json:"total_allocated_bytes"`
	GCCycles       uint32         `json:"gc_cycles"`
	Snapshots      []SnapshotJSON `json:"snapshots"`
}

// SnapshotJSON represents a snapshot in JSON format
type SnapshotJSON struct {
	Timestamp   string  `json:"timestamp"`
	Label       string  `json:"label"`
	HeapBytes   uint64  `json:"heap_bytes"`
	HeapMB      float64 `json:"heap_mb"`
	HeapObjects uint64  `json:"heap_objects"`
}

// JSON returns the report as a JSON string
func (r *MemoryReport) JSON() (string, error) {
	report := ReportJSON{
		Duration:       r.Duration.String(),
		PeakHeapBytes:  r.PeakHeapAlloc(),
		FinalHeapBytes: r.FinalHeapAlloc(),
		PeakHeapMB:     float64(r.PeakHeapAlloc()) / (1024 * 1024),
		FinalHeapMB:    float64(r.FinalHeapAlloc()) / (1024 * 1024),
		PeakObjects:    r.PeakHeapObjects(),
		TotalAllocated: r.TotalAllocated(),
		GCCycles:       r.GCCycles(),
		Snapshots:      make([]SnapshotJSON, len(r.Snapshots)),
	}

	for i, s := range r.Snapshots {
		report.Snapshots[i] = SnapshotJSON{
			Timestamp:   s.Timestamp.Format(time.RFC3339Nano),
			Label:       s.Label,
			HeapBytes:   s.HeapAlloc,
			HeapMB:      float64(s.HeapAlloc) / (1024 * 1024),
			HeapObjects: s.HeapObjects,
		}
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FormatBytes formats bytes into human-readable format
func FormatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

// ParseBytes parses a human-readable byte string back to uint64
func ParseBytes(s string) (uint64, error) {
	var value float64
	var unit string
	_, err := fmt.Sscanf(s, "%f %s", &value, &unit)
	if err != nil {
		// Try without space
		_, err = fmt.Sscanf(s, "%f%s", &value, &unit)
		if err != nil {
			return 0, fmt.Errorf("invalid byte format: %s", s)
		}
	}

	multipliers := map[string]uint64{
		"B":  1,
		"KB": 1024,
		"MB": 1024 * 1024,
		"GB": 1024 * 1024 * 1024,
		"TB": 1024 * 1024 * 1024 * 1024,
	}

	mult, ok := multipliers[strings.ToUpper(unit)]
	if !ok {
		return 0, fmt.Errorf("unknown unit: %s", unit)
	}

	return uint64(value * float64(mult)), nil
}

// CompareReports compares two memory reports and returns a comparison string
func CompareReports(before, after *MemoryReport, labels ...string) string {
	beforeLabel := "Before"
	afterLabel := "After"
	if len(labels) >= 2 {
		beforeLabel = labels[0]
		afterLabel = labels[1]
	}

	peakDelta := int64(after.PeakHeapAlloc()) - int64(before.PeakHeapAlloc())
	finalDelta := int64(after.FinalHeapAlloc()) - int64(before.FinalHeapAlloc())
	objectsDelta := int64(after.PeakHeapObjects()) - int64(before.PeakHeapObjects())

	peakPct := float64(peakDelta) / float64(before.PeakHeapAlloc()) * 100
	finalPct := float64(finalDelta) / float64(before.FinalHeapAlloc()) * 100

	var b strings.Builder
	b.WriteString("Memory Comparison\n")
	b.WriteString(fmt.Sprintf("%-20s %15s %15s %15s %10s\n", "", beforeLabel, afterLabel, "Delta", "Change"))
	b.WriteString(strings.Repeat("-", 75) + "\n")
	b.WriteString(fmt.Sprintf("%-20s %15s %15s %15s %9.1f%%\n",
		"Peak Heap",
		FormatBytes(before.PeakHeapAlloc()),
		FormatBytes(after.PeakHeapAlloc()),
		formatDelta(peakDelta),
		peakPct))
	b.WriteString(fmt.Sprintf("%-20s %15s %15s %15s %9.1f%%\n",
		"Final Heap",
		FormatBytes(before.FinalHeapAlloc()),
		FormatBytes(after.FinalHeapAlloc()),
		formatDelta(finalDelta),
		finalPct))
	b.WriteString(fmt.Sprintf("%-20s %15d %15d %15d\n",
		"Peak Objects",
		before.PeakHeapObjects(),
		after.PeakHeapObjects(),
		objectsDelta))

	return b.String()
}

func formatDelta(delta int64) string {
	if delta >= 0 {
		return "+" + FormatBytes(uint64(delta))
	}
	return "-" + FormatBytes(uint64(-delta))
}
