package controldisplay

import (
	"context"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/powerpipe/internal/dashboardexecute"
	"io"
	"strings"
)

// NullFormatter is to be used when no output is expected. It always returns a `io.Reader` which
// reads an empty string
type NullFormatter struct {
	FormatterBase
}

func (j *NullFormatter) Format(ctx context.Context, tree *dashboardexecute.DisplayExecutionTree_SNAP) (io.Reader, error) {
	return strings.NewReader(""), nil
}

func (j *NullFormatter) FileExtension() string {
	// will not be called
	return ""
}

func (j *NullFormatter) Name() string {
	return constants.OutputFormatNone
}
