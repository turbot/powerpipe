package display

import (
	"github.com/turbot/pipe-fittings/printers"
	"strings"
)

type PrintableHclResource[T printers.Listable] struct {
	Items []T
}

func NewPrintableHclResource[T printers.Listable](items []T) *PrintableHclResource[T] {
	return &PrintableHclResource[T]{
		Items: items,
	}
}

func (p PrintableHclResource[T]) GetItems() []T {
	return p.Items
}

func (p PrintableHclResource[T]) GetTable() (*printers.Table, error) {
	var rows []printers.TableRow
	for _, item := range p.Items {
		row := item.GetListData()
		cleanRow(*row)
		rows = append(rows, *row)

	}
	if len(rows) == 0 {
		return printers.NewTable(), nil
	}
	columns := rows[0].Columns
	t := printers.NewTable().WithData(rows, columns)
	return t, nil
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
