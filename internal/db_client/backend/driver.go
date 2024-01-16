package backend

var Drivers = map[DBClientBackendType]string{
	PostgresDBClientBackend: "pgx",
	DuckDBClientBackend:     "duckdb",
	MySQLDBClientBackend:    "mysql",
	SqliteDBClientBackend:   "sqlite3",
}
