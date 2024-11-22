package resources

import (
	"fmt"
	"github.com/turbot/pipe-fittings/modconfig"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/pipe-fittings/cty_helpers"
	"github.com/turbot/pipe-fittings/printers"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/zclconf/go-cty/cty"
)

// Benchmark is a struct representing the Benchmark resource
type ControlBenchmark struct {
	modconfig.ResourceWithMetadataImpl
	modconfig.ModTreeItemImpl
	DashboardLeafNodeImpl

	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`

	// child names as NamedItem structs - used to allow setting children via the 'children' property
	ChildNames       modconfig.NamedItemList `cty:"child_names" json:"-"`
	ChildNameStrings []string                `cty:"child_name_strings" json:"children,omitempty"`

	// dashboard specific properties
	Inputs []*DashboardInput `cty:"inputs" json:"inputs,omitempty"`

	Base *ControlBenchmark `hcl:"base" json:"-"`
}

func NewRootBenchmarkWithChildren(mod *modconfig.Mod, children []modconfig.ModTreeItem) modconfig.HclResource {
	fullName := fmt.Sprintf("%s.%s.%s", mod.ShortName, "benchmark", "root")
	benchmark := &ControlBenchmark{
		ModTreeItemImpl: modconfig.ModTreeItemImpl{
			HclResourceImpl: modconfig.HclResourceImpl{
				ShortName:       "root",
				FullName:        fullName,
				UnqualifiedName: fmt.Sprintf("%s.%s", "benchmark", "root"),
				BlockType:       "benchmark",
			},
			Mod: mod,
		},
	}

	benchmark.AddChild(children...)
	return benchmark
}

func NewBenchmark(block *hcl.Block, mod *modconfig.Mod, shortName string) modconfig.HclResource {
	benchmark := &ControlBenchmark{
		ModTreeItemImpl: modconfig.NewModTreeItemImpl(block, mod, shortName),
	}
	benchmark.SetAnonymous(block)
	return benchmark
}

func (b *ControlBenchmark) Equals(other *ControlBenchmark) bool {
	if other == nil {
		return false
	}

	return !b.Diff(other).HasChanges()
}

// OnDecoded implements HclResource
func (b *ControlBenchmark) OnDecoded(block *hcl.Block, _ modconfig.ModResourcesProvider) hcl.Diagnostics {
	b.SetBaseProperties()

	return nil
}

func (b *ControlBenchmark) String() string {
	// build list of children's names
	var children []string
	for _, child := range b.GetChildren() {
		children = append(children, child.Name())
	}
	// build list of parents names
	var parents []string
	for _, p := range b.GetParents() {
		parents = append(parents, p.Name())
	}
	sort.Strings(children)
	return fmt.Sprintf(`
	 -----
	 Name: %s
	 Title: %s
	 Description: %s
	 Parent: %s
	 Children:
	   %s
	`,
		b.FullName,
		types.SafeString(b.Title),
		types.SafeString(b.Description),
		strings.Join(parents, "\n    "),
		strings.Join(children, "\n    "))
}

// GetChildControls return a flat list of controls underneath the benchmark in the tree
func (b *ControlBenchmark) GetChildControls() []*Control {
	var res []*Control
	for _, child := range b.GetChildren() {
		if control, ok := child.(*Control); ok {
			res = append(res, control)
		} else if benchmark, ok := child.(*ControlBenchmark); ok {
			res = append(res, benchmark.GetChildControls()...)
		}
	}
	return res
}

func (b *ControlBenchmark) Diff(other *ControlBenchmark) *modconfig.ModTreeItemDiffs {
	res := &modconfig.ModTreeItemDiffs{
		Item: b,
		Name: b.Name(),
	}

	if !utils.SafeStringsEqual(b.Description, other.Description) {
		res.AddPropertyDiff("Description")
	}
	if !utils.SafeStringsEqual(b.Documentation, other.Documentation) {
		res.AddPropertyDiff("Documentation")
	}
	if !utils.SafeStringsEqual(b.Title, other.Title) {
		res.AddPropertyDiff("Title")
	}
	if len(b.Tags) != len(other.Tags) {
		res.AddPropertyDiff("Tags")
	} else {
		for k, v := range b.Tags {
			if otherVal := other.Tags[k]; v != otherVal {
				res.AddPropertyDiff("Tags")
			}
		}
	}

	if b.Type != other.Type {
		res.AddPropertyDiff("Type")
	}

	if len(b.ChildNameStrings) != len(other.ChildNameStrings) {
		res.AddPropertyDiff("Childen")
	} else {
		myChildNames := b.ChildNameStrings
		sort.Strings(myChildNames)
		otherChildNames := other.ChildNameStrings
		sort.Strings(otherChildNames)
		if strings.Join(myChildNames, ",") != strings.Join(otherChildNames, ",") {
			res.AddPropertyDiff("Childen")
		}
	}

	res.Merge(dashboardLeafNodeDiff(b, other))
	return res
}

func (b *ControlBenchmark) WalkResources(resourceFunc func(resource modconfig.ModTreeItem) (bool, error)) error {
	for _, child := range b.GetChildren() {
		continueWalking, err := resourceFunc(child)
		if err != nil {
			return err
		}
		if !continueWalking {
			break
		}

		if childContainer, ok := child.(*ControlBenchmark); ok {
			if err := childContainer.WalkResources(resourceFunc); err != nil {
				return err
			}
		}
	}
	return nil
}

// CtyValue implements CtyValueProvider
func (b *ControlBenchmark) CtyValue() (cty.Value, error) {
	return cty_helpers.GetCtyValue(b)
}

func (b *ControlBenchmark) SetBaseProperties() {
	if b.Base == nil {
		return
	}
	// copy base into the HclResourceImpl 'base' property so it is accessible to all nested structs
	b.HclResourceImpl.SetBase(b.Base)
	// call into parent nested struct SetBaseProperties
	b.ModTreeItemImpl.SetBaseProperties()

	if b.Width == nil {
		b.Width = b.Base.Width
	}

	if b.Display == nil {
		b.Display = b.Base.Display
	}

	if len(b.GetChildren()) == 0 {
		b.Children = b.Base.Children
		b.ChildNameStrings = b.Base.ChildNameStrings
		b.ChildNames = b.Base.ChildNames
	}
}

// GetShowData implements printers.Showable
func (b *ControlBenchmark) GetShowData() *printers.RowData {
	res := printers.NewRowData(
		printers.FieldValue{Name: "Children", Value: b.ChildNameStrings},
	)
	res.Merge(b.ModTreeItemImpl.GetShowData())
	return res
}

// GetListData implements printers.Listable
func (b *Benchmark) GetListData() *printers.RowData {
	res := b.ModTreeItemImpl.GetListData()
	// add type
	res.AddField(printers.NewFieldValue("TYPE", "control"))

	return res
}
