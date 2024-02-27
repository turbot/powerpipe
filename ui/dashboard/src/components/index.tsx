import DashboardListEmptyCallToAction from "./DashboardListEmptyCallToAction";
import { ComponentsMap } from "@powerpipe/types";
import { getComponent } from "@powerpipe/components/dashboards";

const buildComponentsMap = (overrides = {}): ComponentsMap => {
  const SnapshotHeader = getComponent("snapshot_header");

  return {
    DashboardListEmptyCallToAction,
    SnapshotHeader,
    ...overrides,
  };
};

export { buildComponentsMap };
