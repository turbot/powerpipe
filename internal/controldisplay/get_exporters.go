package controldisplay

import (
	"github.com/turbot/pipe-fittings/v2/export"
)

// GetExporters returns an array of ControlExporters corresponding to the available output formats
func GetExporters() ([]export.Exporter, error) {
	formatResolver, err := NewFormatResolver()
	if err != nil {
		return nil, err
	}
	exporters := formatResolver.controlExporters()
	return exporters, nil
}
