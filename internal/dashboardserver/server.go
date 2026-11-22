package dashboardserver

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/spf13/viper"
	typeHelpers "github.com/turbot/go-kit/types"
	"github.com/turbot/pipe-fittings/v2/backend"
	"github.com/turbot/pipe-fittings/v2/connection"
	"github.com/turbot/pipe-fittings/v2/constants"
	"github.com/turbot/pipe-fittings/v2/error_helpers"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/pipe-fittings/v2/schema"
	"github.com/turbot/pipe-fittings/v2/steampipeconfig"
	"github.com/turbot/powerpipe/internal/dashboardevents"
	"github.com/turbot/powerpipe/internal/dashboardexecute"
	"github.com/turbot/powerpipe/internal/initialisation"
	"github.com/turbot/powerpipe/internal/workspace"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
	"gopkg.in/olahol/melody.v1"
)

type Server struct {
	mutex                   *sync.Mutex
	dashboardClients        map[string]*DashboardClientInfo
	webSocket               *melody.Melody
	workspace               *workspace.PowerpipeWorkspace
	defaultDatabase         connection.ConnectionStringProvider
	defaultSearchPathConfig backend.SearchPathConfig
}

func NewServer(ctx context.Context, initData *initialisation.InitData, webSocket *melody.Melody) (*Server, error) {
	OutputWait(ctx, "Starting WorkspaceEvents Server")

	var dashboardClients = make(map[string]*DashboardClientInfo)

	var mutex = &sync.Mutex{}

	w := initData.Workspace

	server := &Server{
		mutex:                   mutex,
		dashboardClients:        dashboardClients,
		webSocket:               webSocket,
		workspace:               w,
		defaultDatabase:         initData.DefaultDatabase,
		defaultSearchPathConfig: initData.DefaultSearchPathConfig,
	}

	w.RegisterDashboardEventHandler(ctx, server.HandleDashboardEvent)
	// tell the clients when a load starts and ends, at once and not through the dashboard events: the handler of an
	// earlier event can wait for the reload to end, and the start would then be told after the end
	w.SetLoadStateListener(server.broadcastLoadState)

	// --watch (or watch in the workspace) turns the watchers off: the server then loads the workspace once, and a
	// restart picks up changes
	err := startWatchers(ctx, w, viper.GetBool(constants.ArgWatch))
	OutputMessage(ctx, "WorkspaceEvents loaded")

	return server, err
}

// startWatchers starts the watchers that reload the workspace when its mod or variable files change, and when the
// installed mods change. With watch false it starts neither, so nothing reloads the workspace.
func startWatchers(ctx context.Context, w *workspace.PowerpipeWorkspace, watch bool) error {
	if !watch {
		slog.Info("not watching the workspace for changes")
		return nil
	}
	err := w.SetupWatcher(ctx, func(c context.Context, e error) {})
	if err == nil {
		// also reload when the installed mods change (e.g. after `powerpipe mod update`)
		if modWatchErr := w.SetupModInstallWatcher(ctx); modWatchErr != nil {
			// not fatal: the server still reloads for changes to mod and variable files
			slog.Warn("failed to watch the installed mods - the server will not reload when they change", "error", modWatchErr)
			OutputMessage(ctx, fmt.Sprintf("Warning: failed to watch the installed mods for changes: %s", modWatchErr))
		}
	}
	return err
}

// Start starts the API server
// it returns a channel which is signalled when the API server terminates
func (s *Server) Start(ctx context.Context) chan struct{} {
	s.InitAsync(ctx)
	return startAPIAsync(ctx, s.webSocket)
}

// Shutdown stops the API server
func (s *Server) Shutdown(ctx context.Context) {
	slog.Debug("Server shutdown")

	if s.webSocket != nil {
		slog.Debug("closing websocket")
		if err := s.webSocket.Close(); err != nil {
			error_helpers.ShowErrorWithMessage(ctx, err, "Websocket shutdown failed")
		}
		slog.Debug("closed websocket")
	}

	slog.Debug("Server shutdown complete")
}

