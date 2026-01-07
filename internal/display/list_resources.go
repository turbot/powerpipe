package display

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"

	"github.com/turbot/pipe-fittings/v2/constants"
	"github.com/turbot/pipe-fittings/v2/error_helpers"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/pipe-fittings/v2/printers"
	"github.com/turbot/pipe-fittings/v2/schema"
	"github.com/turbot/pipe-fittings/v2/workspace"
	localcmdconfig "github.com/turbot/powerpipe/internal/cmdconfig"
	localconstants "github.com/turbot/powerpipe/internal/constants"
	"github.com/turbot/powerpipe/internal/powerpipeconfig"
	"github.com/turbot/powerpipe/internal/resourceindex"
	"github.com/turbot/powerpipe/internal/resources"
	pworkspace "github.com/turbot/powerpipe/internal/workspace"
)

func ListResources[T modconfig.ModTreeItem](cmd *cobra.Command) {
	ctx := cmd.Context()
	modLocation := viper.GetString(constants.ArgModLocation)

	// Check if we should use lazy loading for this resource type
	if shouldUseLazyLoadingForList[T]() {
		listResourcesLazy[T](ctx, cmd, modLocation)
		return
	}

	// Fall back to eager loading for Variables and Mods
	listResourcesEager[T](ctx, cmd, modLocation)
}

// shouldUseLazyLoadingForList returns true if the resource type supports lazy loading.
// Variables and Mods need eager loading because:
// - Variables need their values resolved
// - Mods have special handling
func shouldUseLazyLoadingForList[T modconfig.ModTreeItem]() bool {
	// Check if workspace preload is enabled - this forces eager loading
	if localconstants.WorkspacePreloadEnabled() {
		slog.Debug("Workspace preload enabled - using eager loading for list")
		return false
	}

	var empty T
	switch any(empty).(type) {
	case *modconfig.Variable, *modconfig.Mod:
		// Variables need resolved values, Mods have special handling
		return false
	default:
		return true
	}
}

// listResourcesLazy uses the resource index for fast listing.
func listResourcesLazy[T modconfig.ModTreeItem](ctx context.Context, cmd *cobra.Command, modLocation string) {
	slog.Debug("Using lazy loading for list command")

	// Load lazy workspace
	lw, err := pworkspace.LoadLazy(ctx, modLocation,
		pworkspace.WithPipelingConnections(powerpipeconfig.GlobalConfig.PipelingConnections),
	)
	if err != nil {
		// Fall back to eager loading on error
		slog.Debug("Lazy loading failed, falling back to eager", "error", err)
		listResourcesEager[T](ctx, cmd, modLocation)
		return
	}
	defer lw.Close()

	if !lw.ModfileExists() {
		error_helpers.FailOnError(localconstants.ErrorNoModDefinition{})
	}

	// Get the resource type name
	typeName := getResourceTypeName[T]()

	// Get entries from the index
	index := lw.GetIndex()
	entries := filterIndexEntries(index, typeName, cmd)

	// Convert to ListableIndexEntry for display
	listable := make([]*ListableIndexEntry, 0, len(entries))
	for _, entry := range entries {
		listable = append(listable, NewListableIndexEntry(entry))
	}

	// Print the results
	printIndexListResult(ctx, cmd, listable)
}

// getResourceTypeName returns the schema block type name for a resource type.
func getResourceTypeName[T modconfig.ModTreeItem]() string {
	var empty T
	switch any(empty).(type) {
	case *resources.Benchmark:
		return schema.BlockTypeBenchmark
	case *resources.DetectionBenchmark:
		return "detection_benchmark"
	case *resources.Control:
		return schema.BlockTypeControl
	case *resources.Detection:
		return "detection"
	case *resources.Dashboard:
		return schema.BlockTypeDashboard
	case *resources.Query:
		return schema.BlockTypeQuery
	case *modconfig.Variable:
		return schema.BlockTypeVariable
	default:
		// Fall back to using the generic type converter
		return resources.GenericTypeToBlockType[T]()
	}
}

