package db_client

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/spf13/viper"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/constants"
	"strings"
	"time"

	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/powerpipe/internal/db_client/backend"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
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

func (c *DbClient) establishConnectionPool(ctx context.Context) error {
	utils.LogTime("db_client.establishConnectionPool start")
	defer utils.LogTime("db_client.establishConnectionPool end")

	driverName := getDriverNameFromConnectionString(c.connectionString)
	connectionString := getUseableConnectionString(driverName, c.connectionString)

	var db *sql.DB
	if driverName == "pgx" {
		connector, err := NewPgxConnector(connectionString, c.afterConnectFunc) //searchPath)
		if err != nil {
			return sperr.WrapWithMessage(err, "Unable to parse connection string")
		}

		db = sql.OpenDB(connector)
		// resolve the required search path
		if err := c.resolveDesiredSearchPath(ctx, db); err != nil {
			return err
		}

	} else {
		var err error
		db, err = sql.Open(driverName, connectionString)
		if err != nil {
			return err
		}

		// open a connection and ping it
		// (no need for postgres as we have set the search path)
		conn, err := db.Conn(ctx)
		if err != nil {
			return err
		}
		defer func() {
			_ = conn.Close()
		}()
		err = conn.PingContext(ctx)
		if err != nil {
			return err
		}
	}

	db.SetConnMaxIdleTime(MaxConnIdleTime)
	db.SetConnMaxLifetime(MaxConnLifeTime)
	db.SetMaxOpenConns(MaxDbConnections())

	c.db = db
	return nil
}

func (c *DbClient) afterConnectFunc(ctx context.Context, conn driver.Conn) error {
	if len(c.requiredSearchPath) == 0 {
		return nil
	}
	connPc, ok := conn.(driver.ConnPrepareContext)
	if !ok {
		return fmt.Errorf("stdlib driver does not implement ConnPrepareContext")
	}
	ps, err := connPc.PrepareContext(ctx, "SET search_path TO "+strings.Join(c.requiredSearchPath, ","))
	if err != nil {
		return err
	}
	ec, ok := ps.(driver.StmtExecContext)
	if !ok {
		return fmt.Errorf("prepared statement does not implement StmtExecContext")
	}
	defer ps.Close()

	_, err = ec.ExecContext(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *DbClient) getSearchPath(ctx context.Context, db *sql.DB) ([]string, error) {
	// Get a connection from the database
	conn, err := db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Execute the query
	row := conn.QueryRowContext(ctx, "SHOW search_path")
	if row.Err() != nil {
		return nil, row.Err()
	}

	var searchPath string
	// Scan the result into the searchPath variable
	err = row.Scan(&searchPath)
	if err != nil {
		return nil, err
	}

	// Split the search path into individual paths
	searchPaths := strings.Split(searchPath, ",")
	// Trim spaces from each path
	for i, path := range searchPaths {
		searchPaths[i] = strings.TrimSpace(path)
	}

	return searchPaths, nil
}

func (c *DbClient) resolveDesiredSearchPath(ctx context.Context, db *sql.DB) error {
	// if we have not retrieved the original search path do it now - we do this once per client
	if c.requiredSearchPath != nil {
		return nil
	}

	if viper.IsSet(constants.ArgSearchPath) {
		c.requiredSearchPath = cleanSearchPath(viper.GetStringSlice(constants.ArgSearchPath))
		return nil
	}

	if viper.IsSet(constants.ArgSearchPathPrefix) {
		originalSearchPath, err := c.getSearchPath(ctx, db)
		if err != nil {
			return err
		}

		searchPathPrefix := cleanSearchPath(viper.GetStringSlice(constants.ArgSearchPathPrefix))
		c.requiredSearchPath = append(searchPathPrefix, originalSearchPath...)
		return nil
	}

	return nil
}

func cleanSearchPath(searchPath []string) []string {
	return helpers.RemoveFromStringSlice(searchPath, "")
}
