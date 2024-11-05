package resources

import (
	"github.com/hashicorp/hcl/v2"
	typehelpers "github.com/turbot/go-kit/types"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/cty_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/printers"
	"github.com/turbot/pipe-fittings/schema"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/zclconf/go-cty/cty"
)

// DashboardTable is a struct representing a leaf dashboard node
type DashboardTable struct {
	modconfig.ResourceWithMetadataImpl
	QueryProviderImpl

	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`

	Width      *int                             `cty:"width" hcl:"width"  json:"width,omitempty"`
	Type       *string                          `cty:"type" hcl:"type"  json:"type,omitempty"`
	ColumnList DashboardTableColumnList         `cty:"column_list" hcl:"column,block" json:"columns,omitempty"`
	Columns    map[string]*DashboardTableColumn `cty:"columns" snapshot:"columns"`
	Display    *string                          `cty:"display" hcl:"display" json:"display,omitempty" snapshot:"display"`
	Base       *DashboardTable                  `hcl:"base" json:"-"`
}

func NewDashboardTable(block *hcl.Block, mod *modconfig.Mod, shortName string) modconfig.HclResource {
	t := &DashboardTable{
		QueryProviderImpl: NewQueryProviderImpl(block, mod, shortName),
	}
	t.SetAnonymous(block)
	return t
}

// NewQueryDashboardTable creates a Table to wrap a query.
// This is used in order to execute queries as dashboards
func NewQueryDashboardTable(qp QueryProvider) (*DashboardTable, error) {
	parsedName, err := modconfig.ParseResourceName(constants.SnapshotQueryTableName)
	if err != nil {
		return nil, err
	}
	fullName := parsedName.ToFullName()

	c := &DashboardTable{
		QueryProviderImpl: QueryProviderImpl{
			RuntimeDependencyProviderImpl: RuntimeDependencyProviderImpl{
				ModTreeItemImpl: modconfig.ModTreeItemImpl{
					HclResourceImpl: modconfig.HclResourceImpl{
						ShortName:       parsedName.Name,
						FullName:        fullName,
						UnqualifiedName: parsedName.ToResourceName(),
						Title:           utils.ToStringPointer(qp.GetTitle()),
						BlockType:       schema.BlockTypeTable,
					},
					Database: qp.GetDatabase(),
					Mod:      qp.(modconfig.ModItem).GetMod(),
				},
			},
			Query:  qp.GetQuery(),
			SQL:    qp.GetSQL(),
			Params: qp.GetParams(),
			Args:   qp.GetArgs(),
		},
	}
	return c, nil
}

func (t *DashboardTable) Equals(other *DashboardTable) bool {
	diff := t.Diff(other)
	return !diff.HasChanges()
}

// OnDecoded implements HclResource
func (t *DashboardTable) OnDecoded(block *hcl.Block, resourceMapProvider modconfig.ModResourcesProvider) hcl.Diagnostics {
	t.SetBaseProperties()
	// populate columns map
	if len(t.ColumnList) > 0 {
		t.Columns = make(map[string]*DashboardTableColumn, len(t.ColumnList))
		for _, c := range t.ColumnList {
			t.Columns[c.Name] = c
		}
	}
	return t.QueryProviderImpl.OnDecoded(block, resourceMapProvider)
}

func (t *DashboardTable) Diff(other *DashboardTable) *modconfig.ModTreeItemDiffs {
	res := &modconfig.ModTreeItemDiffs{
		Item: t,
		Name: t.Name(),
	}

	if !utils.SafeStringsEqual(t.Type, other.Type) {
		res.AddPropertyDiff("Type")
	}

	if len(t.ColumnList) != len(other.ColumnList) {
		res.AddPropertyDiff("Columns")
	} else {
		for i, c := range t.Columns {
			if !c.Equals(other.Columns[i]) {
				res.AddPropertyDiff("Columns")
			}
		}
	}

	res.PopulateChildDiffs(t, other)
	res.Merge(t.QueryProviderImpl.Diff(other))
	res.Merge(dashboardLeafNodeDiff(t, other))

	return res
}

// GetWidth implements DashboardLeafNode
func (t *DashboardTable) GetWidth() int {
	if t.Width == nil {
		return 0
	}
	return *t.Width
}

// GetDisplay implements DashboardLeafNode
func (t *DashboardTable) GetDisplay() string {
	return typehelpers.SafeString(t.Display)
}

// GetDocumentation implements DashboardLeafNode, ModTreeItem
func (*DashboardTable) GetDocumentation() string {
	return ""
}

// GetType implements DashboardLeafNode
func (t *DashboardTable) GetType() string {
	return typehelpers.SafeString(t.Type)
}

// CtyValue implements CtyValueProvider
func (t *DashboardTable) CtyValue() (cty.Value, error) {
	return cty_helpers.GetCtyValue(t)
}

func (t *DashboardTable) SetBaseProperties() {
	if t.Base == nil {
		return
	}
	// copy base into the HclResourceImpl 'base' property so it is accessible to all nested structs
	t.HclResourceImpl.SetBase(t.Base)

	// call into parent nested struct SetBaseProperties
	t.QueryProviderImpl.SetBaseProperties()

	if t.Width == nil {
		t.Width = t.Base.Width
	}

	if t.Type == nil {
		t.Type = t.Base.Type
	}

	if t.Display == nil {
		t.Display = t.Base.Display
	}

	if t.ColumnList == nil {
		t.ColumnList = t.Base.ColumnList
	} else {
		t.ColumnList.Merge(t.Base.ColumnList)
	}
}

// GetShowData implements printers.Showable
func (t *DashboardTable) GetShowData() *printers.RowData {
	res := printers.NewRowData(
		printers.NewFieldValue("Width", t.Width),
		printers.NewFieldValue("Type", t.Type),
		printers.NewFieldValue("Display", t.Display),
		printers.NewFieldValue("Columns", t.ColumnList),
	)
	// merge fields from base, putting base fields first
	res.Merge(t.QueryProviderImpl.GetShowData())
	return res
}
