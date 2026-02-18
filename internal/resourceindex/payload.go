package resourceindex

// BuildAvailableDashboardsPayload builds the dashboard list payload
// without loading full resources. This enables the dashboard UI to
// show available dashboards and benchmarks with minimal memory usage.
func (idx *ResourceIndex) BuildAvailableDashboardsPayload() *AvailableDashboardsPayload {
	payload := &AvailableDashboardsPayload{
		Action:     "available_dashboards",
		Dashboards: make(map[string]DashboardInfo),
		Benchmarks: make(map[string]BenchmarkInfo),
	}

	// Build dashboard list from index
	for _, entry := range idx.Dashboards() {
		// Ensure mod tag is set for grouping
		tags := entry.Tags
		if tags == nil {
			tags = make(map[string]string)
		}
		// Add mod tag if not already present (mod_full_name without "github.com/" prefix)
		if _, exists := tags["mod"]; !exists && entry.ModFullName != "" {
			tags["mod"] = entry.ModFullName
		}

		payload.Dashboards[entry.Name] = DashboardInfo{
			Title:       entry.Title,
			FullName:    entry.Name,
			ShortName:   entry.ShortName,
			Tags:        tags,
			ModFullName: entry.ModFullName,
		}
	}

	// Build benchmark list with hierarchy from index
	benchmarkTrunks := make(map[string][][]string)

	for _, entry := range idx.Benchmarks() {
		// Ensure mod tag is set for grouping
		tags := entry.Tags
		if tags == nil {
			tags = make(map[string]string)
		}
		// Add mod tag if not already present
		if _, exists := tags["mod"]; !exists && entry.ModFullName != "" {
			tags["mod"] = entry.ModFullName
		}

		info := BenchmarkInfo{
			Title:         entry.Title,
			FullName:      entry.Name,
			ShortName:     entry.ShortName,
			BenchmarkType: entry.BenchmarkType,
			Tags:          tags,
			IsTopLevel:    entry.IsTopLevel,
			ModFullName:   entry.ModFullName,
		}

		// Top-level benchmarks have their own trunk
		if entry.IsTopLevel {
			benchmarkTrunks[entry.Name] = [][]string{{entry.Name}}
		}

		// Build children recursively from index with cycle detection
		visiting := make(map[string]bool)
		visiting[entry.Name] = true
		info.Children = idx.buildBenchmarkChildren(entry, entry.IsTopLevel,
			[]string{entry.Name}, benchmarkTrunks, visiting)

		payload.Benchmarks[entry.Name] = info
	}

	// Apply trunks
	for name, trunks := range benchmarkTrunks {
		if info, ok := payload.Benchmarks[name]; ok {
			info.Trunks = trunks
			payload.Benchmarks[name] = info
		}
	}

	return payload
}

func (idx *ResourceIndex) buildBenchmarkChildren(parent *IndexEntry,
	recordTrunk bool, trunk []string, trunks map[string][][]string, visiting map[string]bool) []BenchmarkInfo {

	var children []BenchmarkInfo

	for _, childEntry := range idx.GetChildren(parent.Name) {
		// Only include benchmark children (not controls)
		if childEntry.Type != "benchmark" && childEntry.Type != "detection_benchmark" {
			continue
		}

		// Cycle detection: skip if we're already visiting this node in the current path
		if visiting[childEntry.Name] {
			// Circular reference detected - skip to prevent infinite recursion
			continue
		}

		childTrunk := append([]string{}, trunk...)
		childTrunk = append(childTrunk, childEntry.Name)

		if recordTrunk {
			trunks[childEntry.Name] = append(trunks[childEntry.Name], childTrunk)
		}

		// Mark as visiting before recursion
		visiting[childEntry.Name] = true

		// Ensure mod tag is set for child benchmarks too
		childTags := childEntry.Tags
		if childTags == nil {
			childTags = make(map[string]string)
		}
		if _, exists := childTags["mod"]; !exists && childEntry.ModFullName != "" {
			childTags["mod"] = childEntry.ModFullName
		}

		info := BenchmarkInfo{
			Title:         childEntry.Title,
			FullName:      childEntry.Name,
			ShortName:     childEntry.ShortName,
			BenchmarkType: childEntry.BenchmarkType,
			Tags:          childTags,
			Children:      idx.buildBenchmarkChildren(childEntry, recordTrunk, childTrunk, trunks, visiting),
		}

		// Unmark after recursion completes (allows the same node to appear in different branches)
		delete(visiting, childEntry.Name)

		children = append(children, info)
	}

	return children
}

// AvailableDashboardsPayload is the payload sent to the dashboard UI
// containing available dashboards and benchmarks.
type AvailableDashboardsPayload struct {
	Action     string                   `json:"action"`
	Dashboards map[string]DashboardInfo `json:"dashboards"`
	Benchmarks map[string]BenchmarkInfo `json:"benchmarks"`
	Snapshots  map[string]string        `json:"snapshots,omitempty"`
}

// DashboardInfo contains dashboard metadata for the UI list.
type DashboardInfo struct {
	Title       string            `json:"title,omitempty"`
	FullName    string            `json:"full_name"`
	ShortName   string            `json:"short_name"`
	Tags        map[string]string `json:"tags,omitempty"`
	ModFullName string            `json:"mod_full_name,omitempty"`
	Database    string            `json:"database,omitempty"`
}

// BenchmarkInfo contains benchmark metadata for the UI list.
type BenchmarkInfo struct {
	Title         string            `json:"title,omitempty"`
	FullName      string            `json:"full_name"`
	ShortName     string            `json:"short_name"`
	BenchmarkType string            `json:"benchmark_type,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
	IsTopLevel    bool              `json:"is_top_level,omitempty"`
	Trunks        [][]string        `json:"trunks,omitempty"`
	Children      []BenchmarkInfo   `json:"children,omitempty"`
	ModFullName   string            `json:"mod_full_name,omitempty"`
}
