package snapshot

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/r3labs/diff/v3"
	"github.com/turbot/go-kit/helpers"
)

// Store the keyed rows for each panel type
var PanelDataKeyMap = map[string][]string{
	"card":                {},
	"chart":               {}, // Empty array indicates diffing is required, but the keys need to be determined from the column properties
	"control":             {"resource"},
	"detection":           {},
	"detection_benchmark": {},
	"edge":                {"from_id", "to_id"},
	"flow":                {"id", "from_id", "to_id"},
	"graph":               {"id", "from_id", "to_id"},
	"hierarchy":           {"id", "from_id", "to_id"},
	"node":                {"id"},
	"table":               {},
}

var PanelKeyRegex = regexp.MustCompile(`((?:[^:;\\]|\\.)+):((?:[^:;\\]|\\.)+)(?:;|$)`)

type DiffPaths struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
}

type Snapshot struct {
	SchemaVersion string                 `json:"schema_version" diff:"schema_version"`
	Inputs        map[string]interface{} `json:"inputs" diff:"inputs"`
	Variables     map[string]interface{} `json:"variables" diff:"variables"`
	StartTime     string                 `json:"start_time" diff:"start_time"`
	EndTime       string                 `json:"end_time" diff:"end_time"`
	Layout        Layout                 `json:"layout" diff:"layout"`
	Metadata      map[string]interface{} `json:"metadata" diff:"metadata"`
	Panels        map[string]Panel       `json:"panels" diff:"panels"`
}

type Layout struct {
	Name      string   `json:"name" diff:"name"`
	PanelType string   `json:"panel_type" diff:"panel_type"`
	Children  []Layout `json:"children" diff:"children"`
}

type Panel struct {
	Dashboard   string                 `json:"dashboard" diff:"dashboard"`
	Description string                 `json:"description" diff:"description"`
	Name        string                 `json:"name" diff:"name"`
	PanelType   string                 `json:"panel_type" diff:"panel_type"`
	Status      string                 `json:"status" diff:"status"`
	Tags        map[string]string      `json:"tags" diff:"tags"`
	Title       string                 `json:"title" diff:"title"`
	Summary     map[string]interface{} `json:"summary" diff:"summary"`
	Properties  map[string]interface{} `json:"properties" diff:"properties"`
	Data        PanelData              `json:"data" diff:"data"`
}

type PanelData struct {
	Columns []interface{} `json:"columns" diff:"columns"`
	Rows    []interface{} `json:"rows" diff:"rows"`
}

type PanelColumnProperties struct {
	Name     string `json:"name" diff:"name"`
	Display  string `json:"display,omitempty" diff:"display"`
	DiffMode string `json:"diff_mode,omitempty" diff:"diff_mode"`
}

type PanelDataDiffer struct{}

// Whether this differ should be used to match a specific type
func (d *PanelDataDiffer) Match(a, b reflect.Value) bool {
	return diff.AreType(a, b, reflect.TypeOf(PanelData{}))
}

// The actual diff function, where you also append to the changelog
// using your custom format
func (d *PanelDataDiffer) Diff(diffType diff.DiffType, diffFunc diff.DiffFunc, cl *diff.Changelog, path []string, current, previous reflect.Value, parentData interface{}) error {
	// The parent data here refers to the panel as a whole, cast into panel data to determine the panel type
	panel, ok := parentData.(Panel)
	if !ok {
		return fmt.Errorf("The parent data is not of type Panel")
	}

	// Get the keys for the panel data based on the type of the panel being processed
	panelKeys, ok := PanelDataKeyMap[panel.PanelType]
	if !ok {
		return nil
	}
	// If keys for the panel type are not defined, we need to extract the key information from the panel properties section
	if len(panelKeys) == 0 {
		// Unmarshal Panel Properties
		var panelColumnProperties map[string]PanelColumnProperties
		if _, ok := panel.Properties["columns"]; ok {
			columnPropertiesRaw, _ := json.Marshal(panel.Properties["columns"])
			_ = json.Unmarshal(columnPropertiesRaw, &panelColumnProperties)
		}
		// Loop through the panel properties to extract the keys for the panel data
		for column, properties := range panelColumnProperties {
			if properties.DiffMode == "key" {
				panelKeys = append(panelKeys, column)
			}
		}
	}

	// If the panel type is `card` and there are no keys defined, we extract the keys from the column data
	if panel.PanelType == "card" && len(panelKeys) == 0 {
		currentPanelData, _ := current.Interface().(PanelData)
		if len(currentPanelData.Columns) > 1 {
			panelKeys = append(panelKeys, "label")
		}
	}

	// Proceed with the diffing process based on the keys extracted
	generatePanelDataDiff(cl, panel, path, current, previous, panelKeys)

	return nil
}

