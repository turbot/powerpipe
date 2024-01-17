package backend

import (
	"time"
)

const (
	DefaultMaxConnLifeTime  = 10 * time.Minute
	DefaultMaxConnIdleTime  = 1 * time.Minute
	DefaultMaxOpenConns     = 10
	DefaultSearchPath       = ""
	DefaultSearchPathPrefix = ""
)

type PoolConfig struct {
	MaxConnLifeTime time.Duration
	MaxConnIdleTime time.Duration
	MaxOpenConns    int
}

type SearchPathConfig struct {
	SearchPath       string
	SearchPathPrefix string
}

type ConnectConfig struct {
	PoolConfig       *PoolConfig
	SearchPathConfig *SearchPathConfig
}

func newConnectConfig(opts []ConnectOption) *ConnectConfig {
	c := &ConnectConfig{
		PoolConfig: &PoolConfig{
			MaxConnLifeTime: DefaultMaxConnLifeTime,
			MaxConnIdleTime: DefaultMaxConnIdleTime,
			MaxOpenConns:    DefaultMaxOpenConns,
		},
		SearchPathConfig: &SearchPathConfig{
			SearchPath:       DefaultSearchPath,
			SearchPathPrefix: DefaultSearchPathPrefix,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type ConnectOption func(*ConnectConfig)

func WithPoolConfig(config *PoolConfig) ConnectOption {
	return func(c *ConnectConfig) {
		c.PoolConfig = config
	}
}

// WithSearchPathConfig sets the search path to use when connecting to the database.
// If a prefix is also set, the search path will be resolved to the first matching
// schema in the search path. Only applies if the backend is postgres
func WithSearchPathConfig(config *SearchPathConfig) ConnectOption {
	return func(c *ConnectConfig) {
		c.SearchPathConfig = config
	}
}
