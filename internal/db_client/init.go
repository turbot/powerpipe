package db_client

import (
	// database connection drivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/marcboeker/go-duckdb/v2"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	// use this to configure the drivers, if necessary
	// also can be used to register custom drivers
	//
	// also can be used to defined custom connectors that can be used to connect to the Backend
}
