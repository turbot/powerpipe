import {
  BasePrimitiveProps,
  ExecutablePrimitiveProps,
} from "@powerpipe/components/dashboards/common";
import { ComponentType } from "react";
import {
  CategoryMap,
  NodeAndEdgeProperties,
} from "@powerpipe/components/dashboards/common/types";
import { NodeAndEdgeData } from "@powerpipe/components/dashboards/graphs/types";
import { PanelDefinition } from "@powerpipe/types";

export type BaseHierarchyProps = PanelDefinition &
  BasePrimitiveProps &
  ExecutablePrimitiveProps;

export type HierarchyProperties = NodeAndEdgeProperties;

export type HierarchyProps = BaseHierarchyProps & {
  categories: CategoryMap;
  data?: NodeAndEdgeData;
  display_type?: HierarchyType;
  properties?: NodeAndEdgeProperties;
};

export type HierarchyType = "table" | "tree";

export type IHierarchy = {
  type: HierarchyType;
  component: ComponentType<any>;
};
