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
      alarm_diff: 0,
      ok: 0,
      ok_diff: 0,
      info: 0,
      info_diff: 0,
      skip: 0,
      skip_diff: 0,
      error: 0,
      error_diff: 0,
      __diff: "none",
    };
    if (this._result.status_diff === "alarm") {
      summary.alarm_diff += 1;
      summary.__diff = "updated";
    }
    if (this._result.status_diff === "error") {
      summary.error_diff += 1;
      summary.__diff = "updated";
    }
    if (this._result.status_diff === "ok") {
      summary.ok_diff += 1;
      summary.__diff = "updated";
    }
    if (this._result.status_diff === "info") {
      summary.info_diff += 1;
      summary.__diff = "updated";
    }
    if (this._result.status_diff === "skip") {
      summary.skip_diff += 1;
      summary.__diff = "updated";
    }

    if (this._result.status === "alarm") {
      summary.alarm += 1;
      if (!this._result.status_diff) {
        summary.alarm_diff += 1;
      }
    }
    if (this._result.status === "error") {
      summary.error += 1;
      if (!this._result.status_diff) {
        summary.error_diff += 1;
      }
    }
    if (this._result.status === "ok") {
      summary.ok += 1;
      if (!this._result.status_diff) {
        summary.ok_diff += 1;
      }
    }
    if (this._result.status === "info") {
      summary.info += 1;
      if (!this._result.status_diff) {
        summary.info_diff += 1;
      }
    }
    if (this._result.status === "skip") {
      summary.skip += 1;
      if (!this._result.status_diff) {
        summary.skip_diff += 1;
      }
    }
    return summary;
  }

  get status(): CheckNodeStatus {
    // If we have results, this node is complete
    return "complete";
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
