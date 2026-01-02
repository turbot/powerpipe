package resourceloader

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/powerpipe/internal/resourcecache"
	"github.com/turbot/powerpipe/internal/resourceindex"
	"github.com/turbot/powerpipe/internal/resources"
)

// Loader provides on-demand loading of resources using the resource index
// and caching parsed resources for reuse.
type Loader struct {
	mu sync.RWMutex

	index   *resourceindex.ResourceIndex
	cache   *resourcecache.ResourceCache
	mod     *modconfig.Mod
	modPath string

	// Statistics
	loadCount int64
	parseTime int64

	// Resource provider for reference resolution
	resourceProvider modconfig.ResourceProvider
}

// NewLoader creates a resource loader.
func NewLoader(
	index *resourceindex.ResourceIndex,
	cache *resourcecache.ResourceCache,
	mod *modconfig.Mod,
	modPath string,
) *Loader {
	return &Loader{
		index:   index,
		cache:   cache,
		mod:     mod,
		modPath: modPath,
	}
}

// SetResourceProvider sets the resource provider for reference resolution during parsing.
func (l *Loader) SetResourceProvider(provider modconfig.ResourceProvider) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.resourceProvider = provider
}

// Load retrieves a resource by name, loading from disk if not cached.
func (l *Loader) Load(ctx context.Context, name string) (modconfig.HclResource, error) {
	// Check cache first
	if resource, ok := l.cache.GetResource(name); ok {
		return resource, nil
	}

	// Load from disk
	return l.loadFromDisk(ctx, name)
}

// LoadDashboard loads a dashboard with all its children.
func (l *Loader) LoadDashboard(ctx context.Context, name string) (*resources.Dashboard, error) {
	resource, err := l.Load(ctx, name)
	if err != nil {
		return nil, err
	}

	dash, ok := resource.(*resources.Dashboard)
	if !ok {
		return nil, fmt.Errorf("resource %s is not a dashboard", name)
	}

	// Load children recursively
	if err := l.loadChildren(ctx, dash); err != nil {
		return nil, err
	}

	return dash, nil
}

// LoadBenchmark loads a benchmark with all its children.
func (l *Loader) LoadBenchmark(ctx context.Context, name string) (modconfig.ModTreeItem, error) {
	resource, err := l.Load(ctx, name)
	if err != nil {
		return nil, err
	}

	bench, ok := resource.(modconfig.ModTreeItem)
	if !ok {
		return nil, fmt.Errorf("resource %s is not a benchmark", name)
	}

	if err := l.loadBenchmarkChildren(ctx, bench); err != nil {
		return nil, err
	}

	return bench, nil
}

// loadFromDisk parses a resource from its source file and caches it.
func (l *Loader) loadFromDisk(ctx context.Context, name string) (modconfig.HclResource, error) {
	entry, ok := l.index.Get(name)
	if !ok {
		return nil, fmt.Errorf("resource not found in index: %s", name)
	}

	start := time.Now()
	resource, err := l.parseResource(ctx, entry)
	if err != nil {
		return nil, fmt.Errorf("parsing resource %s: %w", name, err)
	}

	l.cache.PutResource(name, resource)

	atomic.AddInt64(&l.loadCount, 1)
	atomic.AddInt64(&l.parseTime, time.Since(start).Nanoseconds())

	return resource, nil
}

// loadChildren recursively loads all children of a ModTreeItem.
func (l *Loader) loadChildren(ctx context.Context, parent modconfig.ModTreeItem) error {
	for _, child := range parent.GetChildren() {
		if child == nil {
			continue
		}

		childName := child.Name()
		if _, ok := l.cache.GetResource(childName); !ok {
			if _, err := l.Load(ctx, childName); err != nil {
				// Child may be inline (defined within parent), continue
				continue
			}
		}

		if err := l.loadChildren(ctx, child); err != nil {
			return err
		}
	}
	return nil
}

