package main

import (
	"context"
	"github.com/turbot/go-kit/files"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/powerpipe/internal/version"
	"io"
	"log"
	"os"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/powerpipe/internal/cmd"
)

var exitCode int

func main() {
	ctx := context.Background()
	utils.LogTime("main start")

	// TODO add logger - discard logs for now
	log.SetOutput(io.Discard)

	// set app specific constants defined in pipe-fittings
	setAppConstants()

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
	cmd.InitCmd()

	// execute the command
	exitCode = cmd.Execute()
}

// set app specific constants defined in pipe-fittings
func setAppConstants() {
	// set the default install dir
	installDir, err := files.Tildefy("~/.powerpipe")
	if err != nil {
		panic(err)
	}
	constants.DefaultInstallDir = installDir

	constants.AppName = "powerpipe"
	constants.ClientConnectionAppNamePrefix = "powerpipe_client"
	constants.ServiceConnectionAppNamePrefix = "powerpipe_service"
	constants.ClientSystemConnectionAppNamePrefix = "powerpipe_client_system"
	constants.AppVersion = version.PowerpipeVersion
	// default to local steampipe service
	constants.DefaultWorkspaceDatabase = "postgres://steampipe@127.0.0.1:9193/steampipe"

}
