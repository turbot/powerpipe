package initialisation

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/v2/app_specific"
	"github.com/turbot/pipe-fittings/v2/backend"
	"github.com/turbot/pipe-fittings/v2/connection"
	"github.com/turbot/pipe-fittings/v2/constants"
	"github.com/turbot/pipe-fittings/v2/error_helpers"
	"github.com/turbot/pipe-fittings/v2/export"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/pipe-fittings/v2/modinstaller"
	"github.com/turbot/pipe-fittings/v2/plugin"
	"github.com/turbot/pipe-fittings/v2/statushooks"
	"github.com/turbot/pipe-fittings/v2/utils"
	"github.com/turbot/powerpipe/internal/cmdconfig"
	localconstants "github.com/turbot/powerpipe/internal/constants"
	"github.com/turbot/powerpipe/internal/dashboardexecute"
	"github.com/turbot/powerpipe/internal/db_client"
	"github.com/turbot/powerpipe/internal/powerpipeconfig"
	"github.com/turbot/powerpipe/internal/resources"
	"github.com/turbot/powerpipe/internal/timing"
	"github.com/turbot/powerpipe/internal/workspace"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
	"github.com/turbot/steampipe-plugin-sdk/v5/telemetry"
)

// clientResult holds the async result of database client creation
type clientResult struct {
	client           *db_client.DbClient
	csp              connection.ConnectionStringProvider
	searchPathConfig backend.SearchPathConfig
	err              error
}

type InitData struct {
	Workspace *workspace.PowerpipeWorkspace
	Result    *InitResult

	ShutdownTelemetry func()
	ExportManager     *export.Manager
	Targets           []modconfig.ModTreeItem
	DefaultClient     *db_client.DbClient

	DefaultDatabase         connection.ConnectionStringProvider
	DefaultSearchPathConfig backend.SearchPathConfig

	// LazyWorkspace is set when lazy loading is enabled
	// It wraps the PowerpipeWorkspace with on-demand resource loading
	LazyWorkspace *workspace.LazyWorkspace
}

func NewErrorInitData(err error) *InitData {
	return &InitData{
		Result: &InitResult{
			ErrorAndWarnings: error_helpers.NewErrorsAndWarning(err),
		},
	}
}

func NewInitData[T modconfig.ModTreeItem](ctx context.Context, cmd *cobra.Command, cmdArgs ...string) *InitData {
	defer timing.Track("NewInitData")()

	modLocation := viper.GetString(constants.ArgModLocation)

	i := &InitData{
		Result:        &InitResult{},
		ExportManager: export.NewManager(),
	}

	// Only use lazy loading for benchmark types
	// For other types (query, dashboard), lazy loading doesn't support target resolution
	var empty T
	_, isBenchmark := any(empty).(*resources.Benchmark)
	_, isDetectionBenchmark := any(empty).(*resources.DetectionBenchmark)
	useLazyLoading := isLazyLoadEnabled() && (isBenchmark || isDetectionBenchmark)

	if useLazyLoading {
		// Use lazy loading for benchmark commands - faster startup and lower memory
		slog.Debug("Loading workspace with lazy loading enabled")
		lw, err := workspace.LoadLazy(ctx, modLocation,
			workspace.WithPipelingConnections(powerpipeconfig.GlobalConfig.PipelingConnections),
		)
		if err != nil {
			return NewErrorInitData(fmt.Errorf("failed to load lazy workspace: %s", error_helpers.HandleCancelError(err).Error()))
		}
		i.LazyWorkspace = lw
		i.Workspace = lw.PowerpipeWorkspace
	} else {
		// Standard eager loading for non-benchmark commands or when lazy loading is disabled
		w, errAndWarnings := workspace.LoadWorkspacePromptingForVariables(ctx,
			modLocation,
			// pass connections
			workspace.WithPipelingConnections(powerpipeconfig.GlobalConfig.PipelingConnections),
			// disable late binding
			workspace.WithLateBinding(false),
		)
		if errAndWarnings.GetError() != nil {
			return NewErrorInitData(fmt.Errorf("failed to load workspace: %s", error_helpers.HandleCancelError(errAndWarnings.GetError()).Error()))
		}
		i.Workspace = w
		i.Result.Warnings = errAndWarnings.Warnings
	}

	if !i.Workspace.ModfileExists() && commandRequiresModfile(cmd, cmdArgs) {
		return NewErrorInitData(localconstants.ErrorNoModDefinition{})
	}

	// resolve target resources
	var targets []modconfig.ModTreeItem
	var err error

	if useLazyLoading {
		// For lazy loading (benchmark commands only), resolve targets using LoadBenchmarkForExecution
		// This bypasses the standard ResolveTargets which requires resources to be loaded in the workspace
		targets, err = i.resolveTargetsForLazyLoading(ctx, cmdArgs)
		if err != nil {
			i.Result.Error = err
			return i
		}
	} else {
		// Standard target resolution for eager loading
		targets, err = cmdconfig.ResolveTargets[T](cmdArgs, i.Workspace)
		if err != nil {
			i.Result.Error = err
			return i
		}
	}

	i.Targets = targets

	// now do the actual initialisation
	i.Init(ctx, cmdArgs...)

	return i
}

