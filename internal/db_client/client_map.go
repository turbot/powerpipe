package db_client

import (
	"context"
	"github.com/turbot/pipe-fittings/v2/backend"
	"sync"
)

type ClientMapOption func(*ClientMap)

type ClientMap struct {
	clients    map[string]*DbClient
	clientsMut sync.RWMutex
}

func NewClientMap() *ClientMap {
	return &ClientMap{
		clients: make(map[string]*DbClient),
	}
}

func (e *ClientMap) Add(client *DbClient, searchPathConfig backend.SearchPathConfig) *ClientMap {
	e.clientsMut.Lock()
	defer e.clientsMut.Unlock()

	// build map key to store client
	key := buildClientMapKey(client.connectionString, searchPathConfig)
	e.clients[key] = client

	e.clients[client.connectionString] = client
	return e
}

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

// Get returns an existing db client for the given connection string
// if no client is stored for the string, it returns null
func (e *ClientMap) Get(connectionString string, searchPathConfig backend.SearchPathConfig) *DbClient {
	key := buildClientMapKey(connectionString, searchPathConfig)

	// get read lock
	e.clientsMut.RLock()
	client := e.clients[key]
	e.clientsMut.RUnlock()

	return client
}

// GetOrCreate returns a db client for the given connection string
// if clients map already contains a client for this connection string, use that
// otherwise create a new client and add to the map
func (e *ClientMap) GetOrCreate(ctx context.Context, connectionString string, searchPathConfig backend.SearchPathConfig) (*DbClient, error) {
	key := buildClientMapKey(connectionString, searchPathConfig)

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
	var opts []backend.ConnectOption
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

func buildClientMapKey(connectionString string, config backend.SearchPathConfig) string {
	return connectionString + config.String()
}
