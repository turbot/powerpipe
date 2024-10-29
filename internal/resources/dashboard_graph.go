package resources

import (
	"fmt"
	"github.com/turbot/pipe-fittings/modconfig"

	"github.com/hashicorp/hcl/v2"
	typehelpers "github.com/turbot/go-kit/types"
	"github.com/turbot/pipe-fittings/cty_helpers"
	"github.com/turbot/pipe-fittings/printers"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/zclconf/go-cty/cty"
)

// DashboardGraph is a struct representing a leaf dashboard node
type DashboardGraph struct {
	modconfig.ResourceWithMetadataImpl
	QueryProviderImpl
	WithProviderImpl

	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`

	Nodes     DashboardNodeList `cty:"node_list" json:"nodes,omitempty"`
	Edges     DashboardEdgeList `cty:"edge_list" json:"edges,omitempty"`
	NodeNames []string          `snapshot:"nodes"`
	EdgeNames []string          `snapshot:"edges"`

	Categories map[string]*DashboardCategory `cty:"categories" json:"categories,omitempty" snapshot:"categories"`
	Direction  *string                       `cty:"direction" hcl:"direction" json:"direction,omitempty" snapshot:"direction"`

	// these properties are JSON serialised by the parent LeafRun
	Width   *int    `cty:"width" hcl:"width"  json:"width,omitempty"`
	Type    *string `cty:"type" hcl:"type"  json:"type,omitempty"`
	Display *string `cty:"display" hcl:"display" json:"display,omitempty"`

	Base *DashboardGraph `hcl:"base" json:"-"`
}

func NewDashboardGraph(block *hcl.Block, mod *modconfig.Mod, shortName string) modconfig.HclResource {
	g := &DashboardGraph{
		Categories:        make(map[string]*DashboardCategory),
		QueryProviderImpl: NewQueryProviderImpl(block, mod, shortName),
	}
	g.SetAnonymous(block)
	return g
}

func (g *DashboardGraph) Equals(other *DashboardGraph) bool {
	diff := g.Diff(other)
	return !diff.HasChanges()
}

// OnDecoded implements HclResource
func (g *DashboardGraph) OnDecoded(block *hcl.Block, resourceMapProvider modconfig.ResourceMapsProvider) hcl.Diagnostics {
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

// GetWidth implements DashboardLeafNode
func (g *DashboardGraph) GetWidth() int {
	if g.Width == nil {
		return 0
	}
	return *g.Width
}

// GetDisplay implements DashboardLeafNode
func (g *DashboardGraph) GetDisplay() string {
	return typehelpers.SafeString(g.Display)
}

// GetType implements DashboardLeafNode
func (g *DashboardGraph) GetType() string {
	return typehelpers.SafeString(g.Type)
}

// GetEdges implements NodeAndEdgeProvider
func (g *DashboardGraph) GetEdges() DashboardEdgeList {
	return g.Edges
}

// GetNodes implements NodeAndEdgeProvider
func (g *DashboardGraph) GetNodes() DashboardNodeList {
	return g.Nodes
}

// SetEdges implements NodeAndEdgeProvider
func (g *DashboardGraph) SetEdges(edges DashboardEdgeList) {
	g.Edges = edges
}

// SetNodes implements NodeAndEdgeProvider
func (g *DashboardGraph) SetNodes(nodes DashboardNodeList) {
	g.Nodes = nodes
}

// AddCategory implements NodeAndEdgeProvider
func (g *DashboardGraph) AddCategory(category *DashboardCategory) hcl.Diagnostics {
	categoryName := category.ShortName
	if _, ok := g.Categories[categoryName]; ok {
		return hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("%s has duplicate category %s", g.Name(), categoryName),
			Subject:  category.GetDeclRange(),
		}}
	}
	g.Categories[categoryName] = category
	return nil
}

// AddChild implements NodeAndEdgeProvider
func (g *DashboardGraph) AddChild(child modconfig.HclResource) hcl.Diagnostics {
	var diags hcl.Diagnostics
	switch c := child.(type) {
	case *DashboardNode:
		g.Nodes = append(g.Nodes, c)
	case *DashboardEdge:
		g.Edges = append(g.Edges, c)
	default:
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("DashboardGraph does not support children of type %s", child.GetBlockType()),
			Subject:  g.GetDeclRange(),
		})
		return diags
	}
	// set ourselves as parent
	err := child.(modconfig.ModTreeItem).AddParent(g)
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
