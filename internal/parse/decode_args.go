package parse

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/v2/hclhelpers"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/pipe-fittings/v2/parse"
	"github.com/turbot/pipe-fittings/v2/schema"
	"github.com/turbot/powerpipe/internal/resources"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

func DecodeArgs(attr *hcl.Attribute, evalCtx *hcl.EvalContext, resource resources.QueryProvider) (*resources.QueryArgs, []*resources.RuntimeDependency, hcl.Diagnostics) {
	var runtimeDependencies []*resources.RuntimeDependency
	var args = resources.NewQueryArgs()
	var diags hcl.Diagnostics

	v, valDiags := attr.Expr.Value(evalCtx)
	ty := v.Type()
	// determine which diags are runtime dependencies (which we allow) and which are not
	if valDiags.HasErrors() {
		for _, diag := range diags {
			dependency := parse.DiagsToDependency(diag)
			if dependency == nil || !dependency.IsRuntimeDependency() {
				diags = append(diags, diag)
			}
		}
	}
	// now diags contains all diags which are NOT runtime dependencies
	if diags.HasErrors() {
		return nil, nil, diags
	}

	var err error

	switch {
	case ty.IsObjectType():
		var argMap map[string]any
		argMap, runtimeDependencies, err = ctyObjectToArgMap(attr, v, evalCtx)
		if err == nil {
			err = args.SetArgMap(argMap)
		}
	case ty.IsTupleType():
		var argList []any
		argList, runtimeDependencies, err = ctyTupleToArgArray(attr, v)
		if err == nil {
			err = args.SetArgList(argList)
		}
	default:
		err = fmt.Errorf("'params' property must be either a map or an array")
	}

	if err != nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("%s has invalid parameter config", resource.Name()),
			Detail:   err.Error(),
			Subject:  &attr.Range,
		})
	}
	return args, runtimeDependencies, diags
}

func ctyTupleToArgArray(attr *hcl.Attribute, val cty.Value) ([]any, []*resources.RuntimeDependency, error) {
	// convert the attribute to a slice
	values := val.AsValueSlice()

	// build output array
	res := make([]any, len(values))
	var runtimeDependencies []*resources.RuntimeDependency

	for idx, v := range values {
		// if the value is unknown, this is a runtime dependency
		if !v.IsKnown() {
			runtimeDependency, err := identifyRuntimeDependenciesFromArray(attr, idx, schema.AttributeTypeArgs)
			if err != nil {
				return nil, nil, err
			}

			runtimeDependencies = append(runtimeDependencies, runtimeDependency)
		} else {
			// decode the value into a go type
			val, err := hclhelpers.CtyToGo(v)
			if err != nil {
				err := fmt.Errorf("invalid value provided for arg #%d: %v", idx, err)
				return nil, nil, err
			}

			res[idx] = val
		}
	}
	return res, runtimeDependencies, nil
}

func ctyObjectToArgMap(attr *hcl.Attribute, val cty.Value, evalCtx *hcl.EvalContext) (map[string]any, []*resources.RuntimeDependency, error) {
	res := make(map[string]any)
	var runtimeDependencies []*resources.RuntimeDependency
	it := val.ElementIterator()
	for it.Next() {
		k, v := it.Element()

		// decode key
		var key string
		if err := gocty.FromCtyValue(k, &key); err != nil {
			return nil, nil, err
		}

		// if the value is unknown, this is a runtime dependency
		if !v.IsKnown() {
			runtimeDependency, err := identifyRuntimeDependenciesFromObject(attr, key, schema.AttributeTypeArgs, evalCtx)
			if err != nil {
				return nil, nil, err
			}
			runtimeDependencies = append(runtimeDependencies, runtimeDependency)
		} else if getWrappedUnknownVal(v) {
			runtimeDependency, err := identifyRuntimeDependenciesFromObject(attr, key, schema.AttributeTypeArgs, evalCtx)
			if err != nil {
				return nil, nil, err
			}
			runtimeDependencies = append(runtimeDependencies, runtimeDependency)
		} else {
			// decode the value into a go type
			val, err := hclhelpers.CtyToGo(v)
			if err != nil {
				err := fmt.Errorf("invalid value provided for param '%s': %v", key, err)
				return nil, nil, err
			}
			res[key] = val
		}
	}

	return res, runtimeDependencies, nil
}

// TACTICAL - is the cty value an array with a single unknown value
func getWrappedUnknownVal(v cty.Value) bool {
	ty := v.Type()

	switch {

	case ty.IsTupleType():
		values := v.AsValueSlice()
		if len(values) == 1 && !values[0].IsKnown() {
			return true
		}
	}
	return false
}

func identifyRuntimeDependenciesFromObject(attr *hcl.Attribute, targetProperty, parentProperty string, evalCtx *hcl.EvalContext) (*resources.RuntimeDependency, error) {
	// find the expression for this key
	argsExpr, ok := attr.Expr.(*hclsyntax.ObjectConsExpr)
	if !ok {
		return nil, fmt.Errorf("could not extract runtime dependency for arg %s", targetProperty)
	}
	for _, item := range argsExpr.Items {
		nameCty, valDiags := item.KeyExpr.Value(evalCtx)
		if valDiags.HasErrors() {
			return nil, fmt.Errorf("could not extract runtime dependency for arg %s", targetProperty)
		}
		var name string
		if err := gocty.FromCtyValue(nameCty, &name); err != nil {
			return nil, err
		}
		if name == targetProperty {
			dep, err := getRuntimeDepFromExpression(item.ValueExpr, targetProperty, parentProperty)
			if err != nil {
				return nil, err
			}

			return dep, nil
		}
	}
	return nil, fmt.Errorf("could not extract runtime dependency for arg %s - not found in attribute map", targetProperty)
}

