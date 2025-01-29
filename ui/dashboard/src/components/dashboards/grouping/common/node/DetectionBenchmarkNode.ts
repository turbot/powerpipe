import DetectionHierarchyNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionHierarchyNode";
import { DetectionNode } from "../index";

class DetectionBenchmarkNode extends DetectionHierarchyNode {
  private readonly _documentation: string | undefined;

  constructor(
    sort: string,
    name: string,
    title: string | undefined,
    documentation: string | undefined,
    children?: DetectionNode[],
  ) {
    super("benchmark", name, title || name, sort, children || []);
    this._documentation = documentation;
  }

  get documentation(): string | undefined {
    return this._documentation;
  }
}

export default DetectionBenchmarkNode;
