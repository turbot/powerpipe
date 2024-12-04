package resources

import (
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
	DashboardLeafNodeImpl
	// NOTE: we must have cty tag on at least one property otherwise gohcl.DecodeExpression panics
	NodeAndEdgeProviderImpl `cty:"node_and_edge_provider"`

	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`

	Base *DashboardFlow `hcl:"base" json:"-"`
}

func NewDashboardFlow(block *hcl.Block, mod *modconfig.Mod, shortName string) modconfig.HclResource {
	f := &DashboardFlow{
		NodeAndEdgeProviderImpl: NewNodeAndEdgeProviderImpl(block, mod, shortName),
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
