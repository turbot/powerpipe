import {
  CheckNodeStatus,
  GroupingNodeType,
  CheckSummary,
  CheckNode,
  CheckResult,
  CheckSeveritySummary,
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

class ControlResultNode implements CheckNode {
  private readonly _result: CheckResult;

  constructor(result: CheckResult) {
    this._result = result;
  }

  get sort(): string {
    return "0";
  }

  get name(): string {
    return `${this._result.control.name}-${this._result.resource}`;
  }

  get title(): string {
    return this._result.reason;
  }

  get result(): CheckResult {
    return this._result;
  }

  get type(): GroupingNodeType {
    return "result";
  }

  get summary(): CheckSummary {
    const summary = {
      alarm: 0,
      ok: 0,
      info: 0,
      skip: 0,
      skipped: 0,
      error: 0,
      invalid: 0,
      muted: 0,
      tbd: 0,
    };
    switch (this._result.status) {
      case "alarm":
        summary.alarm += 1;
        break;
      case "error":
        summary.error += 1;
        break;
      case "invalid":
        summary.invalid += 1;
        break;
      case "ok":
        summary.ok += 1;
        break;
      case "info":
        summary.info += 1;
        break;
      case "skip":
        summary.skip += 1;
        break;
      case "skipped":
        summary.skipped += 1;
        break;
      case "muted":
        summary.muted += 1;
        break;
      case "tbd":
        summary.tbd += 1;
        break;
    }
    return summary;
  }

  get status(): CheckNodeStatus {
    // If we have results, this node is complete
    return "complete";

    // for (const child of this.children || []) {
    //   if (child.status === "running") {
    //     return "running";
    //   }
    // }
    // return "complete";
  }

  get severity_summary(): CheckSeveritySummary {
    const summary = {};
    if (this._result.control.severity) {
      summary[this._result.control.severity] =
        this._result.status === "alarm" ? 1 : 0;
    }
    return summary;
  }
}

export default ControlResultNode;
