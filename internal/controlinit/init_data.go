package controlinit

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/backend"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/statushooks"
	"github.com/turbot/pipe-fittings/workspace"
	localcmdconfig "github.com/turbot/powerpipe/internal/cmdconfig"
	"github.com/turbot/powerpipe/internal/controldisplay"
	"github.com/turbot/powerpipe/internal/db_client"
	"github.com/turbot/powerpipe/internal/initialisation"
)

type CheckTarget interface {
	*modconfig.Benchmark | *modconfig.Control
}

type InitData struct {
	initialisation.InitData
	OutputFormatter controldisplay.Formatter
	ControlFilter   workspace.ResourceFilter
	Client          *db_client.DbClient
}

// NewInitData returns a new InitData object
// It also starts an asynchronous population of the object
// InitData.Done closes after asynchronous initialization completes
func NewInitData[T CheckTarget](ctx context.Context, args []string) *InitData {
	typeName := localcmdconfig.GetGenericTypeName[T]()

	statushooks.SetStatus(ctx, "Loading workspace")

	initData := initialisation.NewInitData(ctx, typeName, args...)

	// create InitData, but do not initialize yet, since 'viper' is not completely setup
	i := &InitData{
		InitData: *initData,
	}
	if i.Result.Error != nil {
		return i
	}

	w := i.Workspace
	if !w.ModfileExists() {
		i.Result.Error = workspace.ErrorNoModDefinition
	}

	if viper.GetString(constants.ArgOutput) == constants.OutputFormatNone {
		// set progress to false
		viper.Set(constants.ArgProgress, false)
	}
	// set color schema
	err := initialiseCheckColorScheme()
	if err != nil {
		i.Result.Error = err
		return i
	}

	if len(w.GetResourceMaps().Controls)+len(w.GetResourceMaps().Benchmarks) == 0 {
		i.Result.AddWarnings("no controls or benchmarks found in current workspace")
	}

	if err := controldisplay.EnsureTemplates(); err != nil {
		i.Result.Error = err
		return i
	}

	if len(viper.GetStringSlice(constants.ArgExport)) > 0 {
		if err := i.registerCheckExporters(); err != nil {
			i.Result.Error = err
			return i
		}

		// validate required export formats
		if err := i.ExportManager.ValidateExportFormat(viper.GetStringSlice(constants.ArgExport)); err != nil {
			i.Result.Error = err
			return i
		}
	}

	output := viper.GetString(constants.ArgOutput)
	formatter, err := parseOutputArg(output)
	if err != nil {
		i.Result.Error = err
		return i
	}
	i.OutputFormatter = formatter

	i.setControlFilter()

	// set the default database and search patch config, based on the target resource
	defaultSearchPathConfig, defaultDatabase := db_client.GetDefaultDatabaseConfig()
	database, searchPathConfig, err := db_client.GetDatabaseConfigForResource(initData.Target, initData.Workspace.Mod, defaultDatabase, defaultSearchPathConfig)
	if err != nil {
		i.Result.Error = err
		return i
	}

	// create client
	var opts []backend.ConnectOption
	if !searchPathConfig.Empty() {
		opts = append(opts, backend.WithSearchPathConfig(searchPathConfig))
	}
	client, err := db_client.NewDbClient(ctx, database, opts...)
	if err != nil {
		i.Result.Error = err
		return i
	}

	i.Client = client

	return i
}

func (i *InitData) setControlFilter() {
	if viper.IsSet(constants.ArgTag) {
		// if '--tag' args were used, derive the whereClause from them
		tags := viper.GetStringSlice(constants.ArgTag)
		i.ControlFilter = workspace.ResourceFilterFromTags(tags)
	} else if viper.IsSet(constants.ArgWhere) {
		// if a 'where' arg was used, execute this sql to get a list of  control names
		// use this list to build a name map used to determine whether to run a particular control
		i.ControlFilter = workspace.ResourceFilter{
			Where: viper.GetString(constants.ArgWhere),
		}
	}
}

// register exporters for each of the supported check formats
func (i *InitData) registerCheckExporters() error {
	exporters, err := controldisplay.GetExporters()
	error_helpers.FailOnErrorWithMessage(err, "failed to load exporters")

	// register all exporters
	return i.RegisterExporters(exporters...)
}

// parseOutputArg parses the --output flag value and returns the Formatter that can format the data
func parseOutputArg(arg string) (formatter controldisplay.Formatter, err error) {
	formatResolver, err := controldisplay.NewFormatResolver()
	if err != nil {
		return nil, err
	}

	return formatResolver.GetFormatter(arg)
}

func initialiseCheckColorScheme() error {
	// TODO kai remove themes and use standard color codes
	theme := "dark"
	if !viper.GetBool(constants.ConfigKeyIsTerminalTTY) {
		// enforce plain output for non-terminals
		theme = "plain"
	}
	themeDef, ok := controldisplay.ColorSchemes[theme]
	if !ok {
		return fmt.Errorf("invalid theme '%s'", theme)
	}
	scheme, err := controldisplay.NewControlColorScheme(themeDef)
	if err != nil {
		return err
	}
	controldisplay.ControlColors = scheme
	return nil
}
