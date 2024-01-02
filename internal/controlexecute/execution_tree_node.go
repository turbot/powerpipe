package controlexecute

import (
	"github.com/turbot/pipe-fittings/steampipeconfig"
)

// ExecutionTreeNode is implemented by all control execution tree nodes
type ExecutionTreeNode interface {
	IsExecutionTreeNode()
	GetChildren() []ExecutionTreeNode
	GetName() string
	AsTreeNode() *steampipeconfig.SnapshotTreeNode
}
