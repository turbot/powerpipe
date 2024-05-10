package cmdconfig

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/turbot/go-kit/files"
	"github.com/turbot/pipe-fittings/app_specific"
	"github.com/turbot/pipe-fittings/cmdconfig"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/filepaths"
)

// SetAppSpecificConstants sets app specific constants defined in pipe-fittings
func SetAppSpecificConstants() {
	app_specific.AppName = "powerpipe"

	// set all app specific env var keys
	app_specific.SetAppSpecificEnvVarKeys("POWERPIPE_")

	// version
	versionString := viper.GetString("main.version")
	app_specific.AppVersion = semver.MustParse(versionString)

	// set the default install dir
	defaultInstallDir, err := files.Tildefy("~/.powerpipe")
	error_helpers.FailOnError(err)
	app_specific.DefaultInstallDir = defaultInstallDir
	defaultPipesInstallDir, err := files.Tildefy("~/.pipes")
	filepaths.DefaultPipesInstallDir = defaultPipesInstallDir
	error_helpers.FailOnError(err)

	// set the default config path
	globalConfigPath := filepath.Join(defaultInstallDir, "config")
	// check whether install-dir env has been set - if so, respect it
	if envInstallDir, ok := os.LookupEnv(app_specific.EnvInstallDir); ok {
		globalConfigPath = filepath.Join(envInstallDir, "config")
		app_specific.InstallDir = envInstallDir
	} else {
		app_specific.InstallDir = defaultInstallDir
	}
	app_specific.DefaultConfigPath = strings.Join([]string{".", globalConfigPath}, ":")

	app_specific.DefaultVarsFileName = "powerpipe.ppvars"
	app_specific.LegacyDefaultVarsFileName = "steampipe.spvars"
	// default to local steampipe service
	app_specific.DefaultDatabase = "postgres://steampipe@127.0.0.1:9193/steampipe"

	// extensions

	// NOTE: where we support multiple extensions, ensure the default is the FIRST in the list
	app_specific.ModDataExtensions = []string{".pp", ".sp"}
	app_specific.VariablesExtensions = []string{".ppvars", ".spvars"}
	app_specific.AutoVariablesExtensions = []string{".auto.ppvars", ".auto.spvars"}

	app_specific.ConfigExtension = ".ppc"
	app_specific.WorkspaceIgnoreFile = ".powerpipeignore"
	app_specific.WorkspaceDataDir = ".powerpipe"
	// EnvInputVarPrefix is the prefix for environment variables that represent values for input variables.
	app_specific.EnvInputVarPrefix = "PP_VAR_"

	// set the command pre and post hooks
	cmdconfig.CustomPreRunHook = preRunHook
	cmdconfig.CustomPostRunHook = postRunHook

	// Version check
	app_specific.VersionCheckHost = "hub.powerpipe.io"
	app_specific.VersionCheckPath = "api/cli/version/latest"
	app_specific.EnvProfile = "POWERPIPE_PROFILE"
}
