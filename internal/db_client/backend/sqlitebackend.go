package backend

import (
	"context"
	"database/sql"
	"strings"
)

type SqliteBackend struct {
	originalConnectionString string
	rowreader                RowReader
}

// Connect implements Backend.
func (s *SqliteBackend) Connect(context.Context, ...ConnectOption) (*sql.DB, error) {
	connString := s.originalConnectionString
	connString = strings.TrimSpace(connString) // remove any leading or trailing whitespace
	connString = strings.TrimPrefix(connString, "sqlite://")
	return sql.Open("sqlite3", connString)
}

// GetType implements Backend.
func (s *SqliteBackend) GetType() DBClientBackendType {
	return SqliteDBClientBackend
}

// RowReader implements Backend.
func (s *SqliteBackend) RowReader() RowReader {
	return s.rowreader
}

func NewSqliteBackend(ctx context.Context, connString string) Backend {
	return &SqliteBackend{
		rowreader: NewSqliteRowReader(),
	}
}

type sqliteRowReader struct {
	PassThruRowReader
}

func NewSqliteRowReader() *sqliteRowReader {
	return &sqliteRowReader{
		// use the generic row reader - there's no real difference between sqlite and generic
		PassThruRowReader: *NewPassThruRowReader(),
	}
}
