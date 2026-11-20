package workspace

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/turbot/pipe-fittings/v2/filepaths"
)

// variables, so the tests can shorten them
var (
	// modInstallQuietPeriod is how long the installed mods must stay unchanged before the workspace is reloaded.
	// A mod command writes many files over a short time (the mods, the lock file, the mod file, then prunes),
	// so the reload waits for the command to settle and runs once.
	modInstallQuietPeriod = 1 * time.Second
	// modInstallMaxWait caps how long a steady stream of changes can hold the reload back
	modInstallMaxWait = 30 * time.Second
)

// modInstallWatcher reloads the workspace when the installed mods change.
//
// The pipe-fittings workspace watcher only sees mod and variable files in the workspace folder and skips hidden
// folders, so it never sees `.mod.cache.json` or the folders under `.powerpipe/mods`, which is all that changes
// when `powerpipe mod update` moves a dependency to a newer version.
//
// It watches folders, not files: the workspace folder (for the lock file), the data folder (for `mods` appearing)
// and every folder under the mods folder. Folders that do not exist yet are picked up when they are created.
type modInstallWatcher struct {
	lockPath string
	dataDir  string
	modsDir  string

	quietPeriod time.Duration
	maxWait     time.Duration
	reload      func(ctx context.Context)

	watch     *fsnotify.Watcher
	closeOnce sync.Once
	closed    chan struct{}
	done      chan struct{}
}

func newModInstallWatcher(workspacePath string, reload func(ctx context.Context)) (*modInstallWatcher, error) {
	// a buffer, so the kernel's queue is drained while a reload is running
	watch, err := fsnotify.NewBufferedWatcher(1024)
	if err != nil {
		return nil, err
	}
	modsDir := filepath.Clean(filepaths.WorkspaceModPath(workspacePath))
	m := &modInstallWatcher{
		lockPath:    filepath.Clean(filepaths.WorkspaceLockPath(workspacePath)),
		dataDir:     filepath.Dir(modsDir),
		modsDir:     modsDir,
		quietPeriod: modInstallQuietPeriod,
		maxWait:     modInstallMaxWait,
		reload:      reload,
		watch:       watch,
		closed:      make(chan struct{}),
		done:        make(chan struct{}),
	}
	if err := m.addWatches(filepath.Dir(m.lockPath)); err != nil {
		_ = watch.Close()
		return nil, fmt.Errorf("failed to watch workspace folder %s: %w", filepath.Dir(m.lockPath), err)
	}
	// the data and mods folders may not exist yet: watch what is there
	m.addExisting()
	return m, nil
}

// start runs the watcher until the context is cancelled or Close is called
func (m *modInstallWatcher) start(ctx context.Context) {
	go m.run(ctx)
}

// Close stops the watcher, and waits for any reload in progress to finish. It is safe to call more than once.
func (m *modInstallWatcher) Close() {
	m.closeOnce.Do(func() {
		close(m.closed)
		_ = m.watch.Close()
	})
	<-m.done
}

func (m *modInstallWatcher) run(ctx context.Context) {
	defer close(m.done)

	var (
		firstChange time.Time
		fire        <-chan time.Time
	)
	schedule := func() {
		now := time.Now()
		if fire == nil {
			firstChange = now
		}
		wait := m.quietPeriod
		// do not let a steady stream of changes hold the reload back for ever
		if left := m.maxWait - now.Sub(firstChange); left < wait {
			wait = max(left, 0)
		}
		fire = time.After(wait)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-m.closed:
			return
		case ev, ok := <-m.watch.Events:
			if !ok {
				return
			}
			slog.Debug("mod install watcher event", "event", ev.String())
			if m.isNewFolder(ev) {
				m.addExisting()
			}
			if m.isModChange(ev) {
				schedule()
			}
		case err, ok := <-m.watch.Errors:
			if !ok {
				return
			}
			slog.Warn("mod install watcher error", "error", err)
			if errors.Is(err, fsnotify.ErrEventOverflow) {
				// events were lost: reload to be sure we have not missed a change
				m.addExisting()
				schedule()
			}
		case <-fire:
			fire = nil
			// do not reload once closed or cancelled
			select {
			case <-ctx.Done():
				return
			case <-m.closed:
				return
			default:
			}
			slog.Debug("installed mods changed - reloading workspace")
			m.reload(ctx)
		}
	}
}

