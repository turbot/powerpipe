import path from "path";
import useDashboardWebSocket, {
  SocketActions,
} from "@powerpipe/hooks/useDashboardWebSocket";
import useDashboardWebSocketEventHandler from "@powerpipe/hooks/useDashboardWebSocketEventHandler";
import useDeepCompareEffect from "use-deep-compare-effect";
import useGlobalContextNavigate from "@powerpipe/hooks/useGlobalContextNavigate";
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
import { useDashboardDatetimeRange } from "@powerpipe/hooks/useDashboardDatetimeRange";
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
  autoNavigate = true,
}: {
  children: ReactNode;
  eventHooks?: {};
  socketUrlFactory?: () => Promise<string>;
  autoNavigate?: boolean;
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
  const { range, supportsTimeRange } = useDashboardDatetimeRange();
  const { searchPathPrefix, supportsSearchPath } = useDashboardSearchPath();
  const { search } = useGlobalContextNavigate();

  useEffect(() => {
    if (
      !!selectedDashboardRef.current &&
      selectedDashboardRef.current !== dashboard_name
    ) {
      closeSidePanel();
    }
  }, [dashboard_name]);

  const clearDashboard = (withNavigate = true) => {
    // Clear any existing executions
    sendMessage({
      action: SocketActions.CLEAR_DASHBOARD,
    });
    if (withNavigate) {
      navigate(`${rootPathname}${search ? `?${search}` : ""}`, {
        replace: true,
      });
    }
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
    clearDashboard(autoNavigate);
  }, [dispatch, autoNavigate, pathname, rootPathname, search]);

  const loadSnapshot = useCallback(
    (
      executionCompleteEvent: DashboardExecutionCompleteEvent,
      snapshotFileName: string,
      withNavigate = true,
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

      // Set the date range parameters
      if (executionCompleteEvent.snapshot.metadata?.datetime_range) {
        snapshotSearchParams.set(
          "datetime_range",
          JSON.stringify(
            executionCompleteEvent.snapshot.metadata?.datetime_range,
          ),
        );
      }

      // Set the search path prefix parameters
      if (executionCompleteEvent.snapshot.metadata?.search_path_prefix) {
        snapshotSearchParams.set(
          "search_path_prefix",
          (
            executionCompleteEvent.snapshot.metadata?.search_path_prefix || []
          ).join(","),
        );
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
      if (withNavigate) {
        navigate(
          path.join(
            rootPathname,
            `/snapshot/${snapshotFileName}${snapshotSearchParamsString ? `?${snapshotSearchParamsString}` : ""}`,
          ),
        );
      }

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
      navigate(`${rootPathname}${search ? `?${search}` : ""}`, {
        replace: true,
      });
      setLastChangedInput(null);
      return;
    }

    const dashboardMessage: any = {
      payload: {
        dashboard: {
          full_name: dashboard.full_name,
        },
        inputs,
      },
    };

    if (supportsTimeRange && range.from) {
      dashboardMessage.payload.datetime_range =
        dashboardMessage.payload.datetime_range || {};
      dashboardMessage.payload.datetime_range.from = range.from;
    }
    if (supportsTimeRange && range.to) {
      dashboardMessage.payload.datetime_range =
        dashboardMessage.payload.datetime_range || {};
      dashboardMessage.payload.datetime_range.to = range.to;
    }

    if (supportsSearchPath && !!searchPathPrefix.length) {
      dashboardMessage.payload.search_path_prefix = searchPathPrefix;
    }

    if (
      lastChangedInput &&
      selectedDashboardRef.current === dashboardFullName
    ) {
      dashboardMessage.action = SocketActions.INPUT_CHANGED;
      dashboardMessage.payload.changed_input = lastChangedInput;
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
    supportsTimeRange,
    supportsSearchPath,
    dataMode,
    dashboard_name,
    inputs,
    range.from,
    range.to,
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
