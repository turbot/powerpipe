import Icon from "components/Icon";
import { classNames } from "utils/styles";
import { DashboardDataModeCLISnapshot } from "types";
import { registerComponent } from "components/dashboards";
import { useDashboard } from "hooks/useDashboard";

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
        </div>
      </div>
    </>
  );
};

registerComponent("snapshot_header", SnapshotHeader);

export default SnapshotHeader;
