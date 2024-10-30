package workspace

import (
	"context"
	"github.com/turbot/powerpipe/internal/resources"
	"log/slog"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/powerpipe/internal/dashboardevents"
)

var EventCount int64 = 0

func (w *PowerpipeWorkspace) PublishDashboardEvent(ctx context.Context, e dashboardevents.DashboardEvent) {
	if w.dashboardEventChan != nil {
		var doneChan = make(chan struct{})
		go func() {
			// send an event onto the event bus
			w.dashboardEventChan <- e
			atomic.AddInt64(&EventCount, 1)
			close(doneChan)
		}()
		select {
		case <-doneChan:
		case <-time.After(1 * time.Second):
			slog.Debug("timeout sending dashboard event", "event", reflect.TypeOf(e).String(), "buffered events", EventCount)
		case <-ctx.Done():
			slog.Debug("context cancelled sending dashboard event")
		}
	}
}

// RegisterDashboardEventHandler starts the event handler goroutine if necessary and
// adds the event handler to our list
func (w *PowerpipeWorkspace) RegisterDashboardEventHandler(ctx context.Context, handler dashboardevents.DashboardEventHandler) {
	// if no event channel has been created we need to start the event handler goroutine
	if w.dashboardEventChan == nil {
		// create a fairly large channel buffer
		w.dashboardEventChan = make(chan dashboardevents.DashboardEvent, 256)
		go w.handleDashboardEvent(ctx)
	}
	// now add the handler to our list
	w.dashboardEventHandlers = append(w.dashboardEventHandlers, handler)
}

// UnregisterDashboardEventHandlers clears all event handlers
// used when generating multiple snapshots
func (w *PowerpipeWorkspace) UnregisterDashboardEventHandlers() {
	w.dashboardEventHandlers = nil
}

// this function is run as a goroutine to call registered event handlers for all received events
func (w *PowerpipeWorkspace) handleDashboardEvent(ctx context.Context) {
	for {
		e := <-w.dashboardEventChan
		atomic.AddInt64(&EventCount, -1)
		if e == nil {
			slog.Debug("handleDashboardEvent nil event received - exiting")
			w.dashboardEventChan = nil
			return
		}

		for _, handler := range w.dashboardEventHandlers {
			handler(ctx, e)
		}
	}
}

