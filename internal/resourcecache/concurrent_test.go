package resourcecache

// Concurrent access stress tests for ResourceCache and Cache.
//
// These tests exercise concurrent access patterns to find race conditions,
// deadlocks, and data corruption issues.
//
// Run with race detector:
//   go test -race -timeout 120s ./internal/resourcecache/... -run Concurrent

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// Concurrent Cache Access Tests
// =============================================================================

// TestConcurrent_CacheAccess tests heavy concurrent reads and writes with LRU eviction.
func TestConcurrent_CacheAccess(t *testing.T) {
	// Create cache with small size to trigger evictions
	config := CacheConfig{
		MaxMemoryBytes: 10 * 1024, // 10KB to trigger evictions
		MaxEntries:     100,
	}
	cache := New(config)

	const numGoroutines = 50
	const opsPerGoroutine = 200
	var wg sync.WaitGroup
	var panicCount int32

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					atomic.AddInt32(&panicCount, 1)
				}
			}()

			for j := 0; j < opsPerGoroutine; j++ {
				key := fmt.Sprintf("key_%d_%d", id, j%50) // Reuse some keys

				// Mix of operations
				switch j % 4 {
				case 0:
					// Put
					cache.Put(key, fmt.Sprintf("value_%d_%d", id, j))
				case 1:
					// Get
					_, _ = cache.Get(key)
				case 2:
					// Put (another write)
					cache.Put(key, fmt.Sprintf("updated_%d_%d", id, j))
				case 3:
					// Get (another read)
					_, _ = cache.Get(key)
				}
			}
		}(i)
	}

	// Wait with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(30 * time.Second):
		t.Fatal("Test timed out - possible deadlock")
	}

	assert.Equal(t, int32(0), panicCount, "No panics should occur during concurrent access")

	// Verify cache is still functional
	stats := cache.Stats()
	t.Logf("Cache stats: entries=%d, hits=%d, misses=%d, evictions=%d",
		stats.Entries, stats.Hits, stats.Misses, stats.Evictions)
}

// TestConcurrent_ResourceCacheAccess tests ResourceCache with concurrent resource operations.
func TestConcurrent_ResourceCacheAccess(t *testing.T) {
	config := CacheConfig{
		MaxMemoryBytes: 50 * 1024 * 1024, // 50MB
	}
	cache := NewResourceCache(config)

	const numGoroutines = 50
	const opsPerGoroutine = 100
	var wg sync.WaitGroup
	var panicCount int32

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					atomic.AddInt32(&panicCount, 1)
				}
			}()

			for j := 0; j < opsPerGoroutine; j++ {
				key := fmt.Sprintf("resource_%d_%d", id, j%30)

				switch j % 5 {
				case 0:
					cache.Put(key, &testItem{name: key})
				case 1:
					_, _ = cache.Get(key)
				case 2:
					// Use underlying cache directly (GetResource requires HclResource)
					_, _ = cache.Get(key)
				case 3:
					cache.Put(key, &testItem{name: key + "_updated"})
				case 4:
					_ = cache.Stats()
				}
			}
		}(i)
	}

	wg.Wait()

	assert.Equal(t, int32(0), panicCount, "No panics should occur")
}

// TestConcurrent_CacheEviction tests cache eviction under heavy concurrent load.
func TestConcurrent_CacheEviction(t *testing.T) {
	// Very small cache to force constant eviction
	config := CacheConfig{
		MaxMemoryBytes: 1024, // 1KB
		MaxEntries:     10,
	}
	cache := New(config)

	const numGoroutines = 30
	const opsPerGoroutine = 500
	var wg sync.WaitGroup
	var panicCount int32

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					atomic.AddInt32(&panicCount, 1)
				}
			}()

			for j := 0; j < opsPerGoroutine; j++ {
				// Use unique keys to force eviction
				key := fmt.Sprintf("evict_key_%d_%d", id, j)
				cache.Put(key, fmt.Sprintf("value_%d_%d", id, j))

				// Sometimes read back
				if j%3 == 0 {
					_, _ = cache.Get(key)
				}
			}
		}(i)
	}

	wg.Wait()

	assert.Equal(t, int32(0), panicCount, "No panics during heavy eviction")

	stats := cache.Stats()
	assert.Greater(t, stats.Evictions, int64(0), "Evictions should have occurred")
	t.Logf("Evictions: %d", stats.Evictions)
}

