import path from "path";
import useDashboardWebSocket, {
  SocketActions,
} from "@powerpipe/hooks/useDashboardWebSocket";
import useDashboardWebSocketEventHandler from "@powerpipe/hooks/useDashboardWebSocketEventHandler";
import useDeepCompareEffect from "use-deep-compare-effect";
import {
  createContext,
  ReactNode,
  useCallback,
  useContext,
  useEffect,
  useRef,
} from "react";
import {
  DashboardActions,
  DashboardDataModeCLISnapshot,
  DashboardDataModeLive,
  DashboardExecutionCompleteEvent,
} from "@powerpipe/types";
import { useDashboardInputs } from "@powerpipe/hooks/useDashboardInputs";
import { useDashboardPanelDetail } from "@powerpipe/hooks/useDashboardPanelDetail";
import { useDashboardSearchPath } from "@powerpipe/hooks/useDashboardSearchPath";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { useLocation, useNavigate, useParams } from "react-router-dom";

interface IDashboardExecutionContext {
  loadSnapshot: (
    executionCompleteEvent: DashboardExecutionCompleteEvent,
    snapshotFileName: string,
  ) => void;
}

const DashboardExecutionContext =
  createContext<IDashboardExecutionContext | null>(null);

export const DashboardExecutionProvider = ({
  children,
  eventHooks,
  socketUrlFactory,
}: {
  children: ReactNode;
  eventHooks?: {};
  socketUrlFactory?: () => Promise<string>;
}) => {
  const navigate = useNavigate();
  const { pathname } = useLocation();
  const { dashboard_name } = useParams();
  const {
    availableDashboardsLoaded,
    dashboards,
    dataMode,
    dispatch,
    rootPathname,
    selectedDashboard,
    snapshotFileName,
  } = useDashboardState();
  const selectedDashboardRef = useRef<string | null>(
    selectedDashboard?.full_name || null,
  );
  const { selectPanel, closeSidePanel } = useDashboardPanelDetail();
  const { eventHandler } = useDashboardWebSocketEventHandler(
    dispatch,
    eventHooks,
  );
  const { send: sendMessage } = useDashboardWebSocket(
    dataMode,
    eventHandler,
    socketUrlFactory,
  );
  const { inputs, lastChangedInput, setLastChangedInput } =
    useDashboardInputs();
  const { searchPathPrefix } = useDashboardSearchPath();

  useEffect(() => {
    if (
      !!selectedDashboardRef.current &&
      selectedDashboardRef.current !== dashboard_name
    ) {
      closeSidePanel();
    }
  }, [dashboard_name]);

  const clearDashboard = () => {
    // Clear any existing executions
    sendMessage({
      action: SocketActions.CLEAR_DASHBOARD,
    });
    navigate(
      `../${!!searchPathPrefix.length ? `?search_path_prefix=${searchPathPrefix}` : ""}`,
      { replace: true },
    );
    setLastChangedInput(null);
    closeSidePanel();
    dispatch({
      type: DashboardActions.CLEAR_DASHBOARD,
    });
  };

  useEffect(() => {
    if (pathname !== rootPathname) {
      return;
    }
    clearDashboard();
  }, [dispatch, pathname, searchPathPrefix]);

  const loadSnapshot = useCallback(
    (
      executionCompleteEvent: DashboardExecutionCompleteEvent,
      snapshotFileName: string,
    ) => {
      // Clear any existing executions
      sendMessage({
        action: SocketActions.CLEAR_DASHBOARD,
      });

      // Build the inputs search params
      const snapshotSearchParams = new URLSearchParams();
      for (const [key, value] of Object.entries(
        executionCompleteEvent.snapshot.inputs || {},
      )) {
        snapshotSearchParams.set(key, value);
      }

      // Build the filter & group search params
      const filtersByPanel = {};
      const groupingsByPanel = {};
      const tableSettingsByPanel = {};
      for (const [panel, panelViewSettings] of Object.entries(
        executionCompleteEvent.snapshot.metadata?.view || {},
      )) {
        if (panelViewSettings.filter_by) {
          filtersByPanel[panel] = panelViewSettings.filter_by;
        }
        if (panelViewSettings.group_by) {
          groupingsByPanel[panel] = panelViewSettings.group_by;
        }
        if (panelViewSettings.table) {
          tableSettingsByPanel[panel] = panelViewSettings.table;
        }
      }
      if (Object.keys(filtersByPanel).length) {
        snapshotSearchParams.set("where", JSON.stringify(filtersByPanel));
      }
      if (Object.keys(groupingsByPanel).length) {
        snapshotSearchParams.set("grouping", JSON.stringify(groupingsByPanel));
      }
      if (Object.keys(tableSettingsByPanel).length) {
        snapshotSearchParams.set("table", JSON.stringify(tableSettingsByPanel));
      }

      const snapshotSearchParamsString = snapshotSearchParams.toString();

      // Navigate to the snapshot page
      navigate(
        path.join(
          rootPathname,
          `/snapshot/${snapshotFileName}${snapshotSearchParamsString ? `?${snapshotSearchParamsString}` : ""}`,
        ),
      );

      dispatch({
        type: DashboardActions.LOAD_SNAPSHOT,
        executionCompleteEvent,
        snapshotFileName,
      });
    },
    [dispatch, rootPathname, navigate, sendMessage],
  );

  const executeDashboard = (dashboardFullName: string | null | undefined) => {
    if (
      dataMode === DashboardDataModeCLISnapshot &&
      snapshotFileName &&
      pathname.startsWith("/snapshot/")
    ) {
      return;
    } else if (dataMode === DashboardDataModeCLISnapshot && dashboard_name) {
      dispatch({
        type: DashboardActions.SET_DATA_MODE,
        dataMode: DashboardDataModeLive,
      });
      setLastChangedInput(null);
      return;
    } else if (!dashboardFullName) {
      setLastChangedInput(null);
      return;
    }

    const dashboard = dashboards.find(
      (dashboard) => dashboard.full_name === dashboardFullName,
    );

    selectPanel(null);

    // If the dashboard we're viewing no longer exists, go back to the main page
    if (!dashboard) {
      dispatch({
        type: DashboardActions.SELECT_DASHBOARD,
        dashboard: null,
      });
      navigate(
        `../${!!searchPathPrefix.length ? `?search_path_prefix=${searchPathPrefix}` : ""}`,
        { replace: true },
      );
      setLastChangedInput(null);
      return;
    }

    const { "input.detection_range": detectionRange, ...rest } = inputs || {};
    let detectionFrom, detectionTo;
    if (detectionRange) {
      try {
        const parsed = JSON.parse(detectionRange);
        detectionFrom = parsed.from;
        detectionTo = parsed.to;
      } catch (err) {
        console.error("Parse error", err);
      }
    }

    const dashboardMessage: any = {
      payload: {
        dashboard: {
          full_name: dashboard.full_name,
        },
        input_values: { inputs: rest },
      },
    };

    if (detectionFrom) {
      dashboardMessage.payload.input_values.detection_time_ranges =
        dashboardMessage.payload.input_values.detection_time_ranges || {};
      dashboardMessage.payload.input_values.detection_time_ranges.from =
        detectionFrom;
    }
    if (detectionTo) {
      dashboardMessage.payload.input_values.detection_time_ranges =
        dashboardMessage.payload.input_values.detection_time_ranges || {};
      dashboardMessage.payload.input_values.detection_time_ranges.to =
        detectionTo;
    }

    if (!!searchPathPrefix.length) {
      dashboardMessage.payload.search_path_prefix = searchPathPrefix;
    }

    if (
      lastChangedInput &&
      selectedDashboardRef.current === dashboardFullName
    ) {
      dashboardMessage.action = SocketActions.INPUT_CHANGED;
      dashboardMessage.changed_input = lastChangedInput;
      sendMessage(dashboardMessage);
    } else {
      selectedDashboardRef.current = dashboard.full_name;
      // Ensure the dashboard is selected
      dispatch({
        type: DashboardActions.SELECT_DASHBOARD,
        dashboard,
      });
      // Clear any existing executions
      sendMessage({
        action: SocketActions.CLEAR_DASHBOARD,
      });
      dashboardMessage.action = SocketActions.SELECT_DASHBOARD;
      sendMessage(dashboardMessage);
    }
  };

  useDeepCompareEffect(() => {
    // We don't need to "execute" if we're in snapshot mode
    if (!availableDashboardsLoaded) {
      return;
    }
    executeDashboard(dashboard_name);
  }, [
    availableDashboardsLoaded,
    dataMode,
    dashboard_name,
    inputs,
    searchPathPrefix,
  ]);

  return (
    <DashboardExecutionContext.Provider value={{ loadSnapshot }}>
      {children}
    </DashboardExecutionContext.Provider>
  );
};

export const useDashboardExecution = () => {
  const context = useContext(DashboardExecutionContext);
  if (!context) {
    throw new Error(
      "useDashboardExecution must be used within a DashboardExecutionProvider",
    );
  }
  return context;
};
