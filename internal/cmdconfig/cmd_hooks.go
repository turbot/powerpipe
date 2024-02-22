package cmdconfig

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/app_specific"
	"github.com/turbot/pipe-fittings/cloud"
	"github.com/turbot/pipe-fittings/cmdconfig"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/filepaths"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/steampipeconfig"
	"github.com/turbot/pipe-fittings/task"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/powerpipe/internal/logger"
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

	logger.Initialize()

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
	loader, err := cmdconfig.GetWorkspaceProfileLoader[*modconfig.SteampipeWorkspaceProfile]()
	if err != nil {
		return error_helpers.NewErrorsAndWarning(err)
	}

	var cmd = viper.Get(constants.ConfigKeyActiveCommand).(*cobra.Command)

	// set-up viper with defaults from the env and default workspace profile

	cmdconfig.BootstrapViper(loader, cmd,
		cmdconfig.WithConfigDefaults(configDefaults()),
		cmdconfig.WithDirectoryEnvMappings(dirEnvMappings()))

	if err != nil {
		return error_helpers.NewErrorsAndWarning(err)
	}

	// set global containing the configured install dir (create directory if needed)
	ensureInstallDirs()

	// set the rest of the defaults from ENV
	// ENV takes precedence over any default configuration
	cmdconfig.SetDefaultsFromEnv(envMappings())

	// if an explicit workspace profile was set, add to viper as highest precedence default
	// NOTE: if install_dir/mod_location are set these will already have been passed to viper by BootstrapViper
	// since the "ConfiguredProfile" is passed in through a cmdline flag, it will always take precedence
	if loader.ConfiguredProfile != nil {
		cmdconfig.SetDefaultsFromConfig(loader.ConfiguredProfile.ConfigMap(cmd))
	}

	// NOTE: we need to resolve the token separately
	// - that is because we need the resolved value of ArgPipesHost in order to load any saved token
	// and we cannot get this until the other config has been resolved
	err = setPipesTokenDefault(loader)
	if err != nil {
		return error_helpers.NewErrorsAndWarning(err)
	}

	// now validate all config values have appropriate values
	return validateConfig()
}

func setPipesTokenDefault(loader *steampipeconfig.WorkspaceProfileLoader[*modconfig.SteampipeWorkspaceProfile]) error {
	/*
	   saved cloud token
	   pipes_token in default workspace
	   explicit env var (PIPES_TOKEN ) wins over
	   pipes_token in specific workspace
	*/
	// set viper defaults in order of increasing precedence
	// 1) saved cloud token
	savedToken, err := cloud.LoadToken()
	if err != nil {
		return err
	}
	if savedToken != "" {
		viper.SetDefault(constants.ArgPipesToken, savedToken)
	}
	// 2) default profile cloud token
	if loader.DefaultProfile.PipesToken != nil {
		viper.SetDefault(constants.ArgPipesToken, *loader.DefaultProfile.PipesToken)
	}
	// 3) env var (PIPES_TOKEN )
	cmdconfig.SetDefaultFromEnv(constants.EnvPipesToken, constants.ArgPipesToken, cmdconfig.EnvVarTypeString)

	// 4) explicit workspace profile
	if p := loader.ConfiguredProfile; p != nil && p.PipesToken != nil {
		viper.SetDefault(constants.ArgPipesToken, *p.PipesToken)
	}
	return nil
}

// now validate  config values have appropriate values
// (currently validates telemetry)
func validateConfig() error_helpers.ErrorAndWarnings {
	var res = error_helpers.ErrorAndWarnings{}
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

func ensureInstallDirs() {
	installDir := viper.GetString(constants.ArgInstallDir)
	pipesInstallDir := viper.GetString(constants.ArgPipesInstallDir)

	slog.Debug("ensureInstallDir", "installDir", installDir, "pipesInstallDir", pipesInstallDir)
	if _, err := os.Stat(installDir); os.IsNotExist(err) {
		slog.Debug("creating install dir")
		err = os.MkdirAll(installDir, 0755)
		error_helpers.FailOnErrorWithMessage(err, fmt.Sprintf("could not create installation directory: %s", installDir))
	}
	if _, err := os.Stat(pipesInstallDir); os.IsNotExist(err) {
		slog.Debug("creating pipes install dir")
		err = os.MkdirAll(pipesInstallDir, 0755)
		error_helpers.FailOnErrorWithMessage(err, fmt.Sprintf("could not create pipes installation directory: %s", installDir))
	}

	// store as InstallDir
	app_specific.InstallDir = installDir
	filepaths.PipesInstallDir = installDir
}
