package snapshot

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/r3labs/diff/v3"
)

type DiffPaths struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
}

func Diff(paths DiffPaths) (map[string]interface{}, error) {
	previousSnap, err := loadSnapshot(paths.Previous)
	if err != nil {
		return nil, fmt.Errorf("failed to load previous snapshot: %w", err)
	}
	currentSnap, err := loadSnapshot(paths.Current)
	if err != nil {
		return nil, fmt.Errorf("failed to load current snapshot: %w", err)
	}

	diffSnap := maps.Clone(currentSnap)

	changeLog, err := diff.Diff(previousSnap, currentSnap)
	if err != nil {
		return nil, fmt.Errorf("failed to diff snapshots: %w", err)
	}

	err = updateDiffSnap(changeLog, &diffSnap)

	return diffSnap, nil

}

func loadSnapshot(path string) (map[string]interface{}, error) {
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

	var snapshot map[string]interface{}
	err = json.Unmarshal(bytes, &snapshot)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal snapshot: %w", err)
	}

	return snapshot, nil
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

func updateDiffSnap(changeLog diff.Changelog, diffSnap *map[string]interface{}) error {

	for _, change := range changeLog {
		var err error
		topLevel := change.Path[0]
		switch topLevel {
		case "layout":
			switch change.Type {
			case "create":
				err = addKeyValueAtPath(*diffSnap, change.Path, "__diff", "inserted")
			case "delete":
				err = addKeyValueAtPath(*diffSnap, change.Path, "__diff", "deleted")
				// TODO: #snapshot we need to extract the deleted layout item and inject into the diff
			case "update":
				err = addKeyValueAtPath(*diffSnap, change.Path, "__diff", "updated")
			default:
				continue
			}
			if err != nil {
				return fmt.Errorf("failed to update diff snapshot: %w", err)
			}
		case "panels":
			switch change.Type {
			case "create":
				err = addKeyValueAtPath(*diffSnap, change.Path, "__diff", "inserted")
			case "delete":
				fmt.Println("delete")
			case "update":
				updatePath := change.Path[:len(change.Path)-1]
				updateKey := fmt.Sprintf("%s_diff", change.Path[len(change.Path)-1])
				err = addKeyValueAtPath(*diffSnap, updatePath, "__diff", "updated")
				err = addKeyValueAtPath(*diffSnap, updatePath, updateKey, change.From)
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
				typedCurrent[key] = value
				return nil
			case []interface{}:
				index, err := strconv.Atoi(p)
				if err != nil || index < 0 || index >= len(typedCurrent) {
					return fmt.Errorf("invalid index at path element '%s'", p)
				}

				if targetMap, ok := typedCurrent[index].(map[string]interface{}); ok {
					targetMap[key] = value
					return nil
				}

				return fmt.Errorf("expected map at index %d, got %T", index, typedCurrent[index])
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
