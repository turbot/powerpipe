package benchmarkrun_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/turbot/powerpipe/internal/benchmarkrun"
	"github.com/turbot/powerpipe/internal/cmdconfig"
	"github.com/turbot/powerpipe/internal/db_client"
	"github.com/turbot/powerpipe/internal/workspace"
)

var setAppConstants sync.Once

const runnerTestMod = `
mod "local" {
  title = "benchmark run test"
}

benchmark "quick" {
  children = [control.ok_rows, control.alarm_row, control.with_args]
}

benchmark "slow" {
  children = [control.slow]
}

control "ok_rows" {
  sql = "select 'r' || value as resource, 'ok' as status, 'reason ' || value as reason from (select 1 as value union all select 2)"
}

control "alarm_row" {
  sql = "select 'r' as resource, 'alarm' as status, 'bad <b>' as reason, 'us-east-1' as region"
}

control "with_args" {
  query = query.status_for
  args  = { status = "ok" }
}

query "status_for" {
  sql = "select 'r' as resource, $1 as status, 'from a query' as reason"
  param "status" {
    default = "error"
  }
}

control "slow" {
  sql = "with recursive c(x) as (select 1 union all select x + 1 from c where x < 100000000) select 'r' as resource, 'ok' as status, count(*) as reason from c"
}
`

// NewTestRunner returns a runner for a test mod, with a SQLite database
func NewTestRunner(t *testing.T) *benchmarkrun.Runner {
	return newTestRunner(t, benchmarkrun.IdleTimeout, benchmarkrun.MaxRunningRuns)
}

func newTestRunner(t *testing.T, idleTimeout time.Duration, maxRunning int) *benchmarkrun.Runner {
	t.Helper()
	setAppConstants.Do(func() {
		viper.Set("main.version", "0.0.0")
		cmdconfig.SetAppSpecificConstants()
	})
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "mod.pp"), []byte(runnerTestMod), 0o600); err != nil {
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
	return benchmarkrun.NewRunner(w, client, idleTimeout, maxRunning)
}

