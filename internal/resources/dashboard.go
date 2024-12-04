package resources

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/spf13/viper"
	"github.com/stevenle/topsort"
	typehelpers "github.com/turbot/go-kit/types"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/cty_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/schema"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/zclconf/go-cty/cty"
)

const rootRuntimeDependencyNode = "rootRuntimeDependencyNode"
const RuntimeDependencyDashboardScope = "self"

// Dashboard is a struct representing the Dashboard  resource
type Dashboard struct {
	modconfig.ResourceWithMetadataImpl
	modconfig.ModTreeItemImpl
	WithProviderImpl
	DashboardLeafNodeImpl

	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`

	Inputs  []*DashboardInput `cty:"inputs" json:"inputs,omitempty"`
	UrlPath string            `cty:"url_path"  json:"url_path,omitempty"`
	Base    *Dashboard        `hcl:"base" json:"-"`
	// store children in a way which can be serialised via cty
	ChildNames []string `cty:"children" json:"children,omitempty"`
	// map of all inputs in our resource tree
	selfInputsMap          map[string]*DashboardInput
	runtimeDependencyGraph *topsort.Graph
}

func NewDashboard(block *hcl.Block, mod *modconfig.Mod, shortName string) modconfig.HclResource {
	d := &Dashboard{
		ModTreeItemImpl: modconfig.NewModTreeItemImpl(block, mod, shortName),
	}
	d.SetAnonymous(block)
	d.setUrlPath()

	return d
}

// NewQueryDashboard creates a dashboard to wrap a query/control
// this is used for snapshot generation
func NewQueryDashboard(qp QueryProvider) (*Dashboard, error) {
	parsedName, title, err := getQueryDashboardName(qp)
	if err != nil {
		return nil, err
	}
	fullName := parsedName.ToFullName()

	// for query dashboard use generated title, for control use original title
	if qp.GetBlockType() != schema.BlockTypeQuery {
		title = qp.GetTitle()
	}

	var dashboard = &Dashboard{
		ModTreeItemImpl: modconfig.ModTreeItemImpl{
			HclResourceImpl: modconfig.HclResourceImpl{
				ShortName:       parsedName.Name,
				FullName:        fullName,
				UnqualifiedName: parsedName.ToResourceName(),
				Title:           utils.ToStringPointer(title),
				Description:     utils.ToStringPointer(qp.GetDescription()),
				Documentation:   utils.ToStringPointer(qp.GetDocumentation()),
				Tags:            qp.GetTags(),
				BlockType:       schema.BlockTypeDashboard,
				DeclRange:       *qp.GetDeclRange(),
			},
			Mod: qp.(modconfig.ModItem).GetMod(),
		},
	}

	dashboard.setUrlPath()

	table, err := NewQueryDashboardTable(qp)
	if err != nil {
		return nil, err
	}
	dashboard.Children = []modconfig.ModTreeItem{table}

	return dashboard, nil
}

func getQueryDashboardName(qp QueryProvider) (*modconfig.ParsedResourceName, string, error) {
	var sql string
	if q := qp.GetQuery(); q != nil {
		sql = typehelpers.SafeString(q.GetSQL())
	} else {
		sql = typehelpers.SafeString(qp.GetSQL())
	}
	hash, err := utils.Base36Hash(sql, 8)
	if err != nil {
		return nil, "", err
	}
	dashboardName := fmt.Sprintf("custom.dashboard.sql_%s", hash)

	parsed, err := modconfig.ParseResourceName(dashboardName)
	if err != nil {
		return nil, "", err
	}
	title := getQueryDashboardTitle(hash)
	return parsed, title, nil
}

func getQueryDashboardTitle(queryHash string) string {
	if titleArg := viper.GetString(constants.ArgSnapshotTitle); titleArg != "" {
		return titleArg
	}
	return fmt.Sprintf("Custom query [%s]", queryHash)
}

func (d *Dashboard) setUrlPath() {
	d.UrlPath = fmt.Sprintf("/%s", d.FullName)
}

func (d *Dashboard) Equals(other *Dashboard) bool {
	diff := d.Diff(other)
	return !diff.HasChanges()
}

// OnDecoded implements HclResource
func (d *Dashboard) OnDecoded(block *hcl.Block, _ modconfig.ModResourcesProvider) hcl.Diagnostics {
	diags := d.SetBaseProperties()
	if diags.HasErrors() {
		return diags
	}
	children := d.GetChildren()
	d.ChildNames = make([]string, len(children))
	for i, child := range children {
		d.ChildNames[i] = child.Name()
	}

	return nil
}

func (d *Dashboard) Diff(other *Dashboard) *modconfig.ModTreeItemDiffs {
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

	if len(d.Tags) != len(other.Tags) {
		res.AddPropertyDiff("Tags")
	} else {
		for k, v := range d.Tags {
			if otherVal := other.Tags[k]; v != otherVal {
				res.AddPropertyDiff("Tags")
			}
		}
	}

	if !utils.SafeStringsEqual(d.Description, other.Description) {
		res.AddPropertyDiff("Description")
	}

	if !utils.SafeStringsEqual(d.Documentation, other.Documentation) {
		res.AddPropertyDiff("Documentation")
	}

	res.PopulateChildDiffs(d, other)
	return res
}

func (d *Dashboard) AddChild(child modconfig.ModTreeItem) hcl.Diagnostics {
	var diags hcl.Diagnostics
	d.ModTreeItemImpl.AddChild(child)

	switch c := child.(type) {
	case *DashboardInput:
		d.Inputs = append(d.Inputs, c)
	case *DashboardWith:
		err := d.AddWith(c)
		if err != nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "failed to add 'with' block  to dashboard",
				Detail:   err.Error(),
				Subject:  &c.DeclRange,
			},
			)
		}
	}
	return diags
}

func (d *Dashboard) WalkResources(resourceFunc func(resource modconfig.HclResource) (bool, error)) error {
	for _, child := range d.GetChildren() {
		continueWalking, err := resourceFunc(child.(modconfig.HclResource))
		if err != nil {
			return err
		}
		if !continueWalking {
			break
		}

		if container, ok := child.(*DashboardContainer); ok {
			if err := container.WalkResources(resourceFunc); err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *Dashboard) ValidateRuntimeDependencies(workspace modconfig.ModResourcesProvider) error {
	d.runtimeDependencyGraph = topsort.NewGraph()
	// add root node - this will depend on all other nodes
	d.runtimeDependencyGraph.AddNode(rootRuntimeDependencyNode)

	// define a walk function which determines whether the resource has runtime dependencies and if so,
	// add to the graph
	resourceFunc := func(resource modconfig.HclResource) (bool, error) {
		wp, ok := resource.(WithProvider)
		if !ok {
			// continue walking
			return true, nil
		}

		if err := d.validateRuntimeDependenciesForResource(resource, workspace); err != nil {
			return false, err
		}

		// if the query provider has any 'with' blocks, add these dependencies as well
		for _, with := range wp.GetWiths() {
			if err := d.validateRuntimeDependenciesForResource(with, workspace); err != nil {
				return false, err
			}
		}

		// continue walking
		return true, nil
	}
	if err := d.WalkResources(resourceFunc); err != nil {
		return err
	}

	// ensure that dependencies can be resolved
	if _, err := d.runtimeDependencyGraph.TopSort(rootRuntimeDependencyNode); err != nil {
		return fmt.Errorf("runtime depedencies cannot be resolved: %s", err.Error())
	}
	return nil
}

func (d *Dashboard) validateRuntimeDependenciesForResource(resource modconfig.HclResource, workspace modconfig.ModResourcesProvider) error {
	// TODO  [node_reuse] re-add parse time validation https://github.com/turbot/steampipe/issues/2925
	return nil
	//rdp := resource.(RuntimeDependencyProvider)
	//// WHAT ABOUT CHILDREN
	//if len(runtimeDependencies) == 0 {
	//	return nil
	//}
	//name := resource.Name()
	//if !d.runtimeDependencyGraph.ContainsNode(name) {
	//	d.runtimeDependencyGraph.AddNode(name)
	//}
	//
	//for _, dependency := range runtimeDependencies {
	//	// try to resolve the dependency source resource
	//	if err := dependency.ValidateSource(d, workspace); err != nil {
	//		return err
	//	}
	//	if err := d.runtimeDependencyGraph.AddEdge(rootRuntimeDependencyNode, name); err != nil {
	//		return err
	//	}
	//	depString := dependency.String()
	//	if !d.runtimeDependencyGraph.ContainsNode(depString) {
	//		d.runtimeDependencyGraph.AddNode(depString)
	//	}
	//	if err := d.runtimeDependencyGraph.AddEdge(name, dependency.String()); err != nil {
	//		return err
	//	}
	//}
	//return nil
}

func (d *Dashboard) GetInput(name string) (*DashboardInput, bool) {
	input, found := d.selfInputsMap[name]
	return input, found
}

func (d *Dashboard) GetInputs() map[string]*DashboardInput {
	return d.selfInputsMap
}

func (d *Dashboard) InitInputs() hcl.Diagnostics {
	// add all our direct child inputs to a map
	// (we must do this before adding child container inputs to detect dupes)
	duplicates := d.setInputMap()

	//  add child containers and dashboard inputs
	resourceFunc := func(resource modconfig.HclResource) (bool, error) {
		if container, ok := resource.(*DashboardContainer); ok {
			for _, i := range container.Inputs {
				// check we do not already have this input
				if _, ok := d.selfInputsMap[i.UnqualifiedName]; ok {
					duplicates = append(duplicates, i.Name())

				}
				d.Inputs = append(d.Inputs, i)
				d.selfInputsMap[i.UnqualifiedName] = i
			}
		}
		// continue walking
		return true, nil
	}
	if err := d.WalkResources(resourceFunc); err != nil {
		return hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("Dashboard '%s' WalkResources failed", d.Name()),
			Detail:   err.Error(),
			Subject:  &d.DeclRange,
		}}
	}

	if len(duplicates) > 0 {
		return hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("Dashboard '%s' contains duplicate input names for: %s", d.Name(), strings.Join(duplicates, ",")),
			Subject:  &d.DeclRange,
		}}
	}

	var diags hcl.Diagnostics
	//  ensure they inputs not have cyclical dependencies
	if err := d.validateInputDependencies(d.Inputs); err != nil {
		return hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("Failed to resolve input dependency order for dashboard '%s'", d.Name()),
			Detail:   err.Error(),
			Subject:  &d.DeclRange,
		}}
	}
	// now 'claim' all inputs and add to mod
	for _, input := range d.Inputs {
		input.SetDashboard(d)
		moreDiags := d.Mod.AddResource(input)
		diags = append(diags, moreDiags...)
	}

	return diags
}

// populate our input map
func (d *Dashboard) setInputMap() []string {
	var duplicates []string
	d.selfInputsMap = make(map[string]*DashboardInput)
	for _, i := range d.Inputs {
		if _, ok := d.selfInputsMap[i.UnqualifiedName]; ok {
			duplicates = append(duplicates, i.UnqualifiedName)
		} else {
			d.selfInputsMap[i.UnqualifiedName] = i
		}
	}
	return duplicates
}

// CtyValue implements CtyValueProvider
func (d *Dashboard) CtyValue() (cty.Value, error) {
	return cty_helpers.GetCtyValue(d)
}

func (d *Dashboard) SetBaseProperties() hcl.Diagnostics {
	var diags hcl.Diagnostics
	if d.Base == nil {
		return diags
	}
	// copy base into the HclResourceImpl 'base' property so it is accessible to all nested structs
	d.HclResourceImpl.SetBase(d.Base)
	// call into parent nested struct SetBaseProperties
	d.ModTreeItemImpl.SetBaseProperties()

	if d.Width == nil {
		d.Width = d.Base.Width
	}

	if len(d.GetChildren()) == 0 {
		d.Children = d.Base.Children
		d.ChildNames = d.Base.ChildNames
	}

	return d.addBaseInputs(d.Base.Inputs)
}

func (d *Dashboard) addBaseInputs(baseInputs []*DashboardInput) hcl.Diagnostics {
	var diags hcl.Diagnostics
	if len(baseInputs) == 0 {
		return diags
	}
	// rebuild Inputs and children
	inheritedInputs := make([]*DashboardInput, 0, len(baseInputs))
	inheritedChildren := make([]modconfig.ModTreeItem, 0, len(baseInputs))

	for i, baseInput := range baseInputs {
		input := baseInput.Clone()
		input.SetDashboard(d)
		// add to mod
		moreDiags := d.Mod.AddResource(input)
		diags = append(diags, moreDiags...)
		// add to our inputs
		inheritedInputs[i] = input
		inheritedChildren[i] = input
	}

	if !diags.HasErrors() {
		// add inputs to beginning of our existing inputs (if any)
		d.Inputs = append(inheritedInputs, d.Inputs...)
		// add inputs to beginning of our children
		d.Children = append(inheritedChildren, d.Children...)
		d.setInputMap()
	}

	return diags
}

// ensure that dependencies between inputs are resolveable
func (d *Dashboard) validateInputDependencies(inputs []*DashboardInput) error {
	dependencyGraph := topsort.NewGraph()
	rootDependencyNode := "dashboard"
	dependencyGraph.AddNode(rootDependencyNode)
	for _, i := range inputs {
		for _, runtimeDep := range i.GetRuntimeDependencies() {
			depName := runtimeDep.PropertyPath.ToResourceName()
			to := depName
			from := i.UnqualifiedName
			if !dependencyGraph.ContainsNode(from) {
				dependencyGraph.AddNode(from)
			}
			if !dependencyGraph.ContainsNode(to) {
				dependencyGraph.AddNode(to)
			}
			if err := dependencyGraph.AddEdge(from, to); err != nil {
				return err
			}
			if err := dependencyGraph.AddEdge(rootDependencyNode, i.UnqualifiedName); err != nil {
				return err
			}
		}
	}

	// now verify we can get a dependency order
	_, err := dependencyGraph.TopSort(rootDependencyNode)
	return err
}
