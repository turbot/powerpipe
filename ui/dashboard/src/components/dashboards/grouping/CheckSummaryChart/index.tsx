import IntegerDisplay from "../../../IntegerDisplay";
import { CheckNodeStatus, CheckSummary } from "../common";
import { classNames } from "@powerpipe/utils/styles";

type ProgressBarGroupProps = {
  children: JSX.Element | JSX.Element[];
  className?: string;
};

type ProgressBarProps = {
  className?: string;
  percent: number;
};

type CheckSummaryChartProps = {
  status: CheckNodeStatus;
  summary: CheckSummary;
  summaryDiff: CheckSummary;
  firstChildSummaries: CheckSummary[];
};

type AlertProgressBarGroupTotalProps = {
  className?: string;
  summary: CheckSummary;
};

type AlertProgressBarGroupDiffTotalProps = AlertProgressBarGroupTotalProps & {
  summaryDiff?: CheckSummary;
};

type NonAlertProgressBarGroupTotalProps = {
  className?: string;
  summary: CheckSummary;
};

type NonAlertProgressBarGroupDiffTotalProps =
  NonAlertProgressBarGroupTotalProps & {
    summaryDiff?: CheckSummary;
  };

type ProgressBarGroupTotalProps = {
  className?: string;
  total: number;
  direction?: "up" | "down";
};

const getWidth = (x, y) => {
  const percent = (x / (x + y)) * 100;
  return percent >= 0.5 ? Math.round(percent) : 1;
};

const ProgressBarGroupTotal = ({
  className,
  total,
  direction,
}: ProgressBarGroupTotalProps) => (
  <span className={classNames(className, "text-right text-sm font-semibold")}>
    {direction === "up" && "↑"}
    {direction === "down" && "↓"}
    <IntegerDisplay num={total} withTitle={false} />
  </span>
);

const getSummaryTitle = (summary: CheckSummary): string => {
  const titleParts: string[] = [];
  if (summary.error) {
    titleParts.push(`Error: ${summary.error.toLocaleString()}`);
  }
  if (summary.alarm) {
    titleParts.push(`Alarm: ${summary.alarm.toLocaleString()}`);
  }
  if (summary.ok) {
    titleParts.push(`OK: ${summary.ok.toLocaleString()}`);
  }
  if (summary.info) {
    titleParts.push(`Info: ${summary.info.toLocaleString()}`);
  }
  if (summary.skip) {
    titleParts.push(`Skipped: ${summary.skip.toLocaleString()}`);
  }
  if (titleParts.length === 0) {
    return "";
  }
  return titleParts.join(`
`);
};

const AlertProgressBarGroupTotal = ({
  className,
  summary,
}: AlertProgressBarGroupTotalProps) => {
  const alertTotal = summary.error + summary.alarm;
  const newClassName = classNames(
    className,
    alertTotal > 0 ? "text-alert" : "text-foreground-lightest",
  );
  return <ProgressBarGroupTotal className={newClassName} total={alertTotal} />;
};

const NonAlertProgressBarGroupTotal = ({
  className,
  summary,
}: NonAlertProgressBarGroupTotalProps) => {
  const nonAlertTotal = summary.ok + summary.info + summary.skip;
  let textClassName;
  if (nonAlertTotal === 0) {
    textClassName = "text-foreground-lightest";
  } else if (summary.skip > summary.info && summary.skip > summary.ok) {
    textClassName = "text-black-scale-5";
  } else if (summary.info > summary.ok && summary.info >= summary.skip) {
    textClassName = "text-info";
  } else {
    textClassName = "text-ok";
  }

  const newClassName = classNames(className, textClassName);
  return (
    <ProgressBarGroupTotal className={newClassName} total={nonAlertTotal} />
  );
};

const AlertProgressBarGroupDiffTotal = ({
  className,
  summary,
}: AlertProgressBarGroupDiffTotalProps) => {
  if (!summary) {
    return null;
  }
  const diffTotal = summary.error + summary.alarm;

  if (diffTotal === 0) {
    return null;
  }

  const newClassName = classNames(
    className,
    diffTotal > 0 ? "text-alert" : "text-foreground-lightest",
  );
  return (
    <ProgressBarGroupTotal
      className={newClassName}
      total={diffTotal}
      direction={diffTotal > 0 ? "up" : "down"}
    />
  );
};

