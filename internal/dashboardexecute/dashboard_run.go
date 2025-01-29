package dashboardexecute

import (
	"context"
	"fmt"
	"github.com/turbot/powerpipe/internal/resources"
	"log/slog"

	"github.com/turbot/pipe-fittings/v2/schema"
	"github.com/turbot/pipe-fittings/v2/steampipeconfig"
	"github.com/turbot/powerpipe/internal/dashboardtypes"
)

// DashboardRun is a struct representing a container run
type DashboardRun struct {
	runtimeDependencyPublisherImpl

	parent    dashboardtypes.DashboardParent
	dashboard *resources.Dashboard
}

func (r *DashboardRun) AsTreeNode() *steampipeconfig.SnapshotTreeNode {
	res := &steampipeconfig.SnapshotTreeNode{
		Name:     r.Name,
		NodeType: r.NodeType,
		Children: make([]*steampipeconfig.SnapshotTreeNode, 0, len(r.children)),
	}

	for _, c := range r.children {
		// NOTE: exclude with runs
		if c.GetNodeType() != schema.BlockTypeWith {
			res.Children = append(res.Children, c.AsTreeNode())
		}
	}

	return res
}

func NewDashboardRun(dashboard *resources.Dashboard, parent dashboardtypes.DashboardParent, executionTree *DashboardExecutionTree) (*DashboardRun, error) {
	r := &DashboardRun{
		parent:    parent,
		dashboard: dashboard,
	}
	// create RuntimeDependencyPublisherImpl- this handles 'with' run creation and resolving runtime dependency resolution
	// (we must create after creating the run as it requires a ref to the run)
	r.runtimeDependencyPublisherImpl = newRuntimeDependencyPublisherImpl(dashboard, parent, r, executionTree)
	// add r into execution tree BEFORE creating child runs or initialising runtime depdencies
	// - this is so child runs can find this dashboard run
	executionTree.runs[r.Name] = r

	// set inputs map on RuntimeDependencyPublisherImpl BEFORE creating child runs
	r.inputs = dashboard.GetInputs()

	// after setting inputs, init runtime dependencies. this creates with runs and adds them to our children
	err := r.initWiths()
	if err != nil {
		return nil, err
	}

	err = r.createChildRuns(executionTree)
	if err != nil {
		return nil, err
	}

	// create buffered channel for children to report their completion
	r.createChildCompleteChan()

	return r, nil
}

// Initialise implements DashboardTreeRun
func (r *DashboardRun) Initialise(ctx context.Context) {
	// initialise our children
	if err := r.initialiseChildren(ctx); err != nil {
		r.SetError(ctx, err)
	}
}

// Execute implements DashboardTreeRun
// execute all children and wait for them to complete
func (r *DashboardRun) Execute(ctx context.Context) {
	r.executeChildrenAsync(ctx)

	// try to set status as running (will be set to blocked if any children are blocked)
	r.setRunning(ctx)

	// wait for children to complete
	err := <-r.waitForChildrenAsync(ctx)
	if err == nil {
		slog.Debug("DashboardRun all children complete, success", "name", r.Name)
		// set complete status on dashboard
		r.SetComplete(ctx)
	} else {
		slog.Debug("DashboardRun all children complete, error", "name", r.Name, "error", err.Error())
		r.SetError(ctx, err)
	}
}

// IsSnapshotPanel implements SnapshotPanel
func (*DashboardRun) IsSnapshotPanel() {}

// GetInput searches for an input with the given name
func (r *DashboardRun) GetInput(name string) (*resources.DashboardInput, bool) {
	return r.dashboard.GetInput(name)
}

// GetInputsDependingOn returns a list o DashboardInputs which have a runtime dependency on the given input
func (r *DashboardRun) GetInputsDependingOn(changedInputName string) []string {
	var res []string
	for _, input := range r.dashboard.Inputs {
		if input.DependsOnInput(changedInputName) {
			res = append(res, input.UnqualifiedName)
		}
	}
	return res
}

func (r *DashboardRun) createChildRuns(executionTree *DashboardExecutionTree) error {
	// ask our resource for its children
	children := r.dashboard.GetChildren()

	for _, child := range children {
		var childRun dashboardtypes.DashboardTreeRun
		var err error
		switch i := child.(type) {
		case *resources.DashboardWith:
			// ignore as with runs are created by RuntimeDependencyPublisherImpl
			continue
		case *resources.Dashboard:
			childRun, err = NewDashboardRun(i, r, executionTree)
			if err != nil {
				return err
			}
		case *resources.DashboardContainer:
			childRun, err = NewDashboardContainerRun(i, r, executionTree)
			if err != nil {
				return err
			}
		case *resources.DetectionBenchmark:
			childRun, err = NewDetectionBenchmarkRun(i, r, executionTree)
			if err != nil {
				return err
			}

		case *resources.Benchmark, *resources.Control:
			childRun, err = NewCheckRun(i.(resources.DashboardLeafNode), r, executionTree)
			if err != nil {
				return err
			}
		case *resources.DashboardInput:
			// NOTE: clone the input to avoid mutating the original
			// TODO remove the need for this when we refactor input values resolution
			// TODO https://github.com/turbot/steampipe/issues/2864

			// TACTICAL: as this is a runtime dependency,  set the run name to the 'scoped name'
			// this is to match the name in the panel dependendencies
			inputRunName := fmt.Sprintf("%s.%s", r.DashboardName, i.UnqualifiedName)
			childRun, err = NewLeafRun(i.Clone(), r, executionTree, withName(inputRunName))
			if err != nil {
				return err
			}

		default:
			// ensure this item is a DashboardLeafNode
			leafNode, ok := i.(resources.DashboardLeafNode)
			if !ok {
				return fmt.Errorf("child %s does not implement DashboardLeafNode", i.Name())
			}

			childRun, err = NewLeafRun(leafNode, r, executionTree)
			if err != nil {
				return err
			}
		}

		// should never happen - container children must be either container or counter
		if childRun == nil {
			continue
		}

		// if our child has not completed, we have not completed
		if childRun.GetRunStatus() == dashboardtypes.RunInitialized {
			r.Status = dashboardtypes.RunInitialized
		}
		r.children = append(r.children, childRun)
	}
	return nil
}
