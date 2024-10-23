package display

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/modconfig/dashboard"
	"golang.org/x/exp/maps"

	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/printers"
	"github.com/turbot/pipe-fittings/schema"
	"github.com/turbot/pipe-fittings/workspace"
	localcmdconfig "github.com/turbot/powerpipe/internal/cmdconfig"
	localconstants "github.com/turbot/powerpipe/internal/constants"
	"github.com/turbot/powerpipe/internal/powerpipeconfig"
)

func ListResources[T modconfig.ModTreeItem](cmd *cobra.Command) {
	ctx := cmd.Context()

	modLocation := viper.GetString(constants.ArgModLocation)
	// build options to specify which blocks we need to load (based on type T
	opts := getListLoadWorkspaceOpts[T]()
	w, errAndWarnings := workspace.LoadWorkspacePromptingForVariables(ctx, modLocation, opts...)
	error_helpers.FailOnError(errAndWarnings.GetError())

	// get resource filter depending on resource type and output type
	resourceFilter := getListResourceFilter[T](w)
	resources, err := workspace.FilterWorkspaceResourcesOfType[T](w, resourceFilter)
	if err != nil {
		error_helpers.ShowErrorWithMessage(ctx, err, "failed to filter resources")
		return
	}

	if !w.ModfileExists() {
		error_helpers.FailOnError(localconstants.ErrorNoModDefinition{})
	}

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

func getListResourceFilter[T modconfig.ModTreeItem](w *workspace.Workspace) workspace.ResourceFilter {
	var res = workspace.ResourceFilter{}

	var empty T
	if _, ok := any(empty).(*dashboard.Benchmark); ok {

		// if T is benchmark, and if output is pretty or plain, only show top level benchmarks
		if viper.GetString(constants.ArgOutput) == constants.OutputFormatPretty || viper.GetString(constants.ArgOutput) == constants.OutputFormatPlain {
			// build a lookup of mod names to filter on
			var modNames = map[string]struct{}{}
			for _, mod := range w.Mods {
				modNames[mod.Name()] = struct{}{}
			}

			// add a predicate which returns true only if the resources parent is one of these mods
			res.WherePredicate = func(item modconfig.HclResource) bool {
				mti, ok := item.(modconfig.ModTreeItem)
				if !ok {
					return false
				}

				parents := mti.GetParents()
				if len(parents) == 0 {
					return false
				}
				_, inTargetMod := modNames[parents[0].Name()]
				return inTargetMod
			}
		}
	}

	return res
}

// build LoadWorkspaceOptions to specify which blocks we need to load (based on type T)
func getListLoadWorkspaceOpts[T modconfig.ModTreeItem]() []workspace.LoadWorkspaceOption {
	var empty T
	var opts = []workspace.LoadWorkspaceOption{
		// pass connections
		workspace.WithPipelingConnections(powerpipeconfig.GlobalConfig.PipelingConnections),
		// disable late binding
		workspace.WithLateBinding(false),
		workspace.WithVariableValidation(false),
	}
	switch any(empty).(type) {
	case *modconfig.Mod:
		opts = append(opts, workspace.WithBlockType([]string{schema.BlockTypeMod}))
	case *modconfig.Variable:
		opts = append(opts, workspace.WithBlockType([]string{schema.BlockTypeVariable}))
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
