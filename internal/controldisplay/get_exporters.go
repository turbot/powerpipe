package controldisplay

import (
	"github.com/turbot/pipe-fittings/v2/export"
	"github.com/turbot/pipe-fittings/v2/modconfig"
)

// GetExporters returns an array of ControlExporters corresponding to the available output formats
func GetExporters(target modconfig.ModTreeItem) ([]export.Exporter, error) {
	formatResolver, err := NewFormatResolver(target)
	if err != nil {
		return nil, err
	}
	exporters := formatResolver.controlExporters()
	return exporters, nil
}
