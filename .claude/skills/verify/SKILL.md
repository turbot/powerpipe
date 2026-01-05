---
name: verify
description: Verify code changes locally before pushing. Runs build, lint, unit tests, and relevant acceptance tests.
---

# Verify Skill

Run local verification of code changes before pushing to CI. This is faster than waiting for GitHub Actions and catches most issues immediately.

## Quick Verification

For most changes, run these steps:

```bash
# 1. Build the binary
make build

# 2. Run linting
golangci-lint run

# 3. Run Go unit tests
go test -short -timeout 120s ./...
```

## Acceptance Tests

### First-Time Setup

Initialize BATS submodules (only needed once):

```bash
git submodule update --init --recursive
```

### Running Tests

Run specific test files from the acceptance directory:

```bash
cd tests/acceptance
./run-local.sh tag_filtering.bats     # Tag filtering tests
./run-local.sh check.bats             # Check/benchmark tests
./run-local.sh mod.bats               # Mod management tests
./run-local.sh dashboard.bats         # Dashboard tests
./run-local.sh                        # ALL tests (slow, ~10+ min)
```

### Available Test Files

| Test File | What It Tests |
|-----------|---------------|
| `check.bats` | Benchmark/control execution, exports |
| `tag_filtering.bats` | Tag-based control filtering |
| `mod.bats` | Mod install, update, dependencies |
| `dashboard.bats` | Dashboard commands |
| `backend.bats` | Database backend connections |
| `snapshot.bats` | Snapshot creation/loading |
| `var_resolution.bats` | Variable resolution |
| `config_precedence.bats` | Config file precedence |

## Direct CLI Testing

For quick manual verification of specific functionality:

```bash
# Test against acceptance test mods
cd tests/acceptance/test_data/mods/<mod_name>
powerpipe benchmark run <benchmark_name> --progress=false --output=json

# Examples:
cd tests/acceptance/test_data/mods/tag_filtering_mod
powerpipe benchmark run tag_filtering_benchmark --tag deprecated=true --progress=false --output=json

cd tests/acceptance/test_data/mods/check_all_mod
powerpipe benchmark run all --progress=false
```

## Verification Workflow

1. Make code changes
2. `make build` - rebuild binary with changes
3. `go test ./...` - run unit tests
4. `./tests/acceptance/run-local.sh <relevant>.bats` - run related acceptance tests
5. Commit and push
6. Verify CI passes on GitHub

## When to Run What

| Change Type | Verification Steps |
|-------------|-------------------|
| Any Go code | `make build` + `go test ./...` |
| Control/benchmark logic | + `check.bats`, `tag_filtering.bats` |
| Mod loading/workspace | + `mod.bats`, `check.bats` |
| Dashboard features | + `dashboard.bats` |
| Config/variable handling | + `config_precedence.bats`, `var_resolution.bats` |
| Major refactoring | Run all tests locally before pushing |

## Troubleshooting

### bats: command not found

```bash
git submodule update --init --recursive
```

### Tests pass locally but fail in CI

- Check if steampipe service is running: `steampipe service status`
- The `run-local.sh` script starts/stops steampipe automatically
- CI uses Linux; local Mac differences are rare but possible

### Slow test runs

- Run only relevant test files, not all tests
- Use `--progress=false` for CLI commands to reduce output