func (s *Server) HandleDashboardEvent(ctx context.Context, event dashboardevents.DashboardEvent) {
	var payloadError error
	var payload []byte
	defer func() {
		if payloadError != nil {
			// we don't expect the build functions to ever error during marshalling
			// this is because the data getting marshalled are not expected to have go specific
			// properties/data in them
			OutputError(ctx, sperr.WrapWithMessage(payloadError, "error building payload for '%s'", reflect.TypeOf(event).String()))
		}
	}()

	switch e := event.(type) {

	case *dashboardevents.WorkspaceError:
		slog.Debug("WorkspaceError event", "error", e.Error)
		payload, payloadError = buildWorkspaceErrorPayload(e)
		if payloadError != nil {
			return
		}
		_ = s.webSocket.Broadcast(payload)
		OutputError(ctx, e.Error)

	case *dashboardevents.ExecutionStarted:
		slog.Debug("ExecutionStarted event", "session ", e.Session, "dashboard", e.Root.GetName())
		payload, payloadError = buildExecutionStartedPayload(e)
		if payloadError != nil {
			return
		}
		s.writePayloadToSession(e.Session, payload)
		OutputWait(ctx, fmt.Sprintf("WorkspaceEvents execution started: %s", e.Root.GetName()))

	case *dashboardevents.ExecutionError:
		slog.Debug("execution error event")
		payload, payloadError = buildExecutionErrorPayload(e)
		if payloadError != nil {
			return
		}

		s.writePayloadToSession(e.Session, payload)
		OutputError(ctx, e.Error)

	case *dashboardevents.ExecutionComplete:
		slog.Debug("execution complete event")
		payload, payloadError = buildExecutionCompletePayload(e)
		if payloadError != nil {
			return
		}
		dashboardName := e.Root.GetName()
		s.writePayloadToSession(e.Session, payload)
		OutputReady(ctx, fmt.Sprintf("Execution complete: %s", dashboardName))

	case *dashboardevents.ControlComplete:
		slog.Debug("ControlComplete event", "session", e.Session, "control", e.Control.GetControlId())
		payload, payloadError = buildControlCompletePayload(e)
		if payloadError != nil {
			return
		}
		s.writePayloadToSession(e.Session, payload)

	case *dashboardevents.ControlError:
		slog.Debug("ControlError event", "session", e.Session, "control", e.Control.GetControlId())
		payload, payloadError = buildControlErrorPayload(e)
		if payloadError != nil {
			return
		}
		s.writePayloadToSession(e.Session, payload)

	case *dashboardevents.LeafNodeUpdated:
		payload, payloadError = buildLeafNodeUpdatedPayload(e)
		if payloadError != nil {
			return
		}
		s.writePayloadToSession(e.Session, payload)

	case *dashboardevents.DashboardChanged:
		slog.Debug("DashboardChanged event")
		deletedDashboards := e.DeletedDashboards
		newDashboards := e.NewDashboards

		changedBenchmarks := e.ChangedBenchmarks
		changedCategories := e.ChangedCategories
		changedContainers := e.ChangedContainers
		changedControls := e.ChangedControls
		changedCards := e.ChangedCards
		changedCharts := e.ChangedCharts
		changedDashboards := e.ChangedDashboards
		changedDetections := e.ChangedDetections
		changedDetectionsBenchmarks := e.ChangedDetectionBenchmarks
		changedEdges := e.ChangedEdges
		changedFlows := e.ChangedFlows
		changedGraphs := e.ChangedGraphs
		changedHierarchies := e.ChangedHierarchies
		changedImages := e.ChangedImages
		changedInputs := e.ChangedInputs
		changedNodes := e.ChangedNodes
		changedTables := e.ChangedTables
		changedTexts := e.ChangedTexts

		// If nothing has changed, ignore
		if len(deletedDashboards) == 0 &&
			len(newDashboards) == 0 &&
			len(changedBenchmarks) == 0 &&
			len(changedCategories) == 0 &&
			len(changedContainers) == 0 &&
			len(changedControls) == 0 &&
			len(changedCards) == 0 &&
			len(changedCharts) == 0 &&
			len(changedDashboards) == 0 &&
			len(changedDetections) == 0 &&
			len(changedDetectionsBenchmarks) == 0 &&
			len(changedEdges) == 0 &&
			len(changedFlows) == 0 &&
			len(changedGraphs) == 0 &&
			len(changedHierarchies) == 0 &&
			len(changedImages) == 0 &&
			len(changedInputs) == 0 &&
			len(changedNodes) == 0 &&
			len(changedTables) == 0 &&
			len(changedTexts) == 0 {
			return
		}

		for k, v := range s.dashboardClients {
			slog.Debug("WorkspaceEvents", "client", k, "event", typeHelpers.SafeString(v.Dashboard))
		}

		// If) any deleted/new/changed dashboards, emit an available dashboards message to clients
		if len(deletedDashboards) != 0 || len(newDashboards) != 0 || len(changedDashboards) != 0 || len(changedBenchmarks) != 0 {
			OutputMessage(ctx, "Available Dashboards updated")

			// Emit dashboard metadata event in case there is a new mod - else the UI won't know about this mod
			payload, payloadError = s.buildServerMetadataPayload(s.workspace.GetModResources(), &steampipeconfig.PipesMetadata{})
			if payloadError != nil {
				return
			}
			_ = s.webSocket.Broadcast(payload)

			// Emit available dashboards event
			workspaceResources := s.workspace.GetPowerpipeModResources()
			payload, payloadError = buildAvailableDashboardsPayload(workspaceResources)
			if payloadError != nil {
				return
			}
			_ = s.webSocket.Broadcast(payload)
		}

		var dashboardsBeingWatched []string

		dashboardClients := s.getDashboardClients()
		for _, dashboardClientInfo := range dashboardClients {
			dashboardName := typeHelpers.SafeString(dashboardClientInfo.Dashboard)
			if dashboardClientInfo.Dashboard != nil {
				if slices.Contains(dashboardsBeingWatched, dashboardName) {
					continue
				}
				dashboardsBeingWatched = append(dashboardsBeingWatched, dashboardName)
			}
		}

		var changedDashboardNames []string
		var newDashboardNames []string

		// Process the changed items and make a note of the dashboard(s) they're in
		changedDashboardNames = append(changedDashboardNames, getDashboardsInterestedInResourceChanges(dashboardsBeingWatched, changedDashboardNames, changedBenchmarks)...)
		changedDashboardNames = append(changedDashboardNames, getDashboardsInterestedInResourceChanges(dashboardsBeingWatched, changedDashboardNames, changedCategories)...)
		changedDashboardNames = append(changedDashboardNames, getDashboardsInterestedInResourceChanges(dashboardsBeingWatched, changedDashboardNames, changedContainers)...)
		changedDashboardNames = append(changedDashboardNames, getDashboardsInterestedInResourceChanges(dashboardsBeingWatched, changedDashboardNames, changedControls)...)
		changedDashboardNames = append(changedDashboardNames, getDashboardsInterestedInResourceChanges(dashboardsBeingWatched, changedDashboardNames, changedCards)...)
		changedDashboardNames = append(changedDashboardNames, getDashboardsInterestedInResourceChanges(dashboardsBeingWatched, changedDashboardNames, changedCharts)...)
		changedDashboardNames = append(changedDashboardNames, getDashboardsInterestedInResourceChanges(dashboardsBeingWatched, changedDashboardNames, changedDetections)...)
		changedDashboardNames = append(changedDashboardNames, getDashboardsInterestedInResourceChanges(dashboardsBeingWatched, changedDashboardNames, changedDetectionsBenchmarks)...)
		changedDashboardNames = append(changedDashboardNames, getDashboardsInterestedInResourceChanges(dashboardsBeingWatched, changedDashboardNames, changedEdges)...)
		changedDashboardNames = append(changedDashboardNames, getDashboardsInterestedInResourceChanges(dashboardsBeingWatched, changedDashboardNames, changedFlows)...)
		changedDashboardNames = append(changedDashboardNames, getDashboardsInterestedInResourceChanges(dashboardsBeingWatched, changedDashboardNames, changedGraphs)...)
		changedDashboardNames = append(changedDashboardNames, getDashboardsInterestedInResourceChanges(dashboardsBeingWatched, changedDashboardNames, changedHierarchies)...)
		changedDashboardNames = append(changedDashboardNames, getDashboardsInterestedInResourceChanges(dashboardsBeingWatched, changedDashboardNames, changedImages)...)
		changedDashboardNames = append(changedDashboardNames, getDashboardsInterestedInResourceChanges(dashboardsBeingWatched, changedDashboardNames, changedInputs)...)
		changedDashboardNames = append(changedDashboardNames, getDashboardsInterestedInResourceChanges(dashboardsBeingWatched, changedDashboardNames, changedNodes)...)
		changedDashboardNames = append(changedDashboardNames, getDashboardsInterestedInResourceChanges(dashboardsBeingWatched, changedDashboardNames, changedTables)...)
		changedDashboardNames = append(changedDashboardNames, getDashboardsInterestedInResourceChanges(dashboardsBeingWatched, changedDashboardNames, changedTexts)...)

		for _, changedDashboard := range changedDashboards {
			if slices.Contains(changedDashboardNames, changedDashboard.Name) {
				continue
			}
			changedDashboardNames = append(changedDashboardNames, changedDashboard.Name)
		}

		for _, changedDashboardName := range changedDashboardNames {
			sessionMap := s.getDashboardClients()
			for sessionId, dashboardClientInfo := range sessionMap {
				if typeHelpers.SafeString(dashboardClientInfo.Dashboard) == changedDashboardName {

					if changedResource := s.getResource(changedDashboardName); changedResource != nil {
						err := dashboardexecute.Executor.ExecuteDashboard(ctx, sessionId, changedResource, dashboardClientInfo.DashboardInputs, s.workspace)
						if err != nil {
							OutputError(ctx, sperr.WrapWithMessage(err, "error executing dashboard"))
						}
					}
				}
			}

		}

		// Special case - if we previously had a workspace error, any previously existing dashboards
		// will come in here as new, so we need to check if any of those new dashboards are being watched.
		// If so, execute them
		for _, newDashboard := range newDashboards {
			if slices.Contains(newDashboardNames, newDashboard.Name()) {
				continue
			}
			newDashboardNames = append(newDashboardNames, newDashboard.Name())
		}

		sessionMap := s.getDashboardClients()
		for _, newDashboardName := range newDashboardNames {
			for sessionId, dashboardClientInfo := range sessionMap {
				if typeHelpers.SafeString(dashboardClientInfo.Dashboard) == newDashboardName {
					if newDashboard := s.getResource(newDashboardName); newDashboard != nil {
						err := dashboardexecute.Executor.ExecuteDashboard(ctx, sessionId, newDashboard, dashboardClientInfo.DashboardInputs, s.workspace)
						if err != nil {
							OutputError(ctx, sperr.WrapWithMessage(err, "error executing dashboard"))
						}
					}
				}
			}
		}

	case *dashboardevents.InputValuesCleared:
		payload, payloadError = buildInputValuesClearedPayload(e)
		if payloadError != nil {
			return
		}

		dashboardClients := s.getDashboardClients()
		if sessionInfo, ok := dashboardClients[e.Session]; ok {
			for _, clearedInput := range e.ClearedInputs {
				delete(sessionInfo.DashboardInputs.Inputs, clearedInput)
			}
		}
		s.writePayloadToSession(e.Session, payload)
	}
}

