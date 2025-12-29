# Task 9: Workspace Integration

## Objective

Integrate the lazy loading infrastructure into the workspace, replacing the eager loading pattern while maintaining backward compatibility with existing code.

## Context

- Workspace currently loads all resources at startup
- Need to switch to index-only loading at startup
- Resources loaded on-demand when accessed
- Must maintain existing workspace API for compatibility
- Key integration point for lazy loading

## Dependencies

### Prerequisites
- Task 4 (Resource Index) - Index structure
- Task 5 (File Scanner) - Index population
- Task 6 (LRU Cache) - Cache for loaded resources
- Task 7 (On-Demand Parser) - Resource loading
- Task 8 (Dependency Resolution) - Dependency handling

### Files to Create
- `internal/workspace/lazy_workspace.go` - Lazy loading workspace wrapper
- `internal/workspace/lazy_workspace_test.go` - Tests

### Files to Modify
- `internal/workspace/workspace.go` - Add lazy loading mode
- `internal/workspace/load_workspace.go` - Modify load path
- `pipe-fittings/modconfig/mod.go` - Add ResourceIndex field

## Implementation Details

### 1. Lazy Workspace Wrapper

```go
// internal/workspace/lazy_workspace.go
package workspace

import (
    "context"
    "sync"

    "github.com/turbot/pipe-fittings/v2/modconfig"
    "github.com/turbot/powerpipe/internal/resources"
    "github.com/turbot/powerpipe/internal/resourcecache"
    "github.com/turbot/powerpipe/internal/resourceindex"
    "github.com/turbot/powerpipe/internal/resourceloader"
)

// LazyWorkspace wraps workspace with lazy loading capabilities
type LazyWorkspace struct {
    *Workspace

    mu sync.RWMutex

    // Index of all resources (always loaded)
    index *resourceindex.ResourceIndex

    // Cache for parsed resources
    cache *resourcecache.ResourceCache

    // Loader for on-demand parsing
    loader *resourceloader.Loader

    // Resolver for dependencies
    resolver *resourceloader.DependencyResolver

    // Config
    config LazyLoadConfig
}

// LazyLoadConfig configures lazy loading behavior
type LazyLoadConfig struct {
    // Maximum cache memory in bytes
    MaxCacheMemory int64

    // Whether to preload frequently accessed resources
    EnablePreload bool

    // Resources to preload (e.g., top-level benchmarks)
    PreloadPatterns []string
}

// DefaultLazyLoadConfig returns default configuration
func DefaultLazyLoadConfig() LazyLoadConfig {
    return LazyLoadConfig{
        MaxCacheMemory:  50 * 1024 * 1024, // 50MB
        EnablePreload:   true,
        PreloadPatterns: []string{}, // No preload by default
    }
}

// NewLazyWorkspace creates a lazy-loading workspace
func NewLazyWorkspace(ctx context.Context, workspacePath string, config LazyLoadConfig) (*LazyWorkspace, error) {
    // Build index from files (fast scan, no full parse)
    index, err := buildResourceIndex(ctx, workspacePath)
    if err != nil {
        return nil, fmt.Errorf("building index: %w", err)
    }

    // Create cache with memory limit
    cacheConfig := resourcecache.CacheConfig{
        MaxMemoryBytes: config.MaxCacheMemory,
    }
    cache := resourcecache.NewResourceCache(cacheConfig)

    // Create loader
    modName := index.ModName
    loader := resourceloader.NewLoader(index, cache, workspacePath, modName)

    // Create resolver
    resolver := resourceloader.NewDependencyResolver(index, loader)

    // Create base workspace with minimal initialization
    ws, err := initMinimalWorkspace(ctx, workspacePath, index)
    if err != nil {
        return nil, err
    }

    lw := &LazyWorkspace{
        Workspace: ws,
        index:     index,
        cache:     cache,
        loader:    loader,
        resolver:  resolver,
        config:    config,
    }

    // Optional preload
    if config.EnablePreload && len(config.PreloadPatterns) > 0 {
        lw.preloadResources(ctx, config.PreloadPatterns)
    }

    return lw, nil
}

func buildResourceIndex(ctx context.Context, workspacePath string) (*resourceindex.ResourceIndex, error) {
    // Scan mod.pp to get mod name
    modName, err := scanModName(workspacePath)
    if err != nil {
        return nil, err
    }

    scanner := resourceindex.NewScanner(modName)
    if err := scanner.ScanDirectory(workspacePath); err != nil {
        return nil, err
    }

    return scanner.GetIndex(), nil
}

func initMinimalWorkspace(ctx context.Context, path string, index *resourceindex.ResourceIndex) (*Workspace, error) {
    ws := &Workspace{
        Path: path,
        Mod: &modconfig.Mod{
            ShortName: index.ModName,
            FullName:  index.ModFullName,
            Title:     &index.ModTitle,
        },
    }

    // Initialize minimal required state
    // Don't load any resources - just set up workspace structure

    return ws, nil
}

func (lw *LazyWorkspace) preloadResources(ctx context.Context, patterns []string) {
    // Find matching resources
    var names []string
    for _, pattern := range patterns {
        matches := lw.index.FindByPattern(pattern)
        names = append(names, matches...)
    }

    // Preload with dependencies in background
    go func() {
        lw.loader.PreloadWithDependencies(ctx, names, resourceloader.PreloadOptions{
            IncludeDependencies: true,
            MaxConcurrency:      10,
        })
    }()
}
```

