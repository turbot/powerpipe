package cmd

import (
	"context"
	"github.com/turbot/pipe-fittings/modconfig"
	"os"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/app_specific"
	"github.com/turbot/pipe-fittings/cmdconfig"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/statushooks"
	"github.com/turbot/pipe-fittings/utils"
	localconstants "github.com/turbot/powerpipe/internal/constants"
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
		AddPersistentStringFlag(constants.ArgInstallDir, app_specific.DefaultInstallDir, "Path to the installation directory").
		AddPersistentStringFlag(constants.ArgWorkspaceDatabase, app_specific.DefaultWorkspaceDatabase, "Path to the workspace database").
		//// Define the CLI flag parameters for wrapped enum flag.
		//AddPersistentVarFlag(enumflag.New(&outputMode, constants.ArgOutput, types.OutputModeIds, enumflag.EnumCaseInsensitive),
		//	constants.ArgOutput,
		//	"Output format; one of: pretty, plain, yaml, json").
		AddPersistentStringFlag(constants.ArgModLocation, wd, "Path to the mod")

	rootCmd.AddCommand(
		serverCmd(),
		modCmd(),
		resourceCmd[*modconfig.Benchmark](),
		resourceCmd[*modconfig.Control](),
		resourceCmd[*modconfig.Dashboard](),
		resourceCmd[*modconfig.DashboardCard](),
		resourceCmd[*modconfig.DashboardChart](),
		resourceCmd[*modconfig.DashboardContainer](),
		resourceCmd[*modconfig.DashboardFlow](),
		resourceCmd[*modconfig.DashboardGraph](),
		resourceCmd[*modconfig.DashboardHierarchy](),
		resourceCmd[*modconfig.DashboardImage](),
		resourceCmd[*modconfig.DashboardInput](),
		resourceCmd[*modconfig.DashboardTable](),
		resourceCmd[*modconfig.DashboardText](),
		resourceCmd[*modconfig.Query](),
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
