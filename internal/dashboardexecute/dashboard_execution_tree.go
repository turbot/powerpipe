package dashboardexecute

import (
	"context"
	"fmt"
	"golang.org/x/exp/maps"
	"log/slog"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/schema"
	"github.com/turbot/pipe-fittings/steampipeconfig"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/powerpipe/internal/dashboardevents"
	"github.com/turbot/powerpipe/internal/dashboardtypes"
	"github.com/turbot/powerpipe/internal/dashboardworkspace"
	"github.com/turbot/powerpipe/internal/db_client"
)

// DashboardExecutionTree is a structure representing the control result hierarchy
type DashboardExecutionTree struct {
	Root dashboardtypes.DashboardTreeRun

	dashboardName string
	sessionId     string
	// map of clients, keyed by connection string
	clients                 map[string]*db_client.DbClient
	clientsMut              sync.RWMutex
	defaultConnectionString string
	// map of executing runs, keyed by full name
	runs        map[string]dashboardtypes.DashboardTreeRun
	workspace   *dashboardworkspace.WorkspaceEvents
	runComplete chan dashboardtypes.DashboardTreeRun

	// map of subscribers to notify when an input value changes
	cancel      context.CancelFunc
	inputLock   sync.Mutex
	inputValues map[string]any
	id          string
}

func NewDashboardExecutionTree(rootResource modconfig.ModTreeItem, sessionId string, clients map[string]*db_client.DbClient, workspace *dashboardworkspace.WorkspaceEvents) (*DashboardExecutionTree, error) {
	// now populate the DashboardExecutionTree
	executionTree := &DashboardExecutionTree{
		dashboardName: rootResource.Name(),
		sessionId:     sessionId,
		// TODO KAI pass in default connection string? <MISC>
		defaultConnectionString: viper.GetString(constants.ArgWorkspaceDatabase),
		clients:                 clients,
		runs:                    make(map[string]dashboardtypes.DashboardTreeRun),
		workspace:               workspace,
		runComplete:             make(chan dashboardtypes.DashboardTreeRun, 1),
		inputValues:             make(map[string]any),
	}
	executionTree.id = fmt.Sprintf("%p", executionTree)

	// create the root run node (either a report run or a counter run)
	root, err := executionTree.createRootItem(rootResource)
	if err != nil {
		return nil, err
	}

	executionTree.Root = root
	return executionTree, nil
}

func (e *DashboardExecutionTree) createRootItem(rootResource modconfig.ModTreeItem) (dashboardtypes.DashboardTreeRun, error) {
	switch r := rootResource.(type) {
	case *modconfig.Dashboard:
		return NewDashboardRun(r, e, e)
	case *modconfig.Benchmark:
		return NewCheckRun(r, e, e)
	case *modconfig.Query:
		// wrap this in a chart and a dashboard
		dashboard, err := modconfig.NewQueryDashboard(r)
		// TACTICAL - set the execution tree dashboard name from the query dashboard
		e.dashboardName = dashboard.Name()
		if err != nil {
			return nil, err
		}
		return NewDashboardRun(dashboard, e, e)
	case *modconfig.Control:
		// wrap this in a chart and a dashboard
		dashboard, err := modconfig.NewQueryDashboard(r)
		if err != nil {
			return nil, err
		}
		return NewDashboardRun(dashboard, e, e)
	default:
		return nil, fmt.Errorf("type %T cannot be executed as dashboard", r)
	}
}

