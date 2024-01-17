package db_client

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/queryresult"
	"github.com/turbot/pipe-fittings/statushooks"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// execute the query in the given Context
// NOTE: The returned Result MUST be fully read - otherwise the connection will block and will prevent further communication
func (c *DbClient) Execute(ctx context.Context, query string, args ...any) (*queryresult.Result, error) {
	// acquire a connection
	databaseConnection, err := c.db.Conn(ctx)
	if err != nil {
		return nil, err
	}

	// TODO KAI REMOVED <SESSION>
	//sessionResult := c.AcquireSession(ctx)
	//if sessionResult.Error != nil {
	//	return nil, sessionResult.Error
	//}

	// TODO KAI steampipe only <TIMING>
	//// disable statushooks when timing is enabled, because setShouldShowTiming internally calls the readRows funcs which
	//// calls the statushooks.Done, which hides the `Executing query…` spinner, when timing is enabled.
	//timingCtx := statushooks.DisableStatusHooks(ctx)
	//// re-read ArgTiming from viper (in case the .timing command has been run)
	//// (this will refetch ScanMetadataMaxId if timing has just been enabled)
	//c.setShouldShowTiming(timingCtx, sessionResult.Session)

	// define callback to close session when the async execution is complete
	// TODO KAI session close  waited for pg shutdown <SESSION>

	closeSessionCallback := func() { _ = databaseConnection.Close() }
	return c.executeOnConnection(ctx, databaseConnection, closeSessionCallback, query, args...)
}

// execute a query against this client and wait for the result
func (c *DbClient) ExecuteSync(ctx context.Context, query string, args ...any) (*queryresult.SyncQueryResult, error) {
	// acquire a connection
	dbConn, err := c.db.Conn(ctx)
	if err != nil {
		return nil, err
	}

	// TODO KAI REMOVED <SESSION>
	//sessionResult := c.AcquireSession(ctx)
	//if sessionResult.Error != nil {
	//	return nil, sessionResult.Error
	//}

	// TODO KAI STEAMPIPE ONLY <TIMING>
	// set setShouldShowTiming flag
	// (this will refetch ScanMetadataMaxId if timing has just been enabled)
	//c.setShouldShowTiming(ctx, sessionResult.Session)

	if c.BeforeExecuteHook != nil {
		if err := c.BeforeExecuteHook(ctx, dbConn); err != nil {
			return nil, err
		}
	}
	defer func() {

		// TODO KAI we do this in session close - move to steampipe from Session.Close <SESSION>
		//if error_helpers.IsContextCanceled(ctx) {
		//	slog.Debug("DatabaseSession.Close wait for connection cleanup")
		//	select {
		//	case <-time.After(5 * time.Second):
		//		slog.Debug("DatabaseSession.Close timed out waiting for connection cleanup")
		//		// case <-s.Connection.Conn().PgConn().CleanupDone():
		//		// 	slog.Debug("DatabaseSession.Close connection cleanup complete")
		//	}
		//}
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
	for row := range *result.RowChan {
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
	// TODO KAI STEAMPIPE ONLY <TIMING>
	//if c.shouldShowTiming() {
	//	syncResult.TimingResult = <-result.TimingResult
	//}

	return syncResult, err
}

// execute the query in the given Context using the provided DatabaseSession
// executeOnConnection assumes no responsibility over the lifecycle of the DatabaseSession - that is the responsibility of the caller
// NOTE: The returned Result MUST be fully read - otherwise the connection will block and will prevent further communication
func (c *DbClient) executeOnConnection(ctx context.Context, dbConn *sql.Conn, onComplete func(), query string, args ...any) (res *queryresult.Result, err error) {
	if query == "" {
		return queryresult.NewResult(nil), nil
	}

	// TODO KAI be clear about which execute calls we need to call the hook for - simplify? <TIMING>
	if c.BeforeExecuteHook != nil {
		if err := c.BeforeExecuteHook(ctx, dbConn); err != nil {
			return nil, err
		}
	}
	// TODO KAI steampipe only <TIMING>
	//startTime := time.Now()

	// get a context with a timeout for the query to execute within
	// we don't use the cancelFn from this timeout context, since usage will lead to 'pgx'
	// prematurely closing the database connection that this query executed in
	ctxExecute := c.getExecuteContext(ctx)

	var tx *sql.Tx

	defer func() {
		if err != nil {
			err = error_helpers.HandleQueryTimeoutError(err)
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

	result := queryresult.NewResult(colDefs)

	// read the rows in a go routine
	go func() {
		// TODO KAI do in steampipe <TIMING>
		//// define a callback which fetches the timing information
		//// this will be invoked after reading rows is complete but BEFORE closing the rows object (which closes the connection)
		//timingCallback := func() {
		//	c.getQueryTiming(ctxExecute, startTime, session, result.TimingResult)
		//}

		// read in the rows and stream to the query result object
		// TODO kai make callbacks options <TIMING>
		c.readRows(ctxExecute, rows, result, nil, nil)

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

func (c *DbClient) readRows(ctx context.Context, rows *sql.Rows, result *queryresult.Result, onRow, onComplete func()) {
	// defer this, so that these get cleaned up even if there is an unforeseen error
	defer func() {
		// we are done fetching results. time for display. clear the status indication
		statushooks.Done(ctx)
		// TODO KAI STEAMPIPE should pass timingCallback as onComplete <TIMING>
		// call the timing callback BEFORE closing the rows
		if onComplete != nil {
			onComplete()
		}
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

			// TODO KAI STEAMPIPE should pass this as onRow <MISC>
			/*
				// TACTICAL
					// determine whether to stop the spinner as soon as we stream a row or to wait for completion
				if isStreamingOutput() {
				 				statushooks.Done(ctx)
							}
			*/
			if onRow != nil {
				onRow()
			}
			// add hook?

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

func (c *DbClient) rowValues(rows *sql.Rows) ([]any, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	// Create an array of interface{} to store the retrieved values
	values := make([]interface{}, len(columns))
	// create an array to store the pointers to the values
	ptrs := make([]interface{}, len(columns))
	for i := range values {
		ptrs[i] = &values[i]
	}
	// Use a variadic to scan values into the array
	err = rows.Scan(ptrs...)
	if err != nil {
		return nil, err
	}
	return values, rows.Err()
}

func (c *DbClient) readRow(rows *sql.Rows, cols []*queryresult.ColumnDef) ([]any, error) {
	columnValues, err := c.rowValues(rows)
	if err != nil {
		return nil, error_helpers.WrapError(err)
	}
	return c.backend.RowReader().Read(columnValues, cols)
}

func isStreamingOutput() bool {
	outputFormat := viper.GetString(constants.ArgOutput)

	return helpers.StringSliceContains([]string{constants.OutputFormatCSV, constants.OutputFormatLine}, outputFormat)
}

func humanizeRowCount(count int) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", count)
}
