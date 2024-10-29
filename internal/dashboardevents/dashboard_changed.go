package dashboardevents

import (
	"github.com/turbot/pipe-fittings/modconfig"
	powerpipe2 "github.com/turbot/powerpipe/internal/resources"
)

type DashboardChanged struct {
	ChangedDashboards  []*modconfig.ModTreeItemDiffs
	ChangedContainers  []*modconfig.ModTreeItemDiffs
	ChangedControls    []*modconfig.ModTreeItemDiffs
	ChangedBenchmarks  []*modconfig.ModTreeItemDiffs
	ChangedCategories  []*modconfig.ModTreeItemDiffs
	ChangedCards       []*modconfig.ModTreeItemDiffs
	ChangedCharts      []*modconfig.ModTreeItemDiffs
	ChangedFlows       []*modconfig.ModTreeItemDiffs
	ChangedGraphs      []*modconfig.ModTreeItemDiffs
	ChangedHierarchies []*modconfig.ModTreeItemDiffs
	ChangedImages      []*modconfig.ModTreeItemDiffs
	ChangedInputs      []*modconfig.ModTreeItemDiffs
	ChangedTables      []*modconfig.ModTreeItemDiffs
	ChangedTexts       []*modconfig.ModTreeItemDiffs
	ChangedNodes       []*modconfig.ModTreeItemDiffs
	ChangedEdges       []*modconfig.ModTreeItemDiffs

	NewDashboards  []*powerpipe2.Dashboard
	NewContainers  []*powerpipe2.DashboardContainer
	NewControls    []*powerpipe2.Control
	NewBenchmarks  []*powerpipe2.Benchmark
	NewCards       []*powerpipe2.DashboardCard
	NewCategories  []*powerpipe2.DashboardCategory
	NewCharts      []*powerpipe2.DashboardChart
	NewFlows       []*powerpipe2.DashboardFlow
	NewGraphs      []*powerpipe2.DashboardGraph
	NewHierarchies []*powerpipe2.DashboardHierarchy
	NewImages      []*powerpipe2.DashboardImage
	NewInputs      []*powerpipe2.DashboardInput
	NewTables      []*powerpipe2.DashboardTable
	NewTexts       []*powerpipe2.DashboardText
	NewNodes       []*powerpipe2.DashboardNode
	NewEdges       []*powerpipe2.DashboardEdge

	DeletedDashboards  []*powerpipe2.Dashboard
	DeletedContainers  []*powerpipe2.DashboardContainer
	DeletedControls    []*powerpipe2.Control
	DeletedBenchmarks  []*powerpipe2.Benchmark
	DeletedCards       []*powerpipe2.DashboardCard
	DeletedCategories  []*powerpipe2.DashboardCategory
	DeletedCharts      []*powerpipe2.DashboardChart
	DeletedFlows       []*powerpipe2.DashboardFlow
	DeletedGraphs      []*powerpipe2.DashboardGraph
	DeletedHierarchies []*powerpipe2.DashboardHierarchy
	DeletedImages      []*powerpipe2.DashboardImage
	DeletedInputs      []*powerpipe2.DashboardInput
	DeletedTables      []*powerpipe2.DashboardTable
	DeletedTexts       []*powerpipe2.DashboardText
	DeletedNodes       []*powerpipe2.DashboardNode
	DeletedEdges       []*powerpipe2.DashboardEdge
}

// IsDashboardEvent implements DashboardEvent interface
func (*DashboardChanged) IsDashboardEvent() {}

func (c *DashboardChanged) HasChanges() bool {
	return len(c.ChangedDashboards)+
		len(c.ChangedContainers)+
		len(c.ChangedBenchmarks)+
		len(c.ChangedControls)+
		len(c.ChangedCards)+
		len(c.ChangedCategories)+
		len(c.ChangedCharts)+
		len(c.ChangedFlows)+
		len(c.ChangedGraphs)+
		len(c.ChangedHierarchies)+
		len(c.ChangedImages)+
		len(c.ChangedInputs)+
		len(c.ChangedTables)+
		len(c.ChangedTexts)+
		len(c.ChangedNodes)+
		len(c.ChangedEdges)+
		len(c.NewDashboards)+
		len(c.NewContainers)+
		len(c.NewBenchmarks)+
		len(c.NewControls)+
		len(c.NewCards)+
		len(c.NewCategories)+
		len(c.NewCharts)+
		len(c.NewFlows)+
		len(c.NewGraphs)+
		len(c.NewHierarchies)+
		len(c.NewImages)+
		len(c.NewInputs)+
		len(c.NewTables)+
		len(c.NewTexts)+
		len(c.NewNodes)+
		len(c.NewEdges)+
		len(c.DeletedDashboards)+
		len(c.DeletedContainers)+
		len(c.DeletedBenchmarks)+
		len(c.DeletedControls)+
		len(c.DeletedCards)+
		len(c.DeletedCategories)+
		len(c.DeletedCharts)+
		len(c.DeletedFlows)+
		len(c.DeletedGraphs)+
		len(c.DeletedHierarchies)+
		len(c.DeletedImages)+
		len(c.DeletedInputs)+
		len(c.DeletedTables)+
		len(c.DeletedTexts)+
		len(c.DeletedNodes)+
		len(c.DeletedEdges) > 0
}

