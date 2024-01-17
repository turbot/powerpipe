package backend

import (
	"context"
	"database/sql"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

type DuckDBBackend struct {
	originalConnectionString string
	rowreader                RowReader
}

// Connect implements Backend.
func (s *DuckDBBackend) Connect(_ context.Context, options ...ConnectOption) (*sql.DB, error) {
	connString := s.originalConnectionString
	connString = strings.TrimSpace(connString) // remove any leading or trailing whitespace
	connString = strings.TrimPrefix(connString, "duckdb://")

	config := newConnectConfig(options)
	db, err := sql.Open("duckdb", connString)
	if err != nil {
		return nil, sperr.WrapWithMessage(err, "could not connect to duckdb backend")
	}
	db.SetConnMaxIdleTime(config.PoolConfig.MaxConnIdleTime)
	db.SetConnMaxLifetime(config.PoolConfig.MaxConnLifeTime)
	db.SetMaxOpenConns(config.PoolConfig.MaxOpenConns)
	return db, nil
}

// RowReader implements Backend.
func (s *DuckDBBackend) RowReader() RowReader {
	return s.rowreader
}

func NewDuckDBBackend(ctx context.Context, connString string) Backend {
	return &DuckDBBackend{
		originalConnectionString: connString,
		rowreader:                NewDuckDBRowReader(),
	}
}

type duckdbRowReader struct {
	GenericRowReader
}

func NewDuckDBRowReader() *duckdbRowReader {
	return &duckdbRowReader{
		// use the generic row reader - there's no real difference between sqlite and duckdb
		GenericRowReader: *NewPassThruRowReader(),
	}
}
