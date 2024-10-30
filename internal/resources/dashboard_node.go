package resources

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/pipe-fittings/cty_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/zclconf/go-cty/cty"
)

// DashboardNode is a struct representing a leaf dashboard node
type DashboardNode struct {
	modconfig.ResourceWithMetadataImpl
	QueryProviderImpl

	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`

	Category *DashboardCategory `cty:"category" hcl:"category" json:"category,omitempty" snapshot:"category"`
	Base     *DashboardNode     `hcl:"base" json:"-"`
}

func NewDashboardNode(block *hcl.Block, mod *modconfig.Mod, shortName string) modconfig.HclResource {
	n := &DashboardNode{
		QueryProviderImpl: NewQueryProviderImpl(block, mod, shortName),
	}

	n.SetAnonymous(block)
	return n
}

func (n *DashboardNode) Equals(other *DashboardNode) bool {
	diff := n.Diff(other)
	return !diff.HasChanges()
}

// OnDecoded implements HclResource—
func (n *DashboardNode) OnDecoded(_ *hcl.Block, resourceMapProvider modconfig.ModResourcesProvider) hcl.Diagnostics {
	n.SetBaseProperties()

	// when we reference resources (i.e. category),
	// not all properties are retrieved as they are no cty serialisable
	// repopulate category from resourceMapProvider
	if n.Category != nil {
		fullCategory, diags := enrichCategory(n.Category, n, resourceMapProvider)
		if diags.HasErrors() {
			return diags
		}
		n.Category = fullCategory
	}
	return nil
}

func (n *DashboardNode) Diff(other *DashboardNode) *modconfig.ModTreeItemDiffs {
	res := &modconfig.ModTreeItemDiffs{
		Item: n,
		Name: n.Name(),
	}

	if (n.Category == nil) != (other.Category == nil) {
		res.AddPropertyDiff("Category")
	}
	if n.Category != nil && !n.Category.Equals(other.Category) {
		res.AddPropertyDiff("Category")
	}

	res.Merge(n.QueryProviderImpl.Diff(&other.QueryProviderImpl))
	res.Merge(dashboardLeafNodeDiff(n, other))
	res.PopulateChildDiffs(n, other)

	return res
}

// GetWidth implements DashboardLeafNode
func (n *DashboardNode) GetWidth() int {
	return 0
}

// GetDisplay implements DashboardLeafNode
func (n *DashboardNode) GetDisplay() string {
	return ""
}

// GetType implements DashboardLeafNode
func (n *DashboardNode) GetType() string {
	return ""
}

// CtyValue implements CtyValueProvider
func (n *DashboardNode) CtyValue() (cty.Value, error) {
	return cty_helpers.GetCtyValue(n)
}

func (n *DashboardNode) SetBaseProperties() {
	if n.Base == nil {
		return
	}
	// copy base into the HclResourceImpl 'base' property so it is accessible to all nested structs
	n.HclResourceImpl.SetBase(n.Base)

	// call into parent nested struct SetBaseProperties
	n.QueryProviderImpl.SetBaseProperties()

	if n.Title == nil {
		n.Title = n.Base.Title
	}

	if n.SQL == nil {
		n.SQL = n.Base.SQL
	}

	if n.Query == nil {
		n.Query = n.Base.Query
	}

	if n.Args == nil {
		n.Args = n.Base.Args
	}

	if n.Category == nil {
		n.Category = n.Base.Category
	}

	if n.Params == nil {
		n.Params = n.Base.Params
	}
}
