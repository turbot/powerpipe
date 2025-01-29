package dashboardtypes

import (
	"fmt"

	"github.com/turbot/pipe-fittings/v2/queryresult"
	"github.com/turbot/pipe-fittings/v2/utils"
)

type LeafData struct {
	Columns []*queryresult.ColumnDef `json:"columns"`
	Rows    []map[string]interface{} `json:"rows"`
}

func NewLeafData(result *queryresult.SyncQueryResult) (*LeafData, error) {
	leafData := &LeafData{
		Rows:    make([]map[string]interface{}, len(result.Rows)),
		Columns: result.Cols,
	}
	// handle duplicate column names - this checks all column names and ensures they are unique
	// if they are not, assign a unique name to the column
	if err := leafData.ensureUniqueColumnName(); err != nil {
		return nil, err
	}

	for rowIdx, row := range result.Rows {
		rowData := make(map[string]interface{}, len(result.Cols))
		for i, data := range row.(*queryresult.RowResult).Data {
			// get unique column name from column defs
			// (NOTE: this may be either the original column name - if there are no duplicates,
			// or a specially generated unique name if there are duplicates)
			columnName := leafData.Columns[i].Name
			rowData[columnName] = data
		}

		leafData.Rows[rowIdx] = rowData
	}
	return leafData, nil
}

func (leafData *LeafData) ensureUniqueColumnName() error {
	// create a unique name generator
	nameGenerator := utils.NewUniqueNameGenerator()

	for colIdx, col := range leafData.Columns {
		uniqueName, err := nameGenerator.GetUniqueName(col.Name, colIdx)
		if err != nil {
			return fmt.Errorf("error generating unique column name: %w", err)
		}
		// if the column name has changed, store the original name and update the column name to be the unique name
		if uniqueName != col.Name {
			// set the original name first, BEFORE mutating name
			col.OriginalName = col.Name
			col.Name = uniqueName
		}
	}
	return nil
}
