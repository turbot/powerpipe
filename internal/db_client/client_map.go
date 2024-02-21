package db_client

import (
	"context"
	"github.com/turbot/pipe-fittings/backend"
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

// Clone returns a shallow copy of the client map
func (e *ClientMap) Clone() *ClientMap {
	clients := make(map[string]*DbClient)
	for k, v := range e.clients {
		clients[k] = v
	}
	return &ClientMap{
		clients: clients,
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

	for _, client := range e.clients {
		if err := client.Close(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Get returns a db client for the given connection string
// if clients map already contains a client for this connection string, use that
// otherwise create a new client and add to the map
func (e *ClientMap) Get(ctx context.Context, connectionString string, searchPathConfig backend.SearchPathConfig) (*DbClient, error) {
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