// isLazyLoadEnabled checks if lazy loading should be used.
// Lazy loading is enabled by default. Set POWERPIPE_WORKSPACE_PRELOAD=true to disable.
// Also disabled when tag/where filters are used (they require full workspace to iterate).
// Also disabled when dependency mods are installed (complex reference resolution needed).
func isLazyLoadEnabled() bool {
	// Check if workspace preload is enabled - this forces eager loading
	if localconstants.WorkspacePreloadEnabled() {
		slog.Info("Workspace preload enabled via POWERPIPE_WORKSPACE_PRELOAD - using eager loading")
		return false
	}

	// Tag and where filters require iterating over all controls in the workspace,
	// which isn't compatible with lazy loading. Disable lazy loading when these are set.
	if viper.IsSet(constants.ArgTag) || viper.IsSet(constants.ArgWhere) {
		slog.Debug("Tag or where filter set - disabling lazy loading")
		return false
	}

	// Check for dependency mods - they require full workspace loading for proper
	// reference resolution across mod boundaries.
	modLocation := viper.GetString(constants.ArgModLocation)
	modsDir := filepath.Join(modLocation, ".powerpipe", "mods")
	if info, err := os.Stat(modsDir); err == nil && info.IsDir() {
		// Check if there are any actual mods installed
		entries, err := os.ReadDir(modsDir)
		if err == nil && len(entries) > 0 {
			slog.Debug("Dependency mods detected - disabling lazy loading")
			return false
		}
	}

	// Default: lazy loading enabled
	return true
}

func commandRequiresModfile(cmd *cobra.Command, args []string) bool {
	// all commands using initData require a modfile EXCEPT query run if it is a raw sql query
	if utils.CommandFullKey(cmd) != "powerpipe.query.run" {
		return true
	}

	// if the command is query run, and the first argument is a raw sql query, we don't need a modfile
	_, argIsNamedResource := workspace.SqlLooksLikeExecutableResource(args[0])
	return argIsNamedResource
}

func (i *InitData) RegisterExporters(exporters ...export.Exporter) error {
	for _, e := range exporters {
		if err := i.ExportManager.Register(e); err != nil {
			return err
		}
	}

	return nil
}

