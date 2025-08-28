import useDashboardVersionCheck from "./useDashboardVersionCheck";
import { buildComponentsMap } from "@powerpipe/components";
import {
  buildDashboards,
  buildPanelsLog,
  updatePanelsLogFromCompletedPanels,
  updateSelectedDashboard,
  wrapDefinitionInArtificialDashboard,
} from "@powerpipe/utils/state";
import {
  controlsUpdatedEventHandler,
  leafNodesUpdatedEventHandler,
  migratePanelStatuses,
} from "@powerpipe/utils/dashboardEventHandlers";
import {
  createContext,
  ReactNode,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useReducer,
} from "react";
import {
  DashboardActions,
  DashboardDataMode,
  DashboardDataModeCLISnapshot,
  DashboardDataModeCloudSnapshot,
  DashboardDataModeLive,
  DashboardSearch,
  IDashboardContext,
} from "@powerpipe/types";
import {
  EXECUTION_SCHEMA_VERSION_20220929,
  EXECUTION_SCHEMA_VERSION_20221222,
} from "@powerpipe/constants/versions";
import {
  ExecutionCompleteSchemaMigrator,
  ExecutionStartedSchemaMigrator,
} from "@powerpipe/utils/schema";

export const DashboardContext = createContext<IDashboardContext | null>(null);

