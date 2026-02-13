package constants

import "os"

const (
	EnvListen           = "POWERPIPE_LISTEN"
	EnvPort             = "POWERPIPE_PORT"
	EnvBenchmarkTimeout = "POWERPIPE_BENCHMARK_TIMEOUT"
	EnvDashboardTimeout = "POWERPIPE_DASHBOARD_TIMEOUT"
	EnvDisplayWidth     = "POWERPIPE_DISPLAY_WIDTH"
	// EnvConfigDump is an undocumented variable is subject to change in the future
	EnvConfigDump = "POWERPIPE_CONFIG_DUMP"
	// EnvWorkspacePreload forces full workspace loading at startup instead of lazy loading
	// Set to "true" to enable preloading (useful as a fallback if lazy loading causes issues)
	EnvWorkspacePreload = "POWERPIPE_WORKSPACE_PRELOAD"
)

// WorkspacePreloadEnabled checks if POWERPIPE_WORKSPACE_PRELOAD is set to "true".
// When enabled, the workspace will be fully loaded at startup instead of using lazy loading.
// This provides a safety valve for any edge cases discovered with lazy loading.
func WorkspacePreloadEnabled() bool {
	val := os.Getenv(EnvWorkspacePreload)
	return val == "true" || val == "1" || val == "yes"
}
