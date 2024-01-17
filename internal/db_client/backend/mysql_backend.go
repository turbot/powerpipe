package backend

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/turbot/pipe-fittings/queryresult"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

const mysqlConnectionStringPrefix = "mysql://"

type MySQLBackend struct {
	connectionString string
	rowreader        RowReader
}

func NewMySQLBackend(ctx context.Context, connString string) Backend {
	connString = strings.TrimSpace(connString) // remove any leading or trailing whitespace
	connString = strings.TrimPrefix(connString, mysqlConnectionStringPrefix)

	return &MySQLBackend{
		connectionString: connString,
		rowreader:        NewMySqlRowReader(),
	}
}

// Connect implements Backend.
func (s *MySQLBackend) Connect(_ context.Context, options ...ConnectOption) (*sql.DB, error) {
	config := newConnectConfig(options)
	db, err := sql.Open("mysql", s.connectionString)
	if err != nil {
		return nil, sperr.WrapWithMessage(err, "could not connect to duckdb backend")
	}
	db.SetConnMaxIdleTime(config.PoolConfig.MaxConnIdleTime)
	db.SetConnMaxLifetime(config.PoolConfig.MaxConnLifeTime)
	db.SetMaxOpenConns(config.PoolConfig.MaxOpenConns)
	return db, nil
}

// RowReader implements Backend.
func (s *MySQLBackend) RowReader() RowReader {
	return s.rowreader
}

type mysqlRowReader struct {
	BasicRowReader
}

func NewMySqlRowReader() RowReader {
	return &mysqlRowReader{
		BasicRowReader: BasicRowReader{
			CellReader: mysqlReadCell,
		},
	}
}

func mysqlReadCell(columnValue any, col *queryresult.ColumnDef) (result any, err error) {
	if columnValue != nil {
		asStr := string(columnValue.([]byte))
		switch col.DataType {
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "BIGINT", "YEAR":
			result, err = strconv.ParseInt(asStr, 10, 64)
		case "DECIMAL", "NUMERIC", "FLOAT", "DOUBLE":
			result, err = strconv.ParseFloat(asStr, 64)
		case "DATE":
			result, err = time.Parse(time.DateOnly, asStr)
		case "TIME":
			result, err = time.Parse(time.TimeOnly, asStr)
		case "DATETIME", "TIMESTAMP":
			result, err = time.Parse(time.DateTime, asStr)
		case "BIT":
			result = columnValue.([]byte)
		// case "CHAR", "VARCHAR", "TEXT", "BINARY", "VARBINARY", "ENUM", "SET":
		default:
			result = asStr
		}
	}
	return result, err
}
