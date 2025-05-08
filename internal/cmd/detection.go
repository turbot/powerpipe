package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thediveo/enumflag/v2"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/v2/cmdconfig"
	"github.com/turbot/pipe-fittings/v2/constants"
	"github.com/turbot/pipe-fittings/v2/error_helpers"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/pipe-fittings/v2/statushooks"
	"github.com/turbot/pipe-fittings/v2/workspace"
	localcmdconfig "github.com/turbot/powerpipe/internal/cmdconfig"
	localconstants "github.com/turbot/powerpipe/internal/constants"
	"github.com/turbot/powerpipe/internal/controldisplay"
	"github.com/turbot/powerpipe/internal/controlinit"
	"github.com/turbot/powerpipe/internal/dashboardexecute"
	"github.com/turbot/powerpipe/internal/resources"
	"github.com/turbot/steampipe-plugin-sdk/v5/logging"
)

type DetectionTarget interface {
	modconfig.ModTreeItem
	*resources.DetectionBenchmark | *resources.Detection
}

// variable used to assign the output mode flag
var detectionOutputMode = localconstants.DetectionOutputModeText

func detectionRunCmd[T DetectionTarget]() *cobra.Command {
	typeName := resources.GenericTypeToBlockType[T]()

	cmd := &cobra.Command{
		Use:              detectionCmdUse(typeName),
		TraverseChildren: true,
		Args:             cobra.ExactArgs(1),
		Run:              detectionRun[T],
		Short:            detectionCmdShort(typeName),
		Long:             detectionCmdLong(typeName),
	}

	// when running mod install before the detection execution, we use the minimal update strategy
	var updateStrategy = constants.ModUpdateIdMinimal

	cmdconfig.OnCmd(cmd).
		AddCloudFlags().
		AddModLocationFlag().
		AddStringArrayFlag(constants.ArgArg, nil, "Specify the value of a detection argument").
		AddBoolFlag(constants.ArgHeader, true, "Include column headers for csv and table output").
		AddStringFlag(constants.ArgSeparator, ",", "Separator string for csv output").
		AddStringSliceFlag(constants.ArgExport, nil, "Export output to file, supported format: json").
		AddStringFlag(constants.ArgDatabase, "", "Turbot Pipes workspace database", localcmdconfig.Deprecated("see https://powerpipe.io/docs/run#selecting-a-database for the new syntax")).
		AddIntFlag(constants.ArgDatabaseQueryTimeout, localconstants.DatabaseDefaultQueryTimeout, "The query timeout").
		AddBoolFlag(constants.ArgHelp, false, "Help for detection", cmdconfig.FlagOptions.WithShortHand("h")).
		AddBoolFlag(constants.ArgInput, true, "Enable interactive prompts").
		AddIntFlag(constants.ArgMaxParallel, constants.DefaultMaxConnections, "The maximum number of concurrent database connections to open").
		AddBoolFlag(constants.ArgModInstall, true, "Specify whether to install mod dependencies before running the detection").
		AddVarFlag(enumflag.New(&updateStrategy, constants.ArgPull, constants.ModUpdateStrategyIds, enumflag.EnumCaseInsensitive),
			constants.ArgPull,
			fmt.Sprintf("Update strategy; one of: %s", strings.Join(constants.FlagValues(constants.ModUpdateStrategyIds), ", "))).
		AddVarFlag(enumflag.New(&detectionOutputMode, constants.ArgOutput, localconstants.DetectionOutputModeIds, enumflag.EnumCaseInsensitive),
			constants.ArgOutput,
			fmt.Sprintf("Output format; one of: %s", strings.Join(constants.FlagValues(localconstants.DetectionOutputModeIds), ", "))).
		AddBoolFlag(constants.ArgProgress, true, "Display detection execution progress respected when a detection name argument is passed").
		AddBoolFlag(constants.ArgSnapshot, false, "Create snapshot in Turbot Pipes with the default (workspace) visibility").
		AddBoolFlag(constants.ArgShare, false, "Create snapshot in Turbot Pipes with 'anyone_with_link' visibility").
		AddStringFlag(constants.ArgSnapshotTitle, "", "The title to give a snapshot").
		// NOTE: use StringArrayFlag for ArgDetectionInput, not StringSliceFlag
		// Cobra will interpret values passed to a StringSliceFlag as CSV, where args passed to StringArrayFlag are not parsed and used raw
		AddStringArrayFlag(constants.ArgSnapshotTag, nil, "Specify tags to set on the snapshot").
		AddStringFlag(constants.ArgSnapshotLocation, "", "The location to write snapshots - either a local file path or a Turbot Pipes workspace").
		// NOTE: use StringArrayFlag for ArgVariable, not StringSliceFlag
		// Cobra will interpret values passed to a StringSliceFlag as CSV, where args passed to StringArrayFlag are not parsed and used raw
		AddStringArrayFlag(constants.ArgVariable, nil, "Specify the value of a variable").
		AddStringSliceFlag(constants.ArgVarFile, nil, "Specify an .ppvar file containing variable values").
		AddIntFlag(constants.ArgDetectionTimeout, 0, "Set the detection execution timeout")

	return cmd
}

