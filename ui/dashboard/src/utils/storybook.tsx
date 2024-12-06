import "@powerpipe/utils/registerComponents";
import Dashboard from "@powerpipe/components/dashboards/layout/Dashboard";
import { buildComponentsMap } from "@powerpipe/components";
import { DashboardContext } from "@powerpipe/hooks/useDashboard";
import {
  DashboardDataMode,
  DashboardDataModeLive,
  DashboardPanelType,
  DashboardRunState,
  DashboardSearch,
} from "@powerpipe/types";
import { noop } from "@powerpipe/utils/func";
import { useStorybookTheme } from "@powerpipe/hooks/useStorybookTheme";

type PanelStoryDecoratorProps = {
  dataMode?: DashboardDataMode;
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
  dataMode = DashboardDataModeLive,
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
        closePanelDetail: noop,
        dataMode,
        snapshotId: null,
        dispatch: noop,
        error: null,
        dashboards: [],
        dashboardsMap: {},
        selectedPanel: null,
        selectedDashboard: {
          title: "Storybook Dashboard Wrapper",
          full_name: "storybook.dashboard.storybook_dashboard_wrapper",
          short_name: "storybook_dashboard_wrapper",
          type: "dashboard",
          tags: {},
          mod_full_name: "mod.storybook",
          is_top_level: true,
        },
        selectedDashboardInputs: {},
        lastChangedInput: null,
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

        search: stubDashboardSearch,

        breakpointContext: {
          currentBreakpoint: "xl",
          maxBreakpoint: () => true,
          minBreakpoint: () => true,
          width: 0,
        },

        themeContext: {
          theme,
          setTheme: noop,
          wrapperRef,
        },

        components: buildComponentsMap(),
        refetchDashboard: false,
        state: "complete",
        progress: 100,
        render: { headless: false, snapshotCompleteDiv: false },
        snapshot: null,
        snapshotFileName: null,
      }}
    >
      <Dashboard showPanelControls={false} />
    </DashboardContext.Provider>
  );
};
