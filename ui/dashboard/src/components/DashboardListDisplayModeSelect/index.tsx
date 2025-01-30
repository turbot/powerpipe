import { DashboardDisplayMode } from "@powerpipe/types";
import { useDashboardSearch } from "@powerpipe/hooks/useDashboardSearch";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { useEffect, useState } from "react";

const options = [
  { value: "top_level", label: "Top-Level" },
  { value: "all", label: "All" },
];

const DashboardListDisplayModeSelect = ({ onClose }) => {
  const { availableDashboardsLoaded } = useDashboardState();
  const { dashboardsDisplay, updateDashboardsDisplay } = useDashboardSearch();

  const findOption = (v) => options.find((o) => o.value === v);

  const [value, setValue] = useState(() => findOption(dashboardsDisplay));

  useEffect(() => {
    setValue(() => findOption(dashboardsDisplay));
  }, [dashboardsDisplay]);

  if (!availableDashboardsLoaded) {
    return null;
  }

  return (
    <select
      name="dashboardsIncludeNested"
      value={value.value}
      onChange={(e) => {
        const option = options.find((o) => o.value === e.target.value);
        if (!option) {
          return;
        }
        updateDashboardsDisplay(option.value as DashboardDisplayMode);
        onClose();
      }}
      className="w-full block border border-divide rounded-md bg-dashboard"
    >
      {options.map((option) => (
        <option key={option.value} value={option.value}>
          {option.label}
        </option>
      ))}
    </select>
  );
};

export default DashboardListDisplayModeSelect;
