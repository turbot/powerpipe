package dashboardassets

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

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
	tempDir := ociinstaller.NewTempDir(reportAssetsPath)
	defer func() {
		if err := tempDir.Delete(); err != nil {
			slog.Debug("Failed to delete temp dir after installing assets", "tempDir", tempDir, "error", err)
		}
	}()

	statushooks.SetStatus(ctx, "Downloading dashboard server…")
	assetTarGz, err := downloadAssets(ctx, tempDir.Path)
	if err != nil {
		return sperr.WrapWithMessage(err, "could not ensure dashboard assets")
	}
	tarGz, err := os.Open(assetTarGz)
	if err != nil {
		return sperr.WrapWithMessage(err, "could not open dashboard assets archive")
	}
	defer tarGz.Close()

	statushooks.SetStatus(ctx, "Extracting dashboard server…")
	err = ExtractTarGz(ctx, tarGz, reportAssetsPath)
	return sperr.WrapWithMessage(err, "could not open deflate assets archive")
}

func downloadAssets(ctx context.Context, assetsLocation string) (string, error) {
	// download the blobs
	filePath := filepath.Join(assetsLocation, "assets.tar.gz")
	assetUrl, err := resolveDownloadUrl()
	if err != nil {
		return "", sperr.WrapWithMessage(err, "could not resolve dashboard assets download url")
	}

	if err := downloadFile(filePath, assetUrl); err != nil {
		return "", sperr.WrapWithMessage(err, "could not download dashboard assets")
	}
	return filePath, nil
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

// ExtractTarGz extracts a .tar.gz archive to a destination directory.
// this can go into pipe-fittings
func ExtractTarGz(ctx context.Context, gzipStream io.Reader, dest string) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()
		statushooks.SetStatus(ctx, fmt.Sprintf("Extracting %s…", header.Name))

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		target := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(target)
			if err != nil {
				return err
			}
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
