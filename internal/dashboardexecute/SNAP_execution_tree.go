package dashboardexecute

import (
	"github.com/turbot/powerpipe/internal/controlexecute"
	"time"

	"github.com/turbot/powerpipe/internal/controlstatus"
	"github.com/turbot/powerpipe/internal/workspace"
)

// DisplayExecutionTree_SNAP is a structure representing the control execution hierarchy
type DisplayExecutionTree_SNAP struct {
	Root *ResultGroup_SNAP `json:"root"`
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
	// TODO for csv only?
	//// ControlRunInstances is a list of control runs for each parent.
	//ControlRunInstances []*ControlRunInstance `json:"-"`
	//client              *db_client.DbClient
	// an optional map of control names used to filter the controls which are run
	//controlNameFilterMap map[string]struct{}
}

//
//func NewExecutionTree_SNAP(ctx context.Context, w *workspace.PowerpipeWorkspace, targets ...modconfig.ModTreeItem) (*DisplayExecutionTree_SNAP, error) {
//	// now populate the ExecutionTree
//	executionTree := &DisplayExecutionTree_SNAP{
//		Workspace: w,
//		LeafRuns:  make(map[string]LeafRun),
//	}
//
//	var resolvedItem modconfig.ModTreeItem
//	// if only one argument is provided, add this as execution root
//	if len(targets) == 1 {
//		resolvedItem = targets[0]
//	} else {
//		// create a root benchmark with `items` as it's children
//		resolvedItem = resources.NewRootBenchmarkWithChildren(w.Mod, targets).(modconfig.ModTreeItem)
//	}
//
//	// build tree of result groups, starting with a synthetic 'root' node
//
//	root, err := NewRootResultGroup_(ctx, executionTree, resolvedItem)
//	if err != nil {
//		return nil, err
//	}
//	executionTree.Root = root
//
//	// after tree has built, ControlCount will be set - create progress rendered
//	executionTree.Progress = controlstatus.NewControlProgress(len(executionTree.LeafRuns))
//
//	return executionTree, nil
//}

// TODO - for csv only?
//
//// PopulateControlRunInstances creates a list of ControlRunInstances, by expanding the list of control runs for each parent.
//func (e *DisplayExecutionTree_SNAP) PopulateControlRunInstances() {
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
func (*DisplayExecutionTree_SNAP) IsExportSourceData() {}