// isModChange returns whether the event is a change to the lock file or the installed mods
func (m *modInstallWatcher) isModChange(ev fsnotify.Event) bool {
	// a permissions change is not an install
	if ev.Op&^fsnotify.Chmod == 0 {
		return false
	}
	name := filepath.Clean(ev.Name)
	return name == m.lockPath || name == m.modsDir || isUnder(name, m.modsDir)
}

// isNewFolder returns whether the event may have created a folder we should now watch
func (m *modInstallWatcher) isNewFolder(ev fsnotify.Event) bool {
	if !ev.Has(fsnotify.Create) && !ev.Has(fsnotify.Rename) {
		return false
	}
	name := filepath.Clean(ev.Name)
	return name == m.dataDir || name == m.modsDir || isUnder(name, m.modsDir)
}

// addExisting adds a watch for every folder that exists now: the data folder, the mods folder and all below it
func (m *modInstallWatcher) addExisting() {
	if info, err := os.Stat(m.dataDir); err != nil || !info.IsDir() {
		return
	}
	if err := m.addWatches(m.dataDir); err != nil {
		slog.Warn("mod install watcher could not watch folder", "path", m.dataDir, "error", err)
		return
	}
	if info, err := os.Stat(m.modsDir); err != nil || !info.IsDir() {
		return
	}
	// a mod that cannot be watched (e.g. out of watches) is not fatal: the lock file still changes
	err := filepath.WalkDir(m.modsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// folder removed while we walk: skip it
			return nil
		}
		if !d.IsDir() {
			return nil
		}
		if d.Name() == ".git" {
			return filepath.SkipDir
		}
		if err := m.addWatches(path); err != nil {
			slog.Warn("mod install watcher could not watch folder", "path", path, "error", err)
		}
		return nil
	})
	if err != nil {
		slog.Warn("mod install watcher could not walk the mods folder", "path", m.modsDir, "error", err)
	}
}

// addWatches watches a folder. Watching a folder that is already watched changes nothing.
func (m *modInstallWatcher) addWatches(dir string) error {
	return m.watch.Add(dir)
}

func isUnder(path, dir string) bool {
	return strings.HasPrefix(path, dir+string(filepath.Separator))
}

// SetupModInstallWatcher reloads the workspace when the installed mods change (see modInstallWatcher).
// A failure here leaves the server watching mod and variable files only.
func (w *PowerpipeWorkspace) SetupModInstallWatcher(ctx context.Context) error {
	if w.FileWatcherErrorHandler == nil {
		// ReloadModResources calls this on a failed load
		w.FileWatcherErrorHandler = func(context.Context, error) {}
	}
	m, err := newModInstallWatcher(w.Path, w.reloadAfterModChange)
	if err != nil {
		return err
	}
	if w.modInstallWatcher != nil {
		w.modInstallWatcher.Close()
	}
	w.modInstallWatcher = m
	m.start(ctx)
	return nil
}

// reloadAfterModChange reloads the workspace and raises the same hooks as the pipe-fittings file watcher does
func (w *PowerpipeWorkspace) reloadAfterModChange(ctx context.Context) {
	prevModResources, modResources, ew := w.ReloadModResources(ctx)
	if ew.GetError() != nil {
		if w.OnFileWatcherError != nil {
			w.OnFileWatcherError(ctx, ew.GetError())
		}
		return
	}
	if w.OnFileWatcherEvent != nil {
		w.OnFileWatcherEvent(ctx, modResources, prevModResources)
	}
}
