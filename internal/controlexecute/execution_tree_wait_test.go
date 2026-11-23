package controlexecute

import (
	"context"
	"testing"
	"time"

	"golang.org/x/sync/semaphore"
)

// when the context is cancelled (for example by the benchmark timeout) before or during the wait, the wait
// still lasts until the active runs stop, so that their results are not read while they still write them
func TestWaitForActiveRunsToComplete_Cancelled(t *testing.T) {
	for name, cancelAfter := range map[string]time.Duration{
		"cancelled during the wait": 50 * time.Millisecond,
		"cancelled before the wait": 0,
	} {
		t.Run(name, func(t *testing.T) {
			lock := semaphore.NewWeighted(3)
			// one active run
			if err := lock.Acquire(context.Background(), 1); err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if cancelAfter == 0 {
				cancel()
			} else {
				time.AfterFunc(cancelAfter, cancel)
			}
			const runStops = 300 * time.Millisecond
			start := time.Now()
			time.AfterFunc(runStops, func() { lock.Release(1) })

			if err := (&ExecutionTree{}).waitForActiveRunsToComplete(ctx, lock, 3, time.Minute); err != nil {
				t.Fatalf("wait failed: %v", err)
			}
			if elapsed := time.Since(start); elapsed < runStops {
				t.Fatalf("the wait returned after %s, before the active run stopped", elapsed)
			}
		})
	}
}

func TestWaitForActiveRunsToComplete_NotCancelled(t *testing.T) {
	lock := semaphore.NewWeighted(3)
	if err := lock.Acquire(context.Background(), 2); err != nil {
		t.Fatal(err)
	}
	time.AfterFunc(100*time.Millisecond, func() { lock.Release(2) })
	if err := (&ExecutionTree{}).waitForActiveRunsToComplete(context.Background(), lock, 3, time.Minute); err != nil {
		t.Fatalf("wait failed: %v", err)
	}
	// the wait holds every slot
	if lock.TryAcquire(1) {
		t.Fatal("expected the wait to hold every slot")
	}
}

// a run that does not stop after the cancel does not hold the wait for ever: the wait gives up after the timeout
func TestWaitForActiveRunsToComplete_GivesUp(t *testing.T) {
	lock := semaphore.NewWeighted(3)
	// an active run that never stops
	if err := lock.Acquire(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	const timeout = 100 * time.Millisecond
	start := time.Now()
	if err := (&ExecutionTree{}).waitForActiveRunsToComplete(ctx, lock, 3, timeout); err == nil {
		t.Fatal("expected the wait to give up")
	}
	if elapsed := time.Since(start); elapsed < timeout {
		t.Fatalf("the wait gave up after %s, before the timeout", elapsed)
	}
}
