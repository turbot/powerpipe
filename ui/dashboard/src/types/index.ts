import {
  CheckDisplayGroup,
  Filter,
} from "@powerpipe/components/dashboards/grouping/common";
import { DatetimeRange } from "@powerpipe/hooks/useDashboardDatetimeRange";
import {
  KeyValuePairs,
  TableConfig,
} from "@powerpipe/components/dashboards/common/types";
import { LeafNodeData, Width } from "@powerpipe/components/dashboards/common";
import { Ref } from "react";
import { Theme } from "@powerpipe/hooks/useTheme";

export type IDashboardContext = {
  versionMismatchCheck: boolean;
  metadata: ServerMetadata | null;
  availableDashboardsLoaded: boolean;

  dispatch(action: DashboardAction): void;

  dataMode: DashboardDataMode;
  snapshotId: string | null;

  error: any;

  overlayVisible: boolean;

  panelsLog: PanelsLog;
  panelsMap: PanelsMap;

  execution_id: string | null;

  dashboards: AvailableDashboard[];
  dashboardsMap: AvailableDashboardsDictionary;
  dashboardsMetadata: DashboardsMetadataDictionary;
  dashboard: DashboardDefinition | null;

  selectedDashboard: AvailableDashboard | null;

  dashboardTags: DashboardTags;

  breakpointContext: IBreakpointContext;

  components: ComponentsMap;

  rootPathname: string;

  progress: number;
  state: DashboardRunState;

  snapshot: DashboardSnapshot | null;
  snapshotFileName: string | null;
  snapshot_metadata_loaded: boolean;
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
  CONTROL_COMPLETE: "control_complete",
  CONTROL_ERROR: "control_error",
  CONTROLS_UPDATED: "controls_updated",
  SERVER_METADATA: "server_metadata",
  DASHBOARD_METADATA: "dashboard_metadata",
  EXECUTION_COMPLETE: "execution_complete",
  EXECUTION_ERROR: "execution_error",
  EXECUTION_STARTED: "execution_started",
  LEAF_NODE_COMPLETE: "leaf_node_complete",
  LEAF_NODE_UPDATED: "leaf_node_updated",
  LEAF_NODES_COMPLETE: "leaf_nodes_complete",
  LEAF_NODES_UPDATED: "leaf_nodes_updated",
  LOAD_SNAPSHOT: "load_snapshot",
  SELECT_DASHBOARD: "select_dashboard",
  SET_DASHBOARD: "set_dashboard",
  SET_OVERLAY_VISIBLE: "set_overlay_visible",
  SET_SEARCH_PATH_PREFIX: "set_search_path_prefix",
  SET_SNAPSHOT_METADATA_LOADED: "set_snapshot_metadata_loaded",
  SET_DASHBOARD_TAG_KEYS: "set_dashboard_tag_keys",
  SET_DATA_MODE: "set_data_mode",
  WORKSPACE_ERROR: "workspace_error",
};

type DashboardExecutionEventSchemaVersion =
  | "20220614"
  | "20220929"
  | "20221222"
  | "20240130"
  | "20240607"
  | "20241125";

type DashboardExecutionStartedEventSchemaVersion =
  | "20220614"
  | "20221222"
  | "20240130"
  | "20240607"
  | "20241125";

type DashboardExecutionCompleteEventSchemaVersion =
  | "20220614"
  | "20220929"
  | "20221222"
  | "20240130"
  | "20240607"
  | "20241125";

type DashboardSnapshotSchemaVersion =
  | "20220614"
  | "20220929"
  | "20221222"
  | "20240130"
  | "20240607"
  | "20241125";

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

export type DashboardSearchGroupByMode = "mod" | "tag";

export type DashboardDisplayMode = "top_level" | "all";

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
  supports_search_path: boolean;
  supports_time_range: boolean;
};

export type DashboardLayoutNode = {
  name: string;
  panel_type: DashboardPanelType;
  children?: DashboardLayoutNode[];
};

export const DashboardPanelTypeBenchmark = "benchmark";

export type DashboardPanelType =
  | "benchmark"
  | "benchmark_tree"
  | "card"
  | "chart"
  | "container"
  | "control"
  | "dashboard"
  | "detection"
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

export type DashboardSnapshotViewFilterByMetadata = Filter;
export type DashboardSnapshotViewGroupByMetadata = CheckDisplayGroup[];
export type DashboardSnapshotViewTableConfigMetadata = TableConfig;

export type DashboardSnapshotViewMetadata = {
  filter_by?: DashboardSnapshotViewFilterByMetadata;
  group_by?: DashboardSnapshotViewGroupByMetadata;
  table?: DashboardSnapshotViewTableConfigMetadata;
};

export type DashboardSnapshotMetadata = {
  view?: KeyValuePairs<DashboardSnapshotViewMetadata>;
  datetime_range?: DatetimeRange;
  search_path_prefix?: string[];
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

type AvailableDashboardType =
  | "available_dashboard"
  | "benchmark"
  | "control"
  | "dashboard"
  | "detection"
  | "snapshot";

export type AvailableDashboard = {
  full_name: string;
  short_name: string;
  mod_full_name?: string;
  benchmark_type?: "control" | "detection";
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
  supports_search_path: boolean;
  supports_time_range: boolean;
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

export type PanelDefinitionBenchmarkType = "control" | "detection";

export type PanelDefinition = {
  name: string;
  args?: any[];
  benchmark_type?: PanelDefinitionBenchmarkType;
  children?: DashboardLayoutNode[];
  dashboard: string;
  data?: LeafNodeData;
  dependencies?: string[];
  description?: string;
  display?: string;
  display_type?: string;
  documentation?: string;
  error?: string;
  panel_type: DashboardPanelType;
  properties?: PanelProperties;
  source_definition?: string;
  sql?: string;
  status?: DashboardRunState;
  title?: string;
  width?: Width;
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
