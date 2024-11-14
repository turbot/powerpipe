import React, { useState, useRef } from "react";
import { Menu } from "@headlessui/react";
import { SnapshotDataToExecutionCompleteSchemaMigrator } from "@powerpipe/utils/schema";
import { useNavigate } from "react-router-dom";
import { ChevronDownIcon } from "@heroicons/react/solid";
import NeutralButton from "@powerpipe/components/forms/NeutralButton";
import Icon from "@powerpipe/components/Icon";
import useGroupingFilterConfig from "@powerpipe/hooks/useGroupingFilterConfig";
import useCheckGroupingConfig from "@powerpipe/hooks/useCheckGroupingConfig";
import {
  DashboardDataModeCLISnapshot,
  DashboardDataModeLive,
  DashboardActions,
  DashboardSnapshotMetadata,
} from "@powerpipe/types";
import { EXECUTION_SCHEMA_VERSION_20240607 } from "@powerpipe/constants/versions";
import {
  filterToSnapshotMetadata,
  groupingToSnapshotMetadata,
  stripSnapshotDataForExport,
} from "@powerpipe/utils/snapshot";
import { saveAs } from "file-saver";
import { timestampForFilename } from "@powerpipe/utils/date";
import { useDashboard } from "@powerpipe/hooks/useDashboard";
import { validateFilter } from "@powerpipe/components/dashboards/grouping/CheckFilterEditor";

interface SplitButtonProps {
  className?: string;
}

const SplitButton: React.FC<SplitButtonProps> = ({ className }) => {
  const { dashboard, dataMode, selectedDashboard, snapshot } = useDashboard();
  const filterConfig = useGroupingFilterConfig();
  const groupingConfig = useCheckGroupingConfig();
  const { dispatch } = useDashboard();
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const fileInputRefForDiff = useRef<HTMLInputElement | null>(null);
  const navigate = useNavigate();

  const isDashBoard = !selectedDashboard && !snapshot;

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

  const handleDiffOpen = (e: React.ChangeEvent<HTMLInputElement>) => {
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

  const handleFileOpen = (e: React.ChangeEvent<HTMLInputElement>) => {
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

      e.target.value = ""; // Clear the input for repeated use
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
          error: "Unable to load snapshot: " + err.message,
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
        {isDashBoard && (
          <NeutralButton
            type="button"
            className="inline-flex items-center space-x-2 shadow-none rounded-r-none"
          >
            <Icon
              className="inline-block text-foreground-lighter w-5 -mt-0.5"
              icon="heroicons-outline:folder-open"
            />
            <span
              className="hidden lg:block"
              onClick={() => fileInputRef.current?.click()}
            >
              Open
            </span>
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

        {isSnapshot && (
          <NeutralButton
            type="button"
            className="inline-flex items-center space-x-2 shadow-none rounded-r-none"
          >
            <Icon
              className="inline-block text-foreground-lighter w-5 h-5 -mt-0.5"
              icon="difference"
            />
            <span
              className="hidden lg:block"
              onClick={() => fileInputRefForDiff.current?.click()}
            >
              Diff
            </span>
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

        {!isDashBoard && (
          <Menu as="div" className="relative flex border rounded-r-md">
            <Menu.Button
              className="px-2 py-2 text-foreground font-medium text-sm rounded-r-md rounded-r-md"
              aria-haspopup="true"
            >
              {/* <ChevronDownIcon className="w-5 h-5" aria-hidden="true" /> */}
              <Icon
                className="h-5 w-5"
                icon="materialsymbols-solid:arrow_drop_down"
              />
            </Menu.Button>
            <Menu.Items className="absolute top-full left-0 w-40 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 z-50 divide-y divide-gray-100 focus:outline-none">
              <div className="py-1">
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
                          className="inline-block text-foreground-lighter w-5 h-5 -mt-0.5"
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
              </div>
            </Menu.Items>
          </Menu>
        )}
      </div>
    </div>
  );
};

export default SplitButton;
