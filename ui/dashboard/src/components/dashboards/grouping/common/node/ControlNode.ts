import HierarchyNode from "./HierarchyNode";
import { CheckNode } from "../index";

class ControlNode extends HierarchyNode {
  private readonly _documentation: string | undefined;

  constructor(
    sort: string,
    name: string,
    title: string | undefined,
    documentation: string | undefined,
    children?: CheckNode[],
  ) {
    super("control", name, title || name, sort, children || []);
    this._documentation = documentation;
  }

  get documentation(): string | undefined {
    return this._documentation;
  }
}

export default ControlNode;
