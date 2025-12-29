package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/v2/cmdconfig"
)

// LazyLoadFlag is the CLI flag name for enabling lazy loading
const LazyLoadFlag = "lazy-load"

// EnvLazyLoad is the environment variable to control lazy loading globally
const EnvLazyLoad = "POWERPIPE_LAZY_LOAD"

// AddLazyLoadFlag adds the lazy loading flag to a command
func AddLazyLoadFlag(cmd *cobra.Command) {
	cmdconfig.OnCmd(cmd).
		AddBoolFlag(LazyLoadFlag, false, "Enable lazy loading of resources (reduces memory usage, faster startup)")
}

// IsLazyLoadEnabled checks if lazy loading is enabled.
// Priority: 1. CLI flag, 2. Environment variable, 3. Default (false)
func IsLazyLoadEnabled(cmd *cobra.Command) bool {
	// Check if explicitly set via CLI flag
	if cmd.Flags().Changed(LazyLoadFlag) {
		return viper.GetBool(LazyLoadFlag)
	}

	// Check environment variable
	if envVal := os.Getenv(EnvLazyLoad); envVal != "" {
		return envVal == "true" || envVal == "1" || envVal == "yes"
	}

	// Default: lazy loading disabled for backward compatibility
	return false
}
