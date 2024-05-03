package queryresult

import (
	"github.com/turbot/pipe-fittings/v2/queryresult"
	"time"
)

type RowResult struct {
	Data  []interface{}
	Error error
}
type TimingMetadata struct {
	Duration time.Duration
}

type Result struct {
	RowChan *chan *RowResult
	Cols    []*queryresult.ColumnDef
	Timing  *TimingMetadata
}

func NewResult(cols []*queryresult.ColumnDef) *Result {

	rowChan := make(chan *RowResult)
	return &Result{
		RowChan: &rowChan,
		Cols:    cols,
	}
}

// IsExportSourceData implements ExportSourceData
func (*Result) IsExportSourceData() {}

// Close closes the row channel
func (r *Result) Close() {
	close(*r.RowChan)
}

func (r *Result) StreamRow(rowResult []interface{}) {
	*r.RowChan <- &RowResult{Data: rowResult}
}
func (r *Result) StreamError(err error) {
	*r.RowChan <- &RowResult{Error: err}
}

type SyncQueryResult struct {
	Rows []interface{}
	Cols []*queryresult.ColumnDef
}
