package backend

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/turbot/pipe-fittings/queryresult"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

var ErrUnknownBackend = errors.New("unknown backend")

type RowReader interface {
	Read(columnValues []any, cols []*queryresult.ColumnDef) ([]any, error)
}

type Backend interface {
	Connect(context.Context, ...ConnectOption) (*sql.DB, error)
	RowReader() RowReader
}

func FromConnectionString(ctx context.Context, str string) (Backend, error) {
	switch {
	case IsPostgresConnectionString(str):
		return NewPostgresBackend(ctx, str), nil
	case IsMySqlConnectionString(str):
		return NewMySQLBackend(ctx, str), nil
	case IsDuckDBConnectionString(str):
		return NewDuckDBBackend(ctx, str), nil
	case IsSqliteConnectionString(str):
		return NewSqliteBackend(ctx, str), nil
	}
	return nil, sperr.WrapWithMessage(ErrUnknownBackend, "could not evaluate backend: %s", str)
}

func HasBackend(connString string) bool {
	if m, err := FromConnectionString(context.Background(), connString); m != nil && err == nil {
		return true
	}
	return false
}

// IsPostgresConnectionString returns true if the connection string is for postgres
// looks for the postgresql:// or postgres:// prefix
func IsPostgresConnectionString(connString string) bool {
	for _, v := range postgresConnectionStringPrefixes {
		if strings.HasPrefix(connString, v) {
			return true
		}
	}
	return false
}

// IsSqliteConnectionString returns true if the connection string is for sqlite
// looks for the sqlite:// prefix
func IsSqliteConnectionString(connString string) bool {
	return strings.HasPrefix(connString, sqliteConnectionStringPrefix)
}

// IsDuckDBConnectionString returns true if the connection string is for duckdb
// looks for the duckdb:// prefix
func IsDuckDBConnectionString(connString string) bool {
	return strings.HasPrefix(connString, duckDBConnectionStringPrefix)
}

// IsMySqlConnectionString returns true if the connection string is for mysql
// looks for the mysql:// prefix
func IsMySqlConnectionString(connString string) bool {
	return strings.HasPrefix(connString, mysqlConnectionStringPrefix)
}
