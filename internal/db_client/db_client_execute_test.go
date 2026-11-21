package db_client

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/v2/constants"
)

func newTestClient(t *testing.T) *DbClient {
	t.Helper()
	ctx := context.Background()
	client, err := NewDbClient(ctx, "sqlite://"+filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = client.Close(ctx) })
	return client
}

// cancelled queries release their connections. With -race, this also checks that no query uses its
// connection after the connection is closed. SQLite usually returns the result before it sees the cancel,
// and then the result is read to the end.
func TestExecuteOnConnection_Cancelled(t *testing.T) {
	client := newTestClient(t)
	for range 500 {
		conn, err := client.db.Conn(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		// as Execute does, close the connection when the query is done
		result, err := client.executeOnConnection(ctx, conn, func() { _ = conn.Close() }, "select 1")
		if err != nil {
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("expected a cancelled error, got %v", err)
			}
			continue
		}
		// the driver can return rows before it sees the cancellation: then the result must be read to the end
		for range result.RowChan {
		}
	}
	for deadline := time.Now().Add(10 * time.Second); client.db.Stats().InUse > 0; time.Sleep(10 * time.Millisecond) {
		if time.Now().After(deadline) {
			t.Fatalf("%d connections are still in use", client.db.Stats().InUse)
		}
	}
}

func TestExecuteSync(t *testing.T) {
	client := newTestClient(t)
	result, err := client.ExecuteSync(context.Background(), "select 42 as n")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(result.Rows))
	}
}

// a database/sql connector whose queries ignore their context: a query blocks until the test releases it
type blockingConnector struct {
	entered chan struct{}
	release chan struct{}
}

func (b *blockingConnector) Connect(context.Context) (driver.Conn, error) {
	return blockingConn{b}, nil
}
func (b *blockingConnector) Driver() driver.Driver { return blockingDriver{b} }

type blockingDriver struct{ b *blockingConnector }

func (d blockingDriver) Open(string) (driver.Conn, error) { return blockingConn(d), nil }

type blockingConn struct{ b *blockingConnector }

func (c blockingConn) Prepare(string) (driver.Stmt, error) { return nil, errors.New("not supported") }
func (c blockingConn) Close() error                        { return nil }
func (c blockingConn) Begin() (driver.Tx, error)           { return nil, errors.New("not supported") }
func (c blockingConn) QueryContext(ctx context.Context, _ string, _ []driver.NamedValue) (driver.Rows, error) {
	close(c.b.entered)
	<-c.b.release
	return nil, ctx.Err()
}

// a query that keeps running after its context is cancelled holds executeOnConnection until it stops, so the
// caller cannot close the connection while the query still uses it (#1181). This does not need -race.
func TestExecuteOnConnection_WaitsForCancelledQuery(t *testing.T) {
	b := &blockingConnector{entered: make(chan struct{}), release: make(chan struct{})}
	release := sync.OnceFunc(func() { close(b.release) })
	db := sql.OpenDB(b)
	defer db.Close()
	conn, err := db.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	defer release()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := (&DbClient{}).executeOnConnection(ctx, conn, nil, "select 1")
		done <- err
	}()
	// cancel once the query is running: database/sql does not start a query whose context is already cancelled
	<-b.entered
	cancel()
	select {
	case err := <-done:
		t.Fatalf("executeOnConnection returned while the query was still running (%v)", err)
	case <-time.After(200 * time.Millisecond):
	}
	release()
	if err := <-done; !errors.Is(err, context.Canceled) {
		t.Fatalf("expected the context error, got %v", err)
	}
}

// a query that runs for many seconds before it returns its single row
const slowQuery = "with recursive c(x) as (select 1 union all select x+1 from c where x < 300000000) select count(*) from c"

// a query with a lot of rows, to stop reading in the middle of them
const manyRowsQuery = "with recursive c(x) as (select 1 union all select x+1 from c where x < 1000000) select x from c"

func waitForNoConnectionsInUse(t *testing.T, client *DbClient) {
	t.Helper()
	for deadline := time.Now().Add(20 * time.Second); client.db.Stats().InUse > 0; time.Sleep(10 * time.Millisecond) {
		if time.Now().After(deadline) {
			t.Fatalf("%d connections are still in use", client.db.Stats().InUse)
		}
	}
}

// #1260: SQLite reports a cancel through the rows, so a result that nobody reads any more used to block its
// reader on the send of that error for ever and keep the connection. More than the pool holds must not hang.
func TestExecute_CancelledAndUnreadReleasesConnection(t *testing.T) {
	for name, query := range map[string]string{"slow query": slowQuery, "many rows": manyRowsQuery} {
		t.Run(name, func(t *testing.T) {
			client := newTestClient(t)
			for range 12 {
				// the deadline makes a leak fail the test, as Execute would otherwise wait for a free connection
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				result, err := client.Execute(ctx, query)
				if err != nil {
					cancel()
					t.Fatal(err)
				}
				if name == "many rows" {
					// take one row, then stop reading
					if row := <-result.RowChan; row == nil || row.Error != nil {
						t.Fatalf("expected a row, got %v", row)
					}
				}
				time.Sleep(50 * time.Millisecond)
				cancel()
			}
			waitForNoConnectionsInUse(t, client)

			// every connection is free again: a new query runs
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			result, err := client.ExecuteSync(ctx, "select 42 as n")
			if err != nil {
				t.Fatal(err)
			}
			if len(result.Rows) != 1 {
				t.Fatalf("expected 1 row, got %d", len(result.Rows))
			}
		})
	}
}

// a query that the database timeout stops, while the reader is still reading, is still reported as an error
func TestExecute_DatabaseTimeoutIsReportedToReader(t *testing.T) {
	viper.Set(constants.ArgDatabaseQueryTimeout, 1)
	t.Cleanup(func() { viper.Set(constants.ArgDatabaseQueryTimeout, 0) })
	client := newTestClient(t)

	result, err := client.Execute(context.Background(), slowQuery)
	if err != nil {
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expected a timeout error, got %v", err)
		}
		return
	}
	var gotErr error
	for row := range result.RowChan {
		if row.Error != nil {
			gotErr = row.Error
		}
	}
	if gotErr == nil {
		t.Fatal("expected the timeout to reach the reader as an error")
	}
	waitForNoConnectionsInUse(t, client)
}
