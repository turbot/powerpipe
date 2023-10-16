package dashboard

import (
	"fmt"
	"os"
	"path/filepath"

	filehelpers "github.com/turbot/go-kit/files"
	"github.com/turbot/steampipe/pkg/constants"
	"github.com/turbot/steampipe/pkg/error_helpers"
)

// Constants for Config
const (
	DefaultInstallDir = "~/.powerpipe"
	versionFileName   = "versions.json"
)

var PowerpipeDir string

func ensureSteampipeSubDir(dirName string) string {
	subDir := steampipeSubDir(dirName)

	if _, err := os.Stat(subDir); os.IsNotExist(err) {
		err = os.MkdirAll(subDir, 0755)
		error_helpers.FailOnErrorWithMessage(err, fmt.Sprintf("could not create %s directory", dirName))
	}

	return subDir
}

func steampipeSubDir(dirName string) string {
	if PowerpipeDir == "" {
		panic(fmt.Errorf("cannot call any Steampipe directory functions before SteampipeDir is set"))
	}
	return filepath.Join(PowerpipeDir, dirName)
}

// EnsureTemplateDir returns the path to the templates directory (creates if missing)
func EnsureTemplateDir() string {
	return ensureSteampipeSubDir(filepath.Join("check", "templates"))
}

// WorkspaceProfileDir returns the path to the workspace profiles directory
// if  STEAMPIPE_WORKSPACE_PROFILES_LOCATION is set use that
// otherwise look in the config folder
// NOTE: unlike other path functions this accepts the install-dir as arg
// this is because of the slightly complex bootstrapping process required because the
// install-dir may be set in the workspace profile
func WorkspaceProfileDir(installDir string) (string, error) {
	if workspaceProfileLocation, ok := os.LookupEnv(constants.EnvWorkspaceProfileLocation); ok {
		return filehelpers.Tildefy(workspaceProfileLocation)
	}
	return filepath.Join(installDir, "config"), nil

}

func ensureDashboardAssetsDir() string {
	return ensureSteampipeSubDir(filepath.Join("dashboard", "assets"))
}

// LegacyDashboardAssetsDir returns the path to the legacy report assets folder
func LegacyDashboardAssetsDir() string {
	return steampipeSubDir("report")
}

// ReportAssetsVersionFilePath returns the report assets version file path
func ReportAssetsVersionFilePath() string {
	return filepath.Join(ensureDashboardAssetsDir(), versionFileName)
}
