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

// DashboardHierarchy is a struct representing a leaf dashboard node
type DashboardHierarchy struct {
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
	Width      *int                          `cty:"width" hcl:"width"  json:"width,omitempty"`
	Type       *string                       `cty:"type" hcl:"type"  json:"type,omitempty"`
	Display    *string                       `cty:"display" hcl:"display" json:"display,omitempty"`

	Base *DashboardHierarchy `hcl:"base" json:"-"`
}

func NewDashboardHierarchy(block *hcl.Block, mod *modconfig.Mod, shortName string) modconfig.HclResource {
	h := &DashboardHierarchy{
		Categories:        make(map[string]*DashboardCategory),
		QueryProviderImpl: NewQueryProviderImpl(block, mod, shortName),
	}
	h.SetAnonymous(block)
	return h
}

func (h *DashboardHierarchy) Equals(other *DashboardHierarchy) bool {
	diff := h.Diff(other)
	return !diff.HasChanges()
}

// OnDecoded implements HclResource
func (h *DashboardHierarchy) OnDecoded(block *hcl.Block, resourceMapProvider modconfig.ModResourcesProvider) hcl.Diagnostics {
	h.SetBaseProperties()
	if len(h.Nodes) > 0 {
		h.NodeNames = h.Nodes.Names()
	}
	if len(h.Edges) > 0 {
		h.EdgeNames = h.Edges.Names()
	}
	return h.QueryProviderImpl.OnDecoded(block, resourceMapProvider)
}

// TODO [node_reuse] Add DashboardLeafNodeImpl and move this there https://github.com/turbot/steampipe/issues/2926

// GetChildren implements ModTreeItem
func (h *DashboardHierarchy) GetChildren() []modconfig.ModTreeItem {
	// return nodes and edges (if any)
	children := make([]modconfig.ModTreeItem, len(h.Nodes)+len(h.Edges))
	for i, n := range h.Nodes {
		children[i] = n
	}
	offset := len(h.Nodes)
	for i, e := range h.Edges {
		children[i+offset] = e
	}
	return children
}

func (h *DashboardHierarchy) Diff(other *DashboardHierarchy) *modconfig.ModTreeItemDiffs {
	res := &modconfig.ModTreeItemDiffs{
		Item: h,
		Name: h.Name(),
	}

	if !utils.SafeStringsEqual(h.Type, other.Type) {
		res.AddPropertyDiff("Type")
	}

	if len(h.Categories) != len(other.Categories) {
		res.AddPropertyDiff("Categories")
	} else {
		for name, c := range h.Categories {
			if !c.Equals(other.Categories[name]) {
				res.AddPropertyDiff("Categories")
			}
		}
	}

	res.PopulateChildDiffs(h, other)
	res.Merge(h.QueryProviderImpl.Diff(other))
	res.Merge(dashboardLeafNodeDiff(h, other))

	return res
}

// GetWidth implements DashboardLeafNode
func (h *DashboardHierarchy) GetWidth() int {
	if h.Width == nil {
		return 0
	}
	return *h.Width
}

// GetDisplay implements DashboardLeafNode
func (h *DashboardHierarchy) GetDisplay() string {
	return typehelpers.SafeString(h.Display)
}

// GetDocumentation implements DashboardLeafNode, ModTreeItem
func (h *DashboardHierarchy) GetDocumentation() string {
	return ""
}

// GetType implements DashboardLeafNode
func (h *DashboardHierarchy) GetType() string {
	return typehelpers.SafeString(h.Type)
}

// GetEdges implements NodeAndEdgeProvider
func (h *DashboardHierarchy) GetEdges() DashboardEdgeList {
	return h.Edges
}

// GetNodes implements NodeAndEdgeProvider
func (h *DashboardHierarchy) GetNodes() DashboardNodeList {
	return h.Nodes
}

// SetEdges implements NodeAndEdgeProvider
func (h *DashboardHierarchy) SetEdges(edges DashboardEdgeList) {
	h.Edges = edges
}

// SetNodes implements NodeAndEdgeProvider
func (h *DashboardHierarchy) SetNodes(nodes DashboardNodeList) {
	h.Nodes = nodes
}

// AddCategory implements NodeAndEdgeProvider
func (h *DashboardHierarchy) AddCategory(category *DashboardCategory) hcl.Diagnostics {
	categoryName := category.ShortName
	if _, ok := h.Categories[categoryName]; ok {
		return hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("%s has duplicate category %s", h.Name(), categoryName),
			Subject:  category.GetDeclRange(),
		}}
	}
	h.Categories[categoryName] = category
	return nil
}

// AddChild implements NodeAndEdgeProvider
func (h *DashboardHierarchy) AddChild(child modconfig.HclResource) hcl.Diagnostics {
	var diags hcl.Diagnostics
	switch c := child.(type) {
	case *DashboardNode:
		h.Nodes = append(h.Nodes, c)
	case *DashboardEdge:
		h.Edges = append(h.Edges, c)
	default:
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("DashboardHierarchy does not support children of type %s", child.GetBlockType()),
			Subject:  h.GetDeclRange(),
		})
		return diags
	}
	// set ourselves as parent
	err := child.(modconfig.ModTreeItem).AddParent(h)
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
func (h *DashboardHierarchy) CtyValue() (cty.Value, error) {
	return cty_helpers.GetCtyValue(h)
}

func (h *DashboardHierarchy) SetBaseProperties() {
	if h.Base == nil {
		return
	}
	// copy base into the HclResourceImpl 'base' property so it is accessible to all nested structs
	h.HclResourceImpl.SetBase(h.Base)

	// call into parent nested struct SetBaseProperties
	h.QueryProviderImpl.SetBaseProperties()

	if h.Type == nil {
		h.Type = h.Base.Type
	}

	if h.Display == nil {
		h.Display = h.Base.Display
	}

	if h.Width == nil {
		h.Width = h.Base.Width
	}

	if h.Categories == nil {
		h.Categories = h.Base.Categories
	} else {
		h.Categories = utils.MergeMaps(h.Categories, h.Base.Categories)
	}

	if h.Edges == nil {
		h.Edges = h.Base.Edges
	} else {
		h.Edges.Merge(h.Base.Edges)
	}

	if h.Nodes == nil {
		h.Nodes = h.Base.Nodes
	} else {
		h.Nodes.Merge(h.Base.Nodes)
	}
}

// GetShowData implements printers.Showable
func (h *DashboardHierarchy) GetShowData() *printers.RowData {
	res := printers.NewRowData(
		printers.NewFieldValue("Width", h.Width),
		printers.NewFieldValue("Type", h.Type),
		printers.NewFieldValue("Display", h.Display),
		printers.NewFieldValue("Nodes", h.Nodes),
		printers.NewFieldValue("Edges", h.Edges),
	)
	// merge fields from base, putting base fields first
	res.Merge(h.QueryProviderImpl.GetShowData())
	return res
}
