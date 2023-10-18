package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *APIService) RegisterModApiEndpoints(router *gin.RouterGroup) {
	router.POST("/mod", api.statusNotImplemented)
	router.GET("/mod", api.getModHandler)                          // mod init
	router.POST("/mod/:id/dependency", api.statusNotImplemented)   // ["dependency name", "dependency name"]
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
	mods := []Mod{}
	if api.workspace.Mod != nil {
		mods = []Mod{{
			Title:       api.workspace.Mod.Title,
			Description: api.workspace.Mod.Description,
			Id:          &api.workspace.Mod.FullName,
			Name:        &api.workspace.Mod.ShortName,
		}}
	}
	c.JSON(http.StatusOK, GetModResponse{Items: mods})
	return
}
