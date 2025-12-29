# Task 4: Resource Index Design & Implementation

## Objective

Design and implement a lightweight resource index that stores minimal metadata about all resources without parsing full HCL. This index enables the dashboard list and resource lookup without loading full resources into memory.

## Context

- The index is the foundation of lazy loading
- Must be small enough to always fit in memory (~1 KB per 100 resources)
- Must contain enough info for dashboard/benchmark list UI
- Must enable efficient lookup by name for on-demand loading

## Repository

**This task is Powerpipe-only.** No changes to pipe-fittings required.

## Dependencies

### Prerequisites
- Task 1 (Behavior Tests) - Need safety net
- Task 2 (Memory Benchmarks) - Need to measure impact
- Task 3 (Access Patterns) - Need to know what metadata is needed

### Files to Create (powerpipe)
- `internal/resourceindex/index.go` - Core index structures
- `internal/resourceindex/entry.go` - Index entry types
- `internal/resourceindex/index_test.go` - Unit tests

### Files to Modify (powerpipe)
- `internal/workspace/workspace.go` - Add ResourceIndex field to workspace (not mod)

## Implementation Details

### 1. Index Entry Structure

```go
// internal/resourceindex/entry.go
package resourceindex

// IndexEntry contains minimal metadata about a resource
type IndexEntry struct {
    // Identity
    Type      string `json:"type"`       // "dashboard", "query", "control", etc.
    Name      string `json:"name"`       // Full name: "mod.type.shortname"
    ShortName string `json:"short_name"` // Just the short name

    // Display metadata (for UI lists)
    Title       string            `json:"title,omitempty"`
    Description string            `json:"description,omitempty"`
    Tags        map[string]string `json:"tags,omitempty"`

    // Hierarchy (for benchmarks/dashboards)
    IsTopLevel bool     `json:"is_top_level,omitempty"`
    ParentName string   `json:"parent_name,omitempty"`
    ChildNames []string `json:"child_names,omitempty"`

    // Source location (for on-demand loading)
    FileName        string `json:"file_name"`
    StartLine       int    `json:"start_line"`
    EndLine         int    `json:"end_line"`
    ByteOffset      int64  `json:"byte_offset"`  // For efficient seeking
    ByteLength      int    `json:"byte_length"`

    // Mod info
    ModName     string `json:"mod_name"`
    ModFullName string `json:"mod_full_name"`

    // Type-specific metadata
    // For benchmarks
    BenchmarkType string `json:"benchmark_type,omitempty"` // "control" or "detection"

    // For queries/controls
    HasSQL bool `json:"has_sql,omitempty"`

    // For inputs
    DashboardName string `json:"dashboard_name,omitempty"` // For scoped inputs
}

// Size returns approximate memory size of this entry
func (e *IndexEntry) Size() int {
    // Rough estimate: strings + overhead
    size := 100 // base overhead
    size += len(e.Type) + len(e.Name) + len(e.ShortName)
    size += len(e.Title) + len(e.Description)
    size += len(e.FileName) + len(e.ModName)
    for k, v := range e.Tags {
        size += len(k) + len(v)
    }
    for _, c := range e.ChildNames {
        size += len(c)
    }
    return size
}
```

### 2. Index Structure

