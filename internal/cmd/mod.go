package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/app_specific"
	"github.com/turbot/pipe-fittings/backend"
	cmdconfig "github.com/turbot/pipe-fittings/cmdconfig"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/modinstaller"
	"github.com/turbot/pipe-fittings/parse"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/powerpipe/internal/db_client"
)

func modCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "mod [command]",
		Args:  cobra.NoArgs,
		Short: "Powerpipe mod management",
		Long: `Powerpipe mod management.

Mods enable you to run, build, and share dashboards, benchmarks and other resources.

Find pre-built mods in the public registry at https://hub.steampipe.io.

Examples:

    # Create a new mod in the current directory
    powerpipe mod init

    # Install a mod
    powerpipe mod install github.com/turbot/steampipe-mod-aws-compliance
    
    # Update a mod
    powerpipe mod update github.com/turbot/steampipe-mod-aws-compliance
    
    # List installed mods
    powerpipe mod list
    
    # Uninstall a mod
    powerpipe mod uninstall github.com/turbot/steampipe-mod-aws-compliance 
	`,
	}
	cmd.AddCommand(modInstallCmd(),
		modUninstallCmd(),
		modUpdateCmd(),
		modListCmd(),
		showCmd[*modconfig.Mod](),
		modInitCmd(),
	)

	cmd.Flags().BoolP("help", "h", false, "Help for mod")

	return cmd
}

func modInstallCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "install",
		Run:   runModInstallCmd,
		Short: "Install one or more mods and their dependencies",
		Long: `Install one or more mods and their dependencies.

Mods provide an easy way to share Steampipe queries, controls, and benchmarks.
Find mods using the public registry at hub.steampipe.io.

Examples:

  # Install a mod(steampipe-mod-aws-compliance)
  powerpipe mod install github.com/turbot/steampipe-mod-aws-compliance

  # Install a specific version of a mod
  powerpipe mod install github.com/turbot/steampipe-mod-aws-compliance@0.1

  # Install a version of a mod using a semver constraint
  powerpipe mod install github.com/turbot/steampipe-mod-aws-compliance@'^1'

  # Install all mods specified in the mod.pp and their dependencies
  powerpipe mod install

  # Preview what powerpipw mod install will do, without actually installing anything
  powerpipe mod install --dry-run`,
	}

	cmdconfig.OnCmd(cmd).
		AddBoolFlag(constants.ArgDryRun, false, "Show which mods would be installed/updated/uninstalled without modifying them").
		AddStringFlag(constants.ArgDatabase, app_specific.DefaultDatabase, "Turbot Pipes workspace database").
		AddBoolFlag(constants.ArgForce, false, "Install mods even if plugin/cli version requirements are not met (cannot be used with --dry-run)").
		AddBoolFlag(constants.ArgHelp, false, "Help for install", cmdconfig.FlagOptions.WithShortHand("h")).
		AddBoolFlag(constants.ArgPrune, true, "Remove unused dependencies after installation is complete").
		AddModLocationFlag()

	return cmd
}

func runModInstallCmd(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	utils.LogTime("cmd.runModInstallCmd")
	defer func() {
		utils.LogTime("cmd.runModInstallCmd end")
		if r := recover(); r != nil {
			error_helpers.ShowError(ctx, helpers.ToError(r))
			// exitCode = constants.ExitCodeUnknownErrorPanic
		}
	}()

	// try to load the workspace mod definition
	// - if it does not exist, this will return a nil mod and a nil error
	workspacePath := viper.GetString(constants.ArgModLocation)
	workspaceMod, err := parse.LoadModfile(workspacePath)
	error_helpers.FailOnErrorWithMessage(err, "failed to load mod definition")

	// if no mod was loaded, create a default
	if workspaceMod == nil {
		workspaceMod, err = createWorkspaceMod(ctx, cmd, workspacePath)
		if err != nil {
			error_helpers.FailOnError(err)
		}
	}

	// if any mod names were passed as args, convert into formed mod names
	installOpts := modinstaller.NewInstallOpts(workspaceMod, args...)
	installOpts.ModArgs = utils.TrimGitUrls(installOpts.ModArgs)
	installOpts.PluginVersions = getPluginVersions(ctx)

	installData, err := modinstaller.InstallWorkspaceDependencies(ctx, installOpts)
	if err != nil {
		// exitCode = constants.ExitCodeModInstallFailed
		error_helpers.FailOnError(err)
	}

	fmt.Println(modinstaller.BuildInstallSummary(installData)) //nolint:forbidigo // intended output
}

