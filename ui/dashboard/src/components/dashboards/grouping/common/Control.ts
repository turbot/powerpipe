import Benchmark from "./Benchmark";
import {
  AddControlResultsAction,
  CheckDynamicColsMap,
  CheckNode,
  CheckNodeStatus,
  GroupingNodeType,
  CheckResult,
  CheckResultStatus,
  CheckSeverity,
  CheckSeveritySummary,
  CheckSummary,
  CheckTags,
  findDimension,
} from "@powerpipe/components/dashboards/grouping/common";
import { DashboardRunState } from "@powerpipe/types";
import {
  LeafNodeData,
  LeafNodeDataColumn,
  LeafNodeDataRow,
} from "@powerpipe/components/dashboards/common";

class Control implements CheckNode {
  private readonly _sortIndex: string;
  private readonly _group_id: string;
  private readonly _group_title: string | undefined;
  private readonly _group_description: string | undefined;
  private readonly _name: string;
  private readonly _title: string | undefined;
  private readonly _description: string | undefined;
  private readonly _severity: CheckSeverity | undefined;
  private readonly _results: CheckResult[];
  private readonly _summary: CheckSummary;
  private readonly _tags: CheckTags;
  private readonly _status: DashboardRunState;
  private readonly _error: string | undefined;

  constructor(
    sortIndex: string,
    group_id: string,
    group_title: string | undefined,
    group_description: string | undefined,
    name: string,
    title: string | undefined,
    description: string | undefined,
    severity: CheckSeverity | undefined,
    data: LeafNodeData | undefined,
    summary: CheckSummary | undefined,
    tags: CheckTags | undefined,
    status: DashboardRunState,
    error: string | undefined,
    benchmark_trunk: Benchmark[],
    add_control_results: AddControlResultsAction,
  ) {
    this._sortIndex = sortIndex;
    this._group_id = group_id;
    this._group_title = group_title;
    this._group_description = group_description;
    this._name = name;
    this._title = title;
    this._description = description;
    this._severity = severity;
    this._results = this._build_check_results(data);
    this._summary = summary || {
      alarm: 0,
      alarm_diff: 0,
      ok: 0,
      ok_diff: 0,
      info: 0,
      info_diff: 0,
      skip: 0,
      skip_diff: 0,
      error: 0,
      error_diff: 0,
      __diff: "updated",
    };
    this._tags = tags || {};
    this._status = status;
    this._error = error;

    if (
      this._status === "initialized" ||
      this._status === "blocked" ||
      this._status === "running"
    ) {
      add_control_results([this._build_control_loading_node(benchmark_trunk)]);
    } else if (this._error) {
      add_control_results([
        this._build_control_error_node(benchmark_trunk, this._error),
      ]);
    } else if (!this._results || this._results.length === 0) {
      add_control_results([this._build_control_empty_result(benchmark_trunk)]);
    } else {
      add_control_results(
        this._build_control_results(benchmark_trunk, this._results),
      );
    }
  }

  get sort(): string {
    return `${this._sortIndex}-${this.title}`;
  }

  get name(): string {
    return this._name;
  }

  get title(): string {
    return this._title || this._name;
  }

  get severity(): CheckSeverity | undefined {
    return this._severity;
  }

  get severity_summary(): CheckSeveritySummary {
    return {};
  }

  get type(): GroupingNodeType {
    return "control";
  }

  get summary(): CheckSummary {
    let baseSummary = this._summary;
    return {
      alarm: baseSummary.alarm || 0,
      alarm_diff:
        baseSummary.alarm_diff !== undefined
          ? baseSummary.alarm_diff
          : baseSummary.alarm || 0,
      ok: baseSummary.ok || 0,
      ok_diff:
        baseSummary.ok_diff !== undefined
          ? baseSummary.ok_diff
          : baseSummary.ok || 0,
      info: baseSummary.info || 0,
      info_diff:
        baseSummary.info_diff !== undefined
          ? baseSummary.info_diff
          : baseSummary.info || 0,
      skip: baseSummary.skip || 0,
      skip_diff:
        baseSummary.skip_diff !== undefined
          ? baseSummary.skip_diff
          : baseSummary.skip || 0,
      error: baseSummary.error || 0,
      error_diff:
        baseSummary.error_diff !== undefined
          ? baseSummary.error_diff
          : baseSummary.error || 0,
      __diff: baseSummary.__diff || "updated", // Preserve __diff if it exists or default to "updated"
    };
  }

