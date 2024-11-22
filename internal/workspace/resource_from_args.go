package workspace

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/schema"
	"github.com/turbot/pipe-fittings/sperr"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/pipe-fittings/workspace"
	pparse "github.com/turbot/powerpipe/internal/parse"
	"github.com/turbot/powerpipe/internal/resources"
)

// ResolveResourceAndArgsFromSQLString attempts to resolve 'arg' to a resource of type T and (optionally) query args
func ResolveResourceAndArgsFromSQLString[T modconfig.ModTreeItem](sqlString string, w *workspace.Workspace) (modconfig.ModTreeItem, *resources.QueryArgs, error) {
	var err error
	var empty T

	// 1) check if this is a resource
	// if this looks like a named query provider invocation, parse the sql string for arguments
	resource, args, err := extractResourceFromQueryString[T](sqlString, w)
	if err != nil {
		return empty, nil, err
	}

	if resource != nil {
		// success
		return resource, args, nil
	}

	// so we failed to resolve the resource from the input string
	// check whether it _looks_ like a resource name (i.e. mod.type.name OR type.name)
	if name, looksLikeResource := SqlLooksLikeExecutableResource(sqlString); looksLikeResource {
		return empty, nil, fmt.Errorf("'%s' not found in %s (%s)", name, w.Mod.Name(), w.Path)
	}
	switch any(empty).(type) {
	case *resources.Query:
		// if the desired type is a query,  and the sqlString DOES NOT look like a resource name,
		// treat it as a raw query and create a Query to wrap it
		q := createQueryResourceForCommandLineQuery(sqlString, w.Mod)

		// add to the workspace mod so the dashboard execution code can find it
		if err := w.Mod.AddResource(q); err != nil {
			return empty, nil, err
		}

		return q, nil, nil
	default:
		// failed to resolve
		return empty, nil, nil
	}

}

// does the input look like a resource which can be executed as a query
// Note: if anything fails just return nil values
func extractResourceFromQueryString[T modconfig.ModTreeItem](input string, w *workspace.Workspace) (modconfig.ModTreeItem, *resources.QueryArgs, error) {
	// can we extract a resource name from the string
	parsedResourceName, err := extractResourceNameFromQuery[T](input)
	if err != nil {
		return nil, nil, err
	}
	if parsedResourceName == nil {
		return nil, nil, nil
	}

	// TODO HACK
	if parsedResourceName.ItemType == "control_benchmark" {
		parsedResourceName.ItemType = "benchmark"
	}

	// ok we managed to extract a resource name - does this resource exist?
	resource, ok := w.GetResource(parsedResourceName)
	if !ok {
		return nil, nil, nil
	}

	// if the target is not the expected type, fail
	var target modconfig.ModTreeItem
	target, ok = resource.(T)
	if !ok {
		typeName := utils.GetGenericTypeName[T]()
		// TODO WHICH
		// // TODO HACK special case handling for detection benchmarks
		target, ok = resource.(*resources.DetectionBenchmark)
		if !ok {
			return nil, nil, sperr.New("target '%s' is not of the expected type '%s'", resource.GetUnqualifiedName(), typeName)
		}
		// // TODO HACK special case handling for detection benchmarks
		target, ok = resource.(*resources.Benchmark)
		if !ok {
			return nil, nil, sperr.New("target '%s' is not of the expected type '%s'", resource.GetUnqualifiedName(), typeName)
		}
	}

	_, args, err := pparse.ParseQueryInvocation(input)
	if err != nil {
		return nil, nil, err
	}

	// success
	return target, args, nil
}

// convert the given command line query into a query resource and add to workspace
// this is to allow us to use existing dashboard execution code
func createQueryResourceForCommandLineQuery(queryString string, mod *modconfig.Mod) *resources.Query {
	// build name
	shortName := "command_line_query"

	// this is NOT a named query - create the query using RawSql
	q := resources.NewQuery(&hcl.Block{Type: schema.BlockTypeQuery}, mod, shortName).(*resources.Query)
	q.SQL = utils.ToStringPointer(queryString)

	// add empty metadata
	q.SetMetadata(&modconfig.ResourceMetadata{})

	// return the new resource
	return q
}

// attempt top extra a resource name of the given type from the input string
// look at string up the the first open bracket
func extractResourceNameFromQuery[T modconfig.ModTreeItem](input string) (*modconfig.ParsedResourceName, error) {
	// convert the type T into a resource type name
	resourceType := resources.GenericTypeToBlockType[T]()
	// special case handling for variables
	if resourceType == schema.BlockTypeVariable {
		// variables are named var.xxxx, not variable.xxxx
		resourceType = schema.AttributeVar
	}

	// remove parameters from the input string before calling ParseResourceName
	// as parameters may break parsing
	openBracketIdx := strings.Index(input, "(")
	if openBracketIdx != -1 {
		input = input[:openBracketIdx]
	}

	parsedName, err := parseResourceName(input, resourceType)

	// if the typo eis query, do not bubble error up, just return nil parsed name
	// it is expected that this function may fail if a raw query is passed to it
	if err != nil && resourceType == schema.BlockTypeQuery {
		return nil, nil
	}

	return parsedName, err
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

	return parsed, nil
}

func SqlLooksLikeExecutableResource(input string) (string, bool) {
	// remove parameters from the input string before calling ParseResourceName
	// as parameters may break parsing
	openBracketIdx := strings.Index(input, "(")
	if openBracketIdx != -1 {
		input = input[:openBracketIdx]
	}
	parsedName, err := modconfig.ParseResourceName(input)
	if err == nil && helpers.StringSliceContains(schema.QueryProviderBlocks, parsedName.ItemType) {
		return parsedName.ToResourceName(), true
	}
	// do not bubble error up, just return false
	return "", false

}
