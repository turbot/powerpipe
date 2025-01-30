import Detection from "@powerpipe/components/dashboards/grouping/common/Detection";
import padStart from "lodash/padStart";
import {
  AddDetectionResultsAction,
  DetectionNode,
  DetectionNodeStatus,
  DetectionResult,
  DetectionRun,
  DetectionSeveritySummary,
  DetectionSummary,
  GroupingNodeType,
} from "@powerpipe/components/dashboards/grouping/common";
import { DashboardLayoutNode, PanelsMap } from "@powerpipe/types";
import {
  LeafNodeData,
  LeafNodeDataColumn,
  LeafNodeDataRow,
} from "@powerpipe/components/dashboards/common";
import { KeyValuePairs } from "@powerpipe/components/dashboards/common/types";

class DetectionBenchmark implements DetectionNode {
  private readonly _sortIndex: string;
  private readonly _name: string;
  private readonly _title: string;
  private readonly _description?: string;
  private readonly _documentation?: string;
  private readonly _benchmarks: DetectionBenchmark[];
  private readonly _detections: Detection[];
  private readonly _add_detection_results: AddDetectionResultsAction;
  private readonly _all_detection_results: DetectionResult[];

  constructor(
    sortIndex: string,
    name: string,
    title: string | undefined,
    description: string | undefined,
    documentation: string | undefined,
    benchmarks: DashboardLayoutNode[] | undefined,
    detections: DashboardLayoutNode[] | undefined,
    panelsMap: PanelsMap,
    trunk: DetectionBenchmark[],
    add_detection_results?: AddDetectionResultsAction,
  ) {
    this._sortIndex = sortIndex;
    this._all_detection_results = [];
    this._name = name;
    this._title = title || name;
    this._description = description;
    this._documentation = documentation;

    if (!add_detection_results) {
      this._add_detection_results = this.add_detection_results;
    } else {
      this._add_detection_results = add_detection_results;
    }

    const thisTrunk = [...trunk, this];
    const nestedBenchmarks: DetectionBenchmark[] = [];
    const benchmarksToAdd = benchmarks || [];
    const lengthMaxBenchmarkIndex = (benchmarksToAdd.length - 1).toString()
      .length;
    benchmarksToAdd.forEach((nestedBenchmark, benchmarkIndex) => {
      const nestedDefinition = panelsMap[nestedBenchmark.name];
      // @ts-ignore
      const benchmarks = nestedBenchmark.children?.filter(
        (child) => child.panel_type === "benchmark",
      );
      // @ts-ignore
      const detections = nestedBenchmark.children?.filter(
        (child) => child.panel_type === "detection",
      );
      nestedBenchmarks.push(
        new DetectionBenchmark(
          `benchmark-${padStart(
            benchmarkIndex.toString(),
            lengthMaxBenchmarkIndex,
          )}`,
          nestedDefinition.name,
          nestedDefinition.title,
          nestedDefinition.description,
          nestedDefinition.documentation,
          benchmarks,
          detections,
          panelsMap,
          thisTrunk,
          this._add_detection_results,
        ),
      );
    });
    const nestedDetections: Detection[] = [];
    const detectionsToAdd = detections || [];
    const lengthMaxDetectionIndex = (detectionsToAdd.length - 1).toString()
      .length;
    detectionsToAdd.forEach((nestedDetection, detectionIndex) => {
      // @ts-ignore
      const detection = panelsMap[nestedDetection.name] as DetectionRun;
      nestedDetections.push(
        new Detection(
          `detection-${padStart(detectionIndex.toString(), lengthMaxDetectionIndex)}`,
          this._name,
          this._title,
          this._description,
          detection.name,
          detection.title,
          detection.description,
          detection.documentation,
          detection.properties?.severity,
          detection.data,
          detection.summary,
          detection.tags,
          detection.status,
          detection.error,
          thisTrunk,
          this._add_detection_results,
        ),
      );
    });
    this._benchmarks = nestedBenchmarks;
    this._detections = nestedDetections;
  }

  private add_detection_results = (results: DetectionResult[]) => {
    this._all_detection_results.push(...results);
  };