// filterIndexEntries filters index entries based on resource type and output format.
func filterIndexEntries(index *resourceindex.ResourceIndex, typeName string, cmd *cobra.Command) []*resourceindex.IndexEntry {
	allEntries := index.List()
	var filtered []*resourceindex.IndexEntry

	// For benchmarks in pretty/plain output, only show top-level
	showOnlyTopLevel := false
	if typeName == schema.BlockTypeBenchmark || typeName == "detection_benchmark" {
		output := viper.GetString(constants.ArgOutput)
		if output == constants.OutputFormatPretty || output == constants.OutputFormatPlain {
			showOnlyTopLevel = true
		}
	}

	// Get the main mod name for filtering
	mainModName := index.ModName

	for _, entry := range allEntries {
		// Filter by type
		if typeName == schema.BlockTypeBenchmark {
			// For benchmark list, include both benchmarks and detection_benchmarks
			if entry.Type != schema.BlockTypeBenchmark && entry.Type != "detection_benchmark" {
				continue
			}
		} else if entry.Type != typeName {
			continue
		}

		// For pretty/plain output of benchmarks, only show top-level
		if showOnlyTopLevel && !entry.IsTopLevel {
			continue
		}

		// Include resources from main mod and dependencies
		// For dependencies, we want to show them in the list
		filtered = append(filtered, entry)
	}

	// Sort by mod name then resource name for consistent output
	// (sorting is handled by printable_hcl_resource.go)
	_ = mainModName // Used for filtering logic above

	return filtered
}

// printIndexListResult prints the list of index entries.
func printIndexListResult(ctx context.Context, cmd *cobra.Command, items []*ListableIndexEntry) {
	printer, err := printers.GetPrinter[*ListableIndexEntry](cmd)
	if err != nil {
		error_helpers.ShowErrorWithMessage(ctx, err, "failed obtaining printer")
		return
	}

	printableResource := NewPrintableHclResource[*ListableIndexEntry](items)
	err = printer.PrintResource(ctx, printableResource, cmd.OutOrStdout())
	if err != nil {
		error_helpers.ShowErrorWithMessage(ctx, err, "failed when printing")
	}
}

// listResourcesEager uses eager loading for resources that need it (Variables, Mods).
func listResourcesEager[T modconfig.ModTreeItem](ctx context.Context, cmd *cobra.Command, modLocation string) {
	// build options to specify which blocks we need to load (based on type T)
	listOpts := getListLoadWorkspaceOpts[T]()
	w, errAndWarnings := pworkspace.LoadWorkspacePromptingForVariables(ctx, modLocation, listOpts...)
	error_helpers.FailOnError(errAndWarnings.GetError())
	if !w.ModfileExists() {
		error_helpers.FailOnError(localconstants.ErrorNoModDefinition{})
	}

	// get resource filter depending on resource type and output type
	resourceFilter := getListResourceFilter[T](&w.Workspace)
	resourceList, err := workspace.FilterWorkspaceResourcesOfType[T](&w.Workspace, resourceFilter)
	if err != nil {
		error_helpers.ShowErrorWithMessage(ctx, err, "failed to filter resources")
		return
	}

	// TODO K TACTICAL for benchmark list, include detection benchmarks
	// https://github.com/turbot/powerpipe/issues/609
	var empty T
	if _, ok := any(empty).(*resources.Benchmark); !ok {
		printListResult[T](ctx, cmd, resourceList)
	} else {
		// list detection benchmarks
		resourceFilter := getListResourceFilter[*resources.DetectionBenchmark](&w.Workspace)
		detectionBenchmarkList, err := workspace.FilterWorkspaceResourcesOfType[*resources.DetectionBenchmark](&w.Workspace, resourceFilter)
		if err != nil {
			error_helpers.ShowErrorWithMessage(ctx, err, "failed to filter resources")
			return
		}
		// build a separate list of all benchmarks
		var l = make(map[string]modconfig.ModTreeItem)
		for k, v := range resourceList {
			l[k] = v
		}
		for k, v := range detectionBenchmarkList {
			l[k] = v
		}
		printListResult[modconfig.ModTreeItem](ctx, cmd, l)
	}
}

