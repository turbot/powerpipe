import DetectionResultNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionResultNode";
import { DetectionNode } from "@powerpipe/components/dashboards/grouping/common";
import { saveAs } from "file-saver";
import { timestampForFilename } from "@powerpipe/utils/date";
import { useCallback, useState } from "react";
import { useDashboard } from "./useDashboard";
import { usePapaParse } from "react-papaparse";

const useDownloadDetectionData = (
  node: DetectionNode,
  resultNodes?: DetectionResultNode[],
) => {
  const { selectedDashboard } = useDashboard();
  const { jsonToCSV } = usePapaParse();
  const [processing, setProcessing] = useState(false);

  const download = useCallback(async () => {
    if (!resultNodes || resultNodes.length === 0) {
      return;
    }
    setProcessing(true);
    const columns = Array.from(
      new Set(resultNodes.flatMap((item) => item.result?.columns || [])),
    ).filter((c) => !!c);
    let csvRows: any[] = [];

    const jsonbColIndices = columns
      .filter(
        (i) => i.data_type === "VARCHAR[]" || i.data_type.startsWith("STRUCT"),
      )
      .map((i) => columns.indexOf(i)); // would return e.g. [3,6,9]

    for (const resultNode of resultNodes) {
      for (const row of resultNode.result?.rows || []) {
        // Deep copy the row or else it will update
        // the values in query output
        const csvRow: any[] = [];
        columns.forEach((col, index) => {
          csvRow[index] =
            col.name in row
              ? jsonbColIndices.includes(index)
                ? JSON.stringify(row[col.name])
                : row[col.name]
              : null;
        });
        csvRows.push(csvRow);
      }
    }

    const csv = jsonToCSV([columns.map((c) => c.name), ...csvRows]);
    const blob = new Blob([csv], { type: "text/csv;charset=utf-8" });

    saveAs(
      blob,
      `${node.name.replaceAll(".", "_")}_${node.type}_${timestampForFilename(
        Date.now(),
      )}.csv`,
    );
    setProcessing(false);
  }, [node.name, resultNodes, jsonToCSV, selectedDashboard]);

  return { download, processing };
};

export default useDownloadDetectionData;
