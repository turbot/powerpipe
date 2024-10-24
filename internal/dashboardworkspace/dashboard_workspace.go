package dashboardworkspace

import (
	"context"
	"github.com/turbot/pipe-fittings/modconfig/powerpipe"
	"log/slog"

	"github.com/turbot/pipe-fittings/workspace"
	"github.com/turbot/powerpipe/internal/dashboardevents"
)

// WorkspaceEvents is a wrapper around workspace.WorkspaceEvents that adds dashboard specific event handling
type WorkspaceEvents struct {
	*workspace.Workspace
	// event handlers
	dashboardEventHandlers []dashboardevents.DashboardEventHandler
	// channel used to send dashboard events to the handleDashboardEvent goroutine
	dashboardEventChan chan dashboardevents.DashboardEvent
}

func NewWorkspaceEvents(workspace *workspace.Workspace) *WorkspaceEvents {
	w := &WorkspaceEvents{
		Workspace: workspace,
	}

	w.OnFileWatcherError = func(ctx context.Context, err error) {
		w.PublishDashboardEvent(ctx, &dashboardevents.WorkspaceError{Error: err})
	}
	w.OnFileWatcherEvent = func(ctx context.Context, resourceMaps, prevResourceMaps *powerpipe.PowerpipeResourceMaps) {
		w.raiseDashboardChangedEvents(ctx, resourceMaps, prevResourceMaps)
	}
	return w
}
func (w *WorkspaceEvents) Close() {
	w.Workspace.Close()
	if ch := w.dashboardEventChan; ch != nil {
		// NOTE: set nil first
		w.dashboardEventChan = nil
		slog.Debug("closing dashboardEventChan")
		close(ch)
	}
}
