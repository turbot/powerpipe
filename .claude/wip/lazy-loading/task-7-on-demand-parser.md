# Task 7: On-Demand Resource Parser

## Objective

Implement a parser that can load and parse individual resources on-demand using the index metadata (file path, byte offset, line numbers). This is the core mechanism that enables lazy loading.

## Context

- Resources are parsed one at a time when accessed
- Uses byte offset from index to seek directly to resource in file
- Parses only the HCL block for the requested resource
- Must handle resource dependencies (control → query reference)
- Must integrate with cache (check cache first, parse if miss)

## Repository

**This task is primarily Powerpipe.** Minimal or no pipe-fittings changes.

The loader and parser are Powerpipe-specific because:
1. They work with Powerpipe resource types (Dashboard, Control, etc.)
2. They integrate with Powerpipe's cache and index
3. They can reuse existing pipe-fittings decode functions without modification

## Dependencies

### Prerequisites
- Task 4 (Resource Index) - Need index entries with file locations
- Task 5 (File Scanner) - Need byte offsets populated
- Task 6 (LRU Cache) - Need cache to store parsed resources

### Files to Create (powerpipe)
- `internal/resourceloader/loader.go` - On-demand resource loader
- `internal/resourceloader/parser.go` - Single-resource HCL parser
- `internal/resourceloader/loader_test.go` - Loader tests

### Files to Modify
- None in pipe-fittings - reuse existing `parse.Decoder` infrastructure

## Implementation Details

### 1. Resource Loader Interface

```go
// internal/resourceloader/loader.go
package resourceloader

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/turbot/pipe-fittings/v2/modconfig"
    "github.com/turbot/powerpipe/internal/resourcecache"
    "github.com/turbot/powerpipe/internal/resourceindex"
)

// Loader provides on-demand loading of resources
type Loader struct {
    mu sync.RWMutex

    index   *resourceindex.ResourceIndex
    cache   *resourcecache.ResourceCache
    modPath string
    modName string

    // Statistics
    loadCount int64
    parseTime int64
}

// NewLoader creates a resource loader
func NewLoader(index *resourceindex.ResourceIndex, cache *resourcecache.ResourceCache,
    modPath, modName string) *Loader {
    return &Loader{
        index:   index,
        cache:   cache,
        modPath: modPath,
        modName: modName,
    }
}

// Load retrieves a resource by name, loading from disk if not cached
func (l *Loader) Load(ctx context.Context, name string) (modconfig.HclResource, error) {
    // Check cache first
    if resource, ok := l.cache.GetResource(name); ok {
        return resource, nil
    }

    // Load from disk
    return l.loadFromDisk(ctx, name)
}

// LoadDashboard loads a dashboard with all its children
func (l *Loader) LoadDashboard(ctx context.Context, name string) (*modconfig.Dashboard, error) {
    resource, err := l.Load(ctx, name)
    if err != nil {
        return nil, err
    }

    dash, ok := resource.(*modconfig.Dashboard)
    if !ok {
        return nil, fmt.Errorf("resource %s is not a dashboard", name)
    }

    // Load children recursively
    if err := l.loadChildren(ctx, dash); err != nil {
        return nil, err
    }

    return dash, nil
}

// LoadBenchmark loads a benchmark with all its children
func (l *Loader) LoadBenchmark(ctx context.Context, name string) (modconfig.ModTreeItem, error) {
    resource, err := l.Load(ctx, name)
    if err != nil {
        return nil, err
    }

    bench, ok := resource.(modconfig.ModTreeItem)
    if !ok {
        return nil, fmt.Errorf("resource %s is not a benchmark", name)
    }

    if err := l.loadBenchmarkChildren(ctx, bench); err != nil {
        return nil, err
    }

    return bench, nil
}

func (l *Loader) loadFromDisk(ctx context.Context, name string) (modconfig.HclResource, error) {
    entry, ok := l.index.Get(name)
    if !ok {
        return nil, fmt.Errorf("resource not found in index: %s", name)
    }

    start := time.Now()
    resource, err := l.parseResource(ctx, entry)
    if err != nil {
        return nil, fmt.Errorf("parsing resource %s: %w", name, err)
    }

    l.cache.PutResource(name, resource)

    l.mu.Lock()
    l.loadCount++
    l.parseTime += time.Since(start).Nanoseconds()
    l.mu.Unlock()

    return resource, nil
}

func (l *Loader) loadChildren(ctx context.Context, parent modconfig.ModTreeItem) error {
    for _, child := range parent.GetChildren() {
        if child == nil {
            continue
        }

        childName := child.Name()
        if _, ok := l.cache.GetResource(childName); !ok {
            if _, err := l.Load(ctx, childName); err != nil {
                continue // Child may be inline
            }
        }

        if treeItem, ok := child.(modconfig.ModTreeItem); ok {
            if err := l.loadChildren(ctx, treeItem); err != nil {
                return err
            }
        }
    }
    return nil
}

func (l *Loader) loadBenchmarkChildren(ctx context.Context, bench modconfig.ModTreeItem) error {
    entry, ok := l.index.Get(bench.Name())
    if !ok {
        return nil
    }

    for _, childName := range entry.ChildNames {
        child, err := l.Load(ctx, childName)
        if err != nil {
            return fmt.Errorf("loading child %s: %w", childName, err)
        }

        if childTree, ok := child.(modconfig.ModTreeItem); ok {
            if err := l.loadBenchmarkChildren(ctx, childTree); err != nil {
                return err
            }
        }

        if control, ok := child.(*modconfig.Control); ok {
            if err := l.loadControlDependencies(ctx, control); err != nil {
                return err
            }
        }
    }
    return nil
}

func (l *Loader) loadControlDependencies(ctx context.Context, control *modconfig.Control) error {
    if control.Query != nil && control.Query.FullName != "" {
        if _, err := l.Load(ctx, control.Query.FullName); err != nil {
            return fmt.Errorf("loading query %s: %w", control.Query.FullName, err)
        }
    }
    return nil
}

// Preload loads multiple resources in parallel
func (l *Loader) Preload(ctx context.Context, names []string) error {
    var wg sync.WaitGroup
    errChan := make(chan error, len(names))
    sem := make(chan struct{}, 10)

    for _, name := range names {
        if _, ok := l.cache.GetResource(name); ok {
            continue
        }

        wg.Add(1)
        go func(n string) {
            defer wg.Done()
            sem <- struct{}{}
            defer func() { <-sem }()

            if _, err := l.Load(ctx, n); err != nil {
                errChan <- err
            }
        }(name)
    }

    wg.Wait()
    close(errChan)

    for err := range errChan {
        return err
    }
    return nil
}

// Stats returns loader statistics
func (l *Loader) Stats() LoaderStats {
    l.mu.RLock()
    defer l.mu.RUnlock()

    avgTime := time.Duration(0)
    if l.loadCount > 0 {
        avgTime = time.Duration(l.parseTime / l.loadCount)
    }

    return LoaderStats{
        LoadCount:    l.loadCount,
        AvgParseTime: avgTime,
        CacheStats:   l.cache.Stats(),
    }
}

type LoaderStats struct {
    LoadCount    int64
    AvgParseTime time.Duration
    CacheStats   resourcecache.CacheStats
}
```

