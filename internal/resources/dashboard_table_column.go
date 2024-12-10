package resources

import (
	"github.com/turbot/pipe-fittings/sperr"
	"github.com/turbot/pipe-fittings/utils"
)

type DiffMode string

const (
	DiffModeInclude DiffMode = "include"
	DiffModeExclude DiffMode = "exclude"
	DiffModeKey     DiffMode = "key"
)

// all and none are valid values for display and wrap
var validDisplayAndWrap = map[string]struct{}{
	"all":  {},
	"none": {},
}

// include, exclude and key are valid values for diff_mode
var validDiffMode = map[string]struct{}{
	string(DiffModeInclude): {},
	string(DiffModeExclude): {},
	string(DiffModeKey):     {},
}

type DashboardTableColumn struct {
	Name     string    `hcl:"name,label" json:"name" snapshot:"name"`
	Display  *string   `cty:"display" hcl:"display" json:"display,omitempty" snapshot:"display"`
	Wrap     *string   `cty:"wrap" hcl:"wrap" json:"wrap,omitempty" snapshot:"wrap"`
	HREF     *string   `cty:"href" hcl:"href" json:"href,omitempty" snapshot:"href"`
	DiffMode *DiffMode `cty:"diff_mode" hcl:"diff_mode,optional" json:"diff_mode,omitempty" snapshot:"diff_mode"`
}

// Validate checks the validity of the column's properties and sets default values.
func (c *DashboardTableColumn) Validate() error {
	// validate display
	if c.Display != nil {
		if _, ok := validDisplayAndWrap[*c.Display]; !ok {
			return sperr.New(`invalid value for display: %s (allowed values: 'all', 'none')`, *c.Display)
		}
	}

	// validate wrap
	if c.Wrap != nil {
		if _, ok := validDisplayAndWrap[*c.Wrap]; !ok {
			return sperr.New(`invalid value for wrap: %s (allowed values: 'all', 'none')`, *c.Wrap)
		}
	}

	// Set default DiffMode if not set
	if c.DiffMode == nil {
		defaultDiffMode := DiffModeInclude
		c.DiffMode = &defaultDiffMode
	}

	// validate DiffMode
	if _, ok := validDiffMode[string(*c.DiffMode)]; !ok {
		return sperr.New(`invalid value for diff_mode: %s (allowed values: 'include', 'exclude', 'key')`, *c.DiffMode)
	}

	return nil
}

func (c DashboardTableColumn) Equals(other *DashboardTableColumn) bool {
	if other == nil {
		return false
	}

	return utils.SafeStringsEqual(c.Name, other.Name) &&
		utils.SafeStringsEqual(c.Display, other.Display) &&
		utils.SafeStringsEqual(c.Wrap, other.Wrap) &&
		utils.SafeStringsEqual(c.HREF, other.HREF) &&
		utils.SafeStringsEqual(c.DiffMode, other.DiffMode)
}
