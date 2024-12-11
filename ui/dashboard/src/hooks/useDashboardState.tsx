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
  DashboardCliMode,
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

const DashboardContext = createContext<IDashboardContext | null>(null);

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
      // TODO remove once all workspaces are running Powerpipe as this is actually SERVER_METADATA in Powerpipe
      if (state.cliMode === "steampipe") {
        return {
          ...state,
          metadata: {
            mod: {},
            ...action.metadata,
          },
        };
      }
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
    case DashboardActions.SET_CLI_MODE:
      return {
        ...state,
        cliMode: action.cli_mode,
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
    case DashboardActions.SELECT_PANEL:
      return { ...state, selectedPanel: action.panel };
    case DashboardActions.SET_DATA_MODE:
      const newState = {
        ...state,
        dataMode: action.dataMode,
      };
      if (action.dataMode === DashboardDataModeCLISnapshot) {
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
          dataMode: DashboardDataModeCLISnapshot,
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
        dataMode: DashboardDataModeLive,
        dashboard: null,
        execution_id: null,
        panelsMap: {},
        snapshot: null,
        snapshotFileName: null,
        snapshotId: null,
        state: null,
        selectedDashboard: action.dashboard,
        selectedPanel: null,
        lastChangedInput: null,
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
    case DashboardActions.WORKSPACE_ERROR:
      return { ...state, error: action.error };
    case DashboardActions.SHOW_CUSTOMIZE_BENCHMARK_PANEL:
      return {
        ...state,
        filterAndGroupControlPanel: action.panel_name,
      };
    case DashboardActions.HIDE_CUSTOMIZE_BENCHMARK_PANEL: {
      const { filterAndGroupControlPanel, ...rest } = state;
      return {
        ...rest,
      };
    }
    default:
      console.warn(`Unsupported action ${action.type}`, action);
      return state;
  }
};

const getInitialState = (searchParams, defaults: any = {}) => {
  return {
    cliMode: defaults.cliMode || "powerpipe",
    versionMismatchCheck: defaults.versionMismatchCheck,
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
    selectedPanel: null,
    selectedDashboard: null,
    snapshot: null,
    lastChangedInput: null,

    execution_id: null,

    progress: 0,
  };
};

type DashboardStateProviderProps = {
  analyticsContext: any;
  breakpointContext: any;
  children: ReactNode;
  componentOverrides?: {};
  dataMode: DashboardDataMode;
  stateDefaults?: {
    cliMode?: DashboardCliMode;
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
  stateDefaults,
  versionMismatchCheck,
}: DashboardStateProviderProps) => {
  const {
    setMetadata: setAnalyticsMetadata,
    setSelectedDashboard: setAnalyticsSelectedDashboard,
  } = analyticsContext;
  const components = buildComponentsMap(componentOverrides);
  const initialState = useMemo(() => {
    const searchParams = new URLSearchParams(window.location.search);
    return getInitialState(searchParams, {
      ...stateDefaults,
      dataMode,
      versionMismatchCheck,
    });
  }, []);
  const [state, dispatchInner] = useReducer(reducer, initialState);
  useDashboardVersionCheck(state);
  const dispatch = useCallback((action) => {
    // console.log(action.type, action);
    dispatchInner(action);
  }, []);

  console.log(state);

  // Alert analytics
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
    throw new Error(
      "useDashboardExecution must be used within a DashboardExecutionProvider",
    );
  }
  return context;
};
