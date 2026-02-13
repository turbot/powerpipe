package resourcecache

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// sizedItem implements Sizer for testing memory-related behavior
type sizedItem struct {
	key  string
	size int64
}

func (s *sizedItem) Size() int64 {
	return s.size
}

// negativeSizedItem returns a negative size (edge case)
type negativeSizedItem struct{}

func (n *negativeSizedItem) Size() int64 {
	return -100
}

// ============================================================================
// LRU Eviction Tests
// ============================================================================

func TestCache_LRUEvictionOrder(t *testing.T) {
	// Test: Least recently used item evicted first
	config := CacheConfig{
		MaxMemoryBytes: 1000,
		MaxEntries:     3,
	}
	cache := New(config)

	// Add A, B, C to full cache
	cache.Put("A", &sizedItem{key: "A", size: 10})
	cache.Put("B", &sizedItem{key: "B", size: 10})
	cache.Put("C", &sizedItem{key: "C", size: 10})

	// Access A to make it recently used
	_, _ = cache.Get("A")

	// Add D (triggers eviction) - B should be evicted as LRU
	cache.Put("D", &sizedItem{key: "D", size: 10})

	// B should be evicted (least recently used after A was accessed)
	_, ok := cache.Get("B")
	assert.False(t, ok, "B should have been evicted as LRU")

	// A should still be present (was accessed)
	_, ok = cache.Get("A")
	assert.True(t, ok, "A should still be present")

	// C should still be present
	_, ok = cache.Get("C")
	assert.True(t, ok, "C should still be present")

	// D should be present
	_, ok = cache.Get("D")
	assert.True(t, ok, "D should be present")
}

func TestCache_AccessUpdatesRecency(t *testing.T) {
	// Test: Access updates recency
	config := CacheConfig{
		MaxMemoryBytes: 1000,
		MaxEntries:     3,
	}
	cache := New(config)

	// Add A, B, C (A is oldest)
	cache.Put("A", &sizedItem{key: "A", size: 10})
	cache.Put("B", &sizedItem{key: "B", size: 10})
	cache.Put("C", &sizedItem{key: "C", size: 10})

	// Access A - moves to most recent
	_, _ = cache.Get("A")

	// Add D - should evict B (now LRU)
	cache.Put("D", &sizedItem{key: "D", size: 10})

	_, ok := cache.Get("B")
	assert.False(t, ok, "B should have been evicted")

	_, ok = cache.Get("A")
	assert.True(t, ok, "A should still be present after access")
}

func TestCache_UpdateCountsAsAccess(t *testing.T) {
	// Test: Update counts as access
	config := CacheConfig{
		MaxMemoryBytes: 1000,
		MaxEntries:     3,
	}
	cache := New(config)

	// Add A, B, C
	cache.Put("A", &sizedItem{key: "A", size: 10})
	cache.Put("B", &sizedItem{key: "B", size: 10})
	cache.Put("C", &sizedItem{key: "C", size: 10})

	// Update A (should update recency)
	cache.Put("A", &sizedItem{key: "A-updated", size: 15})

	// Add D - should evict B (LRU)
	cache.Put("D", &sizedItem{key: "D", size: 10})

	_, ok := cache.Get("B")
	assert.False(t, ok, "B should have been evicted")

	val, ok := cache.Get("A")
	assert.True(t, ok, "A should still be present after update")
	assert.Equal(t, "A-updated", val.(*sizedItem).key)
}

func TestCache_MultipleEvictions(t *testing.T) {
	// Test: Multiple items evicted to make room for large item
	config := CacheConfig{
		MaxMemoryBytes: 100,
		MaxEntries:     0, // Only memory limit
	}
	cache := New(config)

	// Add 5 items of size 20 each
	for i := 0; i < 5; i++ {
		cache.Put(fmt.Sprintf("item%d", i), &sizedItem{key: fmt.Sprintf("item%d", i), size: 20})
	}

	assert.Equal(t, 5, cache.Len())

	// Add large item that needs multiple evictions
	cache.Put("large", &sizedItem{key: "large", size: 60})

	// Should have evicted enough to fit
	stats := cache.Stats()
	assert.LessOrEqual(t, stats.MemoryBytes, int64(100))
	assert.Greater(t, stats.Evictions, int64(0))

	// Large item should be present
	_, ok := cache.Get("large")
	assert.True(t, ok, "large item should be cached")
}

