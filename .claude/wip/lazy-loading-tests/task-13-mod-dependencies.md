# Task 13: Mod Dependency Edge Case Tests (CRITICAL)

## Objective

Write comprehensive tests for mod dependency handling, focusing on the `require` block parsing, dependency discovery, version resolution, transitive dependencies, and conflict handling. This is a critical area with known complexity.

## Context

Based on deep code analysis, mod dependencies involve:
- `require` blocks in mod.pp parsed by pipe-fittings
- Dependency mods discovered in `.powerpipe/mods/`
- Version strings with semver constraints
- Transitive dependency chains
- **Silent failure on missing dependencies** (HIGH RISK)
- Version conflict resolution via replacement (last wins)

## Dependencies

- Task 2 (multi-mod test fixtures)
- Files to test across multiple packages

## Acceptance Criteria

- [ ] Add tests to `internal/workspace/mod_dependency_test.go`
- [ ] Test require block parsing edge cases
- [ ] Test dependency discovery mechanisms
- [ ] Test version constraint handling
- [ ] Test transitive dependency resolution
- [ ] Test missing dependency behavior
- [ ] Test version conflict scenarios
- [ ] All tests pass with `-race` flag

## Critical Risk Areas Identified

| Risk | Code Location | Issue |
|------|---------------|-------|
| **Silent missing deps** | resolver.go:72-79 | Missing deps skipped without warning |
| **Late cycle detection** | resolver.go:52-59 | Cycles only caught during load |
| **Version conflicts** | require.go:170-192 | Last-wins replacement, no merge |
| **Mod name collision** | index.go:42-58 | Full path → short name mapping |

## Test Cases to Implement

### Require Block Parsing
```go
// Test: Basic require block with mods
func TestModDeps_RequireBlockBasic(t *testing.T)
// require {
//   mod "github.com/turbot/steampipe-mod-aws-insights" {
//     version = "0.21"
//   }
// }

// Test: Multiple mod dependencies
func TestModDeps_RequireBlockMultipleMods(t *testing.T)
// require {
//   mod "github.com/turbot/mod-a" { version = "1.0" }
//   mod "github.com/turbot/mod-b" { version = "2.0" }
// }

// Test: Require with version constraint
func TestModDeps_VersionConstraint(t *testing.T)
// version = ">=1.0.0"
// version = "^1.2.3"
// version = "~1.2.3"

// Test: Require with branch instead of version
func TestModDeps_BranchDependency(t *testing.T)
// branch = "main"

// Test: Require with local path
func TestModDeps_LocalPathDependency(t *testing.T)
// path = "../other-mod"

// Test: Require with tag (non-semver)
func TestModDeps_TagDependency(t *testing.T)
// tag = "release-2024-01"

// Test: Invalid - multiple version specifiers
func TestModDeps_InvalidMultipleSpecifiers(t *testing.T)
// version = "1.0" AND branch = "main" (should error)

// Test: Empty require block
func TestModDeps_EmptyRequireBlock(t *testing.T)
// require { }

// Test: Deprecated steampipe property
func TestModDeps_DeprecatedSteampipeProperty(t *testing.T)
// require { steampipe = "0.20.0" }
// Should convert to block with warning
```

### Dependency Mod Discovery
```go
// Test: Discover mods in .powerpipe/mods
func TestModDeps_DiscoveryBasic(t *testing.T)
// .powerpipe/mods/github.com/org/mod@v1.0.0/mod.pp
// Should be discovered and indexed

// Test: Discover nested dependency mods
func TestModDeps_DiscoveryNested(t *testing.T)
// .powerpipe/mods/mod-a@v1/
//   └── .powerpipe/mods/mod-b@v2/
// Both should be discovered

// Test: Discovery with missing mod.pp
func TestModDeps_DiscoveryMissingModFile(t *testing.T)
// Directory exists but no mod.pp
// Should skip gracefully

// Test: Discovery with invalid mod.pp
func TestModDeps_DiscoveryInvalidModFile(t *testing.T)
// mod.pp has syntax errors
// Should skip with warning

// Test: Discovery with symlinks
func TestModDeps_DiscoverySymlinks(t *testing.T)
// Symlinked mod directories
// Should follow and discover

// Test: Empty .powerpipe/mods directory
func TestModDeps_DiscoveryEmptyModsDir(t *testing.T)
// Directory exists but is empty
// No error

// Test: No .powerpipe/mods directory
func TestModDeps_DiscoveryNoModsDir(t *testing.T)
// Directory doesn't exist
// Graceful handling
```

