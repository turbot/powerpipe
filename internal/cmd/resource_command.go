package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/turbot/pipe-fittings/cmdconfig"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/powerpipe/internal/display"
	"strings"
)

func resourceCmd[T modconfig.HclResource]() *cobra.Command {
	var cmd = &cobra.Command{
		Use:  fmt.Sprintf("%s [command]", getGenericTypeName[T]()),
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			error_helpers.FailOnError(err)
		},
	}

	cmdconfig.OnCmd(cmd)

	cmd.AddCommand(
		listCmd[T](),
		showCmd[T](),
		runCmd[T]())

	return cmd
}

func listCmd[T modconfig.HclResource]() *cobra.Command {
	typeName := getGenericTypeName[T]()
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
	typeName := getGenericTypeName[T]()

	var cmd = &cobra.Command{
		Use:   showCommandUse(typeName),
		Args:  cobra.ExactArgs(1),
		Run:   func(cmd *cobra.Command, args []string) { display.ShowResource[T](cmd) },
		Short: showCommandShortDescription(typeName),
		Long:  showCommandLongDescription(typeName),
	}
	// initialize hooks
	cmdconfig.OnCmd(cmd)

	return cmd
}

func runCmd[T modconfig.HclResource]() *cobra.Command {
	typeName := getGenericTypeName[T]()
	var cmd = &cobra.Command{
		Use:   "list",
		Args:  cobra.NoArgs,
		Run:   func(cmd *cobra.Command, args []string) { display.ListResources[*modconfig.Dashboard](cmd) },
		Short: listCommandShortDescription(typeName),
		Long:  listCommandLongDescription(typeName)}
	// initialize hooks
	cmdconfig.OnCmd(cmd)

	return cmd
}

func showCommandUse(typeName string) string {
	return fmt.Sprintf("\"show <%s-name>", typeName)
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

// return lower case form of type unqualified name
func getGenericTypeName[T any]() string {
	longName := fmt.Sprintf("%T", *new(T))
	split := strings.Split(longName, ".")
	return strings.ToLower(split[len(split)-1])
}
