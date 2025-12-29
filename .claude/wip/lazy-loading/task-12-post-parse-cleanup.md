# Task 12: Post-Parse Cleanup (Remain Fields)

## Objective

Clear the HCL `Remain` fields from all resources after parsing to free memory. This is a quick win that works alongside lazy loading to reduce memory usage.

## Context

- **33 total files** have `Remain hcl.Body` fields across both repositories
- `Remain` holds unparsed HCL AST after parsing
- This AST is never used after initial parsing
- Clearing it can free significant memory (~30-40% reduction)
- This optimization works with both eager and lazy loading

## Repository

**This task spans BOTH repositories:**
- **pipe-fittings**: 8 files with `Remain` fields (base types)
- **powerpipe**: 25 files with `Remain` fields (resource types)

Changes should be coordinated - pipe-fittings PR merged first.

## Dependencies

### Prerequisites
- Task 1 (Behavior Tests) - Ensure cleanup doesn't break anything

### Files to Modify (pipe-fittings) - 8 files
Base types and infrastructure:
- `modconfig/mod.go` - Add ClearRemain method
- `modconfig/variable.go` - Add ClearRemain method
- `modconfig/local.go` - Add ClearRemain method
- `modconfig/resource_metadata.go` - Clear ResourceMetadataRemain
- `modconfig/hcl_resource_impl.go` - Add ClearRemain method
- `modconfig/mod_tree_item_impl.go` - Add ClearRemain method
- `modconfig/resource_with_metadata_impl.go` - Add ClearRemain method
- `parse/mod.go` - Call cleanup after parsing complete

### Files to Create (pipe-fittings)
- `modconfig/remain_cleaner.go` - RemainCleaner interface

### Files to Modify (powerpipe) - 25 files
Resource types:
- `internal/resources/dashboard.go`
- `internal/resources/control.go`
- `internal/resources/query.go`
- `internal/resources/control_benchmark.go`
- `internal/resources/dashboard_card.go`
- `internal/resources/dashboard_chart.go`
- `internal/resources/dashboard_container.go`
- `internal/resources/dashboard_flow.go`
- `internal/resources/dashboard_graph.go`
- `internal/resources/dashboard_hierarchy.go`
- `internal/resources/dashboard_image.go`
- `internal/resources/dashboard_input.go`
- `internal/resources/dashboard_table.go`
- `internal/resources/dashboard_text.go`
- `internal/resources/dashboard_category.go`
- `internal/resources/dashboard_node.go`
- `internal/resources/dashboard_edge.go`
- `internal/resources/detection.go`
- `internal/resources/detection_benchmark.go`
- `internal/resources/node_edge_provider_impl.go`
- `internal/resources/with_provider_impl.go`
- `internal/resources/runtime_dependency_provider_impl.go`
- `internal/resources/query_provider_impl.go`
- `internal/resources/dashboard_leaf_node_impl.go`
- `internal/resources/dashboard_with.go`

### Files to Create (powerpipe)
- `internal/resources/remain_cleaner.go` - Powerpipe-specific cleanup utilities
- `internal/resources/remain_cleaner_test.go` - Tests

## Implementation Details

### 1. Identify All Files with Remain

```bash
# pipe-fittings (8 files)
grep -l "Remain.*hcl.Body" /path/to/pipe-fittings/modconfig/*.go
# Results:
# variable.go, resource_with_metadata_impl.go, resource_metadata.go,
# mod_tree_item_impl.go, mod.go, local.go, interfaces.go, hcl_resource_impl.go

# powerpipe (25 files)
grep -l "Remain.*hcl.Body" /path/to/powerpipe/internal/resources/*.go
# Results:
# dashboard.go, control.go, query.go, control_benchmark.go,
# dashboard_card.go, dashboard_chart.go, dashboard_container.go,
# dashboard_flow.go, dashboard_graph.go, dashboard_hierarchy.go,
# dashboard_image.go, dashboard_input.go, dashboard_table.go,
# dashboard_text.go, dashboard_category.go, dashboard_node.go,
# dashboard_edge.go, detection.go, detection_benchmark.go,
# node_edge_provider_impl.go, with_provider_impl.go,
# runtime_dependency_provider_impl.go, query_provider_impl.go,
# dashboard_leaf_node_impl.go, dashboard_with.go
```

### 2. RemainCleaner Interface

