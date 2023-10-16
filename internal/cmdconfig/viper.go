package cmdconfig

import (
	"fmt"
	"os"

	filehelpers "github.com/turbot/go-kit/files"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/powerpipe/pkg/cmdconfig"
	"github.com/turbot/steampipe/pkg/constants"
)

// bootstrapViper sets up viper with the essential path config (workspace-chdir and install-dir)
func bootstrapViper(cmd *cobra.Command) error {
	// set defaults  for keys which do not have a corresponding command flag
	setBaseDefaults()

	// set defaults from defaultWorkspaceProfile
	SetDefaultsFromConfig(make(map[string]interface{}))

	// set defaults for install dir and mod location from env vars
	// this needs to be done since the workspace profile definitions exist in the
	// default install dir
	setDirectoryDefaultsFromEnv()

	// tildefy all paths in viper
	return tildefyPaths()
}

// tildefyPaths cleans all path config values and replaces '~' with the home directory
func tildefyPaths() error {
	pathArgs := []string{
		constants.ArgModLocation,
		constants.ArgInstallDir,
	}
	var err error
	for _, argName := range pathArgs {
		if argVal := viper.GetString(argName); argVal != "" {
			if argVal, err = filehelpers.Tildefy(argVal); err != nil {
				return err
			}
			if viper.IsSet(argName) {
				// if the value was already set re-set
				viper.Set(argName, argVal)
			} else {
				// otherwise just update the default
				viper.SetDefault(argName, argVal)
			}
		}
	}
	return nil
}

// SetDefaultsFromConfig overrides viper default values from hcl config values
func SetDefaultsFromConfig(configMap map[string]interface{}) {
	for k, v := range configMap {
		viper.SetDefault(k, v)
	}
}

// for keys which do not have a corresponding command flag, we need a separate defaulting mechanism
// any option setting, workspace profile property or env var which does not have a command line
// MUST have a default (unless we want the zero value to take effect)
//
// Do not add keys here which have command line defaults - the way this is setup, this value takes
// precedence over command line default
func setBaseDefaults() {
	defaults := map[string]interface{}{
		// global general options
		constants.ArgTelemetry:   constants.TelemetryInfo,
		constants.ArgUpdateCheck: true,

		// workspace profile
		constants.ArgAutoComplete:  true,
		constants.ArgIntrospection: constants.IntrospectionNone,

		// from global database options
		constants.ArgDatabasePort:         constants.DatabaseDefaultPort,
		constants.ArgDatabaseStartTimeout: constants.DBStartTimeout.Seconds(),
		constants.ArgServiceCacheEnabled:  true,
		constants.ArgCacheMaxTtl:          300,

		// dashboard
		constants.ArgDashboardStartTimeout: constants.DashboardStartTimeout.Seconds(),

		// memory
		constants.ArgMemoryMaxMbPlugin: 1024,
		constants.ArgMemoryMaxMb:       1024,
	}

	for k, v := range defaults {
		viper.SetDefault(k, v)
	}
}

type envMapping struct {
	configVar []string
	varType   cmdconfig.EnvVarType
}

// set default values of INSTALL_DIR and ModLocation from env vars
func setDirectoryDefaultsFromEnv() {
	envMappings := map[string]envMapping{
		constants.EnvInstallDir:     {[]string{constants.ArgInstallDir}, cmdconfig.String},
		constants.EnvWorkspaceChDir: {[]string{constants.ArgModLocation}, cmdconfig.String},
		constants.EnvModLocation:    {[]string{constants.ArgModLocation}, cmdconfig.String},
	}

	for envVar, mapping := range envMappings {
		setConfigFromEnv(envVar, mapping.configVar, mapping.varType)
	}
}

