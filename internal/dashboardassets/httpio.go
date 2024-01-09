package dashboardassets

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

func resolveGithubToken() (string, error) {
	slog.Info("dashboardassets.resolveGithubToken start")
	defer slog.Info("dashboardassets.resolveGithubToken end")

	if token, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
		return token, nil
	}
	return "", sperr.New("while powerpipe is in a private repository, a GITHUB_TOKEN is required in environment to download dashboard assets")
}

func downloadFile(filepath string, url string) error {
	slog.Info("dashboardassets.downloadFile start")
	defer slog.Info("dashboardassets.downloadFile end")

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	token, err := resolveGithubToken()
	if err != nil {
		return err
	}

	// Add authorization header to the req if we have a token
	if len(token) > 0 {
		req.Header.Add("Authorization", "token "+token)
	}

	// Add accept header to the req - we need to send this otherwise github will just send back the JSON body
	req.Header.Add("Accept", "application/octet-stream")

	// Send the request via a client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return sperr.New("bad status: %s", resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// Release represents a GitHub release
type Release struct {
	Name   string   `json:"name"`
	Assets []*Asset `json:"assets"`
}

func (r *Release) getDashboardAsset() *Asset {
	slog.Info("dashboardassets.Release.getDashboardAsset start")
	defer slog.Info("dashboardassets.Release.getDashboardAsset end")

	for _, asset := range r.Assets {
		if asset.Name == "dashboard_ui_assets.tar.gz" {
			return asset
		}
	}
	return nil
}

// Asset represents a release asset
type Asset struct {
	Url                string `json:"url"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Name               string `json:"name"`
}

func getReleases() ([]*Release, error) {
	slog.Info("dashboardassets.getReleases start")
	defer slog.Info("dashboardassets.getReleases end")

	token, err := resolveGithubToken()
	if err != nil {
		return nil, err
	}

	url := "https://api.github.com/repos/turbot/powerpipe/releases"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add authorization header to the req if we have a token
	if len(token) > 0 {
		req.Header.Add("Authorization", "token "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var releases []*Release
	if err := json.Unmarshal(body, &releases); err != nil {
		return nil, err
	}

	return releases, nil
}
