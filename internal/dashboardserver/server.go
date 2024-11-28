package dashboardserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
	"log/slog"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/turbot/go-kit/helpers"
	typeHelpers "github.com/turbot/go-kit/types"
	"github.com/turbot/pipe-fittings/backend"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/schema"
	"github.com/turbot/pipe-fittings/steampipeconfig"
	"github.com/turbot/powerpipe/internal/dashboardevents"
	"github.com/turbot/powerpipe/internal/dashboardexecute"
	"github.com/turbot/powerpipe/internal/workspace"
	"gopkg.in/olahol/melody.v1"
)

type Server struct {
	mutex            *sync.Mutex
	dashboardClients map[string]*DashboardClientInfo
	webSocket        *melody.Melody
	workspace        *workspace.PowerpipeWorkspace
}

func NewServer(ctx context.Context, w *workspace.PowerpipeWorkspace, webSocket *melody.Melody) (*Server, error) {
	OutputWait(ctx, "Starting WorkspaceEvents Server")

	var dashboardClients = make(map[string]*DashboardClientInfo)

	var mutex = &sync.Mutex{}

	server := &Server{
		mutex:            mutex,
		dashboardClients: dashboardClients,
		webSocket:        webSocket,
		workspace:        w,
	}

	w.RegisterDashboardEventHandler(ctx, server.HandleDashboardEvent)

	err := w.SetupWatcher(ctx, func(c context.Context, e error) {})
	OutputMessage(ctx, "WorkspaceEvents loaded")

	return server, err
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
			payload, payloadError = buildServerMetadataPayload(s.workspace.GetModResources(), &steampipeconfig.PipesMetadata{})
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
				if helpers.StringSliceContains(dashboardsBeingWatched, dashboardName) {
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
			if helpers.StringSliceContains(changedDashboardNames, changedDashboard.Name) {
				continue
			}
			changedDashboardNames = append(changedDashboardNames, changedDashboard.Name)
		}

		for _, changedDashboardName := range changedDashboardNames {
			sessionMap := s.getDashboardClients()
			for sessionId, dashboardClientInfo := range sessionMap {
				if typeHelpers.SafeString(dashboardClientInfo.Dashboard) == changedDashboardName {

					if changedResource := s.getResource(changedDashboardName); changedResource != nil {
						_ = dashboardexecute.Executor.ExecuteDashboard(ctx, sessionId, changedResource, dashboardClientInfo.DashboardInputs, s.workspace)
					}
				}
			}

		}

		// Special case - if we previously had a workspace error, any previously existing dashboards
		// will come in here as new, so we need to check if any of those new dashboards are being watched.
		// If so, execute them
		for _, newDashboard := range newDashboards {
			if helpers.StringSliceContains(newDashboardNames, newDashboard.Name()) {
				continue
			}
			newDashboardNames = append(newDashboardNames, newDashboard.Name())
		}

		sessionMap := s.getDashboardClients()
		for _, newDashboardName := range newDashboardNames {
			for sessionId, dashboardClientInfo := range sessionMap {
				if typeHelpers.SafeString(dashboardClientInfo.Dashboard) == newDashboardName {
					if newDashboard := s.getResource(newDashboardName); newDashboard != nil {
						_ = dashboardexecute.Executor.ExecuteDashboard(ctx, sessionId, newDashboard, dashboardClientInfo.DashboardInputs, s.workspace)
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
			// TODO KAI verify we are ok to NOT send the cloud metadata here
			payload, err := buildServerMetadataPayload(s.workspace.GetModResources(), &steampipeconfig.PipesMetadata{})
			if err != nil {
				OutputError(ctx, sperr.WrapWithMessage(err, "error building payload for get_metadata"))
			}
			_ = session.Write(payload)
		case "get_dashboard_metadata":
			slog.Debug("get_dashboard_metadata", "dashboard", request.Payload.Dashboard.FullName)
			dashboard := s.getResource(request.Payload.Dashboard.FullName)
			if dashboard == nil {
				return
			}
			payload, err := buildDashboardMetadataPayload(ctx, dashboard, s.workspace)
			if err != nil {
				OutputError(ctx, sperr.WrapWithMessage(err, "error building payload for get_metadata_details"))
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
			s.setDashboardForSession(sessionId, request.Payload.Dashboard.FullName, request.Payload.InputValues)

			// was a search path passed into the execute command?
			var opts []backend.BackendOption
			if request.Payload.SearchPath != nil || request.Payload.SearchPathPrefix != nil {
				opts = append(opts, backend.WithSearchPathConfig(backend.SearchPathConfig{
					SearchPath:       request.Payload.SearchPath,
					SearchPathPrefix: request.Payload.SearchPathPrefix,
				}))
			}
			_ = dashboardexecute.Executor.ExecuteDashboard(ctx, sessionId, dashboard, request.Payload.InputValues, s.workspace, opts...)

		case "select_snapshot":
			snapshotName := request.Payload.Dashboard.FullName
			s.setDashboardForSession(sessionId, snapshotName, request.Payload.InputValues)
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
			s.setDashboardInputsForSession(sessionId, request.Payload.InputValues)
			_ = dashboardexecute.Executor.OnInputChanged(ctx, sessionId, request.Payload.InputValues, request.Payload.ChangedInput)
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
				if !helpers.StringSliceContains([]string{schema.BlockTypeDashboard, schema.BlockTypeBenchmark}, resourceParts.ItemType) {
					continue
				}

				if helpers.StringSliceContains(existingChangedDashboardNames, nodeName) || helpers.StringSliceContains(changedDashboardNames, nodeName) || !helpers.StringSliceContains(dashboardsBeingWatched, nodeName) {
					continue
				}

				changedDashboardNames = append(changedDashboardNames, nodeName)
			}
		}
	}

	return changedDashboardNames
}
