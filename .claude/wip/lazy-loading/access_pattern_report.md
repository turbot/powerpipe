# Resource Access Pattern Analysis Report

## Executive Summary

This report documents all code paths that access resources from the workspace in Powerpipe, categorized by access pattern type and timing. This analysis informs the design of lazy loading to ensure all access patterns are properly handled.

## Resource Types

The workspace contains the following resource types in `PowerpipeModResources`:

| Resource Type | Map Name | Count in Large Mods |
|---------------|----------|---------------------|
| Benchmark | `ControlBenchmarks` | 100s |
| Control | `Controls` | 1000s |
| Dashboard | `Dashboards` | 10s-100s |
| Query | `Queries` | 100s |
| Detection | `Detections` | 100s |
| DetectionBenchmark | `DetectionBenchmarks` | 10s |
| DashboardContainer | `DashboardContainers` | 100s |
| DashboardCard | `DashboardCards` | 100s |
| DashboardChart | `DashboardCharts` | 100s |
| DashboardTable | `DashboardTables` | 100s |
| DashboardInput | `DashboardInputs`, `GlobalDashboardInputs` | 10s |
| Variable | `Variables` | 10s |
| + 10 other dashboard component types | | |

---

## 1. Access Points Summary

### 1.1 GetModResources() Calls
**Location**: Multiple files
**Pattern**: Returns entire `PowerpipeModResources` struct

| File | Line | Context |
|------|------|---------|
| `internal/initialisation/init_data.go` | 271 | Iterating child mods |
| `internal/dashboardserver/server.go` | 220, 371 | Building server metadata payload |
| `internal/workspace/powerpipe_workspace.go` | 53, 127 | Verify runtime deps, get resources |
| `internal/resources/mod_resources.go` | 12, 15 | Helper function |
| `internal/resources/dashboard_category_helpers.go` | 13 | Category resolution |

### 1.2 GetResource() Calls (Single Resource Lookup)
**Pattern**: Direct lookup by parsed resource name

| File | Line | Context |
|------|------|---------|
| `internal/dashboardserver/server.go` | 516 | Get resource for execution |
| `internal/workspace/powerpipe_workspace.go` | 114 | Find resource by name |
| `internal/workspace/resource_from_args.go` | 72 | Resolve resource from SQL string |
| `internal/dashboardevents/dashboard_changed.go` | 340 | Check if resource exists |

### 1.3 WalkResources() Calls (Full Iteration)
**Pattern**: Iterate over ALL resources

| File | Line | Context |
|------|------|---------|
| `internal/resources/mod_resources.go` | 106 | QueryProviders() |
| `internal/resources/mod_resources.go` | 126 | TopLevelResources() |
| `internal/resources/mod_resources.go` | 634 | WalkResources implementation |
| `internal/dashboardexecute/referenced_variables.go` | 41, 55 | Walk dashboard for variables |
| `internal/resources/dashboard.go` | 208, 255, 328 | Walk dashboard children |
| `internal/resources/control_benchmark.go` | 163 | Walk benchmark children |
| `internal/resources/detection_benchmark.go` | 96 | Walk detection benchmark |
| `internal/resources/dashboard_container.go` | 75 | Walk container children |

### 1.4 Map Iterations (Bulk Access)

#### Dashboards Iteration
| File | Line | Context |
|------|------|---------|
| `internal/dashboardserver/payload.go` | 217 | Build available dashboards list |
| `internal/workspace/powerpipe_workspace.go` | 53 | Verify runtime dependencies |
| `internal/workspace/workspace_events.go` | 81, 261 | Detect changes |
| `internal/resources/mod_resources.go` | 193, 200, 650, 959 | Equals, WalkResources, AddMaps |

#### ControlBenchmarks Iteration
| File | Line | Context |
|------|------|---------|
| `internal/dashboardserver/payload.go` | 230 | Build available benchmarks list |
| `internal/workspace/workspace_events.go` | 121, 286 | Detect changes |
| `internal/resources/mod_resources.go` | 167, 174, 640, 953 | Equals, WalkResources, AddMaps |

#### Controls Iteration
| File | Line | Context |
|------|------|---------|
| `internal/workspace/workspace_events.go` | 131, 291 | Detect changes |
| `internal/resources/mod_resources.go` | 154, 161, 645, 956 | Equals, WalkResources, AddMaps |

#### DetectionBenchmarks Iteration
| File | Line | Context |
|------|------|---------|
| `internal/dashboardserver/payload.go` | 272 | Build available benchmarks |
| `internal/resources/mod_resources.go` | 363, 387, 722 | Equals, WalkResources |