func printListResult[T modconfig.ModTreeItem](ctx context.Context, cmd *cobra.Command, resourceList map[string]T) {
	printer, err := printers.GetPrinter[T](cmd)
	if err != nil {
		error_helpers.ShowErrorWithMessage(ctx, err, "failed obtaining printer")
		return
	}
	printableResource := NewPrintableHclResource[T](maps.Values(resourceList))

	err = printer.PrintResource(ctx, printableResource, cmd.OutOrStdout())
	if err != nil {
		error_helpers.ShowErrorWithMessage(ctx, err, "failed when printing")
		return
	}
}

func getListResourceFilter[T modconfig.ModTreeItem](w *workspace.Workspace) workspace.ResourceFilter {
	var res = workspace.ResourceFilter{}

	var empty T
	if _, ok := any(empty).(*resources.Benchmark); ok {

		// if T is benchmark, and if output is pretty or plain, only show top level benchmarks
		if viper.GetString(constants.ArgOutput) == constants.OutputFormatPretty || viper.GetString(constants.ArgOutput) == constants.OutputFormatPlain {
			// build a lookup of mod names to filter on
			var modNames = map[string]struct{}{}
			for _, mod := range w.Mods {
				modNames[mod.Name()] = struct{}{}
			}

			// add a predicate which returns true only if the resources parent is one of these mods
			res.WherePredicate = func(item modconfig.HclResource) bool {
				mti, ok := item.(modconfig.ModTreeItem)
				if !ok {
					return false
				}

				parents := mti.GetParents()
				if len(parents) == 0 {
					return false
				}
				_, inTargetMod := modNames[parents[0].Name()]
				return inTargetMod
			}
		}
	}

	return res
}

// build LoadWorkspaceOptions to specify which blocks we need to load (based on type T)
func getListLoadWorkspaceOpts[T modconfig.ModTreeItem]() []pworkspace.LoadPowerpipeWorkspaceOption {
	var empty T
	var opts = []pworkspace.LoadPowerpipeWorkspaceOption{
		// pass connections
		pworkspace.WithPipelingConnections(powerpipeconfig.GlobalConfig.PipelingConnections),
		// disable late binding
		pworkspace.WithLateBinding(false),
		pworkspace.WithVariableValidation(false),
	}
	switch any(empty).(type) {
	case *modconfig.Mod:
		opts = append(opts, pworkspace.WithBlockType([]string{schema.BlockTypeMod}))
	case *modconfig.Variable:
		opts = append(opts, pworkspace.WithBlockType([]string{schema.BlockTypeVariable}))
	}
	return opts
}

func ShowResource[T modconfig.ModTreeItem](cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	modLocation := viper.GetString(constants.ArgModLocation)

	// Check if we should use lazy loading for this resource type
	if shouldUseLazyLoadingForShow[T]() {
		showResourceLazy[T](ctx, cmd, modLocation, args)
		return
	}

	// Fall back to eager loading for Variables and Mods
	showResourceEager[T](ctx, cmd, modLocation, args)
}

// shouldUseLazyLoadingForShow returns true if the resource type supports lazy loading for show.
// Currently disabled because:
// 1. Nested dashboard components (cards, charts inside containers) aren't in the index
// 2. Show output needs full resource metadata that lazy loading doesn't provide
// The main performance benefit is from list commands, so this is acceptable.
func shouldUseLazyLoadingForShow[T modconfig.ModTreeItem]() bool {
	// Always use eager loading for show commands to ensure complete output
	return false
}

