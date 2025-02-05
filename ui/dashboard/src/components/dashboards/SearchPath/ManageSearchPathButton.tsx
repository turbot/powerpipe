import Badge from "@powerpipe/components/Badge";
import Icon from "@powerpipe/components/Icon";
import NeutralButton from "@powerpipe/components/forms/NeutralButton";
import SearchPathConfig from "@powerpipe/components/dashboards/SearchPath/SearchPathConfig";
import { createPortal } from "react-dom";
import {
  DashboardDataModeCLISnapshot,
  DashboardDataModeCloudSnapshot,
} from "@powerpipe/types";
import { forwardRef, useEffect, useState } from "react";
import { Popover } from "@headlessui/react";
import { ThemeProvider, ThemeWrapper } from "@powerpipe/hooks/useTheme";
import { useDashboardSearchPath } from "@powerpipe/hooks/useDashboardSearchPath";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { usePopper } from "react-popper";

const PopoverButton = forwardRef((props, ref) => {
  const { metadata, dashboardsMetadata, selectedDashboard } =
    useDashboardState();
  const { searchPathPrefix } = useDashboardSearchPath();

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
  const enabled = hasServerMetadataSearchPath || hasDashboardMetadataSearchPath;

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
      {searchPathPrefix.length > 1 ? (
        <Badge>{searchPathPrefix.length}</Badge>
      ) : null}
    </NeutralButton>
  );
});

const ManageSearchPathButton = () => {
  const [popperElement, setPopperElement] = useState(null);
  const [referenceElement, setReferenceElement] = useState(null);
  const { styles, attributes } = usePopper(referenceElement, popperElement, {
    placement: "bottom-start",
    modifiers: [
      {
        name: "flip",
        options: {
          fallbackPlacements: ["bottom-end"],
        },
      },
    ],
  });
  const { dashboard, dashboardsMetadata, dataMode, metadata } =
    useDashboardState();
  const { searchPathPrefix } = useDashboardSearchPath();
  const [show, setShow] = useState(false);

  useEffect(() => {
    if (
      (dataMode === DashboardDataModeCLISnapshot ||
        dataMode === DashboardDataModeCloudSnapshot) &&
      searchPathPrefix.length
    ) {
      setShow(true);
      return;
    }

    if (!metadata && !dashboardsMetadata && !dashboard) {
      return;
    }
    if (dashboard && dashboard.name in dashboardsMetadata) {
      setShow(!!dashboardsMetadata[dashboard.name]?.supports_search_path);
    } else {
      setShow(!!metadata?.supports_search_path);
    }
  }, [metadata, dashboard, dashboardsMetadata, dataMode, searchPathPrefix]);

  return show ? (
    <Popover className="hidden md:block relative">
      <Popover.Button
        ref={setReferenceElement}
        disabled={
          dataMode === DashboardDataModeCLISnapshot ||
          dataMode === DashboardDataModeCloudSnapshot
        }
        as={PopoverButton}
      />
      <Popover.Panel className="absolute z-10">
        {({ close }) => (
          <>
            {createPortal(
              <ThemeProvider>
                <ThemeWrapper>
                  <div
                    // @ts-ignore
                    ref={setPopperElement}
                    style={{ ...styles.popper }}
                    {...attributes.popper}
                    onClick={(e) => e.stopPropagation()}
                  >
                    <div className="w-screen max-w-md flex-auto overflow-hidden rounded-md bg-dashboard border border-divide shadow-lg ring-1 ring-gray-900/5 p-4">
                      <SearchPathConfig onClose={close} />
                    </div>
                  </div>
                </ThemeWrapper>
              </ThemeProvider>,
              // @ts-ignore as this element definitely exists
              document.getElementById("portals"),
            )}
          </>
        )}
      </Popover.Panel>
    </Popover>
  ) : null;
};

export default ManageSearchPathButton;
