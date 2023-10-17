package parse

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/turbot/go-kit/hcl_helpers"
)

// AddReferences populates the 'References' resource field, used for the introspection tables
func AddReferences(resource entities.HclResource, block *hcl.Block, parseCtx *ModParseContext) hcl.Diagnostics {
	resourceWithMetadata, ok := resource.(entities.ResourceWithMetadata)
	if !ok {
		return nil
	}

	var diags hcl.Diagnostics
	for _, attr := range block.Body.(*hclsyntax.Body).Attributes {
		for _, v := range attr.Expr.Variables() {
			for _, referenceBlockType := range entities.ReferenceBlocks {
				if referenceString, ok := hcl_helpers.ResourceNameFromTraversal(referenceBlockType, v); ok {
					var blockName string
					if len(block.Labels) > 0 {
						blockName = block.Labels[0]
					}
					reference := entities.NewResourceReference(resource, block, referenceString, blockName, attr)

					moreDiags := addResourceMetadata(reference, attr.SrcRange, parseCtx)
					diags = append(diags, moreDiags...)
					resourceWithMetadata.AddReference(reference)
					break
				}
			}
		}
	}
	return diags
}
