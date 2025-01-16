package controldisplay

import (
	"context"
	"fmt"

	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/pipes"
	"github.com/turbot/pipe-fittings/statushooks"
	"github.com/turbot/pipe-fittings/steampipeconfig"
	"github.com/turbot/powerpipe/internal/controlexecute"
	"github.com/turbot/powerpipe/internal/controlstatus"
	"github.com/turbot/powerpipe/internal/dashboardexecute"
	"github.com/turbot/powerpipe/internal/resources"
	"github.com/turbot/powerpipe/internal/workspace"
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

func SnapshotToExecutionTree(ctx context.Context, snapshot *steampipeconfig.SteampipeSnapshot, w *workspace.PowerpipeWorkspace, targets ...modconfig.ModTreeItem) (*dashboardexecute.DetectionBenchmarkDisplayTree, error) {
	// Step 1: Create the execution tree
	tree, err := newDisplayExecutionTree(snapshot, w, targets...)
	if err != nil {
		return nil, err
	}

	// Step 2: Populate summaries
	populateSummaries(tree.Root, snapshot.Layout)

	// Step 3: Populate results
	populateResults(tree.Root, snapshot.Panels)

	// Step 4: Add any additional metadata from the snapshot
	tree.StartTime = snapshot.StartTime
	tree.EndTime = snapshot.EndTime

	return tree, nil
}

func newDisplayExecutionTree(snapshot *steampipeconfig.SteampipeSnapshot, w *workspace.PowerpipeWorkspace, targets ...modconfig.ModTreeItem) (*dashboardexecute.DetectionBenchmarkDisplayTree, error) {
	// now populate the ExecutionTree
	executionTree := &dashboardexecute.DetectionBenchmarkDisplayTree{
		LeafRuns: make(map[string]controlexecute.LeafRun),
	}

	var resolvedItem modconfig.ModTreeItem
	// if only one argument is provided, add this as execution root
	if len(targets) == 1 {
		resolvedItem = targets[0]
	} else {
		// create a root benchmark with `items` as it's children
		resolvedItem = resources.NewRootBenchmarkWithChildren(w.Mod, targets).(modconfig.ModTreeItem)
	}

	// build tree of result groups, starting with a synthetic 'root' node

	root, err := dashboardexecute.NewRootBenchmarkDisplay(resolvedItem)
	if err != nil {
		return nil, err
	}

	// now traverse the layout snapshot layout, find the corresponding items in the snapshot panesl and build the tree
	rootLayout := snapshot.Layout

	// build node for this item
	rootRun := snapshot.Panels[rootLayout.Name]
	if rootRun == nil {
		return nil, fmt.Errorf("rootRun %s not found in panels", rootLayout.Name)
	}

	switch resource := rootRun.(type) {
	case *dashboardexecute.DetectionRun:
		root.AddDetection(resource)
	case *dashboardexecute.BenchmarkRun:
		// create a result group for this item
		benchmarkGroup, err := dashboardexecute.NewDetectionBenchmarkDisplay(resource, root)
		if err != nil {
			return nil, err
		}
		root.AddResultGroup(benchmarkGroup)
	}

	executionTree.Root = root

	// after tree has built, ControlCount will be set - create progress rendered
	executionTree.Progress = controlstatus.NewControlProgress(len(executionTree.LeafRuns))

	return executionTree, nil
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
