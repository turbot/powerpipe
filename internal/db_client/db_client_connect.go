package db_client

import (
	"context"
	"time"

	"github.com/turbot/pipe-fittings/backend"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

const (
	MaxConnLifeTime = 10 * time.Minute
	MaxConnIdleTime = 1 * time.Minute
)

func (c *DbClient) connect(ctx context.Context, opts ...backend.ConnectOption) error {
	utils.LogTime("db_client.establishConnectionPool start")
	defer utils.LogTime("db_client.establishConnectionPool end")

	db, err := c.backend.Connect(ctx, opts...)
	if err != nil {
		return sperr.WrapWithMessage(err, "unable to connect to backend")
	}

	c.db = db
	return nil
}
