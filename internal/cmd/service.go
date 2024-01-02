package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/cmdconfig"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/pipe-fittings/workspace"
	"github.com/turbot/powerpipe/internal/dashboard"
	"github.com/turbot/powerpipe/internal/dashboardassets"
	"github.com/turbot/powerpipe/internal/dashboardserver"
	"github.com/turbot/powerpipe/internal/service/api"
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
		AddBoolFlag(constants.ArgHelp, false, "Help for service start", cmdconfig.FlagOptions.WithShortHand("h")).
		AddBoolFlag(constants.ArgBrowser, true, "Specify whether to launch the browser after starting the powerpipe server")

	return cmd
}

func runServiceStartCmd(cmd *cobra.Command, _ []string) {
	ctx := context.Background()
	ctx, stopFn := signal.NotifyContext(ctx, os.Interrupt)
	defer stopFn()

	// initialise the workspace
	modInitData := dashboard.InitDashboard(ctx)
	error_helpers.FailOnError(modInitData.Result.Error)

	// ensure dashboard assets
	err := dashboardassets.Ensure(ctx)
	if err != nil {
		panic(err)
	}

	// setup a new webSocket service
	webSocket := melody.New()
	// create the dashboardServer
	dashboardServer, err := dashboardserver.NewServer(ctx, modInitData.Client, modInitData.Workspace, webSocket)
	error_helpers.FailOnError(err)

	// send it over to the powerpipe API Server
	powerpipeService, err := api.NewAPIService(ctx, api.WithWebSocket(webSocket), api.WithWorkspace(modInitData.Workspace))
	if err != nil {
		error_helpers.FailOnError(err)
	}
	dashboardServer.InitAsync(ctx)

	// start the API server
	err = powerpipeService.Start()
	if err != nil {
		error_helpers.FailOnError(err)
	}
	// start browser if required
	if viper.GetBool(constants.ArgBrowser) {
		url := buildDashboardURL(9194, modInitData.Workspace)
		if err := utils.OpenBrowser(url); err != nil {
			dashboardserver.OutputWarning(ctx, "Could not start web browser.")
		}
	}
	fmt.Println("server started")
	<-ctx.Done()
}

func buildDashboardURL(serverPort dashboardserver.ListenPort, w *workspace.Workspace) string {
	url := fmt.Sprintf("http://localhost:%d", serverPort)
	if len(w.SourceSnapshots) == 1 {
		for snapshotName := range w.GetResourceMaps().Snapshots {
			url += fmt.Sprintf("/%s", snapshotName)
			break
		}
	}
	return url
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
