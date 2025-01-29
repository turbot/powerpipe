import DetectionBenchmark from "@powerpipe/components/dashboards/grouping/common/DetectionBenchmark";
import useDownloadPanelData from "@powerpipe/hooks/useDownloadPanelData";
import { PanelDefinition } from "@powerpipe/types";
import { useMemo } from "react";

const useDownloadDetectionBenchmarkData = (benchmark: DetectionBenchmark) => {
  const definition = useMemo(
    () =>
      ({
        dashboard: benchmark.name,
        panel_type: "benchmark",
      }) as PanelDefinition,
    [benchmark],
  );
  const { download, processing } = useDownloadPanelData(definition);

  const downloadQueryData = async () => {
    if (!benchmark) {
      return;
    }
    const data = benchmark.get_data_table();
    return download(data);
  };

  return { download: downloadQueryData, processing };
};

export default useDownloadDetectionBenchmarkData;
