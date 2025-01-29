package resources

import (
	"github.com/turbot/pipe-fittings/v2/utils"
)

type DashboardCategoryProperty struct {
	ShortName string  `hcl:"name,label" snapshot:"name" json:"name"`
	Display   *string `cty:"display" hcl:"display" snapshot:"display" json:"display,omitempty"`
	Wrap      *string `cty:"wrap" hcl:"wrap" snapshot:"wrap" json:"wrap,omitempty"`
	HREF      *string `cty:"href" hcl:"href" snapshot:"href" json:"href,omitempty"`
}

func (c DashboardCategoryProperty) Equals(other *DashboardCategoryProperty) bool {
	if other == nil {
		return false
	}

	return utils.SafeStringsEqual(c.ShortName, other.ShortName) &&
		utils.SafeStringsEqual(c.Display, other.Display) &&
		utils.SafeStringsEqual(c.Wrap, other.Wrap) &&
		utils.SafeStringsEqual(c.HREF, other.HREF)
}
