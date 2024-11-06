import {
  GroupingNodeType,
  DetectionNode,
  DetectionNodeStatus,
  DetectionSummary,
  DetectionResult,
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
    };
  }

  get status(): DetectionNodeStatus {
    // This will bubble up through the hierarchy and put all ancestral nodes in a running state
    return "running";
  }
}

export default DetectionRunningNode;
