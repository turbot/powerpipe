package common

// import (
// 	"encoding/base64"

// 	"github.com/gin-gonic/gin"
// 	"github.com/spf13/viper"
// 	"github.com/turbot/flowpipe/internal/types"
// )

// func ListPagingRequest(c *gin.Context) (nextToken string, limit int, err error) {

// 	// Validate and extract paging data from the query string
// 	uri := types.ListRequestQuery{}
// 	if e := c.ShouldBindQuery(&uri); e != nil {
// 		err = e
// 		return
// 	}

// 	// Because limit is optional, and could be zero (which is matched by
// 	// omitempty), specifically catch the zero case here and use the default
// 	// instead.
// 	if uri.Limit == nil {
// 		limit = viper.GetInt("api.list.limit.default")
// 	} else {
// 		limit = *uri.Limit
// 	}

// 	// next_token is a base64 encoded version of the last matched ID. If not
// 	// provided then next_token is "", which means to start at the beginning.
// 	// If the decoding fails, then
// 	if uri.NextToken != "" {
// 		data, e := base64.StdEncoding.DecodeString(uri.NextToken)
// 		if e != nil {
// 			err = e
// 			return
// 		}
// 		nextToken = string(data)
// 	}

// 	return
// }
