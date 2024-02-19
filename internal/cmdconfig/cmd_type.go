package cmdconfig

import (
	"strings"

	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/schema"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/pipe-fittings/workspace"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

// ResolveTarget	extracts a list of targets and (if present) query args from the command line parameters
//   - if no resource type is specified in the name, it is added from the command type
//   - validate the resource type specified in the name matches the command type
//   - verify the resource exists in the workspace
//   - if the command type is 'query', the target may be a query string rather than a resource name
//     in this case, convert into a query and add to workspace (to allow for simple snapshot generation)
func ResolveTarget[T modconfig.ModTreeItem](cmdArgs []string, w *workspace.Workspace) (T, error) {
	var targets []modconfig.ModTreeItem

	typeName := utils.GetGenericTypeName[T]()
	// special case for variable
	if typeName == schema.BlockTypeVariable {
		// variables are named var.xxxx, not variable.xxxx
		typeName = schema.AttributeVar
	}

	var empty T
	targets, argsMap, err := workspace.GetResourcesFromArgs[T](cmdArgs, w)
	if err != nil {
		return empty, err
	}
	if len(targets) == 0 {
		return empty, sperr.New("could not resolve %s '%s'", utils.GetGenericTypeName[T](), cmdArgs[0])
	}
	// we only support a single target - should be enforced by cobra
	if len(targets) != 1 {
		return empty, sperr.New("only a single target is supported")
	}

	target, ok := targets[0].(T)
	if !ok {
		return empty, sperr.New("target '%s' is not of the expected type '%s'", targets[0].GetUnqualifiedName(), typeName)

	}
	queryArgs := argsMap[target.GetUnqualifiedName()]

	// if the command type is 'query', the target may be a query string rather than a resource name

	// now check if any args were specified on the command line using the --arg flag
	// if so verify no args were passed in the resource invocation, e.g. query.my_query("val1","val1"
	commandLineQueryArgs, err := getCommandLineQueryArgs()
	if err != nil {
		return empty, err
	}

	// so args were passed using --arg
	if !commandLineQueryArgs.Empty() {
		// verify no args were passed in the resource invocation, e.g. query.my_query("val1","val1"
		if queryArgs != nil {
			return empty, sperr.New("both command line args and query invocation args are set")
		}
		queryArgs = commandLineQueryArgs

	}

	if queryArgs != nil {
		if qp, ok := any(target).(modconfig.QueryProvider); ok {
			qp.SetArgs(queryArgs)
		} else {
			// args provided but target is not a query provider
			return empty, sperr.New("args provided but target is not a query provider")
		}

	}
	return target, nil
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
