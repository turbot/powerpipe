package dashboardexecute

import (
	"fmt"
	"github.com/turbot/pipe-fittings/modconfig/powerpipe"
	"strings"

	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/powerpipe/internal/dashboardtypes"
	"github.com/turbot/powerpipe/internal/workspace"
)

// GetReferencedVariables builds map of variables values containing only those mod variables which are referenced
// NOTE: we refer to variables in dependency mods in the format which is valid for an SPVARS filer, i.e.
// <mod>.<var-name>
// the VariableValues map will contain these variables with the name format <mod>.var.<var-name>,
// so we must convert the name
func GetReferencedVariables(root dashboardtypes.DashboardTreeRun, w *workspace.WorkspaceEvents) (map[string]string, error) {
	var referencedVariables = make(map[string]string)

	addReferencedVars := func(refs []*modconfig.ResourceReference) {
		for _, ref := range refs {
			parts := strings.Split(ref.To, ".")
			if len(parts) == 2 && parts[0] == "var" {
				varName := parts[1]
				varValueName := varName
				// NOTE: if the ref is NOT for the workspace mod, then use the qualified variable name
				// (e.g. aws_insights.var.v1)
				if refMod := ref.GetMetadata().ModName; refMod != w.Mod.ShortName {
					varValueName = fmt.Sprintf("%s.var.%s", refMod, varName)
					varName = fmt.Sprintf("%s.%s", refMod, varName)
				}
				referencedVariables[varName] = w.VariableValues[varValueName]
			}
		}
	}

	switch r := root.(type) {
	case *DashboardRun:

		err := r.dashboard.WalkResources(
			func(resource modconfig.HclResource) (bool, error) {
				if resourceWithMetadata, ok := resource.(modconfig.ResourceWithMetadata); ok {
					addReferencedVars(resourceWithMetadata.GetReferences())
				}
				return true, nil
			},
		)
		if err != nil {
			return nil, err
		}
	case *CheckRun:
		switch n := r.resource.(type) {
		case *powerpipe.Benchmark:
			err := n.WalkResources(
				func(resource modconfig.ModTreeItem) (bool, error) {
					if resourceWithMetadata, ok := resource.(modconfig.ResourceWithMetadata); ok {
						addReferencedVars(resourceWithMetadata.GetReferences())
					}
					return true, nil
				},
			)
			if err != nil {
				return nil, err
			}
		case *powerpipe.Control:
			addReferencedVars(n.GetReferences())
		}
	}

	return referencedVariables, nil
}
