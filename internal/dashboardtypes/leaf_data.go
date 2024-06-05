package dashboardtypes

import (
	"github.com/turbot/pipe-fittings/queryresult"
	"github.com/turbot/pipe-fittings/utils"
	localqueryresult "github.com/turbot/powerpipe/internal/queryresult"
)

type LeafData struct {
	Columns []*queryresult.ColumnDef `json:"columns"`
	Rows    []map[string]interface{} `json:"rows"`
}

func NewLeafData(result *localqueryresult.SyncQueryResult) *LeafData {
	leafData := &LeafData{
		Rows:    make([]map[string]interface{}, len(result.Rows)),
		Columns: result.Cols,
	}
	// handle duplicate column names
	leafData.ensureUniqueColumnName()

	for rowIdx, row := range result.Rows {
		rowData := make(map[string]interface{}, len(result.Cols))
		for i, data := range row.(*localqueryresult.RowResult).Data {
			// get unique column name from column defs
			columnName := leafData.Columns[i].Name
			rowData[columnName] = data
		}

		leafData.Rows[rowIdx] = rowData
	}
	return leafData
}

func (leafData *LeafData) ensureUniqueColumnName() {
	// create a unique name generator
	nameGenerator := utils.NewUniqueNameGenerator()

	for i, col := range leafData.Columns {
		uniqueName := nameGenerator.GetUniqueName(col.Name)
		if uniqueName != col.Name {
			// if the column name has changed, store the original
			leafData.Columns[i].Name = uniqueName
			leafData.Columns[i].OriginalName = col.Name
		}
	}
}
