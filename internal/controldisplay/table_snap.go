package controldisplay

import (
	"github.com/turbot/powerpipe/internal/dashboardexecute"
	"strings"

	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/constants"
)

type TableRenderer_SNAP struct {
	resultTree *dashboardexecute.DisplayExecutionTree_SNAP

	// screen width
	width             int
	maxFailedControls int
	maxTotalControls  int
}

func NewTableSnapRenderer(resultTree *dashboardexecute.DisplayExecutionTree_SNAP) *TableRenderer_SNAP {
	return &TableRenderer_SNAP{
		resultTree:        resultTree,
		maxFailedControls: resultTree.Root.Summary.Status.FailedCount(),
		maxTotalControls:  resultTree.Root.Summary.Status.TotalCount(),
	}
}

// MinimumWidth is the width we require
// It is determined by the left indent, title, severity, counter and counter graph
func (r TableRenderer_SNAP) MinimumWidth() int {
	minimumWidthRequired := r.maxIndent() + minimumGroupTitleWidth + severityMaxLen + minimumCounterWidth + counterGraphSegments
	return minimumWidthRequired
}

func (r TableRenderer_SNAP) maxIndent() int {
	depth := r.groupDepth(r.resultTree.Root, 0)
	// each indent level is "| " or "+ " (2 characters)
	return depth * 2
}

func (r TableRenderer_SNAP) groupDepth(g *dashboardexecute.ResultGroup_SNAP, myDepth int) int {
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

func (r TableRenderer_SNAP) Render(width int) string {
	r.width = width

	// the buffer to put the output data in
	builder := strings.Builder{}

	builder.WriteString(r.renderResult())
	builder.WriteString("\n")
	builder.WriteString(r.renderSummary())

	return builder.String()
}

func (r TableRenderer_SNAP) renderSummary() string {
	// no need to render the summary when the dry-run flag is set
	if viper.GetBool(constants.ArgDryRun) {
		return ""
	}
	return NewSummaryRenderer(r.resultTree, r.width).Render()
}

func (r TableRenderer_SNAP) renderResult() string {
	return NewGroupRenderer(r.resultTree.Root, nil, r.maxFailedControls, r.maxTotalControls, r.resultTree, r.width).Render()
}
