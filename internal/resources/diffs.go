package resources

import (
	"github.com/turbot/pipe-fittings/modconfig"
	"maps"
)

// TODO replace with dashboardLeafNodeImple.Diff
// need to think about also doing ModTreeItemImple.Diff and HclResourceImpl.Diff
// think about which poroperties to diff
func dashboardLeafNodeDiff(l DashboardLeafNode, r DashboardLeafNode) *modconfig.ModTreeItemDiffs {
	d := &modconfig.ModTreeItemDiffs{}
	if l.Name() != r.Name() {
		d.AddPropertyDiff("Name")
	}
	if l.GetTitle() != r.GetTitle() {
		d.AddPropertyDiff("Title")
	}
	if l.GetWidth() != r.GetWidth() {
		d.AddPropertyDiff("Width")
	}
	if l.GetDisplay() != r.GetDisplay() {
		d.AddPropertyDiff("Display")
	}
	if l.GetDocumentation() != r.GetDocumentation() {
		d.AddPropertyDiff("Documentation")
	}
	if l.GetType() != r.GetType() {
		d.AddPropertyDiff("Type")
	}
	if !maps.Equal(l.GetTags(), r.GetTags()) {
		d.AddPropertyDiff("Tags")
	}
	return d
}
