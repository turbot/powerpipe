package workspace_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/powerpipe/internal/cmdconfig"
	"github.com/turbot/powerpipe/internal/workspace"
)

const loadStateTestMod = `
mod "local" {
  title = "load state test"
}

variable "status" {
  type    = string
  default = "ok"
}

benchmark "b" {
  children = [control.c]
}

control "c" {
  title = "first"
  sql   = "select 'r' as resource, 'ok' as status, 'fine' as reason"
}
`

// loadWatchedTestWorkspace writes the files to a new folder, loads it as powerpipe server does, and starts its
// file watcher
func loadWatchedTestWorkspace(t *testing.T, files map[string]string) (*workspace.PowerpipeWorkspace, string) {
	t.Helper()
	return loadWatchedTestWorkspaceWith(t, files, func(*workspace.PowerpipeWorkspace) {})
}

// loadWatchedTestWorkspaceWith is loadWatchedTestWorkspace, calling setup before the watcher starts
func loadWatchedTestWorkspaceWith(t *testing.T, files map[string]string, setup func(*workspace.PowerpipeWorkspace)) (*workspace.PowerpipeWorkspace, string) {
	t.Helper()
	setAppConstants.Do(func() {
		viper.Set("main.version", "0.0.0")
		cmdconfig.SetAppSpecificConstants()
	})
	dir := t.TempDir()
	for name, content := range files {
		writeFile(t, filepath.Join(dir, name), content)
	}
	w, ew := workspace.Load(context.Background(), dir, workspace.WithLateBinding(false))
	if ew.GetError() != nil {
		t.Fatal(ew.GetError())
	}
	setup(w)
	if err := w.SetupWatcher(context.Background(), func(context.Context, error) {}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(w.Close)
	return w, dir
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

// waitForLoadState waits until the load state meets the condition
func waitForLoadState(t *testing.T, w *workspace.PowerpipeWorkspace, what string, condition func(workspace.LoadState) bool) workspace.LoadState {
	t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	for {
		state := w.LoadState()
		if condition(state) {
			return state
		}
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %s, load state %+v", what, state)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// changeTime returns the change time of the file or folder, as a client reads it
func changeTime(t *testing.T, path string) time.Time {
	t.Helper()
	changed, err := workspace.FileChangeTime(path)
	if err != nil {
		t.Fatal(err)
	}
	return changed
}

// loadedFile is true when the last load finished without an error and read the file or folder as it is now, by
// the rule LoadState gives a client
func loadedFile(t *testing.T, path string) func(workspace.LoadState) bool {
	changed := changeTime(t, path)
	return func(s workspace.LoadState) bool {
		return !s.Loading && s.LastLoad != nil && s.LastLoad.Error == "" && s.LastLoad.NewestFileTime != nil &&
			!s.LastLoad.NewestFileTime.Before(changed)
	}
}

func controlTitle(t *testing.T, w *workspace.PowerpipeWorkspace) string {
	t.Helper()
	parsed, err := modconfig.ParseResourceName("control.c")
	if err != nil {
		t.Fatal(err)
	}
	resource, ok := w.GetResource(parsed)
	if !ok {
		t.Fatal("control.c not found")
	}
	return resource.GetTitle()
}

func TestLoadState_FirstLoad(t *testing.T) {
	before := time.Now()
	w, dir := loadWatchedTestWorkspace(t, map[string]string{"mod.pp": loadStateTestMod})

	state := w.LoadState()
	if state.Loading || state.LoadCount != 1 || state.CurrentLoad != nil || state.LastLoad == nil {
		t.Fatalf("after the first load: %+v", state)
	}
	last := state.LastLoad
	if last.StartTime.Before(before) || last.EndTime == nil || last.EndTime.Before(last.StartTime) || last.Error != "" {
		t.Fatalf("the first load: %+v", last)
	}
	if last.NewestFileTime == nil || last.NewestFileTime.Before(changeTime(t, filepath.Join(dir, "mod.pp"))) || last.NewestFileTime.After(last.StartTime) {
		t.Fatalf("the newest file time of the first load is %v, want the change time of mod.pp", last.NewestFileTime)
	}
	// with no load in progress, the workspace can be held
	release, ok := w.HoldLoaded()
	if !ok {
		t.Fatal("the loaded workspace cannot be held")
	}
	release()
}

// each change to a file is a load, whether or not it changes a resource, and the load state says when it has loaded
func TestLoadState_Reload(t *testing.T) {
	w, dir := loadWatchedTestWorkspace(t, map[string]string{"mod.pp": loadStateTestMod})
	modFile := filepath.Join(dir, "mod.pp")

	// a change that changes no resource
	writeFile(t, modFile, loadStateTestMod+"\n# a comment\n")
	state := waitForLoadState(t, w, "the comment to load", loadedFile(t, modFile))
	// a write can raise more than one file event, so more than one load
	if state.LoadCount < 2 {
		t.Fatalf("a change with no new resource: load count %d, want at least 2", state.LoadCount)
	}
	count := state.LoadCount

	// a change to a resource
	writeFile(t, modFile, loadStateTestMod+"\ncontrol \"d\" {\n  sql = \"select 1\"\n}\n")
	state = waitForLoadState(t, w, "the new control to load", loadedFile(t, modFile))
	if state.LoadCount <= count {
		t.Fatalf("after a second change: load count %d, want more than %d", state.LoadCount, count)
	}
	count = state.LoadCount
	if parsed, _ := modconfig.ParseResourceName("control.d"); parsed != nil {
		if _, ok := w.GetResource(parsed); !ok {
			t.Fatal("the load state says the change loaded, but control.d is not in the workspace")
		}
	}

	// a new file, then deleting it: a deletion changes the folder
	extraFile := filepath.Join(dir, "extra.pp")
	writeFile(t, extraFile, "control \"e\" {\n  sql = \"select 1\"\n}\n")
	state = waitForLoadState(t, w, "the new file to load", loadedFile(t, extraFile))
	parsedE, err := modconfig.ParseResourceName("control.e")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := w.GetResource(parsedE); !ok {
		t.Fatal("the load state says the new file loaded, but control.e is not in the workspace")
	}
	if err := os.Remove(extraFile); err != nil {
		t.Fatal(err)
	}
	state = waitForLoadState(t, w, "the deletion to load", loadedFile(t, dir))
	if _, ok := w.GetResource(parsedE); ok {
		t.Fatal("the load state says the deletion loaded, but control.e is still in the workspace")
	}
	count = state.LoadCount

	// a new variable file
	varFile := filepath.Join(dir, "powerpipe.ppvars")
	writeFile(t, varFile, `status = "alarm"`+"\n")
	state = waitForLoadState(t, w, "the variable file to load", loadedFile(t, varFile))
	if state.LoadCount <= count {
		t.Fatalf("after a new variable file: load count %d, want more than %d", state.LoadCount, count)
	}
}

// a load that fails is a finished load with an error, and the workspace stays as it was last loaded
func TestLoadState_FailedReload(t *testing.T) {
	w, dir := loadWatchedTestWorkspace(t, map[string]string{"mod.pp": loadStateTestMod})
	modFile := filepath.Join(dir, "mod.pp")

	writeFile(t, modFile, loadStateTestMod+"\ncontrol \"broken\" {\n")
	state := waitForLoadState(t, w, "the failed load", func(s workspace.LoadState) bool {
		return !s.Loading && s.LoadCount >= 2 && s.LastLoad.Error != "" && !s.LastLoad.NewestFileTime.Before(changeTime(t, modFile))
	})
	if state.LastLoad.EndTime == nil {
		t.Fatalf("a failed load has no end time: %+v", state.LastLoad)
	}
	if title := controlTitle(t, w); title != "first" {
		t.Fatalf("after a failed load the control has title %q, want the one last loaded", title)
	}

	// the next good change loads, and clears the error
	writeFile(t, modFile, loadStateTestMod)
	waitForLoadState(t, w, "the fixed file to load", loadedFile(t, modFile))
}

// while the workspace reloads, the load state says so and HoldLoaded fails at once; a holder delays the reload
func TestLoadState_Loading(t *testing.T) {
	w, dir := loadWatchedTestWorkspace(t, map[string]string{"mod.pp": loadStateTestMod})
	modFile := filepath.Join(dir, "mod.pp")

	// hold the workspace, so that the reload waits
	release, ok := w.HoldLoaded()
	if !ok {
		t.Fatal("the loaded workspace cannot be held")
	}
	writeFile(t, modFile, loadStateTestMod+"\n# change 1\n")
	state := waitForLoadState(t, w, "the reload to start", func(s workspace.LoadState) bool { return s.Loading })
	if state.CurrentLoad == nil || state.LoadCount != 1 || state.LastLoad == nil || state.LastLoad.EndTime == nil {
		t.Fatalf("while loading: %+v", state)
	}
	// the reload has not run: it waits for the holder
	time.Sleep(200 * time.Millisecond)
	if state := w.LoadState(); !state.Loading || state.LoadCount != 1 {
		t.Fatalf("a held workspace reloaded: %+v", state)
	}
	release()
	state = waitForLoadState(t, w, "the reload after the release", loadedFile(t, modFile))
	count := state.LoadCount

	// stop the next reload inside the load itself, as a slow load of a large workspace would be
	w.LoadLock()
	locked := true
	t.Cleanup(func() {
		if locked {
			w.LoadUnlock()
		}
	})
	writeFile(t, modFile, loadStateTestMod+"\n# change 2\n")
	waitForLoadState(t, w, "the second reload to start", func(s workspace.LoadState) bool { return s.Loading })
	deadline := time.Now().Add(10 * time.Second)
	for {
		start := time.Now()
		release, ok := w.HoldLoaded()
		if elapsed := time.Since(start); elapsed > time.Second {
			t.Fatalf("HoldLoaded took %s during a reload", elapsed)
		}
		if !ok {
			break
		}
		// the reload has not taken the gate yet
		release()
		if time.Now().After(deadline) {
			t.Fatal("HoldLoaded still succeeds during a reload")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if retry := w.LoadState().RetryAfter(time.Now()); retry < time.Second || retry > 30*time.Second {
		t.Fatalf("RetryAfter during a reload is %s, want between 1s and 30s", retry)
	}
	locked = false
	w.LoadUnlock()
	state = waitForLoadState(t, w, "the second reload", loadedFile(t, modFile))
	if state.LoadCount <= count {
		t.Fatalf("after the second reload: load count %d, want more than %d", state.LoadCount, count)
	}
	release, ok = w.HoldLoaded()
	if !ok {
		t.Fatal("the workspace cannot be held after the reload")
	}
	release()
}

func TestLoadState_RetryAfter(t *testing.T) {
	now := time.Now()
	end := now.Add(-time.Minute)
	last := &workspace.LoadInfo{StartTime: end.Add(-20 * time.Second), EndTime: &end}
	tests := map[string]struct {
		state workspace.LoadState
		want  time.Duration
	}{
		"no last load":                  {workspace.LoadState{Loading: true, CurrentLoad: &workspace.LoadInfo{StartTime: now}}, time.Second},
		"the last load took 20s":        {workspace.LoadState{Loading: true, CurrentLoad: &workspace.LoadInfo{StartTime: now.Add(-5 * time.Second)}, LastLoad: last}, 15 * time.Second},
		"longer than the last load":     {workspace.LoadState{Loading: true, CurrentLoad: &workspace.LoadInfo{StartTime: now.Add(-25 * time.Second)}, LastLoad: last}, time.Second},
		"the last load took 10 minutes": {workspace.LoadState{Loading: true, CurrentLoad: &workspace.LoadInfo{StartTime: now}, LastLoad: &workspace.LoadInfo{StartTime: end.Add(-10 * time.Minute), EndTime: &end}}, 30 * time.Second},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if got := test.state.RetryAfter(now); got != test.want {
				t.Fatalf("got %s, want %s", got, test.want)
			}
		})
	}
}

// the reload calls the hooks that tell the dashboard server of a change, as the watcher of pipe-fittings did: with
// the resources before and after a change, and with the error of a load that fails
func TestLoadState_ReloadHooks(t *testing.T) {
	type call struct {
		err                     error
		hasDashboard, hadBefore bool
	}
	calls := make(chan call, 100)
	parsed, err := modconfig.ParseResourceName("dashboard.d")
	if err != nil {
		t.Fatal(err)
	}
	_, dir := loadWatchedTestWorkspaceWith(t, map[string]string{"mod.pp": loadStateTestMod}, func(w *workspace.PowerpipeWorkspace) {
		onEvent, onError := w.OnFileWatcherEvent, w.OnFileWatcherError
		w.OnFileWatcherEvent = func(ctx context.Context, modResources, prevModResources modconfig.ModResources) {
			_, has := modResources.GetResource(parsed)
			_, had := prevModResources.GetResource(parsed)
			calls <- call{hasDashboard: has, hadBefore: had}
			onEvent(ctx, modResources, prevModResources)
		}
		w.OnFileWatcherError = func(ctx context.Context, err error) {
			calls <- call{err: err}
			onError(ctx, err)
		}
	})
	modFile := filepath.Join(dir, "mod.pp")
	waitForCall := func(what string, match func(call) bool) {
		t.Helper()
		timeout := time.After(30 * time.Second)
		for {
			select {
			case c := <-calls:
				if match(c) {
					return
				}
			case <-timeout:
				t.Fatalf("no %s", what)
			}
		}
	}

	writeFile(t, modFile, loadStateTestMod+"\ndashboard \"d\" {\n  title = \"new\"\n}\n")
	waitForCall("change hook with the new dashboard", func(c call) bool { return c.err == nil && c.hasDashboard && !c.hadBefore })

	writeFile(t, modFile, loadStateTestMod+"\ncontrol \"broken\" {\n")
	waitForCall("error hook", func(c call) bool { return c.err != nil })
}
