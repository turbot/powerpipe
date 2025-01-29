package controldisplay

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/v2/constants"
)

type DetectionGroupHeadingRenderer struct {
	title    string
	severity string
	// screen width
	width  int
	indent string
	count  int
}

func NewDetectionGroupHeadingRenderer(title string, count, width int, indent string) *DetectionGroupHeadingRenderer {
	return &DetectionGroupHeadingRenderer{
		title:  title,
		count:  count,
		width:  width,
		indent: indent,
	}
}

func (r DetectionGroupHeadingRenderer) Render() string {
	isDryRun := viper.GetBool(constants.ArgDryRun)

	if r.width <= 0 {
		// this should never happen, since the minimum width is set by the formatter
		slog.Warn("DetectionGroupHeadingRenderer.Render unexpected negative width", "width", r.width)
		return ""
	}

	formattedIndent := fmt.Sprintf("%s", ControlColors.Indent(r.indent))
	indentWidth := helpers.PrintableLength(formattedIndent)

	// for a dry run we do not display the counters or graph
	var severityString, counterString string
	if !isDryRun {
		// only display severity if there are any results
		if r.count > 0 {
			severityString = NewSeverityRenderer(r.severity).Render()
		}
		counterString = NewDetectionCounterRenderer(
			r.count,
			DetectionCounterRendererOptions{
				AddLeadingSpace: true,
			},
		).Render()
	}
	severityWidth := helpers.PrintableLength(severityString)
	counterWidth := helpers.PrintableLength(counterString)

	// figure out how much width we have available for the title
	availableWidth := r.width - counterWidth - severityWidth - indentWidth

	// now availableWidth is all we have - if it is not enough we need to truncate the title
	titleString := NewGroupTitleRenderer(r.title, availableWidth).Render()
	titleWidth := helpers.PrintableLength(titleString)

	// is there any room for a spacer
	spacerWidth := availableWidth - titleWidth
	var spacerString string
	if spacerWidth > 0 && !isDryRun {
		spacerString = NewSpacerRenderer(spacerWidth).Render()
	}

	// now put these all together
	str := fmt.Sprintf("%s%s%s%s%s", formattedIndent, titleString, spacerString, severityString, counterString)
	return str
}
