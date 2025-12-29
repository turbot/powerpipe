package workspace

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hashicorp/hcl/v2"
	typehelpers "github.com/turbot/go-kit/types"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/pipe-fittings/v2/workspace"
	"github.com/turbot/powerpipe/internal/dashboardevents"
	"github.com/turbot/powerpipe/internal/resources"
)

type PowerpipeWorkspace struct {
	workspace.Workspace
	// event handlers
	dashboardEventHandlers []dashboardevents.DashboardEventHandler
	// channel used to send dashboard events to the handleDashboardEvent goroutine
	dashboardEventChan chan dashboardevents.DashboardEvent
}

func NewPowerpipeWorkspace(workspacePath string) *PowerpipeWorkspace {
	w := &PowerpipeWorkspace{
		Workspace: workspace.Workspace{
			Path:              workspacePath,
			VariableValues:    make(map[string]string),
			ValidateVariables: true,
			Mod:               modconfig.NewMod("local", workspacePath, hcl.Range{}),
		},
	}

	w.OnFileWatcherError = func(ctx context.Context, err error) {
		w.PublishDashboardEvent(ctx, &dashboardevents.WorkspaceError{Error: err})
	}
	w.OnFileWatcherEvent = func(ctx context.Context, modResources, prevModResources modconfig.ModResources) {
		w.raiseDashboardChangedEvents(ctx, modResources, prevModResources)
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
	for _, d := range w.Mod.GetModResources().(*resources.PowerpipeModResources).Dashboards {
		if err := d.ValidateRuntimeDependencies(w); err != nil {
			return err
		}
	}
	return nil
}

// ResolveQueryFromQueryProvider resolves the query for the given QueryProvider
func (w *PowerpipeWorkspace) ResolveQueryFromQueryProvider(queryProvider resources.QueryProvider, runtimeArgs *resources.QueryArgs) (*modconfig.ResolvedQuery, error) {
	slog.Debug("ResolveQueryFromQueryProvider", "resourceName", queryProvider.Name())

	query := queryProvider.GetQuery()
	sql := queryProvider.GetSQL()

	params := queryProvider.GetParams()

	// merge the base args with the runtime args
	var err error
	runtimeArgs, err = resources.MergeArgs(queryProvider, runtimeArgs)
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

func (w *PowerpipeWorkspace) GetQueryProvider(queryName string) (resources.QueryProvider, bool) {
	parsedName, err := modconfig.ParseResourceName(queryName)
	if err != nil {
		return nil, false
	}
	// try to find the resource
	if resource, ok := w.GetResource(parsedName); ok {
		// found a resource - is it a query provider
		if qp := resource.(resources.QueryProvider); ok {
			return qp, true
		}
		slog.Debug("GetQueryProviderImpl found a mod resource resource for query but it is not a query provider", "resourceName", queryName)
	}

	return nil, false
}

// GetPowerpipeModResources returns the powerpipe PowerpipeModResources from the workspace, cast to the correct type
func (w *PowerpipeWorkspace) GetPowerpipeModResources() *resources.PowerpipeModResources {
	modResources, ok := w.GetModResources().(*resources.PowerpipeModResources)
	if !ok {
		// should never happen
		panic(fmt.Sprintf("mod.GetModResources() did not return a powerpipe PowerpipeModResources: %T", w.GetModResources()))
	}
	return modResources
}

// IsLazy returns false as this is not a lazy-loading workspace.
func (w *PowerpipeWorkspace) IsLazy() bool {
	return false
}

// LoadDashboard loads a dashboard by name from the already-loaded resources.
func (w *PowerpipeWorkspace) LoadDashboard(ctx context.Context, name string) (*resources.Dashboard, error) {
	modResources := w.GetPowerpipeModResources()
	if dash, ok := modResources.Dashboards[name]; ok {
		return dash, nil
	}
	return nil, fmt.Errorf("dashboard not found: %s", name)
}

// LoadBenchmark loads a benchmark by name from the already-loaded resources.
func (w *PowerpipeWorkspace) LoadBenchmark(ctx context.Context, name string) (modconfig.ModTreeItem, error) {
	modResources := w.GetPowerpipeModResources()
	if bench, ok := modResources.ControlBenchmarks[name]; ok {
		return bench, nil
	}
	if bench, ok := modResources.DetectionBenchmarks[name]; ok {
		return bench, nil
	}
	return nil, fmt.Errorf("benchmark not found: %s", name)
}