func waitForRun(t *testing.T, runner *benchmarkrun.Runner, started benchmarkrun.Run) benchmarkrun.Run {
	t.Helper()
	deadline := time.Now().Add(60 * time.Second)
	for time.Now().Before(deadline) {
		run, err := runner.Get(started.Benchmark, started.ID)
		if err != nil {
			t.Fatal(err)
		}
		if run.Status != benchmarkrun.StatusRunning {
			return run
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("run %s did not finish", started.ID)
	return benchmarkrun.Run{}
}

type testSnapshot struct {
	Panels map[string]struct {
		Status string `json:"status"`
		Data   *struct {
			Rows []map[string]any `json:"rows"`
		} `json:"data"`
	} `json:"panels"`
}

func readSnapshot(t *testing.T, runner *benchmarkrun.Runner, run benchmarkrun.Run) testSnapshot {
	t.Helper()
	var buf bytes.Buffer
	if err := runner.WriteSnapshot(context.Background(), run.Benchmark, run.ID, &buf); err != nil {
		t.Fatal(err)
	}
	var snapshot testSnapshot
	if err := json.Unmarshal(buf.Bytes(), &snapshot); err != nil {
		t.Fatalf("invalid snapshot: %v", err)
	}
	if !bytes.HasSuffix(buf.Bytes(), []byte("}\n")) || strings.Contains(buf.String(), "\n  ") {
		t.Fatalf("expected an unindented snapshot ending with a new line, as an export is")
	}
	return snapshot
}

func TestRunner_Complete(t *testing.T) {
	runner := NewTestRunner(t)
	started, err := runner.Start("benchmark.quick", 0)
	if err != nil {
		t.Fatal(err)
	}
	if started.Status != benchmarkrun.StatusRunning || started.ID == "" || started.Benchmark != "local.benchmark.quick" {
		t.Fatalf("unexpected started run %+v", started)
	}

	run := waitForRun(t, runner, started)
	if run.Status != benchmarkrun.StatusComplete || run.Summary == nil || run.Summary.Ok != 3 || run.Summary.Alarm != 1 || run.EndTime == nil {
		t.Fatalf("unexpected completed run %+v (summary %+v)", run, run.Summary)
	}

	snapshot := readSnapshot(t, runner, run)
	okRows := snapshot.Panels["local.control.ok_rows"].Data.Rows
	if len(okRows) != 2 || okRows[0]["reason"] != "reason 1" {
		t.Fatalf("unexpected ok_rows data %v", okRows)
	}
	alarmRows := snapshot.Panels["local.control.alarm_row"].Data.Rows
	if len(alarmRows) != 1 || alarmRows[0]["region"] != "us-east-1" || alarmRows[0]["reason"] != "bad <b>" {
		t.Fatalf("unexpected alarm_row data %v", alarmRows)
	}

	// the snapshot can be read again until the run is deleted
	readSnapshot(t, runner, run)
	if err := runner.Delete(run.Benchmark, run.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := runner.Get(run.Benchmark, run.ID); !errors.Is(err, benchmarkrun.ErrNotFound) {
		t.Fatalf("expected the run to be removed, got %v", err)
	}
}

// runs share the workspace resources: run with -race, this checks that a run does not change them
func TestRunner_ConcurrentRuns(t *testing.T) {
	runner := NewTestRunner(t)
	var started []benchmarkrun.Run
	for range 8 {
		run, err := runner.Start("benchmark.quick", 0)
		if err != nil {
			t.Fatal(err)
		}
		started = append(started, run)
	}
	for _, s := range started {
		run := waitForRun(t, runner, s)
		if run.Status != benchmarkrun.StatusComplete || run.Summary.Ok != 3 {
			t.Fatalf("unexpected run %+v", run)
		}
		if rows := readSnapshot(t, runner, run).Panels["local.control.ok_rows"].Data.Rows; len(rows) != 2 {
			t.Fatalf("run %s: expected 2 rows, got %d", run.ID, len(rows))
		}
	}
}

// a run is refused while the runner runs its maximum, and starts again when a run finishes
func TestRunner_MaxRunning(t *testing.T) {
	runner := newTestRunner(t, benchmarkrun.IdleTimeout, 2)
	var running []benchmarkrun.Run
	for range 2 {
		run, err := runner.Start("benchmark.slow", 0)
		if err != nil {
			t.Fatal(err)
		}
		running = append(running, run)
	}
	if _, err := runner.Start("benchmark.quick", 0); !errors.Is(err, benchmarkrun.ErrTooManyRuns) {
		t.Fatalf("expected ErrTooManyRuns, got %v", err)
	}
	if err := runner.Delete(running[0].Benchmark, running[0].ID); err != nil {
		t.Fatal(err)
	}
	run, err := runner.Start("benchmark.quick", 0)
	if err != nil {
		t.Fatalf("expected a run to start once a run is deleted, got %v", err)
	}
	if run = waitForRun(t, runner, run); run.Status != benchmarkrun.StatusComplete {
		t.Fatalf("unexpected run %+v", run)
	}
	// a complete run does not count
	if run, err = runner.Start("benchmark.quick", 0); err != nil {
		t.Fatalf("expected a run to start once a run is complete, got %v", err)
	}
	waitForRun(t, runner, run)
	if err := runner.Delete(running[1].Benchmark, running[1].ID); err != nil {
		t.Fatal(err)
	}
}

func TestRunner_DeleteCancelsRun(t *testing.T) {
	runner := NewTestRunner(t)
	run, err := runner.Start("benchmark.slow", 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := runner.WriteSnapshot(context.Background(), run.Benchmark, run.ID, &bytes.Buffer{}); !errors.Is(err, benchmarkrun.ErrNotComplete) {
		t.Fatalf("expected the snapshot of a running run to be unavailable, got %v", err)
	}
	start := time.Now()
	if err := runner.Delete(run.Benchmark, run.ID); err != nil {
		t.Fatal(err)
	}
	if elapsed := time.Since(start); elapsed > 30*time.Second {
		t.Fatalf("delete took %s: the run was not cancelled", elapsed)
	}
	if _, err := runner.Get(run.Benchmark, run.ID); !errors.Is(err, benchmarkrun.ErrNotFound) {
		t.Fatalf("expected the run to be removed, got %v", err)
	}
}

// a run that times out reports an error, and its snapshot can still be read, as the command still exports it
func TestRunner_Timeout(t *testing.T) {
	runner := NewTestRunner(t)
	started, err := runner.Start("benchmark.slow", 200*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	run := waitForRun(t, runner, started)
	if run.Status != benchmarkrun.StatusError || !strings.Contains(run.Error, "deadline exceeded") {
		t.Fatalf("expected a timeout error, got %+v", run)
	}
	if status := readSnapshot(t, runner, run).Panels["local.control.slow"].Status; status != "error" {
		t.Fatalf("expected the control to be in error, got %q", status)
	}
}

// blockingWriter blocks the first write until release is closed
type blockingWriter struct {
	started chan struct{}
	release chan struct{}
	once    sync.Once
	buf     bytes.Buffer
}

func (w *blockingWriter) Write(p []byte) (int, error) {
	w.once.Do(func() {
		close(w.started)
		<-w.release
	})
	return w.buf.Write(p)
}

// a delete during a snapshot write waits for the write, so that the write can read every row
func TestRunner_DeleteWaitsForSnapshotWrite(t *testing.T) {
	runner := NewTestRunner(t)
	started, err := runner.Start("benchmark.quick", 0)
	if err != nil {
		t.Fatal(err)
	}
	waitForRun(t, runner, started)

	w := &blockingWriter{started: make(chan struct{}), release: make(chan struct{})}
	writeErr := make(chan error, 1)
	go func() { writeErr <- runner.WriteSnapshot(context.Background(), started.Benchmark, started.ID, w) }()
	<-w.started

	deleted := make(chan error, 1)
	go func() { deleted <- runner.Delete(started.Benchmark, started.ID) }()
	select {
	case err := <-deleted:
		t.Fatalf("delete returned (%v) before the snapshot write completed", err)
	case <-time.After(200 * time.Millisecond):
	}
	// the run is already removed for new requests
	if _, err := runner.Get(started.Benchmark, started.ID); !errors.Is(err, benchmarkrun.ErrNotFound) {
		t.Fatalf("expected the run to be removed, got %v", err)
	}

	close(w.release)
	if err := <-writeErr; err != nil {
		t.Fatalf("snapshot write failed: %v", err)
	}
	if err := <-deleted; err != nil {
		t.Fatal(err)
	}
	var snapshot testSnapshot
	if err := json.Unmarshal(w.buf.Bytes(), &snapshot); err != nil {
		t.Fatalf("invalid snapshot: %v", err)
	}
	if rows := snapshot.Panels["local.control.ok_rows"].Data.Rows; len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
}

func TestRunner_Errors(t *testing.T) {
	runner := NewTestRunner(t)
	var notFound *benchmarkrun.BenchmarkNotFoundError
	for _, name := range []string{"benchmark.missing", "control.ok_rows"} {
		if _, err := runner.Start(name, 0); !errors.As(err, &notFound) {
			t.Fatalf("%s: expected a benchmark not found error, got %v", name, err)
		}
	}
	if _, err := runner.Get("benchmark.quick", "nope"); !errors.Is(err, benchmarkrun.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
	if err := runner.Delete("benchmark.quick", "nope"); !errors.Is(err, benchmarkrun.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
	if err := runner.WriteSnapshot(context.Background(), "benchmark.quick", "nope", &bytes.Buffer{}); !errors.Is(err, benchmarkrun.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

// a run is found by any name of its benchmark, and only by its benchmark
func TestRunner_RunBelongsToItsBenchmark(t *testing.T) {
	runner := NewTestRunner(t)
	started, err := runner.Start("quick", 0)
	if err != nil {
		t.Fatal(err)
	}
	if started.Benchmark != "local.benchmark.quick" {
		t.Fatalf("expected the full name, got %q", started.Benchmark)
	}
	waitForRun(t, runner, started)
	for _, name := range []string{"local.benchmark.quick", "benchmark.quick", "quick"} {
		if _, err := runner.Get(name, started.ID); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
	for _, name := range []string{"benchmark.slow", "benchmark.missing", "control.ok_rows"} {
		if _, err := runner.Get(name, started.ID); !errors.Is(err, benchmarkrun.ErrNotFound) {
			t.Errorf("%s: expected not found, got %v", name, err)
		}
		if err := runner.Delete(name, started.ID); !errors.Is(err, benchmarkrun.ErrNotFound) {
			t.Errorf("%s: expected delete to find nothing, got %v", name, err)
		}
	}
	if err := runner.Delete("quick", started.ID); err != nil {
		t.Fatal(err)
	}
}

const testIdleTimeout = 300 * time.Millisecond

// waitIdle waits long enough for a run that no client reads to be removed
func waitIdle() { time.Sleep(5 * testIdleTimeout) }

// a running run that no client reads is cancelled and removed
func TestRunner_IdleRunningRunIsRemoved(t *testing.T) {
	runner := newTestRunner(t, testIdleTimeout, benchmarkrun.MaxRunningRuns)
	run, err := runner.Start("benchmark.slow", 0)
	if err != nil {
		t.Fatal(err)
	}
	waitIdle()
	if _, err := runner.Get(run.Benchmark, run.ID); !errors.Is(err, benchmarkrun.ErrNotFound) {
		t.Fatalf("expected the idle run to be removed, got %v", err)
	}
}

// a completed run that no client downloads is removed
func TestRunner_IdleCompleteRunIsRemoved(t *testing.T) {
	runner := newTestRunner(t, testIdleTimeout, benchmarkrun.MaxRunningRuns)
	started, err := runner.Start("benchmark.quick", 0)
	if err != nil {
		t.Fatal(err)
	}
	waitForRun(t, runner, started)
	waitIdle()
	if _, err := runner.Get(started.Benchmark, started.ID); !errors.Is(err, benchmarkrun.ErrNotFound) {
		t.Fatalf("expected the idle run to be removed, got %v", err)
	}
}

// a run stays while a client reads it, for longer than the idle timeout
func TestRunner_ReadKeepsRunAlive(t *testing.T) {
	runner := newTestRunner(t, testIdleTimeout, benchmarkrun.MaxRunningRuns)
	run, err := runner.Start("benchmark.slow", 0)
	if err != nil {
		t.Fatal(err)
	}
	for deadline := time.Now().Add(5 * testIdleTimeout); time.Now().Before(deadline); time.Sleep(testIdleTimeout / 4) {
		if _, err := runner.Get(run.Benchmark, run.ID); err != nil {
			t.Fatalf("the run was removed while a client read it: %v", err)
		}
	}
	if err := runner.Delete(run.Benchmark, run.ID); err != nil {
		t.Fatal(err)
	}
}

// a snapshot write that takes longer than the idle timeout keeps the run
func TestRunner_SnapshotWriteKeepsRunAlive(t *testing.T) {
	runner := newTestRunner(t, testIdleTimeout, benchmarkrun.MaxRunningRuns)
	started, err := runner.Start("benchmark.quick", 0)
	if err != nil {
		t.Fatal(err)
	}
	waitForRun(t, runner, started)

	w := &blockingWriter{started: make(chan struct{}), release: make(chan struct{})}
	writeErr := make(chan error, 1)
	go func() { writeErr <- runner.WriteSnapshot(context.Background(), started.Benchmark, started.ID, w) }()
	<-w.started
	waitIdle()
	close(w.release)
	if err := <-writeErr; err != nil {
		t.Fatalf("snapshot write failed: %v", err)
	}
	if _, err := runner.Get(started.Benchmark, started.ID); err != nil {
		t.Fatalf("the run was removed during the snapshot write: %v", err)
	}
	if err := runner.Delete(started.Benchmark, started.ID); err != nil {
		t.Fatal(err)
	}
}

// `all` runs every benchmark of the workspace mod, as a root benchmark that contains them
func TestRunner_All(t *testing.T) {
	runner := NewTestRunner(t)
	run, err := runner.Start("all", 0)
	if err != nil {
		t.Fatal(err)
	}
	if run.Benchmark != "local.benchmark.root" {
		t.Fatalf("expected the root benchmark, got %q", run.Benchmark)
	}
	for _, name := range []string{"all", "local.benchmark.root"} {
		if _, err := runner.Get(name, run.ID); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
	}
	if _, err := runner.Get("benchmark.quick", run.ID); !errors.Is(err, benchmarkrun.ErrNotFound) {
		t.Fatalf("expected the run of all not to be found under one benchmark, got %v", err)
	}
	if err := runner.Delete("all", run.ID); err != nil {
		t.Fatal(err)
	}
}
