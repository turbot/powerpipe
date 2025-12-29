# Task 6: LRU Cache Implementation

## Objective

Implement a thread-safe, bounded LRU (Least Recently Used) cache for parsed resources. This enables lazy loading while keeping memory bounded regardless of mod size.

## Context

- Cache is the core memory management mechanism for lazy loading
- Must be thread-safe for concurrent dashboard executions
- Must track memory size, not just entry count
- Eviction should happen automatically when memory threshold exceeded
- Cache should work with any resource type

## Dependencies

### Prerequisites
- Task 4 (Resource Index) - Need to understand what metadata is cached vs full resource

### Files to Create
- `internal/resourcecache/cache.go` - LRU cache implementation
- `internal/resourcecache/cache_test.go` - Cache tests
- `internal/resourcecache/metrics.go` - Cache hit/miss metrics

### Files to Modify
- None (standalone module)

## Implementation Details

### 1. Cache Interface

```go
// internal/resourcecache/cache.go
package resourcecache

import (
    "container/list"
    "sync"
    "time"
)

// CacheConfig configures the LRU cache
type CacheConfig struct {
    MaxMemoryBytes int64         // Maximum memory usage (default: 50MB)
    MaxEntries     int           // Maximum entries (0 = unlimited, use memory only)
    TTL            time.Duration // Time-to-live for entries (0 = no expiry)
}

func DefaultConfig() CacheConfig {
    return CacheConfig{
        MaxMemoryBytes: 50 * 1024 * 1024, // 50MB
        MaxEntries:     0,                 // Memory-based eviction
        TTL:            0,                 // No expiry
    }
}

// Sizer interface for items that can report their size
type Sizer interface {
    Size() int64
}

// Cache is a thread-safe LRU cache with memory-based eviction
type Cache struct {
    mu sync.RWMutex

    // LRU list - front is most recently used
    list *list.List
    // Map from key to list element
    items map[string]*list.Element

    config CacheConfig

    // Current stats
    currentMemory int64
    hits          int64
    misses        int64
    evictions     int64
}

// entry holds a cached item
type entry struct {
    key       string
    value     interface{}
    size      int64
    timestamp time.Time
}

// New creates a new LRU cache
func New(config CacheConfig) *Cache {
    return &Cache{
        list:   list.New(),
        items:  make(map[string]*list.Element),
        config: config,
    }
}

// Get retrieves an item from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if elem, ok := c.items[key]; ok {
        e := elem.Value.(*entry)

        // Check TTL if configured
        if c.config.TTL > 0 && time.Since(e.timestamp) > c.config.TTL {
            c.removeElement(elem)
            c.misses++
            return nil, false
        }

        // Move to front (most recently used)
        c.list.MoveToFront(elem)
        c.hits++
        return e.value, true
    }

    c.misses++
    return nil, false
}

// Put adds an item to the cache
func (c *Cache) Put(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()

    // Calculate size
    size := int64(0)
    if sizer, ok := value.(Sizer); ok {
        size = sizer.Size()
    }

    // Update existing entry
    if elem, ok := c.items[key]; ok {
        e := elem.Value.(*entry)
        c.currentMemory -= e.size
        e.value = value
        e.size = size
        e.timestamp = time.Now()
        c.currentMemory += size
        c.list.MoveToFront(elem)
        c.evictIfNeeded()
        return
    }

    // Add new entry
    e := &entry{
        key:       key,
        value:     value,
        size:      size,
        timestamp: time.Now(),
    }
    elem := c.list.PushFront(e)
    c.items[key] = elem
    c.currentMemory += size

    c.evictIfNeeded()
}

// Delete removes an item from the cache
func (c *Cache) Delete(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if elem, ok := c.items[key]; ok {
        c.removeElement(elem)
    }
}

// Clear removes all items from the cache
func (c *Cache) Clear() {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.list.Init()
    c.items = make(map[string]*list.Element)
    c.currentMemory = 0
}

// evictIfNeeded removes least recently used items until under memory limit
func (c *Cache) evictIfNeeded() {
    // Check entry limit
    for c.config.MaxEntries > 0 && c.list.Len() > c.config.MaxEntries {
        c.evictOldest()
    }

    // Check memory limit
    for c.currentMemory > c.config.MaxMemoryBytes && c.list.Len() > 0 {
        c.evictOldest()
    }
}

func (c *Cache) evictOldest() {
    elem := c.list.Back()
    if elem != nil {
        c.removeElement(elem)
        c.evictions++
    }
}

func (c *Cache) removeElement(elem *list.Element) {
    e := elem.Value.(*entry)
    c.list.Remove(elem)
    delete(c.items, e.key)
    c.currentMemory -= e.size
}

// Stats returns cache statistics
func (c *Cache) Stats() CacheStats {
    c.mu.RLock()
    defer c.mu.RUnlock()

    return CacheStats{
        Entries:       c.list.Len(),
        MemoryBytes:   c.currentMemory,
        MaxMemory:     c.config.MaxMemoryBytes,
        Hits:          c.hits,
        Misses:        c.misses,
        Evictions:     c.evictions,
        HitRate:       c.hitRate(),
    }
}

func (c *Cache) hitRate() float64 {
    total := c.hits + c.misses
    if total == 0 {
        return 0
    }
    return float64(c.hits) / float64(total)
}

type CacheStats struct {
    Entries     int
    MemoryBytes int64
    MaxMemory   int64
    Hits        int64
    Misses      int64
    Evictions   int64
    HitRate     float64
}
```

