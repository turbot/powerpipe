package cmdconfig

import (
	"github.com/spf13/cobra"
	"github.com/turbot/pipe-fittings/v2/app_specific"
	"github.com/turbot/pipe-fittings/v2/cmdconfig"
	"github.com/turbot/pipe-fittings/v2/constants"
	"github.com/turbot/pipe-fittings/v2/filepaths"
	localconstants "github.com/turbot/powerpipe/internal/constants"
	"golang.org/x/exp/maps"
)

func configDefaults(cmd *cobra.Command) map[string]any {
	defs := map[string]any{
		// global general options
		constants.ArgTelemetry:       constants.TelemetryInfo,
		constants.ArgUpdateCheck:     true,
		constants.ArgPipesInstallDir: filepaths.DefaultPipesInstallDir,

		// dashboard
		constants.ArgDashboardStartTimeout: constants.DashboardServiceStartTimeout.Seconds(),

		// memory
		constants.ArgMemoryMaxMb: 1024,
	}
	cmdSpecificDefaults, ok := cmdSpecificDefaults()[cmd.Name()]
	if ok {
		maps.Copy(defs, cmdSpecificDefaults)
	}
	return defs
}

// command specific config defaults (keyed by comand name)
func cmdSpecificDefaults() map[string]map[string]any {
	return map[string]map[string]any{
		"server": {
			constants.ArgEnvironment: "release",
		},
	}
}

// environment variable mappings for directory paths which must be set as part of the viper bootstrap process
func dirEnvMappings() map[string]cmdconfig.EnvMapping {
	return map[string]cmdconfig.EnvMapping{
		app_specific.EnvInstallDir:  {ConfigVar: []string{constants.ArgInstallDir}, VarType: cmdconfig.EnvVarTypeString},
		app_specific.EnvModLocation: {ConfigVar: []string{constants.ArgModLocation}, VarType: cmdconfig.EnvVarTypeString},
	}
}

// NOTE: EnvWorkspaceProfile has already been set as a viper default as we have already loaded workspace profiles
// (EnvInstallDir has already been set at same time but we set it again to make sure it has the correct precedence)

// a map of known environment variables to map to viper keys - these are set as part of LoadGlobalConfig
func envMappings() map[string]cmdconfig.EnvMapping {
	return map[string]cmdconfig.EnvMapping{
		app_specific.EnvInstallDir:        {ConfigVar: []string{constants.ArgInstallDir}, VarType: cmdconfig.EnvVarTypeString},
		app_specific.EnvModLocation:       {ConfigVar: []string{constants.ArgModLocation}, VarType: cmdconfig.EnvVarTypeString},
		app_specific.EnvTelemetry:         {ConfigVar: []string{constants.ArgTelemetry}, VarType: cmdconfig.EnvVarTypeString},
		app_specific.EnvUpdateCheck:       {ConfigVar: []string{constants.ArgUpdateCheck}, VarType: cmdconfig.EnvVarTypeBool},
		app_specific.EnvSnapshotLocation:  {ConfigVar: []string{constants.ArgSnapshotLocation}, VarType: cmdconfig.EnvVarTypeString},
		app_specific.EnvDatabase:          {ConfigVar: []string{constants.ArgDatabase}, VarType: cmdconfig.EnvVarTypeString},
		app_specific.EnvDisplayWidth:      {ConfigVar: []string{constants.ArgDisplayWidth}, VarType: cmdconfig.EnvVarTypeInt},
		app_specific.EnvMaxParallel:       {ConfigVar: []string{constants.ArgMaxParallel}, VarType: cmdconfig.EnvVarTypeInt},
		app_specific.EnvQueryTimeout:      {ConfigVar: []string{constants.ArgDatabaseQueryTimeout}, VarType: cmdconfig.EnvVarTypeInt},
		app_specific.EnvCacheTTL:          {ConfigVar: []string{constants.ArgCacheTtl}, VarType: cmdconfig.EnvVarTypeInt},
		app_specific.EnvCacheMaxTTL:       {ConfigVar: []string{constants.ArgCacheMaxTtl}, VarType: cmdconfig.EnvVarTypeInt},
		app_specific.EnvMemoryMaxMb:       {ConfigVar: []string{constants.ArgMemoryMaxMb}, VarType: cmdconfig.EnvVarTypeInt},
		app_specific.EnvMemoryMaxMbPlugin: {ConfigVar: []string{constants.ArgMemoryMaxMbPlugin}, VarType: cmdconfig.EnvVarTypeInt},
		app_specific.EnvCacheEnabled:      {ConfigVar: []string{constants.ArgClientCacheEnabled, constants.ArgServiceCacheEnabled}, VarType: cmdconfig.EnvVarTypeBool},
		// common pipes env vars
		constants.EnvPipesInstallDir: {ConfigVar: []string{constants.ArgPipesInstallDir}, VarType: cmdconfig.EnvVarTypeString},
		constants.EnvPipesHost:       {ConfigVar: []string{constants.ArgPipesHost}, VarType: cmdconfig.EnvVarTypeString},
		constants.EnvPipesToken:      {ConfigVar: []string{constants.ArgPipesToken}, VarType: cmdconfig.EnvVarTypeString},
		// powerpipe specific constants
		localconstants.EnvListen:           {ConfigVar: []string{constants.ArgListen}, VarType: cmdconfig.EnvVarTypeString},
		localconstants.EnvBenchmarkTimeout: {ConfigVar: []string{constants.ArgBenchmarkTimeout}, VarType: cmdconfig.EnvVarTypeInt},
		localconstants.EnvDashboardTimeout: {ConfigVar: []string{constants.ArgDashboardTimeout}, VarType: cmdconfig.EnvVarTypeInt},
	}
}
