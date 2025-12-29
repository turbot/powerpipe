package resourceindex

// IndexEntry contains minimal metadata about a resource for lazy loading.
// This provides enough information for UI lists (available_dashboards)
// without needing to parse and load full HCL resources.
type IndexEntry struct {
	// Identity
	Type      string `json:"type"`       // "dashboard", "query", "control", etc.
	Name      string `json:"name"`       // Full name: "mod.type.shortname"
	ShortName string `json:"short_name"` // Just the short name

	// Display metadata (for UI lists)
	Title       string            `json:"title,omitempty"`
	Description string            `json:"description,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`

	// Hierarchy (for benchmarks/dashboards)
	IsTopLevel bool     `json:"is_top_level,omitempty"`
	ParentName string   `json:"parent_name,omitempty"`
	ChildNames []string `json:"child_names,omitempty"`

	// Source location (for on-demand loading)
	FileName   string `json:"file_name"`
	StartLine  int    `json:"start_line"`
	EndLine    int    `json:"end_line"`
	ByteOffset int64  `json:"byte_offset"` // For efficient seeking
	ByteLength int    `json:"byte_length"`

	// Mod info
	ModName     string `json:"mod_name"`
	ModFullName string `json:"mod_full_name"`
	ModRoot     string `json:"mod_root,omitempty"` // Root directory of the mod (for file() function)

	// Type-specific metadata
	// For benchmarks
	BenchmarkType string `json:"benchmark_type,omitempty"` // "control" or "detection"

	// For queries/controls
	HasSQL   bool   `json:"has_sql,omitempty"`
	QueryRef string `json:"query_ref,omitempty"` // For controls/cards referencing a query

	// For inputs
	DashboardName string   `json:"dashboard_name,omitempty"` // For scoped inputs
	InputNames    []string `json:"input_names,omitempty"`    // For dashboards with scoped inputs
}

// Size returns approximate memory size of this entry in bytes.
// Used for tracking index memory usage.
func (e *IndexEntry) Size() int {
	// Rough estimate: strings + overhead
	size := 100 // base overhead for struct
	size += len(e.Type) + len(e.Name) + len(e.ShortName)
	size += len(e.Title) + len(e.Description)
	size += len(e.FileName) + len(e.ModName) + len(e.ModFullName)
	size += len(e.ParentName) + len(e.BenchmarkType) + len(e.DashboardName)
	size += len(e.QueryRef)
	for k, v := range e.Tags {
		size += len(k) + len(v)
	}
	for _, c := range e.ChildNames {
		size += len(c)
	}
	for _, c := range e.InputNames {
		size += len(c)
	}
	return size
}