  get all_detection_results(): DetectionResult[] {
    return this._all_detection_results;
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

  get type(): GroupingNodeType {
    return "benchmark";
  }

  get documentation(): string | undefined {
    return this._documentation;
  }

  get children(): DetectionNode[] {
    return [...this._benchmarks, ...this._detections];
  }

  get detectionBenchmarks(): DetectionBenchmark[] {
    return this._benchmarks;
  }

  get detections(): Detection[] {
    return this._detections;
  }

  get summary(): DetectionSummary {
    const summary = {
      total: 0,
      error: 0,
    };
    for (const benchmark of this._benchmarks) {
      const nestedSummary = benchmark.summary;
      summary.total += nestedSummary.total;
      summary.error += nestedSummary.error;
    }
    for (const detection of this._detections) {
      const nestedSummary = detection.summary;
      summary.total += nestedSummary.total;
      summary.error += nestedSummary.error;
    }
    return summary;
  }

  get severity_summary(): DetectionSeveritySummary {
    const summary: DetectionSeveritySummary = {};
    for (const benchmark of this._benchmarks) {
      const nestedSummary = benchmark.severity_summary;
      if ("low" in nestedSummary) {
        summary.low = summary.low || 0;
        summary.low += nestedSummary.low || 0;
      }
      if ("medium" in nestedSummary) {
        summary.medium = summary.medium || 0;
        summary.medium += nestedSummary.medium || 0;
      }
      if ("high" in nestedSummary) {
        summary.high = summary.high || 0;
        summary.high += nestedSummary.high || 0;
      }
      if ("critical" in nestedSummary) {
        summary.critical = summary.critical || 0;
        summary.critical += nestedSummary.critical || 0;
      }
    }
    for (const detection of this._detections) {
      const nestedSummary = detection.severity_summary;
      if ("low" in nestedSummary) {
        summary.low = summary.low || 0;
        summary.low += nestedSummary.low || 0;
      }
      if ("medium" in nestedSummary) {
        summary.medium = summary.medium || 0;
        summary.medium += nestedSummary.medium || 0;
      }
      if ("high" in nestedSummary) {
        summary.high = summary.high || 0;
        summary.high += nestedSummary.high || 0;
      }
      if ("critical" in nestedSummary) {
        summary.critical = summary.critical || 0;
        summary.critical += nestedSummary.critical || 0;
      }
    }
    return summary;
  }

  get status(): DetectionNodeStatus {
    for (const benchmark of this._benchmarks) {
      if (benchmark.status === "running") {
        return "running";
      }
    }
    for (const detection of this._detections) {
      if (detection.status === "running") {
        return "running";
      }
    }
    return "complete";
  }

  get_data_table(): LeafNodeData {
    const columns = this.get_data_columns();
    const rows = this.get_data_rows();

    return {
      columns,
      rows,
    };
  }

  get_data_columns(): LeafNodeDataColumn[] {
    const columnMap: KeyValuePairs<LeafNodeDataColumn> = {};
    const columns: LeafNodeDataColumn[] = [];

    this._benchmarks.forEach((benchmark) => {
      const nestedColumns = benchmark.get_data_columns();
      for (const nestedColumn of nestedColumns) {
        if (columnMap[nestedColumn.name]) {
          continue;
        }
        columnMap[nestedColumn.name] = nestedColumn;
        columns.push(nestedColumn);
      }
    });
    this._detections.forEach((detection) => {
      const nestedColumns = detection.get_data_columns();
      for (const nestedColumn of nestedColumns) {
        if (columnMap[nestedColumn.name]) {
          continue;
        }
        columnMap[nestedColumn.name] = nestedColumn;
        columns.push(nestedColumn);
      }
    });
    return columns;
  }

  get_data_rows(): LeafNodeDataRow[] {
    let rows: LeafNodeDataRow[] = [];
    this._benchmarks.forEach((benchmark) => {
      rows = [...rows, ...benchmark.get_data_rows()];
    });
    this._detections.forEach((detection) => {
      rows = [...rows, ...detection.get_data_rows()];
    });
    return rows;
  }
}

export default DetectionBenchmark;
