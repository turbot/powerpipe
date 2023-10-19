package main

import (
	"context"
	"github.com/turbot/pipe-fittings/filepaths"
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

	// set default install dir
	filepaths.DefaultInstallDir = ".powerpipe"

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

// func main() {
// 	dashboard.PowerpipeDir = "~/.Powerpipe"

// 	ctx := context.Background()
// 	ctx, stopFn := signal.NotifyContext(ctx, os.Interrupt)
// 	defer stopFn()

// 	err := dashboard.Ensure(ctx)
// 	if err != nil {
// 		panic(err)
// 	}

// 	server, err := api.NewAPIService(ctx)
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = server.Start()
// 	if err != nil {
// 		panic(err)
// 	}
// 	println("server started")
// 	<-ctx.Done()
// }
