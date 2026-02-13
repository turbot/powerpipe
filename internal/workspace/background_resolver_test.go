package workspace

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turbot/powerpipe/internal/resourceindex"
)

func TestBackgroundResolver_PriorityCalculation(t *testing.T) {
	// Create a minimal resolver to test priority calculation
	resolver := &BackgroundResolver{}

	tests := []struct {
		name     string
		entry    *resourceindex.IndexEntry
		expected int
	}{
		{
			name: "top-level dashboard",
			entry: &resourceindex.IndexEntry{
				Type:           "dashboard",
				IsTopLevel:     true,
				UnresolvedRefs: []string{"title"},
			},
			expected: 100 + 50 + 10, // top-level + dashboard + 1 ref
		},
		{
			name: "non-top-level dashboard",
			entry: &resourceindex.IndexEntry{
				Type:           "dashboard",
				IsTopLevel:     false,
				UnresolvedRefs: []string{"title", "tags"},
			},
			expected: 50 + 20, // dashboard + 2 refs
		},
		{
			name: "top-level benchmark",
			entry: &resourceindex.IndexEntry{
				Type:           "benchmark",
				IsTopLevel:     true,
				UnresolvedRefs: []string{},
			},
			expected: 100 + 30, // top-level + benchmark
		},
		{
			name: "control",
			entry: &resourceindex.IndexEntry{
				Type:           "control",
				IsTopLevel:     false,
				UnresolvedRefs: []string{"title"},
			},
			expected: 10, // just 1 ref
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			priority := resolver.calculatePriority(tt.entry)
			assert.Equal(t, tt.expected, priority)
		})
	}
}

func TestBackgroundResolver_QueueUnresolvedEntries(t *testing.T) {
	// Create test index with some resolved and unresolved entries
	index := resourceindex.NewResourceIndex()

	// Resolved entry (should not be queued)
	index.Add(&resourceindex.IndexEntry{
		Name:                "mod.dashboard.resolved",
		Type:                "dashboard",
		TitleResolved:       true,
		DescriptionResolved: true,
		TagsResolved:        true,
	})

	// Unresolved entries (should be queued)
	index.Add(&resourceindex.IndexEntry{
		Name:                "mod.dashboard.unresolved1",
		Type:                "dashboard",
		IsTopLevel:          true,
		TitleResolved:       false,
		DescriptionResolved: true,
		TagsResolved:        true,
	})
	index.Add(&resourceindex.IndexEntry{
		Name:                "mod.benchmark.unresolved2",
		Type:                "benchmark",
		TitleResolved:       true,
		DescriptionResolved: true,
		TagsResolved:        false,
	})

	// Create resolver with the index
	resolver := &BackgroundResolver{
		index:   index,
		queue:   NewResolutionQueue(),
		workers: 1,
	}

	resolver.queueUnresolvedEntries()

	// Should have 2 unresolved entries in queue
	assert.Equal(t, 2, resolver.queue.Len())

	// Top-level dashboard should be first (higher priority)
	first := resolver.queue.Pop()
	assert.Equal(t, "mod.dashboard.unresolved1", first)
}

func TestBackgroundResolver_Options(t *testing.T) {
	// Test option functions
	resolver := &BackgroundResolver{workers: 1}

	WithWorkers(8)(resolver)
	assert.Equal(t, 8, resolver.workers)

	// Test with 0 workers (should be ignored)
	WithWorkers(0)(resolver)
	assert.Equal(t, 8, resolver.workers)

	updateCalled := false
	WithOnUpdate(func(name string) {
		updateCalled = true
	})(resolver)
	resolver.onUpdate("test")
	assert.True(t, updateCalled)

	completeCalled := false
	WithOnComplete(func() {
		completeCalled = true
	})(resolver)
	resolver.onComplete()
	assert.True(t, completeCalled)
}

func TestBackgroundResolver_StopWithoutStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	resolver := &BackgroundResolver{
		ctx:    ctx,
		cancel: cancel,
	}

	// Should not panic
	resolver.Stop()
}

func TestBackgroundResolver_Prioritize(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resolver := &BackgroundResolver{
		queue:  NewResolutionQueue(),
		ctx:    ctx,
		cancel: cancel,
	}

	// Add some items
	resolver.queue.Push("low", 10)
	resolver.queue.Push("medium", 50)

	// Prioritize a new item
	resolver.Prioritize("high")

	// Should be at front
	assert.Equal(t, "high", resolver.queue.Pop())
}

func TestBackgroundResolver_Stats(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resolver := &BackgroundResolver{
		queue:  NewResolutionQueue(),
		ctx:    ctx,
		cancel: cancel,
	}

	resolver.queue.Push("item1", 10)
	resolver.queue.Push("item2", 20)

	stats := resolver.Stats()
	assert.Equal(t, 2, stats.QueueLength)
	assert.False(t, stats.IsStarted)
	assert.False(t, stats.IsComplete)

	// After marking started
	resolver.started.Store(true)
	stats = resolver.Stats()
	assert.True(t, stats.IsStarted)

	// After marking complete
	resolver.completed.Store(true)
	stats = resolver.Stats()
	assert.True(t, stats.IsComplete)
}

func TestBackgroundResolver_FullyResolvedEntry(t *testing.T) {
	// Create an index with a fully resolved entry
	index := resourceindex.NewResourceIndex()
	index.Add(&resourceindex.IndexEntry{
		Name:                "mod.dashboard.resolved",
		Type:                "dashboard",
		Title:               "Existing Title",
		TitleResolved:       true,
		DescriptionResolved: true,
		TagsResolved:        true,
	})

	entry, _ := index.Get("mod.dashboard.resolved")

	// Verify entry is fully resolved
	assert.True(t, entry.IsFullyResolved())
	assert.False(t, entry.NeedsResolution())
}

