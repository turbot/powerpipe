package log

import (
	"context"

	"github.com/gin-gonic/gin"
)

type loggerContextKey struct{}

func ContextWithLogger(ctx context.Context) context.Context {
	opts := []LoggerOption{}
	// TODO - how do we get these from the command line or config?
	opts = append(opts,
		WithLevelFromEnvironment(),
		WithFormatFromEnvironment(),
		WithColor(true),
	)
	lgr, err := NewLogger(ctx, opts...)
	if err != nil {
		panic(err)
	}

	// Golang context is a parent-child relationship. When we "add" a value in a context, we actually
	// create a new context with a pointer to the parent context. When we do a ctx.Value() the code traverses the
	// parent-child relationship up. We always pass the context at the bottom of the relationship.
	return context.WithValue(ctx, loggerContextKey{}, lgr)
}

func Logger(ctx context.Context) *PowerpipeLogger {

	// is it a gin context?
	ginContext, ok := ctx.(*gin.Context)
	if !ok {
		return ctx.Value(loggerContextKey{}).(*PowerpipeLogger)
	}

	// if it's a gin context we store the logger in the "fplogger" key
	logger, exists := ginContext.Get("fplogger")
	if !exists {
		return nil
	}
	return logger.(*PowerpipeLogger)
}
