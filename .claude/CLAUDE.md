# Powerpipe

Powerpipe is a CLI tool for dashboards, benchmarks, and compliance monitoring. Written in Go with a React dashboard UI.

## Build Commands

```bash
make build              # Build Go binary to /usr/local/bin
make dashboard_assets   # Build React dashboard
make all                # Build both
```

## Testing

```bash
go test ./...           # Run Go unit tests
```

Use `/verify` skill for full local verification workflow including acceptance tests.

## Project Structure

- `internal/cmd/` - CLI commands (check, dashboard, query, etc.)
- `internal/controldisplay/` - Control output formatting
- `internal/controlexecute/` - Control execution logic
- `internal/workspace/` - Workspace loading (lazy and eager modes)
- `internal/resourceindex/` - Fast HCL metadata extraction
- `internal/dashboardserver/` - Dashboard server and WebSocket handling
- `ui/dashboard/` - React dashboard application

## Workspace Loading

Powerpipe uses phased loading to balance fast startup with complete functionality.

### Loading Phases

1. **Phase 1 (Index Build)**: Parse HCL syntax, extract metadata
   - Fast: ~300-500ms for large mods
   - Captures: names, titles, descriptions, tags (literals)
   - Skips: reference resolution, validation

2. **Phase 2 (Background Resolution)**: Resolve dynamic metadata
   - Runs in background goroutines
   - Resolves: variable references, templates, function calls
   - Updates index progressively

3. **Phase 3 (On-Demand Loading)**: Full resource loading
   - Triggered by user interaction (click dashboard, run benchmark)
   - Full HCL parsing with reference resolution
   - Results cached for reuse

### Configuration

| Environment Variable | Description | Default |
|---------------------|-------------|---------|
| `POWERPIPE_WORKSPACE_PRELOAD` | Force eager loading (true/false) | false |

### Key Files

- `internal/workspace/lazy_workspace.go` - Lazy workspace coordinator
- `internal/workspace/background_resolver.go` - Background resolution
- `internal/resourceindex/scanner_hcl.go` - HCL metadata extraction
- `internal/resourceindex/entry.go` - Index entry structure
- `internal/resourceloader/loader.go` - On-demand resource loading

## Skills

| Skill | When to Use |
|-------|-------------|
| `verify` | Local verification before pushing (build, lint, tests) |
| `project-workflow` | Complex multi-task projects, parallel execution |
