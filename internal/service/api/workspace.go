package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	pworkspace "github.com/turbot/powerpipe/internal/workspace"
)

// registerWorkspace adds the workspace API
//
//	GET /workspace  the load state of the workspace -> 200 workspace.LoadState
//
// It answers while the workspace reloads, so that a client can wait for a change to its files to load: see
// workspace.LoadState.
func registerWorkspace(router *gin.RouterGroup, w *pworkspace.PowerpipeWorkspace) {
	router.GET("/workspace", func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.JSON(http.StatusOK, w.LoadState())
	})
}
