# Fundamental Memory Analysis: Why Powerpipe Uses So Much Memory

**Date**: 2025-12-28
**Status**: Deep investigation into root causes

## The Real Problem

After deep analysis, I found the **fundamental architectural issue**: Powerpipe keeps the entire HCL AST and source code in memory after parsing is complete.

### Evidence

**1. Every resource keeps unparsed HCL body** (25 resource types):
```go
// Found in ALL resource types:
Remain hcl.Body `hcl:",remain" json:"-"`
```

Files with `Remain hcl.Body`:
- dashboard.go, dashboard_card.go, dashboard_chart.go, dashboard_container.go
- dashboard_edge.go, dashboard_flow.go, dashboard_graph.go, dashboard_hierarchy.go
- dashboard_image.go, dashboard_input.go, dashboard_node.go, dashboard_table.go
- dashboard_text.go, dashboard_with.go, dashboard_category.go
- control.go, control_benchmark.go, query.go, detection.go, detection_benchmark.go
- query_provider_impl.go, with_provider_impl.go, node_edge_provider_impl.go
- runtime_dependency_provider_impl.go, dashboard_leaf_node_impl.go

**2. Full source code stored in every resource**:
```go
// modconfig/resource_metadata.go:18
SourceDefinition string `json:"source_definition"`
```

**3. All file contents kept in memory**:
```go
// parse/mod_parse_context.go
FileData map[string][]byte  // Entire file contents!
```

**4. cty.Value objects never discarded**:
```go
// Variables keep multiple cty.Value fields
Default cty.Value
Type    cty.Type
Enum    cty.Value
Value   cty.Value
```

## Memory Breakdown for Large Mod (200 dashboards, 400 queries, 500 controls)

| Data Type | Count | Estimated Size | Total |
|-----------|-------|----------------|-------|
| `hcl.Body` (AST) per resource | 1,100+ | 50-200 KB each | 55-220 MB |
| `SourceDefinition` strings | 1,100+ | 1-10 KB each | 1-11 MB |
| `FileData` (raw files) | ~100 files | 3-5 KB each | 0.3-0.5 MB |
| cty.Value objects | thousands | varies | 50-100 MB |
| Maps, slices, pointers | many | varies | 50-100 MB |
| **Total** | | | **~150-430 MB** |

The `hcl.Body` AST is the biggest culprit. Each one contains:
- Token stream from lexer
- AST nodes (Blocks, Attributes, Expressions)
- Source position metadata (file, line, column, byte offset)
- References to parent bodies

---

## Fundamental Solutions

### Option 1: Post-Parse Cleanup (Quick Win)

**Concept**: After parsing completes, nil out data that's no longer needed.

```go
// After mod loading is complete:
func (m *PowerpipeModResources) CompactMemory() {
    m.WalkResources(func(r modconfig.HclResource) (bool, error) {
        // Clear hcl.Body - no longer needed after parsing
        if rc, ok := r.(interface{ ClearRemain() }); ok {
            rc.ClearRemain()
        }
        // Clear source definition if not needed
        if rm, ok := r.(modconfig.ResourceWithMetadata); ok {
            rm.GetMetadata().SourceDefinition = ""
        }
        return true, nil
    })
}

// Add to each resource type:
func (q *Query) ClearRemain() {
    q.Remain = nil
}
```

**Impact**: Could reduce memory by 50-70%
**Effort**: Low
**Risk**: Low (data already extracted during parsing)
**Trade-off**: Lose ability to re-decode or report detailed HCL errors after load

---

### Option 2: Lazy Loading Architecture (Medium-Term)

**Concept**: Don't parse everything upfront. Build lightweight index, parse on-demand.

```
Current Flow:
  startup → parse ALL files → decode ALL resources → keep ALL in memory

New Flow:
  startup → scan files for resource headers → build index → done
  on access → parse specific resource → cache with LRU → return
```

**Implementation**:

```go
// Lightweight index (always in memory)
type ResourceIndex struct {
    mu      sync.RWMutex
    entries map[string]*ResourceIndexEntry
}

type ResourceIndexEntry struct {
    Type       string  // "dashboard", "query", etc.
    Name       string
    Title      string  // For display
    Tags       map[string]string
    FilePath   string
    ByteOffset int64   // Where in file
    ByteLength int     // How many bytes
}

// LRU cache of parsed resources
type ResourceCache struct {
    mu       sync.Mutex
    cache    *lru.Cache  // github.com/hashicorp/golang-lru
    maxSize  int
}

func (c *ResourceCache) Get(name string) (modconfig.HclResource, error) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if cached, ok := c.cache.Get(name); ok {
        return cached.(modconfig.HclResource), nil
    }

    // Parse on demand
    entry := c.index.entries[name]
    resource, err := c.parseResource(entry)
    if err != nil {
        return nil, err
    }

    c.cache.Add(name, resource)
    return resource, nil
}
```

**Memory Model**:
```
Index size: O(number of resources) × ~200 bytes = ~220 KB for 1100 resources
Cache size: O(cache limit) × ~50-200 KB = configurable, e.g., 50 MB max
Total: ~50 MB max instead of ~400 MB
```

**Impact**: 80-90% memory reduction, bounded by cache size
**Effort**: High (significant refactoring)
**Risk**: Medium (changes access patterns)
**Trade-off**: First access to resource is slower

---

### Option 3: Binary Cache Format (Medium-Term)

**Concept**: Parse HCL once, serialize to compact binary. Load binary on subsequent runs.

```
First load:
  HCL files → parse → build resources → serialize to .ppcache

Subsequent loads:
  .ppcache → deserialize → resources (skip HCL entirely)
```

**Implementation**:

```go
// Use efficient serialization (gob, msgpack, or protobuf)
type ModCache struct {
    Version      string
    ModHash      string  // Hash of all source files
    Resources    []CachedResource
}

type CachedResource struct {
    Type        string
    Name        string
    // Only essential fields - no hcl.Body, no SourceDefinition
    Title       string
    Description string
    SQL         string  // For queries
    Args        []CachedArg
    Children    []string  // References by name
    // etc.
}

func LoadModWithCache(modPath string) (*Mod, error) {
    cachePath := filepath.Join(modPath, ".powerpipe", "mod.cache")

    // Check if cache is valid
    if cacheValid(modPath, cachePath) {
        return loadFromCache(cachePath)
    }

    // Parse and cache
    mod, err := parseModFromHCL(modPath)
    if err != nil {
        return nil, err
    }

    saveToCache(mod, cachePath)
    return mod, nil
}

func cacheValid(modPath, cachePath string) bool {
    // Compare mod file hashes with cached hashes
    // Return false if any HCL file changed
}
```

**Impact**:
- First load: Same as current
- Subsequent loads: 10-100x faster, 50-70% less memory
**Effort**: Medium
**Risk**: Low (fallback to HCL parsing if cache invalid)
**Trade-off**: Cache invalidation complexity, disk space

---

### Option 4: Arena Allocation (Go 1.20+)

**Concept**: Allocate parsing structures in arena, copy only essentials to heap, free arena.

```go
import "arena"

func parseModWithArena(modPath string) (*Mod, error) {
    // Create arena for parsing phase
    a := arena.NewArena()
    defer a.Free()  // Frees ALL arena allocations at once

    // Parse in arena (requires arena-aware HCL parser)
    // This is complex because we'd need to modify hcl library

    // Copy only essential data to heap
    mod := copyEssentialsToHeap(parsedMod)

    return mod, nil
}
```

**Impact**: Could reduce peak memory and GC pressure significantly
**Effort**: Very High (requires HCL library modifications)
**Risk**: High (experimental Go feature)
**Trade-off**: Complexity, Go version requirement

---

### Option 5: Streaming/Incremental Parser

**Concept**: Custom parser that extracts only needed data without full AST.