func TestCache_EvictionSequence(t *testing.T) {
	// Test: Verify exact eviction sequence
	config := CacheConfig{
		MaxMemoryBytes: 1000,
		MaxEntries:     3,
	}
	cache := New(config)

	// Insert order: A, B, C
	cache.Put("A", &sizedItem{key: "A", size: 10})
	time.Sleep(time.Millisecond) // Ensure different timestamps
	cache.Put("B", &sizedItem{key: "B", size: 10})
	time.Sleep(time.Millisecond)
	cache.Put("C", &sizedItem{key: "C", size: 10})

	// Now order is: front=C, B, A=back (LRU)

	// Add D - A should be evicted
	cache.Put("D", &sizedItem{key: "D", size: 10})
	_, ok := cache.Get("A")
	assert.False(t, ok, "A should be evicted first")

	// Add E - B should be evicted
	cache.Put("E", &sizedItem{key: "E", size: 10})
	_, ok = cache.Get("B")
	assert.False(t, ok, "B should be evicted second")

	// Add F - C should be evicted
	cache.Put("F", &sizedItem{key: "F", size: 10})
	_, ok = cache.Get("C")
	assert.False(t, ok, "C should be evicted third")

	// D, E, F should remain
	for _, key := range []string{"D", "E", "F"} {
		_, ok = cache.Get(key)
		assert.True(t, ok, "%s should be present", key)
	}
}

// ============================================================================
// Memory Limit Enforcement Tests
// ============================================================================

func TestCache_MemoryLimitEnforced(t *testing.T) {
	// Test: Memory limit prevents overflow
	config := CacheConfig{
		MaxMemoryBytes: 100,
		MaxEntries:     0, // Memory only
	}
	cache := New(config)

	// Add items until limit exceeded
	for i := 0; i < 20; i++ {
		cache.Put(fmt.Sprintf("key%d", i), &sizedItem{
			key:  fmt.Sprintf("key%d", i),
			size: 30,
		})
	}

	stats := cache.Stats()
	assert.LessOrEqual(t, stats.MemoryBytes, int64(100),
		"Memory should never exceed limit")
}

func TestCache_LargeItemHandling(t *testing.T) {
	// Test: Single large item larger than cache limit
	// The cache implementation evicts until memory is under limit,
	// which means an item larger than the limit will be evicted immediately
	config := CacheConfig{
		MaxMemoryBytes: 100,
		MaxEntries:     0,
	}
	cache := New(config)

	// Add some normal items first
	cache.Put("small1", &sizedItem{key: "small1", size: 20})
	cache.Put("small2", &sizedItem{key: "small2", size: 20})

	// Add item larger than total cache
	cache.Put("huge", &sizedItem{key: "huge", size: 200})

	// The huge item will be evicted because it exceeds the cache limit
	// The eviction loop runs until memory <= limit, evicting even the just-added item
	_, ok := cache.Get("huge")
	assert.False(t, ok, "huge item should be evicted (too large for cache)")

	// All items should be evicted
	_, ok = cache.Get("small1")
	assert.False(t, ok, "small1 should be evicted")
	_, ok = cache.Get("small2")
	assert.False(t, ok, "small2 should be evicted")

	// Cache should be empty
	assert.Equal(t, 0, cache.Len())
}

func TestCache_ItemFitsExactly(t *testing.T) {
	// Test: Item that fits exactly at the limit
	config := CacheConfig{
		MaxMemoryBytes: 100,
		MaxEntries:     0,
	}
	cache := New(config)

	// Add item that's exactly the limit size
	cache.Put("exact", &sizedItem{key: "exact", size: 100})

	val, ok := cache.Get("exact")
	assert.True(t, ok, "item at exact limit should be cached")
	assert.Equal(t, "exact", val.(*sizedItem).key)

	// Add another small item - should evict the first
	cache.Put("small", &sizedItem{key: "small", size: 10})

	_, ok = cache.Get("exact")
	assert.False(t, ok, "exact item should be evicted")
	_, ok = cache.Get("small")
	assert.True(t, ok, "small item should remain")
}

func TestCache_MemoryAccounting(t *testing.T) {
	// Test: Memory accounting accuracy
	config := CacheConfig{
		MaxMemoryBytes: 1000,
		MaxEntries:     0,
	}
	cache := New(config)

	// Add items and track expected size
	expectedSize := int64(0)

	sizes := []int64{10, 20, 30, 40, 50}
	for i, size := range sizes {
		cache.Put(fmt.Sprintf("key%d", i), &sizedItem{
			key:  fmt.Sprintf("key%d", i),
			size: size,
		})
		expectedSize += size
	}

	stats := cache.Stats()
	assert.Equal(t, expectedSize, stats.MemoryBytes,
		"Memory accounting should match sum of item sizes")
}

