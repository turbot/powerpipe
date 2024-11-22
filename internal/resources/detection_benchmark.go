package resources

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/turbot/pipe-fittings/cty_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/printers"
	"github.com/turbot/pipe-fittings/schema"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/zclconf/go-cty/cty"
)

type Benchmark struct {
	modconfig.ResourceWithMetadataImpl
	modconfig.ModTreeItemImpl
	DashboardLeafNodeImpl

	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`

	// store children in a way which can be serialised via cty
	ChildNames       modconfig.NamedItemList `cty:"child_names" json:"-"`
	ChildNameStrings []string                `cty:"child_name_strings" json:"children,omitempty"`

	// dashboard specific properties
	Inputs []*DashboardInput `cty:"inputs" json:"inputs,omitempty"`
	Base   *Benchmark        `hcl:"base" json:"-"`
}

// NewWrapperDetectionBenchmark creates a new Benchmark to wrap a detection which we wish to execute
func NewWrapperDetectionBenchmark(detection *Detection) *Benchmark {
	// create a fake block for the wrapper benchmark
	block := &hcl.Block{
		Type:   schema.BlockTypeDetectionBenchmark,
		// TODO KAI WHICH???
		//Type:   schema.BlockTypeBenchmark,
		Labels: []string{detection.ShortName + "_benchmark"},
		Body:   &hclsyntax.Body{SrcRange: detection.DeclRange},
	}
	b := NewDetectionBenchmark(block, detection.Mod, detection.ShortName).(*Benchmark)
	b.AddChild(detection)
	b.ChildNames = append(b.ChildNames, modconfig.NamedItem{Name: detection.UnqualifiedName})
	b.ChildNameStrings = append(b.ChildNameStrings, detection.UnqualifiedName)
	return b
}

func NewDetectionBenchmark(block *hcl.Block, mod *modconfig.Mod, shortName string) modconfig.HclResource {
	c := &Benchmark{
		ModTreeItemImpl: modconfig.NewModTreeItemImpl(block, mod, shortName),
	}
	c.SetAnonymous(block)

	return c
}

func (d *Benchmark) Equals(other *Benchmark) bool {
	if other == nil {
		return false
	}

	diff := d.Diff(other)
	return !diff.HasChanges()
}

// OnDecoded implements HclResource
func (d *Benchmark) OnDecoded(block *hcl.Block, _ modconfig.ModResourcesProvider) hcl.Diagnostics {
	d.SetBaseProperties()
	return nil
}

func (d *Benchmark) Diff(other *Benchmark) *modconfig.ModTreeItemDiffs {
	res := &modconfig.ModTreeItemDiffs{
		Item: d,
		Name: d.Name(),
	}

	if !utils.SafeStringsEqual(d.FullName, other.FullName) {
		res.AddPropertyDiff("Name")
	}

	if !utils.SafeStringsEqual(d.Title, other.Title) {
		res.AddPropertyDiff("Title")
	}

	if !utils.SafeIntEqual(d.Width, other.Width) {
		res.AddPropertyDiff("Width")
	}

	if !utils.SafeStringsEqual(d.Display, other.Display) {
		res.AddPropertyDiff("Display")
	}

	res.PopulateChildDiffs(d, other)
	return res
}

func (d *Benchmark) WalkResources(resourceFunc func(resource modconfig.HclResource) (bool, error)) error {
	for _, child := range d.Children {
		continueWalking, err := resourceFunc(child.(modconfig.HclResource))
		if err != nil {
			return err
		}
		if !continueWalking {
			break
		}

		if childContainer, ok := child.(*Benchmark); ok {
			if err := childContainer.WalkResources(resourceFunc); err != nil {
				return err
			}
		}
	}
	return nil
}

// CtyValue implements CtyValueProvider
func (d *Benchmark) CtyValue() (cty.Value, error) {
	return cty_helpers.GetCtyValue(d)
}

func (d *Benchmark) SetBaseProperties() {
	if d.Base == nil {
		return
	}
	// copy base into the HclResourceImpl 'base' property so it is accessible to all nested structs
	d.HclResourceImpl.SetBase(d.Base)
	// call into parent nested struct SetBaseProperties
	d.ModTreeItemImpl.SetBaseProperties()

	if d.Width == nil {
		d.Width = d.Base.Width
	}

	if d.Display == nil {
		d.Display = d.Base.Display
	}

	if len(d.GetChildren()) == 0 {
		d.Children = d.Base.Children
		d.ChildNameStrings = d.Base.ChildNameStrings
		d.ChildNames = d.Base.ChildNames
	}
}

// GetShowData implements printers.Showable
func (d *Benchmark) GetShowData() *printers.RowData {
	res := printers.NewRowData(
		printers.NewFieldValue("Width", d.Width),
		printers.NewFieldValue("Display", d.Display),
		printers.NewFieldValue("Inputs", d.Inputs),
		printers.NewFieldValue("Children", d.ChildNameStrings),
	)
	// merge fields from base, putting base fields first
	res.Merge(d.ModTreeItemImpl.GetShowData())
	return res
}

// GetListData implements printers.Listable
func (d *DetectionBenchmark) GetListData() *printers.RowData {
	res := d.ModTreeItemImpl.GetListData()
	// Add type
	res.AddField(printers.NewFieldValue("TYPE", "detection"))

	return res
}
