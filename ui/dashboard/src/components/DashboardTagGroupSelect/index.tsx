import sortBy from "lodash/sortBy";
import { DashboardSearchGroupByMode } from "@powerpipe/types";
import { useCallback, useEffect, useMemo, useState } from "react";
import { useDashboardSearch } from "@powerpipe/hooks/useDashboardSearch";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { useParams } from "react-router-dom";

const DashboardTagGroupSelect = ({ onClose }) => {
  const { availableDashboardsLoaded } = useDashboardState();
  const { search, updateGroupBy } = useDashboardSearch();
  const { dashboard_name } = useParams();

  const options = useMemo(() => {
    const o = [
      { groupBy: "mod", tag: "", label: "Mod" },
      {
        groupBy: "tag",
        tag: "category",
        label: "Category",
      },
      {
        groupBy: "tag",
        tag: "service",
        label: "Service",
      },
      {
        groupBy: "tag",
        tag: "type",
        label: "Type",
      },
    ];
    return sortBy(o, ["label"]);
  }, []);

  const findOption = useCallback(
    (groupBy) => {
      if (groupBy.value === "tag") {
        return options.find((o) => o.tag === groupBy.tag);
      }
      return options.find((o) => o.groupBy === "mod");
    },
    [options],
  );

  const [value, setValue] = useState(() => findOption(search.groupBy));

  useEffect(() => {
    setValue(findOption(search.groupBy));
  }, [findOption, search.groupBy]);

  if (
    !availableDashboardsLoaded ||
    !value ||
    (dashboard_name && !search.value)
  ) {
    return null;
  }

  return (
    <div className="flex items-center justify-between space-x-2 p-3">
      <label htmlFor="dashboardsGroupBy">Group by:</label>
      <select
        name="dashboardsGroupBy"
        value={`${value.groupBy}${value.tag ? `:${value.tag}` : ""}`}
        onChange={(e) => {
          const option = options.find((o) => {
            const parts = e.target.value.split(":");
            if (parts.length === 2) {
              return o.groupBy === parts[0] && o.tag === parts[1];
            }
            return o.groupBy === parts[0];
          });
          if (!option) {
            return;
          }
          updateGroupBy(
            option.groupBy as DashboardSearchGroupByMode,
            option.tag,
          );
          onClose();
        }}
        className="block border border-divide rounded-md bg-dashboard"
      >
        {options.map((option) => (
          <option
            key={`${option.groupBy}:${option.tag}`}
            value={`${option.groupBy}${option.tag ? `:${option.tag}` : ""}`}
          >
            {option.label}
          </option>
        ))}
      </select>
    </div>
  );
};

export default DashboardTagGroupSelect;
