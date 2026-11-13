// Package benchmarkrun runs benchmarks inside a long running process (powerpipe server), with no client
// connection. A run stays while a client watches it: a run that no client reads for the idle timeout is
// cancelled and removed.
//
// The runs share the loaded workspace and the database connections of the process. A run therefore costs
// only its execution, and its result rows are kept on disk, so many runs fit in the memory of one process,
// where each `powerpipe benchmark run` process loads its own copy of the workspace.
package benchmarkrun

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"

	"github.com/turbot/pipe-fittings/v2/modconfig"
	pworkspace "github.com/turbot/pipe-fittings/v2/workspace"
	"github.com/turbot/powerpipe/internal/cmdconfig"
	"github.com/turbot/powerpipe/internal/controldisplay"
	"github.com/turbot/powerpipe/internal/controlexecute"
	"github.com/turbot/powerpipe/internal/controlstatus"
	"github.com/turbot/powerpipe/internal/db_client"
	"github.com/turbot/powerpipe/internal/resources"
	"github.com/turbot/powerpipe/internal/workspace"
)

// IdleTimeout is how long powerpipe server keeps a run that no client reads.
// A client that waits for a run reads it more often than this.
const IdleTimeout = time.Minute

// MaxRunningRuns is how many runs powerpipe server runs at once. Each run executes its own controls, so
// more runs would compete for the same database connections and memory.
const MaxRunningRuns = 10

type Status string

const (
	StatusRunning   Status = "running"
	StatusComplete  Status = "complete"
	StatusError     Status = "error"
	StatusCancelled Status = "cancelled"
)

// ErrNotFound is returned for a run id that does not exist, or that is not a run of the given benchmark
var ErrNotFound = errors.New("benchmark run not found")

// ErrNotComplete is returned when the snapshot of a run that has not completed is requested
var ErrNotComplete = errors.New("benchmark run is not complete")

// ErrTooManyRuns is returned by Start when the runner already runs its maximum number of runs
var ErrTooManyRuns = errors.New("too many benchmark runs are running, try again when one has finished")

// BenchmarkNotFoundError is returned by Start for a name that does not resolve to a benchmark
type BenchmarkNotFoundError struct {
	Err error
}

func (e *BenchmarkNotFoundError) Error() string { return e.Err.Error() }
func (e *BenchmarkNotFoundError) Unwrap() error { return e.Err }

// Run is the state of one benchmark run
type Run struct {
	ID string `json:"id"`
	// Benchmark is the full name of the benchmark, e.g. aws_compliance.benchmark.cis_v300
	Benchmark string     `json:"benchmark"`
	Status    Status     `json:"status"`
	Error     string     `json:"error,omitempty"`
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	// Summary counts the result rows of each status, for the whole benchmark
	Summary *controlstatus.StatusSummary `json:"summary,omitempty"`

	tree   *controlexecute.ExecutionTree
	cancel context.CancelFunc
	done   chan struct{}
	// closed when the run is removed
	removed chan struct{}
	// when a client last read the run, and the number of snapshot writes in progress
	lastRead time.Time
	writes   int
	// the snapshot writes in progress: remove waits for them before it removes the result rows
	snapshotWrites *sync.WaitGroup
}

// Runner runs benchmarks of one workspace
type Runner struct {
	workspace   *workspace.PowerpipeWorkspace
	client      *db_client.DbClient
	idleTimeout time.Duration
	maxRunning  int

	mu   sync.Mutex
	runs map[string]*Run
	// the runs that are starting or running
	running int
}

// NewRunner returns a runner that runs at most maxRunning runs at once, and removes a run when no client
// has read it for idleTimeout
func NewRunner(w *workspace.PowerpipeWorkspace, client *db_client.DbClient, idleTimeout time.Duration, maxRunning int) *Runner {
	return &Runner{workspace: w, client: client, idleTimeout: idleTimeout, maxRunning: maxRunning, runs: make(map[string]*Run)}
}

// Start starts a run of the benchmark, and returns once its controls are resolved. The benchmark can be
// named in any form that `powerpipe benchmark run` accepts.
// timeout limits the run, as --benchmark-timeout does for `powerpipe benchmark run`. Zero means no limit.
func (r *Runner) Start(benchmark string, timeout time.Duration) (Run, error) {
	name, targets, err := r.resolve(benchmark)
	if err != nil {
		return Run{}, err
	}

	r.mu.Lock()
	if r.running >= r.maxRunning {
		r.mu.Unlock()
		return Run{}, ErrTooManyRuns
	}
	r.running++
	r.mu.Unlock()

	var ctx context.Context
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	tree, err := controlexecute.NewExecutionTree(ctx, r.workspace, r.client, pworkspace.ResourceFilter{}, targets...)
	if err == nil {
		err = tree.StoreRowsOnDisk()
	}
	if err != nil {
		cancel()
		r.mu.Lock()
		r.running--
		r.mu.Unlock()
		return Run{}, err
	}

	now := time.Now()
	run := &Run{
		ID:             newRunID(),
		Benchmark:      name,
		Status:         StatusRunning,
		StartTime:      now,
		tree:           tree,
		cancel:         cancel,
		done:           make(chan struct{}),
		removed:        make(chan struct{}),
		lastRead:       now,
		snapshotWrites: &sync.WaitGroup{},
	}
	r.mu.Lock()
	r.runs[run.ID] = run
	r.mu.Unlock()

	go r.execute(ctx, run)
	go r.removeWhenIdle(run)
	return r.state(run), nil
}

