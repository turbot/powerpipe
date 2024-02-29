package display

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	localconstants "github.com/turbot/powerpipe/internal/constants"
	"golang.org/x/exp/maps"

	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/printers"
	"github.com/turbot/pipe-fittings/schema"
	"github.com/turbot/pipe-fittings/workspace"
	localcmdconfig "github.com/turbot/powerpipe/internal/cmdconfig"
)

func ListResources[T modconfig.ModTreeItem](cmd *cobra.Command) {
	ctx := cmd.Context()

	modLocation := viper.GetString(constants.ArgModLocation)
	// build options to specify which blocks we need to load (based on type T
	opts := getListLoadWorkspaceOpts[T]()
	w, errAndWarnings := workspace.LoadWorkspacePromptingForVariables(ctx, modLocation, opts...)
	error_helpers.FailOnError(errAndWarnings.GetError())

	if !w.ModfileExists() {
		error_helpers.FailOnError(localconstants.ErrorNoModDefinition{})
	}

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

// build LoadWorkspaceOptions to specify which blocks we need to load (based on type T)
func getListLoadWorkspaceOpts[T modconfig.ModTreeItem]() []workspace.LoadWorkspaceOption {
	var empty T
	var opts = []workspace.LoadWorkspaceOption{workspace.WithVariableValidation(false)}
	switch any(empty).(type) {
	case *modconfig.Mod:
		opts = append(opts, workspace.WithBlockType([]string{schema.BlockTypeMod}))
	}
	return opts
}

func ShowResource[T modconfig.ModTreeItem](cmd *cobra.Command, args []string) {
	ctx := cmd.Context()

	modLocation := viper.GetString(constants.ArgModLocation)
	// build options to specify which blocks we need to load (based on type T
	opts := getListLoadWorkspaceOpts[T]()
	w, errAndWarnings := workspace.LoadWorkspacePromptingForVariables(ctx, modLocation, opts...)
	error_helpers.FailOnError(errAndWarnings.GetError())
	if !w.ModfileExists() {
		error_helpers.FailOnError(localconstants.ErrorNoModDefinition{})
	}

	targets, err := localcmdconfig.ResolveTargets[T](args, w)
	error_helpers.FailOnError(err)

	// show only supports a single target (should be enforced by cobra)
	if len(targets) > 1 {
		// not expected
		error_helpers.FailOnError(fmt.Errorf("show command only supports a single target"))
		return
	}
	target := targets[0].(T)

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
