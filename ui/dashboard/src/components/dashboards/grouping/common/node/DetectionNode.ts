import DetectionHierarchyNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionHierarchyNode";
import { DetectionNode as DetectionNodeType } from "../index";

class DetectionNode extends DetectionHierarchyNode {
  private readonly _documentation: string | undefined;

  constructor(
    sort: string,
    name: string,
    title: string | undefined,
    documentation: string | undefined,
    children?: DetectionNodeType[],
  ) {
    super("detection", name, title || name, sort, children || []);
    this._documentation = documentation;
  }

  get documentation(): string | undefined {
    return this._documentation;
  }
}

export default DetectionNode;
