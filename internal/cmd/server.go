package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/turbot/pipe-fittings/cmdconfig"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	localcmdconfig "github.com/turbot/powerpipe/internal/cmdconfig"
	localconstants "github.com/turbot/powerpipe/internal/constants"
	"github.com/turbot/powerpipe/internal/dashboardassets"
	"github.com/turbot/powerpipe/internal/dashboardserver"
	"github.com/turbot/powerpipe/internal/initialisation"
	"github.com/turbot/powerpipe/internal/service/api"
	"gopkg.in/olahol/melody.v1"
)

func serverCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "server",
		Args:  cobra.NoArgs,
		Run:   runServerCmd,
		Short: "Start Powerpipe dashboard server",
		Long:  "Start Powerpipe dashboard server.",
	}

	// TODO KAI CHECK ARGS (https://github.com/turbot/powerpipe/issues/106)
	cmdconfig.
		OnCmd(cmd).
		AddModLocationFlag().
		AddBoolFlag(constants.ArgHelp, false, "Help for service start", cmdconfig.FlagOptions.WithShortHand("h")).
		AddIntFlag(constants.ArgPort, constants.DashboardServerDefaultPort, "Web server port").
		AddBoolFlag(constants.ArgWatch, true, "Watch mod files for changes when running powerpipe server").
		AddStringFlag(constants.ArgListen, "", "Accept connections from local (localhost only) or network (all interfaces / IP addresses)").
		AddStringSliceFlag(constants.ArgVariable, []string{}, "Specify the value of a variable. Multiple --var arguments may be passed.").
		AddStringFlag(constants.ArgVarFile, "", "Specify a .ppvar file containing variable values.")
	return cmd
}

func runServerCmd(cmd *cobra.Command, _ []string) {
	ctx := context.Background()
	ctx, stopFn := signal.NotifyContext(ctx, os.Interrupt)
	defer stopFn()

	// if diagnostic mode is set, print out config and return
	if _, ok := os.LookupEnv(localconstants.EnvConfigDump); ok {
		localcmdconfig.DisplayConfig()
		return
	}

	// initialise the workspace
	modInitData := initialisation.NewInitData(ctx, "dashboard")
	error_helpers.FailOnError(modInitData.Result.Error)

	// ensure dashboard assets
	err := dashboardassets.Ensure(ctx)
	error_helpers.FailOnError(err)

	// setup a new webSocket service
	webSocket := melody.New()
	// create the dashboardServer
	dashboardServer, err := dashboardserver.NewServer(ctx, modInitData.WorkspaceEvents, webSocket)
	error_helpers.FailOnError(err)

	// send it over to the powerpipe API Server
	powerpipeService, err := api.NewAPIService(ctx, api.WithWebSocket(webSocket), api.WithWorkspace(modInitData.Workspace))
	if err != nil {
		error_helpers.FailOnError(err)
	}
	dashboardServer.InitAsync(ctx)

	//start the API server
	err = powerpipeService.Start()
	if err != nil {
		error_helpers.FailOnError(err)
	}

	dashboardserver.OutputMessage(ctx, "server started")
	<-ctx.Done()
}
