package memprofile

import (
	"context"
	"sync"
	"time"
)

// ContinuousTracker samples memory at regular intervals
type ContinuousTracker struct {
	interval  time.Duration
	snapshots []MemorySnapshot
	mu        sync.Mutex
	cancel    context.CancelFunc
	done      chan struct{}
	running   bool
}

// NewContinuousTracker creates a new continuous memory tracker
// that samples memory at the given interval
func NewContinuousTracker(interval time.Duration) *ContinuousTracker {
	return &ContinuousTracker{
		interval: interval,
		done:     make(chan struct{}),
	}
}

// Start begins continuous memory tracking in a background goroutine
func (t *ContinuousTracker) Start(ctx context.Context) {
	t.mu.Lock()
	if t.running {
		t.mu.Unlock()
		return
	}
	t.running = true
	t.snapshots = nil // Clear any previous snapshots
	t.done = make(chan struct{})
	t.mu.Unlock()

	ctx, t.cancel = context.WithCancel(ctx)

	go func() {
		defer close(t.done)
		ticker := time.NewTicker(t.interval)
		defer ticker.Stop()

		// Take initial snapshot
		t.addSnapshot(TakeSnapshot("start"))

		for {
			select {
			case <-ctx.Done():
				t.addSnapshot(TakeSnapshot("stop"))
				return
			case <-ticker.C:
				t.addSnapshot(TakeSnapshot(time.Now().Format("15:04:05.000")))
			}
		}
	}()
}

func (t *ContinuousTracker) addSnapshot(s MemorySnapshot) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.snapshots = append(t.snapshots, s)
}

// Stop stops continuous tracking and returns a memory report
func (t *ContinuousTracker) Stop() *MemoryReport {
	t.mu.Lock()
	if !t.running {
		t.mu.Unlock()
		return &MemoryReport{Snapshots: t.snapshots}
	}
	t.running = false
	t.mu.Unlock()

	if t.cancel != nil {
		t.cancel()
		<-t.done
	}

	t.mu.Lock()
	snapshots := make([]MemorySnapshot, len(t.snapshots))
	copy(snapshots, t.snapshots)
	t.mu.Unlock()

	var duration time.Duration
	if len(snapshots) >= 2 {
		duration = snapshots[len(snapshots)-1].Timestamp.Sub(snapshots[0].Timestamp)
	}

	return &MemoryReport{
		Duration:  duration,
		Snapshots: snapshots,
	}
}

// Snapshots returns the current snapshots without stopping tracking
func (t *ContinuousTracker) Snapshots() []MemorySnapshot {
	t.mu.Lock()
	defer t.mu.Unlock()
	snapshots := make([]MemorySnapshot, len(t.snapshots))
	copy(snapshots, t.snapshots)
	return snapshots
}

// IsRunning returns true if the tracker is currently running
func (t *ContinuousTracker) IsRunning() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.running
}

// MarkEvent adds a labeled snapshot at the current time
// This is useful for marking significant events during tracking
func (t *ContinuousTracker) MarkEvent(label string) {
	t.addSnapshot(TakeSnapshot(label))
}

// MemoryWatcher provides a simple API for watching memory thresholds
type MemoryWatcher struct {
	threshold   uint64
	interval    time.Duration
	onThreshold func(MemorySnapshot)
	cancel      context.CancelFunc
	done        chan struct{}
	running     bool
	mu          sync.Mutex
}

// NewMemoryWatcher creates a watcher that calls the callback when
// heap allocation exceeds the threshold
func NewMemoryWatcher(threshold uint64, interval time.Duration, callback func(MemorySnapshot)) *MemoryWatcher {
	return &MemoryWatcher{
		threshold:   threshold,
		interval:    interval,
		onThreshold: callback,
		done:        make(chan struct{}),
	}
}

// Start begins watching memory in a background goroutine
func (w *MemoryWatcher) Start(ctx context.Context) {
	w.mu.Lock()
	if w.running {
		w.mu.Unlock()
		return
	}
	w.running = true
	w.done = make(chan struct{})
	w.mu.Unlock()

	ctx, w.cancel = context.WithCancel(ctx)

	go func() {
		defer close(w.done)
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				snapshot := TakeSnapshot("watch")
				if snapshot.HeapAlloc > w.threshold && w.onThreshold != nil {
					w.onThreshold(snapshot)
				}
			}
		}
	}()
}

// Stop stops the memory watcher
func (w *MemoryWatcher) Stop() {
	w.mu.Lock()
	if !w.running {
		w.mu.Unlock()
		return
	}
	w.running = false
	w.mu.Unlock()

	if w.cancel != nil {
		w.cancel()
		<-w.done
	}
}
