package backend

import (
	"github.com/turbot/pipe-fittings/queryresult"
)

func NewPassThruRowReader() *GenericRowReader {
	return &GenericRowReader{
		CellReader: func(columnValue any, col *queryresult.ColumnDef) (any, error) {
			return columnValue, nil
		},
	}
}

// GenericRowReader is a RowReader implementation for generic database/sql driver
type GenericRowReader struct {
	CellReader func(columnValue any, col *queryresult.ColumnDef) (any, error)
}

func (r *GenericRowReader) Read(columnValues []any, cols []*queryresult.ColumnDef) ([]any, error) {
	result := make([]any, len(columnValues))
	for i, columnValue := range columnValues {
		if cellValue, err := r.CellReader(columnValue, cols[i]); err == nil {
			result[i] = cellValue
		}
	}
	return result, nil
}
