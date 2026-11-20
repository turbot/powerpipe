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
	// loadIdentity returns a value that changes each time the workspace is reloaded, by anything. It is called
	// when a change is seen, which waits for any reload in progress to finish. If it has changed when the
	// quiet period ends, a reload that started after the last change has already happened, so there is nothing
	// to reload. `mod install` and `mod uninstall` write the mod file, which the pipe-fittings watcher reloads for,
	// so without this the one command would reload twice. Nil if there is no other watcher.
	loadIdentity func() any

	watch     *fsnotify.Watcher
	closeOnce sync.Once
	closed    chan struct{}
	done      chan struct{}
}

func newModInstallWatcher(workspacePath string, reload func(ctx context.Context), loadIdentity func() any) (*modInstallWatcher, error) {
	// a buffer, so the kernel's queue is drained while a reload is running
	watch, err := fsnotify.NewBufferedWatcher(1024)
	if err != nil {
		return nil, err
	}
	modsDir := filepath.Clean(filepaths.WorkspaceModPath(workspacePath))
	m := &modInstallWatcher{
		lockPath:     filepath.Clean(filepaths.WorkspaceLockPath(workspacePath)),
		dataDir:      filepath.Dir(modsDir),
		modsDir:      modsDir,
		quietPeriod:  modInstallQuietPeriod,
		maxWait:      modInstallMaxWait,
		reload:       reload,
		loadIdentity: loadIdentity,
		watch:        watch,
		closed:       make(chan struct{}),
		done:         make(chan struct{}),
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
		// whether there is a change waiting to be reloaded for
		pending bool
		// the load identity as of the last change
		loadedAtChange any
	)
	noteChange := func() {
		if m.loadIdentity != nil {
			loadedAtChange = m.loadIdentity()
		}
	}
	schedule := func() {
		now := time.Now()
		if !pending {
			firstChange = now
			pending = true
		}
		wait := m.quietPeriod
		// do not let a steady stream of changes hold the reload back for ever
		if left := m.maxWait - now.Sub(firstChange); left < wait {
			wait = max(left, 0)
		}
		fire = time.After(wait)
	}

	// handle returns whether the event is a change to the installed mods
	handle := func(ev fsnotify.Event) bool {
		slog.Debug("mod install watcher event", "event", ev.String())
		m.watchNewFolder(ev)
		if !m.isModChange(ev) {
			return false
		}
		noteChange()
		schedule()
		return true
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
			handle(ev)
		case err, ok := <-m.watch.Errors:
			if !ok {
				return
			}
			slog.Warn("mod install watcher error", "error", err)
			if errors.Is(err, fsnotify.ErrEventOverflow) {
				// events were lost: reload to be sure we have not missed a change
				m.addExisting()
				noteChange()
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
			// changes already queued have not been seen yet: take them in, and wait again if there are any
			if m.drain(handle) {
				continue
			}
			pending = false
			if m.loadIdentity != nil && m.loadIdentity() != loadedAtChange {
				slog.Debug("installed mods changed - workspace has already been reloaded since")
				continue
			}
			slog.Debug("installed mods changed - reloading workspace")
			m.reload(ctx)
		}
	}
}

// isModChange returns whether the event is a change to the lock file or the installed mods.
//
// Only files and folders being created or written count under the mods folder, not their removal: the lock file
// changes when a mod stops being used, and the installer then prunes the mod's folder, which changes nothing the
// workspace loads (it would otherwise be a change after the mod file's reload).
func (m *modInstallWatcher) isModChange(ev fsnotify.Event) bool {
	// a permissions change is not an install
	if ev.Op&^fsnotify.Chmod == 0 {
		return false
	}
	name := filepath.Clean(ev.Name)
	if name == m.lockPath {
		return true
	}
	if name == m.modsDir || isUnder(name, m.modsDir) {
		return ev.Has(fsnotify.Create) || ev.Has(fsnotify.Write)
	}
	return false
}

// drain handles the events that are already waiting, and returns whether any was a change to the installed mods
func (m *modInstallWatcher) drain(handle func(fsnotify.Event) bool) (changed bool) {
	for {
		select {
		case ev, ok := <-m.watch.Events:
			if !ok {
				return changed
			}
			if handle(ev) {
				changed = true
			}
		default:
			return changed
		}
	}
}

// watchNewFolder watches a folder that has just been created in the data or mods folder, and the folders in it
func (m *modInstallWatcher) watchNewFolder(ev fsnotify.Event) {
	if !ev.Has(fsnotify.Create) && !ev.Has(fsnotify.Rename) {
		return
	}
	name := filepath.Clean(ev.Name)
	switch {
	case name == m.dataDir:
		// the mods folder may be in it already
		m.addExisting()
	case name == m.modsDir || isUnder(name, m.modsDir):
		// only this folder's tree: walking the whole of the mods folder for each new folder in it
		// would leave this loop behind a large install
		m.addTree(name)
	}
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
	m.addTree(m.modsDir)
}

// addTree watches a folder and the folders below it. A folder that cannot be watched (for example, out of watches)
// is not fatal: the lock file still changes.
func (m *modInstallWatcher) addTree(root string) {
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// removed while we walk: skip it
			slog.Debug("mod install watcher skipped a folder it could not read", "path", path, "error", err)
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
		slog.Warn("mod install watcher could not walk folder", "path", root, "error", err)
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
	m, err := newModInstallWatcher(w.Path, w.reloadAfterModChange, w.loadIdentity)
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

// loadIdentity returns the loaded mod, which is a new one after every successful reload
func (w *PowerpipeWorkspace) loadIdentity() any {
	w.LoadLock()
	defer w.LoadUnlock()
	return w.Mod
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
