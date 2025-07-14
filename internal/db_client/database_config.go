package db_client

import (
	"log/slog"

	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/v2/backend"
	"github.com/turbot/pipe-fittings/v2/connection"
	"github.com/turbot/pipe-fittings/v2/constants"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/pipe-fittings/v2/steampipeconfig"
	"github.com/turbot/powerpipe/internal/powerpipeconfig"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

// GetDatabaseConfigForResource returns the ConnectionStringProvider and searchPathConfig for a resource
// if no database is set, use the default database, likewise for search path
func GetDatabaseConfigForResource(resource modconfig.ModTreeItem, workspaceMod *modconfig.Mod, defaultDatabase connection.ConnectionStringProvider, defaultSearchPathConfig backend.SearchPathConfig) (connection.ConnectionStringProvider, backend.SearchPathConfig, error) {
	csp := defaultDatabase
	searchPathConfig := defaultSearchPathConfig

	// if there is no default search path, check if the mod has a search path
	// (its database field may refer to a connection with a search path)
	if searchPathConfig.Empty() {
		searchPathConfig.SearchPath = workspaceMod.GetSearchPath()
		searchPathConfig.SearchPathPrefix = workspaceMod.GetSearchPathPrefix()
	}

	// NOTE: if the resource is in a dependency mod, check whether database or search path has been specified for it
	depName := resource.(modconfig.ModItem).GetMod().DependencyName

	if depName != "" {
		// look for this mod in the workspace mod require
		modRequirement := workspaceMod.Require.GetModDependency(depName)
		if modRequirement == nil {
			// not expected
			return csp, searchPathConfig, sperr.New("could not find mod requirement for '%s' in workspace mod %s", depName, workspaceMod.ShortName)
		}

		// if the mod requirement has a search path, prefix or database, set it in viper,
		if modRequirement.Database != nil {
			// if database is overridden, also use overriden search path and search path prefix (even if empty)
			csp = connection.NewConnectionString(*modRequirement.Database)
			searchPathConfig.SearchPath = modRequirement.SearchPath
			searchPathConfig.SearchPathPrefix = modRequirement.SearchPathPrefix
		}
		// if the parent mod has a database set, use it
		if modDb := resource.(modconfig.ModItem).GetMod().GetDatabase(); modDb != nil {
			csp = modDb
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
		csp = resource.GetDatabase()
	}
	// if the resource has a search path set, use it
	if resourceSearchPath := resource.GetSearchPath(); len(resourceSearchPath) > 0 {
		searchPathConfig.SearchPath = resourceSearchPath
	}
	if resourceSearchPathPrefix := resource.GetSearchPathPrefix(); len(resourceSearchPathPrefix) > 0 {
		searchPathConfig.SearchPathPrefix = resourceSearchPathPrefix
	}

	// if the database is a cloud workspace, resolve the connection string
	if steampipeconfig.IsPipesWorkspaceConnectionString(csp) {
		cs, err := csp.GetConnectionString()
		if err != nil {
			return nil, backend.SearchPathConfig{}, err
		}
		csp, err = GetPipesWorkspaceConnectionString(cs)
		if err != nil {
			return nil, backend.SearchPathConfig{}, err
		}
	}

	return csp, searchPathConfig, nil
}

// GetDefaultDatabaseConfig returns the default ConnectionStringProvider and searchPathConfig
// NOTE: if the dashboardUI has overridden the search path, opts wil be passed in to set the overridden value
func GetDefaultDatabaseConfig(mod *modconfig.Mod) (connection.ConnectionStringProvider, backend.SearchPathConfig, error) {
	// resolve the active database and search search path config for the dashboard
	defaultSearchPathConfig := backend.SearchPathConfig{
		SearchPath:       viper.GetStringSlice(constants.ArgSearchPath),
		SearchPathPrefix: viper.GetStringSlice(constants.ArgSearchPathPrefix),
	}

	var csp connection.ConnectionStringProvider
	// has a database arg been set in viper?
	databaseArgs := viper.GetString(constants.ArgDatabase)
	modDatabase := mod.GetDatabase()
	// if there is no mod database specified, see if we have mod dependencies
	// and if so, do they all specify the same database. If so we can use that as default
	if modDatabase == nil {
		modDatabase = getDefaultDatabaseFromModDependencies(mod)
	}

	switch {
	// if database command line was passed, set default
	case databaseArgs != "":
		slog.Info("GetDefaultDatabaseConfig: Using database connection string from command line", "database arg", databaseArgs)
		csp = connection.NewConnectionString(databaseArgs)
	case modDatabase != nil:
		slog.Info("GetDefaultDatabaseConfig: Using database connection string from mod")
		csp = modDatabase
	// if no database is set, use the default connection
	default:
		slog.Info("GetDefaultDatabaseConfig: Using default connection")

		defaultConnection := powerpipeconfig.GlobalConfig.GetDefaultConnection()
		csp = defaultConnection
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
	if steampipeconfig.IsPipesWorkspaceConnectionString(csp) {
		slog.Info("GetDefaultDatabaseConfig: Resolving Pipes workspace connection string")

		cs, err := csp.GetConnectionString()
		if err != nil {
			return nil, backend.SearchPathConfig{}, err
		}
		csp, err = GetPipesWorkspaceConnectionString(cs)
		if err != nil {
			return nil, backend.SearchPathConfig{}, err
		}
	}
	return csp, defaultSearchPathConfig, nil
}

func getDefaultDatabaseFromModDependencies(mod *modconfig.Mod) connection.ConnectionStringProvider {
	var depModDb connection.ConnectionStringProvider

	for _, dep := range mod.Resources.GetMods() {
		if dep.Name() == mod.Name() {
			continue
		}
		thisDb := dep.GetDatabase()
		if thisDb == nil {
			thisDb = powerpipeconfig.GlobalConfig.GetDefaultConnection()
		}
		if depModDb == nil {
			depModDb = thisDb
		} else if depModDb != dep.GetDatabase() {
			// more than one database specified - we cannot determine which to use
			return nil
		}
	}
	// so to get here, all mod dependencies use the same database ( which may be the default database)
	return depModDb
}
