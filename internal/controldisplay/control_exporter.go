package controldisplay

import (
	"context"
	"fmt"
	"io"

	"github.com/turbot/pipe-fittings/v2/contexthelpers"
	"github.com/turbot/pipe-fittings/v2/export"
	"github.com/turbot/powerpipe/internal/controlexecute"
	"github.com/turbot/powerpipe/internal/dashboardexecute"
)

var contextKeyFormatterPurpose = contexthelpers.ContextKey("formatter_purpose")

const formatterPurposeExport = "export"

type ControlExporter struct {
	formatter Formatter
}

func NewControlExporter(formatter Formatter) *ControlExporter {
	return &ControlExporter{formatter}
}

// exportContext tells the formatter it is being used for export
// this is a tactical mechanism used to ensure that exported snapshots are unindented
// whereas display snapshots are indented
func exportContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeyFormatterPurpose, formatterPurposeExport)
}

func (e *ControlExporter) Export(ctx context.Context, input export.ExportSourceData, destPath string) error {
	exportCtx := exportContext(ctx)

	var reader io.Reader
	var err error
	switch t := input.(type) {
	case *dashboardexecute.DetectionBenchmarkDisplayTree:
		reader, err = e.formatter.FormatDetection(exportCtx, t)
	case *controlexecute.ExecutionTree:
		reader, err = FormatTree(exportCtx, e.formatter, t)
	default:
		return fmt.Errorf("ControlExporter input must be ExecutionTree or DetectionBenchmarkDisplayTree")
	}
	if err != nil {
		return err
	}

	return export.Write(destPath, reader)
}

// ExportSnapshot writes the tree to w as a snapshot export does (--export=<file>.pps)
func ExportSnapshot(ctx context.Context, tree *controlexecute.ExecutionTree, w io.Writer) error {
	reader, err := FormatTree(exportContext(ctx), &SnapshotFormatter{}, tree)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, reader)
	return err
}

func (e *ControlExporter) FileExtension() string {
	return e.formatter.FileExtension()
}

func (e *ControlExporter) Name() string {
	return e.formatter.Name()
}

func (e *ControlExporter) Alias() string {
	return e.formatter.Alias()
}
