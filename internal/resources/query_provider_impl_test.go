package resources

import (
	"testing"

	"github.com/turbot/pipe-fittings/v2/connection"
	"github.com/turbot/pipe-fittings/v2/modconfig"
)

func TestQueryProviderImpl_DatabaseInheritance(t *testing.T) {
	// Create a mock mod
	mod := &modconfig.Mod{
		ModTreeItemImpl: modconfig.ModTreeItemImpl{
			HclResourceImpl: modconfig.HclResourceImpl{
				ShortName: "test_mod",
			},
		},
	}

	// Create a query with a database set
	queryWithDB := &Query{
		QueryProviderImpl: QueryProviderImpl{
			RuntimeDependencyProviderImpl: RuntimeDependencyProviderImpl{
				ModTreeItemImpl: modconfig.ModTreeItemImpl{
					HclResourceImpl: modconfig.HclResourceImpl{
						ShortName: "test_query",
					},
					Mod: mod,
				},
			},
		},
	}

	// Set a database on the query
	testDB := connection.NewConnectionString("duckdb:/test.db")
	queryWithDB.SetDatabase(testDB)

	// Create a table that references the query
	table := &DashboardTable{
		QueryProviderImpl: QueryProviderImpl{
			RuntimeDependencyProviderImpl: RuntimeDependencyProviderImpl{
				ModTreeItemImpl: modconfig.ModTreeItemImpl{
					HclResourceImpl: modconfig.HclResourceImpl{
						ShortName: "test_table",
					},
					Mod: mod,
				},
			},
			Query: queryWithDB,
		},
	}

	// Test that the table inherits the database from the referenced query
	inheritedDB := table.GetDatabase()
	if inheritedDB == nil {
		t.Fatal("Expected table to inherit database from referenced query, but got nil")
	}

	// Verify it's the same database
	if inheritedDB != testDB {
		t.Fatal("Expected table to inherit the exact same database instance from referenced query")
	}
}

func TestQueryProviderImpl_DatabaseInheritanceWithOwnDB(t *testing.T) {
	// Create a mock mod
	mod := &modconfig.Mod{
		ModTreeItemImpl: modconfig.ModTreeItemImpl{
			HclResourceImpl: modconfig.HclResourceImpl{
				ShortName: "test_mod",
			},
		},
	}

	// Create a query with a database set
	queryWithDB := &Query{
		QueryProviderImpl: QueryProviderImpl{
			RuntimeDependencyProviderImpl: RuntimeDependencyProviderImpl{
				ModTreeItemImpl: modconfig.ModTreeItemImpl{
					HclResourceImpl: modconfig.HclResourceImpl{
						ShortName: "test_query",
					},
					Mod: mod,
				},
			},
		},
	}

	// Set a database on the query
	queryDB := connection.NewConnectionString("duckdb:/query.db")
	queryWithDB.SetDatabase(queryDB)

	// Create a table that references the query but has its own database
	table := &DashboardTable{
		QueryProviderImpl: QueryProviderImpl{
			RuntimeDependencyProviderImpl: RuntimeDependencyProviderImpl{
				ModTreeItemImpl: modconfig.ModTreeItemImpl{
					HclResourceImpl: modconfig.HclResourceImpl{
						ShortName: "test_table",
					},
					Mod: mod,
				},
			},
			Query: queryWithDB,
		},
	}

	// Set a different database on the table
	tableDB := connection.NewConnectionString("duckdb:/table.db")
	table.SetDatabase(tableDB)

	// Test that the table uses its own database, not the inherited one
	tableDBResult := table.GetDatabase()
	if tableDBResult == nil {
		t.Fatal("Expected table to use its own database, but got nil")
	}

	// Verify it's the table's database, not the query's
	if tableDBResult != tableDB {
		t.Fatal("Expected table to use its own database, not the inherited one")
	}

	if tableDBResult == queryDB {
		t.Fatal("Expected table to use its own database, not the inherited one")
	}
}
