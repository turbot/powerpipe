package backend

import (
	"github.com/turbot/pipe-fittings/queryresult"
)

type RowReader interface {
	Read(columnValues []any, cols []*queryresult.ColumnDef) ([]any, error)
}

func RowReaderFactory(backend DBClientBackendType) (RowReader, error) {
	var reader RowReader
	switch backend {
	case PostgresDBClientBackend:
		// we have special handing of a few types for postgres
		reader = NewPgxRowReader()
	case MySQLDBClientBackend:
		reader = NewMySqlRowReader()
	case SqliteDBClientBackend:
		reader = NewGenericSQLRowReader()
	case DuckDBClientBackend:
		reader = NewDuckDBRowReader()
	default:
		return nil, ErrUnknownBackend
	}
	return reader, nil
}
