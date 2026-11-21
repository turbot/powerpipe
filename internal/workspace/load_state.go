package workspace

import (
	"context"
	"log/slog"
	"os"
	"time"

	filehelpers "github.com/turbot/go-kit/files"
	"github.com/turbot/pipe-fittings/v2/filepaths"
	"github.com/turbot/pipe-fittings/v2/load_mod"
)

// LoadState says whether the workspace is loading, and how its loads went. powerpipe server reloads the workspace
// when one of its mod or variable files changes, and when the installed mods change (the lock file or the folders
// under the mods folder), unless watching is off (--watch=false): then it loads once.
//
// A client that changed, added or deleted files can wait for the change to load. After its last change it reads
// the change time of what it changed with FileChangeTime (ctime: stat -c %z on Linux, stat -f %Fc on macOS): the
// file it wrote, or the folder of a file it deleted. The change has loaded once LastLoad has a NewestFileTime at or
// after that time, and no Error.
type LoadState struct {
	// Loading is true while the workspace loads
	Loading bool `json:"loading"`
	// LoadCount is the number of loads that have finished, including the first load and loads that failed.
	// A load that changes no resource counts too.
	LoadCount int `json:"load_count"`
	// CurrentLoad is the load in progress, if any
	CurrentLoad *LoadInfo `json:"current_load,omitempty"`
	// LastLoad is the last load that finished, if any
	LastLoad *LoadInfo `json:"last_load,omitempty"`
}

// LoadInfo describes one load of the workspace
type LoadInfo struct {
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	// NewestFileTime is the newest change time (see FileChangeTime) of the mod and variable files of the workspace,
	// of its folders, and of its lock file (.mod.cache.json), when the load started. The load read every change made up to then.
	NewestFileTime *time.Time `json:"newest_file_time,omitempty"`
	// Error is why the load failed. After a failed load the workspace stays as it was last loaded.
	Error string `json:"error,omitempty"`
}

func (l *LoadInfo) clone() *LoadInfo {
	if l == nil {
		return nil
	}
	c := *l
	if l.EndTime != nil {
		end := *l.EndTime
		c.EndTime = &end
	}
	if l.NewestFileTime != nil {
		newest := *l.NewestFileTime
		c.NewestFileTime = &newest
	}
	return &c
}

// Minimum and maximum of the wait that RetryAfter suggests
const (
	minRetryAfter = time.Second
	maxRetryAfter = 30 * time.Second
)

// RetryAfter is how long a client should wait before it asks again for something the load in progress holds up:
// the time the last load took, less the time the current load has taken, between 1 and 30 seconds
func (s LoadState) RetryAfter(now time.Time) time.Duration {
	wait := minRetryAfter
	if s.CurrentLoad != nil && s.LastLoad != nil && s.LastLoad.EndTime != nil {
		wait = s.LastLoad.EndTime.Sub(s.LastLoad.StartTime) - now.Sub(s.CurrentLoad.StartTime)
	}
	return min(max(wait, minRetryAfter), maxRetryAfter)
}

// LoadState returns the load state of the workspace
func (w *PowerpipeWorkspace) LoadState() LoadState {
	w.loadStateMu.Lock()
	defer w.loadStateMu.Unlock()
	return LoadState{
		Loading:     w.loadState.Loading,
		LoadCount:   w.loadState.LoadCount,
		CurrentLoad: w.loadState.CurrentLoad.clone(),
		LastLoad:    w.loadState.LastLoad.clone(),
	}
}

// loadsStarted returns the number of loads that have started: those that have finished, and the one in progress
func (w *PowerpipeWorkspace) loadsStarted() int {
	w.loadStateMu.Lock()
	defer w.loadStateMu.Unlock()
	started := w.loadState.LoadCount
	if w.loadState.Loading {
		started++
	}
	return started
}

// HoldLoaded holds the workspace as it is loaded, so that a reload waits to start until release is called. It
// returns false, and holds nothing, while the workspace reloads or a reload waits to start, so that a caller can
// answer at once rather than wait for the reload. Hold it only briefly: a reload waits for every holder.
func (w *PowerpipeWorkspace) HoldLoaded() (release func(), ok bool) {
	if !w.reloadGate.TryRLock() {
		return nil, false
	}
	return w.reloadGate.RUnlock, true
}

