import "@powerpipe/utils/registerComponents";
import Dashboard from "@powerpipe/components/dashboards/layout/Dashboard";
import { buildComponentsMap } from "@powerpipe/components";
import { DashboardContext } from "@powerpipe/hooks/useDashboardState";
import {
  DashboardDataModeLive,
  DashboardPanelType,
  DashboardRunState,
  DashboardSearch,
} from "@powerpipe/types";
import { DashboardInputsProvider } from "@powerpipe/hooks/useDashboardInputs";
import { DashboardDatetimeRangeProvider } from "@powerpipe/hooks/useDashboardDatetimeRange";
import { DashboardPanelDetailProvider } from "@powerpipe/hooks/useDashboardPanelDetail";
import { DashboardSearchPathProvider } from "@powerpipe/hooks/useDashboardSearchPath";
import { DashboardSearchProvider } from "@powerpipe/hooks/useDashboardSearch";
import { DashboardThemeProvider } from "@powerpipe/hooks/useDashboardTheme";
import { noop } from "@powerpipe/utils/func";
import { useStorybookTheme } from "@powerpipe/hooks/useStorybookTheme";

type PanelStoryDecoratorProps = {
  definition: any;
  panelType: DashboardPanelType;
  panels?: {
    [key: string]: any;
  };
  status?: DashboardRunState;
  additionalProperties?: {
    [key: string]: any;
  };
};

const stubDashboardSearch: DashboardSearch = {
  value: "",
  groupBy: { value: "mod", tag: null },
};

export const PanelStoryDecorator = ({
  definition = {},
  panels = {},
  panelType,
  status = "complete",
  additionalProperties = {},
}: PanelStoryDecoratorProps) => {
  const { theme, wrapperRef } = useStorybookTheme();
  const { properties, ...rest } = definition;

  const newPanel = {
    ...rest,
    name: `${panelType}.story`,
    panel_type: panelType,
    properties: {
      ...(properties || {}),
      ...additionalProperties,
    },
    sql: "storybook",
    status,
  };

  return (
    <DashboardThemeProvider
      themeContext={{
        theme,
        setTheme: noop,
        wrapperRef,
      }}
    >
      <DashboardSearchProvider defaultSearch={stubDashboardSearch}>
        <DashboardContext.Provider
          value={{
            versionMismatchCheck: false,
            metadata: {
              mod: {
                title: "Storybook",
                full_name: "mod.storybook",
                short_name: "storybook",
              },
              installed_mods: {},
              telemetry: "none",
            },
            availableDashboardsLoaded: true,
            dataMode: DashboardDataModeLive,
            snapshotId: null,
            dispatch: noop,
            error: null,
            dashboards: [],
            dashboardsMap: {},
            dashboardsMetadata: {
              "storybook.dashboard.storybook_dashboard_wrapper": {
                supports_search_path: false,
              },
            },
            selectedDashboard: {
              title: "Storybook Dashboard Wrapper",
              full_name: "storybook.dashboard.storybook_dashboard_wrapper",
              short_name: "storybook_dashboard_wrapper",
              type: "dashboard",
              tags: {},
              mod_full_name: "mod.storybook",
              is_top_level: true,
            },
            execution_id: null,
            panelsLog: {},
            panelsMap: {
              [newPanel.name]: newPanel,
              ...panels,
            },
            dashboard: {
              artificial: false,
              name: "storybook.dashboard.storybook_dashboard_wrapper",
              children: [newPanel],
              panel_type: "dashboard",
              dashboard: "storybook.dashboard.storybook_dashboard_wrapper",
            },

            dashboardTags: {
              keys: [],
            },

            breakpointContext: {
              currentBreakpoint: "xl",
              maxBreakpoint: () => true,
              minBreakpoint: () => true,
              width: 0,
            },

            components: buildComponentsMap(),
            state: "complete",
            progress: 100,
            snapshot: null,
            snapshotFileName: null,
          }}
        >
          <DashboardPanelDetailProvider>
            <DashboardInputsProvider>
              <DashboardSearchPathProvider>
                <DashboardDatetimeRangeProvider>
                  <Dashboard showPanelControls={false} />
                </DashboardDatetimeRangeProvider>
              </DashboardSearchPathProvider>
            </DashboardInputsProvider>
          </DashboardPanelDetailProvider>
        </DashboardContext.Provider>
      </DashboardSearchProvider>
    </DashboardThemeProvider>
  );
};
