import {
  CheckDisplayGroup,
  CheckFilter,
} from "@powerpipe/components/dashboards/check/common";
import { LeafNodeData, Width } from "@powerpipe/components/dashboards/common";
import { Ref } from "react";
import { Theme } from "@powerpipe/hooks/useTheme";

export type IDashboardContext = {
  cliMode: DashboardCliMode;

  versionMismatchCheck: boolean;
  metadata: ServerMetadata | null;
  availableDashboardsLoaded: boolean;

  closePanelDetail(): void;
  dispatch(action: DashboardAction): void;

  dataMode: DashboardDataMode;
  snapshotId: string | null;

  refetchDashboard: boolean;

  error: any;

  panelsLog: PanelsLog;
  panelsMap: PanelsMap;

  execution_id: string | null;

  dashboards: AvailableDashboard[];
  dashboardsMap: AvailableDashboardsDictionary;
  dashboardsMetadata: DashboardsMetadataDictionary;
  dashboard: DashboardDefinition | null;

  selectedPanel: PanelDefinition | null;
  selectedDashboard: AvailableDashboard | null;
  selectedDashboardInputs: DashboardInputs;
  lastChangedInput: string | null;

  searchPathPrefix: string[];

  dashboardTags: DashboardTags;

  search: DashboardSearch;

  breakpointContext: IBreakpointContext;
  themeContext: IThemeContext;

  components: ComponentsMap;

  progress: number;
  state: DashboardRunState;
  render: {
    headless: boolean;
    snapshotCompleteDiv: boolean;
  };

  snapshot: DashboardSnapshot | null;
  snapshotFileName: string | null;
  snapshot_metadata_loaded: boolean;

  diff?: {
    panelsMap: PanelsMap;
    snapshotFileName: string;
  };

  showCustomizeBenchmarkPanel: boolean;
};

export type IBreakpointContext = {
  currentBreakpoint: string | null;
  maxBreakpoint(breakpointAndDown: string): boolean;
  minBreakpoint(breakpointAndUp: string): boolean;
  width: number;
};

export type IThemeContext = {
  theme: Theme;
  setTheme(theme: string): void;
  wrapperRef: Ref<null>;
};

export const DashboardDataModeLive = "live";
export const DashboardDataModeCLISnapshot = "cli_snapshot";
export const DashboardDataModeCloudSnapshot = "cloud_snapshot";

export type DashboardCliMode = "powerpipe" | "steampipe";

export type DashboardDataMode = "live" | "cli_snapshot" | "cloud_snapshot";

export type PanelDataMode = "diff";

export type SocketURLFactory = () => Promise<string>;

export type IActions = {
  [type: string]: string;
};

export type ReceivedSocketMessagePayload = {
  action: string;
  [key: string]: any;
};

export type ComponentsMap = {
  [name: string]: any;
};

export type PanelLog = {
  error?: string | null;
  executionTime?: number;
  isDependency?: boolean;
  prefix?: string;
  status: DashboardRunState;
  timestamp: string;
  title: string;
};

export type PanelsLog = {
  [name: string]: PanelLog[];
};

export type PanelsMap = {
  [name: string]: PanelDefinition;
};

export type DashboardRunState =
  | "initialized"
  | "blocked"
  | "running"
  | "cancelled"
  | "error"
  | "complete";