const reducer = (state: IDashboardContext, action) => {
  switch (action.type) {
    case DashboardActions.SERVER_METADATA:
      return {
        ...state,
        metadata: {
          mod: {},
          ...action.metadata,
        },
      };
    case DashboardActions.DASHBOARD_METADATA:
      if (!state.selectedDashboard?.full_name) {
        return state;
      }
      return {
        ...state,
        dashboardsMetadata: {
          ...state.dashboardsMetadata,
          [state.selectedDashboard?.full_name]: action.metadata,
        },
      };
    case DashboardActions.AVAILABLE_DASHBOARDS:
      const { dashboards, dashboardsMap } = buildDashboards(
        action.dashboards,
        action.benchmarks,
      );
      const selectedDashboard = updateSelectedDashboard(
        state.selectedDashboard,
        dashboards,
      );
      return {
        ...state,
        error: null,
        availableDashboardsLoaded: true,
        dashboards,
        dashboardsMap,
        selectedDashboard:
          state.dataMode === DashboardDataModeCLISnapshot ||
          state.dataMode === DashboardDataModeCloudSnapshot
            ? state.selectedDashboard
            : selectedDashboard,
        dashboard:
          state.dataMode === DashboardDataModeCLISnapshot ||
          state.dataMode === DashboardDataModeCloudSnapshot
            ? state.dashboard
            : selectedDashboard &&
                state.dashboard &&
                state.dashboard.name === selectedDashboard.full_name
              ? state.dashboard
              : null,
      };
    case DashboardActions.EXECUTION_STARTED: {
      const rootLayoutPanel = action.layout;
      const rootPanel = action.panels[rootLayoutPanel.name];
      let dashboard;
      // For benchmarks and controls that are run directly from a mod,
      // we need to wrap these in an artificial dashboard, so we can treat
      // it just like any other dashboard
      if (rootPanel.panel_type !== "dashboard") {
        dashboard = wrapDefinitionInArtificialDashboard(
          rootPanel,
          action.layout,
        );
      } else {
        dashboard = {
          ...rootPanel,
          ...action.layout,
        };
      }

      const eventMigrator = new ExecutionStartedSchemaMigrator();
      const migratedEvent = eventMigrator.toLatest(action);

      return {
        ...state,
        error: null,
        panelsLog: buildPanelsLog(
          migratedEvent.panels,
          migratedEvent.start_time,
        ),
        panelsMap: migratedEvent.panels,
        dashboard,
        execution_id: migratedEvent.execution_id,
        progress: 0,
        snapshot: null,
        state: "running",
      };
    }
    case DashboardActions.EXECUTION_COMPLETE: {
      // If we're in live mode and not expecting execution events for this ID
      if (
        state.dataMode === DashboardDataModeLive &&
        action.execution_id !== state.execution_id
      ) {
        return state;
      }

      const eventMigrator = new ExecutionCompleteSchemaMigrator();
      const migratedEvent = eventMigrator.toLatest(action);
      const layout = migratedEvent.snapshot.layout;
      const panels = migratedEvent.snapshot.panels;
      const rootLayoutPanel = migratedEvent.snapshot.layout;
      const rootPanel = panels[rootLayoutPanel.name];
      let dashboard;

      if (rootPanel.panel_type !== "dashboard") {
        dashboard = wrapDefinitionInArtificialDashboard(rootPanel, layout);
      } else {
        dashboard = {
          ...rootPanel,
          ...layout,
        };
      }

      const panelsMap = migratePanelStatuses(panels, action.schema_version);

      // Replace the whole dashboard as this event contains everything
      return {
        ...state,
        error: null,
        panelsLog: updatePanelsLogFromCompletedPanels(
          state.panelsLog,
          panels,
          migratedEvent.snapshot.end_time,
        ),
        panelsMap,
        dashboard,
        progress: 100,
        snapshot: migratedEvent.snapshot,
        state: "complete",
      };
    }
    case DashboardActions.EXECUTION_ERROR:
      return { ...state, error: action.error, progress: 100, state: "error" };
    case DashboardActions.CONTROLS_UPDATED:
      return controlsUpdatedEventHandler(action, state);
    case DashboardActions.LEAF_NODES_COMPLETE:
      return leafNodesUpdatedEventHandler(
        action,
        EXECUTION_SCHEMA_VERSION_20220929,
        state,
      );
    case DashboardActions.LEAF_NODES_UPDATED:
      return leafNodesUpdatedEventHandler(
        action,
        EXECUTION_SCHEMA_VERSION_20221222,
        state,
      );
    case DashboardActions.CLEAR_DASHBOARD:
      return {
        ...state,
        dataMode:
          state.dataMode === "cloud_snapshot"
            ? state.dataMode
            : DashboardDataModeLive,
        dashboard: null,
        error: null,
        execution_id: null,
        panelsLog: {},
        panelsMap: {},
        selectedDashboard: null,
        snapshot: null,
        progress: 0,
      };
    case DashboardActions.LOAD_SNAPSHOT:
      const { executionCompleteEvent, snapshotFileName } = action;
      const layout = executionCompleteEvent.snapshot.layout;
      const panels = executionCompleteEvent.snapshot.panels;
      const rootPanel = panels[layout.name];
      let dashboard;

      if (rootPanel.panel_type !== "dashboard") {
        dashboard = wrapDefinitionInArtificialDashboard(rootPanel, layout);
      } else {
        dashboard = {
          ...rootPanel,
          ...layout,
        };
      }

      const panelsMap = migratePanelStatuses(
        panels,
        executionCompleteEvent.schema_version,
      );

      return {
        ...state,
        dashboard,
        dataMode:
          state.dataMode === DashboardDataModeCloudSnapshot
            ? DashboardDataModeCloudSnapshot
            : DashboardDataModeCLISnapshot,
        error: null,
        panelsLog: updatePanelsLogFromCompletedPanels(
          state.panelsLog,
          panels,
          executionCompleteEvent.snapshot.end_time,
        ),
        panelsMap,
        progress: 100,
        snapshot: executionCompleteEvent.snapshot,
        snapshotFileName,
        state: "complete",
      };
    case DashboardActions.SET_DATA_MODE:
      const newState = {
        ...state,
        dataMode: action.dataMode,
      };
      if (
        action.dataMode === DashboardDataModeCLISnapshot ||
        action.dataMode === DashboardDataModeCloudSnapshot
      ) {
        newState.snapshotFileName = action.snapshotFileName;
      } else if (
        state.dataMode !== DashboardDataModeLive &&
        action.dataMode === DashboardDataModeLive
      ) {
        newState.snapshot = null;
        newState.snapshotFileName = null;
        newState.snapshotId = null;
      }
      return newState;
    case DashboardActions.SET_DASHBOARD:
      return {
        ...state,
        dashboard: action.dashboard,
      };
    case DashboardActions.SELECT_DASHBOARD:
      if (action.dashboard && action.dashboard.type === "snapshot") {
        return {
          ...state,
          dataMode:
            state.dataMode === DashboardDataModeCloudSnapshot
              ? DashboardDataModeCloudSnapshot
              : DashboardDataModeCLISnapshot,
          selectedDashboard: action.dashboard,
        };
      }

      if (
        action.dataMode === DashboardDataModeCLISnapshot ||
        action.dataMode === DashboardDataModeCloudSnapshot
      ) {
        return {
          ...state,
          dataMode: action.dataMode,
          selectedDashboard: action.dashboard,
        };
      }

      return {
        ...state,
        dataMode:
          state.dataMode === DashboardDataModeCloudSnapshot
            ? DashboardDataModeCloudSnapshot
            : DashboardDataModeLive,
        dashboard: null,
        execution_id: null,
        panelsMap: {},
        snapshot: null,
        snapshotFileName: null,
        snapshotId: null,
        state: null,
        selectedDashboard: action.dashboard,
      };
    case DashboardActions.SET_DASHBOARD_TAG_KEYS:
      return {
        ...state,
        dashboardTags: {
          ...state.dashboardTags,
          keys: action.keys,
        },
      };
    case DashboardActions.SET_SNAPSHOT_METADATA_LOADED:
      return { ...state, snapshot_metadata_loaded: true };
    case DashboardActions.SET_OVERLAY_VISIBLE:
      return { ...state, overlayVisible: action.value };
    case DashboardActions.WORKSPACE_ERROR:
      return { ...state, error: action.error };
    default:
      console.warn(`Unsupported action ${action.type}`, action);
      return state;
  }
};

