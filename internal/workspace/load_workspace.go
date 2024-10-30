package workspace

import (
	"context"
	"log/slog"
	"time"

	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/pipe-fittings/workspace"
)

func LoadWorkspacePromptingForVariables(ctx context.Context, workspacePath string, opts ...LoadPowerpipeWorkspaceOption) (*PowerpipeWorkspace, error_helpers.ErrorAndWarnings) {
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
	cfg := newLoadPowerpipeWorkspaceConfig()
	for _, o := range opts {
		o(cfg)
	}

	utils.LogTime("w.Load start")
	defer utils.LogTime("w.Load end")

	w = NewPowerpipeWorkspace(workspacePath)
	// check whether the workspace contains a modfile
	// this will determine whether we load files recursively, and create pseudo resources for sql files
	w.SetModfileExists()

	// load the .steampipe ignore file
	if err := w.LoadExclusions(); err != nil {
		return nil, error_helpers.NewErrorsAndWarning(err)
	}

	w.SupportLateBinding = cfg.supportLateBinding
	w.BlockTypeInclusions = cfg.blockTypeInclusions
	w.ValidateVariables = cfg.validateVariables

	// if there is a mod file (or if we are loading resources even with no modfile), load them
	if w.ModfileExists() || !cfg.skipResourceLoadIfNoModfile {
		ew = w.LoadWorkspaceMod(ctx)
	}
	if ew.GetError() != nil {
		return nil, ew
	}

	// verify all runtime dependencies can be resolved
	ew.Error = w.verifyResourceRuntimeDependencies()

	return w, ew
}
