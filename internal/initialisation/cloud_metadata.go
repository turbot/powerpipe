package initialisation

import (
	"context"
	"strings"

	"github.com/turbot/powerpipe/pkg/entities"

	"github.com/spf13/viper"
	"github.com/turbot/powerpipe/pkg/cloud"
	"github.com/turbot/powerpipe/pkg/constants"
	"github.com/turbot/powerpipe/pkg/error_helpers"
)

func getCloudMetadata(ctx context.Context) (*entities.CloudMetadata, error) {
	// TODO PSKR remove hardcoding
	workspaceDatabase := "postgresql://pskrbasu:ee7d-47fc-9672@spipetools-toolstest.usea1.db.pipes.turbot.com:9193/wffk04"
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