export const DashboardActions: IActions = {
  AVAILABLE_DASHBOARDS: "available_dashboards",
  CLEAR_DASHBOARD_INPUTS: "clear_dashboard_inputs",
  CONTROL_COMPLETE: "control_complete",
  CONTROL_ERROR: "control_error",
  CONTROLS_UPDATED: "controls_updated",
  SERVER_METADATA: "server_metadata",
  DASHBOARD_METADATA: "dashboard_metadata",
  DELETE_DASHBOARD_INPUT: "delete_dashboard_input",
  DIFF_SNAPSHOT: "diff_snapshot",
  EXECUTION_COMPLETE: "execution_complete",
  EXECUTION_ERROR: "execution_error",
  EXECUTION_STARTED: "execution_started",
  INPUT_VALUES_CLEARED: "input_values_cleared",
  LEAF_NODE_COMPLETE: "leaf_node_complete",
  LEAF_NODE_UPDATED: "leaf_node_updated",
  LEAF_NODES_COMPLETE: "leaf_nodes_complete",
  LEAF_NODES_UPDATED: "leaf_nodes_updated",
  SELECT_DASHBOARD: "select_dashboard",
  SELECT_PANEL: "select_panel",
  SET_CLI_MODE: "set_cli_mode",
  SET_DASHBOARD: "set_dashboard",
  SET_DASHBOARD_INPUT: "set_dashboard_input",
  SET_DASHBOARD_INPUTS: "set_dashboard_inputs",
  SET_SEARCH_PATH_PREFIX: "set_search_path_prefix",
  SET_SNAPSHOT_METADATA_LOADED: "set_snapshot_metadata_loaded",
  SET_DASHBOARD_SEARCH_VALUE: "set_dashboard_search_value",
  SET_DASHBOARD_SEARCH_GROUP_BY: "set_dashboard_search_group_by",
  SET_DASHBOARD_TAG_KEYS: "set_dashboard_tag_keys",
  SET_DATA_MODE: "set_data_mode",
  SET_REFETCH_DASHBOARD: "set_refetch_dashboard",
  SHOW_CUSTOMIZE_BENCHMARK_PANEL: "show_customize_benchmark_panel",
  HIDE_CUSTOMIZE_BENCHMARK_PANEL: "hide_customize_benchmark_panel",
  WORKSPACE_ERROR: "workspace_error",
};

type DashboardExecutionEventSchemaVersion =
  | "20220614"
  | "20220929"
  | "20221222"
  | "20240130"
  | "20240607";

type DashboardExecutionStartedEventSchemaVersion =
  | "20220614"
  | "20221222"
  | "20240130"
  | "20240607";

type DashboardExecutionCompleteEventSchemaVersion =
  | "20220614"
  | "20220929"
  | "20221222"
  | "20240130"
  | "20240607";

type DashboardSnapshotSchemaVersion =
  | "20220614"
  | "20220929"
  | "20221222"
  | "20240130"
  | "20240607";

export type DashboardExecutionStartedEvent = {
  action: "execution_started";
  execution_id: string;
  inputs: DashboardInputs;
  layout: DashboardLayoutNode;
  panels: PanelsMap;
  variables: DashboardVariables;
  schema_version: DashboardExecutionStartedEventSchemaVersion;
  start_time: string;
};

export type DashboardExecutionEventWithSchema = {
  schema_version: DashboardExecutionEventSchemaVersion;
  [key: string]: any;
};

export type DashboardExecutionCompleteEvent = {
  action: string;
  schema_version: DashboardExecutionCompleteEventSchemaVersion;
  execution_id: string;
  snapshot: DashboardSnapshot;
};

// https://github.com/microsoft/TypeScript/issues/28046
export type ElementType<T extends ReadonlyArray<unknown>> =
  T extends ReadonlyArray<infer ElementType> ? ElementType : never;

const dashboardActions = Object.values(DashboardActions);

export type DashboardActionType = ElementType<typeof dashboardActions>;

export type DashboardAction = {
  type: DashboardActionType;
  [key: string]: any;
};

type DashboardSearchGroupByMode = "mod" | "tag";

type DashboardSearchGroupBy = {
  value: DashboardSearchGroupByMode;
  tag: string | null;
};

export type DashboardSearch = {
  value: string;
  groupBy: DashboardSearchGroupBy;
};

export type DashboardTags = {
  keys: string[];
};

export type SelectedDashboardStates = {
  dashboard_name: string | undefined;
  dataMode: DashboardDataMode;
  refetchDashboard: boolean;
  search: DashboardSearch;
  searchParams: URLSearchParams;
  searchPathPrefix: string[];
  selectedDashboard: AvailableDashboard | null;
  selectedDashboardInputs: DashboardInputs;
};

export type DashboardInputs = {
  [name: string]: string;
};

type DashboardVariables = {
  [name: string]: any;
};

export type ModServerMetadata = {
  title: string;
  full_name: string;
  short_name: string;
};

type InstalledModsServerMetadata = {
  [key: string]: ModServerMetadata;
};

type CliServerMetadata = {
  version: string;
};

export type CloudServerActorMetadata = {
  id: string;
  handle: string;
};

