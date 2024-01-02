package dashboardexecute

import (
	"github.com/turbot/powerpipe/internal/dashboardtypes"
)

type RuntimeDependencyPublishTarget struct {
	transform func(*dashboardtypes.ResolvedRuntimeDependencyValue) *dashboardtypes.ResolvedRuntimeDependencyValue
	channel   chan *dashboardtypes.ResolvedRuntimeDependencyValue
}
