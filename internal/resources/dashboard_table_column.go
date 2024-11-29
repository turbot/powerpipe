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

type DashboardTableColumn struct {
	Name     string    `hcl:"name,label" json:"name" snapshot:"name"`
	Display  *string   `cty:"display" hcl:"display" json:"display,omitempty" snapshot:"display"`
	Wrap     *string   `cty:"wrap" hcl:"wrap" json:"wrap,omitempty" snapshot:"wrap"`
	HREF     *string   `cty:"href" hcl:"href" json:"href,omitempty" snapshot:"href"`
	DiffMode *DiffMode `cty:"diff_mode" hcl:"diff_mode,optional" json:"diff_mode,omitempty" snapshot:"diff_mode"`
}

// Validate checks the validity of the column's properties and sets default values.
func (c *DashboardTableColumn) Validate() error {
	// validate Display
	if c.Display != nil && *c.Display != "all" && *c.Display != "none" {
		return sperr.New(`invalid value for display: %s (allowed values: 'all', 'none')`, *c.Display)
	}

	// validate Wrap
	if c.Wrap != nil && *c.Wrap != "all" && *c.Wrap != "none" {
		return sperr.New(`invalid value for wrap: %s (allowed values: 'all', 'none')`, *c.Wrap)
	}

	// Set default DiffMode if not set
	if c.DiffMode == nil {
		defaultDiffMode := DiffModeInclude
		c.DiffMode = &defaultDiffMode
	}

	// validate DiffMode
	if c.DiffMode != nil && *c.DiffMode != DiffModeInclude && *c.DiffMode != DiffModeExclude && *c.DiffMode != DiffModeKey {
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
