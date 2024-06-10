package db_client

import (
	"context"
	"database/sql"
	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/backend"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/utils"
	"log/slog"
)

// DbClient wraps over `sql.DB` and gives an interface to the database
type DbClient struct {
	connectionString string

	// db handle
	db *sql.DB

	// the Backend
	Backend backend.Backend
}

func NewDbClient(ctx context.Context, connectionString string, opts ...backend.ConnectOption) (_ *DbClient, err error) {
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
	config := backend.NewConnectConfig(opts)
	config.MaxOpenConns = MaxDbConnections()
	// if no search path override passed in as an option, use the viper config
	if config.SearchPathConfig.Empty() {
		config.SearchPathConfig = backend.SearchPathConfig{
			SearchPath:       viper.GetStringSlice(constants.ArgSearchPath),
			SearchPathPrefix: viper.GetStringSlice(constants.ArgSearchPathPrefix),
		}
		slog.Debug("NewDbClient - no search path override passed in, using viper config", "search path", config.SearchPathConfig.String())
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
func (c *DbClient) Close(context.Context) error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}
