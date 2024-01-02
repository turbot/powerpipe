package db_client

import (
	"context"
	"database/sql"
)

// TODO search path functions should first verify the DB is postgres and if not do nothing

// SetRequiredSessionSearchPath checks if either a search-path or search-path-prefix is set in config,
// and if so, set the search path (otherwise fall back to user search path)
// this just sets the required search path for this client
// - when creating a database session, we will actually set the searchPath
func (c *DbClient) SetRequiredSessionSearchPath(ctx context.Context) error {
	// TODO FIX THIS
	//configuredSearchPath := viper.GetStringSlice(constants.ArgSearchPath)
	//searchPathPrefix := viper.GetStringSlice(constants.ArgSearchPathPrefix)
	//
	//// strip empty elements from search path and prefix
	//configuredSearchPath = helpers.RemoveFromStringSlice(configuredSearchPath, "")
	//searchPathPrefix = helpers.RemoveFromStringSlice(searchPathPrefix, "")
	//
	//// default required path to user search path
	//requiredSearchPath := c.UserSearchPath
	//
	//// store custom search path and search path prefix
	//c.SearchPathPrefix = searchPathPrefix
	//
	//// if a search path was passed, use that
	//if len(configuredSearchPath) > 0 {
	//	requiredSearchPath = configuredSearchPath
	//}
	//
	//// add in the prefix if present
	//requiredSearchPath = db_common.AddSearchPathPrefix(searchPathPrefix, requiredSearchPath)
	//
	//requiredSearchPath = db_common.EnsureInternalSchemaSuffix(requiredSearchPath)
	//
	//// if either configuredSearchPath or SearchPathPrefix are set, store requiredSearchPath as CustomSearchPath
	//if len(configuredSearchPath)+len(searchPathPrefix) > 0 {
	//	c.CustomSearchPath = requiredSearchPath
	//} else {
	//	// otherwise clear it
	//	c.CustomSearchPath = nil
	//}

	return nil
}

func (c *DbClient) LoadUserSearchPath(ctx context.Context) error {
	// TODO KAI FIX THIS
	//conn, err := c.ManagementPool.Conn(ctx)
	//if err != nil {
	//	return err
	//}
	//defer conn.Close()
	//return c.LoadUserSearchPathForConnection(ctx, conn)
	return nil
}

func (c *DbClient) LoadUserSearchPathForConnection(ctx context.Context, connection *sql.Conn) error {
	// TODO KAI FIX THIS
	//// load the user search path
	//userSearchPath, err := db_common.GetUserSearchPath(ctx, connection)
	//if err != nil {
	//	return err
	//}
	//// update the cached value
	//c.UserSearchPath = userSearchPath
	return nil
}

// GetRequiredSessionSearchPath implements Client
func (c *DbClient) GetRequiredSessionSearchPath() []string {
	if c.CustomSearchPath != nil {
		return c.CustomSearchPath
	}

	return c.UserSearchPath
}

func (c *DbClient) GetCustomSearchPath() []string {
	return c.CustomSearchPath
}