func (s *Server) InitAsync(ctx context.Context) {
	go func() {
		// Return list of dashboards on connect
		s.webSocket.HandleConnect(func(session *melody.Session) {
			slog.Debug("client connected")
			// the client learns at once whether the workspace is loading, also when it connects during a reload:
			// LoadState does not wait for the reload
			s.addSessionAndSendLoadState(session)
		})

		s.webSocket.HandleDisconnect(func(session *melody.Session) {
			slog.Debug("client disconnected")
			s.clearSession(ctx, session)
		})

		s.webSocket.HandleMessage(s.handleMessageFunc(ctx))
		OutputMessage(ctx, "Initialization complete")
	}()
}

func (s *Server) handleMessageFunc(ctx context.Context) func(session *melody.Session, msg []byte) {
	return func(session *melody.Session, msg []byte) {

		sessionId := s.getSessionId(session)

		var request ClientRequest
		// if we could not decode message - ignore
		err := json.Unmarshal(msg, &request)
		if err != nil {
			slog.Warn("failed to marshal message", "error", err.Error())
			return
		}

		if request.Action != "keep_alive" {
			slog.Debug("handleMessageFunc", "message", string(msg))
		}

		switch request.Action {
		case "get_server_metadata":
			payload, err := s.buildServerMetadataPayload(s.workspace.GetModResources(), &steampipeconfig.PipesMetadata{})
			if err != nil {
				OutputError(ctx, sperr.WrapWithMessage(err, "error building payload for get_metadata"))
			}
			_ = session.Write(payload)
		case "get_available_dashboards":
			payload, err := buildAvailableDashboardsPayload(s.workspace.GetPowerpipeModResources())
			if err != nil {
				OutputError(ctx, sperr.WrapWithMessage(err, "error building payload for get_available_dashboards"))
			}
			_ = session.Write(payload)
		case "select_dashboard":
			dashboard := s.getResource(request.Payload.Dashboard.FullName)
			if dashboard == nil {
				return
			}

			inputValues := request.Payload.InputValues()
			s.setDashboardForSession(sessionId, request.Payload.Dashboard.FullName, inputValues)

			// was a search path passed into the execute command?
			var opts []backend.BackendOption
			if request.Payload.SearchPath != nil || request.Payload.SearchPathPrefix != nil {
				opts = append(opts, backend.WithSearchPathConfig(backend.SearchPathConfig{
					SearchPath:       request.Payload.SearchPath,
					SearchPathPrefix: request.Payload.SearchPathPrefix,
				}))
			}
			err := dashboardexecute.Executor.ExecuteDashboard(ctx, sessionId, dashboard, inputValues, s.workspace, opts...)
			if err != nil {
				OutputError(ctx, sperr.WrapWithMessage(err, "error executing dashboard"))
			}
			slog.Debug("get_dashboard_metadata", "dashboard", request.Payload.Dashboard.FullName)
			payload, err := s.buildDashboardMetadataPayload(dashboard)
			if err != nil {
				OutputError(ctx, sperr.WrapWithMessage(err, "error building payload for get_metadata_details"))
			}
			_ = session.Write(payload)

		case "select_snapshot":
			snapshotName := request.Payload.Dashboard.FullName
			s.setDashboardForSession(sessionId, snapshotName, request.Payload.InputValues())
			snap, err := dashboardexecute.Executor.LoadSnapshot(ctx, sessionId, snapshotName, s.workspace)
			// TACTICAL- handle with error message
			error_helpers.FailOnError(err)
			// error handling???
			payload, err := buildDisplaySnapshotPayload(snap)
			// TACTICAL- handle with error message
			error_helpers.FailOnError(err)

			s.writePayloadToSession(sessionId, payload)
			OutputReady(ctx, fmt.Sprintf("Show snapshot complete: %s", snapshotName))
		case "input_changed":
			inputValues := request.Payload.InputValues()
			s.setDashboardInputsForSession(sessionId, inputValues)
			_ = dashboardexecute.Executor.OnInputChanged(ctx, sessionId, inputValues, request.Payload.ChangedInput)
		case "clear_dashboard":
			s.setDashboardInputsForSession(sessionId, nil)
			dashboardexecute.Executor.CancelExecutionForSession(ctx, sessionId)
		}
	}
}