// TestConcurrent_CacheReadHeavy tests a read-heavy workload.
func TestConcurrent_CacheReadHeavy(t *testing.T) {
	config := CacheConfig{
		MaxMemoryBytes: 50 * 1024 * 1024,
	}
	cache := New(config)

	// Pre-populate cache
	for i := 0; i < 100; i++ {
		cache.Put(fmt.Sprintf("preload_%d", i), fmt.Sprintf("value_%d", i))
	}

	const numGoroutines = 50
	const opsPerGoroutine = 500
	var wg sync.WaitGroup
	var panicCount int32
	var readCount, writeCount int32

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					atomic.AddInt32(&panicCount, 1)
				}
			}()

			for j := 0; j < opsPerGoroutine; j++ {
				key := fmt.Sprintf("preload_%d", j%100)

				// 95% reads, 5% writes
				if j%20 == 0 {
					cache.Put(key, fmt.Sprintf("updated_%d_%d", id, j))
					atomic.AddInt32(&writeCount, 1)
				} else {
					_, _ = cache.Get(key)
					atomic.AddInt32(&readCount, 1)
				}
			}
		}(i)
	}

	wg.Wait()

	assert.Equal(t, int32(0), panicCount, "No panics should occur")
	t.Logf("Reads: %d, Writes: %d", readCount, writeCount)
}

// TestConcurrent_CacheWriteHeavy tests a write-heavy workload.
func TestConcurrent_CacheWriteHeavy(t *testing.T) {
	config := CacheConfig{
		MaxMemoryBytes: 10 * 1024, // Small to trigger evictions
		MaxEntries:     50,
	}
	cache := New(config)

	const numGoroutines = 30
	const opsPerGoroutine = 300
	var wg sync.WaitGroup
	var panicCount int32

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					atomic.AddInt32(&panicCount, 1)
				}
			}()

			for j := 0; j < opsPerGoroutine; j++ {
				key := fmt.Sprintf("write_key_%d_%d", id, j)

				// 80% writes, 20% reads
				if j%5 == 0 {
					_, _ = cache.Get(key)
				} else {
					cache.Put(key, fmt.Sprintf("value_%d_%d", id, j))
				}
			}
		}(i)
	}

	wg.Wait()

	assert.Equal(t, int32(0), panicCount, "No panics during write-heavy load")
}

// TestConcurrent_CacheDelete tests concurrent delete operations.
func TestConcurrent_CacheDelete(t *testing.T) {
	config := CacheConfig{
		MaxMemoryBytes: 50 * 1024 * 1024,
	}
	cache := New(config)

	// Pre-populate
	for i := 0; i < 200; i++ {
		cache.Put(fmt.Sprintf("delete_key_%d", i), fmt.Sprintf("value_%d", i))
	}

	const numGoroutines = 40
	var wg sync.WaitGroup
	var panicCount int32

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					atomic.AddInt32(&panicCount, 1)
				}
			}()

			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("delete_key_%d", (id*10+j)%200)

				switch j % 4 {
				case 0:
					cache.Delete(key)
				case 1:
					_, _ = cache.Get(key)
				case 2:
					cache.Put(key, fmt.Sprintf("new_value_%d_%d", id, j))
				case 3:
					cache.Delete(key)
				}
			}
		}(i)
	}

	wg.Wait()

	assert.Equal(t, int32(0), panicCount, "No panics during concurrent delete")
}

// TestConcurrent_CacheClear tests concurrent Clear operations.
func TestConcurrent_CacheClear(t *testing.T) {
	config := CacheConfig{
		MaxMemoryBytes: 50 * 1024 * 1024,
	}
	cache := New(config)

	const numGoroutines = 20
	var wg sync.WaitGroup
	var panicCount int32

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					atomic.AddInt32(&panicCount, 1)
				}
			}()

			for j := 0; j < 50; j++ {
				key := fmt.Sprintf("clear_key_%d_%d", id, j)

				switch j % 5 {
				case 0:
					cache.Put(key, fmt.Sprintf("value_%d", j))
				case 1:
					_, _ = cache.Get(key)
				case 2:
					cache.Clear()
				case 3:
					_ = cache.Len()
				case 4:
					_ = cache.Stats()
				}
			}
		}(i)
	}

	wg.Wait()

	assert.Equal(t, int32(0), panicCount, "No panics during concurrent clear")
}

// TestConcurrent_CacheStats tests concurrent Stats access.
func TestConcurrent_CacheStats(t *testing.T) {
	config := CacheConfig{
		MaxMemoryBytes: 50 * 1024 * 1024,
	}
	cache := New(config)

	const numGoroutines = 50
	const opsPerGoroutine = 200
	var wg sync.WaitGroup
	var panicCount int32

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					atomic.AddInt32(&panicCount, 1)
				}
			}()

			for j := 0; j < opsPerGoroutine; j++ {
				// Interleave stats with operations
				switch j % 3 {
				case 0:
					cache.Put(fmt.Sprintf("key_%d_%d", id, j), j)
				case 1:
					_ = cache.Stats()
				case 2:
					_ = cache.Len()
				}
			}
		}(i)
	}

	wg.Wait()

	assert.Equal(t, int32(0), panicCount, "No panics during concurrent stats access")
}

