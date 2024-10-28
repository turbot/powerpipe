package db_client

import (
	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/backend"
	"github.com/turbot/pipe-fittings/connection"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/steampipeconfig"
	"github.com/turbot/powerpipe/internal/powerpipeconfig"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

func GetDatabaseConfigForResource(resource modconfig.ModTreeItem, workspaceMod *modconfig.Mod, defaultDatabase string, defaultSearchPathConfig backend.SearchPathConfig) (string, backend.SearchPathConfig, error) {
	database := defaultDatabase
	searchPathConfig := defaultSearchPathConfig

	// if there is no default search path, check if the mod has a search path
	// (it's database field may refer to a connection with a search path)
	if searchPathConfig.Empty() {
		searchPathConfig.SearchPath = workspaceMod.GetSearchPath()
		searchPathConfig.SearchPathPrefix = workspaceMod.GetSearchPathPrefix()
	}

	// NOTE: if the resource is in a dependency mod, check whether database or search path has been specified for it
	depName := resource.(modconfig.ModItem).GetMod().GetDependencyName()

	if depName != "" {
		// look for this mod in the workspace mod require
		modRequirement := workspaceMod.GetRequire().GetModDependency(depName)
		if modRequirement == nil {
			// not expected
			return database, searchPathConfig, sperr.New("could not find mod requirement for '%s' in workspace mod %s", depName, workspace
			 mod.ShortName)
		}

		// if the mod requirement has a search path, prefix or database, set it in viper,
		if modRequirement.Database != nil {
			// if database is overriden, also use overriden search path and search path prefix (even if empty)
			database = *modRequirement.Database
			searchPathConfig.SearchPath = modRequirement.SearchPath
			searchPathConfig.SearchPathPrefix = modRequirement.SearchPathPrefix
		}
		// if the parent mod has a database set, use it
		if modDb := resource.(modconfig.ModItem).GetMod().GetDatabase(); modDb != nil {
			database = *modDb
		}
		if modSearchPath := resource.(modconfig.ModItem).GetMod().GetSearchPath(); len(modSearchPath) > 0 {
			searchPathConfig.SearchPath = modSearchPath
		}
		if modSearchPathPrefix := resource.(modconfig.ModItem).GetMod().GetSearchPathPrefix(); len(modSearchPathPrefix) > 0 {
			searchPathConfig.SearchPathPrefix = modSearchPathPrefix
		}

	}

	// if the resource has a database set, use it
	if resource.GetDatabase() != nil {
		database = *resource.GetDatabase()
	}
	// if the resource has a search path set, use it
	if resourceSearchPath := resource.GetSearchPath(); len(resourceSearchPath) > 0 {
		searchPathConfig.SearchPath = resourceSearchPath
	}
	if resourceSearchPathPrefix := resource.GetSearchPathPrefix(); len(resourceSearchPathPrefix) > 0 {
		searchPathConfig.SearchPathPrefix = resourceSearchPathPrefix
	}

	// if the database is a cloud workspace, resolve the connection string
	if steampipeconfig.IsPipesWorkspaceIdentifier(database) {
		var err error
		database, err = GetCloudWorkspaceConnectionString(database)
		if err != nil {
			return "", backend.SearchPathConfig{}, err
		}
	}

	return database, searchPathConfig, nil
}

// GetDefaultDatabaseConfig builds the default database and searchPathConfig
// NOTE: if the dashboardUI has overridden the search path, opts wil be passed in to set the overridden value
func GetDefaultDatabaseConfig(opts ...backend.ConnectOption) (string, backend.SearchPathConfig, error) {
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

	// if no database is set, use the default connection
	if defaultDatabase == "" {
		defaultConnection := powerpipeconfig.GlobalConfig.GetDefaultConnection()
		defaultDatabase = defaultConnection.GetConnectionString()
		// if no search path has been set, use the default connection
		if defaultSearchPathConfig.Empty() {
			if spp, ok := defaultConnection.(connection.SearchPathProvider); ok {
				defaultSearchPathConfig = backend.SearchPathConfig{
					SearchPath:       spp.GetSearchPath(),
					SearchPathPrefix: spp.GetSearchPathPrefix(),
				}
			}
		}
	}

	// if the database is a cloud workspace, resolve the connection string
	if steampipeconfig.IsPipesWorkspaceIdentifier(defaultDatabase) {
		var err error
		defaultDatabase, err = GetCloudWorkspaceConnectionString(defaultDatabase)
		if err != nil {
			return "", backend.SearchPathConfig{}, err
		}
	}

	return defaultDatabase, defaultSearchPathConfig, nil
}
