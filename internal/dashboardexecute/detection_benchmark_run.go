package dashboardexecute

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/turbot/pipe-fittings/steampipeconfig"
	"github.com/turbot/powerpipe/internal/dashboardtypes"
	"github.com/turbot/powerpipe/internal/resources"
)

// DetectionBenchmarkRun is a struct representing a container run
type DetectionBenchmarkRun struct {
	DashboardParentImpl
	BenchmarkType string `json:"benchmark_type"`

	// TODO KAI WHICH???
	dashboardNode *resources.DetectionBenchmark
	dashboardNode *resources.Benchmark
}

func (r *DetectionBenchmarkRun) AsTreeNode() *steampipeconfig.SnapshotTreeNode {
	res := &steampipeconfig.SnapshotTreeNode{
		Name:     r.Name,
		NodeType: r.NodeType,
		Children: make([]*steampipeconfig.SnapshotTreeNode, len(r.children)),
	}
	for i, c := range r.children {
		res.Children[i] = c.AsTreeNode()
	}
	return res
}

func NewDetectionBenchmarkRun(container *resources.Benchmark, parent dashboardtypes.DashboardParent, executionTree *DashboardExecutionTree) (*DetectionBenchmarkRun, error) {
	children := container.GetChildren()

	r := &DetectionBenchmarkRun{
		dashboardNode: container,
		BenchmarkType: "detection",
	}
	// create NewDashboardTreeRunImpl
	// (we must create after creating the run as it requires a ref to the run)
	r.DashboardParentImpl = newDashboardParentImpl(container, parent, r, executionTree)

	if container.Title != nil {
		r.Title = *container.Title
	}

	if container.Width != nil {
		r.Width = *container.Width
	}
	r.childCompleteChan = make(chan dashboardtypes.DashboardTreeRun, len(children))
	for _, child := range children {
		var childRun dashboardtypes.DashboardTreeRun
		var err error
		switch i := child.(type) {
		case *resources.Benchmark:
			childRun, err = NewDetectionBenchmarkRun(i, r, executionTree)
			if err != nil {
				return nil, err
			}
		case *resources.Detection:
			childRun, err = NewDetectionRun(i, r, executionTree)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("invalid child type %T", i)
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
	// add r into execution tree
	executionTree.runs[r.Name] = r
	return r, nil
}

// Initialise implements DashboardTreeRun
func (r *DetectionBenchmarkRun) Initialise(ctx context.Context) {
	// initialise our children
	if err := r.initialiseChildren(ctx); err != nil {
		r.SetError(ctx, err)
	}
}

// Execute implements DashboardTreeRun
// execute all children and wait for them to complete
func (r *DetectionBenchmarkRun) Execute(ctx context.Context) {
	// execute all children asynchronously
	r.executeChildrenAsync(ctx)

	// try to set status as running (will be set to blocked if any children are blocked)
	r.setRunning(ctx)

	// wait for children to complete
	err := <-r.waitForChildrenAsync(ctx)
	if err == nil {
		slog.Debug("Execute waitForChildrenAsync returned success", "name", r.Name)
		// set complete status on dashboard
		r.SetComplete(ctx)
	} else {
		slog.Debug("Execute waitForChildrenAsync failed", "name", r.Name, "error", err.Error())
		r.SetError(ctx, err)
	}
}

// IsSnapshotPanel implements SnapshotPanel
func (*DetectionBenchmarkRun) IsSnapshotPanel() {}