func (i *InitData) Init(ctx context.Context, args ...string) {
	defer timing.Track("InitData.Init")()

	defer func() {
		if r := recover(); r != nil {
			i.Result.Error = helpers.ToError(r)
		}
		// if there is no error, return context cancellation error (if any)
		if i.Result.Error == nil {
			i.Result.Error = ctx.Err()
		}
	}()

	slog.Info("Initializing...")

	// code after this depends of i.Workspace being defined. make sure that it is
	if i.Workspace == nil {
		i.Result.Error = sperr.WrapWithRootMessage(error_helpers.InvalidStateError, "InitData.Init called before setting up WorkspaceEvents")
		return
	}

	statushooks.SetStatus(ctx, "Initializing")

	// Start database client creation in background - this can run concurrently
	// with telemetry init and mod installation
	clientChan := make(chan clientResult, 1)
	var clientWg sync.WaitGroup
	clientWg.Add(1)

	go func() {
		defer clientWg.Done()
		defer close(clientChan)
		defer timing.Track("db_client.CreateAsync")()

		// Get default database config
		csp, searchPathConfig, err := db_client.GetDefaultDatabaseConfig(i.Workspace.Mod)
		if err != nil {
			clientChan <- clientResult{err: err}
			return
		}

		// Build backend options
		var opts []backend.BackendOption
		if !searchPathConfig.Empty() {
			opts = append(opts, backend.WithSearchPathConfig(searchPathConfig))
		}

		connectionString, err := csp.GetConnectionString()
		if err != nil {
			clientChan <- clientResult{err: err}
			return
		}

		// Create the client
		client, err := db_client.NewDbClient(ctx, connectionString, opts...)
		clientChan <- clientResult{
			client:           client,
			csp:              csp,
			searchPathConfig: searchPathConfig,
			err:              err,
		}
	}()

	// Initialize telemetry concurrently with DB client creation
	func() {
		defer timing.Track("telemetry.Init")()
		shutdownTelemetry, err := telemetry.Init(app_specific.AppName)
		if err != nil {
			i.Result.AddWarnings(err.Error())
		} else {
			i.ShutdownTelemetry = shutdownTelemetry
		}
	}()

	// Install mod dependencies if needed (this defaults to true for dashboard and check commands
	// and will always be false for query command)
	// This also runs concurrently with DB client creation
	if viper.GetBool(constants.ArgModInstall) {
		func() {
			defer timing.Track("modinstaller.InstallWorkspaceDependencies")()
			statushooks.SetStatus(ctx, "Installing workspace dependencies")
			slog.Info("Installing workspace dependencies")
			opts := modinstaller.NewInstallOpts(i.Workspace.Mod)
			// arg pull should always be set (to a default at least) if ArgModInstall is set
			opts.UpdateStrategy = viper.GetString(constants.ArgPull)
			// use force install so that errors are ignored during installation
			// (we are validating prereqs later)
			opts.Force = true
			_, err := modinstaller.InstallWorkspaceDependencies(ctx, opts)
			if err != nil {
				i.Result.Error = err
			}
		}()
		if i.Result.Error != nil {
			// Wait for DB goroutine to complete before returning
			clientWg.Wait()
			return
		}
	}

	// Now wait for database client creation to complete
	statushooks.SetStatus(ctx, "Connecting to database")
	result := <-clientChan

	if result.err != nil {
		i.Result.Error = result.err
		return
	}

	// Store results
	i.DefaultClient = result.client
	i.DefaultDatabase = result.csp
	i.DefaultSearchPathConfig = result.searchPathConfig

	// Validate mod requirements - this needs the client
	func() {
		defer timing.Track("validateModRequirementsRecursively")()
		validationWarnings := validateModRequirementsRecursively(i.Workspace.Mod, i.DefaultClient)
		i.Result.AddWarnings(validationWarnings...)
	}()

	// Create the dashboard executor, passing the default client inside a client map
	func() {
		defer timing.Track("NewDashboardExecutor")()
		clientMap := db_client.NewClientMap().Add(i.DefaultClient, i.DefaultSearchPathConfig)
		dashboardexecute.Executor = dashboardexecute.NewDashboardExecutor(clientMap, i.DefaultDatabase, i.DefaultSearchPathConfig)
	}()
}

func validateModRequirementsRecursively(mod *modconfig.Mod, client *db_client.DbClient) []string {
	var validationErrors []string

	var pluginVersionMap = &plugin.PluginVersionMap{
		Database: client.Backend.ConnectionString(),
		Backend:  client.Backend.Name(),
	}
	// if the backend is steampipe, populate the available plugins
	if steampipeBackend, ok := client.Backend.(*backend.SteampipeBackend); ok {
		pluginVersionMap.AvailablePlugins = steampipeBackend.PluginVersions
	}
	// validate this mod
	for _, err := range mod.ValidateRequirements(pluginVersionMap) {
		validationErrors = append(validationErrors, err.Error())
	}

	// validate dependent mods
	for childDependencyName, childMod := range mod.GetModResources().GetMods() {
		if childDependencyName == "local" || mod.DependencyName == childMod.DependencyName {
			// this is a reference to self - skip (otherwise we will end up with a recursion loop)
			continue
		}
		childValidationErrors := validateModRequirementsRecursively(childMod, client)
		validationErrors = append(validationErrors, childValidationErrors...)
	}

	return validationErrors
}

func (i *InitData) Cleanup(ctx context.Context) {
	if i.ShutdownTelemetry != nil {
		i.ShutdownTelemetry()
	}
	// Close lazy workspace if present (this also closes the embedded PowerpipeWorkspace)
	if i.LazyWorkspace != nil {
		i.LazyWorkspace.Close()
	} else if i.Workspace != nil {
		i.Workspace.Close()
	}
	if i.DefaultClient != nil {
		i.DefaultClient.Close(ctx)
	}
}

