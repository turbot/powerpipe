# Workspace Pod Startup: Powerpipe Server & Mod Installation Analysis

**Date:** 2026-02-18

## Overview

When a workspace pod spins up, Powerpipe-related initialization happens across two
distinct layers: the container infrastructure layer (Powerpipe binary + server) and
the Temporal workflow layer (mod installation, variable generation, snapshots). These
are **intentionally async** -- the pod becomes "ready" (state=enabled) before mods
are fully installed, allowing users to query immediately while mods install in the
background.

---

## End-to-End Lifecycle

```
  User: "Create Workspace"
           |
           v
  +---------------------------+
  | API Server                |
  | RunWorkspaceCreateWorkflow|
  +---------------------------+
           |
           v
  +--------------------------------------------------+
  |     WorkspaceCreateWorkflow (API Worker Queue)    |
  |--------------------------------------------------|
  | 1. WorkspaceCreateK8sObjectsActivity              |
  |    - Renders K8s YAML from template               |
  |    - Creates StatefulSet, Services, Certs, etc    |
  |    - Sets DashboardProvider = "powerpipe"         |
  |    - Configures PowerpipeMaxParallelConnections   |
  |    - Passes PowerpipeQueryTimeout (900s)          |
  |                                                    |
  | 2. CreateDnsRecordsetActivity                      |
  |    - Dashboard DNS (legacy + pipes domain)         |
  |    - Dashboard TLS certs                           |
  |                                                    |
  | 3. SaveWorkspacePublicKeyActivity                  |
  |    - Extracts TLS cert from K8s secret             |
  |                                                    |
  | 4. Schedule cron workflows                         |
  |    - WorkspaceK8sUpdateCron                        |
  |    - WorkspaceStatusWorkflow                       |
  |    - WorkspaceMaintenanceWorkflow                  |
  |    - WorkspaceBackupWorkflows                      |
  |                                                    |
  | 5. TestUserConnectionToWorkspaceActivity           |
  |    - Retries up to 180x until pod reachable        |
  +--------------------------------------------------+
           |
           | K8s pod starts, container boots
           v
```

---

## Container Boot Sequence

The workspace pod uses the `docker/workspace/Dockerfile` image:

```
  +---------------------------------------------+
  |  Workspace Container (Ubuntu 24.04)          |
  |---------------------------------------------|
  |  Installed at build time:                    |
  |    - Steampipe CLI (steampipe service)       |
  |    - Powerpipe CLI (powerpipe binary)        |
  |    - PostgreSQL (via steampipe)              |
  |    - pgbackrest, git, jq, etc               |
  |                                              |
  |  Exposed ports:                              |
  |    9193 - Steampipe PostgreSQL               |
  |    9194 - Powerpipe Dashboard Server         |
  |                                              |
  |  Entrypoint: /opt/steampipe/steampipe-server |
  |    (the Pipes Go binary, runs as worker)     |
  +---------------------------------------------+
```

The Go binary starts as a pod worker (`--type pod`), which triggers the in-pod
Temporal worker and the PodInitWorkflow:

```
  steampipe-server worker --type pod
       |
       |  cmd/worker.go:runPodWorkerService()
       |
       +---> launcher.RunPodInitWorkflow(workspaceID, identityID, spDataMode)
       |
       +---> Register Temporal worker on queue: "workspace-{workspaceID}"
       |     Registered workflows:
       |       - PodInitWorkflow
       |       - PodUpdateWorkflow
       |       - WorkspaceModUpdateWorkflowWithSignalsV2
       |       - WorkspaceConnectionUpdateWorkflowWithSignalsV2
       |       - UpdateWorkspaceDatabaseOptionsWorkflowWithSignals
       |       - UpdateWorkspaceAggregatorWorkflowWithSignals
       |       - DatatankWorkflow
       |       - BackupWorkflows
       |       - ... (many more)
       |
       +---> worker.Run()  // Start polling Temporal
```

---

## PodInitWorkflow: The Orchestrator

**File:** `workflow/pod/index.go:48`
**Runs on:** Pod task queue (`workspace-{workspaceID}`)

This is the critical orchestration point. It launches PARALLEL async workflows
via signals and only BLOCKS on PodUpdateWorkflow:

```
  PodInitWorkflow
       |
       +---> Create Process (workspace.init)
       |
       |     BLOCKING (child workflow):
       +====> PodUpdateWorkflow -----> [waits for completion]
       |                                        |
       |     FIRE-AND-FORGET (signals):         |
       +---> SignalUpdateWorkspaceAggregatorWF   |
       +---> SignalWorkspaceModUpdateWF_V2       |
       +---> SignalWorkspaceConnectionUpdateWF   |
       |                                        |
       |     <---- PodUpdateWorkflow returns ----|
       |
       +---> Update Process -> COMPLETED
```

