import {
  DetectionNode,
  DetectionSummary,
  DetectionNodeStatus,
  GroupingNodeType,
} from "../index";

class DetectionHierarchyNode implements DetectionNode {
  private readonly _type: GroupingNodeType;
  private readonly _name: string;
  private readonly _title: string;
  private readonly _sort: string;
  private readonly _children: DetectionNode[];

  constructor(
    type: GroupingNodeType,
    name: string,
    title: string,
    sort: string,
    children: DetectionNode[],
  ) {
    this._type = type;
    this._name = name;
    this._title = title;
    this._sort = sort;
    this._children = children;
  }

  get type(): GroupingNodeType {
    return this._type;
  }

  get name(): string {
    return this._name;
  }

  get title(): string {
    return this._title;
  }

  get sort(): string {
    return this._sort;
  }

  get children(): DetectionNode[] {
    return this._children;
  }

  get summary(): DetectionSummary {
    const summary = {
      total: 0,
    };
    for (const child of this._children) {
      const nestedSummary = child.summary;
      summary.total += nestedSummary.total;
    }
    return summary;
  }

  get status(): DetectionNodeStatus {
    for (const child of this._children) {
      if (child.status === "running") {
        return "running";
      }
    }
    return "complete";
  }

  merge(other: DetectionNode) {
    // merge(other) -> iterate children of other -> if child exists on me, call me_child.merge(other_child), else add to end of children
    for (const otherChild of other.children || []) {
      // Check for existing child with this name
      const matchingSelfChild = this.children.find(
        (selfChild) => selfChild.name === otherChild.name,
      );

      if (matchingSelfChild) {
        if (!matchingSelfChild.merge) {
          continue;
        }
        // If there's a matching child, merge that child in
        matchingSelfChild.merge(otherChild);
      } else {
        // Else append to my children
        this.children.push(otherChild);
      }
    }
  }
}

export default DetectionHierarchyNode;
