package resources

import (
	"github.com/hashicorp/hcl/v2"
	typehelpers "github.com/turbot/go-kit/types"
)

type DashboardLeafNodeImpl struct {
	DashboardLeafNodeRemain hcl.Body `hcl:",remain" json:"-"`
	Width                   *int     `cty:"width" hcl:"width"  json:"width,omitempty"`
	Type                    *string  `cty:"type" hcl:"type"  json:"type,omitempty"`
	Display                 *string  `cty:"display" hcl:"display" json:"display,omitempty" snapshot:"display"`
}

func (d *DashboardLeafNodeImpl) GetDisplay() string {
	return typehelpers.SafeString(d.Display)
}

func (d *DashboardLeafNodeImpl) GetType() string {
	return typehelpers.SafeString(d.Type)
}

func (d *DashboardLeafNodeImpl) GetWidth() int {
	return typehelpers.IntValue(d.Width)
}
