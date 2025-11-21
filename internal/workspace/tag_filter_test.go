package workspace

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	pfworkspace "github.com/turbot/pipe-fittings/v2/workspace"
	"github.com/turbot/powerpipe/internal/resources"
)

func TestResourceFilterFromTagArgs(t *testing.T) {
	mod := modconfig.NewMod("tag_filter_mod", ".", hcl.Range{})
	controls := map[string]*resources.Control{
		"deprecated_true":  makeTaggedControl(t, mod, "deprecated_true", map[string]string{"deprecated": "true"}),
		"deprecated_false": makeTaggedControl(t, mod, "deprecated_false", map[string]string{"deprecated": "false"}),
		"no_tag":           makeTaggedControl(t, mod, "no_tag", map[string]string{}),
		"other_tag_only":   makeTaggedControl(t, mod, "other_tag_only", map[string]string{"env": "qa"}),
	}
	modResources := resources.NewModResources(mod).(*resources.PowerpipeModResources)
	for _, c := range controls {
		_ = modResources.AddResource(c)
	}
	mod.Resources = modResources

	w := &PowerpipeWorkspace{
		Workspace: pfworkspace.Workspace{
			Mod: mod,
		},
	}

	tests := []struct {
		name      string
		tagArgs   []string
		wantNames map[string]struct{}
	}{
		{
			name:    "equals match",
			tagArgs: []string{"deprecated=true"},
			wantNames: map[string]struct{}{
				"tag_filter_mod.control.deprecated_true": {},
			},
		},
		{
			name:    "not equals includes missing",
			tagArgs: []string{"deprecated!=true"},
			wantNames: map[string]struct{}{
				"tag_filter_mod.control.deprecated_false": {},
				"tag_filter_mod.control.no_tag":           {},
				"tag_filter_mod.control.other_tag_only":   {},
			},
		},
		{
			name:    "not equals multiple values includes missing",
			tagArgs: []string{"deprecated!=true", "deprecated!=false"},
			wantNames: map[string]struct{}{
				"tag_filter_mod.control.no_tag":         {},
				"tag_filter_mod.control.other_tag_only": {},
			},
		},
		{
			name:    "mix equals and not equals honors both",
			tagArgs: []string{"deprecated!=true", "env=qa"},
			wantNames: map[string]struct{}{
				"tag_filter_mod.control.other_tag_only": {},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := ResourceFilterFromTagArgs(tt.tagArgs)
			got, err := pfworkspace.FilterWorkspaceResourcesOfType[*resources.Control](&w.Workspace, filter)
			if err != nil {
				t.Fatalf("FilterWorkspaceResourcesOfType error: %v", err)
			}

			if len(got) != len(tt.wantNames) {
				t.Fatalf("expected %d results, got %d", len(tt.wantNames), len(got))
			}
			for name := range tt.wantNames {
				if _, ok := got[name]; !ok {
					t.Fatalf("expected %s but not present", name)
				}
			}
		})
	}
}

func makeTaggedControl(t *testing.T, mod *modconfig.Mod, name string, tags map[string]string) *resources.Control {
	t.Helper()

	title := "title"
	description := "desc"
	sql := "select 1"
	control := resources.NewControl(&hcl.Block{Type: "control"}, mod, name).(*resources.Control)
	control.Title = &title
	control.Description = &description
	control.SQL = &sql
	control.Tags = tags

	return control
}
