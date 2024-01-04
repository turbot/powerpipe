import DiffSnapshotButton from "../DiffSnapshotButton";
import Icon from "../Icon";
import { classNames } from "../../utils/styles";
import { DashboardDataModeCLISnapshot } from "../../types";
import { useDashboard } from "../../hooks/useDashboard";

const SnapshotHeader = () => {
  const { dataMode, diff, snapshotFileName } = useDashboard();

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
          {(!diff || !diff.snapshotFileName) && <DiffSnapshotButton />}
        </div>
        {!!diff && diff.snapshotFileName && (
          <div className="flex items-center space-x-3">
            <Icon className="h-5 w-5" icon="difference" />
            <span className="font-medium">Diff:</span>
            <span className="text-foreground-lighter">
              {diff.snapshotFileName}
            </span>
          </div>
        )}
      </div>
    </>
  );
};

export default SnapshotHeader;
