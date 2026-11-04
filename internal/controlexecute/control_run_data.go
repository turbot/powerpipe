package controlexecute

import (
	"encoding/json"

	"github.com/turbot/pipe-fittings/v2/queryresult"
)

// ControlRunData is the snapshot format data of a completed control run.
// It holds no rows itself: it converts the run's Rows when it is marshalled.
type ControlRunData struct {
	run              *ControlRun
	dimensionsSchema map[string]*queryresult.ColumnDef
}

// MarshalJSON implements json.Marshaler
func (d *ControlRunData) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.run.Rows.ToLeafData(d.dimensionsSchema))
}
