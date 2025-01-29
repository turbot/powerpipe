package resources

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/pipe-fittings/v2/cty_helpers"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/pipe-fittings/v2/printers"
	"github.com/turbot/pipe-fittings/v2/utils"
	"github.com/zclconf/go-cty/cty"
)

// DashboardText is a struct representing a leaf dashboard node
type DashboardText struct {
	modconfig.ResourceWithMetadataImpl
	modconfig.ModTreeItemImpl
	DashboardLeafNodeImpl

	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`

	Value *string `cty:"value" hcl:"value" snapshot:"value"  json:"value,omitempty"`

	Base *DashboardText `hcl:"base" json:"-"`
	Mod  *modconfig.Mod `cty:"mod" json:"-"`
}

func NewDashboardText(block *hcl.Block, mod *modconfig.Mod, shortName string) modconfig.HclResource {
	t := &DashboardText{
		ModTreeItemImpl: modconfig.NewModTreeItemImpl(block, mod, shortName),
	}
	t.SetAnonymous(block)
	return t
}

func (t *DashboardText) Equals(other *DashboardText) bool {
	diff := t.Diff(other)
	return !diff.HasChanges()
}

// OnDecoded implements HclResource
func (t *DashboardText) OnDecoded(*hcl.Block, modconfig.ModResourcesProvider) hcl.Diagnostics {
	t.SetBaseProperties()
	return nil
}

func (t *DashboardText) Diff(other *DashboardText) *modconfig.ModTreeItemDiffs {
	res := &modconfig.ModTreeItemDiffs{
		Item: t,
		Name: t.Name(),
	}

	if !utils.SafeStringsEqual(t.Type, other.Type) {
		res.AddPropertyDiff("Type")
	}

	if !utils.SafeStringsEqual(t.Value, other.Value) {
		res.AddPropertyDiff("Value")
	}

	res.PopulateChildDiffs(t, other)
	res.Merge(dashboardLeafNodeDiff(t, other))
	return res
}

// GetDocumentation implements ModTreeItem
func (*DashboardText) GetDocumentation() string {
	return ""
}

// CtyValue implements CtyValueProvider
func (t *DashboardText) CtyValue() (cty.Value, error) {
	return cty_helpers.GetCtyValue(t)
}

func (t *DashboardText) SetBaseProperties() {
	if t.Base == nil {
		return
	}
	if t.Title == nil {
		t.Title = t.Base.Title
	}
	if t.Type == nil {
		t.Type = t.Base.Type
	}
	if t.Display == nil {
		t.Display = t.Base.Display
	}
	if t.Value == nil {
		t.Value = t.Base.Value
	}
	if t.Width == nil {
		t.Width = t.Base.Width
	}
}

// GetShowData implements printers.Showable
func (t *DashboardText) GetShowData() *printers.RowData {
	res := printers.NewRowData(
		printers.NewFieldValue("Width", t.Width),
		printers.NewFieldValue("Type", t.Type),
		printers.NewFieldValue("Display", t.Display),
		printers.NewFieldValue("Value", t.Value),
	)
	// merge fields from base, putting base fields first
	res.Merge(t.ModTreeItemImpl.GetShowData())
	return res
}
