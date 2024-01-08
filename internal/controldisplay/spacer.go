package controldisplay

import (
	"fmt"
	"log/slog"
	"strings"
)

type SpacerRenderer struct {
	width int
}

func NewSpacerRenderer(width int) *SpacerRenderer {
	return &SpacerRenderer{width}
}

// Render returns a divider string of format: "....... "
// NOTE: adds a trailing space
func (r SpacerRenderer) Render() string {
	if r.width <= 0 {
		// this should never happen, since the minimum width is set by the formatter
		slog.Warn("SpacerRenderer.Render unexpected negative width", "width", r.width)
		return ""
	}
	// we always have a trailing space
	if r.width == 1 {
		return " "
	}

	// allow for trailing space
	numberOfDots := r.width - 1
	return fmt.Sprintf("%s ", ControlColors.Spacer(strings.Repeat(".", numberOfDots)))
}
