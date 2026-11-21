package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/turbot/powerpipe/internal/workspace"
)

// the server's own quiet period applies: these tests use the real watchers, as the server does

const modInstallDepMod = `
mod "dep" {
  title = "dep"
}

control "c" {
  title = "%s"
  sql   = "select 1"
}
`

// the test mod with a require, as `powerpipe mod install` writes it
var modWithRequire = strings.Replace(benchmarkRunTestMod, "  title = \"benchmark run API test\"\n}",
	"  title = \"benchmark run API test\"\n  require {\n    mod \"github.com/example/dep\" {\n      version = \"*\"\n    }\n  }\n}", 1)

func init() {
	if modWithRequire == benchmarkRunTestMod {
		panic("the test mod has no require added")
	}
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

// installTestDep writes what a mod command writes: the mod's folder and the lock file (an install and an uninstall
// also change the mod file)
func installTestDep(t *testing.T, dir, version string) {
	t.Helper()
	writeTestFile(t, filepath.Join(dir, ".powerpipe", "mods", "github.com", "example", "dep@"+version, "mod.pp"),
		fmt.Sprintf(modInstallDepMod, version))
	writeTestFile(t, filepath.Join(dir, ".mod.cache.json"), fmt.Sprintf(`{
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

func removeTestDep(t *testing.T, dir, version string) {
	t.Helper()
	if err := os.RemoveAll(filepath.Join(dir, ".powerpipe", "mods", "github.com", "example", "dep@"+version)); err != nil {
		t.Fatal(err)
	}
}

// waitForState polls the load state until the condition holds
func waitForState(t *testing.T, router *gin.Engine, what string, condition func(workspace.LoadState) bool) workspace.LoadState {
	t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	for {
		state := getLoadState(t, router)
		if condition(state) {
			return state
		}
		if time.Now().After(deadline) {
			t.Fatalf("%s: the load state is %+v", what, state)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// modCommand makes the changes of a mod command while a load is held up, as a large workspace's load is, and checks
// the server through the reload that follows: the load state shows it loading, a run request is answered with 503
// and Retry-After, and once it ends there is exactly one more load, which has read the lock file.
func modCommand(t *testing.T, router *gin.Engine, w *workspace.PowerpipeWorkspace, dir, what string, changes ...func()) {
	t.Helper()
	before := getLoadState(t, router)
	if before.Loading {
		t.Fatalf("%s: the workspace is loading before the command: %+v", what, before)
	}

	w.LoadLock()
	locked := true
	unlock := func() {
		if locked {
			locked = false
			w.LoadUnlock()
		}
	}
	t.Cleanup(unlock)
	for _, change := range changes {
		change()
	}

	loading := waitForState(t, router, what+": the reload did not start", func(s workspace.LoadState) bool { return s.Loading })
	if loading.LoadCount != before.LoadCount {
		t.Fatalf("%s: the load count rose before the reload ended: %+v", what, loading)
	}
	// a reload holds the workspace for the whole of its load
	deadline := time.Now().Add(30 * time.Second)
	for {
		release, ok := w.HoldLoaded()
		if !ok {
			break
		}
		release()
		if time.Now().After(deadline) {
			t.Fatalf("%s: the reload did not take the workspace", what)
		}
		time.Sleep(10 * time.Millisecond)
	}
	rec := serve(t, router, http.MethodPost, "/api/v0/benchmark/quick/run")
	if rec.Code != http.StatusServiceUnavailable || !strings.Contains(rec.Body.String(), "the workspace is loading") {
		t.Fatalf("%s: a run request during the reload: got %d %s, want 503", what, rec.Code, rec.Body)
	}
	if retry, err := strconv.Atoi(rec.Header().Get("Retry-After")); err != nil || retry < 1 || retry > 30 {
		t.Fatalf("%s: a run request during the reload: Retry-After %q, want 1-30 seconds", what, rec.Header().Get("Retry-After"))
	}

	unlock()
	// the command's last change was a write, and the lock file is the one a client reads the change time of
	state := waitForState(t, router, what+": the reload did not end", func(s workspace.LoadState) bool {
		return !s.Loading && s.LoadCount > before.LoadCount
	})
	if state.LastLoad == nil || state.LastLoad.Error != "" {
		t.Fatalf("%s: the reload failed: %+v", what, state)
	}
	if lock, err := workspace.FileChangeTime(filepath.Join(dir, ".mod.cache.json")); err != nil {
		t.Fatal(err)
	} else if state.LastLoad.NewestFileTime == nil || state.LastLoad.NewestFileTime.Before(lock) {
		t.Fatalf("%s: the load's newest file time %v is before the change time of the lock file %v", what, state.LastLoad.NewestFileTime, lock)
	}

	// longer than the quiet period of the mod install watcher, which would reload again if it were going to
	time.Sleep(2500 * time.Millisecond)
	if after := getLoadState(t, router); after.Loading || after.LoadCount != before.LoadCount+1 {
		t.Fatalf("%s: expected exactly one load more than the %d before, got %+v", what, before.LoadCount, after)
	}
}

// after each of `mod install` (which also adds a require to the mod file, so both watchers see a change), `mod update`
// and `mod uninstall`, the server reloads once, through the same load state and gate as any other reload
func TestWorkspaceAPI_ReloadAfterModCommands(t *testing.T) {
	router, w, dir := newWatchedTestRouterWith(t, true)
	modFile := filepath.Join(dir, "mod.pp")

	modCommand(t, router, w, dir, "install", func() {
		installTestDep(t, dir, "v1.0.0")
		writeTestFile(t, modFile, modWithRequire)
	})
	modCommand(t, router, w, dir, "update", func() {
		installTestDep(t, dir, "v1.1.0")
		removeTestDep(t, dir, "v1.0.0")
	})
	modCommand(t, router, w, dir, "uninstall", func() {
		writeTestFile(t, modFile, benchmarkRunTestMod)
		writeTestFile(t, filepath.Join(dir, ".mod.cache.json"), "{}")
		removeTestDep(t, dir, "v1.1.0")
	})
}

// several mod commands back to back while a reload runs give one more reload after it, not one each
func TestWorkspaceAPI_ModCommandsDuringAReload(t *testing.T) {
	router, w, dir := newWatchedTestRouterWith(t, true)
	modCommand(t, router, w, dir, "install", func() {
		installTestDep(t, dir, "v1.0.0")
		writeTestFile(t, filepath.Join(dir, "mod.pp"), modWithRequire)
	})
	before := getLoadState(t, router)

	w.LoadLock()
	locked := true
	unlock := func() {
		if locked {
			locked = false
			w.LoadUnlock()
		}
	}
	t.Cleanup(unlock)

	installTestDep(t, dir, "v1.1.0")
	waitForState(t, router, "the reload did not start", func(s workspace.LoadState) bool { return s.Loading })
	// the reload is running (held up), and two more commands arrive
	installTestDep(t, dir, "v1.2.0")
	installTestDep(t, dir, "v1.3.0")
	unlock()

	waitForState(t, router, "the reloads did not end", func(s workspace.LoadState) bool {
		return !s.Loading && s.LoadCount >= before.LoadCount+2
	})
	time.Sleep(2500 * time.Millisecond)
	if after := getLoadState(t, router); after.Loading || after.LoadCount != before.LoadCount+2 {
		t.Fatalf("expected the reload that ran, and one more for the commands that came during it: %d before, got %+v", before.LoadCount, after)
	}
	state := getLoadState(t, router)
	lock, err := workspace.FileChangeTime(filepath.Join(dir, ".mod.cache.json"))
	if err != nil {
		t.Fatal(err)
	}
	if state.LastLoad.Error != "" || state.LastLoad.NewestFileTime == nil || state.LastLoad.NewestFileTime.Before(lock) {
		t.Fatalf("the last load has not read the last command: %+v (lock file changed %v)", state.LastLoad, lock)
	}
}
