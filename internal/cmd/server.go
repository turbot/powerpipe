package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/app_specific"
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
		Long: `Run the Powerpipe server, including the dashbaord server and the API. 
		
Powerpipe server runs in the foreground; Press Ctrl-C to exit.`,
	}

	cmdconfig.
		OnCmd(cmd).
		AddModLocationFlag().
		AddBoolFlag(constants.ArgHelp, false, "Help for service start", cmdconfig.FlagOptions.WithShortHand("h")).
		AddIntFlag(constants.ArgPort, dashboardserver.DashboardServerDefaultPort, "Web server port").
		AddBoolFlag(constants.ArgWatch, true, "Watch mod files for changes when running powerpipe server").
		AddStringFlag(constants.ArgListen, string(dashboardserver.ListenTypeLocal), "Accept connections from local (localhost only) or network (all interfaces / IP addresses)").
		AddStringSliceFlag(constants.ArgVariable, []string{}, "Specify the value of a variable. Multiple --var arguments may be passed.").
		AddStringFlag(constants.ArgVarFile, "", "Specify a .ppvar file containing variable values.").
		AddStringFlag(constants.ArgDatabase, app_specific.DefaultDatabase, "Turbot Pipes workspace database")
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

	// retrieve server params
	serverPort := dashboardserver.ListenPort(viper.GetInt(constants.ArgPort))
	error_helpers.FailOnError(serverPort.IsValid())

	serverListen := dashboardserver.ListenType(viper.GetString(constants.ArgListen))
	error_helpers.FailOnError(serverListen.IsValid())

	serverHost := ""
	if serverListen == dashboardserver.ListenTypeLocal {
		serverHost = "127.0.0.1"
	}
	if err := utils.IsPortBindable(serverHost, int(serverPort)); err != nil {
		exitCode = constants.ExitCodeBindPortUnavailable
		error_helpers.FailOnError(err)
	}

	// initialise the workspace
	modInitData := initialisation.NewInitData[*modconfig.Dashboard](ctx, cmd)
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
	powerpipeService, err := api.NewAPIService(ctx, api.WithWebSocket(webSocket), api.WithWorkspace(modInitData.Workspace), api.WithHttpPort(serverPort))
	if err != nil {
		error_helpers.FailOnError(err)
	}
	dashboardServer.InitAsync(ctx)

	//start the API server
	err = powerpipeService.Start()
	if err != nil {
		error_helpers.FailOnError(err)
	}

	dashboardserver.OutputReady(ctx, fmt.Sprintf("Dashboard server started on %d and listening on %s", serverPort, viper.GetString(constants.ArgListen)))
	dashboardserver.OutputMessage(ctx, fmt.Sprintf("Visit http://localhost:%d", serverPort))
	dashboardserver.OutputMessage(ctx, "Press Ctrl+C to exit")

	<-ctx.Done()
}
