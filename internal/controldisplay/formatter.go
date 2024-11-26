package controldisplay

import (
	"context"
	"github.com/turbot/powerpipe/internal/dashboardexecute"
	"io"
)

type Formatter interface {
	Format(ctx context.Context, tree *dashboardexecute.DetectionBenchmarkDisplayTree) (io.Reader, error)
	FileExtension() string
	Name() string
	Alias() string
}

type FormatterBase struct{}

func (*FormatterBase) Alias() string {
	return ""
}
