package backend

import (
	"github.com/turbot/pipe-fittings/queryresult"
)

func NewBasicRowReader() *BasicRowReader {
	return &BasicRowReader{
		CellReader: func(columnValue any, col *queryresult.ColumnDef) (any, error) {
			return columnValue, nil
		},
	}
}

// BasicRowReader is a RowReader implementation for generic database/sql driver
type BasicRowReader struct {
	CellReader func(columnValue any, col *queryresult.ColumnDef) (any, error)
}

func (r *BasicRowReader) Read(columnValues []any, cols []*queryresult.ColumnDef) ([]any, error) {
	result := make([]any, len(columnValues))
	for i, columnValue := range columnValues {
		if cellValue, err := r.CellReader(columnValue, cols[i]); err == nil {
			result[i] = cellValue
		}
	}
	return result, nil
}
