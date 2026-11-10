package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/viper"
	"github.com/turbot/powerpipe/internal/benchmarkrun"
	"github.com/turbot/powerpipe/internal/cmdconfig"
	"github.com/turbot/powerpipe/internal/db_client"
	"github.com/turbot/powerpipe/internal/workspace"
)

const benchmarkRunTestMod = `
mod "local" {
  title = "benchmark run API test"
}

benchmark "quick" {
  children = [control.ok_row]
}

benchmark "slow" {
  children = [control.slow]
}

control "ok_row" {
  sql = "select 'r' as resource, 'ok' as status, 'fine' as reason"
}

control "slow" {
  sql = "with recursive c(x) as (select 1 union all select x + 1 from c where x < 100000000) select 'r' as resource, 'ok' as status, count(*) as reason from c"
}
`

func newBenchmarkRunRouter(t *testing.T) *gin.Engine {
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
	client, err := db_client.NewDbClient(context.Background(), "sqlite://"+filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = client.Close(context.Background()) })

	gin.SetMode(gin.TestMode)
	binding.EnableDecoderDisallowUnknownFields = true
	router := gin.New()
	registerBenchmarkRun(router.Group("/api/:api_version"), benchmarkrun.NewRunner(w, client, benchmarkrun.IdleTimeout))
	return router
}