export type CloudServerIdentityMetadata = {
  id: string;
  handle: string;
  type: "org" | "user";
};

export type CloudServerWorkspaceMetadata = {
  id: string;
  handle: string;
};

type CloudServerMetadata = {
  actor: CloudServerActorMetadata;
  identity: CloudServerIdentityMetadata;
  workspace: CloudServerWorkspaceMetadata;
};

export type ServerMetadata = {
  mod: ModServerMetadata;
  installed_mods?: InstalledModsServerMetadata;
  cli?: CliServerMetadata;
  cloud?: CloudServerMetadata;
  telemetry: "info" | "none";
  search_path: SearchPathMetadata;
};

export type DashboardLayoutNode = {
  name: string;
  panel_type: DashboardPanelType;
  children?: DashboardLayoutNode[];
};

export type DashboardPanelType =
  | "benchmark"
  | "benchmark_tree"
  | "card"
  | "chart"
  | "container"
  | "control"
  | "dashboard"
  | "edge"
  | "error"
  | "flow"
  | "graph"
  | "hierarchy"
  | "image"
  | "input"
  | "node"
  | "table"
  | "text"
  | "with";

export type DashboardSnapshotViewFilterByMetadata = CheckFilter;
export type DashboardSnapshotViewGroupByMetadata = CheckDisplayGroup[];

export type DashboardSnapshotViewMetadata = {
  filter_by?: DashboardSnapshotViewFilterByMetadata;
  group_by?: DashboardSnapshotViewGroupByMetadata;
};

export type DashboardSnapshotMetadata = {
  view?: DashboardSnapshotViewMetadata;
};

export type DashboardSnapshot = {
  schema_version: DashboardSnapshotSchemaVersion;
  layout: DashboardLayoutNode;
  panels: PanelsMap;
  inputs: DashboardInputs;
  variables: DashboardVariables;
  search_path: string[];
  start_time: string;
  end_time: string;
  metadata?: DashboardSnapshotMetadata;
};

type AvailableDashboardTags = {
  [key: string]: string;
};

type AvailableDashboardType = "benchmark" | "dashboard" | "snapshot";

export type AvailableDashboard = {
  full_name: string;
  short_name: string;
  mod_full_name?: string;
  tags: AvailableDashboardTags;
  title: string;
  is_top_level: boolean;
  type: AvailableDashboardType;
  children?: AvailableDashboard[];
  trunks?: string[][];
};

export type AvailableDashboardsDictionary = {
  [key: string]: AvailableDashboard;
};

export type SearchPathMetadata = {
  configured_search_path?: string | null;
  original_search_path?: string[];
  resolved_search_path?: string[];
};

export type DashboardMetadata = {
  database?: string | null;
  search_path: SearchPathMetadata;
};

export type DashboardsMetadataDictionary = {
  [key: string]: DashboardMetadata;
};

export type ContainerDefinition = {
  name: string;
  panel_type?: string;
  data?: LeafNodeData;
  title?: string;
  width?: number;
  children?: (ContainerDefinition | PanelDefinition)[];
};

export type PanelProperties = {
  [key: string]: any;
};

export type DependencyPanelProperties = {
  name: string;
};

export type PanelDefinition = {
  name: string;
  args?: any[];
  display?: string;
  display_type?: string;
  panel_type: DashboardPanelType;
  title?: string;
  description?: string;
  documentation?: string;
  width?: Width;
  sql?: string;
  data?: LeafNodeData;
  source_definition?: string;
  status?: DashboardRunState;
  error?: string;
  properties?: PanelProperties;
  dashboard: string;
  children?: DashboardLayoutNode[];
  dependencies?: string[];
};

export type PanelDependenciesByStatus = {
  [status: string]: PanelDefinition[];
};

export type DashboardDefinition = {
  artificial: boolean;
  name: string;
  panel_type: string;
  title?: string;
  width?: number;
  children?: (ContainerDefinition | PanelDefinition)[];
  dashboard: string;
};

export type DashboardsCollection = {
  dashboards: AvailableDashboard[];
  dashboardsMap: AvailableDashboardsDictionary;
};

export type DashboardDataOptions = {
  dataMode: DashboardDataMode;
  snapshotId?: string;
};

export type DashboardRenderOptions = {
  headless?: boolean;
};