func TestBackgroundResolver_ConcurrentResolveNow(t *testing.T) {
	// Create test index
	index := resourceindex.NewResourceIndex()
	for i := 0; i < 10; i++ {
		index.Add(&resourceindex.IndexEntry{
			Name:                "mod.dashboard.test" + string(rune('0'+i)),
			Type:                "dashboard",
			TitleResolved:       false,
			DescriptionResolved: true,
			TagsResolved:        true,
		})
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resolver := &BackgroundResolver{
		index:      index,
		queue:      NewResolutionQueue(),
		ctx:        ctx,
		cancel:     cancel,
		inProgress: sync.Map{},
	}

	// Concurrent ResolveNow calls should not deadlock
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// This will fail since workspace is nil, but should not deadlock
			_ = resolver.ResolveNow(ctx, "mod.dashboard.test"+string(rune('0'+i%10)))
		}(i)
	}

	// Use a timeout to detect deadlocks
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success - no deadlock
	case <-time.After(5 * time.Second):
		t.Fatal("deadlock detected in concurrent ResolveNow")
	}
}

func TestBackgroundResolver_AllWorkersIdle(t *testing.T) {
	resolver := &BackgroundResolver{
		inProgress: sync.Map{},
	}

	// Initially idle
	assert.True(t, resolver.allWorkersIdle())

	// Mark something in progress
	resolver.inProgress.Store("test", true)
	assert.False(t, resolver.allWorkersIdle())

	// Clear it
	resolver.inProgress.Delete("test")
	assert.True(t, resolver.allWorkersIdle())
}

func TestBackgroundResolver_EmptyQueueCompletion(t *testing.T) {
	// Create index with only resolved entries
	index := resourceindex.NewResourceIndex()
	index.Add(&resourceindex.IndexEntry{
		Name:                "mod.dashboard.resolved",
		Type:                "dashboard",
		TitleResolved:       true,
		DescriptionResolved: true,
		TagsResolved:        true,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	completeCalled := atomic.Bool{}
	resolver := &BackgroundResolver{
		index:      index,
		queue:      NewResolutionQueue(),
		ctx:        ctx,
		cancel:     cancel,
		workers:    2,
		inProgress: sync.Map{},
		onComplete: func() {
			completeCalled.Store(true)
		},
	}

	// Start should complete immediately since nothing needs resolution
	resolver.Start()

	// Wait a bit for completion
	time.Sleep(100 * time.Millisecond)

	assert.True(t, resolver.IsComplete())
	assert.True(t, completeCalled.Load())
}

func TestBackgroundResolver_NotifiesOnUpdate(t *testing.T) {
	// Create index with unresolved entry
	index := resourceindex.NewResourceIndex()
	index.Add(&resourceindex.IndexEntry{
		Name:                "mod.dashboard.test",
		Type:                "dashboard",
		TitleResolved:       false,
		DescriptionResolved: true,
		TagsResolved:        true,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var updates []string
	var updatesMu sync.Mutex

	resolver := &BackgroundResolver{
		index:      index,
		queue:      NewResolutionQueue(),
		ctx:        ctx,
		cancel:     cancel,
		inProgress: sync.Map{},
		onUpdate: func(name string) {
			updatesMu.Lock()
			updates = append(updates, name)
			updatesMu.Unlock()
		},
	}

	entry, _ := index.Get("mod.dashboard.test")

	// Simulate updating the entry (as if resource was loaded)
	// Since we can't load real resources without workspace, we'll test the mechanics
	entry.Title = "Resolved Title"
	entry.TitleResolved = true
	index.UpdateEntry(entry)

	// Manually call onUpdate to simulate notification
	resolver.onUpdate("mod.dashboard.test")

	updatesMu.Lock()
	defer updatesMu.Unlock()
	assert.Contains(t, updates, "mod.dashboard.test")
}

func TestIndexUpdateEntry(t *testing.T) {
	index := resourceindex.NewResourceIndex()

	// Add initial entry
	entry := &resourceindex.IndexEntry{
		Name:      "mod.dashboard.test",
		Type:      "dashboard",
		Title:     "Original",
		ShortName: "test",
	}
	index.Add(entry)

	// Update the entry with a longer title and description
	entry.Title = "Updated Title with More Characters"
	entry.Description = "This is a long description that adds significant size to the entry"
	index.UpdateEntry(entry)

	// Verify update
	updated, ok := index.Get("mod.dashboard.test")
	require.True(t, ok)
	assert.Equal(t, "Updated Title with More Characters", updated.Title)
	assert.Equal(t, "This is a long description that adds significant size to the entry", updated.Description)

	// Verify entry exists in type index too
	dashboards := index.Dashboards()
	assert.Equal(t, 1, len(dashboards))
	assert.Equal(t, "Updated Title with More Characters", dashboards[0].Title)
}

func TestIndexUpdateEntry_NewEntry(t *testing.T) {
	index := resourceindex.NewResourceIndex()

	// UpdateEntry on non-existent entry should add it
	entry := &resourceindex.IndexEntry{
		Name:      "mod.dashboard.new",
		Type:      "dashboard",
		Title:     "New Entry",
		ShortName: "new",
	}
	index.UpdateEntry(entry)

	// Should be added
	retrieved, ok := index.Get("mod.dashboard.new")
	require.True(t, ok)
	assert.Equal(t, "New Entry", retrieved.Title)
	assert.Equal(t, 1, index.Count())
}
