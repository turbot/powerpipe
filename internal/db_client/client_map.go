package db_client

import (
	"context"
	"sync"

	"github.com/turbot/pipe-fittings/v2/backend"
)

// ClientMapOption defines a function type for configuring ClientMap options
type ClientMapOption func(*ClientMap)

// ClientMap provides thread-safe storage and management of database client connections.
// It caches clients to avoid creating new connections for each query, using composite cache keys
// generated from connection string, search path config, and database filters.
type ClientMap struct {
	// clients stores database client instances with composite cache keys for efficient lookup
	clients map[string]*DbClient
	// clientsMut provides thread-safe access using read-write mutex for concurrent operations
	clientsMut sync.RWMutex
}

// NewClientMap creates a new ClientMap instance with an empty clients map for storing database connections.
func NewClientMap() *ClientMap {
	return &ClientMap{
		clients: make(map[string]*DbClient),
	}
}

// Add stores a database client in the map using the provided configuration.
// Generates a composite cache key from connection string, search path, and filters.
// The filter parameter may be nil. Returns the ClientMap for method chaining.
func (e *ClientMap) Add(client *DbClient, searchPathConfig backend.SearchPathConfig, filter *backend.DatabaseFilters) *ClientMap {
	e.clientsMut.Lock()
	defer e.clientsMut.Unlock()

	// build map key to store client
	key := buildClientMapKey(client.connectionString, searchPathConfig, filter)
	e.clients[key] = client

	return e
}

// Close closes all database clients and removes them from the map for cleanup.
func (e *ClientMap) Close(ctx context.Context) error {
	e.clientsMut.Lock()
	defer e.clientsMut.Unlock()

	for name, client := range e.clients {
		if err := client.Close(ctx); err != nil {
			return err
		}
		delete(e.clients, name)
	}

	return nil
}

// Get retrieves an existing database client from the map based on configuration.
// Generates a composite cache key to lookup the client. The filter parameter may be nil.
// Returns nil if no matching client is found.
func (e *ClientMap) Get(connectionString string, searchPathConfig backend.SearchPathConfig, filter *backend.DatabaseFilters) *DbClient {
	key := buildClientMapKey(connectionString, searchPathConfig, filter)

	// get read lock
	e.clientsMut.RLock()
	client := e.clients[key]
	e.clientsMut.RUnlock()

	return client
}

// GetOrCreate retrieves an existing client or creates a new one if it doesn't exist.
// Generates a composite cache key for lookup and uses double-checked locking to prevent
// race conditions during concurrent access. The filter parameter may be nil.
func (e *ClientMap) GetOrCreate(ctx context.Context, connectionString string, searchPathConfig backend.SearchPathConfig, filter *backend.DatabaseFilters) (*DbClient, error) {
	key := buildClientMapKey(connectionString, searchPathConfig, filter)

	// get read lock
	e.clientsMut.RLock()
	client := e.clients[key]
	e.clientsMut.RUnlock()

	if client != nil {
		return client, nil
	}

	// get write lock
	e.clientsMut.Lock()
	defer e.clientsMut.Unlock()

	// try again (race condition)
	client = e.clients[key]
	if client != nil {
		return client, nil
	}

	// if a search path override was passed in, set the opt
	var opts []backend.BackendOption
	if !searchPathConfig.Empty() {
		opts = append(opts, backend.WithSearchPathConfig(searchPathConfig))
	}

	// create client
	client, err := NewDbClient(ctx, connectionString, opts...)
	if err != nil {
		return nil, err
	}

	// write to map
	e.clients[key] = client

	return client, nil
}

// buildClientMapKey creates a unique composite cache key by combining connection string,
// search path config, and filters. Uses pipe separators to ensure proper key differentiation
// and prevent key collisions between different configurations. The filter parameter may be nil.
func buildClientMapKey(connectionString string, config backend.SearchPathConfig, filter *backend.DatabaseFilters) string {
	key := connectionString + "|" + config.String()
	if filter != nil {
		key += "|" + filter.String()
	}
	return key
}
