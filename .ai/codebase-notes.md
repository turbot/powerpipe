# Powerpipe Codebase Guide

## What the project does
- CLI and dashboard server for Powerpipe (“dashboards for DevOps”). Runs dashboards, controls/benchmarks/detections, queries, and serves a browser UI that renders dashboards and streams live execution via WebSocket.
- Targets local mods (HCL-defined resources) and Turbot Pipes cloud workspaces; can install dependencies, validate plugin requirements, and publish/share snapshots.
- Built atop the shared `pipe-fittings` library (app framing, mod parsing, backend/database abstraction, config, status hooks, exporters) with additional Powerpipe-specific resource types and runtime.

## Architecture & tech stack
- Language: Go 1.24, Cobra/Viper for CLI, Gin + Melody for HTTP/WebSocket, SQL drivers (pgx, mysql, duckdb, sqlite), Steampipe plugin SDK for telemetry/error helpers, semver/hcl parsing.
- Frontend: React/TypeScript (CRA + craco, Tailwind, Storybook) in `ui/dashboard`, built assets embedded in Go via `internal/dashboardassets`.
- Core packages live under `internal`:
  - `cmd`: Cobra commands (server, mod install/update/list, login, resource list/show, query run, benchmark/control/detection run, dashboard serve).
  - `initialisation`, `controlinit`: workspace loading, dependency install, DB client wiring, telemetry, exporter registration.
  - `workspace`: Powerpipe workspace wrapper around `pipe-fittings/workspace`, file watching, resource resolution, query/provider discovery, runtime dependency verification.
  - `resources`: definitions for Powerpipe resource types (dashboards, cards/charts/tables/edges/flows/nodes/inputs/etc., controls, detections, benchmarks, queries, variables). Provides add/get/walk/equals logic and resource-specific behaviors (query provider interface, runtime dependencies, argument merge).
  - `dashboardexecute`: orchestrates dashboard/control/detection execution trees, snapshots, input dependency handling, runtime dependency pub/sub, and event emission.
  - `dashboardserver` + `service/api`: Gin HTTP server + Melody WebSocket that serve the dashboard UI assets and stream dashboard events to the browser; handles listen/bind options and local/network modes.
  - `db_client`: database/back-end abstraction (connection parsing, search_path overrides, max connections, client map, cloud/steampipe/postgres/etc.).
  - `controlexecute`, `controlstatus`, `controldisplay`, `display`, `queryresult`: run controls/benchmarks/detections, format/stream status, render CLI/output/export formats, convert snapshots to tabular results.
  - `dashboardassets`: ensures embedded dashboard asset tarball matches current app version (or verifies local assets during dev).
  - `dashboardtypes`, `dashboardevents`, `dashboardassets`, `dashboardtypes`: types for execution states/events/payloads.
  - `powerpipeconfig`: global config + pipeling connection map and cloud connection caching.
  - `parse`: HCL schema/decoder for Powerpipe resource blocks, argument decoding, validation helpers.
  - `cmdconfig`, `constants`, `logger`, `types`: CLI flag helpers, env vars, defaults, build metadata, logging glue with pipe-fittings.
- Other: `dist/` built binaries, `images/` art, `tests/acceptance` harness, `scripts/` build/test helpers, `ui/dashboard` frontend source.