```go
// internal/resourceindex/index.go
package resourceindex

import (
    "sync"
)

// ResourceIndex provides fast lookup of resource metadata
type ResourceIndex struct {
    mu sync.RWMutex

    // All entries keyed by full name
    entries map[string]*IndexEntry

    // Type-specific indexes for efficient filtering
    byType map[string]map[string]*IndexEntry // type -> name -> entry

    // Mod information
    ModName     string
    ModFullName string
    ModTitle    string

    // Statistics
    totalSize int
}

// NewResourceIndex creates an empty index
func NewResourceIndex() *ResourceIndex {
    return &ResourceIndex{
        entries: make(map[string]*IndexEntry),
        byType:  make(map[string]map[string]*IndexEntry),
    }
}

// Add adds an entry to the index
func (idx *ResourceIndex) Add(entry *IndexEntry) {
    idx.mu.Lock()
    defer idx.mu.Unlock()

    idx.entries[entry.Name] = entry

    // Add to type index
    if idx.byType[entry.Type] == nil {
        idx.byType[entry.Type] = make(map[string]*IndexEntry)
    }
    idx.byType[entry.Type][entry.Name] = entry

    idx.totalSize += entry.Size()
}

// Get retrieves an entry by full name
func (idx *ResourceIndex) Get(name string) (*IndexEntry, bool) {
    idx.mu.RLock()
    defer idx.mu.RUnlock()

    entry, ok := idx.entries[name]
    return entry, ok
}

// GetByType retrieves all entries of a specific type
func (idx *ResourceIndex) GetByType(resourceType string) []*IndexEntry {
    idx.mu.RLock()
    defer idx.mu.RUnlock()

    typeMap := idx.byType[resourceType]
    if typeMap == nil {
        return nil
    }

    entries := make([]*IndexEntry, 0, len(typeMap))
    for _, entry := range typeMap {
        entries = append(entries, entry)
    }
    return entries
}

// Dashboards returns all dashboard entries
func (idx *ResourceIndex) Dashboards() []*IndexEntry {
    return idx.GetByType("dashboard")
}

// Benchmarks returns all benchmark entries (control and detection)
func (idx *ResourceIndex) Benchmarks() []*IndexEntry {
    controlBenchmarks := idx.GetByType("benchmark")
    detectionBenchmarks := idx.GetByType("detection_benchmark")
    return append(controlBenchmarks, detectionBenchmarks...)
}

// Queries returns all query entries
func (idx *ResourceIndex) Queries() []*IndexEntry {
    return idx.GetByType("query")
}

// Controls returns all control entries
func (idx *ResourceIndex) Controls() []*IndexEntry {
    return idx.GetByType("control")
}

// TopLevelBenchmarks returns benchmarks that are direct children of mod
func (idx *ResourceIndex) TopLevelBenchmarks() []*IndexEntry {
    var result []*IndexEntry
    for _, entry := range idx.Benchmarks() {
        if entry.IsTopLevel {
            result = append(result, entry)
        }
    }
    return result
}

// GetChildren returns child entries for a parent
func (idx *ResourceIndex) GetChildren(parentName string) []*IndexEntry {
    idx.mu.RLock()
    defer idx.mu.RUnlock()

    parent, ok := idx.entries[parentName]
    if !ok || len(parent.ChildNames) == 0 {
        return nil
    }

    children := make([]*IndexEntry, 0, len(parent.ChildNames))
    for _, childName := range parent.ChildNames {
        if child, ok := idx.entries[childName]; ok {
            children = append(children, child)
        }
    }
    return children
}

// Size returns total approximate memory size of index
func (idx *ResourceIndex) Size() int {
    idx.mu.RLock()
    defer idx.mu.RUnlock()
    return idx.totalSize
}

// Count returns total number of entries
func (idx *ResourceIndex) Count() int {
    idx.mu.RLock()
    defer idx.mu.RUnlock()
    return len(idx.entries)
}

// Stats returns index statistics
func (idx *ResourceIndex) Stats() IndexStats {
    idx.mu.RLock()
    defer idx.mu.RUnlock()

    stats := IndexStats{
        TotalEntries: len(idx.entries),
        TotalSize:    idx.totalSize,
        ByType:       make(map[string]int),
    }

    for typeName, entries := range idx.byType {
        stats.ByType[typeName] = len(entries)
    }

    return stats
}

type IndexStats struct {
    TotalEntries int
    TotalSize    int
    ByType       map[string]int
}
```

### 3. Available Dashboards from Index