// TestConcurrent_CacheKeys tests concurrent Keys access.
func TestConcurrent_CacheKeys(t *testing.T) {
	config := CacheConfig{
		MaxMemoryBytes: 50 * 1024 * 1024,
	}
	cache := New(config)

	// Pre-populate
	for i := 0; i < 100; i++ {
		cache.Put(fmt.Sprintf("key_%d", i), i)
	}

	const numGoroutines = 30
	const opsPerGoroutine = 100
	var wg sync.WaitGroup
	var panicCount int32

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					atomic.AddInt32(&panicCount, 1)
				}
			}()

			for j := 0; j < opsPerGoroutine; j++ {
				switch j % 4 {
				case 0:
					keys := cache.Keys()
					// Iterate over returned slice (should be safe)
					for _, k := range keys {
						_ = k
					}
				case 1:
					cache.Put(fmt.Sprintf("new_key_%d_%d", id, j), j)
				case 2:
					_, _ = cache.Get(fmt.Sprintf("key_%d", j%100))
				case 3:
					cache.Delete(fmt.Sprintf("key_%d", j%100))
				}
			}
		}(i)
	}

	wg.Wait()

	assert.Equal(t, int32(0), panicCount, "No panics during concurrent Keys access")
}

// =============================================================================
// Deadlock Detection Tests
// =============================================================================

// TestConcurrent_CacheNoDeadlock tests that the cache doesn't deadlock under heavy contention.
func TestConcurrent_CacheNoDeadlock(t *testing.T) {
	config := CacheConfig{
		MaxMemoryBytes: 5 * 1024, // Small to force contention
		MaxEntries:     20,
	}
	cache := New(config)

	const numGoroutines = 50
	var wg sync.WaitGroup

	done := make(chan struct{})

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("contention_key_%d", j%10) // High contention on same keys

				switch j % 6 {
				case 0:
					cache.Put(key, j)
				case 1:
					_, _ = cache.Get(key)
				case 2:
					cache.Delete(key)
				case 3:
					_ = cache.Stats()
				case 4:
					_ = cache.Keys()
				case 5:
					_ = cache.Len()
				}
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success - no deadlock
	case <-time.After(30 * time.Second):
		t.Fatal("Test timed out - possible deadlock")
	}
}

// =============================================================================
// Memory Safety Tests
// =============================================================================

// TestConcurrent_CacheValueIntegrity tests that values are not corrupted during concurrent access.
func TestConcurrent_CacheValueIntegrity(t *testing.T) {
	config := CacheConfig{
		MaxMemoryBytes: 50 * 1024 * 1024,
	}
	cache := New(config)

	const numKeys = 20
	const numGoroutines = 30
	const opsPerGoroutine = 100

	// Track expected values
	var expectedValues sync.Map

	var wg sync.WaitGroup
	var corruptionCount int32

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < opsPerGoroutine; j++ {
				key := fmt.Sprintf("integrity_key_%d", j%numKeys)
				expectedValue := fmt.Sprintf("value_from_%d_%d", id, j)

				// Write
				cache.Put(key, expectedValue)
				expectedValues.Store(key, expectedValue)

				// Read back and verify (may have been overwritten)
				if val, ok := cache.Get(key); ok {
					strVal, isStr := val.(string)
					if isStr {
						// Value should be a valid string (not corrupted)
						if len(strVal) == 0 {
							atomic.AddInt32(&corruptionCount, 1)
						}
					}
				}
			}
		}(i)
	}

	wg.Wait()

	assert.Equal(t, int32(0), corruptionCount, "No value corruption should occur")
}

// =============================================================================
// Goroutine Leak Tests
// =============================================================================

// TestConcurrent_CacheNoGoroutineLeaks tests that cache operations don't leak goroutines.
func TestConcurrent_CacheNoGoroutineLeaks(t *testing.T) {
	runtime.GC()
	time.Sleep(100 * time.Millisecond)
	baseline := runtime.NumGoroutine()

	config := CacheConfig{
		MaxMemoryBytes: 10 * 1024,
		MaxEntries:     50,
	}
	cache := New(config)

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				key := fmt.Sprintf("leak_test_%d_%d", id, j)
				cache.Put(key, j)
				_, _ = cache.Get(key)
				cache.Delete(key)
			}
		}(i)
	}

	wg.Wait()

	// Clear and allow goroutines to settle
	cache.Clear()
	runtime.GC()
	time.Sleep(200 * time.Millisecond)

	after := runtime.NumGoroutine()
	tolerance := 5
	if after > baseline+tolerance {
		t.Errorf("Possible goroutine leak: before=%d, after=%d", baseline, after)
	}
}

