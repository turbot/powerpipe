package controldisplay

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/powerpipe/internal/controlexecute"
)

// TODO dimensions???

type DetectionResultRenderer struct {
	displayText    string
	rows           []map[string]any
	colorGenerator *controlexecute.DimensionColorGenerator

	// screen width
	width int
	// if true, only display failed results
	errorsOnly bool
	indent     string
	dimensions []controlexecute.Dimension
}

func NewDetectionResultRenderer(displayText string, dimensions []controlexecute.Dimension, colorGenerator *controlexecute.DimensionColorGenerator, width int, indent string) *DetectionResultRenderer {
	return &DetectionResultRenderer{
		displayText:    displayText,
		dimensions:     dimensions,
		colorGenerator: colorGenerator,
		width:          width,
		errorsOnly:     viper.GetString(constants.ArgOutput) == "brief",
		indent:         indent,
	}
}

func (r DetectionResultRenderer) Render() string {

	formattedIndent := fmt.Sprintf("%s", ControlColors.Indent(r.indent))
	indentWidth := helpers.PrintableLength(formattedIndent)

	// figure out how much width we have available for the  dimensions, allowing the minimum for the reason
	availableWidth := r.width - indentWidth

	// for now give this all to reason
	//availableDimensionWidth := availableWidth
	//var dimensionsString string
	//var dimensionWidth int
	//if availableDimensionWidth > 0 {
	//	dimensionsString = NewDimensionsRenderer(r.dimensions, r.colorGenerator, availableDimensionWidth).Render()
	//	dimensionWidth = helpers.PrintableLength(dimensionsString)
	//	availableWidth -= dimensionWidth
	//}

	availableWidth = availableWidth - len(r.displayText)

	// is there any room for a spacer
	// spacerString := NewSpacerRenderer(availableWidth).Render()

	// now put these all together
	str := fmt.Sprintf("%s%s", formattedIndent, r.displayText)
	return str
}
