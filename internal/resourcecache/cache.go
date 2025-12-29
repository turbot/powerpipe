// Package resourcecache provides a thread-safe LRU cache with memory-based eviction
// for caching parsed HCL resources.
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

// DefaultConfig returns sensible defaults for the cache
func DefaultConfig() CacheConfig {
	return CacheConfig{
		MaxMemoryBytes: 50 * 1024 * 1024, // 50MB
		MaxEntries:     0,                 // Memory-based eviction only
		TTL:            0,                 // No expiry
	}
}

// Sizer interface for items that can report their memory size
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

// New creates a new LRU cache with the given configuration
func New(config CacheConfig) *Cache {
	return &Cache{
		list:   list.New(),
		items:  make(map[string]*list.Element),
		config: config,
	}
}

// Get retrieves an item from the cache.
// Returns the value and true if found, nil and false otherwise.
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

// Put adds an item to the cache, evicting least recently used items if necessary
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

// Len returns the number of items in the cache
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.list.Len()
}

// evictIfNeeded removes least recently used items until under limits
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
		Entries:     c.list.Len(),
		MemoryBytes: c.currentMemory,
		MaxMemory:   c.config.MaxMemoryBytes,
		Hits:        c.hits,
		Misses:      c.misses,
		Evictions:   c.evictions,
		HitRate:     c.hitRate(),
	}
}

func (c *Cache) hitRate() float64 {
	total := c.hits + c.misses
	if total == 0 {
		return 0
	}
	return float64(c.hits) / float64(total)
}

// CacheStats holds cache performance statistics
type CacheStats struct {
	Entries     int     // Number of entries in cache
	MemoryBytes int64   // Current memory usage
	MaxMemory   int64   // Maximum memory limit
	Hits        int64   // Number of cache hits
	Misses      int64   // Number of cache misses
	Evictions   int64   // Number of evicted entries
	HitRate     float64 // Hit rate (0.0 - 1.0)
}

// Keys returns all keys currently in the cache (for iteration/debugging)
func (c *Cache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.items))
	for key := range c.items {
		keys = append(keys, key)
	}
	return keys
}
