package controldisplay

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/powerpipe/internal/controlexecute"
)

type DimensionsRenderer struct {
	dimensions     []controlexecute.Dimension
	colorGenerator *controlexecute.DimensionColorGenerator
	width          int
}

func NewDimensionsRenderer(dimensions []controlexecute.Dimension, colorGenerator *controlexecute.DimensionColorGenerator, width int) *DimensionsRenderer {
	return &DimensionsRenderer{
		dimensions:     dimensions,
		colorGenerator: colorGenerator,
		width:          width,
	}
}

// Render returns the dimensions, truncated to the max length if necessary
func (r DimensionsRenderer) Render() string {
	if r.width <= 0 {
		// this should never happen, since the minimum width is set by the formatter
		slog.Warn("DimensionsRenderer.Render unexpected negative width", "width", r.width)
		return ""
	}
	if len(r.dimensions) == 0 {
		return ""
	}
	// make array of dimension values (including trailing spaces
	var formattedDimensions = make([]string, len(r.dimensions))
	for i, d := range r.dimensions {
		// a result row is printed on a single line, so render line breaks in the value as spaces
		formattedDimensions[i] = singleLine(d.Value)
	}

	var length int
	for length = dimensionsLength(formattedDimensions); length > r.width; {
		// truncate the first dimension
		if helpers.PrintableLength(formattedDimensions[0]) > 0 {
			// truncate the original value, not the already truncated value
			newLength := helpers.PrintableLength(formattedDimensions[0]) - 1
			truncated := helpers.TruncateString(formattedDimensions[0], newLength)
			if truncated == formattedDimensions[0] {
				// truncation made no progress - clear the value rather than loop forever
				truncated = ""
			}
			formattedDimensions[0] = truncated
		} else {
			// so event with all dimensions 1 long, we still do not have enough space
			// remove a dimension from the array
			if len(formattedDimensions) > 2 {
				r.dimensions = r.dimensions[1:]
				formattedDimensions = formattedDimensions[1:]
			} else {
				// there is only 1 dimension - nothing we can do here, give up
				return ""
			}
		}
		// update length
		length = dimensionsLength(formattedDimensions)
	}

	// ok we now have dimensions that fit in the space, color them
	// check whether color is disabled

	coloredDimensions := make([]string, 0, len(r.dimensions))
	for i, v := range formattedDimensions {
		// get the source dimension object
		dimension := r.dimensions[i]

		if len(strings.TrimSpace(dimension.Value)) == 0 {
			// if the value of the dimension is empty, skip it
			continue
		}

		// get the color code - there must be an entry
		dimensionColorFunc := func(val interface{}) aurora.Value {
			// if current theme supports colors, apply coloring
			if ControlColors.UseColor {
				dimensionColor := r.colorGenerator.Map[dimension.Key][dimension.Value]
				return aurora.Index(dimensionColor, val)
			}
			return aurora.Reset(val)
		}

		coloredDimensions = append(coloredDimensions, fmt.Sprintf("%s", dimensionColorFunc(v)))
	}

	return strings.Join(coloredDimensions, " ")
}

// count the total length of the dimensions
func dimensionsLength(dimensionValues []string) int {
	var res int
	for _, v := range dimensionValues {
		res += len(v)
	}
	// allow for spaces between the dimensions
	res += len(dimensionValues) - 1
	return res
}

var lineBreakReplacer = strings.NewReplacer("\r\n", " ", "\n", " ", "\r", " ")

// singleLine replaces line breaks with spaces
func singleLine(value string) string {
	return lineBreakReplacer.Replace(value)
}
