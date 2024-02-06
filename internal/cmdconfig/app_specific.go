package cmdconfig

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/viper"
	"github.com/turbot/go-kit/files"
	"github.com/turbot/pipe-fittings/app_specific"
	"github.com/turbot/pipe-fittings/cmdconfig"
	"github.com/turbot/pipe-fittings/error_helpers"
)

// SetAppSpecificConstants sets app specific constants defined in pipe-fittings
func SetAppSpecificConstants() {
	app_specific.AppName = "powerpipe"

	// set all app specific env var keys
	app_specific.SetAppSpecificEnvVarKeys("POWERPIPE_")

	// version
	versionString := viper.GetString("main.version")
	app_specific.AppVersion = semver.MustParse(versionString)
	app_specific.AutoVariablesExtension = ".auto.ppvars"

	// set the default install dir
	defaultInstallDir, err := files.Tildefy("~/.powerpipe")
	error_helpers.FailOnError(err)
	app_specific.DefaultInstallDir = defaultInstallDir

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
	// default to local steampipe service
	app_specific.DefaultDatabase = "postgres://steampipe@127.0.0.1:9193/steampipe"

	// extensions

	app_specific.ModDataExtensions = []string{".sp", ".pp"}
	app_specific.ModFileName = "mod.sp"
	app_specific.ConfigExtension = ".ppc"
	app_specific.VariablesExtension = ".ppvars"

	app_specific.WorkspaceIgnoreFile = ".powerpipeignore"
	app_specific.WorkspaceDataDir = ".powerpipe"
	// EnvInputVarPrefix is the prefix for environment variables that represent values for input variables.
	app_specific.EnvInputVarPrefix = "PP_VAR_"

	// set the command pre and post hooks
	cmdconfig.CustomPreRunHook = preRunHook
	cmdconfig.CustomPostRunHook = postRunHook

}
