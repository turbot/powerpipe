# Additional Performance Optimization Opportunities

**Date**: 2025-12-28
**Phase**: Post-Task 9 Analysis
**Context**: After achieving 46% faster load time and 63% less memory for large mods

## Current State (After Previous Optimizations)

### Top Memory Allocators (Post-Optimization)
| Allocator | % of Memory | Priority |
|-----------|-------------|----------|
| cty.ObjectVal | 23.82% | HIGH |
| hclsyntax.emitToken | 15.42% | MEDIUM |
| cty.ObjectWithOptionalAttrs | 13.50% | HIGH |
| bufio.Scanner.Scan | 8.03% | LOW |
| gohcl.getFieldTags | 5.28% | MEDIUM |

**Combined CTY allocations: ~37%** - This is the new primary target.

---

## HIGH PRIORITY OPTIMIZATIONS

### 1. CTY Type Reflection Caching (pipe-fittings)

**Location**: `pipe-fittings/cty_helpers/attributes.go`

**Problem**: `GetCtyValue()` calls `GetCtyTypes()` which uses reflection to scan struct fields. This happens every time `CtyValue()` is called on a resource.

**Current Code**:
```go
func GetCtyValue(item interface{}) (cty.Value, error) {
    types, err := GetCtyTypes(item)  // Reflection every call!
    // ...
}
```

**Solution**: Cache cty types by `reflect.Type`:
```go
var ctyTypeCache sync.Map  // map[reflect.Type]map[string]cty.Type

func GetCtyTypes(item interface{}) (map[string]cty.Type, error) {
    t := reflect.TypeOf(item)
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
    }

    // Check cache first
    if cached, ok := ctyTypeCache.Load(t); ok {
        return cached.(map[string]cty.Type), nil
    }

    // Compute and cache
    types := computeCtyTypes(t)
    ctyTypeCache.Store(t, types)
    return types, nil
}
```

**Expected Impact**: 20-30% reduction in cty.ObjectVal allocations
**Effort**: Medium
**Risk**: Low

---

### 2. Connection Value Map Optimization (pipe-fittings)

**Location**: `pipe-fittings/parse/mod_parse_context.go:830-856`

**Problem**: `buildConnectionValueMap()` creates new `cty.ObjectVal` for each connection type on every iteration, using `AsValueMap()` to extract and rebuild.

**Current Code**:
```go
for _, conn := range m.PipelingConnections {
    typeMap := connectionMap[connType].AsValueMap()  // Extract
    typeMap[shortName] = ctyVal
    connectionMap[connType] = cty.ObjectVal(typeMap)  // Recreate!
}
```

**Solution**: Build as Go maps first, create ObjectVal once:
```go
// Build intermediate structure
typeMapBuilders := map[string]map[string]cty.Value{}
for _, conn := range m.PipelingConnections {
    connType := conn.GetConnectionType()
    if typeMapBuilders[connType] == nil {
        typeMapBuilders[connType] = make(map[string]cty.Value)
    }
    typeMapBuilders[connType][conn.GetShortName()] = ctyVal
}

// Convert to cty.ObjectVal once
for connType, typeMap := range typeMapBuilders {
    connectionMap[connType] = cty.ObjectVal(typeMap)
}
```

**Expected Impact**: Reduces ObjectVal allocations by number of connections
**Effort**: Low
**Risk**: Low

---

### 3. Credential Environment Map Caching (pipe-fittings)

**Location**: `pipe-fittings/credential/*.go` (67+ credential methods)

**Problem**: Each credential's `getEnv()` creates new `cty.StringVal` objects every call.

**Current Code**:
```go
func (c *AwsCredential) getEnv() map[string]cty.Value {
    env := map[string]cty.Value{}
    if c.AccessKey != nil {
        env["AWS_ACCESS_KEY_ID"] = cty.StringVal(*c.AccessKey)  // Every call!
    }
    // ...
}
```

