package resources

import (
	"github.com/hashicorp/hcl/v2"
	typehelpers "github.com/turbot/go-kit/types"
	"github.com/turbot/pipe-fittings/cty_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/printers"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/zclconf/go-cty/cty"
)

// DashboardCard is a struct representing a leaf dashboard node
type DashboardCard struct {
	modconfig.ResourceWithMetadataImpl
	QueryProviderImpl

	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`

	Label *string `cty:"label" hcl:"label" snapshot:"label" json:"label,omitempty"`
	Value *string `cty:"value" hcl:"value" snapshot:"value" json:"value,omitempty"`
	Icon  *string `cty:"icon" hcl:"icon" snapshot:"icon" json:"icon,omitempty"`
	HREF  *string `cty:"href" hcl:"href" snapshot:"href" json:"href,omitempty"`

	Width   *int           `cty:"width" hcl:"width"  json:"width,omitempty"`
	Type    *string        `cty:"type" hcl:"type"  json:"type,omitempty"`
	Display *string        `cty:"display" hcl:"display" json:"display,omitempty"`
	Base    *DashboardCard `hcl:"base" json:"-"`
}

func NewDashboardCard(block *hcl.Block, mod *modconfig.Mod, shortName string) modconfig.HclResource {
	c := &DashboardCard{
		QueryProviderImpl: NewQueryProviderImpl(block, mod, shortName),
	}

	c.SetAnonymous(block)
	return c
}

func (c *DashboardCard) Equals(other *DashboardCard) bool {
	diff := c.Diff(other)
	return !diff.HasChanges()
}

// OnDecoded implements HclResource
func (c *DashboardCard) OnDecoded(block *hcl.Block, resourceMapProvider modconfig.ResourceMapsProvider) hcl.Diagnostics {
	c.SetBaseProperties()
	return c.QueryProviderImpl.OnDecoded(block, resourceMapProvider)
}

func (c *DashboardCard) Diff(other *DashboardCard) *modconfig.ModTreeItemDiffs {
	res := &modconfig.ModTreeItemDiffs{
		Item: c,
		Name: c.Name(),
	}

	if !utils.SafeStringsEqual(c.Label, other.Label) {
		res.AddPropertyDiff("Instance")
	}

	if !utils.SafeStringsEqual(c.Value, other.Value) {
		res.AddPropertyDiff("Value")
	}

	if !utils.SafeStringsEqual(c.Type, other.Type) {
		res.AddPropertyDiff("Type")
	}

	if !utils.SafeStringsEqual(c.Icon, other.Icon) {
		res.AddPropertyDiff("Icon")
	}

	if !utils.SafeStringsEqual(c.HREF, other.HREF) {
		res.AddPropertyDiff("HREF")
	}

	res.Merge(c.QueryProviderImpl.Diff(&other.QueryProviderImpl))
	res.PopulateChildDiffs(c, other)
	res.Merge(dashboardLeafNodeDiff(c, other))

	return res
}

// GetWidth implements DashboardLeafNode
func (c *DashboardCard) GetWidth() int {
	if c.Width == nil {
		return 0
	}
	return *c.Width
}

// GetDisplay implements DashboardLeafNode
func (c *DashboardCard) GetDisplay() string {
	return typehelpers.SafeString(c.Display)
}

// GetDocumentation implements DashboardLeafNode, ModTreeItem
func (c *DashboardCard) GetDocumentation() string {
	return ""
}

// GetType implements DashboardLeafNode
func (c *DashboardCard) GetType() string {
	return typehelpers.SafeString(c.Type)
}

// ValidateQuery implements QueryProvider
func (c *DashboardCard) ValidateQuery() hcl.Diagnostics {
	// query is optional - nothing to do
	return nil
}

// CtyValue implements CtyValueProvider
func (c *DashboardCard) CtyValue() (cty.Value, error) {
	return cty_helpers.GetCtyValue(c)
}

func (c *DashboardCard) SetBaseProperties() {
	if c.Base == nil {
		return
	}
	// copy base into the HclResourceImpl 'base' property so it is accessible to all nested structs
	c.HclResourceImpl.SetBase(c.Base)
	// call into parent nested struct SetBaseProperties
	c.QueryProviderImpl.SetBaseProperties()

	if c.Label == nil {
		c.Label = c.Base.Label
	}

	if c.Value == nil {
		c.Value = c.Base.Value
	}

	if c.Type == nil {
		c.Type = c.Base.Type
	}

	if c.Display == nil {
		c.Display = c.Base.Display
	}

	if c.Icon == nil {
		c.Icon = c.Base.Icon
	}

	if c.HREF == nil {
		c.HREF = c.Base.HREF
	}

	if c.Width == nil {
		c.Width = c.Base.Width
	}
}

// GetShowData implements printers.Showable
func (c *DashboardCard) GetShowData() *printers.RowData {
	res := printers.NewRowData(
		printers.NewFieldValue("Label", c.Label),
		printers.NewFieldValue("Value", c.Value),
		printers.NewFieldValue("Icon", c.Icon),
		printers.NewFieldValue("HREF", c.HREF),
		printers.NewFieldValue("Width", c.Width),
		printers.NewFieldValue("Type", c.Type),
		printers.NewFieldValue("Display", c.Display),
	)
	// merge fields from base, putting base fields first
	res.Merge(c.QueryProviderImpl.GetShowData())
	return res
}
