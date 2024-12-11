import {
  KeyValuePairs,
  TableConfig,
} from "@powerpipe/components/dashboards/common/types";
import { TableProperties } from "@powerpipe/components/dashboards/Table";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { useMemo } from "react";
import { useSearchParams } from "react-router-dom";

const useTableConfig = (panelName?: string) => {
  const [searchParams, setSearchParams] = useSearchParams();
  const { panelsMap } = useDashboardState();
  const panel = panelName ? panelsMap[panelName] : undefined;

  const allTables = useMemo(() => {
    const rawTables = searchParams.get("table");
    if (rawTables) {
      let parsedTables: KeyValuePairs<TableConfig>;
      parsedTables = JSON.parse(rawTables);
      return parsedTables;
    } else {
      return {};
    }
  }, [searchParams]);

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

      if (!toSave) {
        delete allTables[panelName];
      } else {
        allTables[panelName] = toSave;
      }

      if (!!Object.keys(allTables).length) {
        newParams.set("table", JSON.stringify(allTables));
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
