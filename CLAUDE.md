# Powerpipe

Powerpipe is a CLI and dashboard server for DevOps. It runs dashboards, controls/benchmarks/detections, and queries, and serves a browser UI that renders dashboards and streams live execution via WebSocket. It targets local mods (HCL-defined resources) and Turbot Pipes cloud workspaces.

**Detailed reference docs** (read on demand, not loaded by default):
- `.ai/dashboard-ui.md` - React/TypeScript frontend architecture, hooks, state, WebSocket protocol
- `.ai/testing-guide.md` - Unit tests, BATS acceptance tests, test data, test generation
- `.ai/environment-reference.md` - All env vars, exit codes, config precedence, file extensions, log levels
- `.ai/pipe-fittings-detail.md` - Full pipe-fittings package map, interfaces, import frequency
- `.ai/codebase-notes.md` - Detailed codebase walkthrough and key workflows

## Architecture Overview

```
┌──────────────────────────────────────────────────────────────────────────┐
│  User: powerpipe server / powerpipe benchmark run aws_compliance.cis    │
└──────────────┬───────────────────────────────────────────────────────────┘
               │
       ┌───────▼──────────┐
       │  Powerpipe CLI    │  ← This repo (turbot/powerpipe)
       │  (Cobra + Go)     │
       └───────┬───────────┘
               │ Loads HCL mod
       ┌───────▼──────────────┐
       │  Workspace / Mod      │  ← pipe-fittings + internal/workspace
       │  (dashboards, controls│
       │   benchmarks, queries)│
       └───────┬──────────────┘
               │ SQL queries
       ┌───────▼──────────────┐
       │  Database Backend     │  (Steampipe/PostgreSQL/MySQL/SQLite/DuckDB)
       └──────────────────────┘
               │
       ┌───────▼──────────────┐
       │  Dashboard UI         │  React 18 / TypeScript (Gin + WebSocket)
       └──────────────────────┘
```

### Initialization and Command Flow

Every CLI command follows the same initialization pattern:

```
┌──────────────┐     ┌────────────────────────────────────────────────┐
│ Cobra Command │────▶│ NewInitData[T]()                               │
│ (cmd/*.go)    │     │  T = Dashboard | Control | Query | Detection   │
└──────────────┘     └──────────────┬─────────────────────────────────┘
                                    │
                     ┌──────────────▼──────────────────┐
                     │ 1. LoadWorkspace()               │
                     │    pipe-fittings Workspace        │
                     │    + PowerpipeModDecoder (parse)   │
                     │    → PowerpipeModResources         │
                     ├─────────────────────────────────────┤
                     │ 2. ResolveTargets[T]()             │
                     │    CLI args → actual resources      │
                     ├─────────────────────────────────────┤
                     │ 3. Init()                           │
                     │    ├─ Install mod dependencies      │
                     │    ├─ GetDefaultDatabaseConfig()    │
                     │    │   resolve: CLI > mod > default │
                     │    ├─ NewDbClient()                 │
                     │    │   → backend.FromConnectionString│
                     │    ├─ ValidateModRequirements()     │
                     │    └─ Create DashboardExecutor      │
                     └─────────────────────────────────────┘
```

### Dashboard Execution Flow

```
ExecuteDashboard(sessionId, rootResource, inputs)
  │
  ▼
┌─────────────────────────────────────────────────────────────┐
│ DashboardExecutionTree                                       │
│  ├─ Root: DashboardRun (container)                           │
│  │    ├─ LeafRun (card)  ──── SQL ───▶ DbClient ──▶ Backend │
│  │    ├─ LeafRun (chart) ──── SQL ───▶ DbClient ──▶ Backend │
│  │    ├─ LeafRun (table) ──── SQL ───▶ DbClient ──▶ Backend │
│  │    └─ LeafRun (graph)                                     │
│  │         ├─ with run (data provider) ── SQL ──▶ Backend    │
│  │         ├─ node run ── SQL(with args) ──────▶ Backend     │
│  │         └─ edge run ── SQL(with args) ──────▶ Backend     │
│  │                                                           │
│  ├─ Events published:                                        │
│  │    ExecutionStarted ──▶ WebSocket ──▶ UI                  │
│  │    LeafNodeUpdated  ──▶ WebSocket ──▶ UI (per panel)      │
│  │    ExecutionComplete ──▶ WebSocket ──▶ UI (snapshot)      │
│  │                                                           │
│  └─ ClientMap: connection pool keyed by connString+searchPath│
└─────────────────────────────────────────────────────────────┘
```

