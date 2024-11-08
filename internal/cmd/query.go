package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thediveo/enumflag/v2"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/cmdconfig"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/export"
	"github.com/turbot/pipe-fittings/steampipeconfig"
	"github.com/turbot/pipe-fittings/workspace"
	localcmdconfig "github.com/turbot/powerpipe/internal/cmdconfig"
	localconstants "github.com/turbot/powerpipe/internal/constants"
	"github.com/turbot/powerpipe/internal/dashboardexecute"
	"github.com/turbot/powerpipe/internal/display"
	"github.com/turbot/powerpipe/internal/initialisation"
	"github.com/turbot/powerpipe/internal/queryresult"
	"github.com/turbot/powerpipe/internal/resources"
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
		AddModLocationFlag().
		// NOTE: use StringArrayFlag for ArgQueryInput, not StringSliceFlag
		// Cobra will interpret values passed to a StringSliceFlag as CSV, where args passed to StringArrayFlag are not parsed and used raw
		AddStringArrayFlag(constants.ArgArg, nil, "Specify the value of a query argument").
		AddStringFlag(constants.ArgDatabase, "", "Turbot Pipes workspace database", localcmdconfig.Deprecated("see https://powerpipe.io/docs/run#selecting-a-database for the new syntax")).
		AddIntFlag(constants.ArgDatabaseQueryTimeout, localconstants.DatabaseDefaultQueryTimeout, "The query timeout").
		AddStringSliceFlag(constants.ArgExport, nil, "Export output to file, supported formats: csv, html, json, md, nunit3, pps (snapshot), asff").
		AddBoolFlag(constants.ArgHeader, true, "Include column headers for csv and table output").
		AddBoolFlag(constants.ArgHelp, false, "Help for query", cmdconfig.FlagOptions.WithShortHand("h")).
		AddBoolFlag(constants.ArgInput, true, "Enable interactive prompts").
		// Define the CLI flag parameters for wrapped enum flag.
		AddVarFlag(enumflag.New(&queryOutputMode, constants.ArgOutput, localconstants.QueryOutputModeIds, enumflag.EnumCaseInsensitive),
			constants.ArgOutput,
			fmt.Sprintf("Output format; one of: %s", strings.Join(constants.FlagValues(localconstants.QueryOutputModeIds), ", "))).
		AddBoolFlag(constants.ArgProgress, true, "Display snapshot upload status").
		AddStringSliceFlag(constants.ArgSearchPath, nil, "Set a custom search_path for the steampipe user for a query session (comma-separated)").
		AddStringSliceFlag(constants.ArgSearchPathPrefix, nil, "Set a prefix to the current search path for a query session (comma-separated)").
		AddStringFlag(constants.ArgSeparator, ",", "Separator string for csv output").
		AddBoolFlag(constants.ArgShare, false, "Create snapshot in Turbot Pipes with 'anyone_with_link' visibility").
		AddBoolFlag(constants.ArgSnapshot, false, "Create snapshot in Turbot Pipes with the default (workspace) visibility").
		AddStringFlag(constants.ArgSnapshotLocation, "", "The location to write snapshots - either a local file path or a Turbot Pipes workspace").
		AddStringArrayFlag(constants.ArgSnapshotTag, nil, "Specify tags to set on the snapshot").
		AddStringFlag(constants.ArgSnapshotTitle, "", "The title to give a snapshot").
		AddBoolFlag(constants.ArgTiming, false, "Turn on the query timer").
		// NOTE: use StringArrayFlag for ArgVariable, not StringSliceFlag
		// Cobra will interpret values passed to a StringSliceFlag as CSV, where args passed to StringArrayFlag are not parsed and used raw
		AddStringArrayFlag(constants.ArgVariable, nil, "Specify the value of a variable").
		AddStringSliceFlag(constants.ArgVarFile, nil, "Specify an .ppvar file containing variable values")

	return cmd
}

func queryRun(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()

	startTime := time.Now()

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

	initData := initialisation.NewInitData[*resources.Query](ctx, cmd, args...)
	// shutdown the service on exit
	defer initData.Cleanup(ctx)
	error_helpers.FailOnError(initData.Result.Error)

	// if there is a usage warning we display it
	initData.Result.DisplayMessages()

	if err := initData.Result.Error; err != nil {
		exitCode = constants.ExitCodeInitializationFailed
		error_helpers.FailOnError(err)
	}

	// register the query exporters if necessary
	if len(viper.GetStringSlice(constants.ArgExport)) > 0 {
		err := initData.RegisterExporters(queryExporters()...)
		error_helpers.FailOnError(err)

		// validate required export formats
		err = initData.ExportManager.ValidateExportFormat(viper.GetStringSlice(constants.ArgExport))
		error_helpers.FailOnError(err)
	}

	// execute query as a snapshot
	target, err := initData.GetSingleTarget()
	if err != nil {
		exitCode = constants.ExitCodeInitializationFailed
		error_helpers.FailOnError(err)
	}

	inputs := dashboardexecute.NewInputValues()
	snap, err := dashboardexecute.GenerateSnapshot(ctx, initData.Workspace, target, inputs)
	if err != nil {
		exitCode = constants.ExitCodeSnapshotCreationFailed
		error_helpers.FailOnError(err)
	}

	// display the result
	switch viper.GetString(constants.ArgOutput) {
	case constants.OutputFormatNone:
		// do nothing
	case constants.OutputFormatSnapshot, constants.OutputFormatPowerpipeSnapshotShort:
		// if the format is snapshot, just dump it out
		jsonOutput, err := json.MarshalIndent(snap, "", "  ")
		if err != nil {
			error_helpers.FailOnErrorWithMessage(err, "failed to display result as snapshot")
		}
		fmt.Println(string(jsonOutput)) //nolint:forbidigo // intentional use of fmt
	default:
		// otherwise convert the snapshot into a query result
		result, err := snapshotToQueryResult(snap, startTime)
		error_helpers.FailOnError(err)
		display.ShowQueryOutput(ctx, result)
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

	return localcmdconfig.ValidateDatabaseArg()
}

func queryExporters() []export.Exporter {
	return []export.Exporter{&export.SnapshotExporter{}}
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

func snapshotToQueryResult(snap *steampipeconfig.SteampipeSnapshot, startTime time.Time) (*queryresult.Result, error) {
	// the table of a snapshot query has a fixed name
	tablePanel, ok := snap.Panels[resources.SnapshotQueryTableName]
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
		res.Close()
	}()

	res.Timing = &queryresult.TimingMetadata{
		Duration: time.Since(startTime),
	}
	return res, nil
}
