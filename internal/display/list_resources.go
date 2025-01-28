package display

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	"github.com/turbot/powerpipe/internal/resources"
	pworkspace "github.com/turbot/powerpipe/internal/workspace"
)

func ListResources[T modconfig.ModTreeItem](cmd *cobra.Command) {
	ctx := cmd.Context()

	modLocation := viper.GetString(constants.ArgModLocation)
	// build options to specify which blocks we need to load (based on type T
	listOpts := getListLoadWorkspaceOpts[T]()
	w, errAndWarnings := pworkspace.LoadWorkspacePromptingForVariables(ctx, modLocation, listOpts...)
	error_helpers.FailOnError(errAndWarnings.GetError())
	if !w.ModfileExists() {
		error_helpers.FailOnError(localconstants.ErrorNoModDefinition{})
	}

	// get resource filter depending on resource type and output type
	resourceFilter := getListResourceFilter[T](&w.Workspace)
	resourceList, err := workspace.FilterWorkspaceResourcesOfType[T](&w.Workspace, resourceFilter)
	if err != nil {
		error_helpers.ShowErrorWithMessage(ctx, err, "failed to filter resources")
		return
	}

	// TODO K TACTICAL for benchmark list, include detection benchmarks
	// https://github.com/turbot/powerpipe/issues/609
	var empty T
	if _, ok := any(empty).(*resources.Benchmark); !ok {
		printListResult[T](ctx, cmd, resourceList)
	} else {
		// list detcection benchmarks
		resourceFilter := getListResourceFilter[*resources.DetectionBenchmark](&w.Workspace)
		detectionBechmarkList, err := workspace.FilterWorkspaceResourcesOfType[*resources.DetectionBenchmark](&w.Workspace, resourceFilter)
		if err != nil {
			error_helpers.ShowErrorWithMessage(ctx, err, "failed to filter resources")
			return
		}
		// build a separate list of all benchmarks
		var l = make(map[string]modconfig.ModTreeItem)
		for k, v := range resourceList {
			l[k] = v
		}
		for k, v := range detectionBechmarkList {
			l[k] = v
		}
		printListResult[modconfig.ModTreeItem](ctx, cmd, l)
	}

}

func printListResult[T modconfig.ModTreeItem](ctx context.Context, cmd *cobra.Command, resourceList map[string]T) {
	printer, err := printers.GetPrinter[T](cmd)
	if err != nil {
		error_helpers.ShowErrorWithMessage(ctx, err, "failed obtaining printer")
		return
	}
	printableResource := NewPrintableHclResource[T](maps.Values(resourceList))

	err = printer.PrintResource(ctx, printableResource, cmd.OutOrStdout())
	if err != nil {
		error_helpers.ShowErrorWithMessage(ctx, err, "failed when printing")
		return
	}
}

func getListResourceFilter[T modconfig.ModTreeItem](w *workspace.Workspace) workspace.ResourceFilter {
	var res = workspace.ResourceFilter{}

	var empty T
	if _, ok := any(empty).(*resources.Benchmark); ok {

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
func getListLoadWorkspaceOpts[T modconfig.ModTreeItem]() []pworkspace.LoadPowerpipeWorkspaceOption {
	var empty T
	var opts = []pworkspace.LoadPowerpipeWorkspaceOption{
		// pass connections
		pworkspace.WithPipelingConnections(powerpipeconfig.GlobalConfig.PipelingConnections),
		// disable late binding
		pworkspace.WithLateBinding(false),
		pworkspace.WithVariableValidation(false),
	}
	switch any(empty).(type) {
	case *modconfig.Mod:
		opts = append(opts, pworkspace.WithBlockType([]string{schema.BlockTypeMod}))
	case *modconfig.Variable:
		opts = append(opts, pworkspace.WithBlockType([]string{schema.BlockTypeVariable}))
	}
	return opts
}

func ShowResource[T modconfig.ModTreeItem](cmd *cobra.Command, args []string) {
	ctx := cmd.Context()

	modLocation := viper.GetString(constants.ArgModLocation)
	// build options to specify which blocks we need to load (based on type T
	opts := getListLoadWorkspaceOpts[T]()
	w, errAndWarnings := pworkspace.LoadWorkspacePromptingForVariables(ctx, modLocation, opts...)
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

	// tactical - show detection benchmarks using the benchmark command
	// TODO once we remove DetectionBenchmarks tand ResolveTargets returns [], this casting will not be needed
	// https://github.com/turbot/powerpipe/issues/609
	if _, ok := any(targets[0]).(*resources.DetectionBenchmark); ok {
		err = showTarget(ctx, cmd, targets[0].(*resources.DetectionBenchmark))
	} else {
		err = showTarget(ctx, cmd, targets[0].(T))
	}

	if err != nil {
		error_helpers.ShowErrorWithMessage(ctx, err, "failed when printing")
		return
	}
}

func showTarget[T modconfig.ModTreeItem](ctx context.Context, cmd *cobra.Command, target T) error {
	printer, err := printers.GetPrinter[T](cmd)
	if err != nil {
		error_helpers.ShowErrorWithMessage(ctx, err, "failed obtaining printer")
		return nil
	}
	printableResource := NewPrintableHclResource[T]([]T{target})

	err = printer.PrintResource(ctx, printableResource, cmd.OutOrStdout())
	return err
}
