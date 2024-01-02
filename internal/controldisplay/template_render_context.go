package controldisplay

import "github.com/turbot/powerpipe/internal/controlexecute"

type TemplateRenderConfig struct {
	RenderHeader bool
	Separator    string
}

type TemplateRenderConstants struct {
	SteampipeVersion string
	WorkingDir       string
}

type TemplateRenderContext struct {
	Constants TemplateRenderConstants
	Config    TemplateRenderConfig
	Data      *controlexecute.ExecutionTree
}
