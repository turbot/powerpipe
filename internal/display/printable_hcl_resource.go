package display

import (
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/printers"
	"strings"
)

type PrintableHclResource[T modconfig.HclResource] struct {
	Items []T
}

func NewPrintableHclResource[T modconfig.HclResource](items []T) *PrintableHclResource[T] {
	return &PrintableHclResource[T]{
		Items: items,
	}
}

func (p PrintableHclResource[T]) GetItems() []T {
	return p.Items
}

func (p PrintableHclResource[T]) GetTable() (printers.Table, error) {
	var tableRows []printers.TableRow
	var columnsDefs []printers.TableColumnDefinition
	for _, item := range p.Items {

		var row printers.TableRow
		// is this a query provider - get query provider columns
		if qp, ok := any(item).(modconfig.QueryProvider); ok {
			qpImpl := qp.GetQueryProviderImpl()
			qpRow, qpColumns := GetAsTableRow(qpImpl)
			row.Cells = append(row.Cells, qpRow.Cells...)
			// if this is the first item, set column defs
			if len(tableRows) == 0 {
				columnsDefs = append(columnsDefs, qpColumns...)
			}
		}

		// get hcl resource columns
		hrImpl := item.GetHclResourceImpl()
		hrRow, hrColumns := GetAsTableRow(hrImpl)
		row.Cells = append(row.Cells, hrRow.Cells...)
		if len(tableRows) == 0 {
			columnsDefs = append(columnsDefs, hrColumns...)
		}

		// now get item specific fields
		itemRow, itemColumns := GetAsTableRow(item)
		row.Cells = append(row.Cells, itemRow.Cells...)
		if len(tableRows) == 0 {
			columnsDefs = append(columnsDefs, itemColumns...)
		}

		// strip tabs and newlines
		cleanRow(row)

		//cleanRow := p.Sanitizer.SanitizeTableRow(row)
		tableRows = append(tableRows, row)
	}

	return printers.NewTable(tableRows, columnsDefs), nil
}

func cleanRow(row printers.TableRow) {
	var charsToRemove = []string{"\t", "\n", "\r"}
	for i, c := range row.Cells {
		str, ok := c.(string)
		if !ok {
			continue
		}

		for _, r := range charsToRemove {
			str = strings.ReplaceAll(str, r, "")
		}
		// TODO tactical column width to 100
		const maxWidth = 100
		if len(str) > maxWidth {
			str = str[:maxWidth] + "…"
		}
		row.Cells[i] = str
	}
}
