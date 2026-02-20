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
	DepChild     DependencyType = iota // Parent contains child
	DepReference                       // Resource references another
	DepQuery                           // Control/card uses query
	DepCategory                        // Resource uses category
	DepInput                           // Dashboard scoped input
)

// Dependency represents a relationship between resources
type Dependency struct {
	From string
	To   string
	Type DependencyType
}

// DependencyResolver resolves resource dependencies for lazy loading.
// It ensures resources are loaded in the correct order and detects circular dependencies.
type DependencyResolver struct {
	mu sync.RWMutex

	index  *resourceindex.ResourceIndex
	loader *Loader

	// Track resolution in progress to detect cycles
	resolving map[string]bool
}

// NewDependencyResolver creates a resolver for the given index and loader.
func NewDependencyResolver(index *resourceindex.ResourceIndex, loader *Loader) *DependencyResolver {
	return &DependencyResolver{
		index:     index,
		loader:    loader,
		resolving: make(map[string]bool),
	}
}

// ResolveWithDependencies loads a resource and all its dependencies in the correct order.
// Returns an error if a circular dependency is detected.
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
	deps := r.GetDependencies(name)

	// Resolve each dependency first
	for _, dep := range deps {
		if err := r.ResolveWithDependencies(ctx, dep.To); err != nil {
			// Check if it's a "not found" error - dependency may be optional or inline
			if ctx.Err() != nil {
				return ctx.Err()
			}
			// Continue for missing dependencies - they may be inline resources
			continue
		}
	}

	// Now load the resource itself
	_, err := r.loader.Load(ctx, name)
	return err
}

// GetDependencies returns all direct dependencies for a resource.
func (r *DependencyResolver) GetDependencies(name string) []Dependency {
	entry, ok := r.index.Get(name)
	if !ok {
		return nil
	}

	var deps []Dependency

	// Child dependencies (for benchmarks, dashboards, containers)
	for _, childName := range entry.ChildNames {
		deps = append(deps, Dependency{
			From: name,
			To:   childName,
			Type: DepChild,
		})
	}

	// Query reference (for controls, cards, tables, etc.)
	if entry.QueryRef != "" {
		deps = append(deps, Dependency{
			From: name,
			To:   entry.QueryRef,
			Type: DepQuery,
		})
	}

	// Input dependencies (for dashboards with scoped inputs)
	for _, inputName := range entry.InputNames {
		deps = append(deps, Dependency{
			From: name,
			To:   inputName,
			Type: DepInput,
		})
	}

	return deps
}

// GetDependencyOrder returns resources in dependency order (topological sort).
// Leaf dependencies (resources with no dependencies) come first.
func (r *DependencyResolver) GetDependencyOrder(names []string) ([]string, error) {
	// Build reverse dependency graph: for each resource, track what depends on it
	// graph[A] = [B, C] means B and C depend on A
	dependsOn := make(map[string][]string)  // resource -> what it depends on
	dependedBy := make(map[string][]string) // resource -> what depends on it
	inDegree := make(map[string]int)

	// Collect all nodes including transitive dependencies
	allNodes := make(map[string]bool)
	var collectNodes func(name string)
	collectNodes = func(name string) {
		if allNodes[name] {
			return
		}
		allNodes[name] = true
		deps := r.GetDependencies(name)
		for _, dep := range deps {
			collectNodes(dep.To)
		}
	}
	for _, name := range names {
		collectNodes(name)
	}

	// Initialize all nodes
	for name := range allNodes {
		dependsOn[name] = nil
		dependedBy[name] = nil
		inDegree[name] = 0
	}

	// Build the graph
	for name := range allNodes {
		deps := r.GetDependencies(name)
		for _, dep := range deps {
			if allNodes[dep.To] {
				dependsOn[name] = append(dependsOn[name], dep.To)
				dependedBy[dep.To] = append(dependedBy[dep.To], name)
				inDegree[name]++ // This node depends on something
			}
		}
	}

	// Topological sort using Kahn's algorithm
	// Start with nodes that have no dependencies (in-degree 0)
	var queue []string
	for name, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, name)
		}
	}

	var result []string
	for len(queue) > 0 {
		// Dequeue
		name := queue[0]
		queue = queue[1:]
		result = append(result, name)

		// For each node that depends on this one, reduce their in-degree
		for _, dependent := range dependedBy[name] {
			inDegree[dependent]--
			if inDegree[dependent] == 0 {
				queue = append(queue, dependent)
			}
		}
	}

	// Check for cycles
	if len(result) != len(allNodes) {
		// Find which nodes are in a cycle
		var inCycle []string
		for name, degree := range inDegree {
			if degree > 0 {
				inCycle = append(inCycle, name)
			}
		}
		return nil, fmt.Errorf("circular dependency detected involving: %v", inCycle)
	}

	return result, nil
}

// GetTransitiveDependencies returns all dependencies of a resource (recursive).
// The returned list is in dependency order - leaf resources first, then their dependents.
func (r *DependencyResolver) GetTransitiveDependencies(name string) []string {
	visited := make(map[string]bool)
	var result []string

	var visit func(n string)
	visit = func(n string) {
		if visited[n] {
			return
		}
		visited[n] = true

		// Visit dependencies first (post-order traversal)
		deps := r.GetDependencies(n)
		for _, dep := range deps {
			visit(dep.To)
		}

		result = append(result, n)
	}

	visit(name)
	return result
}

// GetDependents returns all resources that depend on the given resource.
func (r *DependencyResolver) GetDependents(name string) []string {
	var dependents []string

	for _, entry := range r.index.List() {
		deps := r.GetDependencies(entry.Name)
		for _, dep := range deps {
			if dep.To == name {
				dependents = append(dependents, entry.Name)
				break
			}
		}
	}

	return dependents
}

// HasCircularDependency checks if there is a circular dependency starting from the given resource.
func (r *DependencyResolver) HasCircularDependency(name string) bool {
	visiting := make(map[string]bool) // Currently in recursion stack
	visited := make(map[string]bool)  // Completely processed

	var hasCycle func(n string) bool
	hasCycle = func(n string) bool {
		if visiting[n] {
			return true // Back edge - cycle detected
		}
		if visited[n] {
			return false // Already processed, no cycle
		}

		visiting[n] = true
		deps := r.GetDependencies(n)
		for _, dep := range deps {
			if hasCycle(dep.To) {
				return true
			}
		}
		visiting[n] = false
		visited[n] = true

		return false
	}

	return hasCycle(name)
}

// DependencyGraph represents the complete dependency graph for visualization/debugging.
type DependencyGraph struct {
	Nodes map[string][]Dependency
}

// BuildDependencyGraph builds a complete dependency graph for all resources in the index.
func (r *DependencyResolver) BuildDependencyGraph() *DependencyGraph {
	graph := &DependencyGraph{
		Nodes: make(map[string][]Dependency),
	}

	for _, entry := range r.index.List() {
		deps := r.GetDependencies(entry.Name)
		graph.Nodes[entry.Name] = deps
	}

	return graph
}
