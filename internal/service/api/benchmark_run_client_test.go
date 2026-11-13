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

	"github.com/turbot/pipe-fittings/v2/filepaths"
	"github.com/turbot/powerpipe/internal/benchmarkrun"
	"github.com/turbot/powerpipe/internal/dashboardserver"
)

// a server before v1.6.0 has no run API: it answers the run request with the dashboard page and 200, or with a
// plain text 404 when the dashboard is not installed
func TestBenchmarkRunClient_OldServer(t *testing.T) {
	for name, installDashboard := range map[string]bool{"dashboard page": true, "no dashboard": false} {
		t.Run(name, func(t *testing.T) {
			router := newTestRouterFor(t, &APIService{HTTPListen: dashboardserver.ListenTypeLocal})
			if installDashboard {
				if err := os.WriteFile(filepath.Join(filepaths.EnsureDashboardAssetsDir(), "index.html"), []byte("<html></html>"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			server := httptest.NewServer(router)
			t.Cleanup(server.Close)
			client, err := benchmarkrun.NewClient(server.URL)
			if err != nil {
				t.Fatal(err)
			}
			_, _, err = client.Start(context.Background(), "quick", time.Minute)
			if err == nil || err.Error() != "the powerpipe server at "+server.URL+" cannot run benchmarks: it needs powerpipe v1.6.0 or later" {
				t.Fatalf("expected an error that the server cannot run benchmarks, got %v", err)
			}
		})
	}
}

// the client that `powerpipe benchmark run --host` uses, against the API over HTTP
func TestBenchmarkRunClient(t *testing.T) {
	server := httptest.NewServer(newBenchmarkRunRouter(t))
	t.Cleanup(server.Close)
	client, err := benchmarkrun.NewClient(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	run, runURL, err := client.Start(ctx, "quick", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if run.Benchmark != "local.benchmark.quick" || runURL.String() != server.URL+"/api/v0/benchmark/local.benchmark.quick/run/"+run.ID {
		t.Fatalf("unexpected run %+v at %s", run, runURL)
	}
	for deadline := time.Now().Add(60 * time.Second); run.Status == benchmarkrun.StatusRunning && time.Now().Before(deadline); time.Sleep(20 * time.Millisecond) {
		if run, err = client.Get(ctx, runURL); err != nil {
			t.Fatal(err)
		}
	}
	if run.Status != benchmarkrun.StatusComplete || run.Summary.Ok != 1 {
		t.Fatalf("unexpected run %+v", run)
	}

	body, err := client.Snapshot(ctx, runURL)
	if err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(body)
	_ = body.Close()
	var snapshot struct {
		Panels map[string]any `json:"panels"`
	}
	if err != nil || json.Unmarshal(data, &snapshot) != nil || snapshot.Panels["local.control.ok_row"] == nil {
		t.Fatalf("unexpected snapshot (%v) %.300s", err, data)
	}

	if err := client.Delete(ctx, runURL); err != nil {
		t.Fatal(err)
	}
	var apiErr *benchmarkrun.APIError
	if _, err := client.Get(ctx, runURL); !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusNotFound {
		t.Fatalf("expected not found after delete, got %v", err)
	}
}

func TestBenchmarkRunClient_Errors(t *testing.T) {
	server := httptest.NewServer(newBenchmarkRunRouter(t))
	t.Cleanup(server.Close)
	client, err := benchmarkrun.NewClient(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	// the error message from the server reaches the caller
	var apiErr *benchmarkrun.APIError
	_, _, err = client.Start(ctx, "benchmark.missing", 0)
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusNotFound || !strings.Contains(err.Error(), "benchmark.missing") {
		t.Fatalf("expected a not found error naming the benchmark, got %v", err)
	}

	// a running run has no snapshot yet
	_, runURL, err := client.Start(ctx, "benchmark.slow", 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Snapshot(ctx, runURL); !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusConflict {
		t.Fatalf("expected a conflict for the snapshot of a running run, got %v", err)
	}
	if err := client.Delete(ctx, runURL); err != nil {
		t.Fatal(err)
	}

	// a full runner answers 429
	full := httptest.NewServer(newTestRouterFor(t, &APIService{HTTPListen: dashboardserver.ListenTypeLocal, benchmarkRunner: newTestBenchmarkRunner(t, 1)}))
	t.Cleanup(full.Close)
	fullClient, err := benchmarkrun.NewClient(full.URL)
	if err != nil {
		t.Fatal(err)
	}
	_, slowURL, err := fullClient.Start(ctx, "benchmark.slow", 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = fullClient.Delete(context.Background(), slowURL) })
	_, _, err = fullClient.Start(ctx, "benchmark.quick", 0)
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusTooManyRequests || err.Error() != "the powerpipe server at "+full.URL+" is busy: "+benchmarkrun.ErrTooManyRuns.Error() {
		t.Fatalf("expected a busy error, got %v", err)
	}

	// the server sends 413 for a body over its limit as plain text, before the run API sees the request
	tooLarge := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		_, _ = w.Write([]byte("request too large"))
	}))
	t.Cleanup(tooLarge.Close)
	tooLargeClient, err := benchmarkrun.NewClient(tooLarge.URL)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = tooLargeClient.Start(ctx, "benchmark.quick", 0)
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusRequestEntityTooLarge || err.Error() != "the powerpipe server at "+tooLarge.URL+" refused the request as too large" {
		t.Fatalf("expected a too large error, got %v", err)
	}

	for _, host := range []string{"localhost:9033", "ftp://localhost", "http://", ""} {
		if _, err := benchmarkrun.NewClient(host); err == nil {
			t.Errorf("%q: expected an invalid host error", host)
		}
	}
}
