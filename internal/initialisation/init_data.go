package initialisation

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/app_specific"
	"github.com/turbot/pipe-fittings/backend"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/export"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/modinstaller"
	"github.com/turbot/pipe-fittings/statushooks"
	"github.com/turbot/pipe-fittings/workspace"
	"github.com/turbot/powerpipe/internal/cmdconfig"
	"github.com/turbot/powerpipe/internal/dashboardexecute"
	"github.com/turbot/powerpipe/internal/dashboardworkspace"
	"github.com/turbot/powerpipe/internal/db_client"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
	"github.com/turbot/steampipe-plugin-sdk/v5/telemetry"
)

type InitData[T modconfig.ModTreeItem] struct {
	Workspace       *workspace.Workspace
	WorkspaceEvents *dashboardworkspace.WorkspaceEvents
	Result          *InitResult

	ShutdownTelemetry func()
	ExportManager     *export.Manager
	Target            T
	DefaultClient     *db_client.DbClient
}

func NewErrorInitData[T modconfig.ModTreeItem](err error) *InitData[T] {
	return &InitData[T]{
		Result: &InitResult{
			ErrorAndWarnings: error_helpers.NewErrorsAndWarning(err),
		},
	}
}

func NewInitData[T modconfig.ModTreeItem](ctx context.Context, targetNames ...string) *InitData[T] {
	modLocation := viper.GetString(constants.ArgModLocation)

	w, errAndWarnings := workspace.LoadWorkspacePromptingForVariables(ctx, modLocation)
	if errAndWarnings.GetError() != nil {
		return NewErrorInitData[T](fmt.Errorf("failed to load workspace: %s", error_helpers.HandleCancelError(errAndWarnings.GetError()).Error()))
	}

	i := &InitData[T]{
		Result:        &InitResult{},
		ExportManager: export.NewManager(),
	}

	i.Workspace = w
	i.Result.Warnings = errAndWarnings.Warnings

	// now do the actual initialisation
	i.Init(ctx, targetNames...)

	return i
}

func (i *InitData[T]) RegisterExporters(exporters ...export.Exporter) error {
	for _, e := range exporters {
		if err := i.ExportManager.Register(e); err != nil {
			return err
		}
	}

	return nil
}

func (i *InitData[T]) Init(ctx context.Context, args ...string) {
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

	// code after this depends of i.WorkspaceEvents being defined. make sure that it is
	if i.Workspace == nil {
		i.Result.Error = sperr.WrapWithRootMessage(error_helpers.InvalidStateError, "InitData.Init called before setting up WorkspaceEvents")
		return
	}

	i.resolveTarget(args)
	if i.Result.Error != nil {
		return
	}
	statushooks.SetStatus(ctx, "Initializing")
	i.WorkspaceEvents = dashboardworkspace.NewWorkspaceEvents(i.Workspace)

	// initialise telemetry
	shutdownTelemetry, err := telemetry.Init(app_specific.AppName)
	if err != nil {
		i.Result.AddWarnings(err.Error())
	} else {
		i.ShutdownTelemetry = shutdownTelemetry
	}

	// install mod dependencies if needed
	if viper.GetBool(constants.ArgModInstall) {
		statushooks.SetStatus(ctx, "Installing workspace dependencies")
		slog.Info("Installing workspace dependencies")

		opts := modinstaller.NewInstallOpts(i.Workspace.Mod)
		// use force install so that errors are ignored during installation
		// (we are validating prereqs later)
		opts.Force = true
		_, err := modinstaller.InstallWorkspaceDependencies(ctx, opts)
		if err != nil {
			i.Result.Error = err
			return
		}
	}

	// retrieve cloud metadata
	cloudMetadata, err := getCloudMetadata(ctx)
	if err != nil {
		i.Result.Error = err
		return
	}

	// set cloud metadata (may be nil)
	i.Workspace.CloudMetadata = cloudMetadata

	// create default client
	// set the dashboard database and search patch config
	database, searchPathConfig := db_client.GetDefaultDatabaseConfig()
	// if there is a target, this may change the default database and search path - if it is in a dependency mod
	if !helpers.IsNil(modconfig.ModTreeItem(i.Target)) {
		database, searchPathConfig, err = db_client.GetDatabaseConfigForResource(i.Target, i.Workspace.Mod, database, searchPathConfig)
		if err != nil {
			i.Result.Error = err
			return
		}
	}
	// create client
	var opts []backend.ConnectOption
	if !searchPathConfig.Empty() {
		opts = append(opts, backend.WithSearchPathConfig(searchPathConfig))
	}
	client, err := db_client.NewDbClient(ctx, database, opts...)
	if err != nil {
		i.Result.Error = err
		return
	}
	i.DefaultClient = client

	// validate mod requirements
	validationWarnings := validateModRequirementsRecursively(i.Workspace.Mod, client)
	i.Result.AddWarnings(validationWarnings...)

	// create the dashboard executor, passing the default client inside a client map
	clientMap := db_client.NewClientMap().Add(client, searchPathConfig)
	dashboardexecute.Executor = dashboardexecute.NewDashboardExecutor(clientMap)
}

// resolve target resource, args and any target specific search path
func (i *InitData[T]) resolveTarget(args []string) {

	// resolve target resources
	target, err := cmdconfig.ResolveTarget[T](args, i.Workspace)
	if err != nil {
		i.Result.Error = err
		return
	}

	// we only expect zero or one target (depending on command)  - this should be enforced by Cobra
	i.Target = target

}

func validateModRequirementsRecursively(mod *modconfig.Mod, client *db_client.DbClient) []string {
	var validationErrors []string

	var pluginVersionMap = &modconfig.PluginVersionMap{
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
	for childDependencyName, childMod := range mod.ResourceMaps.Mods {
		if childDependencyName == "local" || mod.DependencyName == childMod.DependencyName {
			// this is a reference to self - skip (otherwise we will end up with a recursion loop)
			continue
		}
		childValidationErrors := validateModRequirementsRecursively(childMod, client)
		validationErrors = append(validationErrors, childValidationErrors...)
	}

	return validationErrors
}

func (i *InitData[T]) Cleanup(ctx context.Context) {
	if i.ShutdownTelemetry != nil {
		i.ShutdownTelemetry()
	}
	if i.Workspace != nil {
		i.Workspace.Close()
	}
}
