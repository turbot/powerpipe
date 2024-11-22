package resources

import (
	"fmt"
	"github.com/turbot/pipe-fittings/modconfig"

	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/pipe-fittings/cty_helpers"
	"github.com/turbot/pipe-fittings/printers"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/zclconf/go-cty/cty"
)

// DashboardFlow is a struct representing a leaf dashboard node
type DashboardFlow struct {
	modconfig.ResourceWithMetadataImpl
	QueryProviderImpl
	WithProviderImpl
	DashboardLeafNodeImpl

	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`

	Nodes     DashboardNodeList `cty:"node_list"  json:"-"`
	Edges     DashboardEdgeList `cty:"edge_list" json:"-"`
	NodeNames []string          `json:"nodes" snapshot:"nodes"`
	EdgeNames []string          `json:"edges" snapshot:"edges"`

	Categories map[string]*DashboardCategory `cty:"categories" json:"categories" snapshot:"categories"`

	Base *DashboardFlow `hcl:"base" json:"-"`
}

func NewDashboardFlow(block *hcl.Block, mod *modconfig.Mod, shortName string) modconfig.HclResource {
	f := &DashboardFlow{
		Categories:        make(map[string]*DashboardCategory),
		QueryProviderImpl: NewQueryProviderImpl(block, mod, shortName),
	}
	f.SetAnonymous(block)
	return f
}

func (f *DashboardFlow) Equals(other *DashboardFlow) bool {
	diff := f.Diff(other)
	return !diff.HasChanges()
}

// OnDecoded implements HclResource
func (f *DashboardFlow) OnDecoded(block *hcl.Block, resourceMapProvider modconfig.ModResourcesProvider) hcl.Diagnostics {
	f.SetBaseProperties()
	if len(f.Nodes) > 0 {
		f.NodeNames = f.Nodes.Names()
	}
	if len(f.Edges) > 0 {
		f.EdgeNames = f.Edges.Names()
	}
	return f.QueryProviderImpl.OnDecoded(block, resourceMapProvider)
}

// TODO [node_reuse] Add DashboardLeafNodeImpl and move this there https://github.com/turbot/steampipe/issues/2926
// GetChildren implements ModTreeItem
func (f *DashboardFlow) GetChildren() []modconfig.ModTreeItem {
	// return nodes and edges (if any)
	children := make([]modconfig.ModTreeItem, len(f.Nodes)+len(f.Edges))
	for i, n := range f.Nodes {
		children[i] = n
	}
	offset := len(f.Nodes)
	for i, e := range f.Edges {
		children[i+offset] = e
	}
	return children
}

func (f *DashboardFlow) Diff(other *DashboardFlow) *modconfig.ModTreeItemDiffs {
	res := &modconfig.ModTreeItemDiffs{
		Item: f,
		Name: f.Name(),
	}

	if !utils.SafeStringsEqual(f.Type, other.Type) {
		res.AddPropertyDiff("Type")
	}

	if len(f.Categories) != len(other.Categories) {
		res.AddPropertyDiff("Categories")
	} else {
		for name, c := range f.Categories {
			if !c.Equals(other.Categories[name]) {
				res.AddPropertyDiff("Categories")
			}
		}
	}

	res.PopulateChildDiffs(f, other)
	res.Merge(f.QueryProviderImpl.Diff(other))
	res.Merge(dashboardLeafNodeDiff(f, other))

	return res
}

// ValidateQuery implements QueryProvider
func (*DashboardFlow) ValidateQuery() hcl.Diagnostics {
	// query is optional - nothing to do
	return nil
}

// GetEdges implements NodeAndEdgeProvider
func (f *DashboardFlow) GetEdges() DashboardEdgeList {
	return f.Edges
}

// GetNodes implements NodeAndEdgeProvider
func (f *DashboardFlow) GetNodes() DashboardNodeList {
	return f.Nodes
}

// SetEdges implements NodeAndEdgeProvider
func (f *DashboardFlow) SetEdges(edges DashboardEdgeList) {
	f.Edges = edges
}

// SetNodes implements NodeAndEdgeProvider
func (f *DashboardFlow) SetNodes(nodes DashboardNodeList) {
	f.Nodes = nodes
}

// AddCategory implements NodeAndEdgeProvider
func (f *DashboardFlow) AddCategory(category *DashboardCategory) hcl.Diagnostics {
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
func (f *DashboardFlow) AddChild(child modconfig.HclResource) hcl.Diagnostics {
	var diags hcl.Diagnostics
	switch c := child.(type) {
	case *DashboardNode:
		f.Nodes = append(f.Nodes, c)
	case *DashboardEdge:
		f.Edges = append(f.Edges, c)
	default:
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("DashboardFlow does not support children of type %s", child.GetBlockType()),
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

// CtyValue implements CtyValueProvider
func (f *DashboardFlow) CtyValue() (cty.Value, error) {
	return cty_helpers.GetCtyValue(f)
}

func (f *DashboardFlow) SetBaseProperties() {
	if f.Base == nil {
		return
	}
	// copy base into the HclResourceImpl 'base' property so it is accessible to all nested structs
	f.HclResourceImpl.SetBase(f.Base)

	// call into parent nested struct SetBaseProperties
	f.QueryProviderImpl.SetBaseProperties()

	if f.Type == nil {
		f.Type = f.Base.Type
	}

	if f.Display == nil {
		f.Display = f.Base.Display
	}

	if f.Width == nil {
		f.Width = f.Base.Width
	}

	if f.Categories == nil {
		f.Categories = f.Base.Categories
	} else {
		f.Categories = utils.MergeMaps(f.Categories, f.Base.Categories)
	}

	if f.Edges == nil {
		f.Edges = f.Base.Edges
	} else {
		f.Edges.Merge(f.Base.Edges)
	}
	if f.Nodes == nil {
		f.Nodes = f.Base.Nodes
	} else {
		f.Nodes.Merge(f.Base.Nodes)
	}
}

// GetShowData implements printers.Showable
func (f *DashboardFlow) GetShowData() *printers.RowData {
	res := printers.NewRowData(
		printers.NewFieldValue("Width", f.Width),
		printers.NewFieldValue("Type", f.Type),
		printers.NewFieldValue("Display", f.Display),
		printers.NewFieldValue("Nodes", f.Nodes),
		printers.NewFieldValue("Edges", f.Edges),
	)
	// merge fields from base, putting base fields first
	res.Merge(f.QueryProviderImpl.GetShowData())
	return res
}