func (c *DashboardChanged) WalkChangedResources(resourceFunc func(item modconfig.ModTreeItem) (bool, error)) error {
	for _, r := range c.ChangedDashboards {
		if continueWalking, err := resourceFunc(r.Item); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.ChangedContainers {
		if continueWalking, err := resourceFunc(r.Item); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.ChangedControls {
		if continueWalking, err := resourceFunc(r.Item); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.ChangedCards {
		if continueWalking, err := resourceFunc(r.Item); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.ChangedCategories {
		if continueWalking, err := resourceFunc(r.Item); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.ChangedCharts {
		if continueWalking, err := resourceFunc(r.Item); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.ChangedFlows {
		if continueWalking, err := resourceFunc(r.Item); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.ChangedGraphs {
		if continueWalking, err := resourceFunc(r.Item); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.ChangedHierarchies {
		if continueWalking, err := resourceFunc(r.Item); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.ChangedImages {
		if continueWalking, err := resourceFunc(r.Item); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.ChangedInputs {
		if continueWalking, err := resourceFunc(r.Item); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.ChangedTables {
		if continueWalking, err := resourceFunc(r.Item); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.ChangedTexts {
		if continueWalking, err := resourceFunc(r.Item); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.NewDashboards {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.NewContainers {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.NewControls {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.NewCards {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.NewCategories {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.NewCharts {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.NewFlows {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.NewGraphs {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.NewHierarchies {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.NewImages {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.NewInputs {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.NewTables {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.NewTexts {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.DeletedContainers {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.DeletedControls {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.DeletedCards {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.DeletedCategories {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.DeletedCharts {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.DeletedFlows {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.DeletedGraphs {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.DeletedHierarchies {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.DeletedImages {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.DeletedInputs {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.DeletedTables {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.DeletedTexts {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}

	return nil
}

func (c *DashboardChanged) SetParentsChanged(item modconfig.ModTreeItem, prevResourceMaps *powerpipe2.ModResources) {
	if prevResourceMaps == nil {
		return
	}

	parents := item.GetParents()
	for _, parent := range parents {
		// if the parent DID NOT exist in the previous resource maps, do nothing
		parsedResourceName, _ := modconfig.ParseResourceName(parent.Name())
		if _, existingResource := prevResourceMaps.GetResource(parsedResourceName); existingResource {
			c.AddChanged(parent)
			c.SetParentsChanged(parent, prevResourceMaps)
		}
	}
}

func (c *DashboardChanged) diffsContain(diffs []*modconfig.ModTreeItemDiffs, item modconfig.ModTreeItem) bool {
	for _, d := range diffs {
		if d.Item.Name() == item.Name() {
			return true
		}
	}
	return false
}

func (c *DashboardChanged) AddChanged(item modconfig.ModTreeItem) {
	diff := &modconfig.ModTreeItemDiffs{
		Name:              item.Name(),
		Item:              item,
		ChangedProperties: []string{"Children"},
	}
	switch item.(type) {
	case *powerpipe2.Dashboard:
		if !c.diffsContain(c.ChangedDashboards, item) {
			c.ChangedDashboards = append(c.ChangedDashboards, diff)
		}
	case *powerpipe2.DashboardContainer:
		if !c.diffsContain(c.ChangedContainers, item) {
			c.ChangedContainers = append(c.ChangedContainers, diff)
		}
	case *powerpipe2.Control:
		if !c.diffsContain(c.ChangedControls, item) {
			c.ChangedControls = append(c.ChangedControls, diff)
		}
	case *powerpipe2.Benchmark:
		if !c.diffsContain(c.ChangedBenchmarks, item) {
			c.ChangedBenchmarks = append(c.ChangedBenchmarks, diff)
		}
	case *powerpipe2.DashboardCard:
		if !c.diffsContain(c.ChangedCards, item) {
			c.ChangedCards = append(c.ChangedCards, diff)
		}
	case *powerpipe2.DashboardCategory:
		if !c.diffsContain(c.ChangedCategories, item) {
			c.ChangedCategories = append(c.ChangedCategories, diff)
		}
	case *powerpipe2.DashboardChart:
		if !c.diffsContain(c.ChangedCharts, item) {
			c.ChangedCharts = append(c.ChangedCharts, diff)
		}
	case *powerpipe2.DashboardHierarchy:
		if !c.diffsContain(c.ChangedHierarchies, item) {
			c.ChangedHierarchies = append(c.ChangedHierarchies, diff)
		}

	case *powerpipe2.DashboardImage:
		if !c.diffsContain(c.ChangedImages, item) {
			c.ChangedImages = append(c.ChangedImages, diff)
		}

	case *powerpipe2.DashboardInput:
		if !c.diffsContain(c.ChangedInputs, item) {
			c.ChangedInputs = append(c.ChangedInputs, diff)
		}

	case *powerpipe2.DashboardTable:
		if !c.diffsContain(c.ChangedTables, item) {
			c.ChangedTables = append(c.ChangedTables, diff)
		}
	case *powerpipe2.DashboardText:
		if !c.diffsContain(c.ChangedTexts, item) {
			c.ChangedTexts = append(c.ChangedTexts, diff)
		}
	}
}
