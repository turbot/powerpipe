package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/turbot/pipe-fittings/app_specific"
	"github.com/turbot/pipe-fittings/constants"
)

func FlowipeLogger() *slog.Logger {
	handlerOptions := &slog.HandlerOptions{
		Level: getLogLevel(),
		//ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		//	sanitized := sanitize.Instance.SanitizeKeyValue(a.Key, a.Value.Any())
		//
		//	return slog.Attr{
		//		Key:   a.Key,
		//		Value: slog.AnyValue(sanitized),
		//	}
		//},
	}

	if handlerOptions.Level == constants.LogLevelOff {
		return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	}
	return slog.New(slog.NewJSONHandler(os.Stderr, handlerOptions))
}

func SetDefaultLogger() {
	logger := FlowipeLogger()
	slog.SetDefault(logger)
}

func getLogLevel() slog.Leveler {
	levelEnv := os.Getenv(app_specific.EnvLogLevel)

	switch strings.ToLower(levelEnv) {
	case "trace":
		return constants.LogLevelTrace
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "off":
		return constants.LogLevelOff
	default:
		return constants.LogLevelOff
	}
}
