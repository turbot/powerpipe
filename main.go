package main

import (
	"context"
	"os"

	"github.com/spf13/viper"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/powerpipe/internal/cmd"
	"github.com/turbot/powerpipe/internal/cmdconfig"
	"github.com/turbot/powerpipe/internal/constants"
)

var exitCode int

var (
	// These variables will be set by GoReleaser
	version = constants.DefaultVersion
	commit  = constants.DefaultCommit
	date    = constants.DefaultDate
	builtBy = constants.DefaultBuiltBy
)

func main() {
	ctx := context.Background()
	utils.LogTime("main start")

	// add the auto-populated version properties into viper
	setVersionProperties()
	// set app specific constants defined in pipe-fittings
	cmdconfig.SetAppSpecificConstants()

	defer func() {
		if r := recover(); r != nil {
			error_helpers.ShowError(ctx, helpers.ToError(r))
			if exitCode == 0 {
				exitCode = 255
			}
		}
		utils.LogTime("main end")
		utils.DisplayProfileData()
		os.Exit(exitCode)
	}()

	// execute the root command
	exitCode = cmd.Execute()
}

func setVersionProperties() {
	viper.SetDefault(constants.ConfigKeyVersion, version)
	viper.SetDefault(constants.ConfigKeyCommit, commit)
	viper.SetDefault(constants.ConfigKeyDate, date)
	viper.SetDefault(constants.ConfigKeyBuiltBy, builtBy)
}
