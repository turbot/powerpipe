import {
  DashboardDataMode,
  DashboardDataModeLive,
  DashboardSearch,
  SocketURLFactory,
} from "@powerpipe/types";
import { DashboardExecutionProvider } from "@powerpipe/hooks/useDashboardExecution";
import { DashboardInputsProvider } from "@powerpipe/hooks/useDashboardInputs";
import { DashboardPanelDetailProvider } from "@powerpipe/hooks/useDashboardPanelDetail";
import { DashboardSearchPathProvider } from "@powerpipe/hooks/useDashboardSearchPath";
import { DashboardSearchProvider } from "@powerpipe/hooks/useDashboardSearch";
import { DashboardStateProvider } from "./useDashboardState";
import { DashboardThemeProvider } from "@powerpipe/hooks/useDashboardTheme";

type DashboardProviderProps = {
  analyticsContext: any;
  breakpointContext: any;
  children: null | JSX.Element | JSX.Element[];
  componentOverrides?: {};
  dataMode?: DashboardDataMode;
  eventHooks?: {};
  rootPathname?: string;
  socketUrlFactory?: SocketURLFactory;
  stateDefaults?: {
    search?: DashboardSearch;
  };
  themeContext: any;
  versionMismatchCheck?: boolean;
};

const DashboardProvider = ({
  analyticsContext,
  breakpointContext,
  children,
  componentOverrides = {},
  dataMode = DashboardDataModeLive,
  eventHooks,
  rootPathname = "/",
  socketUrlFactory,
  stateDefaults = {},
  versionMismatchCheck = false,
  themeContext,
}: DashboardProviderProps) => {
  return (
    <DashboardThemeProvider themeContext={themeContext}>
      <DashboardSearchProvider defaultSearch={stateDefaults?.search}>
        <DashboardStateProvider
          analyticsContext={analyticsContext}
          breakpointContext={breakpointContext}
          componentOverrides={componentOverrides}
          dataMode={dataMode}
          rootPathname={rootPathname}
          stateDefaults={stateDefaults}
          versionMismatchCheck={versionMismatchCheck}
        >
          <DashboardInputsProvider>
            <DashboardSearchPathProvider>
              <DashboardPanelDetailProvider>
                <DashboardExecutionProvider
                  eventHooks={eventHooks}
                  socketUrlFactory={socketUrlFactory}
                >
                  {children}
                </DashboardExecutionProvider>
              </DashboardPanelDetailProvider>
            </DashboardSearchPathProvider>
          </DashboardInputsProvider>
        </DashboardStateProvider>
      </DashboardSearchProvider>
    </DashboardThemeProvider>
  );
};

export { DashboardProvider };