func TestCache_MemoryFreedOnEviction(t *testing.T) {
	// Test: Memory freed on eviction
	config := CacheConfig{
		MaxMemoryBytes: 100,
		MaxEntries:     0,
	}
	cache := New(config)

	// Fill cache to limit
	for i := 0; i < 5; i++ {
		cache.Put(fmt.Sprintf("key%d", i), &sizedItem{
			key:  fmt.Sprintf("key%d", i),
			size: 20,
		})
	}

	stats := cache.Stats()
	assert.Equal(t, int64(100), stats.MemoryBytes)

	// Add another item - should evict and memory should stay bounded
	cache.Put("new", &sizedItem{key: "new", size: 20})

	stats = cache.Stats()
	assert.LessOrEqual(t, stats.MemoryBytes, int64(100),
		"Memory should be freed on eviction")
}

func TestCache_MemoryUpdatedOnReplace(t *testing.T) {
	// Test: Memory updated correctly when replacing item
	config := CacheConfig{
		MaxMemoryBytes: 1000,
		MaxEntries:     0,
	}
	cache := New(config)

	cache.Put("key", &sizedItem{key: "key", size: 50})
	stats1 := cache.Stats()
	assert.Equal(t, int64(50), stats1.MemoryBytes)

	// Replace with larger item
	cache.Put("key", &sizedItem{key: "key", size: 100})
	stats2 := cache.Stats()
	assert.Equal(t, int64(100), stats2.MemoryBytes,
		"Memory should update on replacement")

	// Replace with smaller item
	cache.Put("key", &sizedItem{key: "key", size: 25})
	stats3 := cache.Stats()
	assert.Equal(t, int64(25), stats3.MemoryBytes,
		"Memory should decrease on smaller replacement")
}

func TestCache_MemoryFreedOnDelete(t *testing.T) {
	// Test: Memory freed on explicit delete
	config := CacheConfig{
		MaxMemoryBytes: 1000,
		MaxEntries:     0,
	}
	cache := New(config)

	cache.Put("key1", &sizedItem{key: "key1", size: 50})
	cache.Put("key2", &sizedItem{key: "key2", size: 50})

	stats := cache.Stats()
	assert.Equal(t, int64(100), stats.MemoryBytes)

	cache.Delete("key1")

	stats = cache.Stats()
	assert.Equal(t, int64(50), stats.MemoryBytes,
		"Memory should be freed on delete")
}

// ============================================================================
// Concurrent Access Tests
// ============================================================================

func TestCache_ConcurrentReads(t *testing.T) {
	// Test: 100 goroutines reading same item
	cache := New(DefaultConfig())

	// Pre-populate
	cache.Put("shared", &sizedItem{key: "shared", size: 100})

	var wg sync.WaitGroup
	numGoroutines := 100
	numReads := 1000
	successCount := int64(0)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numReads; j++ {
				val, ok := cache.Get("shared")
				if ok && val.(*sizedItem).key == "shared" {
					atomic.AddInt64(&successCount, 1)
				}
			}
		}()
	}

	wg.Wait()

	// All reads should succeed
	assert.Equal(t, int64(numGoroutines*numReads), successCount,
		"All concurrent reads should succeed")
}

func TestCache_ConcurrentWrites(t *testing.T) {
	// Test: 100 goroutines writing different items
	cache := New(DefaultConfig())

	var wg sync.WaitGroup
	numGoroutines := 100
	numWrites := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numWrites; j++ {
				key := fmt.Sprintf("g%d-k%d", id, j)
				cache.Put(key, &sizedItem{key: key, size: 10})
			}
		}(i)
	}

	wg.Wait()

	// Verify cache is in consistent state
	stats := cache.Stats()
	assert.Greater(t, stats.Entries, 0, "Cache should have entries")

	// Verify we can still read
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("g%d-k%d", i, 99)
		_, ok := cache.Get(key)
		// May or may not be present due to eviction, but should not panic
		_ = ok
	}
}

func TestCache_ConcurrentMixedAccess(t *testing.T) {
	// Test: Mix of reads and writes
	cache := New(CacheConfig{
		MaxMemoryBytes: 10000,
		MaxEntries:     100,
	})

	var wg sync.WaitGroup
	numGoroutines := 50
	numOps := 500

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				key := fmt.Sprintf("key%d", j%50)
				if j%3 == 0 {
					cache.Put(key, &sizedItem{key: key, size: 10})
				} else {
					cache.Get(key)
				}
			}
		}(i)
	}

	wg.Wait()

	// Cache should be in valid state
	stats := cache.Stats()
	assert.GreaterOrEqual(t, stats.Entries, 0)
	assert.LessOrEqual(t, stats.Entries, 100)
}

