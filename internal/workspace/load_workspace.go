package workspace

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/turbot/pipe-fittings/v2/error_helpers"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/pipe-fittings/v2/utils"
	"github.com/turbot/pipe-fittings/v2/workspace"
	"github.com/turbot/powerpipe/internal/dashboardevents"
	"github.com/turbot/powerpipe/internal/resources"
	"github.com/turbot/powerpipe/internal/timing"
)

// WorkspaceProvider is an interface that both PowerpipeWorkspace and LazyWorkspace implement.
// This allows code to work with either workspace type.
type WorkspaceProvider interface {
	// GetResource retrieves a resource by parsed name
	GetResource(parsedName *modconfig.ParsedResourceName) (modconfig.HclResource, bool)
	// Close cleans up the workspace
	Close()
	// IsLazy returns true if this is a lazy-loading workspace
	IsLazy() bool
	// LoadDashboard loads a dashboard by name
	LoadDashboard(ctx context.Context, name string) (*resources.Dashboard, error)
	// LoadBenchmark loads a benchmark by name
	LoadBenchmark(ctx context.Context, name string) (modconfig.ModTreeItem, error)
	// LoadBenchmarkForExecution loads a benchmark with children properly resolved for execution.
	// For lazy workspaces, this resolves child references. For eager workspaces, this is the same as LoadBenchmark.
	LoadBenchmarkForExecution(ctx context.Context, name string) (modconfig.ModTreeItem, error)
}

// DashboardServerWorkspace extends WorkspaceProvider with methods needed by the dashboard server.
// Both PowerpipeWorkspace and LazyWorkspace implement this interface.
type DashboardServerWorkspace interface {
	WorkspaceProvider
	// RegisterDashboardEventHandler registers a handler for dashboard events
	RegisterDashboardEventHandler(ctx context.Context, handler dashboardevents.DashboardEventHandler)
	// SetupWatcher sets up file watching
	SetupWatcher(ctx context.Context, errCallback func(context.Context, error)) error
	// GetModResources returns the mod resources
	GetModResources() modconfig.ModResources
	// PublishDashboardEvent publishes a dashboard event
	PublishDashboardEvent(ctx context.Context, event dashboardevents.DashboardEvent)
}

func LoadWorkspacePromptingForVariables(ctx context.Context, workspacePath string, opts ...LoadPowerpipeWorkspaceOption) (*PowerpipeWorkspace, error_helpers.ErrorAndWarnings) {
	defer timing.Track("LoadWorkspacePromptingForVariables")()
	t := time.Now()
	defer func() {
		slog.Debug("Workspace load complete", "duration (ms)", time.Since(t).Milliseconds())
	}()
	w, errAndWarnings := Load(ctx, workspacePath, opts...)
	if errAndWarnings.GetError() == nil {
		return w, errAndWarnings
	}

	// kif there wqs an error check if it was a missing variable error and if so prompt for variables
	if err := workspace.HandleWorkspaceLoadError(ctx, errAndWarnings.GetError(), workspacePath); err != nil {
		return nil, error_helpers.NewErrorsAndWarning(err)
	}

	// ok we should have all variables now - reload workspace
	return Load(ctx, workspacePath, opts...)
}

// Load_ creates a Workspace and loads the workspace mod

