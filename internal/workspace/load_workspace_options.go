package workspace

import (
	"github.com/turbot/pipe-fittings/v2/connection"
)

type LoadPowerpipeWorkspaceOption func(*LoadPowerpipeWorkspaceConfig)

type LoadPowerpipeWorkspaceConfig struct {
	skipResourceLoadIfNoModfile bool
	pipelingConnections         map[string]connection.PipelingConnection
	blockTypeInclusions         []string
	validateVariables           bool
	supportLateBinding          bool

	// Lazy loading options
	lazyLoad       bool
	lazyLoadConfig LazyLoadConfig
}

func newLoadPowerpipeWorkspaceConfig() *LoadPowerpipeWorkspaceConfig {
	return &LoadPowerpipeWorkspaceConfig{
		pipelingConnections: make(map[string]connection.PipelingConnection),
		validateVariables:   true,
		supportLateBinding:  true,
		lazyLoadConfig:      DefaultLazyLoadConfig(),
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

// WithLazyLoading enables lazy loading mode where resources are loaded on-demand
// instead of all at startup. This provides faster startup and lower memory usage.
func WithLazyLoading(enabled bool) LoadPowerpipeWorkspaceOption {
	return func(m *LoadPowerpipeWorkspaceConfig) {
		m.lazyLoad = enabled
	}
}

// WithLazyLoadConfig sets the lazy loading configuration.
func WithLazyLoadConfig(config LazyLoadConfig) LoadPowerpipeWorkspaceOption {
	return func(m *LoadPowerpipeWorkspaceConfig) {
		m.lazyLoadConfig = config
		m.lazyLoad = true
	}
}
