package controldisplay

import (
	"context"
	"fmt"
	"github.com/turbot/pipe-fittings/modconfig"
	pworkspace "github.com/turbot/pipe-fittings/workspace"
	"github.com/turbot/powerpipe/internal/db_client"
	"github.com/turbot/powerpipe/internal/workspace"

	"github.com/turbot/pipe-fittings/pipes"
	"github.com/turbot/pipe-fittings/statushooks"
	"github.com/turbot/pipe-fittings/steampipeconfig"
	"github.com/turbot/powerpipe/internal/controlexecute"
	"github.com/turbot/powerpipe/internal/dashboardexecute"
	"github.com/turbot/powerpipe/internal/resources"
)

func executionTreeToSnapshot(e *controlexecute.ExecutionTree) (*steampipeconfig.SteampipeSnapshot, error) {
	var dashboardNode resources.DashboardLeafNode
	var panels map[string]steampipeconfig.SnapshotPanel
	var checkRun *dashboardexecute.CheckRun

	// get root benchmark/control
	switch root := e.Root.Children[0].(type) {
	case *controlexecute.ResultGroup:
		var ok bool
		dashboardNode, ok = root.GroupItem.(resources.DashboardLeafNode)
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
		// hardcode benchmark type for now
		BenchmarkType: "control",
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
func SnapshotToExecutionTree(ctx context.Context, snapshot *steampipeconfig.SteampipeSnapshot, w *workspace.PowerpipeWorkspace, client *db_client.DbClient, controlFilter pworkspace.ResourceFilter, targets ...modconfig.ModTreeItem) (*controlexecute.ExecutionTree, error) {
	// Step 1: Create the execution tree
	tree, err := controlexecute.NewExecutionTree(ctx, w, client, controlFilter, targets...)
	if err != nil {
		return nil, err
	}

	// Step 2: Populate summaries
	populateSummaries(tree.Root, snapshot.Layout)

	// Step 3: Populate results
	populateResults(tree.Root, snapshot.Panels)

	// Step 4: Add any additional metadata from the snapshot
	tree.SearchPath = snapshot.SearchPath
	tree.StartTime = snapshot.StartTime
	tree.EndTime = snapshot.EndTime

	return tree, nil
}

func populateSummaries(treeNode controlexecute.ExecutionTreeNode, snapshotNode *steampipeconfig.SnapshotTreeNode) {
	if treeNode == nil || snapshotNode == nil {
		return
	}

	//// Copy the summary (ensure treeNode has a Summary method)
	//if summaryProvider, ok := treeNode.(interface{ SetSummary(summary string) }); ok {
	//	summaryProvider.SetSummary(snapshotNode.Summary)
	//}

	// Recursively populate summaries for children
	treeChildren := treeNode.GetChildren()
	snapshotChildren := snapshotNode.Children
	for i := range treeChildren {
		if i < len(snapshotChildren) {
			populateSummaries(treeChildren[i], snapshotChildren[i])
		}
	}
}

func populateResults(treeNode controlexecute.ExecutionTreeNode, panels map[string]steampipeconfig.SnapshotPanel) {
	if treeNode == nil {
		return
	}

	// If the tree node matches a snapshot panel, populate the result
	if panel, exists := panels[treeNode.GetName()]; exists {
		if resultSetter, ok := treeNode.(interface {
			SetResult(panel steampipeconfig.SnapshotPanel)
		}); ok {
			resultSetter.SetResult(panel)
		}
	}

	// Recursively populate results for children
	for _, child := range treeNode.GetChildren() {
		populateResults(child, panels)
	}
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
