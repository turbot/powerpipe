package workspace

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/powerpipe/internal/resourceindex"
	"github.com/turbot/powerpipe/internal/resources"
)

// UpdateListener is notified when resources are updated during background resolution.
type UpdateListener interface {
	OnResourceUpdated(resourceName string)
	OnResolutionComplete()
}

// BackgroundResolver handles progressive resolution of resource metadata.
// It processes resources that need resolution (variables, templates) in the background,
// updating the index as resolution completes.
type BackgroundResolver struct {
	workspace *LazyWorkspace
	index     *resourceindex.ResourceIndex

	// Resolution state
	queue      *ResolutionQueue
	inProgress sync.Map // map[string]bool - tracks resources currently being resolved

	// Notification callbacks
	onUpdate   func(resourceName string)
	onComplete func()

	// Control
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	workers   int
	started   atomic.Bool
	completed atomic.Bool
}

// ResolverOption configures a BackgroundResolver.
type ResolverOption func(*BackgroundResolver)

// WithWorkers sets the number of worker goroutines.
func WithWorkers(n int) ResolverOption {
	return func(r *BackgroundResolver) {
		if n > 0 {
			r.workers = n
		}
	}
}

// WithOnUpdate sets the callback for resource updates.
func WithOnUpdate(fn func(resourceName string)) ResolverOption {
	return func(r *BackgroundResolver) {
		r.onUpdate = fn
	}
}

// WithOnComplete sets the callback for resolution completion.
func WithOnComplete(fn func()) ResolverOption {
	return func(r *BackgroundResolver) {
		r.onComplete = fn
	}
}

