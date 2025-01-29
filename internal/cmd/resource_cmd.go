package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"
	"github.com/turbot/pipe-fittings/cmdconfig"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/schema"
	"github.com/turbot/pipe-fittings/utils"
	localconstants "github.com/turbot/powerpipe/internal/constants"
	"github.com/turbot/powerpipe/internal/display"
	"github.com/turbot/powerpipe/internal/resources"
)

// variable used to assign the output mode flag
var outputMode = localconstants.OutputModePretty

type ResourceCommandOption func(*ResourceCommandConfig)

type ResourceCommandConfig struct {
	cmd string
}

func newResourceCommandConfig[T modconfig.ModTreeItem]() *ResourceCommandConfig {
	typeName := resources.GenericTypeToBlockType[T]()
	return &ResourceCommandConfig{
		cmd: typeName,
	}
}

func withCmdName(name string) ResourceCommandOption {
	return func(m *ResourceCommandConfig) {
		m.cmd = name
	}
}

func resourceCmd[T modconfig.ModTreeItem](opts ...ResourceCommandOption) *cobra.Command {
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

func listCmd[T modconfig.ModTreeItem]() *cobra.Command {
	typeName := resources.GenericTypeToBlockType[T]()
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
			fmt.Sprintf("Output format; one of: %s", strings.Join(constants.FlagValues(localconstants.OutputModeIds), ", ")))

	return cmd
}

func showCmd[T modconfig.ModTreeItem]() *cobra.Command {
	typeName := resources.GenericTypeToBlockType[T]()

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
			fmt.Sprintf("Output format; one of: %s", strings.Join(constants.FlagValues(localconstants.OutputModeIds), ", ")))

	return cmd
}

// determine which resource commands apply to this resource
func getResourceCommands[T modconfig.ModTreeItem]() []*cobra.Command {
	typeName := resources.GenericTypeToBlockType[T]()

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
	res := []*cobra.Command{
		resourceCmd[*resources.DashboardCard](withCmdName("card")),
		resourceCmd[*resources.DashboardChart](withCmdName("chart")),
		resourceCmd[*resources.DashboardContainer](withCmdName("container")),
		resourceCmd[*resources.DashboardFlow](withCmdName("flow")),
		resourceCmd[*resources.DashboardGraph](withCmdName("graph")),
		resourceCmd[*resources.DashboardHierarchy](withCmdName("hierarchy")),
		resourceCmd[*resources.DashboardImage](withCmdName("image")),
		resourceCmd[*resources.DashboardTable](withCmdName("table")),
		resourceCmd[*resources.Detection](withCmdName("detection")),
		resourceCmd[*resources.DetectionBenchmark](withCmdName("detectionbenchmark")),
		resourceCmd[*resources.DashboardText](withCmdName("text")),
	}

	// set all to hidden
	for _, cmd := range res {
		cmd.Hidden = true
	}
	return res
}

func runCmd[T modconfig.HclResource]() *cobra.Command {
	var empty T

	switch any(empty).(type) {
	case *resources.Query:
		return queryRunCmd()
	case *resources.Dashboard:
		return dashboardRunCmd()
	case *resources.Benchmark:
		return checkCmd[*resources.Benchmark]()
	case *resources.DetectionBenchmark:
		return detectionRunCmd[*resources.DetectionBenchmark]()
	case *resources.Detection:
		return detectionRunCmd[*resources.Detection]()
	case *resources.Control:
		return checkCmd[*resources.Control]()
	}

	return nil
}

// functions to return descriptions and use
func showCommandUse(typeName string) string {
	return fmt.Sprintf("show <%s-name>", typeName)
}

func resourceCommandShortDescription(typeName string) string {
	switch typeName {
	case schema.BlockTypeQuery:
		return "List, view, and run Powerpipe queries"

	case schema.BlockTypeControl:
		return "List, view, and run Powerpipe controls"

	case schema.BlockTypeBenchmark:
		return "List, view, and run Powerpipe benchmarks"

	case schema.BlockTypeDashboard:
		return "List, view, and run Powerpipe dashboards in batch mode"

	case schema.BlockTypeVariable:
		return "Manage Powerpipe variables in the current mod and its direct dependents"

	case schema.BlockTypeDetection:
		return "List, view, and run Powerpipe detections"

	default:
		return fmt.Sprintf("%s commands", typeName)
	}
}

func resourceCommandLongDescription(typeName string) string {
	switch typeName {
	case schema.BlockTypeQuery:
		return `List, view, and run Powerpipe queries and its direct dependents.
		
Run a query from the current mod or its direct dependents or from a Powerpipe server instance or
view details of a query from the current mod or its direct dependents or from a Powerpipe server instance.`

	case schema.BlockTypeControl:
		return `List, view, and run Powerpipe controls and its direct dependents.

Run a control from the current mod or its direct dependents or from a Powerpipe server instance or
view details of a control from the current mod or its direct dependents or from a Powerpipe server instance.`

	case schema.BlockTypeBenchmark:
		return `List, view, and run Powerpipe benchmarks and its direct dependents.
		
Run a benchmark from the current mod or its direct dependents or from a Powerpipe server instance or
view details of a benchmark from the current mod or its direct dependents or from a Powerpipe server instance.`

	case schema.BlockTypeDashboard:
		return `List, view, and run Powerpipe dashboards in batch mode. 
		
To run dashboards interactively, run powerpipe server.`

	case schema.BlockTypeVariable:
		return `List variables from the current mod and its direct dependents or show details of a variable
from the current mod or its direct dependents or from a Powerpipe server instance.`

	case schema.BlockTypeDetection:
		return `List, view, and run Powerpipe detections and its direct dependents.

Run a detection from the current mod or its direct dependents or from a Powerpipe server instance or
view details of a detection from the current mod or its direct dependents or from a Powerpipe server instance.`

	default:
		return fmt.Sprintf("%s commands", typeName)
	}
}

func listCommandShortDescription(typeName string) string {
	return fmt.Sprintf("List %s from the current mod and its direct dependents", utils.Pluralize(typeName, 0))
}
func listCommandLongDescription(typeName string) string {
	return fmt.Sprintf("List %s from the current mod and its direct dependents", utils.Pluralize(typeName, 0))
}

func showCommandShortDescription(typeName string) string {
	return fmt.Sprintf("Show %s from the current mod and its direct dependents", utils.Pluralize(typeName, 0))
}
func showCommandLongDescription(typeName string) string {
	return fmt.Sprintf("Show %s from the current mod and its direct dependents", utils.Pluralize(typeName, 0))
}
