package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *APIService) RegisterModApiEndpoints(router *gin.RouterGroup) {
	router.POST("/mod", api.statusNotImplemented)
	router.GET("/mod", api.statusNotImplemented)                   // mod init
	router.POST("/mod/:id/dependency", api.statusNotImplemented)   // ["dependency name", "dependency name"]
	router.GET("/mod/:id/dependency", api.statusNotImplemented)    // [{name: "dependency name", dependencies: [{name: "dependency name", dependencies:[]}]}]
	router.DELETE("/mod/:id/dependency", api.statusNotImplemented) // mod uninstall
	router.POST("/mod/:id/command", api.statusNotImplemented)      // {command:"update_dependency | update_dependencies", dependency: "dependency name"}
}
func (api *APIService) statusNotImplemented(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}
