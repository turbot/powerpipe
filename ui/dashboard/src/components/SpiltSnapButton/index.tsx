import NeutralButton from "@powerpipe/components/forms/NeutralButton";
import Icon from "@powerpipe/components/Icon";
import useFilterConfig from "@powerpipe/hooks/useFilterConfig";
import useGroupingConfig from "@powerpipe/hooks/useGroupingConfig";
import { ChangeEvent, useRef } from "react";
import {
  DashboardDataModeCLISnapshot,
  DashboardDataModeLive,
  DashboardActions,
  DashboardSnapshotMetadata,
  DashboardDataModeDiff,
} from "@powerpipe/types";
import { EXECUTION_SCHEMA_VERSION_20241125 } from "@powerpipe/constants/versions";
import {
  filterToSnapshotMetadata,
  groupingToSnapshotMetadata,
  stripSnapshotDataForExport,
} from "@powerpipe/utils/snapshot";
import { Menu } from "@headlessui/react";
import { saveAs } from "file-saver";
import { SnapshotDataToExecutionCompleteSchemaMigrator } from "@powerpipe/utils/schema";
import { timestampForFilename } from "@powerpipe/utils/date";
import { useDashboard } from "@powerpipe/hooks/useDashboard";
import { useNavigate } from "react-router-dom";
import { validateFilter } from "@powerpipe/components/dashboards/grouping/FilterEditor";

interface SplitButtonProps {
  className?: string;
}

