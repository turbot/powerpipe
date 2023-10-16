package common

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/turbot/powerpipe/pkg/types"
)

const (
	// APIVersion is the current version of the API
	APIVersion0 = "v0"
	//APIVersion1      = "v1"
	APIVersionLatest = "latest"
)

func apiPath(version string) string {
	return fmt.Sprintf("/api/%s", version)
}

// API Prefix is the API path prefix without a version parameter
func APIPrefix() string {
	return "/api/:api_version"
}

// PathPrefix is the API path prefix for the current version
//func PathPrefix() string {
//	return apiPath(APIVersion0)
//}

// PathPrefixWithVersion is the API path prefix for the requested version
func PathPrefixWithVersion(version string) string {
	return apiPath(version)
}

func ValidateAPIVersion(c *gin.Context) {
	// Parse org handle from URI
	var uri types.APIVersionRequestURI
	if err := c.ShouldBindUri(&uri); err != nil {
		AbortWithError(c, err)
		return
	}
	c.Next()
}
