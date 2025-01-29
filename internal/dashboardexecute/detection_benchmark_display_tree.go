package dashboardexecute

import (
	"time"

	"github.com/turbot/powerpipe/internal/controlexecute"

	"github.com/turbot/powerpipe/internal/controlstatus"
	"github.com/turbot/powerpipe/internal/workspace"
)

// DetectionBenchmarkDisplayTree is a structure representing the control execution hierarchy
type DetectionBenchmarkDisplayTree struct {
	Root *DetectionBenchmarkDisplay `json:"root"`
	// map of all leaf runs, keyed by FULL name
	LeafRuns  map[string]controlexecute.LeafRun `json:"-"`
	StartTime time.Time                         `json:"start_time"`
	EndTime   time.Time                         `json:"end_time"`
	Progress  *controlstatus.ControlProgress    `json:"progress"`
	// map of dimension property name to property value to color map
	DimensionColorGenerator *controlexecute.DimensionColorGenerator `json:"-"`
	// the current session search path
	//SearchPath []string                      `json:"-"`
	Workspace *workspace.PowerpipeWorkspace `json:"-"`
	// TODO K for csv only?
	// TODO need DetectionRunInstances - we need to replicate the control execution tree
	// DetectionRunInstances is a list of detection runs for each parent.
	// DetectionRunInstances []*DetectionRunInstance `json:"-"`
	// an optional map of control names used to filter the controls which are run
	//controlNameFilterMap map[string]struct{}
	// for now just using DetectionRuns
	DetectionRuns []*DetectionRun `json:"-"`
}

// IsExportSourceData implements ExportSourceData
func (*DetectionBenchmarkDisplayTree) IsExportSourceData() {}
