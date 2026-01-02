# Task 6: Cache Behavior Tests

## Objective

Write comprehensive tests for the resource cache, focusing on eviction behavior, memory limits, concurrent access, and invalidation correctness.

## Context

- Cache uses LRU eviction with memory limit (50MB default)
- Critical for lazy loading performance
- Must be thread-safe
- Invalidation affects consistency

## Dependencies

- Task 2 (test fixtures)
- Files to test: `internal/resourcecache/resource_cache.go`, `internal/resourcecache/cache.go`

## Acceptance Criteria

- [ ] Add tests to `internal/resourcecache/cache_edge_test.go`
- [ ] Test LRU eviction behavior precisely
- [ ] Verify memory limit enforcement
- [ ] Test concurrent access patterns
- [ ] Verify invalidation propagation

## Test Cases to Implement

### LRU Eviction
```go
// Test: Least recently used item evicted first
func TestCache_LRUEvictionOrder(t *testing.T)
// Add A, B, C to full cache
// Access A
// Add D (triggers eviction)
// B should be evicted (least recently used)

// Test: Access updates recency
func TestCache_AccessUpdatesRecency(t *testing.T)
// Add A, B, C
// Get A
// A moves to most recent

// Test: Update counts as access
func TestCache_UpdateCountsAsAccess(t *testing.T)
// Update existing item
// Should update recency

// Test: Eviction callback
func TestCache_EvictionCallback(t *testing.T)
// Register eviction callback
// Trigger eviction
// Callback called with evicted item
```

### Memory Limit Enforcement
```go
// Test: Memory limit prevents overflow
func TestCache_MemoryLimitEnforced(t *testing.T)
// Set 1MB limit
// Add items until limit
// Verify eviction keeps total under limit

// Test: Single large item handling
func TestCache_LargeItemHandling(t *testing.T)
// Item larger than cache limit
// Should still be cacheable? Or rejected?

// Test: Memory accounting accuracy
func TestCache_MemoryAccounting(t *testing.T)
// Add items, track sizes
// Verify reported memory matches sum

// Test: Memory freed on eviction
func TestCache_MemoryFreedOnEviction(t *testing.T)
// Add items, trigger eviction
// Verify memory decreases

// Test: Size() method accuracy
func TestCache_SizeMethodAccuracy(t *testing.T)
// Items with Size() method
// Verify sizes used correctly
```

### Concurrent Access
```go
// Test: Concurrent reads
func TestCache_ConcurrentReads(t *testing.T)
// 100 goroutines reading same item
// No races, correct values

// Test: Concurrent writes
func TestCache_ConcurrentWrites(t *testing.T)
// 100 goroutines writing different items
// All items stored correctly

// Test: Mixed read/write
func TestCache_ConcurrentMixedAccess(t *testing.T)
// Mix of reads and writes
// No deadlocks or races

// Test: Concurrent eviction
func TestCache_ConcurrentEviction(t *testing.T)
// Writes trigger concurrent evictions
// No races

// Test: Read during eviction
func TestCache_ReadDuringEviction(t *testing.T)
// Item being read while eviction runs
// Read should succeed or fail cleanly
```

### Invalidation
```go
// Test: Single item invalidation
func TestCache_SingleInvalidation(t *testing.T)
// Invalidate one item
// Only that item removed

// Test: Pattern-based invalidation
func TestCache_PatternInvalidation(t *testing.T)
// Invalidate by prefix/pattern
// Matching items removed

// Test: Full cache clear
func TestCache_FullClear(t *testing.T)
// Clear entire cache
// Verify empty

// Test: Invalidation during access
func TestCache_InvalidationDuringAccess(t *testing.T)
// Concurrent read and invalidation
// Clean behavior (no partial state)

// Test: Re-add after invalidation
func TestCache_ReAddAfterInvalidation(t *testing.T)
// Invalidate then re-add
// New value stored correctly
```

### Edge Cases
```go
// Test: Empty cache operations
func TestCache_EmptyCache(t *testing.T)
// Get from empty cache
// Delete from empty cache
// Stats on empty cache

// Test: Zero size limit
func TestCache_ZeroSizeLimit(t *testing.T)
// Create cache with 0 limit
// Should accept 0 items? Or unlimited?

// Test: Negative size item
func TestCache_NegativeSizeItem(t *testing.T)
// Item with negative Size()
// Should handle gracefully

// Test: Nil value handling
func TestCache_NilValue(t *testing.T)
// Cache nil value
// Should store or reject?

// Test: Key collision
func TestCache_KeyCollision(t *testing.T)
// Same key, different values
// Later value overwrites
```

### ResourceCache Specific
```go
// Test: Type-safe resource retrieval
func TestResourceCache_TypeSafe(t *testing.T)
// Store dashboard, retrieve as dashboard
// Correct type returned

// Test: Wrong type retrieval
func TestResourceCache_WrongType(t *testing.T)
// Store dashboard, try to get as control
// Should fail cleanly

// Test: Resource size calculation
func TestResourceCache_SizeCalculation(t *testing.T)
// Different resource types
// Verify size estimates reasonable
```

### Statistics
```go
// Test: Hit rate tracking
func TestCache_HitRateTracking(t *testing.T)
// Mix of hits and misses
// Stats reflect reality

// Test: Item count accuracy
func TestCache_ItemCount(t *testing.T)
// Add/remove items
// Count always correct

// Test: Stats after clear
func TestCache_StatsAfterClear(t *testing.T)
// Clear cache
// Stats reset appropriately
```

### Performance
```go
// Benchmark: Get performance
func BenchmarkCache_Get(b *testing.B)
// Single-threaded get performance

// Benchmark: Concurrent get
func BenchmarkCache_ConcurrentGet(b *testing.B)
// Multi-threaded get performance

// Benchmark: Eviction overhead
func BenchmarkCache_EvictionOverhead(b *testing.B)
// Measure eviction cost

// Test: Large item count
func TestCache_LargeItemCount(t *testing.T)
// 100,000 small items
// Operations still fast
```

## Implementation Notes

- Use `-race` flag for all concurrent tests
- Create helper to measure actual memory usage
- Test with various item sizes (small, medium, large)
- Consider testing GC interaction

## Output Files

- `internal/resourcecache/cache_edge_test.go`
