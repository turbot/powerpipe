import Icon from "@powerpipe/components/Icon";
import NeutralButton from "@powerpipe/components/forms/NeutralButton";
import useFilterConfig from "@powerpipe/hooks/useFilterConfig";
import useGroupingConfig from "@powerpipe/hooks/useGroupingConfig";
import {
  DashboardDataModeCLISnapshot,
  DashboardSnapshotMetadata,
} from "@powerpipe/types";
import { EXECUTION_SCHEMA_VERSION_20241125 } from "@powerpipe/constants/versions";
import {
  filterToSnapshotMetadata,
  groupingToSnapshotMetadata,
  stripSnapshotDataForExport,
} from "@powerpipe/utils/snapshot";
import { saveAs } from "file-saver";
import { timestampForFilename } from "@powerpipe/utils/date";
import { useDashboard } from "@powerpipe/hooks/useDashboard";
import { validateFilter } from "@powerpipe/components/dashboards/grouping/FilterEditor";

const SaveSnapshotButton = () => {
  const { dashboard, dataMode, selectedDashboard, snapshot } = useDashboard();
  const { allFilters } = useFilterConfig();
  const { allGroupings } = useGroupingConfig();

  const saveSnapshot = () => {
    if (!dashboard || !snapshot) {
      return;
    }
    const streamlinedSnapshot = stripSnapshotDataForExport(snapshot);
    const withMetadata = {
      ...streamlinedSnapshot,
    };

    if (
      !!Object.keys(allFilters).length ||
      !!Object.keys(allGroupings).length
    ) {
      const metadata: DashboardSnapshotMetadata = {
        view: {},
      };
      if (!!Object.keys(allFilters).length) {
        for (const [panel, filter] of Object.entries(allFilters)) {
          if (!validateFilter(filter)) {
            console.warn("Ignoring invalid filter", { panel, filter });
            continue;
          }
          // @ts-ignore
          metadata.view[panel] = metadata.view[panel] || {};
          // @ts-ignore
          metadata.view[panel].filter_by = filterToSnapshotMetadata(filter);
        }
      }
      if (!!Object.keys(allGroupings).length) {
        for (const [panel, grouping] of Object.entries(allGroupings)) {
          // @ts-ignore
          metadata.view[panel] = metadata.view[panel] || {};
          // @ts-ignore
          metadata.view[panel].group_by = groupingToSnapshotMetadata(grouping);
        }
      }
      withMetadata.metadata = metadata;
      withMetadata.schema_version = EXECUTION_SCHEMA_VERSION_20241125;
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
