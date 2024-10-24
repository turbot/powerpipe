package workspace

import (
	"context"
	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/pipe-fittings/modconfig/powerpipe"
	"github.com/turbot/pipe-fittings/workspace"
	"github.com/turbot/powerpipe/internal/dashboardevents"
	"log/slog"
)

type PowerpipeWorkspace struct {
	workspace.WorkspaceBase[*powerpipe.PowerpipeResourceMaps]
	// event handlers
	dashboardEventHandlers []dashboardevents.DashboardEventHandler
	// channel used to send dashboard events to the handleDashboardEvent goroutine
	dashboardEventChan chan dashboardevents.DashboardEvent
}

func NewPowerpipeWorkspace(workspacePath string) *PowerpipeWorkspace {
	w := &PowerpipeWorkspace{
		WorkspaceBase: workspace.WorkspaceBase[*powerpipe.PowerpipeResourceMaps]{
			Path:              workspacePath,
			VariableValues:    make(map[string]string),
			ValidateVariables: true,
			Mod:               powerpipe.NewMod("local", workspacePath, hcl.Range{}),
		},
	}

	w.OnFileWatcherError = func(ctx context.Context, err error) {
		w.PublishDashboardEvent(ctx, &dashboardevents.WorkspaceError{Error: err})
	}
	w.OnFileWatcherEvent = func(ctx context.Context, resourceMaps, prevResourceMaps *powerpipe.PowerpipeResourceMaps) {
		w.raiseDashboardChangedEvents(ctx, resourceMaps, prevResourceMaps)
	}
	return w
}

func (w *PowerpipeWorkspace) Close() {
	w.WorkspaceBase.Close()
	if ch := w.dashboardEventChan; ch != nil {
		// NOTE: set nil first
		w.dashboardEventChan = nil
		slog.Debug("closing dashboardEventChan")
		close(ch)
	}
}

func (w *PowerpipeWorkspace) verifyResourceRuntimeDependencies() error {
	for _, d := range w.Mod.GetResourceMaps().(*powerpipe.PowerpipeResourceMaps).Dashboards {
		if err := d.ValidateRuntimeDependencies(w); err != nil {
			return err
		}
	}
	return nil
}
