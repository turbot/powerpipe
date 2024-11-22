package cmdconfig

import (
	"fmt"
	"golang.org/x/exp/maps"
	"strings"

	"github.com/spf13/viper"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/workspace"
	"github.com/turbot/powerpipe/internal/resources"
	pworkspace "github.com/turbot/powerpipe/internal/workspace"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

// TODO once we remove DetectionBenchmarks this should return []T
// https://github.com/turbot/powerpipe/issues/609
func ResolveTargets[T modconfig.ModTreeItem](cmdArgs []string, w *pworkspace.PowerpipeWorkspace) ([]modconfig.ModTreeItem, error) {
	if len(cmdArgs) == 0 {
		return nil, nil
	}

	// special case handling for `benchmark run all`
	targets, err := handleAllArg[T](cmdArgs, w)
	if err != nil {
		return nil, err
	}
	if targets != nil {
		return targets, nil
	}
	if len(cmdArgs) == 1 {
		return resolveSingleTarget[T](cmdArgs[0], w)
	}
	// the only command which supports multiple targets is benchmark run
	return resolveBenchmarkTargets[T](cmdArgs, w)
}

// resolveSingleTarget extracts a single target and (if present) query args from the command line parameters
//   - if no resource type is specified in the name, it is added from the command type
//   - validate the resource type specified in the name matches the command type
//   - verify the resource exists in the workspace
//   - if the command type is 'query', the target may be a query string rather than a resource name
//     in this case, convert into a query and add to workspace (to allow for simple snapshot generation)
//
// TODO K add unit test
// TODO once we remove DetectionBenvhmarks this should return []T
// https://github.com/turbot/powerpipe/issues/609
func resolveSingleTarget[T modconfig.ModTreeItem](cmdArg string, w *pworkspace.PowerpipeWorkspace) ([]modconfig.ModTreeItem, error) {
	typeName := resources.GenericTypeToBlockType[T]()
	var target modconfig.ModTreeItem
	var queryArgs *resources.QueryArgs
	var err error
	target, queryArgs, err = pworkspace.ResolveResourceAndArgsFromSQLString[T](cmdArg, &w.Workspace)
	if err != nil {
		return nil, err
	}
	if helpers.IsNil(target) {
		return nil, fmt.Errorf("'%s.%s' not found in %s (%s)", typeName, cmdArg, w.Mod.Name(), w.Path)
	}
	// TODO KAI CHECK QUERY ARGS LOGIC HERE
	//if queryArgs != nil {
	//	return nil, sperr.New("benchmarks do not support query args")
	//}

	// ok we managed to resolve

	// now check if any args were specified on the command line using the --arg flag
	// if so verify no args were passed in the resource invocation, e.g. query.my_query("val1","val1"
	commandLineQueryArgs, err := getCommandLineQueryArgs()
	if err != nil {
		return nil, err
	}

	// so args were passed using --arg
	if !commandLineQueryArgs.Empty() {
		// verify no args were passed in the resource invocation, e.g. query.my_query("val1","val2")
		if queryArgs != nil {
			return nil, sperr.New("both command line args and query invocation args are set")
		}
	}
	// set query args for target
	if queryArgs != nil {
		// if the target is a query provider set the args
		// (if the target is a dashboard, which is not a query provider,
		// we read the args from viper separately and use to populate the inputs)
		if qp, ok := any(target).(resources.QueryProvider); ok {
			qp.SetArgs(queryArgs)
		}
	}
	// now set the command line args
	if !commandLineQueryArgs.Empty() {
		// if the target is a query provider set the args
		if qp, ok := any(target).(resources.QueryProvider); ok {
			qp.SetArgs(commandLineQueryArgs)
		}
	}
	return []modconfig.ModTreeItem{target}, nil

}

func resolveBenchmarkTargets[T modconfig.ModTreeItem](cmdArgs []string, w *pworkspace.PowerpipeWorkspace) ([]modconfig.ModTreeItem, error) {
	var targets []modconfig.ModTreeItem
	// so there are multiple targets  - this must be the benchmark command, so we do not expect any args
	// verify T is Benchmark (should be enforced by Cobra)
	var empty T
	if _, isBenchmark := (any(empty)).(*resources.ControlBenchmark); !isBenchmark {
		return nil, sperr.New("multiple targets are only supported for benchmarks")
	}

	// now try to resolve targets
	for _, cmdArg := range cmdArgs {
		target, queryArgs, err := pworkspace.ResolveResourceAndArgsFromSQLString[T](cmdArg, &w.Workspace)
		if err != nil {
			return nil, err
		}
		if helpers.IsNil(target) {
			return nil, fmt.Errorf("'%s' not found in %s (%s)", cmdArg, w.Mod.Name(), w.Path)
		}
		if queryArgs != nil {
			return nil, sperr.New("benchmarks do not support query args")
		}
		targets = append(targets, target)
	}

	return targets, nil
}

func handleAllArg[T modconfig.ModTreeItem](args []string, w *pworkspace.PowerpipeWorkspace) ([]modconfig.ModTreeItem, error) {
	// if there is more than 1 arg, "all" is not valid
	if len(args) > 1 {
		// verify that no other benchmarks/controls are given with an all
		if helpers.StringSliceContains(args, "all") {
			return nil, sperr.New("cannot execute 'all' with other benchmarks")
		}
	}

	isAll := len(args) == 1 && args[0] == "all"
	if !isAll {
		return nil, nil
	}
	var empty T
	if _, isBenchmark := (any(empty)).(*resources.ControlBenchmark); !isBenchmark {
		return nil, nil
	}

	// if the arg is "all", we want to execute all _direct_ children of the Mod
	// but NOT children which come from dependency mods
	filter := workspace.ResourceFilter{
		WherePredicate: func(item modconfig.HclResource) bool {
			mti, ok := item.(modconfig.ModItem)
			if !ok {
				return false
			}
			return mti.GetMod().GetShortName() == w.Mod.ShortName
		},
	}
	// TODO K pass workspace interface instead
	targetsMap, err := workspace.FilterWorkspaceResourcesOfType[T](&w.Workspace, filter)
	if err != nil {
		return nil, err
	}

	targets := ToModTreeItemSlice(maps.Values(targetsMap))

	// make a root item to hold the benchmarks
	resolvedItem := resources.NewRootBenchmarkWithChildren(w.Mod, targets).(modconfig.ModTreeItem)

	return []modconfig.ModTreeItem{resolvedItem}, nil

}

// build a QueryArgs from any args passed using the --args flag
func getCommandLineQueryArgs() (*resources.QueryArgs, error) {
	argTuples := viper.GetStringSlice(constants.ArgArg)
	var res = resources.NewQueryArgs()

	if argTuples == nil {
		return res, nil
	}

	for _, argTuple := range argTuples {
		parts := strings.Split(argTuple, "=")
		switch len(parts) {
		case 1:
			// if there is no '=' this must be a positional arg
			if err := res.AddPositionalArgVal(parts[0]); err != nil {
				return nil, err
			}

		case 2:
			argName := parts[0]
			argValue := parts[1]

			if err := res.SetNamedArgVal(argName, argValue); err != nil {
				return nil, err
			}
		default:
			return nil, sperr.New("invalid arg format: %s", argTuple)
		}
	}
	// we should not have both positional and named args
	if len(res.ArgMap) > 0 && len(res.ArgList) > 0 {
		return nil, sperr.New("cannot mix positional and named args")
	}
	return res, nil

}

func ToModTreeItemSlice[T modconfig.ModTreeItem](items []T) []modconfig.ModTreeItem {
	var res []modconfig.ModTreeItem
	for _, item := range items {
		res = append(res, item)
	}
	return res
}
