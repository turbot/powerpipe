package resources

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/pipe-fittings/cty_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/printers"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/zclconf/go-cty/cty"
)

// DashboardImage is a struct representing a leaf dashboard node
type DashboardImage struct {
	modconfig.ResourceWithMetadataImpl
	QueryProviderImpl
	DashboardLeafNodeImpl

	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`

	Src *string `cty:"src" hcl:"src"  json:"src,omitempty" snapshot:"src"`
	Alt *string `cty:"alt" hcl:"alt"  json:"alt,omitempty" snapshot:"alt"`

	Base *DashboardImage `hcl:"base" json:"-"`
}

func NewDashboardImage(block *hcl.Block, mod *modconfig.Mod, shortName string) modconfig.HclResource {
	i := &DashboardImage{
		QueryProviderImpl: NewQueryProviderImpl(block, mod, shortName),
	}
	i.SetAnonymous(block)
	return i
}

func (i *DashboardImage) Equals(other *DashboardImage) bool {
	diff := i.Diff(other)
	return !diff.HasChanges()
}

// OnDecoded implements HclResource
func (i *DashboardImage) OnDecoded(block *hcl.Block, resourceMapProvider modconfig.ModResourcesProvider) hcl.Diagnostics {
	i.SetBaseProperties()
	return i.QueryProviderImpl.OnDecoded(block, resourceMapProvider)
}

func (i *DashboardImage) Diff(other *DashboardImage) *modconfig.ModTreeItemDiffs {
	res := &modconfig.ModTreeItemDiffs{
		Item: i,
		Name: i.Name(),
	}
	if !utils.SafeStringsEqual(i.Src, other.Src) {
		res.AddPropertyDiff("Src")
	}

	if !utils.SafeStringsEqual(i.Alt, other.Alt) {
		res.AddPropertyDiff("Alt")
	}

	res.PopulateChildDiffs(i, other)
	res.Merge(i.QueryProviderImpl.Diff(other))
	res.Merge(dashboardLeafNodeDiff(i, other))

	return res
}

// GetDocumentation implements ModTreeItem
func (*DashboardImage) GetDocumentation() string {
	return ""
}

// GetType implements DashboardLeafNode
func (*DashboardImage) GetType() string {
	return ""
}

// ValidateQuery implements QueryProvider
func (i *DashboardImage) ValidateQuery() hcl.Diagnostics {
	// query is optional - nothing to do
	return nil
}

// CtyValue implements CtyValueProvider
func (i *DashboardImage) CtyValue() (cty.Value, error) {
	return cty_helpers.GetCtyValue(i)
}

func (i *DashboardImage) SetBaseProperties() {
	if i.Base == nil {
		return
	}
	// copy base into the HclResourceImpl 'base' property so it is accessible to all nested structs
	i.HclResourceImpl.SetBase(i.Base)

	// call into parent nested struct SetBaseProperties
	i.QueryProviderImpl.SetBaseProperties()

	if i.Src == nil {
		i.Src = i.Base.Src
	}

	if i.Alt == nil {
		i.Alt = i.Base.Alt
	}

	if i.Width == nil {
		i.Width = i.Base.Width
	}

	if i.Display == nil {
		i.Display = i.Base.Display
	}
}

// GetShowData implements printers.Showable
func (i *DashboardImage) GetShowData() *printers.RowData {
	res := printers.NewRowData(
		printers.NewFieldValue("Width", i.Width),
		printers.NewFieldValue("Display", i.Display),
		printers.NewFieldValue("Src", i.Src),
		printers.NewFieldValue("Alt", i.Alt),
	)
	// merge fields from base, putting base fields first
	res.Merge(i.QueryProviderImpl.GetShowData())
	return res
}
