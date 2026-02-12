# Environment Variables & Configuration Reference

## Powerpipe Environment Variables

| Variable | Purpose | Default |
|----------|---------|---------|
| `POWERPIPE_LISTEN` | Server listen: `"local"` or `"network"` | `local` |
| `POWERPIPE_PORT` | Server port | `9033` |
| `POWERPIPE_BENCHMARK_TIMEOUT` | Benchmark timeout (seconds) | `300` |
| `POWERPIPE_DASHBOARD_TIMEOUT` | Dashboard query timeout (seconds) | `120` |
| `POWERPIPE_DISPLAY_WIDTH` | CLI output width | Auto-detect |
| `POWERPIPE_INSTALL_DIR` | Installation directory | `~/.powerpipe` |
| `POWERPIPE_MOD_LOCATION` | Mod directory override | Current dir |
| `POWERPIPE_DATABASE` | Default database connection string | - |
| `POWERPIPE_MAX_PARALLEL` | Max parallel DB connections | `10` |
| `POWERPIPE_QUERY_TIMEOUT` | Per-query timeout (seconds) | - |
| `POWERPIPE_SNAPSHOT_LOCATION` | Snapshot upload target | - |
| `POWERPIPE_UPDATE_CHECK` | Enable version check on startup | `true` |
| `POWERPIPE_TELEMETRY` | Telemetry level | `info` |
| `POWERPIPE_CACHE_ENABLED` | Enable query caching | `true` |
| `POWERPIPE_CACHE_TTL` | Cache TTL (seconds) | - |
| `POWERPIPE_CACHE_MAX_TTL` | Cache max TTL (seconds) | - |
| `POWERPIPE_MEMORY_MAX_MB` | Go memory limit (MB) | `1024` |
| `POWERPIPE_MEMORY_MAX_MB_PLUGIN` | Plugin memory limit (MB) | - |
| `POWERPIPE_LOG_LEVEL` | Log level (see below) | `off` |
| `POWERPIPE_CONFIG_PATH` | Colon-separated config search paths | - |
| `POWERPIPE_CONFIG_DUMP` | Debug: dump config as JSON (undocumented) | - |

### Pipes Cloud Variables

| Variable | Purpose |
|----------|---------|
| `PIPES_HOST` | Pipes cloud host |
| `PIPES_TOKEN` | Pipes authentication token |
| `PIPES_INSTALL_DIR` | Pipes installation directory |

### Input Variables

Set HCL variables via `PP_VAR_` prefix:
```bash
PP_VAR_region=us-east-1
PP_VAR_max_age=30
```

### Deprecated Variables

| Variable | Replacement |
|----------|-------------|
| `STEAMPIPE_DIAGNOSTICS_LEVEL` | `PIPES_DIAGNOSTICS_LEVEL` |

## Configuration Precedence (highest to lowest)

1. CLI flags (`--flag`)
2. Environment variables (`POWERPIPE_*`)
3. Explicit workspace profile (`--workspace-profile`)
4. Default workspace profile
5. Config files (`.ppc`)
6. Viper defaults (set in `internal/cmdconfig/mappings.go`)

## Log Levels

Set via `POWERPIPE_LOG_LEVEL` (actual env var name from `app_specific.EnvLogLevel`).

| Value | slog Level | Notes |
|-------|------------|-------|
| `trace` | Custom trace | Most verbose |
| `debug` | `slog.LevelDebug` | |
| `info` | `slog.LevelInfo` | |
| `warn` | `slog.LevelWarn` | |
| `error` | `slog.LevelError` | |
| `off` | Discard | **Default** - no logging |

Logger implementation: Go `log/slog` with JSON handler to stderr. Sensitive values automatically redacted via `sanitize.Instance.SanitizeKeyValue()`.

Performance tracing: `utils.LogTime("label")` logs timestamped markers.

## File Extensions

| Extension | Purpose |
|-----------|---------|
| `.pp` | Mod and resource files (primary) |
| `.sp` | Legacy Steampipe format (still supported) |
| `.ppvars` | Variable files |
| `.spvars` | Legacy variable files |
| `.auto.ppvars` | Auto-loaded variable files |
| `.ppc` | Config files |
| `.powerpipeignore` | Workspace ignore patterns |
| `.mod.cache.json` | Dependency lock file |

## Database Connection String Formats

| Format | Backend |
|--------|---------|
| `steampipe://profile/schema` | Steampipe (PostgreSQL) |
| `postgres://user:pass@host/db` | PostgreSQL |
| `mysql://user:pass@host/db` | MySQL |
| `duckdb:///path/to/db.duckdb` | DuckDB (file) |
| `duckdb://` | DuckDB (in-memory) |
| `sqlite:///path/to/db.sqlite` | SQLite |

## Exit Codes

| Code | Constant | Meaning |
|------|----------|---------|
| 0 | `ExitCodeSuccessful` | Success |
| 1 | `ExitCodeControlsAlarm` | Check/benchmark: alarms found, no errors |
| 2 | `ExitCodeControlsError` | Check/benchmark: control errors found |
| 21 | `ExitCodeSnapshotCreationFailed` | Snapshot creation failed |
| 22 | `ExitCodeSnapshotUploadFailed` | Snapshot upload failed |
| 41 | `ExitCodeQueryExecutionFailed` | Query execution failed |
| 62 | `ExitCodeModInstallFailed` | Mod installation failed |
| 250 | `ExitCodeInitializationFailed` | Workspace/DB initialization failed |
| 251 | `ExitCodeBindPortUnavailable` | Server port binding failed |
| 252 | `ExitCodeNoModFile` | No mod.pp found |
| 254 | `ExitCodeInsufficientOrWrongInputs` | Invalid user input |
| 255 | `ExitCodeUnknownErrorPanic` | Panic recovery (unhandled crash) |

## Error Handling Patterns

- **Panic recovery**: `main.go` has a deferred `recover()` that catches panics, logs error via `error_helpers.ShowError()`, and exits with code 255
- **InitData errors**: Commands call `NewInitData[T]()`, then check `initData.Result.HasError()`. Errors from workspace loading, dependency install, and DB client creation are collected in `Result.ErrorAndWarnings`
- **Error display**: `error_helpers.ShowError(ctx, err)` and `error_helpers.FailOnError()` from pipe-fittings
- **Goroutine errors**: In `dashboardexecute`, child errors propagate via channels to parent nodes. No panic/recover in execution goroutines

## Build Variables (injected via ldflags)

| Variable | Dev Value | Release Value |
|----------|-----------|---------------|
| `main.version` | `0.0.0-dev-{branch}.{timestamp}` | Semver from git tag |
| `main.commit` | `none` | Git commit hash |
| `main.date` | `unknown` | Build timestamp |
| `main.builtBy` | `local` | `goreleaser` |

Dev mode detection: `cmdconfig.IsLocal()` returns `true` when `builtBy == "local"`.
