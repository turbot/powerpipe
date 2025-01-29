package resources

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/pipe-fittings/v2/cty_helpers"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/pipe-fittings/v2/utils"
	"github.com/zclconf/go-cty/cty"
)

type DashboardCategory struct {
	modconfig.ResourceWithMetadataImpl
	modconfig.ModTreeItemImpl

	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`

	// TACTICAL: include a title property (with a different name to the property in HclResourceImpl  for clarity)
	// This is purely to ensure the title is included in the panel properties of snapshots
	// Note: this will be parsed from HCL, but we must set this explicitly in SetBaseProperties if there is a base
	CategoryTitle *string                               `cty:"title" hcl:"title" snapshot:"title" json:"-"`
	CategoryName  string                                `snapshot:"name" json:"-"`
	Color         *string                               `cty:"color" hcl:"color" snapshot:"color" json:"color,omitempty"`
	Depth         *int                                  `cty:"depth" hcl:"depth" snapshot:"depth" json:"depth,omitempty"`
	Icon          *string                               `cty:"icon" hcl:"icon" snapshot:"icon" json:"icon,omitempty"`
	HREF          *string                               `cty:"href" hcl:"href" snapshot:"href" json:"href,omitempty"`
	Fold          *DashboardCategoryFold                `cty:"fold" hcl:"fold,block" snapshot:"fold" json:"fold,omitempty"`
	PropertyList  DashboardCategoryPropertyList         `cty:"property_list" hcl:"property,block" json:"-"`
	Properties    map[string]*DashboardCategoryProperty `cty:"properties" snapshot:"properties" json:"properties,omitempty"`
	PropertyOrder []string                              `cty:"property_order" hcl:"property_order,optional" snapshot:"property_order" json:"property_order,omitempty"`
	Base          *DashboardCategory                    `hcl:"base" json:"base,omitempty"`
}

func NewDashboardCategory(block *hcl.Block, mod *modconfig.Mod, shortName string) modconfig.HclResource {
	c := &DashboardCategory{
		ModTreeItemImpl: modconfig.NewModTreeItemImpl(block, mod, shortName),
	}
	c.SetAnonymous(block)
	return c
}

// OnDecoded implements HclResource
func (c *DashboardCategory) OnDecoded(block *hcl.Block, _ modconfig.ModResourcesProvider) hcl.Diagnostics {
	c.SetBaseProperties()
	// populate properties map
	if len(c.PropertyList) > 0 {
		c.Properties = make(map[string]*DashboardCategoryProperty, len(c.PropertyList))
		for _, p := range c.PropertyList {
			c.Properties[p.ShortName] = p
		}
	}
	c.CategoryName = c.ResourceMetadata.ResourceName
	return nil
}

func (c *DashboardCategory) Equals(other *DashboardCategory) bool {
	if other == nil {
		return false
	}
	return !c.Diff(other).HasChanges()
}

func (c *DashboardCategory) SetBaseProperties() {
	if c.Base == nil {
		return
	}
	// copy base into the HclResourceImpl 'base' property so it is accessible to all nested structs
	c.HclResourceImpl.SetBase(c.Base)
	// call into parent nested struct SetBaseProperties
	c.ModTreeItemImpl.SetBaseProperties()

	// TACTICAL: DashboardCategory overrides the title property to ensure is included in the snapshot
	c.CategoryTitle = c.Title
	c.CategoryName = c.Name()

	if c.Color == nil {
		c.Color = c.Base.Color
	}
	if c.Depth == nil {
		c.Depth = c.Base.Depth
	}
	if c.Icon == nil {
		c.Icon = c.Base.Icon
	}
	if c.HREF == nil {
		c.HREF = c.Base.HREF
	}
	if c.Fold == nil {
		c.Fold = c.Base.Fold
	}

	if c.PropertyList == nil {
		c.PropertyList = c.Base.PropertyList
	} else {
		c.PropertyList.Merge(c.Base.PropertyList)
	}

	if c.PropertyOrder == nil {
		c.PropertyOrder = c.Base.PropertyOrder
	}
}

func (c *DashboardCategory) Diff(other *DashboardCategory) *modconfig.ModTreeItemDiffs {
	res := &modconfig.ModTreeItemDiffs{
		Item: c,
		Name: c.Name(),
	}

	if (c.Fold == nil) != (other.Fold == nil) {
		res.AddPropertyDiff("Fold")
	}
	if c.Fold != nil && !c.Fold.Equals(other.Fold) {
		res.AddPropertyDiff("Fold")
	}

	if len(c.PropertyList) != len(other.PropertyList) {
		res.AddPropertyDiff("Properties")
	} else {
		for i, p := range c.Properties {
			if !p.Equals(other.Properties[i]) {
				res.AddPropertyDiff("Properties")
			}
		}
	}

	if len(c.PropertyOrder) != len(other.PropertyOrder) {
		res.AddPropertyDiff("PropertyOrder")
	} else {
		for i, p := range c.PropertyOrder {
			if p != other.PropertyOrder[i] {
				res.AddPropertyDiff("PropertyOrder")
			}
		}
	}

	if !utils.SafeStringsEqual(c.Name, other.Name) {
		res.AddPropertyDiff("Name")
	}
	if !utils.SafeStringsEqual(c.Title, other.Title) {
		res.AddPropertyDiff("Title")
	}
	if !utils.SafeStringsEqual(c.Color, other.Color) {
		res.AddPropertyDiff("Color")
	}
	if !utils.SafeStringsEqual(c.Depth, other.Depth) {
		res.AddPropertyDiff("Depth")
	}
	if !utils.SafeStringsEqual(c.Icon, other.Icon) {
		res.AddPropertyDiff("Icon")
	}
	if !utils.SafeStringsEqual(c.HREF, other.HREF) {
		res.AddPropertyDiff("HREF")
	}

	return res
}

// CtyValue implements CtyValueProvider
func (c *DashboardCategory) CtyValue() (cty.Value, error) {
	return cty_helpers.GetCtyValue(c)
}
