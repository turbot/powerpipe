import {
  BasePrimitiveProps,
  ExecutablePrimitiveProps,
} from "@powerpipe/components/dashboards/common";
import {
  CategoryMap,
  NodeAndEdgeProperties,
} from "@powerpipe/components/dashboards/common/types";
import { ComponentType } from "react";
import { NodeAndEdgeData } from "@powerpipe/components/dashboards/graphs/types";
import { PanelDefinition } from "@powerpipe/types";

export type BaseFlowProps = PanelDefinition &
  BasePrimitiveProps &
  ExecutablePrimitiveProps;

export type FlowProperties = NodeAndEdgeProperties;

export type FlowProps = BaseFlowProps & {
  categories: CategoryMap;
  data?: NodeAndEdgeData;
  display_type?: FlowType;
  properties?: NodeAndEdgeProperties;
};

export type FlowType = "sankey" | "table";

export type IFlow = {
  type: FlowType;
  component: ComponentType<any>;
};
