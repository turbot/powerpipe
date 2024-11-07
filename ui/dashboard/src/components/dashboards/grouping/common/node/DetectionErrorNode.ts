import {
  DetectionNode,
  DetectionNodeStatus,
  DetectionResult,
  DetectionSeveritySummary,
  DetectionSummary,
  GroupingNodeType,
} from "../index";

class DetectionErrorNode implements DetectionNode {
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
    return this._result.detection.title || this.name;
  }

  get error(): string {
    return this._result.detection.error || "Unknown error";
  }

  get type(): GroupingNodeType {
    return "error";
  }

  get severity_summary(): DetectionSeveritySummary {
    // Bubble up the node's severity - always zero though as we have no results
    const summary = {};
    if (this._result.detection.severity) {
      summary[this._result.detection.severity] = 0;
    }
    return summary;
  }

  get summary(): DetectionSummary {
    return {
      total: 1,
    };
  }

  get status(): DetectionNodeStatus {
    // If a control has gone to error, this node is complete
    return "complete";
  }
}

export default DetectionErrorNode;
