package backend

var Drivers = map[BackendType]string{
	PostgresDBClientBackend: "pgx",
	DuckDBClientBackend:     "duckdb",
	MySQLDBClientBackend:    "mysql",
	SqliteDBClientBackend:   "sqlite3",
}
