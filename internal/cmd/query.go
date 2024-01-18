package cmd

import localconstants "github.com/turbot/powerpipe/internal/constants"

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thediveo/enumflag/v2"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/cmdconfig"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/steampipeconfig"
	"github.com/turbot/pipe-fittings/workspace"
	localcmdconfig "github.com/turbot/powerpipe/internal/cmdconfig"
	"github.com/turbot/powerpipe/internal/dashboardexecute"
	"github.com/turbot/powerpipe/internal/display"
	"github.com/turbot/powerpipe/internal/initialisation"
	"github.com/turbot/powerpipe/internal/queryresult"
	"github.com/turbot/steampipe-plugin-sdk/v5/logging"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

// variable used to assign the output mode flag
var queryOutputMode = localconstants.QueryOutputModeSnapshot

func queryRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "run [flags] [query]",
		TraverseChildren: true,
		Args:             cobra.ExactArgs(1),
		Run:              queryRun,
		Short:            "Run a named query",
		Long: `Runs the named query.

The current mod is the working directory, or the directory specified by the --mod-location flag.`,
	}

	cmdconfig.OnCmd(cmd).
		AddCloudFlags().
		AddWorkspaceDatabaseFlag().
		AddModLocationFlag().
		AddBoolFlag(constants.ArgHelp, false, "Help for query", cmdconfig.FlagOptions.WithShortHand("h")).
		AddStringSliceFlag(constants.ArgSearchPath, nil, "Set a custom search_path for the steampipe user for a query session (comma-separated)").
		AddStringSliceFlag(constants.ArgSearchPathPrefix, nil, "Set a prefix to the current search path for a query session (comma-separated)").
		AddBoolFlag(constants.ArgSnapshot, false, "Create snapshot in Turbot Pipes with the default (workspace) visibility").
		AddBoolFlag(constants.ArgShare, false, "Create snapshot in Turbot Pipes with 'anyone_with_link' visibility").
		AddStringFlag(constants.ArgSnapshotLocation, "", "The location to write snapshots - either a local file path or a Turbot Pipes workspace").
		AddStringFlag(constants.ArgSnapshotTitle, "", "The title to give a snapshot").
		// Define the CLI flag parameters for wrapped enum flag.
		AddVarFlag(enumflag.New(&queryOutputMode, constants.ArgOutput, localconstants.QueryOutputModeIds, enumflag.EnumCaseInsensitive),
			constants.ArgOutput,
			fmt.Sprintf("Output format; one of: %s", strings.Join(localconstants.FlagValues(localconstants.QueryOutputModeIds), ", "))).
		// NOTE: use StringArrayFlag for ArgQueryInput, not StringSliceFlag
		// Cobra will interpret values passed to a StringSliceFlag as CSV, where args passed to StringArrayFlag are not parsed and used raw
		AddStringArrayFlag(constants.ArgArg, nil, "Specify the value of a query argument").
		AddStringArrayFlag(constants.ArgSnapshotTag, nil, "Specify tags to set on the snapshot")

	return cmd
}

