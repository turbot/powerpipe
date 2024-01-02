package cmdconfig

import (
	"github.com/turbot/pipe-fittings/app_specific"
	"github.com/turbot/pipe-fittings/cmdconfig"
	"github.com/turbot/pipe-fittings/constants"
)

var configDefaults = map[string]any{
	// global general options
	constants.ArgTelemetry:   constants.TelemetryInfo,
	constants.ArgUpdateCheck: true,

	// workspace profile
	constants.ArgAutoComplete:  true,
	constants.ArgIntrospection: constants.IntrospectionNone,

	// from global database options
	// TODO KAI NEEDED???
	//constants.ArgDatabasePort:         constants.DatabaseDefaultPort,
	//constants.ArgDatabaseStartTimeout: constants.DBStartTimeout.Seconds(),
	//constants.ArgServiceCacheEnabled:  true,
	//constants.ArgCacheMaxTtl:          300,

	// dashboard
	constants.ArgDashboardStartTimeout: constants.DashboardServiceStartTimeout.Seconds(),

	// memory
	constants.ArgMemoryMaxMbPlugin: 1024,
	constants.ArgMemoryMaxMb:       1024,
}

// environment variable mappings for directory paths which must be set as part of the viper bootstrap process
var dirEnvMappings = map[string]cmdconfig.EnvMapping{
	app_specific.EnvInstallDir:  {[]string{constants.ArgInstallDir}, cmdconfig.EnvVarTypeString},
	app_specific.EnvModLocation: {[]string{constants.ArgModLocation}, cmdconfig.EnvVarTypeString},
}

// NOTE: EnvWorkspaceProfile has already been set as a viper default as we have already loaded workspace profiles
// (EnvInstallDir has already been set at same time but we set it again to make sure it has the correct precedence)

// a map of known environment variables to map to viper keys - these are set as part of LoadGlobalConfig
var envMappings = map[string]cmdconfig.EnvMapping{
	app_specific.EnvInstallDir:    {[]string{constants.ArgInstallDir}, cmdconfig.EnvVarTypeString},
	app_specific.EnvModLocation:   {[]string{constants.ArgModLocation}, cmdconfig.EnvVarTypeString},
	app_specific.EnvIntrospection: {[]string{constants.ArgIntrospection}, cmdconfig.EnvVarTypeString},
	app_specific.EnvTelemetry:     {[]string{constants.ArgTelemetry}, cmdconfig.EnvVarTypeString},
	app_specific.EnvUpdateCheck:   {[]string{constants.ArgUpdateCheck}, cmdconfig.EnvVarTypeBool},
	// EnvPipesHost needs to be defined before EnvCloudHost,
	// so that if EnvCloudHost is defined, it can override EnvPipesHost
	constants.EnvPipesHost:    {[]string{constants.ArgCloudHost}, cmdconfig.EnvVarTypeString},
	app_specific.EnvCloudHost: {[]string{constants.ArgCloudHost}, cmdconfig.EnvVarTypeString},
	// EnvPipesToken needs to be defined before EnvCloudToken,
	// so that if EnvCloudToken is defined, it can override EnvPipesToken
	constants.EnvPipesToken:    {[]string{constants.ArgCloudToken}, cmdconfig.EnvVarTypeString},
	app_specific.EnvCloudToken: {[]string{constants.ArgCloudToken}, cmdconfig.EnvVarTypeString},
	//
	app_specific.EnvSnapshotLocation:  {[]string{constants.ArgSnapshotLocation}, cmdconfig.EnvVarTypeString},
	app_specific.EnvWorkspaceDatabase: {[]string{constants.ArgWorkspaceDatabase}, cmdconfig.EnvVarTypeString},
	app_specific.EnvDisplayWidth:      {[]string{constants.ArgDisplayWidth}, cmdconfig.EnvVarTypeInt},
	app_specific.EnvMaxParallel:       {[]string{constants.ArgMaxParallel}, cmdconfig.EnvVarTypeInt},
	app_specific.EnvQueryTimeout:      {[]string{constants.ArgDatabaseQueryTimeout}, cmdconfig.EnvVarTypeInt},
	app_specific.EnvCacheTTL:          {[]string{constants.ArgCacheTtl}, cmdconfig.EnvVarTypeInt},
	app_specific.EnvCacheMaxTTL:       {[]string{constants.ArgCacheMaxTtl}, cmdconfig.EnvVarTypeInt},
	app_specific.EnvMemoryMaxMb:       {[]string{constants.ArgMemoryMaxMb}, cmdconfig.EnvVarTypeInt},
	app_specific.EnvMemoryMaxMbPlugin: {[]string{constants.ArgMemoryMaxMbPlugin}, cmdconfig.EnvVarTypeInt},

	// we need this value to go into different locations
	app_specific.EnvCacheEnabled: {[]string{constants.ArgClientCacheEnabled, constants.ArgServiceCacheEnabled}, cmdconfig.EnvVarTypeBool},
}
