package controldisplay

import (
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/powerpipe/internal/dashboardexecute"
)

type SummarySeverityRenderer struct {
	resultTree *dashboardexecute.DisplayExecutionTree_SNAP
	width      int
}

func NewSummarySeverityRenderer(resultTree *dashboardexecute.DisplayExecutionTree_SNAP, width int) *SummarySeverityRenderer {
	return &SummarySeverityRenderer{
		resultTree: resultTree,
		width:      width,
	}
}

func (r *SummarySeverityRenderer) Render() []string {
	availableWidth := r.width

	// render the critical line
	criticalSeverityRow := NewSummarySeverityRowRenderer(r.resultTree, availableWidth, "critical").Render()
	criticalWidth := helpers.PrintableLength(criticalSeverityRow)
	// if there is a critical line, use this to set the max width
	if criticalWidth > 0 {
		availableWidth = criticalWidth
	}

	// render the high line
	highSeverityRow := NewSummarySeverityRowRenderer(r.resultTree, availableWidth, "high").Render()
	highWidth := helpers.PrintableLength(highSeverityRow)

	// build the severity block
	var strs []string
	if criticalWidth > 0 {
		strs = append(strs, criticalSeverityRow)
	}
	if highWidth > 0 {
		strs = append(strs, highSeverityRow)
	}
	return strs
}
