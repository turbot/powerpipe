import {
  DetectionNodeStatus,
  GroupingNodeType,
  DetectionSummary,
  DetectionNode,
  DetectionResult,
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
