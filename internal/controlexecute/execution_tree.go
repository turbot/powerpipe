package controlexecute

import (
	"context"
	"log/slog"
	"sort"
	"time"

	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/backend"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/workspace"
	"github.com/turbot/powerpipe/internal/controlstatus"
	"github.com/turbot/powerpipe/internal/db_client"
	"golang.org/x/sync/semaphore"
)

// ExecutionTree is a structure representing the control execution hierarchy
type ExecutionTree struct {
	Root *ResultGroup `json:"root"`
	// flat list of all control runs
	ControlRuns []*ControlRun                  `json:"-"`
	StartTime   time.Time                      `json:"start_time"`
	EndTime     time.Time                      `json:"end_time"`
	Progress    *controlstatus.ControlProgress `json:"progress"`
	// map of dimension property name to property value to color map
	DimensionColorGenerator *DimensionColorGenerator `json:"-"`
	// the current session search path
	SearchPath []string             `json:"-"`
	Workspace  *workspace.Workspace `json:"-"`
	client     *db_client.DbClient
	// an optional map of control names used to filter the controls which are run
	controlNameFilterMap map[string]struct{}
}

func NewExecutionTree(ctx context.Context, workspace *workspace.Workspace, client *db_client.DbClient, controlFilter workspace.ResourceFilter, targets ...modconfig.ModTreeItem) (*ExecutionTree, error) {
	// now populate the ExecutionTree
	executionTree := &ExecutionTree{
		Workspace: workspace,
		client:    client,
	}

	// if backend supports search path, get it
	if sp, ok := client.Backend.(backend.SearchPathProvider); ok {
		executionTree.SearchPath = sp.RequiredSearchPath()
	}

	// if a "--where" or "--tag" parameter was passed, build a map of control names used to filter the controls to run
	err := executionTree.populateControlFilterMap(controlFilter)
	if err != nil {
		return nil, err
	}

	var resolvedItem modconfig.ModTreeItem
	// if only one argument is provided, add this as execution root
	if len(targets) == 1 {
		resolvedItem = targets[0]
	} else {
		// create a root benchmark with `items` as it's children
		resolvedItem = modconfig.NewRootBenchmarkWithChildren(workspace.Mod, targets).(modconfig.ModTreeItem)
	}

	// build tree of result groups, starting with a synthetic 'root' node
	executionTree.Root, err = NewRootResultGroup(ctx, executionTree, resolvedItem)
	if err != nil {
		return nil, err
	}

	// after tree has built, ControlCount will be set - create progress rendered
	executionTree.Progress = controlstatus.NewControlProgress(len(executionTree.ControlRuns))

	return executionTree, nil
}

// IsExportSourceData implements ExportSourceData
func (*ExecutionTree) IsExportSourceData() {}

// AddControl checks whether control should be included in the tree
// if so, creates a ControlRun, which is added to the parent group
func (e *ExecutionTree) AddControl(ctx context.Context, control *modconfig.Control, group *ResultGroup) error {
	// note we use short name to determine whether to include a control
	if e.ShouldIncludeControl(control.Name()) {
		// create new ControlRun with treeItem as the parent
		controlRun, err := NewControlRun(control, group, e)
		if err != nil {
			return err
		}
		// add it into the group
		group.addControl(controlRun)

		// also add it into the execution tree control run list
		e.ControlRuns = append(e.ControlRuns, controlRun)
	}
	return nil
}

func (e *ExecutionTree) Execute(ctx context.Context) error {
	slog.Debug("begin ExecutionTree.Execute")
	defer slog.Debug("end ExecutionTree.Execute")
	e.StartTime = time.Now()
	e.Progress.Start(ctx)

	defer func() {
		e.EndTime = time.Now()
		e.Progress.Finish(ctx)
	}()

	// the number of goroutines parallel to start
	var maxParallelGoRoutines int64 = constants.DefaultMaxConnections
	if viper.IsSet(constants.ArgMaxParallel) {
		maxParallelGoRoutines = viper.GetInt64(constants.ArgMaxParallel)
	}

	// to limit the number of parallel controls go routines started
	parallelismLock := semaphore.NewWeighted(maxParallelGoRoutines)

	// just execute the root - it will traverse the tree
	e.Root.execute(ctx, e.client, parallelismLock)

	if err := e.waitForActiveRunsToComplete(ctx, parallelismLock, maxParallelGoRoutines); err != nil {
		slog.Warn("timed out waiting for active runs to complete")
	}

	// now build map of dimension property name to property value to color map
	e.DimensionColorGenerator, _ = NewDimensionColorGenerator(4, 27)
	e.DimensionColorGenerator.populate(e)

	return nil
}

func (e *ExecutionTree) waitForActiveRunsToComplete(ctx context.Context, parallelismLock *semaphore.Weighted, maxParallelGoRoutines int64) error {
	waitCtx := ctx
	// if the context was already cancelled, we must creat ea new one to use  when waiting to acquire the lock
	if ctx.Err() != nil {
		// use a Background context - since the original context has been cancelled
		// this lets us wait for the active control queries to cancel
		c, cancel := context.WithTimeout(context.Background(), constants.ControlQueryCancellationTimeoutSecs*time.Second)
		waitCtx = c
		defer cancel()
	}
	// wait till we can acquire all semaphores - meaning that all active runs have finished
	return parallelismLock.Acquire(waitCtx, maxParallelGoRoutines)
}

func (e *ExecutionTree) populateControlFilterMap(controlFilter workspace.ResourceFilter) error {
	// if we derived or were passed a where clause, run the filter
	if controlFilter.Empty() {
		return nil
	}

	slog.Debug("filtering controls with", "controlFilter", controlFilter)
	var err error
	e.controlNameFilterMap, err = e.getControlMapFromFilter(controlFilter)
	if err != nil {
		return err
	}

	return nil
}

func (e *ExecutionTree) ShouldIncludeControl(controlName string) bool {
	if e.controlNameFilterMap == nil {
		return true
	}
	_, ok := e.controlNameFilterMap[controlName]
	return ok
}

// Get a map of control names from the introspection table steampipe_control
// This is used to implement the 'where' control filtering
func (e *ExecutionTree) getControlMapFromFilter(controlFilter workspace.ResourceFilter) (map[string]struct{}, error) {
	var res = make(map[string]struct{})
	controls, err := workspace.FilterWorkspaceResourcesOfType[*modconfig.Control](e.Workspace, controlFilter)
	if err != nil {
		return nil, err
	}

	for _, control := range controls {
		res[control.Name()] = struct{}{}
	}

	return res, nil
}

func (e *ExecutionTree) GetAllTags() []string {
	// map keep track which tags have been added as columns
	tagColumnMap := make(map[string]bool)
	var tagColumns []string
	for _, r := range e.ControlRuns {
		if r.Control.Tags != nil {
			for tag := range r.Control.Tags {
				if !tagColumnMap[tag] {
					tagColumns = append(tagColumns, tag)
					tagColumnMap[tag] = true
				}
			}
		}
	}
	sort.Strings(tagColumns)
	return tagColumns
}