### 2. Resource Cache Wrapper

```go
// internal/resourcecache/resource_cache.go
package resourcecache

import (
    "github.com/turbot/pipe-fittings/v2/modconfig"
)

// ResourceCache specializes Cache for HCL resources
type ResourceCache struct {
    cache *Cache
}

// NewResourceCache creates a cache for parsed resources
func NewResourceCache(config CacheConfig) *ResourceCache {
    return &ResourceCache{
        cache: New(config),
    }
}

// GetResource retrieves a parsed resource by full name
func (rc *ResourceCache) GetResource(name string) (modconfig.HclResource, bool) {
    val, ok := rc.cache.Get(name)
    if !ok {
        return nil, false
    }
    return val.(modconfig.HclResource), true
}

// PutResource caches a parsed resource
func (rc *ResourceCache) PutResource(name string, resource modconfig.HclResource) {
    rc.cache.Put(name, resource)
}

// GetDashboard retrieves a dashboard by name
func (rc *ResourceCache) GetDashboard(name string) (*modconfig.Dashboard, bool) {
    val, ok := rc.cache.Get(name)
    if !ok {
        return nil, false
    }
    if dash, ok := val.(*modconfig.Dashboard); ok {
        return dash, true
    }
    return nil, false
}

// PutDashboard caches a dashboard
func (rc *ResourceCache) PutDashboard(name string, dashboard *modconfig.Dashboard) {
    rc.cache.Put(name, dashboard)
}

// Similar methods for other resource types...

// Stats returns cache statistics
func (rc *ResourceCache) Stats() CacheStats {
    return rc.cache.Stats()
}

// Clear clears the cache
func (rc *ResourceCache) Clear() {
    rc.cache.Clear()
}

// Invalidate removes a specific resource
func (rc *ResourceCache) Invalidate(name string) {
    rc.cache.Delete(name)
}

// InvalidateAll removes all resources matching a predicate
func (rc *ResourceCache) InvalidateAll(predicate func(string) bool) {
    rc.cache.mu.Lock()
    defer rc.cache.mu.Unlock()

    var toDelete []string
    for key := range rc.cache.items {
        if predicate(key) {
            toDelete = append(toDelete, key)
        }
    }

    for _, key := range toDelete {
        if elem, ok := rc.cache.items[key]; ok {
            rc.cache.removeElement(elem)
        }
    }
}
```

### 3. Cache Metrics