const NonAlertProgressBarGroupDiffTotal = ({
  className,
  summary,
}: NonAlertProgressBarGroupDiffTotalProps) => {
  const nonAlertDiffTotal = summary.ok + summary.info + summary.skip;

  console.log({ summary, nonAlertDiffTotal });

  if (nonAlertDiffTotal === 0) {
    return null;
  }

  let textClassName;
  if (summary.skip > summary.info && summary.skip > summary.ok) {
    textClassName = "text-black-scale-5";
  } else if (summary.info > summary.ok && summary.info >= summary.skip) {
    textClassName = "text-info";
  } else {
    textClassName = "text-ok";
  }

  const newClassName = classNames(className, textClassName);
  return (
    <ProgressBarGroupTotal
      className={newClassName}
      total={nonAlertDiffTotal}
      direction={nonAlertDiffTotal > 0 ? "up" : "down"}
    />
  );
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

export const getCheckSummaryChartPercent = (value, total) => {
  if (!value) {
    return 0;
  }
  const percentOfTotal = value / total;
  const rounded = Math.floor(percentOfTotal * 100);
  return Math.max(rounded, 3);
};

const CheckSummaryChart = ({
  status,
  summary,
  summaryDiff,
  firstChildSummaries,
}: CheckSummaryChartProps) => {
  let maxAlerts = 0;
  let maxNonAlerts = 0;
  for (const firstChildSummary of firstChildSummaries) {
    const currentMaxAlerts = firstChildSummary.error + firstChildSummary.alarm;
    const currentMaxNonAlerts =
      firstChildSummary.ok + firstChildSummary.info + firstChildSummary.skip;
    if (currentMaxAlerts > maxAlerts) {
      maxAlerts = currentMaxAlerts;
    }
    if (currentMaxNonAlerts > maxNonAlerts) {
      maxNonAlerts = currentMaxNonAlerts;
    }
  }
  // const [alarm, error, ok, info, skip] = ensureMinPercentages(name, [
  //   summary.alarm,
  //   summary.error,
  //   summary.ok,
  //   summary.info,
  //   summary.skip,
  // ]);
  let alertsWidth = getWidth(maxAlerts, maxNonAlerts);
  let nonAlertsWidth = getWidth(maxNonAlerts, maxAlerts);
  if (alertsWidth > nonAlertsWidth) {
    alertsWidth -= 2;
  } else {
    nonAlertsWidth -= 2;
  }

  return (
    <div className="flex items-center" title={getSummaryTitle(summary)}>
      <div className="my-auto px-0" style={{ width: `${alertsWidth}%` }}>
        <ProgressBarGroup className="flex-row-reverse">
          <ProgressBar
            className={classNames(
              "border border-alert",
              status === "running" ? "summary-chart-alarm-animate" : "bg-alert",
            )}
            percent={getCheckSummaryChartPercent(summary.alarm, maxAlerts)}
          />
          <ProgressBar
            className={classNames(
              "border border-alert",
              status === "running" ? "summary-chart-error-animate" : null,
            )}
            percent={getCheckSummaryChartPercent(summary.error, maxAlerts)}
          />
          <AlertProgressBarGroupTotal className="mr-2" summary={summary} />
          <AlertProgressBarGroupDiffTotal
            className="mr-2"
            summary={summaryDiff}
          />
        </ProgressBarGroup>
      </div>
      <div
        className={classNames(
          "h-6 w-0 border-l border-black-scale-4",
          status === "running" ? "subtle-ping" : null,
        )}
      />
      <div className="my-auto px-0" style={{ width: `${nonAlertsWidth}%` }}>
        <ProgressBarGroup>
          <ProgressBar
            className={classNames(
              "border border-ok",
              status === "running" ? "summary-chart-ok-animate" : "bg-ok",
            )}
            percent={getCheckSummaryChartPercent(summary.ok, maxNonAlerts)}
          />
          <ProgressBar
            className={classNames(
              "border border-info",
              status === "running" ? "summary-chart-info-animate" : "bg-info",
            )}
            percent={getCheckSummaryChartPercent(summary.info, maxNonAlerts)}
          />
          <ProgressBar
            className={classNames(
              "border border-skip",
              status === "running" ? "summary-chart-skip-animate" : "bg-skip",
            )}
            percent={getCheckSummaryChartPercent(summary.skip, maxNonAlerts)}
          />
          <NonAlertProgressBarGroupTotal className="ml-2" summary={summary} />
          <NonAlertProgressBarGroupDiffTotal
            className="ml-2"
            summary={summaryDiff}
          />
        </ProgressBarGroup>
      </div>
    </div>
  );
};

export default CheckSummaryChart;
