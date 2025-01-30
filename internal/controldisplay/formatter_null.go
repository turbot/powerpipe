package controldisplay

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/turbot/pipe-fittings/v2/constants"
	"github.com/turbot/powerpipe/internal/controlexecute"
	"github.com/turbot/powerpipe/internal/dashboardexecute"
)

// NullFormatter is to be used when no output is expected. It always returns a `io.Reader` which
// reads an empty string
type NullFormatter struct {
	FormatterBase
}

func (*NullFormatter) FormatDetection(context.Context, *dashboardexecute.DetectionBenchmarkDisplayTree) (io.Reader, error) {
	return nil, fmt.Errorf("NullFormatter does not support FormatDetection")
}

func (j *NullFormatter) Format(ctx context.Context, tree *controlexecute.ExecutionTree) (io.Reader, error) {
	return strings.NewReader(""), nil
}

func (j *NullFormatter) FileExtension() string {
	// will not be called
	return ""
}

func (j *NullFormatter) Name() string {
	return constants.OutputFormatNone
}
