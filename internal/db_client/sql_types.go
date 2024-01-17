package db_client

import (
	"database/sql"

	"github.com/turbot/pipe-fittings/queryresult"
)

func fieldDescriptionsToColumns(fieldDescriptions []*sql.ColumnType, connection *sql.Conn) []*queryresult.ColumnDef {
	cols := make([]*queryresult.ColumnDef, len(fieldDescriptions))

	for i, f := range fieldDescriptions {
		cols[i] = &queryresult.ColumnDef{
			Name:     f.Name(),
			DataType: f.DatabaseTypeName(),
		}
	}
	return cols
}
