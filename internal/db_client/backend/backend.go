package backend

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/turbot/pipe-fittings/queryresult"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

type RowReader interface {
	Read(columnValues []any, cols []*queryresult.ColumnDef) ([]any, error)
}

type Backend interface {
	GetType() DBClientBackendType
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

var ErrUnknownBackend = errors.New("unknown backend")

func GetBackendFromConnectionString(ctx context.Context, connectionString string) (DBClientBackendType, error) {
	switch {
	case IsPostgresConnectionString(connectionString):
		return PostgresDBClientBackend, nil
	case IsMySqlConnectionString(connectionString):
		return MySQLDBClientBackend, nil
	case IsDuckDBConnectionString(connectionString):
		return DuckDBClientBackend, nil
	case IsSqliteConnectionString(connectionString):
		return SqliteDBClientBackend, nil
	}
	return UnknownClientBackend, sperr.WrapWithMessage(ErrUnknownBackend, "could not evaluate backend: %s", connectionString)
}

// IsPostgresConnectionString returns true if the connection string is for postgres
// looks for the postgresql:// or postgres:// prefix
func IsPostgresConnectionString(connString string) bool {
	return strings.HasPrefix(connString, "postgresql://") || strings.HasPrefix(connString, "postgres://")
}

// IsSqliteConnectionString returns true if the connection string is for sqlite
// looks for the sqlite:// prefix
func IsSqliteConnectionString(connString string) bool {
	return strings.HasPrefix(connString, "sqlite://")
}

// IsDuckDBConnectionString returns true if the connection string is for duckdb
// looks for the duckdb:// prefix
func IsDuckDBConnectionString(connString string) bool {
	return strings.HasPrefix(connString, "duckdb://")
}

// IsMySqlConnectionString returns true if the connection string is for mysql
// looks for the mysql:// prefix
func IsMySqlConnectionString(connString string) bool {
	return strings.HasPrefix(connString, "mysql://")
}