### Version Path Handling
```go
// Test: Extract version from path
func TestModDeps_VersionFromPath(t *testing.T)
// github.com/turbot/mod@v1.2.3 → version = "1.2.3"

// Test: Handle prerelease versions
func TestModDeps_PrereleaseVersion(t *testing.T)
// @v1.2.3-beta.1

// Test: Handle build metadata
func TestModDeps_BuildMetadataVersion(t *testing.T)
// @v1.2.3+build.456

// Test: Multiple @ symbols (edge case)
func TestModDeps_MultipleAtSymbols(t *testing.T)
// github.com/org@name/mod@v1.0.0
// Should handle correctly

// Test: No version in path
func TestModDeps_NoVersionInPath(t *testing.T)
// github.com/turbot/mod (no @version)
// Should work as local/dev mod

// Test: Version with v prefix variations
func TestModDeps_VersionPrefixVariations(t *testing.T)
// @v1.0.0 vs @1.0.0 (with/without v)
```

### Mod Name Mapping
```go
// Test: Full path to short name mapping
func TestModDeps_NameMapping(t *testing.T)
// github.com/turbot/steampipe-mod-aws-insights → aws_insights

// Test: Name mapping with hyphens
func TestModDeps_NameMappingHyphens(t *testing.T)
// my-mod-name → my_mod_name

// Test: Name mapping collision detection
func TestModDeps_NameMappingCollision(t *testing.T)
// Two mods that would map to same short name
// mod-name and mod_name both → mod_name?

// Test: Resolve by short name
func TestModDeps_ResolveByShortName(t *testing.T)
// Reference: aws_insights.control.test
// Resolves to correct mod

// Test: Resolve by full path
func TestModDeps_ResolveByFullPath(t *testing.T)
// Reference: github.com/turbot/mod.control.test
// Resolves correctly

// Test: Ambiguous name resolution
func TestModDeps_AmbiguousNameResolution(t *testing.T)
// Two mods with similar names
// Resolution should be deterministic
```

### Transitive Dependencies
```go
// Test: Two-level transitive dependency
func TestModDeps_TransitiveTwoLevel(t *testing.T)
// Main → DepA → DepB
// DepB resources accessible from Main

// Test: Three-level transitive dependency
func TestModDeps_TransitiveThreeLevel(t *testing.T)
// Main → A → B → C
// All resources accessible

// Test: Diamond dependency pattern
func TestModDeps_DiamondDependency(t *testing.T)
// Main → A, Main → B, A → C, B → C
// C loaded once, accessible via both paths

// Test: Wide transitive (many deps at one level)
func TestModDeps_TransitiveWide(t *testing.T)
// Main → [A, B, C, D, E, F, G, H, I, J]
// All 10 dependencies loaded

// Test: Transitive with version differences
func TestModDeps_TransitiveVersionDiff(t *testing.T)
// Main → A@1.0, Main → B, B → A@2.0
// Which version of A is used?

// Test: Transitive dependency ordering
func TestModDeps_TransitiveOrdering(t *testing.T)
// Dependencies loaded in correct order
// Leaf dependencies first
```

### Missing Dependency Handling (CRITICAL)
```go
// Test: Missing direct dependency
func TestModDeps_MissingDirectDep(t *testing.T)
// require { mod "nonexistent" { version = "1.0" } }
// What error is shown?

// Test: Missing transitive dependency
func TestModDeps_MissingTransitiveDep(t *testing.T)
// Main → A, A requires B, B not installed
// Error should identify missing B

// Test: Silent skip in lazy loading
func TestModDeps_SilentSkipLazyMode(t *testing.T)
// Missing dependency in lazy mode
// Currently skipped silently - verify behavior
// Should at least log warning

// Test: Reference to missing mod resource
func TestModDeps_ReferenceMissingModResource(t *testing.T)
// Control refs dep_mod.query.nonexistent
// Clear error message

// Test: Partial dependency availability
func TestModDeps_PartialDependency(t *testing.T)
// Some deps available, some missing
// Available ones work, missing ones error clearly

// Test: Error message includes mod name
func TestModDeps_ErrorIncludesModName(t *testing.T)
// Missing dependency error
// Should say which mod is missing
```

