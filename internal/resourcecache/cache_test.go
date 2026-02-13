package resourcecache

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testItem implements Sizer for testing
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

func TestCache_Keys(t *testing.T) {
	cache := New(DefaultConfig())

	cache.Put("key1", &testItem{name: "item1"})
	cache.Put("key2", &testItem{name: "item2"})
	cache.Put("key3", &testItem{name: "item3"})

	keys := cache.Keys()
	assert.Len(t, keys, 3)
	assert.Contains(t, keys, "key1")
	assert.Contains(t, keys, "key2")
	assert.Contains(t, keys, "key3")
}

func TestCache_Len(t *testing.T) {
	cache := New(DefaultConfig())

	assert.Equal(t, 0, cache.Len())

	cache.Put("key1", &testItem{name: "item1"})
	assert.Equal(t, 1, cache.Len())

	cache.Put("key2", &testItem{name: "item2"})
	assert.Equal(t, 2, cache.Len())

	cache.Delete("key1")
	assert.Equal(t, 1, cache.Len())
}

func TestCache_Stats(t *testing.T) {
	config := CacheConfig{
		MaxMemoryBytes: 1024,
		MaxEntries:     10,
	}
	cache := New(config)

	// Add some items
	cache.Put("key1", &testItem{name: "item1", data: []byte("data1")})
	cache.Put("key2", &testItem{name: "item2", data: []byte("data2")})

	// Generate hits
	cache.Get("key1")
	cache.Get("key2")

	// Generate miss
	cache.Get("nonexistent")

	stats := cache.Stats()
	assert.Equal(t, 2, stats.Entries)
	assert.Greater(t, stats.MemoryBytes, int64(0))
	assert.Equal(t, int64(1024), stats.MaxMemory)
	assert.Equal(t, int64(2), stats.Hits)
	assert.Equal(t, int64(1), stats.Misses)
	assert.InDelta(t, 0.666, stats.HitRate, 0.01)
}

func TestCache_NonSizerItems(t *testing.T) {
	cache := New(DefaultConfig())

	// Put item that doesn't implement Sizer
	cache.Put("key1", "simple string")

	val, ok := cache.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "simple string", val)

	// Size should be 0 for non-Sizer items
	stats := cache.Stats()
	assert.Equal(t, int64(0), stats.MemoryBytes)
}

func TestResourceCache_ResourceTypes(t *testing.T) {
	rc := NewResourceCache(DefaultConfig())

	// Simulate caching different resource types using the underlying cache
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

	// Control should remain
	_, ok = rc.cache.Get("mod.control.test")
	assert.True(t, ok)
}

func TestResourceCache_BasicOperations(t *testing.T) {
	rc := NewResourceCache(DefaultConfig())

	// Test with testItem (not HclResource, but tests the cache mechanics)
	rc.cache.Put("test.resource", &testItem{name: "test"})

	val, ok := rc.cache.Get("test.resource")
	assert.True(t, ok)
	assert.Equal(t, "test", val.(*testItem).name)

	// Stats
	stats := rc.Stats()
	assert.Equal(t, 1, stats.Entries)

	// Clear
	rc.Clear()
	stats = rc.Stats()
	assert.Equal(t, 0, stats.Entries)
}

func TestResourceCache_Invalidate(t *testing.T) {
	rc := NewResourceCache(DefaultConfig())

	rc.cache.Put("resource1", &testItem{name: "r1"})
	rc.cache.Put("resource2", &testItem{name: "r2"})

	// Invalidate one
	rc.Invalidate("resource1")

	_, ok := rc.cache.Get("resource1")
	assert.False(t, ok)

	_, ok = rc.cache.Get("resource2")
	assert.True(t, ok)
}

func TestMetricsCollector_Report(t *testing.T) {
	cache := New(DefaultConfig())
	mc := NewMetricsCollector(cache, time.Second)

	// Add some data
	cache.Put("key1", &testItem{name: "item1", data: make([]byte, 1000)})
	cache.Get("key1")
	cache.Get("nonexistent")

	report := mc.Report()
	assert.Contains(t, report, "Resource Cache Metrics")
	assert.Contains(t, report, "Entries:")
	assert.Contains(t, report, "Memory:")
	assert.Contains(t, report, "Hit Rate:")
	assert.Contains(t, report, "Hits:")
	assert.Contains(t, report, "Misses:")
}

func TestMetricsCollector_JSON(t *testing.T) {
	cache := New(DefaultConfig())
	mc := NewMetricsCollector(cache, time.Second)

	cache.Put("key1", &testItem{name: "item1"})

	jsonData, err := mc.JSON()
	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), `"Entries":1`)
}

func TestMetricsCollector_Collect(t *testing.T) {
	cache := New(DefaultConfig())
	mc := NewMetricsCollector(cache, time.Millisecond)

	// Collect some samples
	mc.Collect()
	cache.Put("key1", &testItem{name: "item1"})
	mc.Collect()
	cache.Put("key2", &testItem{name: "item2"})
	mc.Collect()

	samples := mc.Samples()
	assert.Len(t, samples, 3)
	assert.Equal(t, 0, samples[0].Entries)
	assert.Equal(t, 1, samples[1].Entries)
	assert.Equal(t, 2, samples[2].Entries)
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{100, "100 B"},
		{1023, "1023 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1024 * 1024, "1.0 MB"},
		{50 * 1024 * 1024, "50.0 MB"},
		{1024 * 1024 * 1024, "1.0 GB"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.bytes), func(t *testing.T) {
			result := formatBytes(tt.bytes)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Benchmarks

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

func BenchmarkCache_ConcurrentGet(b *testing.B) {
	cache := New(DefaultConfig())

	// Pre-populate cache
	for i := 0; i < 1000; i++ {
		cache.Put(fmt.Sprintf("key%d", i), &testItem{
			name: fmt.Sprintf("item%d", i),
			data: make([]byte, 100),
		})
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Get(fmt.Sprintf("key%d", i%1000))
			i++
		}
	})
}

func BenchmarkCache_ConcurrentPut(b *testing.B) {
	cache := New(DefaultConfig())

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Put(fmt.Sprintf("key%d", i%1000), &testItem{
				name: fmt.Sprintf("item%d", i),
				data: make([]byte, 100),
			})
			i++
		}
	})
}
