package ociinstaller

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	shared_ociinstaller "github.com/turbot/pipe-fittings/ociinstaller"
	"github.com/turbot/steampipe/pkg/constants"
)

// InstallAssets installs the Steampipe report server assets
func InstallAssets(ctx context.Context, assetsLocation string) error {
	tempDir := shared_ociinstaller.NewTempDir(assetsLocation)
	defer func() {
		if err := tempDir.Delete(); err != nil {
			log.Printf("[TRACE] Failed to delete temp dir '%s' after installing assets: %s", tempDir, err)
		}
	}()

	// download the blobs
	imageDownloader := shared_ociinstaller.NewOciDownloader()
	image, err := imageDownloader.Download(ctx, shared_ociinstaller.NewSteampipeImageRef(constants.DashboardAssetsImageRef), shared_ociinstaller.ImageTypeAssets, tempDir.Path)
	if err != nil {
		return err
	}

	// install the files
	if err = installAssetsFiles(image, tempDir.Path, assetsLocation); err != nil {
		return err
	}

	return nil
}

func installAssetsFiles(image *shared_ociinstaller.SteampipeImage, tempdir string, destination string) error {
	fileName := image.Assets.ReportUI
	sourcePath := filepath.Join(tempdir, fileName)
	if err := shared_ociinstaller.MoveFolderWithinPartition(sourcePath, destination); err != nil {
		return fmt.Errorf("could not install %s to %s", sourcePath, destination)
	}
	return nil
}
