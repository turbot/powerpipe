import {
  CheckNodeStatus,
  GroupingNodeType,
  DetectionSummary,
  DetectionResult,
  DetectionNode,
} from "../index";

// class ControlResultNode extends HierarchyNode {
//   private readonly _result: CheckResult;
//
//   constructor(
//     result: CheckResult,
//     sort: string,
//     type: GroupingNodeType,
//     key: string,
//     value: string,
//     children?: CheckNode[],
//   ) {
//     super(type, `${key}=${value}`, value, sort, children || []);
//     this._result = result;
//   }

class DetectionResultNode implements DetectionNode {
  private readonly _result: DetectionResult;

  constructor(result: DetectionResult) {
    this._result = result;
  }

  get sort(): string {
    return "0";
  }

  get name(): string {
    return `${this._result.detection.name}-${this._result.resource}`;
  }

  get title(): string {
    return this._result.reason;
  }

  get result(): DetectionResult {
    return this._result;
  }

  get type(): GroupingNodeType {
    return "result";
  }

  get summary(): DetectionSummary {
    return {
      total: 1,
    };
  }

  get status(): CheckNodeStatus {
    // If we have results, this node is complete
    return "complete";
  }
}

export default DetectionResultNode;
