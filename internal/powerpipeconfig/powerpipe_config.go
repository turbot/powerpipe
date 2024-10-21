package powerpipeconfig

import (
	"log/slog"
	"strings"
	"sync"

	"github.com/turbot/pipe-fittings/app_specific_connection"
	"github.com/turbot/pipe-fittings/connection"
	"github.com/turbot/powerpipe/internal/constants"
)

type PowerpipeConfig struct {
	// todo only config folder??
	ConfigPaths []string

	PipelingConnections map[string]connection.PipelingConnection

	DefaultConnection connection.ConnectionStringProvider
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
	defaultConnectionName := strings.TrimPrefix(constants.DefaultConnection, "connection.")

	return &PowerpipeConfig{
		PipelingConnections:       defaultPipelingConnections,
		cloudConnectionStringLock: &sync.RWMutex{},
		DefaultConnection:         defaultPipelingConnections[defaultConnectionName].(connection.ConnectionStringProvider),
		cloudConnectionStrings:    make(map[string]string),
	}
}

func (c *PowerpipeConfig) updateResources(other *PowerpipeConfig) {
	c.PipelingConnections = other.PipelingConnections

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
