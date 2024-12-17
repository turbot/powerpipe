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

const getSummaryTitle = (
  summary: DetectionSummary,
  severitySummary: DetectionSeveritySummary,
): string => {
  const titleParts: string[] = [];
  if (summary.error) {
    titleParts.push(`Error: ${summary.error.toLocaleString()}`);
  }
  if (severitySummary.critical) {
    titleParts.push(`Critical: ${severitySummary.critical.toLocaleString()}`);
  }
  if (severitySummary.high) {
    titleParts.push(`High: ${severitySummary.high.toLocaleString()}`);
  }
  if (severitySummary.medium) {
    titleParts.push(`Medium: ${severitySummary.medium.toLocaleString()}`);
  }
  if (severitySummary.low) {
    titleParts.push(`Low: ${severitySummary.low.toLocaleString()}`);
  }
  if (severitySummary.none) {
    titleParts.push(`None: ${severitySummary.none.toLocaleString()}`);
  }
  if (titleParts.length === 0) {
    return "";
  }
  return titleParts.join(`
`);
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

  console.log({ total: summary.total, error: summary.error });

  // const maxFirstChildTotalDigits = maxFirstChildTotal.toString().length;
  // const summaryTotalDigits = summary.total.toString().length;

  let isRunning = status === "running";
  const hasTotal = summary.total > 0;

  return (
    <div
      className="flex items-center justify-end space-x-3"
      title={getSummaryTitle(summary, severitySummary)}
    >
      {isRunning && !hasTotal && <LoadingIndicator className="w-5 h-5" />}
      {/*{!isRunning && !hasTotal && (*/}
      {/*  <Icon*/}
      {/*    className="block h-5 w-5 text-ok fill-text-ok"*/}
      {/*    icon="materialsymbols-solid:check_circle"*/}
      {/*  />*/}
      {/*)}*/}
      {hasTotal && (
        <div className="flex w-full">
          <ProgressBarGroup className="flex-grow justify-end">
            <ProgressBar
              className={classNames(
                "border border-alert",
                isRunning ? "summary-chart-error-animate" : null,
              )}
              percent={getDetectionSummaryChartPercent(
                summary.error,
                maxFirstChildTotal,
              )}
              // percent={30}
            />
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
      <ProgressBarGroupTotal total={summary.total} />
    </div>
  );
};

export default DetectionSummaryChart;
