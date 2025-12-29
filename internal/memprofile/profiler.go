// Package memprofile provides memory profiling utilities for tracking and analyzing
// memory usage during workspace loading and resource access.
package memprofile

import (
	"runtime"
	"time"
)

// MemorySnapshot captures memory state at a point in time
type MemorySnapshot struct {
	Timestamp    time.Time
	Label        string
	HeapAlloc    uint64 // Bytes allocated on heap
	HeapInuse    uint64 // Bytes in use by heap
	HeapObjects  uint64 // Number of allocated objects
	TotalAlloc   uint64 // Cumulative bytes allocated
	NumGC        uint32 // Number of GC cycles
	GCPauseTotal uint64 // Total GC pause time (ns)
}

// TakeSnapshot captures current memory state
func TakeSnapshot(label string) MemorySnapshot {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return MemorySnapshot{
		Timestamp:    time.Now(),
		Label:        label,
		HeapAlloc:    m.HeapAlloc,
		HeapInuse:    m.HeapInuse,
		HeapObjects:  m.HeapObjects,
		TotalAlloc:   m.TotalAlloc,
		NumGC:        m.NumGC,
		GCPauseTotal: m.PauseTotalNs,
	}
}

// ForceGCAndSnapshot runs garbage collection and returns snapshot
func ForceGCAndSnapshot(label string) MemorySnapshot {
	runtime.GC()
	runtime.GC() // Run twice to ensure finalization
	return TakeSnapshot(label)
}

// MemoryDelta calculates the difference between two snapshots
func (s MemorySnapshot) Delta(before MemorySnapshot) MemoryDelta {
	return MemoryDelta{
		HeapAllocDelta:   int64(s.HeapAlloc) - int64(before.HeapAlloc),
		HeapInuseDelta:   int64(s.HeapInuse) - int64(before.HeapInuse),
		HeapObjectsDelta: int64(s.HeapObjects) - int64(before.HeapObjects),
		TotalAllocDelta:  s.TotalAlloc - before.TotalAlloc,
		GCCyclesDelta:    s.NumGC - before.NumGC,
		Duration:         s.Timestamp.Sub(before.Timestamp),
	}
}

// MemoryDelta represents the change between two memory snapshots
type MemoryDelta struct {
	HeapAllocDelta   int64
	HeapInuseDelta   int64
	HeapObjectsDelta int64
	TotalAllocDelta  uint64
	GCCyclesDelta    uint32
	Duration         time.Duration
}

// MemoryTracker tracks memory over time with labeled snapshots
type MemoryTracker struct {
	snapshots []MemorySnapshot
	start     time.Time
}

// NewMemoryTracker creates a new memory tracker
func NewMemoryTracker() *MemoryTracker {
	return &MemoryTracker{
		start: time.Now(),
	}
}

// Snapshot captures a memory snapshot with the given label
func (t *MemoryTracker) Snapshot(label string) {
	t.snapshots = append(t.snapshots, TakeSnapshot(label))
}

// SnapshotAfterGC forces GC and captures a memory snapshot
func (t *MemoryTracker) SnapshotAfterGC(label string) {
	t.snapshots = append(t.snapshots, ForceGCAndSnapshot(label))
}

// Snapshots returns all captured snapshots
func (t *MemoryTracker) Snapshots() []MemorySnapshot {
	return t.snapshots
}

// Report generates a memory report from the tracked snapshots
func (t *MemoryTracker) Report() *MemoryReport {
	return &MemoryReport{
		Duration:  time.Since(t.start),
		Snapshots: t.snapshots,
	}
}

// Clear removes all snapshots and resets the start time
func (t *MemoryTracker) Clear() {
	t.snapshots = nil
	t.start = time.Now()
}