func TestCache_ConcurrentEviction(t *testing.T) {
	// Test: Writes trigger concurrent evictions
	config := CacheConfig{
		MaxMemoryBytes: 100,
		MaxEntries:     10,
	}
	cache := New(config)

	var wg sync.WaitGroup
	numGoroutines := 20
	numWrites := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numWrites; j++ {
				key := fmt.Sprintf("g%d-k%d", id, j)
				cache.Put(key, &sizedItem{key: key, size: 20})
			}
		}(i)
	}

	wg.Wait()

	stats := cache.Stats()
	assert.LessOrEqual(t, stats.Entries, 10, "Entry limit should be respected")
	assert.Greater(t, stats.Evictions, int64(0), "Evictions should have occurred")
}

func TestCache_ReadDuringEviction(t *testing.T) {
	// Test: Item being read while eviction runs
	config := CacheConfig{
		MaxMemoryBytes: 100,
		MaxEntries:     5,
	}
	cache := New(config)

	// Pre-populate
	for i := 0; i < 5; i++ {
		cache.Put(fmt.Sprintf("key%d", i), &sizedItem{
			key:  fmt.Sprintf("key%d", i),
			size: 10,
		})
	}

	var wg sync.WaitGroup

	// Readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				key := fmt.Sprintf("key%d", j%5)
				// Should either return value or not found, never panic
				_, _ = cache.Get(key)
			}
		}()
	}

	// Writers causing evictions
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("new%d-%d", id, j)
				cache.Put(key, &sizedItem{key: key, size: 20})
			}
		}(i)
	}

	wg.Wait()

	// Should complete without deadlock or panic
	stats := cache.Stats()
	assert.Greater(t, stats.Evictions, int64(0))
}

func TestCache_ConcurrentClear(t *testing.T) {
	// Test: Clear during concurrent access
	config := CacheConfig{
		MaxMemoryBytes: 10000,
		MaxEntries:     0,
	}
	cache := New(config)

	var wg sync.WaitGroup

	// Writers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("key%d-%d", id, j)
				cache.Put(key, &sizedItem{key: key, size: 10})
			}
		}(i)
	}

	// Clearer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			cache.Clear()
			time.Sleep(time.Millisecond)
		}
	}()

	wg.Wait()

	// Should complete without deadlock
	_ = cache.Stats()
}

// ============================================================================
// Invalidation Tests
// ============================================================================

func TestCache_SingleInvalidation(t *testing.T) {
	// Test: Invalidate one item
	rc := NewResourceCache(DefaultConfig())

	rc.Put("resource1", &sizedItem{key: "r1", size: 10})
	rc.Put("resource2", &sizedItem{key: "r2", size: 10})
	rc.Put("resource3", &sizedItem{key: "r3", size: 10})

	// Invalidate one
	rc.Invalidate("resource2")

	// Only resource2 should be gone
	_, ok := rc.Get("resource1")
	assert.True(t, ok, "resource1 should remain")

	_, ok = rc.Get("resource2")
	assert.False(t, ok, "resource2 should be invalidated")

	_, ok = rc.Get("resource3")
	assert.True(t, ok, "resource3 should remain")
}

func TestCache_PatternInvalidation(t *testing.T) {
	// Test: Invalidate by prefix/pattern
	rc := NewResourceCache(DefaultConfig())

	// Add resources with different prefixes
	rc.Put("mod1.dashboard.test1", &sizedItem{key: "d1", size: 10})
	rc.Put("mod1.dashboard.test2", &sizedItem{key: "d2", size: 10})
	rc.Put("mod1.query.test1", &sizedItem{key: "q1", size: 10})
	rc.Put("mod2.dashboard.test1", &sizedItem{key: "d3", size: 10})

	// Invalidate all dashboards in mod1
	rc.InvalidateAll(func(key string) bool {
		return strings.HasPrefix(key, "mod1.dashboard")
	})

	// mod1 dashboards should be gone
	_, ok := rc.Get("mod1.dashboard.test1")
	assert.False(t, ok)
	_, ok = rc.Get("mod1.dashboard.test2")
	assert.False(t, ok)

	// Query and mod2 dashboard should remain
	_, ok = rc.Get("mod1.query.test1")
	assert.True(t, ok)
	_, ok = rc.Get("mod2.dashboard.test1")
	assert.True(t, ok)
}

func TestCache_FullClear(t *testing.T) {
	// Test: Clear entire cache
	rc := NewResourceCache(DefaultConfig())

	for i := 0; i < 100; i++ {
		rc.Put(fmt.Sprintf("resource%d", i), &sizedItem{
			key:  fmt.Sprintf("r%d", i),
			size: 10,
		})
	}

	require.Equal(t, 100, rc.Len())

	rc.Clear()

	assert.Equal(t, 0, rc.Len())
	stats := rc.Stats()
	assert.Equal(t, 0, stats.Entries)
	assert.Equal(t, int64(0), stats.MemoryBytes)
}

