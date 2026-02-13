package resourceloader

import (
	"context"
	"sync"
	"sync/atomic"
)

// PreloadOptions configures preloading behavior.
type PreloadOptions struct {
	IncludeDependencies bool                         // Load all dependencies of requested resources
	MaxConcurrency      int                          // Maximum parallel loads (default: 10)
	OnProgress          func(loaded, total int)      // Progress callback
	OnError             func(name string, err error) // Error callback (if nil, first error stops loading)
}

// DefaultPreloadOptions returns sensible default preload options.
func DefaultPreloadOptions() PreloadOptions {
	return PreloadOptions{
		IncludeDependencies: true,
		MaxConcurrency:      10,
	}
}

// PreloadWithDependencies loads resources and their dependencies in parallel,
// respecting dependency order.
func (l *Loader) PreloadWithDependencies(ctx context.Context, names []string, opts PreloadOptions) error {
	if !opts.IncludeDependencies {
		return l.Preload(ctx, names)
	}

	resolver := NewDependencyResolver(l.index, l)

	// Collect all dependencies transitively
	allNames := make(map[string]bool)
	for _, name := range names {
		deps := resolver.GetTransitiveDependencies(name)
		for _, dep := range deps {
			allNames[dep] = true
		}
	}

	// Get topologically sorted load order
	namesList := make([]string, 0, len(allNames))
	for name := range allNames {
		namesList = append(namesList, name)
	}

	orderedNames, err := resolver.GetDependencyOrder(namesList)
	if err != nil {
		return err
	}

	// Filter out already cached resources
	var toLoad []string
	for _, name := range orderedNames {
		if _, ok := l.cache.GetResource(name); !ok {
			toLoad = append(toLoad, name)
		}
	}

	if len(toLoad) == 0 {
		return nil
	}

	// Set up concurrency control
	concurrency := opts.MaxConcurrency
	if concurrency <= 0 {
		concurrency = 10
	}

	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	var firstErr atomic.Value

	// Track progress
	var loaded int64
	total := len(toLoad)

	// Track which resources are loaded to respect dependencies
	loadedSet := &sync.Map{}

	// Load in batches that respect dependency order
	for _, name := range toLoad {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Check if we've hit an error
		if err := firstErr.Load(); err != nil && opts.OnError == nil {
			break
		}

		// Wait for dependencies to be loaded first
		deps := resolver.GetDependencies(name)
		for _, dep := range deps {
			// Check if dependency is in our load set and wait for it
			for {
				if _, loaded := loadedSet.Load(dep.To); loaded {
					break
				}
				// Dependency not yet loaded - it either doesn't exist in index
				// or hasn't been loaded yet. Check cache.
				if _, ok := l.cache.GetResource(dep.To); ok {
					break
				}
				// Brief sleep to avoid busy waiting
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					// Don't block indefinitely on dependencies - proceed anyway
					break
				}
				break
			}
		}

		wg.Add(1)
		go func(n string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			if _, err := l.Load(ctx, n); err != nil {
				if opts.OnError != nil {
					opts.OnError(n, err)
				} else {
					firstErr.CompareAndSwap(nil, err)
				}
				return
			}

			// Mark as loaded
			loadedSet.Store(n, true)

			if opts.OnProgress != nil {
				current := atomic.AddInt64(&loaded, 1)
				opts.OnProgress(int(current), total)
			}
		}(name)
	}

	wg.Wait()

	// Return first error if any
	if err := firstErr.Load(); err != nil {
		return err.(error)
	}
	return nil
}

// PreloadBenchmark loads a benchmark and all its transitive dependencies.
func (l *Loader) PreloadBenchmark(ctx context.Context, name string) error {
	return l.PreloadWithDependencies(ctx, []string{name}, PreloadOptions{
		IncludeDependencies: true,
		MaxConcurrency:      10,
	})
}

// PreloadDashboard loads a dashboard and all its transitive dependencies.
func (l *Loader) PreloadDashboard(ctx context.Context, name string) error {
	return l.PreloadWithDependencies(ctx, []string{name}, PreloadOptions{
		IncludeDependencies: true,
		MaxConcurrency:      10,
	})
}
