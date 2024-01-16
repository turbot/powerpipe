package backend

//go:generate go run golang.org/x/tools/cmd/stringer -type=BackendType

type BackendType int

const (
	UnknownClientBackend BackendType = iota
	PostgresDBClientBackend
	MySQLDBClientBackend
	SqliteDBClientBackend
	DuckDBClientBackend
)