// unsure what this is actually for, but you must implement it either way
func (d *PanelDataDiffer) InsertParentDiffer(dfunc func(path []string, a, b reflect.Value, p interface{}) error) {
	return
}

func generatePanelDataDiff(cl *diff.Changelog, panel Panel, path []string, current, previous reflect.Value, panelRowKeys []string) {
	if current.Kind() == reflect.Invalid {
		cl.Add(diff.CREATE, path, nil, previous.Interface())
		return
	}
	if previous.Kind() == reflect.Invalid {
		cl.Add(diff.DELETE, path, current.Interface(), nil)
		return
	}

	// Cast the current and previous panel data
	var currentPanelData, previousPanelData PanelData
	previousPanelData, _ = current.Interface().(PanelData)
	currentPanelData, _ = previous.Interface().(PanelData)

	// Section to process columns
	// `name` would be the primary key for columns in every case
	currentColumnMap := make(map[string]interface{})
	for i := 0; i < len(currentPanelData.Columns); i++ {
		col := reflect.ValueOf(currentPanelData.Columns[i])
		currentColumnMap[col.Interface().(map[string]interface{})["name"].(string)] = col.Interface()
	}

	previousColumnMap := make(map[string]interface{})
	for i := 0; i < len(previousPanelData.Columns); i++ {
		col := reflect.ValueOf(previousPanelData.Columns[i])
		previousColumnMap[col.Interface().(map[string]interface{})["name"].(string)] = col.Interface()
	}

	for currentColumnKey, currentColumnObj := range currentColumnMap {
		// Column not present in previous
		if previousColumnObj, ok := previousColumnMap[currentColumnKey]; !ok {
			currentPath := append(path, "columns", fmt.Sprintf("name:%s", currentColumnKey))
			cl.Add(diff.CREATE, currentPath, nil, currentColumnObj)
		} else {
			// Column present in previous
			// Check to see if any of the column attributes have changed
			currentColumn := reflect.ValueOf(currentColumnObj).Interface().(map[string]interface{})
			previousColumn := reflect.ValueOf(previousColumnObj).Interface().(map[string]interface{})
			for key := range currentColumn {
				currentPath := append(path, "columns", fmt.Sprintf("name:%s", currentColumnKey), key)
				// Column attribute present in previous
				if _, ok := previousColumn[key]; ok {
					// If values are not same, generate a diff
					if currentColumn[key] != previousColumn[key] {
						cl.Add(diff.UPDATE, currentPath, previousColumn[key], currentColumn[key])
					}
					// Delete from the previous map
					delete(previousColumn, key)
				} else {
					cl.Add(diff.CREATE, currentPath, nil, currentColumn[key])
				}
			}
			// Any additional attributes in the previous map can be treated as new
			for key, value := range previousColumn {
				previousPath := append(path, "columns", fmt.Sprintf("name:%s", currentColumnKey), key)
				cl.Add(diff.DELETE, previousPath, value, nil)
			}

			delete(previousColumnMap, currentColumnKey)
		}
	}

	// Column added to previous
	for previousColumnKey, previousColumnObj := range previousColumnMap {
		previousPath := append(path, "columns", fmt.Sprintf("name:%s", previousColumnKey))
		cl.Add(diff.DELETE, previousPath, previousColumnObj, nil)
	}

	// End section to process columns

	// Section to process rows
	if len(panelRowKeys) > 0 {
		currentRowMap := make(map[string]interface{})
		for i := 0; i < len(currentPanelData.Rows); i++ {
			// Extract the row map from the interface
			row := reflect.ValueOf(currentPanelData.Rows[i]).Interface().(map[string]interface{})
			// Frame the key for the row based on the keys of the panel
			rowKey := ""
			for _, key := range panelRowKeys {
				switch row[key].(type) {
				case string:
					rowKey += fmt.Sprintf(";%s:%s", escape(key), escape(row[key].(string)))
				case int:
					rowKey += fmt.Sprintf(";%s:%d", escape(key), row[key].(int))
				case float64:
					rowKey += fmt.Sprintf(";%s:%f", escape(key), row[key].(float64))
				default:
					rowKey += fmt.Sprintf(";%s:%s", escape(key), escape(fmt.Sprintf("%v", row[key])))
				}
			}
			if rowKey != "" {
				rowKey = rowKey[1:]
				currentRowMap[rowKey] = row
			}
		}

		previousRowMap := make(map[string]interface{})
		for i := 0; i < len(previousPanelData.Rows); i++ {
			// Extract the row map from the interface
			row := reflect.ValueOf(previousPanelData.Rows[i]).Interface().(map[string]interface{})
			// Frame the key for the row based on the keys of the panel
			rowKey := ""
			for _, key := range panelRowKeys {
				switch row[key].(type) {
				case string:
					rowKey += fmt.Sprintf(";%s:%s", escape(key), escape(row[key].(string)))
				case int:
					rowKey += fmt.Sprintf(";%s:%d", escape(key), row[key].(int))
				case float64:
					rowKey += fmt.Sprintf(";%s:%f", escape(key), row[key].(float64))
				default:
					rowKey += fmt.Sprintf(";%s:%s", escape(key), escape(fmt.Sprintf("%v", row[key])))
				}
			}
			if rowKey != "" {
				rowKey = rowKey[1:]
				previousRowMap[rowKey] = row
			}
		}

		for currentRowKey, currentRowObj := range currentRowMap {
			// Row not present in previous
			if previousRowObj, ok := previousRowMap[currentRowKey]; !ok {
				currentPath := append(path, "rows", currentRowKey)
				cl.Add(diff.CREATE, currentPath, nil, currentRowObj)
			} else {
				// Row present in previous
				// Check to see if any of the column attributes have changed
				currentRow := reflect.ValueOf(currentRowObj).Interface().(map[string]interface{})
				previousRow := reflect.ValueOf(previousRowObj).Interface().(map[string]interface{})
				for key := range currentRow {
					currentPath := append(path, "rows", currentRowKey, key)
					// Row attribute present in previous
					if _, ok := previousRow[key]; ok {
						// If values are not same, generate a diff
						if currentRow[key] != previousRow[key] {
							cl.Add(diff.UPDATE, currentPath, previousRow[key], currentRow[key])
						}
						// Delete from the previous map
						delete(previousRow, key)
					} else {
						cl.Add(diff.UPDATE, currentPath, nil, currentRow[key])
					}
				}
				// Any additional attributes in the previous map can be treated as new
				for key, value := range previousRow {
					previousPath := append(path, "rows", currentRowKey, key)
					cl.Add(diff.UPDATE, previousPath, value, nil)
				}

				delete(previousRowMap, currentRowKey)
			}
		}

		// Row added to previous
		for previousRowKey, previousRowObj := range previousRowMap {
			previousPath := append(path, "rows", previousRowKey)
			cl.Add(diff.DELETE, previousPath, previousRowObj, nil)
		}
	} else {
		// For the best effort diffing, we will generate a hash for each row of the source and previous
		// data for the common columns only
		// Step 1: Determine the list of common columns which we will hash against between current and previous
		// Frame a map of columns in current & previous
		currentColumns := make(map[string]bool)
		previousColumns := make(map[string]bool)
		hashColumns := make(map[string]bool)
		for i := 0; i < len(currentPanelData.Columns); i++ {
			column := reflect.ValueOf(currentPanelData.Columns[i]).Interface().(map[string]interface{})
			currentColumns[column["name"].(string)] = true
		}
		for i := 0; i < len(previousPanelData.Columns); i++ {
			column := reflect.ValueOf(currentPanelData.Columns[i]).Interface().(map[string]interface{})
			previousColumns[column["name"].(string)] = true
		}
		// Iterate over the current map and check whether its present in the previous map
		// If present, add it to the hash columns
		for currentColumn := range currentColumnMap {
			if _, ok := previousColumns[currentColumn]; ok {
				hashColumns[currentColumn] = true
			}
		}
		// Step 2: Generate the hash for each row in current and previous row data
		currentRowHashes := make(map[string]int)
		for i := 0; i < len(currentPanelData.Rows); i++ {
			row := reflect.ValueOf(currentPanelData.Rows[i]).Interface().(map[string]interface{})
			currentRowHashes[hashRow(row, hashColumns)] = i
		}
		previousRowHashes := make(map[string]int)
		for i := 0; i < len(previousPanelData.Rows); i++ {
			row := reflect.ValueOf(previousPanelData.Rows[i]).Interface().(map[string]interface{})
			previousRowHashes[hashRow(row, hashColumns)] = i
		}
		// Step 3: Loop through the current hashes and check whether they exist in the previous hash or not
		// If they exist in the previous hash, no changes needs to be reported
		// If they dont exist in the previous hash, mark the current row as inserted and add an entry for the previous row as deleted
		currentPath := append(path, "rows")
		for currentHash, currentIndex := range currentRowHashes {
			if _, ok := previousRowHashes[currentHash]; !ok {
				indexString := strconv.Itoa(currentIndex)
				cl.Add(diff.CREATE, append(currentPath, indexString), nil, reflect.ValueOf(currentPanelData.Rows[currentIndex]).Interface().(map[string]interface{}))
			} else {
				delete(previousRowHashes, currentHash)
			}
		}
		// The remaining entries in the previous hash are not present in the current hash and hence need to be marked as deleted
		for _, previousIndex := range previousRowHashes {
			indexString := strconv.Itoa(previousIndex)
			cl.Add(diff.DELETE, append(currentPath, indexString), reflect.ValueOf(previousPanelData.Rows[previousIndex]).Interface().(map[string]interface{}), nil)
		}
	}
}

