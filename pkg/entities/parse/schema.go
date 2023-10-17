package parse

import (
	"github.com/hashicorp/hcl/v2"
)

// cache resource schemas
var resourceSchemaCache = make(map[string]*hcl.BodySchema)

var ConfigBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{},
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       entities.BlockTypeConnection,
			LabelNames: []string{"name"},
		},
		{
			Type:       entities.BlockTypePlugin,
			LabelNames: []string{"name"},
		},
		{
			Type:       entities.BlockTypeOptions,
			LabelNames: []string{"type"},
		},
		{
			Type:       entities.BlockTypeWorkspaceProfile,
			LabelNames: []string{"name"},
		},
	},
}
var PluginBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{},
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       entities.BlockTypeRateLimiter,
			LabelNames: []string{"name"},
		},
	},
}

var WorkspaceProfileBlockSchema = &hcl.BodySchema{

	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "options",
			LabelNames: []string{"type"},
		},
	},
}

var ConnectionBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{
			Name:     "plugin",
			Required: true,
		},
		{
			Name: "type",
		},
		{
			Name: "connections",
		},
		{
			Name: "import_schema",
		},
	},
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "options",
			LabelNames: []string{"type"},
		},
	},
}

// WorkspaceBlockSchema is the top level schema for all workspace resources
var WorkspaceBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{},
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       string(entities.BlockTypeMod),
			LabelNames: []string{"name"},
		},
		{
			Type:       entities.BlockTypeVariable,
			LabelNames: []string{"name"},
		},
		{
			Type:       entities.BlockTypeQuery,
			LabelNames: []string{"name"},
		},
		{
			Type:       entities.BlockTypeControl,
			LabelNames: []string{"name"},
		},
		{
			Type:       entities.BlockTypeBenchmark,
			LabelNames: []string{"name"},
		},
		{
			Type:       entities.BlockTypeDashboard,
			LabelNames: []string{"name"},
		},
		{
			Type:       entities.BlockTypeCard,
			LabelNames: []string{"name"},
		},
		{
			Type:       entities.BlockTypeChart,
			LabelNames: []string{"name"},
		},
		{
			Type:       entities.BlockTypeFlow,
			LabelNames: []string{"name"},
		},
		{
			Type:       entities.BlockTypeGraph,
			LabelNames: []string{"name"},
		},
		{
			Type:       entities.BlockTypeHierarchy,
			LabelNames: []string{"name"},
		},
		{
			Type:       entities.BlockTypeImage,
			LabelNames: []string{"name"},
		},
		{
			Type:       entities.BlockTypeInput,
			LabelNames: []string{"name"},
		},
		{
			Type:       entities.BlockTypeTable,
			LabelNames: []string{"name"},
		},
		{
			Type:       entities.BlockTypeText,
			LabelNames: []string{"name"},
		},
		{
			Type:       entities.BlockTypeNode,
			LabelNames: []string{"name"},
		},
		{
			Type:       entities.BlockTypeEdge,
			LabelNames: []string{"name"},
		},
		{
			Type: entities.BlockTypeLocals,
		},
		{
			Type:       entities.BlockTypeCategory,
			LabelNames: []string{"name"},
		},
	},
}

// DashboardBlockSchema is only used to validate the blocks of a Dashboard
var DashboardBlockSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       entities.BlockTypeInput,
			LabelNames: []string{"name"},
		},
		{
			Type:       entities.BlockTypeParam,
			LabelNames: []string{"name"},
		},
		{
			Type: entities.BlockTypeWith,
		},
		{
			Type: entities.BlockTypeContainer,
		},
		{
			Type: entities.BlockTypeCard,
		},
		{
			Type: entities.BlockTypeChart,
		},
		{
			Type: entities.BlockTypeBenchmark,
		},
		{
			Type: entities.BlockTypeControl,
		},
		{
			Type: entities.BlockTypeFlow,
		},
		{
			Type: entities.BlockTypeGraph,
		},
		{
			Type: entities.BlockTypeHierarchy,
		},
		{
			Type: entities.BlockTypeImage,
		},
		{
			Type: entities.BlockTypeTable,
		},
		{
			Type: entities.BlockTypeText,
		},
	},
}

// DashboardContainerBlockSchema is only used to validate the blocks of a DashboardContainer
var DashboardContainerBlockSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       entities.BlockTypeInput,
			LabelNames: []string{"name"},
		},
		{
			Type:       entities.BlockTypeParam,
			LabelNames: []string{"name"},
		},
		{
			Type: entities.BlockTypeContainer,
		},
		{
			Type: entities.BlockTypeCard,
		},
		{
			Type: entities.BlockTypeChart,
		},
		{
			Type: entities.BlockTypeBenchmark,
		},
		{
			Type: entities.BlockTypeControl,
		},
		{
			Type: entities.BlockTypeFlow,
		},
		{
			Type: entities.BlockTypeGraph,
		},
		{
			Type: entities.BlockTypeHierarchy,
		},
		{
			Type: entities.BlockTypeImage,
		},
		{
			Type: entities.BlockTypeTable,
		},
		{
			Type: entities.BlockTypeText,
		},
	},
}

var BenchmarkBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "children"},
		{Name: "description"},
		{Name: "documentation"},
		{Name: "tags"},
		{Name: "title"},
		// for report benchmark blocks
		{Name: "width"},
		{Name: "base"},
		{Name: "type"},
		{Name: "display"},
	},
}

// QueryProviderBlockSchema schema for all blocks satisfying QueryProvider interface
// NOTE: these are just the blocks/attributes that are explicitly decoded
// other query provider properties are implicitly decoded using tags
var QueryProviderBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "args"},
	},
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "param",
			LabelNames: []string{"name"},
		},
		{
			Type:       "with",
			LabelNames: []string{"name"},
		},
	},
}

// NodeAndEdgeProviderSchema is used to decode graph/hierarchy/flow
// (EXCEPT categories)
var NodeAndEdgeProviderSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "args"},
	},
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "param",
			LabelNames: []string{"name"},
		},
		{
			Type:       "category",
			LabelNames: []string{"name"},
		},
		{
			Type:       "with",
			LabelNames: []string{"name"},
		},
		{
			Type: entities.BlockTypeNode,
		},
		{
			Type: entities.BlockTypeEdge,
		},
	},
}

var ParamDefBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "description"},
		{Name: "default"},
	},
}

var VariableBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{
			Name: "description",
		},
		{
			Name: "default",
		},
		{
			Name: "type",
		},
		{
			Name: "sensitive",
		},
	},
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type: "validation",
		},
	},
}
