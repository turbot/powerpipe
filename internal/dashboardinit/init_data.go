package dashboardinit

import (
	"github.com/turbot/pipe-fittings/initialisation"
	"github.com/turbot/powerpipe/internal/dashboardworkspace"
)

// InitData is a wrapper around initialisation.InitData that adds dashboard specific initialisation data
type InitData struct {
	*initialisation.InitData

	DashboardWorkspace *dashboardworkspace.Workspace
}

func NewInitData(i *initialisation.InitData) *InitData {

	return &InitData{
		InitData:           i,
		DashboardWorkspace: dashboardworkspace.NewWorkspace(i.Workspace),
	}
}

func NewErrorDashboardInitData(err error) *InitData {
	return &InitData{
		InitData: initialisation.NewErrorInitData(err),
	}
}