func GenerateDiff(paths DiffPaths) (map[string]interface{}, error) {
	// Load the previous and current snapshot
	previousSnap, err := loadSnapshot(paths.Previous)
	if err != nil {
		return nil, fmt.Errorf("failed to load previous snapshot: %w", err)
	}
	currentSnap, err := loadSnapshot(paths.Current)
	if err != nil {
		return nil, fmt.Errorf("failed to load current snapshot: %w", err)
	}

	// Initialize the diff snapshot by marshalling the current snapshot to it
	diffSnap := make(map[string]interface{})
	currentSnapBytes, err := json.Marshal(currentSnap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal current snapshot: %w", err)
	}
	err = json.Unmarshal(currentSnapBytes, &diffSnap)

	// Initialize the differ with custom differ for panel data
	differ, err := diff.NewDiffer(diff.CustomValueDiffers(&PanelDataDiffer{}))
	if err != nil {
		return nil, fmt.Errorf("failed to diff snapshots: %w", err)
	}
	differ.SliceOrdering = false
	changeLog, err := differ.Diff(previousSnap, currentSnap)
	if err != nil {
		return nil, fmt.Errorf("failed to diff snapshots: %w", err)
	}

	err = updateDiffSnap(changeLog, diffSnap)

	return diffSnap, nil

}

