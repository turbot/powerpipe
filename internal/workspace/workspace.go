package workspace

import (
	"context"
	"fmt"
	"github.com/hashicorp/hcl/v2"
	typehelpers "github.com/turbot/go-kit/types"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/modconfig/powerpipe"
	"github.com/turbot/pipe-fittings/workspace"
	"github.com/turbot/powerpipe/internal/dashboardevents"
	"log/slog"
)

type PowerpipeWorkspace struct {
	workspace.Workspace[*powerpipe.ModResources]
	// event handlers
	dashboardEventHandlers []dashboardevents.DashboardEventHandler
	// channel used to send dashboard events to the handleDashboardEvent goroutine
	dashboardEventChan chan dashboardevents.DashboardEvent
}

func NewPowerpipeWorkspace(workspacePath string) *PowerpipeWorkspace {
	w := &PowerpipeWorkspace{
		Workspace: workspace.Workspace[*powerpipe.ModResources]{
			Path:              workspacePath,
			VariableValues:    make(map[string]string),
			ValidateVariables: true,
			Mod:               modconfig.NewMod("local", workspacePath, hcl.Range{}),
		},
	}

	w.OnFileWatcherError = func(ctx context.Context, err error) {
		w.PublishDashboardEvent(ctx, &dashboardevents.WorkspaceError{Error: err})
	}
	w.OnFileWatcherEvent = func(ctx context.Context, resourceMaps, prevResourceMaps modconfig.ResourceMapsI) {
		w.raiseDashboardChangedEvents(ctx, resourceMaps, prevResourceMaps)
	}
	return w
}

func (w *PowerpipeWorkspace) Close() {
	w.Workspace.Close()
	if ch := w.dashboardEventChan; ch != nil {
		// NOTE: set nil first
		w.dashboardEventChan = nil
		slog.Debug("closing dashboardEventChan")
		close(ch)
	}
}

func (w *PowerpipeWorkspace) verifyResourceRuntimeDependencies() error {
	for _, d := range w.Mod.GetResourceMaps().(*powerpipe.ModResources).Dashboards {
		if err := d.ValidateRuntimeDependencies(w); err != nil {
			return err
		}
	}
	return nil
}

// ResolveQueryFromQueryProvider resolves the query for the given QueryProvider
func (w *PowerpipeWorkspace) ResolveQueryFromQueryProvider(queryProvider powerpipe.QueryProvider, runtimeArgs *powerpipe.QueryArgs) (*powerpipe.ResolvedQuery, error) {
	slog.Debug("ResolveQueryFromQueryProvider", "resourceName", queryProvider.Name())

	query := queryProvider.GetQuery()
	sql := queryProvider.GetSQL()

	params := queryProvider.GetParams()

	// merge the base args with the runtime args
	var err error
	runtimeArgs, err = powerpipe.MergeArgs(queryProvider, runtimeArgs)
	if err != nil {
		return nil, err
	}

	// determine the source for the query
	// - this will either be the control itself or any named query the control refers to
	// either via its SQL proper ty (passing a query name) or Query property (using a reference to a query object)

	// if a query is provided, use that to resolve the sql
	if query != nil {
		return w.ResolveQueryFromQueryProvider(query, runtimeArgs)
	}

	// must have sql is there is no query
	if sql == nil {
		return nil, fmt.Errorf("%s does not define  either a 'sql' property or a 'query' property\n", queryProvider.Name())
	}

	queryProviderSQL := typehelpers.SafeString(sql)
	slog.Debug("control defines inline SQL")

	// if the SQL refers to a named query, this is the same as if the 'Query' property is set
	if namedQueryProvider, ok := w.GetQueryProvider(queryProviderSQL); ok {
		// in this case, it is NOT valid for the query provider to define its own Param definitions
		if params != nil {
			return nil, fmt.Errorf("%s has an 'SQL' property which refers to %s, so it cannot define 'param' blocks", queryProvider.Name(), namedQueryProvider.Name())
		}
		return w.ResolveQueryFromQueryProvider(namedQueryProvider, runtimeArgs)
	}

	// so the  sql is NOT a named query
	return queryProvider.GetResolvedQuery(runtimeArgs)

}