### 1.5 GetChildren() Calls (Tree Traversal)
**Pattern**: Access child resources from parent

| File | Line | Context |
|------|------|---------|
| `internal/dashboardserver/payload.go` | 153, 178 | Add benchmark children to payload |
| `internal/controlexecute/result_group.go` | 130 | Build execution tree |
| `internal/dashboardexecute/dashboard_run.go` | 120 | Create child runs |
| `internal/dashboardexecute/leaf_run.go` | 116 | Create leaf node children |
| `internal/dashboardexecute/container_run.go` | 33 | Create container children |
| `internal/dashboardexecute/check_run.go` | 163 | Process benchmark children |
| `internal/resources/control_benchmark.go` | 80, 108, 164, 204 | Walk/validate benchmark |
| `internal/resources/dashboard.go` | 137, 209, 398 | Walk/validate dashboard |

### 1.6 QueryProviders() Calls
**Pattern**: Get all query providers

| File | Line | Context |
|------|------|---------|
| `internal/resources/mod_resources.go` | 94 | Return all query providers via WalkResources |

---

## 2. Access Pattern Categories

### 2.1 Startup Access (MUST complete before server ready)

These patterns require access to resources during initialization:

#### A. Available Dashboards/Benchmarks List
```
Flow: buildAvailableDashboardsPayload()
Location: internal/dashboardserver/payload.go:201-315

1. Get topLevelResources from mod
2. Iterate ALL Dashboards:
   - Access: FullName, Title, ShortName, Tags, Mod
3. Iterate ALL ControlBenchmarks:
   - Access: FullName, Title, ShortName, Tags, GetParents(), GetChildren()
   - Recursively access all child benchmarks via addBenchmarkChildren()
4. Iterate ALL DetectionBenchmarks:
   - Same pattern as ControlBenchmarks

Resources Accessed: ALL dashboards, ALL benchmarks (recursively)
Access Pattern: Full iteration + tree traversal
Lazy Loading Impact: CRITICAL - must maintain index of dashboard/benchmark names
```

#### B. Server Metadata
```
Flow: buildServerMetadataPayload()
Location: internal/dashboardserver/payload.go:26-95

1. Get workspaceResources
2. Iterate Mods map:
   - Access: GetFullName(), GetTitle(), ShortName

Resources Accessed: All mods
Access Pattern: Iteration
Lazy Loading Impact: LOW - mods are small, can stay eager
```

#### C. Runtime Dependency Verification
```
Flow: verifyResourceRuntimeDependencies()
Location: internal/workspace/powerpipe_workspace.go:52-58

1. Iterate ALL Dashboards
2. For each: call ValidateRuntimeDependencies()
   - Walks all children recursively

Resources Accessed: ALL dashboards + their children
Access Pattern: Full iteration + tree traversal
Lazy Loading Impact: HIGH - either defer verification or maintain dependency graph
```

### 2.2 On-Demand Access (when user requests)

These patterns only access resources when the user initiates an action:

#### A. Dashboard Execution
```
Flow: ExecuteDashboard -> NewDashboardRun
Location: internal/dashboardexecute/dashboard_run.go:39-69

1. Look up dashboard by name (GetResource)
2. Get dashboard inputs
3. Create child runs:
   - Call dashboard.GetChildren()
   - For each child:
     - If Dashboard: recurse
     - If Container: NewDashboardContainerRun
     - If Benchmark/Control: NewCheckRun
     - If Input: NewLeafRun
     - Else: NewLeafRun

Resources Accessed: Single dashboard + ALL descendants
Access Pattern: Tree traversal from root
Lazy Loading Impact: Load dashboard, then resolve children on-demand
```

#### B. Benchmark/Control Execution
```
Flow: NewResultGroup
Location: internal/controlexecute/result_group.go:103-151

1. Look up benchmark/control by name
2. For tree items, iterate GetChildren():
   - If Benchmark: recurse NewResultGroup
   - If Control: AddControl to execution tree
3. Each control:
   - resolveControlQuery() -> ResolveQueryFromQueryProvider()
   - May reference Query resource

Resources Accessed: Single benchmark + descendants + query references
Access Pattern: Tree traversal + reference resolution
Lazy Loading Impact: Load benchmark, resolve children + queries on-demand
```

#### C. Query Resolution
```
Flow: ResolveQueryFromQueryProvider()
Location: internal/workspace/powerpipe_workspace.go:62-105

1. Get SQL/Query from QueryProvider
2. If Query reference exists:
   - Look up named query (GetQueryProvider -> GetResource)
   - Recurse to resolve
3. If SQL is a query name:
   - Look up named query provider
   - Recurse to resolve

Resources Accessed: QueryProvider + referenced Query
Access Pattern: Reference resolution (chain)
Lazy Loading Impact: Must resolve query references on-demand
```

