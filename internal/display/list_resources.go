package display

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/workspace"
	"golang.org/x/exp/maps"
)

func ListResources[T modconfig.HclResource](cmd *cobra.Command) {
	ctx := cmd.Context()

	modLocation := viper.GetString(constants.ArgModLocation)

	w, errAndWarnings := workspace.LoadWorkspacePromptingForVariables(ctx, modLocation)
	error_helpers.FailOnError(errAndWarnings.GetError())

	resources := workspace.GetWorkspaceResourcesOfType[T](w)

	printer, err := GetPrinter[T](cmd)
	if err != nil {
		error_helpers.ShowErrorWithMessage(ctx, err, "failed obtaining printer")
		return
	}
	printableResource := NewPrintableHclResource[T](maps.Values(resources))

	err = printer.PrintResource(ctx, printableResource, cmd.OutOrStdout())
	if err != nil {
		error_helpers.ShowErrorWithMessage(ctx, err, "failed when printing")
		return
	}
}

func ShowResource[T modconfig.HclResource](cmd *cobra.Command) {
	ctx := cmd.Context()

	modLocation := viper.GetString(constants.ArgModLocation)

	w, errAndWarnings := workspace.LoadWorkspacePromptingForVariables(ctx, modLocation)
	error_helpers.FailOnError(errAndWarnings.GetError())

	dashboards := workspace.GetWorkspaceResourcesOfType[T](w)

	printer, err := GetPrinter[T](cmd)
	if err != nil {
		error_helpers.ShowErrorWithMessage(ctx, err, "failed obtaining printer")
		return
	}
	printableResource := NewPrintableHclResource[T](maps.Values(dashboards))

	err = printer.PrintResource(ctx, printableResource, cmd.OutOrStdout())
	if err != nil {
		error_helpers.ShowErrorWithMessage(ctx, err, "failed when printing")
		return
	}
}
