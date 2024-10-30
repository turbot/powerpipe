package resources

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/go-kit/helpers"
	typehelpers "github.com/turbot/go-kit/types"
	"github.com/turbot/pipe-fittings/cty_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/printers"
	"github.com/turbot/pipe-fittings/schema"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/zclconf/go-cty/cty"
)

type QueryProviderImpl struct {
	RuntimeDependencyProviderImpl
	QueryProviderRemain hcl.Body `hcl:",remain" json:"-"`

	SQL       *string               `cty:"sql" hcl:"sql" json:"sql,omitempty"`
	Query     *Query                `cty:"query" hcl:"query" json:"-"`
	Args      *QueryArgs            `cty:"args" json:"args,omitempty"`
	Params    []*modconfig.ParamDef `cty:"params" json:"params,omitempty"`
	QueryName *string               `json:"query,omitempty"`

	//nolint:unused // TODO: unused function
	withs               []*DashboardWith
	disableCtySerialise bool
	// flags to indicate if params and args were inherited from base resource
	argsInheritedFromBase   bool
	paramsInheritedFromBase bool
}

func NewQueryProviderImpl(block *hcl.Block, mod *modconfig.Mod, shortName string) QueryProviderImpl {
	return QueryProviderImpl{
		RuntimeDependencyProviderImpl: RuntimeDependencyProviderImpl{
			ModTreeItemImpl: modconfig.NewModTreeItemImpl(block, mod, shortName),
		},
	}
}

// GetParams implements QueryProvider
func (q *QueryProviderImpl) GetParams() []*modconfig.ParamDef {
	return q.Params
}

// GetArgs implements QueryProvider
func (q *QueryProviderImpl) GetArgs() *QueryArgs {
	return q.Args

}

// GetSQL implements QueryProvider
func (q *QueryProviderImpl) GetSQL() *string {
	return q.SQL
}

// GetQuery implements QueryProvider
func (q *QueryProviderImpl) GetQuery() *Query {
	return q.Query
}

// SetArgs implements QueryProvider
func (q *QueryProviderImpl) SetArgs(args *QueryArgs) {
	q.Args = args
}

// SetParams implements QueryProvider
func (q *QueryProviderImpl) SetParams(params []*modconfig.ParamDef) {
	q.Params = params
}

// ValidateQuery implements QueryProvider
// returns an error if neither sql or query are set
// it is overridden by resource types for which sql is optional
func (q *QueryProviderImpl) ValidateQuery() hcl.Diagnostics {
	var diags hcl.Diagnostics
	// Top level resources (with the exceptions of controls and queries) are never executed directly,
	// only used as base for a nested resource.
	// Therefore only nested resources, controls and queries MUST have sql or a query defined
	queryRequired := !q.IsTopLevel() ||
		helpers.StringSliceContains([]string{schema.BlockTypeQuery, schema.BlockTypeControl}, q.GetBlockType())

	if !queryRequired {
		return nil
	}

	if queryRequired && q.Query == nil && q.SQL == nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("%s does not define a query or SQL", q.Name()),
			Subject:  q.GetDeclRange(),
		})
	}
	return diags
}

// RequiresExecution implements QueryProvider
func (q *QueryProviderImpl) RequiresExecution(queryProvider QueryProvider) bool {
	return queryProvider.GetQuery() != nil || queryProvider.GetSQL() != nil
}

// GetResolvedQuery return the SQL and args to run the query
func (q *QueryProviderImpl) GetResolvedQuery(runtimeArgs *QueryArgs) (*ResolvedQuery, error) {
	argsArray, err := ResolveArgs(q, runtimeArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve args for %s: %s", q.Name(), err.Error())
	}
	sql := typehelpers.SafeString(q.GetSQL())
	// we expect there to be sql on the query provider, NOT a Query
	if sql == "" {
		return nil, fmt.Errorf("getResolvedQuery faiuled - no sql set for '%s'", q.Name())
	}

	return &ResolvedQuery{
		Name:       q.Name(),
		ExecuteSQL: sql,
		RawSQL:     sql,
		Args:       argsArray,
	}, nil
}

// MergeParentArgs merges our args with our parent args (ours take precedence)
func (q *QueryProviderImpl) MergeParentArgs(queryProvider QueryProvider, parent QueryProvider) (diags hcl.Diagnostics) {
	parentArgs := parent.GetArgs()
	if parentArgs == nil {
		return nil
	}

	args, err := parentArgs.Merge(queryProvider.GetArgs(), parent)
	if err != nil {
		return hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  err.Error(),
			Subject:  parent.(modconfig.HclResource).GetDeclRange(),
		}}
	}

	queryProvider.SetArgs(args)
	return nil
}

// GetQueryProviderImpl implements QueryProvider
func (q *QueryProviderImpl) GetQueryProviderImpl() *QueryProviderImpl {
	return q
}

// ParamsInheritedFromBase implements QueryProvider
// determine whether our params were inherited from base resource
func (q *QueryProviderImpl) ParamsInheritedFromBase() bool {
	return q.paramsInheritedFromBase
}