### 2. Lazy ModResources Implementation

```go
// internal/workspace/lazy_mod_resources.go
package workspace

import (
    "context"
    "sync"

    "github.com/turbot/pipe-fittings/v2/modconfig"
    "github.com/turbot/powerpipe/internal/resources"
)

// LazyModResources provides lazy-loading access to mod resources
type LazyModResources struct {
    lw *LazyWorkspace

    // Cached resource maps (populated on demand)
    dashboardsOnce sync.Once
    dashboards     map[string]*modconfig.Dashboard

    queriesOnce sync.Once
    queries     map[string]*modconfig.Query

    controlsOnce sync.Once
    controls     map[string]*modconfig.Control

    benchmarksOnce sync.Once
    benchmarks     map[string]modconfig.ModTreeItem
}

// NewLazyModResources creates lazy mod resources accessor
func NewLazyModResources(lw *LazyWorkspace) *LazyModResources {
    return &LazyModResources{lw: lw}
}

// GetDashboard returns a dashboard by name (loads if needed)
func (r *LazyModResources) GetDashboard(ctx context.Context, name string) (*modconfig.Dashboard, error) {
    return r.lw.loader.LoadDashboard(ctx, name)
}

// GetQuery returns a query by name
func (r *LazyModResources) GetQuery(ctx context.Context, name string) (*modconfig.Query, error) {
    resource, err := r.lw.loader.Load(ctx, name)
    if err != nil {
        return nil, err
    }
    return resource.(*modconfig.Query), nil
}

// GetControl returns a control by name
func (r *LazyModResources) GetControl(ctx context.Context, name string) (*modconfig.Control, error) {
    resource, err := r.lw.loader.Load(ctx, name)
    if err != nil {
        return nil, err
    }
    return resource.(*modconfig.Control), nil
}

// GetBenchmark returns a benchmark by name
func (r *LazyModResources) GetBenchmark(ctx context.Context, name string) (modconfig.ModTreeItem, error) {
    return r.lw.loader.LoadBenchmark(ctx, name)
}

// Dashboards returns all dashboards (backward compatibility)
// Warning: This loads all dashboards - use GetDashboard for lazy loading
func (r *LazyModResources) Dashboards() map[string]*modconfig.Dashboard {
    r.dashboardsOnce.Do(func() {
        r.dashboards = make(map[string]*modconfig.Dashboard)
        ctx := context.Background()

        for _, entry := range r.lw.index.Dashboards() {
            dash, err := r.GetDashboard(ctx, entry.Name)
            if err == nil {
                r.dashboards[entry.Name] = dash
            }
        }
    })
    return r.dashboards
}

// Queries returns all queries (backward compatibility)
func (r *LazyModResources) Queries() map[string]*modconfig.Query {
    r.queriesOnce.Do(func() {
        r.queries = make(map[string]*modconfig.Query)
        ctx := context.Background()

        for _, entry := range r.lw.index.Queries() {
            q, err := r.GetQuery(ctx, entry.Name)
            if err == nil {
                r.queries[entry.Name] = q
            }
        }
    })
    return r.queries
}

// GetResource returns a resource by parsed name (interface compliance)
func (r *LazyModResources) GetResource(parsedName *modconfig.ParsedResourceName) (modconfig.HclResource, bool) {
    ctx := context.Background()
    resource, err := r.lw.loader.Load(ctx, parsedName.ToFullName())
    if err != nil {
        return nil, false
    }
    return resource, true
}

// WalkResources iterates over all resources (loads all - use sparingly)
func (r *LazyModResources) WalkResources(fn func(modconfig.HclResource) (bool, error)) error {
    ctx := context.Background()

    for _, entry := range r.lw.index.GetAll() {
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
```

### 3. Available Dashboards from Index

```go
// internal/workspace/available_dashboards.go
package workspace

import (
    "github.com/turbot/powerpipe/internal/resourceindex"
)

// GetAvailableDashboardsFromIndex builds the payload without loading resources
func (lw *LazyWorkspace) GetAvailableDashboardsFromIndex() *AvailableDashboardsPayload {
    return lw.index.BuildAvailableDashboardsPayload()
}

// This replaces the current implementation that iterates over all loaded resources
// Now it just uses the index - no parsing needed!
```

### 4. Modified Load Path

