package constants

import (
	"fmt"

	"github.com/turbot/pipe-fittings/v2/constants"
)

type ErrorNoModDefinition struct{}

func (e ErrorNoModDefinition) Error() string {
	return fmt.Sprintf("This command requires a mod definition file (mod.pp) - could not find in the current directory tree.\n\nYou can either clone a mod repository or install a mod using %s and run this command from the cloned/installed mod directory.\nPlease refer to: https://powerpipe.io/docs/build#powerpipe-mods", constants.Bold("powerpipe mod install"))
}
