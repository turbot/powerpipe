package resourceindex

import (
	"sync"
)

// ResourceIndex provides fast lookup of resource metadata for lazy loading.
// It stores minimal information about all resources without requiring full HCL parsing.
type ResourceIndex struct {
	mu sync.RWMutex

	// All entries keyed by full name
	entries map[string]*IndexEntry

	// Type-specific indexes for efficient filtering
	byType map[string]map[string]*IndexEntry // type -> name -> entry

	// Mod information
	ModName     string
	ModFullName string
	ModTitle    string

	// Mod name mappings: full path -> short name
	// e.g., "github.com/turbot/steampipe-mod-aws-insights" -> "aws_insights"
	modNameMap map[string]string

	// Mod title mappings: full path -> title
	// e.g., "github.com/turbot/steampipe-mod-aws-insights" -> "AWS Insights"
	modTitleMap map[string]string

	// Statistics
	totalSize int
}

// NewResourceIndex creates an empty index.
func NewResourceIndex() *ResourceIndex {
	return &ResourceIndex{
		entries:     make(map[string]*IndexEntry),
		byType:      make(map[string]map[string]*IndexEntry),
		modNameMap:  make(map[string]string),
		modTitleMap: make(map[string]string),
	}
}

// RegisterModName registers a mapping from a full mod path to its short name.
// e.g., "github.com/turbot/steampipe-mod-aws-insights" -> "aws_insights"
func (idx *ResourceIndex) RegisterModName(fullPath, shortName string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	idx.modNameMap[fullPath] = shortName
}

// RegisterModTitle registers a mapping from a full mod path to its title.
// e.g., "github.com/turbot/steampipe-mod-aws-insights" -> "AWS Insights"
func (idx *ResourceIndex) RegisterModTitle(fullPath, title string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	idx.modTitleMap[fullPath] = title
}

// ResolveModName converts a full mod path to its short name.
// If no mapping exists, returns the input unchanged.
func (idx *ResourceIndex) ResolveModName(modName string) string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	if shortName, ok := idx.modNameMap[modName]; ok {
		return shortName
	}
	return modName
}

// GetModNameMap returns a copy of the mod name mappings.
// This is used by lazy workspace to build dependency mod metadata.
func (idx *ResourceIndex) GetModNameMap() map[string]string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	// Return a copy to prevent concurrent modification
	result := make(map[string]string, len(idx.modNameMap))
	for k, v := range idx.modNameMap {
		result[k] = v
	}
	return result
}

// GetModTitleMap returns a copy of the mod title mappings.
// This is used by lazy workspace to build dependency mod metadata.
func (idx *ResourceIndex) GetModTitleMap() map[string]string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	// Return a copy to prevent concurrent modification
	result := make(map[string]string, len(idx.modTitleMap))
	for k, v := range idx.modTitleMap {
		result[k] = v
	}
	return result
}

// Add adds an entry to the index.
func (idx *ResourceIndex) Add(entry *IndexEntry) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	idx.entries[entry.Name] = entry

	// Add to type index
	if idx.byType[entry.Type] == nil {
		idx.byType[entry.Type] = make(map[string]*IndexEntry)
	}
	idx.byType[entry.Type][entry.Name] = entry

	idx.totalSize += entry.Size()
}

// Get retrieves an entry by full name.
func (idx *ResourceIndex) Get(name string) (*IndexEntry, bool) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	entry, ok := idx.entries[name]
	return entry, ok
}

// GetByType retrieves all entries of a specific type.
func (idx *ResourceIndex) GetByType(resourceType string) []*IndexEntry {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	typeMap := idx.byType[resourceType]
	if typeMap == nil {
		return nil
	}

	entries := make([]*IndexEntry, 0, len(typeMap))
	for _, entry := range typeMap {
		entries = append(entries, entry)
	}
	return entries
}

// List returns all entries in the index.
func (idx *ResourceIndex) List() []*IndexEntry {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	entries := make([]*IndexEntry, 0, len(idx.entries))
	for _, entry := range idx.entries {
		entries = append(entries, entry)
	}
	return entries
}