const getInitialState = (defaults: any = {}) => {
  return {
    availableDashboardsLoaded: false,
    metadata: null,
    dashboards: [],
    dashboardTags: {
      keys: [],
    },
    dataMode: defaults.dataMode || DashboardDataModeLive,
    snapshotId: defaults.snapshotId ? defaults.snapshotId : null,
    error: null,
    panelsLog: {},
    panelsMap: {},
    dashboard: null,
    dashboardsMetadata: {},
    overlayVisible: false,
    rootPathname: defaults.rootPathname || "/",
    selectedDashboard: null,
    snapshot: null,
    execution_id: null,
    progress: 0,
    versionMismatchCheck: defaults.versionMismatchCheck,
  };
};

type DashboardStateProviderProps = {
  analyticsContext: any;
  breakpointContext: any;
  children: ReactNode;
  componentOverrides?: {};
  dataMode: DashboardDataMode;
  rootPathname: string;
  stateDefaults?: {
    search?: DashboardSearch;
  };
  versionMismatchCheck: boolean;
};

export const DashboardStateProvider = ({
  children,
  analyticsContext,
  breakpointContext,
  componentOverrides = {},
  dataMode,
  rootPathname,
  stateDefaults,
  versionMismatchCheck,
}: DashboardStateProviderProps) => {
  const {
    setMetadata: setAnalyticsMetadata,
    setSelectedDashboard: setAnalyticsSelectedDashboard,
  } = analyticsContext;
  const components = buildComponentsMap(componentOverrides);
  const initialState = useMemo(() => {
    return getInitialState({
      ...stateDefaults,
      dataMode,
      rootPathname,
      versionMismatchCheck,
    });
  }, []);
  const [state, dispatchInner] = useReducer(reducer, initialState);
  useDashboardVersionCheck(state);
  const dispatch = useCallback((action) => {
    // console.log(action.type, action);
    dispatchInner(action);
  }, []);

  // Set up analytics
  useEffect(() => {
    setAnalyticsMetadata(state.metadata);
  }, [state.metadata, setAnalyticsMetadata]);

  useEffect(() => {
    setAnalyticsSelectedDashboard(state.selectedDashboard);
  }, [state.selectedDashboard, setAnalyticsSelectedDashboard]);

  return (
    <DashboardContext.Provider
      value={{ ...state, breakpointContext, components, dispatch }}
    >
      {children}
    </DashboardContext.Provider>
  );
};

export const useDashboardState = () => {
  const context = useContext(DashboardContext);
  if (!context) {
    throw new Error("useDashboardState must be used within a DashboardContext");
  }
  return context;
};