func TestCache_InvalidationDuringAccess(t *testing.T) {
	// Test: Concurrent read and invalidation
	rc := NewResourceCache(DefaultConfig())

	// Pre-populate
	for i := 0; i < 100; i++ {
		rc.Put(fmt.Sprintf("resource%d", i), &sizedItem{
			key:  fmt.Sprintf("r%d", i),
			size: 10,
		})
	}

	var wg sync.WaitGroup

	// Readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 500; j++ {
				key := fmt.Sprintf("resource%d", j%100)
				_, _ = rc.Get(key)
			}
		}()
	}

	// Invalidator
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			rc.Invalidate(fmt.Sprintf("resource%d", i*2))
		}
	}()

	wg.Wait()

	// Should complete without panic or deadlock
	assert.Less(t, rc.Len(), 100, "Some items should be invalidated")
}

func TestCache_ReAddAfterInvalidation(t *testing.T) {
	// Test: Invalidate then re-add
	rc := NewResourceCache(DefaultConfig())

	rc.Put("resource", &sizedItem{key: "original", size: 10})

	val, _ := rc.Get("resource")
	assert.Equal(t, "original", val.(*sizedItem).key)

	rc.Invalidate("resource")
	_, ok := rc.Get("resource")
	assert.False(t, ok)

	// Re-add with new value
	rc.Put("resource", &sizedItem{key: "new", size: 20})

	val, ok = rc.Get("resource")
	assert.True(t, ok)
	assert.Equal(t, "new", val.(*sizedItem).key)
}

func TestCache_InvalidateNonExistent(t *testing.T) {
	// Test: Invalidate non-existent key
	rc := NewResourceCache(DefaultConfig())

	rc.Put("exists", &sizedItem{key: "e", size: 10})

	// Should not panic
	rc.Invalidate("nonexistent")

	// Existing item should be unaffected
	_, ok := rc.Get("exists")
	assert.True(t, ok)
}

// ============================================================================
// Edge Case Tests
// ============================================================================

func TestCache_EmptyCache(t *testing.T) {
	// Test: Empty cache operations
	cache := New(DefaultConfig())

	// Get from empty cache
	val, ok := cache.Get("nonexistent")
	assert.False(t, ok)
	assert.Nil(t, val)

	// Delete from empty cache
	cache.Delete("nonexistent") // Should not panic

	// Stats on empty cache
	stats := cache.Stats()
	assert.Equal(t, 0, stats.Entries)
	assert.Equal(t, int64(0), stats.MemoryBytes)
	assert.Equal(t, int64(0), stats.Hits)
	assert.Equal(t, int64(1), stats.Misses) // One miss from Get
	assert.Equal(t, float64(0), stats.HitRate)

	// Len on empty cache
	assert.Equal(t, 0, cache.Len())

	// Keys on empty cache
	keys := cache.Keys()
	assert.Empty(t, keys)
}

func TestCache_ZeroSizeLimit(t *testing.T) {
	// Test: Create cache with 0 memory limit
	// With MaxMemoryBytes=0, any item with size > 0 will be immediately evicted
	// because currentMemory (> 0) will always exceed MaxMemoryBytes (0)
	config := CacheConfig{
		MaxMemoryBytes: 0, // Zero limit
		MaxEntries:     0,
	}
	cache := New(config)

	// Item with non-zero size will be evicted immediately
	cache.Put("key", &sizedItem{key: "key", size: 100})
	_, ok := cache.Get("key")
	assert.False(t, ok, "item with size should be evicted when limit is 0")

	// Items with zero size CAN be cached (their memory doesn't exceed limit)
	cache.Put("zerosize", &sizedItem{key: "zerosize", size: 0})
	val, ok := cache.Get("zerosize")
	assert.True(t, ok, "zero-size items should be cacheable")
	assert.Equal(t, "zerosize", val.(*sizedItem).key)

	// Non-Sizer items (size=0 by default) can also be cached
	cache.Put("string", "simple string")
	val, ok = cache.Get("string")
	assert.True(t, ok, "non-Sizer items (size 0) should be cacheable")
	assert.Equal(t, "simple string", val)
}

func TestCache_NegativeSizeItem(t *testing.T) {
	// Test: Item with negative Size()
	cache := New(DefaultConfig())

	// Add item with negative size
	cache.Put("negative", &negativeSizedItem{})

	val, ok := cache.Get("negative")
	assert.True(t, ok)
	assert.NotNil(t, val)

	// Memory should reflect the negative size (edge case)
	stats := cache.Stats()
	assert.Equal(t, int64(-100), stats.MemoryBytes)
}

