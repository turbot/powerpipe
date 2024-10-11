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

	var cloudMetadata *steampipeconfig.CloudMetadata

	// TODO KAI we need to handle this being set for any resource

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
		// now update the connection string in viper database flag with the cloud metadata conneciotn string
		viper.Set(constants.ArgDatabase, cloudMetadata.ConnectionString)
	}

	return cloudMetadata, nil
}
