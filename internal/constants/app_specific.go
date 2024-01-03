package constants

import (
	"github.com/turbot/go-kit/files"
	"github.com/turbot/pipe-fittings/app_specific"
	internalversion "github.com/turbot/powerpipe/internal/version"
	"os"
	"path/filepath"
	"strings"
)

// SetAppSpecificConstants sets app specific constants defined in pipe-fittings
func SetAppSpecificConstants() {
	// set the default install dir
	app_specific.AppVersion = internalversion.PowerpipeVersion
	app_specific.AutoVariablesExtension = ".auto.ppvars"
	app_specific.ClientConnectionAppNamePrefix = "powerpipe_client"
	app_specific.ClientSystemConnectionAppNamePrefix = "powerpipe_client_system"
	// set the default install dir
	defaultInstallDir, err := files.Tildefy("~/.powerpipe")
	if err != nil {
		panic(err)
	}
	app_specific.DefaultInstallDir = defaultInstallDir
	// TODO KAI CHECK THIS
	// set the default config path
	globalConfigPath := filepath.Join(defaultInstallDir, "config")
	// check whether install-dir env has been set - if so, respect it
	if envInstallDir, ok := os.LookupEnv(app_specific.EnvInstallDir); ok {
		globalConfigPath = filepath.Join(envInstallDir, "config")
		app_specific.InstallDir = envInstallDir
	} else {
		/*
			NOTE:
			If InstallDir is settable outside of default & env var, need to add
			the following code to end of initGlobalConfig in init.go
			app_specific.InstallDir = viper.GetString(constants.ArgInstallDir) at end of
		*/
		app_specific.InstallDir = defaultInstallDir
	}
	app_specific.DefaultConfigPath = strings.Join([]string{".", globalConfigPath}, ":")

	app_specific.DefaultVarsFileName = "powerpipe.ppvars"
	// default to local steampipe service
	app_specific.DefaultWorkspaceDatabase = "postgres://steampipe@127.0.0.1:9193/steampipe"
	app_specific.ModDataExtension = ".sp"
	app_specific.ModFileName = "mod.sp"
	app_specific.ServiceConnectionAppNamePrefix = "powerpipe_service"
	app_specific.ConfigExtension = ".ppc"
	app_specific.VariablesExtension = ".ppvars"
	app_specific.WorkspaceIgnoreFile = ".powerpipeignore"
	app_specific.WorkspaceDataDir = ".powerpipe"
	app_specific.EnvAppPrefix = "POWERPIPE_"
	// EnvInputVarPrefix is the prefix for environment variables that represent values for input variables.
	app_specific.EnvInputVarPrefix = "PP_VAR_"

}