func TestCache_NilValue(t *testing.T) {
	// Test: Cache nil value
	cache := New(DefaultConfig())

	// Should be able to store nil
	cache.Put("nilkey", nil)

	val, ok := cache.Get("nilkey")
	assert.True(t, ok, "nil value should be retrievable")
	assert.Nil(t, val)

	// Len should reflect the entry
	assert.Equal(t, 1, cache.Len())
}

func TestCache_KeyCollision(t *testing.T) {
	// Test: Same key, different values - later value overwrites
	cache := New(DefaultConfig())

	cache.Put("key", &sizedItem{key: "first", size: 10})
	val, _ := cache.Get("key")
	assert.Equal(t, "first", val.(*sizedItem).key)

	cache.Put("key", &sizedItem{key: "second", size: 20})
	val, _ = cache.Get("key")
	assert.Equal(t, "second", val.(*sizedItem).key)

	// Should only have 1 entry
	assert.Equal(t, 1, cache.Len())
}

func TestCache_EmptyKey(t *testing.T) {
	// Test: Empty string key
	cache := New(DefaultConfig())

	cache.Put("", &sizedItem{key: "empty-key", size: 10})

	val, ok := cache.Get("")
	assert.True(t, ok)
	assert.Equal(t, "empty-key", val.(*sizedItem).key)

	cache.Delete("")
	_, ok = cache.Get("")
	assert.False(t, ok)
}

func TestCache_VeryLongKey(t *testing.T) {
	// Test: Very long key
	cache := New(DefaultConfig())

	longKey := strings.Repeat("x", 10000)
	cache.Put(longKey, &sizedItem{key: "long", size: 10})

	val, ok := cache.Get(longKey)
	assert.True(t, ok)
	assert.Equal(t, "long", val.(*sizedItem).key)
}

func TestCache_SpecialCharacterKeys(t *testing.T) {
	// Test: Keys with special characters
	cache := New(DefaultConfig())

	specialKeys := []string{
		"key with spaces",
		"key\twith\ttabs",
		"key\nwith\nnewlines",
		"key/with/slashes",
		"key.with.dots",
		"mod.local.dashboard.test",
		"日本語キー",
		"🔑emoji🔑",
	}

	for _, key := range specialKeys {
		cache.Put(key, &sizedItem{key: key, size: 10})
	}

	for _, key := range specialKeys {
		val, ok := cache.Get(key)
		assert.True(t, ok, "key %q should be retrievable", key)
		assert.Equal(t, key, val.(*sizedItem).key)
	}
}

// ============================================================================
// ResourceCache Specific Tests
// ============================================================================

func TestResourceCache_TypeSafe(t *testing.T) {
	// Test: Type-safe resource retrieval
	rc := NewResourceCache(DefaultConfig())

	// Use generic Put/Get for testing
	rc.Put("test.dashboard", &sizedItem{key: "dashboard", size: 10})
	rc.Put("test.query", &sizedItem{key: "query", size: 10})

	val, ok := rc.Get("test.dashboard")
	assert.True(t, ok)
	assert.Equal(t, "dashboard", val.(*sizedItem).key)

	val, ok = rc.Get("test.query")
	assert.True(t, ok)
	assert.Equal(t, "query", val.(*sizedItem).key)
}

func TestResourceCache_WrongType(t *testing.T) {
	// Test: Store one type, try to get as different type
	rc := NewResourceCache(DefaultConfig())

	// Store a string
	rc.Put("test.resource", "a string value")

	// Get returns interface{}, type assertion is caller's responsibility
	val, ok := rc.Get("test.resource")
	assert.True(t, ok)
	assert.Equal(t, "a string value", val)

	// Trying to cast to wrong type should fail gracefully
	_, isSized := val.(*sizedItem)
	assert.False(t, isSized)
}

func TestResourceCache_GetResourceWithNonHclResource(t *testing.T) {
	// Test: GetResource with non-HclResource value
	rc := NewResourceCache(DefaultConfig())

	// Put a non-HclResource value
	rc.Put("test.resource", &sizedItem{key: "not-hcl", size: 10})

	// GetResource should return nil, false
	resource, ok := rc.GetResource("test.resource")
	assert.False(t, ok)
	assert.Nil(t, resource)
}

func TestResourceCache_KeysAndLen(t *testing.T) {
	// Test: Keys and Len methods
	rc := NewResourceCache(DefaultConfig())

	rc.Put("key1", &sizedItem{key: "k1", size: 10})
	rc.Put("key2", &sizedItem{key: "k2", size: 10})
	rc.Put("key3", &sizedItem{key: "k3", size: 10})

	assert.Equal(t, 3, rc.Len())

	keys := rc.Keys()
	assert.Len(t, keys, 3)
	assert.Contains(t, keys, "key1")
	assert.Contains(t, keys, "key2")
	assert.Contains(t, keys, "key3")
}

