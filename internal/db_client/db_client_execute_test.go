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
