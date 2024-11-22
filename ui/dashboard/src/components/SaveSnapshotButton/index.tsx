import Icon from "@powerpipe/components/Icon";
import NeutralButton from "@powerpipe/components/forms/NeutralButton";
import useGroupingFilterConfig from "@powerpipe/hooks/useGroupingFilterConfig";
import useCheckGroupingConfig from "@powerpipe/hooks/useCheckGroupingConfig";
import {
  DashboardDataModeCLISnapshot,
  DashboardSnapshotMetadata,
} from "@powerpipe/types";
import { EXECUTION_SCHEMA_VERSION_20240607 } from "@powerpipe/constants/versions";
import {
  filterToSnapshotMetadata,
  stripSnapshotDataForExport,
} from "@powerpipe/utils/snapshot";
import { saveAs } from "file-saver";
import { timestampForFilename } from "@powerpipe/utils/date";
import { useDashboard } from "@powerpipe/hooks/useDashboard";
import { validateFilter } from "@powerpipe/components/dashboards/grouping/CheckFilterEditor";

const SaveSnapshotButton = () => {
  const { dashboard, dataMode, selectedDashboard, snapshot } = useDashboard();
  const filterConfig = useGroupingFilterConfig();
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
      // If a benchmark
      if (
        dashboard.artificial &&
        !!filterConfig &&
        validateFilter(filterConfig)
      ) {
        // @ts-ignore
        metadata.view.filter_by = filterToSnapshotMetadata(filterConfig);
      }
      if (!!groupingConfig) {
        // @ts-ignore
        // TODO @mike re-include this
        // metadata.view.group_by = groupingToSnapshotMetadata(groupingConfig);
      }
      withMetadata.metadata = metadata;
      withMetadata.schema_version = EXECUTION_SCHEMA_VERSION_20240607;
    }

    const blob = new Blob([JSON.stringify(withMetadata)], {
      type: "application/json",
    });
    saveAs(blob, `${dashboard.name}.${timestampForFilename(Date.now())}.pps`);
  };

  if (
    dataMode === DashboardDataModeCLISnapshot ||
    (!selectedDashboard && !snapshot)
  ) {
    return null;
  }

  return (
    <NeutralButton
      className="inline-flex items-center space-x-2"
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
