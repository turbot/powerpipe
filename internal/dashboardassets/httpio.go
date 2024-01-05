package dashboardassets

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/viper"
	"github.com/turbot/powerpipe/internal/constants"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

func resolveGithubToken() (string, error) {
	if token, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
		return token, nil
	}
	return "", sperr.New("while powerpipe is in a private repository, a GITHUB_TOKEN is required in environment to download dashboard assets")
}

func downloadFile(filepath string, url string) error {
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

	// Add authorization header to the req
	req.Header.Add("Authorization", "token "+token)

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
	Name   string  `json:"name"`
	Assets []Asset `json:"assets"`
}

// Asset represents a release asset
type Asset struct {
	Url                string `json:"url"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Name               string `json:"name"`
}

func resolveDownloadUrl() (string, error) {
	// this block of code is here to support downloading of assets till this is in a private repository
	// then we go public - just delete it
	target := viper.GetString(constants.ConfigKeyVersion)
	if viper.GetString(constants.ConfigKeyBuiltBy) == constants.LocalBuild {
		target = "latest"
	}

	// get the assets for the target release
	assets, err := getReleaseAssets(target)
	if err != nil {
		return "", sperr.WrapWithMessage(err, "could not fetch release assets")
	}

	// default to latest
	url := "https://github.com/turbot/steampipe/releases/latest/download/dashboard_ui_assets.tar.gz"
	for _, asset := range assets {
		if asset.Name == "dashboard_ui_assets.tar.gz" {
			url = asset.Url
			break
		}
	}
	return url, nil

	// url := fmt.Sprintf("https://github.com/turbot/powerpipe/releases/download/%s/dashboard_ui_assets.tar.gz", viper.GetString(constants.ConfigKeyVersion))
	// if viper.GetString(constants.ConfigKeyBuiltBy) == constants.LocalBuild {
	// 	url = "https://github.com/turbot/steampipe/releases/latest/download/dashboard_ui_assets.tar.gz"
	// }
}

// getReleaseAssets fetches the assets of the latest release
func getReleaseAssets(version string) ([]Asset, error) {
	token, err := resolveGithubToken()
	if err != nil {
		return nil, err
	}

	url := "https://api.github.com/repos/turbot/powerpipe/releases"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "token "+token)
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

	fmt.Println(string(body))

	var releases []*Release
	if err := json.Unmarshal(body, &releases); err != nil {
		return nil, err
	}

	theRelease := releases[0]
	for _, r := range releases {
		// if version is specified, then we need to find the release with that version
		if version != "" {
			if r.Name == version {
				theRelease = r
				break
			}
			continue
		}

		thisSemver, err := semver.NewVersion(r.Name)
		if err != nil {
			return nil, err
		}

		knownReleaseSemver, err := semver.NewVersion(theRelease.Name)
		if err != nil {
			return nil, err
		}

		if thisSemver.GreaterThan(knownReleaseSemver) {
			theRelease = r
		}
	}

	return theRelease.Assets, nil
}
