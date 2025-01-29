package dashboardevents

import (
	"github.com/turbot/pipe-fittings/v2/utils"
	"time"

	"github.com/turbot/pipe-fittings/v2/steampipeconfig"
	"github.com/turbot/powerpipe/internal/dashboardtypes"
)

type ExecutionComplete struct {
	Root             dashboardtypes.DashboardTreeRun
	Session          string
	ExecutionId      string
	Panels           map[string]steampipeconfig.SnapshotPanel
	Inputs           map[string]interface{}
	Variables        map[string]string
	SearchPath       []string
	SearchPathPrefix []string
	DateTimeRange    utils.TimeRange
	StartTime        time.Time
	EndTime          time.Time
}

// IsDashboardEvent implements DashboardEvent interface
func (*ExecutionComplete) IsDashboardEvent() {}
