package resourcecache

import (
	"github.com/turbot/pipe-fittings/v2/modconfig"
)

// ResourceCache specializes Cache for HCL resources with type-safe access methods
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
	if res, ok := val.(modconfig.HclResource); ok {
		return res, true
	}
	return nil, false
}

// PutResource caches a parsed resource
func (rc *ResourceCache) PutResource(name string, resource modconfig.HclResource) {
	rc.cache.Put(name, resource)
}

// Get retrieves any cached value by key (generic access)
func (rc *ResourceCache) Get(name string) (interface{}, bool) {
	return rc.cache.Get(name)
}

// Put caches any value by key (generic access)
func (rc *ResourceCache) Put(name string, value interface{}) {
	rc.cache.Put(name, value)
}

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

// Len returns the number of items in the cache
func (rc *ResourceCache) Len() int {
	return rc.cache.Len()
}

// Keys returns all keys currently in the cache
func (rc *ResourceCache) Keys() []string {
	return rc.cache.Keys()
}
