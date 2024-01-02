package db_client

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/powerpipe/internal/db_client/backend"
)

// define func type for StartQuery
type startQueryFunc func(ctx context.Context, dbConn *sql.Conn, query string, args ...any) (*sql.Rows, error)

// DbClient wraps over `sql.DB` and gives an interface to the database
type DbClient struct {
	connectionString string

	// connection UserPool for user initiated queries
	UserPool *sql.DB

	// connection used to run system/plumbing queries (connection state, server settings)
	ManagementPool *sql.DB

	// function to start the query - defaults to startquery
	// steampipe overrides this with startQueryWithRetries
	startQueryFunc startQueryFunc

	// the backend type of the dbclient backend
	backend backend.DBClientBackendType

	// a reader which can be used to read rows from a pgx.Rows object
	rowReader backend.RowReader

	// TODO KAI new hook <TIMING>
	BeforeExecuteHook func(context.Context, *sql.Conn) error

	// if a custom search path or a prefix is used, store it here
	CustomSearchPath []string
	SearchPathPrefix []string
	// the default user search path
	UserSearchPath []string
}

func NewDbClient(ctx context.Context, connectionString string, opts ...ClientOption) (_ *DbClient, err error) {
	utils.LogTime("db_client.NewDbClient start")
	defer utils.LogTime("db_client.NewDbClient end")

	backendType, err := backend.GetBackendFromConnectionString(ctx, connectionString)
	if err != nil {
		return nil, err
	}

	client := &DbClient{
		connectionString: connectionString,
		backend:          backendType,
		rowReader:        backend.RowReaderFactory(backendType),
	}

	// set the start query func
	client.startQueryFunc = client.StartQuery

	defer func() {
		if err != nil {
			// try closing the client
			_ = client.Close(ctx)
		}
	}()

	config := clientConfig{}
	for _, o := range opts {
		o(&config)
	}

	if err := client.establishConnectionPool(ctx, config); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *DbClient) closePools() {
	if c.UserPool != nil {
		c.UserPool.Close()
	}
	if c.ManagementPool != nil {
		c.ManagementPool.Close()
	}
}

func (c *DbClient) GetConnectionString() string {
	return c.connectionString
}

// RegisterNotificationListener has an empty implementation
// NOTE: we do not (currently) support notifications from remote connections
func (c *DbClient) RegisterNotificationListener(func(notification *pgconn.Notification)) {}

// Close closes the connection to the database and shuts down the backend
func (c *DbClient) Close(context.Context) error {
	slog.Debug("DbClient.Close user pool")
	c.closePools()

	return nil
}

// TODO KAI STEAMPIPE ONLY <MISC>
// Unimplemented (sql.DB does not have a mechanism to reset pools) - refreshDbClient terminates the current connection and opens up a new connection to the service.
func (c *DbClient) ResetPools(ctx context.Context) {
	slog.Debug("db_client.ResetPools start")
	defer slog.Debug("db_client.ResetPools end")

	// c.UserPool.Reset()
	// c.ManagementPool.Reset()
}
