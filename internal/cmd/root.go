package cmd

import (
	"context"
	"os"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/app_specific"
	"github.com/turbot/pipe-fittings/cmdconfig"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/statushooks"
	"github.com/turbot/pipe-fittings/utils"
	localconstants "github.com/turbot/powerpipe/internal/constants"
	"github.com/turbot/powerpipe/internal/resources"
)

var exitCode int

// Build the cobra command that handles our command line tool.
func rootCommand() *cobra.Command {
	// Define our command
	rootCmd := &cobra.Command{
		Use:     "powerpipe [--version] [--help] COMMAND [args]",
		Short:   localconstants.PowerpipeShortDescription,
		Long:    localconstants.PowerpipeLongDescription,
		Version: viper.GetString("main.version"),
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			error_helpers.FailOnError(err)
		},
	}

	utils.LogTime("cmd.root.InitCmd start")
	defer utils.LogTime("cmd.root.InitCmd end")

	rootCmd.SetVersionTemplate("Powerpipe v{{.Version}}\n")

	// set the current working directory
	wd, err := os.Getwd()
	error_helpers.FailOnError(err)

	cmdconfig.
		OnCmd(rootCmd).
		AddPersistentStringFlag(constants.ArgConfigPath, "", "Colon separated list of paths to search for workspace files, in order of decreasing precedence").
		AddPersistentStringFlag(constants.ArgInstallDir, app_specific.DefaultInstallDir, "Path to the installation directory").
		AddPersistentStringFlag(constants.ArgModLocation, wd, "Path to the workspace working directory").
		AddPersistentStringFlag(constants.ArgWorkspaceProfile, "default", "Sets the Powerpipe workspace profile")

	rootCmd.AddCommand(
		serverCmd(),
		modCmd(),
		loginCmd(),
		resourceCmd[*resources.Benchmark](),
		resourceCmd[*resources.DetectionBenchmark](),
		resourceCmd[*resources.Control](),
		resourceCmd[*resources.Dashboard](),
		resourceCmd[*resources.Query](),
		resourceCmd[*modconfig.Variable](),
	)

	// disable auto completion generation, since we don't want to support
	// powershell yet - and there's no way to disable powershell in the default generator
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd
}

func Execute() int {
	rootCmd := rootCommand()
	utils.LogTime("cmd.root.Execute start")
	defer utils.LogTime("cmd.root.Execute end")

	ctx := createRootContext()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		exitCode = -1
	}
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
	return ctx
}