**Solution**: Cache env map in credential struct:
```go
type AwsCredential struct {
    // ... existing fields ...
    cachedEnv     map[string]cty.Value
    cachedEnvOnce sync.Once
}

func (c *AwsCredential) getEnv() map[string]cty.Value {
    c.cachedEnvOnce.Do(func() {
        c.cachedEnv = c.buildEnv()
    })
    return c.cachedEnv
}
```

**Expected Impact**: Eliminates repeated cty.StringVal creation for credentials
**Effort**: Medium (67+ credential types)
**Risk**: Low (need to invalidate cache if credential changes)

---

### 4. gohcl Field Tag Caching (hashicorp/hcl or pipe-fittings)

**Location**: `github.com/hashicorp/hcl/v2/gohcl/schema.go:128-184`

**Problem**: `getFieldTags()` has NO caching and is called every `DecodeBody()`.

**Current Code** (in hashicorp/hcl):
```go
func getFieldTags(ty reflect.Type) *fieldTags {
    ct := ty.NumField()
    for i := 0; i < ct; i++ {
        field := ty.Field(i)
        tag := field.Tag.Get("hcl")  // Reflection every call!
    }
}
```

**Solution Options**:
1. **Fork hcl or submit PR** - Add caching to `getFieldTags()`
2. **Wrap in pipe-fittings** - Pre-cache schemas before decoding

**Expected Impact**: 50-80% reduction in getFieldTags allocations (5.28% → ~2%)
**Effort**: High (external dependency)
**Risk**: Medium

---

## MEDIUM PRIORITY OPTIMIZATIONS

### 5. PowerpipeModResources Map Pre-allocation

**Location**: `internal/resources/mod_resources.go:63-91`

**Problem**: 22+ maps created with `make(map[string]...)` without capacity hints.

**Current Code**:
```go
func emptyPowerpipeModResources() *PowerpipeModResources {
    return &PowerpipeModResources{
        Controls:          make(map[string]*Control),  // No capacity
        ControlBenchmarks: make(map[string]*Benchmark),
        // ... 20 more maps
    }
}
```

**Solution**: Accept expected sizes or use a builder pattern:
```go
func emptyPowerpipeModResources(opts ...ModResourcesOption) *PowerpipeModResources {
    cfg := defaultModResourcesConfig()
    for _, opt := range opts {
        opt(&cfg)
    }
    return &PowerpipeModResources{
        Controls:          make(map[string]*Control, cfg.expectedControls),
        ControlBenchmarks: make(map[string]*Benchmark, cfg.expectedBenchmarks),
        // ...
    }
}
```

**Expected Impact**: Reduces map growth allocations
**Effort**: Low
**Risk**: Low

---

### 6. Equals() Function Optimization

**Location**: `internal/resources/mod_resources.go:131-425`

**Problem**: ~300 line function doing bidirectional iteration for 22+ maps. O(n²) comparison.

**Solution**: Use content hashing for quick inequality checks:
```go
type PowerpipeModResources struct {
    // ... existing fields ...
    contentHash uint64  // Computed on modification
}

func (m *PowerpipeModResources) Equals(o modconfig.ModResources) bool {
    other := o.(*PowerpipeModResources)
    // Quick check first
    if m.contentHash != other.contentHash {
        return false
    }
    // Full comparison only if hashes match (unlikely for different content)
    return m.deepEquals(other)
}
```

**Expected Impact**: Fast inequality detection for workspace change detection
**Effort**: Medium
**Risk**: Low

---

### 7. Incremental HCL Parsing (pipe-fittings)

**Location**: `pipe-fittings/parse/parser.go`, `workspace/workspace.go`

**Problem**: Full workspace re-parse on any file change. No HCL AST caching.

**Current Flow**:
1. File changes → watcher triggers
2. `LoadWorkspaceMod()` called → ALL files re-parsed
3. ALL resources rebuilt from scratch

