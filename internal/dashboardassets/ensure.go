package dashboardassets

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/spf13/viper"
	filehelpers "github.com/turbot/go-kit/files"
	"github.com/turbot/pipe-fittings/app_specific"
	"github.com/turbot/pipe-fittings/filepaths"
	"github.com/turbot/pipe-fittings/ociinstaller"
	"github.com/turbot/pipe-fittings/statushooks"
	"github.com/turbot/steampipe-plugin-sdk/v5/logging"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

func Ensure(ctx context.Context) error {
	logging.LogTime("dashboardassets.Ensure start")
	defer logging.LogTime("dashboardassets.Ensure end")

	// if this is a local build, just check there is an assets folder
	if viper.GetString("main.builtBy") == "local" {
		if filehelpers.DirectoryExists(filepaths.EnsureDashboardAssetsDir()) {
			return nil
		}
		return sperr.New("Dashboard assets not found - run `make dashboard-assets`")
	}
	// load report assets versions.json
	versionFile, err := loadReportAssetVersionFile()
	if err != nil {
		return err
	}

	if versionFile.Version == app_specific.AppVersion.String() {
		return nil
	}

	statushooks.SetStatus(ctx, "Installing dashboard server…")

	reportAssetsPath := filepaths.EnsureDashboardAssetsDir()

	return ociinstaller.InstallAssets(ctx, reportAssetsPath)
}

type ReportAssetsVersionFile struct {
	Version string `json:"version"`
}

func loadReportAssetVersionFile() (*ReportAssetsVersionFile, error) {
	versionFilePath := filepaths.ReportAssetsVersionFilePath()
	if !filehelpers.FileExists(versionFilePath) {
		return &ReportAssetsVersionFile{}, nil
	}

	file, _ := os.ReadFile(versionFilePath)
	var versionFile ReportAssetsVersionFile
	if err := json.Unmarshal(file, &versionFile); err != nil {
		slog.Error("Error while reading dashboard assets version file", "error", err)
		return nil, err
	}

	return &versionFile, nil

}
