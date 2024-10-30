package resources

import (
	"github.com/hashicorp/hcl/v2"
	typehelpers "github.com/turbot/go-kit/types"
	"github.com/turbot/pipe-fittings/cty_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/printers"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/zclconf/go-cty/cty"
)

// DashboardChart is a struct representing a leaf dashboard node
type DashboardChart struct {
	modconfig.ResourceWithMetadataImpl
	QueryProviderImpl

	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`

	Width      *int                             `cty:"width" hcl:"width" json:"width,omitempty"`
	Type       *string                          `cty:"type" hcl:"type" json:"type,omitempty"`
	Display    *string                          `cty:"display" hcl:"display" json:"display,omitempty"`
	Legend     *DashboardChartLegend            `cty:"legend" hcl:"legend,block" snapshot:"legend" json:"legend,omitempty"`
	SeriesList DashboardChartSeriesList         `cty:"series_list" hcl:"series,block" json:"series,omitempty"`
	Axes       *DashboardChartAxes              `cty:"axes" hcl:"axes,block" snapshot:"axes" json:"axes,omitempty"`
	Grouping   *string                          `cty:"grouping" hcl:"grouping" snapshot:"grouping" json:"grouping,omitempty"`
	Transform  *string                          `cty:"transform" hcl:"transform" snapshot:"transform" json:"transform,omitempty"`
	Series     map[string]*DashboardChartSeries `cty:"series" snapshot:"series"`
	Base       *DashboardChart                  `hcl:"base" json:"-"`
}

func NewDashboardChart(block *hcl.Block, mod *modconfig.Mod, shortName string) modconfig.HclResource {
	c := &DashboardChart{
		QueryProviderImpl: NewQueryProviderImpl(block, mod, shortName),
	}

	c.SetAnonymous(block)
	return c
}

func (c *DashboardChart) Equals(other *DashboardChart) bool {
	diff := c.Diff(other)
	return !diff.HasChanges()
}

// OnDecoded implements HclResource
func (c *DashboardChart) OnDecoded(block *hcl.Block, resourceMapProvider modconfig.ModResourcesProvider) hcl.Diagnostics {
	c.SetBaseProperties()
	// populate series map
	if len(c.SeriesList) > 0 {
		c.Series = make(map[string]*DashboardChartSeries, len(c.SeriesList))
		for _, s := range c.SeriesList {
			s.OnDecoded()
			c.Series[s.Name] = s
		}
	}
	return c.QueryProviderImpl.OnDecoded(block, resourceMapProvider)
}

func (c *DashboardChart) Diff(other *DashboardChart) *modconfig.ModTreeItemDiffs {
	res := &modconfig.ModTreeItemDiffs{
		Item: c,
		Name: c.Name(),
	}

	if !utils.SafeStringsEqual(c.Type, other.Type) {
		res.AddPropertyDiff("Type")
	}

	if !utils.SafeStringsEqual(c.Grouping, other.Grouping) {
		res.AddPropertyDiff("Grouping")
	}

	if !utils.SafeStringsEqual(c.Transform, other.Transform) {
		res.AddPropertyDiff("Transform")
	}

	if len(c.SeriesList) != len(other.SeriesList) {
		res.AddPropertyDiff("Series")
	} else {
		for i, s := range c.Series {
			if !s.Equals(other.Series[i]) {
				res.AddPropertyDiff("Series")
			}
		}
	}

	if c.Legend != nil {
		if !c.Legend.Equals(other.Legend) {
			res.AddPropertyDiff("Legend")
		}
	} else if other.Legend != nil {
		res.AddPropertyDiff("Legend")
	}

	if c.Axes != nil {
		if !c.Axes.Equals(other.Axes) {
			res.AddPropertyDiff("Axes")
		}
	} else if other.Axes != nil {
		res.AddPropertyDiff("Axes")
	}

	res.PopulateChildDiffs(c, other)
	c.QueryProviderImpl.Diff(other)
	res.Merge(dashboardLeafNodeDiff(c, other))

	return res
}

// GetWidth implements DashboardLeafNode
func (c *DashboardChart) GetWidth() int {
	if c.Width == nil {
		return 0
	}
	return *c.Width
}

// GetDisplay implements DashboardLeafNode
func (c *DashboardChart) GetDisplay() string {
	return typehelpers.SafeString(c.Display)
}

// GetType implements DashboardLeafNode
func (c *DashboardChart) GetType() string {
	return typehelpers.SafeString(c.Type)
}

// CtyValue implements CtyValueProvider
func (c *DashboardChart) CtyValue() (cty.Value, error) {
	return cty_helpers.GetCtyValue(c)
}

func (c *DashboardChart) SetBaseProperties() {
	if c.Base == nil {
		return
	}
	// copy base into the HclResourceImpl 'base' property so it is accessible to all nested structs
	c.HclResourceImpl.SetBase(c.Base)
	// call into parent nested struct SetBaseProperties
	c.QueryProviderImpl.SetBaseProperties()

	if c.Type == nil {
		c.Type = c.Base.Type
	}

	if c.Display == nil {
		c.Display = c.Base.Display
	}

	if c.Axes == nil {
		c.Axes = c.Base.Axes
	} else {
		c.Axes.Merge(c.Base.Axes)
	}

	if c.Grouping == nil {
		c.Grouping = c.Base.Grouping
	}

	if c.Transform == nil {
		c.Transform = c.Base.Transform
	}

	if c.Legend == nil {
		c.Legend = c.Base.Legend
	} else {
		c.Legend.Merge(c.Base.Legend)
	}

	if c.SeriesList == nil {
		c.SeriesList = c.Base.SeriesList
	} else {
		c.SeriesList.Merge(c.Base.SeriesList)
	}

	if c.Width == nil {
		c.Width = c.Base.Width
	}
}

// GetShowData implements printers.Showable
func (c *DashboardChart) GetShowData() *printers.RowData {
	res := printers.NewRowData(
		printers.NewFieldValue("Width", c.Width),
		printers.NewFieldValue("Type", c.Type),
		printers.NewFieldValue("Display", c.Display),
		printers.NewFieldValue("Grouping", c.Grouping),
		printers.NewFieldValue("Transform", c.Transform),
		printers.NewFieldValue("Legend", c.Legend),
		printers.NewFieldValue("Series", c.Series),
		printers.NewFieldValue("Axes", c.Axes),
	)
	// merge fields from base, putting base fields first
	res.Merge(c.QueryProviderImpl.GetShowData())
	return res
}
