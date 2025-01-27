package controldisplay

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/powerpipe/internal/dashboardexecute"
)

type DetectionSummaryRenderer struct {
	resultTree *dashboardexecute.DetectionBenchmarkDisplayTree
	width      int
}

func NewDetectionSummaryRenderer(resultTree *dashboardexecute.DetectionBenchmarkDisplayTree, width int) *DetectionSummaryRenderer {
	return &DetectionSummaryRenderer{
		resultTree: resultTree,
		width:      width,
	}
}

func (r DetectionSummaryRenderer) Render() string {
	// TODO: get count from a produced summary
	count := getRowCount(r.resultTree.Root, 0)
	//count := r.resultTree.Root.Summary.Count

	// use alarm colour
	txtColorFunction := ControlColors.StatusColors["alarm"]
	countString := r.getPrintableNumber(count, txtColorFunction)

	statusStr := fmt.Sprintf("%s ", txtColorFunction("COUNT"))
	spaceAvailableForSpacer := r.width - (helpers.PrintableLength(statusStr) + helpers.PrintableLength(countString))
	spacer := NewSpacerRenderer(spaceAvailableForSpacer)

	return fmt.Sprintf(
		"%s%s%s",
		statusStr,
		spacer.Render(),
		countString,
	)
}

func (r DetectionSummaryRenderer) getPrintableNumber(number int, cf colorFunc) string {
	p := message.NewPrinter(language.English)
	s := p.Sprintf("%d", number)
	return fmt.Sprintf("%s ", cf(s))
}

func getRowCount(group *dashboardexecute.DetectionBenchmarkDisplay, count int) int {
	for _, subGroup := range group.Groups {
		count += getRowCount(subGroup, count)
	}
	for _, run := range group.DetectionRuns {
		count += len(run.Data.Rows)
	}
	return count
}
