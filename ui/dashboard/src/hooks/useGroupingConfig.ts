import { DisplayGroup } from "@powerpipe/components/dashboards/grouping/common";
import { KeyValuePairs } from "@powerpipe/components/dashboards/common/types";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { useMemo } from "react";
import { useSearchParams } from "react-router-dom";

const useGroupingConfig = (panelName?: string) => {
  const [searchParams, setSearchParams] = useSearchParams();
  const { panelsMap } = useDashboardState();
  const panel = panelName ? panelsMap[panelName] : undefined;

  const allGroupings = useMemo(() => {
    const rawGroupings = searchParams.get("grouping");
    if (rawGroupings) {
      let parsedGroups: KeyValuePairs<DisplayGroup[]>;
      parsedGroups = JSON.parse(rawGroupings);
      return parsedGroups;
    } else {
      return {};
    }
  }, [searchParams]);

  // The grouping a panel falls back to when nothing is saved. Exposed so the
  // editor's Reset can restore it explicitly rather than relying on a re-render
  // that does not happen when the panel is already on the default.
  const defaultGrouping = useMemo(() => {
    if (!panel) {
      return [] as DisplayGroup[];
    }
    if (
      (!panel.benchmark_type || panel.benchmark_type === "control") &&
      (panel.panel_type === "benchmark" || panel.panel_type === "control")
    ) {
      return [
        { type: "benchmark" },
        { type: "control" },
        { type: "result" },
      ] as DisplayGroup[];
    }
    if (
      (panel.benchmark_type === "detection" &&
        panel.panel_type === "benchmark") ||
      panel.panel_type === "detection"
    ) {
      return [
        { type: "benchmark" },
        { type: "detection" },
        { type: "result" },
      ] as DisplayGroup[];
    }
    return [] as DisplayGroup[];
  }, [panel]);

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

  return { allGroupings, grouping, defaultGrouping, update };
};

export default useGroupingConfig;