func (e *DashboardExecutionTree) Execute(ctx context.Context) {
	startTime := time.Now()

	// TODO KAI WHAT DO WE DO HERE <SEARCH PATH>
	//searchPath := e.defaultClient().GetRequiredSessionSearchPath()

	// store context
	cancelCtx, cancel := context.WithCancel(ctx)
	e.cancel = cancel
	workspace := e.workspace

	// perform any necessary initialisation
	// (e.g. check run creates the control execution tree)
	e.Root.Initialise(cancelCtx)
	if e.Root.GetError() != nil {
		return
	}

	// TODO KAI HOW TO HANDLE SEARCH PATHS?
	// TODO should we always wait even with non custom search path?
	// if there is a custom search path, wait until the first connection of each plugin has loaded
	//if customSearchPath := e.client.GetCustomSearchPath(); customSearchPath != nil {
	//	if err := connection_sync.WaitForSearchPathSchemas(ctx, e.client, customSearchPath); err != nil {
	//		e.Root.SetError(ctx, err)
	//		return
	//	}
	//}

	panels := e.BuildSnapshotPanels()
	// build map of those variables referenced by the dashboard run
	referencedVariables := GetReferencedVariables(e.Root, e.workspace)

	immutablePanels, err := utils.JsonCloneToMap(panels)
	if err != nil {
		e.SetError(ctx, err)
		return
	}
	workspace.PublishDashboardEvent(ctx, &dashboardevents.ExecutionStarted{
		Root:        e.Root,
		Session:     e.sessionId,
		ExecutionId: e.id,
		Panels:      immutablePanels,
		Inputs:      e.inputValues,
		Variables:   referencedVariables,
		StartTime:   startTime,
	})
	defer func() {

		e := &dashboardevents.ExecutionComplete{
			Root:        e.Root,
			Session:     e.sessionId,
			ExecutionId: e.id,
			Panels:      panels,
			Inputs:      e.inputValues,
			Variables:   referencedVariables,
			// search path elements are quoted (for consumption by postgres)
			// unquote them
			// TOSO STEAMPIPE ONLY
			SearchPath: nil, //utils.UnquoteStringArray(searchPath),
			StartTime:  startTime,
			EndTime:    time.Now(),
		}
		workspace.PublishDashboardEvent(ctx, e)
	}()

	slog.Debug("begin DashboardExecutionTree.Execute")
	defer slog.Debug("end DashboardExecutionTree.Execute")

	if e.GetRunStatus().IsFinished() {
		// there must be no nodes to execute
		slog.Debug("execution tree already complete")
		return
	}

	// execute synchronously
	e.Root.Execute(cancelCtx)
}

// GetRunStatus returns the stats of the Root run
func (e *DashboardExecutionTree) GetRunStatus() dashboardtypes.RunStatus {
	return e.Root.GetRunStatus()
}

// SetError sets the error on the Root run
func (e *DashboardExecutionTree) SetError(ctx context.Context, err error) {
	e.Root.SetError(ctx, err)
}

// GetName implements DashboardParent
// use mod short name - this will be the root name for all child runs
func (e *DashboardExecutionTree) GetName() string {
	return e.workspace.Mod.ShortName
}

// GetParent implements DashboardTreeRun
func (e *DashboardExecutionTree) GetParent() dashboardtypes.DashboardParent {
	return nil
}

// GetNodeType implements DashboardTreeRun
func (*DashboardExecutionTree) GetNodeType() string {
	panic("should never call for DashboardExecutionTree")
}

func (e *DashboardExecutionTree) SetInputValues(inputValues map[string]any) {
	slog.Debug("SetInputValues")
	e.inputLock.Lock()
	defer e.inputLock.Unlock()

	// we only support inputs if root is a dashboard (NOT a benchmark)
	runtimeDependencyPublisher, ok := e.Root.(RuntimeDependencyPublisher)
	if !ok {
		// should never happen
		slog.Warn("SetInputValues called but root WorkspaceEvents run is not a RuntimeDependencyPublisher", "root", e.Root.GetName())
		return
	}

	for name, value := range inputValues {
		slog.Debug("DashboardExecutionTree SetInput", "name", name, "value", value)
		e.inputValues[name] = value
		// publish runtime dependency
		runtimeDependencyPublisher.PublishRuntimeDependencyValue(name, &dashboardtypes.ResolvedRuntimeDependencyValue{Value: value})
	}
}

// ChildCompleteChan implements DashboardParent
func (e *DashboardExecutionTree) ChildCompleteChan() chan dashboardtypes.DashboardTreeRun {
	return e.runComplete
}

