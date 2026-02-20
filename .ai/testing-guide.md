# Testing Guide

## Unit Tests (Go)

### Running

```bash
go test ./...                              # All unit tests
go test ./internal/controldisplay          # Specific package
go test ./internal/resources -run TestName # Specific test
go test -v -cover ./...                    # Verbose with coverage
```

### Test Files (17 total)

```
internal/snapshot/snapshot_tag_test.go
internal/controldisplay/group_counter_graph_test.go
internal/controldisplay/formatter_template_test.go
internal/controldisplay/template_functions_test.go
internal/controldisplay/spacer_test.go
internal/controldisplay/result_status_test.go
internal/controldisplay/group_title_test.go
internal/controldisplay/group_counter_test.go
internal/controldisplay/result_reason_test.go
internal/controldisplay/formatter_test.go
internal/controlexecute/dimension_color_map_test.go
internal/resources/query_provider_impl_test.go
internal/resources/query_args_test.go
internal/workspace/lazy_fallback_test.go
internal/workspace/tag_filter_test.go
internal/workspace/resources_of_type_test.go
internal/parse/query_invocation_test.go
```

### Patterns

Standard Go `testing` package only (no testify or other assertion libraries).

**Table-driven tests** (most common):
```go
tests := map[string]struct {
    input    string
    expected string
}{
    "case name": {input: "foo", expected: "bar"},
}
for name, tt := range tests {
    t.Run(name, func(t *testing.T) {
        result := functionUnderTest(tt.input)
        if result != tt.expected {
            t.Fatalf("expected %s, got %s", tt.expected, result)
        }
    })
}
```

**Assertions**: `reflect.DeepEqual()` for complex types, `t.Fatalf()` / `t.Errorf()` for failures.

**Mocks**: Minimal test implementations of interfaces, created inline:
- `makeTaggedControl()` in `tag_filter_test.go` creates controls with specific tags
- `testFormatter` in `formatter_test.go` implements `Formatter` interface

**Cleanup**: `defer os.RemoveAll()` for temp directories.

### Key Test Files

| File | What it tests | Scenarios |
|------|--------------|-----------|
| `resources/query_args_test.go` | `ResolveArgs()` | 30+ scenarios: named/positional args, defaults, runtime overrides, error cases |
| `workspace/tag_filter_test.go` | `ResourceFilterFromTagArgs()` | Tag-based resource filtering with mock controls |
| `controldisplay/formatter_test.go` | `NewFormatResolver()` | Export format resolution (snapshot, csv, json, asff, nunit3) |
| `parse/query_invocation_test.go` | Query parsing | SQL string vs named query resolution |

---

## Acceptance Tests (BATS)

### Running

```bash
tests/acceptance/run-local.sh              # All tests (starts/stops steampipe)
tests/acceptance/run-local.sh check.bats   # Single test file
```

### Framework

**BATS** (Bash Automated Testing System) with three libraries in `tests/acceptance/lib/`:
- `bats-core` - Core framework (test runner)
- `bats-assert` - 19 assertion functions
- `bats-support` - Support utilities

Output format: TAP (Test Anything Protocol).

### Prerequisites

- **Steampipe service** must be running (`run-local.sh` handles this automatically)
- Required CLI tools: `jq`, `jd` (JSON diff), `sed`, `openssl`, `cksum`, and standard Unix tools
- Optional (CI): `SPIPETOOLS_PG_CONN_STRING`, `SPIPETOOLS_TOKEN` environment variables

### Test Configuration

| Setting | Value |
|---------|-------|
| `BATS_TEST_TIMEOUT` | 180 seconds |
| `BATS_NO_PARALLELIZE_WITHIN_FILE` | true (serial execution) |
| `STEAMPIPE_CONNECTION_WATCHER` | false |
| `STEAMPIPE_INTROSPECTION` | info |
| `POWERPIPE_DISPLAY_WIDTH` | 100 (for rendering tests) |

### Test Files (15 files, 233 tests)

