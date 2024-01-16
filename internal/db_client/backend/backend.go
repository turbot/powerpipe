package backend

import (
	"context"
	"errors"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

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
