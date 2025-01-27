import SearchInput from "../SearchInput";
import useDebouncedEffect from "@powerpipe/hooks/useDebouncedEffect";
import { useDashboardSearch } from "@powerpipe/hooks/useDashboardSearch";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { useEffect, useState } from "react";

const DashboardSearch = () => {
  const {
    availableDashboardsLoaded,
    breakpointContext: { minBreakpoint },
    metadata,
  } = useDashboardState();
  const { search, updateSearchValue } = useDashboardSearch();
  const [innerValue, setInnerValue] = useState(search.value);

  useEffect(() => {
    setInnerValue(() => search.value);
  }, [search.value]);

  useDebouncedEffect(
    () => {
      if (search.value === innerValue) {
        return;
      }
      updateSearchValue(innerValue);
    },
    250,
    [search.value, innerValue, updateSearchValue],
  );

  return (
    <div className="w-full sm:w-56 md:w-72 lg:w-96">
      <SearchInput
        disabled={!metadata || !availableDashboardsLoaded}
        placeholder={minBreakpoint("sm") ? "Search dashboards…" : "Search…"}
        value={innerValue}
        setValue={setInnerValue}
      />
    </div>
  );
};

export default DashboardSearch;
