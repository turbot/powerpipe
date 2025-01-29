package resources

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/pipe-fittings/v2/modconfig"
)

type NodeAndEdgeProviderImpl struct {
	WithProviderImpl
	QueryProviderImpl

	NodeEdgeProviderRemain hcl.Body          `hcl:",remain" json:"-"`
	Nodes                  DashboardNodeList `cty:"node_list"  json:"nodes,omitempty"`
	Edges                  DashboardEdgeList `cty:"edge_list" json:"edges,omitempty"`
	NodeNames              []string          `json:"-" snapshot:"nodes"`
	EdgeNames              []string          `json:"-" snapshot:"edges"`

	Categories map[string]*DashboardCategory `cty:"categories" json:"categories" snapshot:"categories"`
}

func NewNodeAndEdgeProviderImpl(block *hcl.Block, mod *modconfig.Mod, name string) NodeAndEdgeProviderImpl {
	return NodeAndEdgeProviderImpl{
		Categories:        make(map[string]*DashboardCategory),
		QueryProviderImpl: NewQueryProviderImpl(block, mod, name),
	}
}

// GetEdges implements NodeAndEdgeProvider
func (f *NodeAndEdgeProviderImpl) GetEdges() DashboardEdgeList {
	return f.Edges
}

// GetNodes implements NodeAndEdgeProvider
func (f *NodeAndEdgeProviderImpl) GetNodes() DashboardNodeList {
	return f.Nodes
}

// SetEdges implements NodeAndEdgeProvider
func (f *NodeAndEdgeProviderImpl) SetEdges(edges DashboardEdgeList) {
	f.Edges = edges
}

// SetNodes implements NodeAndEdgeProvider
func (f *NodeAndEdgeProviderImpl) SetNodes(nodes DashboardNodeList) {
	f.Nodes = nodes
}

// AddCategory implements NodeAndEdgeProvider
func (f *NodeAndEdgeProviderImpl) AddCategory(category *DashboardCategory) hcl.Diagnostics {
	categoryName := category.ShortName
	if _, ok := f.Categories[categoryName]; ok {
		return hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("%s has duplicate category %s", f.Name(), categoryName),
			Subject:  category.GetDeclRange(),
		}}
	}
	f.Categories[categoryName] = category
	return nil
}

// AddChild implements NodeAndEdgeProvider
func (f *NodeAndEdgeProviderImpl) AddChild(child modconfig.HclResource) hcl.Diagnostics {
	var diags hcl.Diagnostics
	switch c := child.(type) {
	case *DashboardNode:
		f.Nodes = append(f.Nodes, c)
	case *DashboardEdge:
		f.Edges = append(f.Edges, c)
	default:
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("NodeAndEdgeProviderImpl does not support children of type %s", child.GetBlockType()),
			Subject:  f.GetDeclRange(),
		})
		return diags
	}
	// set ourselves as parent
	err := child.(modconfig.ModTreeItem).AddParent(f)
	if err != nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "failed to add parent to ModTreeItem",
			Detail:   err.Error(),
			Subject:  child.GetDeclRange(),
		})
	}

	return diags
}
