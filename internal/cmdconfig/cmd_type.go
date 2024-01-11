package cmdconfig

import (
	"fmt"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/workspace"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
	"strings"
)

// GetGenericTypeName returns lower case form of type unqualified name
func GetGenericTypeName[T any]() string {
	longName := fmt.Sprintf("%T", *new(T))
	split := strings.Split(longName, ".")
	return strings.ToLower(split[len(split)-1])
}

func ResolveTargetArgs(args []string, commandTargetType string, w *workspace.Workspace) ([]modconfig.ModTreeItem, error) {
	var targets []modconfig.ModTreeItem
	for _, targetName := range args {
		target, err := resolveResourceName(targetName, commandTargetType, w)
		if err != nil {
			return nil, err
		}
		targets = append(targets, target)
	}
	return targets, nil
}

// resolveResourceName parses targetName to verify it is a named resource
// - if no resource type is specified in the name, it is added from the command type
// - validate the resource type specified in the name matches the command type
// - verify the resource exists in the workspace
func resolveResourceName(targetName string, commandTargetType string, w *workspace.Workspace) (modconfig.ModTreeItem, error) {

	parsed, err := parseResourceName(targetName, commandTargetType)
	if err != nil {
		return nil, err
	}
	resource, found := w.GetResource(parsed)
	if !found {
		return nil, fmt.Errorf("target resource %s does not exist in mod %s", parsed.ToFullName(), w.Mod.Name())
	}
	// 	safe cast to be on the safe side
	modTreeItem, ok := resource.(modconfig.ModTreeItem)
	if !ok {
		// not expected! this is a coding error, we should never get here
		return nil, fmt.Errorf("target resource %s is not a mod tree item", parsed.ToFullName())
	}
	return modTreeItem, nil
}

func parseResourceName(targetName string, commandTargetType string) (*modconfig.ParsedResourceName, error) {
	parsed := &modconfig.ParsedResourceName{}
	parts := strings.Split(targetName, ".")

	switch len(parts) {
	case 0:
		return nil, sperr.New("empty name passed to resolveResourceName")
	case 1:
		// if no type was specified, deduce the type from the check command used
		parsed.Name = parts[0]
		parsed.ItemType = commandTargetType
	case 2:
		parsed.ItemType = parts[0]
		parsed.Name = parts[1]
	case 3:
		parsed.Mod = parts[0]
		parsed.ItemType = parts[1]
		parsed.Name = parts[2]
	default:
		return nil, sperr.New("invalid name passed to ParseResourceName")
	}

	// now validate the resource type matches the commandTargetType
	if parsed.ItemType == "" {
		parsed.ItemType = commandTargetType
	}
	if parsed.ItemType != commandTargetType {
		return nil, sperr.New(fmt.Sprintf("invalid resource type %s - expected %s", parsed.ItemType, commandTargetType))
	}

	return parsed, nil
}
