package backend

import (
	"strconv"
	"time"

	"github.com/turbot/pipe-fittings/queryresult"
)

type mysqlRowReader struct {
	genericSQLRowReader
}

func NewMySqlRowReader() RowReader {
	return &mysqlRowReader{
		genericSQLRowReader: genericSQLRowReader{
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
