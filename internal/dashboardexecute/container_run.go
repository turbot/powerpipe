package dashboardexecute

import (
	"context"
	"fmt"
	"github.com/turbot/powerpipe/internal/resources"
	"log/slog"

	"github.com/turbot/pipe-fittings/steampipeconfig"
	"github.com/turbot/powerpipe/internal/dashboardtypes"
)

// DashboardContainerRun is a struct representing a container run
type DashboardContainerRun struct {
	DashboardParentImpl

	dashboardNode *resources.DashboardContainer
}

func (r *DashboardContainerRun) AsTreeNode() *steampipeconfig.SnapshotTreeNode {
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

func NewDashboardContainerRun(container *resources.DashboardContainer, parent dashboardtypes.DashboardParent, executionTree *DashboardExecutionTree) (*DashboardContainerRun, error) {
	children := container.GetChildren()

	r := &DashboardContainerRun{dashboardNode: container}
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
		case *resources.DashboardContainer:
			childRun, err = NewDashboardContainerRun(i, r, executionTree)
			if err != nil {
				return nil, err
			}
		case *resources.Dashboard:
			childRun, err = NewDashboardRun(i, r, executionTree)
			if err != nil {
				return nil, err
			}
		case *resources.Benchmark, *resources.Control:
			childRun, err = NewCheckRun(i.(resources.DashboardLeafNode), r, executionTree)
			if err != nil {
				return nil, err
			}

		default:
			// ensure this item is a DashboardLeafNode
			leafNode, ok := i.(resources.DashboardLeafNode)
			if !ok {
				return nil, fmt.Errorf("child %s does not implement DashboardLeafNode", i.Name())
			}

			childRun, err = NewLeafRun(leafNode, r, executionTree)
			if err != nil {
				return nil, err
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
	// add r into execution tree
	executionTree.runs[r.Name] = r
	return r, nil
}

// Initialise implements DashboardTreeRun
func (r *DashboardContainerRun) Initialise(ctx context.Context) {
	// initialise our children
	if err := r.initialiseChildren(ctx); err != nil {
		r.SetError(ctx, err)
	}
}

// Execute implements DashboardTreeRun
// execute all children and wait for them to complete
func (r *DashboardContainerRun) Execute(ctx context.Context) {
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
func (*DashboardContainerRun) IsSnapshotPanel() {}
