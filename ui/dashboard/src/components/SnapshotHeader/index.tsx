import Icon from "@powerpipe/components/Icon";
import SnapshotDiffButton from "@powerpipe/components/SnapshotDiffButton";
import { classNames } from "@powerpipe/utils/styles";
import { DashboardDataModeCLISnapshot } from "@powerpipe/types";
import { registerComponent } from "@powerpipe/components/dashboards";
import { useDashboard } from "@powerpipe/hooks/useDashboard";

const SnapshotHeader = () => {
  const { dataMode, snapshotFileName } = useDashboard();

  if (dataMode !== DashboardDataModeCLISnapshot) {
    return null;
  }

  return (
    <>
      <div className={classNames("space-y-2")}>
        <div className="flex items-center space-x-3">
          <Icon className="h-5 w-5" icon="photo_camera" />
          <span className="font-medium">Snapshot:</span>
          <span className="text-foreground-lighter">{snapshotFileName}</span>
          <SnapshotDiffButton />
        </div>
      </div>
    </>
  );
};

registerComponent("snapshot_header", SnapshotHeader);

export default SnapshotHeader;
