package cmd

import (
	"log/slog"

	"github.com/turbot/powerpipe/internal/constants"
)

// IsLazyLoadEnabled checks if lazy loading is enabled.
// Lazy loading is enabled by default. Set POWERPIPE_WORKSPACE_PRELOAD=true to disable.
func IsLazyLoadEnabled() bool {
	// Check if workspace preload is enabled - this forces eager loading
	if constants.WorkspacePreloadEnabled() {
		slog.Info("Workspace preload enabled via POWERPIPE_WORKSPACE_PRELOAD - using eager loading")
		return false
	}

	// Default: lazy loading enabled
	return true
}
