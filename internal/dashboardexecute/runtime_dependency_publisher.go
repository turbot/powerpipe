package dashboardexecute

import (
	"github.com/turbot/pipe-fittings/modconfig/powerpipe"
	"github.com/turbot/powerpipe/internal/dashboardtypes"
)

type RuntimeDependencyPublisher interface {
	dashboardtypes.DashboardTreeRun
	ProvidesRuntimeDependency(dependency *powerpipe.RuntimeDependency) bool
	SubscribeToRuntimeDependency(name string, opts ...RuntimeDependencyPublishOption) chan *dashboardtypes.ResolvedRuntimeDependencyValue
	PublishRuntimeDependencyValue(name string, result *dashboardtypes.ResolvedRuntimeDependencyValue)
	GetWithRuns() map[string]*LeafRun
}
