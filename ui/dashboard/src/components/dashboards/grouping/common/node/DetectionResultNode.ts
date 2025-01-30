import {
  CheckNodeStatus,
  DetectionNode,
  DetectionResult,
  DetectionSeveritySummary,
  DetectionSummary,
  GroupingNodeType,
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

  get severity_summary(): DetectionSeveritySummary {
    const summary: DetectionSeveritySummary = {};
    if (this._result.detection.severity) {
      summary[this._result.detection.severity] =
        this._result?.rows?.length || 0;
    } else {
      summary["none"] = 0;
    }
    return summary;
  }

  get summary(): DetectionSummary {
    return {
      total: this._result?.rows?.length || 0,
      error: this._result.error ? 1 : 0,
    };
  }

  get status(): CheckNodeStatus {
    // If we have results, this node is complete
    return "complete";
  }
}

export default DetectionResultNode;
