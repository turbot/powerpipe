package resources

import (
	"fmt"
	"github.com/turbot/pipe-fittings/modconfig"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/go-kit/types"
	typehelpers "github.com/turbot/go-kit/types"
	"github.com/turbot/pipe-fittings/cty_helpers"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/zclconf/go-cty/cty"
)

// Query is a struct representing the Query resource
type Query struct {
	modconfig.ResourceWithMetadataImpl
	QueryProviderImpl

	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`

	// only here as otherwise gocty.ImpliedType panics
	Unused string `cty:"unused" json:"-"`
}

func NewQuery(block *hcl.Block, mod *modconfig.Mod, shortName string) modconfig.HclResource {
	// queries cannot be anonymous
	return &Query{
		QueryProviderImpl: NewQueryProviderImpl(block, mod, shortName),
	}
}

func (q *Query) Equals(other *Query) bool {
	res := q.ShortName == other.ShortName &&
		q.FullName == other.FullName &&
		typehelpers.SafeString(q.Description) == typehelpers.SafeString(other.Description) &&
		typehelpers.SafeString(q.Documentation) == typehelpers.SafeString(other.Documentation) &&
		typehelpers.SafeString(q.SQL) == typehelpers.SafeString(other.SQL) &&
		typehelpers.SafeString(q.Title) == typehelpers.SafeString(other.Title)
	if !res {
		return res
	}

	// tags
	if q.Tags == nil {
		if other.Tags != nil {
			return false
		}
	} else {
		// we have tags
		if other.Tags == nil {
			return false
		}
		for k, v := range q.Tags {
			if otherVal, ok := (other.Tags)[k]; !ok && v != otherVal {
				return false
			}
		}
	}

	// params
	if len(q.Params) != len(other.Params) {
		return false
	}
	for i, p := range q.Params {
		if !p.Equals(other.Params[i]) {
			return false
		}
	}

	return true
}

func (q *Query) String() string {
	res := fmt.Sprintf(`
  -----
  Name: %s
  Title: %s
  Description: %s
  SQL: %s
`, q.FullName, types.SafeString(q.Title), types.SafeString(q.Description), types.SafeString(q.SQL))

	// add param defs if there are any
	if len(q.Params) > 0 {
		var paramDefsStr = make([]string, len(q.Params))
		for i, def := range q.Params {
			paramDefsStr[i] = def.String()
		}
		res += fmt.Sprintf("Params:\n\t%s\n  ", strings.Join(paramDefsStr, "\n\t"))
	}
	return res
}

// OnDecoded implements HclResource
func (q *Query) OnDecoded(*hcl.Block, modconfig.ModResourcesProvider) hcl.Diagnostics {
	return nil
}

// CtyValue implements CtyValueProvider
func (q *Query) CtyValue() (cty.Value, error) {
	return cty_helpers.GetCtyValue(q)
}

func (q *Query) Diff(other *Query) *modconfig.ModTreeItemDiffs {
	res := &modconfig.ModTreeItemDiffs{
		Item: q,
		Name: q.Name(),
	}

	if !utils.SafeStringsEqual(q.FullName, other.FullName) {
		res.AddPropertyDiff("Name")
	}

	res.PopulateChildDiffs(q, other)
	res.Merge(q.QueryProviderImpl.Diff(other))

	return res
}