### 2. Single-Resource Parser

```go
// internal/resourceloader/parser.go
package resourceloader

import (
    "bufio"
    "context"
    "fmt"
    "io"
    "os"
    "strings"

    "github.com/hashicorp/hcl/v2"
    "github.com/hashicorp/hcl/v2/hclparse"
    "github.com/turbot/pipe-fittings/v2/modconfig"
    "github.com/turbot/powerpipe/internal/resourceindex"
)

// parseResource parses a single resource from its file
func (l *Loader) parseResource(ctx context.Context, entry *resourceindex.IndexEntry) (modconfig.HclResource, error) {
    blockContent, err := l.readResourceBlock(entry)
    if err != nil {
        return nil, fmt.Errorf("reading block: %w", err)
    }

    parser := hclparse.NewParser()
    file, diags := parser.ParseHCL(blockContent, entry.FileName)
    if diags.HasErrors() {
        return nil, fmt.Errorf("parsing HCL: %s", diags.Error())
    }

    return l.decodeResourceBlock(ctx, entry.Type, entry.ShortName, file.Body)
}

// readResourceBlock reads just the bytes for a single resource block
func (l *Loader) readResourceBlock(entry *resourceindex.IndexEntry) ([]byte, error) {
    file, err := os.Open(entry.FileName)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    // Use byte offset if available
    if entry.ByteOffset > 0 && entry.ByteLength > 0 {
        if _, err := file.Seek(entry.ByteOffset, io.SeekStart); err != nil {
            return nil, fmt.Errorf("seeking: %w", err)
        }
        content := make([]byte, entry.ByteLength)
        if _, err := io.ReadFull(file, content); err != nil {
            return nil, fmt.Errorf("reading: %w", err)
        }
        return content, nil
    }

    // Fallback to line-based reading
    return l.readByLines(entry)
}

func (l *Loader) readByLines(entry *resourceindex.IndexEntry) ([]byte, error) {
    file, err := os.Open(entry.FileName)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := bufio.NewReader(file)
    var content strings.Builder
    lineNum := 0

    for {
        line, err := reader.ReadString('\n')
        lineNum++

        if lineNum >= entry.StartLine && lineNum <= entry.EndLine {
            content.WriteString(line)
        }

        if lineNum >= entry.EndLine || err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }
    }

    return []byte(content.String()), nil
}

// decodeResourceBlock decodes HCL body into appropriate resource type
func (l *Loader) decodeResourceBlock(ctx context.Context, blockType, name string,
    body hcl.Body) (modconfig.HclResource, error) {

    switch blockType {
    case "dashboard":
        return l.decodeDashboard(name, body)
    case "query":
        return l.decodeQuery(name, body)
    case "control":
        return l.decodeControl(name, body)
    case "benchmark":
        return l.decodeBenchmark(name, body)
    case "card", "chart", "container", "flow", "graph", "hierarchy",
         "image", "input", "table", "text", "category", "node", "edge":
        return l.decodeDashboardComponent(blockType, name, body)
    case "detection", "detection_benchmark":
        return l.decodeDetection(blockType, name, body)
    case "variable":
        return l.decodeVariable(name, body)
    default:
        return nil, fmt.Errorf("unknown resource type: %s", blockType)
    }
}

// Decoder implementations follow the same pattern:
// 1. Create resource with NewXxx()
// 2. Decode attributes from body
// 3. Do NOT store Remain - this is key for memory savings
// 4. Return resource

func (l *Loader) decodeDashboard(name string, body hcl.Body) (*modconfig.Dashboard, error) {
    // Use existing parsing infrastructure where possible
    // Key difference: we don't keep the Remain field
    // Implementation details depend on modconfig internals
    return nil, fmt.Errorf("TODO: implement with modconfig integration")
}

func (l *Loader) decodeQuery(name string, body hcl.Body) (*modconfig.Query, error) {
    return nil, fmt.Errorf("TODO: implement with modconfig integration")
}

func (l *Loader) decodeControl(name string, body hcl.Body) (*modconfig.Control, error) {
    return nil, fmt.Errorf("TODO: implement with modconfig integration")
}

func (l *Loader) decodeBenchmark(name string, body hcl.Body) (modconfig.ModTreeItem, error) {
    return nil, fmt.Errorf("TODO: implement with modconfig integration")
}

func (l *Loader) decodeDashboardComponent(blockType, name string, body hcl.Body) (modconfig.HclResource, error) {
    return nil, fmt.Errorf("TODO: implement with modconfig integration")
}

func (l *Loader) decodeDetection(blockType, name string, body hcl.Body) (modconfig.HclResource, error) {
    return nil, fmt.Errorf("TODO: implement with modconfig integration")
}

func (l *Loader) decodeVariable(name string, body hcl.Body) (*modconfig.Variable, error) {
    return nil, fmt.Errorf("TODO: implement with modconfig integration")
}
```

