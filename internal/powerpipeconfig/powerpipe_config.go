package powerpipeconfig

import (
	"github.com/turbot/powerpipe/internal/constants"
	"log/slog"
	"sync"

	"github.com/turbot/pipe-fittings/v2/app_specific_connection"
	"github.com/turbot/pipe-fittings/v2/connection"
)

type PowerpipeConfig struct {
	// todo only config folder??
	ConfigPaths []string

	PipelingConnections map[string]connection.PipelingConnection

	// cache the connection strings for cloud workspaces (is this ok???
	cloudConnectionStrings map[string]string
	// lock
	cloudConnectionStringLock *sync.RWMutex
}

func NewPowerpipeConfig() *PowerpipeConfig {
	defaultPipelingConnections, err := app_specific_connection.DefaultPipelingConnections()
	if err != nil {
		slog.Error("Unable to create default pipeling connections", "error", err)
		return nil
	}

	// populate default connection

	return &PowerpipeConfig{
		PipelingConnections:       defaultPipelingConnections,
		cloudConnectionStringLock: &sync.RWMutex{},

		cloudConnectionStrings: make(map[string]string),
	}
}

func (c *PowerpipeConfig) GetDefaultConnection() connection.ConnectionStringProvider {
	return c.PipelingConnections[constants.DefaultConnection].(connection.ConnectionStringProvider)
}

func (c *PowerpipeConfig) SetDefaultConnection(defaultConnection connection.PipelingConnection) {
	c.PipelingConnections[constants.DefaultConnection] = defaultConnection
}

func (c *PowerpipeConfig) Equals(other *PowerpipeConfig) bool {

	if len(c.PipelingConnections) != len(other.PipelingConnections) {
		return false
	}

	for k, v := range c.PipelingConnections {
		if _, ok := other.PipelingConnections[k]; !ok {
			return false
		}

		if !other.PipelingConnections[k].Equals(v) {
			return false
		}
	}

	return true
}

func (c *PowerpipeConfig) GetCloudConnectionString(workspace string) (string, bool) {
	c.cloudConnectionStringLock.RLock()
	defer c.cloudConnectionStringLock.RUnlock()

	connStr, ok := c.cloudConnectionStrings[workspace]
	return connStr, ok
}

func (c *PowerpipeConfig) SetCloudConnectionString(workspace, connStr string) {
	c.cloudConnectionStringLock.Lock()
	defer c.cloudConnectionStringLock.Unlock()

	c.cloudConnectionStrings[workspace] = connStr
}
