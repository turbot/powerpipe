import IntegerDisplay from "../../../IntegerDisplay";
import { CheckNodeStatus, CheckSummary } from "../common";
import { classNames } from "@powerpipe/utils/styles";
import { ReactNode } from "react";

type ProgressBarGroupProps = {
  children: ReactNode;
  className?: string;
};

type ProgressBarProps = {
  className?: string;
  percent: number;
};

type CheckSummaryChartProps = {
  status: CheckNodeStatus;
  summary: CheckSummary;
  firstChildSummaries: CheckSummary[];
};

type AlertProgressBarGroupTotalProps = {
  className?: string;
  summary: CheckSummary;
};

type NonAlertProgressBarGroupTotalProps = {
  className?: string;
  summary: CheckSummary;
};

type ProgressBarGroupTotalProps = {
  className?: string;
  total: number;
};

const getWidth = (x, y) => {
  const percent = (x / (x + y)) * 100;
  return percent >= 0.5 ? Math.round(percent) : 1;
};

export const ProgressBarGroupTotal = ({
  className,
  total,
}: ProgressBarGroupTotalProps) => (
  <span className={classNames(className, "text-right text-sm font-semibold")}>
    {total > 0 ? <IntegerDisplay num={total} withTitle={false} /> : "0"}
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
  if (summary.invalid) {
    titleParts.push(`Invalid: ${summary.invalid.toLocaleString()}`);
  }
  if (summary.muted) {
    titleParts.push(`Muted: ${summary.muted.toLocaleString()}`);
  }
  if (summary.tbd) {
    titleParts.push(`TBD: ${summary.tbd.toLocaleString()}`);
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
  if (summary.skipped) {
    titleParts.push(`Skipped: ${summary.skipped.toLocaleString()}`);
  }
  if (titleParts.length === 0) {
    return "";
  }
  return titleParts.join(`\n`);
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
  const nonAlertTotal =
    summary.ok +
    summary.info +
    summary.muted +
    summary.skip +
    summary.skipped +
    summary.tbd;
  const okTotal = summary.ok;
  const infoTotal = summary.info;
  const skipTotal =
    summary.muted + summary.skip + summary.skipped + summary.tbd;

  let textClassName: string;
  if (nonAlertTotal === 0) {
    textClassName = "text-foreground-lightest";
  } else if (skipTotal > infoTotal && skipTotal > okTotal) {
    textClassName = "text-black-scale-5";
  } else if (infoTotal > okTotal && infoTotal >= skipTotal) {
    textClassName = "text-info";
  } else {
    textClassName = "text-ok";
  }

  const newClassName = classNames(className, textClassName);
  return (
    <ProgressBarGroupTotal className={newClassName} total={nonAlertTotal} />
  );
};

export const ProgressBarGroup = ({
  children,
  className,
}: ProgressBarGroupProps) => (
  <div className={classNames("flex h-3 items-center", className)}>
    {children}
  </div>
);

export const ProgressBar = ({ className, percent }: ProgressBarProps) => {
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
  firstChildSummaries,
}: CheckSummaryChartProps) => {
  let maxAlerts = 0;
  let maxNonAlerts = 0;
  for (const firstChildSummary of firstChildSummaries) {
    const currentMaxAlerts =
      firstChildSummary.error +
      firstChildSummary.alarm +
      (firstChildSummary.invalid || 0);
    const currentMaxNonAlerts =
      firstChildSummary.ok +
      firstChildSummary.info +
      firstChildSummary.skip +
      firstChildSummary.skipped +
      (firstChildSummary.muted || 0) +
      (firstChildSummary.tbd || 0);
    if (currentMaxAlerts > maxAlerts) {
      maxAlerts = currentMaxAlerts;
    }
    if (currentMaxNonAlerts > maxNonAlerts) {
      maxNonAlerts = currentMaxNonAlerts;
    }
  }
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
          <ProgressBar
            className={classNames(
              "border border-invalid",
              status === "running"
                ? "summary-chart-invalid-animate"
                : "bg-invalid",
            )}
            percent={getCheckSummaryChartPercent(summary.invalid, maxAlerts)}
          />
          <AlertProgressBarGroupTotal className="mr-2" summary={summary} />
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
            percent={getCheckSummaryChartPercent(
              summary.skip || summary.skipped,
              maxNonAlerts,
            )}
          />
          <ProgressBar
            className={classNames(
              "border border-muted",
              status === "running" ? "summary-chart-muted-animate" : "bg-muted",
            )}
            percent={getCheckSummaryChartPercent(summary.muted, maxNonAlerts)}
          />
          <ProgressBar
            className={classNames(
              "border border-tbd",
              status === "running" ? "summary-chart-tbd-animate" : "bg-tbd",
            )}
            percent={getCheckSummaryChartPercent(summary.tbd, maxNonAlerts)}
          />
          <NonAlertProgressBarGroupTotal className="ml-2" summary={summary} />
        </ProgressBarGroup>
      </div>
    </div>
  );
};

export default CheckSummaryChart;
