package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/modinstaller"
	"github.com/turbot/pipe-fittings/versionmap"
)

func (api *APIService) RegisterModApiEndpoints(router *gin.RouterGroup) {
	router.POST("/mod", api.statusNotImplemented)
	router.GET("/mod", api.getModHandler)                          // mod init
	router.POST("/mod/:id/dependency", api.installModHandler)      // ["dependency name", "dependency name"]
	router.GET("/mod/:id/dependency", api.statusNotImplemented)    // [{name: "dependency name", dependencies: [{name: "dependency name", dependencies:[]}]}]
	router.DELETE("/mod/:id/dependency", api.statusNotImplemented) // mod uninstall
	router.POST("/mod/:id/command", api.statusNotImplemented)      // {command:"update_dependency | update_dependencies", dependency: "dependency name"}
}
func (api *APIService) statusNotImplemented(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

type Mod struct {
	Id          *string `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Title       *string `json:"title"`
}
type GetModResponse struct {
	Items []Mod `json:"items"`
}

func (api *APIService) getModHandler(c *gin.Context) {
	var mods []Mod
	if api.workspace.Mod != nil {
		mods = []Mod{{
			Title:       api.workspace.Mod.Title,
			Description: api.workspace.Mod.Description,
			Id:          &api.workspace.Mod.FullName,
			Name:        &api.workspace.Mod.ShortName,
		}}
	}
	c.JSON(http.StatusOK, GetModResponse{Items: mods})

}

type InstallModResponse struct {
	Installed   *versionmap.DependencyVersionMap `json:"installed"`
	Uninstalled *versionmap.DependencyVersionMap `json:"uninstalled"`
	Downgraded  *versionmap.DependencyVersionMap `json:"downgraded"`
	Upgraded    *versionmap.DependencyVersionMap `json:"upgraded"`
}

type InstallModRequest struct {
	Names  []string `json:"names"`
	DryRun *bool    `json:"dry_run"`
	Force  *bool    `json:"force"`
}

// TODO all API endpoints which mutates needs locks
func (api *APIService) installModHandler(c *gin.Context) {
	input := InstallModRequest{}
	if err := c.ShouldBindJSON(&input); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
	}

	mod := api.workspace.Mod

	installData, err := modinstaller.InstallWorkspaceDependencies(c.Request.Context(), &modinstaller.InstallOpts{
		WorkspaceMod: mod,
		ModArgs:      input.Names,
		DryRun:       false,
		Force:        false,
	})
	if err != nil {
		error_helpers.FailOnError(err)
	}

	response := InstallModResponse{
		Installed:   &installData.Installed,
		Uninstalled: &installData.Uninstalled,
		Downgraded:  &installData.Downgraded,
		Upgraded:    &installData.Upgraded,
	}

	c.JSON(http.StatusOK, response)
}