// ============================================================================
// Statistics Tests
// ============================================================================

func TestCache_HitRateTracking(t *testing.T) {
	// Test: Hit rate tracking accuracy
	cache := New(DefaultConfig())

	cache.Put("key1", &sizedItem{key: "key1", size: 10})
	cache.Put("key2", &sizedItem{key: "key2", size: 10})

	// 3 hits
	cache.Get("key1")
	cache.Get("key1")
	cache.Get("key2")

	// 2 misses
	cache.Get("nonexistent1")
	cache.Get("nonexistent2")

	stats := cache.Stats()
	assert.Equal(t, int64(3), stats.Hits)
	assert.Equal(t, int64(2), stats.Misses)
	assert.InDelta(t, 0.6, stats.HitRate, 0.001) // 3/5 = 0.6
}

func TestCache_ItemCountAccuracy(t *testing.T) {
	// Test: Item count always correct
	cache := New(DefaultConfig())

	assert.Equal(t, 0, cache.Len())

	// Add items
	for i := 0; i < 10; i++ {
		cache.Put(fmt.Sprintf("key%d", i), &sizedItem{key: fmt.Sprintf("k%d", i), size: 10})
		assert.Equal(t, i+1, cache.Len())
	}

	// Remove items
	for i := 0; i < 5; i++ {
		cache.Delete(fmt.Sprintf("key%d", i))
		assert.Equal(t, 9-i, cache.Len())
	}
}

func TestCache_StatsAfterClear(t *testing.T) {
	// Test: Stats after clear
	cache := New(DefaultConfig())

	// Add items and generate some hits/misses
	for i := 0; i < 10; i++ {
		cache.Put(fmt.Sprintf("key%d", i), &sizedItem{key: fmt.Sprintf("k%d", i), size: 10})
	}
	cache.Get("key1")
	cache.Get("nonexistent")

	// Clear cache
	cache.Clear()

	stats := cache.Stats()
	assert.Equal(t, 0, stats.Entries)
	assert.Equal(t, int64(0), stats.MemoryBytes)
	// Note: hits/misses are NOT reset by Clear() based on the implementation
	assert.Equal(t, int64(1), stats.Hits)
	assert.Equal(t, int64(1), stats.Misses)
}

func TestCache_EvictionCount(t *testing.T) {
	// Test: Eviction count accuracy
	config := CacheConfig{
		MaxMemoryBytes: 100,
		MaxEntries:     5,
	}
	cache := New(config)

	// Add 10 items - should trigger 5 evictions
	for i := 0; i < 10; i++ {
		cache.Put(fmt.Sprintf("key%d", i), &sizedItem{key: fmt.Sprintf("k%d", i), size: 10})
	}

	stats := cache.Stats()
	assert.Equal(t, 5, stats.Entries)
	assert.Equal(t, int64(5), stats.Evictions)
}

// ============================================================================
// TTL Tests
// ============================================================================

func TestCache_TTLExpiryEdge(t *testing.T) {
	// Test: TTL expiry edge behavior
	config := CacheConfig{
		MaxMemoryBytes: 1024 * 1024,
		TTL:            50 * time.Millisecond,
	}
	cache := New(config)

	cache.Put("key", &sizedItem{key: "k", size: 10})

	// Should exist immediately
	_, ok := cache.Get("key")
	assert.True(t, ok)

	// Wait for expiry
	time.Sleep(100 * time.Millisecond)

	// Should be expired
	_, ok = cache.Get("key")
	assert.False(t, ok)
}

func TestCache_TTLRefreshOnAccess(t *testing.T) {
	// Test: TTL is NOT refreshed on access (based on implementation)
	config := CacheConfig{
		MaxMemoryBytes: 1024 * 1024,
		TTL:            100 * time.Millisecond,
	}
	cache := New(config)

	cache.Put("key", &sizedItem{key: "k", size: 10})

	// Access multiple times before expiry
	time.Sleep(30 * time.Millisecond)
	cache.Get("key")
	time.Sleep(30 * time.Millisecond)
	cache.Get("key")
	time.Sleep(30 * time.Millisecond)
	cache.Get("key")

	// Total time ~90ms, but TTL is from original Put time
	// After another 20ms, should be expired
	time.Sleep(20 * time.Millisecond)

	_, ok := cache.Get("key")
	assert.False(t, ok, "item should expire based on original timestamp")
}

