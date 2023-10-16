package initialisation

import (
	"context"
	"github.com/turbot/powerpipe/pkg/entities"
	"strings"

	"github.com/spf13/viper"
	"github.com/turbot/powerpipe/pkg/cloud"
	"github.com/turbot/powerpipe/pkg/constants"
	"github.com/turbot/powerpipe/pkg/error_helpers"
)

func getCloudMetadata(ctx context.Context) (*entities.CloudMetadata, error) {
	workspaceDatabase := viper.GetString(constants.ArgWorkspaceDatabase)
	if workspaceDatabase == "local" {
		// local database - nothing to do here
		return nil, nil
	}
	connectionString := workspaceDatabase

	var cloudMetadata *entities.CloudMetadata

	// so a backend was set - is it a connection string or a database name
	workspaceDatabaseIsConnectionString := strings.HasPrefix(workspaceDatabase, "postgresql://") || strings.HasPrefix(workspaceDatabase, "postgres://")
	if !workspaceDatabaseIsConnectionString {
		// it must be a database name - verify the cloud token was provided
		cloudToken := viper.GetString(constants.ArgCloudToken)
		if cloudToken == "" {
			return nil, error_helpers.MissingCloudTokenError
		}

		// so we have a database and a token - build the connection string and set it in viper
		var err error
		if cloudMetadata, err = cloud.GetCloudMetadata(ctx, workspaceDatabase, cloudToken); err != nil {
			return nil, err
		}
		// read connection string out of cloudMetadata
		connectionString = cloudMetadata.ConnectionString
	}

	// now set the connection string in viper
	viper.Set(constants.ArgConnectionString, connectionString)

	return cloudMetadata, nil
}
