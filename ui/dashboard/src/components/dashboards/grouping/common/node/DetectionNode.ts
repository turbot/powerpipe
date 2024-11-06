import DetectionHierarchyNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionHierarchyNode";
import { DetectionNode } from "../index";

class ControlNode extends DetectionHierarchyNode {
  constructor(
    sort: string,
    name: string,
    title: string | undefined,
    children?: DetectionNode[],
  ) {
    super("detection", name, title || name, sort, children || []);
  }
}

export default ControlNode;