### 2.3 Background/Event Access (file watcher events)

#### A. Dashboard Changed Detection
```
Flow: raiseDashboardChangedEvents()
Location: internal/workspace/workspace_events.go:72-350+

1. Compare ALL resources between old and new ModResources:
   - For each resource type, iterate old map
   - For each: check if exists in new, diff if exists
   - Track deleted resources
2. Then iterate new map for additions

Resources Accessed: ALL resources (both old and new)
Access Pattern: Full iteration (comparison)
Lazy Loading Impact: HIGH - need efficient diff mechanism
```

### 2.4 Reference Resolution Patterns

#### A. Control -> Query
```
Control.SQL may be:
1. Inline SQL (no reference)
2. Query name (string reference to Query resource)
3. Query block reference (object reference)

Resolution: ResolveQueryFromQueryProvider()
```

#### B. Dashboard -> Children
```
Dashboard.GetChildren() returns:
- DashboardWith
- Dashboard (nested)
- DashboardContainer
- DetectionBenchmark
- Benchmark
- Control
- DashboardInput
- DashboardLeafNode types (Card, Chart, Table, etc.)
```

#### C. Benchmark -> Children
```
Benchmark.GetChildren() returns:
- Benchmark (nested)
- Control
```

#### D. Category References
```
Location: internal/resources/dashboard_category_helpers.go:13
Pattern: Access DashboardCategories map for category resolution
```

---

## 3. Critical Flows for Lazy Loading

### 3.1 Flow: Server Startup -> Available Dashboards List

```
┌─────────────────────────────────────────────────────────────────┐
│ Server.Start()                                                  │
├─────────────────────────────────────────────────────────────────┤
│ 1. InitAsync() triggered                                        │
│ 2. buildAvailableDashboardsPayload() called                     │
│    │                                                            │
│    ├─► Iterate ALL Dashboards                                   │
│    │   └─► Access: FullName, Title, ShortName, Tags, Mod        │
│    │                                                            │
│    ├─► Iterate ALL ControlBenchmarks                            │
│    │   ├─► Access: FullName, Title, ShortName, Tags, Parents    │
│    │   └─► RECURSIVE: GetChildren() -> addBenchmarkChildren()   │
│    │       └─► For each child Benchmark: recurse                │
│    │                                                            │
│    └─► Iterate ALL DetectionBenchmarks                          │
│        └─► Same pattern as ControlBenchmarks                    │
│                                                                 │
│ OUTPUT: JSON payload with all dashboard/benchmark names         │
└─────────────────────────────────────────────────────────────────┘

LAZY LOADING REQUIREMENT:
- Need NAME INDEX: Map of resource names -> file locations
- Load: name, title, tags, parent/child relationships (metadata only)
- DO NOT load: full resource definitions, SQL, params
```

### 3.2 Flow: Dashboard Execution

```
┌─────────────────────────────────────────────────────────────────┐
│ ExecuteDashboard(sessionId, dashboardName, inputs)              │
├─────────────────────────────────────────────────────────────────┤
│ 1. GetResource(dashboardName)                                   │
│    └─► Look up in Dashboards map                                │
│                                                                 │
│ 2. NewDashboardRun(dashboard)                                   │
│    │                                                            │
│    ├─► dashboard.GetInputs()                                    │
│    │                                                            │
│    ├─► initWiths() - create with runs                           │
│    │                                                            │
│    └─► createChildRuns()                                        │
│        └─► dashboard.GetChildren()                              │
│            │                                                    │
│            ├─► DashboardWith: skip (handled by initWiths)       │
│            │                                                    │
│            ├─► Dashboard: NewDashboardRun() (recursive)         │
│            │                                                    │
│            ├─► DashboardContainer: NewDashboardContainerRun()   │
│            │   └─► container.GetChildren() (recursive)          │
│            │                                                    │
│            ├─► Benchmark/Control: NewCheckRun()                 │
│            │                                                    │
│            ├─► DashboardInput: NewLeafRun()                     │
│            │                                                    │
│            └─► DashboardLeafNode: NewLeafRun()                  │
│                └─► resolveSQLAndArgs()                          │
│                    └─► May reference Query resource             │
│                                                                 │
│ 3. Execute child runs in parallel                               │
│    └─► Each LeafRun: executeQuery()                             │
│        └─► ResolveQueryFromQueryProvider()                      │
│            └─► May look up named Query                          │
└─────────────────────────────────────────────────────────────────┘

LAZY LOADING REQUIREMENT:
- Load dashboard fully when execution requested
- Load children as they are traversed
- Resolve Query references during execution
- CRITICAL: Children must be fully resolved before execution
```

