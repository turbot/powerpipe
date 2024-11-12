package snapshot

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type DiffPaths struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
}

func Diff(paths DiffPaths) ([]byte, error) {
	// Validate paths
	if err := validateDiffPath(paths.Previous); err != nil {
		return nil, fmt.Errorf("invalid previous path: %s", err)
	}
	if err := validateDiffPath(paths.Current); err != nil {
		return nil, fmt.Errorf("invalid current path: %s", err)
	}

	previousSnap, err := loadSnapshot(paths.Previous)
	if err != nil {
		return nil, fmt.Errorf("failed to load previous snapshot: %w", err)
	}
	currentSnap, err := loadSnapshot(paths.Current)
	if err != nil {
		return nil, fmt.Errorf("failed to load current snapshot: %w", err)
	}

	diffSnap := maps.Clone(currentSnap)

	slog.Debug("previousSnap", "previousSnap", previousSnap)
	slog.Debug("currentSnap", "currentSnap", currentSnap)
	slog.Debug("diffSnap", "diffSnap", diffSnap)

	// TODO: Create New SnapshotDiff struct
	// TODO: Iterate Panels, Compare, & Update SnapshotDiff
	// TODO: Marshal SnapshotDiff to JSON and return

	out, err := json.Marshal(diffSnap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal diff snapshot: %w", err)
	}
	return out, nil

}

func validateDiffPath(path string) error {
	// check if parses as URL
	_, err := url.ParseRequestURI(path)
	if err == nil {
		return nil
	}

	// check if valid file path && exists
	if _, err = os.Stat(path); err == nil || os.IsNotExist(err) {
		absPath, pathErr := filepath.Abs(path)
		if pathErr == nil && absPath != "" {
			return nil
		}
	}

	return err
}

func loadSnapshot(path string) (map[string]interface{}, error) {
	var bytes []byte
	var err error
	var u *url.URL
	// Check if the path is a URL
	u, err = url.ParseRequestURI(path)
	if err == nil && u.Scheme != "" && u.Host != "" {
		// Load content from URL
		resp, httpErr := http.Get(path)
		if httpErr != nil {
			return nil, fmt.Errorf("failed to fetch URL content: %w", httpErr)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("received non-200 HTTP status: %d", resp.StatusCode)
		}

		bytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
	} else {
		// Load content from file
		bytes, err = os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read file content: %w", err)
		}
	}

	var snapshot map[string]interface{}
	err = json.Unmarshal(bytes, &snapshot)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal snapshot: %w", err)
	}

	return snapshot, nil
}