### Runtime Dependency Resolution

Dashboard panels can depend on `with` blocks and `input` values via a publish-subscribe pattern:

```
DashboardRun (publisher)
  │
  ├── publishes input values ────────────────────────┐
  │                                                   │
  ├── LeafRun: "with.vpc_ids" (publisher)             │
  │    └─ executes SQL, publishes result rows         │
  │         │                                         │
  │         ▼                                         ▼
  ├── LeafRun: "chart.instances" (subscriber)         │
  │    ├─ subscribes to "with.vpc_ids" channel        │
  │    ├─ subscribes to "input.region" channel ◄──────┘
  │    ├─ waits for all dependency values
  │    ├─ resolveSQLAndArgs() substitutes into SQL
  │    └─ executeQuery() runs final SQL
```

### Dashboard Server Event Flow

```
┌──────────┐  WebSocket   ┌──────────────┐  Event Bus  ┌──────────────┐
│ React UI  │◄────────────▶│ Melody Server │◄───────────▶│  Executor    │
│ (browser) │              │ (Gin + WS)    │             │              │
└──────────┘              └──────────────┘             └──────────────┘

Client → Server:
  {action: "request_dashboard", payload: {dashboard: "aws.dashboard.vpc"}}
  {action: "update_input", payload: {inputs: {region: "us-east-1"}}}

Server → Client:
  {action: "execution_started",  panels: {...}, layout: {...}}
  {action: "leaf_node_updated",  dashboard_node: {...}}
  {action: "execution_complete", snapshot: {...}}
  {action: "workspace_error",    error: {...}}
```

## Pipe-Fittings Integration

Pipe-fittings (`github.com/turbot/pipe-fittings/v2`) is the shared foundation library. Powerpipe imports 30+ sub-packages across 185+ Go files. The dependency is unidirectional.

> Full package map, interface list, and import frequency: `.ai/pipe-fittings-detail.md`

Powerpipe wraps and extends three key abstractions:

```
pipe-fittings                          powerpipe
─────────────                          ─────────
workspace.Workspace          ────▶  PowerpipeWorkspace
  (mod loading, file watcher)          + dashboard event handlers/channel
                                       + PublishDashboardEvent()

modconfig.ModResources       ────▶  PowerpipeModResources
  (interface for resource             + Dashboards, Cards, Charts, Tables, ...
   collections)                       + Controls, Benchmarks, Detections, Queries

parse.DecoderImpl            ────▶  PowerpipeModDecoder
  (base HCL decoder)                  + custom decode per resource type
                                       + dynamic schema generation
```

## Resource Type Hierarchy

```
modconfig.HclResource (pipe-fittings)
  │
  ├── modconfig.Mod, Variable, Local
  │
  └── powerpipe resources (internal/resources/)
       │
       ├── QueryProvider (interface: has SQL/Query/Args/Params)
       │    ├── Query, Control, Detection
       │    ├── DashboardCard, Chart, Table, Image, Text
       │    └── NodeAndEdgeProvider (extends QueryProvider)
       │         ├── DashboardGraph, DashboardFlow, DashboardHierarchy
       │
       ├── WithProvider (has 'with' blocks)
       ├── DashboardLeafNode (renderable panel)
       ├── Benchmark (tree of Controls)
       ├── DetectionBenchmark (tree of Detections)
       ├── Dashboard, DashboardContainer, DashboardInput
       └── DashboardCategory, DashboardEdge, DashboardNode
```

## Database Layer

