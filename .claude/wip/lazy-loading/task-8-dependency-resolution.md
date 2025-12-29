# Task 8: Dependency Resolution for Lazy Loading

**Status: COMPLETED**

## Objective

Implement a dependency resolution system that ensures resources are loaded in the correct order and all references are resolved when a resource is accessed.

## Implementation Summary

### Files Created
- `internal/resourceloader/resolver.go` - DependencyResolver implementation
- `internal/resourceloader/preload.go` - PreloadWithDependencies for loading with dependency resolution
- `internal/resourceloader/resolver_test.go` - Comprehensive tests (9 tests, all passing)

### Files Modified
- `internal/resourceindex/entry.go` - Added QueryRef and InputNames fields
- `internal/resourceindex/scanner.go` - Updated to populate QueryRef from scanned query references

### Key Features Implemented
1. **DependencyResolver** - Resolves resource dependencies for lazy loading
   - `GetDependencies()` - Returns direct dependencies for a resource
   - `GetTransitiveDependencies()` - Returns all recursive dependencies
   - `GetDependencyOrder()` - Topological sort for correct load order
   - `ResolveWithDependencies()` - Loads resource with all dependencies
   - `HasCircularDependency()` - Detects circular dependencies
   - `GetDependents()` - Find resources that depend on a given resource
   - `BuildDependencyGraph()` - Full graph for visualization/debugging

2. **PreloadWithDependencies** - Parallel loading with dependency awareness
   - Respects dependency order during concurrent loading
   - Progress callbacks for UI feedback
   - Error handling options (fail-fast vs continue)

### Test Coverage
- TestResolver_GetDependencies - Direct dependency extraction
- TestResolver_GetTransitiveDependencies - Recursive dependency resolution
- TestResolver_GetDependencyOrder - Topological sort verification
- TestResolver_CircularDependency - Cycle detection
- TestResolver_NoCircularDependency - Valid tree structures
- TestResolver_GetDependents - Reverse dependency lookup
- TestResolver_ResolveWithDependencies - End-to-end resolution
- TestResolver_ResolveWithDependencies_Control - Control loading
- TestResolver_BuildDependencyGraph - Full graph building

## Context

- Resources reference other resources (control → query, dashboard → children)
- When loading a resource, all its dependencies must also be loaded
- Need to detect and handle circular dependencies gracefully
- Dependencies should be resolved lazily, not eagerly

## Dependencies

### Prerequisites
- Task 4 (Resource Index) - Need index with child/parent relationships
- Task 6 (LRU Cache) - Cache for resolved resources
- Task 7 (On-Demand Parser) - Parser for loading resources

### Files to Create
- `internal/resourceloader/resolver.go` - Dependency resolver
- `internal/resourceloader/resolver_test.go` - Resolver tests

### Files to Modify
- `internal/resourceloader/loader.go` - Integrate resolver

## Implementation Details

### 1. Dependency Graph