func getPluginVersions(ctx context.Context) *modconfig.PluginVersionMap {
	defaultDatabase, _ := db_client.GetDefaultDatabaseConfig()

	client, err := db_client.NewDbClient(ctx, defaultDatabase)
	if err != nil {
		// do not show warning if --force is set
		if !viper.GetBool(constants.ArgForce) {
			error_helpers.ShowWarning("Could not connect to database - plugin validation will not be performed")
		}
		return nil
	}

	steampipeBackend, isSteampipe := client.Backend.(*backend.SteampipeBackend)
	if !isSteampipe {
		slog.Warn("database does not support Steampipe plugins so plugin requirements cannot be validated", "database", client.Backend.Name())
		return nil
	}
	return modconfig.NewPluginVersionMap(client.Backend.Name(), client.Backend.ConnectionString(), steampipeBackend.PluginVersions)
}

func modUninstallCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "uninstall",
		Run:   runModUninstallCmd,
		Short: "Uninstall a mod and its dependencies",
		Long: `Uninstall a mod and its dependencies.

Example:
  
  # Uninstall a mod
  powerpipe mod uninstall github.com/turbot/steampipe-mod-azure-compliance`,
	}

	cmdconfig.OnCmd(cmd).
		AddBoolFlag(constants.ArgPrune, true, "Remove unused dependencies after uninstallation is complete").
		AddBoolFlag(constants.ArgDryRun, false, "Show which mods would be uninstalled without modifying them").
		AddBoolFlag(constants.ArgHelp, false, "Help for uninstall", cmdconfig.FlagOptions.WithShortHand("h")).
		AddModLocationFlag()

	return cmd
}

func runModUninstallCmd(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	utils.LogTime("cmd.runModInstallCmd")
	defer func() {
		utils.LogTime("cmd.runModInstallCmd end")
		if r := recover(); r != nil {
			error_helpers.ShowError(ctx, helpers.ToError(r))
			exitCode = constants.ExitCodeUnknownErrorPanic
		}
	}()

	// try to load the workspace mod definition
	// - if it does not exist, this will return a nil mod and a nil error
	workspaceMod, err := parse.LoadModfile(viper.GetString(constants.ArgModLocation))
	error_helpers.FailOnErrorWithMessage(err, "failed to load mod definition")
	if workspaceMod == nil {
		//nolint:forbidigo // acceptable output
		fmt.Println("No mods installed.")
		return
	}
	opts := modinstaller.NewInstallOpts(workspaceMod, args...)
	trimGitUrls(opts)
	installData, err := modinstaller.UninstallWorkspaceDependencies(ctx, opts)
	error_helpers.FailOnError(err)
	//nolint:forbidigo // acceptable
	fmt.Println(modinstaller.BuildUninstallSummary(installData))
}

func modUpdateCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "update",
		Run:   runModUpdateCmd,
		Short: "Update one or more mods and their dependencies",
		Long: `Update one or more mods and their dependencies.

Example:

  # Update a mod to the latest version allowed by its current constraint
  powerpipe mod update github.com/turbot/steampipe-mod-aws-compliance

  # Update all mods specified in the mod.pp and their dependencies to the latest versions that meet their constraints, and install any that are missing
  powerpipe mod update`,
	}

	cmdconfig.OnCmd(cmd).
		AddBoolFlag(constants.ArgPrune, true, "Remove unused dependencies after update is complete").
		AddBoolFlag(constants.ArgForce, false, "Update mods even if plugin/cli version requirements are not met (cannot be used with --dry-run)").
		AddBoolFlag(constants.ArgDryRun, false, "Show which mods would be updated without modifying them").
		AddBoolFlag(constants.ArgHelp, false, "Help for update", cmdconfig.FlagOptions.WithShortHand("h")).
		AddModLocationFlag()

	return cmd
}

func runModUpdateCmd(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	utils.LogTime("cmd.runModUpdateCmd")
	defer func() {
		utils.LogTime("cmd.runModUpdateCmd end")
		if r := recover(); r != nil {
			error_helpers.ShowError(ctx, helpers.ToError(r))
			exitCode = constants.ExitCodeUnknownErrorPanic
		}
	}()

	// try to load the workspace mod definition
	// - if it does not exist, this will return a nil mod and a nil error
	workspaceMod, err := parse.LoadModfile(viper.GetString(constants.ArgModLocation))
	error_helpers.FailOnErrorWithMessage(err, "failed to load mod definition")
	if workspaceMod == nil {
		//nolint:forbidigo // acceptable output
		fmt.Println("No mods installed.")
		return
	}

	opts := modinstaller.NewInstallOpts(workspaceMod, args...)
	trimGitUrls(opts)
	installData, err := modinstaller.InstallWorkspaceDependencies(ctx, opts)
	error_helpers.FailOnError(err)

	//nolint:forbidigo // acceptable
	fmt.Println(modinstaller.BuildInstallSummary(installData))
}

