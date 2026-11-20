package workspace_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/powerpipe/internal/cmdconfig"
	"github.com/turbot/powerpipe/internal/resources"
	"github.com/turbot/powerpipe/internal/workspace"
)

// the tests use a short quiet period; the settle time is a few times that
const (
	testQuietPeriod = 100 * time.Millisecond
	testSettle      = 600 * time.Millisecond
)

func setupAppConstants() {
	setAppConstants.Do(func() {
		viper.Set("main.version", "0.0.0")
		cmdconfig.SetAppSpecificConstants()
	})
}

// setTimings uses the short quiet period, and the maximum wait if given; the returned function restores them
func setTimings(maxWait ...time.Duration) func() {
	wait := 30 * time.Second
	if len(maxWait) > 0 {
		wait = maxWait[0]
	}
	return workspace.SetModInstallTimings(testQuietPeriod, wait)
}

// startTestWatcher watches a temp workspace folder, counting reloads
func startTestWatcher(t *testing.T, maxWait ...time.Duration) (dir string, reloads *atomic.Int32, m *workspace.ModInstallWatcher) {
	t.Helper()
	setupAppConstants()
	t.Cleanup(setTimings(maxWait...))
	dir = t.TempDir()
	reloads = &atomic.Int32{}
	m, err := workspace.NewModInstallWatcher(dir, func(context.Context) { reloads.Add(1) })
	if err != nil {
		t.Fatal(err)
	}
	m.Start(context.Background())
	t.Cleanup(m.Close)
	return dir, reloads, m
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

// waitForReloads waits until the count reaches want, then waits a settle time to catch any extra reload
func waitForReloads(t *testing.T, reloads *atomic.Int32, want int32) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for reloads.Load() < want && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(testSettle)
	if got := reloads.Load(); got != want {
		t.Fatalf("expected %d reload(s), got %d", want, got)
	}
}

// a mod command, as the installer does it: the mod folders, the lock file, the mod file, then a prune
func fakeModCommand(t *testing.T, dir, version string) {
	t.Helper()
	modDir := filepath.Join(dir, ".powerpipe", "mods", "github.com", "example", "dep@"+version)
	for i := 0; i < 40; i++ {
		writeFile(t, filepath.Join(modDir, fmt.Sprintf("dir%d", i%4), fmt.Sprintf("file%d.pp", i)), "x")
	}
	writeFile(t, filepath.Join(dir, ".mod.cache.json"), version)
	writeFile(t, filepath.Join(dir, "mod.pp"), version)
}

func TestModInstallWatcher_OneReloadPerModCommand(t *testing.T) {
	dir, reloads, _ := startTestWatcher(t)

	// first install: neither .powerpipe nor mods exists yet
	fakeModCommand(t, dir, "v1.0.0")
	waitForReloads(t, reloads, 1)

	// an update: a new version folder (in a mods folder the watcher has not seen before), then the old one pruned
	fakeModCommand(t, dir, "v1.1.0")
	if err := os.RemoveAll(filepath.Join(dir, ".powerpipe", "mods", "github.com", "example", "dep@v1.0.0")); err != nil {
		t.Fatal(err)
	}
	waitForReloads(t, reloads, 2)

	// an uninstall: the mods and the lock file are removed
	if err := os.RemoveAll(filepath.Join(dir, ".powerpipe", "mods")); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(dir, ".mod.cache.json")); err != nil {
		t.Fatal(err)
	}
	waitForReloads(t, reloads, 3)
}

func TestModInstallWatcher_LockFileAloneReloads(t *testing.T) {
	dir, reloads, _ := startTestWatcher(t)
	writeFile(t, filepath.Join(dir, ".mod.cache.json"), "{}")
	waitForReloads(t, reloads, 1)
	// an update that only moves the installed version rewrites the lock file in place
	writeFile(t, filepath.Join(dir, ".mod.cache.json"), `{"a":1}`)
	waitForReloads(t, reloads, 2)
}

