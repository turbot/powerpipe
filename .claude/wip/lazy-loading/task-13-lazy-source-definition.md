# Task 13: Lazy Source Definition Loading

## Objective

Replace eager `SourceDefinition` storage with lazy loading from file, using existing `ResourceMetadata` (FileName, StartLineNumber, EndLineNumber).

## Context

- Every resource stores its full HCL source in `SourceDefinition`
- This is used only for the UI "View Source" feature (panel detail)
- `ResourceMetadata` already has FileName + line numbers
- We can load source on-demand from file instead of storing it
- This can save significant memory for large mods

## Repository

**This task primarily changes pipe-fittings** since `ResourceMetadata` is defined there.

- **pipe-fittings**: Core change to `ResourceMetadata.SourceDefinition`
- **powerpipe**: Only needs to call the new lazy-loading method

The change is backward compatible - existing code that reads `SourceDefinition` directly will still work (it just won't benefit from lazy loading until migrated).

## Dependencies

### Prerequisites
- Task 12 (Post-Parse Cleanup) - Cleanup infrastructure

### Files to Modify (pipe-fittings)
- `modconfig/resource_metadata.go` - Add GetSourceDefinition() method, make field private

### Files to Create (pipe-fittings)
- `modconfig/source_loader.go` - Lazy source loading utility
- `modconfig/source_loader_test.go` - Tests

### Files to Modify (powerpipe)
- `internal/dashboardexecute/dashboard_tree_run_impl.go` - Use GetSourceDefinition() method
- Any other code that accesses SourceDefinition directly

## Implementation Details

### 1. Source Loader

```go
// pipe-fittings/modconfig/source_loader.go
package modconfig

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "strings"
    "sync"
)

// SourceLoader provides lazy loading of HCL source definitions
type SourceLoader struct {
    mu sync.RWMutex

    // Cache of loaded sources (optional)
    cache map[string]string

    // Whether to cache loaded sources
    enableCache bool

    // Maximum cache size
    maxCacheSize int
}

// DefaultSourceLoader is the default source loader instance
var DefaultSourceLoader = NewSourceLoader(false, 0)

// NewSourceLoader creates a source loader
func NewSourceLoader(enableCache bool, maxCacheSize int) *SourceLoader {
    return &SourceLoader{
        cache:        make(map[string]string),
        enableCache:  enableCache,
        maxCacheSize: maxCacheSize,
    }
}

// LoadSource loads source definition from file using metadata
func (sl *SourceLoader) LoadSource(meta *ResourceMetadata) (string, error) {
    if meta == nil || meta.FileName == "" {
        return "", nil
    }

    cacheKey := sl.cacheKey(meta)

    // Check cache
    if sl.enableCache {
        sl.mu.RLock()
        if cached, ok := sl.cache[cacheKey]; ok {
            sl.mu.RUnlock()
            return cached, nil
        }
        sl.mu.RUnlock()
    }

    // Load from file
    source, err := sl.loadFromFile(meta)
    if err != nil {
        return "", err
    }

    // Cache if enabled
    if sl.enableCache && sl.maxCacheSize > 0 {
        sl.mu.Lock()
        if len(sl.cache) < sl.maxCacheSize {
            sl.cache[cacheKey] = source
        }
        sl.mu.Unlock()
    }

    return source, nil
}

func (sl *SourceLoader) loadFromFile(meta *ResourceMetadata) (string, error) {
    file, err := os.Open(meta.FileName)
    if err != nil {
        return "", fmt.Errorf("opening source file: %w", err)
    }
    defer file.Close()

    reader := bufio.NewReader(file)
    var content strings.Builder
    lineNum := 0

    for {
        line, err := reader.ReadString('\n')
        lineNum++

        if lineNum >= meta.StartLineNumber && lineNum <= meta.EndLineNumber {
            content.WriteString(line)
        }

        if lineNum >= meta.EndLineNumber || err == io.EOF {
            break
        }

        if err != nil {
            return "", fmt.Errorf("reading source file: %w", err)
        }
    }

    return content.String(), nil
}

func (sl *SourceLoader) cacheKey(meta *ResourceMetadata) string {
    return fmt.Sprintf("%s:%d:%d", meta.FileName, meta.StartLineNumber, meta.EndLineNumber)
}

// ClearCache clears the source cache
func (sl *SourceLoader) ClearCache() {
    sl.mu.Lock()
    sl.cache = make(map[string]string)
    sl.mu.Unlock()
}

// GetSourceDefinition is a convenience function using the default loader
func GetSourceDefinition(meta *ResourceMetadata) string {
    source, err := DefaultSourceLoader.LoadSource(meta)
    if err != nil {
        return "" // Return empty on error
    }
    return source
}
```

### 2. ResourceMetadata Enhancement

```go
// pipe-fittings/modconfig/resource_metadata.go modifications

type ResourceMetadata struct {
    ResourceName      string `json:"resource_name"`
    FileName          string `json:"file_name"`
    StartLineNumber   int    `json:"start_line_number"`
    EndLineNumber     int    `json:"end_line_number"`
    Anonymous         bool   `json:"anonymous,omitempty"`

    // SourceDefinition is deprecated - use GetSourceDefinition()
    // Kept for backward compatibility during transition
    sourceDefinition string
}

// GetSourceDefinition returns the HCL source for this resource
// It loads lazily from file if not already cached
func (m *ResourceMetadata) GetSourceDefinition() string {
    // If we have a cached value, return it
    if m.sourceDefinition != "" {
        return m.sourceDefinition
    }

    // Otherwise load from file
    return GetSourceDefinition(m)
}

// SetSourceDefinition sets the source definition (for backward compatibility)
func (m *ResourceMetadata) SetSourceDefinition(source string) {
    m.sourceDefinition = source
}

// ClearSourceDefinition clears the cached source to free memory
func (m *ResourceMetadata) ClearSourceDefinition() {
    m.sourceDefinition = ""
}
```

### 3. Modify Resource Types

```go
// pipe-fittings/modconfig/dashboard.go modifications

type Dashboard struct {
    // ... existing fields

    // Remove eager SourceDefinition storage
    // ResourceMetadata contains file location for lazy loading
}

// GetSourceDefinition returns the HCL source (lazy loaded)
func (d *Dashboard) GetSourceDefinition() string {
    if d.ResourceMetadata == nil {
        return ""
    }
    return d.ResourceMetadata.GetSourceDefinition()
}

// Similar modifications for all resource types:
// - query.go
// - control.go
// - benchmark.go
// - card.go, chart.go, etc.
```

### 4. Clear SourceDefinition After Parse

```go
// pipe-fittings/parse/run_context.go additions

// ClearSourceDefinitions clears all cached source definitions
func (r *RunContext) ClearSourceDefinitions() {
    if r.Mod == nil {
        return
    }

    r.Mod.WalkResources(func(resource HclResource) (bool, error) {
        if meta := resource.GetMetadata(); meta != nil {
            meta.ClearSourceDefinition()
        }
        return true, nil
    })
}

// FinalizeParsing clears memory after parsing
func (r *RunContext) FinalizeParsing() {
    r.ClearAllRemainFields()
    r.ClearSourceDefinitions()
}
```

### 5. UI Integration (Panel Detail)

```go
// internal/dashboardexecute/dashboard_tree_run_impl.go modifications

// buildPanelDetailPayload builds the panel detail response
func buildPanelDetailPayload(resource HclResource) *PanelDetailPayload {
    payload := &PanelDetailPayload{
        Name: resource.Name(),
        Type: resource.GetBlockType(),
    }

    // Get source definition lazily
    if meta := resource.GetMetadata(); meta != nil {
        payload.SourceDefinition = meta.GetSourceDefinition()
    }

    return payload
}
```

### 6. Tests

```go
// pipe-fittings/modconfig/source_loader_test.go
package modconfig

import (
    "os"
    "path/filepath"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestSourceLoader_LoadSource(t *testing.T) {
    // Create test file
    tmpDir := t.TempDir()
    testFile := filepath.Join(tmpDir, "test.pp")

    content := `dashboard "test" {
    title = "Test Dashboard"
    description = "A test"
}

query "test_query" {
    sql = "SELECT 1"
}
`
    require.NoError(t, os.WriteFile(testFile, []byte(content), 0644))

    loader := NewSourceLoader(false, 0)

    // Load dashboard source (lines 1-4)
    meta := &ResourceMetadata{
        FileName:        testFile,
        StartLineNumber: 1,
        EndLineNumber:   4,
    }

    source, err := loader.LoadSource(meta)
    require.NoError(t, err)

    assert.Contains(t, source, "dashboard \"test\"")
    assert.Contains(t, source, "title = \"Test Dashboard\"")
    assert.NotContains(t, source, "query \"test_query\"")
}

func TestSourceLoader_WithCache(t *testing.T) {
    tmpDir := t.TempDir()
    testFile := filepath.Join(tmpDir, "test.pp")
    os.WriteFile(testFile, []byte("query \"test\" { sql = \"SELECT 1\" }"), 0644)

    loader := NewSourceLoader(true, 100)

    meta := &ResourceMetadata{
        FileName:        testFile,
        StartLineNumber: 1,
        EndLineNumber:   1,
    }

    // First load
    source1, err := loader.LoadSource(meta)
    require.NoError(t, err)

    // Second load (from cache)
    source2, err := loader.LoadSource(meta)
    require.NoError(t, err)

    assert.Equal(t, source1, source2)
}

func TestResourceMetadata_GetSourceDefinition(t *testing.T) {
    tmpDir := t.TempDir()
    testFile := filepath.Join(tmpDir, "test.pp")
    os.WriteFile(testFile, []byte("query \"test\" { sql = \"SELECT 1\" }"), 0644)

    meta := &ResourceMetadata{
        FileName:        testFile,
        StartLineNumber: 1,
        EndLineNumber:   1,
    }

    source := meta.GetSourceDefinition()
    assert.Contains(t, source, "query \"test\"")
}

func TestResourceMetadata_ClearSourceDefinition(t *testing.T) {
    meta := &ResourceMetadata{
        sourceDefinition: "cached source",
    }

    assert.Equal(t, "cached source", meta.GetSourceDefinition())

    meta.ClearSourceDefinition()

    // After clear, should return empty (no file to load from in this test)
    assert.Equal(t, "", meta.GetSourceDefinition())
}

func TestLazySourceLoading_MemoryReduction(t *testing.T) {
    modPath := setupLargeMod(t, 200)

    // Load with eager source
    runtime.GC()
    var beforeEager runtime.MemStats
    runtime.ReadMemStats(&beforeEager)

    wsEager, _ := LoadWorkspaceEagerSource(context.Background(), modPath)
    _ = wsEager

    runtime.GC()
    var afterEager runtime.MemStats
    runtime.ReadMemStats(&afterEager)

    eagerMem := afterEager.HeapAlloc - beforeEager.HeapAlloc

    // Load with lazy source
    runtime.GC()
    var beforeLazy runtime.MemStats
    runtime.ReadMemStats(&beforeLazy)

    wsLazy, _ := LoadWorkspaceLazySource(context.Background(), modPath)
    _ = wsLazy

    runtime.GC()
    var afterLazy runtime.MemStats
    runtime.ReadMemStats(&afterLazy)

    lazyMem := afterLazy.HeapAlloc - beforeLazy.HeapAlloc

    t.Logf("Eager source: %d bytes", eagerMem)
    t.Logf("Lazy source: %d bytes", lazyMem)

    assert.Less(t, lazyMem, eagerMem, "Lazy source should use less memory")
}
```

## Acceptance Criteria

- [ ] SourceLoader can load source from file using metadata
- [ ] ResourceMetadata.GetSourceDefinition() works lazily
- [ ] Source is loaded on-demand when UI requests panel detail
- [ ] Optional caching for frequently accessed sources
- [ ] ClearSourceDefinition() frees cached memory
- [ ] Parser clears source definitions after parsing
- [ ] UI "View Source" still works correctly
- [ ] Memory reduction measurable
- [ ] All behavior tests pass
- [ ] Graceful handling of missing files

## Notes

- File I/O is fast enough for occasional panel detail requests
- Consider caching if profiling shows repeated loads
- Watch for file changes during long-running sessions
- May need to handle file encoding (UTF-8 assumed)
- Error handling should be graceful (show empty source on error)
