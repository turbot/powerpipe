package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/turbot/pipe-fittings/cmdconfig"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/powerpipe/internal/dashboardassets"
	"github.com/turbot/powerpipe/internal/dashboardserver"
	"github.com/turbot/powerpipe/internal/initialisation"
	"github.com/turbot/powerpipe/internal/service/api"
	"gopkg.in/olahol/melody.v1"
	"os"
	"os/signal"
)

func serverCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "server",
		Args:  cobra.NoArgs,
		Run:   runServerCmd,
		Short: "Start Powerpipe dashboard server",
		Long:  "Start Powerpipe dashboard server.",
	}

	// TODO KAI CHECK ARGS
	cmdconfig.
		OnCmd(cmd).
		AddModLocationFlag().
		AddBoolFlag(constants.ArgHelp, false, "Help for service start", cmdconfig.FlagOptions.WithShortHand("h")).
		AddIntFlag(constants.ArgDashboardPort, constants.DashboardServerDefaultPort, "Dashboard server port")

	return cmd
}

func runServerCmd(cmd *cobra.Command, _ []string) {
	ctx := context.Background()
	ctx, stopFn := signal.NotifyContext(ctx, os.Interrupt)
	defer stopFn()

	// TODO KAI do we need a client?? I don't think so
	// add option
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
