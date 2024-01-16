package backend

import (
	"github.com/turbot/pipe-fittings/queryresult"
)

func NewPassThruRowReader() *PassThruRowReader {
	return &PassThruRowReader{
		CellReader: func(columnValue any, col *queryresult.ColumnDef) (any, error) {
			return columnValue, nil
		},
	}
}

// PassThruRowReader is a RowReader implementation for generic database/sql driver
type PassThruRowReader struct {
	CellReader func(columnValue any, col *queryresult.ColumnDef) (any, error)
}

func (r *PassThruRowReader) Read(columnValues []any, cols []*queryresult.ColumnDef) ([]any, error) {
	result := make([]any, len(columnValues))
	for i, columnValue := range columnValues {
		if cellValue, err := r.CellReader(columnValue, cols[i]); err == nil {
			result[i] = cellValue
		}
	}
	return result, nil
}
