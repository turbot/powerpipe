import { CheckFilter } from "@powerpipe/components/dashboards/check/common";
import { useMemo } from "react";
import { useSearchParams } from "react-router-dom";

const useCheckFilterConfig = (): CheckFilter => {
  const [searchParams] = useSearchParams();
  const defaultFilter = {
    operator: "and",
    // @ts-ignore
    expressions: [{ operator: "equal" }],
  } as CheckFilter;
  return useMemo(() => {
    const rawFilters = searchParams.get("where");
    if (rawFilters) {
      try {
        let parsedFilters: CheckFilter;
        parsedFilters = JSON.parse(rawFilters);
        return parsedFilters;
      } catch (error) {
        console.error("Error parsing where filters", error);
        return defaultFilter;
      }
    } else {
      return defaultFilter;
    }
  }, [searchParams]);
};

export default useCheckFilterConfig;
