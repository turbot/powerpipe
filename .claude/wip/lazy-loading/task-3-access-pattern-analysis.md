# Task 3: Resource Access Pattern Analysis

## Objective

Analyze and document all code paths that access resources from the workspace to ensure lazy loading handles every case correctly.

## Context

- Lazy loading changes how resources are accessed
- Code that assumes all resources are immediately available will break
- Need to identify: what resources are accessed, when, and in what order
- This analysis informs the design of dependency resolution (Task 8)

## Dependencies

### Prerequisites
- None (this is a foundation task)

### Files to Create
- `.claude/wip/lazy-loading/access_pattern_report.md` - Analysis results
- `internal/workspace/access_tracking_test.go` - Tests with access tracking

### Files to Review (not modify)
- All files in `internal/workspace/`
- All files in `internal/dashboardserver/`
- All files in `internal/dashboardexecute/`
- All files in `internal/controlexecute/`
- All files in `internal/cmd/`

## Implementation Details

### 1. Identify All Resource Access Points

Search the codebase for patterns that access resources:

```bash
# Map access patterns
grep -rn "GetModResources()" internal/
grep -rn "\.Dashboards\[" internal/
grep -rn "\.Controls\[" internal/
grep -rn "\.Queries\[" internal/
grep -rn "\.ControlBenchmarks\[" internal/
grep -rn "WalkResources" internal/
grep -rn "GetResource(" internal/
grep -rn "QueryProviders()" internal/
```

### 2. Categorize Access Patterns

Document each access point with:
- **Location**: File and line number
- **Context**: What operation triggers this access
- **Pattern**: Single resource, iteration, or bulk access
- **Timing**: Startup, on-demand, or background
- **Dependencies**: What other resources are accessed together

### 3. Access Pattern Categories

#### A. Startup Access (must complete before server ready)
- Available dashboards list
- Available benchmarks list
- Server metadata
- Variable resolution

#### B. On-Demand Access (when user requests)
- Dashboard execution
- Benchmark execution
- Query execution
- Resource inspection

#### C. Bulk Access (iteration over all resources)
- WalkResources callbacks
- Resource validation
- Snapshot generation

#### D. Reference Resolution (resource A needs resource B)
- Control → Query reference
- Dashboard → Child panels
- Benchmark → Child controls/benchmarks
- Category references
- Input dependencies

### 4. Document Key Access Flows

```markdown
## Flow: Dashboard List (Server Startup)

1. Server starts
2. `buildAvailableDashboardsPayload()` called
3. Iterates `workspaceResources.Dashboards`
4. For each dashboard: access Name, Title, Tags, ModFullName
5. Iterates `workspaceResources.ControlBenchmarks`
6. For each benchmark:
   - Access Name, Title, Tags, IsTopLevel, Parents
   - Recursively access children via `GetChildren()`
7. Returns JSON payload

**Resources Accessed**: All dashboards, all benchmarks
**Access Pattern**: Full iteration
**Lazy Loading Impact**: Need dashboard/benchmark index for list

---

## Flow: Dashboard Execution

1. User selects dashboard
2. `ExecuteDashboard()` called with dashboard name
3. Dashboard resource looked up by name
4. Dashboard children accessed (`GetChildren()`)
5. For each child:
   - If QueryProvider: resolve SQL and query references
   - If Container: recursively access children
6. Create execution tree
7. Execute queries

**Resources Accessed**: Single dashboard + all descendants
**Access Pattern**: Tree traversal from root
**Lazy Loading Impact**: Load dashboard and resolve children on-demand

---

## Flow: Benchmark Execution

1. User runs benchmark
2. Benchmark resource looked up by name
3. Benchmark children accessed (`GetChildren()`)
4. For each child control:
   - Access Query reference
   - Access SQL
   - Access Args/Params
5. Create control execution tree
6. Execute all controls

**Resources Accessed**: Single benchmark + all descendant controls + their queries
**Access Pattern**: Tree traversal with query resolution
**Lazy Loading Impact**: Load benchmark, resolve children, resolve query refs
```

### 5. Create Access Tracking Tests

