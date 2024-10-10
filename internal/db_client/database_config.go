package db_client

import (
	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/backend"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

func GetDatabaseConfigForResource(resource modconfig.ModTreeItem, workspaceMod *modconfig.Mod, defaultDatabase string, defaultSearchPathConfig backend.SearchPathConfig) (string, backend.SearchPathConfig, error) {
	database := defaultDatabase
	searchPathConfig := defaultSearchPathConfig

	// NOTE: if the resource is in a dependency mod, check whether database or search path has been specified for it
	depName := resource.GetMod().DependencyName

	if depName != "" {
		// look for this mod in the workspace mod require
		modRequirement := workspaceMod.Require.GetModDependency(depName)
		if modRequirement == nil {
			// not expected
			return database, searchPathConfig, sperr.New("could not find mod requirement for '%s' in workspace mod %s", depName, workspaceMod.ShortName)
		}

		// if the mod requirement has a search path, prefix or database, set it in viper,
		// overriding whatever value sth, use it
		// TODO should we only respect overridden search path and search path prefix if the db is overriden?
		//if modRequirement.OriginalSearchPath != nil {
		//	searchPathConfig.OriginalSearchPath = modRequirement.OriginalSearchPath
		//}
		//if modRequirement.SearchPathPrefix != nil {
		//	searchPathConfig.SearchPathPrefix = modRequirement.SearchPathPrefix
		//}
		if modRequirement.Database != nil {
			// if database is overriden, also use overriden search path and search path prefix (even if empty)
			database = *modRequirement.Database
			searchPathConfig.SearchPath = modRequirement.SearchPath
			searchPathConfig.SearchPathPrefix = modRequirement.SearchPathPrefix
		}
	}

	return database, searchPathConfig, nil

}

// GetDefaultDatabaseConfig builds the default database and searchPathConfig for the dashboard execution tree
// NOTE: if the dashboardUI has overridden the search path, opts wil be passed in to set the overridden value
func GetDefaultDatabaseConfig(opts ...backend.ConnectOption) (string, backend.SearchPathConfig) {
	var cfg backend.ConnectConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	// resolve the active database and search search path config for the dashboard
	defaultSearchPathConfig := backend.SearchPathConfig{
		SearchPath:       viper.GetStringSlice(constants.ArgSearchPath),
		SearchPathPrefix: viper.GetStringSlice(constants.ArgSearchPathPrefix),
	}
	// has the search path been overridden?
	if !cfg.SearchPathConfig.Empty() {
		defaultSearchPathConfig = cfg.SearchPathConfig
	}
	defaultDatabase := viper.GetString(constants.ArgDatabase)
	return defaultDatabase, defaultSearchPathConfig
}
