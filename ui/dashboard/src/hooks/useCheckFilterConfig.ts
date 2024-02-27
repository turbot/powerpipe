import { CheckFilter } from "@powerpipe/components/dashboards/check/common";
import { useMemo } from "react";
import { useSearchParams } from "react-router-dom";

const useCheckFilterConfig = (): CheckFilter => {
  const [searchParams] = useSearchParams();
  return useMemo(() => {
    const rawFilters = searchParams.get("where");
    if (rawFilters) {
      try {
        let parsedFilters: CheckFilter;
        parsedFilters = JSON.parse(rawFilters);
        return parsedFilters;
      } catch (error) {
        console.error("Error parsing where filters", error);
        return {
          operator: "and",
          // @ts-ignore
          expressions: [{ operator: "equal" }],
        };
      }
    } else {
      return {
        operator: "and",
        // @ts-ignore
        expressions: [{ operator: "equal" }],
      };
    }
  }, [searchParams]);
};

export default useCheckFilterConfig;
