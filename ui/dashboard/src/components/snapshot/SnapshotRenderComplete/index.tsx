import {
  DashboardDataModeCLISnapshot,
  DashboardDataModeCloudSnapshot,
} from "@powerpipe/types";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";

const SnapshotRenderComplete = () => {
  const { dataMode, state } = useDashboardState();

  if (
    dataMode === DashboardDataModeCLISnapshot ||
    dataMode === DashboardDataModeCloudSnapshot ||
    state !== "complete"
  ) {
    return null;
  }

  return <div id="snapshot-complete" className="hidden" />;
};

export default SnapshotRenderComplete;
