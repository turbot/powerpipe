import { createContext, ReactNode, useContext, useMemo } from "react";
import { DashboardSearch, DashboardSearchGroupByMode } from "@powerpipe/types";
import { useSearchParams } from "react-router-dom";

interface IDashboardSearchContext {
  search: DashboardSearch;
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
  const rawSearch = searchParams.get("search");
  const rawGroupBy = searchParams.get("group_by");
  const rawTag = searchParams.get("tag");

  const search = useMemo(() => {
    return {
      value: searchParams.get("search") || "",
      groupBy: {
        value:
          (searchParams.get("group_by") as DashboardSearchGroupByMode) ||
          defaultSearch?.groupBy?.value ||
          "tag",
        tag:
          searchParams.get("tag") || defaultSearch?.groupBy?.tag || "service",
      },
    };
  }, [defaultSearch, rawSearch, rawGroupBy, rawTag]);

  const updateSearchValue = (value: string | undefined) => {
    if (value) {
      searchParams.set("search", value);
    } else {
      searchParams.delete("search");
    }
    setSearchParams(searchParams);
  };

  const updateGroupBy = (value: DashboardSearchGroupByMode, tag?: string) => {
    searchParams.set("group_by", value);
    if (tag) {
      searchParams.set("tag", tag);
    } else {
      searchParams.delete("tag");
    }
    setSearchParams(searchParams);
  };

  return (
    <DashboardSearchContext.Provider
      value={{ search, updateSearchValue, updateGroupBy }}
    >
      {children}
    </DashboardSearchContext.Provider>
  );
};

export const useDashboardSearch = () => {
  const context = useContext(DashboardSearchContext);
  if (!context) {
    throw new Error(
      "useDashboardExecution must be used within a DashboardExecutionProvider",
    );
  }
  return context;
};
