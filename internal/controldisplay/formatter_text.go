package controldisplay

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/powerpipe/internal/controlexecute"
	"github.com/turbot/powerpipe/internal/dashboardexecute"
)

const MaxColumns = 200

type TextFormatter struct {
	FormatterBase
}

func (tf TextFormatter) FormatDetection(ctx context.Context, tree *dashboardexecute.DetectionBenchmarkDisplayTree) (io.Reader, error) {
	renderer := NewDetectionTableRenderer(tree)
	widthConstraint := utils.NewRangeConstraint(renderer.MinimumWidth(), MaxColumns)
	renderedText := renderer.Render(widthConstraint.Constrain(GetMaxCols()))
	res := strings.NewReader(fmt.Sprintf("\n%s\n", renderedText))
	return res, nil
}

func (tf TextFormatter) Format(ctx context.Context, tree *controlexecute.ExecutionTree) (io.Reader, error) {
	renderer := NewTableRenderer(tree)
	widthConstraint := utils.NewRangeConstraint(renderer.MinimumWidth(), MaxColumns)
	renderedText := renderer.Render(widthConstraint.Constrain(GetMaxCols()))
	res := strings.NewReader(fmt.Sprintf("\n%s\n", renderedText))
	return res, nil
}

func (tf TextFormatter) FileExtension() string {
	return constants.TextExtension
}

func (tf TextFormatter) Name() string {
	return constants.OutputFormatText
}

func (tf TextFormatter) Alias() string {
	return constants.OutputFormatBrief
}