// loadBenchmarkChildren loads all children of a benchmark recursively.
func (l *Loader) loadBenchmarkChildren(ctx context.Context, bench modconfig.ModTreeItem) error {
	entry, ok := l.index.Get(bench.Name())
	if !ok {
		return nil
	}

	for _, childName := range entry.ChildNames {
		child, err := l.Load(ctx, childName)
		if err != nil {
			return fmt.Errorf("loading child %s: %w", childName, err)
		}

		if childTree, ok := child.(modconfig.ModTreeItem); ok {
			if err := l.loadBenchmarkChildren(ctx, childTree); err != nil {
				return err
			}
		}

		// Load control query dependencies
		if control, ok := child.(*resources.Control); ok {
			if err := l.loadControlDependencies(ctx, control); err != nil {
				return err
			}
		}

		// Load detection query dependencies
		if detection, ok := child.(*resources.Detection); ok {
			if err := l.loadDetectionDependencies(ctx, detection); err != nil {
				return err
			}
		}
	}
	return nil
}

// loadControlDependencies loads the query referenced by a control.
func (l *Loader) loadControlDependencies(ctx context.Context, control *resources.Control) error {
	if control.Query != nil && control.Query.Name() != "" {
		if _, err := l.Load(ctx, control.Query.Name()); err != nil {
			return fmt.Errorf("loading query %s: %w", control.Query.Name(), err)
		}
	}
	return nil
}

// loadDetectionDependencies loads the query referenced by a detection.
func (l *Loader) loadDetectionDependencies(ctx context.Context, detection *resources.Detection) error {
	if detection.Query != nil && detection.Query.Name() != "" {
		if _, err := l.Load(ctx, detection.Query.Name()); err != nil {
			return fmt.Errorf("loading query %s: %w", detection.Query.Name(), err)
		}
	}
	return nil
}

// Preload loads multiple resources in parallel.
func (l *Loader) Preload(ctx context.Context, names []string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(names))
	sem := make(chan struct{}, 10) // Limit concurrency

	for _, name := range names {
		if _, ok := l.cache.GetResource(name); ok {
			continue // Already cached
		}

		wg.Add(1)
		go func(n string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			if _, err := l.Load(ctx, n); err != nil {
				errChan <- err
			}
		}(name)
	}

	wg.Wait()
	close(errChan)

	// Return first error if any
	for err := range errChan {
		return err
	}
	return nil
}

// PreloadByType loads all resources of a given type in parallel.
func (l *Loader) PreloadByType(ctx context.Context, resourceType string) error {
	entries := l.index.GetByType(resourceType)
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.Name
	}
	return l.Preload(ctx, names)
}

// Stats returns loader statistics.
func (l *Loader) Stats() LoaderStats {
	loadCount := atomic.LoadInt64(&l.loadCount)
	parseTime := atomic.LoadInt64(&l.parseTime)

	avgTime := time.Duration(0)
	if loadCount > 0 {
		avgTime = time.Duration(parseTime / loadCount)
	}

	return LoaderStats{
		LoadCount:    loadCount,
		AvgParseTime: avgTime,
		CacheStats:   l.cache.Stats(),
	}
}

// LoaderStats contains loader performance statistics.
type LoaderStats struct {
	LoadCount    int64
	AvgParseTime time.Duration
	CacheStats   resourcecache.CacheStats
}

// Index returns the underlying resource index.
func (l *Loader) Index() *resourceindex.ResourceIndex {
	return l.index
}

// Cache returns the underlying resource cache.
func (l *Loader) Cache() *resourcecache.ResourceCache {
	return l.cache
}

// Mod returns the mod for this loader.
func (l *Loader) Mod() *modconfig.Mod {
	return l.mod
}

// Clear clears the loader cache.
func (l *Loader) Clear() {
	l.cache.Clear()
	atomic.StoreInt64(&l.loadCount, 0)
	atomic.StoreInt64(&l.parseTime, 0)
}

// Invalidate removes a specific resource from the cache.
func (l *Loader) Invalidate(name string) {
	l.cache.Invalidate(name)
}
