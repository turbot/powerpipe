package initialisation

import (
	"context"
	"fmt"
	"log/slog"
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

	if !w.ModfileExists() && commandRequiresModfile(cmd, cmdArgs) {
		return NewErrorInitData(localconstants.ErrorNoModDefinition{})
	}
	i := &InitData{
		Result:        &InitResult{},
		ExportManager: export.NewManager(),
	}

	i.Workspace = w
	i.Result.Warnings = errAndWarnings.Warnings

	// resolve target resources
	targets, err := cmdconfig.ResolveTargets[T](cmdArgs, w)
	if err != nil {
		i.Result.Error = err
		return i
	}
	i.Targets = targets

	// now do the actual initialisation
	i.Init(ctx, cmdArgs...)

	return i
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
	if i.Workspace != nil {
		i.Workspace.Close()
	}
	if i.DefaultClient != nil {
		i.DefaultClient.Close(ctx)
	}
}

// GetSingleTarget validates there is only a single target and returns it
func (i *InitData) GetSingleTarget() (modconfig.ModTreeItem, error) {

	// cobra should validate this
	if len(i.Targets) != 1 {

		return nil, sperr.New("expected a single target")
	}
	return i.Targets[0], nil
}
