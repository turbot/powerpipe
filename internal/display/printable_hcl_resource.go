package display

import (
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/printers"
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
	for _, item := range p.Items {
		rwm, _ := any(item).(modconfig.ResourceWithMetadata)

		cells := []any{
			rwm.GetMetadata().ModName,
			item.Name(),
			item.GetDescription(),
		}
		tableRows = append(tableRows, printers.TableRow{Cells: cells})
	}

	return printers.NewTable(tableRows, p.getColumns()), nil
}

func (PrintableHclResource[T]) getColumns() (columns []printers.TableColumnDefinition) {
	return []printers.TableColumnDefinition{
		{
			Name:        "MOD",
			Type:        "string",
			Description: "Mod name",
		},
		{
			Name:        "NAME",
			Type:        "string",
			Description: "Resource name",
		},
		{
			Name:        "DESCRIPTION",
			Type:        "string",
			Description: "Resource description",
		},
	}
}
