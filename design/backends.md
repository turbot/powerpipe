# Backends

## Overview

## Resolving Database and search path

When running Powerpipe, the configured `database` and search path config (`search_path` and `search_path_prefix`) are resolved from the following sources (in order of increasing precedence):
- Environment
- workspace profile
- command line

Additionally, the search path and database can be set in various other places

1) Dashboard execution scoped search path config
When executing a dashboard, the dashboard UI can provide search path or search path prefix. This overrides any setting in viper for the scope of this execution (but may be overridden by higher precedence overrides below)
2) dependency mod database settings
`database` and `search_path*` can be set for a dependency mod in the mod `require` block. Any resource which is executed from this mod (even if a child of a resource from another mod) will respect these settings (as long as the resource does not override the connection string) 
3) resource level settings
   `database` and `search_path*` can be set for individual resources (implemented in ModTreeItem). (currently just connection string but should change for consistency?) 

### Resolution approach
When building the execution tree:
- the execution tree is populated with the 'active' database and search path config, resolved from viper config and the search path passed int he dashboard execution request
- each leaf node is populated with the resolved database and search path config, taking into account:
  - dependency mod database config
  - resource database config

### Executing
When executing a leaf node we need to acquire a client specific to the database and search path for the leaf node.

The execution tree has a map of clients, keyed by a database and search path config.  




    
The `DashboardExecutionTree` has a map of clients, keyed by execution. When starting an execution, 
this map is populated with a single client using the default connection string. 

The default connection string is resolved from the following sources (in order of increasing precedence):
    - Environment (`POWERPIPE_DATABASE`)
    - WorkspaceProfile (`Database`)
    - Command Line (`--database`)
    - Mod `require` block (for dependency mods)

## Search Path
The PostgresBackend (and so also the SteampipeBackend) supports setting the search path when acquiring a connection

This is done using a `requiredSearchPath` property. Is this is non-null, after establishing a connection the backend sets the search path to the desired value`

Whenever a PostgresBackend is created, `requiredSearchPath` is resolved using the configured search_path and search_path_prefix.
search_path and search_path_prefix can be set in the following ways (in order of increasing precedence):
- WorkspaceProfile
- Command Line
- Mod `require` block (for dependency mods)
- Dashboard server event parameter passed with `execute_ddasboard` event


## Dashboard Execution
Dashboard executed by calling `Executor.ExecuteDashboard`

This is called from:
    
### CLI commands
`GenerateSnapshot` is called from `dashboard run` and `query run` commands

### Dashboard Server
Executor.ExecuteDashboard is called from the event handlers for `DashboardChanged` and `select_dashboard`

If called from dashboard changed, the existing client map is reused