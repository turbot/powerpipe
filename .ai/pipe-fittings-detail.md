# Pipe-Fittings Detailed Reference

Pipe-fittings (`github.com/turbot/pipe-fittings/v2`) is the shared foundation library for Powerpipe, Flowpipe, Steampipe, and Tailpipe.

## Package Map

```
pipe-fittings/
├── modconfig/         HCL resource model (Mod, Variable, HclResource, ModTreeItem)
├── parse/             HCL parsing engine (Decoder, DecoderImpl, ModParseContext)
├── load_mod/          High-level mod loading orchestration
├── workspace/         Workspace abstraction (load mods, file watcher, variables)
├── backend/           Database backend interface (Postgres, MySQL, DuckDB, SQLite)
├── connection/        40+ connection types (AWS, Azure, GCP, Slack, GitHub, ...)
├── credential/        30+ credential types with lazy resolution
├── modinstaller/      Transactional mod dependency installation (git-based)
├── schema/            HCL block type constants and schema definitions
├── queryresult/       Query result streaming (Result, ColumnDef, RowResult)
├── export/            Export framework (JSON, CSV, snapshot)
├── printers/          CLI output formatting (table, JSON, YAML)
├── cmdconfig/         Cobra/Viper CLI flag binding (CmdBuilder)
├── steampipeconfig/   Snapshot types (SteampipeSnapshot, SnapshotPanel)
├── statushooks/       Progress reporting callbacks
├── constants/         Application-wide constants and env vars
├── utils/             General utilities (file, hash, JSON, collections)
├── error_helpers/     Error wrapping and formatting
├── hclhelpers/        HCL/CTY type conversions
├── cty_helpers/       CTY value helpers
├── filter/            PEG-based filter → SQL WHERE clause
├── filepaths/         Cross-platform path utilities
├── pipes/             Turbot Pipes cloud integration
├── plugin/            Steampipe plugin management
├── versionmap/        Semantic version resolution and lock files
├── ociinstaller/      OCI image downloading for plugins
├── secrets/           Secret detection and authentication
├── inputvars/         Variable validation and collection
├── options/           Execution options (query, dashboard, check)
├── perr/              Structured error types
├── cache/             TTL-based connection caching
└── sanitize/          Input sanitization
```

## Key Interfaces

```
modconfig.ModTreeItem          Base interface for all executable resources
modconfig.HclResource          Base interface for all HCL-defined resources
modconfig.Mod                  Root module container
modconfig.Variable             Input variable with type validation
modconfig.ParamDef             Parameter definition for queries
modconfig.ModResources         Interface for resource collections

backend.Backend                Database backend (Connect, RowReader, Name)
backend.SearchPathProvider     Search path for Steampipe/PostgreSQL
backend.SearchPathConfig       Search path + prefix configuration
backend.BackendOption          Configuration option for backends

connection.ConnectionStringProvider   Get connection string for a resource
connection.SearchPathProvider         Search path from connection config

parse.Decoder                  HCL block decoder interface
parse.DecoderImpl              Base decoder with pluggable decode funcs
parse.ModParseContext          Parsing state for a single mod

queryresult.SyncQueryResult    Complete query result (cols + rows)
queryresult.ColumnDef          Column metadata (name, type)
queryresult.ResultStreamer      Streaming result handler

steampipeconfig.SteampipeSnapshot   Dashboard execution snapshot
steampipeconfig.SnapshotPanel       Individual panel data
steampipeconfig.SnapshotTreeNode    Execution tree structure

workspace.Workspace            Base workspace (mod loading, file watcher)

export.Manager                 Export format registry
export.Exporter                Export implementation (JSON, CSV, snapshot)

statushooks.*                  Progress reporting callbacks

printers.Showable              Objects supporting show command
printers.Listable              Objects supporting list command
```

## How Powerpipe Extends Pipe-Fittings

