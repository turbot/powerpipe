import DetectionHierarchyNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionHierarchyNode";
import { DetectionNode } from "../index";

class RootNode extends DetectionHierarchyNode {
  constructor(children?: DetectionNode[]) {
    super("root", "root", "Root", "root", children || []);
  }
}

export default RootNode;