```go
// internal/workspace/access_tracking_test.go
package workspace

import (
    "testing"
    "sync"
)

// AccessTracker records resource accesses for analysis
type AccessTracker struct {
    mu       sync.Mutex
    accesses []AccessRecord
}

type AccessRecord struct {
    ResourceType string
    ResourceName string
    AccessType   string // "get", "iterate", "children"
    CallerFunc   string
    CallerFile   string
    CallerLine   int
}

func (t *AccessTracker) Record(resourceType, resourceName, accessType string) {
    t.mu.Lock()
    defer t.mu.Unlock()

    // Get caller info
    _, file, line, _ := runtime.Caller(2)

    t.accesses = append(t.accesses, AccessRecord{
        ResourceType: resourceType,
        ResourceName: resourceName,
        AccessType:   accessType,
        CallerFile:   file,
        CallerLine:   line,
    })
}

func (t *AccessTracker) Report() string {
    t.mu.Lock()
    defer t.mu.Unlock()

    var b strings.Builder
    for _, a := range t.accesses {
        b.WriteString(fmt.Sprintf("%s.%s [%s] from %s:%d\n",
            a.ResourceType, a.ResourceName, a.AccessType,
            filepath.Base(a.CallerFile), a.CallerLine))
    }
    return b.String()
}

func TestAccessPatterns_AvailableDashboards(t *testing.T) {
    tracker := &AccessTracker{}

    // Load workspace with tracking
    ws := loadWorkspaceWithTracking(t, "behavior_test_mod", tracker)

    // Build available dashboards payload (simulated)
    resources := ws.GetModResources().(*resources.PowerpipeModResources)

    // Track dashboard iteration
    for name := range resources.Dashboards {
        tracker.Record("dashboard", name, "iterate")
    }

    // Track benchmark iteration
    for name, bench := range resources.ControlBenchmarks {
        tracker.Record("benchmark", name, "iterate")
        trackChildAccess(tracker, bench)
    }

    t.Log("Access Pattern Report:\n" + tracker.Report())

    // Verify expected access patterns
    assert.True(t, tracker.HasAccess("dashboard", "*", "iterate"))
    assert.True(t, tracker.HasAccess("benchmark", "*", "iterate"))
}

func TestAccessPatterns_DashboardExecution(t *testing.T) {
    tracker := &AccessTracker{}
    ws := loadWorkspaceWithTracking(t, "behavior_test_mod", tracker)

    // Execute a dashboard
    resources := ws.GetModResources().(*resources.PowerpipeModResources)
    dash := resources.Dashboards["behavior_test.dashboard.main"]

    tracker.Record("dashboard", dash.FullName, "get")

    // Track child access
    trackChildAccess(tracker, dash)

    t.Log("Access Pattern Report:\n" + tracker.Report())

    // Verify: should NOT iterate all dashboards, only access one
    dashIterates := tracker.CountAccesses("dashboard", "*", "iterate")
    assert.Equal(t, 0, dashIterates, "Should not iterate all dashboards")
}

func trackChildAccess(tracker *AccessTracker, resource modconfig.ModTreeItem) {
    for _, child := range resource.GetChildren() {
        tracker.Record(child.GetBlockType(), child.Name(), "children")
        if treeItem, ok := child.(modconfig.ModTreeItem); ok {
            trackChildAccess(tracker, treeItem)
        }
    }
}
```

### 6. Generate Analysis Report

Create comprehensive report documenting:

1. **All Access Points** - Every place in code that accesses resources
2. **Access Patterns** - How resources are accessed (single, iterate, tree)
3. **Timing** - When access happens (startup, on-demand)
4. **Dependencies** - What resources need other resources
5. **Lazy Loading Requirements** - What must change for each pattern

## Acceptance Criteria

- [ ] All resource access points identified in codebase
- [ ] Access patterns categorized (startup, on-demand, bulk, reference)
- [ ] Key flows documented with sequence of accesses
- [ ] Access tracking tests created
- [ ] Report generated: `.claude/wip/lazy-loading/access_pattern_report.md`
- [ ] Dependencies between resources documented
- [ ] Recommendations for lazy loading design included
- [ ] Edge cases identified (circular refs, missing refs, etc.)

## Notes

- This analysis directly informs Task 8 (Dependency Resolution)
- Focus on what resources are accessed together
- Identify any patterns that assume all resources are loaded
- Document any global iterations that could be avoided