**Key insight:** The comment in code says it all:
> "Don't check the output for mod update workflow because it's a signal workflow
> and we don't need to 'wait' for its result -- just let it do its thing."

---

## PodUpdateWorkflow: Making the Pod Usable

**File:** `workflow/pod/index.go:180`
**Runs on:** Pod task queue
**This is the BLOCKING path** - workspace won't be "enabled" until this completes.

```
  PodUpdateWorkflow
       |
       +---> [PARALLEL child workflows on API queue]
       |     +---> GetPgBouncerAuthUserForWorkspaceWorkflow
       |     +---> GetWorkspaceUsersWorkflow
       |
       +---> UpdatePgBouncerUsersActivity
       |       - Decrypts user credentials
       |       - Renders SQL from add-roles.sql template
       |       - Executes against localhost:9193
       |
       +---> CheckSteampipeDbUpActivity  <-- RETRY x180
       |       - Pings localhost:9193
       |       - Backoff coefficient = 1 (constant interval)
       |       - Waits until Steampipe PostgreSQL is responsive
       |
       +---> GetWorkspaceDatabaseOptionsWorkflow
       |       - Gets search_path config from API DB
       |
       +---> UpdatePodDbUsersActivityV2
       |       - Creates pgexporter role
       |       - Creates workspace users in steampipe_users group
       |       - Sets passwords and search_path per user
       |       - Cleans up orphaned users
       |
       +---> WorkspaceGetMetadataActivity
       |       - Reads /tmp/pp-version.txt  (Powerpipe version)
       |       - Reads /tmp/sp-version.txt  (Steampipe version)
       |       - Reads /tmp/spc-version.txt (API version)
       |
       +---> WorkspaceMetadataUpdateWorkflow
              - Sets State = "enabled"
              - Saves version metadata to DB
              - WORKSPACE IS NOW USABLE
```

---

## Async: Mod Installation Flow

**File:** `workflow/mod/mod_update_signals_v2.go`
**Pattern:** Signal-with-start (coalescing)
**Workflow ID:** `WorkspaceModUpdateWorkflowWithSignalsV2-{workspaceID}`
**Signal name:** `UpdateModWorkspaceV2-{workspaceID}`

This is the fully async mod installation pipeline:

```
  Signal received (from PodInit or API mod add/update)
       |
       |  Coalescing: if multiple signals pending,
       |  only the LATEST one is processed
       |
       +---> GetModsForWorkspaceWorkflow (API queue)
       |       - Fetches workspace mod list from database
       |       - Returns []WorkspaceMod with paths, constraints, integrations
       |
       +---> GenerateIntegrationTokenWorkflow (API queue)
       |       - For private GitHub/GitLab mods
       |       - Creates installation access tokens
       |       - Encrypts tokens for pod consumption
       |
       +---> UpdateWorkspaceModActivity (ON POD)
       |       |
       |       |  Phase 1: Uninstall removed mods
       |       +---> Read .mod.cache.json
       |       +---> For each mod in cache but NOT in workspace:
       |       |       powerpipe mod uninstall {mod-name} \
       |       |         --mod-location /opt/steampipe/workspace
       |       |
       |       |  Phase 2: Install/update mods
       |       +---> For each workspace mod:
       |       |     +---> [if archive type]:
       |       |     |       Download from GCS via signed URL
       |       |     |       Extract to .local/{modID}/
       |       |     |       powerpipe mod install {local-path} \
       |       |     |         --pull latest \
       |       |     |         --mod-location /opt/steampipe/workspace
       |       |     |
       |       |     +---> [if repository type]:
       |       |     |       Set POWERPIPE_GIT_TOKEN if integration
       |       |     |       powerpipe mod install {path}@{constraint} \
       |       |     |         --pull latest \
       |       |     |         --mod-location /opt/steampipe/workspace
       |       |     |
       |       |     +---> On success: state = "installed"
       |       |     +---> On failure: state = "error" (continues)
       |       |
       |       +---> Read updated .mod.cache.json
       |       +---> Return installed mod versions/commits
       |
       +---> GetWorkspaceModVariableListActivity (ON POD)
       |       powerpipe variable list --output json \
       |         --mod-location /opt/steampipe/workspace
       |
       +---> UpdateModInstalledDataWorkflow (API queue)
       |       - Updates DB with installed versions, commits
       |       - Syncs mod variable metadata
       |
       +---> GenerateSPVarsWorkflow (ON POD)
       |       - Writes powerpipe.ppvars file with variable values
       |       - Format: {mod_alias}.{var_name} = {value}
       |
       +---> Update all processes -> COMPLETED/FAILED
```

