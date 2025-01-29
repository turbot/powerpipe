package resources

import (
	"github.com/turbot/pipe-fittings/v2/modconfig"

	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/pipe-fittings/v2/cty_helpers"
	"github.com/turbot/pipe-fittings/v2/printers"
	"github.com/turbot/pipe-fittings/v2/utils"
	"github.com/zclconf/go-cty/cty"
)

// DashboardHierarchy is a struct representing a leaf dashboard node
type DashboardHierarchy struct {
	modconfig.ResourceWithMetadataImpl
	DashboardLeafNodeImpl
	// NOTE: we must have cty tag on at least one property otherwise gohcl.DecodeExpression panics
	NodeAndEdgeProviderImpl `cty:"node_and_edge_provider"`

	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`

	Base *DashboardHierarchy `hcl:"base" json:"-"`
}

func NewDashboardHierarchy(block *hcl.Block, mod *modconfig.Mod, shortName string) modconfig.HclResource {
	h := &DashboardHierarchy{
		NodeAndEdgeProviderImpl: NewNodeAndEdgeProviderImpl(block, mod, shortName),
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
