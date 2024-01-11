package dashboardassets

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	filehelpers "github.com/turbot/go-kit/files"
	"github.com/turbot/pipe-fittings/app_specific"
	"github.com/turbot/pipe-fittings/filepaths"
	"github.com/turbot/pipe-fittings/statushooks"
	localcmdconfig "github.com/turbot/powerpipe/internal/cmdconfig"
	"github.com/turbot/steampipe-plugin-sdk/v5/logging"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

var (
	//go:embed *
	staticFS embed.FS
)

const (
	embeddedAssetArchiveName = "assets.tar.gz"
)

func Ensure(ctx context.Context) error {
	logging.LogTime("dashboardassets.Ensure start")
	defer logging.LogTime("dashboardassets.Ensure end")

	// if we are running in development, we don't need to download assets
	// let's just make sure that the assets exist at all (error out if not)
	if localcmdconfig.IsLocal() {
		// verify that the assets exists
		// to verify, read the version file and verify that it has content
		versionFile, err := LoadDashboardAssetVersion()
		if err != nil {
			return sperr.WrapWithMessage(err, "could not load dashboard assets version file")
		}
		if versionFile.Version == "" {
			return sperr.New("during development, dashboard assets must be present when running powerpipe dashboard")
		}
		return nil
	}

	if installedAssetsMatchAppVersion() {
		// nothing to do here
		return nil
	}
	reportAssetsPath := filepaths.EnsureDashboardAssetsDir()

	tarGz, err := staticFS.Open(embeddedAssetArchiveName)
	if err != nil {
		return sperr.WrapWithMessage(err, "could not open embedded dashboard assets archive")
	}
	defer tarGz.Close()

	err = extractTarGz(ctx, tarGz, reportAssetsPath)
	if err != nil {
		return sperr.WrapWithMessage(err, "could not extract embedded dashboard assets archive")
	}
	err = updateAssetVersionFile()
	if err != nil {
		return sperr.WrapWithMessage(err, "could not update dashboard assets version file")
	}

	return nil
}

func updateAssetVersionFile() error {
	versionFile := ReportAssetsVersion{
		Version: app_specific.AppVersion.String(),
	}

	versionFileJSON, err := json.Marshal(versionFile)
	if err != nil {
		return sperr.WrapWithMessage(err, "could not marshal dashboard assets version file")
	}

	versionFilePath := filepaths.ReportAssetsVersionFilePath()
	err = os.WriteFile(versionFilePath, versionFileJSON, 0600)
	if err != nil {
		return sperr.WrapWithMessage(err, "could not write dashboard assets version file")
	}

	return nil
}

func installedAssetsMatchAppVersion() bool {
	versionFile, err := LoadDashboardAssetVersion()
	if err != nil {
		return false
	}

	return versionFile.Version == app_specific.AppVersion.String()
}

type ReportAssetsVersion struct {
	Version string `json:"version"`
}

func LoadDashboardAssetVersion() (*ReportAssetsVersion, error) {
	versionFilePath := filepaths.ReportAssetsVersionFilePath()
	if !filehelpers.FileExists(versionFilePath) {
		return &ReportAssetsVersion{}, nil
	}

	file, _ := os.ReadFile(versionFilePath)
	var versionFile ReportAssetsVersion
	if err := json.Unmarshal(file, &versionFile); err != nil {
		slog.Error("Error while reading dashboard assets version file", "error", err)
		return nil, err
	}

	return &versionFile, nil
}

// extractTarGz extracts a .tar.gz archive to a destination directory.
// this can go into pipe-fittings
// TODO::Binaek - move this to pipe-fittings
func extractTarGz(ctx context.Context, gzipStream io.Reader, dest string) error {
	slog.Info("dashboardassets.extractTarGz start")
	defer slog.Info("dashboardassets.extractTarGz end")

	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}
	uncompressedStream.Close()

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		//nolint:gosec // known archive
		target := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			statushooks.SetStatus(ctx, fmt.Sprintf("Extracting %s…", header.Name))
			outFile, err := os.Create(target)
			if err != nil {
				return err
			}
			//nolint:gosec // known archive
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		default:
			return sperr.New("ExtractTarGz: uknown type: %b in %s", header.Typeflag, header.Name)
		}
	}
}