func (s *Server) clearSession(ctx context.Context, session *melody.Session) {
	if strings.ToUpper(os.Getenv("DEBUG")) == "TRUE" {
		return
	}

	sessionId := s.getSessionId(session)

	dashboardexecute.Executor.CancelExecutionForSession(ctx, sessionId)

	s.deleteDashboardClient(sessionId)
}

// addSessionAndSendLoadState adds the session and sends it the load state of the workspace, as one step with
// broadcastLoadState: a client is sent the load states in the order they were read, so that a load state that ends a
// reload is never followed by an older one that says loading
func (s *Server) addSessionAndSendLoadState(session *melody.Session) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.dashboardClients[s.getSessionId(session)] = &DashboardClientInfo{Session: session}

	payload, err := buildWorkspaceLoadStatePayload(s.workspace.LoadState())
	if err != nil {
		slog.Warn("failed to build the workspace load state payload", "error", err.Error())
		return
	}
	_ = session.Write(payload)
}

// broadcastLoadState sends the load state of the workspace, as it is now, to every connected client
func (s *Server) broadcastLoadState() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	payload, err := buildWorkspaceLoadStatePayload(s.workspace.LoadState())
	if err != nil {
		slog.Warn("failed to build the workspace load state payload", "error", err.Error())
		return
	}
	for _, client := range s.dashboardClients {
		_ = client.Session.Write(payload)
	}
}

