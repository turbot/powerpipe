package backend

import (
	"database/sql"
	"time"
)

const (
	MaxConnLifeTime = 10 * time.Minute
	MaxConnIdleTime = 1 * time.Minute
)

type connectconfig struct {
	afterConnectFunc func(*sql.Conn) error
	maxConnLifeTime  time.Duration // Add MaxConnLifeTime field
	maxConnIdleTime  time.Duration // Add MaxConnIdleTime field
	maxOpenConns     int           // Add MaxOpenConns field
}

func newConnectConfig(opts []ConnectOption) *connectconfig {
	c := &connectconfig{
		maxConnLifeTime: MaxConnLifeTime,
		maxConnIdleTime: MaxConnIdleTime,
	}
	c.apply(opts)
	return c
}

type ConnectOption func(*connectconfig)

// WithAfterConnect sets the after connect function
func WithAfterConnect(afterConnectFunc func(*sql.Conn) error) ConnectOption {
	return func(c *connectconfig) {
		c.afterConnectFunc = afterConnectFunc
	}
}

func WithConnMaxIdleTime(maxConnIdleTime time.Duration) ConnectOption {
	return func(c *connectconfig) {
		c.maxConnIdleTime = maxConnIdleTime
	}
}
func WithConnMaxLifetime(maxConnLifeTime time.Duration) ConnectOption {
	return func(c *connectconfig) {
		c.maxConnLifeTime = MaxConnLifeTime
	}
}
func WithMaxOpenConns(maxOpenConns int) ConnectOption {
	return func(c *connectconfig) {
		c.maxOpenConns = maxOpenConns
	}
}

func (c *connectconfig) apply(opts []ConnectOption) {
	for _, opt := range opts {
		opt(c)
	}
}
