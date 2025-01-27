import { useDashboardDatetimeRange } from "@powerpipe/hooks/useDashboardDatetimeRange";
import { useDashboardSearchPath } from "@powerpipe/hooks/useDashboardSearchPath";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { useSearchParams } from "react-router-dom";

const useGlobalContextNavigate = () => {
  const [existingSearchParams] = useSearchParams();
  const { metadata } = useDashboardState();
  const { range } = useDashboardDatetimeRange();
  const { searchPathPrefix } = useDashboardSearchPath();
  const serverSupportsSearchPath = metadata?.supports_search_path;
  const serverSupportsTimeRange = metadata?.supports_time_range;

  const urlSearchParams = new URLSearchParams();
  if (metadata && serverSupportsTimeRange) {
    urlSearchParams.set("datetime_range", JSON.stringify(range));
  } else if (metadata && !serverSupportsTimeRange) {
    urlSearchParams.delete("datetime_range");
  } else if (existingSearchParams.has("datetime_range")) {
    urlSearchParams.set("datetime_range", JSON.stringify(range));
  }

  if (metadata && serverSupportsSearchPath && searchPathPrefix.length) {
    urlSearchParams.set("search_path_prefix", searchPathPrefix.join(","));
  } else if (
    metadata &&
    (!serverSupportsSearchPath || !searchPathPrefix.length)
  ) {
    urlSearchParams.delete("search_path_prefix");
  } else if (existingSearchParams.has("search_path_prefix")) {
    urlSearchParams.set("search_path_prefix", searchPathPrefix.join(","));
  }

  return { search: urlSearchParams.toString() };
};

export default useGlobalContextNavigate;
