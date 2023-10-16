package cmd

import (
	"context"
	"os"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	log "github.com/turbot/powerpipe/internal/logging"
	"github.com/turbot/powerpipe/pkg/utils"
	"github.com/turbot/steampipe/pkg/statushooks"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "powerpipe [--version] [--help] COMMAND [args]",
	Version: "0.0.1",
	Short:   "Powerpipe",
}

var exitCode int

func InitCmd() {
	utils.LogTime("cmd.root.InitCmd start")
	defer utils.LogTime("cmd.root.InitCmd end")

	AddCommands()
	// disable auto completion generation, since we don't want to support
	// powershell yet - and there's no way to disable powershell in the default generator
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func AddCommands() {
	rootCmd.AddCommand(
		modCmd(),
		serviceCmd(),
	)
}

func Execute() int {
	utils.LogTime("cmd.root.Execute start")
	defer utils.LogTime("cmd.root.Execute end")

	ctx := createRootContext()

	rootCmd.ExecuteContext(ctx)
	return exitCode
}

// create the root context - add a status renderer
func createRootContext() context.Context {
	statusRenderer := statushooks.NullHooks
	// if the client is a TTY, inject a status spinner
	if isatty.IsTerminal(os.Stdout.Fd()) {
		statusRenderer = statushooks.NewStatusSpinnerHook()
	}

	ctx := statushooks.AddStatusHooksToContext(context.Background(), statusRenderer)
	ctx = log.ContextWithLogger(ctx)
	return ctx
}