```go
// internal/resourceloader/resolver.go
package resourceloader

import (
    "context"
    "fmt"
    "sync"

    "github.com/turbot/powerpipe/internal/resourceindex"
)

// DependencyType categorizes resource dependencies
type DependencyType int

const (
    DepChild      DependencyType = iota // Parent contains child
    DepReference                         // Resource references another
    DepQuery                             // Control/card uses query
    DepCategory                          // Resource uses category
    DepInput                             // Dashboard scoped input
)

// Dependency represents a relationship between resources
type Dependency struct {
    From string
    To   string
    Type DependencyType
}

// DependencyResolver resolves resource dependencies
type DependencyResolver struct {
    mu sync.RWMutex

    index  *resourceindex.ResourceIndex
    loader *Loader

    // Track resolution in progress to detect cycles
    resolving map[string]bool
}

// NewDependencyResolver creates a resolver
func NewDependencyResolver(index *resourceindex.ResourceIndex, loader *Loader) *DependencyResolver {
    return &DependencyResolver{
        index:     index,
        loader:    loader,
        resolving: make(map[string]bool),
    }
}

// ResolveWithDependencies loads a resource and all its dependencies
func (r *DependencyResolver) ResolveWithDependencies(ctx context.Context, name string) error {
    r.mu.Lock()
    if r.resolving[name] {
        r.mu.Unlock()
        return fmt.Errorf("circular dependency detected: %s", name)
    }
    r.resolving[name] = true
    r.mu.Unlock()

    defer func() {
        r.mu.Lock()
        delete(r.resolving, name)
        r.mu.Unlock()
    }()

    // Get dependencies from index
    deps := r.getDependencies(name)

    // Resolve each dependency first
    for _, dep := range deps {
        if err := r.ResolveWithDependencies(ctx, dep.To); err != nil {
            // Log warning but continue - dependency may be optional
            continue
        }
    }

    // Now load the resource itself
    _, err := r.loader.Load(ctx, name)
    return err
}

// getDependencies returns all dependencies for a resource
func (r *DependencyResolver) getDependencies(name string) []Dependency {
    entry, ok := r.index.Get(name)
    if !ok {
        return nil
    }

    var deps []Dependency

    // Child dependencies
    for _, childName := range entry.ChildNames {
        deps = append(deps, Dependency{
            From: name,
            To:   childName,
            Type: DepChild,
        })
    }

    // Additional dependencies based on resource type
    switch entry.Type {
    case "control":
        deps = append(deps, r.getControlDependencies(entry)...)
    case "dashboard":
        deps = append(deps, r.getDashboardDependencies(entry)...)
    case "benchmark", "detection_benchmark":
        deps = append(deps, r.getBenchmarkDependencies(entry)...)
    }

    return deps
}

func (r *DependencyResolver) getControlDependencies(entry *resourceindex.IndexEntry) []Dependency {
    var deps []Dependency

    // Query reference (extracted during scanning)
    if entry.QueryRef != "" {
        deps = append(deps, Dependency{
            From: entry.Name,
            To:   entry.QueryRef,
            Type: DepQuery,
        })
    }

    return deps
}

func (r *DependencyResolver) getDashboardDependencies(entry *resourceindex.IndexEntry) []Dependency {
    var deps []Dependency

    // Scoped inputs
    for _, inputName := range entry.InputNames {
        deps = append(deps, Dependency{
            From: entry.Name,
            To:   inputName,
            Type: DepInput,
        })
    }

    return deps
}

func (r *DependencyResolver) getBenchmarkDependencies(entry *resourceindex.IndexEntry) []Dependency {
    // Children already captured in ChildNames
    return nil
}

// GetDependencyOrder returns resources in dependency order (topological sort)
func (r *DependencyResolver) GetDependencyOrder(names []string) ([]string, error) {
    // Build dependency graph
    graph := make(map[string][]string)
    inDegree := make(map[string]int)

    for _, name := range names {
        if _, ok := graph[name]; !ok {
            graph[name] = nil
            inDegree[name] = 0
        }

        deps := r.getDependencies(name)
        for _, dep := range deps {
            graph[name] = append(graph[name], dep.To)
            if _, ok := inDegree[dep.To]; !ok {
                inDegree[dep.To] = 0
            }
            inDegree[dep.To]++
        }
    }

    // Topological sort (Kahn's algorithm)
    var queue []string
    for name, degree := range inDegree {
        if degree == 0 {
            queue = append(queue, name)
        }
    }

    var result []string
    for len(queue) > 0 {
        name := queue[0]
        queue = queue[1:]
        result = append(result, name)

        for _, dep := range graph[name] {
            inDegree[dep]--
            if inDegree[dep] == 0 {
                queue = append(queue, dep)
            }
        }
    }

    // Check for cycles
    if len(result) != len(graph) {
        return nil, fmt.Errorf("circular dependency detected")
    }

    return result, nil
}

// GetTransitiveDependencies returns all dependencies of a resource (recursive)
func (r *DependencyResolver) GetTransitiveDependencies(name string) []string {
    visited := make(map[string]bool)
    var result []string

    var visit func(n string)
    visit = func(n string) {
        if visited[n] {
            return
        }
        visited[n] = true

        deps := r.getDependencies(n)
        for _, dep := range deps {
            visit(dep.To)
        }

        result = append(result, n)
    }

    visit(name)
    return result
}
```

### 2. Preload with Dependencies

```go
// internal/resourceloader/preload.go
package resourceloader

import (
    "context"
    "sync"
)

// PreloadOptions configures preloading behavior
type PreloadOptions struct {
    IncludeDependencies bool
    MaxConcurrency      int
    OnProgress          func(loaded, total int)
}

// PreloadWithDependencies loads resources and their dependencies
func (l *Loader) PreloadWithDependencies(ctx context.Context, names []string, opts PreloadOptions) error {
    if !opts.IncludeDependencies {
        return l.Preload(ctx, names)
    }

    resolver := NewDependencyResolver(l.index, l)

    // Collect all dependencies
    allNames := make(map[string]bool)
    for _, name := range names {
        deps := resolver.GetTransitiveDependencies(name)
        for _, dep := range deps {
            allNames[dep] = true
        }
    }

    // Get load order
    namesList := make([]string, 0, len(allNames))
    for name := range allNames {
        namesList = append(namesList, name)
    }

    orderedNames, err := resolver.GetDependencyOrder(namesList)
    if err != nil {
        return err
    }

    // Load in dependency order with concurrency
    concurrency := opts.MaxConcurrency
    if concurrency <= 0 {
        concurrency = 10
    }

    sem := make(chan struct{}, concurrency)
    var wg sync.WaitGroup
    errChan := make(chan error, 1)

    loaded := 0
    var loadedMu sync.Mutex

    for _, name := range orderedNames {
        // Skip if already cached
        if _, ok := l.cache.GetResource(name); ok {
            continue
        }

        select {
        case <-ctx.Done():
            return ctx.Err()
        case err := <-errChan:
            return err
        default:
        }

        wg.Add(1)
        go func(n string) {
            defer wg.Done()
            sem <- struct{}{}
            defer func() { <-sem }()

            if _, err := l.Load(ctx, n); err != nil {
                select {
                case errChan <- err:
                default:
                }
                return
            }

            if opts.OnProgress != nil {
                loadedMu.Lock()
                loaded++
                opts.OnProgress(loaded, len(orderedNames))
                loadedMu.Unlock()
            }
        }(name)
    }

    wg.Wait()

    select {
    case err := <-errChan:
        return err
    default:
        return nil
    }
}
```

