package dashboard

import (
	"context"
	"encoding/json"
	"github.com/turbot/powerpipe/internal/ociinstaller"
	"log"
	"os"

	"github.com/turbot/go-kit/files"
	"github.com/turbot/powerpipe/internal/version"
	"github.com/turbot/powerpipe/pkg/statushooks"
)

func Ensure(ctx context.Context) error {
	println("dashboardassets.Ensure start")
	defer println("dashboardassets.Ensure end")

	// load report assets versions.json
	versionFile, err := loadReportAssetVersionFile()
	if err != nil {
		return err
	}

	if versionFile.Version == version.VersionString {
		return nil
	}

	statushooks.SetStatus(ctx, "Installing dashboard server…")

	reportAssetsPath := ensureDashboardAssetsDir()

	// remove the legacy report folder, if it exists
	if _, err := os.Stat(LegacyDashboardAssetsDir()); !os.IsNotExist(err) {
		os.RemoveAll(LegacyDashboardAssetsDir())
	}

	return ociinstaller.InstallAssets(ctx, reportAssetsPath)
}

type ReportAssetsVersionFile struct {
	Version string `json:"version"`
}

func loadReportAssetVersionFile() (*ReportAssetsVersionFile, error) {
	versionFilePath := ReportAssetsVersionFilePath()
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
