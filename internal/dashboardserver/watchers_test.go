package dashboardserver

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/v2/constants"
	"github.com/turbot/powerpipe/internal/cmdconfig"
	"github.com/turbot/powerpipe/internal/initialisation"
	"github.com/turbot/powerpipe/internal/workspace"
	"gopkg.in/olahol/melody.v1"
)

const watcherTestMod = `
mod "local" {
  title = "watcher test"
}

variable "status" {
  type    = string
  default = "ok"
}
`

// the file watcher has a short delay and the mod install watcher a quiet period of a second: a reload is looked
// for up to loadWait, and its absence is judged after quietWait
const (
	loadWait  = 10 * time.Second
	quietWait = 3 * time.Second
)

func loadWatcherTestWorkspace(t *testing.T) (*workspace.PowerpipeWorkspace, string) {
	t.Helper()
	viper.Set("main.version", "0.0.0")
	cmdconfig.SetAppSpecificConstants()
	dir := t.TempDir()
	writeTestFile(t, filepath.Join(dir, "mod.pp"), watcherTestMod)
	w, ew := workspace.Load(context.Background(), dir, workspace.WithLateBinding(false))
	if ew.GetError() != nil {
		t.Fatal(ew.GetError())
	}
	t.Cleanup(w.Close)
	return w, dir
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

// waitForLoadCount waits until the workspace has loaded more than loads times, or the wait is over
func waitForLoadCount(w *workspace.PowerpipeWorkspace, loads int, wait time.Duration) int {
	deadline := time.Now().Add(wait)
	for w.LoadState().LoadCount <= loads && time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
	}
	return w.LoadState().LoadCount
}

// changes the workspace as a person editing files does, and as mod install does (the lock file)
func changeWorkspace(t *testing.T, dir string, round string) {
	t.Helper()
	writeTestFile(t, filepath.Join(dir, "mod.pp"), watcherTestMod+"\n# "+round+"\n")
	writeTestFile(t, filepath.Join(dir, "dashboard_"+round+".pp"), "dashboard \"d"+round+"\" {\n  title = \"d\"\n}\n")
	writeTestFile(t, filepath.Join(dir, ".mod.cache.json"), round)
}

// startTestServer starts a server on the workspace with the watch setting given, as the --watch flag or watch in the
// workspace sets it
func startTestServer(t *testing.T, w *workspace.PowerpipeWorkspace, watch bool) {
	t.Helper()
	t.Cleanup(func() { viper.Set(constants.ArgWatch, nil) })
	viper.Set(constants.ArgWatch, watch)
	if _, err := NewServer(context.Background(), &initialisation.InitData{Workspace: w}, melody.New()); err != nil {
		t.Fatal(err)
	}
}

// with watching off nothing changed in the workspace folder reloads the server: not a changed, added or deleted
// file, nor the lock file that mod install writes
func TestNewServer_WatchOff(t *testing.T) {
	w, dir := loadWatcherTestWorkspace(t)
	before := w.LoadState().LoadCount
	startTestServer(t, w, false)

	changeWorkspace(t, dir, "one")
	if got := waitForLoadCount(w, before, quietWait); got != before {
		t.Fatalf("the server reloaded the workspace with watch false: load count %d, was %d", got, before)
	}
	if err := os.Remove(filepath.Join(dir, "dashboard_one.pp")); err != nil {
		t.Fatal(err)
	}
	if got := waitForLoadCount(w, before, quietWait); got != before {
		t.Fatalf("the server reloaded the workspace after a file was deleted with watch false: load count %d, was %d", got, before)
	}
}

// with watching on, a change to the files and one to the installed mods each reload the server, as before
func TestNewServer_WatchOn(t *testing.T) {
	w, dir := loadWatcherTestWorkspace(t)
	before := w.LoadState().LoadCount
	startTestServer(t, w, true)

	writeTestFile(t, filepath.Join(dir, "dashboard_one.pp"), "dashboard \"done\" {\n  title = \"d\"\n}\n")
	afterFile := waitForLoadCount(w, before, loadWait)
	if afterFile <= before {
		t.Fatalf("a changed file did not reload the workspace: load count %d", afterFile)
	}

	// the lock file is what mod install writes, and only the mod install watcher sees it
	time.Sleep(quietWait)
	settled := w.LoadState().LoadCount
	writeTestFile(t, filepath.Join(dir, ".mod.cache.json"), "two")
	if got := waitForLoadCount(w, settled, loadWait); got <= settled {
		t.Fatalf("a mod install did not reload the workspace: load count %d", got)
	}
}