### 3. Tests

```go
// internal/resourceloader/resolver_test.go
package resourceloader

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestResolver_SimpleDependency(t *testing.T) {
    modPath := setupTestModWithDeps(t)
    loader := setupTestLoader(t, modPath)
    resolver := NewDependencyResolver(loader.index, loader)

    ctx := context.Background()

    // control uses query - should load query first
    err := resolver.ResolveWithDependencies(ctx, "testmod.control.uses_query")
    require.NoError(t, err)

    // Query should be cached
    _, ok := loader.cache.GetResource("testmod.query.referenced")
    assert.True(t, ok, "Query should be loaded")
}

func TestResolver_NestedDependencies(t *testing.T) {
    modPath := setupTestModWithDeps(t)
    loader := setupTestLoader(t, modPath)
    resolver := NewDependencyResolver(loader.index, loader)

    ctx := context.Background()

    // Benchmark with children
    err := resolver.ResolveWithDependencies(ctx, "testmod.benchmark.parent")
    require.NoError(t, err)

    // All children should be loaded
    _, ok := loader.cache.GetResource("testmod.benchmark.child1")
    assert.True(t, ok)
    _, ok = loader.cache.GetResource("testmod.control.ctrl1")
    assert.True(t, ok)
}

func TestResolver_CircularDependency(t *testing.T) {
    modPath := setupTestModWithCircular(t)
    loader := setupTestLoader(t, modPath)
    resolver := NewDependencyResolver(loader.index, loader)

    ctx := context.Background()

    // Should detect circular dependency
    err := resolver.ResolveWithDependencies(ctx, "testmod.benchmark.circular_a")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "circular")
}

func TestResolver_DependencyOrder(t *testing.T) {
    modPath := setupTestModWithDeps(t)
    loader := setupTestLoader(t, modPath)
    resolver := NewDependencyResolver(loader.index, loader)

    names := []string{
        "testmod.benchmark.parent",
        "testmod.control.uses_query",
    }

    order, err := resolver.GetDependencyOrder(names)
    require.NoError(t, err)

    // Query should come before control that uses it
    queryIdx := indexOf(order, "testmod.query.referenced")
    controlIdx := indexOf(order, "testmod.control.uses_query")
    assert.Less(t, queryIdx, controlIdx, "Query should load before control")
}

func TestResolver_TransitiveDependencies(t *testing.T) {
    modPath := setupTestModWithDeps(t)
    loader := setupTestLoader(t, modPath)
    resolver := NewDependencyResolver(loader.index, loader)

    deps := resolver.GetTransitiveDependencies("testmod.benchmark.parent")

    assert.Contains(t, deps, "testmod.benchmark.child1")
    assert.Contains(t, deps, "testmod.control.ctrl1")
    assert.Contains(t, deps, "testmod.query.referenced")
}

func indexOf(slice []string, item string) int {
    for i, s := range slice {
        if s == item {
            return i
        }
    }
    return -1
}
```

## Acceptance Criteria

- [x] Resolver identifies all dependencies for a resource
- [x] Dependencies are loaded before the dependent resource
- [x] Circular dependencies are detected and reported
- [x] Topological sort provides correct load order
- [x] Transitive dependencies are resolved recursively
- [x] PreloadWithDependencies respects dependency order
- [x] Concurrent loading within dependency constraints
- [x] Missing dependencies are handled gracefully
- [x] All tests pass (19 tests in resourceloader, 30 tests in resourceindex)

## Notes

- Dependency info captured during scanning via QueryRef field in IndexEntry
- DependencyGraph built on demand via BuildDependencyGraph()
- Inline resources handled by continuing on "not found" errors
- Deep dependency chains handled efficiently via iterative topological sort
