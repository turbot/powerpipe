package display

import (
	"github.com/turbot/pipe-fittings/printers"
	"reflect"
	"strings"
)

// TagColumn is the tag used to specify the column name and type in the introspection tables
const TagColumn = "column"

// ColumnTag is a struct used to display column info in introspection tables
type ColumnTag struct {
	Column string
	// the introspected go type
	ColumnType string
}

func newColumnTag(field reflect.StructField) (*ColumnTag, bool) {
	columnTag, ok := field.Tag.Lookup(TagColumn)
	if !ok {
		return nil, false
	}
	split := strings.Split(columnTag, ",")
	if len(split) != 2 {
		return nil, false
	}
	column := split[0]
	columnType := split[1]
	return &ColumnTag{column, columnType}, true
}

// GetAsTableRow returns the item as a table row
func GetAsTableRow(item interface{}) (printers.TableRow, []printers.TableColumnDefinition) {
	var columnDefs []printers.TableColumnDefinition
	var row = printers.TableRow{}
	t := reflect.TypeOf(item)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	val := reflect.ValueOf(item)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	for i := 0; i < val.NumField(); i++ {
		fieldName := val.Type().Field(i).Name
		field, _ := t.FieldByName(fieldName)
		columnTag, ok := newColumnTag(field) // Assuming newColumnTag is a defined function
		if !ok {
			continue
		}
		fieldVal := val.Field(i)
		if fieldVal.Kind() == reflect.Pointer {
			if !fieldVal.IsZero() {
				fieldVal = fieldVal.Elem()
			}
		}

		var v any
		if fieldVal.IsZero() {
			// todo handle different types
			v = ""
		} else {
			v = fieldVal.Interface() // This line retrieves the field value
		}
		row.Cells = append(row.Cells, v)

		columnDefs = append(columnDefs, printers.TableColumnDefinition{
			Name: columnTag.Column,
			Type: columnTag.ColumnType,
			// TODO KAI do wew reall yneed the description
		})
	}
	return row, columnDefs
}