// showResourceLazy uses lazy loading to show a single resource.
func showResourceLazy[T modconfig.ModTreeItem](ctx context.Context, cmd *cobra.Command, modLocation string, args []string) {
	slog.Debug("Using lazy loading for show command")

	// Load lazy workspace
	lw, err := pworkspace.LoadLazy(ctx, modLocation,
		pworkspace.WithPipelingConnections(powerpipeconfig.GlobalConfig.PipelingConnections),
	)
	if err != nil {
		// Fall back to eager loading on error
		slog.Debug("Lazy loading failed, falling back to eager", "error", err)
		showResourceEager[T](ctx, cmd, modLocation, args)
		return
	}
	defer lw.Close()

	if !lw.ModfileExists() {
		error_helpers.FailOnError(localconstants.ErrorNoModDefinition{})
	}

	// Parse the target name
	if len(args) != 1 {
		error_helpers.FailOnError(fmt.Errorf("show command requires exactly one argument"))
		return
	}

	// Load the resource on-demand
	typeName := getResourceTypeName[T]()
	fullName := resolveResourceName(lw.GetIndex(), args[0], typeName)

	// Use the lazy workspace's Load method to get the full resource
	resource, err := lw.LoadResource(ctx, fullName)
	if err != nil {
		error_helpers.FailOnError(fmt.Errorf("failed to load resource %s: %w", fullName, err))
		return
	}

	// Cast to the expected type
	target, ok := resource.(T)
	if !ok {
		// Try detection benchmark special case
		if detBench, ok := resource.(*resources.DetectionBenchmark); ok {
			err = showTarget(ctx, cmd, detBench)
			if err != nil {
				error_helpers.ShowErrorWithMessage(ctx, err, "failed when printing")
			}
			return
		}
		error_helpers.FailOnError(fmt.Errorf("resource %s is not of expected type", fullName))
		return
	}

	err = showTarget(ctx, cmd, target)
	if err != nil {
		error_helpers.ShowErrorWithMessage(ctx, err, "failed when printing")
	}
}

// resolveResourceName resolves a short name to a full resource name.
func resolveResourceName(index *resourceindex.ResourceIndex, arg string, typeName string) string {
	// If already fully qualified, return as-is
	if len(splitResourceName(arg)) == 3 {
		return arg
	}

	// Try to find in index
	modName := index.ModName

	// Handle "type.name" format
	parts := splitResourceName(arg)
	if len(parts) == 2 {
		return fmt.Sprintf("%s.%s.%s", modName, parts[0], parts[1])
	}

	// Handle just "name" - prepend mod and type
	return fmt.Sprintf("%s.%s.%s", modName, typeName, arg)
}

// splitResourceName splits a resource name by dots.
func splitResourceName(name string) []string {
	var parts []string
	current := ""
	for _, c := range name {
		if c == '.' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

// showResourceEager uses eager loading to show a single resource.
func showResourceEager[T modconfig.ModTreeItem](ctx context.Context, cmd *cobra.Command, modLocation string, args []string) {
	// build options to specify which blocks we need to load (based on type T)
	opts := getListLoadWorkspaceOpts[T]()
	w, errAndWarnings := pworkspace.LoadWorkspacePromptingForVariables(ctx, modLocation, opts...)
	error_helpers.FailOnError(errAndWarnings.GetError())
	if !w.ModfileExists() {
		error_helpers.FailOnError(localconstants.ErrorNoModDefinition{})
	}

	targets, err := localcmdconfig.ResolveTargets[T](args, w)
	error_helpers.FailOnError(err)

	// show only supports a single target (should be enforced by cobra)
	if len(targets) > 1 {
		// not expected
		error_helpers.FailOnError(fmt.Errorf("show command only supports a single target"))
		return
	}

	// tactical - show detection benchmarks using the benchmark command
	// TODO once we remove DetectionBenchmarks and ResolveTargets returns [], this casting will not be needed
	// https://github.com/turbot/powerpipe/issues/609
	if _, ok := any(targets[0]).(*resources.DetectionBenchmark); ok {
		err = showTarget(ctx, cmd, targets[0].(*resources.DetectionBenchmark))
	} else {
		err = showTarget(ctx, cmd, targets[0].(T))
	}

	if err != nil {
		error_helpers.ShowErrorWithMessage(ctx, err, "failed when printing")
		return
	}
}

func showTarget[T modconfig.ModTreeItem](ctx context.Context, cmd *cobra.Command, target T) error {
	printer, err := printers.GetPrinter[T](cmd)
	if err != nil {
		error_helpers.ShowErrorWithMessage(ctx, err, "failed obtaining printer")
		return nil
	}
	printableResource := NewPrintableHclResource[T]([]T{target})

	err = printer.PrintResource(ctx, printableResource, cmd.OutOrStdout())
	return err
}
