package dto

import (
	"github.com/turbot/pipe-fittings/versionmap"
)

type Mod struct {
	Id          *string `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Title       *string `json:"title"`
}
type GetModResponse struct {
	Items []Mod `json:"items"`
}

type InstallModResponse struct {
	ModDependencyPath string                           `json:"mod_dependency_path"`
	Installed         *versionmap.DependencyVersionMap `json:"installed"`
	Uninstalled       *versionmap.DependencyVersionMap `json:"uninstalled"`
	Downgraded        *versionmap.DependencyVersionMap `json:"downgraded"`
	Upgraded          *versionmap.DependencyVersionMap `json:"upgraded"`
}

type InstallModRequest struct {
	Names  []string `json:"names"`
	DryRun *bool    `json:"dry_run"`
	Force  *bool    `json:"force"`
}
