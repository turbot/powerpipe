import { createContext, ReactNode, useCallback, useContext } from "react";
import {
  DashboardNestingDisplayMode,
  DashboardSearch,
  DashboardSearchGroupByMode,
} from "@powerpipe/types";
import { useSearchParams } from "react-router-dom";

interface IDashboardSearchContext {
  search: DashboardSearch;
  nestedDashboards: DashboardNestingDisplayMode;
  updateNestedDashboards: (value: DashboardNestingDisplayMode) => void;
  updateSearchValue: (value: string | undefined) => void;
  updateGroupBy: (value: DashboardSearchGroupByMode, tag?: string) => void;
}

interface DashboardSearchProviderProps {
  children: ReactNode;
  defaultSearch: DashboardSearch | undefined;
}

const DashboardSearchContext = createContext<IDashboardSearchContext | null>(
  null,
);

export const DashboardSearchProvider = ({
  children,
  defaultSearch,
}: DashboardSearchProviderProps) => {
  const [searchParams, setSearchParams] = useSearchParams();
  const search = {
    value: searchParams.get("search") || "",
    groupBy: {
      value:
        (searchParams.get("group_by") as DashboardSearchGroupByMode) ||
        defaultSearch?.groupBy?.value ||
        "tag",
      tag: searchParams.get("tag") || defaultSearch?.groupBy?.tag || "service",
    },
  };
  const nestedDashboards = (searchParams.get("nested_dashboards") ||
    "exclude") as DashboardNestingDisplayMode;

  const updateSearchValue = useCallback(
    (value: string | undefined) => {
      if (value) {
        searchParams.set("search", value);
      } else {
        searchParams.delete("search");
      }
      setSearchParams(searchParams);
    },
    [searchParams],
  );

  const updateGroupBy = useCallback(
    (value: DashboardSearchGroupByMode, tag?: string) => {
      searchParams.set("group_by", value);
      if (tag) {
        searchParams.set("tag", tag);
      } else {
        searchParams.delete("tag");
      }
      setSearchParams(searchParams);
    },
    [searchParams],
  );

  const updateNestedDashboards = useCallback(
    (value: DashboardNestingDisplayMode) => {
      searchParams.set("nested_dashboards", value);
      setSearchParams(searchParams);
    },
    [searchParams],
  );

  return (
    <DashboardSearchContext.Provider
      value={{
        nestedDashboards,
        search,
        updateNestedDashboards,
        updateSearchValue,
        updateGroupBy,
      }}
    >
      {children}
    </DashboardSearchContext.Provider>
  );
};

export const useDashboardSearch = () => {
  const context = useContext(DashboardSearchContext);
  if (!context) {
    throw new Error(
      "useDashboardSearch must be used within a DashboardSearchContext",
    );
  }
  return context;
};
