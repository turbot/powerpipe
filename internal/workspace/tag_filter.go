package workspace

import (
	"strings"

	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/pipe-fittings/v2/workspace"
)

// tagFilter captures allowed and disallowed values for a tag key.
type tagFilter struct {
	equals    []string
	notEquals []string
}

// ResourceFilterFromTagArgs builds a workspace.ResourceFilter from CLI --tag args,
// supporting both equality (key=value) and inequality (key!=value) semantics.
// For inequality, resources without the tag must be included (only resources
// with the tag set to the disallowed value are excluded). This aligns CLI
// behavior with user expectations (e.g., env!=staging should include items
// with no env tag).
func ResourceFilterFromTagArgs(tags []string) workspace.ResourceFilter {
	tagFilters := map[string]*tagFilter{}

	for _, arg := range tags {
		var key, value string
		var isNotEquals bool

		if parts := strings.SplitN(arg, "!=", 2); len(parts) == 2 {
			key, value = parts[0], parts[1]
			isNotEquals = true
		} else if parts := strings.SplitN(arg, "=", 2); len(parts) == 2 {
			key, value = parts[0], parts[1]
		} else {
			// malformed tag; skip and let validation elsewhere handle errors (consistent with prior behavior)
			continue
		}

		tf := tagFilters[key]
		if tf == nil {
			tf = &tagFilter{}
			tagFilters[key] = tf
		}
		if isNotEquals {
			tf.notEquals = append(tf.notEquals, value)
		} else {
			tf.equals = append(tf.equals, value)
		}
	}

	return workspace.ResourceFilter{
		WherePredicate: func(item modconfig.HclResource) bool {
			itemTags := item.GetTags()

			for key, filter := range tagFilters {
				tagValue, present := itemTags[key]

				// Equality rules: tag must be present AND match one of the allowed values.
				if len(filter.equals) > 0 {
					if !present || !contains(filter.equals, tagValue) {
						return false
					}
				}

				// Inequality rules: exclude only if tag present AND value is disallowed.
				if len(filter.notEquals) > 0 && present && contains(filter.notEquals, tagValue) {
					return false
				}
			}

			return true
		},
	}
}

func contains(values []string, candidate string) bool {
	for _, v := range values {
		if v == candidate {
			return true
		}
	}
	return false
}
