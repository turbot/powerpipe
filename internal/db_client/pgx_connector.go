package db_client

import (
	"context"
	"database/sql/driver"
	"fmt"
	"github.com/jackc/pgx/v5/stdlib"
)

type PgxConnector struct {
	driver.Connector
	AfterConnectFunc func(context.Context, driver.Conn) error
}

func NewPgxConnector(dataSourceName string, afterConnectFunc func(context.Context, driver.Conn) error) (*PgxConnector, error) {
	driverCtx, ok := stdlib.GetDefaultDriver().(driver.DriverContext)
	if !ok {
		return nil, fmt.Errorf("stdlib driver does not implement DriverContext")
	}
	c, err := driverCtx.OpenConnector(dataSourceName)
	if err != nil {
		return nil, err
	}

	connector := &PgxConnector{
		Connector:        c,
		AfterConnectFunc: afterConnectFunc,
	}

	return connector, nil
}

func (c *PgxConnector) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := c.Connector.Connect(ctx)
	if err != nil {
		return nil, err

	}

	if c.AfterConnectFunc != nil {
		if err = c.AfterConnectFunc(ctx, conn); err != nil {
			return nil, err
		}
	}

	// set search path
	return conn, nil
}
