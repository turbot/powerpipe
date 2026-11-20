package benchmarkrun_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/turbot/powerpipe/internal/benchmarkrun"
	"github.com/turbot/powerpipe/internal/cmdconfig"
	"github.com/turbot/powerpipe/internal/db_client"
	"github.com/turbot/powerpipe/internal/workspace"
)

const reloadTestMod = `
mod "local" {
  title = "benchmark run reload test"
}

variable "status" {
  type    = string
  default = "ok"
}

benchmark "vars" {
  children = [control.from_var]
}

control "from_var" {
  query = query.status_for
  args  = { status = var.status }
}

query "status_for" {
  sql = "select 'r' as resource, $1 as status, 'from a variable' as reason"
  param "status" {}
}
`

// newWatchedTestRunner returns a runner for the reload test mod, with its workspace watched for changes as
// powerpipe server watches it
func newWatchedTestRunner(t *testing.T) (*benchmarkrun.Runner, *workspace.PowerpipeWorkspace, string) {
	t.Helper()
	setAppConstants.Do(func() {
		viper.Set("main.version", "0.0.0")
		cmdconfig.SetAppSpecificConstants()
	})
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "mod.pp"), []byte(reloadTestMod), 0o600); err != nil {
		t.Fatal(err)
	}
	w, ew := workspace.Load(context.Background(), dir, workspace.WithLateBinding(false))
	if ew.GetError() != nil {
		t.Fatal(ew.GetError())
	}
	if err := w.SetupWatcher(context.Background(), func(context.Context, error) {}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(w.Close)
	client, err := db_client.NewDbClient(context.Background(), "sqlite://"+filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = client.Close(context.Background()) })
	return benchmarkrun.NewRunner(w, client, benchmarkrun.IdleTimeout, benchmarkrun.MaxRunningRuns), w, dir
}

// within fails the test if f takes longer than a second: a request during a reload must not wait for it
func within[T any](t *testing.T, what string, f func() T) T {
	t.Helper()
	done := make(chan T, 1)
	go func() { done <- f() }()
	select {
	case result := <-done:
		return result
	case <-time.After(time.Second):
		t.Fatalf("%s did not answer within a second", what)
		panic("unreachable")
	}
}

// waitUntilReloadHoldsWorkspace waits until a reload has started and holds the workspace, as it does for the whole
// of a slow load
func waitUntilReloadHoldsWorkspace(t *testing.T, w *workspace.PowerpipeWorkspace) {
	t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		release, ok := w.HoldLoaded()
		if !ok {
			return
		}
		release()
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("the reload did not start")
}

