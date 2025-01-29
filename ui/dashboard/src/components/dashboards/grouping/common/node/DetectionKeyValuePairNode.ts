import DetectionHierarchyNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionHierarchyNode";
import { GroupingNodeType, DetectionNode } from "../index";

class DetectionKeyValuePairNode extends DetectionHierarchyNode {
  private readonly _key: string;
  private readonly _value: string;

  constructor(
    sort: string,
    type: GroupingNodeType,
    key: string,
    value: string,
    children?: DetectionNode[],
  ) {
    super(type, `${key}=${value}`, value, sort, children || []);
    this._key = key;
    this._value = value;
  }

  get key(): string {
    return this._key;
  }

  get value(): string {
    return this._value;
  }
}

export default DetectionKeyValuePairNode;
