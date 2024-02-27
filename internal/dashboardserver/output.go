package dashboardserver

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"log/slog"
	"os"
)

const (
	errorPrefix   = "[ Error   ]"
	warningPrefix = "[ Warning   ]"
	messagePrefix = "[ Message ]"
	readyPrefix   = "[ Ready   ]"
	waitPrefix    = "[ Wait    ]"
)

func output(_ context.Context, prefix string, msg interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, "%s %v\n", prefix, msg)
}

func OutputMessage(ctx context.Context, msg string) {
	output(ctx, applyColor(messagePrefix, color.HiGreenString), msg)
	slog.Info(msg)
}

func OutputWarning(ctx context.Context, msg string) {
	output(ctx, applyColor(messagePrefix, color.RedString), msg)
	slog.Warn(msg)
}

func OutputError(ctx context.Context, err error) {
	output(ctx, applyColor(errorPrefix, color.RedString), err)
	slog.Error("Error", "error", err)
}

func OutputReady(ctx context.Context, msg string) {
	output(ctx, applyColor(readyPrefix, color.GreenString), msg)
	slog.Info(msg)
}

func OutputWait(ctx context.Context, msg string) {
	output(ctx, applyColor(waitPrefix, color.CyanString), msg)
	slog.Info(msg)
}

func applyColor(str string, color func(format string, a ...interface{}) string) string {
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		return str
	} else {
		return color(str)
	}
}
