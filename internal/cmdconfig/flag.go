package cmdconfig

import (
	"github.com/spf13/cobra"
	"github.com/turbot/pipe-fittings/v2/cmdconfig"
)

// TODO temp until new pipe-fittings release
func Deprecated(msg string) cmdconfig.FlagOption {
	return func(c *cobra.Command, name, _ string) {
		c.Flag(name).Deprecated = msg
	}
}
