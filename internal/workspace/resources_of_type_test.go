package workspace

import (
	"github.com/turbot/pipe-fittings/modconfig/powerpipe"
	"github.com/turbot/pipe-fittings/workspace"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/utils"
)

func makeControl(mod *modconfig.Mod, name, title, description, sql string, tags map[string]string) *powerpipe.Control {
	control := powerpipe.NewControl(&hcl.Block{Type: "control"}, mod, name).(*powerpipe.Control)
	control.Title = &title
	control.Description = &description
	control.Tags = tags
	control.SQL = &sql
	return control
}

type testCase[T modconfig.HclResource] struct {
	name   string
	filter workspace.ResourceFilter
	want   map[string]struct{}
}

func TestFilterWorkspaceResourcesOfType(t *testing.T) {
	// Set the AppSpecificNewResourceMapsFunc to the Powerpipe NewResourceMaps function
	modconfig.AppSpecificNewResourceMapsFunc = powerpipe.NewModResources

	var mod = powerpipe.NewMod("test_mod", ".", hcl.Range{})
	mod.ResourceMaps = &powerpipe.ModResources{
		Benchmarks: map[string]*powerpipe.Benchmark{},
		Controls: map[string]*powerpipe.Control{
			"control1":  makeControl(mod, "control1", "Control 1", "Control 1 description", "SELECT * FROM table1", map[string]string{"t1": "val1_foo", "t2": "val2_foo", "t3": "val3_foo"}),
			"control2a": makeControl(mod, "control2a", "Control 2", "Control 2a description", "SELECT id FROM table2", map[string]string{"t1": "val1_foo", "t2": "val2_foo", "t3": "val3_foo_a"}),
			"control2b": makeControl(mod, "control2b", "Control 2", "Control 2b description", "SELECT * FROM table2", map[string]string{"t1": "val1_foo", "t2": "val2_foo", "t3": "val3_foo_b"}),
			"control3":  makeControl(mod, "control3", "Control 3", "Control 3 description", "SELECT * FROM table3", map[string]string{"t1": "val1_bar", "t2": "val2_bar", "t3": "val3_bar"}),
			"control4":  makeControl(mod, "control4", "Control 4", "Control 4 description", "SELECT * FROM table4", map[string]string{"t1": "val1_bar", "t2": "val2_foo", "t3": "val3_bar"}),
		},
	}
	var w = &PowerpipeWorkspace{
		Workspace: *workspace.Workspace{
			Mod: mod,
		},
	}

	controlTests := []testCase[*powerpipe.Control]{
		{
			name: `where "name = 'control1'"`,
			filter: workspace.ResourceFilter{
				Where: "name = 'control1'",
			},

			want: map[string]struct{}{
				"test_mod.control.control1": {},
			},
		},
		{
			name: `where "name != 'control1'"`,
			filter: workspace.ResourceFilter{
				Where: "name != 'control1'",
			},
			want: map[string]struct{}{
				"test_mod.control.control2a": {},
				"test_mod.control.control2b": {},
				"test_mod.control.control3":  {},
				"test_mod.control.control4":  {},
			},
		},
		{
			name: `where "name like 'control2%'"`,
			filter: workspace.ResourceFilter{
				Where: `name like 'control2%'`,
			},
			want: map[string]struct{}{
				"test_mod.control.control2a": {},
				"test_mod.control.control2b": {},
			},
		},
		{
			name: `where "name ilike 'ConTrol2%'"`,
			filter: workspace.ResourceFilter{
				Where: `name ilike 'ConTrol2%'`,
			},
			want: map[string]struct{}{
				"test_mod.control.control2a": {},
				"test_mod.control.control2b": {},
			},
		},
		{
			name: `where "name not like 'control2%'"`,
			filter: workspace.ResourceFilter{
				Where: `name not like 'control2%'`,
			},
			want: map[string]struct{}{
				"test_mod.control.control1": {},
				"test_mod.control.control3": {},
				"test_mod.control.control4": {},
			},
		},
		{
			name: `tags t1=val1_foo t2=val2_foo`,
			filter: workspace.ResourceFilter{
				Tags: map[string][]string{
					"t1": {"val1_foo"},
					"t2": {"val2_foo"},
				},
			},
			want: map[string]struct{}{
				"test_mod.control.control1":  {},
				"test_mod.control.control2a": {},
				"test_mod.control.control2b": {},
			},
		},
		{
			name: `tags t1=val1_bar t2=val2_bar`,
			filter: workspace.ResourceFilter{
				Tags: map[string][]string{
					"t1": {"val1_bar"},
					"t2": {"val2_bar"},
				},
			},
			want: map[string]struct{}{
				"test_mod.control.control3": {},
			},
		},
		{
			name: `tags t3=val3_foo t3=val3_bar`,
			filter: workspace.ResourceFilter{
				Tags: map[string][]string{
					"t3": {"val3_foo", "val3_bar"},
				},
			},
			want: map[string]struct{}{
				"test_mod.control.control1": {},
				"test_mod.control.control3": {},
				"test_mod.control.control4": {},
			},
		},
		{
			name: `tags t1=val1_foo t2=something_else [NO MATCHES]`,
			filter: workspace.ResourceFilter{
				Tags: map[string][]string{
					"t1": {"val1_foo"},
					"t2": {"something_else"},
				},
			},
			want: map[string]struct{}{},
		},
	}
	//var testFilter = "name like 'control1'"
	var testFilter = ""

	executeTests[*powerpipe.Control](t, controlTests, testFilter, w)
}

func executeTests[T modconfig.HclResource](t *testing.T, controlTests []testCase[*powerpipe.Control], testFilter string, w *PowerpipeWorkspace) {
	for _, tt := range controlTests {
		// apply test filter if specified
		if testFilter != "" && tt.name != testFilter {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {

			got, err := workspace.FilterWorkspaceResourcesOfType[T](w, tt.filter)
			if err != nil {
				t.Fatalf("FilterWorkspaceResourcesOfType() test '%s' error = %v", tt.name, err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("FilterWorkspaceResourcesOfType() test '%s' got %d %s, wanted %d",
					tt.name,
					len(got), utils.Pluralize("result", len(got)),
					len(tt.want))
			}
			for k := range got {
				if _, found := tt.want[k]; !found {
					t.Errorf("FilterWorkspaceResourcesOfType() test '%s' got %s but this was not expected", tt.name, k)
				}
			}
		})
	}
}
