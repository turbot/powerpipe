package cmdconfig

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/viper"
	filehelpers "github.com/turbot/go-kit/files"
	"github.com/turbot/pipe-fittings/connection"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/pipes"
	"github.com/turbot/pipe-fittings/steampipeconfig"
	"github.com/turbot/powerpipe/internal/powerpipeconfig"
)

// ValidateDatabaseArg checks if the database arg is a connection reference and resolves it if so
func ValidateDatabaseArg() error {
	databaseArg := viper.GetString(constants.ArgDatabase)
	if databaseArg == "" {
		return nil
	}
	if strings.HasPrefix(databaseArg, "connection.") {
		conn, ok := powerpipeconfig.GlobalConfig.PipelingConnections[strings.TrimPrefix(databaseArg, "connection.")]
		if !ok {
			return fmt.Errorf("connection '%s' not found", databaseArg)
		}

		csp, ok := conn.(connection.ConnectionStringProvider)
		if !ok {
			// unexpected - all registered connections should implement this interface
			return fmt.Errorf("connection '%s' does not implement connection.ConnectionStringProvider", databaseArg)
		}
		connectionString := csp.GetConnectionString()
		// update viper Database arg with the connection string
		viper.Set(constants.ArgDatabase, connectionString)
		// if no search path is set, set it to the connection's default search path
		if spp, ok := conn.(connection.SearchPathProvider); ok {
			if viper.GetString(constants.ArgSearchPath) == "" {
				viper.Set(constants.ArgSearchPath, spp.GetSearchPath())
			}

			if viper.GetString(constants.ArgSearchPathPrefix) == "" {
				viper.Set(constants.ArgSearchPathPrefix, spp.GetSearchPathPrefix())
			}
		}
	}

	return nil
}

func ValidateSnapshotArgs(ctx context.Context) error {
	// only 1 of 'share' and 'snapshot' may be set
	share := viper.GetBool(constants.ArgShare)
	snapshot := viper.GetBool(constants.ArgSnapshot)
	if share && snapshot {
		return fmt.Errorf("only 1 of 'share' and 'snapshot' may be set")
	}

	// if neither share or snapshot are set, nothing more to do
	if !share && !snapshot {
		return nil
	}

	token := viper.GetString(constants.ArgPipesToken)

	// determine whether snapshot location is a cloud workspace or a file location
	// if a file location, check it exists
	if err := validateSnapshotLocation(ctx, token); err != nil {
		return err
	}

	// TODO K this is probably broken as we do not really use database arg
	// if workspace-database or snapshot-location are a cloud workspace handle, cloud token must be set
	requireCloudToken := steampipeconfig.IsPipesWorkspaceIdentifier(viper.GetString(constants.ArgDatabase)) ||
		steampipeconfig.IsPipesWorkspaceIdentifier(viper.GetString(constants.ArgSnapshotLocation))

	// verify cloud token and workspace has been set
	if requireCloudToken && token == "" {
		return error_helpers.MissingCloudTokenError()
	}

	// should never happen as there is a default set
	if viper.GetString(constants.ArgPipesHost) == "" {
		return fmt.Errorf("to share snapshots, cloud host must be set")
	}

	return validateSnapshotTags()
}

func validateSnapshotLocation(ctx context.Context, cloudToken string) error {
	snapshotLocation := viper.GetString(constants.ArgSnapshotLocation)

	// if snapshot location is not set, set to the users default
	if snapshotLocation == "" {
		if cloudToken == "" {
			return error_helpers.MissingCloudTokenError()
		}
		return setSnapshotLocationFromDefaultWorkspace(ctx, cloudToken)
	}

	// if it is NOT a workspace handle, assume it is a local file location:
	// tildefy it and ensure it exists
	if !steampipeconfig.IsPipesWorkspaceIdentifier(snapshotLocation) {
		var err error
		snapshotLocation, err = filehelpers.Tildefy(snapshotLocation)
		if err != nil {
			return err
		}

		// write back to viper
		viper.Set(constants.ArgSnapshotLocation, snapshotLocation)

		if !filehelpers.DirectoryExists(snapshotLocation) {
			return fmt.Errorf("snapshot location %s does not exist", snapshotLocation)
		}
	}
	return nil
}

func setSnapshotLocationFromDefaultWorkspace(ctx context.Context, cloudToken string) error {
	workspaceHandle, err := pipes.GetUserWorkspaceHandle(ctx, cloudToken)
	if err != nil {
		return err
	}

	viper.Set(constants.ArgSnapshotLocation, workspaceHandle)
	return nil
}

func validateSnapshotTags() error {
	tags := viper.GetStringSlice(constants.ArgSnapshotTag)
	for _, tagStr := range tags {
		if len(strings.Split(tagStr, "=")) != 2 {
			return fmt.Errorf("snapshot tags must be specified '--%s key=value'", constants.ArgSnapshotTag)
		}
	}
	return nil
}
