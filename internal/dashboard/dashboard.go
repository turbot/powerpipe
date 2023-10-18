package dashboard

import (
	"context"
	"fmt"

	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/workspace"
	"github.com/turbot/powerpipe/internal/initialisation"
)

func InitDashboard(ctx context.Context) *initialisation.InitData {
	// initialise
	initData := getInitData(ctx)
	if initData.Result.Error != nil {
		return initData
	}

	// there must be a mod-file
	if !initData.Workspace.ModfileExists() {
		error_helpers.ShowWarning("Could not find mod definition file in the current directory tree.")
	}

	return initData
}

func getInitData(ctx context.Context) *initialisation.InitData {
	w, errAndWarnings := workspace.LoadWorkspacePromptingForVariables(ctx)
	if errAndWarnings.GetError() != nil {
		return initialisation.NewErrorInitData(fmt.Errorf("failed to load workspace: %s", error_helpers.HandleCancelError(errAndWarnings.GetError()).Error()))
	}

	i := initialisation.NewInitData()
	i.Workspace = w
	i.Result.Warnings = errAndWarnings.Warnings
	i.Init(ctx, constants.InvokerDashboard)

	return i
}