| File | Purpose |
|------|---------|
| `check.bats` | Control/benchmark execution, exit codes, output formats |
| `dashboard.bats` | Dashboard execution and output |
| `dashboard_parsing_validation.bats` | HCL parsing validation |
| `database_precedence.bats` | Database config priority resolution |
| `config_precedence.bats` | Configuration priority |
| `config_path.bats` | Config path resolution |
| `backend.bats` | DuckDB and other backends |
| `mod.bats` | Mod install/dependency (**auto-generated, do NOT edit**) |
| `mod_install.bats` | Mod installation specific |
| `params_and_args.bats` | Parameter and argument resolution |
| `var_resolution.bats` | Variable resolution |
| `snapshot.bats` | Snapshot/export output |
| `resource_show_outputs.bats` | Resource display |
| `sp_files.bats` | Legacy .sp file support |
| `tag_filtering.bats` | Tag-based filtering |

### Test Pattern

```bash
load "$LIB_BATS_ASSERT/load.bash"
load "$LIB_BATS_SUPPORT/load.bash"

@test "control run produces correct JSON output" {
  cd $FUNCTIONALITY_TEST_MOD
  run powerpipe control run my_control --output json
  assert_equal $status 0
  assert_output --partial '"status": "ok"'
}

@test "benchmark with alarms exits with code 1" {
  cd $CHECK_ALL_MOD
  run powerpipe benchmark run all_alarms
  assert_equal $status 1
}
```

### Common Assertions

| Function | Purpose |
|----------|---------|
| `assert_equal $status <code>` | Check exit code |
| `assert_success` | Exit code 0 |
| `assert_failure` | Non-zero exit code |
| `assert_output "<text>"` | Full stdout match |
| `assert_output --partial "<text>"` | Partial stdout match |
| `assert_line` | Check specific output line |
| `assert_regex` | Regex match |

### Test Data

**53 test mods** in `tests/acceptance/test_data/mods/`:
- Core: `functionality_test_mod`, `check_all_mod`
- Control rendering: `control_rendering_test_mod`
- Dashboard variants: `dashboard_cards`, `dashboard_graphs`, `dashboard_inputs`, `dashboard_withs`
- Parsing: `dashboard_parsing_*` (3 variants)
- Database: `duckdb_mod`, `mod_with_db*` (7 variants)
- Variables: `test_workspace_mod_var_*` (4 variants)
- Special: `failure_test_mod`, `mod_with_blank_dimension_value`, `mod_with_list_param`

**Expected outputs** in `test_data/templates/`:
- `expected_check_*.{csv,json,html,md,xml}` - Control output formats
- `expected_sps_*.json` - Dashboard snapshots
- `expected_*_title.txt` - Formatted display output

**Helper scripts** in `test_data/scripts/`:
- `update_top_level_mod_version.sh` - Simulates mod version update
- `update_top_level_mod_commit.sh` - Simulates mod retagging

### Test Generation

`mod.bats` is auto-generated. Do NOT edit directly.

```bash
make build-tests   # Regenerate from JSON test cases
```

**Source**: `test_data/source_files/mod_test_cases.json`
**Template**: `test_data/templates/mod_test_template.bats.tmpl`
**Generator**: `tests/acceptance/test_generator/generate.go`

### Environment Variables (set by `run.sh`)

```
BATS_PATH          → lib/bats-core/bin/bats
LIB_BATS_ASSERT    → lib/bats-assert
LIB_BATS_SUPPORT   → lib/bats-support
TEST_DATA_DIR      → test_data/templates
MODS_DIR           → test_data/mods
SNAPSHOTS_DIR      → test_data/snapshots
WORKSPACE_DIR      → test_data/mods/sample_workspace
FUNCTIONALITY_TEST_MOD    → test_data/mods/functionality_test_mod
CHECK_ALL_MOD             → test_data/mods/check_all_mod
CONTROL_RENDERING_TEST_MOD → test_data/mods/control_rendering_test_mod
CONFIG_PARSING_TEST_MOD   → test_data/mods/config_parsing_test_mod
```

---

## Frontend Tests

```bash
cd ui/dashboard
yarn test            # Jest + React Testing Library
yarn storybook       # Component stories on http://localhost:6006
```

Test files follow `*.test.ts` / `*.test.tsx` pattern.
Storybook stories follow `*.stories.@(js|jsx|ts|tsx)` pattern.