func (s *Server) setDashboardInputsForSession(sessionId string, inputs *dashboardexecute.InputValues) {
	dashboardClients := s.getDashboardClients()
	if sessionInfo, ok := dashboardClients[sessionId]; ok {
		sessionInfo.DashboardInputs = inputs
	}
}

func (s *Server) getSessionId(session *melody.Session) string {
	return fmt.Sprintf("%p", session)
}

// functions providing locked access to member properties

func (s *Server) setDashboardForSession(sessionId string, dashboardName string, inputs *dashboardexecute.InputValues) *DashboardClientInfo {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	dashboardClientInfo := s.dashboardClients[sessionId]
	dashboardClientInfo.Dashboard = &dashboardName
	dashboardClientInfo.DashboardInputs = inputs

	return dashboardClientInfo
}

func (s *Server) writePayloadToSession(sessionId string, payload []byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if sessionInfo, ok := s.dashboardClients[sessionId]; ok {
		_ = sessionInfo.Session.Write(payload)
	}
}

func (s *Server) getDashboardClients() map[string]*DashboardClientInfo {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.dashboardClients
}

func (s *Server) deleteDashboardClient(sessionId string) {
	s.mutex.Lock()
	delete(s.dashboardClients, sessionId)
	s.mutex.Unlock()
}

