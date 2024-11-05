import DetectionHierarchyNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionHierarchyNode";
import { DetectionNode } from "../index";

class DetectionBenchmarkNode extends DetectionHierarchyNode {
  constructor(
    sort: string,
    name: string,
    title: string | undefined,
    children?: DetectionNode[],
  ) {
    super("detection_benchmark", name, title || name, sort, children || []);
  }
}

export default DetectionBenchmarkNode;
