package dashboardexecute

import (
	"context"
	"fmt"
	"golang.org/x/exp/maps"
	"log/slog"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/backend"
	"github.com/turbot/pipe-fittings/connection"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/schema"
	"github.com/turbot/pipe-fittings/steampipeconfig"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/powerpipe/internal/dashboardevents"
	"github.com/turbot/powerpipe/internal/dashboardtypes"
	"github.com/turbot/powerpipe/internal/db_client"
	"github.com/turbot/powerpipe/internal/resources"
	"github.com/turbot/powerpipe/internal/workspace"
)

// DashboardExecutionTree is a structure representing the control result hierarchy
type DashboardExecutionTree struct {
	Root dashboardtypes.DashboardTreeRun

	dashboardName string
	sessionId     string
	// map of clients, keyed by connection string - we will close this at end of execution
	clientMap *db_client.ClientMap
	// map of server-managed clients, keyed by connection string - we will NOT close this
	defaultClientMap *db_client.ClientMap
	// map of executing runs, keyed by full name
	runs      map[string]dashboardtypes.DashboardTreeRun
	workspace *workspace.PowerpipeWorkspace

	runComplete chan dashboardtypes.DashboardTreeRun
	// map of subscribers to notify when an input value changes
	cancel      context.CancelFunc
	inputLock   sync.Mutex
	inputValues map[string]any
	id          string
	// active database and search path config (unless overridden at the resource level)
	database           connection.ConnectionStringProvider
	searchPathConfig   backend.SearchPathConfig
	DetectionTimeRange utils.TimeRange
}

func newDashboardExecutionTree(rootResource modconfig.ModTreeItem, sessionId string, workspace *workspace.PowerpipeWorkspace, inputs *InputValues, defaultClientMap *db_client.ClientMap, opts ...backend.BackendOption) (*DashboardExecutionTree, error) {
	// now populate the DashboardExecutionTree
	executionTree := &DashboardExecutionTree{
		dashboardName:    rootResource.Name(),
		sessionId:        sessionId,
		defaultClientMap: defaultClientMap,
		clientMap:        db_client.NewClientMap(),
		runs:             make(map[string]dashboardtypes.DashboardTreeRun),
		workspace:        workspace,
		runComplete:      make(chan dashboardtypes.DashboardTreeRun, 1),
		inputValues:      make(map[string]any),
	}
	executionTree.id = fmt.Sprintf("%p", executionTree)

	// set the dashboard database and search patch config
	defaultDatabase, defaultSearchPathConfig, err := db_client.GetDefaultDatabaseConfig(opts...)
	if err != nil {
		return nil, err
	}
	database, searchPathConfig, err := db_client.GetDatabaseConfigForResource(rootResource, workspace.Mod, defaultDatabase, defaultSearchPathConfig)
	if err != nil {
		return nil, err
	}
	executionTree.database = database
	executionTree.searchPathConfig = searchPathConfig

	// create the root run node (either a report run or a counter run)
	root, err := executionTree.createRootItem(rootResource)
	if err != nil {
		return nil, err
	}

	executionTree.Root = root

	// if inputs have been passed, set them first
	executionTree.SetInputValues(inputs)

	// add a client for the active database and search path
	_, err = executionTree.getClient(context.Background(), database, searchPathConfig)
	if err != nil {
		return nil, err
	}

	return executionTree, nil
}

