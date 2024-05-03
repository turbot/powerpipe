package controlexecute

import (
	"github.com/turbot/pipe-fittings/v2/steampipeconfig"
)

// ExecutionTreeNode is implemented by all control execution tree nodes
type ExecutionTreeNode interface {
	IsExecutionTreeNode()
	GetChildren() []ExecutionTreeNode
	GetName() string
	AsTreeNode() *steampipeconfig.SnapshotTreeNode
}
