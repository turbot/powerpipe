import Benchmark from "./Benchmark";
import DetectionBenchmark from "@powerpipe/components/dashboards/grouping/common/DetectionBenchmark";
import {
  BasePrimitiveProps,
  ExecutablePrimitiveProps,
  LeafNodeData,
  LeafNodeDataColumn,
  LeafNodeDataRow,
} from "../../common";
import { DashboardRunState } from "@powerpipe/types";

export type GroupingNodeType =
  | "benchmark"
  | "control"
  | "control_tag"
  | "detection_benchmark"
  | "detection"
  | "detection_tag"
  | "dimension"
  | "empty_result"
  | "error"
  | "reason"
  | "resource"
  | "result"
  | "running"
  | "root"
  | "severity"
  | "status";

export type CheckNode = {
  sort: string;
  name: string;
  title: string;
  type: GroupingNodeType;
  severity?: CheckSeverity;
  severity_summary: CheckSeveritySummary;
  status: CheckNodeStatus;
  summary: CheckSummary;
  children?: CheckNode[];
  data?: LeafNodeData;
  error?: string;
  merge?: (other: CheckNode) => void;
};

export type DetectionNode = {
  sort: string;
  name: string;
  title: string;
  type: GroupingNodeType;
  severity?: DetectionSeverity;
  severity_summary: DetectionSeveritySummary;
  status: CheckNodeStatus;
  summary: DetectionSummary;
  children?: DetectionNode[];
  data?: LeafNodeData;
  error?: string;
  merge?: (other: DetectionNode) => void;
};

export type CheckNodeStatus = "running" | "complete";

export type DetectionNodeStatus = "running" | "complete";

export type CheckSeverity = "none" | "low" | "medium" | "high" | "critical";

export type CheckSeveritySummary =
  | {}
  | {
      [key in CheckSeverity]: number;
    };

export type DetectionSeverity = "none" | "low" | "medium" | "high" | "critical";

export type DetectionSeveritySummary =
  | {}
  | {
      [key in DetectionSeverity]: number;
    };

export type CheckSummary = {
  alarm: number;
  ok: number;
  info: number;
  skip: number;
  error: number;
};

export type DetectionSummary = {
  total: number;
};

export type CheckDynamicValueMap = {
  [dimension: string]: boolean;
};

export type CheckDynamicColsMap = {
  dimensions: CheckDynamicValueMap;
  tags: CheckDynamicValueMap;
};

export type DetectionDynamicColsMap = {
  dimensions: CheckDynamicValueMap;
  tags: CheckDynamicValueMap;
};

export type CheckTags = {
  [key: string]: string;
};

export type DetectionTags = {
  [key: string]: string;
};

export type CheckResultDimension = {
  key: string;
  value: string;
};

export type DetectionResultDimension = {
  key: string;
  value: string;
};

export enum CheckResultStatus {
  alarm = "alarm",
  ok = "ok",
  info = "info",
  skip = "skip",
  error = "error",
  empty = "empty",
}

export type CheckResultType = "loading" | "error" | "empty" | "result";

export type CheckResult = {
  dimensions: CheckResultDimension[];
  tags: CheckTags;
  control: CheckNode;
  benchmark_trunk: Benchmark[];
  status: CheckResultStatus;
  reason: string;
  resource: string;
  severity?: CheckSeverity;
  error?: string;
  type: CheckResultType;
};

export type DetectionResult = {
  rows?: LeafNodeDataRow[];
  columns?: LeafNodeDataColumn[];
  dimension_columns?: LeafNodeDataColumn[];
  dimensions: DetectionResultDimension[];
  tags: CheckTags;
  detection: DetectionNode;
  detection_benchmark_trunk: DetectionBenchmark[];
  severity?: DetectionSeverity;
  status: CheckResultStatus;
  reason: string;
  resource: string;
  error?: string;
  type: CheckResultType;
};

type CheckControlRunProperties = {
  severity?: CheckSeverity | undefined;
};

export type CheckControlRun = {
  data: LeafNodeData;
  description?: string;
  error?: string;
  name: string;
  panel_type: "control";
  properties?: CheckControlRunProperties;
  severity?: CheckSeverity | undefined;
  status: DashboardRunState;
  summary: CheckSummary;
  tags?: CheckTags;
  title?: string;
};

export type DetectionRun = {
  data: LeafNodeData;
  description?: string;
  error?: string;
  name: string;
  panel_type: "detection";
  severity?: DetectionSeverity | undefined;
  status: DashboardRunState;
  summary: DetectionSummary;
  tags?: CheckTags;
  title?: string;
};

export type CheckDisplayGroupType =
  | "benchmark"
  | "control"
  | "control_tag"
  | "result"
  | "dimension"
  | "reason"
  | "resource"
  | "severity"
  | "status"
  | string;

export type DetectionDisplayGroupType =
  | "detection_benchmark"
  | "detection"
  | "detection_tag"
  | "result"
  | "dimension"
  | string;

export type DisplayGroupType =
  | CheckDisplayGroupType
  | DetectionDisplayGroupType;

export type CheckDisplayGroup = {
  type: DisplayGroupType;
  value?: string | undefined;
};

export type DetectionDisplayGroup = {
  type: DetectionDisplayGroupType;
  value?: string | undefined;
};

type BaseOperator = "and" | "equal" | "not_equal";
export type FilterOperator = BaseOperator;

export type Filter = {
  operator: FilterOperator;
  type?: FilterType;
  key?: string;
  value?: string;
  title?: string;
  context?: string;
  expressions?: Filter[];
};

export type FilterType = CheckDisplayGroupType | DetectionDisplayGroupType;

export type BenchmarkTreeProps = BasePrimitiveProps &
  ExecutablePrimitiveProps & {
    properties: {
      grouping: CheckNode;
      first_child_summaries: CheckSummary[];
    };
  };

export type DetectionBenchmarkTreeProps = BasePrimitiveProps &
  ExecutablePrimitiveProps & {
    properties: {
      grouping: DetectionNode;
      first_child_summaries: DetectionSummary[];
    };
  };

export type AddControlResultsAction = (results: CheckResult[]) => void;

export type AddDetectionResultsAction = (results: DetectionResult[]) => void;

export const findDimension = (
  dimensions?: CheckResultDimension[],
  key?: string,
) => {
  if (!dimensions || !key) {
    return undefined;
  }
  return dimensions.find((d) => d.key === key);
};
