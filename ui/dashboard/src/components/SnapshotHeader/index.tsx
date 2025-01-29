import { DashboardDataModeCLISnapshot } from "@powerpipe/types";
import { registerComponent } from "@powerpipe/components/dashboards";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";

const SnapshotHeader = () => {
  const { dataMode, rootPathname, snapshotFileName } = useDashboardState();

  if (dataMode !== DashboardDataModeCLISnapshot || rootPathname !== "/") {
    return null;
  }

  return (
    <div className="flex items-center space-x-2">
      <span className="font-medium">Snapshot:</span>
      <span className="text-foreground-lighter">{snapshotFileName}</span>
    </div>
  );
};

registerComponent("snapshot_header", SnapshotHeader);

export default SnapshotHeader;
