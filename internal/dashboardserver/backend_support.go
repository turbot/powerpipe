package dashboardserver

import (
	"github.com/turbot/pipe-fittings/v2/connection"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"log/slog"
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
