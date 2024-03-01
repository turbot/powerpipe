package initialisation

import (
	"context"

	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/backend"
	"github.com/turbot/pipe-fittings/cloud"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/steampipeconfig"
)

func getCloudMetadata(ctx context.Context) (*steampipeconfig.CloudMetadata, error) {
	database := viper.GetString(constants.ArgDatabase)
	connectionString := database

	var cloudMetadata *steampipeconfig.CloudMetadata

	// so a backend was set - is it a connection string or a database name
	workspaceDatabaseIsConnectionString := backend.HasBackend(database)
	if !workspaceDatabaseIsConnectionString {
		// it must be a database name - verify the cloud token was provided
		cloudToken := viper.GetString(constants.ArgPipesToken)
		if cloudToken == "" {
			return nil, error_helpers.MissingCloudTokenError()
		}

		// so we have a database and a token - build the connection string and set it in viper
		var err error
		if cloudMetadata, err = cloud.GetCloudMetadata(ctx, database, cloudToken); err != nil {
			return nil, err
		}
		// read connection string out of cloudMetadata
		connectionString = cloudMetadata.ConnectionString
	}

	// now set the connection string in viper
	viper.Set(constants.ArgDatabase, connectionString)

	return cloudMetadata, nil
}
