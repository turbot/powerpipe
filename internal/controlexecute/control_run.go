package controlexecute

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	typehelpers "github.com/turbot/go-kit/types"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/queryresult"
	"github.com/turbot/pipe-fittings/schema"
	"github.com/turbot/pipe-fittings/statushooks"
	"github.com/turbot/pipe-fittings/steampipeconfig"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/powerpipe/internal/controlstatus"
	"github.com/turbot/powerpipe/internal/dashboardtypes"
	"github.com/turbot/powerpipe/internal/db_client"
	localqueryresult "github.com/turbot/powerpipe/internal/queryresult"
	"github.com/turbot/powerpipe/internal/snapshot"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc"
)

// ControlRun is a struct representing the execution of a control run. It will contain one or more result items (i.e. for one or more resources).
type ControlRun struct {
	// properties from control
	ControlId     string            `json:"-"`
	FullName      string            `json:"name"`
	Title         string            `json:"title,omitempty"`
	Description   string            `json:"description,omitempty"`
	Documentation string            `json:"documentation,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
	Display       string            `json:"display,omitempty"`
	Type          string            `json:"display_type,omitempty"`

	// this will be serialised under 'properties'
	Severity string `json:"-"`

	// "control"
	NodeType string `json:"panel_type"`

	// the control being run
	Control *modconfig.Control `json:"-"`
	// this is populated by retrieving Control properties with the snapshot tag
	Properties map[string]any `json:"properties,omitempty"`

	// control summary
	Summary   *controlstatus.StatusSummary `json:"summary"`
	RunStatus dashboardtypes.RunStatus     `json:"status"`
	// result rows
	Rows ResultRows `json:"-"`

	// the results in snapshot format
	Data *dashboardtypes.LeafData `json:"data"`

	// a list of distinct dimension keys from the results of this control
	DimensionKeys []string `json:"-"`

	// execution duration
	Duration time.Duration `json:"-"`
	// parent result group
	Parents []*ResultGroup `json:"-"`
	// execution tree
	Tree *ExecutionTree `json:"-"`
	// save run error as string for JSON export
	RunErrorString string `json:"error,omitempty"`
	runError       error
	// the query result stream
	queryResult *localqueryresult.Result
	rowMap      map[string]ResultRows
	stateLock   sync.Mutex
	doneChan    chan bool
	attempts    int
	startTime   time.Time
}

type ResultRowInstance struct {
	ResultRow
	ControlRun *ControlRunInstance `json:"-"`
}

type ControlRunInstance struct {
	ControlRun
	Group *ResultGroup `json:"-"`
	Rows  []*ResultRowInstance
}

func (cr *ControlRun) CloneAndSetParent(parent *ResultGroup) ControlRunInstance {
	res := ControlRunInstance{
		ControlRun: *cr, //nolint:govet
		Group:      parent,
	}
	for _, r := range cr.Rows {
		res.Rows = append(res.Rows, &ResultRowInstance{
			ResultRow:  *r,
			ControlRun: &res,
		})
	}
	return res //nolint:govet
}

func NewControlRun(control *modconfig.Control, group *ResultGroup, executionTree *ExecutionTree) (*ControlRun, error) {
	controlId := control.Name()

	// only show qualified control names for controls from dependent mods
	if control.Mod.Name() == executionTree.Workspace.Mod.Name() {
		controlId = control.UnqualifiedName
	}

	res := &ControlRun{
		Control:       control,
		ControlId:     controlId,
		FullName:      control.Name(),
		Description:   control.GetDescription(),
		Documentation: control.GetDocumentation(),
		Tags:          control.GetTags(),
		Display:       control.GetDisplay(),
		Type:          control.GetType(),

		Severity:   typehelpers.SafeString(control.Severity),
		Title:      typehelpers.SafeString(control.Title),
		rowMap:     make(map[string]ResultRows),
		Summary:    &controlstatus.StatusSummary{},
		Tree:       executionTree,
		RunStatus:  dashboardtypes.RunInitialized,
		Parents:    []*ResultGroup{group},
		NodeType:   schema.BlockTypeControl,
		doneChan:   make(chan bool, 1),
		Properties: make(map[string]any),
	}
	if err := res.populateProperties(); err != nil {
		return nil, err
	}

	return res, nil
}

// GetControlId implements ControlRunStatusProvider
func (r *ControlRun) GetControlId() string {
	r.stateLock.Lock()
	defer r.stateLock.Unlock()
	return r.ControlId
}

// GetRunStatus implements ControlRunStatusProvider
func (r *ControlRun) GetRunStatus() dashboardtypes.RunStatus {
	r.stateLock.Lock()
	defer r.stateLock.Unlock()
	return r.RunStatus
}

// GetStatusSummary implements ControlRunStatusProvider
func (r *ControlRun) GetStatusSummary() *controlstatus.StatusSummary {
	r.stateLock.Lock()
	defer r.stateLock.Unlock()
	return r.Summary
}

func (r *ControlRun) Finished() bool {
	return r.GetRunStatus().IsFinished()
}

// MatchTag returns the value corresponding to the input key. Returns 'false' if not found
func (r *ControlRun) MatchTag(key string, value string) bool {
	val, found := r.Control.GetTags()[key]
	return found && (val == value)
}

func (r *ControlRun) GetError() error {
	return r.runError
}

// IsSnapshotPanel implements SnapshotPanel
func (*ControlRun) IsSnapshotPanel() {}

// IsExecutionTreeNode implements ExecutionTreeNode
func (*ControlRun) IsExecutionTreeNode() {}

// GetChildren implements ExecutionTreeNode
func (*ControlRun) GetChildren() []ExecutionTreeNode { return nil }

// GetName implements ExecutionTreeNode
func (r *ControlRun) GetName() string { return r.Control.Name() }

// AsTreeNode implements ExecutionTreeNode
func (r *ControlRun) AsTreeNode() *steampipeconfig.SnapshotTreeNode {
	res := &steampipeconfig.SnapshotTreeNode{
		Name:     r.Control.Name(),
		NodeType: r.NodeType,
	}
	return res
}

func (r *ControlRun) setError(ctx context.Context, err error) {
	if err == nil {
		return
	}
	// if finished, we dont set the error, this can happen because the same control run might have multiple parents
	if r.Finished() {
		slog.Debug("not setting the control run to error, already finished", "name", r.Control.Name(), "status", r.RunStatus, "error", err)
		return
	}

	if err.Error() == context.DeadlineExceeded.Error() {
		// had the control started?
		if r.RunStatus == dashboardtypes.RunRunning {
			r.runError = fmt.Errorf("control execution timed out after running for %0.2fs", time.Since(r.startTime).Seconds())
		} else {
			r.runError = fmt.Errorf("execution timed out before control started")
		}
	} else {
		r.runError = error_helpers.TransformErrorToSteampipe(err)
	}
	r.RunErrorString = r.runError.Error()
	// update error count
	r.Summary.Error++
	if error_helpers.IsContextCancelledError(err) {
		r.setRunStatus(ctx, dashboardtypes.RunCanceled)
	} else {
		r.setRunStatus(ctx, dashboardtypes.RunError)
	}

}
func (r *ControlRun) skip(ctx context.Context) {
	r.setRunStatus(ctx, dashboardtypes.RunComplete)
}

func (r *ControlRun) execute(ctx context.Context, client *db_client.DbClient) {
	utils.LogTime("ControlRun.execute start")
	defer utils.LogTime("ControlRun.execute end")

	slog.Debug("begin ControlRun.Start", "name", r.Control.Name())
	defer slog.Debug("end ControlRun.Start", "name", r.Control.Name())

	// check if the status has been set to running
	if !r.trySetStateRunning(ctx) {
		slog.Debug("control status has been set to running", "name", r.Control.Name())
		return
	}

	control := r.Control

	startTime := time.Now()

	// function to cleanup and update status after control run completion
	defer func() {
		r.Duration = time.Since(startTime)
		// update all our parents with our status - this will be passed all the way up the execution tree
		for _, parent := range r.Parents {
			parent.updateSummary(r.Summary)
			parent.onChildDone()
			if len(r.Severity) != 0 {
				parent.updateSeverityCounts(r.Severity, r.Summary)
			}
		}
	}()

	// update the current running control in the Progress renderer
	r.Tree.Progress.OnControlStart(ctx, r)
	defer func() {
		// update Progress
		if r.GetRunStatus() == dashboardtypes.RunError {
			r.Tree.Progress.OnControlError(ctx, r)
		} else {
			r.Tree.Progress.OnControlComplete(ctx, r)
		}
	}()

	// resolve the control query
	resolvedQuery, err := r.resolveControlQuery(control)
	if err != nil {
		r.setError(ctx, err)
		return
	}

	controlExecutionCtx := r.getControlQueryContext(ctx)

	// execute the control query
	// NOTE no need to pass an OnComplete callback - we are already closing our session after waiting for results
	slog.Debug("execute start", "name", r.Control.Name())
	queryResult, err := client.Execute(controlExecutionCtx, resolvedQuery.ExecuteSQL, resolvedQuery.Args...)
	slog.Debug("execute finish", "name", r.Control.Name())

	if err != nil {
		r.attempts++

		// is this an rpc EOF error - meaning that the plugin somehow crashed
		if grpc.IsGRPCConnectivityError(err) {
			if r.attempts < constants.MaxControlRunAttempts {
				slog.Debug("control query failed with plugin connectivity error - retrying…", "name", r.Control.Name(), "error", err)
				// recurse into this function to retry using the original context - which Execute will use to create it's own timeout context
				r.execute(ctx, client)
				return
			} else {
				slog.Debug("control query failed with plugin connectivity error - NOT retrying…", "name", r.Control.Name(), "error", err)
			}
		}
		r.setError(ctx, err)
		return
	}

	r.queryResult = queryResult

	// now wait for control completion
	slog.Debug("wait result", "name", r.Control.Name())
	r.waitForResults(ctx)
	slog.Debug("finish result", "name", r.Control.Name())
}

// create a context with status updates disabled (we do not want to show 'loading' results)
func (r *ControlRun) getControlQueryContext(ctx context.Context) context.Context {
	// disable the status spinner to hide 'loading' results)
	newCtx := statushooks.DisableStatusHooks(ctx)

	return newCtx
}

func (r *ControlRun) resolveControlQuery(control *modconfig.Control) (*modconfig.ResolvedQuery, error) {
	resolvedQuery, err := r.Tree.Workspace.ResolveQueryFromQueryProvider(control, nil)
	if err != nil {
		return nil, fmt.Errorf(`cannot run %s - failed to resolve query "%s": %s`, control.Name(), typehelpers.SafeString(control.SQL), err.Error())
	}
	return resolvedQuery, nil
}

func (r *ControlRun) waitForResults(ctx context.Context) {
	defer func() {
		dimensionsSchema := r.getDimensionSchema()
		// convert the data to snapshot format
		r.Data = r.Rows.ToLeafData(dimensionsSchema)
	}()

	for {
		select {
		case <-ctx.Done():
			r.setError(ctx, ctx.Err())
			return
		case row := <-*r.queryResult.RowChan:
			// nil row means control run is complete
			if row == nil {
				// nil row means we are done
				r.setRunStatus(ctx, dashboardtypes.RunComplete)
				r.createdOrderedResultRows()
				return
			}
			// create a result row
			result, err := NewResultRow(r, row, r.queryResult.Cols)
			if err != nil {
				r.setError(ctx, err)
				return
			}
			r.addResultRow(result)
		case <-r.doneChan:
			return
		}
	}
}

func (r *ControlRun) getDimensionSchema() map[string]*queryresult.ColumnDef {
	var dimensionsSchema = make(map[string]*queryresult.ColumnDef)

	for _, row := range r.Rows {
		for _, dim := range row.Dimensions {
			if _, ok := dimensionsSchema[dim.Key]; !ok {
				// add to map
				dimensionsSchema[dim.Key] = &queryresult.ColumnDef{
					Name:     dim.Key,
					DataType: dim.SqlType,
				}
				// also add to DimensionKeys
				r.DimensionKeys = append(r.DimensionKeys, dim.Key)
			}
		}
	}
	// add keys to group
	for _, parent := range r.Parents {
		parent.addDimensionKeys(r.DimensionKeys...)
	}
	return dimensionsSchema
}

// add the result row to our results and update the summary with the row status
func (r *ControlRun) addResultRow(row *ResultRow) {
	// update results
	r.rowMap[row.Status] = append(r.rowMap[row.Status], row)

	// update summary
	switch row.Status {
	case constants.ControlOk:
		r.Summary.Ok++
	case constants.ControlAlarm:
		r.Summary.Alarm++
	case constants.ControlSkip:
		r.Summary.Skip++
	case constants.ControlInfo:
		r.Summary.Info++
	case constants.ControlError:
		r.Summary.Error++
	}
}

// populate ordered list of rows
func (r *ControlRun) createdOrderedResultRows() {
	statusOrder := []string{constants.ControlError, constants.ControlAlarm, constants.ControlInfo, constants.ControlOk, constants.ControlSkip}
	for _, status := range statusOrder {
		r.Rows = append(r.Rows, r.rowMap[status]...)
	}
}

func (r *ControlRun) setRunStatus(ctx context.Context, status dashboardtypes.RunStatus) {
	r.stateLock.Lock()
	r.RunStatus = status
	r.stateLock.Unlock()

	if r.Finished() {
		// close the doneChan - we don't need it anymore
		close(r.doneChan)
	}
}

func (r *ControlRun) trySetStateRunning(ctx context.Context) bool {
	// lock the statuslock
	r.stateLock.Lock()

	defer r.stateLock.Unlock()

	// check the status - if we are not ready(i.e we are running already), return
	if r.RunStatus != dashboardtypes.RunInitialized {
		slog.Debug("control run is not initialized", "name", r.Control.Name(), "status", r.RunStatus)
		return false
	}

	// set status to running and set start time
	r.RunStatus = dashboardtypes.RunRunning
	r.startTime = time.Now()

	return true
}

func (r *ControlRun) populateProperties() error {
	if r.Control == nil {
		return nil
	}
	properties, err := snapshot.GetAsSnapshotPropertyMap(r.Control)
	if err != nil {
		return err

	}
	r.Properties = properties
	return nil
}
