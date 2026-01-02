# Task 8: Cross-Mod Dependency Tests

## Objective

Write comprehensive tests for scenarios where resources in one mod reference resources in another mod (dependencies), focusing on lazy loading correctness.

## Context

- Mods can depend on other mods via `.powerpipe/mods/` directory
- Resources reference across mod boundaries
- Mod name mapping needed for resolution
- Scanner must handle dependency mod structure
- This is where version conflicts could appear

## Dependencies

- Task 2 (multi-mod test fixture)
- Tasks 3-6 (unit tests should pass)
- Files to test: `internal/workspace/lazy_workspace.go`, `internal/resourceindex/scanner.go`

## Acceptance Criteria

- [ ] Add tests to `internal/workspace/cross_mod_test.go`
- [ ] Test cross-mod resource resolution
- [ ] Verify scanner handles dependency mod structure
- [ ] Test mod name collision handling
- [ ] Verify payload includes all mods correctly

## Test Cases to Implement

### Cross-Mod Resource Resolution
```go
// Test: Main mod control references dep mod query
func TestCrossMod_ControlRefsDepModQuery(t *testing.T)
// Main: control.test refs dep_mod.query.helper
// Resolution succeeds
// Execution uses correct query

// Test: Main mod benchmark includes dep mod benchmark
func TestCrossMod_BenchmarkIncludesDepBenchmark(t *testing.T)
// Main: benchmark.all has child dep_mod.benchmark.subset
// Children resolved correctly

// Test: Main mod dashboard uses dep mod card
func TestCrossMod_DashboardUsesDepModCard(t *testing.T)
// Dashboard references external card
// Card loaded from dep mod

// Test: Transitive cross-mod reference
func TestCrossMod_TransitiveReference(t *testing.T)
// Main → Dep1 → Dep2
// Three-level resolution works
```

### Scanner Cross-Mod Handling
```go
// Test: Scanner discovers dependency mods
func TestCrossMod_ScannerDiscoversMods(t *testing.T)
// Workspace with .powerpipe/mods/
// All mods scanned and indexed

// Test: Mod name extracted from nested path
func TestCrossMod_ModNameFromPath(t *testing.T)
// Path: .powerpipe/mods/github.com/org/mod@v1.2.3/
// Mod name: mod (or full path?)

// Test: Version stripped from path
func TestCrossMod_VersionStripping(t *testing.T)
// @v1.2.3 stripped for mod registration
// No version conflicts

// Test: Resources indexed with correct mod name
func TestCrossMod_ResourceModNameCorrect(t *testing.T)
// Dep mod resource's mod_full_name set correctly
// References resolve using this name
```

### Mod Name Mapping
```go
// Test: RegisterModName called for deps
func TestCrossMod_ModNameRegistration(t *testing.T)
// Verify mod name mapping populated
// Both full path and short name work

// Test: Short name resolution
func TestCrossMod_ShortNameResolution(t *testing.T)
// Reference by short name (aws_insights)
// Resolves to correct mod

// Test: Full name resolution
func TestCrossMod_FullNameResolution(t *testing.T)
// Reference by full name (github.com/turbot/...)
// Resolves correctly

// Test: Name collision detection
func TestCrossMod_NameCollision(t *testing.T)
// Two mods with same short name
// Detection or deterministic choice
```

### Payload Generation
```go
// Test: Available dashboards includes all mods
func TestCrossMod_AvailableDashboardsIncludesAllMods(t *testing.T)
// Payload has dashboards from main and deps
// Each has correct mod_full_name

// Test: Benchmark hierarchy spans mods
func TestCrossMod_BenchmarkHierarchySpansMods(t *testing.T)
// Main benchmark with dep mod children
// Hierarchy built correctly
// Trunks include cross-mod paths

// Test: Mod list in server metadata
func TestCrossMod_ModListInMetadata(t *testing.T)
// Server metadata lists all installed mods
// Main mod and dependencies
```

### Lazy Loading Cross-Mod
```go
// Test: Lazy load cross-mod resource
func TestCrossMod_LazyLoadCrossModResource(t *testing.T)
// Request resource from dep mod
// Loaded on-demand correctly

// Test: Cross-mod dependency resolution in lazy mode
func TestCrossMod_LazyDependencyResolution(t *testing.T)
// Control in main needs query in dep
// Query loaded when control accessed

// Test: Index contains cross-mod references
func TestCrossMod_IndexContainsCrossModRefs(t *testing.T)
// Index entry for control has dep mod query ref
// Reference name fully qualified
```

### Error Handling
```go
// Test: Missing dependency mod
func TestCrossMod_MissingDepMod(t *testing.T)
// Reference to non-existent mod
// Graceful error

// Test: Circular cross-mod dependency
func TestCrossMod_CircularCrossModDep(t *testing.T)
// Main → Dep → Main (cycle)
// Should this be allowed?

// Test: Invalid dep mod structure
func TestCrossMod_InvalidDepModStructure(t *testing.T)
// Dep mod without mod.pp
// Handled gracefully
```

### Execution Cross-Mod
```go
// Test: Execute benchmark with cross-mod controls
func TestCrossMod_ExecuteBenchmarkWithCrossModControls(t *testing.T)
// Benchmark children from dep mod
// All controls execute

// Test: Execute dashboard with cross-mod references
func TestCrossMod_ExecuteDashboardWithCrossModRefs(t *testing.T)
// Dashboard using dep mod queries
// All queries resolve and execute
```

### Edge Cases
```go
// Test: Same resource name in different mods
func TestCrossMod_SameNameDifferentMods(t *testing.T)
// main.query.test and dep.query.test
// Correct one selected based on context

// Test: Deeply nested dependency mods
func TestCrossMod_DeeplyNestedDeps(t *testing.T)
// mod/mods/dep1/mods/dep2/mods/dep3
// All levels discovered

// Test: Self-referential mod (reference own resources)
func TestCrossMod_SelfReference(t *testing.T)
// Mod references own resources
// Works correctly (baseline)
```

## Test Fixture Requirements

Create multi-mod fixture:
```
internal/testdata/mods/multi-mod/
├── mod.pp                          # Main mod
├── controls.pp                     # Main controls (some ref dep)
├── benchmarks.pp                   # Main benchmarks (some include dep)
├── dashboards.pp                   # Main dashboards (some use dep)
└── .powerpipe/mods/
    ├── github.com/test/dep-mod@v1.0.0/
    │   ├── mod.pp
    │   ├── queries.pp              # Queries referenced by main
    │   ├── controls.pp             # Controls included by main
    │   └── benchmarks.pp           # Benchmarks included by main
    └── github.com/test/dep-mod-2@v2.0.0/
        ├── mod.pp
        └── utilities.pp            # Second dep mod
```

## Implementation Notes

- Test both eager and lazy modes
- Verify index structure matches expectations
- Check that execution works end-to-end
- Consider testing with real steampipe mods (subset)

## Output Files

- `internal/workspace/cross_mod_test.go`