```go
// internal/resourcecache/metrics.go
package resourcecache

import (
    "encoding/json"
    "fmt"
    "strings"
    "time"
)

// MetricsCollector tracks cache performance over time
type MetricsCollector struct {
    cache    *Cache
    interval time.Duration
    samples  []CacheStats
    maxSamples int
}

func NewMetricsCollector(cache *Cache, interval time.Duration) *MetricsCollector {
    return &MetricsCollector{
        cache:      cache,
        interval:   interval,
        maxSamples: 1000,
    }
}

func (m *MetricsCollector) Collect() {
    stats := m.cache.Stats()
    m.samples = append(m.samples, stats)

    // Keep bounded
    if len(m.samples) > m.maxSamples {
        m.samples = m.samples[len(m.samples)-m.maxSamples:]
    }
}

func (m *MetricsCollector) Report() string {
    stats := m.cache.Stats()

    var b strings.Builder
    b.WriteString("=== Resource Cache Metrics ===\n")
    b.WriteString(fmt.Sprintf("Entries:      %d\n", stats.Entries))
    b.WriteString(fmt.Sprintf("Memory:       %s / %s (%.1f%%)\n",
        formatBytes(stats.MemoryBytes),
        formatBytes(stats.MaxMemory),
        float64(stats.MemoryBytes)/float64(stats.MaxMemory)*100))
    b.WriteString(fmt.Sprintf("Hit Rate:     %.1f%%\n", stats.HitRate*100))
    b.WriteString(fmt.Sprintf("Hits:         %d\n", stats.Hits))
    b.WriteString(fmt.Sprintf("Misses:       %d\n", stats.Misses))
    b.WriteString(fmt.Sprintf("Evictions:    %d\n", stats.Evictions))

    return b.String()
}

func (m *MetricsCollector) JSON() ([]byte, error) {
    return json.Marshal(m.cache.Stats())
}

func formatBytes(b int64) string {
    const unit = 1024
    if b < unit {
        return fmt.Sprintf("%d B", b)
    }
    div, exp := int64(unit), 0
    for n := b / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
```

### 4. Cache Tests

