package display

import (
	"encoding/json"
	"os"
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

	// Add NAME field - always use full qualified name for consistency with v1.4.2
	res.AddField(printers.NewFieldValue("NAME", l.entry.Name))

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
		"resource_name":  l.entry.ShortName,
		"mod_name":       l.entry.ModName,
		"file_name":      l.entry.FileName,
		"is_anonymous":   false,
		"auto_generated": false, // User-defined resources are never auto-generated
	}

	// Add path - the hierarchical path to this resource
	// Format: [["mod.mod_name", "parent_name", ..., "qualified_name"]]
	// Controls can have multiple paths if they're children of multiple benchmarks
	// Use pre-computed paths if available
	if len(l.entry.Paths) > 0 {
		data["path"] = l.entry.Paths
	} else {
		data["path"] = l.buildAllPaths()
	}

	// Add url_path for dashboards
	if l.entry.Type == "dashboard" {
		data["url_path"] = "/" + l.entry.Name
	}

	// Add control-specific fields
	if l.entry.Type == "control" {
		data["args"] = map[string]interface{}{} // Controls always have args (empty if none)
		if l.entry.SQL != "" {
			data["sql"] = l.entry.SQL
		}
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

	// Add children for benchmarks (from explicit children attribute)
	if len(l.entry.ChildNames) > 0 {
		data["children"] = l.entry.ChildNames
	}

	// Add references for benchmarks
	if len(l.entry.References) > 0 {
		data["references"] = l.entry.References
	}

	// Add source_definition if we can read it from file
	if sourceDef := l.readSourceDefinition(); sourceDef != "" {
		data["source_definition"] = sourceDef
	}

	// Note: benchmark "type" field is not included in JSON output for list command
	// (eager mode doesn't include it either - it's only used for pretty/plain output)

	return json.Marshal(data)
}

// readSourceDefinition reads the source definition from the file using byte offsets.
func (l *ListableIndexEntry) readSourceDefinition() string {
	// Need byte offset and length to read source
	if l.entry.ByteOffset == 0 && l.entry.ByteLength == 0 {
		return ""
	}
	if l.entry.FileName == "" {
		return ""
	}

	// Read the source from file
	file, err := os.Open(l.entry.FileName)
	if err != nil {
		return ""
	}
	defer file.Close()

	// Seek to the block start
	_, err = file.Seek(l.entry.ByteOffset, 0)
	if err != nil {
		return ""
	}

	// Read the block content
	buf := make([]byte, l.entry.ByteLength)
	n, err := file.Read(buf)
	if err != nil || n == 0 {
		return ""
	}

	return strings.TrimSpace(string(buf[:n]))
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

// buildAllPaths constructs all hierarchical paths to this resource.
// Controls can be children of multiple benchmarks, so they may have multiple paths.
// Format: [["mod.mod_name", "parent_name", ..., "resource_name"], ...]
func (l *ListableIndexEntry) buildAllPaths() [][]string {
	// If there are multiple parents, build a path for each
	if len(l.entry.ParentNames) > 1 {
		paths := make([][]string, 0, len(l.entry.ParentNames))
		for _, parentName := range l.entry.ParentNames {
			path := []string{l.entry.ModFullName, parentName, l.entry.Name}
			paths = append(paths, path)
		}
		return paths
	}

	// Single path case
	path := []string{l.entry.ModFullName}

	// If this resource has a parent, include it in the path
	if l.entry.ParentName != "" {
		path = append(path, l.entry.ParentName)
	}

	// Add this resource's name
	path = append(path, l.entry.Name)

	return [][]string{path}
}
