package controldisplay

import (
	"context"
	"github.com/turbot/powerpipe/internal/controlexecute"
	"github.com/turbot/powerpipe/internal/dashboardexecute"
	"io"
)

type Formatter interface {
	FormatDetection(ctx context.Context, tree *dashboardexecute.DetectionBenchmarkDisplayTree) (io.Reader, error)
	Format(ctx context.Context, tree *controlexecute.ExecutionTree) (io.Reader, error)
	FileExtension() string
	Name() string
	Alias() string
}

// RowStreamer is implemented by formatters that read the rows of a control run one at a time
// (see ControlRun.EachRow), so that the rows can stay on disk
type RowStreamer interface {
	StreamsRows()
}

// StreamsRows returns whether the formatter reads the rows of a control run one at a time
func StreamsRows(f Formatter) bool {
	_, ok := f.(RowStreamer)
	return ok
}

// FormatTree formats the tree. Unless the formatter reads the rows one at a time, it first loads them into memory.
func FormatTree(ctx context.Context, f Formatter, tree *controlexecute.ExecutionTree) (io.Reader, error) {
	if !StreamsRows(f) {
		if err := tree.LoadRows(); err != nil {
			return nil, err
		}
	}
	return f.Format(ctx, tree)
}

type FormatterBase struct{}

func (*FormatterBase) Alias() string {
	return ""
}
