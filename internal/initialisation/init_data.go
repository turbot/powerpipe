package initialisation

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/app_specific"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/export"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/modinstaller"
	"github.com/turbot/pipe-fittings/statushooks"
	"github.com/turbot/pipe-fittings/workspace"
	"github.com/turbot/powerpipe/internal/cmdconfig"
	"github.com/turbot/powerpipe/internal/dashboardworkspace"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
	"github.com/turbot/steampipe-plugin-sdk/v5/telemetry"
)

type InitData struct {
	Workspace       *workspace.Workspace
	WorkspaceEvents *dashboardworkspace.WorkspaceEvents
	Result          *InitResult

	ShutdownTelemetry func()
	ExportManager     *export.Manager
	Target            modconfig.ModTreeItem
	QueryArgs         map[string]*modconfig.QueryArgs
}

func NewErrorInitData(err error) *InitData {
	return &InitData{
		Result: &InitResult{
			ErrorAndWarnings: error_helpers.NewErrorsAndWarning(err),
		},
	}
}

func NewInitData(ctx context.Context, targetType string, targetNames ...string) *InitData {
	modLocation := viper.GetString(constants.ArgModLocation)

	w, errAndWarnings := workspace.LoadWorkspacePromptingForVariables(ctx, modLocation)
	if errAndWarnings.GetError() != nil {
		return NewErrorInitData(fmt.Errorf("failed to load workspace: %s", error_helpers.HandleCancelError(errAndWarnings.GetError()).Error()))
	}

	i := &InitData{
		Result:        &InitResult{},
		ExportManager: export.NewManager(),
	}

	i.Workspace = w
	i.Result.Warnings = errAndWarnings.Warnings

	// now do the actual initialisation
	i.Init(ctx, targetType, targetNames...)

	return i
}

func (i *InitData) RegisterExporters(exporters ...export.Exporter) error {
	for _, e := range exporters {
		if err := i.ExportManager.Register(e); err != nil {
			return err
		}
	}

	return nil
}

func (i *InitData) Init(ctx context.Context, targetType string, args ...string) {
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

	i.resolveTarget(args, targetType)
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

	// validate mod requirements
	// For now we only validate CLI version
	// TODO add validation for required plugins for steampipe backend
	validationWarnings := validateModRequirementsRecursively(i.Workspace.Mod)
	i.Result.AddWarnings(validationWarnings...)

}

// resolve target resource, args and any target specific search path
func (i *InitData) resolveTarget(args []string, targetType string) {
	// resolve target resources
	targets, queryArgs, err := cmdconfig.ResolveTargets(args, targetType, i.Workspace)
	if err != nil {
		i.Result.Error = err
		return
	}
	i.QueryArgs = queryArgs

	if len(targets) == 0 {
		// no targets found
		return
	}

	// we only expect a single target - this should be enforced by Cobra
	if len(targets) > 1 {
		i.Result.Error = sperr.New("expected a single execution target, got %d", len(targets))
		return
	}
	i.Target = targets[0]

}

func validateModRequirementsRecursively(mod *modconfig.Mod) []string {
	var validationErrors []string

	// validate this mod
	for _, err := range mod.ValidateRequirements() {
		validationErrors = append(validationErrors, err.Error())
	}

	// validate dependent mods
	for childDependencyName, childMod := range mod.ResourceMaps.Mods {
		if childDependencyName == "local" || mod.DependencyName == childMod.DependencyName {
			// this is a reference to self - skip (otherwise we will end up with a recursion loop)
			continue
		}
		childValidationErrors := validateModRequirementsRecursively(childMod)
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
}
