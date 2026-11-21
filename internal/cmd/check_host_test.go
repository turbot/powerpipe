package cmd

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/v2/app_specific"
	"github.com/turbot/pipe-fittings/v2/constants"
	"github.com/turbot/powerpipe/internal/benchmarkrun"
	localconstants "github.com/turbot/powerpipe/internal/constants"
	"github.com/turbot/powerpipe/internal/resources"
)

// a run that another client deletes while it runs is answered with 404 from then on
func TestWaitForHostRun_Removed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"status":404,"title":"Not Found","detail":"benchmark run not found"}`))
	}))
	t.Cleanup(server.Close)
	client, err := benchmarkrun.NewClient(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	runURL, err := url.Parse(server.URL + "/api/v0/benchmark/b/run/1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := waitForHostRun(context.Background(), client, runURL); !errors.Is(err, errHostRunRemoved) {
		t.Fatalf("expected the run removed error, got %v", err)
	}
}

// a server that fails after it has started the snapshot response ends the connection without the final chunk, and
// the client must report that instead of keeping the part it received
func TestHostSnapshotAbortedMidStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"panels":{"a":`))
		w.(http.Flusher).Flush()
		panic(http.ErrAbortHandler)
	}))
	t.Cleanup(server.Close)
	client, err := benchmarkrun.NewClient(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	runURL, err := url.Parse(server.URL + "/api/v0/benchmark/b/run/1")
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if err := writeHostSnapshot(context.Background(), client, runURL, &out); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("expected an unexpected EOF writing the snapshot, got %v", err)
	}

	dest := filepath.Join(t.TempDir(), "run.pps")
	if err := (&hostSnapshotExporter{}).Export(context.Background(), hostSnapshot{client: client, runURL: runURL}, dest); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("expected an unexpected EOF exporting the snapshot, got %v", err)
	}
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		t.Fatalf("expected no partial snapshot file, stat gave %v", err)
	}
}

func TestValidateHostArgs(t *testing.T) {
	tests := map[string]struct {
		args    []string
		env     map[string]string
		wantErr string
	}{
		"snapshot export":             {args: []string{"b", "--output=none", "--export=out.pps"}},
		"snapshot export by format":   {args: []string{"b", "--output=none", "--export=pps"}},
		"snapshot output":             {args: []string{"b", "--output=snapshot"}},
		"short snapshot output":       {args: []string{"b", "--output=pps"}},
		"allowed flags":               {args: []string{"b", "--output=none", "--snapshot-title=t", "--progress=false", "--benchmark-timeout=60"}},
		"default text output":         {args: []string{"b"}, wantErr: "--output=snapshot or --output=none"},
		"csv output":                  {args: []string{"b", "--output=csv"}, wantErr: "--output=snapshot or --output=none"},
		"csv export":                  {args: []string{"b", "--output=none", "--export=out.csv"}, wantErr: "only snapshot (pps) export"},
		"two benchmarks":              {args: []string{"a", "b", "--output=none"}, wantErr: "one benchmark at a time"},
		"local mod":                   {args: []string{"b", "--output=none", "--mod-location=/tmp"}, wantErr: "--mod-location cannot be used"},
		"local database and filters":  {args: []string{"b", "--output=none", "--database=x", "--where=y"}, wantErr: "--database, --where cannot be used"},
		"upload":                      {args: []string{"b", "--output=none", "--share"}, wantErr: "--share cannot be used"},
		"variable in the environment": {args: []string{"b", "--output=none"}, env: map[string]string{"PP_VAR_region": "us-east-1", "PP_VAR_age": "30"}, wantErr: "PP_VAR_age, PP_VAR_region cannot be used with --host"},
		"variable and flag":           {args: []string{"b", "--output=none", "--var=a=b"}, env: map[string]string{"PP_VAR_region": "us-east-1"}, wantErr: "--var, PP_VAR_region cannot be used with --host"},
		"other environment settings":  {args: []string{"b", "--output=none"}, env: map[string]string{"POWERPIPE_DATABASE": "x", "POWERPIPE_MAX_PARALLEL": "2", "PP_VARIABLE": "x"}},
		"host without a scheme":       {args: []string{"b", "--output=none", "--host=localhost:9033"}, wantErr: "invalid host"},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			viper.Reset()
			t.Cleanup(viper.Reset)
			checkOutputMode = localconstants.CheckOutputModeText
			setEnvVarPrefix(t)
			for k, v := range test.env {
				t.Setenv(k, v)
			}

			cmd := checkCmd[*resources.Benchmark]()
			if err := cmd.ParseFlags(append([]string{"--host=http://localhost:9033"}, test.args...)); err != nil {
				t.Fatal(err)
			}
			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				if err := viper.BindPFlag(flag.Name, flag); err != nil {
					t.Fatal(err)
				}
			})

			_, _, err := validateHostArgs(cmd, cmd.Flags().Args())
			switch {
			case test.wantErr == "" && err != nil:
				t.Fatalf("unexpected error: %v", err)
			case test.wantErr != "" && (err == nil || !strings.Contains(err.Error(), test.wantErr)):
				t.Fatalf("expected an error containing %q, got %v", test.wantErr, err)
			}
		})
	}
}