// Remove removes an entry from the index by name.
func (idx *ResourceIndex) Remove(name string) bool {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	entry, ok := idx.entries[name]
	if !ok {
		return false
	}

	// Remove from type index
	if typeMap := idx.byType[entry.Type]; typeMap != nil {
		delete(typeMap, name)
	}

	// Update size
	idx.totalSize -= entry.Size()

	// Remove from main entries
	delete(idx.entries, name)
	return true
}

// UpdateEntry updates an existing entry in the index.
// This is used by background resolution to update metadata after resolving variables.
// The entry's Size is recalculated after the update.
func (idx *ResourceIndex) UpdateEntry(entry *IndexEntry) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	oldEntry, ok := idx.entries[entry.Name]
	if !ok {
		// Entry doesn't exist, add it instead
		idx.entries[entry.Name] = entry
		if idx.byType[entry.Type] == nil {
			idx.byType[entry.Type] = make(map[string]*IndexEntry)
		}
		idx.byType[entry.Type][entry.Name] = entry
		idx.totalSize += entry.Size()
		return
	}

	// Update total size (subtract old, add new)
	idx.totalSize -= oldEntry.Size()
	idx.totalSize += entry.Size()

	// The entry pointer should be the same (we're updating in place),
	// but update the maps just in case
	idx.entries[entry.Name] = entry
	idx.byType[entry.Type][entry.Name] = entry
}

// Dashboards returns all dashboard entries.
func (idx *ResourceIndex) Dashboards() []*IndexEntry {
	return idx.GetByType("dashboard")
}

// Benchmarks returns all benchmark entries (control and detection).
func (idx *ResourceIndex) Benchmarks() []*IndexEntry {
	controlBenchmarks := idx.GetByType("benchmark")
	detectionBenchmarks := idx.GetByType("detection_benchmark")
	return append(controlBenchmarks, detectionBenchmarks...)
}

// Queries returns all query entries.
func (idx *ResourceIndex) Queries() []*IndexEntry {
	return idx.GetByType("query")
}

// Controls returns all control entries.
func (idx *ResourceIndex) Controls() []*IndexEntry {
	return idx.GetByType("control")
}

// TopLevelBenchmarks returns benchmarks that are direct children of mod.
func (idx *ResourceIndex) TopLevelBenchmarks() []*IndexEntry {
	var result []*IndexEntry
	for _, entry := range idx.Benchmarks() {
		if entry.IsTopLevel {
			result = append(result, entry)
		}
	}
	return result
}

// GetChildren returns child entries for a parent.
func (idx *ResourceIndex) GetChildren(parentName string) []*IndexEntry {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	parent, ok := idx.entries[parentName]
	if !ok || len(parent.ChildNames) == 0 {
		return nil
	}

	children := make([]*IndexEntry, 0, len(parent.ChildNames))
	for _, childName := range parent.ChildNames {
		if child, ok := idx.entries[childName]; ok {
			children = append(children, child)
		}
	}
	return children
}

// Size returns total approximate memory size of index in bytes.
func (idx *ResourceIndex) Size() int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.totalSize
}

// Count returns total number of entries.
func (idx *ResourceIndex) Count() int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return len(idx.entries)
}

// Types returns a list of all resource types in the index.
func (idx *ResourceIndex) Types() []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	types := make([]string, 0, len(idx.byType))
	for typeName := range idx.byType {
		types = append(types, typeName)
	}
	return types
}

// Stats returns index statistics.
func (idx *ResourceIndex) Stats() IndexStats {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	stats := IndexStats{
		TotalEntries: len(idx.entries),
		TotalSize:    idx.totalSize,
		ByType:       make(map[string]int),
	}

	for typeName, entries := range idx.byType {
		stats.ByType[typeName] = len(entries)
	}

	return stats
}

// IndexStats contains statistics about the index.
type IndexStats struct {
	TotalEntries int
	TotalSize    int
	ByType       map[string]int
}