### 3.3 Flow: Benchmark Execution

```
┌─────────────────────────────────────────────────────────────────┐
│ NewExecutionTree(benchmarkName)                                 │
├─────────────────────────────────────────────────────────────────┤
│ 1. GetResource(benchmarkName)                                   │
│    └─► Look up in ControlBenchmarks map                         │
│                                                                 │
│ 2. NewRootResultGroup(benchmark)                                │
│    │                                                            │
│    └─► If Control: AddControl()                                 │
│        If Benchmark: NewResultGroup() (recursive)               │
│                                                                 │
│ 3. NewResultGroup(benchmark)                                    │
│    └─► benchmark.GetChildren()                                  │
│        │                                                        │
│        ├─► Benchmark: NewResultGroup() (recursive)              │
│        │                                                        │
│        └─► Control: AddControl()                                │
│            └─► NewControlRun(control)                           │
│                                                                 │
│ 4. Execute controls in parallel                                 │
│    └─► ControlRun.execute()                                     │
│        └─► resolveControlQuery(control)                         │
│            └─► ResolveQueryFromQueryProvider()                  │
│                └─► If control.Query: recurse                    │
│                └─► If control.SQL is query name: look up        │
└─────────────────────────────────────────────────────────────────┘

LAZY LOADING REQUIREMENT:
- Load benchmark when execution requested
- Load child benchmarks/controls during tree construction
- Resolve Query references during control execution
```

---

## 4. Dependency Graph

### 4.1 Resource Dependencies

```
Dashboard
├── requires: DashboardInput (global and local)
├── requires: DashboardWith
├── children: Container, Card, Chart, Table, Text, Image, etc.
│             Benchmark, Control, Detection, DetectionBenchmark
└── may reference: Query (via children)

Benchmark
├── children: Benchmark (nested), Control
└── may reference: Query (via Control)

Control
├── may reference: Query (via SQL property or Query property)
└── may reference: Param definitions

Query
└── standalone (no dependencies)

DetectionBenchmark
├── children: DetectionBenchmark (nested), Detection
└── may reference: Query

Detection
└── may reference: Query
```

### 4.2 Resolution Order

For lazy loading, resources must be resolved in this order:

1. **Variables** - First (may be referenced by any resource)
2. **Queries** - Early (referenced by controls, detections, leaf nodes)
3. **Controls/Detections** - After queries
4. **Benchmarks/DetectionBenchmarks** - After controls (reference children)
5. **Dashboard Components** - After queries
6. **Containers** - After components (contain children)
7. **Dashboards** - Last (reference everything)

---

## 5. Lazy Loading Requirements

### 5.1 Must Load Eagerly
- **Resource Index**: Names, titles, tags for all dashboards/benchmarks
- **Parent-Child Relationships**: Benchmark/Dashboard hierarchy structure
- **Variables**: May be used anywhere

### 5.2 Can Load On-Demand
- **Full Resource Definitions**: SQL, params, display properties
- **Query Contents**: SQL statements
- **Dashboard Component Details**: Everything except name

### 5.3 Index Requirements

```go
// Proposed index structure
type ResourceIndex struct {
    // Quick lookup by full name
    ResourceMeta map[string]*ResourceMetadata

    // Parent-child relationships
    Children map[string][]string  // parent -> child names

    // File locations for lazy loading
    FileLocations map[string]ResourceLocation
}

type ResourceMetadata struct {
    FullName    string
    ShortName   string
    Type        string  // "dashboard", "benchmark", "control", etc.
    Title       string
    Tags        map[string]string
    ModFullName string
    IsTopLevel  bool    // For benchmarks
}

type ResourceLocation struct {
    FilePath   string
    StartLine  int
    EndLine    int
}
```

### 5.4 Critical Patterns to Preserve

1. **GetResource()**: Must work with lazy loading
   - Check index first
   - Load from file if not in memory

2. **GetChildren()**: Must resolve child names to resources
   - Return stubs initially
   - Load full resources when accessed

3. **WalkResources()**: May need to trigger loading
   - Consider if full walk is needed
   - Could walk index instead for some operations

4. **QueryProviders()**: Currently walks ALL resources
   - Consider caching or index-based approach

---

## 6. Edge Cases

### 6.1 Circular References
- Dashboard can contain Dashboard (nested)
- Benchmark can contain Benchmark (nested)
- **Risk**: Infinite loops during lazy loading
- **Mitigation**: Track visited resources during resolution