// while the server reloads its workspace it answers a run request with 503 and Retry-After: the client waits as
// asked and tries again
func TestStartHostRun_RetriesWhileTheServerLoads(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if requests.Add(1) == 1 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"status":503,"title":"Service Unavailable","detail":"the workspace is loading a change to its files, try again when it has loaded"}`))
			return
		}
		w.Header().Set("Location", "/api/v0/benchmark/local.benchmark.b/run/1")
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"id":"1","benchmark":"local.benchmark.b","status":"running"}`))
	}))
	t.Cleanup(server.Close)
	client, err := benchmarkrun.NewClient(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	run, runURL, err := startHostRun(context.Background(), client, "b", 0)
	if err != nil {
		t.Fatal(err)
	}
	if run.ID != "1" || runURL.Path != "/api/v0/benchmark/local.benchmark.b/run/1" || requests.Load() != 2 {
		t.Fatalf("got run %+v at %v after %d requests", run, runURL, requests.Load())
	}
	if elapsed := time.Since(start); elapsed < time.Second {
		t.Fatalf("the client tried again after %s, before the Retry-After of 1s", elapsed)
	}
}

// a 503 with no Retry-After, such as from a proxy, is an error; a cancelled command stops waiting
func TestStartHostRun_Unavailable(t *testing.T) {
	unavailable := func(retryAfter string) *benchmarkrun.Client {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			if retryAfter != "" {
				w.Header().Set("Retry-After", retryAfter)
			}
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		t.Cleanup(server.Close)
		client, err := benchmarkrun.NewClient(server.URL)
		if err != nil {
			t.Fatal(err)
		}
		return client
	}

	var apiErr *benchmarkrun.APIError
	if _, _, err := startHostRun(context.Background(), unavailable(""), "b", 0); !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("a 503 with no Retry-After: got %v, want the 503 error", err)
	}

	client := unavailable("30")
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	start := time.Now()
	if _, _, err := startHostRun(ctx, client, "b", 0); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("a cancelled wait: got %v, want the context error", err)
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Fatalf("a cancelled wait took %s", elapsed)
	}
}

func TestHostIgnoredEnvVars(t *testing.T) {
	setEnvVarPrefix(t)
	got := hostIgnoredEnvVars([]string{"PP_VAR_b=1", "PATH=/bin", "PP_VAR_a=", "PP_VAR_=x", "PP_VAR_noequals", "pp_var_c=1", "=C:=C:\\"})
	if strings.Join(got, ",") != "PP_VAR_a,PP_VAR_b" {
		t.Fatalf("got %v", got)
	}
}

// a variable set in the environment stops a run on a server before the server is contacted, with the exit code
// --var gets
func TestRunBenchmarkOnHost_RefusesEnvVarBeforeContactingServer(t *testing.T) {
	for name, extra := range map[string][]string{"env var": nil, "var flag": {"--var=region=x"}} {
		t.Run(name, func(t *testing.T) {
			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				requests.Add(1)
				w.WriteHeader(http.StatusInternalServerError)
			}))
			defer server.Close()

			viper.Reset()
			t.Cleanup(viper.Reset)
			checkOutputMode = localconstants.CheckOutputModeText
			setEnvVarPrefix(t)
			if extra == nil {
				t.Setenv("PP_VAR_region", "us-east-1")
			}
			prevExit := exitCode
			t.Cleanup(func() { exitCode = prevExit })
			exitCode = 0

			cmd := checkCmd[*resources.Benchmark]()
			if err := cmd.ParseFlags(append([]string{"--host=" + server.URL, "--output=none", "b"}, extra...)); err != nil {
				t.Fatal(err)
			}
			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				if err := viper.BindPFlag(flag.Name, flag); err != nil {
					t.Fatal(err)
				}
			})

			runBenchmarkOnHost(context.Background(), cmd, cmd.Flags().Args())
			if exitCode != constants.ExitCodeInsufficientOrWrongInputs {
				t.Fatalf("exit code %d, want %d", exitCode, constants.ExitCodeInsufficientOrWrongInputs)
			}
			if n := requests.Load(); n != 0 {
				t.Fatalf("the server was contacted %d times", n)
			}
		})
	}
}

// setEnvVarPrefix sets the prefix main sets for variable values in the environment (cmdconfig.SetAppSpecificConstants)
func setEnvVarPrefix(t *testing.T) {
	t.Helper()
	prev := app_specific.EnvInputVarPrefix
	app_specific.EnvInputVarPrefix = "PP_VAR_"
	t.Cleanup(func() { app_specific.EnvInputVarPrefix = prev })
}
