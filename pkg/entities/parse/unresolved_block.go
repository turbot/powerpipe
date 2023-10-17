package parse

import (
	"fmt"
	"github.com/turbot/go-kit/hcl_helpers"
	"github.com/turbot/powerpipe/pkg/entities"
	"strings"

	"github.com/hashicorp/hcl/v2"
)

type unresolvedBlock struct {
	Name         string
	Block        *hcl.Block
	DeclRange    hcl.Range
	Dependencies map[string]*entities.ResourceDependency
}

func newUnresolvedBlock(block *hcl.Block, name string, dependencies map[string]*entities.ResourceDependency) *unresolvedBlock {
	return &unresolvedBlock{
		Name:         name,
		Block:        block,
		Dependencies: dependencies,
		DeclRange:    hcl_helpers.BlockRange(block),
	}
}

func (b unresolvedBlock) String() string {
	depStrings := make([]string, len(b.Dependencies))
	idx := 0
	for _, dep := range b.Dependencies {
		depStrings[idx] = fmt.Sprintf(`%s -> %s`, b.Name, dep.String())
		idx++
	}
	return strings.Join(depStrings, "\n")
}
