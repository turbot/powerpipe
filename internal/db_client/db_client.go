package db_client

import (
	"context"
	"database/sql"
	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/v2/constants"

	"github.com/turbot/pipe-fittings/v2/backend"
	"github.com/turbot/pipe-fittings/v2/utils"
)

// DbClient wraps over `sql.DB` and gives an interface to the database
type DbClient struct {
	connectionString string

	// db handle
	db *sql.DB

	// the Backend
	Backend backend.Backend
}

func NewDbClient(ctx context.Context, connectionString string, opts ...backend.BackendOption) (_ *DbClient, err error) {
	utils.LogTime("db_client.NewDbClient start")
	defer utils.LogTime("db_client.NewDbClient end")

	b, err := backend.FromConnectionString(ctx, connectionString)
	if err != nil {
		return nil, err
	}

	client := &DbClient{
		connectionString: connectionString,
		Backend:          b,
	}

	defer func() {
		if err != nil {
			// try closing the client
			_ = client.Close(ctx)
		}
	}()

	// process options - search path may have been passed in
	config := backend.NewBackendConfig(opts)
	config.MaxOpenConns = MaxDbConnections()
	// if no search path override passed in as an option, use the viper config
	if config.SearchPathConfig.Empty() {
		config.SearchPathConfig = backend.SearchPathConfig{
			SearchPath:       viper.GetStringSlice(constants.ArgSearchPath),
			SearchPathPrefix: viper.GetStringSlice(constants.ArgSearchPathPrefix),
		}
	}

	if err := client.connect(ctx, backend.WithConfig(config)); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *DbClient) GetConnectionString() string {
	return c.connectionString
}

// Close closes the connection to the database and shuts down the Backend
func (c *DbClient) Close(_ context.Context) error {

	if c.db != nil {
		if err :=  c.db.Close(); err != nil {
			return err
		}
	}
	// if the backend has a cleanup method, call it
	if cleaner, ok := c.Backend.(interface { Cleanup() error }); ok {
		if err := cleaner.Cleanup(); err != nil {
			return err
		}
	}
	return nil
}
