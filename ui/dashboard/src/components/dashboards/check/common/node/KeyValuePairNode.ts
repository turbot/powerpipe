import HierarchyNode from "./HierarchyNode";
import { CheckNodeType, CheckNode } from "../index";

class KeyValuePairNode extends HierarchyNode {
  private readonly _key: string;
  private readonly _value: string;

  constructor(
    sort: string,
    type: CheckNodeType,
    key: string,
    value: string,
    children?: CheckNode[],
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

export default KeyValuePairNode;
