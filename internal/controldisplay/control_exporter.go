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

func (e *ControlExporter) Export(ctx context.Context, input export.ExportSourceData, destPath string) error {

	// tell the formatter it is being used for export
	// this is a tactical mechanism used to ensure that exported snapshots are unindented
	// whereas display snapshots are indented
	exportCtx := context.WithValue(ctx, contextKeyFormatterPurpose, formatterPurposeExport)

	var reader io.Reader
	var err error
	switch t := input.(type) {
	case *dashboardexecute.DetectionBenchmarkDisplayTree:
		reader, err = e.formatter.FormatDetection(exportCtx, t)
	case *controlexecute.ExecutionTree:
		reader, err = e.formatter.Format(exportCtx, t)
	default:
		return fmt.Errorf("ControlExporter input must be ExecutionTree or DetectionBenchmarkDisplayTree")
	}
	if err != nil {
		return err
	}

	return export.Write(destPath, reader)
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
