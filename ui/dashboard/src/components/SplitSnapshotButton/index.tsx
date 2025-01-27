import Icon from "@powerpipe/components/Icon";
import NeutralButton from "@powerpipe/components/forms/NeutralButton";
import useFilterConfig from "@powerpipe/hooks/useFilterConfig";
import useGroupingConfig from "@powerpipe/hooks/useGroupingConfig";
import useTableConfig from "@powerpipe/hooks/useTableConfig";
import { ChangeEvent, useRef } from "react";
import { classNames } from "@powerpipe/utils/styles";
import { DashboardActions, DashboardSnapshotMetadata } from "@powerpipe/types";
import { EXECUTION_SCHEMA_VERSION_20241125 } from "@powerpipe/constants/versions";
import {
  filterToSnapshotMetadata,
  groupingToSnapshotMetadata,
  stripSnapshotDataForExport,
  tableConfigToSnapshotMetadata,
} from "@powerpipe/utils/snapshot";
import { KeyValuePairs } from "@powerpipe/components/dashboards/common/types";
import { Menu } from "@headlessui/react";
import { noop } from "@powerpipe/utils/func";
import { saveAs } from "file-saver";
import { SnapshotDataToExecutionCompleteSchemaMigrator } from "@powerpipe/utils/schema";
import { timestampForFilename } from "@powerpipe/utils/date";
import { useDashboardExecution } from "@powerpipe/hooks/useDashboardExecution";
import { useDashboardInputs } from "@powerpipe/hooks/useDashboardInputs";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { validateFilter } from "@powerpipe/components/dashboards/grouping/FilterEditor";

const useSaveSnapshot = () => {
  const { dashboard, snapshot } = useDashboardState();
  const { inputs } = useDashboardInputs();
  const { allFilters } = useFilterConfig();
  const { allGroupings } = useGroupingConfig();
  const { allTables } = useTableConfig();

  return () => {
    if (!dashboard || !snapshot) {
      return;
    }
    const streamlinedSnapshot = stripSnapshotDataForExport(snapshot);
    const withMetadata = {
      ...streamlinedSnapshot,
    };

    if (
      !!Object.keys(allFilters).length ||
      !!Object.keys(allGroupings).length ||
      !!Object.keys(allTables).length
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
      if (!!Object.keys(allTables).length) {
        for (const [panel, tableConfig] of Object.entries(allTables)) {
          // @ts-ignore
          metadata.view[panel] = metadata.view[panel] || {};
          // @ts-ignore
          metadata.view[panel].table =
            tableConfigToSnapshotMetadata(tableConfig);
        }
      }
      withMetadata.metadata = metadata;
    }

    if (!!Object.keys(inputs).length) {
      withMetadata.inputs = inputs;
    }

    withMetadata.schema_version = EXECUTION_SCHEMA_VERSION_20241125;

    const blob = new Blob([JSON.stringify(withMetadata)], {
      type: "application/json",
    });
    saveAs(blob, `${dashboard.name}.${timestampForFilename(Date.now())}.pps`);
  };
};

const useOpenSnapshot = () => {
  const { dispatch } = useDashboardState();
  const { loadSnapshot } = useDashboardExecution();

  return (e: ChangeEvent<HTMLInputElement>) => {
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
        loadSnapshot(migratedEvent, fileName);
      } catch (err: any) {
        dispatch({
          type: DashboardActions.WORKSPACE_ERROR,
          error: "Unable to load snapshot:" + err.message,
        });
      }
    };

    // Max 500 MB size
    if (!files[0] || files[0].size > 500 * 1000 * 1000) {
      dispatch({
        type: DashboardActions.WORKSPACE_ERROR,
        error: "Snapshot exceeds the maximum supported size of 500 MB.",
      });
      return;
    }

    fr.readAsText(files[0]);
  };
};

interface SplitSnapshotAction {
  icon: string;
  disabled?: boolean;
  label: string;
  title: string;
  action: () => void;
}

const SplitSnapshotButton = () => {
  const { dashboard, selectedDashboard, snapshot } = useDashboardState();
  const openSnapshot = useOpenSnapshot();
  const saveSnapshot = useSaveSnapshot();
  const openSnapshotRef = useRef<HTMLInputElement | null>(null);

  let actions: KeyValuePairs<SplitSnapshotAction> = {
    open: {
      icon: "folder_open",
      // label: "Open snapshot…",
      label: "Open…",
      title: "Open snapshot…",
      action: () => openSnapshotRef.current?.click(),
    },
    save: {
      icon: "photo_camera",
      disabled: !dashboard || !snapshot,
      label: "Snap",
      title: "Take snapshot.",
      action: saveSnapshot,
    },
  };

  let defaultAction: SplitSnapshotAction;
  let otherActions: SplitSnapshotAction[] = [];
  if (selectedDashboard) {
    defaultAction = actions.save;
    otherActions = [actions.open];
  } else {
    defaultAction = actions.open;
  }

  return (
    <>
      <div className="inline-flex">
        <NeutralButton
          type="button"
          disabled={defaultAction.disabled}
          className={classNames(
            "relative inline-flex items-center space-x-2 focus:z-10",
            otherActions.length ? "rounded-r-none" : null,
          )}
          onClick={defaultAction.disabled ? noop : defaultAction.action}
        >
          <Icon icon={defaultAction.icon} className="h-5 w-5" />
          <span>{defaultAction.label}</span>
        </NeutralButton>
        {!!otherActions.length && (
          <Menu as="div" className="relative -ml-px block">
            <Menu.Button
              as={NeutralButton}
              className="relative flex items-center rounded-l-none h-full focus:z-10 py-2 px-1"
              size="manual"
              onClick={noop}
            >
              <span className="sr-only">Open options</span>
              <Icon icon="keyboard_arrow_down" className="h-5 w-5" />
            </Menu.Button>
            <Menu.Items className="absolute right-0 z-10 bg-dashboard-panel mt-px min-w-32 origin-top-right rounded-md transition focus:outline-none data-[closed]:scale-95 data-[closed]:transform data-[closed]:opacity-0 data-[enter]:duration-100 data-[leave]:duration-75 data-[enter]:ease-out data-[leave]:ease-in">
              <div>
                {otherActions.map((otherAction, idx) => (
                  <Menu.Item key={idx}>
                    <div
                      className={classNames(
                        "flex items-center space-x-2 p-2 cursor-pointer hover:bg-black-scale-2",
                        otherAction.disabled
                          ? "disabled:bg-dashboard disabled:text-light"
                          : "cursor-pointer",
                      )}
                      onClick={
                        otherAction.disabled ? undefined : otherAction.action
                      }
                    >
                      <Icon icon={otherAction.icon} className="h-5 w-5" />
                      <span className="hidden lg:block">
                        {otherAction.label}
                      </span>
                    </div>
                  </Menu.Item>
                ))}
              </div>
            </Menu.Items>
          </Menu>
        )}
      </div>
      <input
        ref={openSnapshotRef}
        accept="application/json,.pps,.sps"
        className="hidden"
        id="open-snapshot"
        name="open-snapshot"
        type="file"
        onChange={openSnapshot}
      />
    </>
  );
};

export default SplitSnapshotButton;
