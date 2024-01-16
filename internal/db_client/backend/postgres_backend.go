package backend

import (
	"context"
	"database/sql"
	"fmt"
	"net/netip"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/turbot/pipe-fittings/queryresult"
	"github.com/turbot/pipe-fittings/utils"
)

type PostgresBackend struct {
	originalConnectionString string
	rowreader                RowReader
}

// Connect implements Backend.
func (s *PostgresBackend) Connect(context.Context, ...ConnectOption) (*sql.DB, error) {
	connString := s.originalConnectionString
	return sql.Open("postgres", connString)
}

// GetType implements Backend.
func (s *PostgresBackend) GetType() BackendType {
	return PostgresDBClientBackend
}

// RowReader implements Backend.
func (s *PostgresBackend) RowReader() RowReader {
	return s.rowreader
}

func NewPostgresBackend(ctx context.Context, connString string) Backend {
	return &PostgresBackend{
		rowreader: NewPgxRowReader(),
	}
}

func NewPgxRowReader() *pgxRowReader {
	return &pgxRowReader{
		GenericRowReader: GenericRowReader{
			CellReader: pgxReadCell,
		},
	}
}

// pgxRowReader is a RowReader implementation for the pgx database/sql driver
type pgxRowReader struct {
	GenericRowReader
}

func pgxReadCell(columnValue any, col *queryresult.ColumnDef) (any, error) {
	var result any
	if columnValue != nil {
		result = columnValue

		// add special handling for some types
		switch col.DataType {
		case "_TEXT":
			if arr, ok := columnValue.([]interface{}); ok {
				elements := utils.Map(arr, func(e interface{}) string { return e.(string) })
				result = strings.Join(elements, ",")
			}
		case "INET":
			if inet, ok := columnValue.(netip.Prefix); ok {
				result = strings.TrimSuffix(inet.String(), "/32")
			}
		case "UUID":
			if bytes, ok := columnValue.([16]uint8); ok {
				if u, err := uuid.FromBytes(bytes[:]); err == nil {
					result = u
				}
			}
		case "TIME":
			if t, ok := columnValue.(pgtype.Time); ok {
				result = time.UnixMicro(t.Microseconds).UTC().Format("15:04:05")
			}
		case "INTERVAL":
			if interval, ok := columnValue.(pgtype.Interval); ok {
				var sb strings.Builder
				years := interval.Months / 12
				months := interval.Months % 12
				if years > 0 {
					sb.WriteString(fmt.Sprintf("%d %s ", years, utils.Pluralize("year", int(years))))
				}
				if months > 0 {
					sb.WriteString(fmt.Sprintf("%d %s ", months, utils.Pluralize("mon", int(months))))
				}
				if interval.Days > 0 {
					sb.WriteString(fmt.Sprintf("%d %s ", interval.Days, utils.Pluralize("day", int(interval.Days))))
				}
				if interval.Microseconds > 0 {
					d := time.Duration(interval.Microseconds) * time.Microsecond
					formatStr := time.Unix(0, 0).UTC().Add(d).Format("15:04:05")
					sb.WriteString(formatStr)
				}
				result = sb.String()
			}

		case "NUMERIC":
			if numeric, ok := columnValue.(pgtype.Numeric); ok {
				if f, err := numeric.Float64Value(); err == nil {
					result = f.Float64
				}
			}
		}
	}
	return result, nil
}
