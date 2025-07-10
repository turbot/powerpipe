package dashboardexecute

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	filehelpers "github.com/turbot/go-kit/files"
	"github.com/turbot/pipe-fittings/v2/backend"
	"github.com/turbot/pipe-fittings/v2/connection"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/pipe-fittings/v2/utils"
	"github.com/turbot/powerpipe/internal/dashboardevents"
	"github.com/turbot/powerpipe/internal/dashboardtypes"
	"github.com/turbot/powerpipe/internal/db_client"
	"github.com/turbot/powerpipe/internal/workspace"
)

type DashboardExecutor struct {
	// map of executions, keyed by session id
	executions    map[string]*DashboardExecutionTree
	executionLock sync.Mutex
	// is this an interactive execution
	// i.e. inputs may be specified _after_ execution starts
	// false when running a single dashboard in batch mode
	interactive bool
	// store the default client which is created during initData creation
	// - this is to avoid creating a new client for each dashboard execution if the database/search path is NOT overridden
	defaultClient           *db_client.ClientMap
	defaultDatabase         connection.ConnectionStringProvider
	defaultSearchPathConfig backend.SearchPathConfig
}

func NewDashboardExecutor(defaultClient *db_client.ClientMap, defaultDatabase connection.ConnectionStringProvider, defaultSearchPathConfig backend.SearchPathConfig) *DashboardExecutor {
	return &DashboardExecutor{
		executions: make(map[string]*DashboardExecutionTree),
		// default to interactive execution
		interactive:             true,
		defaultClient:           defaultClient,
		defaultDatabase:         defaultDatabase,
		defaultSearchPathConfig: defaultSearchPathConfig,
	}
}

var Executor *DashboardExecutor

func (e *DashboardExecutor) ExecuteDashboard(ctx context.Context, sessionId string, rootResource modconfig.ModTreeItem, inputs *InputValues, workspace *workspace.PowerpipeWorkspace, opts ...backend.BackendOption) (err error) {
	var executionTree *DashboardExecutionTree
	defer func() {
		if err == nil && ctx.Err() != nil {
			if ctx.Err() != nil && ctx.Err().Error() == context.DeadlineExceeded.Error() {
				err = fmt.Errorf("execution timed out")
			} else {
				err = ctx.Err()
			}
		}
		// if there was an error executing, send an ExecutionError event
		if err != nil {
			errorEvent := &dashboardevents.ExecutionError{
				Error:     err,
				Session:   sessionId,
				Timestamp: time.Now(),
			}
			workspace.PublishDashboardEvent(ctx, errorEvent)
		}
	}()

	// reset any existing executions for this session
	e.CancelExecutionForSession(ctx, sessionId)

	// now create a new execution
	executionTree, err = e.newDashboardExecutionTree(rootResource, sessionId, workspace, inputs, opts...)
	if err != nil {
		return err
	}

	// add to execution map
	e.setExecution(sessionId, executionTree)

	// if inputs must be provided before execution (i.e. this is a batch dashboard execution),
	// verify all required inputs are provided
	if err = e.validateInputs(executionTree, inputs.Inputs); err != nil {
		return err
	}

	go executionTree.Execute(ctx)

	return nil
}

// if inputs must be provided before execution (i.e. this is a batch dashboard execution),
// verify all required inputs are provided
func (e *DashboardExecutor) validateInputs(executionTree *DashboardExecutionTree, inputs map[string]interface{}) error {
	if e.interactive {
		// interactive dashboard execution - no need to validate
		return nil
	}
	var missingInputs []string
	for _, inputName := range executionTree.InputRuntimeDependencies() {
		if _, ok := inputs[inputName]; !ok {
			missingInputs = append(missingInputs, inputName)
		}
	}
	if missingCount := len(missingInputs); missingCount > 0 {
		return fmt.Errorf("%s '%s' must be provided using '--arg name=value'", utils.Pluralize("input", missingCount), strings.Join(missingInputs, ","))
	}

	return nil
}