// resolve a resource from the name
func (s *Server) getResource(name string) modconfig.ModTreeItem {
	parsedResourceName, err := modconfig.ParseResourceName(name)
	if err != nil {
		slog.Warn("failed to parse changed resource name", "error", err.Error())
		return nil
	}

	resource, ok := s.workspace.GetResource(parsedResourceName)
	if !ok {
		slog.Warn("changed resource not found in workspace", "resource", name)
		return nil
	}
	return resource.(modconfig.ModTreeItem)
}

func getDashboardsInterestedInResourceChanges(dashboardsBeingWatched []string, existingChangedDashboardNames []string, changedItems []*modconfig.ModTreeItemDiffs) []string {
	var changedDashboardNames []string

	for _, changedItem := range changedItems {
		paths := changedItem.Item.GetPaths()
		for _, nodePath := range paths {
			for _, nodeName := range nodePath {
				resourceParts, _ := modconfig.ParseResourceName(nodeName)
				// We only care about changes from these resource types
				if !slices.Contains([]string{schema.BlockTypeDashboard, schema.BlockTypeBenchmark}, resourceParts.ItemType) {
					continue
				}

				if slices.Contains(existingChangedDashboardNames, nodeName) || slices.Contains(changedDashboardNames, nodeName) || !slices.Contains(dashboardsBeingWatched, nodeName) {
					continue
				}

				changedDashboardNames = append(changedDashboardNames, nodeName)
			}
		}
	}

	return changedDashboardNames
}
