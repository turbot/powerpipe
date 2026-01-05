# Powerpipe

Powerpipe is a CLI tool for dashboards, benchmarks, and compliance monitoring. Written in Go with a React dashboard UI.

## Build Commands

```bash
make build              # Build Go binary to /usr/local/bin
make dashboard_assets   # Build React dashboard
make all                # Build both
```

## Testing

### Quick Verification (Local)

Always verify changes locally before pushing:

```bash
# 1. Build the binary
make build

# 2. Run linting
golangci-lint run

# 3. Run Go unit tests
go test -short -timeout 120s ./...
```

### Acceptance Tests (Local)

BATS acceptance tests require submodule initialization (once):

```bash
git submodule update --init --recursive
```

Run individual test files:

```bash
cd tests/acceptance
./run-local.sh tag_filtering.bats     # Run specific test file
./run-local.sh check.bats             # Run check tests
./run-local.sh                        # Run all tests (slow)
```

### Test Workflow

1. Make changes
2. `make build` - rebuild binary
3. `go test ./...` - run unit tests
4. `./tests/acceptance/run-local.sh <test>.bats` - run relevant acceptance tests
5. Push and verify CI passes

### Direct CLI Testing

For quick manual verification:

```bash
cd tests/acceptance/test_data/mods/<mod_name>
powerpipe benchmark run <benchmark_name> --progress=false --output=json
```

## Project Structure

- `internal/cmd/` - CLI commands (check, dashboard, query, etc.)
- `internal/controldisplay/` - Control output formatting
- `internal/controlexecute/` - Control execution logic
- `ui/dashboard/` - React dashboard application

## Skills

| Skill | When to Use |
|-------|-------------|
| `project-workflow` | Complex multi-task projects, parallel execution |
