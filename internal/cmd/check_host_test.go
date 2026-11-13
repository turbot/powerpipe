package cmd

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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

func TestValidateHostArgs(t *testing.T) {
	tests := map[string]struct {
		args    []string
		wantErr string
	}{
		"snapshot export":            {args: []string{"b", "--output=none", "--export=out.pps"}},
		"snapshot export by format":  {args: []string{"b", "--output=none", "--export=pps"}},
		"snapshot output":            {args: []string{"b", "--output=snapshot"}},
		"short snapshot output":      {args: []string{"b", "--output=pps"}},
		"allowed flags":              {args: []string{"b", "--output=none", "--snapshot-title=t", "--progress=false", "--benchmark-timeout=60"}},
		"default text output":        {args: []string{"b"}, wantErr: "--output=snapshot or --output=none"},
		"csv output":                 {args: []string{"b", "--output=csv"}, wantErr: "--output=snapshot or --output=none"},
		"csv export":                 {args: []string{"b", "--output=none", "--export=out.csv"}, wantErr: "only snapshot (pps) export"},
		"two benchmarks":             {args: []string{"a", "b", "--output=none"}, wantErr: "one benchmark at a time"},
		"local mod":                  {args: []string{"b", "--output=none", "--mod-location=/tmp"}, wantErr: "--mod-location cannot be used"},
		"local database and filters": {args: []string{"b", "--output=none", "--database=x", "--where=y"}, wantErr: "--database, --where cannot be used"},
		"upload":                     {args: []string{"b", "--output=none", "--share"}, wantErr: "--share cannot be used"},
		"host without a scheme":      {args: []string{"b", "--output=none", "--host=localhost:9033"}, wantErr: "invalid host"},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			viper.Reset()
			t.Cleanup(viper.Reset)
			checkOutputMode = localconstants.CheckOutputModeText

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
