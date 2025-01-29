package resources

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/pipe-fittings/v2/cty_helpers"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/zclconf/go-cty/cty"
)

// DashboardEdge is a struct representing a leaf dashboard node
type DashboardEdge struct {
	modconfig.ResourceWithMetadataImpl
	QueryProviderImpl
	DashboardLeafNodeImpl

	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`

	Category *DashboardCategory `cty:"category" hcl:"category" snapshot:"category" json:"category,omitempty"`
	Base     *DashboardEdge     `hcl:"base" json:"-"`
}

func NewDashboardEdge(block *hcl.Block, mod *modconfig.Mod, shortName string) modconfig.HclResource {
	e := &DashboardEdge{
		QueryProviderImpl: NewQueryProviderImpl(block, mod, shortName),
	}

	e.SetAnonymous(block)
	return e
}

func (e *DashboardEdge) Equals(other *DashboardEdge) bool {
	diff := e.Diff(other)
	return !diff.HasChanges()
}

// OnDecoded implements HclResource
func (e *DashboardEdge) OnDecoded(_ *hcl.Block, resourceMapProvider modconfig.ModResourcesProvider) hcl.Diagnostics {
	e.SetBaseProperties()

	// when we reference resources (i.e. category),
	// not all properties are retrieved as they are no cty serialisable
	// repopulate category from resourceMapProvider
	if e.Category != nil {
		fullCategory, diags := enrichCategory(e.Category, e, resourceMapProvider)
		if diags.HasErrors() {
			return diags
		}
		e.Category = fullCategory
	}
	return nil
}

func (e *DashboardEdge) Diff(other *DashboardEdge) *modconfig.ModTreeItemDiffs {
	res := &modconfig.ModTreeItemDiffs{
		Item: e,
		Name: e.Name(),
	}
	if (e.Category == nil) != (other.Category == nil) {
		res.AddPropertyDiff("Category")
	}

	if e.Category != nil && !e.Category.Equals(other.Category) {
		res.AddPropertyDiff("Category")
	}

	res.PopulateChildDiffs(e, other)
	res.Merge(e.QueryProviderImpl.Diff(other))
	res.Merge(dashboardLeafNodeDiff(e, other))

	return res
}

// GetDocumentation implements DashboardLeafNode
func (e *DashboardEdge) GetDocumentation() string {
	return ""
}

// CtyValue implements CtyValueProvider
func (e *DashboardEdge) CtyValue() (cty.Value, error) {
	return cty_helpers.GetCtyValue(e)
}

func (e *DashboardEdge) SetBaseProperties() {
	if e.Base == nil {
		return
	}
	// copy base into the HclResourceImpl 'base' property so it is accessible to all nested structs
	e.HclResourceImpl.SetBase(e.Base)

	// call into parent nested struct SetBaseProperties
	e.QueryProviderImpl.SetBaseProperties()

	if e.Category == nil {
		e.Category = e.Base.Category
	}
}
