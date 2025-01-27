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
  if (
    (metadata && serverSupportsTimeRange) ||
    existingSearchParams.has("datetime_range")
  ) {
    urlSearchParams.set("datetime_range", JSON.stringify(range));
  } else if (metadata) {
    urlSearchParams.delete("datetime_range");
  }

  if (
    (metadata && serverSupportsSearchPath && searchPathPrefix.length) ||
    existingSearchParams.has("search_path_prefix")
  ) {
    urlSearchParams.set("search_path_prefix", searchPathPrefix.join(","));
  } else if (
    metadata &&
    (!serverSupportsSearchPath || !searchPathPrefix.length)
  ) {
    urlSearchParams.delete("search_path_prefix");
  }

  return { search: urlSearchParams.toString() };

  // const wrappedNavigate = useCallback(
  //   (path: string, replace: boolean = false) => {
  //     const urlSearchParams = new URLSearchParams();
  //     if (range.to) {
  //       urlSearchParams.set("datetime_range", JSON.stringify(range));
  //     }
  //     if (searchPathPrefix.length) {
  //       urlSearchParams.set("search_path_prefix", searchPathPrefix.join(","));
  //     }
  //     navigate(
  //       `${path}${urlSearchParams.size ? urlSearchParams.toString() : ""}`,
  //       { replace },
  //     );
  //   },
  //   [navigate, range, searchPathPrefix],
  // );
};

export default useGlobalContextNavigate;
