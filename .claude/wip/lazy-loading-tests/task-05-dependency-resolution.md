# Task 5: Dependency Resolution Edge Case Tests

## Objective

Write comprehensive tests for the dependency resolver, focusing on complex dependency graphs, circular detection, and loading order.

## Context

- Dependency resolver handles: children, references, query deps, categories, inputs
- Must detect circular dependencies without infinite loops
- Loading order affects correctness
- Cross-mod dependencies add complexity

## Dependencies

- Task 2 (test fixtures)
- Files to test: `internal/resourceloader/resolver.go`, `internal/resourceloader/loader.go`

## Acceptance Criteria

- [ ] Add tests to `internal/resourceloader/resolver_edge_test.go`
- [ ] Test complex dependency graphs
- [ ] Verify circular detection is robust
- [ ] Test cross-mod dependency resolution
- [ ] All edge cases documented

## Test Cases to Implement

### Complex Dependency Graphs
```go
// Test: Diamond dependency pattern
func TestResolver_DiamondDependency(t *testing.T)
// A → B, A → C, B → D, C → D
// D should be loaded once

// Test: Deep dependency chain
func TestResolver_DeepChain(t *testing.T)
// A → B → C → D → E → F → G (7+ levels)
// Loading order: G, F, E, D, C, B, A

// Test: Wide dependency (many deps)
func TestResolver_WideDependencies(t *testing.T)
// A depends on B1, B2, B3, ..., B50
// All loaded before A

// Test: Mixed width and depth
func TestResolver_MixedGraph(t *testing.T)
// Realistic benchmark hierarchy with controls and queries
```

### Circular Dependency Detection
```go
// Test: Simple cycle (A → B → A)
func TestResolver_SimpleCycle(t *testing.T)
// Should detect and return error

// Test: Long cycle (A → B → C → D → A)
func TestResolver_LongCycle(t *testing.T)
// Should detect before stack overflow

// Test: Cycle in subtree
func TestResolver_CycleInSubtree(t *testing.T)
// A → B → C, C → D → C (cycle not involving A)

// Test: Multiple cycles
func TestResolver_MultipleCycles(t *testing.T)
// A → B → A, C → D → C (independent cycles)

// Test: Self-referential resource
func TestResolver_SelfReference(t *testing.T)
// A → A (direct self-reference)

// Test: Cycle detection doesn't false positive
func TestResolver_NoCycleFalsePositive(t *testing.T)
// Complex DAG that looks like cycle but isn't
```

### Missing Dependencies
```go
// Test: Missing direct dependency
func TestResolver_MissingDirectDep(t *testing.T)
// Control refs missing query
// Should error gracefully

// Test: Missing transitive dependency
func TestResolver_MissingTransitiveDep(t *testing.T)
// A → B → C (C missing)
// Error should mention C

// Test: Optional dependency missing
func TestResolver_OptionalDepMissing(t *testing.T)
// Category reference missing
// Should skip gracefully

// Test: Inline resource as dependency
func TestResolver_InlineDependency(t *testing.T)
// Card defined inline in dashboard
// Should not try to load separately
```

### Cross-Mod Dependencies
```go
// Test: Reference to dependency mod resource
func TestResolver_CrossModReference(t *testing.T)
// Main mod control refs dep mod query
// Should resolve correctly

// Test: Cross-mod circular detection
func TestResolver_CrossModCycle(t *testing.T)
// Main.A → Dep.B → Main.C → Main.A

// Test: Mod name resolution
func TestResolver_ModNameResolution(t *testing.T)
// Short name vs full name resolution
```

### Dependency Types
```go
// Test: Child dependency resolution
func TestResolver_ChildDependencies(t *testing.T)
// Benchmark children loaded

// Test: Query reference dependency
func TestResolver_QueryReferenceDep(t *testing.T)
// Control's query loaded before control

// Test: Category dependency
func TestResolver_CategoryDep(t *testing.T)
// Graph categories loaded

// Test: Input dependency
func TestResolver_InputDependency(t *testing.T)
// Dashboard scoped inputs

// Test: Mixed dependency types
func TestResolver_MixedDepTypes(t *testing.T)
// Resource with children AND query AND category
```

### Loading Order
```go
// Test: Topological sort correctness
func TestResolver_TopologicalSort(t *testing.T)
// Verify each resource loaded after deps

// Test: Order stability
func TestResolver_OrderStability(t *testing.T)
// Same graph always produces same order

// Test: Parallel-safe order
func TestResolver_ParallelSafeOrder(t *testing.T)
// Order allows parallel loading of independent deps
```

### Performance
```go
// Test: Large graph resolution time
func TestResolver_LargeGraphPerformance(t *testing.T)
// 1000 node graph resolves in <1s

// Test: Memory usage on large graph
func TestResolver_LargeGraphMemory(t *testing.T)
// No excessive memory for visited tracking
```

### Error Messages
```go
// Test: Circular dependency error message
func TestResolver_CircularErrorMessage(t *testing.T)
// Should show cycle path: A → B → C → A

// Test: Missing dependency error message
func TestResolver_MissingDepErrorMessage(t *testing.T)
// Should name the missing resource and requester
```

## Test Fixture Requirements

Create dependency graph fixtures:
- `internal/testdata/mods/dependency-graphs/diamond/`
- `internal/testdata/mods/dependency-graphs/deep-chain/`
- `internal/testdata/mods/dependency-graphs/circular/`
- `internal/testdata/mods/dependency-graphs/missing-deps/`

Or generate programmatically using index directly.

## Implementation Notes

- Use index entries directly for unit tests (no file I/O)
- Create graph builder helper for complex tests
- Visualize graphs in test comments for clarity
- Test both `GetDependencyOrder()` and `ResolveWithDependencies()`

## Output Files

- `internal/resourceloader/resolver_edge_test.go`
