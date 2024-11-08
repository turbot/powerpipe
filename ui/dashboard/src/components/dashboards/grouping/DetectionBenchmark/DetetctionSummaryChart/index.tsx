import IntegerDisplay from "../../../../IntegerDisplay";
import { CheckNodeStatus, DetectionSummary } from "../../common";
import { classNames } from "@powerpipe/utils/styles";
import Icon from "@powerpipe/components/Icon";

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

type AlertProgressBarGroupTotalProps = {
  className?: string;
  summary: DetectionSummary;
};

type ProgressBarGroupTotalProps = {
  className?: string;
  total: number;
};

const getWidth = (x, y) => {
  const percent = (x / (x + y)) * 100;
  return percent >= 0.5 ? Math.round(percent) : 1;
};

const ProgressBarGroupTotal = ({
  className,
  total,
}: ProgressBarGroupTotalProps) => (
  <span className={classNames(className, "text-right text-sm font-semibold")}>
    {total > 0 ? <IntegerDisplay num={total} withTitle={false} /> : "0"}
  </span>
);
const AlertProgressBarGroupTotal = ({
  className,
  summary,
}: AlertProgressBarGroupTotalProps) => {
  const alertTotal = summary.total;
  const newClassName = classNames(
    className,
     "text-foreground-lightest",
  );
  return <ProgressBarGroupTotal className={newClassName} total={alertTotal} />;
};

const ProgressBarGroup = ({ children, className }: ProgressBarGroupProps) => (
  <div className={classNames("flex h-3 items-center", className)}>
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
    const currentMaxAlerts = firstChildSummary.total
    if (currentMaxAlerts > maxAlerts) {
      maxAlerts = currentMaxAlerts;
    }
  }
 
  let alertsWidth = getWidth(maxAlerts, maxNonAlerts);
 
  return (
    <div className="flex items-center justify-between">
      <div className="my-auto px-0" style={{ width: `${alertsWidth}%` }}>
        <ProgressBarGroup >
          <ProgressBar
            className={classNames(
              "border border-alert",
               "bg-alert",
            )}
            percent={getDetectionSummaryChartPercent(summary.total, maxAlerts)}
          />
        
        </ProgressBarGroup>
      </div>
      {summary.total === 0 && (
        <div className="flex justify-end w-full pr-4">
          <Icon
            className="h-6 w-6 text-ok fill-text-ok"
            icon="materialsymbols-solid:check_circle"
          />
        </div>
      )}
    <AlertProgressBarGroupTotal className="mr-2" summary={summary} />
    </div>
  );
};

export default DetectionSummaryChart;
