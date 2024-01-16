package backend

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/turbot/pipe-fittings/queryresult"
)

type MySQLBackend struct {
	originalConnectionString string
	rowreader                RowReader
}

// Connect implements Backend.
func (s *MySQLBackend) Connect(context.Context, ...ConnectOption) (*sql.DB, error) {
	connString := s.originalConnectionString
	connString = strings.TrimSpace(connString) // remove any leading or trailing whitespace
	connString = strings.TrimPrefix(connString, "mysql://")
	return sql.Open("mysql", connString)
}

// GetType implements Backend.
func (s *MySQLBackend) GetType() BackendType {
	return MySQLDBClientBackend
}

// RowReader implements Backend.
func (s *MySQLBackend) RowReader() RowReader {
	return s.rowreader
}

func NewMySQLBackend(ctx context.Context, connString string) Backend {
	return &MySQLBackend{
		rowreader: NewMySqlRowReader(),
	}
}

type mysqlRowReader struct {
	GenericRowReader
}

func NewMySqlRowReader() RowReader {
	return &mysqlRowReader{
		GenericRowReader: GenericRowReader{
			CellReader: mysqlReadCell,
		},
	}
}

func mysqlReadCell(columnValue any, col *queryresult.ColumnDef) (any, error) {
	var result any
	if columnValue != nil {
		asStr := string(columnValue.([]byte))
		switch col.DataType {
		case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "BIGINT", "YEAR":
			r, _ := strconv.ParseInt(asStr, 10, 64)
			result = r
		case "DECIMAL", "NUMERIC", "FLOAT", "DOUBLE":
			r, _ := strconv.ParseFloat(asStr, 64)
			result = r
		case "DATE":
			tim, _ := time.Parse(time.DateOnly, asStr)
			result = tim
		case "TIME":
			tim, _ := time.Parse(time.TimeOnly, asStr)
			result = tim
		case "DATETIME", "TIMESTAMP":
			tim, _ := time.Parse(time.DateTime, asStr)
			result = tim
		case "BIT":
			result = columnValue.([]byte)
		// case "CHAR", "VARCHAR", "TEXT", "BINARY", "VARBINARY", "ENUM", "SET":
		default:
			result = asStr
		}
	}
	return result, nil
}