func loadSnapshot(path string) (*Snapshot, error) {
	var bytes []byte
	var err error

	source := determineSource(path)

	switch source {
	case "url":
		bytes, err = loadSnapshotFromUrl(path)
	case "file":
		bytes, err = os.ReadFile(path)
	case "snapshot":
		bytes, err = loadSnapshotFromJson(path)
	default:
		return nil, fmt.Errorf("expected url, filePath or json, got %v", path)
	}

	var snapshot Snapshot
	err = json.Unmarshal(bytes, &snapshot)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal snapshot: %w", err)
	}

	return &snapshot, nil
}

func determineSource(path string) string {
	u, err := url.ParseRequestURI(path)
	if err == nil && u.Scheme != "" && u.Host != "" {
		return "url"
	}

	if _, err = os.Stat(path); err == nil || os.IsNotExist(err) {
		absPath, pathErr := filepath.Abs(path)
		if pathErr == nil && absPath != "" {
			return "file"
		}
	}

	return "snapshot"
}

func loadSnapshotFromUrl(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL content: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 HTTP status: %d", resp.StatusCode)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return bytes, nil
}

func loadSnapshotFromJson(s string) ([]byte, error) {
	var js json.RawMessage
	err := json.Unmarshal([]byte(s), &js)
	return js, err
}

