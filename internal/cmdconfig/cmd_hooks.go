package cmdconfig

import (
	"context"
	"fmt"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/task"
	"github.com/turbot/powerpipe/internal/logger"
	"log/slog"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/app_specific"
	"github.com/turbot/pipe-fittings/cmdconfig"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

var waitForTasksChannel chan struct{}
var tasksCancelFn context.CancelFunc

// postRunHook is a function that is executed after the PostRun of every command handler
func postRunHook(_ *cobra.Command, _ []string) error {
	utils.LogTime("cmdhook.postRunHook start")
	defer utils.LogTime("cmdhook.postRunHook end")

	if waitForTasksChannel != nil {
		// wait for the async tasks to finish
		select {
		case <-time.After(100 * time.Millisecond):
			tasksCancelFn()
			return nil
		case <-waitForTasksChannel:
			return nil
		}
	}
	return nil
}

// postRunHook is a function that is executed before the PreRun of every command handler
func preRunHook(cmd *cobra.Command, args []string) error {
	utils.LogTime("cmdhook.preRunHook start")
	defer utils.LogTime("cmdhook.preRunHook end")

	viper.Set(constants.ConfigKeyActiveCommand, cmd)
	viper.Set(constants.ConfigKeyActiveCommandArgs, args)
	viper.Set(constants.ConfigKeyIsTerminalTTY, isatty.IsTerminal(os.Stdout.Fd()))

	// set up the global viper config with default values from
	// config files and ENV variables
	ew := initGlobalConfig()
	// display any warnings
	ew.ShowWarnings()
	// check for error
	error_helpers.FailOnError(ew.Error)

	logger.SetDefaultLogger()

	// runScheduledTasks skips running tasks if this instance is the plugin manager
	waitForTasksChannel = runScheduledTasks(cmd.Context(), cmd, args)

	// set the max memory if specified
	setMemoryLimit()
	return nil
}

func setMemoryLimit() {
	maxMemoryBytes := viper.GetInt64(constants.ArgMemoryMaxMb) * 1024 * 1024
	if maxMemoryBytes > 0 {
		// set the max memory
		debug.SetMemoryLimit(maxMemoryBytes)
	}
}

// runScheduledTasks runs the task runner and returns a channel which is closed when
// task run is complete
//
// runScheduledTasks skips running tasks if this instance is the plugin manager
func runScheduledTasks(ctx context.Context, cmd *cobra.Command, args []string) chan struct{} {
	updateCheck := viper.GetBool(constants.ArgUpdateCheck)
	// for now the only scheduled task we support is update check so if that is disabled, do nothing
	if !updateCheck {
		return nil
	}

	taskUpdateCtx, cancelFn := context.WithCancel(ctx)
	tasksCancelFn = cancelFn

	return task.RunTasks(
		taskUpdateCtx,
		cmd,
		args,
		// pass the config value in rather than runRasks querying viper directly - to avoid concurrent map access issues
		// (we can use the update-check viper config here, since initGlobalConfig has already set it up
		// with values from the config files and ENV settings - update-check cannot be set from the command line)
		task.WithUpdateCheck(updateCheck),
	)
}

// initConfig reads in config file and ENV variables if set.
func initGlobalConfig() error_helpers.ErrorAndWarnings {
	utils.LogTime("cmdconfig.initGlobalConfig start")
	defer utils.LogTime("cmdconfig.initGlobalConfig end")

	// load workspace profile from the configured install dir
	loader, err := cmdconfig.GetWorkspaceProfileLoader[*modconfig.PowerpipeWorkspaceProfile]()
	error_helpers.FailOnError(err)

	var cmd = viper.Get(constants.ConfigKeyActiveCommand).(*cobra.Command)

	// set-up viper with defaults from the env and default workspace profile

	cmdconfig.BootstrapViper(loader, cmd,
		cmdconfig.WithConfigDefaults(configDefaults()),
		cmdconfig.WithDirectoryEnvMappings(dirEnvMappings()))

	if err != nil {
		return error_helpers.NewErrorsAndWarning(err)
	}

	// set global containing the configured install dir (create directory if needed)
	ensureInstallDir(viper.GetString(constants.ArgInstallDir))

	// set the rest of the defaults from ENV
	// ENV takes precedence over any default configuration
	cmdconfig.SetDefaultsFromEnv(envMappings())

	// NOTE: we need to resolve the token separately
	// - that is because we need the resolved value of ArgCloudHost in order to load any saved token
	// and we cannot get this until the other config has been resolved
	// err = setCloudTokenDefault()
	// if err != nil {
	// 	loadConfigErrorsAndWarnings.Error = err
	// 	return loadConfigErrorsAndWarnings
	// }

	// now validate all config values have appropriate values
	ew := validateConfig()
	error_helpers.FailOnErrorWithMessage(ew.Error, "failed to validate config")

	// loadConfigErrorsAndWarnings.Merge(ew)

	return error_helpers.NewErrorsAndWarning(nil)
}

// now validate  config values have appropriate values
// (currently validates telemetry)
func validateConfig() *error_helpers.ErrorAndWarnings {
	var res = &error_helpers.ErrorAndWarnings{}
	telemetry := viper.GetString(constants.ArgTelemetry)
	if !helpers.StringSliceContains(constants.TelemetryLevels, telemetry) {
		res.Error = sperr.New(`invalid value of 'telemetry' (%s), must be one of: %s`, telemetry, strings.Join(constants.TelemetryLevels, ", "))
		return res
	}
	if _, legacyDiagnosticsSet := os.LookupEnv(plugin.EnvLegacyDiagnosticsLevel); legacyDiagnosticsSet {
		res.AddWarning(fmt.Sprintf("Environment variable %s is deprecated - use %s", plugin.EnvLegacyDiagnosticsLevel, plugin.EnvDiagnosticsLevel))
	}
	res.Error = plugin.ValidateDiagnosticsEnvVar()

	return res
}

// create a hclog logger with the level specified by the SP_LOG env var
// func createLogger(logBuffer *bytes.Buffer, cmd *cobra.Command) {
// 	if task.IsPluginManagerCmd(cmd) {
// 		// nothing to do here - plugin manager sets up it's own logger
// 		// refer https://github.com/turbot/steampipe/blob/710a96d45fd77294de8d63d77bf78db65133e5ca/cmd/plugin_manager.go#L102
// 		return
// 	}

// 	level := sdklogging.LogLevel()
// 	var logDestination io.Writer
// 	if len(filepaths.SteampipeDir) == 0 {
// 		// write to the buffer - this is to make sure that we don't lose logs
// 		// till the time we get the log directory
// 		logDestination = logBuffer
// 	} else {
// 		logDestination = logging.NewRotatingLogWriter(filepaths.EnsureLogDir(), "steampipe")

// 		// write out the buffered contents
// 		_, _ = logDestination.Write(logBuffer.Bytes())
// 	}

// 	hcLevel := hclog.LevelFromString(level)

// 	options := &hclog.LoggerOptions{
// 		// make the name unique so that logs from this instance can be filtered
// 		Name:       fmt.Sprintf("steampipe [%s]", runtime.ExecutionID),
// 		Level:      hcLevel,
// 		Output:     logDestination,
// 		TimeFn:     func() time.Time { return time.Now().UTC() },
// 		TimeFormat: "2006-01-02 15:04:05.000 UTC",
// 	}
// 	logger := sdklogging.NewLogger(options)
// 	log.SetOutput(logger.StandardWriter(&hclog.StandardLoggerOptions{InferLevels: true}))
// 	log.SetPrefix("")
// 	log.SetFlags(0)

// 	// if the buffer is empty then this is the first time the logger is getting setup
// 	// write out a banner
// 	if logBuffer.Len() == 0 {
// 		// pump in the initial set of logs
// 		// this will also write out the Execution ID - enabling easy filtering of logs for a single execution
// 		// we need to do this since all instances will log to a single file and logs will be interleaved
// 		slog.Info("********************************************************\n")
// 		slog.Info("**%16s%20s%16s**\n", " ", fmt.Sprintf("Steampipe [%s]", runtime.ExecutionID), " ")
// 		slog.Info("********************************************************\n")
// 		slog.Info("AppVersion:   v%s\n", version.VersionString)
// 		slog.Info("Log level: %s\n", sdklogging.LogLevel())
// 		slog.Info("Log date: %s\n", time.Now().Format("2006-01-02"))
// 		//
// 	}
// }

func ensureInstallDir(installDir string) {
	slog.Debug("ensureInstallDir", "installDir", installDir)
	if _, err := os.Stat(installDir); os.IsNotExist(err) {
		slog.Debug("creating install dir")
		err = os.MkdirAll(installDir, 0755)
		error_helpers.FailOnErrorWithMessage(err, fmt.Sprintf("could not create installation directory: %s", installDir))
	}

	// store as PowerpipeDir
	app_specific.InstallDir = installDir
}

// displayDeprecationWarnings shows the deprecated warnings in a formatted way
func displayDeprecationWarnings(errorsAndWarnings *error_helpers.ErrorAndWarnings) {
	if len(errorsAndWarnings.Warnings) > 0 {
		fmt.Println(color.YellowString(fmt.Sprintf("\nDeprecation %s:", utils.Pluralize("warning", len(errorsAndWarnings.Warnings)))))
		for _, warning := range errorsAndWarnings.Warnings {
			fmt.Printf("%s\n\n", warning)
		}
		fmt.Println("For more details, see https://steampipe.io/docs/reference/config-files/workspace")
		fmt.Println()
	}
}
