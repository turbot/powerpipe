package controldisplay

import (
	"github.com/turbot/powerpipe/internal/dashboardexecute"
	"strings"

	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/v2/constants"
)

type DetectionTableRenderer struct {
	resultTree *dashboardexecute.DetectionBenchmarkDisplayTree

	// screen width
	width int
}

func NewDetectionTableRenderer(resultTree *dashboardexecute.DetectionBenchmarkDisplayTree) *DetectionTableRenderer {
	return &DetectionTableRenderer{
		resultTree: resultTree,
	}
}

// MinimumWidth is the width we require
// It is determined by the left indent, title, severity, counter and counter graph
func (r DetectionTableRenderer) MinimumWidth() int {
	minimumWidthRequired := r.maxIndent() + minimumGroupTitleWidth + severityMaxLen + minimumCounterWidth + counterGraphSegments
	return minimumWidthRequired
}

func (r DetectionTableRenderer) maxIndent() int {
	depth := r.groupDepth(r.resultTree.Root, 0)
	// each indent level is "| " or "+ " (2 characters)
	return depth * 2
}

func (r DetectionTableRenderer) groupDepth(g *dashboardexecute.DetectionBenchmarkDisplay, myDepth int) int {
	if len(g.Groups) == 0 {
		return 0
	}
	maxDepth := 0
	for _, subGroup := range g.Groups {
		branchDepth := r.groupDepth(subGroup, myDepth+1)
		if branchDepth > maxDepth {
			maxDepth = branchDepth
		}
	}
	return myDepth + maxDepth
}

func (r DetectionTableRenderer) Render(width int) string {
	r.width = width

	// the buffer to put the output data in
	builder := strings.Builder{}

	builder.WriteString(r.renderResult())
	builder.WriteString("\n")
	builder.WriteString(r.renderSummary())

	return builder.String()
}

func (r DetectionTableRenderer) renderSummary() string {
	// no need to render the summary when the dry-run flag is set
	if viper.GetBool(constants.ArgDryRun) {
		return ""
	}
	return NewDetectionSummaryRenderer(r.resultTree, r.width).Render()
}

func (r DetectionTableRenderer) renderResult() string {
	return NewDetectionGroupRenderer(r.resultTree.Root, nil, r.resultTree, r.width).Render()
}
