package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/turbot/pipe-fittings/v2/perr"
	"github.com/turbot/powerpipe/internal/benchmarkrun"
	"github.com/turbot/powerpipe/internal/service/api/common"
)

// BenchmarkRunRequest is the optional body of a request that starts a benchmark run
type BenchmarkRunRequest struct {
	// TimeoutSeconds limits the run, as --benchmark-timeout does. Zero means no limit.
	TimeoutSeconds int `json:"timeout_seconds" binding:"gte=0"`
}

// registerBenchmarkRun adds the benchmark run API. :benchmark_name is the name of the benchmark in any form
// that `powerpipe benchmark run` accepts, e.g. aws_compliance.benchmark.cis_v300. The URL of a run uses the
// full name of its benchmark.
//
//	POST   /benchmark/:benchmark_name/run                   start a run              -> 202 run, Location: the URL of the run; 429 when the runner is full
//	GET    /benchmark/:benchmark_name/run/:run_id           the state of the run     -> 200 run
//	GET    /benchmark/:benchmark_name/run/:run_id/snapshot  the snapshot of the run  -> 200 snapshot (as --export=<file>.pps), 409 while running
//	DELETE /benchmark/:benchmark_name/run/:run_id           cancel and remove a run  -> 204
//
// A run that no client reads for the idle timeout of the runner is cancelled and removed.
func registerBenchmarkRun(router *gin.RouterGroup, runner *benchmarkrun.Runner) {
	router.POST("/benchmark/:benchmark_name/run", func(c *gin.Context) {
		// a web page can send a cross-site POST without a CORS preflight only with a form or text content type,
		// so requiring JSON stops a page in a browser from starting runs on a server it can reach, e.g. on localhost
		if c.ContentType() != binding.MIMEJSON {
			common.AbortWithError(c, perr.BadRequestWithMessage("the request must have Content-Type: application/json"))
			return
		}
		var req BenchmarkRunRequest
		// the body is optional
		if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
			// the request size limiter has already responded to a body that is too large
			if c.IsAborted() {
				return
			}
			// a body that is not valid JSON, or has an unknown field, is a bad request too
			var validationErrors validator.ValidationErrors
			if !errors.As(err, &validationErrors) {
				err = perr.BadRequestWithMessage(err.Error())
			}
			common.AbortWithError(c, err)
			return
		}
		run, err := runner.Start(c.Param("benchmark_name"), time.Duration(req.TimeoutSeconds)*time.Second)
		if err != nil {
			var notFound *benchmarkrun.BenchmarkNotFoundError
			switch {
			case errors.As(err, &notFound):
				common.AbortWithError(c, perr.NotFoundWithMessage(err.Error()))
			case errors.Is(err, benchmarkrun.ErrTooManyRuns):
				common.AbortWithError(c, perr.TooManyRequestsWithMessage(err.Error()))
			default:
				common.AbortWithError(c, perr.InternalWithMessage(err.Error()))
			}
			return
		}
		c.Header("Location", fmt.Sprintf("/api/%s/benchmark/%s/run/%s", c.Param("api_version"), url.PathEscape(run.Benchmark), run.ID))
		c.JSON(http.StatusAccepted, run)
	})

	router.GET("/benchmark/:benchmark_name/run/:run_id", func(c *gin.Context) {
		run, err := runner.Get(c.Param("benchmark_name"), c.Param("run_id"))
		if err != nil {
			abortWithRunError(c, err)
			return
		}
		c.JSON(http.StatusOK, run)
	})

	router.GET("/benchmark/:benchmark_name/run/:run_id/snapshot", func(c *gin.Context) {
		benchmark, id := c.Param("benchmark_name"), c.Param("run_id")
		if _, err := runner.Get(benchmark, id); err != nil {
			abortWithRunError(c, err)
			return
		}
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		if err := runner.WriteSnapshot(c.Request.Context(), benchmark, id, c.Writer); err != nil {
			// before anything is written, report the error. After that, the client sees a broken response.
			if !c.Writer.Written() {
				abortWithRunError(c, err)
				return
			}
			_ = c.Error(err)
		}
	})

	router.DELETE("/benchmark/:benchmark_name/run/:run_id", func(c *gin.Context) {
		if err := runner.Delete(c.Param("benchmark_name"), c.Param("run_id")); err != nil {
			abortWithRunError(c, err)
			return
		}
		c.Status(http.StatusNoContent)
	})
}

func abortWithRunError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, benchmarkrun.ErrNotFound):
		common.AbortWithError(c, perr.NotFoundWithMessage(err.Error()))
	case errors.Is(err, benchmarkrun.ErrNotComplete):
		common.AbortWithError(c, perr.ConflictWithMessage(err.Error()))
	default:
		common.AbortWithError(c, perr.InternalWithMessage(err.Error()))
	}
}
