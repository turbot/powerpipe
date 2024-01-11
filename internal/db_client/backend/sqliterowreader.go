package backend

type sqliteRowReader struct {
	genericSQLRowReader
}

func NewSqliteRowReader() *sqliteRowReader {
	return &sqliteRowReader{
		// use the generic row reader - there's no real difference between sqlite and generic
		genericSQLRowReader: *NewGenericSQLRowReader(),
	}
}
