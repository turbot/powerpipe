package resources

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/pipe-fittings/v2/modconfig"
)

type WithProvider interface {
	AddWith(with *DashboardWith) hcl.Diagnostics
	GetWiths() []*DashboardWith
	GetWith(string) (*DashboardWith, bool)
}

// NodeAndEdgeProvider must be implemented by any dashboard leaf node which supports edges and nodes
// (DashboardGraph, DashboardFlow, DashboardHierarchy)
// TODO [node_reuse] add NodeAndEdgeProviderImpl https://github.com/turbot/steampipe/issues/2918
type NodeAndEdgeProvider interface {
	QueryProvider
	WithProvider
	GetEdges() DashboardEdgeList
	SetEdges(DashboardEdgeList)
	GetNodes() DashboardNodeList
	SetNodes(DashboardNodeList)
	AddCategory(category *DashboardCategory) hcl.Diagnostics
	AddChild(child modconfig.HclResource) hcl.Diagnostics
}

// RuntimeDependencyProvider is implemented by all QueryProviders and Dashboard
type RuntimeDependencyProvider interface {
	modconfig.ModTreeItem
	AddRuntimeDependencies([]*RuntimeDependency)
	GetRuntimeDependencies() map[string]*RuntimeDependency
}

// QueryProvider must be implemented by resources which have query/sql
type QueryProvider interface {
	RuntimeDependencyProvider
	GetArgs() *QueryArgs
	GetParams() []*modconfig.ParamDef
	GetSQL() *string
	GetQuery() *Query
	SetArgs(*QueryArgs)
	SetParams([]*modconfig.ParamDef)
	GetResolvedQuery(*QueryArgs) (*modconfig.ResolvedQuery, error)
	RequiresExecution(QueryProvider) bool
	ValidateQuery() hcl.Diagnostics
	MergeParentArgs(QueryProvider, QueryProvider) hcl.Diagnostics
	GetQueryProviderImpl() *QueryProviderImpl
	ParamsInheritedFromBase() bool
	ArgsInheritedFromBase() bool
}

// DashboardLeafNode must be implemented by resources may be a leaf node in the dashboard execution tree
type DashboardLeafNode interface {
	modconfig.ModTreeItem
	modconfig.ResourceWithMetadata
	GetDisplay() string
	GetType() string
	GetWidth() int
}
