package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/dashboardserver"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/powerpipe/internal/cmdconfig"
	"github.com/turbot/powerpipe/internal/dashboard"
	"github.com/turbot/powerpipe/internal/service/api"
	exported_commandconfig "github.com/turbot/powerpipe/pkg/cmdconfig"
	"gopkg.in/olahol/melody.v1"
)

func serviceCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "service [command]",
		Args:  cobra.NoArgs,
		Short: "Powerpipe service management",
		Long: `Powerpipe service management.

Run Powerpipe as a local service, exposing it as a database endpoint for
connection from any compatible database client.`,
	}

	cmd.AddCommand(serviceStartCmd())
	cmd.Flags().BoolP(constants.ArgHelp, "h", false, "Help for service")
	return cmd
}

func serviceStartCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "start",
		Args:  cobra.NoArgs,
		Run:   runServiceStartCmd,
		Short: "Start Powerpipe in service mode",
		Long: `Start the Powerpipe service.

Run Powerpipe as a local service, exposing it as a database endpoint for
connection from any compatible database client.`,
	}

	cmdconfig.
		OnCmd(cmd).
		AddModLocationFlag().
		AddBoolFlag(constants.ArgHelp, false, "Help for service start", exported_commandconfig.FlagOptions.WithShortHand("h")).
		AddStringFlag(constants.ArgInstallDir, dashboard.DefaultInstallDir, "The default install directory")

	return cmd
}

func runServiceStartCmd(cmd *cobra.Command, _ []string) {
	dashboard.PowerpipeDir = "~/.powerpipe"

	ctx := context.Background()
	ctx, stopFn := signal.NotifyContext(ctx, os.Interrupt)
	defer stopFn()

	err := dashboard.Ensure(ctx)
	if err != nil {
		panic(err)
	}

	// setup a new webSocket service
	webSocket := melody.New()
	modInitData := dashboard.InitDashboard(ctx)
	error_helpers.FailOnError(modInitData.Result.Error)
	// create the dashboardServer
	dashboardServer, err := dashboardserver.NewServer(ctx, modInitData.Client, modInitData.Workspace, webSocket)
	error_helpers.FailOnError(err)

	// send it over to the API Server
	powerpipeService, err := api.NewAPIService(ctx, api.WithWebSocket(webSocket))
	if err != nil {
		panic(err)
	}
	dashboardServer.InitAsync(ctx)

	// start the API server
	err = powerpipeService.Start()
	if err != nil {
		panic(err)
	}
	fmt.Println("server started")
	<-ctx.Done()
}

// func StartDashboardServer(ctx context.Context, serverPort dashboardserver.ListenPort, serverListen dashboardserver.ListenType) {
// 	// create context for the dashboard execution
// 	dashboardCtx, cancel := context.WithCancel(ctx)
// 	contexthelpers.StartCancelHandler(cancel)

// 	// ensure dashboard assets are present and extract if not
// 	err := dashboard.Ensure(dashboardCtx)
// 	error_helpers.FailOnError(err)

// 	// load the workspace
// 	initData := initDashboard(dashboardCtx)
// 	defer initData.Cleanup(dashboardCtx)
// 	if initData.Result.Error != nil {
// 		exitCode = constants.ExitCodeInitializationFailed
// 		error_helpers.FailOnError(initData.Result.Error)
// 	}

// 	// if there is a usage warning we display it
// 	initData.Result.DisplayMessage = dashboardserver.OutputMessage
// 	initData.Result.DisplayWarning = dashboardserver.OutputWarning
// 	initData.Result.DisplayMessages()

// 	// create the server
// 	server, err := dashboardserver.NewServer(dashboardCtx, initData.Client, initData.Workspace)
// 	error_helpers.FailOnError(err)

// 	// start the server asynchronously - this returns a chan which is signalled when the internal API server terminates
// 	doneChan := server.Start(dashboardCtx)

// 	// cleanup
// 	defer server.Shutdown(dashboardCtx)

// 	// server has started - update state file/start browser, as required
// 	onServerStarted(dashboardCtx, serverPort, serverListen, initData.Workspace)

// 	// wait for API server to terminate
// 	<-doneChan

// log.Println("[TRACE] runDashboardCmd exiting")
// }