func Load(ctx context.Context, workspacePath string, opts ...LoadPowerpipeWorkspaceOption) (w *PowerpipeWorkspace, ew error_helpers.ErrorAndWarnings) {
	defer timing.Track("workspace.Load")()

	cfg := newLoadPowerpipeWorkspaceConfig()
	for _, o := range opts {
		o(cfg)
	}

	utils.LogTime("w.Load start")
	defer utils.LogTime("w.Load end")

	w = NewPowerpipeWorkspace(workspacePath)
	// check whether the workspace contains a modfile
	// this will determine whether we load files recursively, and create pseudo resources for sql files
	func() {
		defer timing.Track("workspace.SetModfileExists")()
		w.SetModfileExists()
	}()

	// load the .steampipe ignore file
	func() {
		defer timing.Track("workspace.LoadExclusions")()
		if err := w.LoadExclusions(); err != nil {
			ew = error_helpers.NewErrorsAndWarning(err)
		}
	}()
	if ew.GetError() != nil {
		return nil, ew
	}

	w.SupportLateBinding = cfg.supportLateBinding
	w.BlockTypeInclusions = cfg.blockTypeInclusions
	w.ValidateVariables = cfg.validateVariables
	w.PipelingConnections = cfg.pipelingConnections

	// if there is a mod file (or if we are loading resources even with no modfile), load them
	if w.ModfileExists() || !cfg.skipResourceLoadIfNoModfile {
		func() {
			defer timing.Track("workspace.LoadWorkspaceMod")()
			ew = w.LoadWorkspaceMod(ctx)
		}()
	}
	if ew.GetError() != nil {
		return nil, ew
	}

	// verify all runtime dependencies can be resolved
	func() {
		defer timing.Track("workspace.verifyResourceRuntimeDependencies")()
		ew.Error = w.verifyResourceRuntimeDependencies()
	}()

	return w, ew
}

// LoadLazy creates a LazyWorkspace that loads resources on-demand.
// This provides faster startup and lower memory usage than the standard Load().
func LoadLazy(ctx context.Context, workspacePath string, opts ...LoadPowerpipeWorkspaceOption) (*LazyWorkspace, error) {
	defer timing.Track("workspace.LoadLazy")()

	t := time.Now()
	defer func() {
		slog.Debug("Lazy workspace load complete", "duration (ms)", time.Since(t).Milliseconds())
	}()

	cfg := newLoadPowerpipeWorkspaceConfig()
	for _, o := range opts {
		o(cfg)
	}

	utils.LogTime("w.LoadLazy start")
	defer utils.LogTime("w.LoadLazy end")

	lw, err := NewLazyWorkspace(ctx, workspacePath, cfg.lazyLoadConfig)
	if err != nil {
		return nil, err
	}

	// Set modfile exists flag on the embedded workspace
	// This is needed for commands that require a modfile
	lw.PowerpipeWorkspace.SetModfileExists()

	// Start background resolution for variable references and templates
	lw.StartBackgroundResolution()

	// For dashboard server scenarios (like Pipes), wait briefly for top-level resources
	// to resolve their tags/metadata. This ensures the initial available_dashboards
	// message has populated tags instead of empty objects.
	// We use a short timeout (1s) so we don't block too long, but long enough for
	// most top-level resources to resolve (they're prioritized by the background resolver).
	if !lw.WaitForResolution(1000) {
		// Timeout reached, but continue anyway - partial resolution is better than none
		slog.Info("initial background resolution timeout reached, continuing with partial metadata")
	}

	return lw, nil
}

// LoadAuto loads a workspace, using lazy loading if enabled in options.
// Returns a WorkspaceProvider interface that can be either type.
// If lazy loading fails due to duplicate mod versions (diamond dependency),
// it automatically falls back to eager loading.
func LoadAuto(ctx context.Context, workspacePath string, opts ...LoadPowerpipeWorkspaceOption) (WorkspaceProvider, error_helpers.ErrorAndWarnings) {
	cfg := newLoadPowerpipeWorkspaceConfig()
	for _, o := range opts {
		o(cfg)
	}

	if cfg.lazyLoad {
		lw, err := LoadLazy(ctx, workspacePath, opts...)
		if err != nil {
			// Check if this is a duplicate mod versions error - if so, fall back to eager loading
			if errors.Is(err, ErrDuplicateModVersions) {
				slog.Info("Falling back to eager loading due to duplicate mod versions")
				return Load(ctx, workspacePath, opts...)
			}
			return nil, error_helpers.NewErrorsAndWarning(err)
		}
		return lw, error_helpers.ErrorAndWarnings{}
	}

	return Load(ctx, workspacePath, opts...)
}
