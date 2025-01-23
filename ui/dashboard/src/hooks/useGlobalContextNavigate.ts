import { useDashboardSearchPath } from "@powerpipe/hooks/useDashboardSearchPath";
import { useDashboardDatetimeRange } from "@powerpipe/hooks/useDashboardDatetimeRange";
// import { useNavigate } from "react-router-dom";
// import { useCallback } from "react";

const useGlobalContextNavigate = () => {
  // const navigate = useNavigate();
  const { range } = useDashboardDatetimeRange();
  const { searchPathPrefix } = useDashboardSearchPath();

  const urlSearchParams = new URLSearchParams();
  if (range.from) {
    urlSearchParams.set("datetime_range", JSON.stringify(range));
  }
  if (searchPathPrefix.length) {
    urlSearchParams.set("search_path_prefix", searchPathPrefix.join(","));
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
