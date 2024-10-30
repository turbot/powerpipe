package resources

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/pipe-fittings/cty_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/zclconf/go-cty/cty"
)

// DashboardWith is a struct representing a leaf dashboard node
type DashboardWith struct {
	modconfig.ResourceWithMetadataImpl
	QueryProviderImpl

	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`
}

func NewDashboardWith(block *hcl.Block, mod *modconfig.Mod, shortName string) modconfig.HclResource {
	// with blocks cannot be anonymous
	return &DashboardWith{
		QueryProviderImpl: NewQueryProviderImpl(block, mod, shortName),
	}
}

func (w *DashboardWith) Equals(other *DashboardWith) bool {
	diff := w.Diff(other)
	return !diff.HasChanges()
}

// OnDecoded implements HclResource
func (w *DashboardWith) OnDecoded(_ *hcl.Block, _ modconfig.ModResourcesProvider) hcl.Diagnostics {
	return nil
}

func (w *DashboardWith) Diff(other *DashboardWith) *modconfig.ModTreeItemDiffs {
	res := &modconfig.ModTreeItemDiffs{
		Item: w,
		Name: w.Name(),
	}

	res.Merge(w.QueryProviderImpl.Diff(other))

	return res
}

// GetWidth implements DashboardLeafNode
func (*DashboardWith) GetWidth() int {
	return 0
}

// GetDisplay implements DashboardLeafNode
func (*DashboardWith) GetDisplay() string {
	return ""
}

// GetType implements DashboardLeafNode
func (*DashboardWith) GetType() string {
	return ""
}

// CtyValue implements CtyValueProvider
func (w *DashboardWith) CtyValue() (cty.Value, error) {
	return cty_helpers.GetCtyValue(w)
}
