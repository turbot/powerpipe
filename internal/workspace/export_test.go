package workspace

import (
	"context"
	"time"
)

// test access to the mod install watcher

type ModInstallWatcher = modInstallWatcher

func NewModInstallWatcher(workspacePath string, reload func(ctx context.Context)) (*ModInstallWatcher, error) {
	return newModInstallWatcher(workspacePath, reload)
}

func (m *modInstallWatcher) Start(ctx context.Context) { m.start(ctx) }

// SetModInstallTimings sets the timings for watchers created after this call, and returns a function to restore them
func SetModInstallTimings(quietPeriod, maxWait time.Duration) (restore func()) {
	q, m := modInstallQuietPeriod, modInstallMaxWait
	modInstallQuietPeriod, modInstallMaxWait = quietPeriod, maxWait
	return func() { modInstallQuietPeriod, modInstallMaxWait = q, m }
}