// =============================================================================
// Resource Cache Specific Tests
// =============================================================================

// TestConcurrent_ResourceCacheInvalidate tests concurrent Invalidate operations.
func TestConcurrent_ResourceCacheInvalidate(t *testing.T) {
	config := CacheConfig{
		MaxMemoryBytes: 50 * 1024 * 1024,
	}
	cache := NewResourceCache(config)

	// Pre-populate using underlying cache (PutResource requires HclResource)
	for i := 0; i < 100; i++ {
		cache.Put(fmt.Sprintf("resource_%d", i), &testItem{name: fmt.Sprintf("res_%d", i)})
	}

	const numGoroutines = 30
	var wg sync.WaitGroup
	var panicCount int32

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					atomic.AddInt32(&panicCount, 1)
				}
			}()

			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("resource_%d", (id+j)%100)

				switch j % 4 {
				case 0:
					cache.Invalidate(key)
				case 1:
					_, _ = cache.Get(key)
				case 2:
					cache.Put(key, &testItem{name: key})
				case 3:
					cache.Invalidate(key)
				}
			}
		}(i)
	}

	wg.Wait()

	assert.Equal(t, int32(0), panicCount, "No panics during concurrent invalidate")
}

// TestConcurrent_ResourceCacheInvalidateAll tests concurrent InvalidateAll operations.
func TestConcurrent_ResourceCacheInvalidateAll(t *testing.T) {
	config := CacheConfig{
		MaxMemoryBytes: 50 * 1024 * 1024,
	}
	cache := NewResourceCache(config)

	const numGoroutines = 20
	var wg sync.WaitGroup
	var panicCount int32

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					atomic.AddInt32(&panicCount, 1)
				}
			}()

			for j := 0; j < 50; j++ {
				key := fmt.Sprintf("batch_resource_%d_%d", id, j)

				switch j % 5 {
				case 0:
					cache.Put(key, &testItem{name: key})
				case 1:
					_, _ = cache.Get(key)
				case 2:
					// Invalidate all resources matching predicate
					cache.InvalidateAll(func(k string) bool {
						return len(k) > 0 && k[0] == 'b' // Match keys starting with 'b'
					})
				case 3:
					_ = cache.Len()
				case 4:
					_ = cache.Keys()
				}
			}
		}(i)
	}

	wg.Wait()

	assert.Equal(t, int32(0), panicCount, "No panics during concurrent InvalidateAll")
}

// =============================================================================
// Bursty Load Tests
// =============================================================================

// TestConcurrent_CacheBursty tests cache under bursty load patterns.
func TestConcurrent_CacheBursty(t *testing.T) {
	config := CacheConfig{
		MaxMemoryBytes: 20 * 1024,
		MaxEntries:     100,
	}
	cache := New(config)

	const numBursts = 5
	const goroutinesPerBurst = 30
	const opsPerGoroutine = 50
	var panicCount int32

	for burst := 0; burst < numBursts; burst++ {
		var wg sync.WaitGroup

		for i := 0; i < goroutinesPerBurst; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				defer func() {
					if r := recover(); r != nil {
						atomic.AddInt32(&panicCount, 1)
					}
				}()

				for j := 0; j < opsPerGoroutine; j++ {
					key := fmt.Sprintf("burst_%d_%d_%d", burst, id, j)
					cache.Put(key, j)
					_, _ = cache.Get(key)
				}
			}(i)
		}

		wg.Wait()

		// Idle between bursts
		time.Sleep(50 * time.Millisecond)
	}

	assert.Equal(t, int32(0), panicCount, "No panics during bursty load")
}

// =============================================================================
// TTL Tests (if TTL is configured)
// =============================================================================

// TestConcurrent_CacheTTL tests concurrent access with TTL expiration.
func TestConcurrent_CacheTTL(t *testing.T) {
	config := CacheConfig{
		MaxMemoryBytes: 50 * 1024 * 1024,
		TTL:            100 * time.Millisecond, // Short TTL
	}
	cache := New(config)

	const numGoroutines = 20
	var wg sync.WaitGroup
	var panicCount int32

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					atomic.AddInt32(&panicCount, 1)
				}
			}()

			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("ttl_key_%d_%d", id, j%20)

				switch j % 3 {
				case 0:
					cache.Put(key, j)
				case 1:
					// This may return miss due to TTL expiry
					_, _ = cache.Get(key)
				case 2:
					// Small sleep to allow some TTL expiries
					if j%10 == 0 {
						time.Sleep(10 * time.Millisecond)
					}
				}
			}
		}(i)
	}

	wg.Wait()

	assert.Equal(t, int32(0), panicCount, "No panics with TTL")
}

// Note: testItem type is defined in cache_test.go and available for concurrent tests
