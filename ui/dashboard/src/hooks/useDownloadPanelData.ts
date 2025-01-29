import { LeafNodeData } from "@powerpipe/components/dashboards/common";
import { PanelDefinition } from "@powerpipe/types";
import { saveAs } from "file-saver";
import { timestampForFilename } from "@powerpipe/utils/date";
import { useCallback, useState } from "react";
import { useDashboardState } from "./useDashboardState";
import { usePapaParse } from "react-papaparse";

const useDownloadPanelData = (definition: PanelDefinition) => {
  const { selectedDashboard } = useDashboardState();
  const { jsonToCSV } = usePapaParse();
  const [processing, setProcessing] = useState(false);

  const downloadQueryData = useCallback(
    async (data: LeafNodeData | undefined) => {
      if (!data) {
        return;
      }
      setProcessing(true);
      const colNames = data.columns.map((c) => c.name);
      let csvRows: any[] = [];

      const jsonbColIndices = data.columns
        .filter(
          (i) =>
            i.data_type === "JSONB" ||
            i.data_type === "VARCHAR[]" ||
            i.data_type.startsWith("STRUCT"),
        )
        .map((i) => data.columns.indexOf(i)); // would return e.g. [3,6,9]

      for (const row of data.rows) {
        // Deep copy the row or else it will update
        // the values in query output
        const csvRow: any[] = [];
        colNames.forEach((col, index) => {
          csvRow[index] = jsonbColIndices.includes(index)
            ? JSON.stringify(row[col])
            : row[col];
        });
        csvRows.push(csvRow);
      }

      const csv = jsonToCSV([colNames, ...csvRows]);
      const blob = new Blob([csv], { type: "text/csv;charset=utf-8" });

      saveAs(
        blob,
        `${(
          definition.dashboard ||
          selectedDashboard?.full_name ||
          ""
        ).replaceAll(".", "_")}_${definition.panel_type}_${timestampForFilename(
          Date.now(),
        )}.csv`,
      );
      setProcessing(false);
    },
    [definition, jsonToCSV, selectedDashboard],
  );

  return { download: downloadQueryData, processing };
};

export default useDownloadPanelData;
