package cmd

import (
	"github.com/spf13/cobra"

	"github.com/turbot/pipe-fittings/cmdconfig"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/powerpipe/internal/snapshot"
)

func snapshotCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "snapshot [command]",
		Args:  cobra.NoArgs,
		Short: "Powerpipe snapshot management",
		Long: `Powerpipe snapshot management

Examples:

    # Diff two snapshots
	powerpipe snapshot diff https://site.com/old_snapshot.pps /path/to/new_snapshot.pps
`,
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			error_helpers.FailOnError(err)
		},
	}

	cmdconfig.OnCmd(cmd)

	cmd.AddCommand(
		snapshotDiffCmd(),
	)

	return cmd
}

func snapshotDiffCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:  "diff [old-snapshot] [new-snapshot]",
		Run:  runSnapshotDiffCmd,
		Args: cobra.ExactArgs(2),
		// TODO: #command add command descriptions short/long
	}

	return cmd
}

func runSnapshotDiffCmd(cmd *cobra.Command, args []string) {
	bytes, err := snapshot.Diff(snapshot.DiffPaths{
		Previous: args[0],
		Current:  args[1],
	})
	error_helpers.FailOnError(err)

	// TODO: #cli DO THIS PROPERLY!
	cmd.Println(string(bytes))
}
