package display

import (
	"github.com/turbot/pipe-fittings/modconfig"
	"slices"
	"strings"

	"github.com/turbot/pipe-fittings/printers"
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
	// split rows into top level mod resources and dependency mod resources
	// show the top level resources first

	var rows, depRows []printers.TableRow
	var columns []string
	for _, item := range p.Items {
		row := item.GetListData().GetRow()
		if len(columns) == 0 {
			columns = row.Columns
		}

		cleanRow(*row)

		if isDependencyResource(item) {
			depRows = append(depRows, *row)

		} else {
			rows = append(rows, *row)
		}
	}
	if len(rows)+len(depRows) == 0 {
		return printers.NewTable(), nil
	}

	// sort output based on column 0
	sortFunc := func(a, b printers.TableRow) int {
		return strings.Compare(a.Cells[0].(string), b.Cells[0].(string))
	}
	slices.SortFunc(rows, sortFunc)
	slices.SortFunc(depRows, sortFunc)

	t := printers.NewTable().WithData(append(rows, depRows...), columns)
	return t, nil
}

func isDependencyResource(item printers.Listable) bool {
	// is this a ModTreeItem - we expect it will be
	mti, ok := item.(modconfig.ModTreeItem)
	if !ok {
		return false
	}
	return mti.IsDependencyResource()
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