func detectionRun[T DetectionTarget](cmd *cobra.Command, args []string) {
	detectionRunWithInitData[T](cmd, nil, args)
}

// tactical - to support callint benchmark run and calling either control or dashboard execution flow,
// we must support calling this from the check command, AFTER the initdata has been fetched
func detectionRunWithInitData[T DetectionTarget](cmd *cobra.Command, initData *controlinit.InitData, args []string) {
	ctx := cmd.Context()

	// there can only be a single arg - cobra will validate
	detectionName := args[0]

	var err error
	logging.LogTime("detectionRun start")
	defer func() {
		logging.LogTime("detectionRun end")
		if r := recover(); r != nil {
			err = helpers.ToError(r)
			error_helpers.ShowError(ctx, err)

		}
		setExitCodeForDetectionError(err)
	}()

	// first check whether a single detection name has been passed as an arg
	error_helpers.FailOnError(validateDetectionArgs(ctx))

	// if diagnostic mode is set, print out config and return
	if _, ok := os.LookupEnv(localconstants.EnvConfigDump); ok {
		localcmdconfig.DisplayConfig()
		return
	}
	// create context for the detection execution
	ctx = createSnapshotContext(ctx, detectionName)

	statushooks.SetStatus(ctx, "Initializing…")
	if initData == nil {
		initData = controlinit.NewInitData[T](ctx, cmd, detectionName)
	}

	if initData.Result.Error != nil {
		exitCode = constants.ExitCodeInitializationFailed
		error_helpers.ShowError(ctx, initData.Result.Error)
		return
	}
	defer initData.Cleanup(ctx)

	if len(viper.GetStringSlice(constants.ArgExport)) > 0 {
		// validate required export formats
		err = initData.ExportManager.ValidateExportFormat(viper.GetStringSlice(constants.ArgExport))
		error_helpers.FailOnError(err)
	}

	statushooks.Done(ctx)

	// if there is a usage warning we display it
	initData.Result.DisplayMessages()

	// so a detection name was specified - just call GenerateSnapshot
	target, err := initData.GetSingleTarget()
	error_helpers.FailOnError(err)

	inputs := dashboardexecute.NewInputValues()

	snap, err := dashboardexecute.GenerateSnapshot(ctx, initData.Workspace, target, inputs)
	error_helpers.FailOnError(err)

	tree, err := controldisplay.SnapshotToExecutionTree(ctx, snap, initData.Workspace, target)
	error_helpers.FailOnError(err)

	err = displayDetectionResults(ctx, tree, initData.OutputFormatter)
	error_helpers.FailOnError(err)

	// display the snapshot result (if needed)
	//displaySnapshot(snap)

	// upload the snapshot (if needed)
	err = publishSnapshotIfNeeded(ctx, snap)
	if err != nil {
		exitCode = constants.ExitCodeSnapshotUploadFailed
		error_helpers.FailOnErrorWithMessage(err, fmt.Sprintf("failed to publish snapshot to %s", viper.GetString(constants.ArgSnapshotLocation)))
	}

	// export the result (if needed)
	exportArgs := viper.GetStringSlice(constants.ArgExport)
	exportMsg, err := initData.ExportManager.DoExport(ctx, snap.FileNameRoot, tree, exportArgs)
	error_helpers.FailOnErrorWithMessage(err, "failed to export snapshot")

	// print the location where the file is exported
	if len(exportMsg) > 0 && viper.GetBool(constants.ArgProgress) {
		//nolint:forbidigo // Intentional UI output
		fmt.Printf("\n%s\n", strings.Join(exportMsg, "\n"))
	}
}

// validate the args and extract a detection name, if provided
func validateDetectionArgs(ctx context.Context) error {
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

func setExitCodeForDetectionError(err error) {
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

func detectionCmdUse(typeName string) string {
	return fmt.Sprintf("run [flags] [%s]", typeName)
}

func detectionCmdShort(typeName string) string {
	return fmt.Sprintf("Execute one or more %ss", typeName)
}
func detectionCmdLong(typeName string) string {
	return fmt.Sprintf(`Execute one or more %ss.

You may specify one or more %ss to run, separated by a space.`, typeName, typeName)
}