func (e *DashboardExecutor) LoadSnapshot(ctx context.Context, sessionId, snapshotName string, w *workspace.PowerpipeWorkspace) (map[string]any, error) {
	// find snapshot path in workspace
	snapshotPath, ok := w.GetPowerpipeModResources().Snapshots[snapshotName]
	if !ok {
		return nil, fmt.Errorf("snapshot %s not found in %s (%s)", snapshotName, w.Mod.Name(), w.Path)
	}

	if !filehelpers.FileExists(snapshotPath) {
		return nil, fmt.Errorf("snapshot %s not does not exist", snapshotPath)
	}

	snapshotContent, err := os.ReadFile(snapshotPath)
	if err != nil {
		return nil, err
	}

	// deserialize the snapshot as an interface map
	// we cannot deserialize into a SteampipeSnapshot struct
	// (without custom derserialisation code) as the Panels property is an interface
	snap := map[string]any{}

	err = json.Unmarshal(snapshotContent, &snap)
	if err != nil {
		return nil, err
	}

	return snap, nil
}

func (e *DashboardExecutor) OnInputChanged(ctx context.Context, sessionId string, inputs *InputValues, changedInput string) error {
	// find the execution
	executionTree, found := e.executions[sessionId]
	if !found {
		return fmt.Errorf("no dashboard running for session %s", sessionId)
	}

	// get the previous value of this input
	inputPrevValue := executionTree.inputValues[changedInput]
	// if there are any dependent inputs, set their value to nil and send an event to the UI
	// first see if any other inputs rely on the one which was just changed
	dependentInputs := e.clearDependentInputs(executionTree.Root, changedInput, inputs.Inputs)
	if len(dependentInputs) > 0 {
		event := &dashboardevents.InputValuesCleared{
			ClearedInputs: dependentInputs,
			Session:       executionTree.sessionId,
			ExecutionId:   executionTree.id,
		}
		executionTree.workspace.PublishDashboardEvent(ctx, event)
	}

	// has the time range changed
	timeRangeChanged := !inputs.DateTimeRange.Equals(executionTree.DateTimeRange)
	currentRunFinished := executionTree.GetRunStatus().IsFinished()
	prevInputsExist := inputPrevValue != nil

	// input has changed - should we immediately re-execute?

	// we should re-execute if:
	// - the execution has completed - reexecute
	// - the time range has changed - reexecute
	// - the input value was NOT previously nil
	// (i.e. this is really a CHANGE of input not just the first time the inputs have been set)
	// NOTE: if the previous input value is nil and we are currently executing we do not need to re-execute
	// as the current execution will be waiting for the inputs to be available
	if currentRunFinished || timeRangeChanged || prevInputsExist {
		return e.ExecuteDashboard(
			ctx,
			sessionId,
			executionTree.Root.GetResource(),
			inputs,
			executionTree.workspace)
	}

	// ok we we are NOT re-executing - just set the inputs
	executionTree.SetInputValues(inputs)

	return nil
}

func (e *DashboardExecutor) clearDependentInputs(root dashboardtypes.DashboardTreeRun, changedInput string, inputs map[string]any) []string {
	dependentInputs := root.GetInputsDependingOn(changedInput)
	clearedInputs := dependentInputs
	if len(dependentInputs) > 0 {
		for _, inputName := range dependentInputs {
			if inputs[inputName] != nil {
				// clear the input value
				inputs[inputName] = nil
				childDependentInputs := e.clearDependentInputs(root, inputName, inputs)
				clearedInputs = append(clearedInputs, childDependentInputs...)
			}
		}
	}

	return clearedInputs
}

func (e *DashboardExecutor) CancelExecutionForSession(_ context.Context, sessionId string) {
	// find the execution
	executionTree, found := e.getExecution(sessionId)
	if !found {
		// nothing to do
		return
	}

	// cancel if in progress
	executionTree.Cancel()
	// remove from execution tree
	e.removeExecution(sessionId)
}

// find the execution for the given session id
func (e *DashboardExecutor) getExecution(sessionId string) (*DashboardExecutionTree, bool) {
	e.executionLock.Lock()
	defer e.executionLock.Unlock()

	executionTree, found := e.executions[sessionId]
	return executionTree, found
}

func (e *DashboardExecutor) setExecution(sessionId string, executionTree *DashboardExecutionTree) {
	e.executionLock.Lock()
	defer e.executionLock.Unlock()

	e.executions[sessionId] = executionTree
}

func (e *DashboardExecutor) removeExecution(sessionId string) {
	e.executionLock.Lock()
	defer e.executionLock.Unlock()

	delete(e.executions, sessionId)
}