// ChildStatusChanged implements DashboardParent
func (*DashboardExecutionTree) ChildStatusChanged(context.Context) {}

func (e *DashboardExecutionTree) Cancel() {
	// if we have not completed, and already have a cancel function - cancel
	if e.GetRunStatus().IsFinished() || e.cancel == nil {
		slog.Debug("DashboardExecutionTree Cancel NOT cancelling", "status", e.GetRunStatus(), "cancel func", e.cancel)
		return
	}

	slog.Debug("DashboardExecutionTree Cancel  - calling cancel")
	e.cancel()

	// if there are any children, wait for the execution to complete
	if !e.Root.RunComplete() {
		<-e.runComplete
	}

	slog.Debug("DashboardExecutionTree Cancel - all children complete")
}

func (e *DashboardExecutionTree) BuildSnapshotPanels() map[string]steampipeconfig.SnapshotPanel {
	// just build from e.runs
	res := map[string]steampipeconfig.SnapshotPanel{}

	for name, run := range e.runs {
		res[name] = run.(steampipeconfig.SnapshotPanel)
		// special case handling for check runs
		if checkRun, ok := run.(*CheckRun); ok {
			checkRunChildren := checkRun.BuildSnapshotPanels(res)
			for k, v := range checkRunChildren {
				res[k] = v
			}
		}
	}
	return res
}

// InputRuntimeDependencies returns the names of all inputs which are runtime dependencies
func (e *DashboardExecutionTree) InputRuntimeDependencies() []string {
	var deps = map[string]struct{}{}
	for _, r := range e.runs {
		if leafRun, ok := r.(*LeafRun); ok {
			for _, r := range leafRun.runtimeDependencies {
				if r.Dependency.PropertyPath.ItemType == schema.BlockTypeInput {
					deps[r.Dependency.SourceResourceName()] = struct{}{}
				}
			}
		}
	}
	return maps.Keys(deps)
}

// GetChildren implements DashboardParent
func (e *DashboardExecutionTree) GetChildren() []dashboardtypes.DashboardTreeRun {
	return []dashboardtypes.DashboardTreeRun{e.Root}
}

// ChildrenComplete implements DashboardParent
func (e *DashboardExecutionTree) ChildrenComplete() bool {
	return e.Root.RunComplete()
}

// Tactical: Empty implementations of DashboardParent functions
// TODO remove need for this

func (e *DashboardExecutionTree) Initialise(ctx context.Context) {
	panic("should never call for DashboardExecutionTree")
}

func (e *DashboardExecutionTree) GetTitle() string {
	panic("should never call for DashboardExecutionTree")
}

func (e *DashboardExecutionTree) GetError() error {
	panic("should never call for DashboardExecutionTree")
}

func (e *DashboardExecutionTree) SetComplete(ctx context.Context) {
	panic("should never call for DashboardExecutionTree")
}

func (e *DashboardExecutionTree) RunComplete() bool {
	panic("should never call for DashboardExecutionTree")
}

func (e *DashboardExecutionTree) GetInputsDependingOn(s string) []string {
	panic("should never call for DashboardExecutionTree")
}

func (*DashboardExecutionTree) AsTreeNode() *steampipeconfig.SnapshotTreeNode {
	panic("should never call for DashboardExecutionTree")
}

func (*DashboardExecutionTree) GetResource() modconfig.DashboardLeafNode {
	panic("should never call for DashboardExecutionTree")
}

//nolint:unused // TODO: unused function
func (e *DashboardExecutionTree) defaultClient() *db_client.DbClient {
	return e.clients[e.defaultConnectionString]
}

func (e *DashboardExecutionTree) getClient(ctx context.Context, connectionString string) (*db_client.DbClient, error) {
	e.clientsMut.RLock()
	client := e.clients[connectionString]
	e.clientsMut.RUnlock()

	if client != nil {
		return client, nil
	}
	e.clientsMut.Lock()
	defer e.clientsMut.Unlock()

	client, err := db_client.NewDbClient(ctx, connectionString)
	if err != nil {
		return nil, err
	}
	return client, nil
}
