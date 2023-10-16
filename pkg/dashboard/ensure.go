package dashboard

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/turbot/go-kit/files"
	"github.com/turbot/powerpipe/pkg/version"
	"github.com/turbot/steampipe-plugin-sdk/logging"
	"github.com/turbot/steampipe/pkg/filepaths"
	"github.com/turbot/steampipe/pkg/ociinstaller"
	"github.com/turbot/steampipe/pkg/statushooks"
)

func Ensure(ctx context.Context) error {
	logging.LogTime("dashboardassets.Ensure start")
	defer logging.LogTime("dashboardassets.Ensure end")

	// load report assets versions.json
	versionFile, err := loadReportAssetVersionFile()
	if err != nil {
		return err
	}

	if versionFile.Version == version.VersionString {
		return nil
	}

	statushooks.SetStatus(ctx, "Installing dashboard server…")

	reportAssetsPath := filepaths.EnsureDashboardAssetsDir()

	// remove the legacy report folder, if it exists
	if _, err := os.Stat(filepaths.LegacyDashboardAssetsDir()); !os.IsNotExist(err) {
		os.RemoveAll(filepaths.LegacyDashboardAssetsDir())
	}

	return ociinstaller.InstallAssets(ctx, reportAssetsPath)
}

type ReportAssetsVersionFile struct {
	Version string `json:"version"`
}

func loadReportAssetVersionFile() (*ReportAssetsVersionFile, error) {
	versionFilePath := filepaths.ReportAssetsVersionFilePath()
	if !files.FileExists(versionFilePath) {
		return &ReportAssetsVersionFile{}, nil
	}

	file, _ := os.ReadFile(versionFilePath)
	var versionFile ReportAssetsVersionFile
	if err := json.Unmarshal(file, &versionFile); err != nil {
		log.Println("[ERROR]", "Error while reading dashboard assets version file", err)
		return nil, err
	}

	return &versionFile, nil

}
