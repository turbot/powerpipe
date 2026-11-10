package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/turbot/powerpipe/internal/benchmarkrun"
)

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

	for _, host := range []string{"localhost:9033", "ftp://localhost", "http://", ""} {
		if _, err := benchmarkrun.NewClient(host); err == nil {
			t.Errorf("%q: expected an invalid host error", host)
		}
	}
}
