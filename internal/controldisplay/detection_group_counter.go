package controldisplay

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type DetectionCounterRendererOptions struct {
	AddLeadingSpace bool
}

type DetectionCounterRenderer struct {
	count           int
	addLeadingSpace bool
}

func DetectionCounterRendererMinWidth() int {
	return 8
}

func NewDetectionCounterRenderer(count int, options DetectionCounterRendererOptions) *DetectionCounterRenderer {
	return &DetectionCounterRenderer{
		count:           count,
		addLeadingSpace: options.AddLeadingSpace,
	}
}

func (r DetectionCounterRenderer) Render() string {
	p := message.NewPrinter(language.English)

	// TODO K fixed width, commas etc.
	// get strings for fails and total - format with commas for thousands
	countString := p.Sprintf("%d", r.count)

	return countString
}
