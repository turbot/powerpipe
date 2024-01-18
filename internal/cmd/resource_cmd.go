package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/turbot/pipe-fittings/cmdconfig"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	localcmdconfig "github.com/turbot/powerpipe/internal/cmdconfig"
	"github.com/turbot/powerpipe/internal/display"
)

func resourceCmd[T modconfig.HclResource]() *cobra.Command {
	typeName := localcmdconfig.GetGenericTypeName[T]()

	var cmd = &cobra.Command{
		Use:   fmt.Sprintf("%s [command]", typeName),
		Args:  cobra.NoArgs,
		Short: resourceCommandShortDescription(typeName),
		Long:  resourceCommandLongDescription(typeName),
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			error_helpers.FailOnError(err)
		},
	}

	cmdconfig.OnCmd(cmd)

	cmd.AddCommand(getResourceCommands[T]()...)
	return cmd
}

func listCmd[T modconfig.HclResource]() *cobra.Command {
	typeName := localcmdconfig.GetGenericTypeName[T]()
	var cmd = &cobra.Command{
		Use:   "list",
		Args:  cobra.NoArgs,
		Run:   func(cmd *cobra.Command, args []string) { display.ListResources[T](cmd) },
		Short: listCommandShortDescription(typeName),
		Long:  listCommandLongDescription(typeName)}
	// initialize hooks
	cmdconfig.OnCmd(cmd)

	return cmd
}

func showCmd[T modconfig.HclResource]() *cobra.Command {
	typeName := localcmdconfig.GetGenericTypeName[T]()

	var cmd = &cobra.Command{
		Use:   showCommandUse(typeName),
		Args:  cobra.ExactArgs(1),
		Run:   func(cmd *cobra.Command, args []string) { display.ShowResource[T](cmd, args) },
		Short: showCommandShortDescription(typeName),
		Long:  showCommandLongDescription(typeName),
	}
	// initialize hooks
	cmdconfig.OnCmd(cmd)

	return cmd
}

// determine which resource commands apply to this resource
func getResourceCommands[T modconfig.HclResource]() []*cobra.Command {
	var res = []*cobra.Command{listCmd[T](), showCmd[T]()}

	// only some resources support run
	if runCommand := runCmd[T](); runCommand != nil {
		res = append(res, runCommand)
	}
	return res
}

func runCmd[T modconfig.HclResource]() *cobra.Command {
	var empty T

	switch any(empty).(type) {
	case *modconfig.Query:
		return queryRunCmd()
	case *modconfig.Dashboard:
		return dashboardRunCmd()
	case *modconfig.Benchmark:
		return checkCmd[*modconfig.Benchmark]()
	case *modconfig.Control:
		return checkCmd[*modconfig.Control]()
	}

	return nil
}

// functions to return descriptions and use
func showCommandUse(typeName string) string {
	return fmt.Sprintf("show <%s-name>", typeName)
}

func resourceCommandShortDescription(typeName string) string {
	// TODO
	return fmt.Sprintf("%s commands", typeName)
}
func resourceCommandLongDescription(typeName string) string {
	// TODO
	return fmt.Sprintf("%s commands", typeName)
}

func listCommandShortDescription(typeName string) string {
	return fmt.Sprintf("List %ss from the current mod and its direct dependents", typeName)
}
func listCommandLongDescription(typeName string) string {
	return fmt.Sprintf("List %ss from the current mod and its direct dependents", typeName)
}

func showCommandShortDescription(typeName string) string {
	return fmt.Sprintf("Show %ss from the current mod and its direct dependents", typeName)
}
func showCommandLongDescription(typeName string) string {
	return fmt.Sprintf("Show %ss from the current mod and its direct dependents", typeName)
}
