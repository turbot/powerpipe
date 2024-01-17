package backend

import (
	"context"
	"database/sql"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

type SqliteBackend struct {
	originalConnectionString string
	rowreader                RowReader
}

// Connect implements Backend.
func (s *SqliteBackend) Connect(_ context.Context, options ...ConnectOption) (*sql.DB, error) {
	connString := s.originalConnectionString
	connString = strings.TrimSpace(connString) // remove any leading or trailing whitespace
	connString = strings.TrimPrefix(connString, "sqlite://")

	config := newConnectConfig(options)
	db, err := sql.Open("sqlite3", connString)
	if err != nil {
		return nil, sperr.WrapWithMessage(err, "could not connect to duckdb backend")
	}
	db.SetConnMaxIdleTime(config.PoolConfig.MaxConnIdleTime)
	db.SetConnMaxLifetime(config.PoolConfig.MaxConnLifeTime)
	db.SetMaxOpenConns(config.PoolConfig.MaxOpenConns)
	return db, nil
}

// RowReader implements Backend.
func (s *SqliteBackend) RowReader() RowReader {
	return s.rowreader
}

func NewSqliteBackend(ctx context.Context, connString string) Backend {
	return &SqliteBackend{
		originalConnectionString: connString,
		rowreader:                NewSqliteRowReader(),
	}
}

type sqliteRowReader struct {
	GenericRowReader
}

func NewSqliteRowReader() *sqliteRowReader {
	return &sqliteRowReader{
		// use the generic row reader - there's no real difference between sqlite and generic
		GenericRowReader: *NewPassThruRowReader(),
	}
}