**Mod installation is resilient to failures:**
- Individual mod install failures don't abort the whole workflow
- Failed mods get `state = "error"` with the stderr as `state_reason`
- Successfully installed mods still get recorded
- A cron workflow (`GenerateWorkspaceModUpdateCronWorkflow`) periodically
  re-attempts mod updates

---

## Async: Connection Configuration Flow

**File:** `workflow/workspace/connection_update_signals_v2.go`
**Pattern:** Signal-with-start (coalescing)
**Workflow ID:** `WorkspaceConnectionUpdateWorkflowWithSignalsV2-{workspaceID}`

```
  Signal received (from PodInit or connection add/update)
       |
       +---> GetWorkspaceConnectionsWorkflowV2 (API queue)
       |       - Fetches encrypted connection config from DB
       |
       +---> UpdatePodConnectionsActivity (ON POD)
       |       - Decrypts connection payload
       |       - Writes /opt/steampipe/config/connections.spc
       |       - If temp credentials: schedule connection update cron
       |       - If no temp creds: cancel connection update cron
       |
       +---> SignalUpdateWorkspaceDatabaseOptionsWorkflow
              - Updates search_path and other DB options
              - Writes database options config on pod
```

---

## Async: Aggregator Configuration Flow

**File:** `workflow/aggregator/index.go`
**Pattern:** Signal-with-start (coalescing)
**Workflow ID:** `UpdateWorkspaceAggregatorWorkflowWithSignals-{workspaceID}`

Configures Steampipe aggregator connections (plugin-level data sources).
Similar signal-coalescing pattern.

---

## Timeline View

```
  T=0s     WorkspaceCreateWorkflow starts
  T=1s     |  K8s StatefulSet created
  T=2s     |  DNS records configured
  T=3s     |  Public key saved
  T=5s     |  Cron workflows scheduled
  T=5s     |  Waiting for pod connectivity (up to 180 retries)
           |
  T=~10s   |  Pod container starts
           |  +--> steampipe-server worker --type pod
           |  +--> Steampipe service starts (PostgreSQL on 9193)
           |  +--> Powerpipe server starts (Dashboard on 9194)
           |  +--> Version files written to /tmp/
           |
  T=~15s   |  PodInitWorkflow triggers
           |  |
           |  |  +-- PodUpdateWorkflow (BLOCKING) -----------+
           |  |  |   Waiting for Steampipe DB...              |
           |  |  |   Configure PgBouncer users                |
           |  |  |   Create DB users, set search_path         |
           |  |  |   Read version metadata                    |
           |  |  +-- Complete: State -> ENABLED -------> T=~45s
           |  |
           |  |  +-- ModUpdateWorkflow (ASYNC) -----------+
           |  |  |   Fetch mod list from API DB            |
           |  |  |   Generate integration tokens           |
           |  |  |   powerpipe mod install (each mod)      |
           |  |  |   powerpipe variable list               |
           |  |  |   Write powerpipe.ppvars                |
           |  |  +-- Complete ----------------------> T=~60-300s
           |  |
           |  |  +-- ConnectionUpdateWorkflow (ASYNC) ---+
           |  |  |   Fetch connections from API DB        |
           |  |  |   Write connections.spc                |
           |  |  |   Update database options              |
           |  |  +-- Complete ----------------------> T=~30s
           |  |
           |  |  +-- AggregatorUpdateWorkflow (ASYNC) ---+
           |  |  |   Fetch aggregator config              |
           |  |  |   Write aggregator .spc files          |
           |  |  +-- Complete ----------------------> T=~30s
           |
  T=~45s   TestUserConnectionToWorkspaceActivity succeeds
           WorkspaceCreateWorkflow COMPLETE
           |
  T=~60s+  Mods still installing in background...
           |
  T=~300s  All async workflows complete
           WORKSPACE FULLY OPERATIONAL
```

---

## Key Architectural Points

### 1. Powerpipe Is Used as CLI, Not a Long-Running Server Process

Despite the port 9194 exposure for the dashboard server, the core Powerpipe
interactions from the Go codebase are all CLI invocations:

