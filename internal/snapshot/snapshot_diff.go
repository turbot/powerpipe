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
	return nil
}
