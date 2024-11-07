import DetectionHierarchyNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionHierarchyNode";
import { DetectionNode as DetectionNodeType } from "../index";

class DetectionNode extends DetectionHierarchyNode {
  constructor(
    sort: string,
    name: string,
    title: string | undefined,
    children?: DetectionNodeType[],
  ) {
    super("detection", name, title || name, sort, children || []);
  }
}

export default DetectionNode;