// waitUntilLoaded waits until the last load has read the file as it is now
func waitUntilLoaded(t *testing.T, w *workspace.PowerpipeWorkspace, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		s := w.LoadState()
		if !s.Loading && s.LastLoad != nil && s.LastLoad.NewestFileTime != nil && !s.LastLoad.NewestFileTime.Before(info.ModTime()) {
			if s.LastLoad.Error != "" {
				t.Fatalf("the load failed: %s", s.LastLoad.Error)
			}
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("%s did not load, load state %+v", path, w.LoadState())
}

// while the workspace reloads, a run request is answered at once; once the change has loaded, a run uses it
func TestRunner_Reload(t *testing.T) {
	runner, w, dir := newWatchedTestRunner(t)

	first, err := runner.Start("vars", 0)
	if err != nil {
		t.Fatal(err)
	}
	if run := waitForRun(t, runner, first); run.Status != benchmarkrun.StatusComplete || run.Summary.Ok != 1 {
		t.Fatalf("the run before the change: %+v", run)
	}

	// a load that takes as long as the test wants, as a large workspace does
	w.LoadLock()
	locked := true
	t.Cleanup(func() {
		if locked {
			w.LoadUnlock()
		}
	})
	varFile := filepath.Join(dir, "powerpipe.ppvars")
	if err := os.WriteFile(varFile, []byte(`status = "alarm"`+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	waitUntilReloadHoldsWorkspace(t, w)

	var loading *benchmarkrun.WorkspaceLoadingError
	err = within(t, "a run request during a reload", func() error {
		_, err := runner.Start("vars", 0)
		return err
	})
	if !errors.As(err, &loading) || loading.RetryAfter < time.Second || loading.RetryAfter > 30*time.Second {
		t.Fatalf("a run request during a reload: got %v, want a WorkspaceLoadingError with a retry of 1-30s", err)
	}
	// the run URL has the full name of the benchmark, which needs no workspace
	type result struct {
		run benchmarkrun.Run
		err error
	}
	read := within(t, "reading a run during a reload", func() result {
		run, err := runner.Get(first.Benchmark, first.ID)
		return result{run, err}
	})
	if read.err != nil || read.run.ID != first.ID {
		t.Fatalf("reading a run during a reload: got %+v, %v", read.run, read.err)
	}
	// another name of the benchmark needs the workspace to resolve
	err = within(t, "reading a run by a short name during a reload", func() error {
		_, err := runner.Get("vars", first.ID)
		return err
	})
	if !errors.As(err, &loading) {
		t.Fatalf("reading a run by a short name during a reload: got %v, want a WorkspaceLoadingError", err)
	}

	locked = false
	w.LoadUnlock()
	waitUntilLoaded(t, w, varFile)

	second, err := runner.Start("vars", 0)
	if err != nil {
		t.Fatalf("a run after the reload: %v", err)
	}
	if run := waitForRun(t, runner, second); run.Status != benchmarkrun.StatusComplete || run.Summary.Alarm != 1 || run.Summary.Ok != 0 {
		t.Fatalf("the run after the variable changed to alarm: %+v (summary %+v)", run, run.Summary)
	}
	if _, err := runner.Get("vars", first.ID); err != nil {
		t.Fatalf("reading a run by a short name after the reload: %v", err)
	}
}

// waitUntilLoadedSince waits until a load that started after the change has finished
func waitUntilLoadedSince(t *testing.T, w *workspace.PowerpipeWorkspace, changed time.Time) {
	t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		s := w.LoadState()
		if !s.Loading && s.LastLoad != nil && s.LastLoad.StartTime.After(changed) {
			if s.LastLoad.Error != "" {
				t.Fatalf("the load failed: %s", s.LastLoad.Error)
			}
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("no load started after the change, load state %+v", w.LoadState())
}

// a file copied with its old modification time (cp -p, rsync -t, tar), and a deleted file, have loaded once a
// load that started after the change has finished
func TestRunner_ReloadOfFileWithOlderModTime(t *testing.T) {
	runner, w, dir := newWatchedTestRunner(t)
	runStatus := func(what string) *benchmarkrun.Run {
		t.Helper()
		started, err := runner.Start("vars", 0)
		if err != nil {
			t.Fatalf("%s: %v", what, err)
		}
		run := waitForRun(t, runner, started)
		if run.Status != benchmarkrun.StatusComplete {
			t.Fatalf("%s: %+v", what, run)
		}
		return &run
	}

	varFile := filepath.Join(dir, "powerpipe.ppvars")
	if err := os.WriteFile(varFile, []byte(`status = "alarm"`+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	old := time.Now().Add(-time.Hour)
	if err := os.Chtimes(varFile, old, old); err != nil {
		t.Fatal(err)
	}
	changed := time.Now()
	// the newest file time of the first load is already after the file's, so it cannot tell
	if s := w.LoadState(); s.LastLoad.NewestFileTime == nil || s.LastLoad.NewestFileTime.Before(old) {
		t.Fatalf("the first load: %+v", s.LastLoad)
	}
	waitUntilLoadedSince(t, w, changed)
	if run := runStatus("the run after the copied variable file loaded"); run.Summary.Alarm != 1 {
		t.Fatalf("the run after the copied variable file loaded: summary %+v, want an alarm", run.Summary)
	}

	if err := os.Remove(varFile); err != nil {
		t.Fatal(err)
	}
	changed = time.Now()
	waitUntilLoadedSince(t, w, changed)
	if run := runStatus("the run after the variable file was deleted"); run.Summary.Ok != 1 {
		t.Fatalf("the run after the variable file was deleted: summary %+v, want ok", run.Summary)
	}
}