// ArgsInheritedFromBase implements QueryProvider
// determine whether our args were inherited from base resource
func (q *QueryProviderImpl) ArgsInheritedFromBase() bool {
	return q.argsInheritedFromBase
}

// CtyValue implements CtyValueProvider
func (q *QueryProviderImpl) CtyValue() (cty.Value, error) {
	if q.disableCtySerialise {
		return cty.Zero, nil
	}
	return cty_helpers.GetCtyValue(q)
}

func (q *QueryProviderImpl) SetBaseProperties() {
	q.RuntimeDependencyProviderImpl.SetBaseProperties()
	if q.SQL == nil {
		q.SQL = q.getBaseImpl().SQL
	}
	if q.Query == nil {
		q.Query = q.getBaseImpl().Query
	}
	if q.Args == nil {
		q.Args = q.getBaseImpl().Args
		q.argsInheritedFromBase = true
	}
	if q.Params == nil {
		q.Params = q.getBaseImpl().Params
		q.paramsInheritedFromBase = true
	}
}

func (q *QueryProviderImpl) getBaseImpl() *QueryProviderImpl {
	return q.GetBase().(QueryProvider).GetQueryProviderImpl()
}

func (q *QueryProviderImpl) OnDecoded(block *hcl.Block, _ modconfig.ModResourcesProvider) hcl.Diagnostics {
	q.populateQueryName()

	return nil
}

func (q *QueryProviderImpl) populateQueryName() {
	if q.Query != nil {
		q.QueryName = &q.Query.FullName
	}
}

// GetShowData implements printers.Showable
func (q *QueryProviderImpl) GetShowData() *printers.RowData {

	res := printers.NewRowData(
		printers.NewFieldValue("SQL", q.SQL),
		printers.NewFieldValue("Query", q.Query),
		printers.NewFieldValue("Args", q.Args),
		printers.NewFieldValue("Params", q.Params),
	)
	// merge fields from base, putting base fields first
	res.Merge(q.RuntimeDependencyProviderImpl.GetShowData())
	return res
}

func (q *QueryProviderImpl) Diff(other QueryProvider) *modconfig.ModTreeItemDiffs {
	d := &modconfig.ModTreeItemDiffs{
		Item: q,
		Name: q.Name(),
	}
	// sql
	if !utils.SafeStringsEqual(q.GetSQL(), other.GetSQL()) {
		d.AddPropertyDiff("SQL")
	}

	// args
	if lArgs := q.GetArgs(); lArgs == nil {
		if other.GetArgs() != nil {
			d.AddPropertyDiff("Args")
		}
	} else {
		// we have args
		if rArgs := other.GetArgs(); rArgs == nil {
			d.AddPropertyDiff("Args")
		} else if !lArgs.Equals(rArgs) {
			d.AddPropertyDiff("Args")
		}
	}

	// query
	if lQuery := q.GetQuery(); lQuery == nil {
		if other.GetQuery() != nil {
			d.AddPropertyDiff("Query")
		}
	} else {
		// we have query
		if rQuery := other.GetQuery(); rQuery == nil {
			d.AddPropertyDiff("Query")
		} else if !lQuery.Equals(rQuery) {
			d.AddPropertyDiff("Query")
		}
	}

	// params
	lParams := q.GetParams()
	rParams := other.GetParams()
	if len(lParams) != len(rParams) {
		d.AddPropertyDiff("Params")
	} else {
		for i, lParam := range lParams {
			if !lParam.Equals(rParams[i]) {
				d.AddPropertyDiff("Params")
			}
		}
	}

	// with
	if lwp, ok := any(q).(WithProvider); ok {
		rwp := other.(WithProvider)
		lWiths := lwp.GetWiths()
		rWiths := rwp.GetWiths()
		if len(lWiths) != len(rWiths) {
			d.AddPropertyDiff("With")
		} else {
			for i, lWith := range lWiths {
				if !lWith.Equals(rWiths[i]) {
					d.AddPropertyDiff("With")
				}
			}
		}

		// have BASE withs changed
		lbase := q.GetBase()
		rbase := other.GetBase()
		var lbaseWiths []*DashboardWith
		var rbaseWiths []*DashboardWith
		if lbase != nil {
			lbaseWiths = lbase.(WithProvider).GetWiths()
		}
		if rbase != nil {
			rbaseWiths = rbase.(WithProvider).GetWiths()
		}
		if len(lbaseWiths) != len(rbaseWiths) {
			d.AddPropertyDiff("With")
		} else {
			for i, lBaseWith := range lbaseWiths {
				if !lBaseWith.Equals(rbaseWiths[i]) {
					d.AddPropertyDiff("With")
				}
			}
		}
	}

	return d
}

func (b *QueryProviderImpl) GetNestedStructs() []modconfig.CtyValueProvider {
	// return all nested structs - this is used to get the nested structs for the cty serialisation
	return append([]modconfig.CtyValueProvider{&b.RuntimeDependencyProviderImpl}, b.RuntimeDependencyProviderImpl.GetNestedStructs()...)
}