func TestModInstallWatcher_ChangeInExistingInstalledModReloads(t *testing.T) {
	setupAppConstants()
	dir := t.TempDir()
	modFile := filepath.Join(dir, ".powerpipe", "mods", "github.com", "example", "dep@v1.0.0", "sub", "mod.pp")
	writeFile(t, modFile, "a")

	defer setTimings()()
	reloads := &atomic.Int32{}
	m, err := workspace.NewModInstallWatcher(dir, func(context.Context) { reloads.Add(1) })
	if err != nil {
		t.Fatal(err)
	}
	m.Start(context.Background())
	defer m.Close()

	writeFile(t, modFile, "b")
	waitForReloads(t, reloads, 1)
}

func TestModInstallWatcher_IgnoresOtherFiles(t *testing.T) {
	dir, reloads, _ := startTestWatcher(t)

	// not the installed mods: the workspace's own files (the pipe-fittings watcher reloads for these),
	// other files in the data folder, and the installer's shadow folder
	writeFile(t, filepath.Join(dir, "mod.pp"), "a")
	writeFile(t, filepath.Join(dir, "README.md"), "a")
	writeFile(t, filepath.Join(dir, ".powerpipe", "internal", "state.json"), "a")
	writeFile(t, filepath.Join(dir, ".powerpipe", ".mods.123", "github.com", "example", "dep@v1.0.0", "mod.pp"), "a")
	writeFile(t, filepath.Join(dir, ".powerpipe", ".mods.123", "other.pp"), "b")
	writeFile(t, filepath.Join(dir, "other", ".mod.cache.json"), "a")

	waitForReloads(t, reloads, 0)
}

func TestModInstallWatcher_ChmodDoesNotReload(t *testing.T) {
	dir, reloads, _ := startTestWatcher(t)
	lock := filepath.Join(dir, ".mod.cache.json")
	// the file is created before the watcher sees it (a reload), then only its mode changes
	writeFile(t, lock, "{}")
	waitForReloads(t, reloads, 1)
	if err := os.Chmod(lock, 0o644); err != nil {
		t.Fatal(err)
	}
	waitForReloads(t, reloads, 1)
}

func TestModInstallWatcher_NoReloadAfterClose(t *testing.T) {
	dir, reloads, m := startTestWatcher(t)

	// a change is pending when the watcher closes
	writeFile(t, filepath.Join(dir, ".mod.cache.json"), "{}")
	time.Sleep(testQuietPeriod / 4)
	m.Close()
	m.Close() // safe to call twice
	waitForReloads(t, reloads, 0)

	// and changes after closing do nothing
	writeFile(t, filepath.Join(dir, ".mod.cache.json"), "{\"a\":1}")
	waitForReloads(t, reloads, 0)
}

func TestModInstallWatcher_ContextCancelStops(t *testing.T) {
	setupAppConstants()
	dir := t.TempDir()
	defer setTimings()()
	reloads := &atomic.Int32{}
	m, err := workspace.NewModInstallWatcher(dir, func(context.Context) { reloads.Add(1) })
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	m.Start(ctx)
	cancel()
	m.Close() // returns once the goroutine has stopped
	writeFile(t, filepath.Join(dir, ".mod.cache.json"), "{}")
	waitForReloads(t, reloads, 0)
}

func TestModInstallWatcher_MaxWaitCapsASteadyStream(t *testing.T) {
	dir, reloads, _ := startTestWatcher(t, 300*time.Millisecond)
	// changes more often than the quiet period never go quiet: the reload must still happen
	stop := time.Now().Add(1500 * time.Millisecond)
	for i := 0; time.Now().Before(stop); i++ {
		writeFile(t, filepath.Join(dir, ".mod.cache.json"), fmt.Sprint(i))
		time.Sleep(testQuietPeriod / 4)
	}
	if reloads.Load() == 0 {
		t.Fatal("expected a reload during a steady stream of changes")
	}
}

// ---- with a real workspace ----

const depModTemplate = `
mod "dep" {
  title = "dep"
}

control "c" {
  title = "%s"
  sql   = "select 1"
}
`