// setDefaultsFromEnv sets default values from env vars
func setDefaultsFromEnv() {
	// NOTE: EnvWorkspaceProfile has already been set as a viper default as we have already loaded workspace profiles
	// (EnvInstallDir has already been set at same time but we set it again to make sure it has the correct precedence)

	// a map of known environment variables to map to viper keys
	envMappings := map[string]envMapping{
		constants.EnvInstallDir:    {[]string{constants.ArgInstallDir}, cmdconfig.String},
		constants.EnvModLocation:   {[]string{constants.ArgModLocation}, cmdconfig.String},
		constants.EnvIntrospection: {[]string{constants.ArgIntrospection}, cmdconfig.String},
		constants.EnvTelemetry:     {[]string{constants.ArgTelemetry}, cmdconfig.String},
		constants.EnvUpdateCheck:   {[]string{constants.ArgUpdateCheck}, cmdconfig.Bool},
		// PIPES_HOST needs to be defined before STEAMPIPE_CLOUD_HOST,
		// so that if STEAMPIPE_CLOUD_HOST is defined, it can override PIPES_HOST
		constants.EnvPipesHost: {[]string{constants.ArgCloudHost}, cmdconfig.String},
		constants.EnvCloudHost: {[]string{constants.ArgCloudHost}, cmdconfig.String},
		// PIPES_TOKEN needs to be defined before STEAMPIPE_CLOUD_TOKEN,
		// so that if STEAMPIPE_CLOUD_TOKEN is defined, it can override PIPES_TOKEN
		constants.EnvPipesToken: {[]string{constants.ArgCloudToken}, cmdconfig.String},
		constants.EnvCloudToken: {[]string{constants.ArgCloudToken}, cmdconfig.String},
		//
		constants.EnvSnapshotLocation:  {[]string{constants.ArgSnapshotLocation}, cmdconfig.String},
		constants.EnvWorkspaceDatabase: {[]string{constants.ArgWorkspaceDatabase}, cmdconfig.String},
		constants.EnvDisplayWidth:      {[]string{constants.ArgDisplayWidth}, cmdconfig.Int},
		constants.EnvMaxParallel:       {[]string{constants.ArgMaxParallel}, cmdconfig.Int},
		constants.EnvQueryTimeout:      {[]string{constants.ArgDatabaseQueryTimeout}, cmdconfig.Int},
		constants.EnvCacheTTL:          {[]string{constants.ArgCacheTtl}, cmdconfig.Int},
		constants.EnvCacheMaxTTL:       {[]string{constants.ArgCacheMaxTtl}, cmdconfig.Int},
		constants.EnvMemoryMaxMb:       {[]string{constants.ArgMemoryMaxMb}, cmdconfig.Int},
		constants.EnvMemoryMaxMbPlugin: {[]string{constants.ArgMemoryMaxMbPlugin}, cmdconfig.Int},

		// we need this value to go into different locations
		constants.EnvCacheEnabled: {[]string{
			constants.ArgClientCacheEnabled,
			constants.ArgServiceCacheEnabled,
		}, cmdconfig.Bool},
	}

	for envVar, v := range envMappings {
		setConfigFromEnv(envVar, v.configVar, v.varType)
	}
}

func setConfigFromEnv(envVar string, configs []string, varType cmdconfig.EnvVarType) {
	for _, configVar := range configs {
		SetDefaultFromEnv(envVar, configVar, varType)
	}
}

func SetDefaultFromEnv(k string, configVar string, varType cmdconfig.EnvVarType) {
	if val, ok := os.LookupEnv(k); ok {
		switch varType {
		case cmdconfig.String:
			viper.SetDefault(configVar, val)
		case cmdconfig.Bool:
			if boolVal, err := types.ToBool(val); err == nil {
				viper.SetDefault(configVar, boolVal)
			}
		case cmdconfig.Int:
			if intVal, err := types.ToInt64(val); err == nil {
				viper.SetDefault(configVar, intVal)
			}
		default:
			// must be an invalid value in the map above
			panic(fmt.Sprintf("invalid env var mapping type: %s", varType))
		}
	}
}
