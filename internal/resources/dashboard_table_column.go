package resources

import "github.com/turbot/pipe-fittings/utils"

type DashboardTableColumn struct {
	Name    string  `hcl:"name,label" json:"name" snapshot:"name"`
	Display *string `cty:"display" hcl:"display" json:"display,omitempty" snapshot:"display"`
	Wrap    *string `cty:"wrap" hcl:"wrap" json:"wrap,omitempty" snapshot:"wrap"`
	HREF    *string `cty:"href" hcl:"href" json:"href,omitempty" snapshot:"href"`

	Attribute *string `cty:"attribute" hcl:"attribute" json:"attribute,omitempty" snapshot:"attribute"`
}

func (c DashboardTableColumn) Equals(other *DashboardTableColumn) bool {
	if other == nil {
		return false
	}

	return utils.SafeStringsEqual(c.Name, other.Name) &&
		utils.SafeStringsEqual(c.Display, other.Display) &&
		utils.SafeStringsEqual(c.Wrap, other.Wrap) &&
		utils.SafeStringsEqual(c.HREF, other.HREF) &&
		utils.SafeStringsEqual(c.Attribute, other.Attribute)
}
