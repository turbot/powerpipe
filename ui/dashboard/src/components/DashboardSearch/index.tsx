import SearchInput from "../SearchInput";
import useDebouncedEffect from "@powerpipe/hooks/useDebouncedEffect";
import { DashboardActions } from "@powerpipe/types";
import { useDashboard } from "@powerpipe/hooks/useDashboard";
import { useEffect, useState } from "react";

const DashboardSearch = () => {
  const {
    availableDashboardsLoaded,
    breakpointContext: { minBreakpoint },
    dispatch,
    search,
    metadata,
  } = useDashboard();
  const [innerValue, setInnerValue] = useState(search.value);

  useEffect(() => {
    setInnerValue(() => search.value);
  }, [search.value]);

  const updateSearchValue = (value) =>
    dispatch({ type: DashboardActions.SET_DASHBOARD_SEARCH_VALUE, value });
  useDebouncedEffect(() => updateSearchValue(innerValue), 250, [innerValue]);

  return (
    <div className="w-full sm:w-56 md:w-72 lg:w-96">
      <SearchInput
        disabled={!metadata || !availableDashboardsLoaded}
        placeholder={minBreakpoint("sm") ? "Search dashboards..." : "Search..."}
        value={innerValue}
        setValue={setInnerValue}
      />
    </div>
  );
};

export default DashboardSearch;
