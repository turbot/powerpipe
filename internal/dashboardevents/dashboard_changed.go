package dashboardevents

import (
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/powerpipe/internal/resources"
)

type DashboardChanged struct {
	ChangedDashboards          []*modconfig.ModTreeItemDiffs
	ChangedContainers          []*modconfig.ModTreeItemDiffs
	ChangedControls            []*modconfig.ModTreeItemDiffs
	ChangedBenchmarks          []*modconfig.ModTreeItemDiffs
	ChangedCategories          []*modconfig.ModTreeItemDiffs
	ChangedCards               []*modconfig.ModTreeItemDiffs
	ChangedCharts              []*modconfig.ModTreeItemDiffs
	ChangedDetections          []*modconfig.ModTreeItemDiffs
	ChangedDetectionBenchmarks []*modconfig.ModTreeItemDiffs
	ChangedFlows               []*modconfig.ModTreeItemDiffs
	ChangedGraphs              []*modconfig.ModTreeItemDiffs
	ChangedHierarchies         []*modconfig.ModTreeItemDiffs
	ChangedImages              []*modconfig.ModTreeItemDiffs
	ChangedInputs              []*modconfig.ModTreeItemDiffs
	ChangedTables              []*modconfig.ModTreeItemDiffs
	ChangedTexts               []*modconfig.ModTreeItemDiffs
	ChangedNodes               []*modconfig.ModTreeItemDiffs
	ChangedEdges               []*modconfig.ModTreeItemDiffs

	NewDashboards          []*resources.Dashboard
	NewContainers          []*resources.DashboardContainer
	NewControls            []*resources.Control
	NewBenchmarks          []*resources.ControlBenchmark
	NewCards               []*resources.DashboardCard
	NewCategories          []*resources.DashboardCategory
	NewCharts              []*resources.DashboardChart
	NewDetections          []*resources.Detection
	NewDetectionBenchmarks []*resources.DetectionBenchmark
	NewFlows               []*resources.DashboardFlow
	NewGraphs              []*resources.DashboardGraph
	NewHierarchies         []*resources.DashboardHierarchy
	NewImages              []*resources.DashboardImage
	NewInputs              []*resources.DashboardInput
	NewTables              []*resources.DashboardTable
	NewTexts               []*resources.DashboardText
	NewNodes               []*resources.DashboardNode
	NewEdges               []*resources.DashboardEdge

	DeletedDashboards          []*resources.Dashboard
	DeletedContainers          []*resources.DashboardContainer
	DeletedControls            []*resources.Control
	DeletedBenchmarks          []*resources.ControlBenchmark
	DeletedCards               []*resources.DashboardCard
	DeletedCategories          []*resources.DashboardCategory
	DeletedCharts              []*resources.DashboardChart
	DeletedDetections          []*resources.Detection
	DeletedDetectionBenchmarks []*resources.DetectionBenchmark
	DeletedFlows               []*resources.DashboardFlow
	DeletedGraphs              []*resources.DashboardGraph
	DeletedHierarchies         []*resources.DashboardHierarchy
	DeletedImages              []*resources.DashboardImage
	DeletedInputs              []*resources.DashboardInput
	DeletedTables              []*resources.DashboardTable
	DeletedTexts               []*resources.DashboardText
	DeletedNodes               []*resources.DashboardNode
	DeletedEdges               []*resources.DashboardEdge
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
		len(c.ChangedDetections)+
		len(c.ChangedDetectionBenchmarks)+
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
		len(c.NewDetections)+
		len(c.NewDetectionBenchmarks)+
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
		len(c.DeletedDetections)+
		len(c.DeletedDetectionBenchmarks)+
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
	for _, r := range c.ChangedDetections {
		if continueWalking, err := resourceFunc(r.Item); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range c.ChangedDetectionBenchmarks {
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

func (c *DashboardChanged) SetParentsChanged(item modconfig.ModTreeItem, prevModResources *resources.PowerpipeModResources) {
	if prevModResources == nil {
		return
	}

	parents := item.GetParents()
	for _, parent := range parents {
		// if the parent DID NOT exist in the previous resource maps, do nothing
		parsedResourceName, _ := modconfig.ParseResourceName(parent.Name())
		if _, existingResource := prevModResources.GetResource(parsedResourceName); existingResource {
			c.AddChanged(parent)
			c.SetParentsChanged(parent, prevModResources)
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
	case *resources.Dashboard:
		if !c.diffsContain(c.ChangedDashboards, item) {
			c.ChangedDashboards = append(c.ChangedDashboards, diff)
		}
	case *resources.DashboardContainer:
		if !c.diffsContain(c.ChangedContainers, item) {
			c.ChangedContainers = append(c.ChangedContainers, diff)
		}
	case *resources.Control:
		if !c.diffsContain(c.ChangedControls, item) {
			c.ChangedControls = append(c.ChangedControls, diff)
		}
	case *resources.ControlBenchmark:
		if !c.diffsContain(c.ChangedBenchmarks, item) {
			c.ChangedBenchmarks = append(c.ChangedBenchmarks, diff)
		}
	case *resources.DashboardCard:
		if !c.diffsContain(c.ChangedCards, item) {
			c.ChangedCards = append(c.ChangedCards, diff)
		}
	case *resources.DashboardCategory:
		if !c.diffsContain(c.ChangedCategories, item) {
			c.ChangedCategories = append(c.ChangedCategories, diff)
		}
	case *resources.DashboardChart:
		if !c.diffsContain(c.ChangedCharts, item) {
			c.ChangedCharts = append(c.ChangedCharts, diff)
		}
	case *resources.DashboardHierarchy:
		if !c.diffsContain(c.ChangedHierarchies, item) {
			c.ChangedHierarchies = append(c.ChangedHierarchies, diff)
		}

	case *resources.DashboardImage:
		if !c.diffsContain(c.ChangedImages, item) {
			c.ChangedImages = append(c.ChangedImages, diff)
		}

	case *resources.DashboardInput:
		if !c.diffsContain(c.ChangedInputs, item) {
			c.ChangedInputs = append(c.ChangedInputs, diff)
		}

	case *resources.DashboardTable:
		if !c.diffsContain(c.ChangedTables, item) {
			c.ChangedTables = append(c.ChangedTables, diff)
		}
	case *resources.Detection:
		if !c.diffsContain(c.ChangedDetections, item) {
			c.ChangedDetections = append(c.ChangedDetections, diff)
		}
	case *resources.DetectionBenchmark:
		if !c.diffsContain(c.ChangedDetectionBenchmarks, item) {
			c.ChangedDetectionBenchmarks = append(c.ChangedDetectionBenchmarks, diff)
		}
	case *resources.DashboardText:
		if !c.diffsContain(c.ChangedTexts, item) {
			c.ChangedTexts = append(c.ChangedTexts, diff)
		}
	}
}