**Solution**: Implement file-level caching:
```go
type ParsedFileCache struct {
    mu       sync.RWMutex
    files    map[string]*CachedFile  // path → parsed result
}

type CachedFile struct {
    ModTime  time.Time
    Body     hcl.Body
    Diags    hcl.Diagnostics
}

func (c *ParsedFileCache) GetOrParse(path string, data []byte) (hcl.Body, hcl.Diagnostics) {
    c.mu.RLock()
    if cached, ok := c.files[path]; ok {
        if fileInfo.ModTime == cached.ModTime {
            c.mu.RUnlock()
            return cached.Body, cached.Diags
        }
    }
    c.mu.RUnlock()

    // Parse and cache
    body, diags := parseFile(path, data)
    c.mu.Lock()
    c.files[path] = &CachedFile{ModTime: fileInfo.ModTime, Body: body, Diags: diags}
    c.mu.Unlock()
    return body, diags
}
```

**Expected Impact**: 80-90% faster workspace reload when 1 file changes
**Effort**: High
**Risk**: Medium (cache invalidation complexity)

---

### 8. Execution Tree Object Pooling

**Location**: `internal/dashboardexecute/*.go`

**Problem**: Every dashboard execution creates new:
- `DashboardExecutionTree`
- `LeafRun` (per leaf node)
- `CheckRun`, `DashboardRun`, etc.
- Multiple maps per run

**Solution**: Use sync.Pool for frequently allocated objects:
```go
var leafRunPool = sync.Pool{
    New: func() interface{} {
        return &LeafRun{
            Properties: make(map[string]any, 10),
        }
    },
}

func NewLeafRun(resource resources.DashboardLeafNode, ...) (*LeafRun, error) {
    r := leafRunPool.Get().(*LeafRun)
    r.Reset()  // Clear previous state
    r.Resource = resource
    // ...
    return r, nil
}

func (r *LeafRun) Release() {
    r.Reset()
    leafRunPool.Put(r)
}
```

**Expected Impact**: Reduced GC pressure for repeated dashboard executions
**Effort**: Medium
**Risk**: Medium (must ensure proper reset)

---

### 9. Available Dashboards Payload Caching (Task 8)

**Location**: `internal/dashboardserver/payload.go`

**Problem**: `buildAvailableDashboardsPayload()` rebuilds the entire payload on every WebSocket request.

**Solution**: Cache the payload and invalidate on workspace changes:
```go
type Server struct {
    // ... existing fields ...
    cachedPayload     []byte
    cachedPayloadOnce sync.Once
    cachedPayloadMu   sync.RWMutex
}

func (s *Server) getAvailableDashboardsPayload() ([]byte, error) {
    s.cachedPayloadMu.RLock()
    if s.cachedPayload != nil {
        defer s.cachedPayloadMu.RUnlock()
        return s.cachedPayload, nil
    }
    s.cachedPayloadMu.RUnlock()

    // Build and cache
    s.cachedPayloadMu.Lock()
    defer s.cachedPayloadMu.Unlock()
    payload, err := buildAvailableDashboardsPayload(s.workspaceResources)
    if err != nil {
        return nil, err
    }
    s.cachedPayload = payload
    return payload, nil
}

// Invalidate on workspace change
func (s *Server) onWorkspaceChanged() {
    s.cachedPayloadMu.Lock()
    s.cachedPayload = nil
    s.cachedPayloadMu.Unlock()
}
```

**Expected Impact**: Instant response for dashboard list requests
**Effort**: Low
**Risk**: Low

---

### 10. Benchmark Trunk Slice Reuse

**Location**: `internal/dashboardserver/payload.go:151-174`

**Problem**: `addBenchmarkChildren()` allocates new slice for each recursion level:
```go
childTrunk := make([]string, len(trunk)+1)  // New allocation per level
copy(childTrunk, trunk)
```