// startLoad records the start of a load. newestFileTime is the newest change time of the files and folders of the
// workspace, read before the load reads them.
func (w *PowerpipeWorkspace) startLoad(newestFileTime *time.Time) {
	w.loadStateMu.Lock()
	defer w.loadStateMu.Unlock()
	w.loadState.Loading = true
	w.loadState.CurrentLoad = &LoadInfo{StartTime: time.Now(), NewestFileTime: newestFileTime}
}

// setNewestFileTime sets the newest change time of the load in progress
func (w *PowerpipeWorkspace) setNewestFileTime(newestFileTime *time.Time) {
	w.loadStateMu.Lock()
	defer w.loadStateMu.Unlock()
	if w.loadState.CurrentLoad != nil {
		w.loadState.CurrentLoad.NewestFileTime = newestFileTime
	}
}

// endLoad records the end of the load in progress
func (w *PowerpipeWorkspace) endLoad(err error) {
	w.loadStateMu.Lock()
	defer w.loadStateMu.Unlock()
	load := w.loadState.CurrentLoad
	if load == nil {
		// endLoad always follows startLoad
		return
	}
	end := time.Now()
	load.EndTime = &end
	if err != nil {
		load.Error = err.Error()
	}
	w.loadState.Loading = false
	w.loadState.CurrentLoad = nil
	w.loadState.LastLoad = load
	w.loadState.LoadCount++
}

// watchedFiles returns the options that list the files the file watcher watches: the mod and variable files of
// the workspace, less its exclusions
func (w *PowerpipeWorkspace) watchedFiles(ctx context.Context) (*filehelpers.ListOptions, error) {
	if w.watchedFileOptions != nil {
		return w.watchedFileOptions, nil
	}
	// the workspace does not export its exclusions, which its parse context lists the files with
	parseCtx, err := w.GetParseContext(ctx)
	if err != nil {
		return nil, err
	}
	w.watchedFileOptions = &filehelpers.ListOptions{
		Flags:   filehelpers.FilesRecursive,
		Include: filehelpers.InclusionsFromExtensions(load_mod.GetModFileExtensions()),
		Exclude: parseCtx.ListOptions.Exclude,
	}
	return w.watchedFileOptions, nil
}

// FileChangeTime returns the time the file or folder last changed, as NewestFileTime counts it: its status
// change time (ctime), which a write, a rename or a deletion in a folder moves to the present. A copy that keeps
// the modification time of its source (cp -p, rsync -t, tar) cannot set it, as it can set the modification time.
// On a system with no status change time, it is the modification time.
func FileChangeTime(path string) (time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return changeTime(info), nil
}

// newestFileTime returns the newest change time of the files the file watcher watches, of their folders, and of the
// lock file, or nil if they cannot be listed. A folder changes when a file in it is added, renamed or deleted. The
// lock file is the one file `powerpipe mod install`, `update` and `uninstall` always write, and it is not a mod file.
func (w *PowerpipeWorkspace) newestFileTime(ctx context.Context) *time.Time {
	opts, err := w.watchedFiles(ctx)
	if err != nil {
		slog.Debug("cannot list the files of the workspace", "error", err)
		return nil
	}
	files, err := filehelpers.ListFilesWithContext(ctx, w.Path, opts)
	if err != nil {
		slog.Debug("cannot list the files of the workspace", "error", err)
		return nil
	}
	// the folders the file watcher watches
	folders, err := filehelpers.ListFilesWithContext(ctx, w.Path, &filehelpers.ListOptions{
		Flags:   filehelpers.DirectoriesRecursive,
		Exclude: opts.Exclude,
	})
	if err != nil {
		slog.Debug("cannot list the folders of the workspace", "error", err)
		return nil
	}
	var newest time.Time
	paths := append(append(files, w.Path, filepaths.WorkspaceLockPath(w.Path)), folders...)
	for _, path := range paths {
		changed, err := FileChangeTime(path)
		if err != nil {
			// removed since it was listed
			continue
		}
		if changed.After(newest) {
			newest = changed
		}
	}
	if newest.IsZero() {
		return nil
	}
	return &newest
}
