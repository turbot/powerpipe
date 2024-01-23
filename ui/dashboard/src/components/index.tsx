import DashboardListEmptyCallToAction from "./DashboardListEmptyCallToAction";
import { ComponentsMap } from "../types";
import { getComponent } from "components/dashboards";

const buildComponentsMap = (overrides = {}): ComponentsMap => {
  const SnapshotHeader = getComponent("snapshot_header");
  
  return {
    DashboardListEmptyCallToAction,
    SnapshotHeader,
    ...overrides,
  };
};

export { buildComponentsMap };
