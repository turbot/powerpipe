package api

//
//import (
//	"github.com/turbot/pipe-fittings/workspace"
//	"net/http"
//
//	"github.com/gin-gonic/gin"
//	"github.com/turbot/pipe-fittings/error_helpers"
//	"github.com/turbot/pipe-fittings/modinstaller"
//	"github.com/turbot/pipe-fittings/parse"
//	"github.com/turbot/pipe-fittings/versionmap"
//)
//
//type Mod struct {
//	Id          *string `json:"id"`
//	Name        *string `json:"name"`
//	Description *string `json:"description"`
//	Title       *string `json:"title"`
//}
//type GetModResponse struct {
//	Items []Mod `json:"items"`
//}
//
//func (api *APIService) getModHandler(c *gin.Context) {
//	mods := []Mod{}
//	if api.workspace.Mod != nil {
//		mods = []Mod{{
//			Title:       api.workspace.Mod.Title,
//			Description: api.workspace.Mod.Description,
//			Id:          &api.workspace.Mod.FullName,
//			Name:        &api.workspace.Mod.ShortName,
//		}}
//	}
//	c.JSON(http.StatusOK, GetModResponse{Items: mods})
//	return
//}
//
//type InstallModResponse struct {
//	Installed   *versionmap.DependencyVersionMap `json:"installed"`
//	Uninstalled *versionmap.DependencyVersionMap `json:"uninstalled"`
//	Downgraded  *versionmap.DependencyVersionMap `json:"downgraded"`
//	Upgraded    *versionmap.DependencyVersionMap `json:"upgraded"`
//}
//
//type InstallModRequest struct {
//	Names  []string `json:"names"`
//	DryRun *bool    `json:"dry_run"`
//	Force  *bool    `json:"force"`
//}
//
//// TODO all API endpoints which mutates needs locks
//func (api *APIService) installModHandler(c *gin.Context) {
//	input := InstallModRequest{}
//	if err := c.ShouldBindJSON(&input); err != nil {
//		c.AbortWithError(http.StatusBadRequest, err)
//	}
//
//	mod := api.workspace.Mod
//
//	if !parse.ModfileExists(api.workspace.Path) {
//		m, err := workspace.CreateWorkspaceMod(c, api.workspace.Path)
//		if err != nil {
//			error_helpers.FailOnError(err)
//		}
//		mod = m
//	}
//
//	installData, err := modinstaller.InstallWorkspaceDependencies(c.Request.Context(), &modinstaller.InstallOpts{
//		WorkspaceMod: mod,
//		ModArgs:      input.Names,
//		DryRun:       false,
//		Force:        false,
//		GitUrlMode:   "https",
//	})
//	if err != nil {
//		error_helpers.FailOnError(err)
//	}
//
//	response := InstallModResponse{
//		Installed:   &installData.Installed,
//		Uninstalled: &installData.Uninstalled,
//		Downgraded:  &installData.Downgraded,
//		Upgraded:    &installData.Upgraded,
//	}
//
//	c.JSON(http.StatusOK, response)
//}
