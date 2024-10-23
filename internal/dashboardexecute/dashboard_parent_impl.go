package dashboardexecute

import (
	"context"
	"fmt"
	"github.com/turbot/pipe-fittings/modconfig/dashboard"
	"log/slog"
	"sync"

	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/schema"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/powerpipe/internal/dashboardtypes"
)

type DashboardParentImpl struct {
	DashboardTreeRunImpl
	children          []dashboardtypes.DashboardTreeRun
	childCompleteChan chan dashboardtypes.DashboardTreeRun
	// are we blocked by a child run
	blockedByChild  bool
	childStatusLock *sync.Mutex
}

func newDashboardParentImpl(resource dashboard.DashboardLeafNode, parent dashboardtypes.DashboardParent, run dashboardtypes.DashboardTreeRun, executionTree *DashboardExecutionTree) DashboardParentImpl {
	return DashboardParentImpl{
		DashboardTreeRunImpl: NewDashboardTreeRunImpl(resource, parent, run, executionTree),
		childStatusLock:      new(sync.Mutex),
	}
}

func (r *DashboardParentImpl) initialiseChildren(ctx context.Context) error {
	var errors []error
	for _, child := range r.children {
		child.Initialise(ctx)

		if err := child.GetError(); err != nil {
			errors = append(errors, err)
		}
	}

	return error_helpers.CombineErrors(errors...)

}

// GetChildren implements DashboardTreeRun
func (r *DashboardParentImpl) GetChildren() []dashboardtypes.DashboardTreeRun {
	return r.children
}

// ChildrenComplete implements DashboardTreeRun
func (r *DashboardParentImpl) ChildrenComplete() bool {
	for _, child := range r.children {
		if !child.RunComplete() {
			slog.Debug("ChildrenComplete but child NOT complete state ", "parent", r.Name, "child", child.GetName(), "state", child.GetRunStatus())
			return false
		}
	}

	return true
}

func (r *DashboardParentImpl) ChildCompleteChan() chan dashboardtypes.DashboardTreeRun {
	return r.childCompleteChan
}
func (r *DashboardParentImpl) createChildCompleteChan() {
	// create buffered child complete chan
	if childCount := len(r.children); childCount > 0 {
		r.childCompleteChan = make(chan dashboardtypes.DashboardTreeRun, childCount)
	}
}

// if this leaf run has children (including with runs) execute them asynchronously
func (r *DashboardParentImpl) executeChildrenAsync(ctx context.Context) {
	for _, c := range r.children {
		go c.Execute(ctx)
	}
}

// if this leaf run has with runs execute them asynchronously
func (r *DashboardParentImpl) executeWithsAsync(ctx context.Context) {
	for _, c := range r.children {
		if c.GetNodeType() == schema.BlockTypeWith {
			go c.Execute(ctx)
		}
	}
}

func (r *DashboardParentImpl) waitForChildrenAsync(ctx context.Context) chan error {
	slog.Debug("waitForChildrenAsync", "name", r.Name)
	var doneChan = make(chan error)
	if len(r.children) == 0 {
		slog.Debug("waitForChildrenAsync - no children so we're done", "name", r.Name)
		// if there are no children, return a closed channel so we do not wait
		close(doneChan)
		return doneChan
	}

	go func() {
		// wait for children to complete
		var errors []error
		for !(r.ChildrenComplete()) {
			completeChild := <-r.childCompleteChan
			slog.Debug("waitForChildrenAsync got child complete", "parent", r.Name, "child", completeChild.GetName())
			if completeChild.GetRunStatus().IsError() {
				errors = append(errors, completeChild.GetError())
				slog.Debug("child  has error", "parent", r.Name, "child", completeChild.GetName(), "error", completeChild.GetError())
			}
		}

		slog.Debug("ALL children and withs complete", "name", r.Name, "errors", errors)

		// so all children have completed - check for errors
		var err error
		if len(errors) > 0 {
			err = fmt.Errorf("%d %s failed with an error", len(errors), utils.Pluralize("child", len(errors)))
		}

		// if context is cancelled, just return context cancellation error
		if ctx.Err() != nil {
			if ctx.Err().Error() == context.DeadlineExceeded.Error() {
				err = fmt.Errorf("execution timed out")
			} else {
				err = ctx.Err()
			}
		}

		doneChan <- err
	}()

	return doneChan
}

func (r *DashboardParentImpl) ChildStatusChanged(ctx context.Context) {
	// this function may be called asyncronously by children
	r.childStatusLock.Lock()
	defer r.childStatusLock.Unlock()

	// if we are currently blocked by a child or we are currently in running state,
	// call setRunning() to determine whether any of our children are now blocked
	if r.blockedByChild || r.GetRunStatus() == dashboardtypes.RunRunning {
		slog.Debug("ChildStatusChanged - calling setRunning to see if we are still running", "parent", r.Name, "status", r.GetRunStatus(), "child", r.blockedByChild)

		// try setting our status to running again
		r.setRunning(ctx)
	}
}

// override DashboardTreeRunImpl) setStatus(
func (r *DashboardParentImpl) setRunning(ctx context.Context) {
	// if the run is already complete (for example, canceled), do nothing
	if r.GetRunStatus().IsFinished() {
		slog.Debug("setRunning - run already terminated - NOT setting running", "name", r.Name, "current state", r.GetRunStatus())
		return
	}

	status := dashboardtypes.RunRunning
	// if we are trying to set status to running, check if any of our children are blocked,
	// and if so set our status to blocked

	// if any children are blocked, we are blocked
	for _, c := range r.children {
		if c.GetRunStatus() == dashboardtypes.RunBlocked {
			status = dashboardtypes.RunBlocked
			r.blockedByChild = true
			break
		}
		// to get here, no children can be blocked - clear blockedByChild
		r.blockedByChild = false
	}

	// set status if it has changed
	if status != r.GetRunStatus() {
		slog.Debug("setRunning - setting state, blockedByChild", "name", r.Name, "state", status, "child", r.blockedByChild)
		r.DashboardTreeRunImpl.setStatus(ctx, status)
	} else {
		slog.Debug("setRunning - state unchanged , blockedByChild", "name", r.Name, "state", status, "child", r.blockedByChild)
	}
}
