package resources

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/pipe-fittings/cty_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/printers"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/zclconf/go-cty/cty"
)

// Detection is a struct representing a leaf dashboard node
type Detection struct {
	modconfig.ResourceWithMetadataImpl
	QueryProviderImpl
	DashboardLeafNodeImpl

	// required to allow partial decoding
	Remain         hcl.Body `hcl:",remain" json:"-"`
	Severity       *string  `cty:"severity" hcl:"severity"  snapshot:"severity" json:"severity,omitempty"`
	DisplayColumns []string `cty:"display_columns" hcl:"display_columns,optional" json:"display_columns,omitempty" snapshot:"display_columns"`

	Base *Detection `hcl:"base" json:"-"`
}

func NewDetection(block *hcl.Block, mod *modconfig.Mod, shortName string) modconfig.HclResource {
	t := &Detection{
		QueryProviderImpl: NewQueryProviderImpl(block, mod, shortName),
	}
	t.SetAnonymous(block)
	return t
}

func (t *Detection) Equals(other *Detection) bool {
	diff := t.Diff(other)
	return !diff.HasChanges()
}

// OnDecoded implements HclResource
func (t *Detection) OnDecoded(block *hcl.Block, resourceMapProvider modconfig.ModResourcesProvider) hcl.Diagnostics {
	t.SetBaseProperties()
	return t.QueryProviderImpl.OnDecoded(block, resourceMapProvider)
}

func (t *Detection) Diff(other *Detection) *modconfig.ModTreeItemDiffs {
	res := &modconfig.ModTreeItemDiffs{
		Item: t,
		Name: t.Name(),
	}

	if !utils.SafeStringsEqual(t.Type, other.Type) {
		res.AddPropertyDiff("Type")
	}

	res.PopulateChildDiffs(t, other)
	res.Merge(t.QueryProviderImpl.Diff(other))
	res.Merge(dashboardLeafNodeDiff(t, other))

	return res
}

// GetDocumentation implements ModTreeItem
func (*Detection) GetDocumentation() string {
	return ""
}

// CtyValue implements CtyValueProvider
func (t *Detection) CtyValue() (cty.Value, error) {
	return cty_helpers.GetCtyValue(t)
}

func (t *Detection) SetBaseProperties() {
	if t.Base == nil {
		return
	}
	// copy base into the HclResourceImpl 'base' property so it is accessible to all nested structs
	t.HclResourceImpl.SetBase(t.Base)

	// call into parent nested struct SetBaseProperties
	t.QueryProviderImpl.SetBaseProperties()

	if t.Width == nil {
		t.Width = t.Base.Width
	}

	if t.Type == nil {
		t.Type = t.Base.Type
	}

	if t.Display == nil {
		t.Display = t.Base.Display
	}
}

// GetShowData implements printers.Showable
func (t *Detection) GetShowData() *printers.RowData {
	res := printers.NewRowData(
		printers.NewFieldValue("Width", t.Width),
		printers.NewFieldValue("Type", t.Type),
		printers.NewFieldValue("Display", t.Display),
	)
	// merge fields from base, putting base fields first
	res.Merge(t.QueryProviderImpl.GetShowData())
	return res
}
