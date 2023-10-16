package main

import (
	"context"
	"os"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/powerpipe/pkg/cmd"
	"github.com/turbot/powerpipe/pkg/utils"
	"github.com/turbot/steampipe/pkg/error_helpers"
)

var exitCode int

func main() {
	ctx := context.Background()
	utils.LogTime("main start")
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
