import {
  KeyValuePairs,
  TableConfig,
} from "@powerpipe/components/dashboards/common/types";
import { TableProperties } from "@powerpipe/components/dashboards/Table";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { useMemo } from "react";
import { useSearchParams } from "react-router-dom";

// Table settings come from the URL, so a link that was edited by hand, cut
// short, or kept after the mod changed can hold anything. Keep only the
// well-formed entries: a table whose entry is unreadable falls back to its
// default columns, and a bad link never takes the dashboard down.
export const parseTableConfigs = (
  raw: string | null,
): KeyValuePairs<TableConfig> => {
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

  const tables: KeyValuePairs<TableConfig> = {};
  for (const [panelName, config] of Object.entries(parsed)) {
    if (!config || typeof config !== "object" || Array.isArray(config)) {
      continue;
    }
    const columns = (config as { display_columns?: unknown }).display_columns;
    if (
      !Array.isArray(columns) ||
      !columns.every((c) => typeof c === "string")
    ) {
      continue;
    }
    tables[panelName] = { display_columns: columns };
  }
  return tables;
};

const useTableConfig = (panelName?: string) => {
  const [searchParams, setSearchParams] = useSearchParams();
  const { panelsMap } = useDashboardState();
  const panel = panelName ? panelsMap[panelName] : undefined;

  const allTables = useMemo(
    () => parseTableConfigs(searchParams.get("table")),
    [searchParams],
  );

  const table = useMemo(() => {
    if (!panel) {
      return { display_columns: [] } as TableConfig;
    }

    const found = allTables[panel.name];
    if (found) {
      return found;
    } else if (
      (panel.panel_type === "table" || panel.panel_type === "detection") &&
      (panel.properties as TableProperties)?.display_columns?.length
    ) {
      return {
        display_columns: (panel.properties as TableProperties).display_columns,
      };
    } else {
      return { display_columns: [] } as TableConfig;
    }
  }, [allTables, panel]);

  const update = (toSave: TableConfig | null) => {
    setSearchParams((previous) => {
      const newParams = new URLSearchParams(previous);

      if (!panelName) {
        return newParams;
      }

      // Start from the settings the link holds now, minus anything
      // unreadable, so one table's change keeps the other tables' settings.
      const current = parseTableConfigs(newParams.get("table"));
      if (!toSave) {
        delete current[panelName];
      } else {
        current[panelName] = toSave;
      }

      if (!!Object.keys(current).length) {
        newParams.set("table", JSON.stringify(current));
        return newParams;
      } else {
        newParams.delete("table");
        return newParams;
      }
    });
  };

  return { allTables, table, update };
};

export default useTableConfig;
