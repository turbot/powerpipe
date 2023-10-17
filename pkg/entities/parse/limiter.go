package parse

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
)

func DecodeLimiter(block *hcl.Block) (*entities.RateLimiter, hcl.Diagnostics) {
	var limiter = &entities.RateLimiter{
		// populate name from label
		Name: block.Labels[0],
	}
	diags := gohcl.DecodeBody(block.Body, nil, limiter)
	if !diags.HasErrors() {
		limiter.OnDecoded(block)
	}

	return limiter, diags
}
