package powerpipeconfig

import (
	"context"
	"github.com/turbot/pipe-fittings/app_specific"
	"github.com/turbot/powerpipe/internal/constants"
	"log/slog"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	filehelpers "github.com/turbot/go-kit/files"
	"github.com/turbot/go-kit/filewatcher"
	"github.com/turbot/pipe-fittings/app_specific_connection"
	"github.com/turbot/pipe-fittings/connection"
)

type PowerpipeConfig struct {
	// todo only config folder??
	ConfigPaths []string

	PipelingConnections map[string]connection.PipelingConnection

	// TODO KAI do we need file watching
	watcher                 *filewatcher.FileWatcher
	fileWatcherErrorHandler func(context.Context, error)

	// Hooks
	OnFileWatcherError func(context.Context, error)
	OnFileWatcherEvent func(context.Context, *PowerpipeConfig)

	loadLock          *sync.Mutex
	DefaultConnection connection.ConnectionStringProvider
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
		PipelingConnections: defaultPipelingConnections,
		loadLock:            &sync.Mutex{},
		DefaultConnection:   defaultPipelingConnections[defaultConnectionName].(connection.ConnectionStringProvider),
	}
}

func (c *PowerpipeConfig) updateResources(other *PowerpipeConfig) {
	c.loadLock.Lock()
	defer c.loadLock.Unlock()

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

func (c *PowerpipeConfig) SetupWatcher(ctx context.Context, errorHandler func(context.Context, error)) error {
	watcherOptions := &filewatcher.WatcherOptions{
		Directories: c.ConfigPaths,
		Include:     filehelpers.InclusionsFromExtensions([]string{app_specific.ConfigExtension}),
		ListFlag:    filehelpers.FilesRecursive,
		EventMask:   fsnotify.Create | fsnotify.Remove | fsnotify.Rename | fsnotify.Write,
		// we should look into passing the callback function into the underlying watcher
		// we need to analyze the kind of errors that come out from the watcher and
		// decide how to handle them
		// OnError: errCallback,
		OnChange: func(events []fsnotify.Event) {
			c.handleFileWatcherEvent(ctx)
		},
	}
	watcher, err := filewatcher.NewWatcher(watcherOptions)
	if err != nil {
		return err
	}
	c.watcher = watcher

	// start the watcher
	watcher.Start()

	// set the file watcher error handler, which will get called when there are parsing errors
	// after a file watcher event
	c.fileWatcherErrorHandler = errorHandler

	return nil
}

func (c *PowerpipeConfig) handleFileWatcherEvent(ctx context.Context) {
	slog.Debug("PowerpipeConfig handleFileWatcherEvent")

	newConfig, errAndWarnings := LoadPowerpipeConfig(c.ConfigPaths...)

	if errAndWarnings.GetError() != nil {
		// call error hook
		if c.OnFileWatcherError != nil {
			c.OnFileWatcherError(ctx, errAndWarnings.Error)
		}

		// Flag on workspace?
		return
	}

	if !newConfig.Equals(c) {
		c.updateResources(newConfig)

		// call hook
		if c.OnFileWatcherEvent != nil {
			c.OnFileWatcherEvent(ctx, newConfig)
		}
	}

}
