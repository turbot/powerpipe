package resources

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/stevenle/topsort"
	typehelpers "github.com/turbot/go-kit/types"
	"github.com/turbot/pipe-fittings/cty_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/printers"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/zclconf/go-cty/cty"
)

type DetectionBenchmark struct {
	modconfig.ResourceWithMetadataImpl
	modconfig.ModTreeItemImpl

	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`

	Width   *int              `cty:"width" hcl:"width"  json:"width,omitempty"`
	Display *string           `cty:"display" hcl:"display" json:"display,omitempty"`
	Inputs  []*DashboardInput `cty:"inputs" json:"inputs,omitempty"`
	// store children in a way which can be serialised via cty
	ChildNames []string `cty:"children" json:"children,omitempty"`

	//nolint:unused // TODO: unused attribute
	runtimeDependencyGraph *topsort.Graph
}

func NewDetectionBenchmark(block *hcl.Block, mod *modconfig.Mod, shortName string) modconfig.HclResource {
	c := &DetectionBenchmark{
		ModTreeItemImpl: modconfig.NewModTreeItemImpl(block, mod, shortName),
	}
	c.SetAnonymous(block)

	return c
}

func (c *DetectionBenchmark) Equals(other *DetectionBenchmark) bool {
	diff := c.Diff(other)
	return !diff.HasChanges()
}

// OnDecoded implements HclResource
func (c *DetectionBenchmark) OnDecoded(block *hcl.Block, _ modconfig.ModResourcesProvider) hcl.Diagnostics {
	c.ChildNames = make([]string, len(c.GetChildren()))
	for i, child := range c.GetChildren() {
		c.ChildNames[i] = child.Name()
	}
	return nil
}

// GetWidth implements DashboardLeafNode
func (c *DetectionBenchmark) GetWidth() int {
	if c.Width == nil {
		return 0
	}
	return *c.Width
}

// GetDisplay implements DashboardLeafNode
func (c *DetectionBenchmark) GetDisplay() string {
	return typehelpers.SafeString(c.Display)
}

// GetType implements DashboardLeafNode
func (c *DetectionBenchmark) GetType() string {
	return ""
}

func (c *DetectionBenchmark) Diff(other *DetectionBenchmark) *modconfig.ModTreeItemDiffs {
	res := &modconfig.ModTreeItemDiffs{
		Item: c,
		Name: c.Name(),
	}

	if !utils.SafeStringsEqual(c.FullName, other.FullName) {
		res.AddPropertyDiff("Name")
	}

	if !utils.SafeStringsEqual(c.Title, other.Title) {
		res.AddPropertyDiff("Title")
	}

	if !utils.SafeIntEqual(c.Width, other.Width) {
		res.AddPropertyDiff("Width")
	}

	if !utils.SafeStringsEqual(c.Display, other.Display) {
		res.AddPropertyDiff("Display")
	}

	res.PopulateChildDiffs(c, other)
	return res
}

func (c *DetectionBenchmark) WalkResources(resourceFunc func(resource modconfig.HclResource) (bool, error)) error {
	for _, child := range c.Children {
		continueWalking, err := resourceFunc(child.(modconfig.HclResource))
		if err != nil {
			return err
		}
		if !continueWalking {
			break
		}

		if childContainer, ok := child.(*DetectionBenchmark); ok {
			if err := childContainer.WalkResources(resourceFunc); err != nil {
				return err
			}
		}
	}
	return nil
}

// CtyValue implements CtyValueProvider
func (c *DetectionBenchmark) CtyValue() (cty.Value, error) {
	return cty_helpers.GetCtyValue(c)
}

// GetShowData implements printers.Showable
func (c *DetectionBenchmark) GetShowData() *printers.RowData {
	res := printers.NewRowData(
		printers.NewFieldValue("Width", c.Width),
		printers.NewFieldValue("Display", c.Display),
		printers.NewFieldValue("Inputs", c.Inputs),
		printers.NewFieldValue("Children", c.ChildNames),
	)
	// merge fields from base, putting base fields first
	res.Merge(c.ModTreeItemImpl.GetShowData())
	return res
}
