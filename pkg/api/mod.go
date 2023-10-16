package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *APIService) RegisterModApiEndpoints(router *gin.RouterGroup) {
	router.POST("/mod", api.createRootMod)
	router.GET("/mod/:id", api.getRootMod)
	router.POST("/mod/:id/dependency", api.addRootModDependencies)
	router.GET("/mod/:id/dependency", api.getRootModDependencies)
	router.DELETE("/mod/:id/dependency", api.removeRootModDependencies)
}

func (api *APIService) createRootMod(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (api *APIService) getRootMod(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (api *APIService) addRootModDependencies(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (api *APIService) removeRootModDependencies(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (api *APIService) getRootModDependencies(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}