func (e *DashboardExecutionTree) createRootItem(rootResource modconfig.ModTreeItem) (dashboardtypes.DashboardTreeRun, error) {
	switch r := rootResource.(type) {
	case *resources.Dashboard:
		return NewDashboardRun(r, e, e)
	case *resources.Benchmark:
		return NewCheckRun(r, e, e)
	case *resources.Detection:
		// create a wrapper for the detection
		benchmark := resources.NewWrapperDetectionBenchmark(r)
		return NewDetectionBenchmarkRun(benchmark, e, e)

	case *resources.DetectionBenchmark:
		return NewDetectionBenchmarkRun(r, e, e)
	case *resources.Query:
		// wrap this in a chart and a dashboard
		dashboard, err := resources.NewQueryDashboard(r)
		// TACTICAL - set the execution tree dashboard name from the query dashboard
		e.dashboardName = dashboard.Name()
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

	// setup a cancel context with timeout and start cancel handler
	var cancel context.CancelFunc
	// if a dashboard timeout was specified, use that
	if executionTimeout := viper.GetInt(constants.ArgDashboardTimeout); executionTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(executionTimeout)*time.Second)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}

	e.cancel = cancel
	workspace := e.workspace

	// if the default database backend supports search path, retrieve it
	defaultClient, err := e.getClient(ctx, e.database, e.searchPathConfig)
	if err != nil {
		e.SetError(ctx, err)
		return
	}
	var searchPath []string
	if sp, ok := defaultClient.Backend.(backend.SearchPathProvider); ok {
		searchPath = sp.RequiredSearchPath()
	}

	// perform any necessary initialisation
	// (e.g. check run creates the control execution tree)
	e.Root.Initialise(ctx)
	if e.Root.GetError() != nil {
		return
	}

	panels := e.BuildSnapshotPanels()
	// build map of those variables referenced by the dashboard run
	referencedVariables, err := GetReferencedVariables(e.Root, e.workspace)
	if err != nil {
		e.SetError(ctx, err)
		return
	}
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

		ev := &dashboardevents.ExecutionComplete{
			Root:        e.Root,
			Session:     e.sessionId,
			ExecutionId: e.id,
			Panels:      panels,
			Inputs:      e.inputValues,
			Variables:   referencedVariables,
			SearchPath:  searchPath,
			StartTime:   startTime,
			EndTime:     time.Now(),
		}

		workspace.PublishDashboardEvent(ctx, ev)
	}()

	slog.Debug("begin DashboardExecutionTree.Execute")
	defer slog.Debug("end DashboardExecutionTree.Execute")

	if e.GetRunStatus().IsFinished() {
		// there must be no nodes to execute
		slog.Debug("execution tree already complete")
		return
	}

	// execute synchronously
	e.Root.Execute(ctx)

	// now close any clients created just for this run
	e.clientMap.Close(ctx)
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

func (e *DashboardExecutionTree) SetInputValues(inputValues *InputValues) {
	slog.Debug("SetInputValues")
	e.inputLock.Lock()
	defer e.inputLock.Unlock()

	if inputValues == nil {
		slog.Warn("SetInputValues - inputValues is nil")
		return
	}
	// set the input values and publish the runtime dependencies (if root implements RuntimeDependencyPublisher)
	runtimeDependencyPublisher, _ := e.Root.(RuntimeDependencyPublisher)
	for name, value := range inputValues.Inputs {
		slog.Debug("DashboardExecutionTree SetInput", "name", name, "value", value)
		e.inputValues[name] = value
		// publish runtime dependency
		if runtimeDependencyPublisher != nil {
			runtimeDependencyPublisher.PublishRuntimeDependencyValue(name, &dashboardtypes.ResolvedRuntimeDependencyValue{Value: value})
		}
	}

	// TODO K
	// TACTICAL
	// if a time range has been passed, set the detection time range
	// (if time range is not set, From and To will be nil - this is expected and handled)
	e.DetectionTimeRange = inputValues.DetectionTimeRange
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

func (*DashboardExecutionTree) GetResource() resources.DashboardLeafNode {
	panic("should never call for DashboardExecutionTree")
}

// function to get a client from one of the client maps
func (e *DashboardExecutionTree) getClient(ctx context.Context, csp connection.ConnectionStringProvider, searchPathConfig backend.SearchPathConfig) (*db_client.DbClient, error) {
	// ask the provider for the connection string, passing the filter
	// TODO check connection type is tailpipe???
	filter := &connection.TailpipeDatabaseFilters{From: e.DetectionTimeRange.From, To: e.DetectionTimeRange.To}
	cs, err := csp.GetConnectionString(connection.WithFilter(filter))
	if err != nil {
		return nil, err
	}
	// if the default map already contains a client for this connection string, use that
	if client := e.defaultClientMap.Get(cs, searchPathConfig); client != nil {
		return client, nil
	}

	// otherwise get or create one
	client, err := e.clientMap.GetOrCreate(ctx, cs, searchPathConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}
