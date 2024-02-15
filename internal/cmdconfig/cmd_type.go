package cmdconfig

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/parse"
	"github.com/turbot/pipe-fittings/schema"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/pipe-fittings/workspace"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

// GetGenericTypeName returns lower case form of type unqualified name
func GetGenericTypeName[T any]() string {
	longName := fmt.Sprintf("%T", *new(T))
	split := strings.Split(longName, ".")
	return strings.ToLower(split[len(split)-1])
}

// ResolveTargets	extracts a list of targets and (if present) query args from the command line parameters
//   - if no resource type is specified in the name, it is added from the command type
//   - validate the resource type specified in the name matches the command type
//   - verify the resource exists in the workspace
//   - if the command type is 'query', the target may be a query string rather than a resource name
//     in this case, convert into a query and add to workspace (to allow for simple snapshot generation)
func ResolveTargets(cmdArgs []string, commandTargetType string, w *workspace.Workspace) ([]modconfig.ModTreeItem, map[string]*modconfig.QueryArgs, error) {
	var targets []modconfig.ModTreeItem
	var queryArgsMap = map[string]*modconfig.QueryArgs{}

	for _, targetName := range cmdArgs {
		// try to parse args out of the invocation (only query supported at present - control too?)
		// for example:
		//		query.my_query("val1","val1")
		// 		query.my_query(my_arg1 => "test", my_arg2 => "test2")
		targetName, argsValues, err := parse.ParseQueryInvocation(targetName)
		if err != nil {
			return nil, nil, err
		}
		if !argsValues.Empty() {
			queryArgsMap[targetName] = argsValues
		}

		target, err := resolveResourceName(targetName, commandTargetType, w)
		if err != nil {
			// if a query resource is not found, treat as a query string
			// for all other resources fail
			if commandTargetType != "query" {
				return nil, nil, err
			}

			// special case handling for query - the arg may be a query string rather than a resource name
			// if a manual query is being run (i.e. not a named query), convert into a query and add to workspace
			// this is to allow us to use existing dashboard execution code
			target, err = ensureSnapshotQueryResource(targetName, w)
			if err != nil {
				return nil, nil, err
			}
			// fall through to add target
		}
		targets = append(targets, target)
	}

	// now check if any args were specified on the command line using the --arg flag
	// if so verify no args were passed in the resource invocation, e.g. query.my_query("val1","val1"
	commandLineQueryArgs, err := getCommandLineQueryArgs()
	if err != nil {
		return nil, nil, err
	}

	// so args were passed using --arg
	if !commandLineQueryArgs.Empty() {
		// verify no args were passed in the resource invocation, e.g. query.my_query("val1","val1"
		if len(queryArgsMap) > 0 {
			return nil, nil, sperr.New("both command line args and query invocation args are set")
		}
		// there must only be 1 target - this should be enforced by cobra as the only
		// command which accept the `--arg` flag (`query run` and `control run`  accept a single argument
		if len(targets) != 1 {
			return nil, nil, sperr.New("'--arg' can only be used with a single target")
		}

		queryArgsMap[targets[0].GetUnqualifiedName()] = commandLineQueryArgs
	}

	return targets, queryArgsMap, nil
}

// convert the given command line query into a query resource and add to workspace
// this is to allow us to use existing dashboard execution code
func ensureSnapshotQueryResource(queryString string, w *workspace.Workspace) (*modconfig.Query, error) {
	// TODO KAI file root???
	// build name
	shortName := "command_line_query"

	// this is NOT a named query - create the query using RawSql
	q := modconfig.NewQuery(&hcl.Block{Type: schema.BlockTypeQuery}, w.Mod, shortName).(*modconfig.Query)
	q.SQL = utils.ToStringPointer(queryString)
	// TODO KAI handle args
	//q.SetArgs(resolvedQuery.QueryArgs())
	// add empty metadata
	q.SetMetadata(&modconfig.ResourceMetadata{})

	// add to the workspace mod so the dashboard execution code can find it
	if err := w.Mod.AddResource(q); err != nil {
		return nil, err
	}
	// return the new resource name
	return q, nil
}

// build a QueryArgs from any args passed using the --args flag
func getCommandLineQueryArgs() (*modconfig.QueryArgs, error) {
	argTuples := viper.GetStringSlice(constants.ArgArg)
	var res = modconfig.NewQueryArgs()

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

// resolveResourceName parses targetName to verify it is a named resource
// - if no resource type is specified in the name, it is added from the command type
// - validate the resource type specified in the name matches the command type
// - verify the resource exists in the workspace
func resolveResourceName(targetName string, commandTargetType string, w *workspace.Workspace) (modconfig.ModTreeItem, error) {
	parsed, err := parseResourceName(targetName, commandTargetType)
	if err != nil {
		return nil, err
	}
	resource, found := w.GetResource(parsed)
	if !found {
		return nil, fmt.Errorf("target resource %s does not exist in mod %s", parsed.ToFullName(), w.Mod.Name())
	}
	// 	safe cast to be on the safe side
	modTreeItem, ok := resource.(modconfig.ModTreeItem)
	if !ok {
		// not expected! this is a coding error, we should never get here
		return nil, fmt.Errorf("target resource %s is not a mod tree item", parsed.ToFullName())
	}
	return modTreeItem, nil
}

func parseResourceName(targetName string, commandTargetType string) (*modconfig.ParsedResourceName, error) {
	parsed := &modconfig.ParsedResourceName{}
	parts := strings.Split(targetName, ".")

	switch len(parts) {
	case 0:
		return nil, sperr.New("empty name passed to resolveResourceName")
	case 1:
		// if no type was specified, deduce the type from the check command used
		parsed.Name = parts[0]
		parsed.ItemType = commandTargetType
	case 2:
		parsed.ItemType = parts[0]
		parsed.Name = parts[1]
	case 3:
		parsed.Mod = parts[0]
		parsed.ItemType = parts[1]
		parsed.Name = parts[2]
	default:
		return nil, sperr.New("invalid name passed to ParseResourceName")
	}

	// now validate the resource type matches the commandTargetType
	if parsed.ItemType == "" {
		parsed.ItemType = commandTargetType
	}
	if parsed.ItemType != commandTargetType {
		return nil, sperr.New(fmt.Sprintf("invalid resource type %s - expected %s", parsed.ItemType, commandTargetType))
	}

	return parsed, nil
}
