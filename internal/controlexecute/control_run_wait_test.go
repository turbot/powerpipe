package controlexecute

import (
	"context"
	"testing"

	"github.com/turbot/powerpipe/internal/controlstatus"
	"github.com/turbot/powerpipe/internal/dashboardtypes"
	localqueryresult "github.com/turbot/powerpipe/internal/queryresult"
)

func newWaitingRun(result *localqueryresult.Result) *ControlRun {
	return &ControlRun{
		RunStatus:   dashboardtypes.RunRunning,
		rowMap:      map[string]ResultRows{},
		Summary:     &controlstatus.StatusSummary{},
		queryResult: result,
		doneChan:    make(chan bool, 1),
	}
}

// when the row reader stops because the run's context ended, it closes the result without an error row:
// the run must then be cancelled, never complete with the rows it had. Both cases of the select are ready, so repeat.
func TestWaitForResults_CancelledWhileRowsEnd(t *testing.T) {
	for range 200 {
		result := localqueryresult.NewResult(nil)
		result.Close()
		run := newWaitingRun(result)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		run.waitForResults(ctx)

		if got := run.GetRunStatus(); got != dashboardtypes.RunCanceled {
			t.Fatalf("expected the run to be cancelled, got %v", got)
		}
	}
}

// the rows end with the context alive: the run is complete, as before
func TestWaitForResults_Complete(t *testing.T) {
	result := localqueryresult.NewResult(nil)
	result.Close()
	run := newWaitingRun(result)

	run.waitForResults(context.Background())

	if got := run.GetRunStatus(); got != dashboardtypes.RunComplete {
		t.Fatalf("expected the run to be complete, got %v", got)
	}
}
