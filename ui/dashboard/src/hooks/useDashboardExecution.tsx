import useDashboardWebSocket, {
  SocketActions,
} from "@powerpipe/hooks/useDashboardWebSocket";
import useDashboardWebSocketEventHandler from "@powerpipe/hooks/useDashboardWebSocketEventHandler";
import useDeepCompareEffect from "use-deep-compare-effect";
import { createContext, ReactNode, useContext, useEffect } from "react";
import {
  DashboardActions,
  DashboardDataModeCLISnapshot,
  DashboardDataModeLive,
} from "@powerpipe/types";
import { useDashboardSearchPath } from "@powerpipe/hooks/useDashboardSearchPath";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { useLocation, useNavigate, useParams } from "react-router-dom";
import { useDashboardInputs } from "@powerpipe/hooks/useDashboardInputs";

interface IDashboardExecutionContext {
  executeDashboard: (dashboardFullName: string | null | undefined) => void;
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
  const { dashboards, dataMode, dispatch } = useDashboardState();
  const { eventHandler } = useDashboardWebSocketEventHandler(
    dispatch,
    eventHooks,
  );
  const { send: sendMessage } = useDashboardWebSocket(
    dataMode,
    eventHandler,
    socketUrlFactory,
  );
  const { inputs } = useDashboardInputs();
  const { searchPathPrefix } = useDashboardSearchPath();

  console.log({ inputs, searchPathPrefix });

  useEffect(() => {
    if (pathname !== "/" || dataMode === DashboardDataModeLive) {
      return;
    }
    dispatch({
      type: DashboardActions.SET_DATA_MODE,
      dataMode: DashboardDataModeLive,
    });
  }, [dispatch, pathname, dataMode]);

  const executeDashboard = (dashboardFullName: string | null | undefined) => {
    const dashboard = dashboards.find(
      (dashboard) => dashboard.full_name === dashboardFullName,
    );

    // If the dashboard we're viewing no longer exists, go back to the main page
    if (!dashboard) {
      dispatch({
        type: DashboardActions.SELECT_DASHBOARD,
        dashboard: null,
        recordInputsHistory: false,
      });
      navigate("../", { replace: true });
      return;
    }
    dispatch({
      type: DashboardActions.SELECT_DASHBOARD,
      dashboard,
    });

    // Clear any existing executions
    sendMessage({
      action: SocketActions.CLEAR_DASHBOARD,
    });

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

    const selectDashboardMessage: any = {
      action: SocketActions.SELECT_DASHBOARD,
      payload: {
        dashboard: {
          full_name: dashboard.full_name,
        },
        input_values: { inputs: rest },
      },
    };

    if (detectionFrom) {
      selectDashboardMessage.payload.input_values.detection_time_ranges =
        selectDashboardMessage.payload.input_values.detection_time_ranges || {};
      selectDashboardMessage.payload.input_values.detection_time_ranges.from =
        detectionFrom;
    }
    if (detectionTo) {
      selectDashboardMessage.payload.input_values.detection_time_ranges =
        selectDashboardMessage.payload.input_values.detection_time_ranges || {};
      selectDashboardMessage.payload.input_values.detection_time_ranges.to =
        detectionTo;
    }

    if (!!searchPathPrefix.length) {
      selectDashboardMessage.payload.search_path_prefix = searchPathPrefix;
    }

    console.log("Executing", {
      dashboardFullName,
      selectDashboardMessage,
      inputs,
      dashboard,
    });
    sendMessage(selectDashboardMessage);
  };

  useDeepCompareEffect(() => {
    // We don't need to "execute" if we're in snapshot mode
    if (dataMode === DashboardDataModeCLISnapshot) {
      return;
    }
    executeDashboard(dashboard_name);
  }, [dataMode, dashboard_name, inputs]);

  return (
    <DashboardExecutionContext.Provider value={{ executeDashboard }}>
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
