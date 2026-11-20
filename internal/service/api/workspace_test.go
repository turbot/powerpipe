package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/turbot/powerpipe/internal/benchmarkrun"
	"github.com/turbot/powerpipe/internal/cmdconfig"
	"github.com/turbot/powerpipe/internal/dashboardserver"
	"github.com/turbot/powerpipe/internal/db_client"
	"github.com/turbot/powerpipe/internal/workspace"
)

// newWatchedTestRouter returns the router of `powerpipe server` for the benchmark run test mod, with the workspace
// watched for changes as the server watches it
func newWatchedTestRouter(t *testing.T) (*gin.Engine, *workspace.PowerpipeWorkspace, string) {
	t.Helper()
	viper.Set("main.version", "0.0.0")
	cmdconfig.SetAppSpecificConstants()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "mod.pp"), []byte(benchmarkRunTestMod), 0o600); err != nil {
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
	runner := benchmarkrun.NewRunner(w, client, benchmarkrun.IdleTimeout, benchmarkrun.MaxRunningRuns)
	router := newTestRouterFor(t, &APIService{HTTPListen: dashboardserver.ListenTypeLocal, workspace: w, benchmarkRunner: runner})
	return router, w, dir
}

// serve sends the request to the router, and fails the test if the router takes more than a second to answer
func serve(t *testing.T, router *gin.Engine, method, path string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader("{}"))
	req.Host = "localhost:9033"
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	done := make(chan struct{})
	go func() {
		router.ServeHTTP(rec, req)
		close(done)
	}()
	select {
	case <-done:
		return rec
	case <-time.After(time.Second):
		t.Fatalf("%s %s did not answer within a second", method, path)
		return nil
	}
}

func getLoadState(t *testing.T, router *gin.Engine) workspace.LoadState {
	t.Helper()
	rec := serve(t, router, http.MethodGet, "/api/v0/workspace")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/v0/workspace: got %d %s", rec.Code, rec.Body)
	}
	var state workspace.LoadState
	if err := json.Unmarshal(rec.Body.Bytes(), &state); err != nil {
		t.Fatalf("invalid load state %q: %v", rec.Body, err)
	}
	return state
}

// a client asks the server for its load state, and while the workspace reloads, a run request is answered at once
// with 503 and Retry-After
func TestWorkspaceAPI_Reload(t *testing.T) {
	router, w, dir := newWatchedTestRouter(t)
	modFile := filepath.Join(dir, "mod.pp")

	rec := serve(t, router, http.MethodGet, "/api/v0/workspace")
	for _, field := range []string{`"loading":false`, `"load_count":1`, `"last_load":{"start_time":`, `"end_time":`, `"newest_file_time":`} {
		if !strings.Contains(rec.Body.String(), field) {
			t.Fatalf("the load state %s has no %s", rec.Body, field)
		}
	}
	if rec.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("the load state can be cached: Cache-Control %q", rec.Header().Get("Cache-Control"))
	}

	// a load that takes as long as the test wants, as a large workspace does
	w.LoadLock()
	locked := true
	t.Cleanup(func() {
		if locked {
			w.LoadUnlock()
		}
	})
	if err := os.WriteFile(modFile, []byte(benchmarkRunTestMod+"\nbenchmark \"added\" {\n  children = [control.ok_row]\n}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(30 * time.Second)
	for {
		state := getLoadState(t, router)
		if state.Loading {
			if state.CurrentLoad == nil || state.LoadCount != 1 {
				t.Fatalf("the load state while loading: %+v", state)
			}
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("the reload did not start")
		}
		time.Sleep(10 * time.Millisecond)
	}

	// once the reload holds the workspace, as it does for the whole of a load, a run request is answered at once
	for {
		release, ok := w.HoldLoaded()
		if !ok {
			break
		}
		release()
		if time.Now().After(deadline) {
			t.Fatal("the reload did not take the workspace")
		}
		time.Sleep(10 * time.Millisecond)
	}
	rec = serve(t, router, http.MethodPost, "/api/v0/benchmark/quick/run")
	if rec.Code != http.StatusServiceUnavailable || !strings.Contains(rec.Body.String(), "the workspace is loading") {
		t.Fatalf("a run request during a reload: got %d %s, want 503", rec.Code, rec.Body)
	}
	if retry, err := strconv.Atoi(rec.Header().Get("Retry-After")); err != nil || retry < 1 || retry > 30 {
		t.Fatalf("a run request during a reload: Retry-After %q, want 1-30 seconds", rec.Header().Get("Retry-After"))
	}
	if state := getLoadState(t, router); !state.Loading {
		t.Fatalf("the load state during the reload: %+v", state)
	}

	locked = false
	w.LoadUnlock()
	changed, err := workspace.FileChangeTime(modFile)
	if err != nil {
		t.Fatal(err)
	}
	for {
		state := getLoadState(t, router)
		if !state.Loading && state.LastLoad.NewestFileTime != nil && !state.LastLoad.NewestFileTime.Before(changed) {
			if state.LoadCount < 2 || state.LastLoad.Error != "" {
				t.Fatalf("the load state after the reload: %+v", state)
			}
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("the change did not load: %+v", state)
		}
		time.Sleep(10 * time.Millisecond)
	}

	// the benchmark the change added runs
	rec = serve(t, router, http.MethodPost, "/api/v0/benchmark/added/run")
	if rec.Code != http.StatusAccepted {
		t.Fatalf("a run of the added benchmark: got %d %s, want 202", rec.Code, rec.Body)
	}
	serve(t, router, http.MethodDelete, rec.Header().Get("Location"))
}

// a server with no workspace has no workspace API, as before
func TestWorkspaceAPI_NoWorkspace(t *testing.T) {
	router := newTestRouter(t, dashboardserver.ListenTypeLocal)
	rec := serve(t, router, http.MethodGet, "/api/v0/workspace")
	if strings.Contains(rec.Body.String(), "load_count") {
		t.Fatalf("a server with no workspace answered with a load state: %s", rec.Body)
	}
}
