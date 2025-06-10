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
import { KeyValuePairs } from "@powerpipe/components/dashboards/common/types";
import { validateFilter } from "@powerpipe/components/dashboards/grouping/FilterEditor";

export type GroupingNodeType =
  | "benchmark"
  | "control"
  | "control_tag"
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
  documentation?: string;
  merge?: (other: CheckNode) => void;
};

export type DetectionNode = {
  sort: string;
  name: string;
  title: string;
  documentation?: string;
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

export type DetectionSeveritySummary = {
  none?: number;
  low?: number;
  medium?: number;
  high?: number;
  critical?: number;
};

export type CheckSummary = {
  alarm: number;
  ok: number;
  info: number;
  skip: number;
  skipped: number;
  error: number;
  invalid: number;
  muted: number;
  tbd: number;
};

export type DetectionSummary = {
  total: number;
  error: number;
};

export type CheckDynamicValueMap = {
  [dimension: string]: boolean;
};

export type CheckDynamicColsMap = {
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
  invalid = "invalid",
  error = "error",
  ok = "ok",
  info = "info",
  muted = "muted",
  skip = "skip",
  skipped = "skipped",
  empty = "empty",
  tbd = "tbd",
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
  benchmark_trunk: DetectionBenchmark[];
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

type DetectionRunProperties = {
  severity?: DetectionSeverity | undefined;
};

export type DetectionRun = {
  data: LeafNodeData;
  description?: string;
  documentation?: string;
  error?: string;
  name: string;
  panel_type: "detection";
  properties?: DetectionRunProperties;
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
  | "benchmark"
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

export type DisplayGroup = CheckDisplayGroup | DetectionDisplayGroup;

type BaseOperator = "and" | "equal" | "not_equal" | "in" | "not_in";
export type FilterOperator = BaseOperator;

export type Filter = {
  operator: FilterOperator;
  type?: FilterType;
  key?: string;
  value?: string | string[];
  title?: string;
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

export const summaryCardFilterPath = ({
  allFilters,
  expressions,
  panelName,
  pathname,
  search,
  dimension,
  metrics,
}: {
  allFilters: KeyValuePairs<Filter>;
  expressions: Filter[] | undefined;
  panelName: string;
  pathname: string;
  search: string;
  dimension: string;
  metrics: string[];
}) => {
  const expressionHasFilter = !!expressions?.find(
    (expr) => expr.type === dimension,
  );
  let newFilter: Filter;
  if (expressionHasFilter) {
    newFilter = {
      operator: "and",
      expressions: expressions?.filter((expr) => expr.type !== dimension),
    } as Filter;
    if (validateFilter(newFilter)) {
      const newParams = new URLSearchParams(search);
      const newFilters = { ...allFilters, [panelName]: newFilter };
      const asJson = JSON.stringify(newFilters);
      newParams.set("where", asJson);
      return `${pathname}?${newParams.toString()}`;
    } else {
      const newParams = new URLSearchParams(search);
      const newFilters = { ...allFilters };
      delete newFilters[panelName];
      if (!Object.keys(allFilters).length) {
        newParams.delete("where");
      } else {
        const asJson = JSON.stringify(newFilters);
        newParams.set("where", asJson);
      }
      return `${pathname}${newParams.toString() ? `?${newParams.toString()}` : ""}`;
    }
  } else {
    newFilter = {
      operator: "and",
      expressions: expressions
        ?.filter((expr) => !!expr.type)
        .concat({
          type: dimension,
          value: metrics,
          operator: "in",
        }),
    } as Filter;
    if (validateFilter(newFilter)) {
      const newParams = new URLSearchParams(search);
      const newFilters = { ...allFilters, [panelName]: newFilter };
      const asJson = JSON.stringify(newFilters);
      newParams.set("where", asJson);
      return `${pathname}?${newParams.toString()}`;
    } else {
      const newParams = new URLSearchParams(search);
      return `${pathname}${newParams.toString() ? `?${newParams.toString()}` : ""}`;
    }
  }
};

const isEqual = (a: any, b: any, isJson: boolean) => {
  return isJson ? JSON.stringify(a) === JSON.stringify(b) : a === b;
};

export const applyFilter = (
  filter: Filter,
  value: any,
  isJson: boolean = false,
) => {
  // Perform operation based on the filter operator
  switch (filter.operator) {
    case "equal":
      return isEqual(value, filter.value, isJson);
    case "not_equal":
      return !isEqual(value, filter.value, isJson);
    case "in":
      return (
        Array.isArray(filter.value) &&
        filter.value.some((v) => isEqual(value, v, isJson))
      );
    case "not_in":
      return (
        Array.isArray(filter.value) &&
        !filter.value.some((v) => isEqual(value, v, isJson))
      );
    default:
      // If an unknown operator is provided, return false
      return false;
  }
};
