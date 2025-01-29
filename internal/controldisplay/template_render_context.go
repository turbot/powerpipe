package controldisplay

type TemplateRenderConfig struct {
	RenderHeader bool
	Separator    string
}

type TemplateRenderConstants struct {
	PowerpipeVersion string
	WorkingDir       string
}

type TemplateRenderContext struct {
	Constants TemplateRenderConstants
	Config    TemplateRenderConfig
	Data      any
}
