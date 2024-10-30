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

	modResources := r.(*resources.PowerpipeModResources)
	prevModResources := p.(*resources.PowerpipeModResources)

	// TODO reports can we use a PowerpipeModResources diff function to do all of this - we are duplicating logic

	// first detect changes to existing resources and deletions
	for name, prev := range prevModResources.Dashboards {
		if current, ok := modResources.Dashboards[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedDashboards = append(event.ChangedDashboards, diff)
			}
		} else {
			event.DeletedDashboards = append(event.DeletedDashboards, prev)
		}
	}
	for name, prev := range prevModResources.DashboardContainers {
		if current, ok := modResources.DashboardContainers[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedContainers = append(event.ChangedContainers, diff)
			}
		} else {
			event.DeletedContainers = append(event.DeletedContainers, prev)
		}
	}
	for name, prev := range prevModResources.DashboardCards {
		if current, ok := modResources.DashboardCards[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedCards = append(event.ChangedCards, diff)
			}
		} else {
			event.DeletedCards = append(event.DeletedCards, prev)
		}
	}
	for name, prev := range prevModResources.DashboardCharts {
		if current, ok := modResources.DashboardCharts[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedCharts = append(event.ChangedCharts, diff)
			}
		} else {
			event.DeletedCharts = append(event.DeletedCharts, prev)
		}
	}
	for name, prev := range prevModResources.Benchmarks {
		if current, ok := modResources.Benchmarks[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedBenchmarks = append(event.ChangedBenchmarks, diff)
			}
		} else {
			event.DeletedBenchmarks = append(event.DeletedBenchmarks, prev)
		}
	}
	for name, prev := range prevModResources.Controls {
		if current, ok := modResources.Controls[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedControls = append(event.ChangedControls, diff)
			}
		} else {
			event.DeletedControls = append(event.DeletedControls, prev)
		}
	}
	for name, prev := range prevModResources.DashboardFlows {
		if current, ok := modResources.DashboardFlows[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedFlows = append(event.ChangedFlows, diff)
			}
		} else {
			event.DeletedFlows = append(event.DeletedFlows, prev)
		}
	}
	for name, prev := range prevModResources.DashboardGraphs {
		if current, ok := modResources.DashboardGraphs[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedGraphs = append(event.ChangedGraphs, diff)
			}
		} else {
			event.DeletedGraphs = append(event.DeletedGraphs, prev)
		}
	}
	for name, prev := range prevModResources.DashboardHierarchies {
		if current, ok := modResources.DashboardHierarchies[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedHierarchies = append(event.ChangedHierarchies, diff)
			}
		} else {
			event.DeletedHierarchies = append(event.DeletedHierarchies, prev)
		}
	}
	for name, prev := range prevModResources.DashboardImages {
		if current, ok := modResources.DashboardImages[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedImages = append(event.ChangedImages, diff)
			}
		} else {
			event.DeletedImages = append(event.DeletedImages, prev)
		}
	}
	for name, prev := range prevModResources.DashboardNodes {
		if current, ok := modResources.DashboardNodes[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedNodes = append(event.ChangedNodes, diff)
			}
		} else {
			event.DeletedNodes = append(event.DeletedNodes, prev)
		}
	}
	for name, prev := range prevModResources.DashboardEdges {
		if current, ok := modResources.DashboardEdges[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedEdges = append(event.ChangedEdges, diff)
			}
		} else {
			event.DeletedEdges = append(event.DeletedEdges, prev)
		}
	}
	for name, prev := range prevModResources.GlobalDashboardInputs {
		if current, ok := modResources.GlobalDashboardInputs[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedInputs = append(event.ChangedInputs, diff)
			}
		} else {
			event.DeletedInputs = append(event.DeletedInputs, prev)
		}
	}
	for name, prevInputsForDashboard := range prevModResources.DashboardInputs {
		if currentInputsForDashboard, ok := modResources.DashboardInputs[name]; ok {
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
	for name, prev := range prevModResources.DashboardTables {
		if current, ok := modResources.DashboardTables[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedTables = append(event.ChangedTables, diff)
			}
		} else {
			event.DeletedTables = append(event.DeletedTables, prev)
		}
	}
	for name, prev := range prevModResources.DashboardCategories {
		if current, ok := modResources.DashboardCategories[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedCategories = append(event.ChangedCategories, diff)
			}
		} else {
			event.DeletedCategories = append(event.DeletedCategories, prev)
		}
	}
	for name, prev := range prevModResources.DashboardTexts {
		if current, ok := modResources.DashboardTexts[name]; ok {
			diff := prev.Diff(current)
			if diff.HasChanges() {
				event.ChangedTexts = append(event.ChangedTexts, diff)
			}
		} else {
			event.DeletedTexts = append(event.DeletedTexts, prev)
		}
	}

	// now detect new resources
	for name, p := range modResources.Dashboards {
		if _, ok := prevModResources.Dashboards[name]; !ok {
			event.NewDashboards = append(event.NewDashboards, p)
		}
	}
	for name, p := range modResources.DashboardContainers {
		if _, ok := prevModResources.DashboardContainers[name]; !ok {
			event.NewContainers = append(event.NewContainers, p)
		}
	}
	for name, p := range modResources.DashboardCards {
		if _, ok := prevModResources.DashboardCards[name]; !ok {
			event.NewCards = append(event.NewCards, p)
		}
	}
	for name, p := range modResources.DashboardCategories {
		if _, ok := prevModResources.DashboardCategories[name]; !ok {
			event.NewCategories = append(event.NewCategories, p)
		}
	}
	for name, p := range modResources.DashboardCharts {
		if _, ok := prevModResources.DashboardCharts[name]; !ok {
			event.NewCharts = append(event.NewCharts, p)
		}
	}
	for name, p := range modResources.Benchmarks {
		if _, ok := prevModResources.Benchmarks[name]; !ok {
			event.NewBenchmarks = append(event.NewBenchmarks, p)
		}
	}
	for name, p := range modResources.Controls {
		if _, ok := prevModResources.Controls[name]; !ok {
			event.NewControls = append(event.NewControls, p)
		}
	}
	for name, p := range modResources.DashboardFlows {
		if _, ok := prevModResources.DashboardFlows[name]; !ok {
			event.NewFlows = append(event.NewFlows, p)
		}
	}
	for name, p := range modResources.DashboardGraphs {
		if _, ok := prevModResources.DashboardGraphs[name]; !ok {
			event.NewGraphs = append(event.NewGraphs, p)
		}
	}
	for name, p := range modResources.DashboardHierarchies {
		if _, ok := prevModResources.DashboardHierarchies[name]; !ok {
			event.NewHierarchies = append(event.NewHierarchies, p)
		}
	}
	for name, p := range modResources.DashboardImages {
		if _, ok := prevModResources.DashboardImages[name]; !ok {
			event.NewImages = append(event.NewImages, p)
		}
	}
	for name, p := range modResources.DashboardNodes {
		if _, ok := prevModResources.DashboardNodes[name]; !ok {
			event.NewNodes = append(event.NewNodes, p)
		}
	}
	for name, p := range modResources.DashboardEdges {
		if _, ok := prevModResources.DashboardEdges[name]; !ok {
			event.NewEdges = append(event.NewEdges, p)
		}
	}
	for name, p := range modResources.GlobalDashboardInputs {
		if _, ok := prevModResources.GlobalDashboardInputs[name]; !ok {
			event.NewInputs = append(event.NewInputs, p)
		}
	}

	for name, currentInputsForDashboard := range modResources.DashboardInputs {
		if prevInputsForDashboard, ok := prevModResources.DashboardInputs[name]; ok {
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

	for name, p := range modResources.DashboardTables {
		if _, ok := prevModResources.DashboardTables[name]; !ok {
			event.NewTables = append(event.NewTables, p)
		}
	}
	for name, p := range modResources.DashboardTexts {
		if _, ok := prevModResources.DashboardTexts[name]; !ok {
			event.NewTexts = append(event.NewTexts, p)
		}
	}

	if event.HasChanges() {
		// for every changed resource, set parents as changed, up the tree
		f := func(item modconfig.ModTreeItem) (bool, error) {
			event.SetParentsChanged(item, prevModResources)
			return true, nil
		}
		err := event.WalkChangedResources(f)
		if err != nil {
			slog.Error("error walking changed resources", "error", err)
		}
		w.PublishDashboardEvent(ctx, event)
	}
}
