import DetectionHierarchyNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionHierarchyNode";
import { DetectionNode as DetectionNodeType } from "../index";
import { LeafNodeData } from "@powerpipe/components/dashboards/common";

class DetectionNode extends DetectionHierarchyNode {
  private readonly _documentation: string | undefined;
  private readonly _data: LeafNodeData | undefined;

  constructor(
    sort: string,
    name: string,
    title: string | undefined,
    documentation: string | undefined,
    data: LeafNodeData | undefined,
    children?: DetectionNodeType[],
  ) {
    super("detection", name, title || name, sort, children || []);
    this._documentation = documentation;
    this._data = data;
  }

  get documentation(): string | undefined {
    return this._documentation;
  }

  get data(): LeafNodeData | undefined {
    return this._data;
  }
}

export default DetectionNode;
