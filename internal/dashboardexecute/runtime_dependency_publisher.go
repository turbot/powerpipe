package dashboardexecute

import (
	"github.com/turbot/powerpipe/internal/dashboardtypes"
	"github.com/turbot/powerpipe/internal/resources"
)

type RuntimeDependencyPublisher interface {
	dashboardtypes.DashboardTreeRun
	ProvidesRuntimeDependency(dependency *resources.RuntimeDependency) bool
	SubscribeToRuntimeDependency(name string, opts ...RuntimeDependencyPublishOption) chan *dashboardtypes.ResolvedRuntimeDependencyValue
	PublishRuntimeDependencyValue(name string, result *dashboardtypes.ResolvedRuntimeDependencyValue)
	GetWithRuns() map[string]*LeafRun
}
