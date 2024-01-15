package backend

type duckdbRowReader struct {
	genericSQLRowReader
}

func NewDuckDBRowReader() *duckdbRowReader {
	return &duckdbRowReader{
		// use the generic row reader - to start with
		genericSQLRowReader: *NewGenericSQLRowReader(),
	}
}
