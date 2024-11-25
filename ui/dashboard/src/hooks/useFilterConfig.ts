import { Filter } from "@powerpipe/components/dashboards/grouping/common";
import { useMemo } from "react";
import { useSearchParams } from "react-router-dom";

const defaultFilter = {
  operator: "and",
  // @ts-ignore
  expressions: [{ operator: "equal" }],
} as Filter;

const useFilterConfig = (): Filter => {
  const [searchParams] = useSearchParams();

  return useMemo(() => {
    const rawFilters = searchParams.get("where");
    if (rawFilters) {
      try {
        let parsedFilters: Filter;
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

export default useFilterConfig;