// IsLazy returns true if lazy loading is enabled for this init data
func (i *InitData) IsLazy() bool {
	return i.LazyWorkspace != nil
}

// GetWorkspaceProvider returns the workspace as a WorkspaceProvider interface,
// which can be either a LazyWorkspace or PowerpipeWorkspace
func (i *InitData) GetWorkspaceProvider() workspace.WorkspaceProvider {
	if i.LazyWorkspace != nil {
		return i.LazyWorkspace
	}
	return i.Workspace
}

// GetSingleTarget validates there is only a single target and returns it
func (i *InitData) GetSingleTarget() (modconfig.ModTreeItem, error) {

	// cobra should validate this
	if len(i.Targets) != 1 {

		return nil, sperr.New("expected a single target")
	}
	return i.Targets[0], nil
}

// resolveTargetsForLazyLoading resolves benchmark targets for lazy loading mode.
// This bypasses the standard ResolveTargets which requires resources to be loaded in the workspace.
// Instead, it uses the lazy workspace's index and LoadBenchmarkForExecution to resolve and load benchmarks.
func (i *InitData) resolveTargetsForLazyLoading(ctx context.Context, cmdArgs []string) ([]modconfig.ModTreeItem, error) {
	if i.LazyWorkspace == nil {
		return nil, fmt.Errorf("resolveTargetsForLazyLoading called but LazyWorkspace is nil")
	}

	if len(cmdArgs) == 0 {
		return nil, nil
	}

	// Handle "all" argument
	if len(cmdArgs) == 1 && cmdArgs[0] == "all" {
		return i.resolveAllBenchmarksForLazyLoading(ctx)
	}

	targets := make([]modconfig.ModTreeItem, 0, len(cmdArgs))
	for _, arg := range cmdArgs {
		// Parse the resource name from the argument
		fullName, err := i.parseTargetName(arg)
		if err != nil {
			return nil, err
		}

		slog.Debug("Loading benchmark for execution (lazy)", "name", fullName)
		benchmark, err := i.LazyWorkspace.LoadBenchmarkForExecution(ctx, fullName)
		if err != nil {
			return nil, fmt.Errorf("failed to load benchmark %s: %w", fullName, err)
		}

		targets = append(targets, benchmark)
	}

	return targets, nil
}

// parseTargetName parses a target argument into a full resource name.
// Handles formats: "name", "type.name", "mod.type.name"
func (i *InitData) parseTargetName(arg string) (string, error) {
	parts := strings.Split(arg, ".")
	modName := i.LazyWorkspace.GetIndex().ModName

	switch len(parts) {
	case 1:
		// Just name, assume benchmark type
		return fmt.Sprintf("%s.benchmark.%s", modName, parts[0]), nil
	case 2:
		// type.name
		return fmt.Sprintf("%s.%s.%s", modName, parts[0], parts[1]), nil
	case 3:
		// mod.type.name - use as-is
		return arg, nil
	default:
		return "", fmt.Errorf("invalid target name format: %s", arg)
	}
}

// resolveAllBenchmarksForLazyLoading resolves all top-level benchmarks for the "all" argument.
func (i *InitData) resolveAllBenchmarksForLazyLoading(ctx context.Context) ([]modconfig.ModTreeItem, error) {
	// Get all top-level benchmarks from the index
	index := i.LazyWorkspace.GetIndex()
	modName := index.ModName

	var childTargets []modconfig.ModTreeItem

	for _, entry := range index.List() {
		// Only include top-level benchmarks from the current mod
		if entry.Type == "benchmark" && entry.IsTopLevel {
			// Check if this benchmark belongs to the current mod
			if strings.HasPrefix(entry.Name, modName+".") {
				slog.Debug("Loading benchmark for execution (lazy, all)", "name", entry.Name)
				benchmark, err := i.LazyWorkspace.LoadBenchmarkForExecution(ctx, entry.Name)
				if err != nil {
					return nil, fmt.Errorf("failed to load benchmark %s: %w", entry.Name, err)
				}
				childTargets = append(childTargets, benchmark)
			}
		}
	}

	if len(childTargets) == 0 {
		return nil, nil
	}

	// Create a root benchmark to hold all the benchmarks (same as handleAllArg does)
	resolvedItem := resources.NewRootBenchmarkWithChildren(i.Workspace.Mod, childTargets).(modconfig.ModTreeItem)
	return []modconfig.ModTreeItem{resolvedItem}, nil
}
