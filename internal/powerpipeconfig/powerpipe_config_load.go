package powerpipeconfig

import (
	"github.com/hashicorp/hcl/v2"
	filehelpers "github.com/turbot/go-kit/files"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/v2/app_specific"
	"github.com/turbot/pipe-fittings/v2/error_helpers"
	"github.com/turbot/pipe-fittings/v2/filepaths"
	"github.com/turbot/pipe-fittings/v2/parse"
	"github.com/turbot/pipe-fittings/v2/schema"
	"log/slog"
)

var GlobalConfig *PowerpipeConfig

type loadConfigOptions struct {
	include []string
}

func LoadPowerpipeConfig(configPaths ...string) (res *PowerpipeConfig, errorsAndWarnings error_helpers.ErrorAndWarnings) {
	errorsAndWarnings = error_helpers.NewErrorsAndWarning(nil)
	defer func() {
		if r := recover(); r != nil {
			errorsAndWarnings = error_helpers.NewErrorsAndWarning(helpers.ToError(r))
		}
	}()

	connectionConfigExtensions := []string{app_specific.ConfigExtension}

	include := filehelpers.InclusionsFromExtensions(connectionConfigExtensions)
	loadOptions := &loadConfigOptions{include: include}

	res = NewPowerpipeConfig()

	lastErrorLength := 0

	for {
		var diags hcl.Diagnostics
		for i := len(configPaths) - 1; i >= 0; i-- {
			configPath := configPaths[i]
			moreDiags := res.loadPowerpipeConfigBlocks(configPath, loadOptions)
			if len(moreDiags) > 0 {
				diags = append(diags, moreDiags...)
			}
		}

		if len(diags) == 0 {
			break
		}

		if len(diags) > 0 && lastErrorLength == len(diags) {
			return nil, error_helpers.DiagsToErrorsAndWarnings("Failed to load Powerpipe config", diags)
		}

		lastErrorLength = len(diags)
	}

	return res, errorsAndWarnings
}

func (c *PowerpipeConfig) loadPowerpipeConfigBlocks(configPath string, opts *loadConfigOptions) hcl.Diagnostics {
	configPaths, err := filehelpers.ListFiles(configPath, &filehelpers.ListOptions{
		Flags:   filehelpers.FilesFlat,
		Include: opts.include,
		Exclude: []string{filepaths.WorkspaceLockFileName},
	})

	if err != nil {
		slog.Warn("failed to get config file paths", "error", err)
		return hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "failed to get config file paths",
				Detail:   err.Error(),
			},
		}
	}

	if len(configPaths) == 0 {
		return hcl.Diagnostics{}
	}

	fileData, diags := parse.LoadFileData(configPaths...)
	if diags.HasErrors() {
		slog.Warn("failed to load all config files", "error", err)
		return diags
	}

	body, diags := parse.ParseHclFiles(fileData)
	if diags.HasErrors() {
		return diags
	}

	// do a partial decode
	content, diags := body.Content(parse.PowerpipeConfigBlockSchema)
	if diags.HasErrors() {
		return diags
	}

	for _, block := range content.Blocks {
		switch block.Type {
		case schema.BlockTypeConnection:

			conn, moreDiags := parse.DecodePipelingConnection(configPath, block)
			if len(moreDiags) > 0 {
				diags = append(diags, moreDiags...)
				slog.Debug("failed to decode connection block")
				continue
			}
			c.PipelingConnections[conn.Name()] = conn
		}
	}

	if len(diags) > 0 {
		return diags
	}

	return diags
}
