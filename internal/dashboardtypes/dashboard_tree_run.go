package dashboardtypes

import (
	"context"

	"github.com/turbot/pipe-fittings/steampipeconfig"
	"github.com/turbot/powerpipe/internal/resources"
)

// DashboardTreeRun is an interface implemented by all dashboard run nodes
type DashboardTreeRun interface {
	Initialise(ctx context.Context)
	Execute(ctx context.Context)
	GetName() string
	GetTitle() string
	GetRunStatus() RunStatus
	SetError(context.Context, error)
	GetError() error
	GetParent() DashboardParent
	SetComplete(context.Context)
	RunComplete() bool
	GetInputsDependingOn(string) []string
	GetNodeType() string
	AsTreeNode() *steampipeconfig.SnapshotTreeNode
	GetResource() resources.DashboardLeafNode
}
