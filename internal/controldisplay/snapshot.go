package controldisplay

import (
	"context"
	"fmt"
	"github.com/turbot/pipe-fittings/steampipeconfig"

	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/cloud"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/statushooks"
	"github.com/turbot/powerpipe/internal/controlexecute"
	"github.com/turbot/powerpipe/internal/dashboardexecute"
	"github.com/turbot/powerpipe/internal/dashboardworkspace"
)

func executionTreeToSnapshot(e *controlexecute.ExecutionTree) (*steampipeconfig.SteampipeSnapshot, error) {
	var dashboardNode modconfig.DashboardLeafNode
	var panels map[string]steampipeconfig.SnapshotPanel
	var checkRun *dashboardexecute.CheckRun

	// get root benchmark/control
	switch root := e.Root.Children[0].(type) {
	case *controlexecute.ResultGroup:
		var ok bool
		dashboardNode, ok = root.GroupItem.(modconfig.DashboardLeafNode)
		if !ok {
			return nil, fmt.Errorf("invalid node found in control execution tree - cannot cast '%s' to a DashboardLeafNode", root.GroupItem.Name())
		}
	case *controlexecute.ControlRun:
		dashboardNode = root.Control
	}

	// TACTICAL create a check run to wrap the execution tree
	checkRun = &dashboardexecute.CheckRun{Root: e.Root.Children[0]}
	checkRun.DashboardTreeRunImpl = dashboardexecute.NewDashboardTreeRunImpl(dashboardNode, nil, checkRun, nil)

	// populate the panels
	panels = checkRun.BuildSnapshotPanels(make(map[string]steampipeconfig.SnapshotPanel))

	// create the snapshot
	res := &steampipeconfig.SteampipeSnapshot{
		SchemaVersion: fmt.Sprintf("%d", steampipeconfig.SteampipeSnapshotSchemaVersion),
		Panels:        panels,
		Layout:        checkRun.Root.AsTreeNode(),
		Inputs:        map[string]interface{}{},
		Variables:     dashboardexecute.GetReferencedVariables(checkRun, dashboardworkspace.NewWorkspaceEvents(e.Workspace)),
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

	message, err := cloud.PublishSnapshot(ctx, snapshot, shouldShare)
	if err != nil {
		return err
	}
	if viper.GetBool(constants.ArgProgress) {
		statushooks.Done(ctx)
		fmt.Println(message) //nolint:forbidigo // acceptable
	}
	return nil

}
