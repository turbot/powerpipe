package dashboardserver

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"github.com/turbot/pipe-fittings/filepaths"
)

var logSink io.Writer

const (
	errorPrefix   = "[ Error   ]"
	messagePrefix = "[ Message ]"
	readyPrefix   = "[ Ready   ]"
	waitPrefix    = "[ Wait    ]"
)

func initLogSink() {
	logName := fmt.Sprintf("dashboard-%s.log", time.Now().Format("2006-01-02"))
	logPath := filepath.Join(filepaths.EnsureLogDir(), logName)
	f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("failed to open dashboard manager log file: %s\n", err.Error()) //nolint:forbidigo // TODO: better way to log error?
		os.Exit(3)
	}
	logSink = f
}

func output(_ context.Context, prefix string, msg interface{}) {
	if logSink == nil {
		logSink = os.Stdout
	}
	_, _ = fmt.Fprintf(logSink, "%s %v\n", prefix, msg)
}

func OutputMessage(ctx context.Context, msg string) {
	output(ctx, applyColor(messagePrefix, color.HiGreenString), msg)
}

func OutputWarning(ctx context.Context, msg string) {
	output(ctx, applyColor(messagePrefix, color.RedString), msg)
}

func OutputError(ctx context.Context, err error) {
	output(ctx, applyColor(errorPrefix, color.RedString), err)
}

func outputReady(ctx context.Context, msg string) {
	output(ctx, applyColor(readyPrefix, color.GreenString), msg)
}

func OutputWait(ctx context.Context, msg string) {
	output(ctx, applyColor(waitPrefix, color.CyanString), msg)
}

func applyColor(str string, color func(format string, a ...interface{}) string) string {
	// TODO check streampipe logic is service mode
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		return str
	} else {
		return color((str))
	}
}