```
                    backend.Backend (pipe-fittings)
                           │
          ┌────────────────┼────────────────┐
          │                │                │
  ┌───────▼──────┐ ┌──────▼───────┐ ┌──────▼──────┐
  │ Steampipe/   │ │   DuckDB     │ │  MySQL      │
  │ PostgreSQL   │ │              │ │             │
  └──────────────┘ └──────────────┘ └─────────────┘
                   ┌──────▼──────┐
                   │  SQLite     │
                   └─────────────┘

DbClient (powerpipe wrapper around sql.DB + Backend)
  ├── ExecuteSync(sql, args) → SyncQueryResult (blocks until complete)
  └── Execute(sql, args) → Result with RowChan (streams rows)

ClientMap (connection pool, keyed by connString + searchPathConfig)
  ├── Default clients: shared, long-lived
  └── Session clients: per-execution, closed after

Database resolution priority (highest → lowest):
  1. Resource-level database/search_path
  2. Mod-level database/search_path
  3. Dependency mod database
  4. CLI --database flag
  5. Default connection
```

## Repository Map

```
powerpipe/
├── main.go                          # Entry point
├── internal/
│   ├── cmd/                         # Cobra commands (server, mod, query, check, login, etc.)
│   ├── initialisation/              # Workspace loading, dependency install, DB client wiring
│   ├── controlinit/                 # Formatter/executor wiring for control-like resources
│   ├── workspace/                   # Powerpipe workspace wrapper around pipe-fittings
│   ├── resources/                   # Powerpipe resource types (dashboards, controls, etc.)
│   ├── dashboardexecute/            # Dashboard/control/detection execution engine
│   ├── dashboardserver/             # Gin HTTP + Melody WebSocket server
│   ├── service/api/                 # Gin setup, middleware, static asset serving
│   ├── db_client/                   # Database/backend abstraction
│   ├── controlexecute/              # Control/benchmark/detection execution
│   ├── controlstatus/               # Control status formatting
│   ├── controldisplay/              # CLI output rendering
│   ├── dashboardassets/             # Embedded dashboard UI assets
│   ├── dashboardtypes/              # Execution state types
│   ├── dashboardevents/             # Event payloads for UI updates
│   ├── powerpipeconfig/             # Global config, pipeling connections
│   ├── parse/                       # HCL schema/decoder for Powerpipe resources
│   ├── cmdconfig/                   # CLI flag helpers
│   ├── constants/                   # App constants, env vars
│   └── ...
├── ui/dashboard/                    # React/TypeScript dashboard UI
├── tests/acceptance/                # BATS acceptance test suite
├── scripts/                         # Build/test helper scripts
└── .ai/                             # Detailed reference docs
```

## Development Guide

### Building

```bash
make build             # Dev build with version ldflags → /usr/local/bin/powerpipe
go build -o powerpipe  # Simple build without version injection
```

Dashboard assets (required for `powerpipe server` in dev):
```bash
cd ui/dashboard && yarn install && yarn build && make zip
```

### Testing

```bash
go test ./...                            # Unit tests (17 files, Go testing package)
tests/acceptance/run-local.sh            # Acceptance tests (15 BATS files, 233 tests)
tests/acceptance/run-local.sh check.bats # Single test file
```

> Full testing guide with patterns, test data, and BATS details: `.ai/testing-guide.md`

### Local Development with Related Repos

```
pipe-fittings          (shared library)
       ↑
powerpipe              (depends on pipe-fittings)
```

Powerpipe's `go.mod` has a **commented-out replace directive**:

```go
//replace github.com/turbot/pipe-fittings/v2 => ../pipe-fittings
```

Uncomment for local development against `../pipe-fittings`. **Re-comment before committing.**

#### Cross-Repo Change Workflow

1. Make the change in `pipe-fittings` first
2. Uncomment the replace directive in `powerpipe`
3. Build and test locally
4. Publish dependency (merge + tag)
5. `go get github.com/turbot/pipe-fittings/v2@v2.x.x`
6. Re-comment the replace directive
7. Commit and PR

### Key Directories for Common Tasks

| Task | Where to Look |
|------|--------------|
| Fix a CLI command | `internal/cmd/` + relevant `internal/` package |
| Fix dashboard execution | `internal/dashboardexecute/` |
| Fix control/benchmark run | `internal/controlexecute/`, `internal/controlinit/` |
| Fix query execution | `internal/cmd/query.go`, `internal/initialisation/` |
| Fix dashboard server | `internal/dashboardserver/`, `internal/service/api/` |
| Fix workspace loading | `internal/workspace/`, `internal/initialisation/` |
| Fix resource definitions | `internal/resources/` |
| Fix HCL parsing | `internal/parse/` |
| Fix DB client | `internal/db_client/` |
| Fix dashboard UI | `ui/dashboard/` (see `.ai/dashboard-ui.md`) |
| Fix display/export | `internal/controldisplay/`, `internal/display/` |
| Fix config | `internal/powerpipeconfig/`, `internal/cmdconfig/` |
| Add new resource type | `internal/resources/`, `internal/parse/mod_decoder.go`, `internal/parse/schema.go` |
| Change mod operations | pipe-fittings `modinstaller/` |
| Change event system | `internal/workspace/workspace_events.go`, `internal/dashboardevents/` |
| Change database config | `internal/db_client/database_config.go` |

