import {
  DetectionNode,
  DetectionNodeStatus,
  DetectionResult,
  DetectionSeveritySummary,
  DetectionSummary,
  GroupingNodeType,
} from "../index";

class DetectionRunningNode implements DetectionNode {
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

  get type(): GroupingNodeType {
    return "running";
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

  get status(): DetectionNodeStatus {
    // This will bubble up through the hierarchy and put all ancestral nodes in a running state
    return "running";
  }
}

export default DetectionRunningNode;