// list
func modListCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list",
		Run:   runModListCmd,
		Short: "List currently installed mods",
		Long: `List currently installed mods.
		
Example:

  # List installed mods
  steampipe mod list`,
	}

	cmdconfig.OnCmd(cmd).
		AddBoolFlag(constants.ArgHelp, false, "Help for list", cmdconfig.FlagOptions.WithShortHand("h")).
		AddModLocationFlag()
	return cmd
}

func runModListCmd(cmd *cobra.Command, _ []string) {
	ctx := cmd.Context()
	utils.LogTime("cmd.runModListCmd")
	defer func() {
		utils.LogTime("cmd.runModListCmd end")
		if r := recover(); r != nil {
			error_helpers.ShowError(ctx, helpers.ToError(r))
			exitCode = constants.ExitCodeUnknownErrorPanic
		}
	}()

	// try to load the workspace mod definition
	// - if it does not exist, this will return a nil mod and a nil error
	workspaceMod, err := parse.LoadModfile(viper.GetString(constants.ArgModLocation))
	error_helpers.FailOnErrorWithMessage(err, "failed to load mod definition")
	if workspaceMod == nil {
		//nolint:forbidigo // acceptable
		fmt.Println("No mods installed.")
		return
	}

	opts := modinstaller.NewInstallOpts(workspaceMod)
	installer, err := modinstaller.NewModInstaller(opts)
	error_helpers.FailOnError(err)

	treeString := installer.GetModList()
	if len(strings.Split(treeString, "\n")) > 1 {
		//nolint:forbidigo // acceptable
		fmt.Println()
	}
	//nolint:forbidigo // acceptable
	fmt.Println(treeString)
}

func modInitCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "init",
		Run:   runModInitCmd,
		Short: "Initialize the current directory with a mod.pp file",
		Long: `Initialize the current directory with a mod.pp file.
		
Example:

  # Initialize the current directory with a mod.pp file
  powerpipe mod init`,
	}

	cmdconfig.OnCmd(cmd).
		AddBoolFlag(constants.ArgHelp, false, "Help for init", cmdconfig.FlagOptions.WithShortHand("h")).
		AddModLocationFlag()
	return cmd
}

func runModInitCmd(cmd *cobra.Command, args []string) {
	utils.LogTime("cmd.runModInitCmd")
	ctx := cmd.Context()

	defer func() {
		utils.LogTime("cmd.runModInitCmd end")
		if r := recover(); r != nil {
			error_helpers.ShowError(ctx, helpers.ToError(r))
			exitCode = constants.ExitCodeUnknownErrorPanic
		}
	}()
	workspacePath := viper.GetString(constants.ArgModLocation)
	if _, err := createWorkspaceMod(ctx, cmd, workspacePath); err != nil {
		exitCode = constants.ExitCodeModInitFailed
		error_helpers.FailOnError(err)
	}
}

func createWorkspaceMod(ctx context.Context, cmd *cobra.Command, workspacePath string) (*modconfig.Mod, error) {
	if !modinstaller.ValidateModLocation(ctx, workspacePath) {
		return nil, fmt.Errorf("mod %s cancelled", cmd.Name())
	}

	if _, exists := parse.ModFileExists(workspacePath); exists {
		error_helpers.ShowWarning("Working folder already contains a mod definition file")
		return nil, nil
	}
	mod := modconfig.CreateDefaultMod(workspacePath)
	if err := mod.Save(); err != nil {
		return nil, err
	}
	// only print message for mod init (not for mod install)
	if cmd.Name() == "init" {
		fmt.Printf("Created mod definition file '%s'\n", app_specific.DefaultModFilePath(workspacePath)) //nolint:forbidigo // acceptable
	}

	// load up the written mod file so that we get the updated
	// block ranges
	mod, err := parse.LoadModfile(workspacePath)
	if err != nil {
		return nil, err
	}

	return mod, nil
}

// Modifies(trims) the URL if contains http ot https in arguments
func trimGitUrls(opts *modinstaller.InstallOpts) {
	for i, url := range opts.ModArgs {
		opts.ModArgs[i] = strings.TrimPrefix(url, "http://")
		opts.ModArgs[i] = strings.TrimPrefix(url, "https://")
	}
}