## External dependencies & inter-repo coupling
- Heavy reliance on [`github.com/turbot/pipe-fittings`](https://github.com/turbot/pipe-fittings) (replace to `../pipe-fittings` available). Key packages used:
  - `cmdconfig`: wraps Cobra/Viper flag setup (`OnCmd(...).Add*Flag`), cloud/profile flags, and shared flag names.
  - `constants`: CLI arg names, config keys, output formats, exit codes.
  - `app_specific` / `app_specific_connection`: app defaults (names, versions, install directories, mod filenames) and default pipeling connections used for DB resolution.
  - `modconfig`, `workspace`, `parse`, `schema`: HCL resource parsing, mod/resource modeling, resource lookup/resolution, variable/lifecycle validation.
  - `modinstaller`: installing/updating/pruning mod dependencies; install summaries.
  - `backend`, `connection`, `search_path`, `steampipeconfig`: backend selection from connection strings, search_path handling, plugin version map for requirement validation, telemetry bootstrapping.
  - `filepaths`, `utils`: asset and config directories, timers/log helpers, filename/path helpers.
  - `statushooks`: status spinner/hooks injected into contexts (CLI, asset extraction).
  - `error_helpers`, `sperr`, `contexthelpers`: panic/error wrapping, user-friendly failure output, cancel handlers.
  - `export`, `querydisplay`, `queryresult`, `printers`: exporter registry and CLI output renderers leveraged by control/query flows.
  - `plugin`: plugin version maps used when checking mod requirements against a steampipe backend.
- Steampipe plugin SDK (`steampipe-plugin-sdk/v5`) for telemetry/logging, error wrapping (`sperr`), and plugin version mapping when validating requirements.
- Gin + Melody provide HTTP/WebSocket for the dashboard server; gzip/size/middleware for rate limiting.
- UI build artifacts are bundled via Go `embed` in `internal/dashboardassets` and served by the server. Frontend depends on pipeline events emitted by `dashboardexecute` via `dashboardserver`.
- Turbot Pipes cloud integration via connection strings/pipeling connections in `powerpipeconfig`; snapshot publishing/export options in CLI rely on shared exporters from pipe-fittings.

## Key workflows
- **CLI entry**: `main.go` sets version metadata, app-specific constants, captures panics, then `cmd.Execute()` (Cobra). Root command wires global flags for config/install dir/mod location/workspace profile.
- **Workspace initialization** (`initialisation.NewInitData`): load workspace (prompt for missing vars) using `workspace.Load`, install mod dependencies (`modinstaller`) if required, construct default DB client from workspace/default connection, validate runtime requirements (plugins, dependencies), build a `dashboardexecute.Executor` with a client map. Handles telemetry lifecycle and captures errors/warnings via `InitResult`.
- **Control/benchmark/detection run** (`checkCmd`): validate args, initialize via `controlinit.NewInitData` (adds formatters/exporters, filters by tags/where, ensures templates), build execution trees (`getExecutionTrees` in `controlexecute`), run controls with status spinner and aggregate alarms/errors; exit codes reflect alarms/errors vs runtime failures. Supports snapshot/export/share options.
- **Query run**: resolves target (raw SQL or named query), reuses dashboard snapshot generation for execution, displays snapshot or converts to table output, supports snapshot publish/export and search_path overrides.
- **Server start** (`serverCmd`): validate port/listen, init workspace, ensure dashboard assets, create WebSocket, instantiate `dashboardserver.Server` (registers workspace event handler), start Gin HTTP API (`service/api`) serving static assets and `/ws` socket, watch files for HCL changes and push `DashboardChanged` events.
- **Dashboard execution**: `dashboardexecute.Executor` manages per-session execution trees, validates required inputs (batch vs interactive), emits events (`ExecutionStarted/ControlComplete/ControlError/LeafNodeUpdated/InputValuesCleared/ExecutionComplete`) consumed by `dashboardserver` and forwarded to UI via WebSocket. Supports loading snapshots from workspace files and runtime dependency propagation between panels.
- **Workspace watch & runtime dependencies**: `workspace.PowerpipeWorkspace` extends pipe-fittings Workspace with dashboard event channels; file watcher publishes `DashboardChanged` events for UI hot-reload. Runtime dependency providers/subscribers managed in `dashboardexecute` and resource implementations.
- **Assets**: `dashboardassets.Ensure` checks embedded asset version; in dev (`cmdconfig.IsLocal`) requires local built assets (`make dashboard_assets`).

## Top-level folder summaries
- `internal/cmd`: Cobra command definitions. Each command wires flags via `cmdconfig` (from pipe-fittings), validates args, initializes workspace (`initialisation`/`controlinit`), and delegates to executor/display packages. Commands include `server`, `mod` (install/update/list/init), `login`, resource list/show, `query` run, `check` for controls/benchmarks/detections, and `dashboard` run.
- `internal/initialisation` & `controlinit`: Shared initialization across commands; handles telemetry, dependency install, DB client creation, requirement validation, target resolution, exporter registration, and cleanup. `controlinit` adds formatter/executor wiring for control-like resources.
- `internal/workspace`: Wrapper around pipe-fittings workspace with Powerpipe-specific resource handling. Loads mod files/HCL resources, tracks whether modfile exists, resolves query providers (inline SQL vs named queries), validates runtime dependencies for dashboards, exposes resource lookup/filtering, and manages dashboard event channels/watchers.
- `internal/resources`: HCL resource models and helpers. Defines `PowerpipeModResources` maps, `AddResource`/`GetResource`/`WalkResources`, equality checks for change detection, runtime dependency provider/subscriber interfaces, query provider implementation, dashboard graph/table/input/card/etc. definitions, benchmarking types, and arg merging/with-provider logic.
- `internal/dashboardexecute`: Execution engine for dashboards/controls/detections. Maintains per-session execution trees, client map selection (default vs overridden DB/search_path), input value tracking and dependency invalidation, runtime dependency pub/sub, snapshot generation (`GenerateSnapshot`), and conversion to display trees. Emits `dashboardevents` for server/UI.
- `internal/dashboardserver`: Handles dashboard events from workspace/executor and forwards payloads over WebSocket; tracks connected clients; starts API server to serve assets and websocket endpoint; outputs readiness/status messages. Includes `ListenPort/ListenType` helpers, payload builders, and message routing per session.
- `internal/service/api`: Gin setup, middleware (gzip, rate limiting, request size limits), static asset serving, websocket hookup, and minimal public endpoints (`/service`). Uses validators for versioning. Accepts listen configuration from CLI.
- `internal/db_client`: DB client wrapper over `pipe-fittings/backend` and `database/sql`; manages connection strings, search_path overrides, max connections, cloud connection helpers, and client maps for multi-database execution.
- `internal/controlexecute`, `controlstatus`, `controldisplay`, `display`, `queryresult`: Runs control trees, formats status/spinner, renders CLI output (text/csv/json/html/md/nunit/snapshot), exports, and converts snapshots to table results. Also handles dimension color maps and result grouping.
- `internal/dashboardtypes`, `dashboardevents`, `dashboardassets`, `dashboardtypes`, `dashboardevents`: Shared types and event payloads for execution state and UI updates.
- `internal/powerpipeconfig`: Global config holder, default pipeling connections, and cache for cloud workspace connection strings; ensures default connections can be overridden.
- `internal/parse`: HCL schema/decoder for Powerpipe resources/args, resource validation, and query invocation parsing.
- `ui/dashboard`: React/Tailwind dashboard UI (CRA + craco). Entry `src/index.tsx` wraps the app in Theme/Breakpoint/Analytics/ErrorBoundary providers; routes are in `App.tsx` (root + snapshot routes). `registerComponents.ts` registers all dashboard components from `@powerpipe/components` (charts, flows, graphs, inputs, layout, snapshot header); custom dashboard pieces live in `src/components/dashboards/*` and shared widgets in `src/components/*`. The core runtime is context-driven:
  - `useDashboardState.tsx`: central reducer/context; handles available dashboards, execution started/complete/error events, panel logs/maps, selection, data modes (live/cli_snapshot/cloud_snapshot), schema migrations, version mismatch checks, and analytics hooks.
  - `useDashboardExecution.tsx`: wires the live event loop; sets up websocket via `useDashboardWebSocket.ts` (react-use-websocket) to `/ws`, dispatches CLEAR/SELECT/INPUT_CHANGED actions (control inputs, datetime range, search_path_prefix), navigates to snapshots, and loads snapshot metadata into state.
  - `useDashboardWebSocketEventHandler.ts` (not shown above) decodes server events (`available_dashboards`, `execution_started`, `execution_complete`, `controls_updated`, `leaf_nodes_*`, etc.) and dispatches reducer actions; schema migrations for panel statuses live in `utils/dashboardEventHandlers.ts`.
  - Other contexts cover inputs, datetime range, search path, search, panel detail, theme, breakpoints, and analytics; combined in `useDashboard.tsx` provider.
  - Layout renderer `components/dashboards/layout/Dashboard/index.tsx` builds grids, titles, progress, side panel, and split-pane resizing; uses `Children`/`Grid` and side-panel components to show selected panels; `DashboardList` lists available dashboards; `DashboardHeader` contains search/toggles; `WorkspaceErrorModal` handles workspace parse/build errors from server events.
  - Data modes: live websocket streaming; CLI snapshot mode (renders a saved snapshot and prevents re-execution); cloud snapshot mode.
  - Routing: `/` live dashboards, `/:dashboard_name` for direct selection, `/snapshot/:dashboard_name` for snapshot rendering with inputs/time range/search-path encoded in query params.
  - Styles: Tailwind classes + `styles/index.css`; icons via `src/icons`; theme switcher via `hooks/useTheme`.
  - Tests/storybook: component stories; `spec.md` documents the mod/dashboard HCL spec; build assets via `make dashboard_assets` (embedded into Go binary by `dashboardassets`).
- `tests/acceptance`: Acceptance test harness/specs for CLI behaviors.
- `dist`, `scripts`, `images`, `Makefile`: Build artifacts, packaging scripts, release assets, helper targets (including building dashboard assets).

## Constraints, assumptions, and edge cases
- Modfile requirement: many commands need a mod definition; query run bypasses if executing raw SQL. Workspace loading prompts for missing variables; fails fast if runtime dependencies (e.g., referenced queries) are missing.
- Ports/listen: server validates bindability; supports `local` (localhost) vs `network` listen. Exits with bind error codes if occupied.
- Input handling: batch mode requires all required dashboard inputs upfront; interactive mode allows runtime input changes with dependent inputs cleared and re-run.
- Runtime dependencies: dashboards/components can depend on others; executor manages pub/sub to avoid stale data. Validation errors surface during init.
- Exports/snapshots: exporters must be registered/validated; snapshot upload handles share/visibility tags/location. Exit codes differ for execution failure vs upload failure vs alarms.
- Database/search_path: CLI flags override defaults; max connections enforced; non-steampipe backends skip plugin validation.
- Dev mode: when `IsLocal()` is true, embedded assets are not used; local dashboard assets must exist or `powerpipe server` fails.

## How this repo links to `pipe-fittings`
- Uses `pipe-fittings` for: app constants/flags/config keys, workspace/model parsing, backend selection and connection strings, mod installer/updater, exporter registry, status hooks/spinners, error helpers, serialization, telemetry bootstrapping, and utility helpers (`utils.LogTime`, `filepaths`, etc.).
- Optional local development uses `replace github.com/turbot/pipe-fittings/v2 => ../pipe-fittings` in `go.mod`.
- Powerpipe-specific layers mostly exist in `internal/resources`, `workspace`, `dashboardexecute`, `controlexecute`, and `dashboardserver`, composed atop pipe-fittings primitives.

## Patterns and conventions
- Cobra/Viper for CLI: flags defined via `cmdconfig.OnCmd`, defaults pulled from pipe-fittings constants/environment, output controlled via `--output`.
- Generic target handling: commands parameterized over resource types (`checkCmd[T]`, `resourceCmd[T]`) to reuse resolution and display logic.
- Event-driven dashboard runtime: file watcher + execution events feed websocket payloads to UI; `dashboardevents` types drive payload builder pattern.
- Snapshot-first execution: query/control/detection executions often render to snapshots and then adapt to display/export formats.
- Extensive use of helper packages for display/export (color schemes, templates) and backend abstraction to keep commands thin.

## Entry points
- CLI: `main.go` → `internal/cmd/root.go`.
- Server: `internal/cmd/server.go` → `dashboardserver.Server` + `service/api`.
- Execution engines: `dashboardexecute.Executor` (dashboards), `controlexecute` (controls/detections).
- Workspace loading: `internal/initialisation` + `internal/workspace`.
- Frontend: `ui/dashboard/src/index.tsx` (CRA entry), assets embedded via `internal/dashboardassets`.

## CHEAT SHEET FOR FUTURE PROMPTS
- **TL;DR architecture**: Go CLI that loads a Powerpipe mod (HCL) via pipe-fittings, initializes DB/backend clients, runs dashboards/controls/detections/queries to produce snapshots, and serves a React UI over Gin+WebSocket for interactive dashboards. Execution events flow from `dashboardexecute` → `dashboardevents` → `dashboardserver` → browser.
- **Glossary**:  
  - Workspace: mod directory + variables + dependencies; loaded by `internal/workspace`.  
  - Mod resources: dashboards/cards/charts/tables/flows/nodes/inputs, controls/detections/benchmarks, queries/variables. Stored in `resources.PowerpipeModResources`.  
  - Executor: `dashboardexecute.Executor` manages runs and input dependencies.  
  - Snapshot: serialized execution result; can be exported/shared.  
  - Backend/DB: resolved via `db_client`/`pipe-fittings/backend` (steampipe/postgres/mysql/sqlite/duckdb).  
  - Dashboard server: Gin HTTP + Melody WS serving assets and events.
- **Core workflows**:  
  - Run control/benchmark/detection: use `checkCmd` flow (init workspace, build execution trees, run via `controlexecute`, export/share snapshot).  
  - Run query: `query run <name|sql>` creates snapshot then renders/export.  
  - Serve dashboards: `server` command initializes workspace, ensures assets, starts API/WebSocket, watches files for live reload.  
  - Manage mods: `mod install/update/list/init` via `modinstaller` (pipe-fittings).
- **Key files**: `main.go`, `internal/cmd/root.go`, `internal/cmd/server.go`, `internal/cmd/check.go`, `internal/cmd/query.go`, `internal/initialisation/*`, `internal/workspace/*`, `internal/resources/*`, `internal/dashboardexecute/*`, `internal/dashboardserver/*`, `internal/service/api/*`, `internal/db_client/*`, `internal/dashboardassets/ensure.go`, `ui/dashboard/*`.
- **Important constraints**: modfile usually required; dev mode needs local dashboard assets; batch dashboards need all inputs; server port/listen validated; runtime dependencies must resolve; exporters require format validation.
- **How to ask Codex for changes**: specify the command or package (e.g., “update dashboardexecute to add new event type”, “add CLI flag in internal/cmd/query.go”), mention any new resource schema changes needed in `internal/parse`/`resources`, and if frontend changes are required, point to `ui/dashboard` components and ensure dashboard assets build target is updated if packaging.

Refer to https://www.powerpipe.io/docs and https://www.powerpipe.io/docs/powerpipe-hcl for CLI and HCL semantics; the repo mirrors those behaviors.
