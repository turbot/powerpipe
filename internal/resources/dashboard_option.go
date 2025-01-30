package resources

import "github.com/turbot/pipe-fittings/v2/utils"

// DashboardInputOption is a struct representing dashboard input option
type DashboardInputOption struct {
	Name  string  `hcl:"name,label" json:"name" snapshot:"name"`
	Label *string `cty:"label" hcl:"label" json:"label,omitempty" snapshot:"label"`
}

func (o DashboardInputOption) Equals(other *DashboardInputOption) bool {
	return utils.SafeStringsEqual(o.Name, other.Name) && utils.SafeStringsEqual(o.Label, other.Label)
}
