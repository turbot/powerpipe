package cmdconfig

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/v2/app_specific"
	"github.com/turbot/pipe-fittings/v2/cmdconfig"
	"github.com/turbot/pipe-fittings/v2/connection"
	"github.com/turbot/pipe-fittings/v2/constants"
	"github.com/turbot/pipe-fittings/v2/error_helpers"
	"github.com/turbot/pipe-fittings/v2/filepaths"
	"github.com/turbot/pipe-fittings/v2/parse"
	"github.com/turbot/pipe-fittings/v2/pipes"
	"github.com/turbot/pipe-fittings/v2/task"
	"github.com/turbot/pipe-fittings/v2/utils"
	"github.com/turbot/pipe-fittings/v2/workspace_profile"
	"github.com/turbot/powerpipe/internal/logger"
	"github.com/turbot/powerpipe/internal/powerpipeconfig"
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
	loader, err := cmdconfig.GetWorkspaceProfileLoader[*workspace_profile.PowerpipeWorkspaceProfile]()
	if err != nil {
		return error_helpers.NewErrorsAndWarning(err)
	}

	var cmd = viper.Get(constants.ConfigKeyActiveCommand).(*cobra.Command)

	var config, ew = powerpipeconfig.LoadPowerpipeConfig(filepaths.EnsureConfigDir())
	if ew.GetError() != nil {
		return ew
	}

	powerpipeconfig.GlobalConfig = config

	// set-up viper with defaults from the env and default workspace profile

	cmdconfig.BootstrapViper(loader, cmd,
		cmdconfig.WithConfigDefaults(configDefaults(cmd)),
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
	wp := loader.ConfiguredProfile
	if wp != nil {
		cmdconfig.SetDefaultsFromConfig(wp.ConfigMap(cmd))
	}

	// now env vars have been processed, set filepaths.PipesInstallDir
	filepaths.PipesInstallDir = viper.GetString(constants.ArgPipesInstallDir)

	// NOTE: we need to resolve the token separately
	// - that is because we need the resolved value of ArgPipesHost in order to load any saved token
	// and we cannot get this until the other config has been resolved
	err = setPipesTokenDefault(loader)
	if err != nil {
		return error_helpers.NewErrorsAndWarning(err)
	}

	// if the configured workspace is a cloud workspace, create cloud metadata and set the default connection
	if wp != nil && wp.IsCloudWorkspace() {
		pipesMetadata, ew := wp.GetPipesMetadata()
		if ew.GetError() != nil {
			return ew
		}
		// create new default connection
		defaultConnection := connection.NewSteampipePgConnection("default", hcl.Range{}).(*connection.SteampipePgConnection)
		defaultConnection.ConnectionString = &pipesMetadata.ConnectionString
		// TODO temp for now we must call validate to populate the defaults
		_ = defaultConnection.Validate()
		config.SetDefaultConnection(defaultConnection)
	}

	// now validate all config values have appropriate values
	return validateConfig(loader.GetActiveWorkspaceProfile())
}

func setPipesTokenDefault(loader *parse.WorkspaceProfileLoader[*workspace_profile.PowerpipeWorkspaceProfile]) error {
	/*
	   saved cloud token
	   pipes_token in default workspace
	   explicit env var (PIPES_TOKEN ) wins over
	   pipes_token in specific workspace
	*/
	// set viper defaults in order of increasing precedence
	// 1) saved cloud token
	savedToken, err := pipes.LoadToken()
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

// now validate config values have appropriate values
func validateConfig(activeWorkspace *workspace_profile.PowerpipeWorkspaceProfile) error_helpers.ErrorAndWarnings {
	var res = error_helpers.ErrorAndWarnings{}
	telemetry := viper.GetString(constants.ArgTelemetry)
	if !helpers.StringSliceContains(constants.TelemetryLevels, telemetry) {
		res.Error = sperr.New(`invalid value of 'telemetry' (%s), must be one of: %s`, telemetry, strings.Join(constants.TelemetryLevels, ", "))
		return res
	}
	if _, legacyDiagnosticsSet := os.LookupEnv(plugin.EnvLegacyDiagnosticsLevel); legacyDiagnosticsSet {
		res.AddWarning(fmt.Sprintf("Environment variable %s is deprecated - use %s", plugin.EnvLegacyDiagnosticsLevel, plugin.EnvDiagnosticsLevel))
	}

	// database deprecation warnings
	if _, dbEnvSet := os.LookupEnv(app_specific.EnvDatabase); dbEnvSet {
		res.AddWarning(fmt.Sprintf("Environment variable %s is deprecated, see https://powerpipe.io/docs/run#selecting-a-database for the new syntax", app_specific.EnvDatabase))
	}
	// check active workspace profile
	if activeWorkspace != nil && activeWorkspace.Database != nil {
		res.AddWarning(fmt.Sprintf("workspace property 'database' is deprecated, see https://powerpipe.io/docs/run#selecting-a-database for the new syntax. (%s:%d-%d)", activeWorkspace.DeclRange.Filename, activeWorkspace.DeclRange.Start.Line, activeWorkspace.DeclRange.End.Line))
	}

	res.Error = plugin.ValidateDiagnosticsEnvVar()

	return res
}

// create ~/.steampipe if needed
func ensureInstallDirs() {
	installDir := viper.GetString(constants.ArgInstallDir)

	slog.Debug("ensureInstallDir", "installDir", installDir)
	if _, err := os.Stat(installDir); os.IsNotExist(err) {
		slog.Debug("creating install dir")
		err = os.MkdirAll(installDir, 0755)
		error_helpers.FailOnErrorWithMessage(err, fmt.Sprintf("could not create installation directory: %s", installDir))
	}

	// store as app_specific.InstallDir
	app_specific.InstallDir = installDir
}