func (w *PowerpipeWorkspace) raiseDashboardChangedEvents(ctx context.Context, r, p modconfig.ModResources) {
	event := &dashboardevents.DashboardChanged{}

	resourceMaps := r.(*resources.PowerpipeModResources)
	prevResourceMaps := p.(*resources.PowerpipeModResources)

	// TODO reports can we use a PowerpipeModResources diff function to do all of this - we are duplicating logic

	// first detect changes to existing resources and deletions
	for name, prev := range prevResourceMaps.Dashboards {
		if current, ok := resourceMaps.Dashboards[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedDashboards = append(event.ChangedDashboards, diff)
			}
		} else {
			event.DeletedDashboards = append(event.DeletedDashboards, prev)
		}
	}
	for name, prev := range prevResourceMaps.DashboardContainers {
		if current, ok := resourceMaps.DashboardContainers[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedContainers = append(event.ChangedContainers, diff)
			}
		} else {
			event.DeletedContainers = append(event.DeletedContainers, prev)
		}
	}
	for name, prev := range prevResourceMaps.DashboardCards {
		if current, ok := resourceMaps.DashboardCards[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedCards = append(event.ChangedCards, diff)
			}
		} else {
			event.DeletedCards = append(event.DeletedCards, prev)
		}
	}
	for name, prev := range prevResourceMaps.DashboardCharts {
		if current, ok := resourceMaps.DashboardCharts[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedCharts = append(event.ChangedCharts, diff)
			}
		} else {
			event.DeletedCharts = append(event.DeletedCharts, prev)
		}
	}
	for name, prev := range prevResourceMaps.Benchmarks {
		if current, ok := resourceMaps.Benchmarks[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedBenchmarks = append(event.ChangedBenchmarks, diff)
			}
		} else {
			event.DeletedBenchmarks = append(event.DeletedBenchmarks, prev)
		}
	}
	for name, prev := range prevResourceMaps.Controls {
		if current, ok := resourceMaps.Controls[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedControls = append(event.ChangedControls, diff)
			}
		} else {
			event.DeletedControls = append(event.DeletedControls, prev)
		}
	}
	for name, prev := range prevResourceMaps.DashboardFlows {
		if current, ok := resourceMaps.DashboardFlows[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedFlows = append(event.ChangedFlows, diff)
			}
		} else {
			event.DeletedFlows = append(event.DeletedFlows, prev)
		}
	}
	for name, prev := range prevResourceMaps.DashboardGraphs {
		if current, ok := resourceMaps.DashboardGraphs[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedGraphs = append(event.ChangedGraphs, diff)
			}
		} else {
			event.DeletedGraphs = append(event.DeletedGraphs, prev)
		}
	}
	for name, prev := range prevResourceMaps.DashboardHierarchies {
		if current, ok := resourceMaps.DashboardHierarchies[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedHierarchies = append(event.ChangedHierarchies, diff)
			}
		} else {
			event.DeletedHierarchies = append(event.DeletedHierarchies, prev)
		}
	}
	for name, prev := range prevResourceMaps.DashboardImages {
		if current, ok := resourceMaps.DashboardImages[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedImages = append(event.ChangedImages, diff)
			}
		} else {
			event.DeletedImages = append(event.DeletedImages, prev)
		}
	}
	for name, prev := range prevResourceMaps.DashboardNodes {
		if current, ok := resourceMaps.DashboardNodes[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedNodes = append(event.ChangedNodes, diff)
			}
		} else {
			event.DeletedNodes = append(event.DeletedNodes, prev)
		}
	}
	for name, prev := range prevResourceMaps.DashboardEdges {
		if current, ok := resourceMaps.DashboardEdges[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedEdges = append(event.ChangedEdges, diff)
			}
		} else {
			event.DeletedEdges = append(event.DeletedEdges, prev)
		}
	}
	for name, prev := range prevResourceMaps.GlobalDashboardInputs {
		if current, ok := resourceMaps.GlobalDashboardInputs[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedInputs = append(event.ChangedInputs, diff)
			}
		} else {
			event.DeletedInputs = append(event.DeletedInputs, prev)
		}
	}
	for name, prevInputsForDashboard := range prevResourceMaps.DashboardInputs {
		if currentInputsForDashboard, ok := resourceMaps.DashboardInputs[name]; ok {
			for name, prev := range prevInputsForDashboard {
				if current, ok := currentInputsForDashboard[name]; ok {
					diff := prev.Diff(current)
					if diff.HasChanges() {
						event.ChangedInputs = append(event.ChangedInputs, diff)
					}
				} else {
					event.DeletedInputs = append(event.DeletedInputs, prev)
				}
			}
		} else {
			for _, prev := range prevInputsForDashboard {
				event.DeletedInputs = append(event.DeletedInputs, prev)
			}
		}
	}
	for name, prev := range prevResourceMaps.DashboardTables {
		if current, ok := resourceMaps.DashboardTables[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedTables = append(event.ChangedTables, diff)
			}
		} else {
			event.DeletedTables = append(event.DeletedTables, prev)
		}
	}
	for name, prev := range prevResourceMaps.DashboardCategories {
		if current, ok := resourceMaps.DashboardCategories[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedCategories = append(event.ChangedCategories, diff)
			}
		} else {
			event.DeletedCategories = append(event.DeletedCategories, prev)
		}
	}
	for name, prev := range prevResourceMaps.DashboardTexts {
		if current, ok := resourceMaps.DashboardTexts[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedTexts = append(event.ChangedTexts, diff)
			}
		} else {
			event.DeletedTexts = append(event.DeletedTexts, prev)
		}
	}

	// now detect new resources
	for name, p := range resourceMaps.Dashboards {
		if _, ok := prevResourceMaps.Dashboards[name]; !ok {
			event.NewDashboards = append(event.NewDashboards, p)
		}
	}
	for name, p := range resourceMaps.DashboardContainers {
		if _, ok := prevResourceMaps.DashboardContainers[name]; !ok {
			event.NewContainers = append(event.NewContainers, p)
		}
	}
	for name, p := range resourceMaps.DashboardCards {
		if _, ok := prevResourceMaps.DashboardCards[name]; !ok {
			event.NewCards = append(event.NewCards, p)
		}
	}
	for name, p := range resourceMaps.DashboardCategories {
		if _, ok := prevResourceMaps.DashboardCategories[name]; !ok {
			event.NewCategories = append(event.NewCategories, p)
		}
	}
	for name, p := range resourceMaps.DashboardCharts {
		if _, ok := prevResourceMaps.DashboardCharts[name]; !ok {
			event.NewCharts = append(event.NewCharts, p)
		}
	}
	for name, p := range resourceMaps.Benchmarks {
		if _, ok := prevResourceMaps.Benchmarks[name]; !ok {
			event.NewBenchmarks = append(event.NewBenchmarks, p)
		}
	}
	for name, p := range resourceMaps.Controls {
		if _, ok := prevResourceMaps.Controls[name]; !ok {
			event.NewControls = append(event.NewControls, p)
		}
	}
	for name, p := range resourceMaps.DashboardFlows {
		if _, ok := prevResourceMaps.DashboardFlows[name]; !ok {
			event.NewFlows = append(event.NewFlows, p)
		}
	}
	for name, p := range resourceMaps.DashboardGraphs {
		if _, ok := prevResourceMaps.DashboardGraphs[name]; !ok {
			event.NewGraphs = append(event.NewGraphs, p)
		}
	}
	for name, p := range resourceMaps.DashboardHierarchies {
		if _, ok := prevResourceMaps.DashboardHierarchies[name]; !ok {
			event.NewHierarchies = append(event.NewHierarchies, p)
		}
	}
	for name, p := range resourceMaps.DashboardImages {
		if _, ok := prevResourceMaps.DashboardImages[name]; !ok {
			event.NewImages = append(event.NewImages, p)
		}
	}
	for name, p := range resourceMaps.DashboardNodes {
		if _, ok := prevResourceMaps.DashboardNodes[name]; !ok {
			event.NewNodes = append(event.NewNodes, p)
		}
	}
	for name, p := range resourceMaps.DashboardEdges {
		if _, ok := prevResourceMaps.DashboardEdges[name]; !ok {
			event.NewEdges = append(event.NewEdges, p)
		}
	}
	for name, p := range resourceMaps.GlobalDashboardInputs {
		if _, ok := prevResourceMaps.GlobalDashboardInputs[name]; !ok {
			event.NewInputs = append(event.NewInputs, p)
		}
	}

	for name, currentInputsForDashboard := range resourceMaps.DashboardInputs {
		if prevInputsForDashboard, ok := prevResourceMaps.DashboardInputs[name]; ok {
			for name, current := range currentInputsForDashboard {
				if _, ok := prevInputsForDashboard[name]; !ok {
					event.NewInputs = append(event.NewInputs, current)
				}
			}
		} else {
			// all new
			for _, current := range currentInputsForDashboard {
				event.NewInputs = append(event.NewInputs, current)
			}
		}
	}

	for name, p := range resourceMaps.DashboardTables {
		if _, ok := prevResourceMaps.DashboardTables[name]; !ok {
			event.NewTables = append(event.NewTables, p)
		}
	}
	for name, p := range resourceMaps.DashboardTexts {
		if _, ok := prevResourceMaps.DashboardTexts[name]; !ok {
			event.NewTexts = append(event.NewTexts, p)
		}
	}

	if event.HasChanges() {
		// for every changed resource, set parents as changed, up the tree
		f := func(item modconfig.ModTreeItem) (bool, error) {
			event.SetParentsChanged(item, prevResourceMaps)
			return true, nil
		}
		err := event.WalkChangedResources(f)
		if err != nil {
			slog.Error("error walking changed resources", "error", err)
		}
		w.PublishDashboardEvent(ctx, event)
	}
}
