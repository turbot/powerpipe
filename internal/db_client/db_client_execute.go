package db_client

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"time"

	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/v2/constants"
	"github.com/turbot/pipe-fittings/v2/error_helpers"
	"github.com/turbot/pipe-fittings/v2/queryresult"
	"github.com/turbot/pipe-fittings/v2/statushooks"
	localqueryresult "github.com/turbot/powerpipe/internal/queryresult"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Execute executes the query in the given Context
// NOTE: The returned Result MUST be fully read - otherwise the connection will block and will prevent further communication
func (c *DbClient) Execute(ctx context.Context, query string, args ...any) (*localqueryresult.Result, error) {
	// acquire a connection
	databaseConnection, err := c.db.Conn(ctx)
	if err != nil {
		return nil, err
	}

	// define callback to close session when the async execution is complete
	closeSessionCallback := func() { _ = databaseConnection.Close() }
	return c.executeOnConnection(ctx, databaseConnection, closeSessionCallback, query, args...)
}

// ExecuteSync executes a query against this client and wait for the result
func (c *DbClient) ExecuteSync(ctx context.Context, query string, args ...any) (*queryresult.SyncQueryResult, error) {
	// acquire a connection
	dbConn, err := c.db.Conn(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		dbConn.Close()

	}()
	return c.executeSyncOnConnection(ctx, dbConn, query, args...)
}

// execute a query against this client and wait for the result
func (c *DbClient) executeSyncOnConnection(ctx context.Context, dbConn *sql.Conn, query string, args ...any) (*queryresult.SyncQueryResult, error) {
	if query == "" {
		return &queryresult.SyncQueryResult{}, nil
	}

	result, err := c.executeOnConnection(ctx, dbConn, nil, query, args...)
	if err != nil {
		return nil, error_helpers.WrapError(err)
	}

	syncResult := &queryresult.SyncQueryResult{Cols: result.Cols}
	for row := range result.RowChan {
		select {
		case <-ctx.Done():
		default:
			// save the first row error to return
			if row.Error != nil && err == nil {
				err = error_helpers.WrapError(row.Error)
			}
			syncResult.Rows = append(syncResult.Rows, row)
		}
	}
	return syncResult, err
}

// execute the query in the given Context using the provided DatabaseSession
// executeOnConnection assumes no responsibility over the lifecycle of the DatabaseSession - that is the responsibility of the caller
// NOTE: The returned Result MUST be fully read - otherwise the connection will block and will prevent further communication
func (c *DbClient) executeOnConnection(ctx context.Context, dbConn *sql.Conn, onComplete func(), query string, args ...any) (res *localqueryresult.Result, err error) {
	if query == "" {
		return localqueryresult.NewResult(nil), nil
	}

	// get a context with a timeout for the query to execute within
	// we don't use the cancelFn from this timeout context, since usage will lead to 'pgx'
	// prematurely closing the database connection that this query executed in
	ctxExecute := c.getExecuteContext(ctx)

	var tx *sql.Tx

	defer func() {
		if err != nil {
			// stop spinner in case of error
			statushooks.Done(ctxExecute)
			// error - rollback transaction if we have one
			if tx != nil {
				_ = tx.Rollback()
			}
			// in case of error call the onComplete callback
			if onComplete != nil {
				onComplete()
			}
		}
	}()

	// start query

	rows, err := c.StartQuery(ctxExecute, dbConn, query, args...)
	if err != nil {
		return
	}

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return
	}
	colDefs := fieldDescriptionsToColumns(colTypes, dbConn)

	result := localqueryresult.NewResult(colDefs)

	// read the rows in a go routine
	go func() {
		// read in the rows and stream to the query result object
		c.readRows(ctxExecute, rows, result)

		// call the completion callback - if one was provided
		if onComplete != nil {
			onComplete()
		}
	}()

	return result, nil
}

func (c *DbClient) getExecuteContext(ctx context.Context) context.Context {
	queryTimeout := time.Duration(viper.GetInt(constants.ArgDatabaseQueryTimeout)) * time.Second
	// if timeout is zero, do not set a timeout
	if queryTimeout == 0 {
		return ctx
	}
	// create a context with a deadline
	shouldBeDoneBy := time.Now().Add(queryTimeout)
	//nolint:govet //we don't use this cancel fn because, pgx prematurely cancels the PG connection when this cancel gets called in 'defer'
	newCtx, _ := context.WithDeadline(ctx, shouldBeDoneBy)

	return newCtx
}

// StartQuery runs query in a goroutine, so we can check for cancellation
// in case the client becomes unresponsive and does not respect context cancellation
func (c *DbClient) StartQuery(ctx context.Context, dbConn *sql.Conn, query string, args ...any) (rows *sql.Rows, err error) {
	doneChan := make(chan bool)
	go func() {
		// start asynchronous query
		//nolint: sqlclosecheck // rows is closed in readRows
		rows, err = dbConn.QueryContext(ctx, query, args...)
		close(doneChan)
	}()

	select {
	case <-doneChan:
	case <-ctx.Done():
		err = ctx.Err()
	}
	return
}

func (c *DbClient) readRows(ctx context.Context, rows *sql.Rows, result *localqueryresult.Result) {
	// defer this, so that these get cleaned up even if there is an unforeseen error
	defer func() {
		// we are done fetching results. time for display. clear the status indication
		statushooks.Done(ctx)
		// close the sql rows object
		rows.Close()
		if err := rows.Err(); err != nil {
			result.StreamError(err)
		}
		// close the channels in the result object
		result.Close()

	}()

	rowCount := 0
Loop:
	for rows.Next() {
		select {
		case <-ctx.Done():
			statushooks.SetStatus(ctx, "Cancelling query")
			break Loop
		default:
			rowResult, err := c.readRow(rows, result.Cols)
			if err != nil {
				// the error will be streamed in the defer
				break Loop
			}

			if isStreamingOutput() {
				statushooks.Done(ctx)
			}

			result.StreamRow(rowResult)

			// update the status message with the count of rows that have already been fetched
			// this will not show if the spinner is not active
			statushooks.SetStatus(ctx, fmt.Sprintf("Loading results: %3s", humanizeRowCount(rowCount)))
			rowCount++
		}
	}
}

func (c *DbClient) rowValues(rows *sql.Rows, cols []*queryresult.ColumnDef) ([]any, error) {
	// Create an array of any to store the retrieved values
	values := make([]any, len(cols))

	// create an array to store the pointers to the values
	ptrs := make([]any, len(cols))
	for i := range values {
		ptrs[i] = &values[i]
	}
	// Use a variadic to scan values into the array
	err := rows.Scan(ptrs...)
	if err != nil {
		return nil, err
	}

	return values, rows.Err()
}

func (c *DbClient) readRow(rows *sql.Rows, cols []*queryresult.ColumnDef) ([]any, error) {
	columnValues, err := c.rowValues(rows, cols)
	if err != nil {
		return nil, error_helpers.WrapError(err)
	}
	return c.Backend.RowReader().Read(columnValues, cols)
}

func isStreamingOutput() bool {
	outputFormat := viper.GetString(constants.ArgOutput)

	return slices.Contains([]string{constants.OutputFormatCSV, constants.OutputFormatLine}, outputFormat)
}

func humanizeRowCount(count int) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", count)
}
