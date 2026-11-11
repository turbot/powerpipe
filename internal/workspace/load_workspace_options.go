package workspace

import (
	"github.com/turbot/pipe-fittings/v2/connection"
	"github.com/turbot/pipe-fittings/v2/parse"
)

type LoadPowerpipeWorkspaceOption func(*LoadPowerpipeWorkspaceConfig)

type LoadPowerpipeWorkspaceConfig struct {
	skipResourceLoadIfNoModfile bool
	pipelingConnections         map[string]connection.PipelingConnection
	blockTypeInclusions         []string
	validateVariables           bool
	supportLateBinding          bool
	resourceTargets             []parse.ResourceTarget
}

func newLoadPowerpipeWorkspaceConfig() *LoadPowerpipeWorkspaceConfig {
	return &LoadPowerpipeWorkspaceConfig{
		pipelingConnections: make(map[string]connection.PipelingConnection),
		validateVariables:   true,
		supportLateBinding:  true,
	}
}

func WithPipelingConnections(pipelingConnections map[string]connection.PipelingConnection) LoadPowerpipeWorkspaceOption {
	return func(m *LoadPowerpipeWorkspaceConfig) {
		m.pipelingConnections = pipelingConnections
	}
}

func WithLateBinding(enabled bool) LoadPowerpipeWorkspaceOption {
	return func(m *LoadPowerpipeWorkspaceConfig) {
		m.supportLateBinding = enabled
	}
}

func WithBlockType(blockTypeInclusions []string) LoadPowerpipeWorkspaceOption {
	return func(m *LoadPowerpipeWorkspaceConfig) {
		m.blockTypeInclusions = blockTypeInclusions
	}
}

// WithResourceTargets loads only the resources that the targets need (see parse.ModParseContext.SetResourceTargets)
func WithResourceTargets(targets []parse.ResourceTarget) LoadPowerpipeWorkspaceOption {
	return func(m *LoadPowerpipeWorkspaceConfig) {
		m.resourceTargets = targets
	}
}

func WithVariableValidation(enabled bool) LoadPowerpipeWorkspaceOption {
	return func(m *LoadPowerpipeWorkspaceConfig) {
		m.validateVariables = enabled
	}
}

// TODO this is only needed as Pipe fittings tests rely on loading workspaces without modfiles
func WithSkipResourceLoadIfNoModfile(enabled bool) LoadPowerpipeWorkspaceOption {
	return func(m *LoadPowerpipeWorkspaceConfig) {
		m.skipResourceLoadIfNoModfile = enabled
	}
}
