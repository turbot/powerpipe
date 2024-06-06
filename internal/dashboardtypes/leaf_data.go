package dashboardtypes

import (
	"fmt"
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
	// handle duplicate column names - this checks all column names and ensures they are unique
	// if they are not, assign a unique name to the column
	leafData.ensureUniqueColumnName()

	for rowIdx, row := range result.Rows {
		rowData := make(map[string]interface{}, len(result.Cols))
		for i, data := range row.(*localqueryresult.RowResult).Data {
			// get unique column name from column defs
			// (NOTE: this may be either the origina column name - if there are no duplicates,
			// or a specially generated unique name if there are duplicates)
			columnName := leafData.Columns[i].GetUniqueName()
			rowData[columnName] = data
		}

		leafData.Rows[rowIdx] = rowData
	}
	return leafData
}

func (leafData *LeafData) ensureUniqueColumnName() error {
	// create a unique name generator
	nameGenerator := utils.NewUniqueNameGenerator()

	for i, col := range leafData.Columns {
		uniqueName, err := nameGenerator.GetUniqueName(col.Name)
		if err != nil {
			return fmt.Errorf("error generating unique column name: %w", err)
		}
		if uniqueName != col.Name {
			// if the column name has changed, store the modified name as UniqueName
			leafData.Columns[i].UniqueName = uniqueName
		}
	}
	return nil
}
