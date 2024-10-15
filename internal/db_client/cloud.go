package db_client

import (
	"context"
	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/cloud"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/powerpipe/internal/powerpipeconfig"
)

func GetCloudWorkspaceConnectionString(workspace string) (string, error) {
	// have we already retrieved this
	if connectionString, ok := powerpipeconfig.GlobalConfig.GetCloudConnectionString(workspace); ok {
		return connectionString, nil
	}

	token := viper.GetString(constants.ArgPipesToken)
	if token == "" {
		return "", error_helpers.MissingCloudTokenError()
	}
	cloudMetadata, err := cloud.GetCloudMetadata(context.Background(), workspace, token)
	if err != nil {
		return "", err
	}
	// cache
	powerpipeconfig.GlobalConfig.SetCloudConnectionString(workspace, cloudMetadata.ConnectionString)

	return cloudMetadata.ConnectionString, nil
}
