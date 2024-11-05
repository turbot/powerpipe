package parse

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/schema"
	"github.com/turbot/powerpipe/internal/resources"
)

// GetResourceSchema adds any app specific blocks to the existing resource schema
func GetResourceSchema(resource modconfig.HclResource, res *hcl.BodySchema) *hcl.BodySchema {
	// special cases for manually parsed attributes and blocks
	switch resource.GetBlockType() {
	case schema.BlockTypeDashboard, schema.BlockTypeContainer:
		res.Blocks = append(res.Blocks,
			hcl.BlockHeaderSchema{Type: schema.BlockTypeDashboard},
			hcl.BlockHeaderSchema{Type: schema.BlockTypeCard},
			hcl.BlockHeaderSchema{Type: schema.BlockTypeChart},
			hcl.BlockHeaderSchema{Type: schema.BlockTypeContainer},
			hcl.BlockHeaderSchema{Type: schema.BlockTypeDetection},
			hcl.BlockHeaderSchema{Type: schema.BlockTypeDetectionBenchmark},
			hcl.BlockHeaderSchema{Type: schema.BlockTypeFlow},
			hcl.BlockHeaderSchema{Type: schema.BlockTypeGraph},
			hcl.BlockHeaderSchema{Type: schema.BlockTypeHierarchy},
			hcl.BlockHeaderSchema{Type: schema.BlockTypeImage},
			hcl.BlockHeaderSchema{Type: schema.BlockTypeInput},
			hcl.BlockHeaderSchema{Type: schema.BlockTypeTable},
			hcl.BlockHeaderSchema{Type: schema.BlockTypeText},
			hcl.BlockHeaderSchema{Type: schema.BlockTypeWith},
		)
	case schema.BlockTypeDetectionBenchmark:
		res.Blocks = append(res.Blocks,
			hcl.BlockHeaderSchema{Type: schema.BlockTypeDetectionBenchmark},
			hcl.BlockHeaderSchema{Type: schema.BlockTypeDetection},
		)
	case schema.BlockTypeQuery:
		// remove `Query` from attributes
		var querySchema = &hcl.BodySchema{}
		for _, a := range res.Attributes {
			if a.Name != schema.AttributeTypeQuery {
				querySchema.Attributes = append(querySchema.Attributes, a)
			}
		}
		res = querySchema
	}

	if _, ok := resource.(resources.QueryProvider); ok {
		res.Blocks = append(res.Blocks, hcl.BlockHeaderSchema{Type: schema.BlockTypeParam})
		// if this is NOT query, add args
		if resource.GetBlockType() != schema.BlockTypeQuery {
			res.Attributes = append(res.Attributes, hcl.AttributeSchema{Name: schema.AttributeTypeArgs})
		}
	}
	if _, ok := resource.(resources.NodeAndEdgeProvider); ok {
		res.Blocks = append(res.Blocks,
			hcl.BlockHeaderSchema{Type: schema.BlockTypeCategory},
			hcl.BlockHeaderSchema{Type: schema.BlockTypeNode},
			hcl.BlockHeaderSchema{Type: schema.BlockTypeEdge})
	}
	if _, ok := resource.(resources.WithProvider); ok {
		res.Blocks = append(res.Blocks, hcl.BlockHeaderSchema{Type: schema.BlockTypeWith})
	}
	return res
}
