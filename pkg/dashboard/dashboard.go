package dashboard

import (
	"context"
	"fmt"

	"github.com/turbot/pipes-fittings/error_helpers"
	"github.com/turbot/powerpipe/internal/initialisation"
	"github.com/turbot/powerpipe/internal/workspace"
	"github.com/turbot/steampipe/pkg/constants"
)

func InitDashboard(ctx context.Context) *initialisation.InitData {
	// initialise
	initData := getInitData(ctx)
	if initData.Result.Error != nil {
		return initData
	}

	// there must be a mod-file
	if !initData.Workspace.ModfileExists() {
		initData.Result.Error = workspace.ErrorNoModDefinition
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
