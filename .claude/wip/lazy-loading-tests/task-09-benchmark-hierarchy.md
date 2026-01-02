# Task 9: Benchmark Hierarchy Tests

## Objective

Write comprehensive tests for benchmark hierarchy handling, focusing on the trunk building logic, parent-child relationships, and deep/wide hierarchies.

## Context

- Benchmarks have parent-child relationships
- "Trunks" track paths from top-level to each benchmark
- Recent bugs in this area (mod_full_name, trunks missing)
- Deep hierarchies stress recursion
- Wide hierarchies stress memory/performance

## Dependencies

- Task 2 (hierarchy test fixtures)
- Tasks 3-6 (unit tests should pass)
- Files to test: `internal/resourceindex/payload.go`, `internal/workspace/lazy_mod_resources.go`

## Acceptance Criteria

- [ ] Add tests to `internal/resourceindex/hierarchy_test.go`
- [ ] Test deep hierarchies (10+ levels)
- [ ] Test wide hierarchies (100+ children)
- [ ] Verify trunk building correctness
- [ ] Test parent-child relationship propagation

## Test Cases to Implement

### Basic Hierarchy
```go
// Test: Single benchmark with controls
func TestHierarchy_SingleBenchmarkWithControls(t *testing.T)
// Benchmark has 3 control children
// Verify children list correct
// Verify controls have parent set

// Test: Two-level benchmark hierarchy
func TestHierarchy_TwoLevelBenchmark(t *testing.T)
// Top → Child → Controls
// Verify relationships both directions

// Test: Sibling benchmarks
func TestHierarchy_SiblingBenchmarks(t *testing.T)
// Top → Child1, Child2
// Both children have same parent
// Parent has both children
```

### Deep Hierarchies
```go
// Test: 10-level deep hierarchy
func TestHierarchy_Deep10Levels(t *testing.T)
// B0 → B1 → B2 → ... → B9 → Control
// All relationships correct
// No stack overflow

// Test: 20-level deep hierarchy
func TestHierarchy_Deep20Levels(t *testing.T)
// Stress test for deep recursion
// Should handle without issues

// Test: Deep trunk building
func TestHierarchy_DeepTrunkBuilding(t *testing.T)
// 10-level hierarchy
// Bottom benchmark has trunk: [[B0], [B0,B1], ..., [B0,...,B9]]
// Trunk length matches depth
```

### Wide Hierarchies
```go
// Test: 100 children benchmark
func TestHierarchy_Wide100Children(t *testing.T)
// Benchmark with 100 control children
// All children tracked
// Performance acceptable

// Test: 500 children benchmark
func TestHierarchy_Wide500Children(t *testing.T)
// Stress test for wide benchmarks
// Memory usage bounded

// Test: Mixed width at different levels
func TestHierarchy_MixedWidth(t *testing.T)
// Top: 10 children
// Each child: 20 grandchildren
// 200 total leaf controls
```

### Trunk Building
```go
// Test: Top-level benchmark trunk
func TestHierarchy_TopLevelTrunk(t *testing.T)
// Top-level benchmark
// Trunk: [[benchmark_name]]

// Test: Child benchmark trunk includes parent
func TestHierarchy_ChildTrunkIncludesParent(t *testing.T)
// Parent → Child
// Child trunk: [[Parent, Child]]

// Test: Multiple trunks for diamond pattern
func TestHierarchy_DiamondPatternTrunks(t *testing.T)
// Top → A, Top → B, A → Child, B → Child
// Child should have two trunks:
// [[Top, A, Child], [Top, B, Child]]

// Test: Trunk order consistency
func TestHierarchy_TrunkOrderConsistent(t *testing.T)
// Multiple runs produce same trunk order
// Deterministic output
```