  get error(): string | undefined {
    return this._error;
  }

  get status(): CheckNodeStatus {
    switch (this._status) {
      case "initialized":
      case "blocked":
      case "running":
        return "running";
      default:
        return "complete";
    }
  }

  get results(): CheckResult[] {
    return this._results;
  }

  get tags(): CheckTags {
    return this._tags;
  }

  get_dynamic_cols(): CheckDynamicColsMap {
    const dimensionKeysMap = {
      dimensions: {},
      tags: {},
    };

    Object.keys(this._tags).forEach((t) => (dimensionKeysMap.tags[t] = true));

    if (this._results.length === 0) {
      return dimensionKeysMap;
    }
    for (const result of this._results) {
      for (const dimension of result.dimensions || []) {
        dimensionKeysMap.dimensions[dimension.key] = true;
      }
    }
    return dimensionKeysMap;
  }

  get_data_rows(tags: string[], dimensions: string[]): LeafNodeDataRow[] {
    let rows: LeafNodeDataRow[] = [];
    this._results.forEach((result) => {
      const row: LeafNodeDataRow = {
        group_id: this._group_id,
        title: this._group_title ? this._group_title : null,
        description: this._group_description ? this._group_description : null,
        control_id: this._name,
        control_title: this._title ? this._title : null,
        control_description: this._description ? this._description : null,
        severity: this._severity ? this._severity : null,
        reason: result.reason,
        resource: result.resource,
        status: result.status,
      };

      tags.forEach((tag) => {
        const val = this._tags[tag];
        row[tag] = val === undefined ? null : val;
      });

      dimensions.forEach((dimension) => {
        const val = findDimension(result.dimensions, dimension);
        row[dimension] = val === undefined ? null : val.value;
      });

      rows.push(row);
    });
    return rows;
  }

  private _build_control_loading_node = (
    benchmark_trunk: Benchmark[],
  ): CheckResult => {
    return {
      type: "loading",
      dimensions: [],
      tags: this.tags,
      control: this,
      reason: "",
      resource: "",
      status: CheckResultStatus.ok,
      benchmark_trunk,
    };
  };

  private _build_control_error_node = (
    benchmark_trunk: Benchmark[],
    error: string,
  ): CheckResult => {
    return {
      type: "error",
      error,
      dimensions: [],
      tags: this.tags,
      control: this,
      reason: "",
      resource: "",
      status: CheckResultStatus.error,
      benchmark_trunk,
    };
  };

  private _build_control_empty_result = (
    benchmark_trunk: Benchmark[],
  ): CheckResult => {
    return {
      type: "empty",
      error: undefined,
      dimensions: [],
      tags: this.tags,
      control: this,
      reason: "",
      resource: "",
      status: CheckResultStatus.empty,
      benchmark_trunk,
    };
  };

  private _build_control_results = (
    benchmark_trunk: Benchmark[],
    results: CheckResult[],
  ): CheckResult[] => {
    return results.map((r) => ({
      ...r,
      type: "result",
      severity: this.severity,
      tags: this.tags,
      benchmark_trunk,
      control: this,
    }));
  };

  private _build_check_results = (data?: LeafNodeData): CheckResult[] => {
    if (!data || !data.columns || !data.rows) {
      return [];
    }
    const results: CheckResult[] = [];
    const dimensionColumns: LeafNodeDataColumn[] = [];
    for (const col of data.columns) {
      if (
        col.name === "reason" ||
        col.name === "resource" ||
        col.name === "status"
      ) {
        continue;
      }
      dimensionColumns.push(col);
    }
    for (const row of data.rows) {
      const result = {
        reason: row.reason,
        resource: row.resource,
        status: row.status,
        status_diff: row.status_diff,
        dimensions: dimensionColumns.map((col) => ({
          key: col.name,
          value: row[col.name],
        })),
      };
      // @ts-ignore
      results.push(result);
    }
    return results;
  };
}

export default Control;