func getRuntimeDepFromExpression(expr hcl.Expression, targetProperty, parentProperty string) (*resources.RuntimeDependency, error) {
	isArray, propertyPath, err := modconfig.PropertyPathFromExpression(expr)
	if err != nil {
		return nil, err
	}

	if propertyPath.ItemType == schema.BlockTypeInput {
		// tactical: validate input dependency
		if err := validateInputRuntimeDependency(propertyPath); err != nil {
			return nil, err
		}
	}
	ret := &resources.RuntimeDependency{
		PropertyPath:       propertyPath,
		ParentPropertyName: parentProperty,
		TargetPropertyName: &targetProperty,
		IsArray:            isArray,
	}
	return ret, nil
}

func identifyRuntimeDependenciesFromArray(attr *hcl.Attribute, idx int, parentProperty string) (*resources.RuntimeDependency, error) {
	// find the expression for this key
	argsExpr, ok := attr.Expr.(*hclsyntax.TupleConsExpr)
	if !ok {
		return nil, fmt.Errorf("could not extract runtime dependency for arg #%d", idx)
	}
	for i, item := range argsExpr.Exprs {
		if i == idx {
			isArray, propertyPath, err := modconfig.PropertyPathFromExpression(item)
			if err != nil {
				return nil, err
			}
			// tactical: validate input dependency
			if propertyPath.ItemType == schema.BlockTypeInput {
				if err := validateInputRuntimeDependency(propertyPath); err != nil {
					return nil, err
				}
			}
			ret := &resources.RuntimeDependency{
				PropertyPath:        propertyPath,
				ParentPropertyName:  parentProperty,
				TargetPropertyIndex: &idx,
				IsArray:             isArray,
			}

			return ret, nil
		}
	}
	return nil, fmt.Errorf("could not extract runtime dependency for arg %d - not found in attribute list", idx)
}

// tactical - if runtime dependency is an input, validate it is of correct format
// TODO - include this with the main runtime dependency validation, when it is rewritten https://github.com/turbot/steampipe/issues/2925
func validateInputRuntimeDependency(propertyPath *modconfig.ParsedPropertyPath) error {
	// input references must be of form self.input.<input_name>.value
	if propertyPath.Scope != resources.RuntimeDependencyDashboardScope {
		return fmt.Errorf("could not resolve runtime dependency resource %s", propertyPath.Original)
	}
	return nil
}

func DecodeParam(block *hcl.Block, parseCtx *parse.ModParseContext) (*modconfig.ParamDef, []*resources.RuntimeDependency, hcl.Diagnostics) {
	def := modconfig.NewParamDef(block)
	var runtimeDependencies []*resources.RuntimeDependency
	content, diags := block.Body.Content(parse.ParamDefBlockSchema)

	if attr, exists := content.Attributes["description"]; exists {
		moreDiags := gohcl.DecodeExpression(attr.Expr, parseCtx.EvalCtx, &def.Description)
		diags = append(diags, moreDiags...)
	}
	if attr, exists := content.Attributes["default"]; exists {
		defaultValue, deps, moreDiags := decodeParamDefault(attr, parseCtx, def.UnqualifiedName)
		diags = append(diags, moreDiags...)
		if !helpers.IsNil(defaultValue) {
			err := def.SetDefault(defaultValue)
			if err != nil {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "invalid default config for " + def.UnqualifiedName,
					Detail:   err.Error(),
					Subject:  &attr.Range,
				})
				return nil, nil, diags
			}
		}
		runtimeDependencies = deps
	}
	return def, runtimeDependencies, diags
}

func decodeParamDefault(attr *hcl.Attribute, parseCtx *parse.ModParseContext, paramName string) (any, []*resources.RuntimeDependency, hcl.Diagnostics) {
	v, diags := attr.Expr.Value(parseCtx.EvalCtx)

	if v.IsKnown() {
		// convert the raw default into a string representation
		val, err := hclhelpers.CtyToGo(v)
		if err != nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("%s has invalid default config", paramName),
				Detail:   err.Error(),
				Subject:  &attr.Range,
			})
			return nil, nil, diags
		}
		return val, nil, nil
	}

	// so value not known - is there a runtime dependency?

	// check for a runtime dependency
	runtimeDependency, err := getRuntimeDepFromExpression(attr.Expr, "default", paramName)
	if err != nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("%s has invalid parameter default config", paramName),
			Detail:   err.Error(),
			Subject:  &attr.Range,
		})
		return nil, nil, diags
	}
	if runtimeDependency == nil {
		// return the original diags
		return nil, nil, diags
	}

	// so we have a runtime dependency
	return nil, []*resources.RuntimeDependency{runtimeDependency}, nil
}
