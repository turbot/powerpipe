package display

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/printers"
	"golang.org/x/exp/maps"

	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/schema"
	"github.com/turbot/pipe-fittings/workspace"
	localcmdconfig "github.com/turbot/powerpipe/internal/cmdconfig"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

func ListResources[T modconfig.HclResource](cmd *cobra.Command) {
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

func ShowResource[T modconfig.HclResource](cmd *cobra.Command, args []string) {
	ctx := cmd.Context()

	modLocation := viper.GetString(constants.ArgModLocation)

	w, errAndWarnings := workspace.LoadWorkspacePromptingForVariables(ctx, modLocation)
	error_helpers.FailOnError(errAndWarnings.GetError())

	typeName := localcmdconfig.GetGenericTypeName[T]()
	// special case for variable
	if typeName == schema.BlockTypeVariable {
		// variables are named var.xxxx, not variable.xxxx
		typeName = schema.AttributeVar
	}

	targets, _, err := localcmdconfig.ResolveTargets(args, typeName, w)
	error_helpers.FailOnError(err)

	var target T = targets[0].(T)
	// we expect a single target - this will be enforced by cobra
	if len(targets) != 1 {
		error_helpers.FailOnError(sperr.New("expected a single target, got %d", len(targets)))
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
