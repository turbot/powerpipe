package dashboard

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/turbot/pipes-fittings/error_helpers"
)

// Constants for Config
const (
	// TODO KAI MOVE THIS
	DefaultInstallDir = "~/.powerpipe"
	versionFileName   = "versions.json"
)

var PowerpipeDir string

func ensurePowerpipeSubDir(dirName string) string {
	subDir := powerpipeSubDir(dirName)

	if _, err := os.Stat(subDir); os.IsNotExist(err) {
		err = os.MkdirAll(subDir, 0755)
		error_helpers.FailOnErrorWithMessage(err, fmt.Sprintf("could not create %s directory", dirName))
	}

	return subDir
}

func powerpipeSubDir(dirName string) string {
	if PowerpipeDir == "" {
		panic(fmt.Errorf("cannot call any Powerpipt directory functions before PowerpipeDir is set"))
	}
	return filepath.Join(PowerpipeDir, dirName)
}

// EnsureTemplateDir returns the path to the templates directory (creates if missing)
func EnsureTemplateDir() string {
	return ensurePowerpipeSubDir(filepath.Join("check", "templates"))
}

func ensureDashboardAssetsDir() string {
	return ensurePowerpipeSubDir(filepath.Join("dashboard", "assets"))
}

// LegacyDashboardAssetsDir returns the path to the legacy report assets folder
func LegacyDashboardAssetsDir() string {
	return powerpipeSubDir("report")
}

// ReportAssetsVersionFilePath returns the report assets version file path
func ReportAssetsVersionFilePath() string {
	return filepath.Join(ensureDashboardAssetsDir(), versionFileName)
}