func TestCache_TTLRefreshOnUpdate(t *testing.T) {
	// Test: TTL IS refreshed on update (Put)
	config := CacheConfig{
		MaxMemoryBytes: 1024 * 1024,
		TTL:            100 * time.Millisecond,
	}
	cache := New(config)

	cache.Put("key", &sizedItem{key: "v1", size: 10})

	// Wait 80ms
	time.Sleep(80 * time.Millisecond)

	// Update the item - should refresh TTL
	cache.Put("key", &sizedItem{key: "v2", size: 10})

	// Wait another 80ms (160ms total from first Put, but only 80ms from update)
	time.Sleep(80 * time.Millisecond)

	// Should still be present because TTL was refreshed
	val, ok := cache.Get("key")
	assert.True(t, ok, "item should still be present after TTL refresh")
	assert.Equal(t, "v2", val.(*sizedItem).key)
}

// ============================================================================
// Performance Benchmarks
// ============================================================================

func BenchmarkCache_GetHit(b *testing.B) {
	cache := New(DefaultConfig())

	// Pre-populate
	for i := 0; i < 1000; i++ {
		cache.Put(fmt.Sprintf("key%d", i), &sizedItem{
			key:  fmt.Sprintf("item%d", i),
			size: 100,
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(fmt.Sprintf("key%d", i%1000))
	}
}

func BenchmarkCache_GetMiss(b *testing.B) {
	cache := New(DefaultConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(fmt.Sprintf("missing%d", i))
	}
}

func BenchmarkCache_PutNew(b *testing.B) {
	config := CacheConfig{
		MaxMemoryBytes: 1024 * 1024 * 1024, // 1GB to avoid eviction
		MaxEntries:     0,
	}
	cache := New(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Put(fmt.Sprintf("key%d", i), &sizedItem{
			key:  fmt.Sprintf("item%d", i),
			size: 100,
		})
	}
}

func BenchmarkCache_PutUpdate(b *testing.B) {
	cache := New(DefaultConfig())
	cache.Put("key", &sizedItem{key: "initial", size: 100})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Put("key", &sizedItem{
			key:  fmt.Sprintf("item%d", i),
			size: 100,
		})
	}
}

func BenchmarkCache_ConcurrentGetParallel(b *testing.B) {
	cache := New(DefaultConfig())

	// Pre-populate
	for i := 0; i < 1000; i++ {
		cache.Put(fmt.Sprintf("key%d", i), &sizedItem{
			key:  fmt.Sprintf("item%d", i),
			size: 100,
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

func BenchmarkCache_ConcurrentPutParallel(b *testing.B) {
	config := CacheConfig{
		MaxMemoryBytes: 1024 * 1024 * 100, // 100MB
		MaxEntries:     0,
	}
	cache := New(config)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Put(fmt.Sprintf("key%d", i%10000), &sizedItem{
				key:  fmt.Sprintf("item%d", i),
				size: 100,
			})
			i++
		}
	})
}

func BenchmarkCache_EvictionOverhead(b *testing.B) {
	config := CacheConfig{
		MaxMemoryBytes: 1000,
		MaxEntries:     10, // Force frequent eviction
	}
	cache := New(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Put(fmt.Sprintf("key%d", i), &sizedItem{
			key:  fmt.Sprintf("item%d", i),
			size: 100,
		})
	}
}

func TestCache_LargeItemCount(t *testing.T) {
	// Test: 100,000 small items
	config := CacheConfig{
		MaxMemoryBytes: 1024 * 1024 * 100, // 100MB
		MaxEntries:     0,
	}
	cache := New(config)

	numItems := 100000
	start := time.Now()

	// Insert items
	for i := 0; i < numItems; i++ {
		cache.Put(fmt.Sprintf("key%d", i), &sizedItem{
			key:  fmt.Sprintf("item%d", i),
			size: 10, // Small items
		})
	}

	insertTime := time.Since(start)
	t.Logf("Inserted %d items in %v", numItems, insertTime)

	// Verify count
	assert.Equal(t, numItems, cache.Len())

	// Test lookup performance
	start = time.Now()
	for i := 0; i < 10000; i++ {
		cache.Get(fmt.Sprintf("key%d", i))
	}
	lookupTime := time.Since(start)
	t.Logf("10000 lookups in %v", lookupTime)

	// Operations should be reasonably fast (less than 1 second for each phase)
	assert.Less(t, insertTime, 5*time.Second, "Insertion should be fast")
	assert.Less(t, lookupTime, time.Second, "Lookups should be fast")
}

func BenchmarkCache_Delete(b *testing.B) {
	cache := New(DefaultConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cache.Put("key", &sizedItem{key: "item", size: 100})
		b.StartTimer()
		cache.Delete("key")
	}
}

func BenchmarkCache_Clear(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cache := New(DefaultConfig())
		for j := 0; j < 1000; j++ {
			cache.Put(fmt.Sprintf("key%d", j), &sizedItem{key: fmt.Sprintf("item%d", j), size: 100})
		}
		b.StartTimer()
		cache.Clear()
	}
}