### Parent-Child Relationships
```go
// Test: SetParentNames populates correctly
func TestHierarchy_SetParentNames(t *testing.T)
// After scanning, children have parent set
// Scanner post-processing works

// Test: Parent-child bidirectional
func TestHierarchy_BidirectionalRelationship(t *testing.T)
// Parent.Children contains Child
// Child.ParentName == Parent.Name

// Test: Orphan detection
func TestHierarchy_OrphanDetection(t *testing.T)
// Benchmark not in any children list
// Should be marked as top-level

// Test: Child not in index
func TestHierarchy_MissingChild(t *testing.T)
// Benchmark lists child that doesn't exist
// Graceful handling
```

### Top-Level Detection
```go
// Test: Correctly identifies top-level
func TestHierarchy_CorrectTopLevel(t *testing.T)
// Only truly top-level marked
// Children not marked top-level

// Test: Multiple top-level benchmarks
func TestHierarchy_MultipleTopLevel(t *testing.T)
// Several independent benchmark trees
// Each root is top-level

// Test: Flat benchmarks (no hierarchy)
func TestHierarchy_FlatBenchmarks(t *testing.T)
// Benchmarks with only control children
// All are top-level (no benchmark children)
```

### Detection Benchmarks
```go
// Test: Detection benchmark hierarchy
func TestHierarchy_DetectionBenchmarks(t *testing.T)
// detection_benchmark type
// Similar hierarchy handling

// Test: Mixed benchmark types
func TestHierarchy_MixedBenchmarkTypes(t *testing.T)
// Regular benchmarks and detection benchmarks
// Separate handling but consistent
```

### Payload Generation
```go
// Test: BenchmarkInfo structure
func TestHierarchy_BenchmarkInfoStructure(t *testing.T)
// Verify all fields populated:
// FullName, Title, ModFullName, Trunks, Children, IsTopLevel

// Test: Recursive children in payload
func TestHierarchy_PayloadRecursiveChildren(t *testing.T)
// 3-level hierarchy
// Payload has nested BenchmarkInfo

// Test: Payload from index matches eager
func TestHierarchy_PayloadMatchesEager(t *testing.T)
// Build from index
// Build from eager workspace
// Same structure

// Test: Large hierarchy payload performance
func TestHierarchy_LargePayloadPerformance(t *testing.T)
// 500 benchmarks, 2000 controls
// Payload builds in <1s
```

### Edge Cases
```go
// Test: Empty children array
func TestHierarchy_EmptyChildren(t *testing.T)
// Benchmark with children = []
// Handled correctly (no children)

// Test: Single-item children
func TestHierarchy_SingleChild(t *testing.T)
// Benchmark with one child
// Arrays of length 1 handled

// Test: Duplicate in children
func TestHierarchy_DuplicateChild(t *testing.T)
// Same child listed twice
// Deduplicated or error?

// Test: Circular hierarchy reference
func TestHierarchy_CircularHierarchy(t *testing.T)
// A → B → A (should be prevented)
// Detected during trunk building
```

### Index vs Full Parse Consistency
```go
// Test: Hierarchy from index matches parse
func TestHierarchy_IndexMatchesParse(t *testing.T)
// Build hierarchy from scanner
// Build from full parse
// Compare: same children, same parents

// Test: Trunk from index matches parse
func TestHierarchy_TrunkMatchesParse(t *testing.T)
// Trunks built from index
// Trunks built from eager
// Must match
```

## Test Fixture Requirements

Create hierarchy fixtures:
```
internal/testdata/mods/hierarchies/
├── simple/           # Basic 2-3 level hierarchy
├── deep/             # 10+ level depth
├── wide/             # 100+ children
├── diamond/          # Multiple paths to same node
└── mixed/            # Combination of patterns
```

## Implementation Notes

- Use generated hierarchies for scale tests
- Visualize expected trunk structure in comments
- Test both index-based and full-parse paths
- Profile memory for wide hierarchies

## Output Files

- `internal/resourceindex/hierarchy_test.go`
- `internal/resourceindex/payload_test.go` (additions)
