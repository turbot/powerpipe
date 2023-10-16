package log

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	//nolint:depguard // Wrapper for Zap
	"go.uber.org/zap"

	//nolint:depguard // Wrapper for Zap
	"go.uber.org/zap/zapcore"

	"github.com/spf13/viper"
	"github.com/turbot/powerpipe/internal/sanitize"
	"github.com/turbot/powerpipe/pkg/constants"
)

type PowerpipeLogger struct {

	// Level is the logging level to use for output
	Level zapcore.Level

	// Special handling for the "trace" level
	TraceLevel string

	// Format is the logging format to use for output: json or console
	Format string

	// Color is whether to use color in the console output
	Color bool

	// Zap is the Zap logger instance
	Zap   *zap.Logger
	Sugar *zap.SugaredLogger
}

// LoggerOption defines a type of function to configures the Logger.
type LoggerOption func(*PowerpipeLogger) error

// NewLogger creates a new Logger.
func NewLogger(ctx context.Context, opts ...LoggerOption) (*PowerpipeLogger, error) {
	// Defaults
	c := &PowerpipeLogger{
		Level:  zapcore.InfoLevel,
		Format: "console",
	}
	// Set options
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return c, err
		}
	}

	err := c.Initialize()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func WithColor(enabled bool) LoggerOption {
	return func(c *PowerpipeLogger) error {
		c.Color = enabled
		return nil
	}
}

func WithLevelFromEnvironment() LoggerOption {
	return func(c *PowerpipeLogger) error {
		traceLevelStr := strings.ToLower(os.Getenv("POWERPIPE_TRACE_LEVEL"))
		if traceLevelStr != "" {
			var err error
			logLevel, err := zapcore.ParseLevel(traceLevelStr)
			if err == nil {
				c.Level = logLevel
				c.TraceLevel = traceLevelStr
			}
		} else {
			// Get the desired logging level from the POWERPIPE_LOG_LEVEL environment variable
			logLevelStr := strings.ToLower(os.Getenv("POWERPIPE_LOG_LEVEL"))
			// If the POWERPIPE_LOG_LEVEL environment variable is set, parse its value to determine the logging level
			if logLevelStr != "" {
				var err error
				logLevel, err := zapcore.ParseLevel(logLevelStr)
				if err == nil {
					c.Level = logLevel
				}
			}
		}
		return nil
	}
}

func WithFormatFromEnvironment() LoggerOption {
	return func(c *PowerpipeLogger) error {
		// Get the desired logging format from the POWERPIPE_LOG_FORMAT environment variable
		logFormat := strings.ToLower(os.Getenv("POWERPIPE_LOG_FORMAT"))
		switch logFormat {
		case "json", "console":
			c.Format = logFormat
		}
		return nil
	}
}

func (c *PowerpipeLogger) Initialize() error {

	// Configure the logging output
	var encoder zapcore.Encoder
	if c.Format == "json" {
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	} else {
		ec := zap.NewDevelopmentEncoderConfig()
		if c.Color {
			ec.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}
		encoder = zapcore.NewConsoleEncoder(ec)
	}
	consoleDebugging := zapcore.Lock(os.Stdout)
	//consoleErrors := zapcore.Lock(os.Stderr)

	// Configure the logging level based on the POWERPIPE_LOG_LEVEL environment variable
	atomicLevel := zap.NewAtomicLevelAt(c.Level)

	// Create the Zap logger instance
	core := zapcore.NewTee(
		//zapcore.NewCore(encoder, consoleErrors, atomicLevel),
		zapcore.NewCore(encoder, consoleDebugging, atomicLevel),
	)

	c.Zap = zap.New(core)
	c.Sugar = c.Zap.Sugar()

	// Do not Desugar() it's expensive (according to Zap themselves)
	// Zap suggested that we have 2 logger instances
	_, err := zap.RedirectStdLogAt(c.Zap, zapcore.DebugLevel)
	if err != nil {
		return err
	}

	return nil
}

func (c *PowerpipeLogger) Sync() error {
	return c.Zap.Sync()
}

func (c *PowerpipeLogger) Error(msg string, keysAndValues ...interface{}) {
	sanitizedKeysAndValues := sanitize.SanitizeLogEntries(keysAndValues)
	c.Sugar.Errorw(msg, sanitizedKeysAndValues...)
}

func (c *PowerpipeLogger) Warn(msg string, keysAndValues ...interface{}) {
	sanitizedKeysAndValues := sanitize.SanitizeLogEntries(keysAndValues)
	c.Sugar.Warnw(msg, sanitizedKeysAndValues...)
}

func (c *PowerpipeLogger) Info(msg string, keysAndValues ...interface{}) {
	sanitizedKeysAndValues := sanitize.SanitizeLogEntries(keysAndValues)
	c.Sugar.Infow(msg, sanitizedKeysAndValues...)
}

func (c *PowerpipeLogger) Debug(msg string, keysAndValues ...interface{}) {
	sanitizedKeysAndValues := sanitize.SanitizeLogEntries(keysAndValues)
	c.Sugar.Debugw(msg, sanitizedKeysAndValues...)
}

func (c *PowerpipeLogger) Trace(msg string, keysAndValues ...interface{}) {
	if c.TraceLevel != "" {
		sanitizedKeysAndValues := sanitize.SanitizeLogEntries(keysAndValues)
		msg = "**** " + msg
		switch c.TraceLevel {
		case "debug":
			c.Sugar.Debugw(msg, sanitizedKeysAndValues...)
		case "info":
			c.Sugar.Infow(msg, sanitizedKeysAndValues...)
		case "warn":
			c.Sugar.Warnw(msg, sanitizedKeysAndValues...)
		}
	}
}

func ExecutionLogger(ctx context.Context, executionID string) *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	cfg.Sampling = nil
	cfg.OutputPaths = []string{path.Join(viper.GetString(constants.ArgLogDir), fmt.Sprintf("%s.jsonl", executionID))}
	return zap.Must(cfg.Build())
}
