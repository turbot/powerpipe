package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"
	"github.com/turbot/pipe-fittings/cmdconfig"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/schema"
	localcmdconfig "github.com/turbot/powerpipe/internal/cmdconfig"
	localconstants "github.com/turbot/powerpipe/internal/constants"
	"github.com/turbot/powerpipe/internal/display"
	"strings"
)

// variable used to assign the output mode flag
var outputMode = localconstants.OutputModePretty

type ResourceCommandOption func(*ResourceCommandConfig)

type ResourceCommandConfig struct {
	cmd string
}

func newResourceCommandConfig[T modconfig.HclResource]() *ResourceCommandConfig {
	typeName := localcmdconfig.GetGenericTypeName[T]()
	return &ResourceCommandConfig{
		cmd: typeName,
	}
}

func withCmdName(name string) ResourceCommandOption {
	return func(m *ResourceCommandConfig) {
		m.cmd = name
	}
}

func resourceCmd[T modconfig.HclResource](opts ...ResourceCommandOption) *cobra.Command {
	cfg := newResourceCommandConfig[T]()
	for _, o := range opts {
		o(cfg)
	}

	var cmd = &cobra.Command{
		Use:   fmt.Sprintf("%s [command]", cfg.cmd),
		Args:  cobra.NoArgs,
		Short: resourceCommandShortDescription(cfg.cmd),
		Long:  resourceCommandLongDescription(cfg.cmd),
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
	cmdconfig.OnCmd(cmd).
		AddVarFlag(enumflag.New(&outputMode, constants.ArgOutput, localconstants.OutputModeIds, enumflag.EnumCaseInsensitive),
			constants.ArgOutput,
			fmt.Sprintf("Output format; one of: %s", strings.Join(localconstants.FlagValues(localconstants.OutputModeIds), ", ")))

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
	cmdconfig.OnCmd(cmd).
		AddVarFlag(enumflag.New(&outputMode, constants.ArgOutput, localconstants.OutputModeIds, enumflag.EnumCaseInsensitive),
			constants.ArgOutput,
			fmt.Sprintf("Output format; one of: %s", strings.Join(localconstants.FlagValues(localconstants.OutputModeIds), ", ")))

	return cmd
}

// determine which resource commands apply to this resource
func getResourceCommands[T modconfig.HclResource]() []*cobra.Command {
	typeName := localcmdconfig.GetGenericTypeName[T]()

	var res = []*cobra.Command{listCmd[T](), showCmd[T]()}

	// only some resources support run
	if runCommand := runCmd[T](); runCommand != nil {
		res = append(res, runCommand)
	}

	// spacial case for dashboard
	if typeName == schema.BlockTypeDashboard {
		res = append(res, dashboardChildCommands()...)
	}

	return res
}

func dashboardChildCommands() []*cobra.Command {
	return []*cobra.Command{
		resourceCmd[*modconfig.DashboardCard](withCmdName("card")),
		resourceCmd[*modconfig.DashboardChart](withCmdName("chart")),
		resourceCmd[*modconfig.DashboardContainer](withCmdName("container")),
		resourceCmd[*modconfig.DashboardFlow](withCmdName("flow")),
		resourceCmd[*modconfig.DashboardGraph](withCmdName("graph")),
		resourceCmd[*modconfig.DashboardHierarchy](withCmdName("hierarchy")),
		resourceCmd[*modconfig.DashboardImage](withCmdName("image")),
		resourceCmd[*modconfig.DashboardInput](withCmdName("input")),
		resourceCmd[*modconfig.DashboardTable](withCmdName("table")),
		resourceCmd[*modconfig.DashboardText](withCmdName("text")),
	}

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
