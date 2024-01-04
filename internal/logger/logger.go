package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/turbot/go-kit/logging"
	"github.com/turbot/pipe-fittings/app_specific"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/filepaths"
)

func SetDefaultLogger() {
	logger := PowerpipeLogger()
	slog.SetDefault(logger)
}

func PowerpipeLogger() *slog.Logger {
	level := getLogLevel()
	if level == constants.LogLevelOff {
		return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	}

	handlerOptions := &slog.HandlerOptions{
		Level: level,

		// TODO KAI SANITIZE
		//ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		//	sanitized := sanitize.Instance.SanitizeKeyValue(a.Key, a.Value.Any())
		//
		//	return slog.Attr{
		//		Key:   a.Key,
		//		Value: slog.AnyValue(sanitized),
		//	}
		//},
	}

	logDestination := logging.NewRotatingLogWriter(filepaths.EnsureLogDir(), app_specific.AppName)
	return slog.New(slog.NewJSONHandler(logDestination, handlerOptions))
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
		return slog.LevelInfo
	}
}
