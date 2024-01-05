package main

import (
	"context"
	"github.com/turbot/powerpipe/internal/cmdconfig"
	"os"

	"github.com/spf13/viper"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/powerpipe/internal/cmd"
)

var exitCode int

var (
	// These variables will be set by GoReleaser.
	version = "0.0.0"
	commit  = "none"
	date    = "unknown"
	builtBy = "local"
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
	viper.SetDefault("main.version", version)
	viper.SetDefault("main.commit", commit)
	viper.SetDefault("main.date", date)
	viper.SetDefault("main.builtBy", builtBy)
}
