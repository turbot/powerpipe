package display

import (
	"encoding/json"
	"strings"

	"github.com/turbot/pipe-fittings/v2/printers"
	"github.com/turbot/powerpipe/internal/resourceindex"
)

// ListableIndexEntry wraps an IndexEntry to implement printers.Listable interface.
// This allows index entries to be displayed in list commands without loading full resources.
type ListableIndexEntry struct {
	entry *resourceindex.IndexEntry
}

// NewListableIndexEntry creates a ListableIndexEntry from an IndexEntry.
func NewListableIndexEntry(entry *resourceindex.IndexEntry) *ListableIndexEntry {
	return &ListableIndexEntry{entry: entry}
}

// GetListData implements printers.Listable
func (l *ListableIndexEntry) GetListData() *printers.RowData {
	res := printers.NewRowData()

	// Add MOD field
	res.AddField(printers.NewFieldValue("MOD", l.entry.ModName))

	// Add NAME field - use short name for local resources, full name for dependencies
	name := l.entry.ShortName
	if l.IsDependencyResource() {
		name = l.entry.Name
	}
	res.AddField(printers.NewFieldValue("NAME", name))

	// Add TYPE field for benchmarks (matches Benchmark.GetListData behavior)
	if l.entry.Type == "benchmark" || l.entry.Type == "detection_benchmark" {
		benchmarkType := "control"
		if l.entry.BenchmarkType == "detection" || l.entry.Type == "detection_benchmark" {
			benchmarkType = "detection"
		}
		res.AddField(printers.NewFieldValue("TYPE", benchmarkType))
	}

	return res
}

// MarshalJSON implements json.Marshaler to provide the expected JSON output format.
func (l *ListableIndexEntry) MarshalJSON() ([]byte, error) {
	// Build a map with the fields expected in JSON output
	data := map[string]interface{}{
		"qualified_name": l.entry.Name,
		"mod_name":       l.entry.ModName,
		"file_name":      l.entry.FileName,
		"is_anonymous":   false,
	}

	// Add optional fields if present
	if l.entry.Title != "" {
		data["title"] = l.entry.Title
	}
	if l.entry.Description != "" {
		data["description"] = l.entry.Description
	}
	if len(l.entry.Tags) > 0 {
		data["tags"] = l.entry.Tags
	}
	// Note: documentation in IndexEntry is a URL field, not the full doc text
	// The full documentation text is only available via show/on-demand loading

	// Add line numbers if available
	if l.entry.StartLine > 0 {
		data["start_line_number"] = l.entry.StartLine
	}
	if l.entry.EndLine > 0 {
		data["end_line_number"] = l.entry.EndLine
	}

	// Add children for benchmarks
	if len(l.entry.ChildNames) > 0 {
		data["children"] = l.entry.ChildNames
	}

	// Add benchmark type if present
	if l.entry.Type == "benchmark" || l.entry.Type == "detection_benchmark" {
		benchmarkType := "control"
		if l.entry.BenchmarkType == "detection" || l.entry.Type == "detection_benchmark" {
			benchmarkType = "detection"
		}
		data["type"] = benchmarkType
	}

	return json.Marshal(data)
}

// IsDependencyResource returns true if this resource is from a dependency mod.
// For index entries, dependency resources have a ModFullName like "github.com/turbot/xxx@version"
// while local mods have ModFullName like "mod.aws_compliance".
func (l *ListableIndexEntry) IsDependencyResource() bool {
	// Local mods have ModFullName starting with "mod."
	// Dependencies have ModFullName like "github.com/turbot/xxx@version"
	if l.entry.ModFullName != "" && !strings.HasPrefix(l.entry.ModFullName, "mod.") {
		return true
	}
	return false
}

// Name returns the full resource name.
func (l *ListableIndexEntry) Name() string {
	return l.entry.Name
}

// GetTitle returns the title for the resource.
func (l *ListableIndexEntry) GetTitle() string {
	return l.entry.Title
}

// GetDescription returns the description for the resource.
func (l *ListableIndexEntry) GetDescription() string {
	return l.entry.Description
}

// GetTags returns the tags for the resource.
func (l *ListableIndexEntry) GetTags() map[string]string {
	return l.entry.Tags
}

// GetEntry returns the underlying IndexEntry.
func (l *ListableIndexEntry) GetEntry() *resourceindex.IndexEntry {
	return l.entry
}

// IsTopLevel returns true if this is a top-level resource.
func (l *ListableIndexEntry) IsTopLevel() bool {
	return l.entry.IsTopLevel
}

// GetType returns the resource type.
func (l *ListableIndexEntry) GetType() string {
	return l.entry.Type
}

// GetModName returns the mod name.
func (l *ListableIndexEntry) GetModName() string {
	return l.entry.ModName
}