// NewBackgroundResolver creates a new background resolver.
func NewBackgroundResolver(ws *LazyWorkspace, opts ...ResolverOption) *BackgroundResolver {
	ctx, cancel := context.WithCancel(context.Background())

	r := &BackgroundResolver{
		workspace: ws,
		index:     ws.GetIndex(),
		queue:     NewResolutionQueue(),
		ctx:       ctx,
		cancel:    cancel,
		workers:   4, // Default worker count
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// Start begins background resolution.
// It queues all entries that need resolution and starts worker goroutines.
func (r *BackgroundResolver) Start() {
	if r.started.Swap(true) {
		return // Already started
	}

	// Queue all entries that need resolution
	r.queueUnresolvedEntries()

	// If nothing to resolve, we're already done
	if r.queue.IsEmpty() {
		r.completed.Store(true)
		if r.onComplete != nil {
			r.onComplete()
		}
		return
	}

	// Start worker goroutines
	for i := 0; i < r.workers; i++ {
		r.wg.Add(1)
		go r.worker(i)
	}

	// Start completion monitor
	r.wg.Add(1)
	go r.completionMonitor()
}

// Stop gracefully stops background resolution.
func (r *BackgroundResolver) Stop() {
	r.cancel()
	r.wg.Wait()
}

// IsComplete returns true if all resolution is complete.
func (r *BackgroundResolver) IsComplete() bool {
	return r.completed.Load()
}

// WaitForComplete blocks until resolution is complete or timeout is reached.
// Returns true if resolution completed, false if timeout was reached.
func (r *BackgroundResolver) WaitForComplete(timeout time.Duration) bool {
	if r.completed.Load() {
		return true
	}

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if r.completed.Load() {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return r.completed.Load()
}

// Prioritize moves a resource to higher priority in the resolution queue.
// If the resource is not in the queue, it's added with high priority.
func (r *BackgroundResolver) Prioritize(resourceName string) {
	r.queue.Prioritize(resourceName, 1000) // Very high priority
}

// ResolveNow immediately resolves a resource, bypassing the queue.
// This is useful for on-demand resolution when user clicks a dashboard.
func (r *BackgroundResolver) ResolveNow(ctx context.Context, resourceName string) error {
	entry, ok := r.index.Get(resourceName)
	if !ok {
		return nil // Resource not found - nothing to resolve
	}

	// If already fully resolved, nothing to do
	if entry.IsFullyResolved() {
		return nil
	}

	// Mark as in progress
	if _, loaded := r.inProgress.LoadOrStore(resourceName, true); loaded {
		// Already being resolved by another goroutine, wait a bit and return
		time.Sleep(50 * time.Millisecond)
		return nil
	}
	defer r.inProgress.Delete(resourceName)

	// Resolve the resource
	r.resolveEntry(ctx, entry)
	return nil
}

// queueUnresolvedEntries adds all entries needing resolution to the queue.
func (r *BackgroundResolver) queueUnresolvedEntries() {
	for _, entry := range r.index.List() {
		if entry.NeedsResolution() {
			priority := r.calculatePriority(entry)
			r.queue.Push(entry.Name, priority)
		}
	}
}

// calculatePriority determines resolution priority.
// Higher priority = resolved sooner.
func (r *BackgroundResolver) calculatePriority(entry *resourceindex.IndexEntry) int {
	priority := 0

	// Top-level resources have higher priority (shown in list)
	if entry.IsTopLevel {
		priority += 100
	}

	// Dashboards higher than benchmarks (more commonly browsed)
	if entry.Type == "dashboard" {
		priority += 50
	} else if entry.Type == "benchmark" || entry.Type == "detection_benchmark" {
		priority += 30
	}

	// Resources with more unresolved fields have higher priority
	priority += len(entry.UnresolvedRefs) * 10

	return priority
}

// worker processes resolution requests from the queue.
func (r *BackgroundResolver) worker(id int) {
	defer r.wg.Done()

	for {
		select {
		case <-r.ctx.Done():
			return
		default:
			resourceName := r.queue.Pop()
			if resourceName == "" {
				// Queue empty, check if we should exit
				select {
				case <-r.ctx.Done():
					return
				case <-time.After(50 * time.Millisecond):
					// Check queue again
					continue
				}
			}

			r.resolveResource(resourceName)
		}
	}
}

// resolveResource resolves a single resource's metadata.
func (r *BackgroundResolver) resolveResource(name string) {
	// Mark as in progress
	if _, loaded := r.inProgress.LoadOrStore(name, true); loaded {
		return // Already being resolved
	}
	defer r.inProgress.Delete(name)

	entry, ok := r.index.Get(name)
	if !ok {
		return
	}

	r.resolveEntry(r.ctx, entry)
}

// resolveEntry resolves metadata for an index entry by loading the full resource.
func (r *BackgroundResolver) resolveEntry(ctx context.Context, entry *resourceindex.IndexEntry) {
	// Skip if workspace or loader is nil (e.g., in tests)
	if r.workspace == nil || r.workspace.loader == nil {
		return
	}

	// Load the full resource
	resource, err := r.workspace.loader.Load(ctx, entry.Name)
	if err != nil {
		slog.Debug("background resolution failed", "resource", entry.Name, "error", err)
		return
	}

	// Extract metadata based on resource type and update the entry
	updated := r.updateEntryFromResource(entry, resource)

	if updated {
		// Notify listeners
		if r.onUpdate != nil {
			r.onUpdate(entry.Name)
		}
	}
}

// updateEntryFromResource extracts resolved metadata from a loaded resource
// and updates the index entry. Returns true if any updates were made.
func (r *BackgroundResolver) updateEntryFromResource(entry *resourceindex.IndexEntry, resource modconfig.HclResource) bool {
	updated := false

	// Try to extract title
	if !entry.TitleResolved {
		if title := r.extractTitle(resource); title != "" {
			entry.Title = title
			entry.TitleResolved = true
			updated = true
		} else {
			// Even if empty, mark as resolved (no title available)
			entry.TitleResolved = true
		}
	}

	// Try to extract description
	if !entry.DescriptionResolved {
		if desc := r.extractDescription(resource); desc != "" {
			entry.Description = desc
			entry.DescriptionResolved = true
			updated = true
		} else {
			entry.DescriptionResolved = true
		}
	}

	// Try to extract tags
	if !entry.TagsResolved {
		tags := r.extractTags(resource)
		if tags != nil {
			entry.Tags = tags
			entry.TagsResolved = true
			entry.UnresolvedRefs = nil
			updated = true
		} else {
			// No tags, but mark as resolved
			entry.Tags = make(map[string]string)
			entry.TagsResolved = true
			entry.UnresolvedRefs = nil
		}
	}

	// Update the index with the modified entry
	if updated {
		r.index.UpdateEntry(entry)
	}

	return updated
}

// extractTitle extracts the title from a resource.
func (r *BackgroundResolver) extractTitle(resource modconfig.HclResource) string {
	switch res := resource.(type) {
	case *resources.Dashboard:
		if res.Title != nil {
			return *res.Title
		}
	case *resources.Benchmark:
		if res.Title != nil {
			return *res.Title
		}
	case *resources.DetectionBenchmark:
		if res.Title != nil {
			return *res.Title
		}
	case *resources.Control:
		if res.Title != nil {
			return *res.Title
		}
	case *resources.Query:
		if res.Title != nil {
			return *res.Title
		}
	}
	return ""
}

// extractDescription extracts the description from a resource.
func (r *BackgroundResolver) extractDescription(resource modconfig.HclResource) string {
	switch res := resource.(type) {
	case *resources.Dashboard:
		if res.Description != nil {
			return *res.Description
		}
	case *resources.Benchmark:
		if res.Description != nil {
			return *res.Description
		}
	case *resources.DetectionBenchmark:
		if res.Description != nil {
			return *res.Description
		}
	case *resources.Control:
		if res.Description != nil {
			return *res.Description
		}
	case *resources.Query:
		if res.Description != nil {
			return *res.Description
		}
	}
	return ""
}

// extractTags extracts tags from a resource.
func (r *BackgroundResolver) extractTags(resource modconfig.HclResource) map[string]string {
	switch res := resource.(type) {
	case *resources.Dashboard:
		return res.Tags
	case *resources.Benchmark:
		return res.Tags
	case *resources.DetectionBenchmark:
		return res.Tags
	case *resources.Control:
		return res.Tags
	case *resources.Query:
		return res.Tags
	}
	return nil
}

// completionMonitor watches for resolution completion.
func (r *BackgroundResolver) completionMonitor() {
	defer r.wg.Done()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-r.ctx.Done():
			return
		case <-ticker.C:
			// Check if queue is empty and nothing is in progress
			if r.queue.IsEmpty() && r.allWorkersIdle() {
				if r.completed.Swap(true) {
					return // Already completed
				}
				if r.onComplete != nil {
					r.onComplete()
				}
				slog.Debug("background resolution complete")
				return
			}
		}
	}
}

// allWorkersIdle returns true if no resources are currently being resolved.
func (r *BackgroundResolver) allWorkersIdle() bool {
	idle := true
	r.inProgress.Range(func(key, value interface{}) bool {
		idle = false
		return false // Stop iteration
	})
	return idle
}

// Stats returns statistics about background resolution.
func (r *BackgroundResolver) Stats() BackgroundResolverStats {
	return BackgroundResolverStats{
		QueueLength: r.queue.Len(),
		IsComplete:  r.completed.Load(),
		IsStarted:   r.started.Load(),
	}
}

// BackgroundResolverStats contains statistics about background resolution.
type BackgroundResolverStats struct {
	QueueLength int
	IsComplete  bool
	IsStarted   bool
}
