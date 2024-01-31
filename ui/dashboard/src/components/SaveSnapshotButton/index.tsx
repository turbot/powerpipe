import Icon from "components/Icon";
import NeutralButton from "components/forms/NeutralButton";
import useCheckFilterConfig from "hooks/useCheckFilterConfig";
import useCheckGroupingConfig from "hooks/useCheckGroupingConfig";
import { DashboardDataModeCLISnapshot, DashboardSnapshotMetadata } from "types";
import {
  filterToSnapshotMetadata,
  groupingToSnapshotMetadata,
  stripSnapshotDataForExport,
} from "utils/snapshot";
import { saveAs } from "file-saver";
import { timestampForFilename } from "utils/date";
import { useDashboard } from "hooks/useDashboard";

const SaveSnapshotButton = () => {
  const { dashboard, dataMode, selectedDashboard, snapshot } = useDashboard();
  const filterConfig = useCheckFilterConfig();
  const groupingConfig = useCheckGroupingConfig();

  const saveSnapshot = () => {
    if (!dashboard || !snapshot) {
      return;
    }
    const streamlinedSnapshot = stripSnapshotDataForExport(snapshot);
    const withMetadata = {
      ...streamlinedSnapshot,
    };

    if (!!filterConfig || !!groupingConfig) {
      const metadata: DashboardSnapshotMetadata = {
        view: {},
      };
      if (!!filterConfig) {
        // @ts-ignore
        metadata.view.filter_by = filterToSnapshotMetadata(filterConfig);
      }
      if (!!groupingConfig) {
        // @ts-ignore
        metadata.view.group_by = groupingToSnapshotMetadata(groupingConfig);
      }
      withMetadata.metadata = metadata;
    }

    const blob = new Blob([JSON.stringify(withMetadata)], {
      type: "application/json",
    });
    saveAs(blob, `${dashboard.name}.${timestampForFilename(Date.now())}.sps`);
  };

  if (
    dataMode === DashboardDataModeCLISnapshot ||
    (!selectedDashboard && !snapshot)
  ) {
    return null;
  }

  return (
    <NeutralButton
      className="inline-flex items-center space-x-1"
      disabled={!dashboard || !snapshot}
      onClick={saveSnapshot}
    >
      <>
        <Icon
          className="inline-block text-foreground-lighter w-5 -mt-0.5"
          icon="heroicons-outline:camera"
        />
        <span className="hidden lg:block">Snap</span>
      </>
    </NeutralButton>
  );
};

export default SaveSnapshotButton;
