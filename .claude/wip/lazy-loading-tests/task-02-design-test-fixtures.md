# Task 2: Design Comprehensive Test Mod Fixtures

## Objective

Create a rich set of test mod fixtures specifically designed to exercise lazy loading edge cases, hierarchies, cross-references, and error conditions.

## Context

- Existing test mods are relatively simple
- Need fixtures that stress edge cases
- Fixtures should be reusable across multiple test files
- Need both static fixtures and generation scripts

## Dependencies

- Task 1 (gap analysis informs fixture design)

## Acceptance Criteria

- [x] Create `/internal/testdata/mods/lazy-loading-tests/` directory structure
- [x] Design and implement test mods for each category below
- [x] Document each fixture's purpose and usage
- [x] Ensure fixtures work with both eager and lazy loading
- [x] Create generation script for scale testing

## Fixture Categories

### 1. Simple Reference (`simple/`)
Basic mod with clear, simple resources for baseline testing.
```
simple/
├── mod.pp
├── queries.pp      # 3 simple queries
├── controls.pp     # 2 controls (one refs query)
├── benchmark.pp    # 1 benchmark with 2 controls
└── dashboard.pp    # 1 dashboard with cards
```

### 2. Deep Hierarchy (`deep-hierarchy/`)
Test deep nesting (10+ levels) for recursion and trunk building.
```
deep-hierarchy/
├── mod.pp
└── benchmarks.pp   # benchmark_0 → benchmark_1 → ... → benchmark_10 → control
```

### 3. Wide Hierarchy (`wide-hierarchy/`)
Test wide benchmarks (100+ children) for performance and memory.
```
wide-hierarchy/
├── mod.pp
├── controls.pp     # 100 controls
└── benchmark.pp    # 1 benchmark with 100 children
```

### 4. Cross References (`cross-refs/`)
Resources referencing each other in complex patterns.
```
cross-refs/
├── mod.pp
├── queries.pp      # Queries with param refs
├── controls.pp     # Controls referencing queries, other controls
├── dashboards.pp   # Dashboards with nested containers
└── README.md       # Diagram of reference graph
```

### 5. Edge Cases (`edge-cases/`)
Unusual but valid configurations.
```
edge-cases/
├── mod.pp
├── unicode-names.pp        # Resources with unicode in names/titles
├── special-chars.pp        # Names with dots, dashes, underscores
├── empty-resources.pp      # Empty children arrays, no descriptions
├── huge-sql.pp             # Very large heredoc SQL (100KB)
├── deeply-nested-json.pp   # Deep JSON in params
├── multiline-strings.pp    # Complex string formatting
└── whitespace-variations.pp # Different indentation patterns
```

### 6. Error Conditions (`error-conditions/`)
Files designed to trigger specific error paths.
```
error-conditions/
├── mod.pp
├── missing-query-ref.pp    # Control refs non-existent query
├── circular-benchmark.pp   # A → B → C → A
├── invalid-child-type.pp   # Benchmark with invalid child type
├── duplicate-names.pp      # Same resource name twice
└── malformed-partial.pp    # Syntactically valid but semantically wrong
```

### 7. Multi-Mod (`multi-mod/`)
Simulate dependency mod structure.
```
multi-mod/
├── main/
│   ├── mod.pp
│   └── .powerpipe/mods/
│       └── github.com/test/dep-mod@v1.0.0/
│           ├── mod.pp
│           ├── queries.pp
│           └── controls.pp
└── README.md
```

### 8. Generated (`generated/`)
Programmatically generated mods at various scales.
```
generated/
├── generate.go           # Generation script
├── small/                # 10 dashboards, 20 controls
├── medium/               # 100 dashboards, 200 controls
├── large/                # 500 dashboards, 1000 controls
└── stress/               # 2000 dashboards (for stress testing)
```

## Generation Script Requirements

Create `scripts/generate_lazy_test_mods.go`:
- Accept parameters: dashboard count, control count, benchmark depth
- Generate deterministic output (seeded random)
- Include realistic metadata (titles, descriptions, tags)
- Create varied reference patterns

## Documentation

Each fixture directory should contain a `README.md` explaining:
- Purpose of the fixture
- What edge cases it tests
- Expected resource counts
- How to use in tests

## Notes

- Fixtures must be valid HCL (parseable by both scanner and full parser)
- Error condition fixtures should fail predictably
- Consider CI/CD: fixtures shouldn't be too large (keep repo reasonable)
- Use `.gitignore` for generated stress mods if needed
