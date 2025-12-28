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
go test ./...           # Run Go tests
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
