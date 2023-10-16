package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterPublicAPI(router *gin.RouterGroup) {
	router.GET("/service", serviceGet)
}

func serviceGet(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "UP",
	})
}
