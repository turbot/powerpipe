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

	// Enhanced metadata fields
	Category      string `json:"category,omitempty"`      // Category for grouping
	Documentation string `json:"documentation,omitempty"` // URL to documentation
	Display       string `json:"display,omitempty"`       // Display type hint
	Width         *int   `json:"width,omitempty"`         // Dashboard width

	// Resolution tracking - indicates if values are fully resolved or need lazy resolution
	TitleResolved       bool     `json:"title_resolved,omitempty"`       // True if title is fully resolved (literal value)
	DescriptionResolved bool     `json:"description_resolved,omitempty"` // True if description is fully resolved
	TagsResolved        bool     `json:"tags_resolved,omitempty"`        // True if all tag values are fully resolved
	UnresolvedRefs      []string `json:"unresolved_refs,omitempty"`      // List of fields/keys that need resolution

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
	size := 120 // base overhead for struct (increased for new fields)
	size += len(e.Type) + len(e.Name) + len(e.ShortName)
	size += len(e.Title) + len(e.Description)
	size += len(e.Category) + len(e.Documentation) + len(e.Display)
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
	for _, r := range e.UnresolvedRefs {
		size += len(r)
	}
	return size
}

// NeedsResolution returns true if any metadata needs semantic resolution.
// This is used to determine if the entry should be queued for background resolution.
func (e *IndexEntry) NeedsResolution() bool {
	return !e.TitleResolved || !e.DescriptionResolved || !e.TagsResolved
}

// GetUnresolvedFields returns a list of field names that need resolution.
func (e *IndexEntry) GetUnresolvedFields() []string {
	fields := []string{}
	if !e.TitleResolved && e.Title == "" {
		fields = append(fields, "title")
	}
	if !e.DescriptionResolved && e.Description == "" {
		fields = append(fields, "description")
	}
	if !e.TagsResolved {
		fields = append(fields, "tags")
	}
	return fields
}

// IsFullyResolved returns true if all metadata is resolved and no lazy resolution is needed.
func (e *IndexEntry) IsFullyResolved() bool {
	return e.TitleResolved && e.DescriptionResolved && e.TagsResolved
}
