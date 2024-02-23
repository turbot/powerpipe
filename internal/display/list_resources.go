package display

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/utils"
	"golang.org/x/exp/maps"

	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/printers"
	"github.com/turbot/pipe-fittings/workspace"
	localcmdconfig "github.com/turbot/powerpipe/internal/cmdconfig"
)

func ListResources[T modconfig.ModTreeItem](cmd *cobra.Command) {
	ctx := cmd.Context()

	modLocation := viper.GetString(constants.ArgModLocation)

	w, errAndWarnings := workspace.LoadWorkspacePromptingForVariables(ctx, modLocation)
	error_helpers.FailOnError(errAndWarnings.GetError())

	resources := workspace.GetWorkspaceResourcesOfType[T](w)

	printer, err := printers.GetPrinter[T](cmd)
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

func ShowResource[T modconfig.ModTreeItem](cmd *cobra.Command, args []string) {
	ctx := cmd.Context()

	modLocation := viper.GetString(constants.ArgModLocation)

	w, errAndWarnings := workspace.LoadWorkspacePromptingForVariables(ctx, modLocation)
	error_helpers.FailOnError(errAndWarnings.GetError())

	target, err := localcmdconfig.ResolveTarget[T](args, w)
	error_helpers.FailOnError(err)
	if helpers.IsNil(target) {
		error_helpers.FailOnError(fmt.Errorf("%s '%s' not found", utils.GetGenericTypeName[T](), args[0]))
		return
	}

	printer, err := printers.GetPrinter[T](cmd)
	if err != nil {
		error_helpers.ShowErrorWithMessage(ctx, err, "failed obtaining printer")
		return
	}
	printableResource := NewPrintableHclResource[T]([]T{target})

	err = printer.PrintResource(ctx, printableResource, cmd.OutOrStdout())
	if err != nil {
		error_helpers.ShowErrorWithMessage(ctx, err, "failed when printing")
		return
	}
}