func writeInstalledDep(t *testing.T, dir, version, title string) {
	t.Helper()
	writeFile(t, filepath.Join(dir, ".powerpipe", "mods", "github.com", "example", "dep@"+version, "mod.pp"),
		fmt.Sprintf(depModTemplate, title))
	writeFile(t, filepath.Join(dir, ".mod.cache.json"), fmt.Sprintf(`{
  "local": {
    "github.com/example/dep": {
      "version": %q,
      "name": "github.com/example/dep",
      "commit": "abc",
      "git_ref": "refs/tags/%s",
      "struct_version": 20240429,
      "alias": "dep"
    }
  }
}`, version[1:], version))
}

func controlTitle(w *workspace.PowerpipeWorkspace) string {
	for name, c := range w.GetPowerpipeModResources().Controls {
		if name == "dep.control.c" {
			return *c.Title
		}
	}
	return ""
}

func TestSetupModInstallWatcher_ReloadsWorkspaceAfterUpdate(t *testing.T) {
	setupAppConstants()
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "mod.pp"), `
mod "local" {
  title = "local"
  require {
    mod "github.com/example/dep" {
      version = "*"
    }
  }
}`)
	writeInstalledDep(t, dir, "v1.0.0", "old title")

	ctx := context.Background()
	w, ew := workspace.Load(ctx, dir)
	if ew.GetError() != nil {
		t.Fatal(ew.GetError())
	}
	defer w.Close()
	if got := controlTitle(w); got != "old title" {
		t.Fatalf("expected the installed v1.0.0 to be loaded, got title %q", got)
	}

	type reload struct {
		err  error
		prev modconfig.ModResources
		cur  modconfig.ModResources
	}
	reloaded := make(chan reload, 10)
	w.OnFileWatcherEvent = func(_ context.Context, cur, prev modconfig.ModResources) {
		reloaded <- reload{cur: cur, prev: prev}
	}
	w.OnFileWatcherError = func(_ context.Context, err error) { reloaded <- reload{err: err} }

	defer setTimings()()
	if err := w.SetupModInstallWatcher(ctx); err != nil {
		t.Fatal(err)
	}

	waitReload := func() reload {
		t.Helper()
		select {
		case r := <-reloaded:
			return r
		case <-time.After(10 * time.Second):
			t.Fatal("the workspace did not reload")
			return reload{}
		}
	}
	noMoreReloads := func() {
		t.Helper()
		select {
		case <-reloaded:
			t.Fatal("more than one reload for one mod command")
		case <-time.After(testSettle):
		}
	}

	// `mod update` moves the dependency to a newer version: only the installed mod and the lock file change
	writeInstalledDep(t, dir, "v1.1.0", "new title")
	if err := os.RemoveAll(filepath.Join(dir, ".powerpipe", "mods", "github.com", "example", "dep@v1.0.0")); err != nil {
		t.Fatal(err)
	}
	r := waitReload()
	if r.err != nil {
		t.Fatal(r.err)
	}
	noMoreReloads()
	if got := controlTitle(w); got != "new title" {
		t.Fatalf("expected the workspace to use the updated mod, got title %q", got)
	}
	// the hook is given the resources from before and after, as for a change to a local file
	if prev := r.prev.(*resources.PowerpipeModResources).Controls["dep.control.c"]; prev == nil || *prev.Title != "old title" {
		t.Fatalf("expected the previous resources to hold the old control, got %v", prev)
	}
	if cur := r.cur.(*resources.PowerpipeModResources).Controls["dep.control.c"]; cur == nil || *cur.Title != "new title" {
		t.Fatalf("expected the current resources to hold the new control, got %v", cur)
	}

	// a broken install is reported as an error through the same hook, and the next good one recovers
	writeFile(t, filepath.Join(dir, ".mod.cache.json"), "not json")
	if r := waitReload(); r.err == nil {
		t.Fatal("expected an error for a lock file that cannot be read")
	}
	noMoreReloads()
	writeInstalledDep(t, dir, "v1.2.0", "newest title")
	if r := waitReload(); r.err != nil {
		t.Fatal(r.err)
	}
	noMoreReloads()
	if got := controlTitle(w); got != "newest title" {
		t.Fatalf("expected the workspace to recover, got title %q", got)
	}
}
