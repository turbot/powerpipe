import DetectionBenchmark from "@powerpipe/components/dashboards/grouping/common/DetectionBenchmark";
import {
  AddDetectionResultsAction,
  CheckNodeStatus,
  CheckResultStatus,
  DetectionNode,
  DetectionResult,
  DetectionSeverity,
  DetectionSeveritySummary,
  DetectionTags,
  DetectionSummary,
  GroupingNodeType,
  DetectionResultDimension,
} from "@powerpipe/components/dashboards/grouping/common";
import { DashboardRunState } from "@powerpipe/types";
import {
  LeafNodeData,
  LeafNodeDataColumn,
  LeafNodeDataRow,
} from "@powerpipe/components/dashboards/common";

class Detection implements DetectionNode {
  private readonly _sortIndex: string;
  private readonly _group_id: string;
  private readonly _group_title: string | undefined;
  private readonly _group_description: string | undefined;
  private readonly _name: string;
  private readonly _title: string | undefined;
  private readonly _description: string | undefined;
  private readonly _documentation: string | undefined;
  private readonly _severity: DetectionSeverity | undefined;
  private readonly _results: DetectionResult[];
  private readonly _summary: DetectionSummary;
  private readonly _tags: DetectionTags;
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
    documentation: string | undefined,
    severity: DetectionSeverity | undefined,
    data: LeafNodeData | undefined,
    summary: DetectionSummary | undefined,
    tags: DetectionTags | undefined,
    status: DashboardRunState,
    error: string | undefined,
    benchmark_trunk: DetectionBenchmark[],
    add_detection_results: AddDetectionResultsAction,
  ) {
    this._sortIndex = sortIndex;
    this._group_id = group_id;
    this._group_title = group_title;
    this._group_description = group_description;
    this._name = name;
    this._title = title;
    this._description = description;
    this._documentation = documentation;
    this._severity = severity;
    this._results = this._build_check_results(data);
    this._summary = summary || {
      total: this.results?.length || 0,
    };
    this._tags = tags || {};
    this._status = status;
    this._error = error;

    if (
      this._status === "initialized" ||
      this._status === "blocked" ||
      this._status === "running"
    ) {
      add_detection_results([
        this._build_detection_loading_node(benchmark_trunk),
      ]);
    } else if (this._error) {
      add_detection_results([
        this._build_detection_error_node(benchmark_trunk, this._error),
      ]);
    } else if (!this._results || this._results.length === 0) {
      add_detection_results([
        this._build_detection_empty_result(benchmark_trunk),
      ]);
    } else {
      add_detection_results(
        this._build_detection_results(benchmark_trunk, this._results),
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

  get documentation(): string | undefined {
    return this._documentation;
  }

  get severity(): DetectionSeverity | undefined {
    return this._severity;
  }

  get severity_summary(): DetectionSeveritySummary {
    return {};
  }

  get type(): GroupingNodeType {
    return "detection";
  }

  get summary(): DetectionSummary {
    return { total: this._results?.length || 0 }; // this._summary;
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

  get results(): DetectionResult[] {
    return this._results;
  }

  get tags(): DetectionTags {
    return this._tags;
  }

  get_data_columns(): LeafNodeDataColumn[] {
    if (this._results.length === 0) {
      return [];
    }
    return this._results[0].columns || [];
  }

  get_data_rows(): LeafNodeDataRow[] {
    if (this._results.length === 0 || !this._results[0].rows?.length) {
      return [];
    }
    return this._results[0].rows;
  }

  private _build_detection_loading_node = (
    benchmark_trunk: DetectionBenchmark[],
  ): DetectionResult => {
    return {
      type: "loading",
      dimensions: [],
      tags: this.tags,
      detection: this,
      reason: "",
      resource: "",
      status: CheckResultStatus.ok,
      benchmark_trunk,
    };
  };

  private _build_detection_error_node = (
    benchmark_trunk: DetectionBenchmark[],
    error: string,
  ): DetectionResult => {
    return {
      type: "error",
      error,
      dimensions: [],
      tags: this.tags,
      detection: this,
      reason: "",
      resource: "",
      status: CheckResultStatus.error,
      benchmark_trunk,
    };
  };

  private _build_detection_empty_result = (
    benchmark_trunk: DetectionBenchmark[],
  ): DetectionResult => {
    return {
      type: "empty",
      error: undefined,
      dimensions: [],
      tags: this.tags,
      detection: this,
      reason: "",
      resource: "",
      status: CheckResultStatus.empty,
      benchmark_trunk,
    };
  };

  private _build_detection_results = (
    benchmark_trunk: DetectionBenchmark[],
    results: DetectionResult[],
  ): DetectionResult[] => {
    return results.map((r) => ({
      ...r,
      type: "result",
      severity: this.severity,
      tags: this.tags,
      benchmark_trunk,
      detection: this,
    }));
  };

  private _build_check_results = (data?: LeafNodeData): DetectionResult[] => {
    if (!data || !data.columns || !data.rows) {
      return [];
    }
    const results: DetectionResult[] = [];
    const dimensionColumns: LeafNodeDataColumn[] = [];
    for (const col of data.columns) {
      if (col.name === "timestamp") {
        continue;
      }
      dimensionColumns.push(col);
    }

    // const recordDimensions = {};
    const dimensions: DetectionResultDimension[] = [];
    for (const row of data.rows) {
      for (const column of dimensionColumns) {
        // recordDimensions[column.name] = recordDimensions[column.name] || {};
        const columnValue = row[column.name];
        // if (!(columnValue in recordDimensions[column.name])) {
        dimensions.push({
          key: column.name,
          value: columnValue,
        });
        //   recordDimensions[column.name][columnValue] = true;
        // }
      }
    }

    const result = {
      rows: data.rows,
      columns: data.columns,
      detection: this,
      dimension_columns: dimensionColumns,
      dimensions,
    };
    // @ts-ignore
    results.push(result);
    return results;
  };
}

export default Detection;