func updateDiffSnap(changeLog diff.Changelog, diffSnap map[string]interface{}) error {

	for _, change := range changeLog {
		var err error
		topLevel := change.Path[0]
		switch topLevel {
		case "layout":
			switch change.Type {
			case "create":
				err = addKeyValueAtPath(diffSnap, change.Path, "__diff", "inserted")
			case "delete":
				createdValue := change.From.(map[string]interface{})
				createdValue["__diff"] = "deleted"
				lastPath := change.Path[len(change.Path)-1]
				err = addKeyValueAtPath(diffSnap, change.Path[:len(change.Path)-1], lastPath, createdValue)
			case "update":
				err = addKeyValueAtPath(diffSnap, change.Path, "__diff", "updated")
			default:
				continue
			}
			if err != nil {
				return fmt.Errorf("failed to update diff snapshot: %w", err)
			}
		case "panels":
			switch change.Type {
			case "create":
				err = addKeyValueAtPath(diffSnap, change.Path, "__diff", "inserted")
			case "delete":
				createdValue := change.From.(map[string]interface{})
				createdValue["__diff"] = "deleted"
				// Check whether the last path matches the panel key regex
				lastPath := change.Path[len(change.Path)-1]
				matches := PanelKeyRegex.Match([]byte(lastPath))
				if matches {
					err = addValueToSliceAtPath(diffSnap, change.Path[:len(change.Path)-1], createdValue)
				} else {
					// If the last path is an integer, we need to add the value to the slice
					// else proceed with adding the key-value pair
					_, err := strconv.Atoi(lastPath)
					if err != nil {
						err = addKeyValueAtPath(diffSnap, change.Path[:len(change.Path)-1], lastPath, createdValue)
					} else {
						err = addValueToSliceAtPath(diffSnap, change.Path[:len(change.Path)-1], createdValue)
					}
				}
			case "update":
				updatePath := change.Path[:len(change.Path)-1]
				updateKey := fmt.Sprintf("%s_diff", change.Path[len(change.Path)-1])
				err = addKeyValueAtPath(diffSnap, updatePath, "__diff", "updated")
				err = addKeyValueAtPath(diffSnap, updatePath, updateKey, change.From)
			default:
				continue
			}
			if err != nil {
				return fmt.Errorf("failed to update diff snapshot: %w", err)
			}
		default:
			continue
		}
	}
	return nil
}

func addKeyValueAtPath(diffSnap map[string]interface{}, path []string, key string, value interface{}) error {
	var current interface{} = diffSnap
	// traverse path
	for i, p := range path {
		// end o path
		if i == len(path)-1 {
			switch typedCurrent := current.(type) {
			case map[string]interface{}:
				if targetMap, ok := typedCurrent[p].(map[string]interface{}); ok {
					targetMap[key] = value
					return nil
				}
				return nil
			case []interface{}:
				isKeyed := PanelKeyRegex.Match([]byte(p))
				if isKeyed {
					matches := PanelKeyRegex.FindAllStringSubmatch(p, -1)
					keyMap := make(map[string]string)
					for _, match := range matches {
						keyMap[unescape(match[1])] = unescape(match[2])
					}
				outerEnd:
					for _, row := range typedCurrent {
						rowMap := row.(map[string]interface{})
						for k, v := range keyMap {
							switch rowMap[k].(type) {
							case string:
								if rowMap[k] != v {
									continue outerEnd
								}
							case int:
								if fmt.Sprintf("%d", rowMap[k]) != v {
									continue outerEnd
								}
							case float64:
								if fmt.Sprintf("%f", rowMap[k]) != v {
									continue outerEnd
								}
							default:
								if fmt.Sprintf("%v", rowMap[k]) != v {
									continue outerEnd
								}
							}
						}
						rowMap[key] = value
					}
				} else {
					// If the next element is an index, we need to find the element at that index and traverse deeper
					index, err := strconv.Atoi(p)
					if err != nil || index < 0 || index >= len(typedCurrent) {
						return fmt.Errorf("invalid index '%s' at path element '%s'", p, p)
					}
					mapValue := typedCurrent[index].(map[string]interface{})
					mapValue[key] = value
				}
				return nil
			default:
				return fmt.Errorf("expected map or slice at path element '%s', got %T", p, current)
			}
		}

		// traverse deeper
		switch typedCurrent := current.(type) {
		case map[string]interface{}:
			if next, ok := typedCurrent[p]; ok {
				current = next
			} else {
				return fmt.Errorf("path element '%s' not found", p)
			}
		case []interface{}:
			// To traverse deeper into a slice of interface{}, we need to either see whether the path contains an index or a keyed row
			isKeyed := PanelKeyRegex.Match([]byte(p))
			if isKeyed {
				matches := PanelKeyRegex.FindAllStringSubmatch(p, -1)
				keyMap := make(map[string]string)
				for _, match := range matches {
					keyMap[unescape(match[1])] = unescape(match[2])
				}
			outerTraverse:
				for _, row := range typedCurrent {
					rowMap := row.(map[string]interface{})
					for k, v := range keyMap {
						switch rowMap[k].(type) {
						case string:
							if rowMap[k] != v {
								continue outerTraverse
							}
						case int:
							if fmt.Sprintf("%d", rowMap[k]) != v {
								continue outerTraverse
							}
						case float64:
							if fmt.Sprintf("%f", rowMap[k]) != v {
								continue outerTraverse
							}
						default:
							if fmt.Sprintf("%v", rowMap[k]) != v {
								continue outerTraverse
							}
						}
					}
					current = row
					break
				}
			} else {
				// If the next element is an index, we need to find the element at that index and traverse deeper
				index, err := strconv.Atoi(p)
				if err != nil || index < 0 || index >= len(typedCurrent) {
					return fmt.Errorf("invalid index '%s' at path element '%s'", p, p)
				}
				current = typedCurrent[index]
			}
		default:
			return fmt.Errorf("expected map or slice at path element '%s', got %T", p, current)
		}
	}
	return fmt.Errorf("failed to traverse path")
}

