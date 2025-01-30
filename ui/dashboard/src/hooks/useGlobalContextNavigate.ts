import { useDashboardDatetimeRange } from "@powerpipe/hooks/useDashboardDatetimeRange";
import { useDashboardSearchPath } from "@powerpipe/hooks/useDashboardSearchPath";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { useSearchParams } from "react-router-dom";
import { DashboardDataModeLive } from "@powerpipe/types";

const useGlobalContextNavigate = (includeDisplayOptions = true) => {
  const [existingSearchParams] = useSearchParams();
  const { dataMode, metadata, selectedDashboard } = useDashboardState();
  const { range } = useDashboardDatetimeRange();
  const { searchPathPrefix } = useDashboardSearchPath();
  const serverSupportsSearchPath = metadata?.supports_search_path;
  const serverSupportsTimeRange = metadata?.supports_time_range;

  const isLive = dataMode === DashboardDataModeLive;

  const urlSearchParams = new URLSearchParams();
  if (isLive && metadata && serverSupportsTimeRange) {
    urlSearchParams.set("datetime_range", JSON.stringify(range));
  } else if (isLive && existingSearchParams.has("datetime_range")) {
    urlSearchParams.set("datetime_range", JSON.stringify(range));
  }

  if (
    isLive &&
    metadata &&
    serverSupportsSearchPath &&
    searchPathPrefix.length
  ) {
    urlSearchParams.set("search_path_prefix", searchPathPrefix.join(","));
  } else if (isLive && existingSearchParams.has("search_path_prefix")) {
    urlSearchParams.set("search_path_prefix", searchPathPrefix.join(","));
  }

  if (
    includeDisplayOptions &&
    !selectedDashboard &&
    existingSearchParams.has("dashboard_display")
  ) {
    urlSearchParams.set(
      "dashboard_display",
      existingSearchParams.get("dashboard_display"),
    );
  }

  if (
    includeDisplayOptions &&
    !selectedDashboard &&
    existingSearchParams.has("group_by")
  ) {
    urlSearchParams.set("group_by", existingSearchParams.get("group_by"));
  }

  if (
    includeDisplayOptions &&
    !selectedDashboard &&
    existingSearchParams.has("tag")
  ) {
    urlSearchParams.set("tag", existingSearchParams.get("tag"));
  }

  if (
    includeDisplayOptions &&
    !selectedDashboard &&
    existingSearchParams.has("search")
  ) {
    urlSearchParams.set("search", existingSearchParams.get("search"));
  }

  return { search: urlSearchParams.toString() };
};

export default useGlobalContextNavigate;
