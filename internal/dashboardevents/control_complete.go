package dashboardevents

import (
	"time"

	"github.com/turbot/powerpipe/internal/controlstatus"
)

type ControlComplete struct {
	Progress    *controlstatus.ControlProgress
	Control     controlstatus.ControlRunStatusProvider
	Name        string
	Session     string
	ExecutionId string
	Timestamp   time.Time
}

// IsDashboardEvent implements DashboardEvent interface
func (*ControlComplete) IsDashboardEvent() {}
