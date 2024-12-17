import Icon from "@powerpipe/components/Icon";
import LoadingIndicator from "@powerpipe/components/dashboards/LoadingIndicator";
import {
  CheckNodeStatus,
  DetectionSeveritySummary,
  DetectionSummary,
} from "../common";
import { classNames } from "@powerpipe/utils/styles";
import {
  ProgressBar,
  ProgressBarGroup,
  ProgressBarGroupTotal,
} from "@powerpipe/components/dashboards/grouping/CheckSummaryChart";

type DetectionSummaryChartProps = {
  status: CheckNodeStatus;
  summary: DetectionSummary;
  severitySummary: DetectionSeveritySummary;
  firstChildSummaries: DetectionSummary[];
};

export const getDetectionSummaryChartPercent = (value, total) => {
  if (!value) {
    return 0;
  }
  const percentOfTotal = value / total;
  const rounded = Math.floor(percentOfTotal * 100);
  return Math.max(rounded, 3);
};

const DetectionSummaryChart = ({
  status,
  summary,
  severitySummary,
  firstChildSummaries,
}: DetectionSummaryChartProps) => {
  let maxFirstChildTotal = 0;

  // Calculate max alerts
  for (const firstChildSummary of firstChildSummaries) {
    const currentMaxAlerts = firstChildSummary.total;
    if (currentMaxAlerts > maxFirstChildTotal) {
      maxFirstChildTotal = currentMaxAlerts;
    }
  }

  const maxFirstChildTotalDigits = maxFirstChildTotal.toString().length;
  const summaryTotalDigits = summary.total.toString().length;

  // if (summary.total > 0) {
  //   console.log({
  //     total: summary.total,
  //     maxFirstChildTotal,
  //     status,
  //     summary,
  //     severitySummary,
  //     firstChildSummaries,
  //   });
  // }

  let isRunning = status === "running";
  const hasTotal = summary.total > 0;

  return (
    <div className="flex items-center justify-end">
      {isRunning && !hasTotal && <LoadingIndicator className="w-5 h-5" />}
      {!isRunning && !hasTotal && (
        <Icon
          className="block h-5 w-5 text-ok fill-text-ok"
          icon="materialsymbols-solid:check_circle"
        />
      )}
      {hasTotal && (
        <div className="flex w-full justify-end">
          <ProgressBarGroup className="flex-grow">
            <ProgressBar
              className={classNames(
                "border border-alert",
                isRunning
                  ? "summary-chart-severity-critical-animate"
                  : "bg-alert",
              )}
              percent={getDetectionSummaryChartPercent(
                severitySummary.critical,
                maxFirstChildTotal,
              )}
              // percent={30}
            />
            <ProgressBar
              className={classNames(
                "border border-orange",
                isRunning ? "summary-chart-severity-high-animate" : "bg-orange",
              )}
              percent={getDetectionSummaryChartPercent(
                severitySummary.high,
                maxFirstChildTotal,
              )}
              // percent={40}
            />
            <ProgressBar
              className={classNames(
                "border border-severity",
                isRunning
                  ? "summary-chart-severity-medium-animate"
                  : "bg-severity",
              )}
              percent={getDetectionSummaryChartPercent(
                severitySummary.medium,
                maxFirstChildTotal,
              )}
              // percent={20}
            />
            <ProgressBar
              className={classNames(
                "border border-info",
                isRunning ? "summary-chart-severity-low-animate" : "bg-info",
              )}
              percent={getDetectionSummaryChartPercent(
                severitySummary.low,
                maxFirstChildTotal,
              )}
              // percent={10}
            />
          </ProgressBarGroup>
        </div>
      )}
      <ProgressBarGroupTotal className="ml-2" total={summary.total} />
    </div>
  );
};

export default DetectionSummaryChart;