const SplitButton = ({ className }: SplitButtonProps) => {
  const { dashboard, dataMode, selectedDashboard, snapshot } = useDashboard();
  const { allFilters } = useFilterConfig();
  const { allGroupings } = useGroupingConfig();
  const { dispatch } = useDashboard();
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const fileInputRefForDiff = useRef<HTMLInputElement | null>(null);
  const navigate = useNavigate();

  const isDashboardList = !selectedDashboard && !snapshot;

  const isLive = selectedDashboard && dataMode === DashboardDataModeLive;

  const isSnapshot = snapshot && dataMode === DashboardDataModeCLISnapshot;

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

  const handleDiffOpen = (e: ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (!files || files.length === 0) {
      return;
    }
    const fr = new FileReader();
    fr.onload = () => {
      if (!fr.result) {
        return;
      }

      e.target.value = "";
      try {
        navigate(`/snapshot/diff`);
        const data = JSON.parse(fr.result.toString());
        dispatch({
          type: DashboardActions.GET_SNAPSHOT_DIFF,
          snapshot: data,
        });
      } catch (err: any) {
        dispatch({
          type: DashboardActions.WORKSPACE_ERROR,
          error: "Unable to load snapshot:" + err.message,
        });
      }
    };
    fr.readAsText(files[0]);
  };

  const handleFileOpen = (e: ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (!files || files.length === 0) {
      return;
    }
    const fileName = files[0].name;
    const fr = new FileReader();
    fr.onload = () => {
      if (!fr.result) {
        return;
      }

      e.target.value = "";
      try {
        const data = JSON.parse(fr.result.toString());
        const eventMigrator =
          new SnapshotDataToExecutionCompleteSchemaMigrator();
        const migratedEvent = eventMigrator.toLatest(data);
        dispatch({
          type: DashboardActions.CLEAR_DASHBOARD_INPUTS,
          recordInputsHistory: false,
        });
        dispatch({
          type: DashboardActions.SELECT_DASHBOARD,
          dashboard: null,
          recordInputsHistory: false,
        });
        navigate(`/snapshot/${fileName}`);
        dispatch({
          type: DashboardActions.SET_DATA_MODE,
          dataMode: DashboardDataModeCLISnapshot,
          snapshotFileName: fileName,
        });
        dispatch({
          type: DashboardActions.EXECUTION_COMPLETE,
          ...migratedEvent,
        });
        dispatch({
          type: DashboardActions.SET_DASHBOARD_INPUTS,
          value: migratedEvent.snapshot.inputs,
          recordInputsHistory: false,
        });
      } catch (err: any) {
        dispatch({
          type: DashboardActions.WORKSPACE_ERROR,
          error: "Unable to load snapshot:" + err.message,
        });
      }
    };
    fr.readAsText(files[0]);
  };

  return (
    <div className={`relative inline-block text-left ${className}`}>
      <div className="flex ">
        {/* Main Button */}
        {/* <NeutralButton
          type="button"
          className="inline-flex items-center space-x-2 shadow-none rounded-r-none"
          onClick={handleMainButtonClick}
        ></NeutralButton> */}
        {isDashboardList && (
          <NeutralButton
            type="button"
            className="inline-flex items-center space-x-2 shadow-none rounded-r-none"
            onClick={() => fileInputRef.current?.click()}
          >
            <Icon
              className="inline-block text-foreground-lighter w-5 -mt-0.5"
              icon="heroicons-outline:folder-open"
            />
            <span className="hidden lg:block">Open</span>
            <input
              ref={fileInputRef}
              accept="application/json, .pps, .sps"
              className="hidden"
              id="open-snapshot"
              name="open-snapshot"
              type="file"
              onChange={handleFileOpen}
            />
          </NeutralButton>
        )}

        {(isSnapshot || dataMode == DashboardDataModeDiff) && (
          <NeutralButton
            type="button"
            className="inline-flex items-center space-x-2 shadow-none rounded-r-none"
            onClick={() => fileInputRefForDiff.current?.click()}
          >
            <Icon
              className="inline-block text-foreground-lighter w-5 h-5 -mt-0.5"
              icon="difference"
            />
            <span className="hidden lg:block">Diff</span>
            <input
              ref={fileInputRefForDiff}
              accept="application/json, .pps, .sps"
              className="hidden"
              id="snapshot-diff"
              name="snapshot-diff"
              type="file"
              onChange={handleDiffOpen}
            />
          </NeutralButton>
        )}

        {isLive && (
          <NeutralButton
            type="button"
            className="inline-flex items-center space-x-2 shadow-none rounded-r-none"
            onClick={saveSnapshot}
          >
            <Icon
              className="inline-block text-foreground-lighter w-5 h-5 -mt-0.5"
              icon="heroicons-outline:camera"
            />
            <span className="hidden lg:block">Snap</span>
          </NeutralButton>
        )}

        {!isDashboardList && (
          <Menu as="div" className="relative flex border rounded-r-md">
            <Menu.Button
              className="px-2 py-2 text-foreground font-medium text-sm rounded-r-md"
              aria-haspopup="true"
            >
              {/* <ChevronDownIcon className="w-5 h-5" aria-hidden="true" /> */}
              <Icon
                className="h-5 w-5"
                icon="materialsymbols-solid:arrow_drop_down"
              />
            </Menu.Button>
            <Menu.Items className="absolute top-full right-0 w-40 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 z-50 divide-y divide-gray-100 focus:outline-none">
              {isLive && (
                <Menu.Item>
                  {({ active }) => (
                    <NeutralButton
                      type="button"
                      className="inline-flex items-center space-x-2 shadow-none rounded-none w-full"
                      onClick={saveSnapshot}
                    >
                      <Icon
                        className="inline-block text-foreground-lighter w-5 h-5 -mt-0.5"
                        icon="heroicons-outline:camera"
                      />
                      <span className="hidden lg:block">Snap</span>
                    </NeutralButton>
                  )}
                </Menu.Item>
              )}
              <Menu.Item>
                {({ active }) => (
                  <>
                    <NeutralButton
                      type="button"
                      className="inline-flex items-center space-x-2 shadow-none rounded-none w-full"
                      onClick={() => {
                        if (fileInputRefForDiff.current) {
                          fileInputRefForDiff.current.click(); // Ensure this triggers the click on the input
                        }
                      }}
                    >
                      <Icon
                        className="inline-block text-foreground-lighter w-5 h-5"
                        icon="difference"
                      />
                      <span className="hidden lg:block">Diff</span>
                    </NeutralButton>
                    <input
                      ref={fileInputRefForDiff}
                      accept="application/json, .pps, .sps"
                      className="hidden"
                      id="snapshot-diff"
                      name="snapshot-diff"
                      type="file"
                      onChange={handleDiffOpen}
                    />
                  </>
                )}
              </Menu.Item>
              <Menu.Item>
                {({ active }) => (
                  <>
                    <NeutralButton
                      type="button"
                      className="inline-flex items-center space-x-2 shadow-none rounded-none w-full"
                      onClick={() => {
                        if (fileInputRef.current) {
                          fileInputRef.current.click(); // Ensure this triggers the click on the input
                        }
                      }}
                    >
                      <Icon
                        className="inline-block text-foreground-lighter w-5 -mt-0.5"
                        icon="heroicons-outline:folder-open"
                      />
                      <span className="hidden lg:block">Open</span>
                    </NeutralButton>
                    {/* Move the input outside of NeutralButton */}
                    <input
                      ref={fileInputRef}
                      accept="application/json, .pps, .sps"
                      className="hidden"
                      id="open-snapshot"
                      name="open-snapshot"
                      type="file"
                      onChange={handleFileOpen}
                    />
                  </>
                )}
              </Menu.Item>
            </Menu.Items>
          </Menu>
        )}
      </div>
    </div>
  );
};

export default SplitButton;
