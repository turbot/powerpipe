package dashboardserver

import (
	"context"
	"log/slog"

	"github.com/turbot/pipe-fittings/v2/backend"
	"github.com/turbot/pipe-fittings/v2/connection"
	"github.com/turbot/pipe-fittings/v2/constants"
	"github.com/turbot/pipe-fittings/v2/modconfig"
)

type backendSupport struct {
	supportsSearchPath bool
	supportsTimeRange  bool
}

// setFromDb sets the backend support based on the database type
func (bs *backendSupport) setFromDb(db connection.ConnectionStringProvider) {
	if db != nil {
		switch db.(type) {
		case *connection.SteampipePgConnection, *connection.PostgresConnection:
			bs.supportsSearchPath = true
		case *connection.TailpipeConnection:
			bs.supportsTimeRange = true
		case *connection.ConnectionString:
			// create a backend and get its name
			// if it is a steampipe or postgres backend, set supportsSearchPath
			// NOTE: a tailpipe connection cannot be specified by a connection string
			// (as the connection string is dynamic and provided by the Tailpipe CLI),
			// so we will not set the supportsTimeRange flag here

			// we do not expect an error from ConnectionString.GetConnectionString
			connectionString, _ := db.GetConnectionString()
			// get the backend name from the connection string
			// NOTE: this does not create the backend and will therefore return postgres for a steampipe backend
			// as we cannot tell the difference purely from	the connection string
			// This is fine as we just want to determine whether a search path is supported
			backendName, _ := backend.NameFromConnectionString(context.Background(), connectionString)
			// set supportsSearchPath if the backend is a steampipe or postgres backend
			bs.supportsSearchPath = backendName == constants.PostgresBackendName
		}
	}
}

func newBackendSupport(database connection.ConnectionStringProvider) *backendSupport {
	bs := &backendSupport{}
	bs.setFromDb(database)
	return bs
}

// determineBackendSupport determines the backend support for a dashboard
// if no resource has a specified database, use the default database to set the backend support
func determineBackendSupport(dashboard modconfig.ModTreeItem, defaultDatabase connection.ConnectionStringProvider) backendSupport {
	bs := determineBackendSupportForResource(dashboard)
	if bs == nil {
		slog.Info("determineBackendSupport - no resource in the tree specifies a database, using default database")
		bs = newBackendSupport(defaultDatabase)
	}

	slog.Info("determineBackendSupport", "supportsSearchPath", bs.supportsSearchPath, "supportsTimeRange", bs.supportsTimeRange)
	return *bs
}

func determineBackendSupportForResource(item modconfig.ModTreeItem) *backendSupport {
	var bs *backendSupport

	// NOT: just check the database on this resource - GetDatabase also checks the parents
	// - there is no need to do that here as we are traversing down the tree
	if db := item.GetModTreeItemImpl().Database; db != nil {
		bs = &backendSupport{}
		bs.setFromDb(db)
		slog.Info("determineBackendSupportForResource - resource has database", "resource", item.Name(), "backendSupport", bs)
	}

	// if we have now set both flags, we can stop - no need to traverse further
	if backendSupportsAll(bs) {
		return bs
	}

	for _, child := range item.GetChildren() {
		childBs := determineBackendSupportForResource(child)
		// merge this with out ba
		bs = mergeBackendSupport(bs, childBs)

		// if we have now set both flags, we can stop - no need to traverse further
		if backendSupportsAll(bs) {
			return bs
		}
	}
	return bs
}

func backendSupportsAll(bs *backendSupport) bool {
	return bs != nil && bs.supportsSearchPath && bs.supportsTimeRange
}

// merge 2 backend support objects
func mergeBackendSupport(bs1, bs2 *backendSupport) *backendSupport {
	if bs1 == nil {
		return bs2
	}
	if bs2 == nil {
		return bs1
	}
	return &backendSupport{
		supportsSearchPath: bs1.supportsSearchPath || bs2.supportsSearchPath,
		supportsTimeRange:  bs1.supportsTimeRange || bs2.supportsTimeRange,
	}
}
