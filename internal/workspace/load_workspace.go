package workspace

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/turbot/pipe-fittings/load_mod"
	"github.com/turbot/powerpipe/pkg/entities/parse"

	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/inputvars"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/statushooks"
	shared_workspace "github.com/turbot/pipe-fittings/workspace"
	"github.com/turbot/terraform-components/terraform"
)

func LoadWorkspacePromptingForVariables(ctx context.Context) (*shared_workspace.Workspace, *error_helpers.ErrorAndWarnings) {
	workspacePath := viper.GetString(constants.ArgModLocation)
	t := time.Now()
	defer func() {
		log.Printf("[TRACE] Workspace load took %dms\n", time.Since(t).Milliseconds())
	}()
	w, errAndWarnings := shared_workspace.Load(ctx, workspacePath)
	if errAndWarnings.GetError() == nil {
		return w, errAndWarnings
	}
	missingVariablesError, ok := errAndWarnings.GetError().(*load_mod.MissingVariableError)
	// if there was an error which is NOT a MissingVariableError, return it
	if !ok {
		return nil, errAndWarnings
	}
	// if there are missing transitive dependency variables, fail as we do not prompt for these
	if len(missingVariablesError.MissingTransitiveVariables) > 0 {
		return nil, errAndWarnings
	}
	// if interactive input is disabled, return the missing variables error
	if !viper.GetBool(constants.ArgInput) {
		return nil, error_helpers.NewErrorsAndWarning(missingVariablesError)
	}
	// so we have missing variables - prompt for them
	// first hide spinner if it is there
	statushooks.Done(ctx)
	if err := promptForMissingVariables(ctx, missingVariablesError.MissingVariables, workspacePath); err != nil {
		log.Printf("[TRACE] Interactive variables prompting returned error %v", err)
		return nil, error_helpers.NewErrorsAndWarning(err)
	}
	// ok we should have all variables now - reload workspace
	return shared_workspace.Load(ctx, workspacePath)
}

func promptForMissingVariables(ctx context.Context, missingVariables []*modconfig.Variable, workspacePath string) error {
	fmt.Println()
	fmt.Println("Variables defined with no value set.")
	for _, v := range missingVariables {
		variableName := v.ShortName
		variableDisplayName := fmt.Sprintf("var.%s", v.ShortName)
		// if this variable is NOT part of the workspace mod, add the mod name to the variable name
		if v.Mod.ModPath != workspacePath {
			variableDisplayName = fmt.Sprintf("%s.var.%s", v.ModName, v.ShortName)
			variableName = fmt.Sprintf("%s.%s", v.ModName, v.ShortName)
		}
		r, err := promptForVariable(ctx, variableDisplayName, v.GetDescription())
		if err != nil {
			return err
		}
		addInteractiveVariableToViper(variableName, r)
	}
	return nil
}

func promptForVariable(ctx context.Context, name, description string) (string, error) {
	uiInput := &inputvars.UIInput{}
	rawValue, err := uiInput.Input(ctx, &terraform.InputOpts{
		Id:          name,
		Query:       name,
		Description: description,
	})

	return rawValue, err
}

func addInteractiveVariableToViper(name string, rawValue string) {
	varMap := viper.GetStringMap(constants.ConfigInteractiveVariables)
	varMap[name] = rawValue
	viper.Set(constants.ConfigInteractiveVariables, varMap)
}

func CreateWorkspaceMod(ctx context.Context, workspacePath string) (*modconfig.Mod, error) {
	if parse.ModfileExists(workspacePath) {
		fmt.Println("Working folder already contains a mod definition file")
		return nil, nil
	}
	mod := modconfig.CreateDefaultMod(workspacePath)
	if err := mod.Save(); err != nil {
		return nil, err
	}

	// load up the written mod file so that we get the updated
	// block ranges
	mod, err := parse.LoadModfile(workspacePath)
	if err != nil {
		return nil, err
	}

	return mod, nil
}
