package cmdconfig

import (
	"github.com/turbot/pipe-fittings/app_specific"
	"github.com/turbot/pipe-fittings/cmdconfig"
	"github.com/turbot/pipe-fittings/constants"
	localconstants "github.com/turbot/powerpipe/internal/constants"
)

func configDefaults() map[string]any {
	return map[string]any{
		// global general options
		constants.ArgTelemetry:   constants.TelemetryInfo,
		constants.ArgUpdateCheck: true,

		// from global database options
		// TODO KAI NEEDED???
		//constants.ArgDatabasePort:         constants.DatabaseDefaultPort,
		//constants.ArgDatabaseStartTimeout: constants.DBStartTimeout.Seconds(),
		//constants.ArgServiceCacheEnabled:  true,
		//constants.ArgCacheMaxTtl:          300,

		// dashboard
		constants.ArgDashboardStartTimeout: constants.DashboardServiceStartTimeout.Seconds(),

		// memory
		constants.ArgMemoryMaxMb: 1024,
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
		app_specific.EnvIntrospection:     {ConfigVar: []string{constants.ArgIntrospection}, VarType: cmdconfig.EnvVarTypeString},
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
		constants.EnvPipesHost:            {ConfigVar: []string{constants.ArgPipesHost}, VarType: cmdconfig.EnvVarTypeString},
		constants.EnvPipesToken:           {ConfigVar: []string{constants.ArgPipesToken}, VarType: cmdconfig.EnvVarTypeString},
		// powerpipe specific constants
		localconstants.EnvListen: {ConfigVar: []string{constants.ArgListen}, VarType: cmdconfig.EnvVarTypeString},
		localconstants.EnvPort:   {ConfigVar: []string{constants.ArgPort}, VarType: cmdconfig.EnvVarTypeInt},
	}
}
