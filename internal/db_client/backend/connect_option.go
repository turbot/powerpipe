package backend

import (
	"time"
)

const (
	MaxConnLifeTime = 10 * time.Minute
	MaxConnIdleTime = 1 * time.Minute
)

type poolConfig struct {
	maxConnLifeTime time.Duration
	maxConnIdleTime time.Duration
	maxOpenConns    int
}

type searchPathConfig struct {
	searchPath       string
	searchPathPrefix string
}

type connectconfig struct {
	poolConfig       *poolConfig
	searchPathConfig *searchPathConfig
}

func newConnectConfig(opts []ConnectOption) *connectconfig {
	c := &connectconfig{
		poolConfig: &poolConfig{
			maxConnLifeTime: MaxConnLifeTime,
			maxConnIdleTime: MaxConnIdleTime,
			maxOpenConns:    0,
		},
		searchPathConfig: &searchPathConfig{
			searchPath:       "",
			searchPathPrefix: "",
		},
	}
	c.apply(opts)
	return c
}

type ConnectOption func(*connectconfig)

func WithPoolConfig(config *poolConfig) ConnectOption {
	return func(c *connectconfig) {
		c.poolConfig = config
	}
}

func WithSearchPathConfig(config *searchPathConfig) ConnectOption {
	return func(c *connectconfig) {
		c.searchPathConfig = config
	}
}

func (c *connectconfig) apply(opts []ConnectOption) {
	for _, opt := range opts {
		opt(c)
	}
}