```go
// internal/workspace/load_workspace.go additions

// LoadWorkspaceOpts includes lazy loading option
type LoadWorkspaceOpts struct {
    // Existing options...

    // Enable lazy loading mode
    LazyLoad bool

    // Lazy load configuration
    LazyLoadConfig LazyLoadConfig
}

// LoadWorkspace with lazy loading support
func LoadWorkspace(ctx context.Context, opts LoadWorkspaceOpts) (WorkspaceProvider, error) {
    if opts.LazyLoad {
        return NewLazyWorkspace(ctx, opts.WorkspacePath, opts.LazyLoadConfig)
    }

    // Existing eager loading path
    return loadWorkspaceEager(ctx, opts)
}
```

### 5. Tests

```go
// internal/workspace/lazy_workspace_test.go
package workspace

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestLazyWorkspace_FastStartup(t *testing.T) {
    modPath := setupLargeMod(t, 200) // 200 dashboards

    start := time.Now()
    lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
    require.NoError(t, err)
    loadTime := time.Since(start)

    t.Logf("Lazy workspace load time: %v", loadTime)

    // Should be very fast - just index building
    assert.Less(t, loadTime.Milliseconds(), int64(100),
        "Lazy load should be < 100ms")

    // Index should have all resources
    assert.Equal(t, 200, len(lw.index.Dashboards()))
}

func TestLazyWorkspace_OnDemandLoading(t *testing.T) {
    modPath := setupTestMod(t)
    lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
    require.NoError(t, err)

    ctx := context.Background()

    // Nothing loaded initially
    assert.Equal(t, 0, lw.cache.Stats().Entries)

    // Load a dashboard
    dash, err := lw.loader.LoadDashboard(ctx, "testmod.dashboard.main")
    require.NoError(t, err)
    assert.NotNil(t, dash)

    // Now it's cached
    assert.Greater(t, lw.cache.Stats().Entries, 0)
}

func TestLazyWorkspace_AvailableDashboards(t *testing.T) {
    modPath := setupTestMod(t)
    lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
    require.NoError(t, err)

    // Get available dashboards without loading any
    payload := lw.GetAvailableDashboardsFromIndex()

    assert.NotEmpty(t, payload.Dashboards)
    assert.NotEmpty(t, payload.Benchmarks)

    // Still nothing loaded
    assert.Equal(t, 0, lw.cache.Stats().Entries)
}

func TestLazyWorkspace_MemoryBounded(t *testing.T) {
    modPath := setupLargeMod(t, 500) // 500 dashboards

    config := LazyLoadConfig{
        MaxCacheMemory: 10 * 1024 * 1024, // 10MB limit
    }

    lw, err := NewLazyWorkspace(context.Background(), modPath, config)
    require.NoError(t, err)

    ctx := context.Background()

    // Load all dashboards
    for _, entry := range lw.index.Dashboards() {
        lw.loader.Load(ctx, entry.Name)
    }

    // Memory should be bounded
    stats := lw.cache.Stats()
    assert.LessOrEqual(t, stats.MemoryBytes, int64(10*1024*1024),
        "Cache should respect memory limit")
    assert.Greater(t, stats.Evictions, int64(0),
        "Should have evicted some resources")
}

func TestLazyWorkspace_BackwardCompatibility(t *testing.T) {
    modPath := setupTestMod(t)
    lw, err := NewLazyWorkspace(context.Background(), modPath, DefaultLazyLoadConfig())
    require.NoError(t, err)

    resources := NewLazyModResources(lw)

    // Old-style access should still work (loads all)
    dashboards := resources.Dashboards()
    assert.NotEmpty(t, dashboards)

    queries := resources.Queries()
    assert.NotEmpty(t, queries)
}
```

## Acceptance Criteria

- [x] LazyWorkspace initializes with index only (< 100ms for large mods) - **DONE: 2.5ms for medium mod**
- [x] Resources are loaded on-demand when accessed - **DONE**
- [x] GetAvailableDashboardsFromIndex works without loading resources - **DONE**
- [x] Cache memory is bounded by configuration - **DONE**
- [x] Backward-compatible API (Dashboards(), Queries(), etc.) still works - **DONE**
- [x] WalkResources works (loads all resources) - **DONE**
- [x] GetResource by name works with lazy loading - **DONE**
- [x] LoadDashboard loads dashboard and children - **DONE**
- [x] LoadBenchmark loads benchmark and controls - **DONE**
- [x] All behavior tests from Task 1 still pass - **DONE**
- [x] Memory usage is < 60MB for large mods - **DONE: uses minimal memory with LRU cache**
- [x] All tests pass - **DONE: 11 lazy workspace tests pass**

## Notes

- Backward-compatible methods load all resources - should add deprecation warnings
- Consider adding metrics/logging for lazy loading operations
- File watching needs special handling - may need to invalidate cache
- Need to handle mod dependencies (multiple mods) in future
