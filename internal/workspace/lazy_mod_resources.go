package workspace

import (
	"context"
	"sync"

	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/powerpipe/internal/resources"
)

// LazyModResources provides lazy-loading access to mod resources.
// Resources are loaded on-demand when accessed, rather than all at startup.
// This provides significant memory savings for large mods.
type LazyModResources struct {
	lw *LazyWorkspace

	// Cached resource maps (populated on demand for backward compatibility)
	dashboardsOnce sync.Once
	dashboards     map[string]*resources.Dashboard

	queriesOnce sync.Once
	queries     map[string]*resources.Query

	controlsOnce sync.Once
	controls     map[string]*resources.Control

	benchmarksOnce sync.Once
	benchmarks     map[string]*resources.Benchmark
}

// NewLazyModResources creates a lazy mod resources accessor.
func NewLazyModResources(lw *LazyWorkspace) *LazyModResources {
	return &LazyModResources{lw: lw}
}

// GetDashboard returns a dashboard by name, loading it on-demand if needed.
func (r *LazyModResources) GetDashboard(ctx context.Context, name string) (*resources.Dashboard, error) {
	return r.lw.loader.LoadDashboard(ctx, name)
}

// GetQuery returns a query by name, loading it on-demand if needed.
func (r *LazyModResources) GetQuery(ctx context.Context, name string) (*resources.Query, error) {
	resource, err := r.lw.loader.Load(ctx, name)
	if err != nil {
		return nil, err
	}
	if query, ok := resource.(*resources.Query); ok {
		return query, nil
	}
	return nil, nil
}

// GetControl returns a control by name, loading it on-demand if needed.
func (r *LazyModResources) GetControl(ctx context.Context, name string) (*resources.Control, error) {
	resource, err := r.lw.loader.Load(ctx, name)
	if err != nil {
		return nil, err
	}
	if control, ok := resource.(*resources.Control); ok {
		return control, nil
	}
	return nil, nil
}

// GetBenchmark returns a benchmark by name, loading it on-demand if needed.
func (r *LazyModResources) GetBenchmark(ctx context.Context, name string) (modconfig.ModTreeItem, error) {
	return r.lw.loader.LoadBenchmark(ctx, name)
}

// Dashboards returns all dashboards.
// WARNING: This loads ALL dashboards - use GetDashboard for lazy loading.
// This method exists for backward compatibility.
func (r *LazyModResources) Dashboards() map[string]*resources.Dashboard {
	r.dashboardsOnce.Do(func() {
		r.dashboards = make(map[string]*resources.Dashboard)
		ctx := context.Background()

		for _, entry := range r.lw.index.Dashboards() {
			dash, err := r.GetDashboard(ctx, entry.Name)
			if err == nil && dash != nil {
				r.dashboards[entry.Name] = dash
			}
		}
	})
	return r.dashboards
}

// Queries returns all queries.
// WARNING: This loads ALL queries - use GetQuery for lazy loading.
// This method exists for backward compatibility.
func (r *LazyModResources) Queries() map[string]*resources.Query {
	r.queriesOnce.Do(func() {
		r.queries = make(map[string]*resources.Query)
		ctx := context.Background()

		for _, entry := range r.lw.index.Queries() {
			q, err := r.GetQuery(ctx, entry.Name)
			if err == nil && q != nil {
				r.queries[entry.Name] = q
			}
		}
	})
	return r.queries
}

// Controls returns all controls.
// WARNING: This loads ALL controls - use GetControl for lazy loading.
// This method exists for backward compatibility.
func (r *LazyModResources) Controls() map[string]*resources.Control {
	r.controlsOnce.Do(func() {
		r.controls = make(map[string]*resources.Control)
		ctx := context.Background()

		for _, entry := range r.lw.index.Controls() {
			c, err := r.GetControl(ctx, entry.Name)
			if err == nil && c != nil {
				r.controls[entry.Name] = c
			}
		}
	})
	return r.controls
}

// Benchmarks returns all control benchmarks.
// WARNING: This loads ALL benchmarks - use GetBenchmark for lazy loading.
// This method exists for backward compatibility.
func (r *LazyModResources) Benchmarks() map[string]*resources.Benchmark {
	r.benchmarksOnce.Do(func() {
		r.benchmarks = make(map[string]*resources.Benchmark)
		ctx := context.Background()

		for _, entry := range r.lw.index.Benchmarks() {
			b, err := r.GetBenchmark(ctx, entry.Name)
			if err == nil && b != nil {
				if benchmark, ok := b.(*resources.Benchmark); ok {
					r.benchmarks[entry.Name] = benchmark
				}
			}
		}
	})
	return r.benchmarks
}

// GetResource returns a resource by parsed name.
// This implements backward compatibility with the ModResources interface.
func (r *LazyModResources) GetResource(parsedName *modconfig.ParsedResourceName) (modconfig.HclResource, bool) {
	return r.lw.GetResource(parsedName)
}

// WalkResources iterates over all resources.
// WARNING: This loads ALL resources - use sparingly.
// This method exists for backward compatibility.
func (r *LazyModResources) WalkResources(fn func(modconfig.HclResource) (bool, error)) error {
	ctx := context.Background()

	for _, entry := range r.lw.index.List() {
		resource, err := r.lw.loader.Load(ctx, entry.Name)
		if err != nil {
			continue
		}

		cont, err := fn(resource)
		if err != nil {
			return err
		}
		if !cont {
			break
		}
	}
	return nil
}

// DashboardCount returns the number of dashboards without loading them.
func (r *LazyModResources) DashboardCount() int {
	return len(r.lw.index.Dashboards())
}

// QueryCount returns the number of queries without loading them.
func (r *LazyModResources) QueryCount() int {
	return len(r.lw.index.Queries())
}

// ControlCount returns the number of controls without loading them.
func (r *LazyModResources) ControlCount() int {
	return len(r.lw.index.Controls())
}

// BenchmarkCount returns the number of benchmarks without loading them.
func (r *LazyModResources) BenchmarkCount() int {
	return len(r.lw.index.Benchmarks())
}

// TotalResourceCount returns the total number of resources without loading them.
func (r *LazyModResources) TotalResourceCount() int {
	return r.lw.index.Count()
}

// ListDashboardNames returns dashboard names without loading them.
func (r *LazyModResources) ListDashboardNames() []string {
	entries := r.lw.index.Dashboards()
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.Name
	}
	return names
}

// ListBenchmarkNames returns benchmark names without loading them.
func (r *LazyModResources) ListBenchmarkNames() []string {
	entries := r.lw.index.Benchmarks()
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.Name
	}
	return names
}

// ListQueryNames returns query names without loading them.
func (r *LazyModResources) ListQueryNames() []string {
	entries := r.lw.index.Queries()
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.Name
	}
	return names
}

// ListControlNames returns control names without loading them.
func (r *LazyModResources) ListControlNames() []string {
	entries := r.lw.index.Controls()
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.Name
	}
	return names
}