```
pipe-fittings                          powerpipe
─────────────                          ─────────
workspace.Workspace          ────▶  PowerpipeWorkspace
  (mod loading, file watcher)          + dashboard event handlers
                                       + dashboard event channel
                                       + PublishDashboardEvent()
                                       + RegisterDashboardEventHandler()
                                       + ResolveQueryFromQueryProvider()

modconfig.ModResources       ────▶  PowerpipeModResources
  (interface for resource             + Dashboards, DashboardCards, Charts,
   collections)                         Tables, Texts, Graphs, Flows, ...
                                       + Controls, Benchmarks
                                       + Detections, DetectionBenchmarks
                                       + Queries
                                       + WalkResources(), AddResource()

parse.DecoderImpl            ────▶  PowerpipeModDecoder
  (base HCL decoder)                  + custom decode for Dashboards
                                       + custom decode for QueryProviders
                                       + custom decode for NodeAndEdgeProviders
                                       + custom decode for Benchmarks
                                       + dynamic schema per resource type
```

## Usage by Powerpipe Component

| Powerpipe Component | Top Pipe-Fittings Packages Used | Purpose |
|---------------------|-------------------------------|---------|
| `internal/resources/` (42 files) | modconfig, utils, cty_helpers, printers, schema | Resource type definitions |
| `internal/cmd/` (11 files) | constants, cmdconfig, error_helpers, modconfig | CLI command handlers |
| `internal/dashboardexecute/` (15 files) | steampipeconfig, modconfig, backend, connection | Dashboard execution engine |
| `internal/controlexecute/` (8 files) | steampipeconfig, modconfig, queryresult, schema | Control/benchmark execution |
| `internal/db_client/` (10 files) | backend, connection, constants, queryresult | Database client layer |
| `internal/cmdconfig/` (8 files) | cmdconfig, error_helpers, connection, parse | CLI config and flag handling |
| `internal/controldisplay/` (20 files) | constants, export, modconfig, filepaths | CLI output rendering |
| `internal/workspace/` (8 files) | workspace, modconfig, error_helpers, utils | Workspace wrapper |
| `internal/parse/` (5 files) | schema, modconfig, hclhelpers, parse | HCL parsing extensions |
| `internal/dashboardserver/` (6 files) | backend, connection, modconfig, steampipeconfig | WebSocket server |
| `internal/initialisation/` (2 files) | modconfig, backend, error_helpers, statushooks | Init orchestration |

## Import Frequency (top 10)

| Sub-Package | Imports | Importance |
|-------------|---------|------------|
| `modconfig` | 74 | CRITICAL - resource data model |
| `utils` | 62 | HIGH - utility functions |
| `constants` | 58 | HIGH - app constants |
| `error_helpers` | 34 | HIGH - error handling |
| `steampipeconfig` | 24 | HIGH - snapshot types |
| `schema` | 24 | HIGH - HCL schema |
| `cty_helpers` | 21 | MEDIUM - HCL type conversion |
| `printers` | 18 | MEDIUM - CLI output |
| `connection` | 17 | MEDIUM - DB connections |
| `statushooks` | 15 | MEDIUM - progress reporting |

## Design Patterns

- **Factory**: `backend.FromConnectionString()`, `modinstaller.NewModInstaller()`
- **Strategy**: Multiple backend implementations, exporters
- **Plugin Architecture**: Pluggable decoders, exporters, secret matchers
- **Visitor**: `WalkResources()` for mod traversal
- **Builder**: `CmdBuilder` for CLI setup, `ModParseContext` for parsing

## Internal Dependency Graph

```
modconfig (resource definitions)
    ↑
parse (HCL parsing)
    ↑
load_mod (orchestration)
    ↑
workspace (runtime container)
    ↑
powerpipe/flowpipe (applications)

connection/credential (data source config)
    ↑
backend (database abstraction)
    ↑
queryresult (execution results)

modinstaller (dependency management)
    ↑
versionmap/ociinstaller (version/artifact handling)
```
