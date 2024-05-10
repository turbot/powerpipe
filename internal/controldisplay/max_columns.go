package controldisplay

import (
	"github.com/karrick/gows"
	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/constants"
)

func GetMaxCols() int {
	colsAvailable, _, _ := gows.GetWinSize()
	// check if STEAMPIPE_DISPLAY_WIDTH env variable is set
	if viper.IsSet(constants.ArgDisplayWidth) {
		colsAvailable = viper.GetInt(constants.ArgDisplayWidth)
	}
	return colsAvailable
}
