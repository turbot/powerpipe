package controlexecute

import (
	"encoding/json"

	"github.com/turbot/pipe-fittings/v2/queryresult"
	"github.com/turbot/pipe-fittings/v2/steampipeconfig"
	"github.com/turbot/powerpipe/internal/dashboardtypes"
)

// ControlRunData is the snapshot format data of a completed control run.
// It holds no rows itself: it converts the run's Rows when it is marshalled.
type ControlRunData struct {
	run     *ControlRun
	columns []*queryresult.ColumnDef
}

func newControlRunData(run *ControlRun) *ControlRunData {
	return &ControlRunData{run: run, columns: leafColumns(run.getDimensionSchema())}
}

// MarshalJSON implements json.Marshaler
func (d *ControlRunData) MarshalJSON() ([]byte, error) {
	rows := make([]map[string]any, len(d.run.Rows))
	for i, row := range d.run.Rows {
		rows[i] = row.leafRow()
	}
	return json.Marshal(&dashboardtypes.LeafData{Columns: d.columns, Rows: rows})
}

// a snapshot export writes the rows of a control run one at a time
var _ steampipeconfig.SnapshotRowsPanel = (*ControlRun)(nil)

// controlRunColumns is ControlRunData without its rows
type controlRunColumns struct {
	Columns []*queryresult.ColumnDef `json:"columns"`
}

// PanelWithoutRows implements steampipeconfig.SnapshotRowsPanel
func (r *ControlRun) PanelWithoutRows() any {
	var data *controlRunColumns
	if r.Data != nil {
		data = &controlRunColumns{Columns: r.Data.columns}
	}
	return struct {
		*ControlRun
		Data *controlRunColumns `json:"data"`
	}{r, data}
}

// EachRow implements steampipeconfig.SnapshotRowsPanel
func (r *ControlRun) EachRow(f func(row any) error) error {
	for _, row := range r.Rows {
		if err := f(row.leafRow()); err != nil {
			return err
		}
	}
	return nil
}
