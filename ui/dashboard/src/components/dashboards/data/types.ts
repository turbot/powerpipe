import { DashboardPanelType, PanelDataMode } from "@powerpipe/types";
import { LeafNodeData } from "@powerpipe/components/dashboards/common";

export interface IDiffProperties {
  data_mode?: PanelDataMode;
  diff_data?: LeafNodeData;
}

export interface IPanelDiff {}

export interface IPanelDataDiff {
  get panel_type(): DashboardPanelType;
  diff(
    properties: IDiffProperties,
    current_data: LeafNodeData,
    previous_data: LeafNodeData,
  ): IPanelDiff;
}
