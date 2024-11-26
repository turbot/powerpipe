package dashboardexecute

import (
	"github.com/turbot/powerpipe/internal/controlexecute"
	"time"

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
	// TODO need DetectionRunInstances
	//// ControlRunInstances is a list of control runs for each parent.
	//ControlRunInstances []*ControlRunInstance `json:"-"`
	// an optional map of control names used to filter the controls which are run
	//controlNameFilterMap map[string]struct{}
}

// TODO - for csv only?
//
//// PopulateControlRunInstances creates a list of ControlRunInstances, by expanding the list of control runs for each parent.
//func (e *DetectionBenchmarkDisplayTree) PopulateControlRunInstances() {
//	var controlRunInstances []*ControlRunInstance
//
//	for _, controlRun := range e.LeafRuns {
//		for _, parent := range controlRun.Parents {
//			flatControlRun := NewControlRunInstance(controlRun, parent)
//			controlRunInstances = append(controlRunInstances, &flatControlRun)
//		}
//	}
//
//	e.ControlRunInstances = controlRunInstances
//}

// IsExportSourceData implements ExportSourceData
func (*DetectionBenchmarkDisplayTree) IsExportSourceData() {}
