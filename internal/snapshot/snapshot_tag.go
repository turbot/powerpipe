package snapshot

import (
	"encoding/json"
	"reflect"

	"github.com/turbot/pipe-fittings/modconfig"
)

// TagColumn is the tag used to specify the column name and type in the introspection tables
const TagColumn = "column"

// SnapshotTag is a struct used to display column info in introspection tables
type SnapshotTag string

func newSnapshotTag(field reflect.StructField) *SnapshotTag {
	columnTag, ok := field.Tag.Lookup(TagColumn)
	if !ok {
		return nil
	}
	var res = SnapshotTag(columnTag)

	return &res
}

func GetAsSnapshotPropertyMap(item any) (map[string]any, error) {
	var res = make(map[string]any)

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
		snapshotTag := newSnapshotTag(field) // Assuming newSnapshotTag is a defined function
		if snapshotTag == nil {
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
			v = nil
		} else {
			v = fieldVal.Interface()
			if val.Kind() == reflect.Struct || val.Kind() == reflect.Slice {
				var target = make(map[string]any)
				jsonBytes, err := json.Marshal(v)
				if err != nil {
					return nil, err
				}
				err = json.Unmarshal(jsonBytes, &target)
				if err != nil {
					return nil, err
				}
				v = target
			}
		}

		if v != nil {
			res[string(*snapshotTag)] = v
		}
	}

	// tactical
	// add in name property from HclResource
	if hr, ok := item.(modconfig.HclResource); ok {
		res["name"] = hr.GetShortName()
	}
	return res, nil
}
