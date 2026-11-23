package api

import (
	"context"
	"encoding/json"
	"errors"
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
	"github.com/turbot/powerpipe/internal/dashboardserver"
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
	runner := newTestBenchmarkRunner(t, benchmarkrun.MaxRunningRuns)
	gin.SetMode(gin.TestMode)
	binding.EnableDecoderDisallowUnknownFields = true
	router := gin.New()
	registerBenchmarkRun(router.Group("/api/:api_version"), runner)
	return router
}

// newTestBenchmarkRunner returns a runner for a test mod, with a SQLite database
func newTestBenchmarkRunner(t *testing.T, maxRunning int) *benchmarkrun.Runner {
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
	return benchmarkrun.NewRunner(w, client, benchmarkrun.IdleTimeout, maxRunning)
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

// failingWriter writes half of the first write to the connection, then fails, as a snapshot export that fails
// part way through does
type failingWriter struct {
	http.ResponseWriter
}

func (w failingWriter) Write(p []byte) (int, error) {
	n, _ := w.ResponseWriter.Write(p[:len(p)/2])
	w.ResponseWriter.(http.Flusher).Flush()
	return n, errors.New("export failed")
}

func (w failingWriter) Flush() { w.ResponseWriter.(http.Flusher).Flush() }

// a snapshot export that fails after the response has started must fail the transfer, not end it as complete
func TestRouter_SnapshotExportFailsMidStream(t *testing.T) {
	runner := newTestBenchmarkRunner(t, 1)
	router := newTestRouterFor(t, &APIService{HTTPListen: dashboardserver.ListenTypeLocal, benchmarkRunner: runner})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/snapshot") {
			w = failingWriter{w}
		}
		router.ServeHTTP(w, r)
	}))
	t.Cleanup(server.Close)
	// a client that does not ask for gzip, so that the failure is not hidden by compression
	client := &http.Client{Transport: &http.Transport{DisableCompression: true}}

	post, err := http.Post(server.URL+"/api/v0/benchmark/quick/run", "application/json", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}
	defer post.Body.Close()
	runURL := server.URL + post.Header.Get("Location")
	complete := false
	for deadline := time.Now().Add(60 * time.Second); !complete && time.Now().Before(deadline); time.Sleep(20 * time.Millisecond) {
		resp, err := client.Get(runURL)
		if err != nil {
			t.Fatal(err)
		}
		var run benchmarkrun.Run
		err = json.NewDecoder(resp.Body).Decode(&run)
		resp.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		complete = run.Status != benchmarkrun.StatusRunning
	}
	if !complete {
		t.Fatal("the run did not finish")
	}

	// the headers are flushed before the export fails, so the request succeeds and the body must fail
	resp, err := client.Get(runURL + "/snapshot")
	if err != nil {
		t.Fatalf("the response did not start: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("snapshot: got status %d, want 200", resp.StatusCode)
	}
	if body, err := io.ReadAll(resp.Body); err == nil {
		t.Fatalf("the transfer of a failed snapshot export ended cleanly after %d bytes, status %d", len(body), resp.StatusCode)
	}
}
