package backend

import (
	"context"
	"database/sql"
	"strings"
)

type DuckDBBackend struct {
	originalConnectionString string
	rowreader                RowReader
}

// Connect implements Backend.
func (s *DuckDBBackend) Connect(context.Context, ...ConnectOption) (*sql.DB, error) {
	connString := s.originalConnectionString
	connString = strings.TrimSpace(connString) // remove any leading or trailing whitespace
	connString = strings.TrimPrefix(connString, "duckdb://")
	return sql.Open("sqlite3", connString)
}

// GetType implements Backend.
func (s *DuckDBBackend) GetType() BackendType {
	return SqliteDBClientBackend
}

// RowReader implements Backend.
func (s *DuckDBBackend) RowReader() RowReader {
	return s.rowreader
}

func NewDuckDBBackend(ctx context.Context, connString string) Backend {
	return &DuckDBBackend{
		rowreader: NewDuckDBRowReader(),
	}
}

type duckdbRowReader struct {
	GenericRowReader
}

func NewDuckDBRowReader() *duckdbRowReader {
	return &duckdbRowReader{
		// use the generic row reader - to start with
		GenericRowReader: *NewPassThruRowReader(),
	}
}