| Command | Purpose | Where |
|---------|---------|-------|
| `powerpipe mod install` | Install/update mods | UpdateWorkspaceModActivity |
| `powerpipe mod uninstall` | Remove mods | UpdateWorkspaceModActivity |
| `powerpipe variable list` | List mod variables | GetWorkspaceModVariableListActivity |
| `powerpipe dashboard run` | Create snapshots | uploadSnapshotV2 |
| `powerpipe benchmark run` | Run benchmarks | uploadSnapshotV2 |
| `powerpipe query run` | Run queries | uploadSnapshotV2 |

The Powerpipe dashboard HTTP server (port 9194) is started by the container
infrastructure (K8s template/sidecar) and handles the live dashboard UI. The Go
workflow code doesn't manage this process.

### 2. Signal-With-Start Pattern for Coalescing

All three async workflows (mods, connections, aggregators) use Temporal's
`SignalWithStartWorkflow` pattern. This means:

- If the workflow is already running, the signal is delivered to it
- If the workflow is not running, it starts a new instance AND delivers the signal
- Multiple rapid signals are coalesced: `if !selector.HasPending()` ensures only
  the latest signal is fully processed

### 3. Two Task Queue Architecture

```
  +----------------------------------+     +---------------------------+
  |  API Worker Queue                |     |  Pod Worker Queue         |
  |  ("workspace-crud")              |     |  ("workspace-{wsID}")     |
  |                                  |     |                           |
  |  - WorkspaceCreateWorkflow       |     |  - PodInitWorkflow        |
  |  - GetModsForWorkspaceWorkflow   |     |  - PodUpdateWorkflow      |
  |  - GetWorkspaceConnectionsWF     |     |  - ModUpdateWF (signals)  |
  |  - UpdateModInstalledDataWF      |     |  - ConnectionUpdateWF     |
  |  - GenerateIntegrationTokenWF    |     |  - AggregatorUpdateWF     |
  |  - WorkspaceMetadataUpdateWF     |     |  - UpdateWorkspaceModAct  |
  +----------------------------------+     |  - UpdatePodConnectionAct |
       ^                                   |  - BackupWorkflows        |
       |  Child workflows to               |  - DatatankWorkflows      |
       |  access Pipes database             +---------------------------+
       +---------- called from pod -------->
```

Child workflows that need database access run on the API queue. Activities
that need pod filesystem access (mod install, connection writing) run on the
pod queue.

### 4. Workspace Directory Structure on Pod

```
  /opt/steampipe/
  ├── steampipe-server          # Pipes Go binary
  ├── config/
  │   └── connections.spc       # Written by ConnectionUpdate workflow
  ├── workspace/                # SP_PP_WORKSPACE_PATH
  │   ├── mod.sp                # Powerpipe mod definition
  │   ├── .mod.cache.json       # Mod cache (versions, commits)
  │   ├── powerpipe.ppvars      # Variable values
  │   └── .local/               # Extracted archive mods
  │       └── {modID}/
  ├── templates/
  │   └── add-roles.sql         # PgBouncer role template
  └── nfs/                      # Shared NFS for plugins

  /tmp/
  ├── pp-version.txt            # Powerpipe version
  ├── sp-version.txt            # Steampipe CLI version
  └── spc-version.txt           # API version
```

---

## Key Source Files

| File | Purpose |
|------|---------|
| `cmd/worker.go:734` | `runPodWorkerService` - pod worker entry point |
| `workflow/pod/index.go:48` | `PodInitWorkflow` - orchestrates pod init |
| `workflow/pod/index.go:180` | `PodUpdateWorkflow` - makes pod usable |
| `workflow/mod/mod_update_signals_v2.go:21` | `WorkspaceModUpdateWorkflowWithSignalsV2` |
| `workflow/mod/index.go:98` | `UpdateWorkspaceModActivity` - runs powerpipe CLI |
| `workflow/workspace/create.go:40` | `WorkspaceCreateWorkflow` - full create flow |
| `workflow/workspace/connection_update_signals_v2.go:19` | Connection signal workflow |
| `workflow/launcher/mod_update_signals_v2.go:15` | Signal launcher for mod updates |
| `workflow/launcher/connection_update_signals_v2.go:14` | Signal launcher for connections |
| `workflow/launcher/aggregator.go:18` | Signal launcher for aggregators |
| `types/mod.go:94` | `GeneratePowerpipeInstallCommand` |
| `util/dashboard.go:3` | `GetDashboardProvider` returns "powerpipe" |
| `tppath/index.go:11` | `PowerpipeWorkspacePath` = /opt/steampipe/workspace |
| `docker/workspace/Dockerfile` | Workspace container image definition |