### Version Conflict Scenarios
```go
// Test: Same mod, same version, multiple requires
func TestModDeps_SameModSameVersion(t *testing.T)
// Both A and B require C@1.0.0
// No conflict, C loaded once

// Test: Same mod, different versions
func TestModDeps_SameModDifferentVersions(t *testing.T)
// A requires C@1.0.0, B requires C@2.0.0
// How is this resolved?

// Test: Version constraint satisfaction
func TestModDeps_VersionConstraintSatisfied(t *testing.T)
// Require >=1.0.0, installed 1.2.0
// Should work

// Test: Version constraint not satisfied
func TestModDeps_VersionConstraintNotSatisfied(t *testing.T)
// Require >=2.0.0, installed 1.2.0
// Should error with clear message

// Test: Last-wins replacement behavior
func TestModDeps_LastWinsReplacement(t *testing.T)
// Multiple require blocks for same mod
// Later one should win
// Document this behavior

// Test: Conflicting constraints
func TestModDeps_ConflictingConstraints(t *testing.T)
// Require >=2.0.0 and <1.5.0 for same mod
// Impossible to satisfy - error message?
```

### Circular Dependencies
```go
// Test: Direct circular between mods
func TestModDeps_CircularDirect(t *testing.T)
// ModA requires ModB, ModB requires ModA
// Should be detected

// Test: Indirect circular through resources
func TestModDeps_CircularIndirect(t *testing.T)
// ModA.control refs ModB.query
// ModB.benchmark includes ModA.control
// Circular at resource level

// Test: Circular detection timing
func TestModDeps_CircularDetectionTiming(t *testing.T)
// When is circular detected?
// At scan time? Load time? Execution time?

// Test: Circular error message
func TestModDeps_CircularErrorMessage(t *testing.T)
// Shows cycle path: A → B → C → A

// Test: Long circular chain
func TestModDeps_CircularLongChain(t *testing.T)
// A → B → C → D → E → A
// Detected without stack overflow
```

### Integration with Lazy Loading
```go
// Test: Lazy loading discovers all dep mods
func TestModDeps_LazyDiscoveryComplete(t *testing.T)
// Index contains resources from all dep mods
// Nothing missed

// Test: Lazy to eager transition with deps
func TestModDeps_LazyEagerTransitionWithDeps(t *testing.T)
// Lazy mode indexes deps
// Eager mode loads deps fully
// Consistent results

// Test: Cross-mod resource access in lazy mode
func TestModDeps_LazyCrossModAccess(t *testing.T)
// GetResource for dep mod resource
// Loaded on demand correctly

// Test: Available dashboards includes dep mods
func TestModDeps_AvailableDashboardsIncludesDeps(t *testing.T)
// Dashboard list shows dep mod dashboards
// With correct mod_full_name

// Test: Benchmark hierarchy crosses mods
func TestModDeps_BenchmarkHierarchyCrossesMods(t *testing.T)
// Main benchmark includes dep mod benchmark
// Hierarchy built correctly
// Trunks span mods
```

### Real-World Patterns
```go
// Test: AWS Insights style dependency
func TestModDeps_AWSInsightsPattern(t *testing.T)
// Main mod with many dashboards
// References shared queries mod
// Common real-world pattern

// Test: Compliance framework pattern
func TestModDeps_ComplianceFrameworkPattern(t *testing.T)
// Framework mod with benchmarks
// Plugin-specific mods with controls
// Cross-mod benchmark assembly

// Test: Local development pattern
func TestModDeps_LocalDevPattern(t *testing.T)
// Main mod with path dependency
// path = "../shared-queries"
// Development workflow
```

## Test Fixture Requirements

Create comprehensive multi-mod fixture:
```
internal/testdata/mods/mod-dependencies/
├── main-mod/
│   ├── mod.pp                    # require { mod "dep-a" {...} }
│   ├── controls.pp               # refs dep-a.query.helper
│   └── benchmarks.pp             # includes dep-a.benchmark.subset
│
├── dep-mods/
│   ├── dep-a@v1.0.0/
│   │   ├── mod.pp
│   │   ├── queries.pp
│   │   └── benchmarks.pp
│   │
│   ├── dep-b@v1.0.0/             # Depends on dep-a
│   │   ├── mod.pp                # require { mod "dep-a" {...} }
│   │   └── controls.pp
│   │
│   ├── dep-conflicting@v1.0.0/   # Version conflict test
│   │   └── mod.pp
│   │
│   └── dep-conflicting@v2.0.0/   # Different version
│       └── mod.pp
│
├── circular-mods/
│   ├── circular-a/
│   │   └── mod.pp                # requires circular-b
│   └── circular-b/
│       └── mod.pp                # requires circular-a
│
└── missing-dep-mod/
    └── mod.pp                    # requires nonexistent mod
```

## Implementation Notes

- Test both scanning/indexing and full loading paths
- Verify error messages are helpful and specific
- Check that warnings are logged for silent skips
- Test with real steampipe mod structures if possible
- Profile memory for large dependency trees

## Output Files

- `internal/workspace/mod_dependency_test.go`
- `internal/resourceindex/mod_discovery_test.go`
- `internal/resourceloader/transitive_deps_test.go`