**Solution**: Use a pre-allocated builder or slice pool:
```go
func addBenchmarkChildren(benchmark *resources.Benchmark, ...) []ModAvailableBenchmark {
    // Use shared builder with capacity
    builder := getTrunkBuilder(maxDepth)
    defer releaseTrunkBuilder(builder)
    return addBenchmarkChildrenWithBuilder(benchmark, builder, ...)
}
```

**Expected Impact**: Reduced allocations for deep benchmark trees
**Effort**: Low
**Risk**: Low

---

## LOW PRIORITY OPTIMIZATIONS

### 11. String Interning for Common Strings

**Problem**: 172+ occurrences of FullName/ShortName access, often creating temporary strings.

**Solution**: Intern common string values:
```go
var stringInternPool sync.Map

func InternString(s string) string {
    if v, ok := stringInternPool.Load(s); ok {
        return v.(string)
    }
    stringInternPool.Store(s, s)
    return s
}
```

**Expected Impact**: Reduced string allocations for repeated names
**Effort**: Low
**Risk**: Low (memory leak if unbounded)

---

### 12. JSON Marshaling Optimization

**Problem**: Every payload builds uses `json.Marshal()` which allocates.

**Solution**: Use json encoder with pooled buffers:
```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func marshalPayload(v interface{}) ([]byte, error) {
    buf := bufferPool.Get().(*bytes.Buffer)
    buf.Reset()
    defer bufferPool.Put(buf)

    enc := json.NewEncoder(buf)
    if err := enc.Encode(v); err != nil {
        return nil, err
    }
    // Return copy since buffer will be reused
    return append([]byte(nil), buf.Bytes()...), nil
}
```

**Expected Impact**: Reduced JSON marshaling allocations
**Effort**: Low
**Risk**: Low

---

## Summary by Location

### Powerpipe Changes
| Optimization | File | Priority | Effort |
|--------------|------|----------|--------|
| Map pre-allocation | mod_resources.go | Medium | Low |
| Equals() hashing | mod_resources.go | Medium | Medium |
| Execution tree pooling | dashboardexecute/*.go | Medium | Medium |
| Payload caching | payload.go | Medium | Low |
| Trunk slice reuse | payload.go | Low | Low |
| JSON buffer pooling | various | Low | Low |

### Pipe-fittings Changes
| Optimization | File | Priority | Effort |
|--------------|------|----------|--------|
| CTY type caching | cty_helpers/attributes.go | HIGH | Medium |
| Connection map optimization | mod_parse_context.go | HIGH | Low |
| Credential env caching | credential/*.go | HIGH | Medium |
| Incremental HCL parsing | parse/parser.go | Medium | High |

### External Dependencies
| Optimization | Dependency | Priority | Effort |
|--------------|------------|----------|--------|
| gohcl field tag caching | hashicorp/hcl | Medium | High |

---

## Recommended Implementation Order

1. **CTY type caching** (pipe-fittings) - Highest impact, medium effort
2. **Connection map optimization** (pipe-fittings) - High impact, low effort
3. **Credential env caching** (pipe-fittings) - High impact, medium effort
4. **Payload caching** (Powerpipe) - Quick win, low effort
5. **Map pre-allocation** (Powerpipe) - Low effort improvement
6. **Execution tree pooling** (Powerpipe) - For dashboard server workloads
7. **Incremental HCL parsing** (pipe-fittings) - For hot reload scenarios

---

## Expected Cumulative Impact

If all HIGH priority optimizations are implemented:

| Metric | Current | Expected | Improvement |
|--------|---------|----------|-------------|
| Large Mod Load | 239.85 ms | ~160-180 ms | 25-35% |
| Large Mod Memory | 413.70 MB | ~280-320 MB | 23-32% |
| Workspace Reload | Full re-parse | Incremental | 80-90% |
| Dashboard List Request | Rebuild | Cached | 95%+ |

Combined with previous optimizations (Task 5-7, getSourceDefinition fix):
- **Total time improvement from baseline**: ~65-70%
- **Total memory improvement from baseline**: ~75-80%
