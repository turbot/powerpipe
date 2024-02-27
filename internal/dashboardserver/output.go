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

var logSinks []io.Writer = []io.Writer{os.Stdout}

const (
	errorPrefix   = "[ Error   ]"
	warningPrefix = "[ Warning   ]"
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
	logSinks = append(logSinks, f)
}

type colorFunc func(format string, a ...interface{}) string

func output(prefix string, color colorFunc, msg any) {
	for _, logSink := range logSinks {
		_, _ = fmt.Fprintf(logSink, "%s %v\n", prefix, applyColor(prefix, color))
	}
}

func OutputMessage(ctx context.Context, msg string) {
	output(messagePrefix, color.HiGreenString, msg)

}

func OutputWarning(ctx context.Context, msg string) {
	output(warningPrefix, color.RedString, msg)
}

func OutputError(ctx context.Context, err error) {
	output(errorPrefix, color.RedString, err)
}

func outputReady(ctx context.Context, msg string) {
	output(readyPrefix, color.GreenString, msg)
}

func OutputWait(ctx context.Context, msg string) {
	output(waitPrefix, color.GreenString, msg)
}

func applyColor(str string, color colorFunc) string {
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		return str
	} else {
		return color(str)
	}
}
