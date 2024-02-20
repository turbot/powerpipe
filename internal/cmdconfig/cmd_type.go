package cmdconfig

import (
	"strings"

	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/modconfig"
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

	var empty T
	if len(cmdArgs) > 1 {
		return empty, sperr.New("only a single target is supported")
	}
	target, queryArgs, err := workspace.ResolveResourceAndArgsFromSQLString[T](cmdArgs[0], w)
	if err != nil {
		return empty, err
	}

	// if the command type is 'query', the target may be a query string rather than a resource name

	// now check if any args were specified on the command line using the --arg flag
	// if so verify no args were passed in the resource invocation, e.g. query.my_query("val1","val1"
	commandLineQueryArgs, err := getCommandLineQueryArgs()
	if err != nil {
		return empty, err
	}

	// so args were passed using --arg
	if !commandLineQueryArgs.Empty() {
		// verify no args were passed in the resource invocation, e.g. query.my_query("val1","val2")
		if queryArgs != nil {
			return empty, sperr.New("both command line args and query invocation args are set")
		}
		queryArgs = commandLineQueryArgs

	}

	if queryArgs != nil {
		// if the target is a query provider set the args
		// (if the target is a dashboard, which i snot a query provider,
		// we read the args from viper separately and use to populate the inputs)
		if qp, ok := any(target).(modconfig.QueryProvider); ok {
			qp.SetArgs(queryArgs)
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