func addValueToSliceAtPath(diffSnap map[string]interface{}, path []string, value interface{}) error {
	var current interface{} = diffSnap

	// traverse path
	for i, p := range path {
		// end o path
		if i == len(path)-1 {
			switch typedCurrent := current.(type) {
			case map[string]interface{}:
				if targetSlice, ok := typedCurrent[p].([]interface{}); ok {
					targetSlice = append(targetSlice, value)
					typedCurrent[p] = targetSlice
					return nil
				}
			case []interface{}:
				typedCurrent = append(typedCurrent, value)
				return nil
			default:
				return fmt.Errorf("expected slice at path element '%s', got %T", p, current)
			}
		}

		// traverse deeper
		switch typedCurrent := current.(type) {
		case map[string]interface{}:
			if next, ok := typedCurrent[p]; ok {
				current = next
			} else {
				return fmt.Errorf("path element '%s' not found", p)
			}
		case []interface{}:
			index, err := strconv.Atoi(p)
			if err != nil || index < 0 || index >= len(typedCurrent) {
				return fmt.Errorf("invalid index '%s' at path element '%s'", p, p)
			}
			current = typedCurrent[index]
		default:
			return fmt.Errorf("expected map or slice at path element '%s', got %T", p, current)
		}
	}
	return fmt.Errorf("failed to traverse path")
}

// hashRow generates a hash for the given row, properly handling blob data.
func hashRow(row map[string]interface{}, hashColumns map[string]bool) string {
	// Sort the keys to ensure consistent ordering
	var keys []string
	for k := range row {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Initialize a hash writer
	hasher := sha256.New()

	// Process each key-value pair
	for _, k := range keys {
		// Skip if the key is not in the hash columns
		if _, ok := hashColumns[k]; !ok {
			continue
		}

		value := row[k]

		// Skip if the value for a key is nil
		if helpers.IsNil(value) {
			continue
		}

		// Check if the value is a slice of bytes (blob data)
		if reflect.TypeOf(value).Kind() == reflect.Slice {
			slice, ok := value.([]byte)
			if ok {
				// Write the raw bytes directly to the hasher
				hasher.Write(slice)
				continue
			}
		}

		// For other data types, use fmt.Sprintf to convert them to strings
		hasher.Write([]byte(fmt.Sprintf("%v=%v;", k, value)))
	}

	// Compute the hash
	hashBytes := hasher.Sum(nil)

	// Convert the hash to a hexadecimal string
	return hex.EncodeToString(hashBytes)
}

func escape(input string) string {
	return strings.ReplaceAll(strings.ReplaceAll(input, `:`, `\:`), `;`, `\;`)
}

func unescape(input string) string {
	return strings.ReplaceAll(strings.ReplaceAll(input, `\:`, `:`), `\;`, `;`)
}
