package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/gin-contrib/gzip"
	size "github.com/gin-contrib/size"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/workspace"
	"github.com/turbot/powerpipe/internal/dashboard"
	"github.com/turbot/powerpipe/internal/service/api/common"
	"gopkg.in/olahol/melody.v1"
)

// @title powerpipe
// @version 0.1.0
// @description Powerpipe is a ...
// @contact.name Support
// @contact.email help@powerpipe.io

// @contact.name   powerpipe
// @contact.url    http://www.powerpipe.io
// @contact.email  info@powerpipe.io

// @license.name  AGPLv3
// @license.url   https://www.gnu.org/licenses/agpl-3.0.en.html

// @host localhost
// @schemes https
// @BasePath /api/v0
// @query.collection.format multi

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// APIService represents the API service.
type APIService struct {
	// Ctx is the context used by the API service.
	ctx context.Context

	httpServer  *http.Server
	httpsServer *http.Server

	HTTPPort string `json:"http_port,omitempty"`

	HTTPSHost string `json:"https_host,omitempty"`
	HTTPSPort string `json:"https_port,omitempty"`

	// Status tracking for the API service.
	Status    string     `json:"status"`
	StartedAt *time.Time `json:"started_at,omitempty"`
	StoppedAt *time.Time `json:"stopped_at,omitempty"`

	apiPrefixGroup *gin.RouterGroup
	router         *gin.Engine
	webSocket      *melody.Melody

	// the loaded workspace
	workspace *workspace.Workspace
}

// APIServiceOption defines a type of function to configures the APIService.
type APIServiceOption func(*APIService) error

func WithWebSocket(webSocket *melody.Melody) APIServiceOption {
	return func(api *APIService) error {
		api.webSocket = webSocket
		return nil
	}
}

func WithWorkspace(workspace *workspace.Workspace) APIServiceOption {
	return func(api *APIService) error {
		api.workspace = workspace
		return nil
	}
}

// NewAPIService creates a new APIService.
func NewAPIService(ctx context.Context, opts ...APIServiceOption) (*APIService, error) {
	// Defaults
	api := &APIService{
		ctx:      ctx,
		Status:   "initialized",
		HTTPPort: fmt.Sprintf("%d", 9194),
	}

	// Set options
	for _, opt := range opts {
		err := opt(api)
		if err != nil {
			return api, err
		}
	}
	return api, nil
}

// Start starts services managed by the Manager.
func (api *APIService) Start() error {
	// if api.mod == nil {
	// 	return sperr.New("cannot start without a defined mod")
	// }
	// Set the gin mode based on our environment, to configure logging etc as appropriate
	gin.SetMode(viper.GetString("environment"))
	binding.EnableDecoderDisallowUnknownFields = true

	// Initialize gin
	router := gin.New()

	apiPrefixGroup := router.Group(common.APIPrefix())
	apiPrefixGroup.Use(common.ValidateAPIVersion)

	// Limit the size of POST requests
	// There doesn't seem a way to set the request size per path, but for now we have
	// no requirement for different limits on different paths. So just set one limit
	// for all request (for now)
	router.Use(size.RequestSizeLimiter(viper.GetInt64("web.request.size_limit")))

	// Create compression middleware - exclude process logs as we handle compression within the API itself
	compressionMiddleware := gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPathsRegexs([]string{"^/api/.+/.*[avatar|\\.jsonl]$"}))
	apiPrefixGroup.Use(compressionMiddleware)
	router.Use(compressionMiddleware)

	// Simple rate limiting:
	// * In memory only, so will not check across web servers
	// * Burst is the initial credits, with fill being added per second (to max of burst)
	//
	// Other option: ulele/limiter
	//
	// In the end decided to use tollbooth even though it doesn't have Redis support because that what was used in SPC
	// so I don't have to learn a new library.
	//
	// ulele/limiter support Redis AND in memory, so we may want to switch to that when we have more functionality in flowpipe
	//
	apiLimiter := tollbooth.NewLimiter(viper.GetFloat64("web.rate.fill"), &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	apiLimiter.SetBurst(viper.GetInt("web.rate.burst"))

	RegisterPublicAPI(apiPrefixGroup)
	api.RegisterModApiEndpoints(apiPrefixGroup)

	// put in handing for the dashboard for the mod
	assetsDirectory := dashboard.GetStaticAssetsDir()
	// respond with the static dashboard assets for / (root)
	router.Use(static.Serve("/", static.LocalFile(assetsDirectory, true)))
	if api.webSocket != nil {
		router.GET("/ws", func(c *gin.Context) {
			if err := api.webSocket.HandleRequest(c.Writer, c.Request); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		})
	}

	// fall through
	router.NoRoute(func(c *gin.Context) {
		// https://stackoverflow.com/questions/49547/how-do-we-control-web-page-caching-across-all-browsers
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
		c.Header("Pragma", "no-cache")                                   // HTTP 1.0.
		c.Header("Expires", "0")                                         // Proxies.
		c.File(path.Join(assetsDirectory, "index.html"))
	})

	api.apiPrefixGroup = apiPrefixGroup
	api.router = router

	// Custom validators for our types
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Return the JSON fieldname in the Tag() field for errors.
		// See https://github.com/go-playground/validator/issues/287
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		// Custom validators using struct field tags
		_ = v.RegisterValidation("flowpipe_api_version", common.APIVersionValidator())
	}

	// Single Page App must catch all routes that are not found, it
	// handles them in a client side router.
	// router.NoRoute(func(c *gin.Context) {
	// 	path := c.Request.URL.Path
	// 	method := c.Request.Method
	// 	if strings.HasPrefix(path, "/api") {
	// 		c.JSON(http.StatusNotFound, gin.H{"error": perr.NotFoundWithMessage(fmt.Sprintf("API Not Found: %s %s.", method, path))})
	// 	} else {
	// 		c.File("./static/index.html")
	// 	}
	// })

	// Server setup with graceful shutdown
	api.httpServer = &http.Server{
		Addr:              fmt.Sprintf("%s:%s", api.HTTPSHost, api.HTTPPort),
		Handler:           router,
		ReadHeaderTimeout: 60 * time.Second,
	}

	api.httpsServer = &http.Server{
		Addr:              fmt.Sprintf("%s:%s", api.HTTPSHost, api.HTTPSPort),
		Handler:           router,
		ReadHeaderTimeout: 60 * time.Second,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	// go func() {
	// 	if err := api.httpsServer.ListenAndServeTLS("./service/certificate/server.crt", "./service/certificate/server.key"); err != nil && err != http.ErrServerClosed {
	// 		log.Fatalf("listen: %s\n", err)
	// 	}
	// }()

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := api.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// api.StartedAt = utils.TimeNow()
	api.Status = "running"

	return nil
}
