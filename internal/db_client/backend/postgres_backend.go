package backend

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/netip"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/spf13/viper"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/queryresult"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

type PostgresBackend struct {
	originalConnectionString string
	rowreader                RowReader

	// if a custom search path or a prefix is used, store the resolved search path
	// NOTE: only applies to postgres backend
	requiredSearchPath []string
}

// Connect implements Backend.
func (s *PostgresBackend) Connect(ctx context.Context, options ...ConnectOption) (*sql.DB, error) {

	connString := s.originalConnectionString
	connector, err := NewPgxConnector(connString, s.afterConnectFunc)
	if err != nil {
		return nil, sperr.WrapWithMessage(err, "Unable to parse connection string")
	}

	config := newConnectConfig(options)
	db := sql.OpenDB(connector)
	db.SetConnMaxIdleTime(config.PoolConfig.MaxConnIdleTime)
	db.SetConnMaxLifetime(config.PoolConfig.MaxConnLifeTime)
	db.SetMaxOpenConns(config.PoolConfig.MaxOpenConns)

	// resolve the required search path
	if err := s.resolveDesiredSearchPath(ctx, db, config.SearchPathConfig); err != nil {
		return nil, err
	}
	return db, nil
}

// RowReader implements Backend.
func (s *PostgresBackend) RowReader() RowReader {
	return s.rowreader
}

func NewPostgresBackend(ctx context.Context, connString string) Backend {
	return &PostgresBackend{
		originalConnectionString: connString,
		rowreader:                NewPgxRowReader(),
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

func (c *PostgresBackend) afterConnectFunc(ctx context.Context, conn driver.Conn) error {
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

func (c *PostgresBackend) getSearchPath(ctx context.Context, db *sql.DB) ([]string, error) {
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

func (c *PostgresBackend) resolveDesiredSearchPath(ctx context.Context, db *sql.DB, searchPath *SearchPathConfig) error {
	// if we have not retrieved the original search path do it now - we do this once per client
	if c.requiredSearchPath != nil {
		return nil
	}

	if len(searchPath.SearchPath) > 0 {
		c.requiredSearchPath = cleanSearchPath(viper.GetStringSlice(constants.ArgSearchPath))
		return nil
	}

	if len(searchPath.SearchPathPrefix) > 0 {
		requiredSearchPath, err := c.constructSearchPathFromPrefix(ctx, db)
		if err != nil {
			return err
		}
		c.requiredSearchPath = requiredSearchPath
	}

	return nil
}

func (c *PostgresBackend) constructSearchPathFromPrefix(ctx context.Context, db *sql.DB) ([]string, error) {
	originalSearchPath, err := c.getSearchPath(ctx, db)
	if err != nil {
		return nil, err
	}

	searchPathPrefix := cleanSearchPath(viper.GetStringSlice(constants.ArgSearchPathPrefix))
	return append(searchPathPrefix, originalSearchPath...), nil

}

func cleanSearchPath(searchPath []string) []string {
	return helpers.RemoveFromStringSlice(searchPath, "")
}
