import { DashboardNestingDisplayMode } from "@powerpipe/types";
import { useDashboardSearch } from "@powerpipe/hooks/useDashboardSearch";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { useEffect, useState } from "react";

const DashboardListIncludeNestedSelect = ({ onClose }) => {
  const { availableDashboardsLoaded } = useDashboardState();
  const { nestedDashboards, updateNestedDashboards } = useDashboardSearch();

  const options = [
    { value: "exclude", label: "Exclude" },
    { value: "include", label: "Include" },
  ];

  const findOption = (v) => options.find((o) => o.value === v);

  const [value, setValue] = useState(() => findOption(nestedDashboards));

  useEffect(() => {
    setValue(() => findOption(nestedDashboards));
  }, [nestedDashboards]);

  if (!availableDashboardsLoaded) {
    return null;
  }

  return (
    <div className="flex items-center justify-between space-x-2 p-3">
      <label htmlFor="dashboardsIncludeNested">Nested dashboards:</label>
      <select
        name="dashboardsIncludeNested"
        value={value.value}
        onChange={(e) => {
          const option = options.find((o) => o.value === e.target.value);
          if (!option) {
            return;
          }
          updateNestedDashboards(option.value as DashboardNestingDisplayMode);
          onClose();
        }}
        className="block border border-divide rounded-md bg-dashboard"
      >
        {options.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
    </div>
  );
};

export default DashboardListIncludeNestedSelect;
