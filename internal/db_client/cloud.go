package db_client

import (
	"context"
	"github.com/turbot/pipe-fittings/v2/connection"

	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/v2/constants"
	"github.com/turbot/pipe-fittings/v2/error_helpers"
	"github.com/turbot/pipe-fittings/v2/pipes"
	"github.com/turbot/powerpipe/internal/powerpipeconfig"
)

func GetPipesWorkspaceConnectionString(workspace string) (connection.ConnectionStringProvider, error) {
	// have we already retrieved this
	if connectionString, ok := powerpipeconfig.GlobalConfig.GetCloudConnectionString(workspace); ok {
		return connection.NewConnectionString(connectionString), nil
	}

	token := viper.GetString(constants.ArgPipesToken)
	if token == "" {
		return nil, error_helpers.MissingCloudTokenError()
	}
	pipesMetadata, err := pipes.GetPipesMetadata(context.Background(), workspace, token)
	if err != nil {
		return nil, err
	}
	// cache
	powerpipeconfig.GlobalConfig.SetCloudConnectionString(workspace, pipesMetadata.ConnectionString)

	return connection.NewConnectionString(pipesMetadata.ConnectionString), nil
}
