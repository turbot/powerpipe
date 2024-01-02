package backend

import (
	"context"
	"strings"
)

func GetBackendFromConnectionString(ctx context.Context, connectionString string) (DBClientBackendType, error) {
	if IsPostgresConnectionString(connectionString) {
		return PostgresDBClientBackend, nil
	} else if IsMySqlConnectionString(connectionString) {
		return MySQLDBClientBackend, nil
	}
	return SqliteDBClientBackend, nil
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

// IsMySqlConnectionString returns true if the connection string is for mysql
// looks for the mysql:// prefix
func IsMySqlConnectionString(connString string) bool {
	return strings.HasPrefix(connString, "mysql://")
}