// resolve returns the full name of the benchmark and the targets to run
func (r *Runner) resolve(benchmark string) (string, []modconfig.ModTreeItem, error) {
	targets, err := cmdconfig.ResolveTargets[*resources.Benchmark]([]string{benchmark}, r.workspace)
	if err == nil && len(targets) == 0 {
		err = fmt.Errorf("'%s' is not a benchmark", benchmark)
	}
	if err != nil {
		return "", nil, &BenchmarkNotFoundError{Err: err}
	}
	// a benchmark name resolves to one target: `all` resolves to a root benchmark that contains every benchmark
	return targets[0].Name(), targets, nil
}

// execute runs the controls of the run. A panic fails this run only, not the process.
func (r *Runner) execute(ctx context.Context, run *Run) {
	var err error
	// runs after finish has read the context error
	defer run.cancel()
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("benchmark run failed: %v", p)
		}
		r.finish(ctx, run, err)
	}()
	ctx = controlstatus.AddControlHooksToContext(ctx, controlstatus.NewStatusControlHooks())
	err = run.tree.Execute(ctx)
}

func (r *Runner) finish(ctx context.Context, run *Run, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.running--
	end := time.Now()
	run.EndTime = &end
	if run.tree.Root != nil && run.tree.Root.Summary != nil {
		summary := run.tree.Root.Summary.Status
		run.Summary = &summary
	}
	switch {
	case err != nil:
		run.Status, run.Error = StatusError, err.Error()
	case errors.Is(ctx.Err(), context.Canceled):
		run.Status = StatusCancelled
	case ctx.Err() != nil:
		run.Status, run.Error = StatusError, ctx.Err().Error()
	default:
		run.Status = StatusComplete
	}
	if run.Status == StatusError {
		// as `powerpipe benchmark run` does, so that the run can still be exported
		run.tree.EnsureConsistentStateForExport()
	}
	slog.Info("benchmark run finished", "id", run.ID, "benchmark", run.Benchmark, "status", run.Status, "error", run.Error)
	close(run.done)
}

// find returns the run with the id, if it is a run of the benchmark. The benchmark can be named in any
// form that Start accepts.
func (r *Runner) find(benchmark, id string) (*Run, error) {
	r.mu.Lock()
	run, ok := r.runs[id]
	r.mu.Unlock()
	if !ok {
		return nil, ErrNotFound
	}
	if benchmark != run.Benchmark {
		if name, _, err := r.resolve(benchmark); err != nil || name != run.Benchmark {
			return nil, ErrNotFound
		}
	}
	return run, nil
}

// Get returns the state of the run
func (r *Runner) Get(benchmark, id string) (Run, error) {
	run, err := r.find(benchmark, id)
	if err != nil {
		return Run{}, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.runs[id] != run {
		return Run{}, ErrNotFound
	}
	run.lastRead = time.Now()
	return *run, nil
}

// WriteSnapshot writes the snapshot of a completed run to w, in the same form as
// `powerpipe benchmark run --export=<file>.pps`. A run that timed out or failed has a snapshot too,
// as the command exports one in that case.
func (r *Runner) WriteSnapshot(ctx context.Context, benchmark, id string, w io.Writer) error {
	run, err := r.find(benchmark, id)
	if err != nil {
		return err
	}
	r.mu.Lock()
	switch {
	case r.runs[id] != run || run.Status == StatusCancelled:
		r.mu.Unlock()
		return ErrNotFound
	case run.Status == StatusRunning:
		run.lastRead = time.Now()
		r.mu.Unlock()
		return ErrNotComplete
	}
	// under the lock, so that a remove after this waits for the write
	run.snapshotWrites.Add(1)
	run.writes++
	r.mu.Unlock()
	defer func() {
		r.mu.Lock()
		run.writes--
		run.lastRead = time.Now()
		r.mu.Unlock()
		run.snapshotWrites.Done()
	}()
	return controldisplay.ExportSnapshot(ctx, run.tree, w)
}

// Delete cancels the run if it is running, and removes it
func (r *Runner) Delete(benchmark, id string) error {
	run, err := r.find(benchmark, id)
	if err != nil {
		return err
	}
	return r.remove(run)
}

// remove cancels the run if it is running, waits for it to stop and for the snapshot writes in progress,
// and removes it and its result rows. It returns ErrNotFound if the run is already removed.
func (r *Runner) remove(run *Run) error {
	r.mu.Lock()
	if r.runs[run.ID] != run {
		r.mu.Unlock()
		return ErrNotFound
	}
	delete(r.runs, run.ID)
	close(run.removed)
	r.mu.Unlock()
	run.cancel()
	<-run.done
	run.snapshotWrites.Wait()
	return run.tree.Close()
}

// removeWhenIdle removes the run when no client has read it for the idle timeout.
// A snapshot write in progress counts as a read.
func (r *Runner) removeWhenIdle(run *Run) {
	timer := time.NewTimer(r.idleTimeout)
	defer timer.Stop()
	for {
		select {
		case <-run.removed:
			return
		case <-timer.C:
		}
		r.mu.Lock()
		idle := time.Since(run.lastRead)
		writing := run.writes > 0
		r.mu.Unlock()
		if writing {
			timer.Reset(r.idleTimeout)
			continue
		}
		if idle < r.idleTimeout {
			timer.Reset(r.idleTimeout - idle)
			continue
		}
		slog.Info("removing a benchmark run that no client has read", "id", run.ID, "benchmark", run.Benchmark, "idle", idle)
		if err := r.remove(run); err != nil && !errors.Is(err, ErrNotFound) {
			slog.Warn("failed to remove a benchmark run", "id", run.ID, "error", err)
		}
		return
	}
}

func (r *Runner) state(run *Run) Run {
	r.mu.Lock()
	defer r.mu.Unlock()
	return *run
}

func newRunID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