```go
// internal/resourcecache/cache_test.go
package resourcecache

import (
    "fmt"
    "sync"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

type testItem struct {
    name string
    data []byte
}

func (t *testItem) Size() int64 {
    return int64(len(t.name) + len(t.data))
}

func TestCache_BasicOperations(t *testing.T) {
    cache := New(DefaultConfig())

    // Put and Get
    cache.Put("key1", &testItem{name: "item1", data: []byte("data1")})

    val, ok := cache.Get("key1")
    assert.True(t, ok)
    assert.Equal(t, "item1", val.(*testItem).name)

    // Miss
    _, ok = cache.Get("nonexistent")
    assert.False(t, ok)

    // Delete
    cache.Delete("key1")
    _, ok = cache.Get("key1")
    assert.False(t, ok)
}

func TestCache_LRUEviction(t *testing.T) {
    config := CacheConfig{
        MaxMemoryBytes: 100,
        MaxEntries:     3,
    }
    cache := New(config)

    // Add 3 items
    cache.Put("a", &testItem{name: "a", data: []byte("data")})
    cache.Put("b", &testItem{name: "b", data: []byte("data")})
    cache.Put("c", &testItem{name: "c", data: []byte("data")})

    // Access 'a' to make it recently used
    cache.Get("a")

    // Add 'd' - should evict 'b' (least recently used)
    cache.Put("d", &testItem{name: "d", data: []byte("data")})

    // 'b' should be gone
    _, ok := cache.Get("b")
    assert.False(t, ok, "b should have been evicted")

    // 'a' should still be present
    _, ok = cache.Get("a")
    assert.True(t, ok, "a should still be present")
}

func TestCache_MemoryEviction(t *testing.T) {
    config := CacheConfig{
        MaxMemoryBytes: 100,
        MaxEntries:     0, // No entry limit
    }
    cache := New(config)

    // Add items until we exceed memory
    for i := 0; i < 10; i++ {
        cache.Put(fmt.Sprintf("key%d", i), &testItem{
            name: fmt.Sprintf("item%d", i),
            data: make([]byte, 20),
        })
    }

    stats := cache.Stats()
    assert.LessOrEqual(t, stats.MemoryBytes, int64(100),
        "Memory should be bounded")
    assert.Greater(t, stats.Evictions, int64(0),
        "Should have evicted items")
}

func TestCache_TTLExpiry(t *testing.T) {
    config := CacheConfig{
        MaxMemoryBytes: 1024 * 1024,
        TTL:            50 * time.Millisecond,
    }
    cache := New(config)

    cache.Put("key1", &testItem{name: "item1"})

    // Should exist immediately
    _, ok := cache.Get("key1")
    assert.True(t, ok)

    // Wait for expiry
    time.Sleep(100 * time.Millisecond)

    // Should be expired
    _, ok = cache.Get("key1")
    assert.False(t, ok)
}

func TestCache_ConcurrentAccess(t *testing.T) {
    cache := New(DefaultConfig())

    var wg sync.WaitGroup
    numGoroutines := 100
    numOperations := 1000

    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for j := 0; j < numOperations; j++ {
                key := fmt.Sprintf("key%d-%d", id, j%10)
                cache.Put(key, &testItem{
                    name: key,
                    data: []byte(fmt.Sprintf("data%d", j)),
                })
                cache.Get(key)
            }
        }(i)
    }

    wg.Wait()

    stats := cache.Stats()
    t.Logf("Final stats: entries=%d, memory=%d, hits=%d, misses=%d",
        stats.Entries, stats.MemoryBytes, stats.Hits, stats.Misses)
}

func TestCache_Update(t *testing.T) {
    cache := New(DefaultConfig())

    // Put initial value
    cache.Put("key1", &testItem{name: "v1", data: []byte("short")})
    stats1 := cache.Stats()

    // Update with larger value
    cache.Put("key1", &testItem{name: "v2", data: []byte("much longer data")})
    stats2 := cache.Stats()

    // Memory should increase
    assert.Greater(t, stats2.MemoryBytes, stats1.MemoryBytes)

    // Value should be updated
    val, _ := cache.Get("key1")
    assert.Equal(t, "v2", val.(*testItem).name)
}

func TestCache_Clear(t *testing.T) {
    cache := New(DefaultConfig())

    for i := 0; i < 100; i++ {
        cache.Put(fmt.Sprintf("key%d", i), &testItem{name: fmt.Sprintf("item%d", i)})
    }

    stats := cache.Stats()
    require.Greater(t, stats.Entries, 0)

    cache.Clear()

    stats = cache.Stats()
    assert.Equal(t, 0, stats.Entries)
    assert.Equal(t, int64(0), stats.MemoryBytes)
}

func TestResourceCache_ResourceTypes(t *testing.T) {
    rc := NewResourceCache(DefaultConfig())

    // Simulate caching different resource types
    rc.cache.Put("mod.dashboard.test", &testItem{name: "dashboard"})
    rc.cache.Put("mod.query.test", &testItem{name: "query"})
    rc.cache.Put("mod.control.test", &testItem{name: "control"})

    // Invalidate all queries
    rc.InvalidateAll(func(key string) bool {
        return strings.Contains(key, ".query.")
    })

    // Query should be gone
    _, ok := rc.cache.Get("mod.query.test")
    assert.False(t, ok)

    // Dashboard should remain
    _, ok = rc.cache.Get("mod.dashboard.test")
    assert.True(t, ok)
}

func BenchmarkCache_Get(b *testing.B) {
    cache := New(DefaultConfig())

    // Pre-populate cache
    for i := 0; i < 1000; i++ {
        cache.Put(fmt.Sprintf("key%d", i), &testItem{
            name: fmt.Sprintf("item%d", i),
            data: make([]byte, 100),
        })
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cache.Get(fmt.Sprintf("key%d", i%1000))
    }
}

func BenchmarkCache_Put(b *testing.B) {
    cache := New(DefaultConfig())

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cache.Put(fmt.Sprintf("key%d", i%1000), &testItem{
            name: fmt.Sprintf("item%d", i),
            data: make([]byte, 100),
        })
    }
}

func BenchmarkCache_ConcurrentMixed(b *testing.B) {
    cache := New(DefaultConfig())

    b.RunParallel(func(pb *testing.PB) {
        i := 0
        for pb.Next() {
            key := fmt.Sprintf("key%d", i%1000)
            if i%4 == 0 {
                cache.Put(key, &testItem{name: key})
            } else {
                cache.Get(key)
            }
            i++
        }
    })
}
```

## Acceptance Criteria

- [ ] LRU cache correctly evicts least recently used items
- [ ] Memory-based eviction keeps cache under configured limit
- [ ] TTL expiry works when configured
- [ ] Cache is thread-safe under concurrent access
- [ ] Hit/miss/eviction metrics are tracked accurately
- [ ] Cache supports updating existing entries
- [ ] Clear operation removes all entries and resets memory
- [ ] ResourceCache wrapper provides type-safe access
- [ ] All benchmarks show acceptable performance (< 1us per operation)
- [ ] All tests pass

## Notes

- Consider adding sharded cache for better concurrency if benchmarks show contention
- Size estimation for resources may need refinement based on actual memory profiling
- Consider adding batch operations (PutMany, GetMany) if needed
- Watch for GC pressure from evictions - may need object pooling