### Branching and Workflow

- **Base branch**: `develop` for all work
- **Main branch**: `main` (releases merge here)
- **Release branch**: `v1.4.x` (or similar version branch)
- **PR titles**: End with `closes #XXXX` for bug fixes
- **Merge-to-develop PRs**: Title must be `Merge branch '<branchname>' into develop`
- **Small PRs**: One logical change per PR

## Gotchas

- **Initialization order**: In `dashboardexecute/dashboard_run.go`, runs must be added to execution tree map **before** creating child runs (children look up parent by name)
- **Dashboard assets in dev**: `dashboardassets.Ensure()` checks at server startup. Missing `ui/dashboard/build/` fails the server. Fix: `cd ui/dashboard && yarn build`
- **`mod.bats` is auto-generated**: Do not edit `tests/acceptance/test_files/mod.bats` directly. Run `make build-tests` to regenerate
- **Steampipe service**: Required for acceptance tests. `run-local.sh` handles this, but manual test runs need `steampipe service start` first
- **Replace directive**: `go.mod` replace for pipe-fittings must remain commented in commits
- **Batch mode variables**: With `--input=false`, missing required variables cause immediate failure with no recovery
- **Runtime dependency deadlocks**: Circular `with` dependencies cause execution to hang. Validation should catch this, but be careful
- **Concurrent Viper access**: Background tasks receive config values as params, not via direct Viper reads, to avoid map access panics
- **Memory limit**: `POWERPIPE_MEMORY_MAX_MB` (default 1024) sets Go's `debug.SetMemoryLimit()` - affects GC, not a hard limit
- **Search path**: Non-existent schemas in search path cause runtime query failures, not parse-time errors

## Release Process

Follow these steps in order to perform a release:

### 1. Changelog
- Draft a changelog entry in `CHANGELOG.md` matching the style of existing entries.
- Use today's date and the next patch version.

### 2. Commit
- Commit message for release changelog changes should be the version number, e.g. `v1.4.3`.

### 3. Release Issue
- Use the `.github/ISSUE_TEMPLATE/release_issue.md` template.
- Title: `Powerpipe v<version>`, label: `release`.

### 4. PRs
1. **Against `develop`**: Title should be `Merge branch '<branchname>' into develop`.
2. **Against `main`**: Title should be `Release Powerpipe v<version>`.
   - Body format:
     ```
     ## Release Issue
     [Powerpipe v<version>](link-to-release-issue)

     ## Checklist
     - [ ] Confirmed that version has been correctly upgraded.
     ```
   - Tag the release issue to the PR (add `release` label).

### 5. powerpipe.io Changelog
- Create a changelog PR in the `turbot/powerpipe.io` repo.
- Branch off `main`, branch name: `pp-<version without dots>` (e.g. `pp-143`).
- Add a file at `content/changelog/<YYYYMMDD>-powerpipe-cli-v<version-with-dashes>.md`.
- Frontmatter format:
  ```
  ---
  title: Powerpipe CLI v<version> - <short summary>
  publishedAt: "<YYYY-MM-DD>T10:00:00"
  permalink: powerpipe-cli-v<version-with-dashes>
  tags: cli
  ---
  ```
- Body should match the changelog content from `CHANGELOG.md`.
- PR title: `Powerpipe CLI v<version>`, base: `main`.

### 6. Deploy powerpipe.io
- After the powerpipe.io changelog PR is merged, trigger the `Deploy powerpipe.io` workflow in `turbot/powerpipe.io` from `main`.

### 7. Close Release Issue
- Check off all items in the release issue checklist as steps are completed.
- Close the release issue once all steps are done.
