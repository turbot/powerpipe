package cmd

import (
	"context"
	"fmt"
	"github.com/turbot/powerpipe/internal/service/api"
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
	"gopkg.in/olahol/melody.v1"
)

func dashboardCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "dashboard",
		Args:  cobra.NoArgs,
		Run:   runDashboardCmd,
		Short: "Start Powerpipe dashboard server",
		Long:  "Start Powerpipe dashboard server.",
	}

	cmdconfig.
		OnCmd(cmd).
		AddModLocationFlag().
		AddBoolFlag(constants.ArgHelp, false, "Help for service start", cmdconfig.FlagOptions.WithShortHand("h")).
		AddBoolFlag(constants.ArgBrowser, true, "Specify whether to launch the browser after starting the dashboard server")

	return cmd
}

func runDashboardCmd(cmd *cobra.Command, _ []string) {
	ctx := context.Background()
	ctx, stopFn := signal.NotifyContext(ctx, os.Interrupt)
	defer stopFn()

	// initialise the workspace
	modInitData := dashboard.InitDashboard(ctx)
	error_helpers.FailOnError(modInitData.Result.Error)

	// ensure dashboard assets
	err := dashboardassets.Ensure(ctx)
	error_helpers.FailOnError(err)

	// setup a new webSocket service
	webSocket := melody.New()
	// create the dashboardServer
	dashboardServer, err := dashboardserver.NewServer(ctx, modInitData.Client, modInitData.WorkspaceEvents, webSocket)
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
	// start browser if required
	if viper.GetBool(constants.ArgBrowser) {
		url := buildDashboardURL(9194, modInitData.Workspace)
		if err := utils.OpenBrowser(url); err != nil {
			dashboardserver.OutputWarning(ctx, "Could not start web browser.")
		}
	}
	dashboardserver.OutputMessage(ctx, "server started")
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

// slog.Debug("runDashboardCmd exiting")
// }
