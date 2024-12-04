package resources

import (
	"fmt"
	"github.com/turbot/pipe-fittings/modconfig"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/go-kit/types"
	typehelpers "github.com/turbot/go-kit/types"
	"github.com/turbot/pipe-fittings/cty_helpers"
	"github.com/turbot/pipe-fittings/printers"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/zclconf/go-cty/cty"
)

// Control is a struct representing the Control resource
type Control struct {
	modconfig.ResourceWithMetadataImpl
	QueryProviderImpl
	DashboardLeafNodeImpl

	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`

	Severity *string `cty:"severity" hcl:"severity"  snapshot:"severity" json:"severity,omitempty"`

	// dashboard specific properties
	Base *Control `hcl:"base" json:"-"`

	parents []modconfig.ModTreeItem
}

func NewControl(block *hcl.Block, mod *modconfig.Mod, shortName string) modconfig.HclResource {
	control := &Control{
		QueryProviderImpl: NewQueryProviderImpl(block, mod, shortName),
	}
	control.Args = NewQueryArgs()
	control.SetAnonymous(block)
	return control
}

func (c *Control) Equals(other *Control) bool {
	res := c.ShortName == other.ShortName &&
		c.FullName == other.FullName &&
		typehelpers.SafeString(c.Description) == typehelpers.SafeString(other.Description) &&
		typehelpers.SafeString(c.Documentation) == typehelpers.SafeString(other.Documentation) &&
		typehelpers.SafeString(c.Severity) == typehelpers.SafeString(other.Severity) &&
		typehelpers.SafeString(c.SQL) == typehelpers.SafeString(other.SQL) &&
		typehelpers.SafeString(c.Title) == typehelpers.SafeString(other.Title)
	if !res {
		return res
	}
	if len(c.Tags) != len(other.Tags) {
		return false
	}
	for k, v := range c.Tags {
		if otherVal := other.Tags[k]; v != otherVal {
			return false
		}
	}

	// args
	if c.Args == nil {
		if other.Args != nil {
			return false
		}
	} else {
		// we have args
		if other.Args == nil {
			return false
		}
		if !c.Args.Equals(other.Args) {
			return false
		}
	}

	// query
	if c.Query == nil {
		if other.Query != nil {
			return false
		}
	} else {
		// we have a query
		if other.Query == nil {
			return false
		}
		if !c.Query.Equals(other.Query) {
			return false
		}
	}

	// params
	if len(c.Params) != len(other.Params) {
		return false
	}
	for i, p := range c.Params {
		if !p.Equals(other.Params[i]) {
			return false
		}
	}

	return true
}

func (c *Control) String() string {
	// build list of parents's names
	parents := c.GetParentNames()
	res := fmt.Sprintf(`
  -----
  Name: %s
  Title: %s
  Description: %s
  SQL: %s
  Parents: %s
`,
		c.FullName,
		types.SafeString(c.Title),
		types.SafeString(c.Description),
		types.SafeString(c.SQL),
		strings.Join(parents, "\n    "))

	// add param defs if there are any
	if len(c.Params) > 0 {
		var paramDefsStr = make([]string, len(c.Params))
		for i, def := range c.Params {
			paramDefsStr[i] = def.String()
		}
		res += fmt.Sprintf("Params:\n\t%s\n  ", strings.Join(paramDefsStr, "\n\t"))
	}

	// add args
	if c.Args != nil && !c.Args.Empty() {
		res += fmt.Sprintf("Args:\n\t%s\n  ", c.Args)
	}
	return res
}

func (c *Control) GetParentNames() []string {
	var parents []string
	for _, p := range c.parents {
		parents = append(parents, p.Name())
	}
	return parents
}

// OnDecoded implements HclResource
func (c *Control) OnDecoded(block *hcl.Block, resourceMapProvider modconfig.ModResourcesProvider) hcl.Diagnostics {
	c.SetBaseProperties()

	return c.QueryProviderImpl.OnDecoded(block, resourceMapProvider)
}

func (c *Control) Diff(other *Control) *modconfig.ModTreeItemDiffs {
	res := &modconfig.ModTreeItemDiffs{
		Item: c,
		Name: c.Name(),
	}

	if !utils.SafeStringsEqual(c.Description, other.Description) {
		res.AddPropertyDiff("Description")
	}
	if !utils.SafeStringsEqual(c.Documentation, other.Documentation) {
		res.AddPropertyDiff("Documentation")
	}
	if !utils.SafeStringsEqual(c.Severity, other.Severity) {
		res.AddPropertyDiff("Severity")
	}
	if len(c.Tags) != len(other.Tags) {
		res.AddPropertyDiff("Tags")
	} else {
		for k, v := range c.Tags {
			if otherVal := other.Tags[k]; v != otherVal {
				res.AddPropertyDiff("Tags")
			}
		}
	}

	res.Merge(dashboardLeafNodeDiff(c, other))
	res.Merge(c.QueryProviderImpl.Diff(other))

	return res
}

// CtyValue implements CtyValueProvider
func (c *Control) CtyValue() (cty.Value, error) {
	return cty_helpers.GetCtyValue(c)
}

func (c *Control) SetBaseProperties() {
	if c.Base == nil {
		return
	}
	// copy base into the HclResourceImpl 'base' property so it is accessible to all nested structs
	c.HclResourceImpl.SetBase(c.Base)
	// call into parent nested struct SetBaseProperties
	c.QueryProviderImpl.SetBaseProperties()

	if c.Severity == nil {
		c.Severity = c.Base.Severity
	}

	if c.Width == nil {
		c.Width = c.Base.Width
	}
	if c.Type == nil {
		c.Type = c.Base.Type
	}
	if c.Display == nil {
		c.Display = c.Base.Display
	}
}

// GetShowData implements printers.Showable
func (c *Control) GetShowData() *printers.RowData {
	res := printers.NewRowData(
		printers.NewFieldValue("Severity", c.Severity),
		printers.NewFieldValue("Width", c.Width),
		printers.NewFieldValue("Type", c.Type),
		printers.NewFieldValue("Display", c.Display),
	)
	// merge fields from base, putting base fields first
	res.Merge(c.QueryProviderImpl.GetShowData())
	return res
}
