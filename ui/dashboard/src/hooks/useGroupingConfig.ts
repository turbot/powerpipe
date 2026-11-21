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

  const grouping = useMemo(() => {
    if (!panel) {
      return [] as DisplayGroup[];
    }

    const found = allGroupings[panel.name];
    if (
      (!panel.benchmark_type || panel.benchmark_type === "control") &&
      (panel.panel_type === "benchmark" || panel.panel_type === "control")
    ) {
      const usable = usableGroups(found, checkDisplayGroupTypes);
      if (usable.length) {
        return usable;
      }
      return [
        { type: "benchmark" },
        { type: "control" },
        { type: "result" },
      ] as DisplayGroup[];
    } else if (
      (panel.benchmark_type === "detection" &&
        panel.panel_type === "benchmark") ||
      panel.panel_type === "detection"
    ) {
      const usable = usableGroups(found, detectionDisplayGroupTypes);
      if (usable.length) {
        return usable;
      }
      return [
        { type: "benchmark" },
        { type: "detection" },
        { type: "result" },
      ] as DisplayGroup[];
    } else {
      return [] as DisplayGroup[];
    }
  }, [allGroupings, panel]);

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

  return { allGroupings, grouping, update };
};

export default useGroupingConfig;