```go
// internal/resourceindex/payload.go
package resourceindex

// BuildAvailableDashboardsFromIndex builds the dashboard list payload
// without loading full resources
func (idx *ResourceIndex) BuildAvailableDashboardsPayload() *AvailableDashboardsPayload {
    payload := &AvailableDashboardsPayload{
        Action:     "available_dashboards",
        Dashboards: make(map[string]DashboardInfo),
        Benchmarks: make(map[string]BenchmarkInfo),
    }

    // Build dashboard list from index
    for _, entry := range idx.Dashboards() {
        payload.Dashboards[entry.Name] = DashboardInfo{
            Title:       entry.Title,
            FullName:    entry.Name,
            ShortName:   entry.ShortName,
            Tags:        entry.Tags,
            ModFullName: entry.ModFullName,
        }
    }

    // Build benchmark list with hierarchy from index
    benchmarkTrunks := make(map[string][][]string)

    for _, entry := range idx.Benchmarks() {
        info := BenchmarkInfo{
            Title:         entry.Title,
            FullName:      entry.Name,
            ShortName:     entry.ShortName,
            BenchmarkType: entry.BenchmarkType,
            Tags:          entry.Tags,
            IsTopLevel:    entry.IsTopLevel,
            ModFullName:   entry.ModFullName,
        }

        // Build children recursively from index
        info.Children = idx.buildBenchmarkChildren(entry, entry.IsTopLevel,
            []string{entry.Name}, benchmarkTrunks)

        payload.Benchmarks[entry.Name] = info
    }

    // Apply trunks
    for name, trunks := range benchmarkTrunks {
        if info, ok := payload.Benchmarks[name]; ok {
            info.Trunks = trunks
            payload.Benchmarks[name] = info
        }
    }

    return payload
}

func (idx *ResourceIndex) buildBenchmarkChildren(parent *IndexEntry,
    recordTrunk bool, trunk []string, trunks map[string][][]string) []BenchmarkInfo {

    var children []BenchmarkInfo

    for _, childEntry := range idx.GetChildren(parent.Name) {
        // Only include benchmark children (not controls)
        if childEntry.Type != "benchmark" && childEntry.Type != "detection_benchmark" {
            continue
        }

        childTrunk := append([]string{}, trunk...)
        childTrunk = append(childTrunk, childEntry.Name)

        if recordTrunk {
            trunks[childEntry.Name] = append(trunks[childEntry.Name], childTrunk)
        }

        info := BenchmarkInfo{
            Title:         childEntry.Title,
            FullName:      childEntry.Name,
            ShortName:     childEntry.ShortName,
            BenchmarkType: childEntry.BenchmarkType,
            Tags:          childEntry.Tags,
            Children:      idx.buildBenchmarkChildren(childEntry, recordTrunk, childTrunk, trunks),
        }

        children = append(children, info)
    }

    return children
}

// Payload types (match existing server types)
type AvailableDashboardsPayload struct {
    Action     string                   `json:"action"`
    Dashboards map[string]DashboardInfo `json:"dashboards"`
    Benchmarks map[string]BenchmarkInfo `json:"benchmarks"`
    Snapshots  map[string]string        `json:"snapshots,omitempty"`
}

type DashboardInfo struct {
    Title       string            `json:"title,omitempty"`
    FullName    string            `json:"full_name"`
    ShortName   string            `json:"short_name"`
    Tags        map[string]string `json:"tags,omitempty"`
    ModFullName string            `json:"mod_full_name,omitempty"`
}

type BenchmarkInfo struct {
    Title         string            `json:"title,omitempty"`
    FullName      string            `json:"full_name"`
    ShortName     string            `json:"short_name"`
    BenchmarkType string            `json:"benchmark_type,omitempty"`
    Tags          map[string]string `json:"tags,omitempty"`
    IsTopLevel    bool              `json:"is_top_level,omitempty"`
    Trunks        [][]string        `json:"trunks,omitempty"`
    Children      []BenchmarkInfo   `json:"children,omitempty"`
    ModFullName   string            `json:"mod_full_name,omitempty"`
}
```

### 4. Unit Tests

