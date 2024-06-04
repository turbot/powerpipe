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

	// create a unique name generator
	nameGenerator := utils.NewUniqueNameGenerator()

	for rowIdx, row := range result.Rows {
		rowData := make(map[string]interface{}, len(result.Cols))
		for i, data := range row.(*localqueryresult.RowResult).Data {
			// ensure column name is unique
			columnName := leafData.Columns[i].Name
			uniqueName := nameGenerator.GetUniqueName(columnName)
			if uniqueName != columnName {
				// if the column name has changed, store the original
				leafData.Columns[i].OriginalName = columnName
				leafData.Columns[i].Name = uniqueName
			}
			rowData[uniqueName] = data
		}

		leafData.Rows[rowIdx] = rowData
	}
	return leafData
}
