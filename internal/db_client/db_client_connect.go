package db_client

import (
	"context"
	"time"

	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/powerpipe/internal/db_client/backend"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

const (
	MaxConnLifeTime = 10 * time.Minute
	MaxConnIdleTime = 1 * time.Minute
)

func (c *DbClient) connect(ctx context.Context) error {
	utils.LogTime("db_client.establishConnectionPool start")
	defer utils.LogTime("db_client.establishConnectionPool end")

	poolConfig := backend.PoolConfig{
		MaxConnIdleTime: MaxConnIdleTime,
		MaxConnLifeTime: MaxConnLifeTime,
		MaxOpenConns:    MaxDbConnections(),
	}

	searchPathConfig := backend.SearchPathConfig{
		SearchPath:       viper.GetStringSlice(constants.ArgSearchPath),
		SearchPathPrefix: viper.GetStringSlice(constants.ArgSearchPathPrefix),
	}

	db, err := c.backend.Connect(ctx, backend.WithPoolConfig(poolConfig), backend.WithSearchPathConfig(searchPathConfig))
	if err != nil {
		return sperr.WrapWithMessage(err, "unable to connect to backend")
	}

	c.db = db
	return nil
}
