import { DisplayGroup } from "@powerpipe/components/dashboards/grouping/common";
import { KeyValuePairs } from "@powerpipe/components/dashboards/common/types";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { useMemo } from "react";
import { useSearchParams } from "react-router-dom";

const useGroupingConfig = (panelName?: string) => {
  const [searchParams, setSearchParams] = useSearchParams();
  const { panelsMap } = useDashboardState();
  const panel = panelName ? panelsMap[panelName] : undefined;

  // Grouping arrives from the URL, so it is untrusted input: shareable links
  // get hand-edited, truncated by chat clients, and outlive the mod they were
  // built for. Anything unusable is DROPPED so the panel falls back to its
  // default - previously each of these replaced the whole dashboard with an
  // error boundary, leaving no way back except editing the URL by hand:
  //
  //   ?grouping={not-valid-json   -> "Expected property name or '}' in JSON..."
  //   [{"type":"banana"}]         -> "Unknown group type banana"
  //   {"panel":"nonsense"}        -> "n.filter is not a function"
  //   [{"type":null}]             -> "Unknown group type null"
  //
  // The sibling `where` (filter) parameter already tolerates junk; this brings
  // grouping into line with it.
  const allGroupings = useMemo(() => {
    const rawGroupings = searchParams.get("grouping");
    if (!rawGroupings) {
      return {};
    }

    let parsed: any;
    try {
      parsed = JSON.parse(rawGroupings);
    } catch (e) {
      return {};
    }

    if (!parsed || typeof parsed !== "object" || Array.isArray(parsed)) {
      return {};
    }

    const knownTypes = new Set([
      "benchmark",
      "control",
      "control_tag",
      "detection",
      "dimension",
      "reason",
      "resource",
      "result",
      "severity",
      "status",
    ]);

    const safe: KeyValuePairs<DisplayGroup[]> = {};
    for (const [panelKey, groups] of Object.entries(parsed)) {
      if (!Array.isArray(groups)) {
        continue;
      }
      const usable = (groups as any[]).filter(
        (g) => !!g && typeof g === "object" && knownTypes.has(g.type),
      ) as DisplayGroup[];
      if (usable.length) {
        safe[panelKey] = usable;
      }
    }
    return safe;
  }, [searchParams]);

  const grouping = useMemo(() => {
    if (!panel) {
      return [] as DisplayGroup[];
    }

    if (
      panel.panel_type !== "benchmark" &&
      panel.panel_type !== "control" &&
      panel.panel_type !== "detection"
    ) {
      return [] as DisplayGroup[];
    }

    const found = allGroupings[panel.name];
    if (found) {
      return found;
    } else if (
      (!panel.benchmark_type || panel.benchmark_type === "control") &&
      (panel.panel_type === "benchmark" || panel.panel_type === "control")
    ) {
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