// request sends a request to the router. path is the full path, e.g. /api/v0/benchmark/x/run
func request(t *testing.T, router *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	var reader io.Reader = http.NoBody
	if body != "" {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func decodeRun(t *testing.T, rec *httptest.ResponseRecorder) benchmarkrun.Run {
	t.Helper()
	var run benchmarkrun.Run
	if err := json.Unmarshal(rec.Body.Bytes(), &run); err != nil {
		t.Fatalf("invalid run JSON %q: %v", rec.Body.String(), err)
	}
	return run
}

// start starts a run and returns the URL of the run, from the Location header
func start(t *testing.T, router *gin.Engine, benchmark, body string) string {
	t.Helper()
	rec := request(t, router, http.MethodPost, "/api/v0/benchmark/"+benchmark+"/run", body)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("start: got %d %s", rec.Code, rec.Body)
	}
	run := decodeRun(t, rec)
	want := "/api/v0/benchmark/local.benchmark." + strings.TrimPrefix(benchmark, "benchmark.") + "/run/" + run.ID
	if location := rec.Header().Get("Location"); location != want {
		t.Fatalf("expected Location %q, got %q", want, location)
	}
	return want
}

func TestBenchmarkRunAPI(t *testing.T) {
	router := newBenchmarkRunRouter(t)
	runURL := start(t, router, "benchmark.quick", `{"timeout_seconds": 60}`)

	var rec *httptest.ResponseRecorder
	var run benchmarkrun.Run
	for deadline := time.Now().Add(60 * time.Second); time.Now().Before(deadline); time.Sleep(20 * time.Millisecond) {
		rec = request(t, router, http.MethodGet, runURL, "")
		if rec.Code != http.StatusOK {
			t.Fatalf("get: got %d %s", rec.Code, rec.Body)
		}
		if run = decodeRun(t, rec); run.Status != benchmarkrun.StatusRunning {
			break
		}
	}
	if run.Status != benchmarkrun.StatusComplete || run.Summary.Ok != 1 || run.Benchmark != "local.benchmark.quick" {
		t.Fatalf("unexpected run %s", rec.Body)
	}

	rec = request(t, router, http.MethodGet, runURL+"/snapshot", "")
	var snapshot struct {
		Panels map[string]any `json:"panels"`
	}
	if rec.Code != http.StatusOK || json.Unmarshal(rec.Body.Bytes(), &snapshot) != nil || snapshot.Panels["local.control.ok_row"] == nil {
		t.Fatalf("snapshot: got %d %.300s", rec.Code, rec.Body)
	}

	if rec = request(t, router, http.MethodDelete, runURL, ""); rec.Code != http.StatusNoContent {
		t.Fatalf("delete: got %d %s", rec.Code, rec.Body)
	}
	for _, path := range []string{runURL, runURL + "/snapshot"} {
		if rec = request(t, router, http.MethodGet, path, ""); rec.Code != http.StatusNotFound {
			t.Fatalf("%s after delete: got %d %s", path, rec.Code, rec.Body)
		}
	}
	if rec = request(t, router, http.MethodDelete, runURL, ""); rec.Code != http.StatusNotFound {
		t.Fatalf("second delete: got %d %s", rec.Code, rec.Body)
	}
}

// the body is optional, the benchmark can be named in short form, and the run is found by any of its names
func TestBenchmarkRunAPI_Names(t *testing.T) {
	router := newBenchmarkRunRouter(t)
	runURL := start(t, router, "quick", "")
	id := runURL[strings.LastIndex(runURL, "/")+1:]
	for _, name := range []string{"local.benchmark.quick", "benchmark.quick", "quick"} {
		if rec := request(t, router, http.MethodGet, "/api/v0/benchmark/"+name+"/run/"+id, ""); rec.Code != http.StatusOK {
			t.Errorf("%s: got %d %s", name, rec.Code, rec.Body)
		}
	}
	// a run is not found under another benchmark
	for _, path := range []string{"/benchmark/benchmark.slow/run/" + id, "/benchmark/benchmark.missing/run/" + id, "/benchmark/benchmark.slow/run/" + id + "/snapshot"} {
		if rec := request(t, router, http.MethodGet, "/api/v0"+path, ""); rec.Code != http.StatusNotFound {
			t.Errorf("%s: got %d %s", path, rec.Code, rec.Body)
		}
	}
	if rec := request(t, router, http.MethodDelete, "/api/v0/benchmark/benchmark.slow/run/"+id, ""); rec.Code != http.StatusNotFound {
		t.Errorf("delete under another benchmark: got %d %s", rec.Code, rec.Body)
	}
	if rec := request(t, router, http.MethodDelete, runURL, ""); rec.Code != http.StatusNoContent {
		t.Fatalf("delete: got %d %s", rec.Code, rec.Body)
	}
}

func TestBenchmarkRunAPI_SnapshotOfRunningRun(t *testing.T) {
	router := newBenchmarkRunRouter(t)
	runURL := start(t, router, "benchmark.slow", "")
	if rec := request(t, router, http.MethodGet, runURL+"/snapshot", ""); rec.Code != http.StatusConflict {
		t.Fatalf("snapshot of a running run: got %d %s", rec.Code, rec.Body)
	}
	if rec := request(t, router, http.MethodDelete, runURL, ""); rec.Code != http.StatusNoContent {
		t.Fatalf("delete of a running run: got %d %s", rec.Code, rec.Body)
	}
}

func TestBenchmarkRunAPI_BadRequests(t *testing.T) {
	router := newBenchmarkRunRouter(t)
	for name, body := range map[string]string{
		"not json":         `{`,
		"unknown field":    `{"benchmark": "benchmark.quick"}`,
		"negative timeout": `{"timeout_seconds": -1}`,
	} {
		if rec := request(t, router, http.MethodPost, "/api/v0/benchmark/benchmark.quick/run", body); rec.Code != http.StatusBadRequest {
			t.Errorf("%s: got %d %s", name, rec.Code, rec.Body)
		}
	}
	// the content types that a web page can send cross-site without a CORS preflight
	for _, contentType := range []string{"", "text/plain", "application/x-www-form-urlencoded", "multipart/form-data; boundary=x"} {
		req := httptest.NewRequest(http.MethodPost, "/api/v0/benchmark/benchmark.quick/run", http.NoBody)
		if contentType != "" {
			req.Header.Set("Content-Type", contentType)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "application/json") {
			t.Errorf("content type %q: got %d %s", contentType, rec.Code, rec.Body)
		}
	}
	// a JSON content type with parameters is fine
	req := httptest.NewRequest(http.MethodPost, "/api/v0/benchmark/benchmark.quick/run", http.NoBody)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusAccepted {
		t.Errorf("json with charset: got %d %s", rec.Code, rec.Body)
	} else if rec := request(t, router, http.MethodDelete, rec.Header().Get("Location"), ""); rec.Code != http.StatusNoContent {
		// the delete waits for the run to stop, so that it does not write to the test's directory after the test
		t.Errorf("delete the run: got %d %s", rec.Code, rec.Body)
	}
	for _, name := range []string{"benchmark.missing", "control.ok_row"} {
		if rec := request(t, router, http.MethodPost, "/api/v0/benchmark/"+name+"/run", ""); rec.Code != http.StatusNotFound {
			t.Errorf("start %s: got %d %s", name, rec.Code, rec.Body)
		}
	}
	if rec := request(t, router, http.MethodGet, "/api/v0/benchmark/benchmark.quick/run/nope", ""); rec.Code != http.StatusNotFound {
		t.Errorf("unknown run: got %d %s", rec.Code, rec.Body)
	}
}