```go
// pipe-fittings/modconfig/remain_cleaner.go
package modconfig

import (
    "github.com/hashicorp/hcl/v2"
)

// RemainCleaner is implemented by resources that hold HCL Remain
type RemainCleaner interface {
    // ClearRemain clears the Remain hcl.Body field to free memory
    ClearRemain()
}

// ClearRemainRecursive clears Remain from a resource and all children
func ClearRemainRecursive(resource HclResource) {
    // Clear this resource
    if cleaner, ok := resource.(RemainCleaner); ok {
        cleaner.ClearRemain()
    }

    // Clear children if ModTreeItem
    if treeItem, ok := resource.(ModTreeItem); ok {
        for _, child := range treeItem.GetChildren() {
            if childRes, ok := child.(HclResource); ok {
                ClearRemainRecursive(childRes)
            }
        }
    }
}

// ClearAllRemain clears Remain from all resources in a mod
func ClearAllRemain(mod *Mod) {
    if mod == nil {
        return
    }

    // Walk all resources
    mod.WalkResources(func(resource HclResource) (bool, error) {
        if cleaner, ok := resource.(RemainCleaner); ok {
            cleaner.ClearRemain()
        }
        return true, nil
    })
}
```

### 3. Add ClearRemain to Each Resource Type

```go
// pipe-fittings/modconfig/dashboard.go
type Dashboard struct {
    // ... existing fields
    Remain hcl.Body `hcl:",remain" json:"-"`
}

// ClearRemain clears the Remain field
func (d *Dashboard) ClearRemain() {
    d.Remain = nil
}

// pipe-fittings/modconfig/query.go
type Query struct {
    // ... existing fields
    Remain hcl.Body `hcl:",remain" json:"-"`
}

func (q *Query) ClearRemain() {
    q.Remain = nil
}

// pipe-fittings/modconfig/control.go
type Control struct {
    // ... existing fields
    Remain hcl.Body `hcl:",remain" json:"-"`
}

func (c *Control) ClearRemain() {
    c.Remain = nil
}

// ... repeat for all 25+ resource types with Remain field:
// benchmark.go, card.go, category.go, chart.go, container.go,
// dashboard_category.go, dashboard_chart.go, dashboard_container.go,
// dashboard_edge.go, dashboard_flow.go, dashboard_graph.go,
// dashboard_hierarchy.go, dashboard_image.go, dashboard_input.go,
// dashboard_node.go, dashboard_table.go, dashboard_text.go,
// detection.go, detection_benchmark.go, edge.go, flow.go,
// graph.go, hierarchy.go, image.go, input.go, node.go,
// table.go, text.go, variable.go
```

### 4. Call Cleanup After Parsing

```go
// pipe-fittings/parse/run_context.go modifications

// FinalizeParsing is called after all HCL parsing is complete
func (r *RunContext) FinalizeParsing() {
    // Clear Remain fields to free HCL AST memory
    r.ClearAllRemainFields()
}

func (r *RunContext) ClearAllRemainFields() {
    if r.Mod == nil {
        return
    }

    modconfig.ClearAllRemain(r.Mod)
}
```

### 5. Integration Point

```go
// pipe-fittings/parse/parser.go modifications

func ParseMod(ctx context.Context, path string, opts ParseModOptions) (*modconfig.Mod, error) {
    // ... existing parsing logic

    runCtx := NewRunContext(...)

    // Parse all files
    if err := runCtx.ParseFiles(); err != nil {
        return nil, err
    }

    // Resolve references
    if err := runCtx.ResolveReferences(); err != nil {
        return nil, err
    }

    // NEW: Clean up Remain fields after parsing is complete
    runCtx.FinalizeParsing()

    return runCtx.Mod, nil
}
```

### 6. Tests

