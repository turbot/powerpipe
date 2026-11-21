package workspace

import (
	"context"
	"log/slog"

	"github.com/fsnotify/fsnotify"
	"github.com/turbot/go-kit/filewatcher"
	"github.com/turbot/pipe-fittings/v2/error_helpers"
	"github.com/turbot/pipe-fittings/v2/modconfig"
)

// SetupWatcher starts the file watcher, which reloads the workspace when one of its mod or variable files changes.
// It replaces the watcher of pipe-fittings, which reloads the workspace the same way but cannot say when a reload
// starts: this one records each load in LoadState, and makes HoldLoaded fail while it reloads.
func (w *PowerpipeWorkspace) SetupWatcher(ctx context.Context, errorHandler func(context.Context, error)) error {
	opts, err := w.watchedFiles(ctx)
	if err != nil {
		return err
	}
	watcher, err := filewatcher.NewWatcher(&filewatcher.WatcherOptions{
		Directories: []string{w.Path},
		Include:     opts.Include,
		Exclude:     opts.Exclude,
		ListFlag:    opts.Flags,
		EventMask:   fsnotify.Create | fsnotify.Remove | fsnotify.Rename | fsnotify.Write,
		// the watcher calls this for one change at a time, after a short delay that gathers the events of a save
		OnChange: func([]fsnotify.Event) {
			w.reload(ctx)
		},
	})
	if err != nil {
		return err
	}
	// set before the watcher starts: a failed reload calls it
	w.FileWatcherErrorHandler = errorHandler
	w.watcher = watcher
	watcher.Start()
	return nil
}

// reload reloads the workspace after a change to its files, and calls the file watcher hooks, as the file watcher
// of pipe-fittings does. Reloads run one at a time: one that is asked for while another runs waits for it.
func (w *PowerpipeWorkspace) reload(ctx context.Context) {
	w.reloadMu.Lock()
	defer w.reloadMu.Unlock()
	w.reloadLocked(ctx)
}

// reloadUnlessLoadedSince reloads the workspace, unless a load that started after the change has finished. started
// is loadsStarted as it was when the change was seen. It waits for any reload in progress first, so a load that
// started after the change and is still running is waited for, not repeated. A load that started before the change
// is no reason to skip, whether it is running or has finished.
func (w *PowerpipeWorkspace) reloadUnlessLoadedSince(ctx context.Context, started int) {
	w.reloadMu.Lock()
	defer w.reloadMu.Unlock()
	if w.LoadState().LoadCount > started {
		slog.Debug("the workspace has already been reloaded since the change")
		return
	}
	w.reloadLocked(ctx)
}

// reloadLocked is reload, with reloadMu held by the caller
func (w *PowerpipeWorkspace) reloadLocked(ctx context.Context) {
	slog.Debug("reloading the workspace after a file change")
	// the load starts before the files are listed, so a change seen after this counts this load as one that
	// started before it (reloadUnlessLoadedSince does not skip for it), and a change seen before it is listed
	w.startLoad(nil)
	// before the load reads the files, so that a change made after this is not reported as loaded
	w.setNewestFileTime(w.newestFileTime(ctx))
	prevModResources, modResources, errAndWarnings := w.reloadModResources(ctx)
	err := errAndWarnings.GetError()
	w.endLoad(err)

	if err != nil {
		if w.OnFileWatcherError != nil {
			w.OnFileWatcherError(ctx, err)
		}
		return
	}
	if w.OnFileWatcherEvent != nil {
		w.OnFileWatcherEvent(ctx, modResources, prevModResources)
	}
}

// reloadModResources reloads the workspace once every holder has released it, and makes HoldLoaded fail until the
// reload ends
func (w *PowerpipeWorkspace) reloadModResources(ctx context.Context) (modconfig.ModResources, modconfig.ModResources, error_helpers.ErrorAndWarnings) {
	w.reloadGate.Lock()
	defer w.reloadGate.Unlock()
	return w.ReloadModResources(ctx)
}
