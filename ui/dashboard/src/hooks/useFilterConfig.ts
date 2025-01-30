import { Filter } from "@powerpipe/components/dashboards/grouping/common";
import { KeyValuePairs } from "@powerpipe/components/dashboards/common/types";
import { useMemo } from "react";
import { useSearchParams } from "react-router-dom";
import { validateFilter } from "@powerpipe/components/dashboards/grouping/FilterEditor";

const defaultFilter = {
  operator: "and",
  // @ts-ignore
  expressions: [{ operator: "equal" }],
} as Filter;

const useFilterConfig = (
  panelName?: string,
): {
  allFilters: KeyValuePairs<Filter>;
  filter: Filter;
  update: (filter: Filter) => void;
} => {
  const [searchParams, setSearchParams] = useSearchParams();

  const allFilters = useMemo(() => {
    const rawFilters = searchParams.get("where");
    if (rawFilters) {
      try {
        let parsedFilters: KeyValuePairs<Filter>;
        parsedFilters = JSON.parse(rawFilters);
        return parsedFilters;
      } catch (error) {
        console.error("Error parsing where filters", error);
        return {};
      }
    } else {
      return {};
    }
  }, [searchParams]);

  const filter = useMemo(() => {
    if (!panelName) {
      return defaultFilter;
    }
    return allFilters[panelName] || defaultFilter;
  }, [allFilters, panelName]);

  const update = (toSave: Filter) => {
    setSearchParams((previous) => {
      const newParams = new URLSearchParams(previous);

      if (!panelName) {
        return newParams;
      }

      if (!validateFilter(toSave)) {
        delete allFilters[panelName];
      } else {
        allFilters[panelName] = toSave;
      }

      if (!!Object.keys(allFilters).length) {
        newParams.set("where", JSON.stringify(allFilters));
        return newParams;
      } else {
        newParams.delete("where");
        return newParams;
      }
    });
  };

  return { allFilters, filter, update };
};

export default useFilterConfig;
