package controlstatus

import "github.com/turbot/powerpipe/internal/dashboardtypes"

// ControlRunStatusProvider is an interface used to allow us to pass a control as the payload of ControlComplete and ControlError events -
type ControlRunStatusProvider interface {
	GetControlId() string
	GetRunStatus() dashboardtypes.RunStatus
	GetStatusSummary() *StatusSummary
}
