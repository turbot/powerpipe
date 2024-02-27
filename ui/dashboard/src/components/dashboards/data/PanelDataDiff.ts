import { CardDataDiff } from "./CardDiff";
import { DashboardPanelType } from "@powerpipe/types";
import {
  IDiffProperties,
  IPanelDataDiff,
  IPanelDiff,
} from "@powerpipe/components/dashboards/data/types";
import { LeafNodeData } from "@powerpipe/components/dashboards/common";

export class PanelDataDiff implements IPanelDataDiff {
  private _panel_type: DashboardPanelType;

  constructor(panel_type: DashboardPanelType) {
    this._panel_type = panel_type;
  }

  calculate(
    properties: IDiffProperties,
    current: LeafNodeData | undefined,
    previous: LeafNodeData | undefined,
  ): IPanelDiff {
    switch (this.panel_type) {
      case "card":
        return new CardDataDiff().calculate(properties, current, previous);
      default:
        return {};
    }
  }

  get panel_type(): DashboardPanelType {
    return this._panel_type;
  }
}
