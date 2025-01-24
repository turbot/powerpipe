import Badge from "@powerpipe/components/Badge";
import Icon from "@powerpipe/components/Icon";
import NeutralButton from "@powerpipe/components/forms/NeutralButton";
import SearchPathConfig from "@powerpipe/components/dashboards/SearchPath/SearchPathConfig";
import {
  DashboardDataModeCLISnapshot,
  DashboardDataModeCloudSnapshot,
} from "@powerpipe/types";
import { forwardRef, useEffect, useState } from "react";
import { Popover } from "@headlessui/react";
import { useDashboardSearchPath } from "@powerpipe/hooks/useDashboardSearchPath";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";

const PopoverButton = forwardRef((props, ref) => {
  const { metadata, dataMode, dashboardsMetadata, selectedDashboard } =
    useDashboardState();
  const { searchPathPrefix } = useDashboardSearchPath();

  if (
    dataMode === DashboardDataModeCLISnapshot ||
    dataMode === DashboardDataModeCloudSnapshot
  ) {
    return null;
  }

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
  const enabled = {
    enabled: hasServerMetadataSearchPath || hasDashboardMetadataSearchPath,
    hasServerMetadataSearchPath,
    hasDashboardMetadataSearchPath,
  };

  return (
    // @ts-ignore
    <NeutralButton
      ref={ref}
      className="inline-flex items-center space-x-2 h-full"
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
        {!searchPathPrefix.length ? (
          <span className="hidden lg:block">Search Path</span>
        ) : null}
        {searchPathPrefix.length >= 1 ? (
          <span className="hidden lg:block">{searchPathPrefix[0]}</span>
        ) : null}
        {!!searchPathPrefix.length && searchPathPrefix.length > 1 ? (
          <Badge>{searchPathPrefix.length}</Badge>
        ) : null}
      </>
    </NeutralButton>
  );
});

const ManageSearchPathButton = () => {
  const { metadata, dashboard, dashboardsMetadata } = useDashboardState();
  const [show, setShow] = useState(false);

  useEffect(() => {
    if (!metadata && !dashboardsMetadata && !dashboard) {
      return;
    }
    if (dashboard && dashboard.name in dashboardsMetadata) {
      setShow(!!dashboardsMetadata[dashboard.name]?.supports_search_path);
    } else {
      setShow(!!metadata?.supports_search_path);
    }
  }, [metadata?.supports_search_path, dashboard, dashboardsMetadata]);

  return show ? (
    <Popover className="hidden md:block relative">
      <Popover.Button as={PopoverButton} />
      <Popover.Panel className="absolute left-1/2 z-10 mt-4 flex w-screen max-w-max -translate-x-1/2 px-4">
        {({ close }) => (
          <div className="w-screen max-w-md flex-auto overflow-hidden rounded-md bg-dashboard border border-divide shadow-lg ring-1 ring-gray-900/5 p-4">
            <SearchPathConfig onClose={close} />
          </div>
        )}
      </Popover.Panel>
    </Popover>
  ) : null;
};

export default ManageSearchPathButton;
