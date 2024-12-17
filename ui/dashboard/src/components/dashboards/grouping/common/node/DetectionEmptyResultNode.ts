import {
  CheckNodeStatus,
  DetectionNode,
  DetectionResult,
  DetectionSeveritySummary,
  DetectionSummary,
  GroupingNodeType,
} from "@powerpipe/components/dashboards/grouping/common";

class DetectionEmptyResultNode implements DetectionNode {
  private readonly _result: DetectionResult;

  constructor(result: DetectionResult) {
    this._result = result;
  }

  get sort(): string {
    return this.title;
  }

  get name(): string {
    return this._result.detection.name;
  }

  get title(): string {
    return "No results";
  }

  get result(): DetectionResult {
    return this._result;
  }

  get type(): GroupingNodeType {
    return "empty_result";
  }

  get summary(): DetectionSummary {
    return {
      total: 0,
      error: 0,
    };
  }

  get severity_summary(): DetectionSeveritySummary {
    // Bubble up the node's severity - always zero though as we have no results
    const summary: DetectionSeveritySummary = {};
    if (this._result.detection.severity) {
      summary[this._result.detection.severity] = 0;
    } else {
      summary["none"] = 0;
    }
    return summary;
  }

  get status(): CheckNodeStatus {
    // If a control has no results, this node is complete
    return "complete";
  }
}

export default DetectionEmptyResultNode;