func queryRun(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()

	var err error
	logging.LogTime("queryRun start")
	defer func() {
		logging.LogTime("queryRun end")
		if r := recover(); r != nil {
			err = helpers.ToError(r)
			error_helpers.ShowError(ctx, err)

		}
		setExitCodeForQueryError(err)
	}()

	// first check whether a single query name has been passed as an arg
	error_helpers.FailOnError(validateQueryArgs(ctx))

	// if diagnostic mode is set, print out config and return
	if _, ok := os.LookupEnv(localconstants.EnvConfigDump); ok {
		localcmdconfig.DisplayConfig()
		return
	}

	initData := initialisation.NewInitData(ctx, "query", args...)
	// shutdown the service on exit
	defer initData.Cleanup(ctx)
	error_helpers.FailOnError(initData.Result.Error)

	// if there is a usage warning we display it
	initData.Result.DisplayMessages()

	// TODO check cancellation
	// start cancel handler to intercept interrupts and cancel the context
	// NOTE: use the initData Cancel function to ensure any initialisation is cancelled if needed
	//contexthelpers.StartCancelHandler(initData.Cancel)

	if err := initData.Result.Error; err != nil {
		exitCode = constants.ExitCodeInitializationFailed
		error_helpers.FailOnError(err)
	}

	// if there are args, put them into viper for retrieval by the dashboard execution
	if initData.Args != nil {
		queryArgs, ok := initData.Args[initData.Targets[0].GetUnqualifiedName()]
		if !ok {
			// not expected
			error_helpers.FailOnError(sperr.New("arguments have been provided but none found for the query"))
		}
		viper.Set(constants.ConfigKeyQueryArgs, queryArgs)
	}
	// so a dashboard name was specified - just call GenerateSnapshot
	snap, err := dashboardexecute.GenerateSnapshot(ctx, initData, nil)
	if err != nil {
		exitCode = constants.ExitCodeSnapshotCreationFailed
		error_helpers.FailOnError(err)
	}

	// TODO KAI FIX ME
	// set the filename root for the snapshot (in case needed)
	//if !existingResource {
	//	snap.FileNameRoot = "query"
	//}

	// display the result
	switch viper.GetString(constants.ArgOutput) {
	case constants.OutputFormatNone:
		// do nothing
	case constants.OutputFormatSnapshot, constants.OutputFormatSnapshotShort:
		// if the format is snapshot, just dump it out
		jsonOutput, err := json.MarshalIndent(snap, "", "  ")
		if err != nil {
			error_helpers.FailOnErrorWithMessage(err, "failed to display result as snapshot")
		}
		fmt.Println(string(jsonOutput)) //nolint:forbidigo // intentional use of fmt
	default:
		// otherwise convert the snapshot into a query result
		result, err := snapshotToQueryResult(snap)
		error_helpers.FailOnErrorWithMessage(err, "failed to display result as snapshot")
		display.ShowOutput(ctx, result, display.WithTimingDisabled())
	}

	// share the snapshot if necessary
	err = publishSnapshotIfNeeded(ctx, snap)
	if err != nil {
		exitCode = constants.ExitCodeSnapshotUploadFailed
		error_helpers.FailOnErrorWithMessage(err, fmt.Sprintf("failed to publish snapshot to %s", viper.GetString(constants.ArgSnapshotLocation)))
	}

	// export the result if necessary
	exportArgs := viper.GetStringSlice(constants.ArgExport)
	exportMsg, err := initData.ExportManager.DoExport(ctx, snap.FileNameRoot, snap, exportArgs)
	error_helpers.FailOnErrorWithMessage(err, "failed to export snapshot")
	// print the location where the file is exported
	if len(exportMsg) > 0 && viper.GetBool(constants.ArgProgress) {
		fmt.Printf("\n")                           //nolint:forbidigo // intentional use of fmt
		fmt.Println(strings.Join(exportMsg, "\n")) //nolint:forbidigo // intentional use of fmt
		fmt.Printf("\n")                           //nolint:forbidigo // intentional use of fmt
	}
}

// validate the args and extract a query name, if provided
func validateQueryArgs(ctx context.Context) error {
	err := localcmdconfig.ValidateSnapshotArgs(ctx)
	if err != nil {
		return err
	}

	if viper.IsSet(constants.ArgSearchPath) && viper.IsSet(constants.ArgSearchPathPrefix) {
		return fmt.Errorf("only one of --search-path or --search-path-prefix may be set")
	}

	// only 1 of 'share' and 'snapshot' may be set
	share := viper.GetBool(constants.ArgShare)
	snapshot := viper.GetBool(constants.ArgSnapshot)
	if share && snapshot {
		return fmt.Errorf("only one of --share or --snapshot may be set")
	}

	return nil
}

func setExitCodeForQueryError(err error) {
	// if exit code already set, leave as is
	if exitCode != 0 || err == nil {
		return
	}

	if errors.Is(err, workspace.ErrorNoModDefinition) {
		exitCode = constants.ExitCodeNoModFile
	} else {
		exitCode = constants.ExitCodeUnknownErrorPanic
	}
}

func snapshotToQueryResult(snap *steampipeconfig.SteampipeSnapshot) (*queryresult.Result, error) {
	// the table of a snapshot query has a fixed name
	tablePanel, ok := snap.Panels[modconfig.SnapshotQueryTableName]
	if !ok {
		return nil, sperr.New("dashboard does not contain table result for query")
	}
	chartRun := tablePanel.(*dashboardexecute.LeafRun)
	if !ok {
		return nil, sperr.New("failed to read query result from snapshot")
	}
	// check for error
	if err := chartRun.GetError(); err != nil {
		return nil, error_helpers.DecodePgError(err)
	}

	res := queryresult.NewResult(chartRun.Data.Columns)

	// start a goroutine to stream the results as rows
	go func() {
		for _, d := range chartRun.Data.Rows {
			// we need to allocate a new slice everytime, since this gets read
			// asynchronously on the other end and we need to make sure that we don't overwrite
			// data already sent
			rowVals := make([]interface{}, len(chartRun.Data.Columns))
			for i, c := range chartRun.Data.Columns {
				rowVals[i] = d[c.Name]
			}
			res.StreamRow(rowVals)
		}
		res.TimingResult <- chartRun.TimingResult
		res.Close()
	}()

	return res, nil
}
