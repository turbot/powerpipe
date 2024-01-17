package backend

import (
	"context"
	"database/sql"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

const duckDBConnectionStringPrefix = "duckdb://"

type DuckDBBackend struct {
	connectionString string
	rowreader        RowReader
}

func NewDuckDBBackend(ctx context.Context, connString string) Backend {
	connString = strings.TrimSpace(connString) // remove any leading or trailing whitespace
	connString = strings.TrimPrefix(connString, duckDBConnectionStringPrefix)
	return &DuckDBBackend{
		connectionString: connString,
		rowreader:        NewDuckDBRowReader(),
	}
}

// Connect implements Backend.
func (s *DuckDBBackend) Connect(_ context.Context, options ...ConnectOption) (*sql.DB, error) {
	config := newConnectConfig(options)
	db, err := sql.Open("duckdb", s.connectionString)
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

type duckdbRowReader struct {
	BasicRowReader
}

func NewDuckDBRowReader() *duckdbRowReader {
	return &duckdbRowReader{
		// use the generic row reader - there's no real difference between sqlite and duckdb
		BasicRowReader: *NewBasicRowReader(),
	}
}
