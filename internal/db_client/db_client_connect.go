package db_client

import (
	"context"
	"database/sql"
	"time"

	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/powerpipe/internal/db_client/backend"
	"github.com/turbot/powerpipe/internal/db_common"
)

const (
	MaxConnLifeTime = 10 * time.Minute
	MaxConnIdleTime = 1 * time.Minute
)

// getDriverNameFromConnectionString returns the driver name for the given connection string
func getDriverNameFromConnectionString(connStr string) string {
	if backend.IsPostgresConnectionString(connStr) {
		return "pgx"
	} else if backend.IsSqliteConnectionString(connStr) {
		return "sqlite3"
	} else if backend.IsMySqlConnectionString(connStr) {
		return "mysql"
	}
	return "Unknown"
}

//type DbConnectionCallback func(context.Context, *sql.Conn) error

func (c *DbClient) establishConnectionPool(ctx context.Context) error {
	utils.LogTime("db_client.establishConnectionPool start")
	defer utils.LogTime("db_client.establishConnectionPool end")

	pool, err := establishConnectionPool(ctx, c.connectionString)
	if err != nil {
		return err
	}

	// TODO - how do we apply the AfterConnect hook here?
	// the after connect hook used to create and populate the introspection tables

	//// apply any overrides
	//// this is used to set the pool size and lifetimes of the connections from up top
	//overrides.userPoolSettings.apply(pool)
	c.UserPool = pool

	return nil
}

func establishConnectionPool(ctx context.Context, connectionString string) (*sql.DB, error) {
	driverName := getDriverNameFromConnectionString(connectionString)
	connectionString = getUseableConnectionString(driverName, connectionString)

	pool, err := sql.Open(driverName, connectionString)
	if err != nil {
		return nil, err
	}
	pool.SetConnMaxIdleTime(MaxConnIdleTime)
	pool.SetConnMaxLifetime(MaxConnLifeTime)
	pool.SetMaxOpenConns(db_common.MaxDbConnections())

	// open a connection and ping it
	conn, err := pool.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = conn.Close()
	}()
	err = conn.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return pool, nil
}
