import { useDashboardDatetimeRange } from "@powerpipe/hooks/useDashboardDatetimeRange";
import { useDashboardSearchPath } from "@powerpipe/hooks/useDashboardSearchPath";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";

const useGlobalContextNavigate = () => {
  const { metadata } = useDashboardState();
  const { range } = useDashboardDatetimeRange();
  const { searchPathPrefix } = useDashboardSearchPath();
  const serverSupportsSearchPath = metadata?.supports_search_path;
  const serverSupportsTimeRange = metadata?.supports_time_range;

  const urlSearchParams = new URLSearchParams();
  if (serverSupportsTimeRange) {
    urlSearchParams.set("datetime_range", JSON.stringify(range));
  } else {
    urlSearchParams.delete("datetime_range");
  }

  if (serverSupportsSearchPath && searchPathPrefix.length) {
    urlSearchParams.set("search_path_prefix", searchPathPrefix.join(","));
  } else {
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
