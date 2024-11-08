import Icon from "@powerpipe/components/Icon";
import LoadingIndicator from "@powerpipe/components/dashboards/LoadingIndicator";
import { CheckNodeStatus, DetectionSummary } from "../common";
import { classNames } from "@powerpipe/utils/styles";

type ProgressBarGroupProps = {
  children: JSX.Element | JSX.Element[];
  className?: string;
};

type ProgressBarProps = {
  className?: string;
  percent: number;
};

type DetectionSummaryChartProps = {
  status: CheckNodeStatus;
  summary: DetectionSummary;
  firstChildSummaries: DetectionSummary[];
};

type ProgressBarGroupTotalProps = {
  className?: string;
  total: number;
};

const getWidth = (x, y) => {
  const percent = (x / (x + y)) * 100;
  return percent >= 0.5 ? Math.round(percent) : 1;
};

const ProgressBarGroup = ({ children, className }: ProgressBarGroupProps) => (
  <div
    className={classNames(
      "flex h-3 items-center justify-end space-x-2",
      className,
    )}
  >
    {children}
  </div>
);

const ProgressBar = ({ className, percent }: ProgressBarProps) => {
  if (!percent) {
    return null;
  }

  return (
    <div
      className={classNames("h-3", className)}
      aria-valuenow={percent}
      aria-valuemin={0}
      aria-valuemax={100}
      style={{ display: "inline-block", width: `${percent}%` }}
    />
  );
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
  firstChildSummaries,
}: DetectionSummaryChartProps) => {
  let maxAlerts = 0;
  let maxNonAlerts = 0;
  for (const firstChildSummary of firstChildSummaries) {
    const currentMaxAlerts = firstChildSummary.total;
    if (currentMaxAlerts > maxAlerts) {
      maxAlerts = currentMaxAlerts;
    }
  }

  let alertsWidth = getWidth(maxAlerts, maxNonAlerts);

  const isRunning = status === "running";
  const hasTotal = summary.total > 0;

  return (
    <div className="flex items-center justify-between space-x-2">
      {isRunning && (
        <div className="mr-2">
          <LoadingIndicator className=" w-5 h-5" />
        </div>
      )}

      {hasTotal ? (
        <div className="my-auto px-0" style={{ width: `${alertsWidth}%` }}>
          <ProgressBarGroup>
            <ProgressBar
              className={classNames(
                "border border-alert",
                isRunning ? "summary-chart-alarm-animate" : "bg-alert",
              )}
              percent={getDetectionSummaryChartPercent(
                summary.total,
                maxAlerts,
              )}
            />
          </ProgressBarGroup>
        </div>
      ) : (
        !isRunning && (
          <div className="flex justify-end w-full pr-4">
            <Icon
              className="h-6 w-6 text-ok fill-text-ok"
              icon="materialsymbols-solid:check_circle"
            />
          </div>
        )
      )}

      <span className="text-sm font-semibold">{summary.total}</span>
    </div>
  );
};

export default DetectionSummaryChart;
