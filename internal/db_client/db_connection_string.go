package db_client

import (
	"strings"

	"github.com/turbot/pipe-fittings/db_client/backend"
)

// getUseableConnectionString returns a connection string that can be used by the database driver
// this is derived from the connection string passed in by the user and the driver name
func getUseableConnectionString(driver string, connString string) string {
	// using this to remove the sqlite3?:// prefix from the connection string
	// this is necessary because the sqlite3 driver doesn't recognize the sqlite3?:// prefix
	connString = strings.TrimPrefix(connString, "sqlite3://")
	connString = strings.TrimPrefix(connString, "sqlite://")

	// case for mysql connection strings
	connString = strings.TrimPrefix(connString, "mysql://")
	return connString
}

func IsConnectionString(connString string) bool {
	isPostgres := backend.IsPostgresConnectionString(connString)
	isSqlite := backend.IsSqliteConnectionString(connString)
	isMysql := backend.IsMySqlConnectionString(connString)
	return isPostgres || isSqlite || isMysql
}
