package main

import (
	"context"
	"github.com/turbot/powerpipe/internal/constants"
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
	constants.SetAppSpecificConstants()

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
