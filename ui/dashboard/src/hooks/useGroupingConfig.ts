import {
  checkDisplayGroupTypes,
  detectionDisplayGroupTypes,
  DisplayGroup,
} from "@powerpipe/components/dashboards/grouping/common";
import { KeyValuePairs } from "@powerpipe/components/dashboards/common/types";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { useMemo } from "react";
import { useSearchParams } from "react-router-dom";

// Groupings come from the URL, so a link that was edited by hand, cut short, or
// made for another mod can hold anything. Keep only well-formed groupings: a
// bad link then falls back to the default grouping instead of crashing the
// dashboard.
export const parseGroupings = (
  raw: string | null,
): KeyValuePairs<DisplayGroup[]> => {
  if (!raw) {
    return {};
  }

  let parsed: unknown;
  try {
    parsed = JSON.parse(raw);
  } catch (e) {
    return {};
  }

  if (!parsed || typeof parsed !== "object" || Array.isArray(parsed)) {
    return {};
  }

  const groupings: KeyValuePairs<DisplayGroup[]> = {};
  for (const [panelName, groups] of Object.entries(parsed)) {
    if (!Array.isArray(groups)) {
      continue;
    }
    const wellFormed = groups.filter(
      (g) => !!g && typeof g === "object" && typeof g.type === "string",
    );
    if (wellFormed.length) {
      groupings[panelName] = wellFormed;
    }
  }
  return groupings;
};

// The saved groups the panel's tree builder can draw, or none.
export const usableGroups = (
  groups: DisplayGroup[] | undefined,
  types: string[],
): DisplayGroup[] => (groups || []).filter((g) => types.includes(g.type));

const useGroupingConfig = (panelName?: string) => {
  const [searchParams, setSearchParams] = useSearchParams();
  const { panelsMap } = useDashboardState();
  const panel = panelName ? panelsMap[panelName] : undefined;

  const allGroupings = useMemo(
    () => parseGroupings(searchParams.get("grouping")),
    [searchParams],
  );

  // The group types the panel's tree builder draws, and the grouping it falls
  // back to when none is saved. The editor's Reset restores that default
  // itself: when the panel is already on it, clearing the saved grouping
  // changes nothing, so nothing re-renders. These depend on the panel's type
  // only, not on the panel's run state, so a running benchmark does not reset
  // the editor while you edit.
  const panelType = panel?.panel_type;
  const benchmarkType = panel?.benchmark_type;
  const { groupTypes, defaultGrouping } = useMemo(() => {
    if (
      (!benchmarkType || benchmarkType === "control") &&
      (panelType === "benchmark" || panelType === "control")
    ) {
      return {
        groupTypes: checkDisplayGroupTypes,
        defaultGrouping: [
          { type: "benchmark" },
          { type: "control" },
          { type: "result" },
        ] as DisplayGroup[],
      };
    }
    if (
      (benchmarkType === "detection" && panelType === "benchmark") ||
      panelType === "detection"
    ) {
      return {
        groupTypes: detectionDisplayGroupTypes,
        defaultGrouping: [
          { type: "benchmark" },
          { type: "detection" },
          { type: "result" },
        ] as DisplayGroup[],
      };
    }
    return {
      groupTypes: [] as string[],
      defaultGrouping: [] as DisplayGroup[],
    };
  }, [panelType, benchmarkType]);

  const grouping = useMemo(() => {
    const usable = usableGroups(
      panelName ? allGroupings[panelName] : undefined,
      groupTypes,
    );
    return usable.length ? usable : defaultGrouping;
  }, [allGroupings, panelName, groupTypes, defaultGrouping]);

  const update = (toSave: DisplayGroup[]) => {
    setSearchParams((previous) => {
      const newParams = new URLSearchParams(previous);

      if (!panelName) {
        return newParams;
      }

      if (!toSave || !toSave.length) {
        delete allGroupings[panelName];
      } else {
        allGroupings[panelName] = toSave;
      }

      if (!!Object.keys(allGroupings).length) {
        newParams.set("grouping", JSON.stringify(allGroupings));
        return newParams;
      } else {
        newParams.delete("grouping");
        return newParams;
      }
    });
  };

  return { allGroupings, grouping, defaultGrouping, update };
};

export default useGroupingConfig;
