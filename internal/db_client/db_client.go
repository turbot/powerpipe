package db_client

import (
	"context"
	"database/sql"

	"github.com/turbot/pipe-fittings/backend"
	"github.com/turbot/pipe-fittings/utils"
)

// DbClient wraps over `sql.DB` and gives an interface to the database
type DbClient struct {
	connectionString string

	// db handle
	db *sql.DB

	// the backend
	backend backend.Backend

	// TODO KAI new hook <TIMING>
	BeforeExecuteHook func(context.Context, *sql.Conn) error
}

func NewDbClient(ctx context.Context, connectionString string) (_ *DbClient, err error) {
	utils.LogTime("db_client.NewDbClient start")
	defer utils.LogTime("db_client.NewDbClient end")

	backend, err := backend.FromConnectionString(ctx, connectionString)
	if err != nil {
		return nil, err
	}

	client := &DbClient{
		connectionString: connectionString,
		backend:          backend,
	}

	defer func() {
		if err != nil {
			// try closing the client
			_ = client.Close(ctx)
		}
	}()

	if err := client.connect(ctx); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *DbClient) GetConnectionString() string {
	return c.connectionString
}

// Close closes the connection to the database and shuts down the backend
func (c *DbClient) Close(context.Context) error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}