```go
// internal/resourceindex/index_test.go
package resourceindex

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestIndex_AddAndGet(t *testing.T) {
    idx := NewResourceIndex()

    entry := &IndexEntry{
        Type:      "dashboard",
        Name:      "mymod.dashboard.test",
        ShortName: "test",
        Title:     "Test Dashboard",
        FileName:  "/path/to/mod/dashboards.pp",
        StartLine: 10,
        EndLine:   50,
    }

    idx.Add(entry)

    // Get by name
    got, ok := idx.Get("mymod.dashboard.test")
    assert.True(t, ok)
    assert.Equal(t, "Test Dashboard", got.Title)

    // Get by type
    dashboards := idx.Dashboards()
    assert.Len(t, dashboards, 1)
    assert.Equal(t, "mymod.dashboard.test", dashboards[0].Name)
}

func TestIndex_BenchmarkHierarchy(t *testing.T) {
    idx := NewResourceIndex()

    // Add parent benchmark
    parent := &IndexEntry{
        Type:       "benchmark",
        Name:       "mymod.benchmark.parent",
        ShortName:  "parent",
        Title:      "Parent Benchmark",
        IsTopLevel: true,
        ChildNames: []string{"mymod.benchmark.child1", "mymod.benchmark.child2"},
    }
    idx.Add(parent)

    // Add children
    idx.Add(&IndexEntry{
        Type:       "benchmark",
        Name:       "mymod.benchmark.child1",
        ShortName:  "child1",
        Title:      "Child 1",
        ParentName: "mymod.benchmark.parent",
    })
    idx.Add(&IndexEntry{
        Type:       "benchmark",
        Name:       "mymod.benchmark.child2",
        ShortName:  "child2",
        Title:      "Child 2",
        ParentName: "mymod.benchmark.parent",
    })

    // Get children
    children := idx.GetChildren("mymod.benchmark.parent")
    assert.Len(t, children, 2)

    // Top level benchmarks
    topLevel := idx.TopLevelBenchmarks()
    assert.Len(t, topLevel, 1)
    assert.Equal(t, "mymod.benchmark.parent", topLevel[0].Name)
}

func TestIndex_Size(t *testing.T) {
    idx := NewResourceIndex()

    // Add many entries
    for i := 0; i < 1000; i++ {
        idx.Add(&IndexEntry{
            Type:      "query",
            Name:      fmt.Sprintf("mymod.query.query_%d", i),
            ShortName: fmt.Sprintf("query_%d", i),
            Title:     fmt.Sprintf("Query %d", i),
            FileName:  "/path/to/queries.pp",
            StartLine: i * 10,
            EndLine:   i*10 + 9,
        })
    }

    // Index should be small
    size := idx.Size()
    t.Logf("Index size for 1000 entries: %d bytes (%.2f KB)", size, float64(size)/1024)

    // Should be less than 500KB for 1000 entries
    assert.Less(t, size, 500*1024, "Index too large")
}

func TestIndex_AvailableDashboardsPayload(t *testing.T) {
    idx := NewResourceIndex()

    // Add dashboards
    idx.Add(&IndexEntry{
        Type:      "dashboard",
        Name:      "mymod.dashboard.main",
        ShortName: "main",
        Title:     "Main Dashboard",
        Tags:      map[string]string{"service": "aws"},
    })

    // Add benchmarks with hierarchy
    idx.Add(&IndexEntry{
        Type:       "benchmark",
        Name:       "mymod.benchmark.cis",
        ShortName:  "cis",
        Title:      "CIS Benchmark",
        IsTopLevel: true,
        ChildNames: []string{"mymod.benchmark.cis_1"},
    })
    idx.Add(&IndexEntry{
        Type:       "benchmark",
        Name:       "mymod.benchmark.cis_1",
        ShortName:  "cis_1",
        Title:      "CIS 1.x",
        ParentName: "mymod.benchmark.cis",
    })

    payload := idx.BuildAvailableDashboardsPayload()

    assert.Len(t, payload.Dashboards, 1)
    assert.Equal(t, "Main Dashboard", payload.Dashboards["mymod.dashboard.main"].Title)

    assert.Len(t, payload.Benchmarks, 2)
    assert.True(t, payload.Benchmarks["mymod.benchmark.cis"].IsTopLevel)
    assert.Len(t, payload.Benchmarks["mymod.benchmark.cis"].Children, 1)
}
```

## Acceptance Criteria

- [ ] IndexEntry struct captures all metadata needed for UI lists
- [ ] ResourceIndex provides O(1) lookup by name
- [ ] Index supports filtering by type (Dashboards(), Benchmarks(), etc.)
- [ ] Index supports hierarchy navigation (GetChildren, TopLevelBenchmarks)
- [ ] Index can build available_dashboards payload without full resources
- [ ] Index memory usage < 1 KB per 100 resources
- [ ] All unit tests pass
- [ ] Index is thread-safe (concurrent read access)

## Notes

- Keep index entries immutable after creation
- Consider adding index serialization for persistence (future optimization)
- String interning could further reduce memory (Task 14)
- ByteOffset/ByteLength enable efficient file seeking in Task 7
