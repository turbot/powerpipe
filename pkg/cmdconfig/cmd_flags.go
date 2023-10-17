package cmdconfig

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/turbot/powerpipe/pkg/error_helpers"
)

var requiredColor = color.New(color.Bold).SprintfFunc()

type FlagOpt func(c *cobra.Command, name string, key string)

// FlagOptions :: shortcut for common flag options
var FlagOptions = struct {
	Required      func() FlagOpt
	Hidden        func() FlagOpt
	Deprecated    func(string) FlagOpt
	NoOptDefVal   func(string) FlagOpt
	WithShortHand func(string) FlagOpt
}{
	Required:      requiredOpt,
	Hidden:        hiddenOpt,
	Deprecated:    deprecatedOpt,
	NoOptDefVal:   noOptDefValOpt,
	WithShortHand: withShortHand,
}

// Helper function to mark a flag as required
func requiredOpt() FlagOpt {
	return func(c *cobra.Command, name, key string) {
		err := c.MarkFlagRequired(key)
		error_helpers.FailOnErrorWithMessage(err, "could not mark flag as required")
		key = fmt.Sprintf("required.%s", key)
		viper.GetViper().Set(key, true)
		u := c.Flag(name).Usage
		c.Flag(name).Usage = fmt.Sprintf("%s %s", u, requiredColor("(required)"))
	}
}

func hiddenOpt() FlagOpt {
	return func(c *cobra.Command, name, _ string) {
		c.Flag(name).Hidden = true
	}
}

func deprecatedOpt(replacement string) FlagOpt {
	return func(c *cobra.Command, name, _ string) {
		c.Flag(name).Deprecated = fmt.Sprintf("please use %s", replacement)
	}
}

func noOptDefValOpt(noOptDefVal string) FlagOpt {
	return func(c *cobra.Command, name, _ string) {
		c.Flag(name).NoOptDefVal = noOptDefVal
	}
}

func withShortHand(shorthand string) FlagOpt {
	return func(c *cobra.Command, name, _ string) {
		c.Flag(name).Shorthand = shorthand
	}
}
