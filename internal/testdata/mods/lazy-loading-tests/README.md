# Lazy Loading Test Fixtures

This directory contains test mod fixtures specifically designed for testing lazy loading functionality in Powerpipe.

## Fixture Categories

### 1. `simple/`
**Purpose:** Baseline testing with clear, simple resources.

**Contents:**
- 3 queries (simple_count, simple_status, control_query)
- 2 controls (inline_sql, uses_query)
- 1 benchmark with 2 controls
- 1 dashboard with cards

**Use Cases:**
- Basic lazy loading sanity tests
- Verifying resource loading works at all
- Reference for expected behavior

### 2. `deep-hierarchy/`
**Purpose:** Test deep nesting (11 levels) for recursion and trunk building.

**Contents:**
- 1 control at the leaf
- 11 benchmarks (root → level_1 → ... → level_10 → control)
- Additional branching hierarchy

**Use Cases:**
- Stack overflow detection in recursive algorithms
- Trunk/path building verification
- Deep nested execution

### 3. `wide-hierarchy/`
**Purpose:** Test wide benchmarks (100+ children) for performance and memory.

**Contents:**
- 100 controls
- 1 benchmark containing all 100 controls

**Use Cases:**
- Memory efficiency testing
- Child iteration performance
- Large children array handling

### 4. `cross-refs/`
**Purpose:** Complex reference patterns between resources.

**Contents:**
- Variables referenced by params
- Locals referenced by queries
- Queries referenced by controls
- Controls referenced by multiple benchmarks
- Dashboards with nested containers referencing queries

**Use Cases:**
- Reference resolution testing
- Dependency graph building
- Cross-reference lazy loading

### 5. `edge-cases/`
**Purpose:** Unusual but valid configurations.

**Contents:**
- `unicode-names.pp` - Resources with emojis, Chinese, Arabic, accented chars
- `special-chars.pp` - Names with underscores, version patterns, numbers
- `empty-resources.pp` - Minimal resources, empty children
- `huge-sql.pp` - Very large heredoc SQL statements
- `deeply-nested-json.pp` - Deep JSON structures in params/locals
- `multiline-strings.pp` - Complex heredoc and multiline formatting
- `whitespace-variations.pp` - Various indentation patterns

**Use Cases:**
- Parser robustness testing
- Edge case handling
- Character encoding

### 6. `error-conditions/`
**Purpose:** Files designed to trigger specific error paths.

**Contents:**
- `missing-query-ref.pp` - Controls referencing non-existent queries
- `circular-benchmark.pp` - Circular benchmark references (A → B → C → A)
- `invalid-child-type.pp` - Benchmarks with invalid child types
- `duplicate-names.pp` - Same resource name defined twice
- `malformed-partial.pp` - Semantically invalid configurations

**Use Cases:**
- Error path testing
- Error message verification
- Graceful failure handling

### 7. `multi-mod/`
**Purpose:** Simulate dependency mod structure.

**Structure:**
```
multi-mod/
└── main/
    ├── mod.pp (requires dep-mod)
    └── .powerpipe/mods/github.com/test/dep-mod@v1.0.0/
        ├── mod.pp
        ├── queries.pp
        └── controls.pp
```

**Use Cases:**
- Cross-mod reference resolution
- Dependency loading
- Mod namespace handling

### 8. `generated/`
**Purpose:** Programmatically generated mods at various scales.

**Sizes:**
- `small/` - 10 dashboards, 20 queries, 30 controls, 5 benchmarks
- `medium/` - 50 dashboards, 100 queries, 150 controls, 25 benchmarks
- `large/` - 200 dashboards, 400 queries, 600 controls, 75 benchmarks
- `stress/` - (generate locally) 500 dashboards, 1000 queries, 2000 controls

**Generation:**
```bash
# Generate stress mod locally (not committed to repo)
go run scripts/generate_lazy_test_mods.go internal/testdata/mods/lazy-loading-tests/generated/stress stress
```

**Use Cases:**
- Performance benchmarking
- Memory usage testing
- Scale testing

## Usage in Tests

```go
package yourpackage_test

import (
    "testing"
    "path/filepath"
)

const lazyTestModsPath = "internal/testdata/mods/lazy-loading-tests"

func TestSimpleLoading(t *testing.T) {
    modPath := filepath.Join(lazyTestModsPath, "simple")
    // ... test code
}

func TestDeepHierarchy(t *testing.T) {
    modPath := filepath.Join(lazyTestModsPath, "deep-hierarchy")
    // ... test code
}

func TestWideHierarchy(t *testing.T) {
    modPath := filepath.Join(lazyTestModsPath, "wide-hierarchy")
    // ... test code
}

func TestErrorConditions(t *testing.T) {
    modPath := filepath.Join(lazyTestModsPath, "error-conditions")
    // Expect specific errors when loading this mod
}
```

## Notes

- All fixtures are valid HCL (parseable by the scanner)
- Error condition fixtures are valid HCL but should fail during resolution
- Fixtures work with both eager and lazy loading
- Generated `stress/` mods are gitignored to keep repo size reasonable
