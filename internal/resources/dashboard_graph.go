package resources

import (
	"github.com/turbot/pipe-fittings/v2/modconfig"

	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/pipe-fittings/v2/cty_helpers"
	"github.com/turbot/pipe-fittings/v2/printers"
	"github.com/turbot/pipe-fittings/v2/utils"
	"github.com/zclconf/go-cty/cty"
)

// DashboardGraph is a struct representing a leaf dashboard node
type DashboardGraph struct {
	modconfig.ResourceWithMetadataImpl
	WithProviderImpl
	DashboardLeafNodeImpl
	// NOTE: we must have cty tag on at least one property otherwise gohcl.DecodeExpression panics
	NodeAndEdgeProviderImpl `cty:"node_and_edge_provider"`

	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`

	Direction *string `cty:"direction" hcl:"direction" json:"direction,omitempty" snapshot:"direction"`

	Base *DashboardGraph `hcl:"base" json:"-"`
}

func NewDashboardGraph(block *hcl.Block, mod *modconfig.Mod, shortName string) modconfig.HclResource {
	g := &DashboardGraph{
		NodeAndEdgeProviderImpl: NewNodeAndEdgeProviderImpl(block, mod, shortName),
	}
	g.SetAnonymous(block)
	return g
}

func (g *DashboardGraph) Equals(other *DashboardGraph) bool {
	diff := g.Diff(other)
	return !diff.HasChanges()
}

// OnDecoded implements HclResource
func (g *DashboardGraph) OnDecoded(block *hcl.Block, resourceMapProvider modconfig.ModResourcesProvider) hcl.Diagnostics {
	g.SetBaseProperties()
	if len(g.Nodes) > 0 {
		g.NodeNames = g.Nodes.Names()
	}
	if len(g.Edges) > 0 {
		g.EdgeNames = g.Edges.Names()
	}
	return g.QueryProviderImpl.OnDecoded(block, resourceMapProvider)
}

// TODO [node_reuse] Add DashboardLeafNodeImpl and move this there https://github.com/turbot/steampipe/issues/2926

// GetChildren implements ModTreeItem
func (g *DashboardGraph) GetChildren() []modconfig.ModTreeItem {
	// return nodes and edges (if any)
	children := make([]modconfig.ModTreeItem, len(g.Nodes)+len(g.Edges))
	for i, n := range g.Nodes {
		children[i] = n
	}
	offset := len(g.Nodes)
	for i, e := range g.Edges {
		children[i+offset] = e
	}
	return children
}

func (g *DashboardGraph) Diff(other *DashboardGraph) *modconfig.ModTreeItemDiffs {
	res := &modconfig.ModTreeItemDiffs{
		Item: g,
		Name: g.Name(),
	}

	if !utils.SafeStringsEqual(g.Type, other.Type) {
		res.AddPropertyDiff("Type")
	}

	if !utils.SafeStringsEqual(g.Direction, other.Direction) {
		res.AddPropertyDiff("Direction")
	}

	if len(g.Categories) != len(other.Categories) {
		res.AddPropertyDiff("Categories")
	} else {
		for name, c := range g.Categories {
			if !c.Equals(other.Categories[name]) {
				res.AddPropertyDiff("Categories")
			}
		}
	}

	res.PopulateChildDiffs(g, other)
	res.Merge(g.QueryProviderImpl.Diff(other))
	res.Merge(dashboardLeafNodeDiff(g, other))

	return res
}

// CtyValue implements CtyValueProvider
func (g *DashboardGraph) CtyValue() (cty.Value, error) {
	return cty_helpers.GetCtyValue(g)
}

func (g *DashboardGraph) SetBaseProperties() {
	if g.Base == nil {
		return
	}
	// copy base into the HclResourceImpl 'base' property so it is accessible to all nested structs
	g.HclResourceImpl.SetBase(g.Base)

	// call into parent nested struct SetBaseProperties
	g.QueryProviderImpl.SetBaseProperties()

	if g.Type == nil {
		g.Type = g.Base.Type
	}

	if g.Display == nil {
		g.Display = g.Base.Display
	}

	if g.Width == nil {
		g.Width = g.Base.Width
	}

	if g.Categories == nil {
		g.Categories = g.Base.Categories
	} else {
		g.Categories = utils.MergeMaps(g.Categories, g.Base.Categories)
	}

	if g.Direction == nil {
		g.Direction = g.Base.Direction
	}

	if g.Edges == nil {
		g.Edges = g.Base.Edges
	} else {
		g.Edges.Merge(g.Base.Edges)
	}

	if g.Nodes == nil {
		g.Nodes = g.Base.Nodes
	} else {
		g.Nodes.Merge(g.Base.Nodes)
	}
}

// GetShowData implements printers.Showable
func (g *DashboardGraph) GetShowData() *printers.RowData {
	res := printers.NewRowData(
		printers.NewFieldValue("Width", g.Width),
		printers.NewFieldValue("Type", g.Type),
		printers.NewFieldValue("Display", g.Display),
		printers.NewFieldValue("Nodes", g.Nodes),
		printers.NewFieldValue("Edges", g.Edges),
		printers.NewFieldValue("Direction", g.Direction),
	)
	// merge fields from base, putting base fields first
	res.Merge(g.QueryProviderImpl.GetShowData())
	return res
}
