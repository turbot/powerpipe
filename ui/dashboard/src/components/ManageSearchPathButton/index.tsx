import Icon from "@powerpipe/components/Icon";
import NeutralButton from "@powerpipe/components/forms/NeutralButton";
import SearchPathConfig from "../dashboards/SearchPathConfig";
import {
  DashboardDataModeCLISnapshot,
  DashboardDataModeCloudSnapshot,
} from "@powerpipe/types";
import { forwardRef, Fragment, useMemo } from "react";
import { Popover, Transition } from "@headlessui/react";
import { useDashboard } from "@powerpipe/hooks/useDashboard";

const PopoverButton = forwardRef((props, ref) => {
  const {
    metadata,
    dataMode,
    dashboardsMetadata,
    selectedDashboard,
    snapshot,
  } = useDashboard();
  const {
    enabled,
    hasServerMetadataSearchPath,
    hasDashboardMetadataSearchPath,
  } = useMemo(() => {
    const hasServerMetadataSearchPath =
      !!metadata?.search_path?.original_search_path &&
      !!metadata.search_path.original_search_path.length;
    const hasDashboardMetadataSearchPath =
      !!selectedDashboard &&
      !!dashboardsMetadata &&
      !!dashboardsMetadata[selectedDashboard.full_name] &&
      !!dashboardsMetadata[selectedDashboard.full_name]?.search_path &&
      !!dashboardsMetadata[selectedDashboard.full_name]?.search_path
        ?.original_search_path &&
      !!dashboardsMetadata[selectedDashboard.full_name]?.search_path
        ?.original_search_path?.length;
    return {
      enabled: hasServerMetadataSearchPath || hasDashboardMetadataSearchPath,
      hasServerMetadataSearchPath,
      hasDashboardMetadataSearchPath,
    };
  }, [selectedDashboard, metadata, dashboardsMetadata, snapshot]);

  if (
    dataMode === DashboardDataModeCLISnapshot ||
    dataMode === DashboardDataModeCloudSnapshot
  ) {
    return null;
  }

  return (
    // @ts-ignore
    <NeutralButton
      ref={ref}
      className="inline-flex items-center space-x-2"
      disabled={!enabled}
      title={
        enabled
          ? undefined
          : !!selectedDashboard &&
              !hasDashboardMetadataSearchPath &&
              !hasServerMetadataSearchPath
            ? "No dashboard search path available"
            : !hasServerMetadataSearchPath
              ? "No server search path available"
              : undefined
      }
      {...props}
    >
      <>
        <Icon
          className="inline-block text-foreground-lighter w-5 h-5"
          icon="sort"
        />
        <span className="hidden lg:block">Search Path</span>
      </>
    </NeutralButton>
  );
});

const ManageSearchPathButton = () => {
  return (
    <Popover className="relative">
      <Popover.Button as={PopoverButton} />
      <Transition
        as={Fragment}
        enter="transition ease-out duration-200"
        enterFrom="opacity-0 translate-y-1"
        enterTo="opacity-100 translate-y-0"
        leave="transition ease-in duration-150"
        leaveFrom="opacity-100 translate-y-0"
        leaveTo="opacity-0 translate-y-1"
      >
        <Popover.Panel className="absolute left-1/2 z-10 mt-4 flex w-screen max-w-max -translate-x-1/2 px-4">
          {({ close }) => (
            <div className="w-screen max-w-md flex-auto overflow-hidden rounded-md bg-dashboard border border-divide shadow-lg ring-1 ring-gray-900/5 p-4">
              <SearchPathConfig onClose={close} />
            </div>
          )}
        </Popover.Panel>
      </Transition>
    </Popover>
  );
};

export default ManageSearchPathButton;
