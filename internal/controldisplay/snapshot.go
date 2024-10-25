package controldisplay

import (
	"context"
	"fmt"
	"github.com/turbot/pipe-fittings/modconfig/powerpipe"

	"github.com/turbot/pipe-fittings/pipes"
	"github.com/turbot/pipe-fittings/statushooks"
	"github.com/turbot/pipe-fittings/steampipeconfig"
	"github.com/turbot/powerpipe/internal/controlexecute"
	"github.com/turbot/powerpipe/internal/dashboardexecute"
)

func executionTreeToSnapshot(e *controlexecute.ExecutionTree) (*steampipeconfig.SteampipeSnapshot, error) {
	var dashboardNode powerpipe.DashboardLeafNode
	var panels map[string]steampipeconfig.SnapshotPanel
	var checkRun *dashboardexecute.CheckRun

	// get root benchmark/control
	switch root := e.Root.Children[0].(type) {
	case *controlexecute.ResultGroup:
		var ok bool
		dashboardNode, ok = root.GroupItem.(powerpipe.DashboardLeafNode)
		if !ok {
			return nil, fmt.Errorf("invalid node found in control execution tree - cannot cast '%s' to a DashboardLeafNode", root.GroupItem.Name())
		}
	case *controlexecute.ControlRun:
		dashboardNode = root.Control
	}

	// TACTICAL create a check run to wrap the execution tree
	checkRun = &dashboardexecute.CheckRun{
		Root:    e.Root.Children[0],
		Summary: e.Root.Summary,
	}
	checkRun.DashboardTreeRunImpl = dashboardexecute.NewDashboardTreeRunImpl(dashboardNode, nil, checkRun, nil)

	// populate the panels
	panels = checkRun.BuildSnapshotPanels(make(map[string]steampipeconfig.SnapshotPanel))

	vars, err := dashboardexecute.GetReferencedVariables(checkRun, e.Workspace)
	if err != nil {
		return nil, err
	}
	// create the snapshot
	res := &steampipeconfig.SteampipeSnapshot{
		SchemaVersion: fmt.Sprintf("%d", steampipeconfig.SteampipeSnapshotSchemaVersion),
		Panels:        panels,
		Layout:        checkRun.Root.AsTreeNode(),
		Inputs:        map[string]interface{}{},
		Variables:     vars,
		SearchPath:    e.SearchPath,
		StartTime:     e.StartTime,
		EndTime:       e.EndTime,
		Title:         dashboardNode.GetTitle(),
		FileNameRoot:  dashboardNode.Name(),
	}
	return res, nil
}

func PublishSnapshot(ctx context.Context, e *controlexecute.ExecutionTree, shouldShare bool) error {
	statushooks.Show(ctx)
	defer statushooks.Done(ctx)

	snapshot, err := executionTreeToSnapshot(e)
	if err != nil {
		return err
	}

	message, err := pipes.PublishSnapshot(ctx, snapshot, shouldShare)
	if err != nil {
		return err
	}
	statushooks.Done(ctx)
	fmt.Println(message) //nolint:forbidigo // acceptable
	return nil

}