### 6.2 Missing References
- Control.SQL references non-existent query
- Dashboard child references deleted resource
- **Risk**: Runtime errors during lazy load
- **Mitigation**: Validate references during initial index build

### 6.3 File Watcher Updates
- Resource modified while lazy-loaded
- New resource added
- Resource deleted
- **Risk**: Stale cached resources
- **Mitigation**: Clear cache on file change, rebuild index

### 6.4 Concurrent Access
- Multiple executions request same resource
- File watcher modifies while execution in progress
- **Risk**: Race conditions
- **Mitigation**: Mutex per resource, copy-on-write for updates

---

## 7. Recommendations

### 7.1 Phase 1: Index-Based Loading
1. Build lightweight index during workspace load
2. Index contains: names, titles, tags, parent-child relationships
3. Full resources loaded on-demand via GetResource()

### 7.2 Phase 2: Smart Prefetching
1. When dashboard requested, prefetch all children
2. When benchmark requested, prefetch all controls + queries
3. Background load commonly-used resources

### 7.3 Phase 3: Incremental Updates
1. File watcher triggers index update only
2. Invalidate cached resources for modified files
3. Re-parse only affected resources

### 7.4 Performance Targets
- Index build: < 100ms for 1000 resources
- Single resource load: < 10ms
- Dashboard execution (cold): < 500ms for 100 child resources
- File watcher update: < 50ms

---

## Appendix A: All Access Point Locations

### Dashboards Map Access
```
internal/dashboardserver/payload.go:217,220
internal/workspace/load_workspace_test.go:67,96,218,239
internal/workspace/powerpipe_workspace.go:53
internal/workspace/workspace_events.go:81,82,261,262
internal/resources/mod_resources.go:193,194,200,201,448,650,785,789,959,960
```

### ControlBenchmarks Map Access
```
internal/dashboardserver/payload.go:230,262,265,267
internal/workspace/load_workspace_test.go:107,122,133,137
internal/workspace/workspace_events.go:121,122,286,287
internal/resources/mod_resources.go:167,168,174,175,440,640,777,781,953,954
```

### Controls Map Access
```
internal/workspace/load_workspace_test.go:101
internal/workspace/workspace_events.go:131,132,291,292
internal/resources/mod_resources.go:154,155,161,162,446,627,645,769,773,956,957
```

### Queries Map Access
```
internal/workspace/load_workspace_test.go:61
internal/resources/mod_resources.go:142,148,149,480,623,761,765,1017
```

### GetChildren() Calls
```
internal/workspace/load_workspace_test.go:69,98,109,124,127,135,139
internal/controlexecute/result_group.go:130,219
internal/controlexecute/control_run.go:191
internal/controldisplay/control.go:43
internal/controldisplay/detection.go:45
internal/controldisplay/group.go:42,160
internal/controldisplay/detection_group.go:41,157
internal/controldisplay/snapshot.go:166
internal/dashboardserver/payload.go:153,178
internal/dashboardserver/backend_support.go:81
internal/dashboardexecute/dashboard_run.go:120
internal/dashboardexecute/leaf_run.go:116
internal/dashboardexecute/container_run.go:33
internal/dashboardexecute/check_run.go:163
internal/dashboardexecute/detection_benchmark_run.go:32
internal/dashboardexecute/detection_benchmark_display.go:102,206
internal/resources/dashboard.go:137,209,398
internal/resources/dashboard_container.go:42,43
internal/resources/dashboard_flow.go:53
internal/resources/dashboard_graph.go:57
internal/resources/dashboard_hierarchy.go:52
internal/resources/control_benchmark.go:80,108,164,204
internal/resources/detection_benchmark.go:137
```

---

## Appendix B: Resource Type Hierarchy

```
modconfig.ModResources (interface)
└── resources.PowerpipeModResources (implementation)
    ├── Dashboards map[string]*Dashboard
    │   └── Dashboard implements:
    │       - ModTreeItem (GetChildren, GetParents)
    │       - ResourceWithMetadata
    │       - DashboardLeafNode
    │       - RuntimeDependencyProvider
    │
    ├── ControlBenchmarks map[string]*Benchmark
    │   └── Benchmark implements:
    │       - ModTreeItem
    │       - ResourceWithMetadata
    │
    ├── Controls map[string]*Control
    │   └── Control implements:
    │       - ModTreeItem
    │       - QueryProvider
    │       - DashboardLeafNode
    │
    ├── Queries map[string]*Query
    │   └── Query implements:
    │       - QueryProvider
    │       - HclResource
    │
    └── [20+ other resource type maps]
```