```go
// pipe-fittings/modconfig/remain_cleaner_test.go
package modconfig

import (
    "testing"

    "github.com/hashicorp/hcl/v2"
    "github.com/stretchr/testify/assert"
)

func TestDashboard_ClearRemain(t *testing.T) {
    dash := &Dashboard{
        Remain: &hcl.BodyContent{}, // Mock body
    }

    assert.NotNil(t, dash.Remain)

    dash.ClearRemain()

    assert.Nil(t, dash.Remain)
}

func TestQuery_ClearRemain(t *testing.T) {
    query := &Query{
        Remain: &hcl.BodyContent{},
    }

    assert.NotNil(t, query.Remain)

    query.ClearRemain()

    assert.Nil(t, query.Remain)
}

func TestClearRemainRecursive(t *testing.T) {
    // Dashboard with children
    card := &Card{
        Remain: &hcl.BodyContent{},
    }
    dash := &Dashboard{
        Remain: &hcl.BodyContent{},
        children: []ModTreeItem{card},
    }

    ClearRemainRecursive(dash)

    assert.Nil(t, dash.Remain)
    assert.Nil(t, card.Remain)
}

func TestClearAllRemain_Mod(t *testing.T) {
    mod := &Mod{
        Dashboards: map[string]*Dashboard{
            "dash1": {Remain: &hcl.BodyContent{}},
            "dash2": {Remain: &hcl.BodyContent{}},
        },
        Queries: map[string]*Query{
            "query1": {Remain: &hcl.BodyContent{}},
        },
        Controls: map[string]*Control{
            "ctrl1": {Remain: &hcl.BodyContent{}},
        },
    }

    ClearAllRemain(mod)

    for _, dash := range mod.Dashboards {
        assert.Nil(t, dash.Remain)
    }
    for _, query := range mod.Queries {
        assert.Nil(t, query.Remain)
    }
    for _, ctrl := range mod.Controls {
        assert.Nil(t, ctrl.Remain)
    }
}

func TestClearRemain_MemoryReduction(t *testing.T) {
    // Create mod with many resources
    mod := createLargeTestMod(t, 100)

    // Measure memory before
    runtime.GC()
    var before runtime.MemStats
    runtime.ReadMemStats(&before)

    // Clear Remain
    ClearAllRemain(mod)

    // Measure memory after
    runtime.GC()
    runtime.GC()
    var after runtime.MemStats
    runtime.ReadMemStats(&after)

    // Should reduce memory
    reduction := before.HeapAlloc - after.HeapAlloc
    t.Logf("Memory reduction: %d bytes", reduction)
    assert.Greater(t, reduction, uint64(0), "Should reduce memory")
}
```

### 7. Memory Verification Test

```go
// internal/workspace/remain_cleanup_test.go
package workspace

import (
    "context"
    "runtime"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestRemainCleanup_MemoryImpact(t *testing.T) {
    modPath := setupLargeMod(t, 200)

    // Load WITHOUT cleanup (for comparison)
    runtime.GC()
    var beforeNoClean runtime.MemStats
    runtime.ReadMemStats(&beforeNoClean)

    wsNoClean, err := LoadWorkspaceNoCleanup(context.Background(), modPath)
    require.NoError(t, err)
    _ = wsNoClean

    runtime.GC()
    var afterNoClean runtime.MemStats
    runtime.ReadMemStats(&afterNoClean)

    noCleanMem := afterNoClean.HeapAlloc - beforeNoClean.HeapAlloc

    // Load WITH cleanup
    runtime.GC()
    var beforeClean runtime.MemStats
    runtime.ReadMemStats(&beforeClean)

    wsClean, err := LoadWorkspace(context.Background(), modPath)
    require.NoError(t, err)
    _ = wsClean

    runtime.GC()
    var afterClean runtime.MemStats
    runtime.ReadMemStats(&afterClean)

    cleanMem := afterClean.HeapAlloc - beforeClean.HeapAlloc

    t.Logf("Without cleanup: %d bytes", noCleanMem)
    t.Logf("With cleanup: %d bytes", cleanMem)
    t.Logf("Reduction: %.1f%%", float64(noCleanMem-cleanMem)/float64(noCleanMem)*100)

    // Should reduce memory by at least 20%
    assert.Less(t, cleanMem, noCleanMem*80/100,
        "Cleanup should reduce memory by at least 20%")
}
```

## Acceptance Criteria

- [ ] All 25+ resource types implement ClearRemain()
- [ ] ClearRemainRecursive clears nested resources
- [ ] ClearAllRemain walks and clears entire mod
- [ ] Parser calls cleanup after parsing complete
- [ ] Memory reduction measurable (> 20% for large mods)
- [ ] All behavior tests still pass
- [ ] No functionality broken by clearing Remain
- [ ] Works with both eager and lazy loading

## Notes

- This is safe because Remain is never accessed after parsing
- Can be enabled/disabled via flag if needed
- Should be called AFTER all reference resolution is complete
- May want to add timing metrics for cleanup phase
- Consider adding warning log if Remain accessed after cleanup