### 3. Tests

```go
// internal/resourceloader/loader_test.go
package resourceloader

import (
    "context"
    "os"
    "path/filepath"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestLoader_LoadQuery(t *testing.T) {
    modPath := setupTestMod(t)
    loader := setupTestLoader(t, modPath)

    ctx := context.Background()
    resource, err := loader.Load(ctx, "testmod.query.simple")
    require.NoError(t, err)
    assert.NotNil(t, resource)
}

func TestLoader_CacheHit(t *testing.T) {
    modPath := setupTestMod(t)
    loader := setupTestLoader(t, modPath)

    ctx := context.Background()

    // First load - miss
    _, err := loader.Load(ctx, "testmod.query.simple")
    require.NoError(t, err)
    assert.Equal(t, int64(1), loader.cache.Stats().Misses)

    // Second load - hit
    _, err = loader.Load(ctx, "testmod.query.simple")
    require.NoError(t, err)
    assert.Equal(t, int64(1), loader.cache.Stats().Hits)
}

func TestLoader_NotFound(t *testing.T) {
    modPath := setupTestMod(t)
    loader := setupTestLoader(t, modPath)

    ctx := context.Background()
    _, err := loader.Load(ctx, "testmod.query.nonexistent")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "not found")
}

func TestLoader_Preload(t *testing.T) {
    modPath := setupTestMod(t)
    loader := setupTestLoader(t, modPath)

    names := []string{"testmod.query.q1", "testmod.query.q2", "testmod.query.q3"}
    err := loader.Preload(context.Background(), names)
    require.NoError(t, err)

    for _, name := range names {
        _, ok := loader.cache.GetResource(name)
        assert.True(t, ok)
    }
}

func setupTestMod(t testing.TB) string {
    tmpDir := t.TempDir()

    os.WriteFile(filepath.Join(tmpDir, "mod.pp"), []byte(`mod "testmod" {}`), 0644)
    os.WriteFile(filepath.Join(tmpDir, "queries.pp"), []byte(`
query "simple" { sql = "SELECT 1" }
query "q1" { sql = "SELECT 1" }
query "q2" { sql = "SELECT 2" }
query "q3" { sql = "SELECT 3" }
`), 0644)

    return tmpDir
}
```

## Acceptance Criteria

- [ ] Loader can load any resource type by name
- [ ] Uses byte offset for efficient file seeking when available
- [ ] Fallback to line-based reading works
- [ ] Loaded resources are cached automatically
- [ ] Cache hits return immediately without disk I/O
- [ ] LoadDashboard loads all children recursively
- [ ] LoadBenchmark loads all children and control queries
- [ ] Preload enables parallel loading
- [ ] Error handling for not found, parse errors
- [ ] Memory is not wasted (no Remain field storage)
- [ ] Performance: < 5ms per resource load from disk
- [ ] All tests pass

## Notes

- Decoder implementations need integration with existing modconfig constructors
- Consider reusing existing HCL decode infrastructure from pipe-fittings
- May need to handle inline resources (defined within parent block)
- Watch for race conditions in concurrent loading