```go
// Custom scanner that extracts resource metadata without full parse
type ResourceScanner struct {
    file    *os.File
    lexer   *hcl.Lexer
}

func (s *ResourceScanner) ScanResources() ([]ResourceHeader, error) {
    var headers []ResourceHeader

    for {
        // Look for block starts: "dashboard", "query", "control", etc.
        blockType, name, startPos := s.scanBlockHeader()
        if blockType == "" {
            break
        }

        // Extract just title and tags without full decode
        title, tags := s.scanBasicAttributes()

        // Skip to end of block
        endPos := s.skipToBlockEnd()

        headers = append(headers, ResourceHeader{
            Type:       blockType,
            Name:       name,
            Title:      title,
            Tags:       tags,
            StartPos:   startPos,
            EndPos:     endPos,
        })
    }

    return headers, nil
}
```

**Impact**: Initial load could be 10x faster and use 90% less memory
**Effort**: Very High (custom parser)
**Risk**: Medium (must handle HCL syntax correctly)
**Trade-off**: Duplicated parsing logic

---

### Option 6: String Interning + Deduplication

**Concept**: Many strings are duplicated. Intern them.

```go
var stringPool = sync.Map{}

func intern(s string) string {
    if v, ok := stringPool.Load(s); ok {
        return v.(string)
    }
    stringPool.Store(s, s)
    return s
}

// Use during resource creation:
resource.Type = intern(resource.Type)  // "dashboard" shared by all dashboards
resource.ModName = intern(resource.ModName)  // Shared by all resources in mod
```

Common strings that get duplicated:
- Resource types: "dashboard", "query", "control", "benchmark"
- Tag keys: "service", "category", "type"
- Mod names: repeated in every resource
- SQL keywords: "SELECT", "FROM", "WHERE"

**Impact**: 10-20% memory reduction
**Effort**: Low
**Risk**: Low
**Trade-off**: Minor (unbounded pool could grow, but strings are usually limited)

---

## Recommended Approach: Phased Implementation

### Phase 1: Quick Wins (1-2 days)
1. **Post-parse cleanup** - Nil out `Remain` and `SourceDefinition`
2. **String interning** - Deduplicate common strings
3. **Map capacity hints** - Pre-size maps

Expected impact: **40-50% memory reduction**

### Phase 2: Caching (1 week)
4. **Binary cache format** - Skip HCL parsing on repeat loads
5. **Available dashboards payload caching** - Cache WebSocket responses

Expected impact: **Additional 20-30% on repeat loads**

### Phase 3: Architecture (2-4 weeks)
6. **Lazy loading with LRU** - Parse on-demand, bounded cache

Expected impact: **80-90% reduction, bounded memory regardless of mod size**

---

## Quick Test: Post-Parse Cleanup

To validate the impact, try this experiment:

```go
// In workspace loading, after mod is fully loaded:
func cleanupAfterLoad(mod *modconfig.Mod) {
    resources := GetModResources(mod)
    resources.WalkResources(func(r modconfig.HclResource) (bool, error) {
        // Use reflection to nil out Remain fields
        v := reflect.ValueOf(r).Elem()
        for i := 0; i < v.NumField(); i++ {
            field := v.Field(i)
            if field.Type() == reflect.TypeOf((*hcl.Body)(nil)).Elem() {
                if field.CanSet() {
                    field.Set(reflect.Zero(field.Type()))
                }
            }
        }
        return true, nil
    })

    // Force GC to see impact
    runtime.GC()
}
```

Add memory profiling before/after to measure impact.

---

## Conclusion

The root cause of high memory usage is **keeping HCL AST data after parsing**. The `hcl.Body` "remain" pattern is convenient for partial decoding but extremely expensive for large mods.

**The minimum viable fix** is to nil out `Remain` fields and clear `SourceDefinition` after parsing completes. This should reduce memory by 40-50% with minimal code changes.

**The ideal long-term solution** is lazy loading with an LRU cache, which would make memory usage bounded regardless of mod size.
