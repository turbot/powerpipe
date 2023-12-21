package constants

import (
	"github.com/turbot/go-kit/files"
	"github.com/turbot/pipe-fittings/app_specific"
	internalversion "github.com/turbot/powerpipe/internal/version"
)

// SetAppSpecificConstants sets app specific constants defined in pipe-fittings
func SetAppSpecificConstants() {
	// set the default install dir
	app_specific.AppVersion = internalversion.PowerpipeVersion
	app_specific.AutoVariablesExtension = ".auto.ppvars"
	app_specific.ClientConnectionAppNamePrefix = "powerpipe_client"
	app_specific.ClientSystemConnectionAppNamePrefix = "powerpipe_client_system"
	// set the default install dir
	installDir, err := files.Tildefy("~/.powerpipe")
	if err != nil {
		panic(err)
	}
	app_specific.DefaultInstallDir = installDir

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
	app_specific.SetAppSpecificEnvVarKeys("POWERPIPE_")
	// EnvInputVarPrefix is the prefix for environment variables that represent values for input variables.
	app_specific.EnvInputVarPrefix = "PP_VAR_"

}
