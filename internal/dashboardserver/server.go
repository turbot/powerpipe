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

	typeHelpers "github.com/turbot/go-kit/types"
	"github.com/turbot/pipe-fittings/v2/backend"
	"github.com/turbot/pipe-fittings/v2/connection"
	"github.com/turbot/pipe-fittings/v2/error_helpers"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/pipe-fittings/v2/schema"
	"github.com/turbot/pipe-fittings/v2/steampipeconfig"
	"github.com/turbot/powerpipe/internal/dashboardevents"
	"github.com/turbot/powerpipe/internal/dashboardexecute"
	"github.com/turbot/powerpipe/internal/initialisation"
	"github.com/turbot/powerpipe/internal/timing"
	"github.com/turbot/powerpipe/internal/workspace"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
	"gopkg.in/olahol/melody.v1"
)

type Server struct {
	mutex                   *sync.Mutex
	dashboardClients        map[string]*DashboardClientInfo
	webSocket               *melody.Melody
	workspace               *workspace.PowerpipeWorkspace
	lazyWorkspace           *workspace.LazyWorkspace
	defaultDatabase         connection.ConnectionStringProvider
	defaultSearchPathConfig backend.SearchPathConfig
}

func NewServer(ctx context.Context, initData *initialisation.InitData, webSocket *melody.Melody) (*Server, error) {
	defer timing.Track("dashboardserver.NewServer")()
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

	// If lazy loading is enabled, set the lazy workspace reference
	if initData.IsLazy() {
		server.lazyWorkspace = initData.LazyWorkspace
		// Register as update listener for background resolution
		// This ensures the dashboard server broadcasts updated metadata when tags/titles are resolved
		server.lazyWorkspace.RegisterUpdateListener(server)

		// Setup file watcher for lazy mode
		// When files change (e.g., mod install), rebuild index and clear cache
		err := server.lazyWorkspace.SetupWatcher(ctx, func(c context.Context, e error) {
			if e != nil {
				slog.Error("File watcher error in lazy mode", "error", e)
				return
			}
			// Rebuild index when files change
			// This will start background resolution, which will call OnResolutionComplete()
			// when done, triggering the broadcast with resolved tags
			server.lazyWorkspace.HandleFileWatcherEvent(c)
		})
		if err != nil {
			return nil, err
		}
		OutputMessage(ctx, "WorkspaceEvents loaded (lazy mode with file watching)")
	} else {
		// Eager mode - setup standard file watcher
		err := w.SetupWatcher(ctx, func(c context.Context, e error) {})
		if err != nil {
			return nil, err
		}
		OutputMessage(ctx, "WorkspaceEvents loaded")
	}

	w.RegisterDashboardEventHandler(ctx, server.HandleDashboardEvent)

	return server, nil
}

// NewServerWithLazyWorkspace creates a server with lazy loading enabled.
// This allows available_dashboards to be served from the index without loading all resources.
func NewServerWithLazyWorkspace(ctx context.Context, lazyWorkspace *workspace.LazyWorkspace, defaultDatabase connection.ConnectionStringProvider, defaultSearchPathConfig backend.SearchPathConfig, webSocket *melody.Melody) (*Server, error) {
	defer timing.Track("dashboardserver.NewServerWithLazyWorkspace")()
	OutputWait(ctx, "Starting WorkspaceEvents Server (lazy mode)")

	var dashboardClients = make(map[string]*DashboardClientInfo)
	var mutex = &sync.Mutex{}

	server := &Server{
		mutex:                   mutex,
		dashboardClients:        dashboardClients,
		webSocket:               webSocket,
		lazyWorkspace:           lazyWorkspace,
		defaultDatabase:         defaultDatabase,
		defaultSearchPathConfig: defaultSearchPathConfig,
	}

	lazyWorkspace.RegisterDashboardEventHandler(ctx, server.HandleDashboardEvent)

	// Register as update listener for background resolution
	// This ensures the dashboard server broadcasts updated metadata when tags/titles are resolved
	lazyWorkspace.RegisterUpdateListener(server)

	// Setup file watcher for lazy mode
	// When files change (e.g., mod install), rebuild index and clear cache
	err := lazyWorkspace.SetupWatcher(ctx, func(c context.Context, e error) {
		if e != nil {
			slog.Error("File watcher error in lazy mode", "error", e)
			return
		}
		// Rebuild index when files change
		// This will start background resolution, which will call OnResolutionComplete()
		// when done, triggering the broadcast with resolved tags
		lazyWorkspace.HandleFileWatcherEvent(c)
	})
	if err != nil {
		return nil, err
	}

	OutputMessage(ctx, "WorkspaceEvents loaded (lazy mode with file watching)")

	return server, nil
}

// isLazyMode returns true if the server is using lazy loading.
func (s *Server) isLazyMode() bool {
	return s.lazyWorkspace != nil
}

// getActiveWorkspace returns the appropriate workspace based on mode.
func (s *Server) getActiveWorkspace() workspace.DashboardServerWorkspace {
	if s.isLazyMode() {
		return s.lazyWorkspace
	}
	return s.workspace
}

// buildAvailableDashboardsPayload builds the available dashboards payload
// using the index in lazy mode, or full resources in eager mode.
func (s *Server) buildAvailableDashboardsPayload() ([]byte, error) {
	if s.isLazyMode() {
		return buildAvailableDashboardsPayloadFromIndex(s.lazyWorkspace)
	}
	return buildAvailableDashboardsPayload(s.workspace.GetPowerpipeModResources())
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
			payload, payloadError = s.buildServerMetadataPayload(s.getActiveWorkspace().GetModResources(), &steampipeconfig.PipesMetadata{})
			if payloadError != nil {
				return
			}
			_ = s.webSocket.Broadcast(payload)

			// Emit available dashboards event
			payload, payloadError = s.buildAvailableDashboardsPayload()
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
						ws, err := s.getWorkspaceForExecution(ctx)
						if err != nil {
							OutputError(ctx, sperr.WrapWithMessage(err, "error loading workspace for execution"))
							continue
						}
						err = dashboardexecute.Executor.ExecuteDashboard(ctx, sessionId, changedResource, dashboardClientInfo.DashboardInputs, ws)
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
						ws, err := s.getWorkspaceForExecution(ctx)
						if err != nil {
							OutputError(ctx, sperr.WrapWithMessage(err, "error loading workspace for execution"))
							continue
						}
						err = dashboardexecute.Executor.ExecuteDashboard(ctx, sessionId, newDashboard, dashboardClientInfo.DashboardInputs, ws)
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

// OnResourceUpdated is called when a resource's metadata is updated during background resolution.
// For now, we don't need to take action on individual resource updates.
func (s *Server) OnResourceUpdated(resourceName string) {
	// Individual updates are not broadcasted to avoid spam
	// The final complete payload is sent in OnResolutionComplete
	slog.Debug("Resource metadata updated", "resource", resourceName)
}

// OnResolutionComplete is called when all background resolution is done.
// This broadcasts an updated available_dashboards payload with resolved tags/titles.
func (s *Server) OnResolutionComplete() {
	slog.Info("Background resolution complete - broadcasting updated dashboard metadata")

	// Build and broadcast updated available_dashboards payload
	// This ensures the frontend gets the resolved tags for proper grouping
	payload, err := s.buildAvailableDashboardsPayload()
	if err != nil {
		slog.Error("Failed to build available_dashboards payload after resolution", "error", err)
		return
	}

	// Broadcast to all connected clients
	if err := s.webSocket.Broadcast(payload); err != nil {
		slog.Error("Failed to broadcast available_dashboards after resolution", "error", err)
	}
}

func (s *Server) InitAsync(ctx context.Context) {
	defer timing.Track("Server.InitAsync")()
	go func() {
		// Return list of dashboards on connect
		s.webSocket.HandleConnect(func(session *melody.Session) {
			slog.Debug("client connected")
			s.addSession(session)
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
			payload, err := s.buildServerMetadataPayload(s.getActiveWorkspace().GetModResources(), &steampipeconfig.PipesMetadata{})
			if err != nil {
				OutputError(ctx, sperr.WrapWithMessage(err, "error building payload for get_metadata"))
			}
			_ = session.Write(payload)
		case "get_available_dashboards":
			payload, err := s.buildAvailableDashboardsPayload()
			if err != nil {
				OutputError(ctx, sperr.WrapWithMessage(err, "error building payload for get_available_dashboards"))
			}
			_ = session.Write(payload)
		case "select_dashboard":
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
			ws, err := s.getWorkspaceForExecution(ctx)
			if err != nil {
				OutputError(ctx, sperr.WrapWithMessage(err, "error loading workspace for execution"))
				return
			}
			// Get the dashboard/benchmark from the eager workspace for execution
			// This ensures query references are properly resolved
			parsedName, err := modconfig.ParseResourceName(request.Payload.Dashboard.FullName)
			if err != nil {
				OutputError(ctx, sperr.WrapWithMessage(err, "error parsing dashboard name"))
				return
			}
			resource, ok := ws.GetResource(parsedName)
			if !ok {
				slog.Warn("dashboard not found in workspace for execution", "name", request.Payload.Dashboard.FullName)
				return
			}
			dashboard := resource.(modconfig.ModTreeItem)
			err = dashboardexecute.Executor.ExecuteDashboard(ctx, sessionId, dashboard, inputValues, ws, opts...)
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
			ws, err := s.getWorkspaceForExecution(ctx)
			if err != nil {
				OutputError(ctx, sperr.WrapWithMessage(err, "error loading workspace for snapshot"))
				return
			}
			snap, err := dashboardexecute.Executor.LoadSnapshot(ctx, sessionId, snapshotName, ws)
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

func (s *Server) addSession(session *melody.Session) {
	sessionId := s.getSessionId(session)

	clientSession := &DashboardClientInfo{
		Session: session,
	}

	s.addDashboardClient(sessionId, clientSession)
}

func (s *Server) setDashboardInputsForSession(sessionId string, inputs *dashboardexecute.InputValues) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if sessionInfo, ok := s.dashboardClients[sessionId]; ok {
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

func (s *Server) addDashboardClient(sessionId string, clientSession *DashboardClientInfo) {
	s.mutex.Lock()
	s.dashboardClients[sessionId] = clientSession
	s.mutex.Unlock()
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

	slog.Debug("getResource: looking up resource",
		"inputName", name,
		"parsedMod", parsedResourceName.Mod,
		"parsedItemType", parsedResourceName.ItemType,
		"parsedName", parsedResourceName.Name,
		"isLazyMode", s.isLazyMode())

	resource, ok := s.getActiveWorkspace().GetResource(parsedResourceName)
	if !ok {
		slog.Warn("changed resource not found in workspace",
			"resource", name,
			"parsedMod", parsedResourceName.Mod,
			"parsedItemType", parsedResourceName.ItemType,
			"parsedName", parsedResourceName.Name)
		return nil
	}
	return resource.(modconfig.ModTreeItem)
}

// getWorkspaceForExecution returns the workspace to use for dashboard execution.
// In lazy mode, this triggers eager loading of the full workspace with proper
// reference resolution (needed for controls that reference queries).
func (s *Server) getWorkspaceForExecution(ctx context.Context) (*workspace.PowerpipeWorkspace, error) {
	if s.isLazyMode() {
		// In lazy mode, load the workspace eagerly for execution
		// This ensures controls have their query references properly resolved
		return s.lazyWorkspace.GetWorkspaceForExecution(ctx)
	}
	return s.workspace, nil
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
